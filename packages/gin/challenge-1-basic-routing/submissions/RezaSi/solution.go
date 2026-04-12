package main

import (
	"errors"
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
	// Create Gin router
	router := gin.Default()

	// Setup routes
	router.GET("/users/search", searchUsers) // Specific route first
	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	// Start server on port 8080
	router.Run(":8080")
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Message: "Users retrieved successfully",
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "User retrieved successfully",
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate user data
	if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Assign ID and add to storage
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "User created successfully",
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate user data
	if err := validateUser(updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Find user and update
	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Keep the original ID
	updatedUser.ID = id
	users[index] = updatedUser

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    updatedUser,
		Message: "User updated successfully",
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Find user and remove
	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Remove user from slice
	users = append(users[:index], users[index+1:]...)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Name parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	results := make([]User, 0)
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			results = append(results, user)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    results,
		Message: "Search completed successfully",
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for i, user := range users {
		if user.ID == id {
			return &user, i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
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
