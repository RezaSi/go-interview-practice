# Implementation Hints

## Getting Started

1. **Create Fiber App**
   ```go
   app := fiber.New()
   ```

2. **Basic Route Structure**
   ```go
   app.Get("/path", func(c *fiber.Ctx) error {
       // Your handler logic
       return c.JSON(response)
   })
   ```

## Step-by-Step Implementation

### 1. Health Check Endpoint (`GET /ping`)
```go
app.Get("/ping", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "message": "pong",
    })
})
```

### 2. Get All Tasks (`GET /tasks`)
```go
app.Get("/tasks", func(c *fiber.Ctx) error {
    tasks := taskStore.GetAll()
    return c.JSON(tasks)
})
```

### 3. Get Task by ID (`GET /tasks/:id`)
```go
app.Get("/tasks/:id", func(c *fiber.Ctx) error {
    // Extract ID from URL parameter
    idStr := c.Params("id")
    
    // Convert string to integer
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid task ID",
        })
    }
    
    // Get task from store
    task, exists := taskStore.GetByID(id)
    if !exists {
        return c.Status(404).JSON(fiber.Map{
            "error": "Task not found",
        })
    }
    
    return c.JSON(task)
})
```

### 4. Create Task (`POST /tasks`)
```go
app.Post("/tasks", func(c *fiber.Ctx) error {
    var newTask Task
    
    // Parse JSON body into struct
    if err := c.BodyParser(&newTask); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid JSON",
        })
    }
    
    // Create task in store
    task := taskStore.Create(newTask.Title, newTask.Description, newTask.Completed)
    
    // Return with 201 Created status
    return c.Status(201).JSON(task)
})
```

### 5. Update Task (`PUT /tasks/:id`)
```go
app.Put("/tasks/:id", func(c *fiber.Ctx) error {
    // Extract and validate ID
    idStr := c.Params("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid task ID",
        })
    }
    
    // Parse update data
    var updateTask Task
    if err := c.BodyParser(&updateTask); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid JSON",
        })
    }
    
    // Update in store
    task, exists := taskStore.Update(id, updateTask.Title, updateTask.Description, updateTask.Completed)
    if !exists {
        return c.Status(404).JSON(fiber.Map{
            "error": "Task not found",
        })
    }
    
    return c.JSON(task)
})
```

### 6. Delete Task (`DELETE /tasks/:id`)
```go
app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
    // Extract and validate ID
    idStr := c.Params("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid task ID",
        })
    }
    
    // Delete from store
    if !taskStore.Delete(id) {
        return c.Status(404).JSON(fiber.Map{
            "error": "Task not found",
        })
    }
    
    // Return 204 No Content
    return c.SendStatus(204)
})
```

### 7. Start the Server
```go
// Start server on port 3000
app.Listen(":3000")
```

## Key Fiber Methods

- **`c.Params(key)`** - Extract URL parameters
- **`c.BodyParser(&struct)`** - Parse JSON request body
- **`c.JSON(data)`** - Send JSON response
- **`c.Status(code)`** - Set HTTP status code
- **`c.SendStatus(code)`** - Send status code only
- **`fiber.Map{}`** - Create JSON object quickly

## Common Patterns

### Error Response
```go
return c.Status(400).JSON(fiber.Map{
    "error": "Error message",
})
```

### Success Response with Data
```go
return c.JSON(data)
```

### Created Response
```go
return c.Status(201).JSON(createdData)
```

### No Content Response
```go
return c.SendStatus(204)
```

## Testing Tips

1. Use the provided test file to verify your implementation
2. Run tests with: `go test -v`
3. Check that all HTTP status codes are correct
4. Ensure JSON responses match expected format

## Troubleshooting

- **Import Error**: Make sure `go.mod` is set up correctly with Fiber dependency
- **JSON Parsing**: Ensure request Content-Type is `application/json`
- **Status Codes**: Remember to set appropriate codes (200, 201, 404, etc.)
- **Concurrency**: The provided TaskStore handles thread safety for you