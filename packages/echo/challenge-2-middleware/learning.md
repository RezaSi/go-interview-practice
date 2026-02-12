# Learning: Advanced Echo Middleware Patterns

## 🌟 **What is Middleware?**

Middleware in Echo are functions that execute during the HTTP request-response cycle. They can:
- **Intercept** requests before they reach your handlers
- **Modify** requests and responses
- **Add** functionality like logging, authentication, CORS
- **Handle** errors and panics globally

### **The Middleware Chain**
```
Request → Middleware1 → Middleware2 → Handler → Middleware2 → Middleware1 → Response
```

## 🔗 **Middleware Execution Flow**

### **Basic Middleware Structure**
```go
func MyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // BEFORE: Code here runs before the request is processed
        fmt.Println("Before request")
        
        err := next(c) // Call the next middleware/handler
        
        // AFTER: Code here runs after the request is processed  
        fmt.Println("After request")
        return err
    }
}
```

### **Registering Middleware**
```go
// Global middleware (applies to all routes)
e.Use(MyMiddleware)

// Route-specific middleware
e.GET("/protected", protectedHandler, AuthMiddleware)

// Group middleware
api := e.Group("/api")
api.Use(APIKeyMiddleware)
```

## 🛠️ **Built-in Middleware**

### **Logger Middleware**
```go
e.Use(middleware.Logger())

// Custom logger format
e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
}))
```

### **Recovery Middleware**
```go
e.Use(middleware.Recover())

// Custom recovery
e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
    StackSize: 1 << 10, // 1 KB
}))
```

### **CORS Middleware**
```go
e.Use(middleware.CORS())

// Custom CORS
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"https://example.com"},
    AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
}))
```

## 🔐 **Authentication Middleware**

### **API Key Authentication**
```go
func APIKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        apiKey := c.Request().Header.Get("X-API-Key")
        
        if apiKey == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "API key required",
            })
        }
        
        if !isValidAPIKey(apiKey) {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "Invalid API key",
            })
        }
        
        return next(c)
    }
}
```

### **JWT Authentication**
```go
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")
        
        if token == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "Authorization header required",
            })
        }
        
        // Validate JWT token
        claims, err := validateJWT(token)
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "error": "Invalid token",
            })
        }
        
        // Store user info in context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        
        return next(c)
    }
}
```

## 📊 **Logging Middleware**

### **Request ID Middleware**
```go
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        requestID := generateRequestID()
        c.Set("request_id", requestID)
        c.Response().Header().Set("X-Request-ID", requestID)
        return next(c)
    }
}
```

### **Custom Logging Middleware**
```go
func CustomLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        
        err := next(c)
        
        duration := time.Since(start)
        requestID := c.Get("request_id")
        
        log.Printf("[%s] %s %s %d %v",
            requestID,
            c.Request().Method,
            c.Request().URL.Path,
            c.Response().Status,
            duration,
        )
        
        return err
    }
}
```

## 🚦 **Rate Limiting Middleware**

### **Simple Rate Limiter**
```go
var requestCounts = make(map[string]int)
var mutex = sync.RWMutex{}

func RateLimitMiddleware(limit int) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            ip := c.RealIP()
            
            mutex.Lock()
            count := requestCounts[ip]
            if count >= limit {
                mutex.Unlock()
                return c.JSON(http.StatusTooManyRequests, map[string]string{
                    "error": "Rate limit exceeded",
                })
            }
            requestCounts[ip]++
            mutex.Unlock()
            
            return next(c)
        }
    }
}
```

## 🔧 **Error Handling Middleware**

### **Global Error Handler**
```go
func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        err := next(c)
        
        if err != nil {
            // Log the error
            log.Printf("Error: %v", err)
            
            // Return consistent error response
            if he, ok := err.(*echo.HTTPError); ok {
                return c.JSON(he.Code, map[string]interface{}{
                    "error":   he.Message,
                    "status":  he.Code,
                })
            }
            
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": "Internal server error",
            })
        }
        
        return nil
    }
}
```

## 🎯 **Middleware Best Practices**

### **Order Matters**
```go
e.Use(middleware.Logger())      // 1. Log all requests
e.Use(middleware.Recover())     // 2. Recover from panics
e.Use(middleware.CORS())        // 3. Handle CORS
e.Use(RequestIDMiddleware)      // 4. Add request IDs
e.Use(RateLimitMiddleware(100)) // 5. Rate limiting
```

### **Context Usage**
```go
// Store data in context
c.Set("user_id", userID)

// Retrieve data from context
userID := c.Get("user_id").(string)
```

### **Conditional Middleware**
```go
func ConditionalMiddleware(condition bool) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if condition {
                // Apply middleware logic
                return middlewareLogic(c, next)
            }
            return next(c)
        }
    }
}
```

## 🏗️ **Middleware Composition**

### **Combining Multiple Middleware**
```go
// Create middleware chain
authChain := []echo.MiddlewareFunc{
    RequestIDMiddleware,
    CustomLoggerMiddleware,
    APIKeyMiddleware,
}

// Apply to route group
api := e.Group("/api", authChain...)
```

### **Middleware Factory Pattern**
```go
func AuthMiddleware(requiredRole string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            userRole := getUserRole(c)
            if userRole != requiredRole {
                return c.JSON(http.StatusForbidden, map[string]string{
                    "error": "Insufficient permissions",
                })
            }
            return next(c)
        }
    }
}

// Usage
e.GET("/admin", adminHandler, AuthMiddleware("admin"))
```

## 🧪 **Testing Middleware**

### **Unit Testing**
```go
func TestAPIKeyMiddleware(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("X-API-Key", "valid-key")
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    
    handler := APIKeyMiddleware(func(c echo.Context) error {
        return c.String(http.StatusOK, "success")
    })
    
    err := handler(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

## 📚 **Advanced Patterns**

### **Middleware with Configuration**
```go
type LoggerConfig struct {
    Format    string
    Output    io.Writer
    SkipPaths []string
}

func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Skip logging for certain paths
            for _, path := range config.SkipPaths {
                if c.Request().URL.Path == path {
                    return next(c)
                }
            }
            
            // Custom logging logic
            return next(c)
        }
    }
}
```

### **Async Middleware**
```go
func AsyncLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        
        err := next(c)
        
        // Log asynchronously
        go func() {
            duration := time.Since(start)
            log.Printf("Request processed in %v", duration)
        }()
        
        return err
    }
}
```

Middleware is the backbone of Echo applications, enabling clean separation of concerns and reusable functionality across your entire application!