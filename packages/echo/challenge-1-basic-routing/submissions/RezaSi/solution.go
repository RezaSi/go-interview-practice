package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Task represents a task in our task management system
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"` // "low", "medium", "high"
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
}

// Response structures for consistent API responses
type TaskResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Task    *Task  `json:"task,omitempty"`
}

type TaskListResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	Tasks  []Task `json:"tasks"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// In-memory storage for tasks (in production, use a database)
var tasks []Task

// resetTasks clears the tasks slice for testing
func resetTasks() {
	tasks = []Task{}
}

func main() {
	// TODO: Create a new Echo instance
	// Hint: Use echo.New() to create a new Echo instance
	e := echo.New()

	// TODO: Set up routes by calling setupRoutes function
	// Hint: Call setupRoutes(e) to configure all the routes
	setupRoutes(e)

	// TODO: Start the server on port 8080
	// Hint: Use e.Logger.Fatal(e.Start(":8080")) to start the server
	e.Logger.Fatal(e.Start(":8080"))
}

func setupRoutes(e *echo.Echo) {
	// TODO: Implement health check endpoint
	// Hint: Use e.GET("/health", healthHandler) to register the health endpoint
	e.GET("/health", healthHandler)

	// TODO: Implement task management endpoints
	// Hint: Register the following routes:
	// - GET /tasks -> getAllTasks
	// - POST /tasks -> createTask
	// - GET /tasks/:id -> getTask
	// - PUT /tasks/:id -> updateTask
	// - DELETE /tasks/:id -> deleteTask
	e.GET("/tasks", getAllTasks)
	e.POST("/tasks", createTask)
	e.GET("/tasks/:id", getTask)
	e.PUT("/tasks/:id", updateTask)
	e.DELETE("/tasks/:id", deleteTask)
}

func healthHandler(c echo.Context) error {
	// TODO: Implement health check handler
	// Return a HealthResponse with:
	// - Status: "healthy"
	// - Message: "Task API is running"
	// - Timestamp: current time in RFC3339 format
	// - Version: "1.0.0"
	// Hint: Use c.JSON(http.StatusOK, response) to return JSON response
	response := HealthResponse{
		Status:    "healthy",
		Message:   "Task API is running",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}
	return c.JSON(http.StatusOK, response)
}

func getAllTasks(c echo.Context) error {
	// TODO: Implement get all tasks handler with filtering
	// 1. Get query parameters for filtering (priority, completed)
	// 2. Filter tasks based on query parameters
	// 3. Return TaskListResponse with filtered tasks
	// Hint: Use c.QueryParam("priority") to get query parameters
	// Hint: Use strconv.ParseBool() to parse boolean values

	priority := c.QueryParam("priority")
	completedStr := c.QueryParam("completed")

	filteredTasks := make([]Task, 0)

	for _, task := range tasks {
		// Filter by priority if specified
		if priority != "" && task.Priority != priority {
			continue
		}

		// Filter by completed status if specified
		if completedStr != "" {
			completed, err := strconv.ParseBool(completedStr)
			if err == nil && task.Completed != completed {
				continue
			}
		}

		filteredTasks = append(filteredTasks, task)
	}

	response := TaskListResponse{
		Status: "success",
		Count:  len(filteredTasks),
		Tasks:  filteredTasks,
	}

	return c.JSON(http.StatusOK, response)
}

func createTask(c echo.Context) error {
	// TODO: Implement create task handler
	// 1. Bind JSON request to Task struct
	// 2. Validate required fields (title must not be empty)
	// 3. Generate ID and set CreatedAt timestamp
	// 4. Add task to tasks slice
	// 5. Return TaskResponse with created task
	// Hint: Use c.Bind(&task) to bind JSON to struct
	// Hint: Use uuid.New().String() to generate unique ID
	// Hint: Use strings.TrimSpace() to check if title is empty

	var task Task

	// Bind JSON request to task struct
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "Invalid JSON format",
		})
	}

	// Validate required fields
	if strings.TrimSpace(task.Title) == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "title is required",
		})
	}

	// Generate ID and set timestamps
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now().Format(time.RFC3339)

	// Add to tasks slice
	tasks = append(tasks, task)

	response := TaskResponse{
		Status:  "success",
		Message: "Task created successfully",
		Task:    &task,
	}

	return c.JSON(http.StatusCreated, response)
}

func getTask(c echo.Context) error {
	// TODO: Implement get single task handler
	// 1. Get task ID from path parameter
	// 2. Find task in tasks slice
	// 3. Return task if found, or 404 error if not found
	// Hint: Use c.Param("id") to get path parameter
	// Hint: Loop through tasks slice to find matching ID

	id := c.Param("id")

	for _, task := range tasks {
		if task.ID == id {
			response := TaskResponse{
				Status: "success",
				Task:   &task,
			}
			return c.JSON(http.StatusOK, response)
		}
	}

	return c.JSON(http.StatusNotFound, ErrorResponse{
		Status:  "error",
		Message: "Task not found",
	})
}

func updateTask(c echo.Context) error {
	// TODO: Implement update task handler
	// 1. Get task ID from path parameter
	// 2. Find task in tasks slice
	// 3. Bind JSON request to update the task
	// 4. Validate required fields
	// 5. Update task in tasks slice
	// 6. Return updated task
	// Hint: Use pointer to task in slice to modify it directly

	id := c.Param("id")

	// Find task index
	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Status:  "error",
			Message: "Task not found",
		})
	}

	var updatedTask Task
	if err := c.Bind(&updatedTask); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "Invalid JSON format",
		})
	}

	// Validate required fields
	if strings.TrimSpace(updatedTask.Title) == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "title is required",
		})
	}

	// Preserve original ID and CreatedAt
	updatedTask.ID = tasks[taskIndex].ID
	updatedTask.CreatedAt = tasks[taskIndex].CreatedAt

	// Update task in slice
	tasks[taskIndex] = updatedTask

	response := TaskResponse{
		Status:  "success",
		Message: "Task updated successfully",
		Task:    &updatedTask,
	}

	return c.JSON(http.StatusOK, response)
}

func deleteTask(c echo.Context) error {
	// TODO: Implement delete task handler
	// 1. Get task ID from path parameter
	// 2. Find task in tasks slice
	// 3. Remove task from slice
	// 4. Return success message
	// Hint: Use slice operations to remove element: tasks = append(tasks[:i], tasks[i+1:]...)

	id := c.Param("id")

	// Find task index
	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Status:  "error",
			Message: "Task not found",
		})
	}

	// Remove task from slice
	tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)

	response := MessageResponse{
		Status:  "success",
		Message: "Task deleted successfully",
	}

	return c.JSON(http.StatusOK, response)
}
