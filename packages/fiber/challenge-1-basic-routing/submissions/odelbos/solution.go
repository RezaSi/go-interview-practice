package main

import (
	"sync"
	"net/http"
	"strconv"

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
	app.Listen(":3000")
}

// setupApp creates and configures the Fiber app with all routes
func setupApp() *fiber.App {
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	app.Get("/tasks", func(c *fiber.Ctx) error {
		return c.JSON(taskStore.GetAll())
	})

	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid ID"})
		}
		task, ok := taskStore.GetByID(id)
		if ! ok {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"msg": "Not found"})
		}
		return c.JSON(task)
	})

	app.Post("/tasks", func(c *fiber.Ctx) error {
		var data Task
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON"})
		}
		task := taskStore.Create(data.Title, data.Description, data.Completed)
		return c.Status(http.StatusCreated).JSON(task)
	})

	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid ID"})
		}

		var data Task
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON"})
		}
		task, ok := taskStore.Update(id, data.Title, data.Description, data.Completed)
		if ! ok {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"msg": "Task not found"})
		}
		return c.JSON(task)
	})

	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid ID"})
		}

		if ! taskStore.Delete(id) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"msg": "Not found"})
		}
		return c.SendStatus(http.StatusNoContent)
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
