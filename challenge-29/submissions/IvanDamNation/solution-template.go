package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Core Rate Limiter Interface
type RateLimiter interface {
	Allow() bool
	AllowN(n int) bool
	Wait(ctx context.Context) error
	WaitN(ctx context.Context, n int) error
	Limit() int
	Burst() int
	Reset()
	GetMetrics() RateLimiterMetrics
}

// Rate Limiter Metrics
type RateLimiterMetrics struct {
	TotalRequests   int64
	AllowedRequests int64
	DeniedRequests  int64
	AverageWaitTime time.Duration
}

type internalMetrics struct {
	totalRequests   atomic.Int64
	allowedRequests atomic.Int64
	deniedRequests  atomic.Int64
	totalWaitTimeNs atomic.Int64
	waitCount       atomic.Int64
}

func (m *internalMetrics) snapshot() RateLimiterMetrics {
	total := m.totalRequests.Load()
	allowed := m.allowedRequests.Load()
	denied := m.deniedRequests.Load()
	totalWait := m.totalWaitTimeNs.Load()
	count := m.waitCount.Load()

	var avgWait time.Duration
	if count > 0 {
		avgWait = time.Duration(totalWait / count)
	}

	return RateLimiterMetrics{
		TotalRequests:   total,
		AllowedRequests: allowed,
		DeniedRequests:  denied,
		AverageWaitTime: avgWait,
	}
}

// Token Bucket Rate Limiter
type TokenBucketLimiter struct {
	mu         sync.Mutex
	rate       int       // tokens per second
	burst      int       // maximum burst capacity
	tokens     float64   // current token count
	lastRefill time.Time // last token refill time
	metrics    internalMetrics
	waitQueue  []chan struct{} // queue for waiting requests
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
func NewTokenBucketLimiter(rate int, burst int) RateLimiter {
	if rate <= 0 || burst <= 0 {
		return nil
	}

	return &TokenBucketLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		metrics:    internalMetrics{},
		waitQueue:  make([]chan struct{}, 0),
	}
}

func (tb *TokenBucketLimiter) Allow() bool {
	return tb.AllowN(1)
}

func (tb *TokenBucketLimiter) AllowN(n int) bool {
	tb.metrics.totalRequests.Add(1)
	if n < 0 {
		tb.metrics.deniedRequests.Add(1)
		return false
	}
	n64 := float64(n)

	tb.mu.Lock()
	isAllowed, newTokens, nextRefill, _ := tb.hasTokensLazy(n64)
	if isAllowed {
		tb.tokens = newTokens - n64
		tb.lastRefill = nextRefill
	}
	tb.mu.Unlock()

	if !isAllowed {
		tb.metrics.deniedRequests.Add(1)
		return false
	}

	tb.metrics.allowedRequests.Add(1)
	return true
}

func (tb *TokenBucketLimiter) Wait(ctx context.Context) error {
	return tb.WaitN(ctx, 1)
}

func (tb *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	tb.metrics.totalRequests.Add(1)

	if n < 0 {
		tb.metrics.deniedRequests.Add(1)
		return fmt.Errorf("request size %d cannot be negative", n)
	}

	if n > tb.burst {
		tb.metrics.deniedRequests.Add(1)
		return fmt.Errorf("request size %d exceeds token bucket burst %v", n, tb.burst)
	}
	n64 := float64(n)

	start := time.Now()

	if n > tb.burst {
		tb.metrics.deniedRequests.Add(1)
		return fmt.Errorf("request size %d exceeds burst limit %d", n, tb.burst)
	}

	tb.mu.Lock()
	isAllowed, newTokens, nextRefill, sleepDuration := tb.hasTokensLazy(n64)
	tb.tokens = newTokens - n64
	tb.lastRefill = nextRefill
	tb.mu.Unlock()

	if isAllowed {
		tb.metrics.allowedRequests.Add(1)
		return nil
	}

	timer := time.NewTimer(sleepDuration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		tb.mu.Lock()

		now := time.Now()
		tb.tokens += now.Sub(tb.lastRefill).Seconds() * float64(tb.rate)
		tb.tokens += n64
		if tb.tokens > float64(tb.burst) {
			tb.tokens = float64(tb.burst)
		}
		tb.lastRefill = now
		tb.mu.Unlock()

		tb.metrics.deniedRequests.Add(1)
		return ctx.Err()
	case <-timer.C:
	}

	tb.metrics.allowedRequests.Add(1)
	tb.metrics.totalWaitTimeNs.Add(int64(time.Since(start)))
	tb.metrics.waitCount.Add(1)
	return nil
}

