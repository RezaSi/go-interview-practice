package main

import (
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
	router := gin.Default()

	// TODO: Setup routes

	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/search", searchUsers)

	router.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Data:    nil,
		})
	}
	u, _ := findUserByID(id)
	if u != nil {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Data:    u,
		})
	} else {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Data:    nil,
		})
	}
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	type Req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	if req.Name == "" || req.Email == "" || req.Age == 0 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
		})
		return
	}
	u := User{
		ID:    4,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	users = append(users, u)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    u,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	type Req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
	}

	u, _ := findUserByID(id)
	if u == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Data:    nil,
		})
		return
	}
	u.Name = req.Name
	u.Email = req.Email
	u.Age = req.Age
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    u,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
	}
	filtered := users[:0]
	for _, u := range users {
		if id != u.ID {
			filtered = append(filtered, u)
		}
	}
	if len(filtered) == len(users) {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	v := c.Query("name")
	if v == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
		})
		return
	}
	filtered := users[:0]
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(v)) {
			filtered = append(filtered, u)
		}
	}
	c.JSON(http.StatusOK, Response{Success: true, Data: filtered})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	var u *User
	for i := range users {
		if id == users[i].ID {
			u = &users[i]
			break
		}
	}
	if u != nil {
		return u, 0
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	return nil
}