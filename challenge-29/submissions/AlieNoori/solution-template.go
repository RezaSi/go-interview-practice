package main

import (
	"context"
	"fmt"
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
	waitedCount     int64
}

// Token Bucket Rate Limiter
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
	return &TokenBucketLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		metrics:    RateLimiterMetrics{},
		waitQueue:  make([]chan struct{}, 0),
	}
}

func (tb *TokenBucketLimiter) Allow() bool {
	return tb.AllowN(1)
}

func (tb *TokenBucketLimiter) AllowN(n int) bool {
	if n <= 0 {
		return true
	}

	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	tb.metrics.TotalRequests += int64(n)

	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)
		tb.metrics.AllowedRequests += int64(n)
		return true
	}

	tb.metrics.DeniedRequests += int64(n)
	return false
}

func (tb *TokenBucketLimiter) Wait(ctx context.Context) error {
	return tb.WaitN(ctx, 1)
}

func (tb *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
	if n <= 0 {
		return nil
	}

	if tb.AllowN(n) {
		return nil
	}

	start := time.Now()

	for {
		tb.mu.Lock()
		tb.refill()

		if tb.tokens >= float64(n) {
			tb.tokens -= float64(n)
			tb.metrics.TotalRequests += int64(n)
			tb.metrics.AllowedRequests += int64(n)
			tb.updateAverageWait(time.Since(start), n)
			tb.mu.Unlock()
			return nil
		}

		deficit := float64(n) - tb.tokens
		waitDuration := time.Duration(deficit / float64(tb.rate) * float64(time.Second))
		tb.mu.Unlock()

		timer := time.NewTimer(waitDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
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
	tb.metrics = RateLimiterMetrics{}
	tb.waitQueue = make([]chan struct{}, 0)
}

func (tb *TokenBucketLimiter) GetMetrics() RateLimiterMetrics {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.metrics
}

func (tb *TokenBucketLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = min(float64(tb.burst), tb.tokens+elapsed*float64(tb.rate))
	tb.lastRefill = now
}

func (tb *TokenBucketLimiter) updateAverageWait(d time.Duration, num int) {
	tb.metrics.waitedCount += int64(num)
	n := tb.metrics.waitedCount
	if n <= 0 {
		return
	}
	if n == 1 {
		tb.metrics.AverageWaitTime = d
		return
	}
	tb.metrics.AverageWaitTime += (d - tb.metrics.AverageWaitTime) / time.Duration(n)
}

// Sliding Window Rate Limiter
type SlidingWindowLimiter struct {
	mu         sync.Mutex
	rate       int
	windowSize time.Duration
	requests   []time.Time // timestamps of recent requests
	metrics    RateLimiterMetrics
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter
func NewSlidingWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	return &SlidingWindowLimiter{
		rate:       rate,
		windowSize: windowSize,
		requests:   make([]time.Time, 0),
		metrics:    RateLimiterMetrics{},
	}
}

func (sw *SlidingWindowLimiter) Allow() bool {
	return sw.AllowN(1)
}

func (sw *SlidingWindowLimiter) AllowN(n int) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()

	sw.cleanList()

	count := len(sw.requests) + n

	sw.metrics.TotalRequests += int64(n)

	if count <= sw.rate {
		for i := 0; i < n; i++ {
			sw.requests = append(sw.requests, now)
		}
		sw.metrics.AllowedRequests += int64(n)
		return true
	}

	sw.metrics.DeniedRequests += int64(n)

	return false
}

func (sw *SlidingWindowLimiter) Wait(ctx context.Context) error {
	return sw.WaitN(ctx, 1)
}

func (sw *SlidingWindowLimiter) WaitN(ctx context.Context, n int) error {
	if n <= 0 {
		return fmt.Errorf("n must be positive")
	}

	start := time.Now()
	for {
		if sw.AllowN(n) {
			sw.updateAverageWait(time.Since(start), n)
			return nil
		}

		now := time.Now()

		sw.mu.Lock()

		sw.cleanList()

		current := len(sw.requests)
		if current+n <= sw.rate {
			sw.mu.Unlock()
			continue
		}

		needToExpire := current + n - sw.rate
		idx := needToExpire - 1
		if idx < 0 || idx >= len(sw.requests) {
			sw.mu.Unlock()
			continue
		}

		expireAt := sw.requests[idx].Add(sw.windowSize)
		wait := expireAt.Sub(now)
		sw.mu.Unlock()

		if wait <= 0 {
			continue
		}

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
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
	sw.requests = make([]time.Time, 0)
	sw.metrics = RateLimiterMetrics{}
}

func (sw *SlidingWindowLimiter) GetMetrics() RateLimiterMetrics {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.metrics
}

func (sw *SlidingWindowLimiter) cleanList() {
	now := time.Now()

	validFrom := 0
	for validFrom < len(sw.requests) && now.Sub(sw.requests[validFrom]) > sw.windowSize {
		validFrom++
	}
	sw.requests = sw.requests[validFrom:]
}

func (sw *SlidingWindowLimiter) updateAverageWait(d time.Duration, num int) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.metrics.waitedCount += int64(num)
	n := sw.metrics.waitedCount
	if n <= 0 {
		return
	}
	if n == 1 {
		sw.metrics.AverageWaitTime = d
		return
	}
	sw.metrics.AverageWaitTime += (d - sw.metrics.AverageWaitTime) / time.Duration(n)
}

// Fixed Window Rate Limiter
type FixedWindowLimiter struct {
	mu           sync.Mutex
	rate         int
	windowSize   time.Duration
	windowStart  time.Time
	requestCount int
	metrics      RateLimiterMetrics
}

// NewFixedWindowLimiter creates a new fixed window rate limiter
func NewFixedWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	return &FixedWindowLimiter{
		rate:         rate,
		windowSize:   windowSize,
		windowStart:  time.Now(),
		requestCount: 0,
		metrics:      RateLimiterMetrics{},
	}
}

