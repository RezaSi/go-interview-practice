package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Package-level errors used by the limiter implementations.
var (
	ErrNilReceiver     = fmt.Errorf("nil receiver")
	ErrTooManyRequests = fmt.Errorf("too many requests")
	ErrBadParam        = fmt.Errorf("bad parameter")
)

// RateLimiter defines the behavior shared by all limiter algorithms.
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

// RateLimiterMetrics stores counters and timing statistics for a limiter.
type RateLimiterMetrics struct {
	TotalRequests   int64
	AllowedRequests int64
	DeniedRequests  int64
	AverageWaitTime time.Duration
}

// TokenBucketLimiter implements the token bucket algorithm.
// Tokens are refilled over time up to the configured burst capacity.
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

// NewTokenBucketLimiter creates a token bucket limiter with the given rate and burst.
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

// Allow reports whether one request can proceed immediately.
func (tb *TokenBucketLimiter) Allow() bool {
	// Guard against nil receivers to keep callers safe.
	if tb == nil {
		return false
	}

	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill the bucket according to time elapsed since the last update.
	now := time.Now()
	elasped := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += float64(elasped) * float64(tb.rate)
	tb.lastRefill = now

	// Never allow the internal bucket to exceed burst capacity.
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}

	// Consume one token when capacity is available.
	if tb.tokens >= 1 {
		tb.tokens -= 1

		// Record a successful request.
		tb.metrics.TotalRequests++
		tb.metrics.AllowedRequests++
		return true
	}

	// Record a rejected request when the bucket is empty.
	tb.metrics.TotalRequests++
	tb.metrics.DeniedRequests++
	return false
}

// AllowN reports whether n requests can proceed immediately.
func (tb *TokenBucketLimiter) AllowN(n int) bool {
	// Reject invalid receivers and request sizes early.
	if tb == nil {
		return false
	}
	if n <= 0 {
		return false
	}

	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill the bucket before checking current capacity.
	now := time.Now()
	elasped := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += float64(elasped) * float64(tb.rate)
	tb.lastRefill = now

	// Keep the bucket bounded by the burst value.
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}

	// Consume n tokens atomically when enough capacity exists.
	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)

		// Count each request in the batch as allowed.
		tb.metrics.TotalRequests += int64(n)
		tb.metrics.AllowedRequests += int64(n)

		return true
	}

	// Count each request in the batch as denied.
	tb.metrics.TotalRequests += int64(n)
	tb.metrics.DeniedRequests += int64(n)
	return false
}

// Wait blocks until one request can be served or the context is canceled.
func (tb *TokenBucketLimiter) Wait(ctx context.Context) error {
	// Preserve a predictable error for nil receivers.
	if tb == nil {
		return fmt.Errorf("token bucket limiter method Wait: %w", ErrNilReceiver)
	}

	start := time.Now()
	if tb.Allow() {
		// A zero or near-zero wait still contributes to the moving average.
		tb.mu.Lock()
		awt := tb.calculateAverageWaitTimeLocked(time.Since(start))
		tb.metrics.AverageWaitTime = awt
		tb.mu.Unlock()
		return nil
	}

	ch := make(chan struct{})

	// Track the waiter so wait time can reflect queued demand.
	tb.mu.Lock()
	if len(tb.waitQueue) >= tb.maxQueue {
		tb.mu.Unlock()
		return fmt.Errorf("token bucket limiter metod Wait: bucket is full: %w", ErrTooManyRequests)
	}
	tb.waitQueue = append(tb.waitQueue, ch)

	tb.mu.Unlock()

	for {
		// Recompute the delay on each pass because other goroutines may
		// consume or refill tokens while this waiter is sleeping.
		tb.mu.Lock()
		waitTime := tb.calculateWaitTimeLocked()
		tb.mu.Unlock()
		select {
		case <-ctx.Done():
			tb.removeCh(ch)
			close(ch) // Clean up the wait queue channel
			return fmt.Errorf("token bucket limiter method Wait: %w", ctx.Err())
		case <-time.After(waitTime):
			if tb.Allow() {
				tb.mu.Lock()
				awt := tb.calculateAverageWaitTimeLocked(time.Since(start))
				tb.metrics.AverageWaitTime = awt
				tb.mu.Unlock()
				// Remove the waiter once the request is granted.
				tb.removeCh(ch)
				close(ch)
				return nil
			}
			tb.mu.Lock()
			awt := tb.calculateAverageWaitTimeLocked(time.Since(start))
			tb.metrics.AverageWaitTime = awt
			tb.mu.Unlock()
		}
	}
}

