package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	// Reset articles data for each test
	articles = []Article{
		{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Web Development with Fiber", Content: "Fiber is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	nextID = 3

	// Create app with middleware (should match main function setup)
	app := fiber.New()

	// Add middleware in correct order (implement these in solution)
	app.Use(ErrorHandlerMiddleware())
	app.Use(RequestIDMiddleware())
	app.Use(LoggingMiddleware())
	app.Use(CORSMiddleware())
	app.Use(RateLimitMiddleware())

	// Public routes
	app.Get("/ping", pingHandler)
	app.Get("/articles", getArticlesHandler)
	app.Get("/articles/:id", getArticleHandler)

	// Protected routes
	protected := app.Group("/", AuthMiddleware())
	protected.Post("/articles", createArticleHandler)
	protected.Put("/articles/:id", updateArticleHandler)
	protected.Delete("/articles/:id", deleteArticleHandler)
	protected.Get("/admin/stats", getStatsHandler)

	return app
}

func TestPingHandler(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/ping", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Check for request ID header
	assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
}

func TestCORSMiddleware(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("OPTIONS", "/ping", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "GET")
}

func TestGetArticles(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/articles", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func TestAuthMiddleware(t *testing.T) {
	app := setupTestApp()

	// Test without API key
	req := httptest.NewRequest("POST", "/articles", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	// Test with valid API key
	body := bytes.NewBufferString(`{"title":"Test Article","content":"Test content","author":"Test Author"}`)
	req = httptest.NewRequest("POST", "/articles", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "admin-key-123")
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestRequestIDGeneration(t *testing.T) {
	app := setupTestApp()

	req1 := httptest.NewRequest("GET", "/ping", nil)
	resp1, _ := app.Test(req1)

	req2 := httptest.NewRequest("GET", "/ping", nil)
	resp2, _ := app.Test(req2)

	id1 := resp1.Header.Get("X-Request-ID")
	id2 := resp2.Header.Get("X-Request-ID")

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
}

func TestCreateArticle(t *testing.T) {
	app := setupTestApp()

	articleData := map[string]interface{}{
		"title":   "New Article",
		"content": "New article content",
		"author":  "Test Author",
	}

	body, _ := json.Marshal(articleData)
	req := httptest.NewRequest("POST", "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "admin-key-123")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var response APIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
}

func TestGetArticleByID(t *testing.T) {
	app := setupTestApp()

	// Test existing article
	req := httptest.NewRequest("GET", "/articles/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test non-existent article
	req = httptest.NewRequest("GET", "/articles/999", nil)
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestUpdateArticle(t *testing.T) {
	app := setupTestApp()

	updateData := map[string]interface{}{
		"title":   "Updated Article",
		"content": "Updated content",
		"author":  "Updated Author",
	}

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/articles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "admin-key-123")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response APIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response.Success)
}

func TestDeleteArticle(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("DELETE", "/articles/1", nil)
	req.Header.Set("X-API-Key", "admin-key-123")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response APIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response.Success)
}

func TestGetStats(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	req.Header.Set("X-API-Key", "admin-key-123")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response APIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response.Success)
}

func TestRateLimiting(t *testing.T) {
	app := setupTestApp()

	// Make multiple requests to trigger rate limiting
	for i := 0; i < 105; i++ {
		req := httptest.NewRequest("GET", "/ping", nil)
		resp, _ := app.Test(req)

		if i >= 100 {
			// Should be rate limited after 100 requests
			assert.Equal(t, 429, resp.StatusCode)
			break
		}
	}
}

