package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
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

// Endpoints
const (
	articlesEndpoint = "/articles"
	articleEndpoint  = "/articles/:id"
	pingEndpoint     = "/ping"
	statsEndpoint    = "/admin/stats"
)

// Response messages
const (
	articleNotFoundErr          = "Article not found"
	missingApiKeyErr            = "Missing API key"
	invalidApiKeyErr            = "Invalid API key"
	rateLimitExceededErr        = "Rate limit exceeded"
	mustBeJsonErr               = "Unsupported Media Type. Content-Type must be application/json"
	errorBuildingCorsHeadersErr = "An error occurred whilst building CORS headers"
	internalServerErrorErr      = "Internal server error"
	invalidRequestBodyErr       = "Invalid request body"
	articleNotSpecifiedErr      = "Article not specified"
	titleRequiredErr            = "title is required"
	contentRequiredErr          = "content is required"
	authorRequiredErr           = "author is required"
)

// Context parameters
const (
	requestIdCtx = "request_id"
	userRole     = "user_role"
)

// Headers
const (
	requestIdHeader          = "X-Request-ID"
	apiKeyHeader             = "X-API-Key"
	rateLimitLimitHeader     = "X-RateLimit-Limit"
	rateLimitRemainingHeader = "X-RateLimit-Remaining"
	rateLimitResetHeader     = "X-RateLimit-Reset"
	corsOriginHeader         = "Access-Control-Allow-Origin"
	corsMethodHeader         = "Access-Control-Allow-Methods"
	corsHeadersHeader        = "Access-Control-Allow-Headers"
	contentTypeHeader        = "Content-Type"
)

// Content-types
const (
	jsonContentType = "application/json"
)

var (
	corsAllowedHeaders = []string{contentTypeHeader, apiKeyHeader, requestIdHeader}
	corsAllowedOrigins = []string{"http://localhost:3000"}
	corsAllowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}

	corsAllowedHeadersVal = strings.Join(corsAllowedHeaders[:], ", ")
	corsAllowedMethodsVal = strings.Join(corsAllowedMethods[:], ", ")

	rateLimits        = make(map[string]*rate.Limiter)
	requestsPerMinute = 100

	defaultPage     = 1
	defaultPageSize = 10
)

func main() {
	r := gin.New()

	r.Use(ErrorHandlerMiddleware(),
		RequestIDMiddleware(),
		LoggingMiddleware(),
		CORSMiddleware(),
		RateLimitMiddleware(),
		ContentTypeMiddleware())

	r.GET(pingEndpoint, ping)
	r.GET(articlesEndpoint, getArticles)
	r.GET(articleEndpoint, getArticle)

	r.Use(AuthMiddleware())
	{
		r.POST(articlesEndpoint, createArticle)
		r.PUT(articleEndpoint, updateArticle)
		r.DELETE(articleEndpoint, deleteArticle)
		r.GET(statsEndpoint, getStats)
	}

	r.Run()
}

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := uuid.New().String()
		c.Set(requestIdCtx, reqId)
		c.Header(requestIdHeader, reqId)

		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()
		duration := end.Sub(start)

		fmt.Printf("[%s] %s %s %d %v %s %s\n", c.GetString(requestIdCtx),
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(),
			duration, c.ClientIP(), c.Request.UserAgent())
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	keys := map[string]string{
		"admin-key-123": "admin",
		"user-key-456":  "user",
	}

	return func(c *gin.Context) {
		header := c.GetHeader(apiKeyHeader)
		if header == "" {
			unauthorised(c, nil, "", missingApiKeyErr)
			c.Abort()
			return
		}

		if role, ok := keys[header]; ok {
			c.Set(userRole, role)

			c.Next()
			return
		}

		unauthorised(c, nil, "", invalidApiKeyErr)
		c.Abort()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, origin := range corsAllowedOrigins {
			c.Header(corsOriginHeader, origin)
		}

		c.Header(corsMethodHeader, corsAllowedMethodsVal)
		c.Header(corsHeadersHeader, corsAllowedHeadersVal)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rateLimiter *rate.Limiter
		var exists bool
		if rateLimiter, exists = rateLimits[c.ClientIP()]; !exists {
			rateLimiter = rate.NewLimiter(rate.Limit(requestsPerMinute/60), requestsPerMinute*2)
			rateLimits[c.ClientIP()] = rateLimiter
		}

		remaining := rateLimiter.Tokens() - 1

		c.Header(rateLimitLimitHeader, strconv.Itoa(requestsPerMinute))
		c.Header(rateLimitRemainingHeader, strconv.Itoa(int(remaining)))
		c.Header(rateLimitResetHeader, strconv.Itoa(int(time.Now().Unix()+60)))

		if !rateLimiter.Allow() {
			tooManyRequests(c, nil, "", rateLimitExceededErr)
			c.Abort()
			return
		}
		c.Next()
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if (c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut) &&
			c.ContentType() != jsonContentType {
			unsupportedMediaType(c, nil, "", mustBeJsonErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		message := ""

		if err, ok := recovered.(string); ok {
			message = err
		} else if err, ok := recovered.(error); ok {
			message = err.Error()
		}

		internalServerError(c, nil, message, internalServerErrorErr)
	})
}

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	ok(c, "pong", "", "")
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	size := c.Query("size")
	page := c.Query("page")

	// Convert size and page to integers
	sizeInt, err := strconv.Atoi(size)
	if err != nil || sizeInt <= 0 {
		sizeInt = defaultPageSize
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = defaultPage
	}

	start := (pageInt - 1) * sizeInt
	end := start + sizeInt

	if start > len(articles) {
		start = len(articles)
	}
	if end > len(articles) {
		end = len(articles)
	}

	ok(c, articles[start:end], "", "")
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {

	var id int
	var err error
	if id, err = getRequestedArticleId(c); err != nil {
		badRequest(c, nil, "", err.Error())
		return
	}

	article, _ := findArticleByID(id)
	if article == nil {
		notFound(c, nil, "", articleNotFoundErr)
		return
	}

	ok(c, article, "", "")
}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		badRequest(c, nil, "", invalidRequestBodyErr)
		return
	}

	if err := validateArticle(newArticle); err != "" {
		badRequest(c, nil, "", err)
		return
	}

	newArticle.ID = nextID
	nextID++
	newArticle.CreatedAt = time.Now()
	newArticle.UpdatedAt = time.Now()
	articles = append(articles, newArticle)

	created(c, newArticle, "", "")
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	id, err := getRequestedArticleId(c)
	if err != nil {
		badRequest(c, nil, "", articleNotSpecifiedErr)
		return
	}
	article, index := findArticleByID(id)
	if article == nil {
		notFound(c, nil, "", articleNotFoundErr)
		return
	}

	var updatedArticle Article
	if err := c.ShouldBindJSON(&updatedArticle); err != nil {
		badRequest(c, nil, "", invalidRequestBodyErr)
		return
	}

	if err := validateArticle(updatedArticle); err != "" {
		badRequest(c, nil, "", err)
		return
	}

	// Update article fields
	articles[index].Title = updatedArticle.Title
	articles[index].Content = updatedArticle.Content
	articles[index].Author = updatedArticle.Author
	articles[index].UpdatedAt = time.Now()

	ok(c, articles[index], "", "")
}

