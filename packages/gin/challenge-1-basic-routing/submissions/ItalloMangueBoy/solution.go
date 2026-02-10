package main

import (
	"fmt"
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
	r.Run("8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
		c.JSON(200, Response{
		Success: true,
		Data:    users,
		Code:    200,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	ID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    400,
		})
		return
	}

	user, _ := findUserByID(int(ID))
	if user == nil {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    user,
		Code:    200,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   "Invalid request body",
			Code:    400,
		})
		return
	}

	if err := validateUser(newUser); err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
			Code:    422,
		})
		return
	}

	newUser.ID = nextID
	nextID++
	
	users = append(users, newUser)

	c.JSON(201, Response{
		Success: true,
		Data:    newUser,
		Code:    201,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated use
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    400,
		})
		return
	}

	var updatedUser User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   "Invalid request body",
			Code:    400,
		})
		return
	}

	if err := validateUser(updatedUser); err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
			Code:    422,
		})
		return
	}

	user, index := findUserByID(ID)
	if index == -1 {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}

	users[index] = updatedUser

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
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    400,
		})
		return
	}

	user, index := findUserByID(ID)
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
		Message: "User deleted successfully",
		Code:    200,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	nameQuery := strings.ToLower(c.Query("name"))
	if nameQuery == "" {
		c.JSON(400, Response{
			Success: false,
			Error:   "Name query parameter is required",
			Code:    400,
		})
		return
	}


	var matchedUsers []User
	for _, user := range users {
		if strings.Contains(
			strings.ToLower(user.Name),
			nameQuery,
		) {
			matchedUsers = append(matchedUsers, user)
		}
	}
	
	if len(matchedUsers) == 0 {
		c.JSON(200, Response{
			Success: true,
			Error:   "No users found",
			Data:    []interface{}{},
			Code:    200,
		})
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    matchedUsers,
		Code:    200,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	for i := range users {
    	if users[i].ID == id {
			return &users[i], i
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
		return fmt.Errorf("Name is required")
	}

	if user.Email == "" {
		return fmt.Errorf("Email is required")
	}

	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("Invalid email format")
	}

	if !strings.Contains(user.Email, ".com") {
		return fmt.Errorf("Invalid email format")
	}

	return nil
}
