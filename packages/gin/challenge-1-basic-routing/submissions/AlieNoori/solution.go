package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id" `
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

var (
	nextID     = 4
	mu         = &sync.Mutex{}
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
)

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
	if err := router.Run(); err != nil {
		panic(err)
	}
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	response := &Response{
		Success: true,
		Data:    users,
		Code:    http.StatusOK,
	}

	c.JSON(http.StatusOK, response)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	userID, err := getUserIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "invalid user ID format",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, idx := findUserByID(userID)

	if idx == -1 {
		c.JSON(http.StatusNotFound, &Response{
			Success: false,
			Message: fmt.Sprintf("user with id %d not found", userID),
			Error:   "not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Data:    user,
		Success: true,
		Code:    http.StatusOK,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var user User
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Message: "invalid json body",
			Error:   "bad request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate required fields
	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Message: "invalid user",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Add user to storage
	mu.Lock()
	user.ID = nextID
	nextID++
	mu.Unlock()
	users = append(users, user)
	// Return created user

	c.JSON(http.StatusCreated, &Response{
		Success: true,
		Data:    user,
		Code:    http.StatusCreated,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	userID, err := getUserIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "invalid user ID format",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Parse JSON request body
	var user User
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Message: "invalid json body",
			Error:   "bad request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Message: "invalid user",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Find and update user
	_, idx := findUserByID(userID)
	if idx == -1 {
		c.JSON(http.StatusNotFound, &Response{
			Success: false,
			Message: fmt.Sprintf("user with id %d not found", userID),
			Error:   "not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	user.ID = userID
	users[idx] = user

	// Return updated user
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    user,
		Code:    http.StatusOK,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	userID, err := getUserIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "invalid user ID format",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Find and remove user
	user, idx := findUserByID(userID)
	if idx == -1 {
		c.JSON(http.StatusNotFound, &Response{
			Success: false,
			Message: fmt.Sprintf("user with id %d not found", userID),
			Error:   "not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	users = append(users[:idx], users[idx+1:]...)

	// Return success message
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    user,
		Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: fmt.Sprintf("name parameter is required"),
			Code:    http.StatusBadRequest,
			Error:   "bad request",
		})
		return
	}

	users := findUsersByName(name)
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    users,
		Code:    http.StatusOK,
	})
}

// getUserIDFromParam retrieve userID from path
func getUserIDFromParam(c *gin.Context) (int, error) {
	userIDStr := c.Param("id")
	return strconv.Atoi(userIDStr)
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	var user *User
	idx := -1

	for i, u := range users {
		if u.ID == id {
			idx = i
			user = &u
			break
		}
	}

	return user, idx
}

// Helper function to find user by Name
func findUsersByName(name string) []*User {
	usrs := []*User{}

	for _, u := range users {
		userName := strings.SplitN(u.Name, " ", 2)
		firstName := userName[0]
		if strings.ToLower(firstName) == strings.ToLower(name) {
			usrs = append(usrs, &u)
		}
	}

	return usrs
}

// Helper function to validate user data
func validateUser(user User) error {
	// Check required fields: Name, Email
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}

	if user.Email == "" {
		return fmt.Errorf("eamil is required")
	}

	// Validate email format (basic check)
	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
