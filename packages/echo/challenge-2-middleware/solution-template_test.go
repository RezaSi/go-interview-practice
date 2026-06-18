package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlogAPI(t *testing.T) {
	// Create a new Echo instance for testing
	e := echo.New()
	setupRoutes(e)

	t.Run("GET /health - should return health status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response HealthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "Blog API is running", response.Message)
		assert.NotEmpty(t, response.Timestamp)
		assert.Equal(t, "1.0.0", response.Version)
	})

	t.Run("POST /posts - should create a new post with valid API key", func(t *testing.T) {
		postJSON := `{"title":"Echo Middleware","content":"Learning about Echo middleware patterns","author":"John Doe","tags":["echo","middleware","go"]}`
		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response PostResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Post created successfully", response.Message)
		assert.NotNil(t, response.Post)
		assert.Equal(t, "Echo Middleware", response.Post.Title)
		assert.Equal(t, "John Doe", response.Post.Author)
		assert.NotEmpty(t, response.Post.ID)
		assert.NotEmpty(t, response.Post.CreatedAt)
	})

	t.Run("POST /posts - should reject request without API key", func(t *testing.T) {
		postJSON := `{"title":"Test Post","content":"Test content","author":"Test Author"}`
		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "API key")
	})

	t.Run("POST /posts - should reject request with invalid API key", func(t *testing.T) {
		postJSON := `{"title":"Test Post","content":"Test content","author":"Test Author"}`
		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "invalid-key")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "Invalid API key")
	})

	t.Run("GET /posts - should return all posts", func(t *testing.T) {
		// First create a post
		postJSON := `{"title":"Test Post","content":"Test content","author":"Test Author"}`
		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Now get all posts
		req = httptest.NewRequest(http.MethodGet, "/posts", nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PostListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.GreaterOrEqual(t, response.Count, 1)
		assert.NotEmpty(t, response.Posts)
	})

	t.Run("GET /posts/:id - should return specific post", func(t *testing.T) {
		// First create a post
		postJSON := `{"title":"Specific Post","content":"Specific content","author":"Specific Author"}`
		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse PostResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		postID := createResponse.Post.ID

		// Now get the specific post
		req = httptest.NewRequest(http.MethodGet, "/posts/"+postID, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PostResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, postID, response.Post.ID)
		assert.Equal(t, "Specific Post", response.Post.Title)
	})

	t.Run("PUT /posts/:id - should update existing post with valid API key", func(t *testing.T) {
		// First create a post
		postJSON := `{"title":"Original Post","content":"Original content","author":"Original Author"}`
		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse PostResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		postID := createResponse.Post.ID

		// Now update the post
		updateJSON := `{"title":"Updated Post","content":"Updated content","author":"Updated Author","tags":["updated","echo"]}`
		req = httptest.NewRequest(http.MethodPut, "/posts/"+postID, strings.NewReader(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PostResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Post updated successfully", response.Message)
		assert.Equal(t, "Updated Post", response.Post.Title)
		assert.Equal(t, "Updated content", response.Post.Content)
	})

	t.Run("DELETE /posts/:id - should delete existing post with valid API key", func(t *testing.T) {
		// First create a post
		postJSON := `{"title":"Post to Delete","content":"Will be deleted","author":"Test Author"}`
		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse PostResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		postID := createResponse.Post.ID

		// Now delete the post
		req = httptest.NewRequest(http.MethodDelete, "/posts/"+postID, nil)
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response MessageResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Post deleted successfully", response.Message)
	})

	t.Run("Rate limiting should work", func(t *testing.T) {
		// Make multiple requests quickly from same IP
		clientIP := "192.168.1.100"

		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodGet, "/posts", nil)
			req.Header.Set("X-Real-IP", clientIP)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if i < 3 {
				// First 3 requests should succeed
				assert.Equal(t, http.StatusOK, rec.Code)
			} else {
				// Subsequent requests should be rate limited
				assert.Equal(t, http.StatusTooManyRequests, rec.Code)

				var response ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, "error", response.Status)
				assert.Contains(t, response.Message, "rate limit")
			}
		}
	})

	t.Run("CORS headers should be present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/posts", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Check CORS headers
		assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Origin"))
		assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Methods"))
		assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("Request ID should be generated and included in response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Check that X-Request-ID header is present
		requestID := rec.Header().Get("X-Request-ID")
		assert.NotEmpty(t, requestID)
		assert.Len(t, requestID, 36) // UUID length
	})

	t.Run("Request validation should work", func(t *testing.T) {
		// Test missing title
		postJSON := `{"content":"Content without title","author":"Test Author"}`
		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "title")

		// Test missing content
		postJSON = `{"title":"Title without content","author":"Test Author"}`
		req = httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-api-key-123")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "content")
	})
}

