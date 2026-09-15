// https://app.gointerview.dev/packages/gin/challenge-1-basic-routing
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
	r.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	response := Response{
		Success: true,
		Data:    users,
		Message: "Successfully",
		Code:    http.StatusOK,
	}
	c.JSON(http.StatusOK, response)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "invalid id",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "user not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "Successfully",
		Code:    http.StatusOK,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "invalid json",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := validateUser(req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "name and email required",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	req.ID = nextID
	nextID++
	users = append(users, req)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    req,
		Message: "Created",
		Code:    http.StatusCreated,
	})

}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "invalid id",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, idx := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "user not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "invalid json",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	if err := validateUser(req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "name and email required",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	req.ID = user.ID
	users[idx] = req

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    req,
		Message: "Updated",
		Code:    http.StatusOK,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "invalid id",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, idx := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "user not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	users = append(users[:idx], users[idx+1:]...)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Successfully delete",
		Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "name query required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	nameLower := strings.ToLower(name)
	matched := []User{} // 非 nil，编码为空数组 []
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), nameLower) {
			matched = append(matched, u)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,    // 有 name 参数时固定为 true
		Data:    matched, // 无匹配为空数组
		Message: "search success",
		Code:    http.StatusOK,
	})

}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	for index := range users {
		if users[index].ID == id {
			return &users[index], index
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	if user.Name == "" || user.Email == "" {
		return errors.New("name and email are required")
	}
	return nil
}
