package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	missingID error = errors.New("ID must not be empty")
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
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
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
	router.PUT("/users/:id", updateUser)
	// DELETE /users/:id - Delete user
	router.DELETE("/users/:id", deleteUser)
	// GET /users/search - Search users by name
	router.GET("/users/search", searchUsers)

	// TODO: Start server on port 8080	
	err := router.Run("localhost:8080")
	if err != nil {
		log.Printf("starting server: %v", err)
		return
	}
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	retrievedUsers := make([]User, len(users))		

	for i, user := range users { // is using loop here optimal? IDE tells me to use copy(from, to)
		retrievedUsers[i] = user
	}

	response := &Response{
		Success: true,
		Data: retrievedUsers,
		Message: "List of all users",
		Error: "",
		Code: 200,
	}

	c.JSON(200, response)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	id, err := retrieveID(c)
	if err != nil {	
			c.JSON(400, &Response{
				Success: false,
				Data: nil,
				Message: "",
				Error: err.Error(),
				Code: 400,
			})	
			return
		/*
		c.JSON(500, &Response{
			Success: false,
			Data: nil,
			Message: "",
			Error: "ID string conversion - Internal Server Error",
			Code: 500,
		})	
		return
		*/
	}
	// Return 404 if user not found
	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(404, &Response{
			Success: false,
			Data: nil,
			Message: "",
			Error: "User with the specified ID not found",
			Code: 404,
		})	
		return
	}

	response := &Response{
		Success: true,
		Data: user,
		Message: "User with specified ID",
		Error: "",
		Code: 200,
	}

	c.JSON(200, response)	
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var user User

	err := c.Bind(&user)
	if err != nil {
		c.JSON(400, &Response{
			Success: false,
			Data: nil,
			Message: "",
			Error: "binding the request body",
			Code: 400,
		})	
		return
	}
	// Validate required fields
	err = validateUser(user)
	if err != nil {
		c.JSON(400, &Response{
			Success: false,
			Data: nil,
			Message: "",
			Error: err.Error(),
			Code: 400,
		})	
		return	
	}
	// Add user to storage
	user.ID = nextID
	nextID++
	users = append(users, user)

	// Return created user
	response := &Response{
		Success: true,
		Data: user,
		Message: "user successfully created",
		Error: "",
		Code: 201,
	}

	c.JSON(201, response)
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, err := retrieveID(c)
	if err != nil {
		if errors.Is(err, missingID) {
			c.JSON(400, &Response{
				Success: false,
				Data: nil,
				Message: "",
				Error: err.Error(),
				Code: 400,
			})	
			return
		}
	}
	// Parse JSON request body
	var newUser User
	newUser.ID = id

	err = c.Bind(&newUser)
	if err != nil {
		c.JSON(400, &Response{
			Success: false,
			Data: nil,
			Message: "",
			Error: "binding the request body",
			Code: 400,
		})	
		return
	}
	// Find and update user
	user, idx := findUserByID(id)
	if user == nil {
		c.JSON(404, &Response{
			Success: false,
			Data: nil,
			Message: "",
			Error: "User with the specified ID not found",
			Code: 404,
		})	
		return
	}

	users[idx] = newUser

	// Return updated user
	c.JSON(200, &Response{
		Success: true,
		Data: users[idx],
		Message: "User updated successfully",
		Error: "",
		Code: 200,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, err := retrieveID(c)
	if err != nil {
		if errors.Is(err, missingID) {
			c.JSON(400, &Response{
				Success: false,
				Data: nil,
				Message: "",
				Error: err.Error(),
				Code: 400,
			})	
			return
		}
	}
	// Find and remove user
	user, idx := findUserByID(id)
	if user == nil {
		c.JSON(404, &Response{
			Success: false,
			Data: nil,
			Message: "",
			Error: "User with the specified ID not found",
			Code: 404,
		})	
		return
	}

	if idx == len(users)-1 {
		users = users[:idx]
	} else {
		users = append(users[:idx], users[idx+1:]...)
	}


	// Return success message
	c.JSON(200, &Response{
		Success: true,
		Data: nil,
		Message: "User deleted successfully",
		Error: "",
		Code: 200,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	name := c.Query("name")
	if name == "" {
		c.JSON(400, &Response{
			Success: false,
			Data: nil,
			Message: "",
			Error: "missing parameter 'name'",
			Code: 400,
		})	
		return
	}
	// Filter users by name (case-insensitive)
	result := []User{}

	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			result = append(result, user)
		} 
	}
	// Return matching users
	response := &Response{
		Success: true,
		Data: result,
		Message: "List of users with a specified name",
		Error: "",
		Code: 200,
	}

	c.JSON(200, response)
}

// Helper function to retrieve ID
func retrieveID(c *gin.Context) (int, error) {
	stringId := c.Param("id")
	if stringId == "" {
		return 0, errors.New("invalid ID")
	}

	id, err := strconv.Atoi(stringId)
	if err != nil {
		return 0, fmt.Errorf("string to int conversion: %w", err)
	}	

	return id, nil
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	for i, user := range users {
		if user.ID == id {	
			return &user, i
		}
	}
	// Return user pointer and index, or nil and -1 if not found
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	if user.Email == "" || user.Name == "" {
		return errors.New("missing required fields")
	}
	// Validate email format (basic check)
	re := regexp.MustCompile(`^[a-zA-z0-9._%+-]+@[a-zA-z.-]+\.[a-zA-z]{3,}`)
	valid := re.MatchString(user.Email)

	if !valid {
		return errors.New("invalid email")
	}

	return nil
}
