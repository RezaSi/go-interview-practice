package main

import (
	"errors"
	"net/http"
	"regexp"
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

//main.go starts herez
func main() {
	// TODO: Create Gin router
	router := gin.Default()

	// TODO: Setup routes
	router.GET("/users", getAllUsers)        // GET /users - Get all users
	router.GET("/users/:id", getUserByID)    // GET /users/:id - Get user by ID
	router.POST("/users", createUser)        // POST /users - Create new user
	router.PUT("/users/:id", updateUser)     // PUT /users/:id - Update user
	router.DELETE("/users/:id", deleteUser)  // DELETE /users/:id - Delete user
	router.GET("/users/search", searchUsers) // GET /users/search - Search users by name

	router.Run(":8080") // TODO: Start server on http://localhost:8080
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	userID, err := strconv.Atoi(c.Param("id"))

	// Handle invalid ID format
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}
	foundUser, index := findUserByID(userID)
	// Return 404 if user not found
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "404: User not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    foundUser,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var newUser User // Validate required fields
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid JSON format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	newUser.ID = nextID
	nextID++
	users = append(users, newUser) // Add user to storage
	// Return created user
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "User creation was done",
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	_, index := findUserByID(userID) //- так вроде правильно
	if index < 0 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid JSON format",
			Code:    http.StatusBadRequest,
		})
		return
	}
	updatedUser.ID = userID

	if err := validateUser(updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	users[index] = updatedUser

	// Parse JSON request body
	// Find and update user
	// Return updated user
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    updatedUser,
		Message: "User updated Succsesfully",
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	if c == nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusBadRequest,
		})
		return
	}
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if userID < 0 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	_, index := findUserByID(userID)
	if index < 0 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	// Find and remove user
	// Return success message
	users = append(users[:index], users[index+1:]...)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted Succsesfully",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	nameQuery := c.Query("name")
	if nameQuery == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Missing 'name' query parameter",
			Code:    http.StatusBadRequest,
		})
		return
	}
	var matchedUsers []User
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(nameQuery)) {
			matchedUsers = append(matchedUsers, u)
		}
	}
	// Filter users by name (case-insensitive)
	if matchedUsers == nil {
		matchedUsers = make([]User, 0) // или []User{}
	}
	// Return matching users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    matchedUsers,
		Message: "Search completed",
		Code:    http.StatusOK,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	for i, u := range users {
		if u.ID == id {
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
	if user.Name == "" {
		return errors.New("Name is required")
	}
	if user.Email == "" {
		return errors.New("Email is required")
	}
	if matched, _ := regexp.MatchString(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, user.Email); !matched {
		return errors.New("Email is invalid")
	}
	//
	return nil
}
