# Hints for Challenge 1: Basic Routing with Echo

## Hint 1: Getting Started with Echo

Echo applications start with creating an instance and setting up routes. The basic pattern is:

```go
e := echo.New()
setupRoutes(e)
e.Logger.Fatal(e.Start(":8080"))
```

The `setupRoutes` function is where you'll register all your HTTP endpoints using methods like `e.GET()`, `e.POST()`, etc.

## Hint 2: Understanding Echo Context
The `echo.Context` is your main interface for handling requests and responses. Key methods you'll use:

```go
func handler(c echo.Context) error {
    // Get path parameters: c.Param("id")
    // Get query parameters: c.QueryParam("priority")
    // Bind JSON to struct: c.Bind(&struct)
    // Return JSON response: c.JSON(statusCode, data)
    return nil
}
```

Every handler function must return an `error` - return `nil` for success or an error for failures.

## Hint 3: Route Registration Pattern
Register routes in the `setupRoutes` function using this pattern:

```go
func setupRoutes(e *echo.Echo) {
    e.GET("/health", healthHandler)
    e.GET("/tasks", getAllTasks)
    e.POST("/tasks", createTask)
    e.GET("/tasks/:id", getTask)     // :id is a path parameter
    e.PUT("/tasks/:id", updateTask)
    e.DELETE("/tasks/:id", deleteTask)
}
```

The `:id` syntax creates a path parameter that you can access with `c.Param("id")`.


## Hint 4: JSON Request Binding
To read JSON from the request body, use the `Bind` method:

```go
var task Task
if err := c.Bind(&task); err != nil {
    return c.JSON(http.StatusBadRequest, ErrorResponse{
        Status:  "error",
        Message: "Invalid JSON format",
    })
}
```

Always check for binding errors and return appropriate error responses.


## Hint 5: Generating UUIDs and Timestamps
For creating new tasks, you'll need to generate unique IDs and timestamps:

```go
import (
    "time"
    "github.com/google/uuid"
)

// Generate unique ID
task.ID = uuid.New().String()

// Set current timestamp in RFC3339 format
task.CreatedAt = time.Now().Format(time.RFC3339)
```

The RFC3339 format is the standard for JSON timestamps.


## Hint 6: Input Validation
Validate required fields before processing:

```go
import "strings"

// Check if title is empty after trimming whitespace
if strings.TrimSpace(task.Title) == "" {
    return c.JSON(http.StatusBadRequest, ErrorResponse{
        Status:  "error",
        Message: "title is required",
    })
}
```

Always trim whitespace when checking for empty strings to handle spaces-only input.


## Hint 7: Query Parameter Filtering
Handle optional query parameters for filtering:

```go
priority := c.QueryParam("priority")
completedStr := c.QueryParam("completed")

// Convert string to boolean for completed filter
var completed bool
var hasCompletedFilter bool
if completedStr != "" {
    var err error
    completed, err = strconv.ParseBool(completedStr)
    if err == nil {
        hasCompletedFilter = true
    }
}

// Apply filters in your loop
for _, task := range tasks {
    if priority != "" && task.Priority != priority {
        continue // Skip this task
    }
    if hasCompletedFilter && task.Completed != completed {
        continue // Skip this task
    }
    // Add to filtered results
}
```


## Hint 8: Finding Tasks by ID
When working with task IDs, you'll need to search through your tasks slice:

```go
id := c.Param("id")

// Find task by ID
var foundTask *Task
taskIndex := -1
for i, task := range tasks {
    if task.ID == id {
        foundTask = &task
        taskIndex = i
        break
    }
}

if foundTask == nil {
    return c.JSON(http.StatusNotFound, ErrorResponse{
        Status:  "error",
        Message: "Task not found",
    })
}
```

Store the index when you need to modify or delete the task later.


## Hint 9: Updating Tasks Correctly
When updating tasks, preserve the original ID and CreatedAt:

```go
var updatedTask Task
if err := c.Bind(&updatedTask); err != nil {
    return c.JSON(http.StatusBadRequest, ErrorResponse{
        Status:  "error", 
        Message: "Invalid JSON format",
    })
}

// Preserve original values
updatedTask.ID = tasks[taskIndex].ID
updatedTask.CreatedAt = tasks[taskIndex].CreatedAt

// Replace the task in the slice
tasks[taskIndex] = updatedTask
```

This ensures the task maintains its identity and creation timestamp.


## Hint 10: Deleting from Slices
To remove a task from the slice, use slice operations:

```go
// Remove task at taskIndex
tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)
```

This creates a new slice without the element at `taskIndex`. The `[:taskIndex]` gets elements before the index, and `[taskIndex+1:]` gets elements after the index.


## 🔍 **Debugging Tips**

### **Common Issues and Solutions**

**Issue: Tests failing with "route not found"**
- Check that all routes are registered in `setupRoutes`
- Verify the exact path matches what the tests expect
- Ensure HTTP methods (GET, POST, etc.) are correct

**Issue: JSON binding not working**
- Verify struct tags match JSON field names exactly
- Check that the request Content-Type is "application/json"
- Ensure the JSON is valid (use a JSON validator)

**Issue: Path parameters returning empty strings**
- Use `c.Param("id")` for path parameters like `/tasks/:id`
- Use `c.QueryParam("name")` for query parameters like `?name=value`
- Check that the route pattern matches your handler registration

**Issue: Tasks not persisting between requests**
- Remember that the `tasks` slice is in-memory storage
- Each test may start with a fresh state
- Make sure you're appending to the global `tasks` slice

### **Testing Your Implementation**
Run tests frequently to catch issues early:

```bash
./run_tests.sh
```

If tests fail, read the error messages carefully - they often tell you exactly what's wrong.

### **Validation Checklist**
Before running tests, verify:
- [ ] All routes are registered with correct HTTP methods
- [ ] Handler functions return `error` (not other types)
- [ ] JSON responses use `c.JSON(statusCode, data)`
- [ ] Path parameters use `c.Param("name")`
- [ ] Query parameters use `c.QueryParam("name")`
- [ ] Required field validation is implemented
- [ ] Error responses follow the expected format
- [ ] UUIDs are generated for new task IDs
- [ ] Timestamps are in RFC3339 format

Good luck! Take it step by step, and don't hesitate to refer back to the learning materials if you get stuck.