package main

import (
	"fmt"
	"log"
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
	r := gin.New()
	public := r.Group("/")
	protected := r.Group("/")

	// TODO: Setup custom middleware in correct order
	// 1. ErrorHandlerMiddleware (first to catch panics)
	r.Use(ErrorHandlerMiddleware())
	// 2. RequestIDMiddleware
	r.Use(RequestIDMiddleware())
	// 3. LoggingMiddleware
	r.Use(LoggingMiddleware())
	// 4. AuthMiddleware (for protected routes)
	protected.Use(AuthMiddleware())
	// 4. CORSMiddleware
	r.Use(CORSMiddleware())
	// 5. RateLimitMiddleware
	r.Use(RateLimitMiddleware())
	// 6. ContentTypeMiddleware
	r.Use(ContentTypeMiddleware())

	// TODO: Setup route groups
	// TODO: Define routes

	// Public routes (no authentication required)
	// Public: GET /ping, GET /articles, GET /articles/:id

	{
		public.GET("/ping", ping)
		public.GET("/articles", getArticles)
		public.GET("/articles/:id", getArticle)
	}

	// Protected routes (require authentication)
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats
	{
		protected.POST("/articles", createArticle)
		protected.PUT("/articles/:id", updateArticle)
		protected.DELETE("/articles/:id", deleteArticle)
		protected.GET("/admin/stats", getStats)
	}

	// TODO: Start server on port 8080
	r.Run(":8080")
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Generate UUID for request ID
		// Use github.com/google/uuid package
		// Store in context as "request_id"
		// Add to response header as "X-Request-ID"

		reqID := uuid.New().String()
		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)

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
		requestID := c.GetString("request_id")
		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		log.Printf("[%s] %s %s %d %v %s %s", requestID, method, path, status, duration, ip, userAgent)
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"
	apiKeysList := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		// TODO: Get API key from X-API-Key header
		// TODO: Validate API key
		// TODO: Set user role in context
		// TODO: Return 401 if invalid or missing
		api_key := c.GetHeader("X-API-Key")

		if _, exists := apiKeysList[api_key]; !exists {
			requestID := c.GetString("request_id")
			c.AbortWithStatusJSON(401, APIResponse{
				Success:   false,
				Error:     "unauthorized",
				RequestID: requestID,
				Message:   "User unauthorized: invalid or missing API key",
			})
			return
		}

		c.Set("user_role", apiKeysList[api_key])

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
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")

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

	limit := rate.Every(time.Minute / 100)
	burst := 100

	visitors := make(map[string]*rate.Limiter)
	var mu sync.Mutex

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()

		limiter, exists := visitors[ip]
		if !exists {
			limiter = rate.NewLimiter(limit, burst)
			visitors[ip] = limiter
		}
		return limiter
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.Header("X-RateLimit-Limit", "100")
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))

			c.AbortWithStatusJSON(429, APIResponse{
				Success:   false,
				Error:     "rate limit exceeded",
				RequestID: c.GetString("request_id"),
				Message:   "Too many requests. Please try again later.",
			})
			return
		}

		c.Header("X-RateLimit-Limit", "100")
		c.Header(
			"X-RateLimit-Remaining",
			strconv.Itoa(int(limiter.Tokens())),
		)
		c.Header(
			"X-RateLimit-Reset",
			strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10),
		)

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
			if c.GetHeader("Content-Type") != "application/json" {
				c.AbortWithStatusJSON(415, APIResponse{
					Success:   false,
					Error:     "unsupported media type",
					RequestID: c.GetString("request_id"),
					Message:   "Content-Type must be application/json",
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
		requestID := c.GetString("request_id")

		var message string
		if err, ok := recovered.(error); ok {
			// evite vazar detalhes em prod
			message = err.Error()
		}

		c.AbortWithStatusJSON(500, APIResponse{
			Success:   false,
			Error: "Internal server error",
			Message: message,
			RequestID: requestID,

		})
	})
}
// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	requestID := c.GetString("request_id")

	c.JSON(200, APIResponse{
		Success:   true,
		Data:      "pong",
		RequestID: requestID,
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// TODO: Implement pagination (optional)
	// TODO: Return articles in standard format
	c.JSON(200, APIResponse{
		Success:   true,
		Data:      articles,
		RequestID: c.GetString("request_id"),
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find article by ID
	// TODO: Return 404 if not found
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, APIResponse{
			Success:   false,
			Error:     "invalid Id format",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	article, index := findArticleByID(id)
	if index == -1 {
		c.JSON(404, APIResponse{
			Success:   false,
			Error:     "article not found",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success:   true,
		Data:      article,
		RequestID: c.GetString("request_id"),
	})
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// TODO: Parse JSON request body
	// TODO: Validate required fields
	// TODO: Add article to storage
	// TODO: Return created article
	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(400, APIResponse{
			Success:   false,
			Error:     "invalid request body",
			RequestID: c.GetString("request_id"),
			Message:   err.Error(),
		})
		return
	}

	if err := validateArticle(&newArticle); err != nil {
		c.JSON(400, APIResponse{
			Success:   false,
			Error:     "validation error",
			RequestID: c.GetString("request_id"),
			Message:   err.Error(),
		})
		return
	}

	articles = append(articles, newArticle)

	c.JSON(201, APIResponse{
		Success:   true,
		Data:      newArticle,
		RequestID: c.GetString("request_id"),
		Message:   "Article created successfully",
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Parse JSON request body
	// TODO: Find and update article
	// TODO: Return updated article
	type ArticleUpdate struct {
		Title   *string `json:"title,omitempty"`
		Content *string `json:"content,omitempty"`
		Author  *string `json:"author,omitempty"`
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, APIResponse{
			Success:   false,
			Error:     "Bad request",
			Message:   "invalid Id format",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	var updates ArticleUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, APIResponse{
			Success:   false,
			Error:     "invalid request body",
			RequestID: c.GetString("request_id"),
			Message:   err.Error(),
		})
		return
	}

	article, index := findArticleByID(id)
	if index == -1 {
		c.JSON(404, APIResponse{
			Success:   false,
			Message:   "article not found",
			Error:     "not found",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	if updates.Title != nil {
		article.Title = *updates.Title
	}
	if updates.Content != nil {
		article.Content = *updates.Content
	}
	if updates.Author != nil {
		article.Author = *updates.Author
	}

	article.UpdatedAt = time.Now()

	c.JSON(200, APIResponse{
		Success:   true,
		Data:      article,
		RequestID: c.GetString("request_id"),
	})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find and remove article
	// TODO: Return success message
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, APIResponse{
			Success:   false,
			Error:     "Bad request",
			Message:   "invalid Id format",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	article, index := findArticleByID(id)
	if index == -1 {
		c.JSON(404, APIResponse{
			Success:   false,
			Message:   "article not found",
			Error:     "not found",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	articles = append(articles[:index], articles[index+1:]...)

	c.JSON(200, APIResponse{
		Success:   true,
		Message:   fmt.Sprintf("Article with ID %d deleted successfully", article.ID),
		RequestID: c.GetString("request_id"),
	})
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	userRole := c.GetString("user_role")
	if userRole != "admin" {
		c.JSON(403, APIResponse{
			Success:   false,
			Error:     "forbidden",
			RequestID: c.GetString("request_id"),
			Message:   "Access denied: admin only",
		})
		return
	}
	
	// TODO: Return mock statistics
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}

	// TODO: Return stats in standard format
	c.JSON(200, APIResponse{
		Success:   true,
		Data:      stats,
		RequestID: c.GetString("request_id"),
	})
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// TODO: Implement article lookup
	// Return article pointer and index, or nil and -1 if not found
	for i := range articles {
		if articles[i].ID == id {
			return &articles[i], i
		}
	}

	return nil, -1
}

// validateArticle validates article data
func validateArticle(article *Article) error {
	// TODO: Implement validation
	// Check required fields: Title, Content, Author
	if article.Title == "" {
		return fmt.Errorf("title is required")
	}

	if article.Content == "" {
		return fmt.Errorf("content is required")
	}

	if article.Author == "" {
		return fmt.Errorf("author is required")
	}

	article.ID = nextID
	nextID++

	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	return nil
}
