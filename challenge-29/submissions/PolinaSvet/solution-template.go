package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
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
	totalWaitTime   time.Duration
}

// Token Bucket Rate Limiter
// 1. =================================================================
type TokenBucketLimiter struct {
	mu         sync.Mutex
	rate       int       // tokens per second
	burst      int       // maximum burst capacity
	tokens     float64   // current token count
	lastRefill time.Time // last token refill time
	metrics    RateLimiterMetrics
	waitQueue  []chan struct{} // queue for waiting requests
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
func NewTokenBucketLimiter(rate int, burst int) RateLimiter {
	// TODO: Implement token bucket rate limiter constructor
	// Initialize with proper rate, burst capacity, and current time
	// Set initial token count to burst capacity
	if rate <= 0 {
		rate = 1
	}

	if burst <= 0 {
		burst = 1
	}

	return &TokenBucketLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		metrics:    RateLimiterMetrics{},
		waitQueue:  make([]chan struct{}, 0),
	}
}

func (tb *TokenBucketLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	tb.tokens += elapsed * float64(tb.rate)

	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}

	tb.lastRefill = now
}

func (tb *TokenBucketLimiter) notifyWaiters() {
	if len(tb.waitQueue) == 0 {
		return
	}

	select {
	case tb.waitQueue[0] <- struct{}{}:
		tb.waitQueue = tb.waitQueue[1:]
	default:
	}

}

func (tb *TokenBucketLimiter) removeFromWaitQueue(ch chan struct{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	for i, waitCh := range tb.waitQueue {
		if waitCh == ch {
			tb.waitQueue = append(tb.waitQueue[:i], tb.waitQueue[i+1:]...)
			break
		}
	}
}

func (tb *TokenBucketLimiter) updateMetrics(waitTime time.Duration, n int) {
	tb.metrics.TotalRequests += int64(n)
	tb.metrics.totalWaitTime += waitTime

	if waitTime > 0 {
		totalRequests := float64(tb.metrics.TotalRequests)
		totalWait := float64(tb.metrics.totalWaitTime)
		tb.metrics.AverageWaitTime = time.Duration(totalWait / totalRequests)
	}

}

func (tb *TokenBucketLimiter) Allow() bool {
	// TODO: Implement Allow method for token bucket
	// 1. Calculate time elapsed since last refill
	// 2. Add tokens based on elapsed time and rate
	// 3. Cap tokens at burst capacity
	// 4. If tokens >= 1, consume one token and return true
	// 5. Update metrics

	return tb.AllowN(1)
}

func (tb *TokenBucketLimiter) AllowN(n int) bool {
	// TODO: Implement AllowN method for token bucket
	// Similar to Allow() but check for n tokens availability

	tb.mu.Lock()
	defer tb.mu.Unlock()

	start := time.Now()
	defer func() {
		tb.updateMetrics(time.Since(start), n)
	}()

	tb.refillTokens()

	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)
		tb.notifyWaiters()
		tb.metrics.AllowedRequests += int64(n)
		return true
	}

	return false
}

func (tb *TokenBucketLimiter) Wait(ctx context.Context) error {
	// TODO: Implement blocking Wait method
	// 1. If Allow() returns true, return immediately
	// 2. Calculate wait time based on token deficit
	// 3. Use context timeout and cancellation
	// 4. Update average wait time metrics
	return tb.WaitN(ctx, 1)
}

func (tb *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
	// TODO: Implement blocking WaitN method
	// Similar to Wait() but for n tokens
	start := time.Now()
	defer func() {
		tb.updateMetrics(time.Since(start), n)
	}()

	if tb.AllowN(n) {
		return nil
	}

	waitCh := make(chan struct{}, 1)
	tb.mu.Lock()
	tb.waitQueue = append(tb.waitQueue, waitCh)
	tb.mu.Unlock()

	select {

	case <-waitCh:
		if tb.AllowN(n) {
			return nil
		}
		return tb.WaitN(ctx, n)

	case <-ctx.Done():
		tb.removeFromWaitQueue(waitCh)
		return ctx.Err()

	case <-time.After(time.Second / time.Duration(tb.rate)):
		return tb.WaitN(ctx, n)
	}

}

func (tb *TokenBucketLimiter) Limit() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.rate
}

func (tb *TokenBucketLimiter) Burst() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.burst
}

