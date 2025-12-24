package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
    "sync"
	"github.com/gin-gonic/gin"
)

// Email regex for validation
var (
    EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    usersMutex = sync.RWMutex{}
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

	router.GET("/users", getAllUsers)
	router.GET("/users/search", searchUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	if err := router.Run(":8080"); err != nil {

		fmt.Println("failed to start server: ", err)
	}
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
    usersMutex.RLock()
    defer usersMutex.RUnlock()
    
	c.JSON(200, Response{
		Success: true,
		Data:    users,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "user id must be int",
		})
		
		return
	}
	
	usersMutex.RLock()
    defer usersMutex.RUnlock()
    
	user, idx := findUserByID(id)
    
	if idx == -1 {
	    
		c.JSON(404, Response{
			Success: false,
			Message: "user not found",
		})

		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    user,
	})

}

// createUser handles POST /users
func createUser(c *gin.Context) {

	var bodyData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	if err := c.ShouldBindJSON(&bodyData); err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "invalid request body",
		})

		return
	}

	user := User{
		Name:  bodyData.Name,
		Email: bodyData.Email,
		Age:   bodyData.Age,
	}

	if err := validateUser(user); err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: fmt.Sprintf("invalid user data: %s", err.Error()),
		})

		return
	}
	
	
	usersMutex.Lock()
    defer usersMutex.Unlock()

	user.ID = nextID
	users   = append(users, user)
	nextID  += 1
    
	c.JSON(201, Response{
		Success: true,
		Data:    user,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "user id must be integer",
		})

		return
	}


	var bodyData struct {
	    ID    int    `json:"id,omitempty"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

   
	if  err = c.ShouldBindJSON(&bodyData); err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "invalid request body",
		})

		return
	}
	
	err = validateUser(User{
        Name: bodyData.Name,
        Email: bodyData.Email,
        Age: bodyData.Age,
    })
    
	if err != nil {
	    	c.JSON(400, Response{
			Success: false,
			Message: fmt.Sprintf("validation error: %s", err.Error()),
		})

		return
	}

    usersMutex.Lock()
    defer usersMutex.Unlock()
    
	_, idx := findUserByID(id)

	if idx == -1 {
	    
		c.JSON(404, Response{
			Success: false,
			Message: "user not found",
		})

		return
	}
    
	users[idx].Name = bodyData.Name
	users[idx].Email = bodyData.Email
	users[idx].Age = bodyData.Age
    
    bodyData.ID = users[idx].ID
    
	c.JSON(200, Response{
		Success: true,
		Data:    bodyData,
	})

}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
    
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, Response{
			Success: false,
			Message: "user id must be integer",
		})

		return
	}
	
	usersMutex.Lock()
	defer usersMutex.Unlock()
	
	_, idx := findUserByID(id)

    if idx == -1 {
        
    	c.JSON(404, Response{
    		Success: false,
    		Message: "user not found",
    	})
    	
    	return
    }
    
	users = append(users[:idx], users[idx+1:]...)
	
	c.JSON(200, Response{
		Success: true,
		Data:    map[string]User{},
	})
    
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
    
	userName := c.Query("name")
	
	if userName == "" {
		c.JSON(400, Response{
			Success: false,
			Message: "user name must be provided",
		})

		return
	}

    result := []User{}
    
    usersMutex.RLock()
    defer usersMutex.RUnlock()
    
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(userName)) {
			
			result = append(result, user)
		}
	}
	
	if len(result) > 0 {
	    c.JSON(200, Response{
				Success: true,
				Data:    result,
			})
			
	} else {
	    c.JSON(200, Response{
		Success: true,
		Data:    result,
		Message: "user not found",
	})

	}


}

// // Helper function to find user by ID
func findUserByID(id int) (User, int) {

	for idx, user := range users {
		if user.ID == id {
			return users[idx], idx

		}
	}
	return User{}, -1
}

// Helper function to validate user data
func validateUser(user User) error {

	if len(user.Name) == 0 || len(user.Email) == 0 {
		return fmt.Errorf("user name and email must be provided, len(name)= %d, len(email) = %d", len(user.Name), len(user.Email))
	}

	if !EmailRX.MatchString(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
