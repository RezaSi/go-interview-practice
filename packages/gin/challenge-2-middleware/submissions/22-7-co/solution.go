package main

import (
	"encoding/json"
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
var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3

// 用于保护 articles 切片的并发访问
var articlesMutex sync.RWMutex

func main() {
	// 📌 关键：使用 gin.New() 创建路由器，不带默认中间件
	// gin.Default() 会自动添加 Logger 和 Recovery 中间件
	// gin.New() 让我们可以完全控制中间件的添加顺序
	r := gin.New()

	// 📌 中间件执行顺序很重要！
	// 中间件按照添加的顺序执行，像洋葱模型：
	// Request -> Middleware1 -> Middleware2 -> Handler -> Middleware2 -> Middleware1 -> Response

	// 1. ErrorHandlerMiddleware (最外层，捕获所有 panic)
	r.Use(ErrorHandlerMiddleware())

	// 2. RequestIDMiddleware (为每个请求生成唯一ID)
	r.Use(RequestIDMiddleware())

	// 3. LoggingMiddleware (记录请求日志)
	r.Use(LoggingMiddleware())

	// 4. CORSMiddleware (处理跨域请求)
	r.Use(CORSMiddleware())

	// 5. RateLimitMiddleware (限制请求频率)
	r.Use(RateLimitMiddleware())

	// 6. ContentTypeMiddleware (验证内容类型)
	r.Use(ContentTypeMiddleware())
	// 7. Sanitize500Middleware (兜底，必须放最后)
	r.Use(Sanitize500Middleware())

	// 📌 路由分组：将相关的路由组织在一起
	// Public routes (公开路由，不需要认证)
	public := r.Group("/")
	{
		public.GET("/ping", ping)               // 健康检查
		public.GET("/articles", getArticles)    // 获取所有文章
		public.GET("/articles/:id", getArticle) // 获取单篇文章
	}

	// Protected routes (受保护路由，需要 API Key 认证)
	protected := r.Group("/")
	protected.Use(AuthMiddleware()) // 只对这个组应用认证中间件
	{
		protected.POST("/articles", createArticle)       // 创建文章
		protected.PUT("/articles/:id", updateArticle)    // 更新文章
		protected.DELETE("/articles/:id", deleteArticle) // 删除文章
		protected.GET("/admin/stats", getStats)          // 管理员统计信息
	}

	// 启动服务器
	log.Println("🚀 Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// ============================================================================
// 中间件函数
// ============================================================================

// RequestIDMiddleware 为每个请求生成唯一的 ID
// 📌 用途：追踪请求，方便调试和日志关联
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成 UUID 作为请求 ID
		requestID := uuid.New().String()

		// 存储到 Gin Context 中，后续的处理器可以访问
		// c.Set() 用于在请求的生命周期内存储数据
		c.Set("request_id", requestID)

		// 添加到响应头，客户端可以看到
		c.Header("X-Request-ID", requestID)

		// 📌 关键：调用 c.Next() 继续执行下一个中间件/处理器
		c.Next()
	}
}

// LoggingMiddleware 记录所有请求的详细信息
// 📌 用途：监控 API 性能，调试问题
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 从 context 获取 request_id
		requestID, _ := c.Get("request_id")

		// 执行请求（调用后续的中间件和处理器）
		c.Next()

		// 📌 c.Next() 之后的代码在请求处理完成后执行
		// 这时我们可以获取响应状态码等信息

		// 计算请求处理时间
		duration := time.Since(startTime)

		// 格式化日志输出
		log.Printf("[%s] %s %s | Status: %d | Duration: %v | IP: %s | UserAgent: %s",
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

// AuthMiddleware 验证 API Key 并设置用户角色
// 📌 用途：保护敏感接口，实现权限控制
func AuthMiddleware() gin.HandlerFunc {
	// 定义有效的 API Key 和对应的角色
	// 实际项目中应该从数据库或配置文件读取
	validAPIKeys := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		// 从请求头获取 API Key
		apiKey := c.GetHeader("X-API-Key")

		// 检查 API Key 是否为空
		if apiKey == "" {
			requestID, _ := c.Get("request_id")
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success:   false,
				Error:     "API Key is required",
				RequestID: fmt.Sprintf("%v", requestID),
			})
			// 📌 关键：调用 c.Abort() 停止后续处理器的执行
			c.Abort()
			return
		}

		// 验证 API Key 是否有效
		role, exists := validAPIKeys[apiKey]
		if !exists {
			requestID, _ := c.Get("request_id")
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success:   false,
				Error:     "Invalid API Key",
				RequestID: fmt.Sprintf("%v", requestID),
			})
			c.Abort()
			return
		}

		// 将用户角色存储到 context 中
		c.Set("user_role", role)

		// 继续执行后续处理器
		c.Next()
	}
}

