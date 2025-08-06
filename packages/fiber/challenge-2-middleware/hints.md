# Hints for Challenge 2: Middleware & Request/Response Handling

## Hint 1: Setting up Fiber App

Start with a clean Fiber app:

```go
app := fiber.New(fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        // Custom error handler
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "error": "Internal server error",
        })
    },
})
```

## Hint 2: Middleware Order

Apply middleware in the correct order:

```go
app.Use(ErrorHandlerMiddleware())
app.Use(RequestIDMiddleware())
app.Use(LoggingMiddleware())
app.Use(CORSMiddleware())
app.Use(RateLimitMiddleware())
```

## Hint 3: Request ID Generation

Use UUID for unique request IDs:

```go
func RequestIDMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        requestID := uuid.New().String()
        c.Locals("request_id", requestID)
        c.Set("X-Request-ID", requestID)
        return c.Next()
    }
}
```

## Hint 4: Custom Logging

Log requests with timing information:

```go
func LoggingMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        requestID := c.Locals("request_id").(string)
        
        err := c.Next()
        
        duration := time.Since(start)
        fmt.Printf("[%s] %s %s %d %v %s\n",
            requestID,
            c.Method(),
            c.Path(),
            c.Response().StatusCode(),
            duration,
            c.IP(),
        )
        
        return err
    }
}
```

## Hint 5: Rate Limiting

Implement simple in-memory rate limiting:

```go
var rateLimitMap = make(map[string][]time.Time)
var rateLimitMutex sync.RWMutex

func RateLimitMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        ip := c.IP()
        now := time.Now()
        
        rateLimitMutex.Lock()
        defer rateLimitMutex.Unlock()
        
        // Clean old entries
        if times, exists := rateLimitMap[ip]; exists {
            var validTimes []time.Time
            for _, t := range times {
                if now.Sub(t) < time.Minute {
                    validTimes = append(validTimes, t)
                }
            }
            rateLimitMap[ip] = validTimes
        }
        
        // Check limit
        if len(rateLimitMap[ip]) >= 100 {
            return c.Status(429).JSON(fiber.Map{
                "success": false,
                "error": "Rate limit exceeded",
            })
        }
        
        // Add current request
        rateLimitMap[ip] = append(rateLimitMap[ip], now)
        
        return c.Next()
    }
}
```

## Hint 6: Route Groups

Organize routes using groups:

```go
// Public routes
public := app.Group("/")
public.Get("/ping", pingHandler)
public.Get("/articles", getArticlesHandler)
public.Get("/articles/:id", getArticleHandler)

// Protected routes
protected := app.Group("/", AuthMiddleware())
protected.Post("/articles", createArticleHandler)
protected.Put("/articles/:id", updateArticleHandler)
protected.Delete("/articles/:id", deleteArticleHandler)
protected.Get("/admin/stats", getStatsHandler)
```