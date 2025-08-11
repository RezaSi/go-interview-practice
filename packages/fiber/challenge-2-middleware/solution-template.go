package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
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
	{ID: 2, Title: "Web Development with Fiber", Content: "Fiber is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3

func main() {
	// TODO: Create Fiber app without default middleware
	// Use fiber.New() with custom config

	// TODO: Setup custom middleware in correct order
	// 1. ErrorHandlerMiddleware (first to catch panics)
	// 2. RequestIDMiddleware
	// 3. LoggingMiddleware
	// 4. CORSMiddleware
	// 5. RateLimitMiddleware

	// TODO: Setup route groups
	// Public routes (no authentication required)
	// Protected routes (require authentication)

	// TODO: Define routes
	// Public: GET /ping, GET /articles, GET /articles/:id
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats

	// TODO: Start server on port 3000
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Generate UUID for request ID
		// Use github.com/google/uuid package
		// Store in context locals as "request_id"
		// Add to response header as "X-Request-ID"

		return c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Capture start time

		err := c.Next()

		// TODO: Calculate duration and log request
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT

		return err
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Set CORS headers
		// Access-Control-Allow-Origin: *
		// Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
		// Access-Control-Allow-Headers: Content-Type,Authorization,X-API-Key
		// Access-Control-Allow-Credentials: true

		// TODO: Handle preflight OPTIONS requests

		return c.Next()
	}
}

// RateLimitMiddleware limits requests per IP
func RateLimitMiddleware() fiber.Handler {
	// TODO: Implement rate limiting using a map or external store
	// Limit: 100 requests per minute per IP
	// Use a sliding window or token bucket algorithm
	// Return 429 Too Many Requests when limit exceeded

	return func(c *fiber.Ctx) error {
		// TODO: Check rate limit for c.IP()
		// If exceeded, return 429 status

		return c.Next()
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Get API key from X-API-Key header
		// Valid keys: "admin-key-123", "user-key-456"
		// Return 401 Unauthorized if missing or invalid
		// Store key type in context for role-based access

		return c.Next()
	}
}

// ErrorHandlerMiddleware provides centralized error handling
func ErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Handle panics and errors
		// Recover from panics and return 500 status
		// Log errors with request ID
		// Return consistent error response format

		return c.Next()
	}
}

// TODO: Implement route handlers

// pingHandler handles health check requests
func pingHandler(c *fiber.Ctx) error {
	// TODO: Return simple pong response with request ID
	return nil
}

// getArticlesHandler returns all articles with pagination
func getArticlesHandler(c *fiber.Ctx) error {
	// TODO: Implement pagination using query parameters
	// ?page=1&limit=10 (default: page=1, limit=10)
	// Return articles with pagination metadata
	return nil
}

// getArticleHandler returns a specific article by ID
func getArticleHandler(c *fiber.Ctx) error {
	// TODO: Get article ID from URL parameter
	// Return 404 if article not found
	return nil
}

// createArticleHandler creates a new article
func createArticleHandler(c *fiber.Ctx) error {
	// TODO: Parse request body into Article struct
	// Validate required fields (title, content, author)
	// Add article to storage with auto-increment ID
	// Return created article
	return nil
}

// updateArticleHandler updates an existing article
func updateArticleHandler(c *fiber.Ctx) error {
	// TODO: Get article ID from URL parameter
	// Parse request body for updates
	// Update article if exists, return 404 if not found
	// Return updated article
	return nil
}

// deleteArticleHandler deletes an article
func deleteArticleHandler(c *fiber.Ctx) error {
	// TODO: Get article ID from URL parameter
	// Remove article from storage
	// Return 404 if article not found
	// Return success message
	return nil
}

// getStatsHandler returns API usage statistics (admin only)
func getStatsHandler(c *fiber.Ctx) error {
	// TODO: Return API statistics
	// Total articles, request count, etc.
	// Only accessible with admin API key
	return nil
}
