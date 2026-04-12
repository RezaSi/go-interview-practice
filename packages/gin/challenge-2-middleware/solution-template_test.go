package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Reset articles data for each test
	articles = []Article{
		{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	nextID = 3

	// Create router with middleware (should match main function setup)
	router := gin.New()

	// Add middleware in correct order
	router.Use(ErrorHandlerMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(LoggingMiddleware())
	router.Use(CORSMiddleware())
	router.Use(RateLimitMiddleware())
	router.Use(ContentTypeMiddleware())

	// Public routes
	router.GET("/ping", ping)
	router.GET("/articles", getArticles)
	router.GET("/articles/:id", getArticle)

	// Protected routes
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/articles", createArticle)
		protected.PUT("/articles/:id", updateArticle)
		protected.DELETE("/articles/:id", deleteArticle)
		protected.GET("/admin/stats", getStats)
	}

	return router
}

// Test Health Check
func TestPing(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotEmpty(t, response.RequestID)
}

// Test Request ID Middleware
func TestRequestIDMiddleware(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	// Check that X-Request-ID header is set
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)

	// Check that request ID is in response
	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, requestID, response.RequestID)
}

// Test CORS Middleware
func TestCORSMiddleware(t *testing.T) {
	router := setupRouter()

	// Test allowed origin
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")

	// Test preflight OPTIONS request
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("OPTIONS", "/articles", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

// Test Rate Limiting Middleware - Basic Check
func TestRateLimitMiddlewareBasic(t *testing.T) {
	router := setupRouter()

	// Make a few requests to test basic rate limiting functionality
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(w, req)

		// Check rate limit headers are present (if implemented)
		limit := w.Header().Get("X-RateLimit-Limit")
		if limit != "" {
			assert.NotEmpty(t, limit)
		}

		remaining := w.Header().Get("X-RateLimit-Remaining")
		if remaining != "" {
			assert.NotEmpty(t, remaining)
		}

		reset := w.Header().Get("X-RateLimit-Reset")
		if reset != "" {
			assert.NotEmpty(t, reset)
		}
	}
}

// Test Content Type Middleware
func TestContentTypeMiddleware(t *testing.T) {
	router := setupRouter()

	// Test POST without JSON content type
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/articles", bytes.NewBufferString("invalid"))
	req.Header.Set("X-API-Key", "admin-key-123")
	req.Header.Set("Content-Type", "text/plain")
	router.ServeHTTP(w, req)

	assert.Equal(t, 415, w.Code)

	// Test POST with correct JSON content type
	articleData := map[string]interface{}{
		"title":   "Test Article",
		"content": "Test content",
		"author":  "Test Author",
	}
	jsonData, _ := json.Marshal(articleData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonData))
	req.Header.Set("X-API-Key", "admin-key-123")
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, 415, w.Code) // Should not be content type error
}

// Test Authentication Middleware
func TestAuthMiddleware(t *testing.T) {
	router := setupRouter()

	// Test without API key
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/articles", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	// Test with invalid API key
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/articles", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "invalid-key")
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	// Test with valid admin API key
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/admin/stats", nil)
	req.Header.Set("X-API-Key", "admin-key-123")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, 401, w.Code) // Should not be unauthorized
}

// Test Public Routes (No Auth Required)
func TestGetArticles(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/articles", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check if articles data is returned
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(data))
}

