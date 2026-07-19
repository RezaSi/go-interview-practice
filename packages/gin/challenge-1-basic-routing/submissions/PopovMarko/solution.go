package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Response represents a standard API response envelope returned by every endpoint.
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
var mu sync.RWMutex

// main configures the routes and starts the HTTP server on port 8080.
func main() {
	router := gin.New()

	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/search", searchUsers)

	router.Run(":8080")
}

// getAllUsers handles GET /users and returns the full list of users.
func getAllUsers(c *gin.Context) {
	if users == nil {
		errorResponse(
			c,
			"database not initialized",
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}
	mu.RLock()
	snapshot := append([]User{}, users...)
	mu.RUnlock()
	successResponse(c, snapshot, "Users retrieved successfully", http.StatusOK)
}

// getUserByID handles GET /users/:id and returns a single user by its ID.
func getUserByID(c *gin.Context) {
	id, err := getUserIdFromPath(c)
	if err != nil {
		errorResponse(
			c, "ID param should be integer", err.Error(), http.StatusBadRequest,
		)
		return
	}
	user, index := findUserByID(id)
	if index == -1 {
		errorResponse(
			c, "User not found", "not found", http.StatusNotFound,
		)
		return
	}
	successResponse(c, user, "User retrieved successfully", http.StatusOK)
}

// createUser handles POST /users, validating and persisting a new user.
func createUser(c *gin.Context) {
	user := User{}

	if err := c.ShouldBindJSON(&user); err != nil {
		errorResponse(c, "Failed to read JSON", err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateUser(user); err != nil {
		errorResponse(c, "User data not valid", err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	user.ID = nextID
	nextID++
	users = append(users, user)
	mu.Unlock()
	successResponse(c, user, "User added successfully", http.StatusCreated)
}

// updateUser handles PUT /users/:id, applying any non-empty fields from the
// request body to the existing user in place.
func updateUser(c *gin.Context) {
	id, err := getUserIdFromPath(c)
	if err != nil {
		errorResponse(c, "User ID should be integer", err.Error(), http.StatusBadRequest)
		return
	}

	userDataToUpdate := User{}
	if err := c.ShouldBindJSON(&userDataToUpdate); err != nil {
		errorResponse(c, "Failed to parse JSON", err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()

	user, _ := findUserByID(id)
	if user == nil {
		errorResponse(c, "User not found", "not found", http.StatusNotFound)
		mu.Unlock()
		return
	}
	candidateUser := *user

	if userDataToUpdate.Email != "" {
		candidateUser.Email = userDataToUpdate.Email
	}
	if userDataToUpdate.Age != 0 {
		candidateUser.Age = userDataToUpdate.Age
	}
	if userDataToUpdate.Name != "" {
		candidateUser.Name = userDataToUpdate.Name
	}
	if err := validateUser(candidateUser); err != nil {
		errorResponse(c, "Bad user data", err.Error(), http.StatusBadRequest)
		return
	}
	user = &candidateUser
	mu.Unlock()
	successResponse(c, candidateUser, "User updated successfully", http.StatusOK)
}

// deleteUser handles DELETE /users/:id and removes the matching user.
func deleteUser(c *gin.Context) {
	id, err := getUserIdFromPath(c)
	if err != nil {
		errorResponse(c, "User ID should be integer", err.Error(), http.StatusBadRequest)
		return
	}
	_, i := findUserByID(id)
	if i == -1 {
		errorResponse(c, "User not found", "not found", http.StatusNotFound)
		return
	}
	mu.Lock()
	users = append(users[:i], users[i+1:]...)
	mu.Unlock()
	successResponse(c, nil, "User deleted", http.StatusOK)
}

// searchUsers handles GET /users/search?name=value and returns all users whose
// name contains the given query, case-insensitively.
func searchUsers(c *gin.Context) {
	name, ok := c.GetQuery("name")
	if !ok || name == "" {
		errorResponse(c, "No query parameter", "bad request", http.StatusBadRequest)
		return
	}
	query := strings.ToLower(name)
	responseUser := []User{}
	mu.RLock()
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), query) {
			responseUser = append(responseUser, user)
		}
	}
	mu.RUnlock()
	successResponse(c, responseUser, "Users retrieved successfully", http.StatusOK)
}

// findUserByID returns a pointer to the stored user with the given ID and its
// index in the slice, or (nil, -1) if no such user exists.
func findUserByID(id int) (*User, int) {
	for i, user := range users {
		if user.ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

// validateUser checks that the required fields are present and the email is
// well-formed, returning an error describing the first problem found.
func validateUser(user User) error {
	if user.Name == "" {
		return errors.New("user name required")
	}
	if user.Email == "" {
		return errors.New("user email required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("user email not valid")
	}
	return nil
}

// getUserIdFromPath extracts the ":id" path parameter and parses it as an integer.
func getUserIdFromPath(c *gin.Context) (int, error) {
	p := c.Param("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		return 0, errors.New("bad request")
	}
	return id, nil
}

// saveUser appends the given user to the store and returns the persisted copy.
func saveUser(user *User) (User, error) {
	mu.Lock()
	defer mu.Unlock()
	if users == nil {
		users = []User{}
	}
	if user == nil {
		return User{}, errors.New("bad user data")
	}
	users = append(users, *user)
	newUser, _ := findUserByID(user.ID)
	if newUser == nil {
		return User{}, errors.New("failed to save user")
	}
	return *newUser, nil
}

// errorResponse writes a standardized failure response with the given message,
// error detail, and HTTP status code.
func errorResponse(c *gin.Context, message, err string, code int) {
	response := Response{
		Success: false,
		Message: message,
		Error:   err,
		Code:    code,
	}
	responseHandler(c, response)
}

// successResponse writes a standardized success response carrying the given
// data, message, and HTTP status code.
func successResponse(c *gin.Context, data interface{}, message string, code int) {
	response := Response{
		Success: true,
		Data:    data,
		Message: message,
		Code:    code,
	}
	responseHandler(c, response)
}

// responseHandler serializes the response as JSON using its embedded status code.
func responseHandler(c *gin.Context, response Response) {
	c.JSON(response.Code, response)
}
