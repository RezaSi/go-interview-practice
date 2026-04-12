# Learning: Advanced Gin Middleware Patterns

## 🌟 **What is Middleware?**

Middleware in Gin is code that runs **before** and **after** your route handlers. Think of it as a chain of functions that can:
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
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // BEFORE: Code here runs before the request is processed
        fmt.Println("Before request")
        
        c.Next() // Call the next middleware/handler
        
        // AFTER: Code here runs after the request is processed  
        fmt.Println("After request")
    }
}
```

### **Middleware Registration**
```go
router := gin.Default()

// Global middleware (applies to all routes)
router.Use(MyMiddleware())

// Group middleware (applies to specific route groups)
api := router.Group("/api")
api.Use(AuthMiddleware())
{
    api.GET("/users", getUsers)
}

// Route-specific middleware
router.GET("/admin", AdminMiddleware(), adminHandler)
```

## 🔐 **Authentication Middleware**

### **API Key Authentication**
```go
func APIKeyAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        
        if apiKey == "" {
            c.JSON(401, gin.H{"error": "API key required"})
            c.Abort() // Stop middleware chain
            return
        }
        
        // Validate API key
        if !isValidAPIKey(apiKey) {
            c.JSON(401, gin.H{"error": "Invalid API key"})
            c.Abort()
            return
        }
        
        // Store user info in context
        c.Set("user_id", getUserIDFromAPIKey(apiKey))
        c.Set("user_role", getUserRole(apiKey))
        
        c.Next() // Continue to next middleware/handler
    }
}
```

### **Role-Based Access Control**
```go
func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("user_role")
        
        if userRole != requiredRole {
            c.JSON(403, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Usage: Require admin role
router.DELETE("/users/:id", RequireRole("admin"), deleteUser)
```

## 📝 **Logging Middleware**

### **Custom Request Logger**
```go
func CustomLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()
        
        // Calculate request duration
        duration := time.Since(start)
        
        // Log request details
        log.Printf("[%s] %s %s %d %v %s",
            c.GetString("request_id"),
            c.Request.Method,
            path,
            c.Writer.Status(),
            duration,
            c.ClientIP(),
        )
    }
}
```

### **Structured Logging with Context**
```go
func StructuredLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        
        entry := map[string]interface{}{
            "request_id": c.GetString("request_id"),
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "status":     c.Writer.Status(),
            "duration":   duration.Milliseconds(),
            "ip":         c.ClientIP(),
            "user_agent": c.Request.UserAgent(),
        }
        
        if c.Writer.Status() >= 400 {
            log.Printf("ERROR: %+v", entry)
        } else {
            log.Printf("INFO: %+v", entry)
        }
    }
}
```

## 🌐 **CORS Middleware**

### **Understanding CORS**
Cross-Origin Resource Sharing (CORS) allows web pages from one domain to access resources from another domain.

### **Custom CORS Implementation**
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Define allowed origins
        allowedOrigins := map[string]bool{
            "http://localhost:3000":  true,
            "https://myapp.com":      true,
        }
        
        if allowedOrigins[origin] {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")
        c.Header("Access-Control-Allow-Credentials", "true")
        
        // Handle preflight requests
        if c.Request.Method == "OPTIONS" {
            c.Status(204)
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## ⏱️ **Rate Limiting Middleware**

### **Simple In-Memory Rate Limiter**
```go
type RateLimiter struct {
    visitors map[string]*visitor
    mu       sync.RWMutex
}

type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

