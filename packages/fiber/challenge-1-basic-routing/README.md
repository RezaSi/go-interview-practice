# Challenge 1: Basic Routing

Build a simple **Task Management API** using Fiber with basic HTTP routing and request handling.

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

Your implementation must pass all the provided tests, which verify:

- ✅ Correct HTTP methods and routes
- ✅ Proper JSON request/response handling
- ✅ Path parameter extraction
- ✅ In-memory data persistence
- ✅ Error handling for invalid requests
- ✅ Appropriate HTTP status codes

## Implementation Notes

- Use Fiber's built-in JSON handling with `c.JSON()`
- Extract path parameters with `c.Params()`
- Parse JSON request bodies with `c.BodyParser()`
- Store tasks in memory (slice or map)
- Return appropriate HTTP status codes

## Getting Started

1. Examine the `solution-template.go` file
2. Implement the TODO sections
3. Run tests with `./run_tests.sh`
4. Iterate until all tests pass

## Success Criteria

- All tests pass
- Clean, readable code
- Proper error handling
- RESTful API design
- Efficient in-memory storage