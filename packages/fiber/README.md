# Fiber Web Development Challenges

Master high-performance web development in Go using the Fiber framework. This package contains 4 progressive challenges that take you from basic HTTP concepts to advanced production-ready patterns with Fiber's Express-inspired API.

## Challenge Overview

### 🎯 [Challenge 1: Basic Routing](./challenge-1-basic-routing/)
**Difficulty:** Beginner | **Duration:** 30-45 minutes

Learn the fundamentals of Fiber by building a simple task management API with basic routing, request handling, and JSON responses.

**Key Skills:**
- Basic Fiber application setup
- Route handlers and HTTP methods
- JSON request/response handling
- Path parameters
- Query parameters

**Topics Covered:**
- `fiber.App` basics
- Route definitions and handlers
- Context handling
- JSON marshaling/unmarshaling
- Error responses

---

### 🚀 [Challenge 2: Middleware & Request/Response Handling](./challenge-2-middleware/)
**Difficulty:** Intermediate | **Duration:** 45-60 minutes

Build an enhanced blog API with comprehensive middleware patterns including logging, authentication, CORS, and rate limiting.

**Key Skills:**
- Custom middleware creation
- Request ID generation and tracking
- Rate limiting implementation
- CORS handling
- Authentication middleware

**Topics Covered:**
- Request/response logging
- API key authentication
- Cross-origin request handling
- Rate limiting per IP
- Centralized error handling

---

### 📦 [Challenge 3: Validation & Error Handling](./challenge-3-validation-errors/)
**Difficulty:** Intermediate | **Duration:** 60-75 minutes

Build a product catalog API with comprehensive input validation, custom validators, and robust error handling.

**Key Skills:**
- Input validation using struct tags
- Custom validator creation
- Bulk operations with partial failures
- Detailed error responses
- Filtering and search functionality

**Topics Covered:**
- Validator package integration
- Custom validation rules
- Error message formatting
- API filtering patterns
- Bulk operation handling

---

### ⚡ [Challenge 4: Authentication & Session Management](./challenge-4-authentication/)
**Difficulty:** Advanced | **Duration:** 75-90 minutes

Build a secure user authentication API with JWT tokens, password hashing, and role-based access control.

**Key Skills:**
- JWT token generation and validation
- Password hashing with bcrypt
- Role-based access control
- Authentication middleware
- Session management

**Topics Covered:**
- User registration and login
- JWT claims and validation
- Password security best practices
- Protected route middleware
- Admin role management

---

## Why Learn Fiber?

**🚀 Performance**: Built on top of Fasthttp for high-performance HTTP handling

**📝 Express-like**: Familiar routing and middleware patterns for JavaScript developers

**🔧 Feature-Rich**: Built-in middleware, validation support, and extensive ecosystem

**🎯 Production Ready**: Used by companies for high-traffic applications

## Learning Path

1. **Start with Challenge 1** if you're new to Fiber or web frameworks in Go
2. **Jump to Challenge 2** if you understand basic HTTP concepts in Go
3. **Challenge 3** focuses on real-world concerns like validation and error handling
4. **Challenge 4** covers advanced features for production applications

## Prerequisites

- **Basic Go knowledge**: Variables, structs, functions, packages
- **HTTP fundamentals**: Understanding of HTTP methods, status codes, headers
- **JSON handling**: Basic familiarity with JSON in Go

## Real-World Applications

These challenges prepare you for building:

- **High-performance REST APIs**
- **Real-time applications** with WebSocket support
- **Microservices** that require fast response times
- **API gateways** handling thousands of requests per second

## Challenge Structure

Each challenge follows a consistent structure:

```
challenge-X-name/
├── README.md              # Challenge description and requirements
├── solution-template.go   # Template with TODOs to implement
├── solution-template_test.go  # Comprehensive test suite
├── run_tests.sh          # Test runner script
├── go.mod                # Go module with dependencies
├── metadata.json         # Challenge metadata
├── SCOREBOARD.md         # Participant scores
├── hints.md              # Implementation hints (when available)
├── learning.md           # Additional learning resources (when available)
└── submissions/          # Participant submission directory
```

## Getting Started

1. **Choose your starting challenge** based on your experience level
2. **Read the README.md** in the challenge directory
3. **Implement the solution** in `solution-template.go`
4. **Test your solution** using `./run_tests.sh`
5. **Submit via PR** to the submissions directory

## Testing Your Solutions

Each challenge includes a comprehensive test suite. To test your solution:

```bash
cd packages/fiber/challenge-X-name/
./run_tests.sh
```

The test script will:
- Prompt for your GitHub username
- Copy your solution to a temporary environment
- Run all tests against your implementation
- Provide detailed feedback on test results

## Common Patterns and Best Practices

### Basic App Setup
```go
app := fiber.New()

app.Get("/", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "message": "Hello, World!",
    })
})

app.Listen(":3000")
```

### Middleware Usage
```go
// Built-in middleware
app.Use(logger.New())
app.Use(cors.New())

// Custom middleware
app.Use(func(c *fiber.Ctx) error {
    // Custom logic
    return c.Next()
})
```

### Error Handling
```go
app.Get("/users/:id", func(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "ID is required",
        })
    }
    
    // Process request
    return c.JSON(response)
})
```

## Resources

- [Official Fiber Documentation](https://docs.gofiber.io/)
- [Fiber GitHub Repository](https://github.com/gofiber/fiber)
- [Fiber Examples](https://github.com/gofiber/recipes)
- [Performance Benchmarks](https://docs.gofiber.io/extra/benchmarks)

---

Ready to build blazing-fast web applications with Fiber? Start with Challenge 1! 🚀