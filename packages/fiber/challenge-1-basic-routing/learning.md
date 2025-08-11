# Learning: Fiber Web Framework Fundamentals

## 🌟 **What is Fiber?**

Fiber is an Express.js inspired web framework built on top of Fasthttp, the fastest HTTP engine for Go. Fiber is designed to ease things up for fast development with zero memory allocation and performance in mind.

### **Why Fiber?**
- **Fast**: Built on Fasthttp, one of the fastest HTTP engines
- **Low Memory**: Zero memory allocation router  
- **Express-like**: If you know Express.js, you already know Fiber
- **Middleware Rich**: 40+ middleware packages available
- **Developer Friendly**: Simple routing, static files, and template engines

## 🏗️ **Core Concepts**

### **1. App Instance**
The app instance is the core of a Fiber application. It handles incoming HTTP requests and routes them to appropriate handlers.

```go
app := fiber.New() // Create new Fiber instance
// or with config
app := fiber.New(fiber.Config{
    Prefork: true,
    CaseSensitive: true,
})
```

### **2. HTTP Methods**
Fiber supports all standard HTTP methods with Express-like syntax:
- **GET**: Retrieve data
- **POST**: Create new resource
- **PUT**: Update entire resource
- **PATCH**: Partial update
- **DELETE**: Remove resource
- **HEAD**: Get headers only
- **OPTIONS**: Check allowed methods

### **3. Context (fiber.Ctx)**
The context carries request data, validates JSON, and renders responses.

```go
func handler(c *fiber.Ctx) error {
    // c contains everything about the HTTP request/response
    return c.JSON(fiber.Map{"message": "Hello"})
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
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
})
```

### **Route Parameters**
```go
app.Get("/tasks/:id", func(c *fiber.Ctx) error {
    id := c.Params("id")
    return c.JSON(fiber.Map{"task_id": id})
})
```

### **Query Parameters**
```go
app.Get("/tasks", func(c *fiber.Ctx) error {
    search := c.Query("search", "")
    return c.JSON(fiber.Map{"search": search})
})
```

## 📝 **Request/Response Handling**

### **JSON Responses**
```go
app.Get("/api/data", func(c *fiber.Ctx) error {
    data := map[string]interface{}{
        "message": "Success",
        "data": []string{"item1", "item2"},
    }
    return c.JSON(data)
})
```

### **JSON Parsing**
```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

app.Post("/users", func(c *fiber.Ctx) error {
    user := new(User)
    if err := c.BodyParser(user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(user)
})
```

## 🔧 **Essential Methods**

### **Context Methods**
- `c.Params(key)` - Get route parameter
- `c.Query(key)` - Get query parameter  
- `c.BodyParser(&struct)` - Parse request body
- `c.JSON(data)` - Send JSON response
- `c.Status(code)` - Set status code
- `c.SendString(text)` - Send plain text

### **Response Helpers**
```go
// Status codes
c.Status(fiber.StatusOK)        // 200
c.Status(fiber.StatusNotFound)  // 404
c.Status(fiber.StatusCreated)   // 201

// JSON responses
c.JSON(fiber.Map{"key": "value"})
c.Status(400).JSON(fiber.Map{"error": "Bad request"})
```

## 🚀 **Performance Benefits**

### **Memory Efficiency**
- Zero memory allocation router
- Fast HTTP parsing
- Low memory footprint

### **Framework Features**
- **Built on Fasthttp**: High-performance HTTP engine for Go
- **Express-like API**: Familiar developer experience for JavaScript developers
- **Rich Ecosystem**: Extensive middleware and community support

## 📚 **Best Practices**

1. **Error Handling**: Always return errors from handlers
2. **Status Codes**: Use appropriate HTTP status codes
3. **JSON Validation**: Parse and validate request bodies
4. **Context Usage**: Use context for request/response operations
5. **Route Organization**: Group related routes together

## 🔗 **Framework Characteristics**

### **Fiber Key Features**
- **API Style**: Express.js inspired syntax and patterns
- **Learning Curve**: Easy for developers familiar with Express.js
- **Memory Usage**: Optimized for low memory footprint
- **Middleware**: Rich ecosystem of built-in and third-party middleware

## 🎯 **Next Steps**

After mastering basic routing, you'll learn:
1. **Middleware** - Request/response processing
2. **Validation** - Input validation and error handling
3. **Authentication** - JWT tokens and security
4. **Advanced Features** - Production-ready patterns