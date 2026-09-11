package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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
var mu sync.RWMutex
var requestCounter atomic.Int64
var startTime = time.Now()

func Middleware() {
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
	api := router.Group("/")
	api.Use(AuthMiddleware())
	{
		api.POST("/articles", createArticle)
		api.PUT("/articles/:id", updateArticle)
		api.DELETE("/articles/:id", deleteArticle)
		api.GET("/admin/stats", RequireRole("admin"), getStats)
	}

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
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		// Use github.com/google/uuid package
		// Store in context as "request_id"
		c.Set("request_id", requestID)
		// Add to response header as "X-Request-ID"
		c.Header("X-Request-ID", requestID)
		requestCounter.Add(1)
		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Capture start time
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		// TODO: Calculate duration and log request
		duration := time.Since(start)
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
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

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"
	userRole := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		// TODO: Get API key from X-API-Key header
		apiKey := c.GetHeader("X-API-Key")
		// TODO: Return 401 if invalid or missing
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Unauthorized",
				Message: "API key required",
			})
			c.Abort()
			return
		}
		// TODO: Validate API key
		role, ok := userRole[apiKey]
		if !ok {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Unauthorized",
				Message: "Invalid API key",
			})
			c.Abort()
			return
		}
		// TODO: Set user role in context
		c.Set("user_role", role)
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
		origin := c.Request.Header.Get("Origin")

		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow headers: Content-Type, X-API-Key, X-Request-ID
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")
		// TODO: Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	}
}

type IPClient struct {
	limiter  *rate.Limiter
	lastSeen atomic.Int64
}

var clients sync.Map
var once sync.Once
// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	// TODO: Implement rate limiting
	// Limit: 100 requests per IP per minute
	// Use golang.org/x/time/rate package
	// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
	// Return 429 if rate limit exceeded
	once.Do(func() {
		go func() {
			for {
				time.Sleep(1 * time.Minute)
				clients.Range(func(key, value any) bool {
					client := value.(*IPClient)
					if time.Since(time.Unix(0, client.lastSeen.Load())) > 3*time.Minute {
						clients.Delete(key)
					}
					return true
				})
			}
		}()
	})

	return func(c *gin.Context) {
		ip := c.ClientIP()

		var client *IPClient
		if val, ok := clients.Load(ip); ok {
			client = val.(*IPClient)
		} else {
			newClient := &IPClient{
				limiter: rate.NewLimiter(rate.Every(1*time.Minute/100), 100),
			}
			actual, _ := clients.LoadOrStore(ip, newClient)
			client = actual.(*IPClient)
		}
		client.lastSeen.Store(time.Now().UnixNano())
		maxLimit := client.limiter.Burst()
		tokensBeforeRequest := int(client.limiter.Tokens())
		tokensMissing := float64(maxLimit - tokensBeforeRequest)
		refillRatePerSecond := float64(client.limiter.Limit())

		var resetTime int64
		if refillRatePerSecond > 0 && tokensMissing > 0 {
			secondToFullyRefill := math.Ceil(tokensMissing / refillRatePerSecond)
			resetTime = time.Now().Add(time.Duration(secondToFullyRefill) * time.Second).Unix()
		} else {
			resetTime = time.Now().Add(1 * time.Minute).Unix()
		}
		c.Header("X-RateLimit-Limit", strconv.Itoa(maxLimit))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
		if !client.limiter.Allow() {
			c.Header("X-RateLimit-Remaining", "0")
			c.JSON(http.StatusTooManyRequests, APIResponse{
				Success: false,
				Error:   "Too Many Requests",
				Message: "rate limit exceeded",
			})
			c.Abort()
			return
		}
		remainingAfterConsumption := int(client.limiter.Tokens())
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remainingAfterConsumption))

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check content type for POST/PUT requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")

			// Must be application/json
			// Return 415 if invalid content type
			if !strings.HasPrefix(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, APIResponse{
					Success: false,
					Error:   "Unsupported Media Type",
					Message: "Content-Type must be application/json",
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
		// TODO: Handle panics gracefully
		// Return consistent error response format
		switch err := recovered.(type) {
		case error:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success:   false,
				Error:     "Internal server error",
				Message:   err.Error(),
				RequestID: c.GetString("request_id"),
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success:   false,
				Error:     "Internal server error",
				Message:   fmt.Sprintf("%v", recovered),
				RequestID: c.GetString("request_id"),
			})
		}
		// Include request ID in response
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		RequestID: c.GetString("request_id"),
		Message:   "pong",
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// TODO: Implement pagination (optional)
	// TODO: Return articles in standard format
	mu.RLock()
	defer mu.RUnlock()
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles,
		RequestID: c.GetString("request_id"),
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Bad Request",
			Message:   "Invalid ID",
			RequestID: c.GetString("request_id"),
		})
		return
	}
	mu.RLock()
	defer mu.RUnlock()
	// TODO: Find article by ID
	article, i := findArticleByID(id)
	// TODO: Return 404 if not found
	if i == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Not Found",
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
	// TODO: Parse JSON request body
	var article Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Bad Request",
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}
	// TODO: Validate required fields
	if err := validateArticle(article); err != nil {
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Success:   false,
			Error:     "Unprocessable Entity",
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}
	// TODO: Add article to storage
	mu.Lock()
	defer mu.Unlock()
	article.ID = nextID
	nextID++
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	articles = append(articles, article)
	// TODO: Return created article
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Data:      article,
		RequestID: c.GetString("request_id"),
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Bad Request",
			Message:   "Invalid ID",
			RequestID: c.GetString("request_id"),
		})
		return
	}
	// TODO: Parse JSON request body
	var newArticle *Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Bad Request",
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}
	if err := validateArticle(*newArticle); err != nil {
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Success:   false,
			Error:     "Unprocessable Entity",
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	// TODO: Find and update article
	_, i := findArticleByID(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Not Found",
			Message:   "ID Not Found",
			RequestID: c.GetString("request_id"),
		})
		return
	}
	// TODO: Return updated article
	newArticle.ID = articles[i].ID
	newArticle.CreatedAt = articles[i].CreatedAt
	newArticle.UpdatedAt = time.Now()
	articles[i] = *newArticle
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      newArticle,
		RequestID: c.GetString("request_id"),
	})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Bad Request",
			Message:   "Invalid ID",
			RequestID: c.GetString("request_id"),
		})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	// TODO: Find and remove article
	_, i := findArticleByID(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Not Found",
			Message:   "ID Not Found",
			RequestID: c.GetString("request_id"),
		})
		return
	}
	// TODO: Return success message
	articles = append(articles[:i], articles[i+1:]...)
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: fmt.Sprintf("Article ID = %d Successfully Deleted", id),
	})
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	// TODO: Return mock statistics
	mu.RLock()
	defer mu.RUnlock()
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": requestCounter.Load(),
		"uptime":         time.Since(startTime).String(),
	}

	// TODO: Return stats in standard format
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      stats,
		RequestID: c.GetString("request_id"),
	})
}

// Helper functions

// check role access
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole := ctx.GetString("user_role")
		if userRole != requiredRole {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Success: false,
				Error:   "Insufficient permissions",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// TODO: Implement article lookup
	for i, article := range articles {
		if article.ID == id {
			return &articles[i], i
		}
	}
	// Return article pointer and index, or nil and -1 if not found
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	// TODO: Implement validation
	// Check required fields: Title, Content, Author
	if article.Title == "" {
		return errors.New("Title is required")
	}
	if article.Content == "" {
		return errors.New("Content is required")
	}
	if article.Author == "" {
		return errors.New("Author is required")
	}
	return nil
}
