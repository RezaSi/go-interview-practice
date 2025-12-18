package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
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
var (
	articles = []Article{
		{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	nextID = 3
	mu     sync.RWMutex
)

func main() {
	// Create Gin router without default middleware
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
	)

	// Setup route groups
	publicEndpoint := router.Group("/")
	protectedEndpoint := router.Group("/")

	// Public routes (no authentication required)
	// Protected routes (require authentication)

	// Define routes
	// Public: GET /ping, GET /articles, GET /articles/:id
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats
	//
	publicEndpoint.Use()
	{
		publicEndpoint.GET("ping", ping)
		articles := publicEndpoint.Group("articles")
		articles.Use()
		{
			articles.GET("/", getArticles)
			articles.GET("/:id", getArticle)
		}
	}

	protectedEndpoint.Use(AuthMiddleware())
	{
		protectedEndpoint.GET("admin/stats", getStats)

		articles := protectedEndpoint.Group("articles")
		articles.Use(ContentTypeMiddleware())
		{
			articles.POST("/", createArticle)
			articles.PUT("/:id", updateArticle)
			articles.DELETE("/:id", deleteArticle)
		}
	}

	// Start server on port 8080
	router.Run(":8081")
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate UUID for request ID
		// Use github.com/google/uuid package
		// Store in context as "request_id"
		// Add to response header as "X-Request-ID"

		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture start time

		start := time.Now()

		c.Next()

		duration := time.Since(start)

		// Calculate duration and log request
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		log.Printf("[%s] %s %s %d %s %s %s",
			c.GetString("request_id"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// "admin-key-123" -> "admin"
	// Define valid API keys and their roles
	// "user-key-456" -> "user"
	//

	admin := "admin-key-123"
	user := "user-key-456"

	return func(c *gin.Context) {
		// Get API key from X-API-Key header

		apiKey := c.GetHeader("X-API-Key")

		// Validate API key

		if apiKey == admin {
			c.Set("Auth", "admin")
		} else if apiKey == user {
			// Set user role in context
			c.Set("Auth", "user")
		} else {
			// Return 401 if invalid or missing
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid Key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		// Allow origins: http://localhost:3000, https://myblog.com
		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		// Allow headers: Content-Type, X-API-Key, X-Request-ID

		origin := c.Request.Header.Get("Origin")

		AllowOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://myblog.com":     true,
		}

		if AllowOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS requests

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	// Implement rate limiting
	// Limit: 100 requests per IP per minute
	// Use golang.org/x/time/rate package
	// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
	// Return 429 if rate limit exceeded

	const (
		burstLimit   = 100
		fillInterval = time.Minute / 100
	)

	var rateLimiters = make(map[string]*rate.Limiter)
	var mu sync.Mutex

	return func(c *gin.Context) {

		clientIP := c.ClientIP()

		mu.Lock()
		limiter, exists := rateLimiters[clientIP]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(fillInterval), burstLimit)
			rateLimiters[clientIP] = limiter
		}
		mu.Unlock()

		allowed := limiter.Allow()
		c.Header("X-RateLimit-Limit", strconv.Itoa(100))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(int(limiter.Tokens())))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(fillInterval).Unix(), 10))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, APIResponse{
				Success: false,
				Error:   "Too Many Requests",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT reques:wq
// ts
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check content type for POST/PUT requests

		// Must be application/json
		// Return 415 if invalid content type
		//
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			contentType := c.GetHeader("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, APIResponse{
					Success: false,
					Error:   "Unsupported MediaType",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Handle panics gracefully
		// Return consistent error response format
		// Include request ID in response

		log.Printf("Panic: %v", recovered)

		requestID := c.GetString("request_id")

		if requestID == "" {

			c.JSON(http.StatusBadGateway, APIResponse{
				Success: false,
				Error:   "RequestID Not Exists",
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Message:   fmt.Sprintf("%v", recovered),
			Error:     "Internal server error",
			RequestID: requestID,
		})
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "Pong",
		RequestID: c.GetString("request_id"),
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// TODO: Implement pagination (optional)

	mu.RLock()
	defer mu.RUnlock()
	// Return articles in standard format
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles,
		RequestID: c.GetString("request_id"),
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// Get article ID from URL parameter
	articlesID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	mu.RLock()
	defer mu.RUnlock()
	// Find article by ID
	article, index := findArticleByID(articlesID)

	// Return 404 if not found
	if index == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Article Not found",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      article,
		RequestID: c.GetString("request_id"),
	})

}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// Parse JSON request body
	var articleContent Article
	err := c.ShouldBindJSON(&articleContent)

	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}
	// Validate required fields
	if err := validateArticle(articleContent); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Add article to storage
	if err := validateArticle(articleContent); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	articleContent.ID = nextID
	articleContent.CreatedAt = time.Now()
	articleContent.UpdatedAt = time.Now()

	mu.Lock()
	defer mu.Unlock()

	articles = append(articles, articleContent)
	nextID++

	// Return created article
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Data:      articleContent,
		RequestID: c.GetString("request_id"),
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// Get article ID from URL parameter
	articlesID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}
	// Parse JSON request body

	var articleUpdate Article
	err = c.ShouldBindJSON(&articleUpdate)

	if err := validateArticle(articleUpdate); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Find and update article
	_, index := findArticleByID(articlesID)

	if index == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Article Not found",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Return updated article

	articles[index].Title = articleUpdate.Title
	articles[index].Content = articleUpdate.Content
	articles[index].Author = articleUpdate.Author
	articles[index].UpdatedAt = time.Now()

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articleUpdate,
		RequestID: c.GetString("request_id"),
	})

}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	articlesID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// TODO: Find and remove article

	mu.Lock()
	defer mu.Unlock()

	_, index := findArticleByID(articlesID)

	if index == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Article Not found",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	articles = append(articles[:index], articles[index+1:]...)

	// TODO: Return success message
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles,
		RequestID: c.GetString("request_id"),
	})

}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	userRole := c.GetString("Auth")

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, APIResponse{
			Success: false,
			Error:   "Permission Deny",
		})
		return
	}

	// TODO: Return mock statistics
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}

	// Return stats in standard format

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      stats,
		RequestID: c.GetString("request_id"),
	})

}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// Implement article lookup
	index := -1

	if id < nextID {

		for i, art := range articles {
			if art.ID == id {
				index = i
				break
			}
		}
	}
	// Return article pointer and index, or nil and -1 if not found
	if index != -1 {
		return &articles[index], index
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	// Implement validation
	// Check required fields: Title, Content, Author
	if article.Title == "" || article.Content == "" || article.Author == "" {
		return fmt.Errorf("Invalid Article")
	}
	return nil
}