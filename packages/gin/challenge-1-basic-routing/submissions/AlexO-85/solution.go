package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id" binding:"omitempty,number"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,number"`
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
	router := gin.Default()

	// TODO: Setup routes
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

	router.Run("localhost:8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	response := Response{
		Success: true,
		Data:    users,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	// Handle invalid ID format
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Success: false,
			Error:   "User ID format is wrong",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	if user, _ := findUserByID(id); user != nil {
		response := Response{
			Success: true,
			Data:    user,
		}
		c.IndentedJSON(http.StatusOK, response)
		return
	}

	response := Response{
		Success: false,
		Error:   "User not found",
		Code:    http.StatusNotFound,
	}

	c.IndentedJSON(http.StatusNotFound, response)
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		response := Response{
			Success: false,
			Error:   "Bas request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	user.ID = nextID
	users = append(users, user)
	nextID++

	response := Response{
		Success: true,
		Data:    user,
	}
	c.IndentedJSON(http.StatusCreated, response)
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	// Handle invalid ID format
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Success: false,
			Error:   "User ID format is wrong",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		response := Response{
			Success: false,
			Error:   "Bas request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	if user, _ := findUserByID(id); user != nil {
		user.Name = newUser.Name
		user.Email = newUser.Email
		user.Age = newUser.Age

		response := Response{
			Success: true,
			Data:    user,
		}
		c.IndentedJSON(http.StatusOK, response)
		return
	}

	response := Response{
		Success: false,
		Error:   "User not found",
		Code:    http.StatusNotFound,
	}

	c.IndentedJSON(http.StatusNotFound, response)
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	// Handle invalid ID format
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Success: false,
			Error:   "User ID format is wrong",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	if user, i := findUserByID(id); user != nil {
		response := Response{
			Success: true,
			Data:    user,
		}
		c.IndentedJSON(http.StatusOK, response)

		users = append(users[:i], users[i+1:]...)

		return
	}

	response := Response{
		Success: false,
		Error:   "User not found",
		Code:    http.StatusNotFound,
	}

	c.IndentedJSON(http.StatusNotFound, response)
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {

	name := strings.ToLower(c.Query("name"))

	if name == "" {
		response := Response{
			Success: false,
			Error:   "Name parameter is missing",
			Code:    http.StatusBadRequest,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	found := []User{}

	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), name) {
			found = append(found, user)
		}
	}

	response := Response{
		Success: true,
		Data:    found,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for i := range users {
		if users[i].ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}