// removeCh removes a waiter channel from the internal queue.
func (tb *TokenBucketLimiter) removeCh(ch chan struct{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	for i, c := range tb.waitQueue {
		if c == ch {
			tb.waitQueue = append(tb.waitQueue[:i], tb.waitQueue[i+1:]...)
			break
		}
	}
}

// calculateWaitTimeLocked estimates the wait needed for the current queue depth.
// The caller must hold tb.mu.
func (tb *TokenBucketLimiter) calculateWaitTimeLocked() time.Duration {
	tokenDef := float64(len(tb.waitQueue)) - tb.tokens
	waitTime := time.Duration(tokenDef / float64(tb.rate) * float64(time.Second))
	return waitTime
}

// calculateAverageWaitTimeLocked updates the rolling average wait duration.
// The caller must hold tb.mu.
func (tb *TokenBucketLimiter) calculateAverageWaitTimeLocked(waitTime time.Duration) time.Duration {
	return (tb.metrics.AverageWaitTime*time.Duration(tb.metrics.AllowedRequests) +
		waitTime) / time.Duration(float64(tb.metrics.AllowedRequests+1))
}

// WaitN blocks until capacity for n requests is available or the context ends.
func (tb *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
	// Reject invalid receivers and request sizes up front.
	if tb == nil {
		return fmt.Errorf("token bucket limiter method Wait: %w", ErrNilReceiver)
	}
	if n <= 0 {
		return fmt.Errorf("token bucket limiter method Wait: %w", ErrBadParam)
	}

	// Reserve per-request wait markers for queue accounting.
	chans := make([]chan struct{}, n)
	for i := 0; i < n; i++ {
		chans[i] = make(chan struct{})

	}

	start := time.Now()
	if tb.AllowN(n) {
		// Successful immediate acquisition still updates wait metrics.
		tb.mu.Lock()
		awt := tb.calculateAverageWaitTimeLocked(time.Since(start))
		tb.metrics.AverageWaitTime = awt
		tb.mu.Unlock()
		return nil
	}
	// Enqueue all requested slots so wait time reflects total demand.
	tb.mu.Lock()
	if len(tb.waitQueue)+n > tb.maxQueue {
		tb.mu.Unlock()
		return fmt.Errorf("token bucket limiter metod Wait: bucket is full: %w", ErrTooManyRequests)
	}
	tb.waitQueue = append(tb.waitQueue, chans...)
	tb.mu.Unlock()

	for {
		// Recalculate on each iteration to react to queue and token changes.
		tb.mu.Lock()
		waitTime := tb.calculateWaitTimeLocked()
		tb.mu.Unlock()
		select {
		case <-ctx.Done():
			for _, ch := range chans {
				tb.removeCh(ch)
				close(ch)
			}
			return fmt.Errorf("token bucket limiter method WaitN: %w", ctx.Err())
		case <-time.After(waitTime):
			if tb.AllowN(n) {
				tb.mu.Lock()
				awt := tb.calculateAverageWaitTimeLocked(time.Since(start))
				tb.metrics.AverageWaitTime = awt
				tb.mu.Unlock()
				for _, ch := range chans {
					tb.removeCh(ch)
					close(ch)
				}
				return nil
			}
			tb.mu.Lock()
			awt := tb.calculateAverageWaitTimeLocked(time.Since(start))
			tb.metrics.AverageWaitTime = awt
			tb.mu.Unlock()
		}
	}
}

// Limit returns the configured token refill rate.
func (tb *TokenBucketLimiter) Limit() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.rate
}

// Burst returns the maximum bucket capacity.
func (tb *TokenBucketLimiter) Burst() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.burst
}

