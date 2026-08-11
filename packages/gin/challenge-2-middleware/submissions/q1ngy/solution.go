package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// Article represents a blog article
type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"Title"`
	Content   string    `json:"Content"`
	Author    string    `json:"Author"`
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
	// Use gin.New() instead of gin.Default()
	g := gin.New()

	// Setup custom middleware in correct order
	// 1. ErrorHandlerMiddleware (first to catch panics)
	g.Use(ErrorHandlerMiddleware())
	// 2. RequestIDMiddleware
	g.Use(RequestIDMiddleware())
	// 3. LoggingMiddleware
	g.Use(LoggingMiddleware())
	// 4. CORSMiddleware
	g.Use(CORSMiddleware())
	// 5. RateLimitMiddleware
	g.Use(RateLimitMiddleware())
	// 6. ContentTypeMiddleware
	g.Use(ContentTypeMiddleware())

	// Public routes (no authentication required)
	public := g.Group("/")

	// Protected routes (require authentication)
	private := g.Group("/", AuthMiddleware())

	// Define routes
	// Public: GET /ping, GET /articles, GET /articles/:id
	public.GET("/ping", ping)
	public.GET("/articles", getArticles)
	public.GET("/articles/:id", getArticle)

	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats
	articlesGroup := private.Group("articles")
	articlesGroup.POST("/", createArticle)
	articlesGroup.PUT("/:id", updateArticle)
	articlesGroup.DELETE("/:id", deleteArticle)

	private.GET("/admin/stats", getStats)

	g.Run(":8080")
}

// Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use github.com/google/uuid package
		// Store in context as "request_id"
		// Add to response header as "X-Request-ID"
		requestId := c.GetHeader("X-Request-ID")
		if requestId == "" {
			requestId = uuid.New().String()
		}

		c.Set("request_id", requestId)
		c.Writer.Header().Set("X-Request-ID", requestId)
		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture start time
		start := time.Now()

		c.Next()

		// Calculate duration and log request
		d := time.Since(start)

		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		requestId, _ := c.Get("request_id")
		logLine := fmt.Sprintf("[%s] %s %s %d %s %s %s\n",
			requestId,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			d,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
		fmt.Print(logLine)
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"
	m := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		// Get API key from X-API-Key header
		key := c.GetHeader("X-API-Key")

		// Validate API key
		// Set user role in context
		// Return 401 if invalid or missing
		role, ok := m[key]
		if !ok {
			c.JSON(http.StatusUnauthorized, APIResponse{Success: false})
		}

		c.Set("role", role)

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow origins: http://localhost:3000, https://myblog.com
		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		// Allow headers: Content-Type, X-API-Key, X-Request-ID
		origin := c.GetHeader("Origin")
		if origin == "http://localhost:3000" || origin == "https://myblog.com" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
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
	var visitors sync.Map
	return func(c *gin.Context) {
		ip := c.ClientIP()

		v, _ := visitors.LoadOrStore(ip, rate.NewLimiter(rate.Every(time.Minute/100), 100))
		limiter := v.(*rate.Limiter)

		c.Writer.Header().Set("X-RateLimit-Limit", "100")
		c.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

		if !limiter.Allow() {
			c.Writer.Header().Set("X-RateLimit-Remaining", "0")
			c.AbortWithStatus(http.StatusTooManyRequests)
		}

		remaining := int(limiter.Tokens())
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		c.Next()
	}
}

// ContentTypeMiddleware validates Content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Content type for POST/PUT requests
		// Must be application/json
		// Return 415 if invalid Content type
		method := c.Request.Method
		if method == "POST" || method == "Put" {
			contentType := c.ContentType()
			if contentType != "application/json" {
				c.JSON(http.StatusUnsupportedMediaType, APIResponse{Success: false})
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
		requestId, _ := c.Get("request_id")
		var errMsg string
		switch v := recovered.(type) {
		case string:
			c.JSON(http.StatusBadRequest, APIResponse{
				Success:   false,
				Data:      nil,
				Message:   errMsg,
				Error:     errMsg,
				RequestID: requestId.(string),
			})
		case error:
			errMsg = v.Error()
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Data:    nil,
				Message: errMsg,
				//Error:     http.StatusText(http.StatusInternalServerError),
				Error:     "Internal server error",
				RequestID: requestId.(string),
			})
		default:
		}

		c.Abort()
	})
}

// Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// Return simple pong response with request ID
	requestId := c.Writer.Header().Get("X-Request-ID")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		RequestID: requestId,
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	requestId, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles,
		RequestID: requestId.(string),
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// Get article ID from URL parameter
	// Find article by ID
	// Return 404 if not found
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic("illegal id")
	}
	article, errCode := findArticleByID(id)
	if errCode != 0 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
		})
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    article,
	})
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// Parse JSON request body
	// Validate required fields
	// Add article to storage
	// Return created article
	type Req struct {
		Title   string `binding:"required"`
		Content string `binding:"required"`
		Author  string `binding:"required"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
	}
	article := Article{
		ID:        3,
		Title:     req.Title,
		Content:   req.Content,
		Author:    req.Author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	articles = append(articles, article)
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    article,
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// Get article ID from URL parameter
	// Parse JSON request body
	// Find and update article
	// Return updated article
	type Req struct {
		Title   string `binding:"required"`
		Content string `binding:"required"`
		Author  string `binding:"required"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
	}
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var article *Article
	for i, _ := range articles {
		pointerObj := &articles[i]
		if id == pointerObj.ID {
			article = pointerObj
		}
	}
	if article == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
		})
		return
	}
	article.Title = req.Title
	article.Author = req.Author
	article.Content = req.Content
	article.UpdatedAt = time.Now()
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    article,
	})

}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// Get article ID from URL parameter
	// Find and remove article
	// Return success message
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	filtered := articles[:0]
	for _, v := range articles {
		if id != v.ID {
			filtered = append(filtered, v)
		}
	}
	articles = filtered
	c.JSON(http.StatusOK, APIResponse{Success: true})

}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// Check if user role is "admin"
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, APIResponse{Success: false})
		return
	}

	// Return mock statistics
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}
	fmt.Println(stats)

	// Return stats in standard format
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// Implement article lookup
	// Return article pointer and index, or nil and -1 if not found
	var article *Article
	for _, v := range articles {
		if id == v.ID {
			article = &v
			return article, 0
		}
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	// TODO: Implement validation
	// Check required fields: Title, Content, Author
	return nil
}
