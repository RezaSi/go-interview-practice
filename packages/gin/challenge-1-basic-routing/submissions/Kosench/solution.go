package main

import (
	"errors"
	"log"
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

	router.GET("/users", getAllUsers)
	router.GET("/users/search", searchUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users", updateUser)
	router.DELETE("/users/:id", deleteUser)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	usersCopy := make([]User, len(users))
	copy(usersCopy, users)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    usersCopy,
		Message: "Users retrieved successfully",
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format",
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
		Message: "User retrieved successfully",
	})

}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
	}

	if err := validateUser(user); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	user.ID = nextID
	nextID++

	users = append(users, user)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
		Message: "User created successfully",
	})
	return
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := validateUser(updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
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

	updatedUser.ID = user.ID
	*user = updatedUser

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "User updated successfully",
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	_, idx := findUserByID(id)
	if idx == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Удаляем элемент из слайса, сохраняя порядок остальных
	users = append(users[:idx], users[idx+1:]...)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Query parameter 'name' is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Case-insensitive поиск
	nameLower := strings.ToLower(name)
	var results []User
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), nameLower) {
			results = append(results, u)
		}
	}

	// Возвращаем пустой слайс, а не nil, чтобы JSON был [], а не null
	if results == nil {
		results = []User{}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    results,
		Message: "Search completed successfully",
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for i, u := range users {
		if u.ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

func handleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	})
}