// Reset restores the token bucket to a full initial state.
func (tb *TokenBucketLimiter) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens = float64(tb.burst)
	tb.lastRefill = time.Now()
	tb.metrics = RateLimiterMetrics{}
	tb.waitQueue = make([]chan struct{}, 0)
}

// GetMetrics returns a snapshot of the current token bucket metrics.
func (tb *TokenBucketLimiter) GetMetrics() RateLimiterMetrics {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.metrics
}

// SlidingWindowLimiter tracks requests within a moving time window.
type SlidingWindowLimiter struct {
	mu         sync.Mutex
	rate       int
	windowSize time.Duration
	requests   []time.Time // timestamps of recent requests
	metrics    RateLimiterMetrics
}

// NewSlidingWindowLimiter creates a sliding window limiter.
func NewSlidingWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	if rate <= 0 || windowSize <= 0 {
		panic("rate and window size must be positive")
	}
	return &SlidingWindowLimiter{
		rate:       rate,
		windowSize: windowSize,
		requests:   make([]time.Time, 0),
		metrics:    RateLimiterMetrics{},
	}
}

// Allow reports whether one request fits in the active window.
func (sw *SlidingWindowLimiter) Allow() bool {
	if sw == nil {
		return false
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	validRequests := make([]time.Time, 0)

	// Drop timestamps that no longer belong to the active window.
	cutoff := now.Add(-sw.windowSize)
	for _, req := range sw.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}

	sw.requests = validRequests

	// Accept the request when the current window is below its limit.
	if len(sw.requests) < sw.rate {
		sw.requests = append(sw.requests, now)
		// Record a successful request.
		sw.metrics.TotalRequests++
		sw.metrics.AllowedRequests++
		return true
	}

	// Record a rejected request when the window is full.
	sw.metrics.TotalRequests++
	sw.metrics.DeniedRequests++
	return false
}

// AllowN reports whether n requests fit in the active window.
func (sw *SlidingWindowLimiter) AllowN(n int) bool {
	if sw == nil {
		return false
	}
	if n <= 0 {
		return false
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	validRequests := make([]time.Time, 0)

	// Drop timestamps that no longer belong to the active window.
	cutoff := now.Add(-sw.windowSize)
	for _, req := range sw.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}

	sw.requests = validRequests

	// Accept the whole batch only when it fits as one atomic operation.
	if len(sw.requests)+n <= sw.rate {
		for i := 0; i < n; i++ {
			sw.requests = append(sw.requests, now)
		}
		// Record all requests in the batch as allowed.
		sw.metrics.TotalRequests += int64(n)
		sw.metrics.AllowedRequests += int64(n)
		return true
	}

	// Record all requests in the batch as denied.
	sw.metrics.TotalRequests += int64(n)
	sw.metrics.DeniedRequests += int64(n)
	return false
}

// Wait blocks until a single request fits in the sliding window.
func (sw *SlidingWindowLimiter) Wait(ctx context.Context) error {
	if sw == nil {
		return fmt.Errorf("sliding window limiter method Wait: %w", ErrNilReceiver)
	}

	start := time.Now()
	if sw.Allow() {
		// Immediate success still contributes to the average wait metric.
		sw.mu.Lock()
		awt := sw.calculateAverageWaitTimeLocked(time.Since(start))
		sw.metrics.AverageWaitTime = awt
		sw.mu.Unlock()
		return nil
	}

	for {
		// Recalculate the window delay every time through the loop because the
		// oldest request may change while other goroutines are active.
		sw.mu.Lock()
		waitTime := sw.calculateWaitTimeLocked()
		sw.mu.Unlock()
		select {
		case <-ctx.Done():
			return fmt.Errorf("sliding window limiter metod Wait: %w", ctx.Err())
		case <-time.After(waitTime):
			if sw.Allow() {
				sw.mu.Lock()
				// Record the observed wait after a successful retry.
				sw.metrics.AverageWaitTime = sw.calculateAverageWaitTimeLocked(time.Since(start))
				sw.mu.Unlock()
				return nil
			}
			// Keep the running average current even when a retry loses the race.
			sw.mu.Lock()
			sw.metrics.AverageWaitTime = sw.calculateAverageWaitTimeLocked(time.Since(start))
			sw.mu.Unlock()
		}
	}
}

