package main

import (
	"fmt"
	"net/http"
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
	// TODO: Create Gin router (Done)
    router := gin.Default()
    
	// TODO: Setup routes
	
    router.GET("/users", getAllUsers)
	// GET /users - Get all users
	
	router.GET("/users/:id", getUserByID)
	// GET /users/:id - Get user by ID
	
	router.POST("/users", createUser)
	// POST /users - Create new user
	
	router.PUT("/users/:id", updateUser)
	// PUT /users/:id - Update user
	
	router.DELETE("/users/:id", deleteUser)
	// DELETE /users/:id - Delete user
	
	router.GET("/users/search", searchUsers)
	// GET /users/search - Search users by name

	// TODO: Start server on port 8080
	router.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	
	c.JSON(http.StatusOK, Response{
	    Success: true,
	    Data: users,
	})
	
}

// getUserByID handles GET /users/:id (Done)
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	
	idStr := c.Param("id")
    // Convert id to int and find user
    id, err := strconv.Atoi(idStr)
    if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid user id format",
		})
		return
	}
	
	user, index := findUserByID(id)
	if index == -1 {
	    c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
		})
		return
	}
	
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// createUser handles POST /users (Done)
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
    if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	
    newUser.ID = len(users) + 1
    users = append(users, newUser)
    c.JSON(http.StatusCreated, Response{ 
		Success: true,
		Data:    newUser,
		Message: "successfully created new user",
	})
}

// updateUser handles PUT /users/:id (Done)
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "wrong user id format",
		})
		return
	}
	
	user, index := findUserByID(id)
	if index == -1 {
	    c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
		})
		return
	}
	
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invaild request data " + err.Error(),
		})
		return
	}
	
	updatedUser.ID = user.ID
	users[index] = updatedUser

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    updatedUser,
		Message: "successfully updated the user",
	})
}

// deleteUser handles DELETE /users/:id (Done)
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "wrong user id format",
		})
		return
	}
	
	_, index := findUserByID(id)
	if index == -1 {
	    c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
		})
		return
	}
	
	users = append(users[:index], users[index+1:]...)
	
	c.JSON(http.StatusOK, Response{
	    Success: true,
	    Message: "successfully deleted the user",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// 1. 取得 name 查詢參數
	nameQuery := c.Query("name")
	if nameQuery == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			// 我修正了你的一個小錯字
			Error:   "請提供 'name' 查詢參數",
		})
		// 重要：在回傳錯誤後，必須加上 return，
		// 否則函式會繼續往下執行，這是常見的錯誤。
		return
	}

	// 2. 準備一個空切片，用來存放符合條件的使用者
	matchedUsers :=  make([]User, 0)
	// 為了進行不分大小寫的比對，我們先把查詢字串轉成小寫
	nameQueryLower := strings.ToLower(nameQuery)

	// 3. 遍歷所有使用者，進行不分大小寫的「包含」比對
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), nameQueryLower) {
			matchedUsers = append(matchedUsers, user)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    matchedUsers,
	})
}

// Helper function to find user by ID (Done)
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found

	for index, user := range users {
        if user.ID == id {
            return &user, index
        }
    }
	return nil, -1
}

// Helper function to validate user data (Done)
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	
	if user.Name == "" {
	    return fmt.Errorf("Missig user's name")
	}
	
	if user.Email == "" {
	    return fmt.Errorf("Missing users's email")
	}
	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("Wrong email format")
	}
	
	return nil
}
