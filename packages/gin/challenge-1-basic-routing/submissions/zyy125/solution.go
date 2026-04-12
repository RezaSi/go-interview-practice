package main

import (
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
	r := gin.Default()
	// TODO: Setup routes
	// GET /users - Get all users
  r.GET("/users", getAllUsers)
	// GET /users/:id - Get user by ID
  r.GET("/users/:id", getUserByID)
	// POST /users - Create new user
  r.POST("users", createUser)
	// PUT /users/:id - Update user
  r.PUT("/users/:id", updateUser)
	// DELETE /users/:id - Delete user
  r.DELETE("/users/:id", deleteUser)
	// GET /users/search - Search users by name
  r.GET("/users/search", searchUsers)
	// TODO: Start server on port 8080
  r.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users

  response := Response{
	Success: true,
	Data: users,
	Code: 200,
  }
  c.JSON(200, response)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
  id := c.Param("id")

  userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{Success: false, Code: 400, Error: "Invalid ID"})
		return
	}

  for _,user := range users {
    if user.ID == userID {
		response := Response{
			Code: 200,
			Data: user,
			Success: true,
		}
      c.JSON(200, response)
      return
    }
  }
  	response := Response{
		Code: 404,	
		Success: false,
		Error: "User not found",
	} 
    c.JSON(404, response)
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if newUser.Name == "" || newUser.Email == "" {
		c.JSON(400, Response{Success: false, Code: 400, Error: "Name and Email required"})
		return
	}

	newUser.ID = len(users) + 1 

	users = append(users, newUser)
	res := Response{
		Code: 201,
		Data: newUser,
		Success: true,
	}
	c.JSON(201, res)
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	idStr := c.Param("id")
	id,_ := strconv.Atoi(idStr)
	var updateUser User
	c.ShouldBindJSON(&updateUser)

	for i, user := range users {
		if user.ID == id {
			updateUser.ID = user.ID
			users[i] = updateUser
			res := Response{
				Success: true,
				Data: updateUser,
				Code: 200,
			}
			c.JSON(200, res)
			return
		}
	}
	c.JSON(404, gin.H{"error": "User not found","code": 404,})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(200, gin.H{
				"success": true,
				"code":200,
			})
			return
		}
	}
	c.JSON(404, gin.H{
		"code": 404,
		"error": "User not found",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")
    if name == "" {
        c.JSON(400, Response{Success: false, Code: 400, Error: "Missing name parameter"})
        return
    }

    results := []User{}
    for _, user := range users {
        if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
            results = append(results, user)
        }
    }

    // 即使没有结果，也要返回 200 和空数组
    c.JSON(200, Response{
        Success: true,
        Code:    200,
        Data:    results,
    })
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found

	for i, user := range users {
		if user.ID == id {
			return &user, i
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