func NewRateLimiter(requests int, duration time.Duration) *RateLimiter {
    rl := &RateLimiter{
        visitors: make(map[string]*visitor),
    }
    
    // Clean up old visitors every minute
    go rl.cleanupVisitors()
    
    return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    v, exists := rl.visitors[ip]
    if !exists {
        limiter := rate.NewLimiter(rate.Every(time.Minute), 100) // 100 requests per minute
        rl.visitors[ip] = &visitor{limiter, time.Now()}
        return limiter
    }
    
    v.lastSeen = time.Now()
    return v.limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        limiter := rl.getVisitor(c.ClientIP())
        
        if !limiter.Allow() {
            c.Header("X-RateLimit-Limit", "100")
            c.Header("X-RateLimit-Remaining", "0")
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## 🆔 **Request ID Middleware**

### **UUID Generation**
```go
import "github.com/google/uuid"

func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if request ID already exists in header
        requestID := c.GetHeader("X-Request-ID")
        
        if requestID == "" {
            // Generate new UUID
            requestID = uuid.New().String()
        }
        
        // Store in context for other middleware/handlers
        c.Set("request_id", requestID)
        
        // Add to response headers
        c.Header("X-Request-ID", requestID)
        
        c.Next()
    }
}
```

## ❌ **Error Handling Middleware**

### **Centralized Error Handler**
```go
type APIError struct {
    StatusCode int    `json:"-"`
    Code       string `json:"code"`
    Message    string `json:"message"`
    Details    string `json:"details,omitempty"`
}

func (e APIError) Error() string {
    return e.Message
}

func ErrorHandler() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        var apiErr APIError
        
        switch err := recovered.(type) {
        case APIError:
            apiErr = err
        case error:
            apiErr = APIError{
                StatusCode: 500,
                Code:       "INTERNAL_ERROR",
                Message:    "Internal server error",
                Details:    err.Error(),
            }
        default:
            apiErr = APIError{
                StatusCode: 500,
                Code:       "PANIC",
                Message:    "Internal server error",
                Details:    fmt.Sprintf("%v", recovered),
            }
        }
        
        c.JSON(apiErr.StatusCode, gin.H{
            "success":    false,
            "error":      apiErr.Message,
            "code":       apiErr.Code,
            "request_id": c.GetString("request_id"),
        })
    })
}
```

## 🔍 **Content Type Validation**

### **JSON Content Type Middleware**
```go
func RequireJSON() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == "POST" || c.Request.Method == "PUT" {
            contentType := c.GetHeader("Content-Type")
            
            if !strings.HasPrefix(contentType, "application/json") {
                c.JSON(415, gin.H{
                    "error":   "Content-Type must be application/json",
                    "code":    "INVALID_CONTENT_TYPE",
                })
                c.Abort()
                return
            }
        }
        
        c.Next()
    }
}
```

## 🔄 **Context Data Sharing**

### **Passing Data Between Middleware**
```go
// Setting data in middleware
func SetUserMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authenticate user...
        user := User{ID: 123, Name: "John Doe"}
        
        c.Set("current_user", user)
        c.Set("user_id", user.ID)
        
        c.Next()
    }
}

// Getting data in handlers
func getUserHandler(c *gin.Context) {
    user, exists := c.Get("current_user")
    if !exists {
        c.JSON(401, gin.H{"error": "User not found"})
        return
    }
    
    currentUser := user.(User)
    c.JSON(200, currentUser)
}
```

## 🏗️ **Middleware Best Practices**

### **1. Order Matters**
```go
router.Use(
    ErrorHandler(),      // First: Catch panics
    RequestIDMiddleware(), // Early: Generate request ID
    CORSMiddleware(),     // Early: Handle CORS
    CustomLogger(),       // Log requests
    RateLimiter(),        // Rate limit
    AuthMiddleware(),     // Authenticate (if needed)
)
```

### **2. Graceful Error Handling**
```go
func SafeMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Middleware panic: %v", r)
                c.JSON(500, gin.H{"error": "Internal server error"})
                c.Abort()
            }
        }()
        
        c.Next()
    }
}
```

### **3. Performance Considerations**
```go
// Cache expensive operations
var onceCache sync.Once
var expensiveData string

func OptimizedMiddleware() gin.HandlerFunc {
    onceCache.Do(func() {
        expensiveData = loadExpensiveData()
    })
    
    return func(c *gin.Context) {
        // Use cached data
        c.Set("data", expensiveData)
        c.Next()
    }
}
```

## 🔗 **Third-Party Middleware**

### **Popular Gin Middleware**
```go
import (
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/gzip"
)

// CORS
router.Use(cors.Default())

// Sessions
store := sessions.NewCookieStore([]byte("secret"))
router.Use(sessions.Sessions("mysession", store))

// Gzip compression
router.Use(gzip.Gzip(gzip.DefaultCompression))
```

## 🧪 **Testing Middleware**

### **Unit Testing Middleware**
```go
func TestAuthMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    router := gin.New()
    router.Use(APIKeyAuth())
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })
    
    // Test without API key
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 401, w.Code)
    
    // Test with valid API key
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/test", nil)
    req.Header.Set("X-API-Key", "valid-key")
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

## 🌍 **Real-World Applications**

### **Production Middleware Stack**
```go
func SetupMiddleware(router *gin.Engine) {
    // Security
    router.Use(SecurityHeaders())
    router.Use(RateLimiter(100, time.Minute))
    
    // Observability
    router.Use(RequestID())
    router.Use(StructuredLogger())
    router.Use(Metrics())
    
    // CORS & Content
    router.Use(CORS())
    router.Use(gzip.Gzip(gzip.DefaultCompression))
    
    // Error handling
    router.Use(ErrorHandler())
    
    // Authentication (for protected routes)
    api := router.Group("/api/v1")
    api.Use(JWTAuth())
}
```

## 📚 **Next Steps**

After mastering middleware, explore:
1. **Custom Validators**: JSON schema validation middleware
2. **Caching**: Response caching middleware
3. **Circuit Breakers**: Fault tolerance patterns
4. **Distributed Tracing**: OpenTelemetry integration
5. **Health Checks**: Endpoint monitoring middleware 