func (tb *TokenBucketLimiter) Reset() {
	// TODO: Reset the rate limiter state
	// Set tokens to burst capacity, reset metrics, clear wait queue
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.tokens = float64(tb.burst)
	tb.lastRefill = time.Now()
	tb.metrics = RateLimiterMetrics{}
	for _, ch := range tb.waitQueue {
		close(ch)
	}
	tb.waitQueue = make([]chan struct{}, 0)
}

func (tb *TokenBucketLimiter) GetMetrics() RateLimiterMetrics {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.metrics.DeniedRequests = tb.metrics.TotalRequests - tb.metrics.AllowedRequests

	return tb.metrics
}

// Sliding Window Rate Limiter
// 2. =================================================================
type SlidingWindowLimiter struct {
	mu         sync.Mutex
	rate       int
	windowSize time.Duration
	requests   []time.Time // timestamps of recent requests
	metrics    RateLimiterMetrics
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter
func NewSlidingWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	// TODO: Implement sliding window rate limiter constructor
	if rate <= 0 {
		rate = 1
	}

	return &SlidingWindowLimiter{
		rate:       rate,
		windowSize: windowSize,
		requests:   make([]time.Time, 0),
		metrics:    RateLimiterMetrics{},
	}
}

func (sw *SlidingWindowLimiter) Allow() bool {
	// TODO: Implement Allow method for sliding window
	// 1. Remove old requests outside the window
	// 2. Check if current request count < rate
	// 3. If allowed, add current timestamp to requests
	// 4. Update metrics
	return sw.AllowN(1)
}

func (sw *SlidingWindowLimiter) calculateWaitTime(n int) (time.Duration, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()

	sw.cleanupOldRequests(now)

	if len(sw.requests)+n <= sw.rate {
		return 0, nil
	}

	requestsToRemove := (len(sw.requests) + n) - sw.rate
	if requestsToRemove <= 0 {
		return 0, nil
	}

	if requestsToRemove > len(sw.requests) {
		return 0, fmt.Errorf("invalid wait time calculation")
	}

	oldestRequestToRemove := sw.requests[requestsToRemove-1]
	waitUntil := oldestRequestToRemove.Add(sw.windowSize)

	if waitUntil.Before(now) {
		return 0, nil
	}

	return waitUntil.Sub(now), nil
}

func (sw *SlidingWindowLimiter) cleanupOldRequests(now time.Time) {

	if len(sw.requests) == 0 {
		return
	}

	windowStart := now.Add(-sw.windowSize)

	firstValidIndex := sw.findFirstValidIndex(windowStart)

	if firstValidIndex > 0 {
		sw.requests = sw.requests[firstValidIndex:]
	}
}

func (sw *SlidingWindowLimiter) findFirstValidIndex(windowStart time.Time) int {
	low, high := 0, len(sw.requests)-1

	for low <= high {
		mid := (low + high) / 2
		if sw.requests[mid].Before(windowStart) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low
}

func (sw *SlidingWindowLimiter) addRequests(timestamp time.Time, n int) {
	if len(sw.requests)+n > cap(sw.requests) {
		newCapacity := len(sw.requests) + n + 10
		newRequests := make([]time.Time, len(sw.requests), newCapacity)
		copy(newRequests, sw.requests)
		sw.requests = newRequests
	}

	for i := 0; i < n; i++ {
		sw.requests = append(sw.requests, timestamp)
	}
}

func (sw *SlidingWindowLimiter) updateMetrics(waitTime time.Duration, n int) {
	sw.metrics.TotalRequests += int64(n)
	sw.metrics.totalWaitTime += waitTime

	if waitTime > 0 {
		totalRequests := float64(sw.metrics.TotalRequests)
		totalWait := float64(sw.metrics.totalWaitTime)
		sw.metrics.AverageWaitTime = time.Duration(totalWait / totalRequests)
	}
}

func (sw *SlidingWindowLimiter) AllowN(n int) bool {
	// TODO: Implement AllowN method for sliding window

	sw.mu.Lock()
	defer sw.mu.Unlock()

	start := time.Now()
	defer func() {
		sw.updateMetrics(time.Since(start), n)
	}()

	now := time.Now()

	sw.cleanupOldRequests(now)

	if len(sw.requests)+n > sw.rate {
		return false
	}

	sw.addRequests(now, n)
	sw.metrics.AllowedRequests += int64(n)
	return true
}

func (sw *SlidingWindowLimiter) Wait(ctx context.Context) error {
	// TODO: Implement blocking Wait method for sliding window
	return sw.WaitN(ctx, 1)
}

func (sw *SlidingWindowLimiter) WaitN(ctx context.Context, n int) error {
	// TODO: Implement blocking WaitN method for sliding window
	start := time.Now()
	defer func() {
		sw.updateMetrics(time.Since(start), n)
	}()

	if sw.AllowN(n) {
		return nil
	}

	waitTime, err := sw.calculateWaitTime(n)
	if err != nil {
		return err
	}

	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		return sw.WaitN(ctx, n)
	case <-ctx.Done():
		return ctx.Err()
	}

}

