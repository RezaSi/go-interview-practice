package main

import (
	"errors"
	"net/http"
	"regexp"
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
	r := gin.Default()
	// TODO: Setup routes
	// GET /users - Get all users
	r.GET("/users", getAllUsers)
	// GET /users/:id - Get user by ID
	r.GET("/users/:id", getUserByID)
	// POST /users - Create new user
	r.POST("/users", createUser)
	// PUT /users/:id - Update user
	r.PUT("/users/:id", updateUser)
	// DELETE /users/:id - Delete user
	r.DELETE("/users/:id", deleteUser)
	// GET /users/search - Search users by name
	r.GET("/users/search", searchUsers)
	// TODO: Start server on port 8080
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format",
			Code:    400,
		})
		return
	}
	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "User retrieved successfully!",
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
	// Add user to storage
	user.ID = nextID
	nextID++
	users = append(users, user)
	// Return created user
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
		Message: "User created successfully!",
		Code:    201,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format",
			Code:    400,
		})
		return
	}
	var updatedUser User
	// Parse JSON request body
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
	// Find and update user
	user, idx := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}
	// Validate updated user data
	if err := validateUser(updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
	// Keep original ID
	updatedUser.ID = id
	// Update user
	users[idx] = updatedUser
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    updatedUser,
		Message: "User updated successfully!",
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format",
			Code:    400,
		})
		return
	}
	// Find and remove user
	user, idx := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}
	// Remove user from storage
	users = append(users[:idx], users[idx+1:]...)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully!",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	name := c.Query("name")

	// Validate that name parameter is provided
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "name query parameter is required",
			Code:    400,
		})
		return
	}

	// Filter users by name (case-insensitive partial match)
	userArr := make([]User, 0)
	nameLower := strings.ToLower(name)
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), nameLower) {
			userArr = append(userArr, user)
		}
	}
	// Return matching users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    userArr,
		Message: "Users retrieved successfully!",
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	for i, user := range users {
		if user.ID == id {
			return &users[i], i
		}
	}
	// Return user pointer and index, or nil and -1 if not found
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	if user.Name == "" {
		return errors.New("user name is required")
	}
	// Check required fields: Name, Email
	if user.Email == "" {
		return errors.New("user email is required")
	}
	// Validate email format (basic check)
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	var emailRegex = regexp.MustCompile(emailRegexPattern)
	if len(user.Email) < 3 || len(user.Email) > 254 {
		return errors.New("user email is invalid")
	}
	if !emailRegex.MatchString(user.Email) {
		return errors.New("user email is invalid")
	}
	return nil
}
