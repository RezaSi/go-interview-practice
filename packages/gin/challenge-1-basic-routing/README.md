# Challenge 1: Basic Routing

Build a simple **User Management API** using Gin with basic HTTP routing and request handling.

## Challenge Requirements

Implement a REST API for managing users with the following endpoints:

- `GET /users` - Get all users
- `GET /users/:id` - Get user by ID
- `POST /users` - Create new user
- `PUT /users/:id` - Update existing user
- `DELETE /users/:id` - Delete user
- `GET /users/search` - Search users by name

## Data Structure

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message,omitempty"`
    Error   string      `json:"error,omitempty"`
    Code    int         `json:"code,omitempty"`
}
```

## Request/Response Examples

**GET /users**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "age": 30
        }
    ]
}
```

**POST /users** (Request body)
```json
{
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "age": 28
}
```

## Testing Requirements

Your solution must pass tests for:
- Get all users returns proper response structure
- Get user by ID returns correct user or 404
- Create user adds new user with auto-incremented ID
- Update user modifies existing user or returns 404
- Delete user removes user or returns 404
- Search users by name (case-insensitive)
- Proper HTTP status codes and response format for all operations 