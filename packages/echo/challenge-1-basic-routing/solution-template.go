package main

import (
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

	// TODO: Set up routes by calling setupRoutes function
	// Hint: Call setupRoutes(e) to configure all the routes

	// TODO: Start the server on port 8080
	// Hint: Use e.Logger.Fatal(e.Start(":8080")) to start the server
}

func setupRoutes(e *echo.Echo) {
	// TODO: Implement health check endpoint
	// Hint: Use e.GET("/health", healthHandler) to register the health endpoint

	// TODO: Implement task management endpoints
	// Hint: Register the following routes:
	// - GET /tasks -> getAllTasks
	// - POST /tasks -> createTask
	// - GET /tasks/:id -> getTask
	// - PUT /tasks/:id -> updateTask
	// - DELETE /tasks/:id -> deleteTask
}

func healthHandler(c echo.Context) error {
	// TODO: Implement health check handler
	// Return a HealthResponse with:
	// - Status: "healthy"
	// - Message: "Task API is running"
	// - Timestamp: current time in RFC3339 format
	// - Version: "1.0.0"
	// Hint: Use c.JSON(http.StatusOK, response) to return JSON response
	return nil
}

func getAllTasks(c echo.Context) error {
	// TODO: Implement get all tasks handler with filtering
	// 1. Get query parameters for filtering (priority, completed)
	// 2. Filter tasks based on query parameters
	// 3. Return TaskListResponse with filtered tasks
	// Hint: Use c.QueryParam("priority") to get query parameters
	// Hint: Use strconv.ParseBool() to parse boolean values
	return nil
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
	return nil
}

func getTask(c echo.Context) error {
	// TODO: Implement get single task handler
	// 1. Get task ID from path parameter
	// 2. Find task in tasks slice
	// 3. Return task if found, or 404 error if not found
	// Hint: Use c.Param("id") to get path parameter
	// Hint: Loop through tasks slice to find matching ID
	return nil
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
	return nil
}

func deleteTask(c echo.Context) error {
	// TODO: Implement delete task handler
	// 1. Get task ID from path parameter
	// 2. Find task in tasks slice
	// 3. Remove task from slice
	// 4. Return success message
	// Hint: Use slice operations to remove element: tasks = append(tasks[:i], tasks[i+1:]...)
	return nil
}