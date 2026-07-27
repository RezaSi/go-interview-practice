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
	// GET /users - Get all users
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name

    router.GET("/users", getAllUsers)
    router.GET("/users/:id", getUserByID)
    router.POST("/users", createUser)
    router.PUT("/users/:id", updateUser)
    router.DELETE("/users/:id", deleteUser)
    router.GET("/users/search", searchUsers)


	// TODO: Start server on port 8080
	
	router.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(200, gin.H{
	    "success":true,
	    "data":users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil{
	    c.JSON(400, gin.H{"error":"Invalid ID"})
	    return
	}
	user, _ := findUserByID(userID)
	if user == nil{
	    c.JSON(404, gin.H{"error":"User with given ID not found"})
	    return
	}
	
	c.JSON(200, gin.H{
	    "success":true,
	    "data":user,
	})
	
	
	
	
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if err := validateUser(user); err != nil{
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user.ID = nextID
    users = append(users, user)
    nextID += 1
    
    c.JSON(201, gin.H{
        "success":true,
	    "data":user,
    })
    
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if err := validateUser(user); err != nil{
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
	id := c.Param("id")
    userID, err := strconv.Atoi(id)       
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid ID"})
        return
    }
    
    
    userDB, _ := findUserByID(userID)
    if userDB == nil{
        c.JSON(404, gin.H{"error": "User with given ID not found"})
        return
    }
	
	userDB.Name = user.Name
	userDB.Email = user.Email
	userDB.Age = user.Age
	
    c.JSON(200, gin.H{
        "success":true,
	    "data":userDB,
    })
	
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	id := c.Param("id")
    userID, err := strconv.Atoi(id)       
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid ID"})
        return
    }
    
    user, index := findUserByID(userID)
    if user == nil{
        c.JSON(404, gin.H{"error": "User with given ID not found"})
        return
    }
    
    users = append(users[:index], users[index+1:]...)
    c.JSON(200, gin.H{
        "success":true,
	    "data":user,
    })
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	
	name := c.Query("name")
	if name == ""{
        c.JSON(400, gin.H{"error": "Invalid Query"})
        return
	}
	found_users := make([]User, 0)
	
	for _, user := range users{
	    if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)){
	        found_users = append(found_users, user)
	    }
	}
	
    c.JSON(200, gin.H{
        "success":true,
	    "data":found_users,
    })
	
	
	
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
    for i := 0; i < len(users); i++ {
        if users[i].ID == id{
            return &users[i], i
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
