# Hints for Challenge 29: Rate Limiter Implementation

## Hint 1: Token Bucket Algorithm
The token bucket allows controlled bursts while maintaining average rate limits:

**Core Concept:**
- **Bucket**: Holds tokens (capacity = burst size)
- **Refill Rate**: Tokens added per second (rate limit)
- **Consumption**: Each request consumes 1+ tokens
- **Burst Handling**: Full bucket allows temporary spikes

**Key Implementation Points:**
- Track current token count and last refill time
- Calculate tokens to add based on elapsed time
- Allow request only if sufficient tokens available

```go
type TokenBucket struct {
    tokens     float64    // Current tokens
    capacity   float64    // Max tokens (burst)
    refillRate float64    // Tokens per second
    lastRefill time.Time  // Last refill timestamp
}

// Refill based on elapsed time
tokensToAdd := elapsed.Seconds() * refillRate
tokens = min(capacity, tokens + tokensToAdd)
```
```

## Hint 2: Token Bucket - Allow and Wait Methods
Implement the core rate limiting logic:
```go
func (tb *TokenBucketLimiter) Allow() bool {
    return tb.AllowN(1)
}

func (tb *TokenBucketLimiter) AllowN(n int) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    
    tb.refill()
    tb.totalReqs++
    
    if tb.tokens >= float64(n) {
        tb.tokens -= float64(n)
        tb.allowedReqs++
        return true
    }
    
    return false
}

func (tb *TokenBucketLimiter) Wait(ctx context.Context) error {
    return tb.WaitN(ctx, 1)
}

func (tb *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
    for {
        if tb.AllowN(n) {
            return nil
        }
        
        select {
        case &lt;-ctx.Done():
            return ctx.Err()
        case &lt;-time.After(time.Millisecond * 10):
            // Small sleep to avoid busy waiting
        }
    }
}

func (tb *TokenBucketLimiter) Limit() int {
    return int(tb.refillRate)
}

func (tb *TokenBucketLimiter) Burst() int {
    return int(tb.capacity)
}

func (tb *TokenBucketLimiter) Reset() {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    tb.tokens = tb.capacity
    tb.lastRefill = time.Now()
    tb.totalReqs = 0
    tb.allowedReqs = 0
}
```

## Hint 3: Sliding Window Algorithm
More precise than fixed windows, avoids boundary effects:

**Key Insight:** 
- Track timestamps of recent requests in a sliding time window
- Before each request, remove timestamps older than window size
- Allow request if remaining count < rate limit

**Advantages over Fixed Window:**
- No burst at window boundaries
- More accurate rate limiting
- Smooths traffic over time

```go
// Clean old requests outside window
cutoff := now.Add(-windowSize)
requests = removeOlderThan(requests, cutoff)

// Check if within rate limit
if len(requests) < rateLimit {
    requests = append(requests, now)
    return true  // Allow
}
```

## Hint 4: Fixed Window Rate Limiter
Implement simple counter-based rate limiting:
```go
type FixedWindowLimiter struct {
    mu          sync.Mutex
    rate        int
    windowSize  time.Duration
    count       int
    windowStart time.Time
    totalReqs   int64
    allowedReqs int64
}

func NewFixedWindowLimiter(rate int, windowSize time.Duration) RateLimiter {
    return &FixedWindowLimiter{
        rate:        rate,
        windowSize:  windowSize,
        windowStart: time.Now(),
    }
}

func (fw *FixedWindowLimiter) resetIfNeeded() {
    now := time.Now()
    if now.Sub(fw.windowStart) >= fw.windowSize {
        fw.count = 0
        fw.windowStart = now
    }
}

func (fw *FixedWindowLimiter) AllowN(n int) bool {
    fw.mu.Lock()
    defer fw.mu.Unlock()
    
    fw.resetIfNeeded()
    fw.totalReqs++
    
    if fw.count+n <= fw.rate {
        fw.count += n
        fw.allowedReqs++
        return true
    }
    
    return false
}

func (fw *FixedWindowLimiter) Wait(ctx context.Context) error {
    return fw.WaitN(ctx, 1)
}

func (fw *FixedWindowLimiter) WaitN(ctx context.Context, n int) error {
    for {
        if fw.AllowN(n) {
            return nil
        }
        
        // Calculate time until next window
        fw.mu.Lock()
        nextWindow := fw.windowStart.Add(fw.windowSize)
        fw.mu.Unlock()
        
        waitTime := time.Until(nextWindow)
        if waitTime <= 0 {
            continue
        }
        
        select {
        case &lt;-ctx.Done():
            return ctx.Err()
        case &lt;-time.After(waitTime):
            // Window has reset, try again
        }
    }
}
```

## Hint 5: Rate Limiter Factory Pattern
Create flexible factory for different rate limiter types:
```go
type RateLimiterConfig struct {
    Algorithm  string
    Rate       int
    Burst      int
    WindowSize time.Duration
}

