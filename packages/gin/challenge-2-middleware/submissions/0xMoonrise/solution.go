package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
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
	{ID: 1,
		Title:     "Getting Started with Go",
		Content:   "Go is a programming language...",
		Author:    "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()},

	{ID: 2,
		Title:     "Web Development with Gin",
		Content:   "Gin is a web framework...",
		Author:    "Jane Smith",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()},
}
var nextID = 3

func main() {
	router := gin.New()

	router.Use(ErrorHandlerMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(LoggingMiddleware())
	router.Use(CORSMiddleware())
	router.Use(RateLimitMiddleware())
	router.Use(ContentTypeMiddleware())
	router.Use(RequestCounterMiddleware())

	router.GET("/ping", ping)
	router.GET("/articles", getArticles)
	router.GET("/articles/:id", getArticle)

	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/articles", createArticle)
		protected.PUT("/articles/:id", updateArticle)
		protected.DELETE("/articles/:id", deleteArticle)
		protected.GET("/admin/stats", getStats)
	}

	addr := net.JoinHostPort("localhost", "8081")
	router.Run(addr)
}

var servedRequests int64

func RequestCounterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		atomic.AddInt64(&servedRequests, 1)
	}
}

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		request_id, _ := c.Get("request_id")
		path := c.Request.URL.Path
		agent := c.Request.UserAgent()
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		entry := fmt.Sprintf("[%s] %s %s %d %s %s %s",
			request_id,
			method,
			path,
			statusCode,
			latency,
			clientIP,
			agent,
		)

		if c.Writer.Status() >= 400 {
			log.Printf("%s%s", "ERROR:", entry)
		} else {
			log.Printf("%-6s%s", "INFO:", entry)
		}
	}
}

type role string

const (
	Admin role = "admin"
	User  role = "user"
)

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {

	roles := make(map[string]role)
	roles["admin-key-123"] = Admin
	roles["user-key-456"] = User

	return func(c *gin.Context) {
		key := c.Request.Header.Get("X-API-Key")
		r, ok := roles[key]

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid API key",
			})
			return
		}

		c.Set("role", r)
		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"https://myapp.com":     true,
		}

		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(requests int, duration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
	}

	go rl.cleanupVisitors()

	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(time.Minute), 100) // 100 requests per minute
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

var rateLimiter = NewRateLimiter(100, time.Minute)

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := rateLimiter.getVisitor(c.ClientIP())

		limit := 100
		reset := time.Now().Add(time.Minute).Unix()

		allowed := limiter.Allow()
		remaining := max(0, int(limiter.Tokens()))

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(reset, 10))

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, APIResponse{
				Success: false,
				Error:   "Rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {

			if c.ContentType() != "application/json" {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, APIResponse{
					Success: false,
					Error:   "Invalid content type",
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
		requestID, exists := c.Get("request_id")
		reqIDStr := ""
		if exists {
			if id, ok := requestID.(string); ok {
				reqIDStr = id
			}
		}

		var errMsg string
		switch v := recovered.(type) {
		case error:
			errMsg = v.Error()
		case string:
			errMsg = v
		default:
			errMsg = fmt.Sprintf("%v", recovered)
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Internal server error",
			Message:   errMsg,
			RequestID: reqIDStr,
		})
	})
}

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	request_id, ok := c.Get("request_id")

	if !ok {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "something went wrong while processing the request.",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		RequestID: request_id.(string),
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	request_id, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "List of articles",
		RequestID: request_id.(string),
		Data:      articles,
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid id",
		})
		return
	}

	if len(articles) < id || id <= 0 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Id not found",
		})
		return
	}
	article, i := findArticleByID(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Id not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    article,
	})
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	var json Article

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid json",
		})
		return
	}

	if err := validateArticle(json); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	json.ID = nextID
	json.CreatedAt = time.Now()
	json.UpdatedAt = time.Now()
	nextID++
	articles = append(articles, json)

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "The article has been created",
		Data:    json,
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid id",
		})
		return
	}

	if len(articles) < id || id <= 0 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Id not found",
		})
		return
	}

	var updatedArticle Article

	if err := c.ShouldBindJSON(&updatedArticle); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid json",
		})
		return
	}

	if err := validateArticle(updatedArticle); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	originalArticle := articles[id-1]
	updatedArticle.ID = originalArticle.ID
	updatedArticle.CreatedAt = originalArticle.CreatedAt
	updatedArticle.UpdatedAt = time.Now()
	articles[id-1] = updatedArticle

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "The article has been updated",
		Data:    updatedArticle,
	})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid id",
		})
		return
	}

	_, i := findArticleByID(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Id not found",
		})
		return
	}

	articles = append(articles[:i], articles[i+1:]...)
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "The article has been deleted",
	})

}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {

	user_role, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "No role was set",
		})
		return
	}

	if user_role != Admin {
		c.JSON(http.StatusForbidden, APIResponse{
			Success: false,
			Error:   "Your role has no privileges",
		})
		return
	}
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": atomic.LoadInt64(&servedRequests),
		"uptime":         "24h",
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Server stats",
		Data:    stats,
	})
}

// Helper functions
// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	for i, article := range articles {
		if article.ID == id {
			return &article, i
		}
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	if article.Title == "" {
		return errors.New("The 'title' field can't be empty")
	}

	if article.Content == "" {
		return errors.New("The 'content' field can't be empty")
	}

	if article.Author == "" {
		return errors.New("The 'author' field can't be empty")
	}

	return nil
}
