package main

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Age   int    `json:"age" binding:"required"`
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
	r := gin.New()

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

	// TODO: Start server on port 8080
	http.ListenAndServe(":8080", r)
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(http.StatusOK, Response{Success: true, Data: users})

}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false})
	}

	user, _ := findUserByID(userId)

	if user == nil {
		c.JSON(http.StatusNotFound, Response{Success: false})
		return
	}

	c.JSON(http.StatusOK, Response{Success: true, Data: user})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false})
		return
	}

	// Validate required fields

	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false})
	}

	// Add user to storage
	user.ID = nextID
	nextID++
	users = append(users, user)

	// Return created user
	c.JSON(http.StatusCreated, Response{Success: true, Data: user})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Success: false})
		return
	}

	// Parse JSON request body
	var inputUser User
	err = c.ShouldBindJSON(&inputUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false})
		return
	}

	// Find and update user
	_, idx := findUserByID(userId)

	if idx == -1 {
		c.JSON(http.StatusNotFound, Response{Success: false})
		return
	}

	inputUser.ID = users[idx].ID
	users[idx] = inputUser

	// Return updated user
	c.JSON(http.StatusOK, Response{Success: true, Data: inputUser})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Success: false})
		return
	}

	// Find and remove user
	_, idx := findUserByID(userId)

	if idx == -1 {
		c.JSON(http.StatusNotFound, Response{Success: false})
		return
	}
	users = slices.Delete(users, idx, idx+1)

	// Return success message
	c.JSON(http.StatusOK, Response{Success: true})

}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	nameStr := strings.ToLower(c.Query("name"))
	if nameStr == "" {
		c.JSON(http.StatusBadRequest, Response{Success: false})
		return
	}

	// Filter users by name (case-insensitive)
	foundUsers := []User{}
	for i := range users {
		if strings.Contains(strings.ToLower(users[i].Name), nameStr) {
			foundUsers = append(foundUsers, users[i])
		}
	}

	// Return matching users
	c.JSON(http.StatusOK, Response{Success: true, Data: foundUsers})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for i := 0; i < len(users); i++ {
		if users[i].ID == id {
			return &users[i], i
		}
	}

	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	if user.Email == "" || user.Name == "" {
		return errors.New("Invalid")
	}
	return nil
}
