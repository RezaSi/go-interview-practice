package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"slices"
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

var requestIdKey = "request_id"
var roleKey = "role"

func responseSuccess(c *gin.Context, data interface{}) APIResponse {
	requestId, _ := c.Get(requestIdKey)
	return APIResponse{
		Success:   true,
		Data:      data,
		Message:   "Success",
		RequestID: requestId.(string),
	}
}
func responseFailure(c *gin.Context, err string) APIResponse {
	requestId, _ := c.Get(requestIdKey)
	return APIResponse{
		Success:   false,
		Error:     err,
		RequestID: requestId.(string),
	}
}
func responseFailureMessage(c *gin.Context, err string, msg string) APIResponse {
	requestId, _ := c.Get(requestIdKey)
	return APIResponse{
		Success:   false,
		Error:     err,
		Message:   msg,
		RequestID: requestId.(string),
	}
}

// In-memory storage
var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var articlesLock sync.RWMutex
var nextID = 3
var nextIdLock sync.Mutex

func withNextId(f func(int)) {
	nextIdLock.Lock()
	defer nextIdLock.Unlock()
	id := nextID
	nextID++
	f(id)
}

func main() {
	// Create Gin router without default middleware
	// Use gin.New() instead of gin.Default()
	r := gin.New()

	// Setup custom middleware in correct order
	// 1. ErrorHandlerMiddleware (first to catch panics)
	r.Use(ErrorHandlerMiddleware())
	// 2. RequestIDMiddleware
	r.Use(RequestIDMiddleware())
	// 3. LoggingMiddleware
	r.Use(LoggingMiddleware())
	// 4. CORSMiddleware
	r.Use(CORSMiddleware())
	// 5. RateLimitMiddleware
	r.Use(RateLimitMiddleware())
	// 6. ContentTypeMiddleware
	r.Use(ContentTypeMiddleware())

	// Setup route groups
	// Public routes (no authentication required)
	public := r.Group("/")
	// Protected routes (require authentication)
	protected := r.Group("/", AuthMiddleware())

	// Define routes
	// Public: GET /ping, GET /articles, GET /articles/:id
	{
		public.GET("/ping", ping)
		public.GET("/articles", getArticles)
		public.GET("/articles/:id", getArticle)
	}
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats
	{
		protected.POST("/articles", createArticle)
		protected.PUT("/articles/:id", updateArticle)
		protected.DELETE("/articles/:id", deleteArticle)
		protected.GET("/admin/stats", getStats)
	}

	// Start server on port 8080
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

// Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate UUID for request ID
		// Use github.com/google/uuid package
		uuid := uuid.New().String()
		// Store in context as "request_id"
		c.Set(requestIdKey, uuid)
		// Add to response header as "X-Request-ID"
		c.Header("X-Request-Id", uuid)

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
		duration := time.Since(start)
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		requestId, _ := c.Get(requestIdKey)
		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		log.Printf("[%s] %s %s %d %v %s %s", requestId, method, path, status, duration, ip, userAgent)
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// Define valid API keys and their roles
	apiKeyAndRoleMap := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		// Get API key from X-API-Key header
		requestApiKey := c.GetHeader("X-Api-Key")
		// Validate API key
		role, ok := apiKeyAndRoleMap[requestApiKey]
		if !ok {
			//Return 401 if invalid or missing
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Set user role in context
		c.Set(roleKey, role)
		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	allowOrigin := []string{
		"http://localhost:3000",
		"https://myblog.com",
	}
	return func(c *gin.Context) {
		// Set CORS headers
		origin := c.GetHeader("Origin")
		if slices.Contains(allowOrigin, origin) {
			// Allow origins: http://localhost:3000, https://myblog.com
			c.Header("Access-Control-Allow-Origin", origin)
			// Allow methods: GET, POST, PUT, DELETE, OPTIONS
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			// Allow headers: Content-Type, X-API-Key, X-Request-ID
			c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")

			// Handle preflight OPTIONS requests
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	// Implement rate limiting
	// Limit: 100 requests per IP per minute
	requestsPerIpPerMinute := 100
	ipAndLimiterMap := new(sync.Map)

	return func(c *gin.Context) {
		// Use golang.org/x/time/rate package
		// this is not "exactly" per minute, but let's do the simple way...
		l, _ := ipAndLimiterMap.LoadOrStore(c.ClientIP(),
			//rate.NewLimiter(rate.Every(time.Minute), requestsPerIpPerMinute),
			rate.NewLimiter(rate.Limit((requestsPerIpPerMinute)/60.0), requestsPerIpPerMinute),
		)
		limiter := l.(*rate.Limiter)
		// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
		c.Header("X-RateLimit-Limit", fmt.Sprint(requestsPerIpPerMinute))
		remain := math.Floor(limiter.Tokens() - 1)
		if remain < 0 {
			remain = 0
		}
		c.Header("X-RateLimit-Remaining", fmt.Sprint(remain))
		// I don't see how I can implement the header with golang.org/x/time/rate...
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Unix()+120, 10))
		// Return 429 if rate limit exceeded
		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check content type for POST/PUT requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			// Must be application/json
			if strings.Index(c.ContentType(), "application/json") != 0 {
				// Return 415 if invalid content type
				c.AbortWithStatus(http.StatusUnsupportedMediaType)
				return
			}
		}

		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	//gin.Recovery()
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Handle panics gracefully
		if recovered != nil {
			// Return consistent error response format
			// Include request ID in response
			var errMsg string
			switch v := recovered.(type) {
			case string:
				errMsg = v
			case error:
				errMsg = v.Error()
			default:
				errMsg = "unknown error"
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, responseFailureMessage(c, "Internal server error", errMsg))
			return
		}
	})
}

// Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// Return simple pong response with request ID
	c.JSON(http.StatusOK, responseSuccess(c, "pong"))
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// Implement pagination (optional)
	pageNumString := c.DefaultQuery("page", "1")
	pageNum, err := strconv.Atoi(pageNumString)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid pageNum: "+err.Error()))
		return
	}
	if pageNum <= 0 {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid pageNum"))
		return
	}
	pageSizeString := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid pageSize: "+err.Error()))
		return
	}
	if pageSize <= 0 {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid pageSize"))
		return
	}

	startIndex := (pageNum - 1) * pageSize
	endIndex := startIndex + pageSize
	articlesLock.RLock()
	defer articlesLock.RUnlock()
	var result []Article
	if startIndex >= len(articles) {
		result = []Article{}
	} else {
		if endIndex >= len(articles) {
			endIndex = len(articles) - 1
		}
		result = articles[startIndex : endIndex+1]
	}
	// Return articles in standard format
	c.JSON(http.StatusOK, responseSuccess(c, result))
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// Get article ID from URL parameter
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid id: "+err.Error()))
	}
	// Find article by ID
	articlesLock.RLock()
	defer articlesLock.RUnlock()
	article, _ := findArticleByID(id)
	if article == nil {
		// Return 404 if not found
		c.JSON(http.StatusNotFound, responseFailure(c, "article not found"))
		return
	}
	c.JSON(http.StatusOK, responseSuccess(c, article))
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// Parse JSON request body
	var articleParam Article
	if err := c.ShouldBind(&articleParam); err != nil {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid body: "+err.Error()))
		return
	}
	// Validate required fields
	if err := validateArticle(articleParam); err != nil {
		c.JSON(http.StatusBadRequest, responseFailure(c, err.Error()))
		return
	}
	// Add article to storage
	var article Article
	now := time.Now()
	withNextId(func(id int) {
		article = Article{
			ID:        id,
			Title:     articleParam.Title,
			Content:   articleParam.Content,
			Author:    articleParam.Author,
			CreatedAt: now,
			UpdatedAt: now,
		}
	})
	articlesLock.Lock()
	defer articlesLock.Unlock()
	articles = append(articles, article)
	// Return created article
	c.JSON(http.StatusCreated, responseSuccess(c, article))
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// Get article ID from URL parameter
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid id: "+err.Error()))
	}
	// Parse JSON request body
	var articleParam Article
	if err := c.ShouldBind(&articleParam); err != nil {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid body: "+err.Error()))
		return
	}
	// Find and update article
	articlesLock.Lock()
	defer articlesLock.Unlock()
	article, _ := findArticleByID(id)
	if article == nil {
		// Return 404 if not found
		c.JSON(http.StatusNotFound, responseFailure(c, "article not found"))
		return
	}
	article.Title = articleParam.Title
	article.Content = articleParam.Content
	article.Author = articleParam.Author
	// Return updated article
	c.JSON(http.StatusOK, responseSuccess(c, article))
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// Get article ID from URL parameter
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseFailure(c, "invalid id: "+err.Error()))
	}
	// Find and remove article
	articlesLock.Lock()
	defer articlesLock.Unlock()
	_, index := findArticleByID(id)
	if index < 0 {
		// Return 404 if not found
		c.JSON(http.StatusNotFound, responseFailure(c, "article not found"))
		return
	}
	articles = append(articles[:index], articles[index+1:]...)
	// Return success message
	c.JSON(http.StatusOK, responseSuccess(c, nil))
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// Check if user role is "admin"
	role, ok := c.Get("role")
	if !ok || role.(string) != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	// Return mock statistics
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}
	// Return stats in standard format
	c.JSON(http.StatusOK, responseSuccess(c, stats))
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	for i := 0; i < len(articles); i++ {
		if articles[i].ID == id {
			return &articles[i], i
		}
	}
	// Return article pointer and index, or nil and -1 if not found
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	// Check required fields: Title, Content, Author
	if article.Title == "" {
		return errors.New("article title is empty")
	}
	if article.Content == "" {
		return errors.New("article content is empty")
	}
	if article.Author == "" {
		return errors.New("article author is empty")
	}
	return nil
}
