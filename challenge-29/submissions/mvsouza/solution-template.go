package main

import (
	"context"
	"errors"
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

func (rl *RateLimiterMetrics) Clone() *RateLimiterMetrics {
	return &RateLimiterMetrics{
		TotalRequests:   rl.TotalRequests,
		AllowedRequests: rl.AllowedRequests,
		DeniedRequests:  rl.DeniedRequests,
		AverageWaitTime: rl.AverageWaitTime,
	}
}

// Token Bucket Rate Limiter
type TokenBucketLimiter struct {
	mu         sync.Mutex
	rate       int       // tokens per second
	burst      int       // maximum burst capacity
	tokens     float64   // current token count
	lastRefill time.Time // last token refill time
	metrics    RateLimiterMetrics
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
func NewTokenBucketLimiter(rate int, burst int) RateLimiter {
	return &TokenBucketLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		metrics:    RateLimiterMetrics{},
	}
}

func (tb *TokenBucketLimiter) Allow() bool {
	return tb.AllowN(1)
}

func (tb *TokenBucketLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * float64(tb.rate)
	tb.lastRefill = now
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}
}

func (tb *TokenBucketLimiter) AllowN(n int) bool {
	if n < 0 {
		return false
	}
	if n == 0 {
		return true
	}
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refillTokens()
	tb.metrics.TotalRequests++
	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)
		tb.metrics.AllowedRequests++
		return true
	}
	tb.metrics.DeniedRequests++
	return false
}

func (tb *TokenBucketLimiter) Wait(ctx context.Context) error {
	return tb.WaitN(ctx, 1)
}

func (tb *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
	if n < 0 {
		return errors.New("n must be non-negative")
	}
	if n == 0 {
		return nil
	}
	if n > tb.Burst() {
		return errors.New("requested tokens exceed burst capacity")
	}
	ticker := time.NewTicker(time.Second / time.Duration(tb.rate))
	defer ticker.Stop()
	for {
		if tb.AllowN(n) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
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
	return &SlidingWindowLimiter{
		rate:       rate,
		windowSize: windowSize,
		requests:   make([]time.Time, 0),
		metrics:    RateLimiterMetrics{},
	}
}

func (sw *SlidingWindowLimiter) shouldPopHead() bool {
	threshhold := time.Now().Add(-1 * sw.windowSize)
	return len(sw.requests) > 0 && sw.requests[0].Before(threshhold)
}

func (sw *SlidingWindowLimiter) cleanWindow() {
	for sw.shouldPopHead() {
		sw.requests = sw.requests[1:]
	}
}

func (sw *SlidingWindowLimiter) hasQuota(n int) bool {
	sw.cleanWindow()
	return len(sw.requests)+n <= sw.rate
}

func (sw *SlidingWindowLimiter) Allow() bool {
	return sw.AllowN(1)
}

func (sw *SlidingWindowLimiter) AllowN(n int) bool {
	if n < 0 {
		return false
	}
	if n == 0 {
		return true
	}
	sw.mu.Lock()
	defer func() {
		sw.mu.Unlock()
	}()
	sw.metrics.TotalRequests++
	if !sw.hasQuota(n) {
		sw.metrics.DeniedRequests++
		return false
	}
	for i := 0; i < n; i++ {
		sw.requests = append(sw.requests, time.Now())
	}
	sw.metrics.AllowedRequests++
	return true
}

func (sw *SlidingWindowLimiter) Wait(ctx context.Context) error {
	// TODO: Implement blocking Wait method for sliding window
	return errors.New("sliding window Wait not implemented")
}

func (sw *SlidingWindowLimiter) WaitN(ctx context.Context, n int) error {
	// TODO: Implement blocking WaitN method for sliding window
	return errors.New("sliding window WaitN not implemented")
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
	return *sw.metrics.Clone()
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

func (fw *FixedWindowLimiter) hasQuota(n int) bool {
	threshhold := time.Now().Add(-1 * fw.windowSize)
	if fw.windowStart.Before(threshhold) {
		fw.windowStart = time.Now()
		fw.requestCount = 0
	}
	return fw.rate >= fw.requestCount+n
}

func (fw *FixedWindowLimiter) AllowN(n int) bool {
	if n < 0 {
		return false
	}
	if n == 0 {
		return true
	}
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.metrics.TotalRequests++
	if !fw.hasQuota(n) {
		fw.metrics.DeniedRequests++
		return false
	}
	fw.metrics.AllowedRequests++
	fw.requestCount += n
	return true
}

func (fw *FixedWindowLimiter) Wait(ctx context.Context) error {
	// TODO: Implement blocking Wait method for fixed window
	return errors.New("fixed window Wait not implemented")
}

func (fw *FixedWindowLimiter) WaitN(ctx context.Context, n int) error {
	// TODO: Implement blocking WaitN method for fixed window
	return errors.New("fixed window WaitN not implemented")
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