func (tb *TokenBucketLimiter) hasTokensLazy(n float64) (bool, float64, time.Time, time.Duration) {
	timeElapsed := time.Since(tb.lastRefill).Seconds()
	newTokens := tb.tokens + timeElapsed*float64(tb.rate)
	if newTokens > float64(tb.burst) {
		newTokens = float64(tb.burst)
	}

	now := time.Now()

	if newTokens < n {
		deficit := n - newTokens
		sleepSec := deficit / float64(tb.rate)
		sleepDuration := time.Duration(sleepSec * float64(time.Second))
		return false, newTokens, now, sleepDuration
	}

	return true, newTokens, now, 0
}

func (tb *TokenBucketLimiter) Limit() int {
	return tb.rate
}

func (tb *TokenBucketLimiter) Burst() int {
	return tb.burst
}

func (tb *TokenBucketLimiter) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.tokens = float64(tb.burst)
	tb.lastRefill = time.Now()

	tb.metrics.allowedRequests.Store(0)
	tb.metrics.deniedRequests.Store(0)
	tb.metrics.totalRequests.Store(0)
	tb.metrics.totalWaitTimeNs.Store(0)
	tb.metrics.waitCount.Store(0)

	tb.waitQueue = make([]chan struct{}, 0)
}

func (tb *TokenBucketLimiter) GetMetrics() RateLimiterMetrics {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.metrics.snapshot()
}

// Sliding Window Rate Limiter
type SlidingWindowLimiter struct {
	mu         sync.Mutex
	rate       int
	windowSize time.Duration
	requests   []time.Time // timestamps of recent requests
	readPtr    int
	reqQSize   int
	metrics    internalMetrics
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter
func NewSlidingWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	if rate <= 0 || windowSize <= 0 {
		return nil
	}

	return &SlidingWindowLimiter{
		rate:       rate,
		windowSize: windowSize,
		requests:   make([]time.Time, rate),
		metrics:    internalMetrics{},
	}
}

func (sw *SlidingWindowLimiter) Allow() bool {
	return sw.AllowN(1)
}

func (sw *SlidingWindowLimiter) AllowN(n int) bool {
	sw.metrics.totalRequests.Add(1)

	if n > sw.rate {
		sw.metrics.deniedRequests.Add(1)
		return false
	}

	now := time.Now()

	sw.mu.Lock()
	isAllowed, _ := sw.checkLazy(n, now)
	if isAllowed {
		for i := 0; i < n; i++ {
			sw.putRequest(now)
		}
	}
	sw.mu.Unlock()

	if !isAllowed {
		sw.metrics.deniedRequests.Add(1)
		return false
	}

	sw.metrics.allowedRequests.Add(1)
	return true
}

func (sw *SlidingWindowLimiter) Wait(ctx context.Context) error {
	return sw.WaitN(ctx, 1)
}

