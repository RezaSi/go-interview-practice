package main

import (
	"errors"
	"regexp"
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
	r := gin.Default()
	//r := gin.New()
	// TODO: Setup routes
	//r.Use(gin.Logger())
	//r.Use(JSONRecovery())
	// GET /users - Get all users
	r.GET("/users", getAllUsers)
	// GET /users/:id - Get user by ID
	r.GET("/users/:id", getUserByID)
	// POST /users - Create new user
	r.POST("/users", createUser)
	// PUT /users/:id - Update user
	r.PUT("/users/:id", updateUser)
	// DELETE /users/:id - Delete user
	r.DELETE("/users/:id", deleteUser)
	// GET /users/search - Search users by name
	r.GET("/users/search", searchUsers)

	// TODO: Start server on port 8080
	r.Run(":8080")
}

// TODO: Implement handler functions
//func JSONRecovery() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		defer func() {
//			if err := recover(); err != nil {
//				c.JSON(400, Response{Code: 400, Error: fmt.Sprint(err)})
//				c.Abort()
//			}
//		}()
//		c.Next()
//	}
//}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(200, Response{true, users, "success", "", 200})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// Handle invalid ID format
	// TODO: Get user by ID
	id, ok := strconv.Atoi(c.Param("id"))
	if ok != nil {
		//panic(errors.New("invalid id"))
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "invalid id",
		})
		return
	}
	// Return 404 if user not found
	u, _ := findUserByID(id)
	if u == nil {
		//c.JSON(404, Response{Code: 404})
		c.JSON(404, Response{
			Success: false,
			Code:    404,
			Error:   "user not found",
		})
		return
	}
	c.JSON(200, Response{true, u, "success", "", 200})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	u := User{}
	if err := c.ShouldBind(&u); err != nil {
		//panic(err)
		c.JSON(400, Response{
			Code:  400,
			Error: err.Error(),
		})
		return
	}
	// Validate required fields
	if err := validateUser(u); err != nil {
		//panic(err)
		c.JSON(400, Response{
			Code:  400,
			Error: err.Error(),
		})
		return
	}
	// Add user to storage
	// Return created user
	u.ID = nextID
	nextID++
	users = append(users, u)
	c.JSON(201, Response{true, u, "success", "", 200})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, ok := strconv.Atoi(c.Param("id"))
	if ok != nil {
		//panic(errors.New("invalid id"))
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "invalid id",
		})
		return
	}
	// Parse JSON request body
	u := User{}
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(400, Response{
			Code:  400,
			Error: err.Error(),
		})
		return
	}
	// Find and update user
	delegate, idx := findUserByID(id)
	if delegate == nil {
		c.JSON(404, Response{Code: 404})
		return
	}
	u.ID = id
	users[idx] = u
	// Return updated user
	c.JSON(200, Response{true, u, "success", "", 200})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	id, ok := strconv.Atoi(c.Param("id"))
	if ok != nil {
		//panic(errors.New("invalid id"))
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "invalid id",
		})
		return
	}
	// Find and remove user
	u, idx := findUserByID(id)
	if u == nil {
		c.JSON(404, Response{Code: 404})
		return
	}
	users = append(users[:idx], users[idx+1:]...)
	// Return success message
	c.JSON(200, Response{true, users, "success", "", 200})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	n := c.Query("name")
	if n == "" {
		c.JSON(400, Response{Code: 400})
		return
	}
	// Filter users by name (case-insensitive)
	var arr []User
	for _, user := range users {
		for _, s := range strings.Split(user.Name, " ") {
			if strings.EqualFold(s, n) {
				arr = append(arr, user)
				break
			}
		}
	}
	if arr == nil || len(arr) == 0 {
		c.JSON(200, Response{Success: true, Data: []string{}, Code: 200})
		return
	}
	// Return matching users
	c.JSON(200, Response{true, arr, "success", "", 200})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	var u *User
	var idx int
	for i, user := range users {
		if user.ID == id {
			u = &user
			idx = i
			break
		}
	}
	if u != nil {
		return u, idx
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	if user.Name == "" || user.Email == "" {
		return errors.New("missing name or email")
	}
	// Validate email format (basic check)
	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
	)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
