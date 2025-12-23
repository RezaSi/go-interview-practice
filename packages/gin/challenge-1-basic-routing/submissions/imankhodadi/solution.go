
package main

import (
	"errors"
	"strconv"
	"strings"

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

func main() {
	router := gin.Default()
	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.GET("/users/search", searchUsers)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.Run(":8080")
}
func getAllUsers(c *gin.Context) {
	c.JSON(200, Response{
		Success: true,
		Data:    users,
		Message: "Users retrieved successfully",
		Code:    200})
	return
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
	} else {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		})
		return
	}
}

func findUserByID(id int) (*User, int) {
	for ind, user := range users {
		if user.ID == id {
			return &user, ind
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
	nextID++
	newUser.ID = nextID
	if err := validateUser(newUser); err != nil {
		c.JSON(400, Response{Success: false, Error: err.Error()})
		return
	}
	users = append(users, newUser)
	c.JSON(201, Response{Success: true, Data: newUser, Message: "User created", Code: 201})
	return
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
		user.ID = newUser.ID
		user.Name = newUser.Name
		user.Age = newUser.Age
		user.Email = newUser.Email
		c.JSON(200, Response{
			Success: true,
			Data:    user,
			Message: "Users updated successfully",
			Code:    200})
	} else {
		c.JSON(404, Response{
			Success: false,
			Error:   "User not found"})
	}
	return
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
		users[ind] = users[len(users)-1]
		users = users[:len(users)-1]
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

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	query_name := c.DefaultQuery("name", "")
	if query_name == "" {
		c.JSON(400, Response{
			Success: false,
			Code:    400,
			Error:   "provide name in url"})
		return
	}
	matched_users := []User{}
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), query_name) {
			matched_users = append(matched_users, user)
		}
	}
	c.JSON(200, Response{
		Success: true,
		Data:    matched_users,
		Message: "matched",
		Code:    200})
	return
}
