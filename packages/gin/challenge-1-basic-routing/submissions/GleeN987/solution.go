package main

import (
	"errors"
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
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name
	router.GET("/users", getAllUsers)
	router.GET("/users/search", searchUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	// TODO: Start server on port 8080
	router.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(200, Response{
		Success: true,
		Data:    users,
		Message: "Request successful",
		Code:    200,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
	for _, user := range users {
		if user.ID == userID {
			c.JSON(200, Response{
				Success: true,
				Data:    user,
				Message: "Request successful",
				Code:    200,
			})
			return
		}
	}
	c.JSON(404, Response{
		Success: false,
		Error:   "User not found",
		Code:    404,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err = validateUser(user)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	user.ID = nextID
	nextID++
	users = append(users, user)
	c.JSON(201, Response{
		Success: true,
		Data:    user,
		Error:   "no errors",
		Code:    201,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
	var userUpdated User
	err = c.ShouldBindJSON(&userUpdated)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
	user, _ := findUserByID(userID)
	if user == nil {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}

	user.Age = userUpdated.Age
	user.Email = userUpdated.Email
	user.Name = userUpdated.Name
	c.JSON(200, Response{
		Success: true,
		Data:    user,
		Code:    200,
	})

}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(404, Response{
			Success: false,
			Error:   err.Error(),
			Code:    404,
		})
	}
	user, index := findUserByID(userID)
	if user == nil {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}

	users = append(users[:index], users[index+1:]...)
	c.JSON(200, Response{
		Success: true,
		Message: "Deleting successful",
		Code:    200,
	})

}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	search := c.Query("name")
	if search == "" {
		c.JSON(400, Response{
			Success: false,
			Error:   "Missing search query parameter",
			Code:    400,
		})
		return
	}
	searchResult := []User{}
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(search)) {
			searchResult = append(searchResult, user)
		}
	}
	c.JSON(200, Response{
		Success: true,
		Data:    searchResult,
		Code:    200,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	for index, user := range users {
		if user.ID == id {
			return &user, index
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	if user.Name == "" {
		return errors.New("Name is required")
	}
	if user.Email == "" {
		return errors.New("Email is required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("Incorrect email format")
	}
	return nil
}
