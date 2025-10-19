package main

import (
    "fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"strconv"

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
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3

func main() {
	router := gin.New()

	// 全域中介軟體 
	router.Use(ErrorHandlerMiddleware(), RequestIDMiddleware(), LoggingMiddleware(), CORSMiddleware())

	// --- 2. 劃分「公開」和「受保護」的區域 ---

	// 公開 API 的群組 (Public Routes)
	// 任何人都可以訪問
	publicRoutes := router.Group("/api")
	{
		publicRoutes.GET("/ping", ping)
		publicRoutes.GET("/articles", getArticles)
		publicRoutes.GET("/articles/:id", getArticle)
	}

	// 受保護 API 的群組 (Protected Routes)
	// 建立一個「獨立」的群組來管理所有需要保護的路由
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(AuthMiddleware())      // 先驗票
	protectedRoutes.Use(RateLimitMiddleware()) // 再做人流管制
	{
		// 所有在這個群組下定義的路由，都會自動應用上面那兩道安檢
		protectedRoutes.POST("/articles", ContentTypeMiddleware(), createArticle)
		protectedRoutes.PUT("/articles/:id", ContentTypeMiddleware(), updateArticle)
		protectedRoutes.DELETE("/articles/:id", deleteArticle)

		// 巢狀的 Admin 群組會「繼承」protectedRoutes 的所有中介軟體
		adminOnly := protectedRoutes.Group("/admin")
		{
			adminOnly.GET("/stats", getStats)
		}
	}

	// 3. 啟動伺服器
	fmt.Println("Server is running on port 8080")
	router.Run(":8080")
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request (Done)
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Generate UUID for request ID
		// Use github.com/google/uuid package
		// Store in context as "request_id"
		// Add to response header as "X-Request-ID"
		requestID := uuid.NewString()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
        
		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information (Done)
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Capture start time
        startTime := time.Now()
        
		c.Next()
		// TODO: Calculate duration and log request
		duration := time.Since(startTime)
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		requestID, _ := c.Get("request_id")
		log.Printf("[%s] %s %s %d %v %s \"%s\"",
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

// AuthMiddleware validates API keys for protected routes (Done)
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"
	validApiKeys := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
	    
		// TODO: Get API key from X-API-Key header
		apiKey := c.GetHeader("X-API-Key")
		
		// TODO: Validate API key
		role, found := validApiKeys[apiKey]
		if !found {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "invalid API Key",
			})
			return
		}
		
		// TODO: Set user role in context
		c.Set("user_role", role)

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests (Done)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Set CORS headers
		// Allow origins: http://localhost:3000, https://myblog.com
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow headers: Content-Type, X-API-Key, X-Request-ID
        c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID, Authorization")
		// TODO: Handle preflight OPTIONS requests
        if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RateLimitMiddleware 實作每個 IP 的請求速率限制，並加上回饋標頭
func RateLimitMiddleware() gin.HandlerFunc {
	// 每個 IP 對應一個速率限制器
	var visitors = make(map[string]*rate.Limiter)
	var mu sync.Mutex

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := visitors[ip]
		if !exists {
			// 需求：每分鐘最多 100 次請求
			limit := rate.Every(time.Minute / 100)
			// 測試案例的邏輯期望初始突發容量 (burst) 就是限制的總數
			burst := 100
			limiter = rate.NewLimiter(limit, burst)
			visitors[ip] = limiter
		}
		mu.Unlock()

		// --- 設定回饋標頭 ---
		// X-RateLimit-Limit: 固定的總限制數
		c.Header("X-RateLimit-Limit", "100")
		// X-RateLimit-Reset: 簡單起見，我們設定為 60 秒後重置
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Unix()+60, 10))

		// 檢查是否還有可用的令牌
		if !limiter.Allow() {
			// 如果被阻止，剩餘次數為 0
			c.Header("X-RateLimit-Remaining", "0")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, APIResponse{
				Success: false,
				Error:   "請求頻率過高，請稍後再試",
			})
			return
		}

		// 如果請求被允許，計算剩餘次數並設定標頭
		// limiter.Tokens() 回傳的是 float64，我們需要轉換成整數
		remaining := int(limiter.Tokens())
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests (Done)
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check content type for POST/PUT requests
		// Must be application/json
		// Return 415 if invalid content type
        if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			// 簡單檢查 Content-Type 是否為 application/json
			if !strings.HasPrefix(contentType, "application/json") {
				// 415 erorr = StatusUnsupportedMediaType
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, APIResponse{
					Success: false,
					Error:   "Unsupported content type. Please use application/json",
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
		requestID, _ := c.Get("request_id")
		
		c.JSON(http.StatusInternalServerError, APIResponse{
		    Success: false,
		    Error: "Internal server error",
		    Message: fmt.Sprintf("%v", recovered),
		    RequestID: fmt.Sprintf("%v", requestID),
		    
		})
	})
}