func TestGetArticleByID(t *testing.T) {
	router := setupRouter()

	// Test valid article ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/articles/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Test invalid article ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/articles/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

// Test Protected Routes
func TestCreateArticle(t *testing.T) {
	router := setupRouter()

	articleData := map[string]interface{}{
		"title":   "New Test Article",
		"content": "This is test content",
		"author":  "Test Author",
	}
	jsonData, _ := json.Marshal(articleData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonData))
	req.Header.Set("X-API-Key", "admin-key-123")
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestUpdateArticle(t *testing.T) {
	router := setupRouter()

	updateData := map[string]interface{}{
		"title":   "Updated Title",
		"content": "Updated content",
		"author":  "Updated Author",
	}
	jsonData, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/articles/1", bytes.NewBuffer(jsonData))
	req.Header.Set("X-API-Key", "admin-key-123")
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestDeleteArticle(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/articles/1", nil)
	req.Header.Set("X-API-Key", "admin-key-123")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify article is deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/articles/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

// Test Admin-Only Routes
func TestGetStatsAdminOnly(t *testing.T) {
	router := setupRouter()

	// Test with user key (should fail if admin-only)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/stats", nil)
	req.Header.Set("X-API-Key", "user-key-456")
	router.ServeHTTP(w, req)

	// Should either be 403 (if role checking implemented) or 200 (if not)
	assert.True(t, w.Code == 200 || w.Code == 403)

	// Test with admin key (should succeed)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/admin/stats", nil)
	req.Header.Set("X-API-Key", "admin-key-123")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// Test Error Handling
func TestErrorHandling(t *testing.T) {
	router := setupRouter()

	// Test invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/articles", bytes.NewBufferString("invalid json"))
	req.Header.Set("X-API-Key", "admin-key-123")
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	// Test invalid article ID format
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/articles/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

// Test Error Middleware
func TestErrorMiddleware(t *testing.T) {
	// Specific router to test the Error middleware.
	// To simulate an internal server error we will use a divide by zero trick and
	// check if the Error middleware is doing his job correctly by returning
	// a code 500 with correct JSON message.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.Use(ErrorHandlerMiddleware())
	router.GET("/div/:a/:b", func(c *gin.Context) {
		a, _ := strconv.Atoi(c.Param("a"))
		b, _ := strconv.Atoi(c.Param("b"))
		val := a / b // No check on 'b' value, we want to allow a div by zero
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: val})
	})

	w := httptest.NewRecorder()
	// a = 5, b = 0     a / b --> booom
	req, _ := http.NewRequest("GET", "/div/5/0", nil)
	router.ServeHTTP(w, req)

	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)

	assert.Equal(t, 500, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, response.Error, "Internal server error")
	assert.Equal(t, response.Message, "runtime error: integer divide by zero")
	assert.NotEmpty(t, response.RequestID)
}

// Test Middleware Integration
func TestMiddlewareIntegration(t *testing.T) {
	router := setupRouter()

	// Test that all middleware work together
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/articles", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Check that multiple middleware effects are present
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))                // RequestID middleware
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin")) // CORS middleware

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotEmpty(t, response.RequestID) // RequestID in response
}

// -----------------------------------------------------------------
// Test Rate Limiting Middleware
//
// -- IMPORTANT --
// The following test must be the last one otherwise it will
// cause rate limit problem to all other tests because it will
// consume all available tokens and reach the rate limit.
// -----------------------------------------------------------------

func TestRateLimitMiddleware(t *testing.T) {
	router := setupRouter()

	// NOTE: Description of the rate limiter to implement
	// Limit: 100 requests per IP per minute
	// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
	// Return 429 if rate limit exceeded

	var limit, remain, reset int = 0, 0, 0
	var err error

	// Do a first request to capture the remaining token value
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	limitStr := w.Header().Get("X-RateLimit-Limit")
	limit, err = strconv.Atoi(limitStr)
	assert.NoError(t, err)
	assert.Equal(t, 100, limit)

	remainStr := w.Header().Get("X-RateLimit-Remaining")
	remain, err = strconv.Atoi(remainStr)
	assert.NoError(t, err)
	assert.Greater(t, remain, 0)

	resetStr := w.Header().Get("X-RateLimit-Reset")
	reset, err = strconv.Atoi(resetStr)
	assert.NoError(t, err)
	assert.Greater(t, reset, int(time.Now().Unix()))

	limitIndex := remain

	for i := 1; i < 102; i++ {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(w, req)

		limitStr = w.Header().Get("X-RateLimit-Limit")
		limit, err = strconv.Atoi(limitStr)
		assert.NoError(t, err)
		assert.Equal(t, 100, limit)

		remainStr = w.Header().Get("X-RateLimit-Remaining")
		remain, err = strconv.Atoi(remainStr)
		assert.NoError(t, err)

		resetStr = w.Header().Get("X-RateLimit-Reset")
		reset, err = strconv.Atoi(resetStr)
		assert.NoError(t, err)
		assert.Greater(t, reset, int(time.Now().Unix()))

		// Allowed requests should succeed (with remain > 0)
		if i < limitIndex {
			assert.Equal(t, w.Code, 200)
			assert.Greater(t, remain, 0)
		}

		// Last allowed request should succeed (with remain == 0)
		if i == limitIndex {
			assert.Equal(t, w.Code, 200)
			assert.Equal(t, remain, 0)
		}

		// All Requests over the rate limit should fail
		if i > limitIndex {
			assert.Equal(t, w.Code, 429)
			assert.Equal(t, remain, 0)
		}
	}
}
