package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang.org/x/time/rate"
)

const (
	RequestIDKey       = "request_id"
	UserRoleKey        = "user_role"
	RequestIDHeaderKey = "X-Request-ID"
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

// LRUCache implements a thread-safe LRU cache for rate limiters
type LRUCache struct {
	mu      sync.RWMutex
	cache   map[string]*lruNode
	head    *lruNode
	tail    *lruNode
	maxSize int
}

// lruNode represents a node in the doubly-linked list
type lruNode struct {
	key   string
	value *rate.Limiter
	next  *lruNode
	prev  *lruNode
}

// NewLRUCache creates a new LRU cache with the given maximum size
func NewLRUCache(maxSize int) *LRUCache {
	return &LRUCache{
		cache:   make(map[string]*lruNode),
		maxSize: maxSize,
	}
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (*rate.Limiter, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// Move to front (most recently used)
	c.moveToFront(node)
	return node.value, true
}

// Put adds a value to the cache
func (c *LRUCache) Put(key string, value *rate.Limiter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If key already exists, update it
	if node, exists := c.cache[key]; exists {
		node.value = value
		c.moveToFront(node)
		return
	}

	// Create new node
	newNode := &lruNode{
		key:   key,
		value: value,
	}

	// Add to front
	c.addToFront(newNode)
	c.cache[key] = newNode

	// Check if we need to evict
	if len(c.cache) > c.maxSize {
		c.evict()
	}
}

// Remove removes a key from the cache
func (c *LRUCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.removeNode(node)
		delete(c.cache, key)
	}
}

// moveToFront moves a node to the front of the list
func (c *LRUCache) moveToFront(node *lruNode) {
	if c.head == node {
		return
	}

	c.removeNode(node)
	c.addToFront(node)
}

// addToFront adds a node to the front of the list
func (c *LRUCache) addToFront(node *lruNode) {
	if c.head == nil {
		c.head = node
		c.tail = node
	} else {
		node.next = c.head
		c.head.prev = node
		c.head = node
	}
}

// removeNode removes a node from the list
func (c *LRUCache) removeNode(node *lruNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}
}

// evict removes the least recently used item
func (c *LRUCache) evict() {
	if c.tail == nil {
		return
	}

	// Remove from cache map
	delete(c.cache, c.tail.key)

	// Remove from list
	if c.tail.prev != nil {
		c.tail.prev.next = nil
		c.tail = c.tail.prev
	} else {
		// Only one node
		c.head = nil
		c.tail = nil
	}
}

