package main

import (
	"errors"
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

func ResponseSuccess(data interface{}, code int) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: "success",
		Code:    code,
	}
}
func ResponseFailure(err string, code int) Response {
	return Response{
		Success: false,
		Error:   err,
		Code:    code,
	}
}

// In-memory storage
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var usersLock sync.RWMutex

var nextID = 4
var nextIdLock sync.Mutex

func withNextId(fn func(int)) {
	nextIdLock.Lock()
	defer nextIdLock.Unlock()
	id := nextID
	nextID++
	fn(id)
}

func main() {
	router := gin.Default()
	// GET /users/search - Search users by name
	router.GET("/users/search", searchUsers)
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

	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	usersLock.RLock()
	defer usersLock.RUnlock()

	c.JSON(http.StatusOK, ResponseSuccess(users, http.StatusOK))
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// Handle invalid ID format
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseFailure(fmt.Sprintf("id is not number: %s", err), http.StatusBadRequest))
		return
	}

	usersLock.RLock()
	defer usersLock.RUnlock()
	// Return 404 if user not found
	var user *User
	if user, _ = findUserByID(id); user == nil {
		c.JSON(http.StatusNotFound, ResponseFailure("user not found", http.StatusNotFound))
		return
	}
	c.JSON(http.StatusOK, ResponseSuccess(user, http.StatusOK))
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// Parse JSON request body
	var userParam User
	if err := c.ShouldBind(&userParam); err != nil {
		c.JSON(http.StatusBadRequest, ResponseFailure(fmt.Sprintf("bad input: %s", err.Error()), http.StatusBadRequest))
		return
	}
	// Validate required fields
	if err := validateUser(userParam, true); err != nil {
		c.JSON(http.StatusBadRequest, ResponseFailure(err.Error(), http.StatusBadRequest))
		return
	}
	// Add user to storage
	usersLock.Lock()
	defer usersLock.Unlock()

	var newUser User
	withNextId(func(nextId int) {
		newUser = User{
			ID:    nextId,
			Name:  userParam.Name,
			Email: userParam.Email,
			Age:   userParam.Age,
		}
		users = append(users, newUser)
	})
	// Return created user
	c.JSON(201, ResponseSuccess(newUser, 201))
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// Get return user ID from path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseFailure(fmt.Sprintf("id is not number: %s", err), http.StatusBadRequest))
		return
	}
	// Parse JSON request body
	var userParam User
	if err := c.ShouldBind(&userParam); err != nil {
		c.JSON(http.StatusBadRequest, ResponseFailure(fmt.Sprintf("bad input: %s", err), http.StatusBadRequest))
		return
	}
	if err := validateUser(userParam, false); err != nil {
		c.JSON(http.StatusBadRequest, ResponseFailure(err.Error(), http.StatusBadRequest))
		return
	}
	// Find and update user
	usersLock.Lock()
	defer usersLock.Unlock()

	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, ResponseFailure("user not found", http.StatusNotFound))
		return
	}
	if userParam.Name != "" {
		user.Name = userParam.Name
	}
	if userParam.Email != "" {
		user.Email = userParam.Email
	}
	if userParam.Age != 0 {
		user.Age = userParam.Age
	}
	// Return updated user
	c.JSON(http.StatusOK, ResponseSuccess(user, 200))
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// Get user ID from path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseFailure(fmt.Sprintf("id is not number: %s", err), http.StatusBadRequest))
		return
	}
	// Find and remove user
	usersLock.Lock()
	defer usersLock.Unlock()

	_, i := findUserByID(id)
	if i < 0 {
		c.JSON(http.StatusNotFound, ResponseFailure("user not found", http.StatusNotFound))
		return
	}
	users = append(users[:i], users[i+1:]...)
	// Return success message
	c.JSON(http.StatusOK, ResponseSuccess(nil, 200))
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// Get name query parameter
	nameParam := c.Query("name")
	if nameParam == "" {
		c.JSON(http.StatusBadRequest, ResponseFailure("name is required", http.StatusBadRequest))
		return
	}
	nameParam = strings.ToLower(nameParam)
	// Filter users by name (case-insensitive)
	nameUsers := make([]User, 0)

	usersLock.RLock()
	defer usersLock.RUnlock()

	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), nameParam) {
			nameUsers = append(nameUsers, user)
		}
	}
	// Return matching users
	c.JSON(http.StatusOK, ResponseSuccess(nameUsers, 200))
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	var user *User
	var i int
	for i = 0; i < len(users); i++ {
		if id == users[i].ID {
			user = &users[i]
			break
		}
	}
	if user == nil {
		return nil, -1
	}
	return user, i
}

// Helper function to validate user data
func validateUser(user User, checkEmpty bool) error {
	// Check required fields: Name, Email
	// Validate email format (basic check)
	if checkEmpty {
		if user.Name == "" {
			return errors.New("user name is empty")
		}
		if user.Email == "" {
			return errors.New("user email is empty")
		}
		if user.Age == 0 {
			return errors.New("user age is zero")
		}
	}
	if user.Age < 0 {
		return errors.New("user age is not positive")
	}
	if user.Email != "" {
		if _, err := mail.ParseAddress(user.Email); err != nil {
			return errors.New("invalid email")
		}
	}
	return nil
}
