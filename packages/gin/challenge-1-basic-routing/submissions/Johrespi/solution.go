package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"errors"
	"strings"
	"net/http"
	"strconv"
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
	r := Response{
	    Success: true,
	    Data: users,
	    Message: "Success",
	    Error: "",
	    Code: http.StatusOK,
	}
	
	c.JSON(http.StatusOK, r)
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
	    return
	}
	
	foundUser, _ := findUserByID(userID)

	if foundUser != nil {
	    c.JSON(http.StatusOK, Response{
            Success: true,
            Data: *foundUser,
            Message: "User found",
            Error: "",
            Code: http.StatusOK,
	    })
        return 
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
    return
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	    return
	}
	
	err := validateUser(user)
	if err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	
	user.ID = nextID
	nextID++
	
	users = append(users, user)
	
	c.JSON(http.StatusCreated, Response{
        Success: true,
        Data: user,
        Message: "User created",
        Error: "",
        Code: http.StatusCreated,
	})
    return
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
	    return
	}
	
	foundUser, _ := findUserByID(userID)
    if foundUser == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    
    var user User 
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) 
        return
    }
    
    if err := validateUser(user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    foundUser.Name = user.Name
    foundUser.Email = user.Email
    foundUser.Age = user.Age
    
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data: *foundUser,
        Message: "User updated",
        Error: "",
        Code: http.StatusOK,
    })
    return
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
	}
	
	_, index := findUserByID(userID)
	
	if index == -1 {
	    c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
	}
	
	users = append(users[:index], users[index+1:]...)

	c.JSON(http.StatusOK, Response{Success: true})
}

func searchUsers(c *gin.Context) {
	queryName := c.Query("name")

	if queryName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search parameter 'name'"})
		return
	}

	searchLower := strings.ToLower(queryName)

	var matchingUsers []User

	for _, user := range users {
		userNameLower := strings.ToLower(user.Name)
		if strings.Contains(userNameLower, searchLower) {
			matchingUsers = append(matchingUsers, user)
		}
	}

	if matchingUsers == nil {
		matchingUsers = []User{}
	}

	response := Response{
		Success: true,
		Data: matchingUsers,
		Message: fmt.Sprintf("Found %d users matching '%s'", len(matchingUsers), queryName),
		Code: http.StatusOK,
	}
	c.JSON(http.StatusOK, response)
}

func findUserByID(id int) (*User, int) {
    for i := range users { 
        if users[i].ID == id {
            return &users[i], i 
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
