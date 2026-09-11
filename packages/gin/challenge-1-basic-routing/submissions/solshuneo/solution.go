package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id" `
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

	// TODO: Setup routes
	// GET /users - Get all users
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name

	// TODO: Start server on port 8080
	user, _ := findUserByID(1)
	fmt.Printf("%v", user)
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(200, Response{
		Success: true,
		Data:    users,
		Code:    200,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	if idUser, err := strconv.Atoi(c.Param("id")); err == nil {
		user, _ := findUserByID(idUser)
		if user != nil {
			c.JSON(200, Response{
				Data:    *user,
				Success: true,
				Code:    200,
			})
		} else {
			c.JSON(404, Response{
				Success: false,
				Code:    404,
			})
		}
	} else {
		c.JSON(400, Response{
			Success: false,
			Code:    400,
		})
	}
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var user User
	if err := c.ShouldBind(&user); err == nil {
		user.ID = nextID
		nextID += 1
		users = append(users, user)
		c.JSON(201, Response{
			Success: true,
			Code:    201,
			Data:    user,
		})
	} else {
		c.JSON(400, Response{
			Success: false,
			Code:    400,
		})
	}
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	var id int
	var err error
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(401, Response{
			Success: false,
			Code:    401,
		})
		return
	}
	var userUpdate User
	if err := c.ShouldBind(&userUpdate); err != nil {
		return
	}
	userUpdate.ID = id
	for index := range users {
		if users[index].ID == userUpdate.ID {
			users[index] = userUpdate

			c.JSON(200, Response{
				Success: true,
				Code:    200,
				Data:    users[index],
			})
			return
		}
	}
	c.JSON(404, Response{
		Success: false,
		Code:    404,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	var (
		id  int
		err error
	)
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		return
	}
	for index, user := range users {
		if user.ID == id {
			if index == len(users)-1 {
				users = users[:index]
			} else {
				users = append(users[:index], users[index+1:]...)
			}
			c.JSON(200, Response{
				Code:    200,
				Success: true,
			})
			return
		}
	}
	c.JSON(404, Response{
		Code: 404,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	_, ok := c.GetQuery("name")
	if !ok {
		c.JSON(400, Response{
			Success: false,
		})
		return
	}
	nameFilter := strings.ToLower(c.Query("name"))
	userRes := []User{}
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), nameFilter) {
			userRes = append(userRes, user)
		}
	}
	c.JSON(200, Response{
		Success: true,
		Data:    userRes,
		Code:    200,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	for index, user := range users {
		if id == user.ID {
			return &user, index
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)

	return nil
}
