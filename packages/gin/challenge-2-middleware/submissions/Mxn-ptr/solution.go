package main

import (
	"errors"
	"log"
	"sort"
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
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	Author    string    `json:"author" binding:"required"`
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
	router := gin.New()

	router.Use(ErrorHandlerMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(LoggingMiddleware())
	router.Use(CORSMiddleware())
	router.Use(RateLimitMiddleware())
	router.Use(ContentTypeMiddleware())

	// TODO: Setup custom middleware in correct order
	// 1. ErrorHandlerMiddleware (first to catch panics)
	// 2. RequestIDMiddleware
	// 3. LoggingMiddleware
	// 4. CORSMiddleware
	// 5. RateLimitMiddleware
	// 6. ContentTypeMiddleware
	publicRoutes := router.Group("")
	{
		publicRoutes.GET("/ping", ping)
		publicRoutes.GET("/articles", getArticles)
		publicRoutes.GET("/article/:id", getArticle)
	}

	privateRoutes := router.Group("")
	privateRoutes.Use(AuthMiddleware())
	{
		privateRoutes.POST("/articles", createArticle)
		privateRoutes.PUT("/article/:id", updateArticle)
		privateRoutes.DELETE("/article/:id", deleteArticle)
	}

	// TODO: Start server on port 8080
	router.Run("localhost:8080")
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)

		c.Header("X-Request-ID", requestID)

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
			"request_id": c.GetString("request_id"),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration.Microseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		log.Printf("[%s] %s %s %v %v %s %s", entry["request_id"], entry["method"], entry["path"], entry["status"], entry["duration"], entry["ip"], entry["user_agent"])
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {

	apiKeys := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {

		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "API key required"})
			return
		}
		role := apiKeys[apiKey]

		if role == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
			return
		}

		c.Set("user_role", role)

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"https://myblog.com":    true,
		}

		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	const maxRequests = 100
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Cleanup stale entries every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(maxRequests)), maxRequests),
			}
		}
		clients[ip].lastSeen = time.Now()
		limiter := clients[ip].limiter
		mu.Unlock()

		allowed := limiter.Allow()
		remaining := int(limiter.Tokens())
		if remaining < 0 {
			remaining = 0
		}
		resetTime := time.Now().Add(time.Minute).Unix()

		c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		if !allowed {
			c.AbortWithStatusJSON(429, APIResponse{
				Success:   false,
				Error:     "Rate limit exceeded. Try again later.",
				RequestID: c.GetString("request_id"),
			})
			return
		}

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")

			if !strings.HasPrefix(contentType, "application/json") {
				c.AbortWithStatusJSON(415, APIResponse{
					Success:   false,
					Error:     "Content-Type must be application/json",
					RequestID: c.GetString("request_id"),
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
		c.AbortWithStatusJSON(500, APIResponse{
			Success:   false,
			Error:     "Internal server error",
			Message:   recovered.(error).Error(),
			RequestID: c.GetString("request_id"),
		})
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	c.JSON(200, APIResponse{
		Success:   true,
		Message:   "pong",
		RequestID: c.GetString("request_id"),
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	requestId := c.GetString("request_id")

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.AbortWithStatusJSON(400, APIResponse{
			Success:   false,
			Error:     "Invalid format page",
			RequestID: requestId,
		})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.AbortWithStatusJSON(400, APIResponse{
			Success:   false,
			Error:     "Invalid format limit",
			RequestID: requestId,
		})
		return
	}
	limitInt = min(limitInt, 100)

	if pageInt < 1 {
		c.AbortWithStatusJSON(422, APIResponse{
			Success:   false,
			Error:     "Invalid value page",
			RequestID: requestId,
		})
		return
	}

	if limitInt < 1 {
		c.AbortWithStatusJSON(422, APIResponse{
			Success:   false,
			Error:     "Invalid value limit",
			RequestID: requestId,
		})
		return
	}
	articlesCopy := append([]Article(nil), articles...)
	sort.Slice(articlesCopy, func(i, j int) bool {
		return articlesCopy[i].ID > articlesCopy[j].ID
	})

	start := (pageInt - 1) * limitInt
	items := []Article{}
	if start < len(articlesCopy) {
		items = articlesCopy[start:min(start+limitInt, len(articlesCopy))]
	}

	c.JSON(200, APIResponse{
		Success:   true,
		Message:   "ArticleList",
		Data:      items,
		RequestID: requestId,
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	requestId := c.GetString("request_id")
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(400, APIResponse{
			Success:   false,
			Error:     "Invalid format id",
			RequestID: requestId,
		})
		return
	}

	if idInt < 1 {
		c.AbortWithStatusJSON(422, APIResponse{
			Success:   false,
			Error:     "Invalid value id",
			RequestID: requestId,
		})
		return
	}

	article, index := findArticleByID(idInt)
	if index == -1 {
		c.AbortWithStatusJSON(404, APIResponse{
			Success:   false,
			Error:     "Article not found",
			RequestID: requestId,
		})
		return
	}

	c.JSON(200, APIResponse{
		Success:   true,
		Message:   "Article found",
		Data:      article,
		RequestID: requestId,
	})
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	var article Article
	requestId := c.GetString("request_id")

	if c.GetString("user_role") != "admin" {
		c.AbortWithStatusJSON(403, APIResponse{
			Success:   false,
			Error:     "Forbidden",
			RequestID: requestId,
		})
		return
	}

	if err := c.ShouldBindJSON(&article); err != nil {
		c.AbortWithStatusJSON(400, APIResponse{
			Success:   false,
			Error:     "Invalid format article",
			RequestID: requestId,
		})
		return
	}

	if err := validateArticle(article); err != nil {
		c.AbortWithStatusJSON(422, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: requestId,
		})
		return
	}

	article.ID = nextID
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	nextID++

	articles = append(articles, article)

	c.JSON(201, APIResponse{
		Success:   true,
		Message:   "Article created",
		Data:      article,
		RequestID: requestId,
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	requestId := c.GetString("request_id")

	if c.GetString("user_role") != "admin" {
		c.AbortWithStatusJSON(403, APIResponse{
			Success:   false,
			Error:     "Forbidden",
			RequestID: requestId,
		})
		return
	}

	var bodyArticle Article

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, APIResponse{
			Success:   false,
			Error:     "Invalid format id",
			RequestID: requestId,
		})
		return
	}

	if idInt < 1 {
		c.AbortWithStatusJSON(422, APIResponse{
			Success:   false,
			Error:     "Invalid value id",
			RequestID: requestId,
		})
		return
	}

	if err := c.ShouldBindJSON(&bodyArticle); err != nil {
		c.AbortWithStatusJSON(400, APIResponse{
			Success:   false,
			Error:     "Invalid format article",
			RequestID: requestId,
		})
		return
	}

	if err := validateArticle(bodyArticle); err != nil {
		c.AbortWithStatusJSON(422, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: requestId,
		})
		return
	}

	article, index := findArticleByID(idInt)
	if index == -1 {
		c.AbortWithStatusJSON(404, APIResponse{
			Success:   false,
			Error:     "Article not found",
			RequestID: requestId,
		})
		return
	}

	article.Title = bodyArticle.Title
	article.Content = bodyArticle.Content
	article.Author = bodyArticle.Author
	article.UpdatedAt = time.Now()

	c.JSON(200, APIResponse{
		Success:   true,
		Message:   "Article updated",
		Data:      article,
		RequestID: requestId,
	})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	requestId := c.GetString("request_id")

	if c.GetString("user_role") != "admin" {
		c.AbortWithStatusJSON(403, APIResponse{
			Success:   false,
			Error:     "Forbidden",
			RequestID: requestId,
		})
		return
	}

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, APIResponse{
			Success:   false,
			Error:     "Invalid format id",
			RequestID: requestId,
		})
		return
	}

	if idInt < 1 {
		c.AbortWithStatusJSON(422, APIResponse{
			Success:   false,
			Error:     "Invalid value id",
			RequestID: requestId,
		})
		return
	}

	_, index := findArticleByID(idInt)
	if index == -1 {
		c.AbortWithStatusJSON(404, APIResponse{
			Success:   false,
			Error:     "Article not found",
			RequestID: requestId,
		})
		return
	}

	articles = append(articles[:index], articles[index+1:]...)

	c.JSON(200, APIResponse{
		Success:   true,
		Message:   "Article deleted",
		RequestID: requestId,
	})
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	requestId := c.GetString("request_id")

	if c.GetString("user_role") != "admin" {
		c.AbortWithStatusJSON(403, APIResponse{
			Success:   false,
			Error:     "Forbidden",
			RequestID: requestId,
		})
		return
	}
	// TODO: Check if user role is "admin"
	// TODO: Return mock statistics
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}

	c.JSON(200, APIResponse{
		Success:   true,
		Message:   "Article found",
		Data:      stats,
		RequestID: requestId,
	})

	// TODO: Return stats in standard format
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	for i, v := range articles {
		if v.ID == id {
			return &articles[i], i
		}
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	if len(article.Title) < 3 || len(article.Title) > 100 {
		return errors.New("Title must be between 3 and 100 characters.")
	} else if len(article.Content) < 3 || len(article.Content) > 100 {
		return errors.New("Content must be between 3 and 100 characters.")
	} else if len(article.Author) < 3 || len(article.Author) > 100 {
		return errors.New("Author must be between 3 and 100 characters.")
	}
	return nil
}