// CORSMiddleware 处理跨域资源共享 (CORS)
// 📌 用途：允许浏览器从不同域名访问 API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源（Origin）
		origin := c.GetHeader("Origin")
		allowedOrigins := []string{
			"http://localhost:3000",
			"https://myblog.com",
		}

		// 检查请求的 Origin 是否在允许列表中
		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			// 设置 CORS 响应头
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 允许的 HTTP 方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")

		// 是否允许携带凭证（cookies）
		c.Header("Access-Control-Allow-Credentials", "true")

		// 预检请求的缓存时间（秒）
		c.Header("Access-Control-Max-Age", "86400")

		// 📌 处理 OPTIONS 预检请求
		// 浏览器在发送跨域请求前，会先发送 OPTIONS 请求询问是否允许
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 实现基于 IP 的速率限制
// 📌 用途：防止 API 被滥用，保护服务器资源
func RateLimitMiddleware() gin.HandlerFunc {
	// 使用 map 存储每个 IP 的限流器
	// key: IP 地址, value: rate.Limiter
	limiters := make(map[string]*rate.Limiter)
	var mu sync.Mutex // 保护 map 的并发访问

	// 限制：每分钟 100 个请求
	// rate.Every(time.Minute) / 100 = 每 0.6 秒允许一个请求
	rateLimit := rate.Every(time.Minute / 100)
	burst := 10 // 突发容量：允许短时间内最多 10 个请求

	return func(c *gin.Context) {
		// 获取客户端 IP
		ip := c.ClientIP()

		// 获取或创建该 IP 的限流器
		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			// 为新 IP 创建限流器
			limiter = rate.NewLimiter(rateLimit, burst)
			limiters[ip] = limiter
		}
		mu.Unlock()

		// 检查是否允许请求
		if !limiter.Allow() {
			// 超过速率限制
			requestID, _ := c.Get("request_id")

			// 设置速率限制响应头
			c.Header("X-RateLimit-Limit", "100")
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

			c.JSON(http.StatusTooManyRequests, APIResponse{
				Success:   false,
				Error:     "Rate limit exceeded. Try again later.",
				RequestID: fmt.Sprintf("%v", requestID),
			})
			c.Abort()
			return
		}

		// 计算剩余令牌数（估算）
		tokens := limiter.Tokens()
		remaining := int(tokens)
		if remaining < 0 {
			remaining = 0
		}

		// 设置速率限制信息头
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

		c.Next()
	}
}

// ContentTypeMiddleware 验证 POST/PUT 请求的 Content-Type
// 📌 用途：确保客户端发送正确格式的数据
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只检查 POST 和 PUT 请求
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")

			// 检查是否为 application/json
			// strings.Contains 因为可能是 "application/json; charset=utf-8"
			if !strings.Contains(contentType, "application/json") {
				requestID, _ := c.Get("request_id")
				c.JSON(http.StatusUnsupportedMediaType, APIResponse{
					Success:   false,
					Error:     "Content-Type must be application/json",
					RequestID: fmt.Sprintf("%v", requestID),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ErrorHandlerMiddleware 捕获 panic 并返回友好的错误信息
// 📌 用途：防止服务器崩溃，优雅地处理错误
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 获取 request ID
		requestID, _ := c.Get("request_id")

		// 记录详细错误（仅服务器日志，包含真实panic信息）
		log.Printf("[ERROR] [%v] Panic recovered: %v", requestID, recovered)

		// 返回统一的错误响应（完全不暴露panic细节）
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(http.StatusInternalServerError)

		recoveredMsg := fmt.Sprint(recovered)
		response := APIResponse{
			Success:   false,
			Message:   recoveredMsg,
			Error:     "Internal server error",
			RequestID: fmt.Sprintf("%v", requestID),
		}

		// 手动序列化JSON
		jsonBytes, _ := json.Marshal(response)
		c.Writer.Write(jsonBytes)
	})
}

// Sanitize500Middleware 兜底清洗 500 响应体，防止泄露 panic 文本
type sanitizeWriter struct {
	gin.ResponseWriter
	status int
	buf    []byte
}

func (w *sanitizeWriter) WriteHeader(code int) {
	w.status = code
	// 延迟写出，由中间件收尾统一处理
}

func (w *sanitizeWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func Sanitize500Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 包装 writer 拦截写入
		sw := &sanitizeWriter{ResponseWriter: c.Writer}
		c.Writer = sw

		c.Next()

		// 确定最终状态码
		status := sw.status
		if status == 0 {
			status = sw.ResponseWriter.Status()
		}

		// 500 时强制覆盖为统一响应
		if status == http.StatusInternalServerError {
			requestID, _ := c.Get("request_id")
			// 使用捕获的缓冲内容作为 panic 信息（如果有），否则给空字符串
			panicMsg := string(sw.buf)
			resp := APIResponse{
				Success:   false,
				Message:   panicMsg,
				Error:     "Internal server error",
				RequestID: fmt.Sprintf("%v", requestID),
			}
			body, _ := json.Marshal(resp)
			w := sw.ResponseWriter
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(body)
			return
		}

		// 非 500：按原样下发缓冲体
		w := sw.ResponseWriter
		if status != 0 {
			w.WriteHeader(status)
		}
		if len(sw.buf) > 0 {
			_, _ = w.Write(sw.buf)
		}
	}
}