func (sw *SlidingWindowLimiter) WaitN(ctx context.Context, n int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	sw.metrics.totalRequests.Add(1)

	if n > sw.rate {
		sw.metrics.deniedRequests.Add(1)
		return fmt.Errorf("request size %d exceeds rate limit %d", n, sw.rate)
	}

	start := time.Now()

	for {
		now := time.Now()
		sw.mu.Lock()
		isAllowed, sleepDuration := sw.checkLazy(n, now)
		if isAllowed {
			for i := 0; i < n; i++ {
				sw.putRequest(now)
			}
			sw.mu.Unlock()
			break
		}
		sw.mu.Unlock()

		if sleepDuration <= 0 {
			sleepDuration = time.Millisecond
		}

		timer := time.NewTimer(sleepDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			sw.metrics.deniedRequests.Add(1)
			return ctx.Err()
		case <-timer.C:
			timer.Stop()
		}
	}

	sw.metrics.allowedRequests.Add(1)
	waitDuration := time.Since(start)
	if waitDuration > 0 {
		sw.metrics.totalWaitTimeNs.Add(int64(waitDuration))
		sw.metrics.waitCount.Add(1)
	}
	return nil
}

func (sw *SlidingWindowLimiter) Limit() int {
	return sw.rate
}

func (sw *SlidingWindowLimiter) Burst() int {
	return sw.rate // sliding window doesn't have burst concept
}

func (sw *SlidingWindowLimiter) Reset() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.requests = make([]time.Time, sw.rate)
	sw.readPtr = 0
	sw.reqQSize = 0

	sw.metrics.allowedRequests.Store(0)
	sw.metrics.deniedRequests.Store(0)
	sw.metrics.totalRequests.Store(0)
	sw.metrics.totalWaitTimeNs.Store(0)
	sw.metrics.waitCount.Store(0)
}

func (sw *SlidingWindowLimiter) GetMetrics() RateLimiterMetrics {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.metrics.snapshot()
}

func (sw *SlidingWindowLimiter) checkLazy(n int, now time.Time) (bool, time.Duration) {
	sw.cut(now)
	return sw.check(n)
}

func (sw *SlidingWindowLimiter) cut(now time.Time) {
	cutoff := now.Add(-sw.windowSize)

	for sw.reqQSize > 0 {
		oldest := sw.requests[sw.readPtr]
		if !oldest.Before(cutoff) {
			break
		}
		sw.advance()
	}
}

func (sw *SlidingWindowLimiter) check(n int) (bool, time.Duration) {
	freeSpace := len(sw.requests) - sw.reqQSize
	if freeSpace < n {
		requiredIndex := sw.calcRequiredIdx(n, freeSpace)
		estimatedTime := sw.requests[requiredIndex].Add(sw.windowSize)
		return false, time.Until(estimatedTime)
	}

	return true, 0
}

func (sw *SlidingWindowLimiter) advance() {
	sw.readPtr = (sw.readPtr + 1) % len(sw.requests)
	sw.reqQSize--
}

func (sw *SlidingWindowLimiter) calcRequiredIdx(nRequests, freeSpace int) int {
	return (sw.readPtr + (nRequests - freeSpace) - 1) % len(sw.requests)
}

func (sw *SlidingWindowLimiter) putRequest(now time.Time) {
	writePtr := (sw.readPtr + sw.reqQSize) % len(sw.requests)
	sw.requests[writePtr] = now
	sw.reqQSize++
}

// Fixed Window Rate Limiter
type FixedWindowLimiter struct {
	mu           sync.Mutex
	rate         int
	windowSize   time.Duration
	windowStart  time.Time
	requestCount int
	metrics      internalMetrics
}

// NewFixedWindowLimiter creates a new fixed window rate limiter
func NewFixedWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	if rate <= 0 || windowSize <= 0 {
		return nil
	}

	return &FixedWindowLimiter{
		rate:         rate,
		windowSize:   windowSize,
		windowStart:  time.Now(),
		requestCount: 0,
		metrics:      internalMetrics{},
	}
}

func (fw *FixedWindowLimiter) Allow() bool {
	return fw.AllowN(1)
}

