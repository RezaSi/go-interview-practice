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

func TestTaskAPI(t *testing.T) {
	// Reset tasks slice before each test
	resetTasks()

	// Create a new Echo instance for testing
	e := echo.New()
	setupRoutes(e)

	t.Run("GET /tasks - should return empty list initially", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response TaskListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, 0, response.Count)
		assert.Empty(t, response.Tasks)
	})

	t.Run("POST /tasks - should create a new task", func(t *testing.T) {
		taskJSON := `{"title":"Learn Echo","description":"Master Echo web framework","priority":"high"}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response TaskResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Task created successfully", response.Message)
		assert.NotNil(t, response.Task)
		assert.Equal(t, "Learn Echo", response.Task.Title)
		assert.Equal(t, "Master Echo web framework", response.Task.Description)
		assert.Equal(t, "high", response.Task.Priority)
		assert.False(t, response.Task.Completed)
		assert.NotEmpty(t, response.Task.ID)
		assert.NotEmpty(t, response.Task.CreatedAt)
	})

	t.Run("POST /tasks - should validate required fields", func(t *testing.T) {
		taskJSON := `{"description":"Missing title"}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "title")
	})

	t.Run("GET /tasks/:id - should return specific task", func(t *testing.T) {
		// First create a task
		taskJSON := `{"title":"Test Task","description":"Test description","priority":"medium"}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse TaskResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		taskID := createResponse.Task.ID

		// Now get the task
		req = httptest.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response TaskResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, taskID, response.Task.ID)
		assert.Equal(t, "Test Task", response.Task.Title)
	})

	t.Run("GET /tasks/:id - should return 404 for non-existent task", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks/nonexistent", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "not found")
	})

	t.Run("PUT /tasks/:id - should update existing task", func(t *testing.T) {
		// First create a task
		taskJSON := `{"title":"Original Task","description":"Original description","priority":"low"}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse TaskResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		taskID := createResponse.Task.ID

		// Now update the task
		updateJSON := `{"title":"Updated Task","description":"Updated description","priority":"high","completed":true}`
		req = httptest.NewRequest(http.MethodPut, "/tasks/"+taskID, strings.NewReader(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response TaskResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Task updated successfully", response.Message)
		assert.Equal(t, "Updated Task", response.Task.Title)
		assert.Equal(t, "Updated description", response.Task.Description)
		assert.Equal(t, "high", response.Task.Priority)
		assert.True(t, response.Task.Completed)
	})

	t.Run("DELETE /tasks/:id - should delete existing task", func(t *testing.T) {
		// First create a task
		taskJSON := `{"title":"Task to Delete","description":"Will be deleted","priority":"medium"}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse TaskResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		taskID := createResponse.Task.ID

		// Now delete the task
		req = httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response MessageResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Task deleted successfully", response.Message)

		// Verify task is deleted
		req = httptest.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("GET /tasks with query parameters", func(t *testing.T) {
		// Reset tasks for this specific test
		resetTasks()

		// Create multiple tasks
		tasks := []string{
			`{"title":"High Priority Task","description":"Important task","priority":"high"}`,
			`{"title":"Medium Priority Task","description":"Normal task","priority":"medium"}`,
			`{"title":"Low Priority Task","description":"Less important","priority":"low"}`,
		}

		for _, taskJSON := range tasks {
			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskJSON))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		// Test filtering by priority
		req := httptest.NewRequest(http.MethodGet, "/tasks?priority=high", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response TaskListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, 1, response.Count)
		assert.Equal(t, "high", response.Tasks[0].Priority)

		// Test filtering by completed status
		req = httptest.NewRequest(http.MethodGet, "/tasks?completed=false", nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		// All tasks should be incomplete initially
		for _, task := range response.Tasks {
			assert.False(t, task.Completed)
		}
	})

	t.Run("GET /health - should return health status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response HealthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "Task API is running", response.Message)
		assert.NotEmpty(t, response.Timestamp)
		assert.Equal(t, "1.0.0", response.Version)
	})

	t.Run("Invalid JSON should return 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"invalid": json}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, strings.ToLower(response.Message), "invalid")
	})
}

func TestTaskStructure(t *testing.T) {
	t.Run("Task struct should have correct fields", func(t *testing.T) {
		task := Task{
			ID:          "test-id",
			Title:       "Test Task",
			Description: "Test Description",
			Priority:    "high",
			Completed:   false,
			CreatedAt:   "2024-01-01T00:00:00Z",
		}

		assert.Equal(t, "test-id", task.ID)
		assert.Equal(t, "Test Task", task.Title)
		assert.Equal(t, "Test Description", task.Description)
		assert.Equal(t, "high", task.Priority)
		assert.False(t, task.Completed)
		assert.Equal(t, "2024-01-01T00:00:00Z", task.CreatedAt)
	})

	t.Run("Response structs should have correct structure", func(t *testing.T) {
		task := Task{ID: "1", Title: "Test"}

		taskResponse := TaskResponse{
			Status:  "success",
			Message: "Task created",
			Task:    &task,
		}

		assert.Equal(t, "success", taskResponse.Status)
		assert.Equal(t, "Task created", taskResponse.Message)
		assert.NotNil(t, taskResponse.Task)

		listResponse := TaskListResponse{
			Status: "success",
			Count:  1,
			Tasks:  []Task{task},
		}

		assert.Equal(t, "success", listResponse.Status)
		assert.Equal(t, 1, listResponse.Count)
		assert.Len(t, listResponse.Tasks, 1)
	})
}
