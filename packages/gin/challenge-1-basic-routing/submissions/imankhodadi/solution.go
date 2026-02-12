package main

import (
	"errors"
	"net/mail"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

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
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
func getAllUsers(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(200, Response{Success: true, Data: users, Message: "All users"})
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{Success: false, Error: "Invalid ID"})
		return
	}
	mu.Lock()
	user, ind := findUserByID(userID)
	var userData *User
	if ind != -1 {
		// Copy user data to avoid holding lock during response
		userCopy := *user
		userData = &userCopy
	}
	mu.Unlock()
	if userData != nil {
		c.JSON(200, Response{Success: true, Data: userData, Message: "Users retrieved successfully"})
		return
	}
	c.JSON(404, Response{Success: false, Error: "User not found"})
}

func findUserByID(id int) (*User, int) {
	for ind, user := range users {
		if user.ID == id {
			return &users[ind], ind
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
	if user.Age < 0 {
		return errors.New("age cannot be negative")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
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

	mu.Lock()
	user, ind := findUserByID(userID)
	if ind != -1 {
		user.Name = newUser.Name
		user.Age = newUser.Age
		user.Email = newUser.Email
		userCopy := *user
		mu.Unlock()
		c.JSON(200, Response{Success: true, Data: userCopy, Message: "User updated successfully"})
		return
	}
	mu.Unlock()
	c.JSON(404, Response{Success: false, Error: "User not found"})

}
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, Response{Success: false, Error: "Invalid ID"})
		return
	}
	mu.Lock()
	_, ind := findUserByID(userID)
	if ind != -1 {
		users[ind] = users[len(users)-1]
		users = users[:len(users)-1]
		//users = append(users[:ind], users[ind+1:]...) if order matters
		mu.Unlock()
		c.JSON(200, Response{Success: true, Message: "User deleted successfully"})
	} else {
		mu.Unlock()
		c.JSON(404, Response{Success: false, Error: "User not found"})
	}
}

// /users/search?name=value
func searchUsers(c *gin.Context) {
	queryName := c.DefaultQuery("name", "")
	if queryName == "" {
		c.JSON(400, Response{Success: false, Error: "provide name in url"})
		return
	}
	mu.Lock()
	usersCopy := make([]User, len(users))
	copy(usersCopy, users)
	mu.Unlock()
	matchedUsers := []User{}
	for _, user := range usersCopy {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(queryName)) {
			matchedUsers = append(matchedUsers, user)
		}
	}
	c.JSON(200, Response{Success: true, Data: matchedUsers, Message: "matched"})
}

