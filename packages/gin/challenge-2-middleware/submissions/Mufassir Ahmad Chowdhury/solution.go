package main

import (
	"time"
    "net/http"
    "log"
    "strings"
    "strconv"
    "errors"
    "slices"
    
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

	// TODO: Setup custom middleware in correct order
	// 1. ErrorHandlerMiddleware (first to catch panics)
	// 2. RequestIDMiddleware
	// 3. LoggingMiddleware
	// 4. CORSMiddleware
	// 5. RateLimitMiddleware
	// 6. ContentTypeMiddleware
// 	router.Use(ErrorHandlerMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(LoggingMiddleware())
	router.Use(CORSMiddleware())
// 	router.Use(RateLimitMiddleware())
	router.Use(ContentTypeMiddleware())
    
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

	router.Run(":8080")
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
        request_id := uuid.New().String()
        c.Set("request_id", request_id)
        c.Header("X-Request-ID", request_id)
        
		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
        c.Next()
        
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
        duration := time.Since(start)
        log.Printf("[%s] %s %s %v", 
            c.Request.Method, 
            c.Request.URL.Path, 
            duration,
            c.ClientIP())


		
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	apiKeyMap := make(map[string]string)
	apiKeyMap["admin-key-123"] = "admin"
	apiKeyMap["user-key-456"] = "user"

	return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" || apiKeyMap[apiKey] == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse {
                Success: false,
                Error: "Not authorized",
            })
            return
        }
        c.Set("user_role", apiKeyMap[apiKey])
        
		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
        
        // Define allowed origins
        allowedOrigins := map[string]bool{
            "http://localhost:3000":  true,
            "https://myapp.com":      true,
        }
        
        if allowedOrigins[origin] {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-KEY, X-Request-ID")

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
		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
            contentType := c.GetHeader("Content-Type")
            
            if !strings.HasPrefix(contentType, "application/json") {
                c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, APIResponse {
                    Success: false,
                    Error: "Content-Type must be application/json",
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
		// TODO: Handle panics gracefully
		// Return consistent error response format
		// Include request ID in response
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
    value, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse {
	    Success: true,
	    Message: "pong",
	    RequestID: value.(string),
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
    value, _ := c.Get("request_id")
	// TODO: Implement pagination (optional)
	c.JSON(http.StatusOK, APIResponse {
	    Success: true,
	    Data: articles,
	    RequestID: value.(string),
	})
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
    value, _ := c.Get("request_id")
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse {
	        Success: false,
	        Error: "Invalid ID", 
	        RequestID: value.(string),
	    })
	    return
	}
	article, _ := findArticleByID(id)
	if article != nil {
    	c.JSON(http.StatusOK, APIResponse {
            Success: true,
            Data: *article,
            RequestID: value.(string),
        })
	    return
	}
	c.JSON(http.StatusNotFound, APIResponse {
	    Success: false,
	    Error: "No article found",
	    RequestID: value.(string),
	})
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	value, _ := c.Get("request_id")
	var article Article
	
	if err := c.ShouldBindJSON(&article); err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse {
	        Success: false,
	        Error: "Bad format",
	        RequestID: value.(string),
	    })
	    return
	}
	err := validateArticle(article)
	if err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse {
	        Success: false,
	        Error: err.Error(),
	        RequestID: value.(string),
	    })
	    return
	}
	articles = append(articles, article)
	c.JSON(http.StatusCreated, APIResponse {
	    Success: true,
	    Data: article,
	    RequestID: value.(string),
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	value, _ := c.Get("request_id")
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse {
	        Success: false, 
	        Error: "Invalid ID",
	        RequestID: value.(string),
	    })
	    return
	}
	var article Article
	
	if err := c.ShouldBindJSON(&article); err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse {
	        Success: false,
	        Error: "Bad format",
	        RequestID: value.(string),
	    })
	    return
	}
	err = validateArticle(article)
	if err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse {
	        Success: false,
	        Error: err.Error(),
	        RequestID: value.(string),
	    })
	    return
	}
	oldArticle, _ := findArticleByID(id)
	if oldArticle != nil {
	    oldArticle = &article
    	c.JSON(http.StatusOK, APIResponse {
            Success: true,
            Data: article,
            RequestID: value.(string),
        })
	    return
	}
	c.JSON(http.StatusNotFound, APIResponse {
	    Success: false,
	    Error: "No article found",
	    RequestID: value.(string),
	})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	value, _ := c.Get("request_id")
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse {
	        Success: false, 
	        Error: "Invalid ID",
	        RequestID: value.(string),
	    })
	    return
	}
	_, index := findArticleByID(id)
	if index != -1 {
	    articles = slices.Delete(articles, index, index+1)
	    
    	c.JSON(http.StatusOK, APIResponse {
            Success: true,
            Message: "Successfully deleted!",
            RequestID: value.(string),
        })
	    return
	}
	c.JSON(http.StatusNotFound, APIResponse {
	    Success: false,
	    Error: "No article found",
	    RequestID: value.(string),
	})
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
    value, _ := c.Get("request_id")
    role, _ := c.Get("user_role")
    
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}
    if role != "admin" {
        c.JSON(http.StatusForbidden, APIResponse {
            Success: false,
            Error: "Admin Only",
            RequestID: value.(string),
        })
        return
    }
    
	c.JSON(http.StatusOK, APIResponse {
	    Success: true,
	    Data: stats,
	    RequestID: value.(string),
	})
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// TODO: Implement article lookup
	// Return article pointer and index, or nil and -1 if not found
	for index, article := range articles {
	    if article.ID == id {
	        return &article, index
	    }
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	if article.Title == "" {
	    return errors.New("No Title Provided")
	}
	if article.Content == "" {
	    return errors.New("No Content Provided")
	}
	if article.Author == "" {
	    return errors.New("No Author Provided")
	}
	return nil
}
