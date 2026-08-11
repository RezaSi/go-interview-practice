package main

import (
	"net/http"
	"strconv"
	"errors"
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
    r := gin.Default()
	// GET /users - Get all users
	r.GET("/users",getAllUsers)
	// GET /users/:id - Get user by ID
	r.GET("/users/:id",getUserByID)
	// POST /users - Create new user
	r.POST("/users",createUser)
	// PUT /users/:id - Update user
	r.PUT("/users/:id",updateUser)
	// DELETE /users/:id - Delete user
	r.DELETE("/users/:id",deleteUser)
	// GET /users/search - Search users by name

	r.Run(":8000")
}


// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK,Response{
	    Success: true,
	    Data:users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// Handle invalid ID format
	// Return 404 if user not found
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err !=nil {
	    c.JSON(http.StatusBadRequest, Response{
	        Success: false,
	        Error: "",
	        Code: http.StatusBadRequest,
	    })
	    return
	}
	
	user, _ := findUserByID(id)
	
	if user == nil {
	    c.JSON(http.StatusNotFound, Response{
	        Success :false,
	        Error:"",
	        Code : http.StatusNotFound,
	    })
	    return
	}
	c.JSON(http.StatusOK, Response{
	    Success: true,
	    Data: user,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// Validate required fields
	// Add user to storage
	// Return created user
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
	    c.JSON(http.StatusBadRequest, Response{
	        Success: false,
	        Error: err.Error(),
	        Code: http.StatusBadRequest,
	    })
	    return
	}
	
	if err := validateUser(newUser); err !=nil {
	    c.JSON(http.StatusBadRequest, Response{
	        Success: false,
	        Error: err.Error(),
	        Code: http.StatusBadRequest,
	    })
	    return
	}
	
	newUser.ID = nextID
	nextID ++
	users = append(users, newUser)
	
	c.JSON(http.StatusCreated, Response{
	    Success: true,
	    Data: newUser,
	    Message:"",
	})
	
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// Parse JSON request body
	// Find and update user
	// Return updated user
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err !=nil {
	    c.JSON(http.StatusBadRequest,Response{
	        Success: false,
	        Error:"",
	        Code: http.StatusBadRequest,
	    })
	    return
	}
	_,index := findUserByID(id)
	if index == -1 {
	    c.JSON(http.StatusNotFound, Response{
	        Success: false,
	        Error: "",
	        Code: http.StatusNotFound,
	    })
	    return
	}
	
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
	    c.JSON(http.StatusBadRequest, Response{
	        Success: false,
	        Error: err.Error(),
	        Code: http.StatusBadRequest,
	    })
	    return
	}
	if err := validateUser(updatedUser); err !=nil {
	    c.JSON(http.StatusBadRequest, Response{
	        Success: false,
	        Error: err.Error(),
	        Code: http.StatusBadRequest,
	    })
	    return
	}
	
	updatedUser.ID = id
	users[index] = updatedUser
	
	c.JSON(http.StatusOK,Response{
	    Success:true,
	    Data: updatedUser,
	    Message: "",
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// Find and remove user
	// Return success 
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
	    c.JSON(http.StatusBadRequest, Response{
	        Success: false,
	        Error: "",
	        Code: http.StatusNotFound,
	    })
	    return 
	}
	_,index := findUserByID(id)
	if index == -1 {
	    c.JSON(http.StatusNotFound,Response{
	        Success: false,
	        Error:"",
	        Code: http.StatusNotFound,
	    })
	    return
	}
	users = append(users[:index],users[index+1:]...)
	c.JSON(http.StatusOK,Response{
	    Success:true,
	    Message:"",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// Filter users by name (case-insensitive)
	// Return matching users
	nameQuery := c.Query("name")
	if strings.TrimSpace(nameQuery)==""{
	    c.JSON(http.StatusBadRequest,Response{
	        Success:false,
	        Error:"",
	        Code: http.StatusBadRequest,
	    })
	    return
	}
	matchingUsers := make([]User,0)
	
	for _,user :=range users {
	    if strings.Contains(strings.ToLower(user.Name),strings.ToLower(nameQuery)){
	        matchingUsers = append(matchingUsers,user)
	    }
	}
	c.JSON(http.StatusOK,Response{
	    Success:true,
	    Data: matchingUsers,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// Return user pointer and index, or nil and -1 if not found
	for i, user :=range users {
	    if user.ID == id {
	        return &users[i],i
	    }
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// Check required fields: Name, Email
	// Validate email format (basic check)
	if strings.TrimSpace(user.Name) == ""{
	    return errors.New("name is required")
	}
	if strings.TrimSpace(user.Email) == ""{
	    return errors.New("Email is required")
	}
	if !strings.Contains(user.Email,"@"){
	    return errors.New("invalid email format")
	}
	return nil
	return nil
}
