package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
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
    router:=gin.Default()

    router.GET("/users",getAllUsers)
    router.GET("/users/:id",getUserByID)
    router.POST("/users",createUser)
    router.PUT("/users/:id",getAllUsers)
    router.DELETE("/users/:id",getAllUsers)
    router.GET("/users/search",getAllUsers)

	// TODO: Setup routes
	// GET /users - Get all users
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name

	// TODO: Start server on port 8080
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(200,gin.H{
	    "success":true,
	    "data":users,
	})
}

func getUserByID(c *gin.Context) {
	sid := c.Param("id")

	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid user id",
		})
		return
	}

	for _, u := range users {
		if u.ID == id {
			c.JSON(200, gin.H{
				"success": true,
				"data":    u,
			})
			return
		}
	}

	// User not found
	c.JSON(404, gin.H{
		"success": false,
		"error":   "user not found",
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
   var u User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}
if strings.TrimSpace(u.Name) == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "name is required",
		})
		return
	}
	// ✅ Assign ID before saving
	u.ID = len(users) + 1

	users = append(users, u)

	c.JSON(201, gin.H{
		"success": true,
		"data":    u,
	})
}

func updateUser(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid user id",
		})
		return
	}

	var body User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	if strings.TrimSpace(body.Name) == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "name is required",
		})
		return
	}

	for i := range users {
		if users[i].ID == id {
			users[i].Name = body.Name

			c.JSON(200, gin.H{
				"success": true,
				"data":    users[i],
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"success": false,
		"error":   "user not found",
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	sid := c.Param("id")

	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid user id",
		})
		return
	}
	
	for i,v:=range users{
	    if v.ID==id{
	        users= append(users[:i], users[i+1:]...)
	        	c.JSON(200, gin.H{
				"success": true,
			})
			return 
	    }
	}
		c.JSON(404, gin.H{
		"success": false,
		"error":   "user not found",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
name := strings.TrimSpace(c.Query("name"))

	// ❌ Missing parameter
	if name == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "name query parameter is required",
		})
		return
	}

	name = strings.ToLower(name)
result := make([]User, 0)

	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), name) {
			result = append(result, u)
		}
	}

	// ✅ Always success=true, even if result is empty
	c.JSON(200, gin.H{
		"success": true,
		"data":    result,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	return nil
}
	  