// WaitN blocks until capacity exists for n requests in the sliding window.
func (sw *SlidingWindowLimiter) WaitN(ctx context.Context, n int) error {
	if sw == nil {
		return fmt.Errorf("sliding window limiter method Wait: %w", ErrNilReceiver)
	}
	if n <= 0 {
		return fmt.Errorf("sliding window limiter method WaitN: %w", ErrBadParam)
	}

	start := time.Now()
	if sw.AllowN(n) {
		// Immediate success still contributes to the average wait metric.
		sw.mu.Lock()
		awt := sw.calculateAverageWaitTimeLocked(time.Since(start))
		sw.metrics.AverageWaitTime = awt
		sw.mu.Unlock()
		return nil
	}

	for {
		// Recompute the delay on every loop because request expirations are dynamic.
		sw.mu.Lock()
		waitTime := sw.calculateWaitTimeLocked()
		sw.mu.Unlock()
		select {
		case <-ctx.Done():
			return fmt.Errorf("sliding window limiter metod Wait: %w", ctx.Err())
		case <-time.After(waitTime):
			if sw.AllowN(n) {
				sw.mu.Lock()
				// Record the observed wait after a successful retry.
				sw.metrics.AverageWaitTime = sw.calculateAverageWaitTimeLocked(time.Since(start))
				sw.mu.Unlock()
				return nil
			}
			// Update the average even if another goroutine consumed the slot first.
			sw.mu.Lock()
			sw.metrics.AverageWaitTime = sw.calculateAverageWaitTimeLocked(time.Since(start))
			sw.mu.Unlock()
		}
	}
}

// Limit returns the maximum number of requests allowed per window.
func (sw *SlidingWindowLimiter) Limit() int {
	return sw.rate
}

// Burst returns the effective burst size for compatibility with the interface.
func (sw *SlidingWindowLimiter) Burst() int {
	return sw.rate // sliding window doesn't have burst concept
}

// Reset clears the active sliding window and metric counters.
func (sw *SlidingWindowLimiter) Reset() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.requests = make([]time.Time, 0)
	sw.metrics = RateLimiterMetrics{}
}

// GetMetrics returns a snapshot of sliding window metrics.
func (sw *SlidingWindowLimiter) GetMetrics() RateLimiterMetrics {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.metrics
}

// calculateWaitTimeLocked returns the duration until the oldest request expires.
// The caller must hold sw.mu.
func (sw *SlidingWindowLimiter) calculateWaitTimeLocked() time.Duration {
	if len(sw.requests) == 0 {
		return 0
	}
	oldest := sw.requests[0]
	waitTime := time.Until(oldest.Add(sw.windowSize))
	if waitTime <= 0 {
		return 0
	}
	return waitTime
}

// calculateAverageWaitTimeLocked updates the rolling average wait duration.
// The caller must hold sw.mu.
func (sw *SlidingWindowLimiter) calculateAverageWaitTimeLocked(wt time.Duration) time.Duration {
	return (sw.metrics.AverageWaitTime*time.Duration(sw.metrics.AllowedRequests) + wt) / time.Duration(float64(sw.metrics.AllowedRequests+1))
}

// FixedWindowLimiter counts requests within discrete, resettable windows.
type FixedWindowLimiter struct {
	mu           sync.Mutex
	rate         int
	windowSize   time.Duration
	windowStart  time.Time
	requestCount int
	metrics      RateLimiterMetrics
}

// NewFixedWindowLimiter creates a fixed window limiter.
func NewFixedWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	if rate <= 0 || windowSize <= 0 {
		panic("rate and window size must be positive")
	}

	return &FixedWindowLimiter{
		rate:         rate,
		windowSize:   windowSize,
		windowStart:  time.Now(),
		requestCount: 0,
		metrics:      RateLimiterMetrics{},
	}
}

