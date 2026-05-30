package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func main() {
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
	if len(users) > 0 {
		c.JSON(200, Response{
			Success: true,
			Data:    users,
		})
	} else {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	id := c.Param("id")
	// Handle invalid ID format
	userID, err := strconv.Atoi(id)
	// Return 404 if user not found
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   "Invalid ID",
		})
		return
	}
	mu.RLock()
	defer mu.RUnlock()
	var u User
	for _, user := range users {
		if userID == user.ID {
			u = user
		}
	}
	if u != (User{}) {
		c.JSON(200, Response{
			Success: true,
			Data:    u,
		})
		return
	}
	c.JSON(404, Response{
		Success: false,
		Error:   "User not found",
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	user := User{}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&user)
	if err != nil {
		c.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	// Validate required fields
	err = validateUser(user)
	if err != nil {
		c.JSON(400, Response{
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
	c.JSON(201, Response{
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
		c.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	// Parse JSON request body
	user := &User{}
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(user)
	if err != nil {
		c.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	if err = validateUser(*user); err != nil {
		c.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	// Find and update user
	u, i := findUserByID(id)
	if i == -1 {
		c.JSON(404, Response{
			Success: false,
			Message: "User not found",
		})
		return
	}
	users[i] = *user
	users[i].ID = id
	// Return updated user
	c.JSON(200, Response{
		Success: true,
		Data:    u,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(404, Response{
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
		c.JSON(404, Response{
			Success: false,
			Message: "User not found",
		})
		return
	}
	users = append(users[:i], users[i+1:]...)
	// Return success message
	c.JSON(200, Response{
		Success: true,
		Message: fmt.Sprintf("User ID = %d Succesfuly Deleted", id),
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	name := c.Query("name")
	if name == "" {
		c.JSON(400, Response{
			Success: false,
			Error:   "Need param name",
		})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	// Filter users by name (case-insensitive)
	newUsers := []User{}
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(name)) {
			newUsers = append(newUsers, u)
		}
	}
	// Return matching users
	c.JSON(200, Response{
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
