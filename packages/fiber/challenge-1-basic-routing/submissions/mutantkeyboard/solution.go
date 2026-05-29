package main

import (
	"sync"
	"strconv"

	"github.com/gofiber/fiber/v3"
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

	// TODO: Start the server on port 3000
	// Hint: Use app.Listen(":3000")
	app.Listen(":3000")
}

// setupApp creates and configures the Fiber app with all routes
func setupApp() *fiber.App {
	// TODO: Create a new Fiber app instance
	app := fiber.New()

	// TODO: Implement health check endpoint
	// GET /ping - should return {"message": "pong"}
	app.Get("/ping", func(c fiber.Ctx) error {
		// TODO: Return JSON response with "pong" message
		
		c.Status(fiber.StatusOK)
	    return c.JSON(fiber.Map{
	        "message" : "pong",
	    })
	})

	// TODO: Implement get all tasks endpoint
	// GET /tasks - should return all tasks as JSON array
	app.Get("/tasks", func(c fiber.Ctx) error {
		// TODO: Get all tasks from store and return as JSON
		// Hint: Use taskStore.GetAll() and c.JSON()
		
        tasks := taskStore.GetAll()
        return c.Status(fiber.StatusOK).JSON(tasks)
        
        

	})

	// TODO: Implement get task by ID endpoint
	// GET /tasks/:id - should return specific task or 404 if not found
	app.Get("/tasks/:id", func(c fiber.Ctx) error {
		// TODO: Extract ID from params, get task from store
		// Return 404 if task not found, otherwise return task as JSON
		// Hint: Use c.Params("id") and strconv.Atoi()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil{
		    return c.SendStatus(fiber.StatusBadRequest)
		}
		task, found := taskStore.GetByID(id)
		if !found {
		    return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusOK).JSON(task)
	})

	// TODO: Implement create task endpoint
	// POST /tasks - should create new task and return it with 201 status
	app.Post("/tasks", func(c fiber.Ctx) error {
		
		task := new(Task)
		if err := c.Bind().Body(task); err != nil {
		    return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.Status(fiber.StatusCreated).JSON(taskStore.Create(task.Title, task.Description, task.Completed))
	})

	// TODO: Implement update task endpoint
	// PUT /tasks/:id - should update existing task or return 404
	app.Put("/tasks/:id", func(c fiber.Ctx) error {
	
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.SendStatus(fiber.StatusBadRequest)
    }

    var input Task
    if err := c.Bind().Body(&input); err != nil {
        return c.SendStatus(fiber.StatusBadRequest)
    }

    updatedTask, found := taskStore.Update(id, input.Title, input.Description, input.Completed)
    if !found {
        return c.SendStatus(fiber.StatusNotFound)
    }

    return c.Status(fiber.StatusOK).JSON(updatedTask)
	})

	// TODO: Implement delete task endpoint
	// DELETE /tasks/:id - should delete task or return 404
	app.Delete("/tasks/:id", func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil{
		    return err
		}
		deleted := taskStore.Delete(id)
		if !deleted {
		    return c.SendStatus(fiber.StatusNotFound)
		}
		
		
		return c.SendStatus(fiber.StatusNoContent)
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
