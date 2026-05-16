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
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
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
	router := gin.Default()

	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/search", searchUsers)

	router.Run(":8080")
}

func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Message: "Successfully retrieved all users",
		Code:    200,
	})
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Message: "Invalid ID format", Error: err.Error(), Code: 400})
		return
	}

	// 【重构】使用 findUserByID
	user, _ := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{Success: false, Message: "User not found", Code: 404})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "User found successfully",
		Code:    200,
	})
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Message: "Invalid JSON format", Error: err.Error(), Code: 400})
		return
	}

	// 【新增】使用 validateUser
	if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Message: "Validation failed", Error: err.Error(), Code: 400})
		return
	}

	newUser.ID = nextID
	nextID++
	users = append(users, newUser)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "User created successfully!",
		Code:    201,
	})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Message: "Invalid ID format", Code: 400})
		return
	}

	// 【重构】使用 findUserByID
	_, index := findUserByID(userID)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{Success: false, Message: "User updated unsuccessfully", Error: "User not found", Code: 404})
		return
	}

	var incomingData User
	if err := c.ShouldBindJSON(&incomingData); err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Message: "Invalid JSON format", Error: err.Error(), Code: 400})
		return
	}

	// 【新增】使用 validateUser
	if err := validateUser(incomingData); err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Message: "Validation failed", Error: err.Error(), Code: 400})
		return
	}

	// 更新逻辑 (保持 ID 不变)
	incomingData.ID = userID
	users[index] = incomingData

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users[index],
		Message: "User updated successfully!",
		Code:    200,
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Message: "Invalid ID format", Code: 400})
		return
	}

	// 【重构】使用 findUserByID
	_, index := findUserByID(userID)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{Success: false, Message: "User deletion unsuccessfully", Error: "User not found", Code: 404})
		return
	}

	users = append(users[:index], users[index+1:]...)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully!",
		Code:    200,
	})
}

// 修复了两个 Bug 的 searchUsers 函数
func searchUsers(c *gin.Context) {
	nameQuery := c.Query("name")

	// 【修复 Bug 2】: 检查参数是否为空
	if nameQuery == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Query parameter 'name' is required",
			Code:    400,
		})
		return
	}

	matchedUsers := []User{}

	for _, user := range users {
		// 【修复 Bug 1】: 使用 Contains 代替 ==，实现模糊搜索
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(nameQuery)) {
			matchedUsers = append(matchedUsers, user)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    matchedUsers,
		Message: "Users retrieved successfully",
		Code:    200,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for i := 0; i < len(users); i++ {
		if users[i].ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	// 简单的 Email 格式检查
	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}