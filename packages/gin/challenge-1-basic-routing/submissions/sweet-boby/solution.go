package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
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
	// TODO: Create Gin router
	router := gin.Default()
	// TODO: Setup routes
	// GET /users - Get all users
	router.GET("/user",getAllUsers)
	// GET /users/:id - Get user by ID
	router.GET("/user:id",getUserByID)	
	// POST /users - Create new user
	router.POST("/user:id",createUser)		
	// PUT /users/:id - Update user
	router.PUT("/user:id",updateUser)		
	// DELETE /users/:id - Delete user
	router.DELETE("/user:id",deleteUser)		
	// GET /users/search - Search users by name
    router.GET("/users/search",searchUsers)	
	// TODO: Start server on port 8080
	router.Run() // listens on 0.0.0.0:8080 by default
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
    c.JSON(200, gin.H{
        "success":true,
        "data" :users,
        "message" :"ok",
        "error"   :"",
        "code"   :200,
    })
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	id := c.Param("id")
    userID, err := strconv.Atoi(id)
    
    if err != nil {
            c.JSON(400, gin.H{
                "success":false,
                "data" :nil,
                "message" :"user not found",
                "error"   :"",
                "code"   :400,
            })
    }
    
    for _, user := range users {
            if user.ID == userID {
                c.JSON(200, gin.H{
                "success":true,
                "data" :user,
                "message" :"ok",
                "error"   :"",
                "code"   :200,
            })
            return
        }
    }
    c.JSON(404, gin.H{
                "success":false,
                "data" :nil,
                "message" :"user not found",
                "error"   :"",
                "code"   :404,
            })
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var newUser User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    
    if newUser.Name == ""{
        c.JSON(400, gin.H{"error": "err"})
        return
    }
    if newUser.Email == ""{
        c.JSON(400, gin.H{"error": "err"})
        return
    }
    if newUser.Age == 0{
        c.JSON(400, gin.H{"error": "err"})
        return
    }
    
    newUser.ID = len(users) + 1
    users = append(users, newUser)
    c.JSON(201, gin.H{
                "success":true,
                "data" :newUser,
                "message" :"ok",
                "error"   :"",
                "code"   :201,
            })
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
        c.JSON(400, gin.H{"error": "err"})
        return
	}
	// Parse JSON request body
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
	    c.JSON(400, gin.H{"error": err.Error()})
        return  
	}
	f := false
	// Find and update user
	for index , u := range users {
	    if u.ID == userID {
	        f = true
	        users[index].Name = user.Name
	        users[index].Email = user.Email
	        users[index].Age = user.Age
	    }
	}
	if f == false{
	    c.JSON(404, gin.H{"error": "err"})
        return  	    
	}
	// Return updated user
	c.JSON(200,gin.H{
	    "success":true,
	    "data" :user,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
        c.JSON(400, gin.H{"error": "err"})
        return
	}
	// Find and remove user
	for index , u := range users {
	    if u.ID == userID {
	        users := append(users[:index],users[index+1:]...)
	        c.JSON(200, gin.H{"success":true,
	            "data":users,
	        })
            return
	    }
	}
	// Return success message
	c.JSON(404, gin.H{"success":false,})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	name := c.Query("name") 
	name = strings.ToLower(name)
	fmt.Println(name)
	if name == ""{
	    c.JSON(400, gin.H{"success":false,
	    "error":"err",
    })
    return
	}
	// Filter users by name (case-insensitive)
	resUsers := []User{}
	for _ , u := range users {
	    uName := strings.ToLower(u.Name)
        if strings.Contains(uName,name) {
            resUsers = append(resUsers,u)
        }
	}
    c.JSON(200, gin.H{"success":true,
        "data":resUsers,
    })
	// Return matching users
    // c.JSON(404, gin.H{"success":false,})
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
