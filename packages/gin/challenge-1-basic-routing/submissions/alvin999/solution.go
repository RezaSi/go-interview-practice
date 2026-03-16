package main

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// In-memory storage
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var nextID = 4

func main() {
	// TODO: Create Gin router
    router := gin.Default()
	// TODO: Setup routes
	// GET /users - Get all users
	router.GET("/users", getAllUsers)
	// GET /users/:id - Get user by ID
	router.GET("/users/:id", getUserByID)
	// POST /users - Create new user
	router.POST("/users", createUser)
	// PUT /users/:id - Update user
	router.PUT("/users/:id", updateUser)
	// DELETE /users/:id - Delete user
	router.DELETE("/users/:id", deleteUser)
	// GET /users/search - Search users by name
    router.GET("/users/search", searchUsers)
	// TODO: Start server on port 8080
	router.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(http.StatusOK, Response{
	    Success:    true,
	    Data:       users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	idStr := c.Param("id")
	// Handle invalid ID format
	id, err := strconv.Atoi(idStr)
	if err != nil {
	    c.JSON(http.StatusBadRequest, Response{Success: false, Error: "Invalid ID format"})
	    return
	}
	// Return 404 if user not found
	user, _ := findUserByID(id)
	if user == nil {
	    c.JSON(http.StatusNotFound, Response{Success: false, Error: "User not found"})
	    return
	}
	
	c.JSON(http.StatusOK, Response{Success: true, Data: user})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
    var newUser User
	// TODO: Parse JSON request body
	// Validate required fields
	if err := c.ShouldBindJSON(&newUser); err != nil {
	    c.JSON(http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
	    return
	}
	if err := validateUser(newUser); err != nil {
	    c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error()})
	    return
	}
	// Add user to storage
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	// Return created user
	
	c.JSON(http.StatusCreated, Response{Success: true, Data: newUser})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, _ := strconv.Atoi(c.Param("id"))
	user, index := findUserByID(id)
	// Parse JSON request body
	if user == nil {
	    c.JSON(http.StatusNotFound, Response{Success: false, Error: "User not found"})
	    return
	}
	// Find and update user
	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
	    c.JSON(http.StatusBadRequest, Response{Success: false, Error: "Invalid body"})
	    return
	}
	// Return updated user
	updateData.ID = id
	users[index] = updateData
	
	c.JSON(http.StatusOK, Response{Success: true, Data: updateData})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, _ := strconv.Atoi(c.Param("id"))
	_, index := findUserByID(id)
	if index == -1 {
	    c.JSON(http.StatusNotFound, Response{Success: false, Error: "User not found"})
	    return
	}
	// Find and remove user
	users = append(users[:index], users[index+1:]...)
	// Return success message
	c.JSON(http.StatusOK, Response{Success: true, Message: "Userf deleted successfully"})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	nameQuery := c.Query("name")
	
	if nameQuery == "" {
	    c.JSON(http.StatusBadRequest, Response{
	        Success: false,
	        Error: "Query parameter 'name' is required",
	    })
	    return
	}
	// Filter users by name (case-insensitive)
	result := []User{}
	
	queryLower := strings.ToLower(nameQuery)
	for _, u := range users {
	    if strings.Contains(strings.ToLower(u.Name), queryLower) {
	        result = append(result, u)
	    }
	}
	// Return matching users
	c.JSON(http.StatusOK, Response{Success: true, Data: result})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	for i, u := range users {
	    if u.ID == id {
	        return &users[i], i
	    }
	}
	// Return user pointer and index, or nil and -1 if not found
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	if user.Name == "" || user.Email == "" {
	    return fmt.Errorf("name and email are required")
	}
	// Validate email format (basic check)
	if !strings.Contains(user.Email, "@") {
	    return fmt.Errorf("invalid email format")
	}
	return nil
}