// TODO: Implement route handlers

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "pong",
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// 1. 從 URL 查詢參數中取得 "page" 和 "limit"。
	// c.DefaultQuery 非常好用，如果使用者沒有提供該參數，它會使用我們給的預設值。
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// 2. 將字串轉換成數字，並做基本驗證。
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // 如果格式不對或頁數小於1，就預設為第 1 頁
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // 如果格式不對或數量小於1，就預設為每頁 10 筆
	}

	// 3. 根據頁數和每頁數量，計算切片的起始和結束索引。
	//    這就是「圖書館員」的計算過程。
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	// 4. 處理邊界情況，防止程式崩潰 (panic)。
	// 如果請求的頁數太大，導致起始索引超過了文章總數...
	if startIndex >= len(articles) {
		// ...就直接回傳一個「空的」文章列表。
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Data:    make([]Article, 0), // 使用 make 確保回傳的是 `[]` 而不是 `null`
		})
		return
	}

	// 如果結束索引超過了文章總數（例如在最後一頁）...
	if endIndex > len(articles) {
		// ...就把結束索引設定為文章總數，避免 slice out of bounds 錯誤。
		endIndex = len(articles)
	}

	// 5. 使用計算好的索引，從完整的 articles 列表中「切」出我們需要的那一頁資料。
	paginatedArticles := articles[startIndex:endIndex]
    
    requestID, _ := c.Get("request_id")
	// 6. 將分頁後的資料，用標準的 APIResponse 格式回傳。
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    paginatedArticles,
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// getArticle handles GET /articles/:id - get article by ID (Done)
func getArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find article by ID
	// TODO: Return 404 if not found
	
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "無效的文章 ID 格式"})
		return
	}
	
	article, _ := findArticleByID(id)
	if article == nil {
	    c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "找不到該文章"})
		return
	}
	
	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
	    Success: true, 
	    Data: article,
	    RequestID: fmt.Sprintf("%v", requestID),
	})
}

// createArticle handles POST /articles - create new article (protected) (Done)
func createArticle(c *gin.Context) {
	// TODO: Parse JSON request body
	// TODO: Validate required fields
	// TODO: Add article to storage
	// TODO: Return created article
	
	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	
	if err := validateArticle(newArticle); err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	
	newArticle.ID = nextID
	nextID++
	newArticle.CreatedAt = time.Now()
	newArticle.UpdatedAt = time.Now()
	
	articles = append(articles, newArticle)
	
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    newArticle,
		Message: "Successfully created new article",
	})
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Parse JSON request body
	// TODO: Find and update article
	// TODO: Return updated article
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "無效的文章 ID 格式"})
		return
	}
	
	oriArticle, index := findArticleByID(id)
	if oriArticle == nil {
		c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "找不到要更新的文章"})
		return
	}

	var updatedArticleData Article
	if err := c.ShouldBindJSON(&updatedArticleData); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	
	if err := validateArticle(updatedArticleData); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	oriArticle.Title = updatedArticleData.Title
	oriArticle.Content = updatedArticleData.Content
	oriArticle.Author = updatedArticleData.Author
	oriArticle.UpdatedAt = time.Now()

	articles[index] = *oriArticle

	c.JSON(http.StatusOK, APIResponse{Success: true, Data: oriArticle, Message: "文章更新成功"})
}

// deleteArticle handles DELETE /articles/:id - delete article (protected) (Done)
func deleteArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find and remove article
	// TODO: Return success message
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
	    c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "無效的文章 ID 格式"})
		return
	}
	
	_, index := findArticleByID(id)
	if id == -1 {
	    c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "找不到要刪除的文章"})
		return
	}
	
	articles = append(articles[:index], articles[index+1:]...)
	c.JSON(http.StatusOK, APIResponse{Success: true, Message: "文章刪除成功"})
}

// getStats handles GET /admin/stats - get API usage statistics (admin only) (Done)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	role, _ := c.Get("user_role")
	
	if role != "admin" {
	    c.JSON(http.StatusForbidden, APIResponse{
			Success: false,
			Error:   "權限不足，僅限管理員訪問",
		})
		return
	}
	
	// TODO: Return mock statistics
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}

	// TODO: Return stats in standard format
	c.JSON(http.StatusOK, APIResponse{
	    Success: true,
	    Data: stats,
	})
}

// Helper functions

// findArticleByID finds an article by ID (Done)
func findArticleByID(id int) (*Article, int) {
	// TODO: Implement article lookup
	// Return article pointer and index, or nil and -1 if not found
	for i, article := range articles {
	    if article.ID == id {
	        return &article, i
	    }
	}
	return nil, -1
}

// validateArticle validates article data (Done)
func validateArticle(article Article) error {
	// TODO: Implement validation
	// Check required fields: Title, Content, Author
	if article.Title == "" {
	    return fmt.Errorf("Title is necessary")
	}
	if article.Content == "" {
	    return fmt.Errorf("Content is necessary")
	}
	if article.Author == "" {
	    return fmt.Errorf("Author is necessary")
	}
	return nil
}
