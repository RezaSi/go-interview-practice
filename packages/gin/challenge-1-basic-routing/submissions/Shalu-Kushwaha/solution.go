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
    route := gin.Default()
    
	// TODO: Setup routes
	// GET /users - Get all users
	route.GET("/users", getAllUsers)
	
	// GET /users/:id - Get user by ID
	route.GET("/users/:id",getUserByID)
	
	// POST /users - Create new user
	route.POST("/users",createUser)
	
	// PUT /users/:id - Update user
	route.PUT("/users/:id",updateUser)
	
	// DELETE /users/:id - Delete user
	route.DELETE("/users/:id",deleteUser)
	
	// GET /users/search - Search users by name
	route.GET("/users/search",searchUsers)

	// TODO: Start server on port 8080
	route.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Message: "Users fetched successfully",
		Code:    http.StatusOK,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}
	
	// Return 404 if user not found
	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "User fetched successfully",
		Code:    http.StatusOK,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var newUser User
	
	// Validate required fields
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
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
	
	// Add user to storage
	users = append(users, newUser)

	// Return created user
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "User created successfully",
		Code:    http.StatusCreated,
	})
}


// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	var updatedUser User

	// Parse request body
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate data
	if err := validateUser(updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Keep same ID
	updatedUser.ID = id

	// Update user
	users[index] = updatedUser

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    updatedUser,
		Message: "User updated successfully",
		Code:    http.StatusOK,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Remove user
	users = append(users[:index], users[index+1:]...)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
		Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))

	// Check missing parameter
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "name query parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	name = strings.ToLower(name)

	// Initialize empty slice
	matchedUsers := []User{}

	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), name) {
			matchedUsers = append(matchedUsers, user)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    matchedUsers,
		Message: "Search completed successfully",
		Code:    http.StatusOK,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for index, user := range users {
		if user.ID == id {
			return &users[index], index
		}
	}

	return nil, -1

}

// Helper function to validate user data
func validateUser(user User) error {
	// Check required fields
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}

	// Basic email validation
	if !strings.Contains(user.Email, "@") ||
		!strings.Contains(user.Email, ".") {
		return errors.New("invalid email format")
	}

	return nil
}
