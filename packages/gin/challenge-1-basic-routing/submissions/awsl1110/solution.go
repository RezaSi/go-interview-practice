package main

import (
    "errors"
    "github.com/gin-gonic/gin"
    "net/http" // Import the http package
    "regexp"
    "strconv" // Import the strconv package
    "strings"
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
    router := gin.Default()
    router.GET("/users", getAllUsers)
    router.GET("/users/:id", getUserByID)
    router.GET("/users/search", searchUsers)
    router.POST("/users", createUser)
    router.PUT("/users/:id", updateUser)
    router.DELETE("/users/:id", deleteUser)
    err := router.Run(":8080")
    if err != nil {
        return
    }
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
    c.JSON(http.StatusOK, Response{Success: true, Data: users})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
    intId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error(), Code: http.StatusBadRequest}) // Use Response struct
        return
    }
    userid, _ := findUserByID(intId)
    if userid == nil { // Handle user not found
        c.JSON(http.StatusNotFound, Response{Success: false, Error: "User not found", Code: http.StatusNotFound})
        return
    }
    c.JSON(http.StatusOK, Response{Success: true, Data: userid})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        // Changed from StatusBadGateway to StatusBadRequest as it's a client-side error
        c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error(), Code: http.StatusBadRequest})
        return
    }
    if err := validateUser(user); err != nil {
        c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error(), Code: http.StatusBadRequest})
        return
    }
    user.ID = nextID
    nextID++
    users = append(users, user)
    c.JSON(http.StatusCreated, // Use StatusCreated for successful creation
        Response{Success: true, Data: user, Message: "User created", Code: http.StatusCreated})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
    intId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error(), Code: http.StatusBadRequest})
        return
    }
    var newUser User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        // Changed from StatusBadGateway to StatusBadRequest as it's a client-side error
        c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error(), Code: http.StatusBadRequest})
        return // Added return here
    }
    userid, oldId := findUserByID(intId)
    if userid == nil {
        c.JSON(http.StatusNotFound, Response{Success: false, Error: "User not found", Code: http.StatusNotFound}) // Use StatusNotFound
        return
    }
    if err := validateUser(newUser); err != nil {
        c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error(), Code: http.StatusBadRequest})
        return
    }
    newUser.ID = intId
    users[oldId] = newUser
    c.JSON(http.StatusOK, Response{Success: true, Data: newUser, Message: "User updated"})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
    intId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error(), Code: http.StatusBadRequest})
        return
    }
    user, id := findUserByID(intId)
    if user == nil {
        c.JSON(http.StatusNotFound, Response{Success: false, Error: "User not found", Code: http.StatusNotFound}) // Use StatusNotFound
        return
    }
    users = append(users[:id], users[id+1:]...)
    c.JSON(http.StatusOK, Response{Success: true, Data: user, Message: "User deleted"})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
    name := c.Query("name")
    if name == "" {
        c.JSON(http.StatusBadRequest, Response{Success: false, Error: "name is empty", Code: http.StatusBadRequest})
        return
    }
    userArr := make([]User, 0)
    for _, user := range users {
        if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
            userArr = append(userArr, user)
        }
    }
    c.JSON(http.StatusOK, Response{Success: true, Data: userArr})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
    for i, user := range users {
        if user.ID == id {
            return &user, i
        }
    }
    return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
    if user.Name == "" {
        return errors.New("user name is required")
    }
    if user.Email == "" {
        return errors.New("user email is required")
    }
    const re = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    emailValidator := regexp.MustCompile(re)
    if !emailValidator.MatchString(user.Email) {
        return errors.New("invalid email")
    }
    return nil
}