// Allow reports whether one request fits in the current fixed window.
func (fw *FixedWindowLimiter) Allow() bool {
	if fw == nil {
		return false
	}

	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	// Start a fresh window when the current one has expired.
	if now.Sub(fw.windowStart) >= fw.windowSize {
		fw.windowStart = now
		fw.requestCount = 0
	}

	// Accept the request when the current window has remaining capacity.
	if fw.requestCount < fw.rate {
		fw.requestCount++
		fw.metrics.TotalRequests++
		fw.metrics.AllowedRequests++
		return true
	}

	fw.metrics.TotalRequests++
	fw.metrics.DeniedRequests++
	return false
}

// AllowN reports whether n requests fit in the current fixed window.
func (fw *FixedWindowLimiter) AllowN(n int) bool {
	if fw == nil {
		return false
	}
	if n <= 0 {
		return false
	}

	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	// Roll the fixed window forward when it has expired.
	if now.Sub(fw.windowStart) >= fw.windowSize {
		fw.windowStart = now
		fw.requestCount = 0
	}

	// Accept the full batch only when enough room remains in this window.
	if fw.requestCount+n <= fw.rate {
		fw.requestCount += n
		fw.metrics.TotalRequests += int64(n)
		fw.metrics.AllowedRequests += int64(n)
		return true
	}

	fw.metrics.TotalRequests += int64(n)
	fw.metrics.DeniedRequests += int64(n)
	return false
}

// Wait blocks until one request can fit in the current or next fixed window.
func (fw *FixedWindowLimiter) Wait(ctx context.Context) error {
	if fw == nil {
		return fmt.Errorf("fixed window limiter method Wait: %w", ErrNilReceiver)
	}

	start := time.Now()
	if fw.Allow() {
		// Immediate success still contributes to the average wait metric.
		fw.mu.Lock()
		awt := fw.calculateAverageWaitTimeLocked(time.Since(start))
		fw.metrics.AverageWaitTime = awt
		fw.mu.Unlock()

		return nil
	}

	for {
		// Recompute the delay because window boundaries move over time.
		fw.mu.Lock()
		waitTime := fw.calculateWaitTimeLocked()
		fw.mu.Unlock()
		select {
		case <-ctx.Done():
			return fmt.Errorf("fixed window limiter method Wait: %w", ctx.Err())
		case <-time.After(waitTime):
			if fw.Allow() {
				fw.mu.Lock()
				// Record the observed delay for the granted request.
				fw.metrics.AverageWaitTime = fw.calculateAverageWaitTimeLocked(time.Since(start))
				fw.mu.Unlock()
				return nil
			}
			// Another goroutine may have claimed the slot first; keep waiting.
			fw.mu.Lock()
			fw.metrics.AverageWaitTime = fw.calculateAverageWaitTimeLocked(time.Since(start))
			fw.mu.Unlock()
		}
	}
}

// WaitN blocks until n requests can fit in the current or next fixed window.
func (fw *FixedWindowLimiter) WaitN(ctx context.Context, n int) error {
	if fw == nil {
		return fmt.Errorf("fixed window limiter method Wait: %w", ErrNilReceiver)
	}
	if n <= 0 {
		return fmt.Errorf("fixed window limiter method WaitN: %w", ErrBadParam)
	}

	start := time.Now()
	if fw.AllowN(n) {
		// Immediate success still contributes to the average wait metric.
		fw.mu.Lock()
		awt := fw.calculateAverageWaitTimeLocked(time.Since(start))
		fw.metrics.AverageWaitTime = awt
		fw.mu.Unlock()

		return nil
	}

	for {
		// Recompute on every loop so retries align with the next window boundary.
		fw.mu.Lock()
		waitTime := fw.calculateWaitTimeLocked()
		fw.mu.Unlock()
		select {
		case <-ctx.Done():
			return fmt.Errorf("fixed window limiter method Wait: %w", ctx.Err())
		case <-time.After(waitTime):
			if fw.AllowN(n) {
				fw.mu.Lock()
				// Record the observed delay for the granted batch.
				fw.metrics.AverageWaitTime = fw.calculateAverageWaitTimeLocked(time.Since(start))
				fw.mu.Unlock()
				return nil
			}
			// Another goroutine may have consumed the reopened window first.
			fw.mu.Lock()
			fw.metrics.AverageWaitTime = fw.calculateAverageWaitTimeLocked(time.Since(start))
			fw.mu.Unlock()
		}
	}
}

