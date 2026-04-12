package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"regexp"
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
	r := gin.Default()
	r.GET("/", getAllUsers)
	r.GET("/users/:id", getUserByID)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)
	r.GET("/users/search", searchUsers)

	// TODO: Setup routes
	// GET /users - Get all users
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name

	// TODO: Start server on port 8080
	c := r.Run(":8080")
	if c != nil {
		panic(c)
	}
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(200, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    404,
		})
		return
	}
	if id < 1 || id > int64(len(users)) {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}
	for i := 0; i < len(users); i++ {
		if users[i].ID == int(id) {
			c.JSON(200, Response{
				Success: true,
				Data:    users[i],
			})
			return
		}
	}
	c.JSON(404, Response{
		Success: true,
		Data:    nil,
	})
	// Handle invalid ID format
	// Return 404 if user not found
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(404, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}
	if user.Name == "" || user.Email == "" || user.Age == 0 {
		c.JSON(400, Response{
			Success: false,
			Error:   "Missing required fields",
			Code:    400,
		})
		return
	}
	user.ID = nextID
	users = append(users, user)
	c.JSON(201, Response{
		Success: true,
		Data:    user,
	})
	return
	// Validate required fields
	// Add user to storage
	// Return created user
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(404, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    404,
		})
		return
	}

	// Parse JSON request body
	var updatedUser User
	if id < 1 || id > int64(len(users)) {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(404, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}
	updatedUser.ID = int(id)
	users[id] = updatedUser
	c.JSON(200, Response{
		Success: true,
		Data:    updatedUser,
	})
	return
	// Find and update user
	// Return updated user
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(404, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    404,
		})
		return
	}
	if id < 1 || id > int64(len(users)) {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
	}
	for i := 0; i < len(users); i++ {
		if users[i].ID == int(id) {
			users = append(users[:i], users[i+1:]...)
			c.JSON(200, Response{
				Success: true,
				Message: "User deleted successfully",
			})
			return
		}

	}
	c.JSON(404, Response{
		Success: false,
		Error:   "User not found",
		Code:    404,
	})
	// Find and remove user
	// Return success message
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")

	// 优化：如果 name 参数为空，可以返回所有用户或一个错误
	if name == "" {
		c.JSON(400, Response{
			Success: false,
			Data:    "Query parameter 'name' is required",
		})
		return
	}

	var matchingUsers []User // 创建一个空切片，用于存放所有匹配的用户

	// 遍历所有用户，进行不区分大小写的模糊搜索
	for _, user := range users {
		// 使用 strings.EqualFold 进行不区分大小写的精确匹配
		// 或者继续使用 strings.Contains 进行模糊搜索，取决于你的需求
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			matchingUsers = append(matchingUsers, user)
		}
	}

	// 无论有没有找到，都返回一个列表
	// 如果 matchingUsers 为空，它就是 []，符合测试的期望
	if len(matchingUsers) == 0 {
		matchingUsers = []User{}
	}
	c.JSON(200, Response{
		Success: true,
		Data:    matchingUsers,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	for i := 0; i < len(users); i++ {
		if users[i].ID == id {
			return &users[i], i
		}
		// Return user pointer and index, or nil and -1 if not found

	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// 1. 检查必填字段
	if user.Name == "" {
		return errors.New("name is a required field")
	}
	if user.Email == "" {
		return errors.New("email is a required field")
	}

	// 2. 验证邮箱格式（基本检查）
	if !isValidEmail(user.Email) {
		return errors.New("email format is invalid")
	}

	// 3. 如果所有检查都通过，返回 nil 表示成功
	return nil
}

// isValidEmail 使用正则表达式进行基本的邮件格式验证
func isValidEmail(email string) bool {
	// 这是一个非常常用且实用的基础邮件正则
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
