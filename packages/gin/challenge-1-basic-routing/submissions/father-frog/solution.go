package main

import (
	"errors"
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
var (
	usersMutex sync.RWMutex
	users      = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
	}
	nextID = 4
)

func main() {
	r := gin.Default()

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

	r.Run()
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	usersMutex.RLock()
	defer usersMutex.RUnlock()
	user, idx := findUserByID(id)
	if idx < 0 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "user not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "user found",
		Code:    http.StatusOK,
		Data:    user,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var inputUser User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "bad user data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate required fields
	if err := validateUser(inputUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "invalid user data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Add user to storage
	usersMutex.Lock()
	defer usersMutex.Unlock()
	inputUser.ID = nextID
	nextID++
	users = append(users, inputUser)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "added user",
		Code:    http.StatusCreated,
		Data:    inputUser,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	var inputUser User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "bad user data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate required fields
	if err := validateUser(inputUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "invalid user data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Find and update user
	usersMutex.Lock()
	defer usersMutex.Unlock()
	_, idx := findUserByID(id)
	if idx < 0 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "user not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	users[idx].Age = inputUser.Age
	users[idx].Email = inputUser.Email
	users[idx].Name = inputUser.Name

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "updated user",
		Code:    http.StatusOK,
		Data:    users[idx],
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	// Find and remove user
	usersMutex.Lock()
	defer usersMutex.Unlock()
	_, idx := findUserByID(id)
	if idx < 0 {
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
		Message: "deleted user",
		Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	// Get name query parameter
	name := c.Query("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "no search name",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Filter users by name (case-insensitive)
	matchedUsers := []User{}
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			matchedUsers = append(matchedUsers, user)
		}
	}

	// Return matching users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    matchedUsers,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for i, user := range users {
		if user.ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	if len(user.Name) == 0 {
		return errors.New("name is required")
	}

	if len(user.Email) == 0 {
		return errors.New("email is required")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return err
	}

	return nil
}

// parseIDParam parses and validates the ID parameter from the URL
func parseIDParam(c *gin.Context) (int, error) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "bad id",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
	}
	return id, err
}