// ============================================================================
// 路由处理函数
// ============================================================================

// ping 处理健康检查请求
func ping(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "pong",
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// getArticles 获取所有文章
func getArticles(c *gin.Context) {
	// 读锁：允许多个并发读取
	articlesMutex.RLock()
	defer articlesMutex.RUnlock()

	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      articles,
		Message:   "Articles retrieved successfully",
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// getArticle 获取单篇文章
func getArticle(c *gin.Context) {
	// 📌 获取 URL 参数
	// 路由定义为 /articles/:id，这里获取 :id 的值
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Invalid article ID",
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 查找文章
	article, _ := findArticleByID(id)
	if article == nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Article not found",
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      article,
		Message:   "Article retrieved successfully",
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// createArticle 创建新文章（需要认证）
func createArticle(c *gin.Context) {
	var article Article

	// 📌 解析 JSON 请求体
	// ShouldBindJSON 会自动验证 JSON 格式
	if err := c.ShouldBindJSON(&article); err != nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Invalid request body: " + err.Error(),
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 验证文章数据
	if err := validateArticle(article); err != nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 写锁：独占访问
	articlesMutex.Lock()
	// 设置文章属性
	article.ID = nextID
	nextID++
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	// 添加到存储
	articles = append(articles, article)
	articlesMutex.Unlock()

	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Data:      article,
		Message:   "Article created successfully",
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// updateArticle 更新文章（需要认证）
func updateArticle(c *gin.Context) {
	// 获取文章 ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Invalid article ID",
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 解析更新数据
	var updatedArticle Article
	if err := c.ShouldBindJSON(&updatedArticle); err != nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Invalid request body: " + err.Error(),
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 验证数据
	if err := validateArticle(updatedArticle); err != nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     err.Error(),
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 查找并更新文章
	articlesMutex.Lock()
	defer articlesMutex.Unlock()

	article, idx := findArticleByID(id)
	if article == nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Article not found",
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 更新字段（保持 ID 和 CreatedAt 不变）
	updatedArticle.ID = id
	updatedArticle.CreatedAt = article.CreatedAt
	updatedArticle.UpdatedAt = time.Now()

	articles[idx] = updatedArticle

	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      updatedArticle,
		Message:   "Article updated successfully",
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// deleteArticle 删除文章（需要认证）
func deleteArticle(c *gin.Context) {
	// 获取文章 ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Error:     "Invalid article ID",
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 查找并删除文章
	articlesMutex.Lock()
	defer articlesMutex.Unlock()

	article, idx := findArticleByID(id)
	if article == nil {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Article not found",
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	// 从切片中删除
	articles = append(articles[:idx], articles[idx+1:]...)

	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "Article deleted successfully",
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// getStats 获取统计信息（仅管理员）
func getStats(c *gin.Context) {
	// 📌 检查用户角色
	role, exists := c.Get("user_role")
	if !exists || role != "admin" {
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusForbidden, APIResponse{
			Success:   false,
			Error:     "Admin access required",
			RequestID: fmt.Sprintf("%v", requestID),
		})
		return
	}

	articlesMutex.RLock()
	totalArticles := len(articles)
	articlesMutex.RUnlock()

	stats := map[string]interface{}{
		"total_articles": totalArticles,
		"total_authors":  2, // 简化示例
		"uptime":         "24h",
		"version":        "1.0.0",
	}

	requestID, _ := c.Get("request_id")
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      stats,
		Message:   "Statistics retrieved successfully",
		RequestID: fmt.Sprintf("%v", requestID),
	})
}

// ============================================================================
// 辅助函数
// ============================================================================

// findArticleByID 根据 ID 查找文章
// 返回文章指针和索引，未找到则返回 nil 和 -1
func findArticleByID(id int) (*Article, int) {
	for i := range articles {
		if articles[i].ID == id {
			return &articles[i], i
		}
	}
	return nil, -1
}

// validateArticle 验证文章数据
func validateArticle(article Article) error {
	if strings.TrimSpace(article.Title) == "" {
		return errors.New("title is required")
	}
	if len(article.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}
	if strings.TrimSpace(article.Content) == "" {
		return errors.New("content is required")
	}
	if strings.TrimSpace(article.Author) == "" {
		return errors.New("author is required")
	}
	return nil
}
