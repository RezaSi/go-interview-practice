# Learning: Advanced Fiber Middleware Patterns

## 🌟 **What is Middleware?**

Middleware functions execute during the request-response cycle and can:
- Execute code before the route handler
- Modify request or response objects
- End the request-response cycle
- Call the next middleware function

### **Middleware in Fiber**
Fiber middleware uses a similar pattern to Express.js with the `c.Next()` function:

```go
func MyMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Before handler logic
        
        err := c.Next() // Execute next middleware/handler
        
        // After handler logic
        
        return err
    }
}
```

## 🔄 **Middleware Execution Order**

Middleware executes in the order it's registered:

```
Request → MW1 → MW2 → MW3 → Handler → MW3 → MW2 → MW1 → Response
```

### **Best Practice Order**
1. **Error Recovery** - Catch panics first
2. **Request ID** - Track requests
3. **Logging** - Log with request ID
4. **CORS** - Handle cross-origin requests
5. **Rate Limiting** - Protect against abuse
6. **Authentication** - Verify users
7. **Route Handlers** - Business logic

## 🛠️ **Essential Middleware Patterns**

### **1. Request ID Middleware**
Track requests across your application:

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

### **2. Logging Middleware**
Monitor request performance:

```go
func LoggingMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        err := c.Next()
        
        duration := time.Since(start)
        log.Printf("%s %s %d %v",
            c.Method(),
            c.Path(),
            c.Response().StatusCode(),
            duration,
        )
        
        return err
    }
}
```

### **3. CORS Middleware**
Enable cross-origin requests:

```go
func CORSMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Set("Access-Control-Allow-Origin", "*")
        c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
        
        if c.Method() == "OPTIONS" {
            return c.SendStatus(204)
        }
        
        return c.Next()
    }
}
```

### **4. Authentication Middleware**
Protect routes with API keys or tokens:

```go
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        apiKey := c.Get("X-API-Key")
        
        if !isValidAPIKey(apiKey) {
            return c.Status(401).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }
        
        c.Locals("api_key", apiKey)
        return c.Next()
    }
}
```

## 📊 **Rate Limiting Strategies**

### **Fixed Window**
Simple but can allow bursts:

```go
func FixedWindowRateLimit() fiber.Handler {
    requests := make(map[string]int)
    lastReset := time.Now()
    
    return func(c *fiber.Ctx) error {
        now := time.Now()
        ip := c.IP()
        
        // Reset window every minute
        if now.Sub(lastReset) >= time.Minute {
            requests = make(map[string]int)
            lastReset = now
        }
        
        requests[ip]++
        if requests[ip] > 100 {
            return c.Status(429).JSON(fiber.Map{
                "error": "Rate limit exceeded",
            })
        }
        
        return c.Next()
    }
}
```

### **Sliding Window**
More accurate but uses more memory:

```go
func SlidingWindowRateLimit() fiber.Handler {
    requests := make(map[string][]time.Time)
    
    return func(c *fiber.Ctx) error {
        now := time.Now()
        ip := c.IP()
        
        // Clean old requests
        var validRequests []time.Time
        for _, reqTime := range requests[ip] {
            if now.Sub(reqTime) < time.Minute {
                validRequests = append(validRequests, reqTime)
            }
        }
        
        if len(validRequests) >= 100 {
            return c.Status(429).JSON(fiber.Map{
                "error": "Rate limit exceeded",
            })
        }
        
        requests[ip] = append(validRequests, now)
        return c.Next()
    }
}
```

## 🔒 **Error Handling Patterns**

### **Centralized Error Handler**
Handle all errors in one place:

```go
func ErrorHandlerMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic recovered: %v", r)
                c.Status(500).JSON(fiber.Map{
                    "success": false,
                    "error": "Internal server error",
                })
            }
        }()
        
        return c.Next()
    }
}
```

### **Custom Error Types**
Create specific error types:

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e APIError) Error() string {
    return e.Message
}

func ValidationError(message string) APIError {
    return APIError{
        Code:    400,
        Message: message,
    }
}
```

## 🎯 **Context and State Management**

### **Storing Data in Context**
Share data between middleware:

```go
// Store data
c.Locals("user_id", 123)
c.Locals("request_start", time.Now())

// Retrieve data
userID := c.Locals("user_id").(int)
startTime := c.Locals("request_start").(time.Time)
```

### **Request Scoped Data**
Keep data tied to specific requests:

```go
type RequestContext struct {
    UserID    int
    RequestID string
    StartTime time.Time
}

func ContextMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        ctx := &RequestContext{
            RequestID: uuid.New().String(),
            StartTime: time.Now(),
        }
        
        c.Locals("ctx", ctx)
        return c.Next()
    }
}
```

## 📈 **Performance Considerations**

### **Efficient Middleware**
- Avoid heavy computations in middleware
- Use connection pooling for external services
- Cache frequently accessed data
- Clean up resources properly

### **Memory Management**
```go
func EfficientMiddleware() fiber.Handler {
    // Initialize outside the handler
    cache := make(map[string]interface{})
    
    return func(c *fiber.Ctx) error {
        // Lightweight operations only
        key := c.Get("Cache-Key")
        if data, exists := cache[key]; exists {
            c.Locals("cached_data", data)
        }
        
        return c.Next()
    }
}
```

## 🔧 **Testing Middleware**

### **Unit Testing**
Test middleware in isolation:

```go
func TestRequestIDMiddleware(t *testing.T) {
    app := fiber.New()
    app.Use(RequestIDMiddleware())
    app.Get("/test", func(c *fiber.Ctx) error {
        return c.SendString("OK")
    })
    
    req := httptest.NewRequest("GET", "/test", nil)
    resp, _ := app.Test(req)
    
    assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
}
```

### **Integration Testing**
Test middleware chains:

```go
func TestMiddlewareChain(t *testing.T) {
    app := fiber.New()
    app.Use(RequestIDMiddleware())
    app.Use(LoggingMiddleware())
    app.Use(AuthMiddleware())
    
    // Test with valid auth
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("X-API-Key", "valid-key")
    
    resp, _ := app.Test(req)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## 🎯 **Best Practices**

1. **Keep middleware focused** - One responsibility per middleware
2. **Order matters** - Place middleware in logical order
3. **Handle errors gracefully** - Don't let middleware crash the app
4. **Use context for sharing** - Store request-scoped data in context
5. **Test thoroughly** - Unit test each middleware
6. **Monitor performance** - Track middleware execution time
7. **Clean up resources** - Release resources in defer statements

## 📚 **Next Steps**

After mastering middleware patterns:
1. **Validation & Error Handling** - Input validation and error responses
2. **Authentication & Authorization** - JWT tokens and role-based access
3. **Database Integration** - Connecting to databases
4. **Testing Strategies** - Comprehensive testing approaches