func (sw *SlidingWindowLimiter) Limit() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	return sw.rate
}

func (sw *SlidingWindowLimiter) Burst() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	return sw.rate // sliding window doesn't have burst concept
}

func (sw *SlidingWindowLimiter) Reset() {
	// TODO: Reset sliding window state
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.requests = make([]time.Time, 0)
	sw.metrics = RateLimiterMetrics{}
}

func (sw *SlidingWindowLimiter) GetMetrics() RateLimiterMetrics {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.metrics.DeniedRequests = sw.metrics.TotalRequests - sw.metrics.AllowedRequests
	return sw.metrics
}

// Fixed Window Rate Limiter
// 3. =================================================================
type FixedWindowLimiter struct {
	mu           sync.Mutex
	rate         int
	windowSize   time.Duration
	windowStart  time.Time
	requestCount int
	metrics      RateLimiterMetrics
}

func (fw *FixedWindowLimiter) updateMetrics(waitTime time.Duration, n int) {
	fw.metrics.TotalRequests += int64(n)
	fw.metrics.totalWaitTime += waitTime

	if waitTime > 0 {
		totalRequests := float64(fw.metrics.TotalRequests)
		totalWait := float64(fw.metrics.totalWaitTime)
		fw.metrics.AverageWaitTime = time.Duration(totalWait / totalRequests)
	}
}

// NewFixedWindowLimiter creates a new fixed window rate limiter
func NewFixedWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	// TODO: Implement fixed window rate limiter constructor
	if rate <= 0 {
		rate = 1
	}

	return &FixedWindowLimiter{
		rate:         rate,
		windowSize:   windowSize,
		windowStart:  time.Now(),
		requestCount: 0,
		metrics:      RateLimiterMetrics{},
	}
}

func (fw *FixedWindowLimiter) Allow() bool {
	// TODO: Implement Allow method for fixed window
	// 1. Check if current time is in a new window
	// 2. If new window, reset counter and window start time
	// 3. If request count < rate, increment and allow
	// 4. Update metrics
	return fw.AllowN(1)
}

func (fw *FixedWindowLimiter) AllowN(n int) bool {
	// TODO: Implement AllowN method for fixed window

	fw.mu.Lock()
	defer fw.mu.Unlock()

	start := time.Now()
	defer func() {
		fw.updateMetrics(time.Since(start), n)
	}()

	now := time.Now()

	if now.Sub(fw.windowStart) >= fw.windowSize {
		fw.requestCount = 0
		fw.windowStart = now
	}

	if fw.requestCount+n <= fw.rate {
		fw.requestCount += n
		fw.metrics.AllowedRequests += int64(n)
		return true
	}

	return false

}

func (fw *FixedWindowLimiter) Wait(ctx context.Context) error {
	// TODO: Implement blocking Wait method for fixed window
	return fw.WaitN(ctx, 1)
}