func TestMiddleware(t *testing.T) {
	t.Run("Logging middleware should work", func(t *testing.T) {
		e := echo.New()
		setupRoutes(e)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		// Logging is tested by observing console output during test runs
	})

	t.Run("Request ID middleware should generate unique IDs", func(t *testing.T) {
		e := echo.New()
		setupRoutes(e)

		// Make multiple requests
		var requestIDs []string
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			requestID := rec.Header().Get("X-Request-ID")
			assert.NotEmpty(t, requestID)
			requestIDs = append(requestIDs, requestID)
		}

		// All request IDs should be unique
		for i := 0; i < len(requestIDs); i++ {
			for j := i + 1; j < len(requestIDs); j++ {
				assert.NotEqual(t, requestIDs[i], requestIDs[j])
			}
		}
	})

	t.Run("API key middleware should validate keys correctly", func(t *testing.T) {
		e := echo.New()
		setupRoutes(e)

		testCases := []struct {
			name           string
			apiKey         string
			expectedStatus int
		}{
			{"Valid API key", "test-api-key-123", http.StatusCreated},
			{"Invalid API key", "invalid-key", http.StatusUnauthorized},
			{"Missing API key", "", http.StatusUnauthorized},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				postJSON := `{"title":"Test","content":"Test content","author":"Test"}`
				req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
				req.Header.Set("Content-Type", "application/json")

				if tc.apiKey != "" {
					req.Header.Set("X-API-Key", tc.apiKey)
				}

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				assert.Equal(t, tc.expectedStatus, rec.Code)
			})
		}
	})
}

func TestPostStructure(t *testing.T) {
	t.Run("Post struct should have correct fields", func(t *testing.T) {
		post := Post{
			ID:        "test-id",
			Title:     "Test Post",
			Content:   "Test Content",
			Author:    "Test Author",
			Tags:      []string{"test", "go"},
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		}

		assert.Equal(t, "test-id", post.ID)
		assert.Equal(t, "Test Post", post.Title)
		assert.Equal(t, "Test Content", post.Content)
		assert.Equal(t, "Test Author", post.Author)
		assert.Equal(t, []string{"test", "go"}, post.Tags)
		assert.Equal(t, "2024-01-01T00:00:00Z", post.CreatedAt)
		assert.Equal(t, "2024-01-01T00:00:00Z", post.UpdatedAt)
	})

	t.Run("Response structs should have correct structure", func(t *testing.T) {
		post := Post{ID: "1", Title: "Test"}

		postResponse := PostResponse{
			Status:  "success",
			Message: "Post created",
			Post:    &post,
		}

		assert.Equal(t, "success", postResponse.Status)
		assert.Equal(t, "Post created", postResponse.Message)
		assert.NotNil(t, postResponse.Post)

		listResponse := PostListResponse{
			Status: "success",
			Count:  1,
			Posts:  []Post{post},
		}

		assert.Equal(t, "success", listResponse.Status)
		assert.Equal(t, 1, listResponse.Count)
		assert.Len(t, listResponse.Posts, 1)
	})
}

