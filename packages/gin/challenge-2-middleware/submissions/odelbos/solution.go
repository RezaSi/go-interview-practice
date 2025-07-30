package main

import (
	"time"
	"net/http"
	"log"
	"strconv"
	"fmt"
	"strings"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// ----------------------------------------------------------------
// STRUCT & DATA
// ----------------------------------------------------------------

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

var (
	rateLimiters = make(map[string]*rate.Limiter)
	rateLimitMutex sync.Mutex
)

// ----------------------------------------------------------------
// Main
// ----------------------------------------------------------------

func main() {
	r := gin.New()

	r.Use(
		RequestIDMiddleware(),
		ErrorHandlerMiddleware(),
		LoggingMiddleware(),
		CORSMiddleware(),
		RateLimitMiddleware(),
		ContentTypeMiddleware(),
	)

	public := r.Group("/")
	{
		public.GET("/ping", ping)
		public.GET("/articles/:id", getArticle)
		public.GET("/articles", getArticles)
	}

	private := r.Group("/").Use(AuthMiddleware())
	{
		private.POST("/articles", createArticle)
		private.PUT("/articles/:id", updateArticle)
		private.DELETE("/articles/:id", deleteArticle)
		private.GET("/admin/stats", getStats)
	}

	r.Run(":8080")
}

// ----------------------------------------------------------------
// Middlewares
// ----------------------------------------------------------------

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String()
		c.Set("request_id", id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.Printf("[%s] %s %s %d %s %s %s",
			c.GetString("request_id"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(time.Now()),
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		role := ""
		switch key {
		case "admin-key-123":
			role = "admin"
		case "user-key-456":
			role = "user"
		default:
			errResponse(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Set("role", role)
		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,X-API-Key,X-Request-ID")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
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

	return func(c *gin.Context) {
		ip := c.ClientIP()
		rateLimitMutex.Lock()
		limiter, ok := rateLimiters[ip]
		if ! ok {
			limiter = rate.NewLimiter(rate.Every(time.Minute / 100), 100)
			rateLimiters[ip] = limiter
		}
		rateLimitMutex.Unlock()

		c.Writer.Header().Set("X-RateLimit-Limit", "100")
		c.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

		if ! limiter.Allow() {
			c.Writer.Header().Set("X-RateLimit-Remaining", "0")
			errResponse(c, http.StatusTooManyRequests, "Rate limit exceeded")
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
			if c.ContentType() != "application/json" {
				errResponse(c, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Internal server error",
			Message:   fmt.Sprintf("%v", recovered),
			RequestID: c.GetString("request_id"),
		})
		c.Abort()
	})
}

// ----------------------------------------------------------------
// Handlers
// ----------------------------------------------------------------

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	okResponse(c, http.StatusOK, "pong", nil)
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	okResponse(c, http.StatusOK, "Articles", articles)
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	article, _ := findArticleByID(id)
	if article == nil {
		errResponse(c, http.StatusNotFound, "Not found")
		return
	}
	okResponse(c, http.StatusOK, "Article", article)
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	var article Article
	if err := c.ShouldBindJSON(&article); err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}
	if err := validateArticle(article); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	article.ID = nextID
	article.CreatedAt = time.Now()
	article.UpdatedAt = article.CreatedAt
	articles = append(articles, article)
	nextID++
	okResponse(c, http.StatusCreated, "Article created", article)
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var articleData Article
	if err := c.ShouldBindJSON(&articleData); err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}
	if err := validateArticle(articleData); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	article, index := findArticleByID(id)
	if article == nil {
		errResponse(c, http.StatusNotFound, "Not found")
		return
	}

	articleData.ID = id
	articleData.CreatedAt = article.CreatedAt
	articleData.UpdatedAt = time.Now()
	articles[index] = articleData
	okResponse(c, http.StatusOK, "Article updated", articleData)
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	_, index := findArticleByID(id)
	if index == -1 {
		errResponse(c, http.StatusNotFound, "Not found")
		return
	}

	articles = slices.Delete(articles, index, index + 1)
	okResponse(c, http.StatusOK, "Article deleted", nil)
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	if c.GetString("role") != "admin" {
		errResponse(c, http.StatusForbidden, "Admin role required")
		return
	}

	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0,
		"uptime":         time.Since(time.Now().Add(-24 * time.Hour)).String(),
	}
	okResponse(c, http.StatusOK, "Statistics", stats)
}

// ----------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	for i, a := range(articles) {
		if a.ID == id {
			return &a, i
		}
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	if strings.TrimSpace(article.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(article.Content) == "" {
		return fmt.Errorf("content is required")
	}
	if strings.TrimSpace(article.Author) == "" {
		return fmt.Errorf("author is required")
	}
	return nil
}

func okResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		RequestID: c.GetString("request_id"),
	})
}

func errResponse(c *gin.Context, status int, msg string) {
	c.JSON(status, APIResponse{
		Success:   false,
		Error:     msg,
		RequestID: c.GetString("request_id"),
	})
}
