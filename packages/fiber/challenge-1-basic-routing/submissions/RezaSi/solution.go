package main

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// Task represents a task in our task management system
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// TaskStore manages our in-memory task storage
type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[int]*Task
	nextID int
}

// NewTaskStore creates a new task store
func NewTaskStore() *TaskStore {
	store := &TaskStore{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}

	// Add some sample tasks
	store.tasks[1] = &Task{ID: 1, Title: "Learn Go", Description: "Complete Go tutorial", Completed: false}
	store.tasks[2] = &Task{ID: 2, Title: "Build API", Description: "Create REST API with Fiber", Completed: false}
	store.nextID = 3

	return store
}

// Global task store
var taskStore = NewTaskStore()

func main() {
	app := setupApp()

	// Start the server on port 3000
	app.Listen(":3000")
}

// setupApp creates and configures the Fiber app with all routes
func setupApp() *fiber.App {
	// Create a new Fiber app instance
	app := fiber.New()

	// Health check endpoint
	// GET /ping - should return {"message": "pong"}
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	// Get all tasks endpoint
	// GET /tasks - should return all tasks as JSON array
	app.Get("/tasks", func(c *fiber.Ctx) error {
		tasks := taskStore.GetAll()
		return c.JSON(tasks)
	})

	// Get task by ID endpoint
	// GET /tasks/:id - should return specific task or 404 if not found
	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		task, exists := taskStore.GetByID(id)
		if !exists {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}

		return c.JSON(task)
	})

	// Create task endpoint
	// POST /tasks - should create new task and return it with 201 status
	app.Post("/tasks", func(c *fiber.Ctx) error {
		var newTask Task
		if err := c.BodyParser(&newTask); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		task := taskStore.Create(newTask.Title, newTask.Description, newTask.Completed)
		return c.Status(201).JSON(task)
	})

	// Update task endpoint
	// PUT /tasks/:id - should update existing task or return 404
	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		var updateTask Task
		if err := c.BodyParser(&updateTask); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		task, exists := taskStore.Update(id, updateTask.Title, updateTask.Description, updateTask.Completed)
		if !exists {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}

		return c.JSON(task)
	})

	// Delete task endpoint
	// DELETE /tasks/:id - should delete task or return 404
	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		if !taskStore.Delete(id) {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}

		return c.SendStatus(204)
	})

	return app
}

// Helper methods for TaskStore

// GetAll returns all tasks
func (ts *TaskStore) GetAll() []*Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	tasks := make([]*Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetByID returns a task by ID
func (ts *TaskStore) GetByID(id int) (*Task, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	task, exists := ts.tasks[id]
	return task, exists
}

// Create adds a new task and returns it
func (ts *TaskStore) Create(title, description string, completed bool) *Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task := &Task{
		ID:          ts.nextID,
		Title:       title,
		Description: description,
		Completed:   completed,
	}

	ts.tasks[ts.nextID] = task
	ts.nextID++

	return task
}

// Update modifies an existing task
func (ts *TaskStore) Update(id int, title, description string, completed bool) (*Task, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return nil, false
	}

	task.Title = title
	task.Description = description
	task.Completed = completed

	return task, true
}

// Delete removes a task by ID
func (ts *TaskStore) Delete(id int) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	_, exists := ts.tasks[id]
	if exists {
		delete(ts.tasks, id)
	}

	return exists
}
