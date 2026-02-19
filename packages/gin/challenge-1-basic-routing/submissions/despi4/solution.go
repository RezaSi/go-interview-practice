package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"regexp"
	"strings"
	"slices"
	"net/http"
	"errors"
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

	// TODO: Setup routes
	// GET /users/search - Search users by name
	router.GET("/users/search", searchUsers)
	// GET /users - Get all users
	router.GET("/users", getAllUsers)
	// GET /users/:id - Get user by ID
	router.GET("/users/:id", getUserByID)
	// POST /users - Create new user
	router.POST("/users", createUser)
	// PUT /users/:id - Update user
	router.PUT("/users/:id", updateUser)
	// DELETE /users/:id - Delete user
	router.DELETE("/users/:id", deleteUser)

	// TODO: Start server on port 8080
	router.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	res := Response{
	    Success: true,
	    Data: users,
	}
	
	c.JSON(http.StatusOK, res)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	
	var res Response
	
	string_id := c.Param("id")
    id, err := strconv.Atoi(string_id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
        return
    }
    
    u, index := findUserByID(id)
    if u == nil && index == -1 {
        res := Response {
            Success: false,
            Message: "Not Found",
            Error: "Not Found",
            Code: http.StatusNotFound,
        }
        
        c.JSON(http.StatusNotFound, res)
        return
    }
    
    res = Response {
        Success: true,
        Data: u,
        Message: "User Found",
    }
    
    c.JSON(http.StatusOK, res)
    return
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	
	var user User
	if err := c.BindJSON(&user); err != nil {
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
    
    res := Response {
	    Success: true, 
	    Data: user,
	    Message: "New User Created",
    }
    
    c.JSON(201, res)
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	
	id_param := c.Param("id")
    
    id, err := strconv.Atoi(id_param)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})    
        return
    }
    
    var updatedata User
    if err = c.ShouldBindJSON(&updatedata); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})    
        return
    }
    
    err = validateUser(updatedata)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    u, index := findUserByID(id)
    if u == nil && index == -1 {
        res := Response {
            Success: false,
            Message: "Not Found",
            Error: "Not Found",
            Code: http.StatusNotFound,
        }
        
        c.JSON(http.StatusNotFound, res)
        return
    }
    
    users[index].Name = updatedata.Name
    users[index].Age = updatedata.Age
    users[index].Email = updatedata.Email
    
    res := Response {
        Success: true,
        Data: users[index],
        Message: "User Updated by id: " + strconv.Itoa(id),
    }
    
    c.JSON(http.StatusOK, res)
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	
	u, index := findUserByID(id)
	if u == nil && index == -1 {
	    res := Response {
            Success: false,
            Message: "Not Found",
            Error: "Not Found",
            Code: http.StatusNotFound,
        }
        
        c.JSON(http.StatusNotFound, res)
        return
	}
	
	users = slices.Delete(users, index, index+1)
	
	res := Response {
	    Success: true,
	    Data: u,
	    Message: "User Deleted from Storage",
	}
	
	c.JSON(http.StatusOK, res)
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")
    
    // 1. Инициализируем как пустой массив, а не nil
    // Это гарантирует JSON ответ [] вместо null
    matchUsers := make([]User, 0)
    
    // 2. Если тесты падают на NoResults, убираем жесткую проверку на 400 ошибку
    // Или проверяем, не требует ли тест вернуть всех пользователей при пустом имени
    if name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
        return
    }
    
    queryName := strings.ToLower(name)
    
    for _, user := range users {
        userName := strings.ToLower(user.Name)
        
        // 3. Используем Contains вместо == (если тесты ищут часть имени)
        // Если тесты требуют строгого соответствия, оставь ==
        if strings.Contains(userName, queryName) {
            matchUsers = append(matchUsers, user)
        }
    }
	
	res := Response {
	    Success: true,
	    Data: matchUsers,
	    Message: "Matched users",
	}
	
    c.JSON(http.StatusOK, res)
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	
	for index, user := range users {
	    if user.ID == id {
	        return &user, index
	    }
	}
	
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	
	var (
	    name_length = len(user.Name)
	    email_length = len(user.Email)
	)
	
	if name_length == 0 || name_length > 50  {
	    if email_length == 0 || email_length > 80 {
	        return errors.New("invalid name, email")
	    }
	    
	    return errors.New("name is invalid")
	} else if email_length == 0 || email_length > 80 {
	    return errors.New("email is invalid")
	}
	
	email_pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	if !email_pattern.MatchString(user.Email) {
	    return errors.New("invalid email format")
	}
	
	return nil
}

