# Challenge 1: Basic Routing

Build a simple **Task Management API** using Echo with basic HTTP routing and request handling.

## Challenge Requirements

Implement a REST API for managing tasks with the following endpoints:

- `GET /health` - Health check endpoint
- `GET /tasks` - Get all tasks (with optional filtering)
- `POST /tasks` - Create a new task
- `GET /tasks/:id` - Get a specific task by ID
- `PUT /tasks/:id` - Update an existing task
- `DELETE /tasks/:id` - Delete a task

## Data Structure

```go
type Task struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Priority    string `json:"priority"`
    Completed   bool   `json:"completed"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type TaskResponse struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
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

type HealthResponse struct {
    Status    string `json:"status"`
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
    Version   string `json:"version"`
}
```

## Request/Response Examples

**GET /health**
```json
{
    "status": "healthy",
    "message": "Task API is running",
    "timestamp": "2024-01-15T10:30:00Z",
    "version": "1.0.0"
}
```

**POST /tasks** (Request body)
```json
{
    "title": "Learn Echo Framework",
    "description": "Master the basics of Echo web development",
    "priority": "high"
}
```

## Testing Requirements

Your solution must pass tests for:
- Health check returns proper status and format
- Get all tasks returns proper response structure
- Get task by ID returns correct task or 404
- Create task adds new task with generated ID and timestamps
- Update task modifies existing task or returns 404
- Delete task removes task or returns 404
- Filter tasks by priority and completion status
- Proper HTTP status codes and response format for all operations
- Input validation for required fields