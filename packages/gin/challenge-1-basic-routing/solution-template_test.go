package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Reset users data for each test
	users = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
	}
	nextID = 4

	router := gin.New()

	// Setup routes
	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/search", searchUsers)

	return router
}

func TestGetAllUsers(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check if users data is returned
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 3, len(data))
}

func TestGetUserByID_Success(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check user data
	userData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "John Doe", userData["name"])
}

func TestGetUserByID_NotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestGetUserByID_InvalidID(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestCreateUser_Success(t *testing.T) {
	router := setupRouter()

	newUser := User{
		Name:  "Alice Johnson",
		Email: "alice@example.com",
		Age:   28,
	}

	jsonData, _ := json.Marshal(newUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check created user data
	userData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Alice Johnson", userData["name"])
	assert.NotZero(t, userData["id"])
}

func TestCreateUser_InvalidData(t *testing.T) {
	router := setupRouter()

	// Missing required fields
	invalidUser := User{
		Age: 28,
	}

	jsonData, _ := json.Marshal(invalidUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestUpdateUser_Success(t *testing.T) {
	router := setupRouter()

	updatedUser := User{
		Name:  "John Updated",
		Email: "john.updated@example.com",
		Age:   31,
	}

	jsonData, _ := json.Marshal(updatedUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check updated user data
	userData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "John Updated", userData["name"])
}

func TestUpdateUser_NotFound(t *testing.T) {
	router := setupRouter()

	updatedUser := User{
		Name:  "Updated Name",
		Email: "updated@example.com",
		Age:   25,
	}

	jsonData, _ := json.Marshal(updatedUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestDeleteUser_Success(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify user is actually deleted
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 404, w2.Code)
}

func TestDeleteUser_NotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestSearchUsers_Success(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/search?name=john", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check search results
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(data)) // Should find John Doe
}

func TestSearchUsers_NoResults(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/search?name=nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check empty results
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 0, len(data))
}

func TestSearchUsers_MissingParameter(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/search", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}
