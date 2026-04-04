package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Custom errors
var (
	ErrNilReceiver     = fmt.Errorf("nil receiver")
	ErrTooManyRequests = fmt.Errorf("too many requests")
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
	maxQueue   int
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
func NewTokenBucketLimiter(rate int, burst int) RateLimiter {
	if rate <= 0 || burst <= 0 {
		panic("ratea and burst must be positive")
	}

	return &TokenBucketLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		metrics:    RateLimiterMetrics{},
		waitQueue:  make([]chan struct{}, 0),
		maxQueue:   1000,
	}
}

// Allow checks if a request can be allowed immediately without blocking
// Call this method under lock in concurrent scenarios
func (tb *TokenBucketLimiter) Allow() bool {
	// Nil receiver check po prevent panics
	if tb == nil {
		return false
	}

	// Generate new tokens based on elapsed time since last refill
	now := time.Now()
	elasped := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += float64(elasped) * float64(tb.rate)
	tb.lastRefill = now

	// Cap tokens at bust capasity
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}

	// Check if at least 1 token is avaliable for the request
	if tb.tokens >= 1 {
		tb.tokens -= 1

		// Update metrics
		tb.metrics.TotalRequests++
		tb.metrics.AllowedRequests++
		return true
	}

	// Update metrics
	tb.metrics.TotalRequests++
	tb.metrics.DeniedRequests++
	return false
}

// Allow checks if a n requests can be allowed immediately without blocking
// Call this method under lock in concurrent scenarios
func (tb *TokenBucketLimiter) AllowN(n int) bool {
	// Nil receiver check po prevent panics
	if tb == nil {
		return false
	}

	// Generate new tokens based on elapsed time since last refill
	now := time.Now()
	elasped := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += float64(elasped) * float64(tb.rate)
	tb.lastRefill = now

	// Cap tokens at bust capasity
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}

	// Check if at least 1 token is avaliable for the request
	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)

		// Update metrics
		tb.metrics.TotalRequests += int64(n)
		tb.metrics.AllowedRequests += int64(n)

		return true
	}

	// Update metrics
	tb.metrics.TotalRequests += int64(n)
	tb.metrics.DeniedRequests += int64(n)
	return false
}

// Wait blocks untill a request can be allowed or context is canceled or deadline is exceeded
// if the context is cancelled or deadline is exceeeded Wait closes the channel in TokenBucketLimiter.waitQueue
// always check the channel for closure !
func (tb *TokenBucketLimiter) Wait(ctx context.Context) error {
	// Nil reciever cheack to prevent panics
	if tb == nil {
		return fmt.Errorf("token bucket limiter method Wait: %w", ErrNilReceiver)
	}

	// Check for context cancellation. Default check if Allow returns true return immediately,
	// otherwise calculate wait time

	// Returns immediately if request allowed
	tb.mu.Lock()
	if tb.Allow() {
		tb.mu.Unlock()
		return nil
	}
	ch := make(chan struct{})

	// Add a new channel to the wait queue for this request
	if len(tb.waitQueue) >= tb.maxQueue {
		tb.mu.Unlock()
		return fmt.Errorf("token bucket limiter metod Wait: bucket is full: %w", ErrTooManyRequests)
	}
	tb.waitQueue = append(tb.waitQueue, ch)

	// Calculate token dificit and correspondeng wait and average wait time
	tokenDef := float64(len(tb.waitQueue)+1) - tb.tokens
	waitTime := time.Duration(tokenDef / float64(tb.rate) * float64(time.Second))
	averageWaitTime := tb.metrics.AverageWaitTime*time.Duration(tb.metrics.AllowedRequests) +
		waitTime/time.Duration(float64(tb.metrics.AllowedRequests+1))

	// Update metrics with average wait time
	// TODO Change this calculation to be more accuarate
	tb.metrics.AverageWaitTime = averageWaitTime

	tb.mu.Unlock()

	select {
	case <-ctx.Done():
		close(ch) // Clean up the wait queue channel
		return fmt.Errorf("token bucket limiter method Wait: %w", ctx.Err())
	case <-ch:
		return nil
	}
}

func (tb *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
	// Nil reciever cheack to prevent panics
	if tb == nil {
		return fmt.Errorf("token bucket limiter method Wait: %w", ErrNilReceiver)
	}

	// Check for context cancellation. Default check if Allow returns true return immediately,
	// otherwise calculate wait time

	// Returns immediately if request allowed
	tb.mu.Lock()
	if tb.AllowN(n) {
		tb.mu.Unlock()
		return nil
	}

	chans := make([]chan struct{}, n)
	for i := 0; i < n; i++ {
		chans[i] = make(chan struct{})

	}

	// Add a new channel to the wait queue for this request
	if len(tb.waitQueue)+n >= tb.maxQueue {
		tb.mu.Unlock()
		return fmt.Errorf("token bucket limiter metod Wait: bucket is full: %w", ErrTooManyRequests)
	}
	tb.waitQueue = append(tb.waitQueue, chans...)

	// Calculate token dificit and correspondeng wait and average wait time
	tokenDef := float64(len(tb.waitQueue)+1) - tb.tokens
	waitTime := time.Duration(tokenDef / float64(tb.rate) * float64(time.Second))
	averageWaitTime := tb.metrics.AverageWaitTime*time.Duration(tb.metrics.AllowedRequests) +
		waitTime/time.Duration(float64(tb.metrics.AllowedRequests+1))

	// Update metrics with average wait time
	// TODO Change this calculation to be more accuarate
	tb.metrics.AverageWaitTime = averageWaitTime

	tb.mu.Unlock()

	for _, ch := range chans {
		go func(ch chan struct{}) {
			select {
			case <-ctx.Done():
				close(ch) // Clean up the wait queue channel
				// return fmt.Errorf("token bucket limiter method Wait: %w", ctx.Err())
				return
			case <-ch:
				// return nil
				return
			}

		}(ch)
	}
	return nil
}

// Limit returns the configured rate limit
func (tb *TokenBucketLimiter) Limit() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.rate
}

// Burst returns the configured burst capacity
func (tb *TokenBucketLimiter) Burst() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.burst
}

// Reset set tokens to burst capacity, reset metrics, clear wait queue
func (tb *TokenBucketLimiter) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens = float64(tb.burst)
	tb.lastRefill = time.Now()
	tb.metrics = RateLimiterMetrics{}
	tb.waitQueue = make([]chan struct{}, 0)
}

// GetMetrics returns current rate limiter metrics
func (tb *TokenBucketLimiter) GetMetrics() RateLimiterMetrics {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.metrics
}

// ===========================
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
