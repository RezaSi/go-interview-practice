package main

import (
    "slices"
    "strconv"
    "net/http"
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

type Task struct {
    ID int `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
    Completed bool `json:"completed"`
}

var nextID = 4

func main() {
    router := gin.New()
    
    router.GET("/ping", pingResponse)

	// TODO: Setup routes
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
	// GET /users/search - Search users by name
    router.GET("/users/search", searchUsers)
    
    router.Run("localhost:8080")
}

// TODO: Implement handler functions
func pingResponse(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, gin.H{"message": "pong"})
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
    success[[]User](c, users)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	param_id := c.Param("id")
	id, err := strconv.Atoi(param_id)
	if err != nil {
	    badRequest(c)
	    return
	}
	
	for _, user := range users {
	    if user.ID == id {
	        success[User](c, user)
	        return
	    }
	}
	notFound(c)
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var user User
	
	if err := c.ShouldBindJSON(&user); err != nil || user.Name == "" || user.Age == 0 || user.Email == "" {
	   badRequest(c)
	   return
	}
	user.ID = len(users) + 1
	users = append(users, user)
	c.JSON(http.StatusCreated, Response{
	    Success: true,
	    Code: http.StatusCreated,
	    Data: user,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	param_id := c.Param("id")
	id, err := strconv.Atoi(param_id)
	if err != nil {
	    badRequest(c)
	    return
	}
	var user User
	if err = c.ShouldBindJSON(&user); err != nil {
	    badRequest(c)
	    return
	}
	for _, u := range users {
	    if u.ID == id {
	        if user.Name != "" {
	            u.Name = user.Name
	        }
	        if user.Email != "" {
	            u.Email = user.Email
	        }
	        if user.Age != 0 {
	            u.Age = user.Age
	        }
	        success[User](c, u)
	        return
	    }
	}
	notFound(c)
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	param_id := c.Param("id")
	id, err := strconv.Atoi(param_id)
	if err != nil {
	    badRequest(c)
	    return
	}
	len_users := len(users)
	users = slices.DeleteFunc(users, func (user User) bool {
	    return user.ID == id
	})
	if len_users == len(users) {
	    notFound(c)
	    return
	}
	c.JSON(http.StatusOK, Response {
	    Success: true,
	    Code: http.StatusOK,
	    Message: "Successfully deleted User",
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
	    badRequest(c)
	    return
	}
	var temp = []User{}
	for _, user := range users {
	    if strings.Contains( strings.ToLower(user.Name), strings.ToLower(name)) == true {
	        temp = append(temp, user)
	    }
	}
    success[[]User](c, temp)
	
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	return nil
}

func badRequest(c *gin.Context) {
    c.JSON(http.StatusBadRequest, Response{
        Success: false,
        Code: http.StatusBadRequest,
        Error: "Invalid Data",
    })
}

func notFound(c *gin.Context) {
    c.JSON(http.StatusNotFound, Response{
        Success: false,
        Code: http.StatusNotFound,
        Error: "User Not Found",
    })
}

func success[T User | []User](c *gin.Context, body T) {
    c.JSON(http.StatusOK, Response {
        Success: true,
        Code: http.StatusOK,
        Data: body,
    })
}