func (fw *FixedWindowLimiter) Allow() bool {
	return fw.AllowN(1)
}

func (fw *FixedWindowLimiter) AllowN(n int) bool {
	now := time.Now()

	fw.mu.Lock()
	defer fw.mu.Unlock()

	windowEnd := fw.windowStart.Add(fw.windowSize)

	if now.After(windowEnd) {
		fw.requestCount = 0
		fw.windowStart = now
	}

	fw.metrics.TotalRequests += int64(n)

	if fw.requestCount+n <= fw.rate {
		fw.requestCount += n
		fw.metrics.AllowedRequests += int64(n)
		return true
	}

	fw.metrics.DeniedRequests += int64(n)

	return false
}

func (fw *FixedWindowLimiter) Wait(ctx context.Context) error {
	return fw.WaitN(ctx, 1)
}

func (fw *FixedWindowLimiter) WaitN(ctx context.Context, n int) error {
	start := time.Now()
	for {
		if n > fw.rate {
			return fmt.Errorf("requested %d exceeds limiter capacity %d", n, fw.rate)
		}
		if fw.AllowN(n) {
			fw.updateAverageWait(time.Since(start), n)
			return nil
		}

		fw.mu.Lock()
		wait := fw.windowStart.Add(fw.windowSize).Sub(time.Now())
		fw.mu.Unlock()

		if wait < 0 {
			wait = 0
		}

		timer := time.NewTimer(wait)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
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
	fw.metrics = RateLimiterMetrics{}
}

func (fw *FixedWindowLimiter) GetMetrics() RateLimiterMetrics {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.metrics
}

func (fw *FixedWindowLimiter) updateAverageWait(d time.Duration, num int) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fw.metrics.waitedCount += int64(num)
	n := fw.metrics.waitedCount
	if n <= 0 {
		return
	}
	if n == 1 {
		fw.metrics.AverageWaitTime = d
		return
	}
	fw.metrics.AverageWaitTime += (d - fw.metrics.AverageWaitTime) / time.Duration(n)
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
