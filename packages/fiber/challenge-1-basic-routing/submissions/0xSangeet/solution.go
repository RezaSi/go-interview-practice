package main

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v3"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type NewTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[int]*Task
	nextID int
}

func NewTaskStore() *TaskStore {
	store := &TaskStore{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
	store.tasks[1] = &Task{ID: 1, Title: "Learn Go", Description: "Complete Go tutorial", Completed: false}
	store.tasks[2] = &Task{ID: 2, Title: "Build API", Description: "Create REST API with Fiber", Completed: false}
	store.nextID = 3
	return store
}

var taskStore = NewTaskStore()

func main() {
	app := setupApp()
	app.Listen(":3000")
}

func setupApp() *fiber.App {
	app := fiber.New()

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	app.Get("/tasks", func(c fiber.Ctx) error {
		return c.JSON(taskStore.GetAll())
	})

	app.Get("/tasks/:id", func(c fiber.Ctx) error {
        id, err := strconv.Atoi(c.Params("id"))
        if err != nil {
            return c.SendStatus(fiber.StatusBadRequest)
        }
        task, ok := taskStore.GetByID(id)
        if !ok {
            return c.SendStatus(fiber.StatusNotFound)
        }
        return c.JSON(task)
    })

	app.Post("/tasks", func(c fiber.Ctx) error {
		task := new(NewTask)
		if err := c.Bind().Body(task); err == nil { 
			return c.Status(201).JSON(
				taskStore.Create(task.Title, task.Description, task.Completed),
			)
		}
		return c.SendStatus(fiber.StatusBadRequest)
	})

	app.Put("/tasks/:id", func(c fiber.Ctx) error {
        id, err := strconv.Atoi(c.Params("id"))
        if err != nil {
            return c.SendStatus(fiber.StatusBadRequest) 
        }
        task := new(Task)
        if err := c.Bind().Body(task); err == nil {
            if updated, ok := taskStore.Update(id, task.Title, task.Description, task.Completed); ok {
                return c.JSON(updated)
            }
        }
        return c.SendStatus(fiber.StatusNotFound)
    })

	app.Delete("/tasks/:id", func(c fiber.Ctx) error {
        id, err := strconv.Atoi(c.Params("id"))
        if err != nil {
            return c.SendStatus(fiber.StatusBadRequest) 
        }
        if ok := taskStore.Delete(id); ok {
            return c.SendStatus(fiber.StatusNoContent)
        }
        return c.SendStatus(fiber.StatusNotFound)
    })

	return app
}

func (ts *TaskStore) GetAll() []*Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	tasks := make([]*Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (ts *TaskStore) GetByID(id int) (*Task, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	task, exists := ts.tasks[id]
	return task, exists
}

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

func (ts *TaskStore) Delete(id int) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	_, exists := ts.tasks[id]
	if exists {
		delete(ts.tasks, id)
	}
	return exists
}