// deleteArticle handles DELETE /articles/:id - delete article (protected)
func deleteArticle(c *gin.Context) {
	id, err := getRequestedArticleId(c)
	if err != nil {
		badRequest(c, nil, "", articleNotSpecifiedErr)
		return
	}

	article, index := findArticleByID(id)
	if article == nil {
		notFound(c, nil, "", articleNotFoundErr)
		return
	}

	// Remove article from slice
	articles = append(articles[:index], articles[index+1:]...)

	ok(c, "", "article deleted successfully", "")
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         "24h",
	}

	ok(c, stats, "", "")
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// Return article pointer and index, or nil and -1 if not found
	index := slices.IndexFunc(articles, func(a Article) bool {
		return a.ID == id
	})

	if index != -1 {
		return &articles[index], index
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) string {
	switch {
	case article.Title == "":
		return titleRequiredErr
	case article.Content == "":
		return contentRequiredErr
	case article.Author == "":
		return authorRequiredErr
	default:
		return ""
	}
}

// getRequestedArticleId retrieves and validates id parameter from URI
func getRequestedArticleId(c *gin.Context) (int, error) {

	id := c.Param("id")
	if id == "" {
		return -1, errors.New("missing article id")
	}

	articleId, err := strconv.Atoi(id)
	if err != nil {
		return -1, errors.New("invalid article id")
	}
	return articleId, nil
}

func ok(c *gin.Context, data any, message string, err string) {
	successResponse(c, http.StatusOK, data, message, err)
}
func created(c *gin.Context, data any, message string, err string) {
	successResponse(c, http.StatusCreated, data, message, err)
}

func badRequest(c *gin.Context, data any, message string, err string) {
	failureResponse(c, http.StatusBadRequest, data, message, err)
}
func unauthorised(c *gin.Context, data any, message string, err string) {
	failureResponse(c, http.StatusUnauthorized, data, message, err)
}
func notFound(c *gin.Context, data any, message string, err string) {
	failureResponse(c, http.StatusNotFound, data, message, err)
}
func unsupportedMediaType(c *gin.Context, data any, message string, err string) {
	failureResponse(c, http.StatusUnsupportedMediaType, data, message, err)
}
func tooManyRequests(c *gin.Context, data any, message string, err string) {
	failureResponse(c, http.StatusTooManyRequests, data, message, err)
}

func internalServerError(c *gin.Context, data any, message string, err string) {
	failureResponse(c, http.StatusInternalServerError, data, message, err)
}

func successResponse(c *gin.Context, statusCode int, data any, message string, err string) {
	response(c, statusCode, true, data, message, err)
}
func failureResponse(c *gin.Context, statusCode int, data any, message string, err string) {
	response(c, statusCode, false, data, message, err)
}
func response(c *gin.Context, statusCode int, success bool, data any, message string, err string) {
	c.JSON(statusCode, newResponse(c, success, data, message, err))
}

func newResponse(c *gin.Context, success bool, data any, message string, err string) APIResponse {
	resp := APIResponse{
		Success:   success,
		RequestID: c.GetString(requestIdCtx),
	}

	if data != nil {
		resp.Data = data
	}

	if message != "" {
		resp.Message = message
	}

	if err != "" {
		resp.Error = err
	}

	return resp
}
