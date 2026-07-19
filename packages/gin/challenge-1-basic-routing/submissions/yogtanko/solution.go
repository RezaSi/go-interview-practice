package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

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
var mu sync.RWMutex

func BasicRouting() {
	// TODO: Create Gin router
	router := gin.Default()
	// TODO: Setup routes

	// GET /users - Get all users
	router.GET("/users", getAllUsers)
	// GET /users/search - Search users by name
	router.GET("/users/search", searchUsers)
	// POST /users - Create new user
	router.POST("/users", createUser)
	// GET /users/:id - Get user by ID
	router.GET("/users/:id", getUserByID)
	// PUT /users/:id - Update user
	router.PUT("/users/:id", updateUser)
	// DELETE /users/:id - Delete user
	router.DELETE("/users/:id", deleteUser)

	// TODO: Start server on port 8080
	router.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()
	// TODO: Return all users
	c.JSON(http.StatusOK, Response{
			Success: true,
			Data:    users,
		})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	id := c.Param("id")
	// Handle invalid ID format
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID",
		})
		return
	}
	mu.RLock()
	defer mu.RUnlock()
	u, _ := findUserByID(userID)
	// Return 404 if user not found
	if u == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    u,
	})

}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	user := User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	// Validate required fields
	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	// Add user to storage
	user.ID = nextID
	nextID++
	users = append(users, user)
	// Return created user
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	// Parse JSON request body
	user := User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	if err = validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	// Find and update user
	_, i := findUserByID(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "User not found",
		})
		return
	}
	users[i] = user
	users[i].ID = id
	// Return updated user
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users[i],
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	// Find and remove user
	_, i := findUserByID(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "User not found",
		})
		return
	}
	users = append(users[:i], users[i+1:]...)
	// Return success message
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("User ID = %d Successfully Deleted", id),
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Need param name",
		})
		return
	}
	mu.RLock()
	defer mu.RUnlock()
	// Filter users by name (case-insensitive)
	newUsers := []User{}
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(name)) {
			newUsers = append(newUsers, u)
		}
	}
	// Return matching users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    newUsers,
	})
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
	// Validate email format (basic check)
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}
