package main

import (
	"log"
	"net/http"
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

type UserUriId struct {
	ID int `uri:"id" binding:"required"`
}

// In-memory storage
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var nextID = 4
var usersMutex sync.RWMutex

func main() {
	router := gin.Default()
	// GET /users - Get all users
	router.GET("/users", getAllUsers)

	// GET /users/:id - Get user by ID
	router.GET("/users/:id", getUserByID)

	// POST /users - Create new user
	router.POST("/users", createUser)

	// PUT /users/:id - Update user
	router.PUT("/users/:id", updateUser)

	// DELETE /users/:id - Delete user
	router.DELETE("/users/:id", deleteUser)

	// GET /users/search - Search users by name
	router.GET("/users/search", searchUsers)

	// Start server on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	// Return all users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// Get user by ID
	var userId UserUriId
	if err := c.ShouldBindUri(&userId); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	usersMutex.RLock()
	user, index := findUserByID(userId.ID)
	usersMutex.RUnlock()

	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    *user,
	})
}

type NewUser struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Age   int    `json:"age" binding:"required"`
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// Parse JSON request body
	var newUser NewUser

	// Validate required fields
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	usersMutex.Lock()

	var user User
	user.ID = nextID
	nextID++
	user.Name = newUser.Name
	user.Email = newUser.Email
	user.Age = newUser.Age
	// Add user to storage
	users = append(users, user)
	defer usersMutex.Unlock()
	// Return created user
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// Get user ID from path
	var userUriId UserUriId
	if err := c.ShouldBindUri(&userUriId); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Parse JSON request body
	var user NewUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Find and update user
	usersMutex.Lock()

	var index int = -1
	for i := 0; i < len(users); i++ {
		if users[i].ID == userUriId.ID {
			index = i
			users[i].Name = user.Name
			users[i].Email = user.Email
			users[i].Age = user.Age
			break
		}
	}
	defer usersMutex.Unlock()
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	// Return updated user
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users[index],
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// Get user ID from path
	var userUriId UserUriId
	if err := c.ShouldBindUri(&userUriId); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Find and remove user
	usersMutex.Lock()
	defer usersMutex.Unlock()

	index := -1
	for i := 0; i < len(users); i++ {
		if users[i].ID == userUriId.ID {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Remove user from storage
	users = append(users[:index], users[index+1:]...)
	// Return success message
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    "User deleted successfully",
	})
}

type SearchUsersQuery struct {
	Name string `form:"name" binding:"required"`
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// Get name query parameter
	var searchUsersQuery SearchUsersQuery
	if err := c.ShouldBindQuery(&searchUsersQuery); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Filter users by name (case-insensitive)
	usersMutex.RLock()

	result := []User{}
	for i := 0; i < len(users); i++ {
		if strings.Contains(strings.ToLower(users[i].Name), strings.ToLower(searchUsersQuery.Name)) {
			result = append(result, users[i])
		}
	}
	defer usersMutex.RUnlock()

	// Return matching users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    result,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// Implement user lookup

	for i := 0; i < len(users); i++ {
		if users[i].ID == id {
			return &users[i], i
		}
	}
	// Return user pointer and index, or nil and -1 if not found
	return nil, -1
}
