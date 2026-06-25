package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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

// In-memory storage for articles and the next ID to assign, guarded by mu.
var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3
var mu sync.RWMutex
var limiters = make(map[string]*rate.Limiter)
var limitersCash = make(map[string]time.Time)
var ttl = time.Duration(time.Minute)
var limitersMu sync.RWMutex
var startedAt = time.Now()
var totalRequests atomic.Uint64

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

	protectedRoute := router.Group("/", AuthMiddleware())
	{
		protectedRoute.POST("/articles", createArticle)
		protectedRoute.PUT("/articles/:id", updateArticle)
		protectedRoute.DELETE("/articles/:id", deleteArticle)
		protectedRoute.GET("/admin/stats", getStats)
	}

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()
	go limitersCleanup(ctx)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		// requestID := c.GetString("RequestID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		totalRequests.Add(1)

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
			log.Printf("INFO %+v", entry)
		}
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const (
			admin = "admin-key-123"
			user  = "user-key-456"
		)
		if admin == "" || user == "" {
			errorResponse(c, "", "Authensication not configured", http.StatusInternalServerError)
			c.Abort()
			return
		}

		key := c.GetHeader("X-API-Key")
		switch key {
		case admin:
			c.Set("user_role", "admin")
		case user:
			c.Set("user_role", "user")
		default:
			errorResponse(c, "invalid key", "invalid API key", http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigines := map[string]bool{
			"http://localhost:3000": true,
			"https://myblog.com":    true,
		}

		if allowedOrigines[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var limit = 100
		ip := c.ClientIP()
		limitersMu.Lock()
		limiter, ok := limiters[ip]
		if !ok {
			limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(limit)), limit*2)
			limiters[ip] = limiter
		}
		limitersCash[ip] = time.Now()

		if !limiter.Allow() {
			limitersMu.Unlock()
			errorResponse(c, "too many requests", "too many requests", http.StatusTooManyRequests)
			c.Abort()
			return
		}
		limitersMu.Unlock()
		l := strconv.Itoa(limit)
		c.Header("X-RateLimit-Limit", l)
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", int(limiter.Tokens())))

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			if !strings.HasPrefix(c.Request.Header.Get("Content-Type"), "application/json") {
				errorResponse(c, "", "unsupported media type", http.StatusUnsupportedMediaType)
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
		errorResponse(c, fmt.Sprintf("%v", recovered), "Internal server error", http.StatusInternalServerError)
		c.Abort()
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
	art := append([]Article{}, articles...)
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
		errorResponse(c, "", "bad request", http.StatusBadRequest)
		return
	}
	mu.RLock()
	articlesCp := append([]Article{}, articles...)
	mu.RUnlock()
	for _, article := range articlesCp {
		if article.ID == id {
			successResponse(c, article, "success", http.StatusOK)
			return
		}
	}
	errorResponse(c, "", "not found", http.StatusNotFound)
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// TODO: Parse JSON request body
	role, _ := c.Get("user_role")

	if role != "admin" {
		errorResponse(c, "not allowed", "forbidden", http.StatusForbidden)
		return
	}

	article := Article{}
	if err := c.ShouldBindJSON(&article); err != nil {
		errorResponse(c, "", "failed to parse JSON", http.StatusBadRequest)
		return
	}
	// TODO: Validate required fields
	if err := validateArticle(article); err != nil {
		errorResponse(c, "", "Invalid article data", http.StatusBadRequest)
		return
	}
	// TODO: Add article to storage
	mu.Lock()
	article.ID = nextID
	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now
	nextID++
	articles = append(articles, article)
	mu.Unlock()
	// TODO: Return created article
	successResponse(c, article, "Successfully added", http.StatusCreated)
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	role, _ := c.Get("user_role")
	if role == "user" {
		errorResponse(c, "not allowed", "forbiden", http.StatusForbidden)
		return
	}
	if role != "user" && role != "admin" {
		errorResponse(c, "not allowed", "not authorised", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, "", "article ID should be integer", http.StatusBadRequest)
		return
	}
	dataForUpdate := Article{}
	if err := c.ShouldBindJSON(&dataForUpdate); err != nil {
		errorResponse(c, "", err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	article, i := findArticleByID(id)
	if i == -1 {
		errorResponse(c, "", "not found", http.StatusNotFound)
		mu.Unlock()
		return
	}
	if dataForUpdate.Author == "" {
		dataForUpdate.Author = article.Author
	}
	if dataForUpdate.Content == "" {
		dataForUpdate.Content = article.Content
	}
	if dataForUpdate.Title == "" {
		dataForUpdate.Title = article.Title
	}

	dataForUpdate.UpdatedAt = time.Now()

	dataForUpdate.CreatedAt = article.CreatedAt

	dataForUpdate.ID = article.ID

	if err := validateArticle(dataForUpdate); err != nil {
		errorResponse(c, "", err.Error(), http.StatusBadRequest)
		mu.Unlock()
		return
	}
	articles[i] = dataForUpdate
	mu.Unlock()
	successResponse(c, dataForUpdate, "success", http.StatusOK)
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	role, _ := c.Get("user_role")
	switch role {
	case "user":
		errorResponse(c, "not allowed", "forbidden", http.StatusForbidden)
		return

	case "admin":
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			errorResponse(c, "", "article ID should be integer", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		_, i := findArticleByID(id)
		if i == -1 {
			errorResponse(c, "", "not found", http.StatusNotFound)
			return
		}
		articles = append(articles[:i], articles[i+1:]...)
		successResponse(c, nil, "deleted successfully", http.StatusOK)

	default:
		errorResponse(c, "not allowed", "not authorize", http.StatusUnauthorized)
		return
	}
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	userRole := c.GetString("user_role")
	mu.RLock()
	artLen := len(articles)
	mu.RUnlock()
	switch userRole {
	case "admin":
		stats := map[string]interface{}{
			"total_articles": artLen,
			"total_requests": totalRequests.Load(), // Could track this in middleware
			"uptime":         time.Since(startedAt).String(),
		}
		successResponse(c, stats, "success", http.StatusOK)
		return
	case "user":
		errorResponse(c, "", "forbidden", http.StatusForbidden)
	default:
		errorResponse(c, "", "not authorize", http.StatusUnauthorized)
	}
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

func limitersCleanup(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			limitersMu.Lock()
			for addr, liveSince := range limitersCash {
				if time.Since(liveSince) > ttl {
					delete(limitersCash, addr)
					delete(limiters, addr)
				}
			}
			limitersMu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func errorResponse(c *gin.Context, msg, err string, status int) {
	response := APIResponse{
		Success:   false,
		Message:   msg,
		Error:     err,
		RequestID: c.GetString("RequestID"),
	}

	c.JSON(status, response)
}

func successResponse(c *gin.Context, data interface{}, msg string, status int) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Message:   msg,
		RequestID: c.GetString("RequestID"),
	}

	c.JSON(status, response)
}