type RateLimiterFactory struct{}

func NewRateLimiterFactory() *RateLimiterFactory {
    return &RateLimiterFactory{}
}

func (f *RateLimiterFactory) CreateLimiter(config RateLimiterConfig) (RateLimiter, error) {
    switch config.Algorithm {
    case "token_bucket":
        if config.Burst <= 0 {
            config.Burst = config.Rate
        }
        return NewTokenBucketLimiter(config.Rate, config.Burst), nil
        
    case "sliding_window":
        if config.WindowSize == 0 {
            config.WindowSize = time.Second
        }
        return NewSlidingWindowLimiter(config.Rate, config.WindowSize), nil
        
    case "fixed_window":
        if config.WindowSize == 0 {
            config.WindowSize = time.Second
        }
        return NewFixedWindowLimiter(config.Rate, config.WindowSize), nil
        
    default:
        return nil, fmt.Errorf("unknown algorithm: %s", config.Algorithm)
    }
}
```

## Hint 6: HTTP Middleware Implementation
Add rate limiting to HTTP handlers:
```go
func RateLimitMiddleware(limiter RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
                w.Header().Set("X-RateLimit-Remaining", "0")
                w.Header().Set("Retry-After", "1")
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
            next.ServeHTTP(w, r)
        })
    }
}

// Per-IP rate limiting middleware
func PerIPRateLimitMiddleware(factory *RateLimiterFactory, config RateLimiterConfig) func(http.Handler) http.Handler {
    limiters := sync.Map{}
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := getClientIP(r)
            
            limiterInterface, _ := limiters.LoadOrStore(ip, func() RateLimiter {
                limiter, _ := factory.CreateLimiter(config)
                return limiter
            }())
            
            limiter := limiterInterface.(RateLimiter)
            
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func getClientIP(r *http.Request) string {
    forwarded := r.Header.Get("X-Forwarded-For")
    if forwarded != "" {
        return strings.Split(forwarded, ",")[0]
    }
    return r.RemoteAddr
}
```

## Hint 7: Metrics and Monitoring
Add metrics collection for rate limiter performance:
```go
type RateLimiterMetrics struct {
    TotalRequests   int64
    AllowedRequests int64
    DeniedRequests  int64
    AverageWaitTime time.Duration
    mu              sync.RWMutex
}

func (m *RateLimiterMetrics) RecordRequest(allowed bool, waitTime time.Duration) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.TotalRequests++
    if allowed {
        m.AllowedRequests++
    } else {
        m.DeniedRequests++
    }
    
    // Update moving average of wait time
    if m.TotalRequests == 1 {
        m.AverageWaitTime = waitTime
    } else {
        // Simple moving average
        m.AverageWaitTime = time.Duration(
            (int64(m.AverageWaitTime)*9 + int64(waitTime)) / 10,
        )
    }
}

func (m *RateLimiterMetrics) GetStats() (total, allowed, denied int64, avgWait time.Duration) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.TotalRequests, m.AllowedRequests, m.DeniedRequests, m.AverageWaitTime
}

func (m *RateLimiterMetrics) SuccessRate() float64 {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    if m.TotalRequests == 0 {
        return 0.0
    }
    return float64(m.AllowedRequests) / float64(m.TotalRequests)
}

// Enhanced rate limiter with metrics
type MetricsRateLimiter struct {
    limiter RateLimiter
    metrics *RateLimiterMetrics
}

func NewMetricsRateLimiter(limiter RateLimiter) *MetricsRateLimiter {
    return &MetricsRateLimiter{
        limiter: limiter,
        metrics: &RateLimiterMetrics{},
    }
}

func (m *MetricsRateLimiter) Allow() bool {
    start := time.Now()
    allowed := m.limiter.Allow()
    waitTime := time.Since(start)
    m.metrics.RecordRequest(allowed, waitTime)
    return allowed
}
```

## Key Rate Limiting Concepts:
- **Token Bucket**: Allows bursts while maintaining average rate
- **Sliding Window**: Precise rate limiting without boundary effects  
- **Fixed Window**: Simple and efficient but allows bursts at boundaries
- **Thread Safety**: Use mutexes for concurrent access
- **Context Cancellation**: Support timeouts in Wait methods
- **HTTP Integration**: Middleware for web applications with proper headers
</rewritten_file> 