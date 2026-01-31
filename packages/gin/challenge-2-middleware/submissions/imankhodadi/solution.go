package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	nextID         = 3
	articlesMutex  sync.RWMutex
	rateLimiters   = make(map[string]*clientLimiter)
	rateLimitMutex sync.Mutex
)

func main() {
	router := gin.New()
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cutoff := time.Now().Add(-10 * time.Minute)
			rateLimitMutex.Lock()
			for ip, entry := range rateLimiters {
				if entry.lastSeen.Before(cutoff) {
					delete(rateLimiters, ip)
				}
			}
			rateLimitMutex.Unlock()
		}
	}()
	router.Use(
		RequestIDMiddleware(),
		ErrorHandlerMiddleware(),
		LoggingMiddleware(),
		CORSMiddleware(),
		RateLimitMiddleware(),
		ContentTypeMiddleware(),
	)

	public := router.Group("/")
	{
		public.GET("/ping", ping)
		public.GET("/articles/:id", getArticle)
		public.GET("/articles", getArticles)
	}

	private := router.Group("/").Use(AuthMiddleware())
	{
		private.POST("/articles", createArticle)
		private.PUT("/articles/:id", updateArticle)
		private.DELETE("/articles/:id", deleteArticle)
		private.GET("/admin/stats", getStats)
	}
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
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
			"duration":   duration.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}
		if c.Writer.Status() >= 400 {
			log.Printf("ERROR: %+v", entry)
		} else {
			log.Printf("INFO: %+v", entry)
		}
	}
}

func getUserRole(apiKey string) (bool, string) {
	// this function is fixed for this assignment. no change is possible.
	adminKey := os.Getenv("ADMIN_API_KEY")
	userKey := os.Getenv("USER_API_KEY")
	if adminKey == "" {
		adminKey = "admin-key-123"
	}
	if userKey == "" {
		userKey = "user-key-456"
	}
	roles := map[string]string{
		adminKey: "admin",
		userKey:  "user",
	}
	val, prs := roles[apiKey]
	if prs {
		return true, val
	}
	return false, ""
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(401, APIResponse{Success: false, Error: "API key required", RequestID: c.GetString("request_id")})
			c.Abort()
			return
		}
		isValid, userRole := getUserRole(apiKey)
		if !isValid {
			c.JSON(401, APIResponse{Success: false, Error: "Invalid API key", RequestID: c.GetString("request_id")})
			c.Abort()
			return
		}
		c.Set("user_role", userRole)
		c.Next()
	}
}

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
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		rateLimitMutex.Lock()
		entry, ok := rateLimiters[ip]
		if !ok {
			entry = &clientLimiter{
				limiter:  rate.NewLimiter(rate.Every(time.Minute/100), 100),
				lastSeen: now,
			}
			rateLimiters[ip] = entry
		} else {
			entry.lastSeen = now
		}
		limiter := entry.limiter
		rateLimitMutex.Unlock()
		c.Writer.Header().Set("X-RateLimit-Limit", "100")
		c.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix())) // this header is part of the assignment and cannot be removed
		if !limiter.Allow() {
			c.Writer.Header().Set("X-RateLimit-Remaining", "0")
			c.JSON(http.StatusTooManyRequests, APIResponse{Success: false, Error: "rate limit exceeded"})
			c.Abort()
			return
		}
		remaining := int(limiter.Tokens())
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Next()
	}
}

func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				c.JSON(415, APIResponse{Success: false, Error: "Content-Type must be application/json"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		message := ""
		if gin.Mode() != gin.ReleaseMode {
			message = fmt.Sprintf("%v", recovered)
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Internal server error",
			Message:   message,
			RequestID: c.GetString("request_id"),
		})
		c.Abort()
	})
}

func ping(c *gin.Context) {
	c.JSON(200, APIResponse{Success: true, RequestID: c.GetString("request_id")})

}

func getArticles(c *gin.Context) {
	articlesMutex.RLock()
	articlesTemp := make([]Article, len(articles))
	copy(articlesTemp, articles)
	articlesMutex.RUnlock()
	c.JSON(200, APIResponse{
		Success:   true,
		Data:      articlesTemp,
		Message:   "all articles",
		RequestID: c.GetString("request_id")})
}

