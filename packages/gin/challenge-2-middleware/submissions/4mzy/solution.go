package main

import (
	"errors"
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
var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3
var articlesMu sync.Mutex

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
	router.Use(ErrorHandlerMiddleware(), RequestIDMiddleware(), LoggingMiddleware(), CORSMiddleware(), RateLimitMiddleware(), ContentTypeMiddleware())
	// TODO: Setup route groups
	// Public routes (no authentication required)
	// Protected routes (require authentication)
	router.GET("/ping", ping)
	router.GET("/articles", getArticles)
	router.GET("/articles/:id", getArticle)
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	protected.POST("/articles", createArticle)
	protected.PUT("/articles/:id", updateArticle)
	protected.DELETE("/articles/:id", deleteArticle)
	protected.GET("/admin/stats", getStats)

	// TODO: Define routes
	// Public: GET /ping, GET /articles, GET /articles/:id
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats
	router.Run(":8080")

	// TODO: Start server on port 8080
}

// TODO: Implement middleware functions

func getRequestID(c *gin.Context) string {
	value, exists := c.Get("request_id")
	if !exists {
		return ""
	}
	requestID, ok := value.(string)
	if !ok {
		return ""
	}
	return requestID
}

func handleError(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Success:   false,
		Error:     message,
		RequestID: getRequestID(c),
	})

}

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Generate UUID for request ID
		// Use github.com/google/uuid package
		// Store in context as "request_id"
		// Add to response header as "X-Request-ID"
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
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

		// TODO: Calculate duration and log request

		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		duration := time.Since(start)
		status := c.Writer.Status()
		requestID := getRequestID(c)
		log.Printf("[%s] %s %s %d %s %s %s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			status,
			duration,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"
	validKeys := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		// TODO: Get API key from X-API-Key header
		apiKey := c.GetHeader("X-API-Key")
		// TODO: Validate API key
		role, ok := validKeys[apiKey]
		if !ok {
			handleError(c, http.StatusUnauthorized, "invalid or missing API key")
			c.Abort()
			return
		}
		// TODO: Set user role in context
		c.Set("user_role", role)
		// TODO: Return 401 if invalid or missing

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := map[string]bool{
		"http://localhost:3000": true,
		"https://myblog.com":    true,
	}
	return func(c *gin.Context) {
		// TODO: Set CORS headers
		// Allow origins: http://localhost:3000, https://myblog.com
		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		// Allow headers: Content-Type, X-API-Key, X-Request-ID
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		// TODO: Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

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
	limiters := make(map[string]*rate.Limiter)
	var mu sync.Mutex
	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(time.Minute/100), 100)
			limiters[ip] = limiter
		}
		mu.Unlock()

		c.Header("X-RateLimit-Limit", "100")
		resetTime := strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10)
		c.Header("X-RateLimit-Reset", resetTime)

		if !limiter.Allow() {
			c.Header("X-RateLimit-Remaining", "0")
			handleError(c, http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}

		remaining := int(limiter.Tokens())
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check content type for POST/PUT requests
		// Must be application/json
		// Return 415 if invalid content type
		method := c.Request.Method
		if method == "POST" || method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				handleError(c, http.StatusUnsupportedMediaType, "invalid content type")
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
		// TODO: Handle panics gracefully
		// Return consistent error response format
		// Include request ID in response
		log.Printf("panic recovered: %v", recovered)
		c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Internal server error",
			RequestID: getRequestID(c),
			Message:   fmt.Sprint(recovered),
		})
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "pong",
		RequestID: getRequestID(c),
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// TODO: Implement pagination (optional)
	// TODO: Return articles in standard format
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles,
		RequestID: getRequestID(c),
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find article by ID
	// TODO: Return 404 if not found
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid id format")
		return
	}

	articlesMu.Lock()
	article, _ := findArticleByID(id)
	if article == nil {
		articlesMu.Unlock()
		handleError(c, http.StatusNotFound, "article not found")
		return
	}
	responseArticle := *article
	articlesMu.Unlock()

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      responseArticle,
		RequestID: getRequestID(c),
	})
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// TODO: Parse JSON request body
	// TODO: Validate required fields
	// TODO: Add article to storage
	// TODO: Return created article
	var article Article

	err := c.ShouldBindJSON(&article)
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if err := validateArticle(article); err != nil {
		handleError(c, http.StatusBadRequest, "validation error")
		return
	}
	articlesMu.Lock()
	article.ID = nextID
	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now
	nextID++
	articles = append(articles, article)
	articlesMu.Unlock()
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Data:      article,
		Message:   "article created",
		RequestID: getRequestID(c),
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid id format")
		return
	}
	var newArticle Article

	if err := c.ShouldBindJSON(&newArticle); err != nil {
		handleError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if err := validateArticle(newArticle); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}
	newArticle.ID = id
	articlesMu.Lock()
	article, index := findArticleByID(id)
	if article == nil {
		articlesMu.Unlock()
		handleError(c, http.StatusNotFound, "article not found")
		return
	}
	newArticle.CreatedAt = article.CreatedAt
	now := time.Now()
	newArticle.UpdatedAt = now
	articles[index] = newArticle
	articlesMu.Unlock()

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      newArticle,
		RequestID: getRequestID(c),
	})
	// TODO: Get article ID from URL parameter
	// TODO: Parse JSON request body
	// TODO: Find and update article
	// TODO: Return updated article
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid id format")
		return
	}
	articlesMu.Lock()
	_, index := findArticleByID(id)
	if index == -1 {
		articlesMu.Unlock()
		handleError(c, http.StatusNotFound, "article not found")
		return
	}
	articles = append(articles[:index], articles[index+1:]...)
	articlesMu.Unlock()

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		RequestID: getRequestID(c),
		Message:   "article deleted",
	})
	// TODO: Get article ID from URL parameter
	// TODO: Find and remove article
	// TODO: Return success message
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	// TODO: Return mock statistics
	role, ok := c.Get("user_role")
	if !ok {
		handleError(c, http.StatusUnauthorized, "authentication required")
		return
	}
	if role != "admin" {
		handleError(c, http.StatusForbidden, "admin access required")
		return
	}
	articlesMu.Lock()
	totalArticles := len(articles)
	articlesMu.Unlock()
	stats := map[string]interface{}{
		"total_articles": totalArticles,
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      stats,
		RequestID: getRequestID(c),
	})
	// TODO: Return stats in standard format
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	for i := range articles {
		if articles[i].ID == id {
			return &articles[i], i
		}
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	// TODO: Implement validation
	// Check required fields: Title, Content, Author
	if strings.TrimSpace(article.Title) == "" {
		return errors.New("'title' field is required")
	}
	if strings.TrimSpace(article.Content) == "" {
		return errors.New("'content' field is required")
	}
	if strings.TrimSpace(article.Author) == "" {
		return errors.New("'author' field is required")
	}
	return nil
}
