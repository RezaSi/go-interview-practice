# Hints for Challenge 2: Middleware & Request/Response Handling

## Hint 1: Understanding Middleware
Middleware in Echo are functions that execute before or after your route handlers. They follow this pattern:

```go
func middlewareFunction(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Do something before the handler
        
        err := next(c) // Call the next middleware/handler
        
        // Do something after the handler
        return err
    }
}
```

Register middleware with `e.Use(middlewareFunction)`.


## Hint 2: Built-in Middleware SetupEcho provides many built-in middleware. Set them up in `setupMiddleware`:

```go
// Logging with custom format
e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: "custom format string with ${variables}",
}))

// Recovery from panics
e.Use(middleware.Recover())

// CORS for cross-origin requests
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
}))
```


## Hint 3: Request ID MiddlewareGenerate unique IDs for request tracking:

```go
func requestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        requestID := uuid.New().String()
        c.Response().Header().Set("X-Request-ID", requestID)
        return next(c)
    }
}
```


## Hint 4: API Key AuthenticationCheck for API keys in headers:

```go
func apiKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        apiKey := c.Request().Header.Get("X-API-Key")
        
        if apiKey == "" {
            return c.JSON(http.StatusUnauthorized, ErrorResponse{
                Status: "error",
                Message: "API key required",
            })
        }
        
        if !validAPIKeys[apiKey] {
            return c.JSON(http.StatusUnauthorized, ErrorResponse{
                Status: "error", 
                Message: "Invalid API key",
            })
        }
        
        return next(c)
    }
}
```


## Hint 5: Rate Limiting ImplementationTrack requests per IP with time windows:

```go
type RateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.RWMutex
    limit    int
    window   time.Duration
}

func (rl *RateLimiter) Allow(clientIP string) bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    
    now := time.Now()
    
    // Clean old requests outside window
    // Check if limit exceeded
    // Add current request if allowed
    
    return allowed
}
```


## Hint 6: Route Groups with MiddlewareApply middleware to specific route groups:

```go
// Health endpoint - no authentication
e.GET("/health", healthHandler)

// API endpoints - with authentication
api := e.Group("/posts")
api.Use(apiKeyMiddleware) // Apply to all routes in group

api.GET("", getAllPosts)
api.POST("", createPost)
// ... other routes
```


## Hint 7: Getting Client IPExtract client IP for rate limiting:

```go
clientIP := c.RealIP()
if clientIP == "" {
    clientIP = c.Request().Header.Get("X-Real-IP")
}
if clientIP == "" {
    clientIP = "unknown"
}
```


## Hint 8: Middleware Order MattersSet up middleware in the correct order:

```go
func setupMiddleware(e *echo.Echo) {
    e.Use(middleware.Logger())    // First - log everything
    e.Use(middleware.Recover())   // Second - catch panics
    e.Use(middleware.CORS())      // Third - handle CORS
    e.Use(requestIDMiddleware)    // Fourth - generate request IDs
    e.Use(rateLimitMiddleware)    // Fifth - rate limiting
    // Route-specific middleware applied in route groups
}
```


## Hint 9: Post ValidationValidate all required fields for blog posts:

```go
if strings.TrimSpace(post.Title) == "" {
    return c.JSON(http.StatusBadRequest, ErrorResponse{
        Status: "error",
        Message: "Title is required",
    })
}

if strings.TrimSpace(post.Content) == "" {
    return c.JSON(http.StatusBadRequest, ErrorResponse{
        Status: "error", 
        Message: "Content is required",
    })
}

if strings.TrimSpace(post.Author) == "" {
    return c.JSON(http.StatusBadRequest, ErrorResponse{
        Status: "error",
        Message: "Author is required", 
    })
}
```


## Hint 10: Timestamp ManagementHandle CreatedAt and UpdatedAt properly:

```go
// For new posts
now := time.Now().Format(time.RFC3339)
post.CreatedAt = now
post.UpdatedAt = now

// For updates - preserve CreatedAt
updatedPost.CreatedAt = posts[index].CreatedAt
updatedPost.UpdatedAt = time.Now().Format(time.RFC3339)
```

## 🔍 Debugging Tips
Middleware not executing: Check the order of `e.Use()` calls and ensure middleware is registered before routes.

Rate limiting not working: Verify IP extraction and time window calculations in your rate limiter.

API key authentication failing: Check header name (`X-API-Key`) and ensure valid keys are in your map.

CORS issues: Make sure CORS middleware allows the required origins, methods, and headers.

Request ID missing: Ensure request ID middleware runs before logging middleware to include IDs in logs.