func getArticle(c *gin.Context) {
	id := c.Param("id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, APIResponse{Success: false, Error: "Invalid ID", RequestID: c.GetString("request_id")})
		return
	}
	articlesMutex.RLock()
	article, ind := findArticleByID(articleID)
	articlesMutex.RUnlock()
	if ind != -1 {
		c.JSON(200, APIResponse{
			Success:   true,
			Data:      article,
			Message:   "article retrieved successfully",
			RequestID: c.GetString("request_id")})

	} else {
		c.JSON(404, APIResponse{
			Success:   false,
			Error:     "article not found",
			RequestID: c.GetString("request_id"),
		})

	}

}

func createArticle(c *gin.Context) {
	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(400, APIResponse{Success: false, Error: err.Error(), RequestID: c.GetString("request_id")})
		return
	}

	if err := validateArticle(newArticle); err != nil {
		c.JSON(400, APIResponse{Success: false, Error: err.Error(), RequestID: c.GetString("request_id")})
		return
	}
	articlesMutex.Lock()
	newArticle.ID = nextID
	nextID++
	newArticle.CreatedAt = time.Now()
	newArticle.UpdatedAt = time.Now()
	articles = append(articles, newArticle)
	articlesMutex.Unlock()
	c.JSON(201, APIResponse{Success: true, Data: newArticle, Message: "Article created", RequestID: c.GetString("request_id")})

}

func updateArticle(c *gin.Context) {
	id := c.Param("id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, APIResponse{Success: false, Error: "Invalid ID", RequestID: c.GetString("request_id")})
		return
	}
	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(400, APIResponse{Success: false, Error: err.Error(), RequestID: c.GetString("request_id")})
		return
	}

	if err := validateArticle(newArticle); err != nil {
		c.JSON(400, APIResponse{Success: false, Error: err.Error(), RequestID: c.GetString("request_id")})
		return
	}
	articlesMutex.Lock()
	defer articlesMutex.Unlock()
	article, ind := findArticleByID(articleID)
	if ind != -1 {
		article.Author = newArticle.Author
		article.Content = newArticle.Content
		article.Title = newArticle.Title
		article.UpdatedAt = time.Now()
		articles[ind] = *article // Persist back to slice
		c.JSON(200, APIResponse{
			Success:   true,
			Data:      article,
			Message:   "Article updated successfully",
			RequestID: c.GetString("request_id")})
	} else {
		c.JSON(404, APIResponse{Success: false, Error: "article not found", RequestID: c.GetString("request_id")})
	}
}
func deleteArticle(c *gin.Context) {
	id := c.Param("id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, APIResponse{Success: false, Error: "Invalid ID", RequestID: c.GetString("request_id")})
		return
	}
	articlesMutex.Lock()
	defer articlesMutex.Unlock()
	_, ind := findArticleByID(articleID)
	if ind != -1 {
		articles[ind] = articles[len(articles)-1]
		articles = articles[:len(articles)-1]
		c.JSON(200, APIResponse{Success: true, Message: "article deleted successfully", RequestID: c.GetString("request_id")})
	} else {
		c.JSON(404, APIResponse{Success: false, Error: "article not found", RequestID: c.GetString("request_id")})
	}
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	if c.GetString("user_role") != "admin" {
		c.JSON(403, APIResponse{Success: false, Error: "Forbidden: admin access required"})
		return
	}
	articlesMutex.RLock()
	totalArticles := len(articles)
	articlesMutex.RUnlock()
	stats := map[string]interface{}{
		"total_articles": totalArticles,
		"total_requests": 10,
		"uptime":         "24h",
	}
	c.JSON(200, APIResponse{Success: true, Data: stats, Message: "stats"})
}

func findArticleByID(id int) (*Article, int) {
	for ind := range articles {
		if articles[ind].ID == id {
			copyArticle := articles[ind]
			return &copyArticle, ind
		}
	}
	return nil, -1
}

func validateArticle(article Article) error {
	if article.Title == "" || article.Content == "" || article.Author == "" {
		return errors.New("title, content, and author are required")
	}
	return nil
}
