package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrEmptyEmail         = errors.New("empty email")
	ErrEmptySearch        = errors.New("no search parameter provided")
	ErrEmptyUsername      = errors.New("empty username")
	ErrInvalidAge         = errors.New("ivalid age")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrInvalidIDFormat    = errors.New("invalid id format")
	ErrInvalidUserFields  = errors.New("invalid user fields")
	ErrUserNotFound       = errors.New("user not found")
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
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
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
	{
		users.GET("search", searchUsers)
		users.GET(":id", getUserByID)
		users.GET("", getAllUsers)
		users.POST("", createUser)
		users.PUT(":id", updateUser)
		users.DELETE(":id", deleteUser)
	}

	router.Run(":8080")
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	u := users

	resp := Response{
		Success: true,
		Message: "users list",
		Data:    u,
		Code:    http.StatusOK,
	}

	c.JSON(http.StatusOK, resp)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: ErrInvalidIDFormat.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	user, found := findUserByID(id)
	if found == -1 {
		c.JSON(http.StatusNotFound, Response{
			Error: ErrUserNotFound.Error(),
			Code:  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "user found",
		Code:    http.StatusOK,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var newUser User
	var err error
	if err = c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: ErrInvalidUserFields.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	if err = validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	newUser.ID = nextID
	users = append(users, newUser)
	nextID++

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "user created",
		Code:    http.StatusCreated,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	var err error
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: ErrInvalidIDFormat.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	user, found := findUserByID(id)
	if found == -1 {
		c.JSON(http.StatusNotFound, Response{
			Error: ErrUserNotFound.Error(),
			Code:  http.StatusNotFound,
		})
		return
	}

	var updateUser User
	if err = c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: ErrInvalidUserFields.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	if err = validateUser(updateUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	user.Age = updateUser.Age
	user.Email = updateUser.Email
	user.Name = updateUser.Name

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "successful user update",
		Code:    http.StatusOK,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	var err error
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: ErrInvalidIDFormat.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	_, found := findUserByID(id)
	if found == -1 {
		c.JSON(http.StatusNotFound, Response{
			Error: ErrUserNotFound.Error(),
			Code:  http.StatusNotFound,
		})
		return
	}

	users = append(users[:found], users[found+1:]...)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "successful user deletion",
		Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	username := strings.ToLower(c.Query("name"))
	if username == "" {
		c.JSON(http.StatusBadRequest, Response{
			Error: ErrEmptySearch.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	usersFound := []User{}
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), username) {
			usersFound = append(usersFound, u)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    usersFound,
		Message: "search results",
		Code:    http.StatusOK,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	var user *User
	index := -1
	for i, u := range users {
		if u.ID == id {
			user = &users[i]
			index = i
			break
		}
	}

	return user, index
}

// Helper function to validate user data
func validateUser(user User) error {
	if user.Name == "" {
		return ErrEmptyUsername
	}

	if user.Age < 0 {
		return ErrInvalidAge
	}

	if user.Email == "" {
		return ErrEmptyEmail
	}

	if !strings.Contains(user.Email, "@") {
		return ErrInvalidEmailFormat
	}

	return nil
}
