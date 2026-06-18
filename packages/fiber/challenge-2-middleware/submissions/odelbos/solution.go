package main

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
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
	{ID: 2, Title: "Web Development with Fiber", Content: "Fiber is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3
var articlesMu sync.RWMutex

var validKeys = map[string]string{
	"admin-key-123": "admin",
	"user-key-456":  "user",
}

var (
	rateLimiters   = make(map[string]*rate.Limiter)
	rateLimitMutex sync.Mutex
)

var (
	requestCount = 0
	statsMu      sync.Mutex
)

func main() {
	app := fiber.New()

	app.Use(RequestIDMiddleware())
	app.Use(ErrorHandlerMiddleware())
	app.Use(LoggingMiddleware())
	app.Use(CORSMiddleware())
	app.Use(RateLimitMiddleware())

	public := app.Group("/")
	public.Get("/ping", pingHandler)
	public.Get("/articles", getArticlesHandler)
	public.Get("/articles/:id", getArticleHandler)

	protected := app.Group("/").Use(AuthMiddleware())
	protected.Post("/articles", createArticleHandler)
	protected.Put("/articles/:id", updateArticleHandler)
	protected.Delete("/articles/:id", deleteArticleHandler)
	protected.Get("/admin/stats", getStatsHandler)

	app.Listen(":3000")
}

// -----------------------------------------------------------
// Middlewares
// -----------------------------------------------------------

func RequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Locals("request_id", requestID)
		c.Set("X-Request-ID", requestID)

		statsMu.Lock()
		requestCount++
		statsMu.Unlock()

		return c.Next()
	}
}

func LoggingMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		log.Printf("[%s] %s %s %d %v %s %s",
			getRequestId(c),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			time.Since(start),
			c.IP(),
			c.Get("User-Agent"),
		)
		return err
	}
}

func CORSMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-API-Key")
		c.Set("Access-Control-Allow-Credentials", "true")

		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	}
}

func RateLimitMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		rateLimitMutex.Lock()
		limiter, ok := rateLimiters[ip]
		if !ok {
			limiter = rate.NewLimiter(rate.Every(time.Minute/100), 100)
			rateLimiters[ip] = limiter
		}
		rateLimitMutex.Unlock()

		c.Set("X-RateLimit-Limit", "100")
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

		if !limiter.Allow() {
			c.Set("X-RateLimit-Remaining", "0")
			return errResponse(c, fiber.StatusTooManyRequests, "Rate limit exceeded")
		}

		remaining := int(limiter.Tokens())
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		return c.Next()
	}
}

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		role, ok := validKeys[apiKey]
		if !ok {
			return errResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}
		c.Locals("role", role)
		return c.Status(fiber.StatusCreated).Next()
	}
}

func ErrorHandlerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Recover from panics and return 500 status
		// Log errors with request ID
		// Return consistent error response format
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[%s] : PANIC : %v", getRequestId(c), r)
				errResponse(c, fiber.StatusInternalServerError, "Internal Server Error")
			}
		}()

		err := c.Next()
		if err != nil {
			log.Printf("[%s]: Error : %v", getRequestId(c), err)

			var code = fiber.StatusInternalServerError
			var message = "Internal Server Error"
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}
			errResponse(c, code, message)
			return nil
		}
		return nil
	}
}

// -----------------------------------------------------------
// Handlers
// -----------------------------------------------------------

func pingHandler(c fiber.Ctx) error {
	return okResponse(c, fiber.StatusOK, "pong", nil)
}

func getArticlesHandler(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	start := (page - 1) * limit
	end := start + limit

	articlesMu.RLock()
	defer articlesMu.RUnlock()

	if start > len(articles) {
		start = len(articles)
	}
	if end > len(articles) {
		end = len(articles)
	}
	return okResponse(c, fiber.StatusOK, "Articles", articles[start:end])
}

func getArticleHandler(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id == 0 {
		return errResponse(c, fiber.StatusBadRequest, "Invalid Id")
	}

	articlesMu.RLock()
	defer articlesMu.RUnlock()

	for _, a := range articles {
		if a.ID == id {
			return okResponse(c, fiber.StatusOK, "", a)
		}
	}
	return errResponse(c, fiber.StatusNotFound, "Not found")
}

func createArticleHandler(c fiber.Ctx) error {
	var article Article
	if err := c.Bind().Body(&article); err != nil {
		return errResponse(c, fiber.StatusBadRequest, "Invalid body")
	}

	if article.Title == "" || article.Content == "" || article.Author == "" {
		return errResponse(c, fiber.StatusBadRequest, "Title, content, author are required")
	}

	articlesMu.Lock()
	defer articlesMu.Unlock()

	article.ID = nextID
	article.CreatedAt = time.Now()
	article.UpdatedAt = article.CreatedAt
	articles = append(articles, article)
	nextID++

	return okResponse(c, fiber.StatusCreated, "Created successfully", article)
}

func updateArticleHandler(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id == 0 {
		return errResponse(c, fiber.StatusBadRequest, "Invalid Id")
	}

	var data Article
	if err := c.Bind().Body(&data); err != nil {
		return errResponse(c, fiber.StatusBadRequest, "Invalid body")
	}

	articlesMu.Lock()
	defer articlesMu.Unlock()

	for i, article := range articles {
		if article.ID == id {
			if data.Title != "" {
				articles[i].Title = data.Title
			}
			if data.Content != "" {
				articles[i].Content = data.Content
			}
			if data.Author != "" {
				articles[i].Author = data.Author
			}
			articles[i].UpdatedAt = time.Now()

			return okResponse(c, fiber.StatusOK, "Updated successfully", articles[i])
		}
	}
	return errResponse(c, fiber.StatusNotFound, "Not found")
}

func deleteArticleHandler(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id == 0 {
		return errResponse(c, fiber.StatusBadRequest, "Invalid Id")
	}

	articlesMu.Lock()
	defer articlesMu.Unlock()

	for i, article := range articles {
		if article.ID == id {
			articles = slices.Delete(articles, i, i+1)
			return okResponse(c, fiber.StatusOK, "Delted successfully", articles[i])
		}
	}
	return errResponse(c, fiber.StatusNotFound, "Not found")
}

func getStatsHandler(c fiber.Ctx) error {
	// Total articles, request count, etc.
	// Only accessible with admin API key
	if c.Locals("role") != "admin" {
		return errResponse(c, fiber.StatusForbidden, "Forbidden")
	}

	statsMu.Lock()
	defer statsMu.Unlock()

	var data = map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": requestCount,
	}
	return okResponse(c, fiber.StatusOK, "Statistics", data)
}

// -----------------------------------------------------------
// Helpers
// -----------------------------------------------------------

func getRequestId(c fiber.Ctx) string {
	id := c.Locals("request_id")
	if s, ok := id.(string); ok {
		return s
	}
	return ""
}

func okResponse(c fiber.Ctx, status int, msg string, data interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		RequestID: getRequestId(c),
	})
}

func errResponse(c fiber.Ctx, status int, err string) error {
	return c.Status(status).JSON(APIResponse{
		Success:   false,
		Error:     err,
		RequestID: getRequestId(c),
	})
}
