# Hints for Challenge 2: Middleware & Request/Response Handling

## Hint 1: Understanding Middleware

Middleware functions run before your route handlers. They follow this pattern:

```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Pre-processing logic here
        
        c.Next() // IMPORTANT: Call this to continue the chain
        
        // Post-processing logic here (runs after handler)
    }
}
```

## Hint 2: Request ID Middleware

Create a middleware that adds a unique request ID to each request:

```go
import "github.com/google/uuid"

func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

## Hint 3: Logging Middleware Structure

Build a logging middleware that tracks request details:

```go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        log.Printf("[%s] %s %s - %v", 
            c.Request.Method, 
            c.Request.URL.Path, 
            c.ClientIP(), 
            duration)
    }
}
```

## Hint 4: CORS Middleware Implementation

Handle Cross-Origin Resource Sharing for web client access:

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Request-ID")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

## Hint 5: Rate Limiting with Memory Store

Implement basic rate limiting using an in-memory store:

```go
import "golang.org/x/time/rate"

var rateLimiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

func RateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
        
        mu.Lock()
        limiter, exists := rateLimiters[clientIP]
        if !exists {
            limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond*2)
            rateLimiters[clientIP] = limiter
        }
        mu.Unlock()
        
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## Hint 6: Setting Up Router with Middleware

Apply middleware to your router in the correct order:

```go
func setupRouter() *gin.Engine {
    router := gin.New() // Start with clean router
    
    // Add middleware in order
    router.Use(LoggingMiddleware())
    router.Use(RequestIDMiddleware())
    router.Use(CORSMiddleware())
    router.Use(RateLimitMiddleware(100)) // 100 requests per second
    
    // Add your routes
    router.GET("/users", getUsers)
    router.POST("/users", createUser)
    
    return router
}
``` 