func (fw *FixedWindowLimiter) AllowN(n int) bool {
	fw.metrics.totalRequests.Add(1)
	if n < 0 {
		fw.metrics.deniedRequests.Add(1)
		return false
	}

	fw.mu.Lock()
	isAllowed := fw.updateWindowAndCheckLazy(n)
	if isAllowed {
		fw.requestCount += n
	}
	fw.mu.Unlock()

	if !isAllowed {
		fw.metrics.deniedRequests.Add(1)
		return false
	}

	fw.metrics.allowedRequests.Add(1)
	return true
}

func (fw *FixedWindowLimiter) Wait(ctx context.Context) error {
	return fw.WaitN(ctx, 1)
}

func (fw *FixedWindowLimiter) WaitN(ctx context.Context, n int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fw.metrics.totalRequests.Add(1)

	if n < 0 {
		fw.metrics.deniedRequests.Add(1)
		return fmt.Errorf("request size %d cannot be negative", n)
	}

	if n > fw.rate {
		fw.metrics.deniedRequests.Add(1)
		return fmt.Errorf("request size %d exceeds rate limit %d", n, fw.rate)
	}

	start := time.Now()

	for {
		fw.mu.Lock()
		isAllowed := fw.updateWindowAndCheckLazy(n)
		if isAllowed {
			fw.requestCount += n
			fw.mu.Unlock()
			break
		}
		timeLeft := fw.windowSize - time.Since(fw.windowStart)
		fw.mu.Unlock()

		if timeLeft <= 0 {
			timeLeft = time.Millisecond
		}

		timer := time.NewTimer(timeLeft)

		select {
		case <-ctx.Done():
			timer.Stop()
			fw.metrics.deniedRequests.Add(1)
			return ctx.Err()
		case <-timer.C:
			timer.Stop()
		}
	}

	fw.metrics.allowedRequests.Add(1)
	waitDuration := time.Since(start)
	if waitDuration > 0 {
		fw.metrics.totalWaitTimeNs.Add(int64(waitDuration))
		fw.metrics.waitCount.Add(1)
	}

	return nil
}

func (fw *FixedWindowLimiter) updateWindowAndCheckLazy(n int) bool {
	elapsed := time.Since(fw.windowStart)
	if elapsed >= fw.windowSize {
		numWindows := elapsed / fw.windowSize
		fw.windowStart = fw.windowStart.Add(numWindows * fw.windowSize)
		fw.requestCount = 0
	}

	if fw.requestCount+n > fw.rate {
		return false
	}
	return true
}

func (fw *FixedWindowLimiter) Limit() int {
	return fw.rate
}

func (fw *FixedWindowLimiter) Burst() int {
	return fw.rate
}

func (fw *FixedWindowLimiter) Reset() {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.windowStart = time.Now()
	fw.requestCount = 0
	fw.metrics.allowedRequests.Store(0)
	fw.metrics.deniedRequests.Store(0)
	fw.metrics.totalRequests.Store(0)
	fw.metrics.totalWaitTimeNs.Store(0)
	fw.metrics.waitCount.Store(0)
}

func (fw *FixedWindowLimiter) GetMetrics() RateLimiterMetrics {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.metrics.snapshot()
}

// Rate Limiter Factory
type RateLimiterFactory struct{}

type RateLimiterConfig struct {
	Algorithm  string        // "token_bucket", "sliding_window", "fixed_window"
	Rate       int           // requests per second
	Burst      int           // maximum burst capacity (for token bucket)
	WindowSize time.Duration // for sliding window and fixed window
}

// NewRateLimiterFactory creates a new rate limiter factory
func NewRateLimiterFactory() *RateLimiterFactory {
	return &RateLimiterFactory{}
}

