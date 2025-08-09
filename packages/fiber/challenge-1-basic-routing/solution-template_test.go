package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	// Reset task store for each test
	taskStore = NewTaskStore()

	// Call the setupApp function from the user's solution
	return setupApp()
}

func TestPingEndpoint(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/ping", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	assert.Equal(t, "pong", response["message"])
}

func TestGetAllTasks(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/tasks", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var tasks []Task
	json.Unmarshal(body, &tasks)
	assert.Len(t, tasks, 2) // Default tasks from NewTaskStore
	assert.Equal(t, "Learn Go", tasks[0].Title)
}

func TestGetTaskByID(t *testing.T) {
	app := setupTestApp()

	// Test existing task
	req := httptest.NewRequest("GET", "/tasks/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var task Task
	json.Unmarshal(body, &task)
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, "Learn Go", task.Title)

	// Test non-existing task
	req = httptest.NewRequest("GET", "/tasks/999", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestCreateTask(t *testing.T) {
	app := setupTestApp()

	newTask := Task{
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
	}

	taskJSON, _ := json.Marshal(newTask)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var createdTask Task
	json.Unmarshal(body, &createdTask)
	assert.Equal(t, 3, createdTask.ID) // Should be next ID
	assert.Equal(t, "Test Task", createdTask.Title)
}

func TestUpdateTask(t *testing.T) {
	app := setupTestApp()

	updateTask := Task{
		Title:       "Updated Task",
		Description: "Updated Description",
		Completed:   true,
	}

	taskJSON, _ := json.Marshal(updateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var updatedTask Task
	json.Unmarshal(body, &updatedTask)
	assert.Equal(t, 1, updatedTask.ID)
	assert.Equal(t, "Updated Task", updatedTask.Title)
	assert.True(t, updatedTask.Completed)

	// Test non-existing task
	req = httptest.NewRequest("PUT", "/tasks/999", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestDeleteTask(t *testing.T) {
	app := setupTestApp()

	// Test deleting existing task
	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)

	// Verify task is deleted
	req = httptest.NewRequest("GET", "/tasks/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	// Test deleting non-existing task
	req = httptest.NewRequest("DELETE", "/tasks/999", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestInvalidJSON(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestInvalidTaskID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/tasks/invalid", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}
