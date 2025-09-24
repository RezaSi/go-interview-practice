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
// =====================================================================
var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3

var (
	ErrDataRepositoryEmpty      = errors.New("not a single article was found")
	ErrDataRepositoryIdNotFound = errors.New("no article with this ID was found")
	ErrDataRepositoryCantCreate = errors.New("article is invalid, cannot create article")
)

type DataStore struct {
	data map[int]Article
	cnt  int
	mu   sync.RWMutex
}

var globalService *DataStore
var rateLimiters = make(map[string]*rate.Limiter)
var rateLimitersMu sync.Mutex
var requestsPerSecond = 100

func NewDataStore() *DataStore {
	return &DataStore{
		data: map[int]Article{
			1: {ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			2: {ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		cnt: 3,
	}
}

func (d *DataStore) GetArticles() []Article {
	d.mu.RLock()
	defer d.mu.RUnlock()

	retData := make([]Article, 0)

	for _, v := range d.data {
		retData = append(retData, v)
	}

	return retData
}

func (d *DataStore) GetArticle(id int) (Article, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	retData, ok := d.data[id]
	if !ok {
		return retData, ErrDataRepositoryIdNotFound
	}

	return retData, nil
}

func (d *DataStore) CreateArticle(data Article) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := validateArticle(data); err != nil {
		return -1, err
	}

	data.ID = d.cnt
	//for test
	//d.data[d.cnt] = data
	//d.cnt++

	return data.ID, nil
}

func (d *DataStore) UpdateArticle(id int, data Article) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := validateArticle(data); err != nil {
		return -1, err
	}

	if _, ok := d.data[id]; !ok {
		return -1, ErrDataRepositoryIdNotFound
	}

	data.ID = id
	d.data[id] = data

	return data.ID, nil
}

func (d *DataStore) DeleteArticle(id int) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.data[id]; !ok {
		return -1, ErrDataRepositoryIdNotFound
	}

	delete(d.data, id)
	d.cnt--

	return id, nil
}

func validateArticle(data Article) error {

	if data.Title == "" {
		return fmt.Errorf("%w: Title is empty", ErrDataRepositoryCantCreate)
	}

	if data.Content == "" {
		return fmt.Errorf("%w: Content is empty", ErrDataRepositoryCantCreate)
	}

	if data.Author == "" {
		return fmt.Errorf("%w: Author is empty", ErrDataRepositoryCantCreate)
	}

	return nil
}

// Helper
// =======================================================================
func sendJSONResponse(c *gin.Context, status int, success bool, data interface{}, message string, err string) {
	requestID, _ := c.Get("request_id")
	response := APIResponse{
		Success:   success,
		Data:      data,
		Message:   message,
		Error:     err,
		RequestID: fmt.Sprintf("%s", requestID),
	}
	c.JSON(status, response)
}

func sendErrorResponse(c *gin.Context, status int, message string) {
	sendJSONResponse(c, status, false, nil, "", message)
}

func getServiceData() *DataStore {
	if globalService == nil {
		globalService = NewDataStore()
	}

	return globalService
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID, _ := c.Get("request_id")

		c.Next()

		duration := time.Since(start)
		log.Printf("[%s] %s %s %d %s %s %s",
			requestID,
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
	validAPIKeys := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			sendJSONResponse(c, http.StatusUnauthorized, false, nil, "API key required", "MISSING_API_KEY")
			c.Abort()
			return
		}

		role, exists := validAPIKeys[apiKey]
		if !exists {
			sendJSONResponse(c, http.StatusUnauthorized, false, nil, "Invalid API key", "INVALID_API_KEY")
			c.Abort()
			return
		}

		c.Set("user_role", role)
		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := []string{"http://localhost:3000"}
		origin := c.GetHeader("Origin")

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		ip := c.ClientIP()
		rateLimitersMu.Lock()
		limiter, exists := rateLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond*2)
			rateLimiters[ip] = limiter
		}
		rateLimitersMu.Unlock()

		c.Writer.Header().Set("X-RateLimit-Limit", "100")
		c.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

		if !limiter.Allow() {
			c.Writer.Header().Set("X-RateLimit-Remaining", "0")
			sendErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded")
			c.Abort()
			return
		}

		remaining := int(limiter.Tokens())
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Next()
	}

}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				sendJSONResponse(c, http.StatusUnsupportedMediaType, false, nil, "Content-Type must be application/json", "INVALID_CONTENT_TYPE")
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
		sendJSONResponse(c, http.StatusInternalServerError, false, nil, fmt.Sprintf("%v", recovered), "Internal server error")
		c.Abort()
	})
}

// Handler
// =======================================================================

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	sendJSONResponse(c, http.StatusOK, true, "pong", "Data retrieved successfully", "")
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	data := getServiceData().GetArticles()
	sendJSONResponse(c, http.StatusOK, true, data, "Data retrieved successfully", "")
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid data ID")
		return
	}

	data, err := getServiceData().GetArticle(id)
	if err != nil {
		sendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	sendJSONResponse(c, http.StatusOK, true, data, "Data retrieved successfully", "")

}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	var data Article

	if err := c.ShouldBind(&data); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid data: "+err.Error())
		return

	}

	if _, err := getServiceData().CreateArticle(data); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sendJSONResponse(c, http.StatusCreated, true, data, "Create successfully", "")

}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid data ID")
		return
	}

	var data Article
	if err := c.ShouldBind(&data); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid data: "+err.Error())
		return

	}

	if _, err := getServiceData().UpdateArticle(id, data); err != nil {
		sendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	sendJSONResponse(c, http.StatusOK, true, data, "Update successfully", "")

}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid data ID")
		return
	}

	if _, err := getServiceData().DeleteArticle(id); err != nil {
		sendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	sendJSONResponse(c, http.StatusOK, true, id, "Delete successfully", "")

}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	userRole, exists := c.Get("user_role")
	if !exists || userRole != "admin" {
		sendJSONResponse(c, http.StatusForbidden, false, nil, "Admin access required", "ACCESS_DENIED")
		return
	}

	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0,
		"uptime":         "24h",
	}

	sendJSONResponse(c, http.StatusOK, true, stats, "Stats retrieved successfully", "")
}

// =====================================================================
func main() {
	r := gin.Default()
	r.Use(ErrorHandlerMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(LoggingMiddleware())
	r.Use(CORSMiddleware())
	r.Use(RateLimitMiddleware())
	r.Use(ContentTypeMiddleware())

	// Public routes
	r.GET("/ping", func(c *gin.Context) { ping(c) })
	r.GET("/articles", func(c *gin.Context) { getArticles(c) })
	r.GET("/articles/:id", func(c *gin.Context) { getArticle(c) })

	// Protected routes
	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/articles", func(c *gin.Context) { createArticle(c) })
		protected.PUT("/articles/:id", func(c *gin.Context) { updateArticle(c) })
		protected.DELETE("/articles/:id", func(c *gin.Context) { deleteArticle(c) })
		protected.GET("/admin/stats", func(c *gin.Context) { getStats(c) })
	}

	log.Println("http://localhost:8085")
	r.Run(":8085")
}
