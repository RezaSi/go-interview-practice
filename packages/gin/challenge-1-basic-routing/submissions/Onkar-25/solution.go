package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
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
	router.GET("/users",getAllUsers)
	// GET /users/:id - Get user by ID
	router.GET("/users/:id",getUserByID)
	// POST /users - Create new user
	router.POST("/users",createUser)
	// PUT /users/:id - Update user
	router.PUT("/users/:id",updateUser)
	// DELETE /users/:id - Delete user
	router.DELETE("/users/:id",deleteUser)
	// GET /users/search - Search users by name
    router.GET("/users/search",searchUsers)
	// TODO: Start server on port 8080
	router.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	response:= Response{
	    Success : true,
	    Data : users,
	    Error : "",
	    Code : 200,
	}
	c.JSON(200, response)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	idParam := c.Param("id")

	userID, err := strconv.Atoi(idParam)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid user ID",
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, index := findUserByID(userID)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "User not found",
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Code:    http.StatusOK,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var newUser User
	if err:= c.ShouldBindJSON(&newUser);err!=nil{
	    c.JSON(400, Response{
			Success: false,
			Message: "Invalid user format",
			Error:   "Invalid user format",
			Code:    http.StatusBadRequest,
		})
        return
	}
	
	// Validate required fields
	err := validateUser(newUser)
	if err != nil{
	    c.JSON(400,  Response{
			Success: false,
			Message: "Invalid user Details",
			Error:   "Invalid user Details",
			Code:    http.StatusBadRequest,
		})
        return
	}
	// Add user to storage
	newUser.ID = len(users) + 1
	users = append(users,newUser)
	// Return created user
	c.JSON(http.StatusCreated, Response{
			Success: true,
		    Data: newUser,
			Code:    http.StatusOK,
		})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid user ID",
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	_, index := findUserByID(userID)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "User not found",
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	// Parse JSON request body
	var updatedDetails User
	if err:= c.ShouldBindJSON(&updatedDetails);err!=nil{
	    c.JSON(400, Response{
			Success: false,
			Message: "Invalid user format",
			Error:   "Invalid user format",
			Code:    http.StatusBadRequest,
		})
        return
	}
	// Find and update user
	users[index].Name = updatedDetails.Name
	users[index].Email = updatedDetails.Email
	users[index].Age = updatedDetails.Age
	// Return updated user
	c.JSON(http.StatusOK,Response{
			Success: true,
		    Data: users[index],
			Code:    http.StatusOK,
		})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid user ID",
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	_, index := findUserByID(userID)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "User not found",
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	// Find and remove user
	users = append(users[:index], users[index+1:]...)
	// Return success message
	c.JSON(http.StatusOK,Response{
			Success: true,
			Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	userName := c.Query("name")
	if userName == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid user Name",
			Error:   "Invalid user Name",
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Filter users by name (case-insensitive)
	Matchingusers := []User{}
	for _,user:= range users{
	    if strings.Contains(strings.ToLower(user.Name),strings.ToLower(userName)){
	        Matchingusers = append(Matchingusers,user)
	    }
	}
	// Return matching users
	c.JSON(http.StatusOK,Response{
			Success: true,
			Data : Matchingusers,
			Code:    http.StatusOK,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	for index, user := range users{
	    if user.ID == id{
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
