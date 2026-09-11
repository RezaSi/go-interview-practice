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
	router := gin.Default()
	// TODO: Setup routes
	// GET /users - Get all users
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name
	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/search", searchUsers)

	// TODO: Start server on port 8080
	router.Run(":8080")

}

// TODO: Implement handler functions

func handleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
		Code:    statusCode,
	})
}

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

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, "cannot parse user from json")
		return
	}
	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	user.ID = nextID
	nextID++
	users = append(users, user)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
		Code:    201,
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
		handleError(c, http.StatusBadRequest, "invalid 'id' format")
		return
	}
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, "cannot parse user from json")
		return
	}
	if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	newUser.ID = id
	updated := false
	for i := 0; i < len(users); i++ {
		if users[i].ID == id {
			users[i] = newUser

			updated = true
			break
		}
	}
	if !updated {
		handleError(c, 404, "user to update not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    newUser,
		Code:    200,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid 'id' format")
		return
	}
	deleted := false
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			deleted = true
			break
		}
	}
	if !deleted {
		handleError(c, 404, "user to delete not found")
	}
	c.JSON(200, Response{
		Success: true,
		Code:    200,
		Message: "user deleted",
	})

}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	nameQuery := c.Query("name")
	if nameQuery == "" {
		handleError(c, http.StatusBadRequest, "name query is invalid")
		return
	}
	filteredUsers := make([]User, 0, len(users))
	for _, v := range users {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(nameQuery)) {
			filteredUsers = append(filteredUsers, v)
		}
	}
	c.JSON(200, Response{
		Success: true,
		Data:    filteredUsers,
		Code:    200,
	})
}

func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid 'id' format")
		return
	}
	user, index := findUserByID(id)
	if index == -1 && user == nil {
		handleError(c, 404, "user not found")
		return
	}
	c.JSON(200, Response{
		Success: true,
		Code:    200,
		Data:    user,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	index := -1
	var pointer *User
	for i, value := range users {
		if value.ID == id {
			pointer = &value
			index = i
			break
		}
	}
	return pointer, index
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	if user.Name == "" {
		return errors.New("field 'Name' is required")
	}
	if user.Email == "" {
		return errors.New("field 'Email' is required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid format of the field 'Email'")
	}
	return nil
}
