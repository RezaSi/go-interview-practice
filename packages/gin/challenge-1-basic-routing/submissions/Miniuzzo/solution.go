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
	router := gin.Default()

	// TODO: Setup routes
	// GET /users - Get all users
	router.GET("/users", getAllUsers)
	// GET /users/:id - Get user by ID
	router.GET("/users/:id", getUserByID)
	// POST /users - Create new user
	router.POST("/users", createUser)
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name

	// TODO: Start server on port 8080
	router.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	res := Response{
	    Success: true,
	    Data: users,
	}
	
	c.JSON(200, res)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	
	if err != nil {
	    c.JSON(400, Response{Success: false, Message:"Invalid ID format", Error: err.Error(),})
	    return
	}
	
	for _, user := range users {
	    if user.ID == id {
	        
	        res := Response{
	            Success: true,
	            Data: user,
	        }
	        
	        c.JSON(200, res)
	        return
	    }
	}
	
	c.JSON(404, Response{Success: false, Message:"User not found",})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	
	var newUser User
	// take parameter from body of http
    if err := c.ShouldBindJSON(&newUser); err != nil { // if i had some problem to parse JSON to struct
        c.JSON(400, Response{
            Success: false, 
            Message: "Invalid JSON format", 
            Error:   err.Error(),
        })
        return
    }
    
    if newUser.Name == "" || newUser.Email == "" {
        c.JSON(400, Response{
            Success: false, 
            Message: "Name and email cannot be empty",
        })
        return
    }
    
    newUser.ID = nextID
    nextID++
    
    // now i need to add user
    users = append(users, newUser)
    
    res := Response{
        Success: true,
        Data:    newUser,
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
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(404, Response{Success: false, Message: "Invalid ID format", Error: err.Error()})
        return
    }
	
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
	    c.JSON(404, Response{
            Success: false, 
            Message: "Invalid JSON format", 
            Error:   err.Error(),
        })
        return
	}
	
	for i := range users {
	    if users[i].ID == id{
	        users[i] = user
	        
	        c.JSON(200, Response{Success:true,Data:users[i],})
	        return
	    }
	}
	
	c.JSON(404, Response{
	    Success: false, 
	    Message: "ID not found", 
    })
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	
	if err != nil {
	    c.JSON(404, Response{
	    Success: false, 
	    Message: "Invalid ID format", 
	    Error: err.Error()})
        return
	}
	
	for i := range users {
	    if users[i].ID == id {
	        users = append(users[:i], users[i+1:]...)
	        
	        c.JSON(200, Response{
	            Success: true, 
	            Message: "User deleted with success!",
	        })
	        
	        return
	    } 
	}
	c.JSON(404, Response{
	    Success: false, 
	    Message: "ID not found"})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	
	urname := c.Query("name")
	matchUsers := []User{}
	
	if urname == "" {
        c.JSON(400, Response{Success: false, Message: "Query parameter 'name' is required"})
        return
    }
	
	for i := range users {
	    if strings.Contains(strings.ToLower(users[i].Name), strings.ToLower(urname)) {
	        matchUsers = append(matchUsers, users[i])
	    }
	}
	
	
	    c.JSON(200, Response{
	        Success: true, 
	        Data: matchUsers})
	
	
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
