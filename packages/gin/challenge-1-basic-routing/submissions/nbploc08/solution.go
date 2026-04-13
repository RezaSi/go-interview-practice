package main

import (
	"errors"
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

	// TODO: Setup routes
	// GET /users - Get all users
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name

	// TODO: Start server on port 8080

	r := gin.Default()
	r.GET("/users", getAllUsers)
	r.GET("/users/:id", getUserByID)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)
	r.GET("/users/search", searchUsers)
	r.Run() // listens on 0.0.0.0:8080 by default
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	resp := Response{
		Success: true,
		Data:    users,
		Message: "ssss",
		Code:    200,
	}
	c.JSON(200, resp)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	user, _ := findUserByID(userID)
	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    user,
		Message: "Get user success",
		Code:    200,
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
	err := validateUser(user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return // ✅ Fix: phải return, không thì tiếp tục append user lỗi
	}
	user.ID = nextID // ✅ Fix: gán ID cho user mới
	nextID++
	users = append(users, user)
	c.JSON(201, Response{Success: true, Data: user, Message: "Create user success", Code: 201}) // ✅ Fix: 201 Created

}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return // ✅ Fix: return để dừng, không chạy tiếp với userID rác
	}
	var userData User

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if vali := validateUser(userData); vali != nil {
		c.JSON(400, gin.H{"error": vali.Error()}) // ✅ Fix: dùng vali.Error() thay vì err.Error() (err ở đây là nil), thêm return
		return
	}
	user, i := findUserByID(userID)
	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	userData.ID = userID // ✅ Fix: giữ nguyên ID từ URL, không để client tự đặt
	users[i] = userData
	c.JSON(200, Response{
		Success: true,
		Data:    userData,
		Message: "Update user success",
		Code:    200,
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
	user, i := findUserByID(userID)
	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return // ✅ Fix: phải return, không thì tiếp tục với i = -1 => panic
	}
	users = append(users[:i], users[i+1:]...) // ✅ Fix: users[i+1:] thay vì users[:i+1], và gán lại users

	c.JSON(200, Response{
		Success: true,
		Message: "Delete user success",
		Code:    200,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	search := c.Query("name")
	search = strings.TrimSpace(search)
	if search == "" {
		c.JSON(400, gin.H{"error": "No search query"})
		return
	}
	var matches []User
	search = strings.ToLower(search)
	for _, user := range users {
		userName := strings.ToLower(user.Name)
		if strings.Contains(userName, search) { // ✅ Fix: Contains thay vì ==, tìm kiếm linh hoạt hơn
			matches = append(matches, user)
		}
	}
	if matches == nil {
		matches = []User{} // trả về mảng rỗng thay vì nil
	}
	c.JSON(200, gin.H{"success": true, "data": matches}) // ✅ Fix: luôn trả 200, kể cả khi không có kết quả

}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	for i, user := range users {
		if user.ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	if user.Name == "" || user.Email == "" || user.Age <= 0 {
		return errors.New("invalid user")
	}
	return nil
}
