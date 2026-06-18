# Challenge 1: MongoDB Connection & CRUD Operations

Build a simple **User Management System** using MongoDB with basic CRUD operations and proper error handling.

## Challenge Requirements

Implement a user management system with the following operations:

- **CreateUser** - Create new user with auto-generated ID
- **GetUser** - Retrieve user by ID
- **UpdateUser** - Update existing user (partial updates)
- **DeleteUser** - Remove user from database
- **ListUsers** - Get all users
- **ConnectMongoDB** - Establish database connection

## Data Structure

```go
type User struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name  string             `bson:"name" json:"name"`
    Email string             `bson:"email" json:"email"`
    Age   int                `bson:"age" json:"age"`
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

**CreateUser Request**
```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
}
```

**CreateUser Response**
```json
{
    "success": true,
    "data": {
        "id": "507f1f77bcf86cd799439011",
        "name": "John Doe",
        "email": "john@example.com",
        "age": 30
    },
    "message": "User created successfully",
    "code": 201
}
```

**GetUser Response**
```json
{
    "success": true,
    "data": {
        "id": "507f1f77bcf86cd799439011",
        "name": "John Doe",
        "email": "john@example.com",
        "age": 30
    },
    "code": 200
}
```

**Error Response**
```json
{
    "success": false,
    "error": "User not found",
    "code": 404
}
```

## Testing Requirements

Your solution must pass tests for:
- MongoDB connection establishment and ping test
- User creation with proper validation and ID generation
- User retrieval by ID with proper error handling
- User updates with partial data support
- User deletion with confirmation
- List all users with proper cursor handling
- Proper error handling for invalid IDs and missing users
- Consistent response format for all operations