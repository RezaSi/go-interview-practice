package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"errors"
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
    router := gin.Default()
	// TODO: Setup routes
	router.GET("/users",getAllUsers)
	router.GET("/users/search", searchUsers)
	router.GET("/users/:id",getUserByID)
	router.POST("/users",createUser)
	router.PUT("/users/:id",updateUser)
	router.DELETE("/users/:id",deleteUser)

	// GET /users - Get all users
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name
    
    router.Run()
	// TODO: Start server on port 8080
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(200, gin.H{
	    "success":true,
	    "data":users,
	    "message": "Users retrieved successfully",
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
	    c.JSON(400,gin.H{"error":"Invalid ID"})
	    return
	}
	// Handle invalid ID format
	for i := 0; i < len(users); i++ {
	    if users[i].ID == userID {
	        c.JSON(200, gin.H{
	            "success":true,
	            "data":users[i],
	            "message": "Users retrieves successfully",
	        })
	        return
	    }
	}
	c.JSON(404,gin.H{"error":"user ID not found"})
	// Return 404 if user not found
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	maxID := 0
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
	    c.JSON(400, gin.H{"error": err.Error()})
	    return
	}
	if err := validateUser(user); err != nil {
	    c.JSON(400, gin.H{"error": err.Error()})
	    return
	}
	
	for _, v := range users {
	    if v.ID > maxID {
	        maxID = v.ID
	    }
	}
	user.ID = maxID
	users = append(users, user)
	c.JSON(201, gin.H{
	    "success":true,
	    "data":user,
	    "message": "Users created successfully",
	})
	// Validate required fields
	// Add user to storage
	// Return created user
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	var user User
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
	    c.JSON(400,gin.H{"error":"Invalid ID"})
	    return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
	    c.JSON(400, gin.H{"error": err.Error()})
	    return
	}
	for i := 0; i < len(users); i++ {
	    if users[i].ID == userID {
	        users[i] = user
	        c.JSON(200, gin.H{
	            "success":true,
	            "data":users[i],
	            "message": "Users updates successfully",
	        })
	        return
	    }
	}
	c.JSON(404,gin.H{"error":"user ID not found"})
	// Parse JSON request body
	// Find and update user
	// Return updated user
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
	    c.JSON(400,gin.H{"error":"Invalid ID"})
	    return
	}
	
	for i := 0; i < len(users); i++ {
	    if users[i].ID == userID {
	        users = append(users[:i], users[i+1:]...)
	        c.JSON(200, gin.H{
	            "success":true,
	            "data":id,
	            "message": "Users deletes successfully",
	        })
	        return
	    }
	    c.JSON(404,gin.H{"error":"user ID not found"})
	}
	// Find and remove user
	// Return success message
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
    name := c.Query("name")
    
    if name == "" {
        c.JSON(400, gin.H{"error":"parameters needed"})
        return
    }
    
    lowername := strings.ToLower(name)
    result := make([]User, 0)
    for _, v := range users {
        lowerCaseName := strings.ToLower(v.Name)
        if strings.Contains(lowerCaseName, lowername) {
            result = append(result, v)
        }
    }
    c.JSON(200, gin.H{
	            "success":true,
	            "data":result,
	            "message": "search successfully",
	        })
	        return
    
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
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
	if user.Name == "" {
	    return errors.New("name is required")
	}
	if user.Email == "" {
	    return errors.New("email is required")
	}
	if !strings.Contains(user.Email, "@") {
	    return errors.New("invalid email format")
	}
	// Check required fields: Name, Email
	// Validate email format (basic check)
	return nil
}
