package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiterMetrics struct {
	TotalRequests   int64
	AllowedRequests int64
	DeniedRequests  int64
	AverageWaitTime time.Duration
}

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

// ----------------------------------------------------------
// Token Bucket
// ----------------------------------------------------------

type TokenBucketLimiter struct {
	mu         sync.Mutex
	rate       int       // tokens per second
	burst      int       // maximum burst capacity
	tokens     float64   // current token count
	lastRefill time.Time // last token refill time
	metrics    RateLimiterMetrics
}

func NewTokenBucketLimiter(rate int, burst int) RateLimiter {
	return &TokenBucketLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		metrics:    RateLimiterMetrics{},
	}
}

func (tb *TokenBucketLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	addTokens := elapsed * float64(tb.rate)
	if addTokens > 0 {
		tb.tokens += addTokens
		if tb.tokens > float64(tb.burst) {
			tb.tokens = float64(tb.burst)
		}
		tb.lastRefill = now
	}
}

func (tb *TokenBucketLimiter) Allow() bool {
	return tb.AllowN(1)
}

func (tb *TokenBucketLimiter) AllowN(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.metrics.TotalRequests++
	tb.refillTokens()
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
	start := time.Now()

	for {
		tb.mu.Lock()
		tb.refillTokens()
		if tb.tokens >= float64(n) {
			tb.tokens -= float64(n)
			tb.metrics.AllowedRequests++
			tb.metrics.TotalRequests++
			tb.mu.Unlock()
			elapsed := time.Since(start)
			tb.updateAvgWait(elapsed)
			return nil
		}
		tb.metrics.DeniedRequests++
		tb.metrics.TotalRequests++
		tb.mu.Unlock()

		select {
		case <-ctx.Done():
			elapsed := time.Since(start)
			tb.updateAvgWait(elapsed)
			return ctx.Err()
		default:
			time.Sleep(5 * time.Millisecond)
	}
	}
}

func (tb *TokenBucketLimiter) Limit() int {
	return tb.rate
}

func (tb *TokenBucketLimiter) Burst() int {
	return tb.burst
}

func (tb *TokenBucketLimiter) GetMetrics() RateLimiterMetrics {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.metrics
}

func (tb *TokenBucketLimiter) updateAvgWait(elapsed time.Duration) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	total := tb.metrics.AllowedRequests + tb.metrics.DeniedRequests
	if total == 0 {
		return
	}
	tb.metrics.AverageWaitTime += (elapsed - tb.metrics.AverageWaitTime) / time.Duration(total)
}

func (tb *TokenBucketLimiter) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.tokens = float64(tb.burst)
	tb.lastRefill = time.Now()
	tb.metrics = RateLimiterMetrics{}
}

// ----------------------------------------------------------
// Sliding Window
// ----------------------------------------------------------

type SlidingWindowLimiter struct {
	mu         sync.Mutex
	rate       int
	windowSize time.Duration
	requests   []time.Time
	metrics    RateLimiterMetrics
}

func NewSlidingWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	return &SlidingWindowLimiter{
		rate:       rate,
		windowSize: windowSize,
		requests:   make([]time.Time, 0, rate),
	}
}

func (sw *SlidingWindowLimiter) Allow() bool {
	return sw.AllowN(1)
}

func (sw *SlidingWindowLimiter) AllowN(n int) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	limit := now.Add(-sw.windowSize)
	newRequests := make([]time.Time, 0, len(sw.requests))
	for _, t := range(sw.requests) {
		if t.After(limit) {
			newRequests = append(newRequests, t)
		}
	}
	sw.requests = newRequests

	sw.metrics.TotalRequests += int64(n)
	if len(sw.requests)+n <= sw.rate {
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
	start := time.Now()
	for {
		if sw.AllowN(n) {
			elapsed := time.Since(start)
			sw.updateAvgWait(elapsed)
			return nil
		}
		select {
		case <-ctx.Done():
			elapsed := time.Since(start)
			sw.updateAvgWait(elapsed)
			return ctx.Err()
		default:
			time.Sleep(5 * time.Millisecond)
	}
	}
}

func (sw *SlidingWindowLimiter) Limit() int {
	return sw.rate
}

func (sw *SlidingWindowLimiter) Burst() int {
	return sw.rate
}

func (sw *SlidingWindowLimiter) GetMetrics() RateLimiterMetrics {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.metrics
}

func (sw *SlidingWindowLimiter) updateAvgWait(elapsed time.Duration) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	total := sw.metrics.AllowedRequests + sw.metrics.DeniedRequests
	if total == 0 {
		return
	}
	sw.metrics.AverageWaitTime += (elapsed - sw.metrics.AverageWaitTime) / time.Duration(total)
}