func (fw *FixedWindowLimiter) WaitN(ctx context.Context, n int) error {
	// TODO: Implement blocking WaitN method for fixed window
	start := time.Now()
	defer func() {
		fw.updateMetrics(time.Since(start), n)
	}()

	if fw.AllowN(n) {
		return nil
	}

	// Calculate time until next window
	now := time.Now()
	nextWindow := fw.windowStart.Add(fw.windowSize)
	waitTime := nextWindow.Sub(now)

	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		return fw.WaitN(ctx, n)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (fw *FixedWindowLimiter) Limit() int {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	return fw.rate
}

func (fw *FixedWindowLimiter) Burst() int {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	return fw.rate
}

func (fw *FixedWindowLimiter) Reset() {
	// TODO: Reset fixed window state
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fw.windowStart = time.Now()
	fw.requestCount = 0
	fw.metrics = RateLimiterMetrics{}
}

func (fw *FixedWindowLimiter) GetMetrics() RateLimiterMetrics {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.metrics.DeniedRequests = fw.metrics.TotalRequests - fw.metrics.AllowedRequests
	return fw.metrics
}

// Rate Limiter Factory
// 4. =================================================================

type RateLimiterConfig struct {
	Algorithm  string        // "token_bucket", "sliding_window", "fixed_window"
	Rate       int           // requests per second
	Burst      int           // maximum burst capacity (for token bucket)
	WindowSize time.Duration // for sliding window and fixed window
}

type RateLimiterFactory struct{}

// NewRateLimiterFactory creates a new rate limiter factory
func NewRateLimiterFactory() *RateLimiterFactory {
	return &RateLimiterFactory{}
}

func (f *RateLimiterFactory) CreateLimiter(config RateLimiterConfig) (RateLimiter, error) {
	// TODO: Implement factory method to create different types of rate limiters
	// Validate configuration and create appropriate limiter type
	switch config.Algorithm {
	case "token_bucket":
		if config.Rate <= 0 || config.Burst <= 0 {
			return nil, fmt.Errorf("invalid token bucket configuration")
		}
		return NewTokenBucketLimiter(config.Rate, config.Burst), nil
	case "sliding_window":
		if config.Rate <= 0 || config.WindowSize <= 0 {
			return nil, fmt.Errorf("invalid sliding window configuration")
		}
		return NewSlidingWindowLimiter(config.Rate, config.WindowSize), nil
	case "fixed_window":
		if config.Rate <= 0 || config.WindowSize <= 0 {
			return nil, fmt.Errorf("invalid fixed window configuration")
		}
		return NewFixedWindowLimiter(config.Rate, config.WindowSize), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", config.Algorithm)
	}
}

// HTTP Middleware for rate limiting
// 5. =================================================================
func RateLimitMiddleware(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter.Allow() {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
				w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Burst())) // Упрощенно
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second).Unix()))

				next.ServeHTTP(w, r)
			} else {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second).Unix()))
				w.Header().Set("Retry-After", "1")

				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "rate limit exceeded", "retry_after": "1s"}`))
			}
		})
	}
}

// Advanced Features (Optional - for extra credit)

// DistributedRateLimiter - Rate limiter that works across multiple instances
type DistributedRateLimiter struct {
	// TODO: Implement distributed rate limiting using Redis or similar
	// This is an advanced feature for extra credit
}

// AdaptiveRateLimiter - Rate limiter that adjusts limits based on system load
type AdaptiveRateLimiter struct {
	// TODO: Implement adaptive rate limiting
	// Monitor system metrics and adjust rate limits dynamically
}

func main() {
	fmt.Println("=== Rate Limiter Factory ===")

	factory := NewRateLimiterFactory()

	testAlgorithms := []struct {
		name   string
		config RateLimiterConfig
	}{
		{
			name: "Token Bucket",
			config: RateLimiterConfig{
				Algorithm: "token_bucket",
				Rate:      5,
				Burst:     10,
			},
		},
		{
			name: "Sliding Window",
			config: RateLimiterConfig{
				Algorithm:  "sliding_window",
				Rate:       3,
				WindowSize: 1 * time.Second,
			},
		},
		{
			name: "Fixed Window",
			config: RateLimiterConfig{
				Algorithm:  "fixed_window",
				Rate:       4,
				WindowSize: 2 * time.Second,
			},
		},
	}

	for _, test := range testAlgorithms {
		fmt.Printf("\n--- Testing %s ---\n", test.name)
		testRateLimiter(factory, test.config)
	}

}

func testRateLimiter(factory *RateLimiterFactory, config RateLimiterConfig) {
	limiter, err := factory.CreateLimiter(config)
	if err != nil {
		log.Fatalf("Failed to create limiter: %v", err)
	}

	fmt.Printf("Config: Rate=%d", config.Rate)
	if config.Burst > 0 {
		fmt.Printf(", Burst=%d", config.Burst)
	}
	if config.WindowSize > 0 {
		fmt.Printf(", Window=%v", config.WindowSize)
	}
	fmt.Println()

	successCount := 0
	totalRequests := 15

	for i := 0; i < totalRequests; i++ {
		if limiter.Allow() {
			successCount++
			fmt.Printf("✓")
		} else {
			fmt.Printf("✗")
		}
		time.Sleep(200 * time.Millisecond)
	}

	metrics := limiter.GetMetrics()
	fmt.Printf("\nResults: %d/%d allowed (%.1f%%)\n",
		metrics.AllowedRequests,
		metrics.TotalRequests,
		float64(metrics.AllowedRequests)/float64(metrics.TotalRequests)*100)

	fmt.Printf("Metrics: Total=%d, Allowed=%d, Denied=%d, AvgWait=%v\n",
		metrics.TotalRequests, metrics.AllowedRequests,
		metrics.DeniedRequests, metrics.AverageWaitTime)
}
