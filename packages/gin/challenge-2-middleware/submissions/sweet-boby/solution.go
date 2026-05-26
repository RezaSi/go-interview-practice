package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"sync"
	"github.com/google/uuid"
	"strings"
	"strconv"
	"errors"
	"fmt"
    "golang.org/x/time/rate"
    "log"
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
	// 2. RequestIDMiddleware
	// 3. LoggingMiddleware
	// 4. CORSMiddleware
	// 5. RateLimitMiddleware
	// 6. ContentTypeMiddleware
    router.Use(
        ErrorHandlerMiddleware(),
        RequestIDMiddleware(),
        LoggingMiddleware(),
        CORSMiddleware(),
        RateLimitMiddleware(),
        ContentTypeMiddleware(),
        )

	// TODO: Setup route groups
	// Public routes (no authentication required)
	// Protected routes (require authentication)


	// TODO: Define routes
	// Public: GET /ping, GET /articles, GET /articles/:id
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats

    router.GET("/ping",ping)
    router.GET("/articles",getArticles)    
    router.GET("/articles/:id",getArticle)

    api := router.Group("/api")
    api.Use(AuthMiddleware())
        api.POST("/articles", createArticle)
        api.PUT("/articles/:id", updateArticle)
        api.DELETE("/articles/:id", deleteArticle)
        api.GET("/admin/stats", getStats)
    
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

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Capture start time

        start := time.Now()
        c.Next()
        duration := time.Since(start)
        log.Printf("[%s] %s %s - %v", 
            c.Request.Method, 
            c.Request.URL.Path, 
            c.ClientIP(), 
            duration)
		// TODO: Calculate duration and log request
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"
    keys := map[string]string{
        "admin-key-123":"admin",
        "user-key-456" : "user",
    }


	return func(c *gin.Context) {
		// TODO: Get API key from X-API-Key header
		// TODO: Validate API key
		// TODO: Set user role in context
		// TODO: Return 401 if invalid or missing
        apiKey := c.GetHeader("X-API-Key")
        
        if apiKey == "" {
            c.JSON(401, gin.H{"error": "API key required"})
            c.Abort() // Stop middleware chain
            return
        }
        key,ok := keys[apiKey];
        // Validate API key
        if ok == false {
            c.JSON(401, gin.H{"error": "Invalid API key"})
            c.Abort()
            return
        }
        
        // Store user info in context
        c.Set("user_role", key)
        
        c.Next() // Continue to next middleware/handler
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Set CORS headers
		// Allow origins: http://localhost:3000, https://myblog.com
		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		// Allow headers: Content-Type, X-API-Key, X-Request-ID

		// TODO: Handle preflight OPTIONS requests
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

var rateLimiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	// TODO: Implement rate limiting
	// Limit: 100 requests per IP per minute
	// Use golang.org/x/time/rate package
	// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
	// Return 429 if rate limit exceeded
	requestsPerSecond := 100
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		mu.Lock()
		limiter, exists := rateLimiters[clientIP]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond*2)
			rateLimiters[clientIP] = limiter
		}
		mu.Unlock()

        remaining := int(limiter.Tokens()) 
		reset := time.Now().Add(time.Minute).Unix()

		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerSecond))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(reset, 10))
		if !limiter.Allow() {

			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check content type for POST/PUT requests
		// Must be application/json
		// Return 415 if invalid content type
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

type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// TODO: Handle panics gracefully
		// Return consistent error response format
		// Include request ID in response
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
			"message":apiErr.Details,
			"request_id": c.GetString("request_id"),
		})
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	data,_ := c.Get("request_id")
	c.JSON(200,gin.H{
	    "success":true,
	    "request_id":data,
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// TODO: Implement pagination (optional)
	// TODO: Return articles in standard format
	data, _ := c.Get("request_id")
	c.JSON(200, gin.H{
	    	    "success":true,
		"data":       articles,
		"request_id": data,
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find article by ID
	// TODO: Return 404 if not found
	id := c.Param("id")
	artID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}

	for _, art := range articles {
		if art.ID == artID {
			c.JSON(200, gin.H{
				"success": true,
				"data":    art,
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"success": false,
	})
	
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// TODO: Parse JSON request body
	// TODO: Validate required fields
	// TODO: Add article to storage
	// TODO: Return created article
	var article Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}

	if article.Title == "" {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}

	if article.Content == "" {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}

	if article.Author == "" {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	articles = append(articles, article)
	c.JSON(201,gin.H{
	    "success": true,
	    "data":article,
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Parse JSON request body
	// TODO: Find and update article
	// TODO: Return updated article
	id := c.Param("id")
	artID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
		})
		return
	}

	var article Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}
	for index, art := range articles {
		if art.ID == artID {
			oldArt := &articles[index]
			oldArt.Title = art.Title
			oldArt.Content = art.Content
			oldArt.UpdatedAt = time.Now()
			c.JSON(200, gin.H{
				"success": true,
				"data":    article,
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"success": false,
	})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find and remove article
	// TODO: Return success message
		id := c.Param("id")
	artID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
		})
		return
	}

	for index, art := range articles {
		if art.ID == artID {
			articles = append(articles[:index], articles[index+1:]...)
			c.JSON(200, gin.H{
				"success": true,
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"success": false,
	})

}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	// TODO: Return mock statistics
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}
	role, ok := c.Get("user_role")

	if ok == false {
		c.JSON(403, gin.H{
			"success": false,
		})
		return
	}

	if role != "admin" {
		c.JSON(403, gin.H{
			"success": false,
		})
		return
	}
	// TODO: Return stats in standard format
	c.JSON(200, gin.H{
		"success": true,
		"data":    stats,
	})

	// TODO: Return stats in standard format
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// TODO: Implement article lookup
	// Return article pointer and index, or nil and -1 if not found
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	// TODO: Implement validation
	// Check required fields: Title, Content, Author
	if article.Title == "" {
		return errors.New("title is nil")
	}
	if article.Content == "" {
		return errors.New("content is nil")
	}
	if article.Author == "" {
		return errors.New("author is nil")
	}
	return nil
}