func (sw *SlidingWindowLimiter) Reset() {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.requests = sw.requests[:0]
	sw.metrics = RateLimiterMetrics{}
}

// ----------------------------------------------------------
// Fixed Window
// ----------------------------------------------------------

type FixedWindowLimiter struct {
	mu           sync.Mutex
	rate         int
	windowSize   time.Duration
	windowStart  time.Time
	requestCount int
	metrics      RateLimiterMetrics
}

func NewFixedWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
	return &FixedWindowLimiter{
		rate:         rate,
		windowSize:   windowSize,
		windowStart:  time.Now(),
		requestCount: 0,
	}
}

func (fw *FixedWindowLimiter) Allow() bool {
	return fw.AllowN(1)
}

func (fw *FixedWindowLimiter) AllowN(n int) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	if now.Sub(fw.windowStart) >= fw.windowSize {
		fw.windowStart = now
		fw.requestCount = 0
	}
	fw.metrics.TotalRequests++
	if fw.requestCount+n <= fw.rate {
		fw.requestCount += n
		fw.metrics.AllowedRequests++
		return true
	} else {
		fw.metrics.DeniedRequests++
		return false
	}
}

func (fw *FixedWindowLimiter) Wait(ctx context.Context) error {
	return fw.WaitN(ctx, 1)
}

func (fw *FixedWindowLimiter) WaitN(ctx context.Context, n int) error {
	start := time.Now()
	for {
		if fw.AllowN(n) {
			elapsed := time.Since(start)
			fw.updateAvgWait(elapsed)
			return nil
		}
		select {
		case <-ctx.Done():
			elapsed := time.Since(start)
			fw.updateAvgWait(elapsed)
			return ctx.Err()
		default:
			time.Sleep(5 * time.Millisecond)
	}
	}
}

func (fw *FixedWindowLimiter) Limit() int {
	return fw.rate
}

func (fw *FixedWindowLimiter) Burst() int {
	return fw.rate
}

func (fw *FixedWindowLimiter) GetMetrics() RateLimiterMetrics {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.metrics
}

func (fw *FixedWindowLimiter) updateAvgWait(elapsed time.Duration) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	total := fw.metrics.AllowedRequests + fw.metrics.DeniedRequests
	if total == 0 {
		return
	}
	fw.metrics.AverageWaitTime += (elapsed - fw.metrics.AverageWaitTime) / time.Duration(total)
}

func (fw *FixedWindowLimiter) Reset() {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fw.windowStart = time.Now()
	fw.requestCount = 0
	fw.metrics = RateLimiterMetrics{}
}

// ----------------------------------------------------------
// Factory
// ----------------------------------------------------------

type RateLimiterFactory struct{}

type RateLimiterConfig struct {
	Algorithm  string
	Rate       int
	Burst      int
	WindowSize time.Duration
}

func NewRateLimiterFactory() *RateLimiterFactory {
	return &RateLimiterFactory{}
}

func (f *RateLimiterFactory) CreateLimiter(cfg RateLimiterConfig) (RateLimiter, error) {
	switch cfg.Algorithm {
	case "token_bucket":
		if cfg.Rate <= 0 || cfg.Burst <= 0 {
			return nil, errors.New("invalid config")
		}
		return NewTokenBucketLimiter(cfg.Rate, cfg.Burst), nil
	case "sliding_window":
		if cfg.Rate <= 0 || cfg.WindowSize <= 0 {
			return nil, errors.New("invalid config")
		}
		return NewSlidingWindowLimiter(cfg.Rate, cfg.WindowSize), nil
	case "fixed_window":
		if cfg.Rate <= 0 || cfg.WindowSize <= 0 {
			return nil, errors.New("invalid config")
		}
		return NewFixedWindowLimiter(cfg.Rate, cfg.WindowSize), nil
	default:
		return nil, errors.New("unknown algorithm")
}
}

// ----------------------------------------------------------
// Middleware
// ----------------------------------------------------------

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

// ----------------------------------------------------------
// Main
// ----------------------------------------------------------

func main() {
	factory := NewRateLimiterFactory()
	cfg := RateLimiterConfig{Algorithm: "token_bucket", Rate: 10, Burst: 5}
	limiter, _ := factory.CreateLimiter(cfg)
	fmt.Println("TokenBucket Allow?", limiter.Allow())
	for i := 0; i < 7; i++ {
		fmt.Printf("Req %d, allowed? %v\n", i, limiter.Allow())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := limiter.Wait(ctx)
	fmt.Println("Wait with ctx:", err)
	metrics := limiter.GetMetrics()
	fmt.Printf("Metrics: %+v\n", metrics)
}
