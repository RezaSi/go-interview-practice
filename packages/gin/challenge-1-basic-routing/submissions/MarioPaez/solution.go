package main

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

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

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var nextID = 4

func main() {
	router := gin.Default()
	router.GET("/users", getAllUsers)
	router.GET("users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/search", searchUsers)

	router.Run("localhost:8080")
}

func getAllUsers(c *gin.Context) {

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}

func getUserByID(c *gin.Context) {

	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "invalid ID",
		})
		return
	}
	user, index := findUserByID(userId)

	if index == -1 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

func createUser(c *gin.Context) {

	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	user.ID = nextID
	nextID++
	users = append(users, user)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
	})
}

func updateUser(c *gin.Context) {

	id := c.Param("id")
	var userToUpdate User

	if err := c.ShouldBindJSON(&userToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID",
		})
		return
	}

	if err := validateUser(userToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	_, index := findUserByID(userId)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Not found",
		})
		return
	}
	users[index].Age = userToUpdate.Age
	users[index].Name = userToUpdate.Name
	users[index].Email = userToUpdate.Email

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    userToUpdate,
	})
}

func deleteUser(c *gin.Context) {

	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID",
		})
		return
	}

	_, index := findUserByID(userId)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Not found",
		})
		return
	}
	users = slices.Delete(users, index, index+1)

	c.JSON(http.StatusOK, Response{
		Success: true,
	})
}

func searchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Name parameter is required",
		})
		return
	}
	var res = make([]User, 0)
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			res = append(res, user)
		}
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    res,
	})
}

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
	email := user.Email
	if email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(email, "@") && !strings.Contains(email, ".") {
		return errors.New("email must be valid")
	}

	return nil
}