// In-memory storage
var (
	articlesMutex sync.RWMutex
	articles      = []Article{
		{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	nextID = 3

	ipLimiters = NewLRUCache(1000)

	keys = map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}
)

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

	// Define and Setup route groups
	// Define routes
	// Public routes (no authentication required)
	// Public: GET /ping, GET /articles, GET /articles/:id
	public := r.Group("/")
	public.GET("/ping", ping)
	public.GET("/articles", getArticles)
	public.GET("/articles/:id", getArticle)

	// Protected routes (require authentication)
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats
	protected := r.Group("/").Use(AuthMiddleware())
	protected.POST("/articles", ContentTypeMiddleware(), createArticle)
	protected.PUT("/articles/:id", ContentTypeMiddleware(), updateArticle)
	protected.DELETE("/articles/:id", deleteArticle)
	protected.GET("/admin/stats", getStats)

	// Start server on port 8080
	r.Run(":8080")
}

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate UUID for request ID
		// Use github.com/google/uuid package
		requestID := uuid.New().String()
		// Store in context as "request_id"
		c.Set(RequestIDKey, requestID)
		// Add to response header as "X-Request-ID"
		c.Header(RequestIDHeaderKey, requestID)

		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture start time
		start := time.Now()
		c.Set("start_time", start)

		c.Next()

		// Calculate duration and log request
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		requestID := c.GetString(RequestIDKey)
		duration := time.Since(c.GetTime("start_time"))
		log.Printf("[%s] %s %s %d %s %s %s", requestID, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration, c.ClientIP(), c.Request.UserAgent())
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from X-API-Key header
		headerKey := c.GetHeader("X-API-Key")

		// Validate API key
		role, ok := keys[headerKey]
		if !ok {
			// Return 401 if invalid or missing
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success:   false,
				RequestID: c.GetString(RequestIDKey),
			})
			return
		}

		// Set user role in context
		c.Set(UserRoleKey, role)

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers

		// Allow origins: http://localhost:3000, https://myblog.com
		origin := c.GetHeader("Origin")
		allowedOrigins := []string{"http://localhost:3000", "https://myblog.com"}
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
				break
			}
		}

		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow headers: Content-Type, X-API-Key, X-Request-ID
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")

		// Handle preflight OPTIONS requests by returning 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	// Limit: 100 requests per IP per minute
	// Use golang.org/x/time/rate package
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter, exists := ipLimiters.Get(clientIP)
		if !exists {
			limiter = rate.NewLimiter(rate.Every(time.Minute/100.0), 100)
			ipLimiters.Put(clientIP, limiter)
		}
		// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
		c.Header("X-RateLimit-Limit", "100")
		// Token bucket rate limiters refill continuously at a constant rate rather than resetting at a fixed time,
		// This header value is an approximation of the reset time.
		c.Header("X-RateLimit-Reset", strconv.Itoa(int(time.Now().Add(time.Minute).UnixMilli())))
		if !limiter.Allow() {
			c.Header("X-RateLimit-Remaining", "0")
			// Return 429 if rate limit exceeded
			c.AbortWithStatusJSON(http.StatusTooManyRequests, APIResponse{
				Success:   false,
				RequestID: c.GetString(RequestIDKey),
			})
			return
		}
		c.Header("X-RateLimit-Remaining", strconv.Itoa(int(math.Round(limiter.Tokens()))))

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check content type for POST/PUT requests
		// Must be application/json
		// Return 415 if invalid content type
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			if c.ContentType() != "application/json" {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, APIResponse{
					Success:   false,
					RequestID: c.GetString(RequestIDKey),
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
		// Handle panics gracefully
		requestID := c.GetString(RequestIDKey)

		// Return consistent error response format
		// Include request ID in response
		c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Internal server error",
			Message:   fmt.Sprintf("%v", recovered),
			RequestID: requestID,
		})
	})
}

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// Return simple pong response with request ID
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      "pong",
		RequestID: c.GetString(RequestIDKey),
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// todo add pagination ? optional
	articlesMutex.RLock()
	defer articlesMutex.RUnlock()
	// Return articles in standard format
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles,
		RequestID: c.GetString(RequestIDKey),
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// Get article ID from URL parameter
	id, err := parseIDParam(c)
	if err != nil {
		return
	}
	// Find article by ID
	articlesMutex.RLock()
	defer articlesMutex.RUnlock()
	article, _ := findArticleByID(id)
	if article == nil {
		// return 404 if article not found
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Message:   "article not found",
			RequestID: c.GetString(RequestIDKey),
		})
		return
	}
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      article,
		RequestID: c.GetString(RequestIDKey),
	})
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// Parse JSON request body
	var inputArticle Article
	if err := c.ShouldBindJSON(&inputArticle); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Message:   err.Error(),
			RequestID: c.GetString(RequestIDKey),
		})
		return
	}

	// Validate required fields
	if err := validateArticle(inputArticle); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Message:   err.Error(),
			RequestID: c.GetString(RequestIDKey),
		})
		return
	}

	// Add article to storage
	now := time.Now()
	inputArticle.CreatedAt = now
	inputArticle.UpdatedAt = now
	articlesMutex.Lock()
	defer articlesMutex.Unlock()
	inputArticle.ID = nextID
	articles = append(articles, inputArticle)
	nextID++

	// Return created article
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Data:      inputArticle,
		RequestID: c.GetString(RequestIDKey),
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// Get article ID from URL parameter
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	// Parse JSON request body
	var inputArticle Article
	if err := c.ShouldBindJSON(&inputArticle); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Message:   err.Error(),
			RequestID: c.GetString(RequestIDKey),
		})
		return
	}

	// Validate required fields
	if err := validateArticle(inputArticle); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Message:   err.Error(),
			RequestID: c.GetString(RequestIDKey),
		})
		return
	}

	// Find and update article
	articlesMutex.Lock()
	defer articlesMutex.Unlock()
	_, idx := findArticleByID(id)
	if idx < 0 {
		// return 404 if article not found
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Message:   "article not found",
			RequestID: c.GetString(RequestIDKey),
		})
		return
	}
	articles[idx].Author = inputArticle.Author
	articles[idx].Content = inputArticle.Content
	articles[idx].Title = inputArticle.Title
	articles[idx].UpdatedAt = time.Now()

	// Return updated article
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles[idx],
		RequestID: c.GetString(RequestIDKey),
	})

}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	// Get article ID from URL parameter
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	// Find and remove article
	articlesMutex.Lock()
	defer articlesMutex.Unlock()
	_, idx := findArticleByID(id)
	if idx < 0 {
		// return 404 if article not found
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Message:   "article not found",
			RequestID: c.GetString(RequestIDKey),
		})
		return
	}
	articles = append(articles[:idx], articles[idx+1:]...)

	// Return success message
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "article deleted",
		RequestID: c.GetString(RequestIDKey),
	})
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// Check if user role is "admin"
	role := c.GetString(UserRoleKey)
	if role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
			Success:   false,
			RequestID: c.GetString(RequestIDKey),
		})
		return
	}

	articlesMutex.RLock()
	totalArticles := len(articles)
	articlesMutex.RUnlock()

	// Return mock statistics
	stats := map[string]interface{}{
		"total_articles": totalArticles,
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}

	// Return stats in standard format
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      stats,
		RequestID: c.GetString(RequestIDKey),
	})
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// Implement article lookup
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
	// Implement validation
	// Check required fields: Title, Content, Author
	if len(article.Title) == 0 {
		return errors.New("article validation failed : title is required")
	}
	if len(article.Content) == 0 {
		return errors.New("article validation failed : content is required")
	}
	if len(article.Author) == 0 {
		return errors.New("article validation failed : author is required")
	}
	return nil
}

// parseIDParam parses and validates the ID parameter from the URL
func parseIDParam(c *gin.Context) (int, error) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Message:   "bad id",
			Error:     err.Error(),
			RequestID: c.GetString(RequestIDKey),
		})
	}
	return id, err
}
