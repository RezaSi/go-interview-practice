package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"errors"
	"fmt"
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
	router.DELETE("users/:id", deleteUser)
	// GET /users/search - Search users by name
	router.GET("/users/search", searchUsers)

	router.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(200, Response{
	    Success: true,
	    Data: users,
	    Code: 200,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
    id := c.Param("id")
    
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "Invalid request",
			Error: err.Error(),
			Code: 400,
		})
		return
	}

	user, index := findUserByID(userId)
	if index == -1 {
		c.JSON(404, Response{
			Success: false,
			Data: user,
			Message: "User not found",
			Error: "User not found",
			Code: 404,
		})

		return
	}

	c.JSON(200, Response{
		Success: true,
		Data: user,
		Code: 200,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "Invalid request",
			Error: err.Error(),
			Code: 400,
		})
		return
	}

	// Validate required fields
	if err := validateUser(user); err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "Invalid request",
			Error: err.Error(),
			Code: 400,
		})
		return
	}

	// Add user to storage
	user.ID = nextID
	nextID += 1
	users = append(users, user)

	// Return created user
	c.JSON(201, Response{
		Success: true,
		Data: user,
		Code: 201,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// Parse JSON request body
	id := c.Param("id")
	
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "Invalid request",
			Error: err.Error(),
			Code: 400,
		})
		return
	}
	
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "Invalid request",
			Error: err.Error(),
			Code: 400,
		})
		return
	}

	// Find and update user
	_, index := findUserByID(userId)
	if index == -1 {
		c.JSON(404, Response{
			Success: false,
			Message: "User not found",
			Error: "User not found",
			Code: 404,
		})
		return
	}

	if user.Name != "" {
		users[index].Name = user.Name
	}
	if user.Email != "" {
		users[index].Email = user.Email
	}
	if user.Age > 0 {
		users[index].Age = user.Age
	}

	// Return updated user
	c.JSON(200, Response{
		Success: true,
		Data: users[index],
		Code: 200,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "Invalid request",
			Error: err.Error(),
			Code: 400,
		})
		return
	}

	_, index := findUserByID(userId)
	if index == -1 {
		c.JSON(404, Response{
			Success: false,
			Message: "User not found",
			Error: "User not found",
			Code: 404,
		})
		return
	}
	
	users = append(users[:index], users[index+1:]...)

	// Return success message
	c.JSON(200, Response{
		Success: true,
		Message: fmt.Sprintf("User ID %v successfully deleted", userId),
		Code: 200,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// Filter users by name (case-insensitive)
	name := c.Query("name")
	if name == "" {
		c.JSON(400, Response{
			Success: false,
			Message: "Invalid request",
			Error: "Name is required",
			Code: 400,
		})
		return
	}
	
	lname := strings.ToLower(name)
	
	usersFound := []*User{}
	for idx, usr := range users {
		if userName := strings.ToLower(usr.Name); strings.Contains(userName, lname) {
			usersFound = append(usersFound, &users[idx])
		}
	}

	// Return matching users
	c.JSON(200, Response{
		Success: true,
		Data: usersFound,
		Code: 200,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
    for idx, usr := range users {
		if (usr.ID == id) {
			return &users[idx], idx
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email format")
    }
	return nil
}
