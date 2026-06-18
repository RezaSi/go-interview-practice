package main

import (
	"errors"
	"slices"
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

	router.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(200, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{
			"error": "id no es valido",
		})
		return
	}

	user, _ := findUserByID(id)

	if user == nil {
		c.JSON(404, Response{
			Success: false,
			Message: "User not found",
		})
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    user,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	err := validateUser(user)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	nextID += 1
	var createdUser User = User{ID: nextID, Name: user.Name, Email: user.Email, Age: user.Age}
	users = append(users, createdUser)

	c.JSON(201, Response{
	    Success: true,
	    Data: createdUser,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   "Id not specified",
			Code:    400,
		})
		return
	}

	index := slices.IndexFunc(users, func(u User) bool {
		return u.ID == id
	})

	if index == -1 {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	var userBody User

	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	if err := validateUser(userBody); err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	users[index] = User{
		ID:    users[index].ID,
		Name:  userBody.Name,
		Email: userBody.Email,
		Age:   userBody.Age,
	}

	c.JSON(200, Response{
		Success: true,
		Data:    users[index],
		Message: "User successfully updated",
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "Invalid id",
		})
		return
	}
	
	user, _ := findUserByID(id)
	
	if user == nil {
	    c.JSON(404, Response{
	        Success: false,
	        Message: "User not found",
	    })
	    return
	}
	
	users = slices.DeleteFunc(users, func(u User) bool {
		return u.ID == id
	})
	
	c.JSON(200, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")

	if name == "" {
		c.JSON(400, Response{
			Success: false,
			Message: "Specify a name",
		})
		return
	}

	var filteredUsers []User = make([]User, 0)

	for _, u := range users {
		isContained := strings.Contains(strings.ToLower(u.Name), strings.ToLower(name))
		if isContained {
			filteredUsers = append(filteredUsers, u)
		}
	}

	c.JSON(200, Response{
		Success: true,
		Data:    filteredUsers,
		Code:    200,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	index := slices.IndexFunc(users, func(u User) bool {
		return u.ID == id
	})

	if index == -1 {
		return nil, index
	}

	return &users[index], index
}

// Helper function to validate user data
func validateUser(user User) error {
	if user.Name == "" {
		return errors.New("Name is required")
	}

	if user.Email == "" {
		return errors.New("Email is required")
	}

	if !strings.Contains(user.Email, "@") {
		return errors.New("Invalid email format")
	}

	return nil
}
