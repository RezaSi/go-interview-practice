package main

import (
	"errors"
	"net/http"
	"regexp"
	"slices"
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
	Success bool    `json:"success"`
	Data    any     `json:"data,omitempty"`
	Message string  `json:"message,omitempty"`
	Error   string  `json:"error,omitempty"`
	Code    int     `json:"code,omitempty"`
}

// In-memory storage
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}

var nextID = 4

// ---------------------------------------------------------------
// Main
// ---------------------------------------------------------------

func main() {
	r := gin.Default()
	r.GET("/users/search", searchUsers)
	r.GET("/users", getAllUsers)
	r.GET("/users/:id", getUserByID)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)
	r.Run(":8080")
}

// ---------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------

func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Success: true, Data: users})
}

func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format",
		})
		return
	}

	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{Success: true, Data: user})
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	newUser.ID = nextID
	nextID++
	users = append(users, newUser)

	c.JSON(http.StatusCreated, Response{Success: true, Data: newUser})
}

func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID",
		})
		return
	}

	var userData User
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if err := validateUser(userData); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Not found",
		})
		return
	}

	userData.ID = id
	users[index] = userData

	c.JSON(http.StatusOK, Response{Success: true, Data: userData})
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID",
		})
		return
	}

	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Not found",
		})
		return
	}

	users = slices.Delete(users, index, index + 1)
	c.JSON(http.StatusOK, Response{Success: true})
}

func searchUsers(c *gin.Context) {
	name := strings.ToLower(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Name parameter is required",
		})
		return
	}

	results := make([]User, 0)
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), name) {
			results = append(results, user)
		}
	}

	c.JSON(http.StatusOK, Response{Success: true, Data: results})
}

// ---------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------

func findUserByID(id int) (*User, int) {
	for i, user := range users {
		if user.ID == id {
			return &user, i
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
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if ! re.MatchString(user.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
