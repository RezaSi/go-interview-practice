package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strconv"
	"strings"
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
	router.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	id, err := strconv.Atoi(c.Param("id"))
	// Handle invalid ID format
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "Invalid ID format",
		})
		return
	}
	user, index := findUserByID(id)
	// Return 404 if user not found
	if index == -1 {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "user not found",
		})
	} else {
		c.JSON(200, gin.H{
			"success": true,
			"data":    user,
		})
	}
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var user User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "Invalid data format",
		})
		return
	}
	// Validate required fields
	err = validateUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Add user to storage
	user.ID = nextID
	nextID += 1
	users = append(users, user)

	// Return created user
	c.JSON(201, gin.H{
		"success": true,
		"data":    user,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid ID",
		})

		return
	}

	// Parse JSON request body
	var user User
	err = c.ShouldBind(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid format",
		})
		return
	}

	// Find and update user
	_, idx := findUserByID(id)
	if idx == -1 {
		c.JSON(404, gin.H{
			"success": false,
			"message": "user not found",
		})
		return
	}
	users[idx].Email = user.Email
	users[idx].Name = user.Name
	users[idx].Age = user.Age

	// Return updated user
	c.JSON(200, gin.H{
		"success": true,
		"data":    users[idx],
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid ID",
		})

		return
	}
	// Find and remove user
	_, idx := findUserByID(id)
	if idx == -1 {
		c.JSON(404, gin.H{
			"success": false,
		})
		return
	} else {
		users = slices.Delete(users, idx, idx+1)
	}

	// Return success message
	c.JSON(200, gin.H{
		"success": true,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	name := c.Query("name")
	if name == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "missing params",
		})
		return
	}
	// Filter users by name (case-insensitive)
	filteredUsers := []User{}
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    filteredUsers,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	for idx, item := range users {
		if item.ID == id {
			return &users[idx], idx
		}
	}
	// Return user pointer and index, or nil and -1 if not found
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	if user.Name == "" {
		return errors.New("Name is required")
	}

	if user.Email == "" {
		return errors.New("Email is required")
	}

	if !strings.Contains(user.Email, "@") {
		return errors.New("Invalid email format")
	}

	// Validate email format (basic check)
	return nil
}