// Limit returns the maximum number of requests permitted per window.
func (fw *FixedWindowLimiter) Limit() int {
	return fw.rate
}

// Burst returns the effective burst size for compatibility with the interface.
func (fw *FixedWindowLimiter) Burst() int {
	return fw.rate
}

// Reset clears the current fixed window and metrics.
func (fw *FixedWindowLimiter) Reset() {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.windowStart = time.Now()
	fw.requestCount = 0
	fw.metrics = RateLimiterMetrics{}
}

// GetMetrics returns a snapshot of fixed window metrics.
func (fw *FixedWindowLimiter) GetMetrics() RateLimiterMetrics {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.metrics
}

// calculateWaitTimeLocked returns the time until the current fixed window ends.
// The caller must hold fw.mu.
func (fw *FixedWindowLimiter) calculateWaitTimeLocked() time.Duration {
	now := time.Now()
	if now.Sub(fw.windowStart) >= fw.windowSize {
		return 0
	}
	waitTime := time.Until(fw.windowStart.Add(fw.windowSize))
	if waitTime <= 0 {
		return 0
	}
	return waitTime
}

// calculateAverageWaitTimeLocked updates the rolling average wait duration.
// The caller must hold fw.mu.
func (fw *FixedWindowLimiter) calculateAverageWaitTimeLocked(wt time.Duration) time.Duration {
	return (fw.metrics.AverageWaitTime*time.Duration(fw.metrics.AllowedRequests) + wt) / time.Duration(float64(fw.metrics.AllowedRequests+1))
}

// RateLimiterFactory creates limiter implementations from configuration.
type RateLimiterFactory struct{}

// RateLimiterConfig selects the algorithm and parameters for a limiter instance.
type RateLimiterConfig struct {
	Algorithm  string        // "token_bucket", "sliding_window", "fixed_window"
	Rate       int           // requests per second
	Burst      int           // maximum burst capacity (for token bucket)
	WindowSize time.Duration // for sliding window and fixed window
}

// NewRateLimiterFactory creates a new factory instance.
func NewRateLimiterFactory() *RateLimiterFactory {
	return &RateLimiterFactory{}
}

// CreateLimiter validates the config and constructs the requested limiter.
func (f *RateLimiterFactory) CreateLimiter(config RateLimiterConfig) (RateLimiter, error) {
	switch config.Algorithm {
	case "token_bucket":
		if config.Rate <= 0 || config.Burst <= 0 {
			return nil, fmt.Errorf("invalid token bucket configuration: rate and burst must be positive %w", ErrBadParam)
		}
		return NewTokenBucketLimiter(config.Rate, config.Burst), nil
	case "sliding_window":
		if config.Rate <= 0 || config.WindowSize <= 0 {
			return nil, fmt.Errorf("Invalid sliding window configuration: rate and window size must be positive %w", ErrBadParam)
		}
		return NewSlidingWindowLimiter(config.Rate, config.WindowSize), nil
	case "fixed_window":
		if config.Rate <= 0 || config.WindowSize <= 0 {
			return nil, fmt.Errorf("Invalid fixed window configuration: rate and window size must be positive %w", ErrBadParam)
		}
		return NewFixedWindowLimiter(config.Rate, config.WindowSize), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", config.Algorithm)
	}
}

// RateLimitMiddleware wraps an HTTP handler with limiter checks.
func RateLimitMiddleware(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter.Allow() {
				next.ServeHTTP(w, r)
			} else {
				// Return a minimal but useful set of rate-limit headers.
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Rate limit exceeded"))
			}
		})
	}
}

// DistributedRateLimiter is a placeholder for a cross-instance limiter design.
type DistributedRateLimiter struct {
}

// AdaptiveRateLimiter is a placeholder for dynamic, load-aware limiting.
type AdaptiveRateLimiter struct {
}

// main prints a minimal message for the standalone challenge program.
func main() {
	fmt.Println("Rate Limiter Challenge - Solution Template")
	fmt.Println("Implement the TODO sections to complete the challenge")

	// Example usage once implemented:
	// limiter := NewTokenBucketLimiter(10, 5)
	// if limiter.Allow() {
	//     fmt.Println("Request allowed")
	// }
}
