package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	router := gin.Default()

	router.GET("/users", getAllUsers)
	router.GET("/user/:id", getUserByID)
	router.GET("/users/search", searchUsers)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	if err := router.Run(":8080"); err != nil {

		fmt.Println("faild to start server: ", err)
	}
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(200, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "user id must be int",
		})
	}

	u, _ := findUserByID(id)

	if u == nil {
		c.JSON(404, Response{
			Success: false,
			Message: "user not found",
		})

		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    u,
	})

}

// createUser handles POST /users
func createUser(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, Response{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	var u struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	err = json.Unmarshal(bodyBytes, &u)
	if err != nil {
		c.JSON(500, Response{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	user := User{
		Name:  u.Name,
		Email: u.Email,
		Age:   u.Age,
	}

	if err = validateUser(user); err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: fmt.Sprintf("invalid user data: %s", err.Error()),
		})

		return
	}

	user.ID = nextID
	users = append(users, user)
	nextID += 1

	c.JSON(201, Response{
		Success: true,
		Data:    user,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "user id must be integer",
		})

		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, Response{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	var u struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	err = json.Unmarshal(bodyBytes, &u)
	if err != nil {
		c.JSON(500, Response{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	user, _ := findUserByID(id)

	if user == nil {
		c.JSON(404, Response{
			Success: false,
			Message: "user not found",
		})

		return
	}

	user.Name = u.Name
	user.Email = u.Email
	user.Age = u.Age

	c.JSON(200, Response{
		Success: true,
		Data:    user,
		// Message: fmt.Sprintf("%s updated", user.Name),
	})

}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "user id must be integer",
		})

		return
	}

	for i, u := range users {

		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(200, Response{
				Success: true,
				Data:    map[string]User{},
			})

			return
		}
	}

	c.JSON(404, Response{
		Success: false,
		Message: "user not found",
	})

}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	userName := c.Query("name")
	if userName == "" {
		c.JSON(400, Response{
			Success: false,
			Message: "user name must be provided",
		})

		return
	}

	for _, user := range users {

		if strings.EqualFold(user.Name, userName) || strings.Contains(strings.ToLower(user.Name), strings.ToLower(userName)) {
			c.JSON(200, Response{
				Success: true,
				Data:    []User{user},
			})

			return

		}
	}

	c.JSON(200, Response{
		Success: true,
		Data:    []string{},
		Message: "user not found",
	})

}

// // Helper function to find user by ID
func findUserByID(id int) (*User, int) {

	for idx, user := range users {
		if user.ID == id {
			return &user, idx

		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {

	if len(user.Name) == 0 || len(user.Email) == 0 {
		return fmt.Errorf("user name and email must be provided, len(name)= %d, len(email) = %d", len(user.Name), len(user.Email))
	}

	var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !EmailRX.MatchString(user.Email) {
		// if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
