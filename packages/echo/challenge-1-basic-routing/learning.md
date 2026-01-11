# Learning: Echo Web Framework Fundamentals

## 🌟 **What is Echo?**

Echo is a high-performance, extensible, minimalist Go web framework. It provides a fast HTTP router with zero memory allocation and a rich set of middleware.

### **Why Echo?**
- **Fast**: Zero memory allocation router with optimized performance
- **Middleware support**: Rich ecosystem with easy custom middleware creation
- **Minimalist**: Clean and simple API with minimal boilerplate
- **Standards compliant**: HTTP/2, IPv6, Unix domain sockets support
- **Developer friendly**: Excellent documentation and intuitive design

## 🏗️ **Core Concepts**

### **1. Echo Instance**
The Echo instance is the core of your web application. It handles incoming HTTP requests and routes them to appropriate handlers.

```go
e := echo.New() // Create new Echo instance
// or with configuration
e := echo.New()
e.HideBanner = true
```

### **2. HTTP Methods**
Echo supports all standard HTTP methods:
- **GET**: Retrieve data
- **POST**: Create new resource
- **PUT**: Update entire resource
- **PATCH**: Partial update
- **DELETE**: Remove resource
- **HEAD**: Get headers only
- **OPTIONS**: Check allowed methods

### **3. Context (echo.Context)**
The context carries request data, validates input, and renders responses.

```go
func handler(c echo.Context) error {
    // c contains everything about the HTTP request/response
    return c.JSON(http.StatusOK, data)
}
```

## 📡 **HTTP Request/Response Cycle**

### **Understanding the Flow**
1. **Client** sends HTTP request
2. **Router** matches URL pattern to handler
3. **Handler** processes request and prepares response
4. **Server** sends response back to client

### **Request Components**
- **Method**: GET, POST, PUT, DELETE
- **URL**: `/tasks/123`
- **Headers**: Content-Type, Authorization
- **Body**: JSON, form data, etc.

### **Response Components**
- **Status Code**: 200, 404, 500, etc.
- **Headers**: Content-Type, Cache-Control
- **Body**: JSON, HTML, plain text

## 🛣️ **Routing Patterns**

### **Static Routes**
```go
e.GET("/tasks", getAllTasks)           // Exact match
e.GET("/health", healthCheck)          // Exact match
```

### **Parameter Routes**
```go
e.GET("/tasks/:id", getTaskByID)       // :id captures any value
e.PUT("/tasks/:id", updateTask)        // Same parameter pattern
```

### **Query Parameters**
```go
// URL: /tasks?priority=high&completed=true
priority := c.QueryParam("priority")         // Get query parameter
completed := c.QueryParam("completed")       // Get another parameter
```

### **Route Registration**
```go
func setupRoutes(e *echo.Echo) {
    e.GET("/health", healthHandler)
    e.GET("/tasks", getAllTasks)
    e.POST("/tasks", createTask)
    e.GET("/tasks/:id", getTask)
    e.PUT("/tasks/:id", updateTask)
    e.DELETE("/tasks/:id", deleteTask)
}
```

## 📨 **Request Handling**

### **Reading JSON Data**
```go
type Task struct {
    Title    string `json:"title"`
    Priority string `json:"priority"`
}

func createTask(c echo.Context) error {
    var task Task
    if err := c.Bind(&task); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid JSON format",
        })
    }
    // Process task...
    return c.JSON(http.StatusCreated, task)
}
```

### **Path Parameters**
```go
func getTask(c echo.Context) error {
    id := c.Param("id")  // Get :id from URL
    // Find task by ID...
    return c.JSON(http.StatusOK, task)
}
```

### **Query Parameters**
```go
func getAllTasks(c echo.Context) error {
    priority := c.QueryParam("priority")
    completed := c.QueryParam("completed")
    
    // Filter tasks based on parameters...
    return c.JSON(http.StatusOK, filteredTasks)
}
```

## 📤 **Response Handling**

### **JSON Responses**
```go
// Success response
return c.JSON(http.StatusOK, map[string]interface{}{
    "status": "success",
    "data":   tasks,
})

// Error response
return c.JSON(http.StatusBadRequest, map[string]string{
    "status":  "error",
    "message": "Invalid input",
})
```

### **HTTP Status Codes**
- **200 OK**: Successful GET, PUT
- **201 Created**: Successful POST
- **400 Bad Request**: Invalid input
- **404 Not Found**: Resource doesn't exist
- **500 Internal Server Error**: Server error

## 🔧 **Error Handling**

### **Basic Error Handling**
```go
func handler(c echo.Context) error {
    if someCondition {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Something went wrong",
        })
    }
    return c.JSON(http.StatusOK, data)
}
```

### **Input Validation**
```go
if strings.TrimSpace(task.Title) == "" {
    return c.JSON(http.StatusBadRequest, map[string]string{
        "error": "Title is required",
    })
}
```

## 🏃 **Getting Started**

### **Basic Echo Application**
```go
package main

import (
    "github.com/labstack/echo/v4"
    "net/http"
)

func main() {
    e := echo.New()
    
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

### **Project Structure**
```
project/
├── main.go          # Application entry point
├── handlers.go      # Route handlers
├── models.go        # Data structures
└── go.mod          # Go module file
```

## 🎯 **Best Practices**

### **Handler Organization**
- Keep handlers focused and simple
- Use separate functions for each route
- Group related handlers together

### **Error Responses**
- Use consistent error response format
- Include helpful error messages
- Return appropriate HTTP status codes

### **JSON Structure**
- Use consistent field naming (snake_case or camelCase)
- Include status fields in responses
- Validate input data before processing

## 📚 **Next Steps**

After mastering basic routing, explore:
- **Middleware**: Add logging, authentication, CORS
- **Validation**: Input validation with struct tags
- **Database**: Connect to databases for data persistence
- **Testing**: Write tests for your handlers
- **Deployment**: Deploy your Echo applications

Echo's simplicity and performance make it an excellent choice for building APIs and web services in Go!