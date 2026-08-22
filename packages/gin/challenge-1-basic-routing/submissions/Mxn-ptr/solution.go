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
	router := gin.Default()

	users := router.Group("/users")
	users.GET("", getAllUsers)
	users.POST("", createUser)
	users.GET("/:id", getUserByID)
	users.DELETE("/:id", deleteUser)
	users.PUT("/:id", updateUser)
	users.GET("/search", searchUsers)

	router.Run()
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(200, Response{
		Success: true,
		Data:    users,
		Message: "Users retrieved successfully",
		Code:    200,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "Invalid ID format",
		})
		return
	}

	user, _ := findUserByID(userId)
	if user == nil {
		c.JSON(404, Response{
			Success: false,
			Code:    404,
			Error:   "User not found",
		})
		return
	}

	c.JSON(200, Response{
		Success: true,
		Code:    200,
		Data:    user,
		Message: "User retrieved successfully",
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "Invalid JSON format",
		})
		return
	}
	if err := validateUser(newUser); err != nil {
		c.AbortWithStatusJSON(400, Response{
			Success: false,
			Code:    400,
			Error:   err.Error(),
		})
		return
	}

	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	c.JSON(201, Response{
		Success: true,
		Data:    newUser,
		Message: "User created successfully",
		Code:    201,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "Invalid ID format",
		})
		return
	}

	user, _ := findUserByID(userId)
	if user == nil {
		c.JSON(404, Response{
			Success: false,
			Code:    404,
			Error:   "User not found",
		})
		return
	}

	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "Invalid JSON format",
		})
		return
	}
	if input.Age != 0 {
		user.Age = input.Age
	}
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	c.JSON(200, Response{
		Success: true,
		Code:    200,
		Data:    user,
		Message: "User modified successfully",
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "Invalid ID format",
		})
		return
	}

	_, pos := findUserByID(userId)
	if pos == -1 {
		c.JSON(404, Response{
			Success: false,
			Code:    404,
			Error:   "User not found",
		})
		return
	}

	users = append(users[:pos], users[pos+1:]...)
	c.JSON(200, Response{
		Success: true,
		Code:    200,
		Message: "User deleted successfully",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "Missing parameter name",
		})
		return
	}
	name = strings.ToLower(name)
	usersFound := []User{}
	for _, v := range users {
		if strings.Contains(strings.ToLower(v.Name), name) {
			usersFound = append(usersFound, v)
		}
	}

	c.JSON(200, Response{
		Success: true,
		Code:    200,
		Data:    usersFound,
		Message: "List of users",
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for i, v := range users {
		if v.ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
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
