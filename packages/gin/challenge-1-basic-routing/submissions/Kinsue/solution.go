package main

import (
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"sync"

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
var mu sync.Mutex

func main() {
	// Create Gin router
	router := gin.Default()

	// Setup routes
	// GET /users - Get all users
	router.GET("/users", getAllUsers)
	// GET /users/:id - Get user by ID
	router.GET("/user/:id", getUserByID)
	// POST /users - Create new user
	router.POST("/users", createUser)
	// PUT /users/:id - Update user
	router.PUT("/users/:id", updateUser)
	// DELETE /users/:id - Delete user
	router.DELETE("/users/:id", deleteUser)
	// GET /users/search - Search users by name
	router.GET("/users/search", searchUsers)

	// Start server on port 8080
	router.Run()
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// Return all users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Bad Request",
		})
		return
	}

	user, index := findUserByID(id)

	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    *user,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Bad Request",
		})
		return
	} else {
		users = append(users, User{
			ID:    nextID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		})

		c.JSON(http.StatusCreated, Response{
			Success: true,
			Data:    users[len(users)-1],
		})

		nextID++
	}
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	var userUpdate User
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	index := -1
	for i, u := range users {
		if u.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User Not Found",
		})
		return
	}

	users[index].Age = userUpdate.Age
	users[index].Email = userUpdate.Email
	users[index].Name = userUpdate.Name

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users[index],
	})

}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// Get user ID from path
	// Find and remove user
	// Return success message

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User Not Found",
		})
		return
	}

	users = append(users[:index], users[index+1:]...)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})

}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "No Query Params",
		})
		return
	}

	var usersFinded []User

	for _, u := range users {
		if strings.Contains(strings.ToUpper(u.Name), strings.ToUpper(name)) {

			usersFinded = append(usersFinded, u)
		}
	}

	if len(usersFinded) > 0 {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Data:    usersFinded,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    make([]interface{}, 0),
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	index := -1
	for i, u := range users {
		if u.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, index
	}
	return &users[index], index
}

// Helper function to validate user data
func validateUser(user User) error {
	// Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)

	if user.Name == "" {
		return fmt.Errorf("Name is empty")
	}

	if user.Email == "" {
		return fmt.Errorf("Email is empty")
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return err
	}
	return nil
}