func (f *RateLimiterFactory) CreateLimiter(config RateLimiterConfig) (RateLimiter, error) {
	switch config.Algorithm {
	case "token_bucket":
		if config.Rate <= 0 || config.Burst <= 0 {
			return nil, fmt.Errorf("invalid token bucket configuration: rate and burst must be positive")
		}
		return NewTokenBucketLimiter(config.Rate, config.Burst), nil
	case "sliding_window":
		if config.Rate <= 0 || config.WindowSize <= 0 {
			return nil, fmt.Errorf("invalid sliding window configuration: rate and window size must be positive")
		}
		return NewSlidingWindowLimiter(config.Rate, config.WindowSize), nil
	case "fixed_window":
		if config.Rate <= 0 || config.WindowSize <= 0 {
			return nil, fmt.Errorf("invalid fixed window configuration: rate and window size must be positive")
		}
		return NewFixedWindowLimiter(config.Rate, config.WindowSize), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", config.Algorithm)
	}
}

// HTTP Middleware for rate limiting
func RateLimitMiddleware(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter.Allow() {
				next.ServeHTTP(w, r)
			} else {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Rate limit exceeded"))
			}
		})
	}
}

// Advanced Features (Optional - for extra credit)

/*
// DistributedRateLimiter - Rate limiter that works across multiple instances
type DistributedRateLimiter struct {
	rdb        *redis.Client
	keyPrefix  string // "ratelimit:user_id:"
	rate       int
	windowSize time.Duration
	metrics    internalMetrics
}

func NewDistributedRateLimiter(rdb *redis.Client, prefix string, rate int, window time.Duration) *DistributedRateLimiter {
	return &DistributedRateLimiter{
		rdb:        rdb,
		keyPrefix:  prefix,
		rate:       rate,
		windowSize: window,
	}
}

func (dr *DistributedRateLimiter) AllowN(ctx context.Context, key string, n int) (bool, error) {
	luaScript := `
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local increment = tonumber(ARGV[3])

		local current = redis.call('GET', key)
		if current and tonumber(current) + increment > limit then
			return 0
		else
			redis.call('INCRBY', key, increment)
			if not current then
				redis.call('PEXPIRE', key, window)
			end
			return 1
		end
	`

	redisKey := dr.keyPrefix + key
	result, err := dr.rdb.Eval(ctx, luaScript, []string{redisKey}, dr.rate, dr.windowSize.Milliseconds(), n).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}
*/

/*
// AdaptiveRateLimiter - Rate limiter that adjusts limits based on system load
type AdaptiveRateLimiter struct {
	mu          sync.Mutex
	baseLimiter *TokenBucketLimiter
	maxRate     int
	minRate     int

	latencyThreshold time.Duration
	stopChan         chan struct{}
}

func NewAdaptiveRateLimiter(base *TokenBucketLimiter, min, max int, threshold time.Duration) *AdaptiveRateLimiter {
	al := &AdaptiveRateLimiter{
		baseLimiter:      base,
		minRate:          min,
		maxRate:          max,
		latencyThreshold: threshold,
		stopChan:         make(chan struct{}),
	}

	go al.startSystemMonitor()
	return al
}

func (al *AdaptiveRateLimiter) startSystemMonitor() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-al.stopChan:
			return
		case <-ticker.C:
			al.mu.Lock()

			metrics := al.baseLimiter.GetMetrics()

			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			currentRate := al.baseLimiter.Limit()

			if metrics.AverageWaitTime > al.latencyThreshold {
				newRate := int(float64(currentRate) * 0.8)
				if newRate < al.minRate {
					newRate = al.minRate
				}
				al.baseLimiter.UpdateRate(newRate)
			} else {
				newRate := currentRate + 5
				if newRate > al.maxRate {
					newRate = al.maxRate
				}
				al.baseLimiter.UpdateRate(newRate)
			}

			al.mu.Unlock()
		}
	}
}

func (al *AdaptiveRateLimiter) Close() {
	close(al.stopChan)
}
*/

// Demo function to show basic usage
func main() {
	fmt.Println("Rate Limiter Challenge - Solution Template")
	fmt.Println("Implement the TODO sections to complete the challenge")

	// Example usage once implemented:
	// limiter := NewTokenBucketLimiter(10, 5)
	// if limiter.Allow() {
	//     fmt.Println("Request allowed")
	// }
}
