package main

import (
	"errors"
	"net/http"
	"net/mail"
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
	r.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(200, Response{
		Success: true,
		Data:    users,
	})

}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	idParam := c.Param("id")
	// TODO: Get user by ID
	id, err := strconv.Atoi(idParam)
	// Handle invalid ID format
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid user ID",
			Code:    400,
		})
		return
	}

	user, _ := findUserByID(id)
	// Return 404 if user not found
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
			Code:    404,
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request body",
			Code:    400,
		})
		return
	}

	// Validate required fields
	if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "illegal request",
			Code:    400,
		})
		return
	}

	newUser.ID = nextID
	nextID++

	// Add user to storage
	users = append(users, newUser)

	// Return created user
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "user created successfully",
		Code:    201,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	idParam := c.Param("id")

	// Parse JSON request body
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid user ID",
			Code:    400,
		})
		return
	}

	user, index := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
			Code:    404,
		})
		return
	}

	// Find and update user
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request body",
			Code:    400,
		})
		return
	}

	if updatedUser.Name == "" || updatedUser.Email == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "name and email are required",
			Code:    400,
		})
		return
	}

	updatedUser.ID = id
	users[index] = updatedUser

	// Return updated user
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    updatedUser,
		Message: "user updated successfully",
		Code:    200,
	})

}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid user ID",
			Code:    400,
		})
		return
	}

	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
			Code:    404,
		})
		return
	}

	users = append(users[:index], users[index+1:]...)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "user deleted successfully",
		Code:    200,
	})

	// Find and remove user
	// Return success message
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	nameQuery := c.Query("name")

	if nameQuery == "" {
		c.JSON(400, Response{
			Success: false,
			Data:    []User{},
		})
		return
	}

	// Filter users by name (case-insensitive)
	nameQuery = strings.ToLower(nameQuery)
	results := []User{}

	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), nameQuery) {
			results = append(results, user)
		}
	}

	// Return matching users
	c.JSON(200, Response{
		Success: true,
		Data:    results,
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
	// Check required fields: Name, Email
	name := user.Name
	email := user.Email
	if name == "" || email == "" {
		return errors.New("参数不能为空")
		// Validate email format (basic check)
	} else if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("邮箱地址格式错误")
	}
	return nil
}
