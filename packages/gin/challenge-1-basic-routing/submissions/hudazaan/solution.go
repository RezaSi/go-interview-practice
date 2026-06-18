package main

import (
	"fmt"
	"net/http"
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

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var nextID = 4

func main() {
	router := gin.Default()

	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/search", searchUsers)

	router.Run(":8080")
}

func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Code:    http.StatusOK,
	})
}

func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Code:    http.StatusOK,
	})
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid JSON input",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	newUser.ID = nextID
	nextID++
	users = append(users, newUser)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "User created successfully",
		Code:    http.StatusCreated,
	})
}

func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, index := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid JSON input",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if updateData.Name != "" {
		users[index].Name = updateData.Name
	}
	if updateData.Email != "" {
		users[index].Email = updateData.Email
	}
	if updateData.Age > 0 {
		users[index].Age = updateData.Age
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users[index],
		Message: "User updated successfully",
		Code:    http.StatusOK,
	})
}

func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, index := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	users = append(users[:index], users[index+1:]...)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
		Code:    http.StatusOK,
	})
}

func searchUsers(c *gin.Context) {
	nameQuery := c.Query("name")
	if nameQuery == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Name query parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	results := []User{} // Initialize as empty slice, not nil
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(nameQuery)) {
			results = append(results, user)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    results,
		Code:    http.StatusOK,
	})
}

func findUserByID(id int) (*User, int) {
	for i, user := range users {
		if user.ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

func validateUser(user User) error {
	if user.Name == "" {
		return fmt.Errorf("Name is required")
	}
	if user.Email == "" {
		return fmt.Errorf("Email is required")
	}
	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("Invalid email format")
	}
	return nil
}