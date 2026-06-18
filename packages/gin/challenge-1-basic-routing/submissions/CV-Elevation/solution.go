package main

import (
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
	{ID: 4, Name: "kuoz", Email: "kuoz@qq.com", Age: 21},
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
	r.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Code:    http.StatusOK,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "无效的id",
		})
		return
	}
	// Handle invalid ID format
	// Return 404 if user not found
	user, code := findUserByID(id)
	if code == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Code:    http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})

}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var user User
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
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Add user to storage
	user.ID = len(users) + 1
	users = append(users, user)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
		Code:    http.StatusCreated,
	})
	// Return created user
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	// Parse JSON request body
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	err = validateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	user.ID = id
	// Find and update user
	for index, userEx := range users {
		if userEx.ID == id {
			users[index] = user
			c.JSON(http.StatusOK, Response{
				Success: true,
				Data:    user,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   "user not found",
	})
	// Return updated user
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	// Find and remove user
	for index, user := range users {
		if user.ID == id {
			users[index] = users[len(users)-1]
			users = users[:len(users)-1]
			c.JSON(http.StatusOK, Response{
				Success: true,
				Data:    user,
			})
			return
		}
	}
	// Return success message
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   "User Not Found",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "name is nil",
			Code:    http.StatusBadRequest,
		})
		return
	}
	name = strings.ToLower(name)
	// Filter users by name (case-insensitive)
	resUsers := make([]User, 0)
	for _, user := range users {
		lowName := strings.ToLower(user.Name)
		if strings.Contains(lowName, name) {
			resUsers = append(resUsers, user)
		}
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    resUsers,
		Code:    http.StatusOK,
	})
	// Return matching users
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	for _, user := range users {
		if user.ID == id {
			return &user, user.ID
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	if strings.TrimSpace(user.Name) == "" {
		return &gin.Error{Err: http.ErrBodyNotAllowed, Type: gin.ErrorTypeBind}
	}
	if strings.TrimSpace(user.Email) == "" {
		return &gin.Error{Err: http.ErrBodyNotAllowed, Type: gin.ErrorTypeBind}
	}
	if user.Age <= 0 {
		return &gin.Error{Err: http.ErrBodyNotAllowed, Type: gin.ErrorTypeBind}
	}
	if !strings.Contains(user.Email, "@") {
		return &gin.Error{Err: http.ErrBodyNotAllowed, Type: gin.ErrorTypeBind}
	}
	return nil
}
