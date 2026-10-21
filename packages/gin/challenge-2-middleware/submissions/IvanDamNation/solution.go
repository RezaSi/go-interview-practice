package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	Success   bool   `json:"success"`
	Data      any    `json:"data,omitempty"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// In-memory storage
var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3

func main() {
	r := gin.New()

	r.Use(RequestIDMiddleware())
	r.Use(ErrorHandlerMiddleware())
	r.Use(LoggingMiddleware())
	r.Use(CORSMiddleware())
	r.Use(RateLimitMiddleware())

	r.GET("/ping", ping)
	r.GET("/articles", getArticles)
	r.GET("/articles/:id", getArticle)

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	protected.Use(ContentTypeMiddleware())

	protected.POST("/articles", createArticle)
	protected.PUT("/articles/:id", updateArticle)
	protected.DELETE("/articles/:id", deleteArticle)

	protected.GET("/admin/stats", getStats)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.New().String()
		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)

		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		reqID := getRequestID(c)
		dur := time.Since(start)

		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		log.Printf("[%s] %s %s %d %s %s %s\n",
			reqID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			dur,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	roles := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		role, exists := roles[key]
		if !exists {
			reqID := getRequestID(c)
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Error:     "Invalid API key",
				RequestID: reqID,
			})
			return
		}

		if role != "admin" && strings.HasPrefix(c.Request.URL.Path, "/admin") {
			reqID := getRequestID(c)
			c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
				Error:     "low permission",
				RequestID: reqID,
			})
			return
		}

		c.Set("role", role)

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}

		c.Header("Access-Control-Allow-Origin",
			origin)
		c.Header("Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers",
			"Origin, Content-Type, Accept, Authorization, X-API-Key, X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}

		c.Next()
	}
}

type ClientRate struct {
	Count     int
	ResetTime time.Time
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	var mu sync.Mutex
	clients := make(map[string]*ClientRate)
	lastCleanup := time.Time{}

	const (
		maxLimit    = 100
		maxLimitStr = "100"
		resetTime   = 1 * time.Minute
	)

	return func(c *gin.Context) {
		ip := c.RemoteIP()
		now := time.Now()

		mu.Lock()
		if now.Sub(lastCleanup) >= resetTime {
			for clientIP, clientRate := range clients {
				if now.After(clientRate.ResetTime) {
					delete(clients, clientIP)
				}
			}
			lastCleanup = now
		}
		rate, exists := clients[ip]
		if !exists || now.After(rate.ResetTime) {
			rate = &ClientRate{
				Count:     1,
				ResetTime: now.Add(resetTime),
			}
			clients[ip] = rate
		} else {
			rate.Count++
		}
		remaining := maxLimit - rate.Count
		tillReset := rate.ResetTime.Unix()
		mu.Unlock()

		displayRemaining := max(remaining, 0)

		c.Header("X-RateLimit-Limit", maxLimitStr)
		c.Header("X-RateLimit-Remaining", strconv.Itoa(displayRemaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(tillReset, 10))

		if remaining < 0 {
			reqID := getRequestID(c)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, APIResponse{
				Error:     "Request limit exceeded",
				RequestID: reqID,
			})
			return
		}

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodPost || method == http.MethodPut {
			ct := c.GetHeader("Content-Type")
			if !strings.HasPrefix(ct, "application/json") {
				reqID := getRequestID(c)
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, APIResponse{
					Error:     "Unsupported Media Type",
					Message:   "Content-Type must be application/json",
					RequestID: reqID,
				})
				return
			}
		}

		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		reqID := getRequestID(c)

		resp := APIResponse{
			Error:     "Internal server error",
			Message:   fmt.Sprintf("%v", recovered),
			RequestID: reqID,
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
	})
}

// Usecase handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	reqID := getRequestID(c)
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "Pong",
		RequestID: reqID,
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	reqID := getRequestID(c)
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles,
		Message:   "Article list",
		RequestID: reqID,
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	reqID := getRequestID(c)
	id := validateID(c)
	if id == -1 {
		return
	}

	ar, found := findArticleByID(id)
	if found == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Error:     "Article not found",
			RequestID: reqID,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "Article found",
		Data:      ar,
		RequestID: reqID,
	})
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	var input Article
	var err error

	reqID := getRequestID(c)
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:     "Invalid payload",
			RequestID: reqID,
		})
		return
	}

	if err = validateArticle(input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:     err.Error(),
			RequestID: reqID,
		})
		return
	}

	now := time.Now()
	input.ID = nextID
	input.CreatedAt = now
	input.UpdatedAt = now

	articles = append(articles, input)
	nextID++

	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Message:   "Article created successfully",
		Data:      input,
		RequestID: reqID,
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	reqID := getRequestID(c)
	id := validateID(c)
	if id == -1 {
		return
	}

	art, found := findArticleByID(id)
	if found == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Error:     "Article not found",
			RequestID: reqID,
		})
		return
	}

	var input Article
	var err error
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:     "Invalid payload",
			RequestID: reqID,
		})
		return
	}

	if err = validateArticle(input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:     err.Error(),
			RequestID: reqID,
		})
		return
	}

	art.Author = input.Author
	art.Content = input.Content
	art.Title = input.Title
	art.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "Article updated successfully",
		Data:      art,
		RequestID: reqID,
	})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	id := validateID(c)
	if id == -1 {
		return
	}

	reqID := getRequestID(c)
	_, index := findArticleByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, APIResponse{
			Error:     "Article not found",
			RequestID: reqID,
		})
		return
	}
	articles = append(articles[:index], articles[index+1:]...)

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "Article deletion successful",
		RequestID: reqID,
	})
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	stats := map[string]any{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}

	reqID := getRequestID(c)
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "API statistics",
		Data:      stats,
		RequestID: reqID,
	})
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

// getRequestID extract request ID
func getRequestID(c *gin.Context) string {
	if val, exists := c.Get("request_id"); exists {
		return val.(string)
	}
	return "NO_REQUEST_ID"
}

// validateArticle validates article data
func validateArticle(article Article) error {
	if article.Title == "" {
		return fmt.Errorf("Title is empty")
	}

	if article.Author == "" {
		return fmt.Errorf("Author is empty")
	}

	if article.Content == "" {
		return fmt.Errorf("Article content is empty")
	}

	return nil
}

func validateID(c *gin.Context) int {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:     "Invalid type id",
			RequestID: getRequestID(c),
		})
		return -1
	}
	return id
}
