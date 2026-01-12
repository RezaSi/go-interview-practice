// Gin Web Framework 1

package main

import (
	"errors"
	"strconv"
	"strings"
    "sync"
    
	"github.com/gin-gonic/gin"
)

// gin.H = map[string]interface{}
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var nextID = 4
var mu sync.Mutex
func main() {
	router := gin.Default()
	router.GET("/users", getAllUsers)
		router.GET("/users/search", searchUsers)
	router.GET("/users/:id", getUserByID)

	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.Run(":8080")
}
func getAllUsers(c *gin.Context) {
    mu.Lock()
	c.JSON(200, Response{Success: true, Data: users, Message: "All users"})
    mu.Unlock()
}
func getUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{Success: false, Error: "Invalid ID"})
		return
	}
	user, ind := findUserByID(userID)
	if ind != -1 {
		c.JSON(200, Response{
			Success: true,
			Data:    user,
			Message: "Users retrieved successfully",
			Code:    200})
		return
	}
	c.JSON(404, Response{Success: false,Error:   "User not found"})
}

func findUserByID(id int) (*User, int) {
	for ind, user := range users {
		if user.ID == id {
			return &users[ind], ind//&users[ind] not &user
		}
	}
	return nil, -1
}
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
func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, Response{Success: false, Error: err.Error()})
		return
	}
	
	if err := validateUser(newUser); err != nil {
		c.JSON(400, Response{Success: false, Error: err.Error()})
		return
	}
	mu.Lock()
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	mu.Unlock()
	c.JSON(201, Response{Success: true, Data: newUser, Message: "User created"})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{Success: false, Error: "Invalid ID"})
		return
	}
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, Response{Success: false, Error: err.Error()})
		return
	}

	if err := validateUser(newUser); err != nil {
		c.JSON(400, Response{Success: false, Error: err.Error()})
		return
	}

	user, ind := findUserByID(userID)
	if ind != -1 {
	      mu.Lock()
	      user.Name = newUser.Name
	      user.Age = newUser.Age
	      user.Email = newUser.Email
	      mu.Unlock()
	      c.JSON(200, Response{ Success: true, Data:    user, Message: "Users updated successfully", Code:    200})
	      return
	}
	c.JSON(404, Response{Success: false, Error:   "User not found"})
	
}
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{Success: false, Error: "Invalid ID", Code: 400})
		return
	}
	_, ind := findUserByID(userID)
	if ind != -1 {
	    mu.Lock()
		users[ind] = users[len(users)-1]
		users = users[:len(users)-1]
		mu.Unlock()
		c.JSON(200, Response{
			Success: true,
			Message: "Users deleted successfully",
			Code:    200})
	} else {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
	}
}

// /users/search?name=value
 func searchUsers(c *gin.Context) {
	queryName := c.DefaultQuery("name", "")
	if queryName == "" {
 		c.JSON(400, Response{Success: false, Code: 400, Error: "provide name in url"})
 		return
 	}
        mu.Lock()
	matchedUsers := []User{}
 	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(queryName)) {
			matchedUsers = append(matchedUsers, user)
 		}
 	}
 	mu.Unlock()
 	c.JSON(200, Response{Success: true, Data:    matchedUsers,Message: "matched",})
 }