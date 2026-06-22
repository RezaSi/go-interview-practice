package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
var mu sync.RWMutex

func main() {
	router := gin.New()

	router.Use(ErrorHandlerMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(LoggingMiddleware())
	router.Use(CORSMiddleware())
	router.Use(RateLimitMiddleware())
	router.Use(ContentTypeMiddleware())

	publicRoute := router.Group("/")
	{
		publicRoute.GET("/ping", ping)
		publicRoute.GET("/articles", getArticles)
		publicRoute.GET("/articles/:id", getArticle)
	}

	protectedRoute := router.Group("/", LoggingMiddleware(), AuthMiddleware())
	{
		protectedRoute.POST("/articles", createArticle)
		protectedRoute.PUT("/articles/:id", updateArticle)
		protectedRoute.DELETE("articles/:id", deleteArticle)
		protectedRoute.GET("/admin/stats", getStats)
	}

	router.Run(":8080")
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// requestID := c.GetHeader("X-Request-ID")
		requestID := c.GetString("RequestID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		entry := map[string]interface{}{
			"requestID":  c.GetString("RequestID"),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		if c.Writer.Status() >= 400 {
			log.Printf("ERROR %+v", entry)
		} else {
			log.Printf("WARNING %+v", entry)
		}
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"
	const (
		admin = "admin-key-123"
		user  = "user-key-456"
	)

	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == admin {
			c.Set("role", "admin")
		}
		if key == user {
			c.Set("role", "user")
		}
		if key == "" {
			errorResponse(c, "invalid API key", http.StatusUnauthorized)
			c.Abort()
		}

		c.Next()
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

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	// TODO: Implement rate limiting
	// Limit: 100 requests per IP per minute
	// Use golang.org/x/time/rate package
	// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
	// Return 429 if rate limit exceeded

	return func(c *gin.Context) {
		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check content type for POST/PUT requests
		// Must be application/json
		// Return 415 if invalid content type

		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// TODO: Handle panics gracefully
		// Return consistent error response format
		// Include request ID in response
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	successResponse(c, "pong", "success", http.StatusOK)
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// TODO: Implement pagination (optional)
	// TODO: Return articles in standard format
	mu.RLock()
	art := articles
	mu.RUnlock()
	successResponse(c, art, "success", http.StatusOK)
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find article by ID
	// TODO: Return 404 if not found
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, "bad request", http.StatusBadRequest)
	}
	mu.RLock()
	articlesCp := articles
	mu.RUnlock()
	for _, article := range articlesCp {
		if article.ID == id {
			successResponse(c, article, "success", http.StatusOK)
			return
		}
	}
	errorResponse(c, "not found", http.StatusNotFound)
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// TODO: Parse JSON request body
	// TODO: Validate required fields
	// TODO: Add article to storage
	// TODO: Return created article
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Parse JSON request body
	// TODO: Find and update article
	// TODO: Return updated article
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find and remove article
	// TODO: Return success message
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	// TODO: Return mock statistics

	// stats := map[string]interface{}{
	// 	"total_articles": len(articles),
	// 	"total_requests": 0, // Could track this in middleware
	// 	"uptime":         "24h",
	// }

	// TODO: Return stats in standard format
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	for i, article := range articles {
		if article.ID == id {
			return &articles[i], i
		}
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	if article.Title == "" {
		return errors.New("title field required")
	}
	if article.Content == "" {
		return errors.New("content filed required")
	}
	if article.Author == "" {
		return errors.New("author filed required")
	}
	return nil
}

func errorResponse(c *gin.Context, err string, status int) {
	response := APIResponse{
		Success:   false,
		Error:     err,
		RequestID: c.GetString("RequestID"),
	}

	responseHandler(c, response, status)
}

func successResponse(c *gin.Context, data interface{}, msg string, status int) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Message:   msg,
		RequestID: c.GetString("RequestID"),
	}

	responseHandler(c, response, status)
}

func responseHandler(c *gin.Context, response interface{}, status int) {
	c.JSON(status, response)
}
