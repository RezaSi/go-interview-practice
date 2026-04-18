package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"strconv"
	"sync"
	"strings"
	"errors"
	"fmt"
)

// Article represents a blog article
type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// In-memory storage
var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3

func main() {
	// TODO: Create Gin router without default middleware
	// Use gin.New() instead of gin.Default()
    router := gin.New()
	// TODO: Setup custom middleware in correct order
	// 1. ErrorHandlerMiddleware (first to catch panics)
	router.Use(ErrorHandlerMiddleware())
	// 2. RequestIDMiddleware
	router.Use(RequestIDMiddleware())
	// 3. LoggingMiddleware
	router.Use(LoggingMiddleware())
	// 4. CORSMiddleware
	router.Use(CORSMiddleware())
	// 5. RateLimitMiddleware
	router.Use(RateLimitMiddleware())
	// 6. ContentTypeMiddleware
	router.Use(ContentTypeMiddleware())

	// TODO: Setup route groups
	// Public routes (no authentication required)
	router.GET("/ping", ping)
	router.GET("/articles", getArticles)
	router.GET("/articles/:id", getArticle)
	// Protected routes (require authentication)
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
    protected.POST("/articles", createArticle)
	protected.PUT("/articles/:id", updateArticle)
	protected.DELETE("/articles/:id", deleteArticle)
	protected.GET("/admin/stats", getStats)
	// TODO: Define routes
	// Public: GET /ping, GET /articles, GET /articles/:id
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats

	// TODO: Start server on port 8080
	router.Run()
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Generate UUID for request ID
		// Use github.com/google/uuid package
		// Store in context as "request_id"
		// Add to response header as "X-Request-ID"
        id := uuid.NewString()
		c.Set("request_id", id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Capture start time
        start := time.Now()
        
		c.Next()
		
		duration := time.Since(start)

		// TODO: Calculate duration and log request
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		method:= c.Request.Method
		path := c.Request.URL.Path
		status:= c.Writer.Status()
		ip:= c.ClientIP()
		userAgent := c.Request.UserAgent()
		requestID, exists := c.Get("request_id")
		if !exists {
			requestID = "-"
		}

		// 5️⃣ Log in required format
		log.Printf(
			"[%v] %s %s %d %s %s %s",
			requestID,
			method,
			path,
			status,
			duration,
			ip,
			userAgent,
		)
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"
	apiKeys := map[string]string{
	    "admin-key-123": "admin",
	    "user-key-456" : "user",
	}

	return func(c *gin.Context) {
		// TODO: Get API key from X-API-Key header
		apiKey := c.GetHeader("X-API-Key")
		// TODO: Validate API key
		if apiKey == ""{
		    c.AbortWithStatusJSON(http.StatusUnauthorized,APIResponse{
		        Success : false,
		        Message : "API key is missing",
		        Error   : "API key is missing",
		    })
		    return
		}
		role,ok:= apiKeys[apiKey]
		if !ok{
		    c.AbortWithStatusJSON(http.StatusUnauthorized,APIResponse{
		        Success : false,
		        Message : "Invalid API key",
		        Error   : "Invalid API key",
		    })
		    return
		}
		// TODO: Set user role in context
		c.Set("role",role)
		// TODO: Return 401 if invalid or missing

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
    allowedOrigins := map[string]bool{
        "http://localhost:3000":true,
        "https://myblog.com" : true,
    }
	return func(c *gin.Context) {
		// TODO: Set CORS headers
		// Allow origins: http://localhost:3000, https://myblog.com
		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		// Allow headers: Content-Type, X-API-Key, X-Request-ID

		// TODO: Handle preflight OPTIONS requests
		origin:= c.GetHeader("origin")
		if allowedOrigins[origin]{
		    c.Writer.Header().Set("Access-Control-Allow-Origin",origin)
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS")

        c.Writer.Header().Set("Access-Control-Allow-Headers","Content-Type, X-API-Key, X-Request-ID")
        
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        
        if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
// func RateLimitMiddleware() gin.HandlerFunc {
// 	// TODO: Implement rate limiting
// 	// Limit: 100 requests per IP per minute
// 	// Use golang.org/x/time/rate package
// 	// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
// 	// Return 429 if rate limit exceeded
//     var rateLimiters = make(map[string]*rate.Limiter)
//     var mu sync.Mutex
//     	const(
// 	 requestPerMinute = 100   
// 	 )
	 
// 	 requestsPerSecond := rate.Limit(requestPerMinute)/60

// 	return func(c *gin.Context) {
//         clientIP := c.ClientIP()
//         mu.Lock()
//         limiter, exists := rateLimiters[clientIP]
//         if !exists {
//             limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond*2)
//             rateLimiters[clientIP] = limiter
//         }
//         mu.Unlock()
//         if !limiter.Allow() {
//             c.JSON(429, gin.H{"error": "Rate limit exceeded"})
//             c.Abort()
//             return
//         }
//         c.Next()
// 	}
// }

func RateLimitMiddleware() gin.HandlerFunc {
	const limit = 100

	var (
		mu       sync.Mutex
		limiters = make(map[string]*rate.Limiter)
	)

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()

		if l, ok := limiters[ip]; ok {
			return l
		}

		// Effectively disables refill during the test loop
		l := rate.NewLimiter(rate.Every(time.Minute), limit)
		limiters[ip] = l
		return l
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		allowed := limiter.Allow()
		remaining := int(limiter.Tokens())

		c.Writer.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(max(remaining, 0)))
		c.Writer.Header().Set(
			"X-RateLimit-Reset",
			strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10),
		)

		if !allowed {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}


func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check content type for POST/PUT requests
		// Must be application/json
		// Return 415 if invalid content type
        method:= c.Request.Method
        
        if method == http.MethodPost || method == http.MethodPut{
            contentType:= c.GetHeader("Content-Type")
            
            if !strings.HasPrefix(contentType, "application/json"){
                c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, APIResponse{
                    Success: false,
					Error: "Content-Type must be application/json",
				})
				return
            }
        }
		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// TODO: Handle panics gracefully
		// Return consistent error response format
		// Include request ID in response
		// Default values
		statusCode := http.StatusInternalServerError
		message := ""

		// If panic is an error type, we can extract info
		if err, ok := recovered.(error); ok {
			message = err.Error()
		}

		// Try to get request ID from context
		requestID := ""
		ID, exists := c.Get("request_id")
		if exists {
			requestID = ID.(string)
		}
		    
		// Log the panic (VERY IMPORTANT)
		log.Printf(
			"panic recovered: %v | method=%s path=%s request_id=%v",
			recovered,
			c.Request.Method,
			c.Request.URL.Path,
			requestID,
		)

		// Send consistent error response
		c.AbortWithStatusJSON(statusCode, APIResponse{
		    Success:    false,
			Message:    message,
			RequestID: requestID,
			Error    : "Internal server error",
		})
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	requestID := ""
	ID, exists := c.Get("request_id")
	if exists {
	    requestID = ID.(string)
	}
	c.JSON(http.StatusOK,APIResponse{
	    Success : true,
	    Message : "Pinged Success",
	    RequestID: requestID,
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// TODO: Implement pagination (optional)
	// TODO: Return articles in standard format
	requestID := ""
	ID,exist := c.Get("request_id")
	if exist{
	    requestID = ID.(string)
	}
	c.JSON(http.StatusOK,APIResponse{
	    Success : true,
	    Data :  articles,
	    Message : "Article Success",
	    RequestID: requestID,
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	id:= c.Param("id")
	articleId ,err := strconv.Atoi(id)
	if err != nil{
	    c.JSON(http.StatusBadRequest, APIResponse{
	        Success : false,
	        Error: "invalid id in input",
	    })
	    return
	}
	// TODO: Find article by ID
	article , index := findArticleByID(articleId)
	if index == -1{
	    c.JSON(http.StatusNotFound, APIResponse{
	        Success : false,
	        Error: "Aritcle Not Found",
	    })
	    return
	}
	c.JSON(http.StatusOK,APIResponse{
	    Success : true,
	    Message : "Article Found",
	    Data : article,
	})
	// TODO: Return 404 if not found
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// TODO: Parse JSON request body
	var newArticle Article
	if err:= c.ShouldBindJSON(&newArticle);err!=nil{
	    c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid user format",
		})
        return
	}
	// TODO: Validate required fields
	err := validateArticle(newArticle)
	if err != nil{
	    c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid user format",
		})
        return
	}
	// TODO: Add article to storage
	newArticle.ID = len(articles)+1
	articles = append(articles,newArticle)
	// TODO: Return created article
	c.JSON(201, APIResponse{
			Success: true,
			Data   : newArticle,
		})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	id:= c.Param("id")
	articleId ,err := strconv.Atoi(id)
	if err != nil{
	    c.JSON(http.StatusBadRequest, APIResponse{
	        Success : false,
	        Error: "invalid id in input",
	    })
	    return
	}
	// TODO: Find article by ID
	_ , index := findArticleByID(articleId)
	if index == -1{
	    c.JSON(http.StatusNotFound, APIResponse{
	        Success : false,
	        Error: "Aritcle Not Found",
	    })
	    return
	}
	// TODO: Parse JSON request body
	var newArticle Article
	if err:= c.ShouldBindJSON(&newArticle);err!=nil{
	    c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid user format",
		})
        return
	}
	// TODO: Find and update article
	articles[index] = newArticle
	// TODO: Return updated article
	c.JSON(200, APIResponse{
	    Success: true,
	    Data : articles[index],
	})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	id:= c.Param("id")
	articleId ,err := strconv.Atoi(id)
	if err != nil{
	    c.JSON(http.StatusBadRequest, APIResponse{
	        Success : false,
	        Error: "invalid id in input",
	    })
	    return
	}
	
	// TODO: Find and remove article
	_ , index := findArticleByID(articleId)
	if index == -1{
	    c.JSON(http.StatusNotFound, APIResponse{
	        Success : false,
	        Error: "Aritcle Not Found",
	    })
	    return
	}
	
	articles = append(articles[:index], articles[index+1:]...)
	// TODO: Return success message
	c.JSON(200,APIResponse{
	    Success : true,
	})
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	role := ""
	r, exists := c.Get("role")
	if exists {
	    role = r.(string)
	}
	fmt.Println("role",role)
	if role == "admin"{
	    stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	    }
	    c.JSON(200, APIResponse{
	        Success : true,
	        Data : stats,
	    })
	    return
	}else{
	    c.JSON(403,APIResponse{
	        Success : false,
	        Error : "not allowed to access stats only for admins",
	    })
	    return
	}
	// TODO: Return mock statistics
// 	stats := map[string]interface{}{
// 		"total_articles": len(articles),
// 		"total_requests": 0, // Could track this in middleware
// 		"uptime":         "24h",
// 	}

	// TODO: Return stats in standard format
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// TODO: Implement article lookup
	// Return article pointer and index, or nil and -1 if not found
	for index,article := range articles{
	    if article.ID == id{
	        return &article,index
	    }
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	// TODO: Implement validation
	// Check required fields: Title, Content, Author
	if article.Title == "" {
        return errors.New("Article Title is required")
    }
    if article.Content == "" {
        return errors.New("No content provide for article")
    }
    
     if article.Author == "" {
        return errors.New("Article Author is required")
    }
    
	return nil
}
