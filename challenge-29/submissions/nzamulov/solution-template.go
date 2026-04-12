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
	// TODO: Implement token bucket rate limiter constructor
	// Initialize with proper rate, burst capacity, and current time
	// Set initial token count to burst capacity
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
	// TODO: Implement Allow method for token bucket
	// 1. Calculate time elapsed since last refill
	// 2. Add tokens based on elapsed time and rate
	// 3. Cap tokens at burst capacity
	// 4. If tokens >= 1, consume one token and return true
	// 5. Update metrics
	var coef = int(time.Second.Milliseconds()) / tb.rate
	var additional = int(time.Now().Sub(tb.lastRefill).Milliseconds()) / coef
	tb.tokens = min(tb.tokens + float64(additional), float64(tb.burst))
	tb.lastRefill = time.Now()
	tb.metrics.TotalRequests++
	if (tb.tokens >= 1) {
	    tb.tokens--
	    tb.metrics.AllowedRequests++
	    return true
	}
	tb.metrics.DeniedRequests++
	return false
}

func (tb *TokenBucketLimiter) AllowN(n int) bool {
	// TODO: Implement AllowN method for token bucket
	// Similar to Allow() but check for n tokens availability
	var coef = int(time.Second.Milliseconds()) / tb.rate
	var additional = int(time.Now().Sub(tb.lastRefill).Milliseconds()) / coef
	tb.tokens = min(tb.tokens + float64(additional), float64(tb.burst))
	tb.lastRefill = time.Now()
	if (tb.tokens >= float64(n)) {
	    tb.tokens -= float64(n)
	    tb.metrics.AllowedRequests += int64(n)
	    return true
	}
	tb.metrics.DeniedRequests += int64(n)
	return false
}

func (tb *TokenBucketLimiter) Wait(ctx context.Context) error {
	// TODO: Implement blocking Wait method
	// 1. If Allow() returns true, return immediately
	// 2. Calculate wait time based on token deficit
	// 3. Use context timeout and cancellation
	// 4. Update average wait time metrics
	if (tb.Allow()) {
	    return nil;
	}
	var coef = int(time.Second.Milliseconds()) / tb.rate
	var refill_diff = int(time.Now().Sub(tb.lastRefill).Milliseconds())
	var timeout = time.Duration(coef - refill_diff)
	
	cwt, cancel := context.WithTimeout(ctx, timeout * time.Millisecond)
	defer cancel()
	
	select {
	    case <-time.After(timeout * time.Millisecond):
	        return nil
	    case <-cwt.Done():
	        return cwt.Err()
	}
	
	return nil
}

func (tb *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
	// TODO: Implement blocking WaitN method
	// Similar to Wait() but for n tokens
	return nil
}

func (tb *TokenBucketLimiter) Limit() int {
	return tb.rate
}

func (tb *TokenBucketLimiter) Burst() int {
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
	tb.waitQueue = make([]chan struct{}, 0)
}

func (tb *TokenBucketLimiter) GetMetrics() RateLimiterMetrics {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.metrics
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
	// TODO: Implement sliding window rate limiter constructor
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
	var i int = 0
	for i < len(sw.requests) {
	    if (time.Now().Sub(sw.requests[i]) <= sw.windowSize) {
	        break
	    }
	    i++
	}
	sw.requests = sw.requests[i:]
	if len(sw.requests) + 1 <= sw.rate {
	    sw.requests = append(sw.requests, time.Now())
	    sw.metrics.AllowedRequests++
	    return true;
	}
	sw.metrics.DeniedRequests++
	return false
}

func (sw *SlidingWindowLimiter) AllowN(n int) bool {
	// TODO: Implement AllowN method for sliding window
	return false
}

func (sw *SlidingWindowLimiter) Wait(ctx context.Context) error {
	// TODO: Implement blocking Wait method for sliding window
	return nil
}

func (sw *SlidingWindowLimiter) WaitN(ctx context.Context, n int) error {
	// TODO: Implement blocking WaitN method for sliding window
	return nil
}

func (sw *SlidingWindowLimiter) Limit() int {
	return sw.rate
}

func (sw *SlidingWindowLimiter) Burst() int {
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
	return sw.metrics
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
	// TODO: Implement fixed window rate limiter constructor
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
	if time.Now().Sub(fw.windowStart) > fw.windowSize {
	    fw.windowStart = time.Now()
	    fw.requestCount = 1
	    fw.metrics.AllowedRequests++
	    return true
	} else if (fw.requestCount + 1 <= fw.rate) {
	    fw.requestCount++
	    fw.metrics.AllowedRequests++
	    return true
	}
	fw.metrics.DeniedRequests++
	return false
}

func (fw *FixedWindowLimiter) AllowN(n int) bool {
	// TODO: Implement AllowN method for fixed window
	return false
}

func (fw *FixedWindowLimiter) Wait(ctx context.Context) error {
	// TODO: Implement blocking Wait method for fixed window
	return nil
}

func (fw *FixedWindowLimiter) WaitN(ctx context.Context, n int) error {
	// TODO: Implement blocking WaitN method for fixed window
	return nil
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
	// TODO: Implement factory method to create different types of rate limiters
	// Validate configuration and create appropriate limiter type
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
	// TODO: Implement HTTP middleware for rate limiting
	// 1. Check if request is allowed using limiter.Allow()
	// 2. If allowed, call next handler
	// 3. If rate limited, return HTTP 429 (Too Many Requests)
	// 4. Add appropriate headers (X-RateLimit-Remaining, X-RateLimit-Reset, etc.)
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
