# Challenge 1: Basic Routing

Build a simple **Task Management API** using Gin with basic HTTP routing and request handling.

## Challenge Requirements

Implement a REST API for managing tasks with the following endpoints:

- `GET /ping` - Health check endpoint (returns "pong")
- `GET /tasks` - Get all tasks  
- `GET /tasks/:id` - Get task by ID
- `POST /tasks` - Create new task
- `PUT /tasks/:id` - Update existing task
- `DELETE /tasks/:id` - Delete task

## Data Structure

```go
type Task struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}
```

## Request/Response Examples

**GET /tasks**
```json
[
    {
        "id": 1,
        "title": "Learn Go",
        "description": "Complete Go tutorial",
        "completed": false
    }
]
```

**POST /tasks** (Request body)
```json
{
    "title": "New Task",
    "description": "Task description",
    "completed": false
}
```

## Testing Requirements

Your solution must pass tests for:
- Health check endpoint returns proper response
- Get all tasks returns task array
- Get task by ID returns correct task or 404
- Create task adds new task with auto-incremented ID
- Update task modifies existing task or returns 404
- Delete task removes task or returns 404
- Proper HTTP status codes for all operations 