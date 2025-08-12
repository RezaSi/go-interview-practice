package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	// Reset users data for each test
	users = []User{
		{ID: 1, Username: "admin", Email: "admin@example.com", Password: "$2a$12$...", Role: "admin", Active: true},
		{ID: 2, Username: "user1", Email: "user1@example.com", Password: "$2a$12$...", Role: "user", Active: true},
	}
	nextUserID = 3

	setupCustomValidator()

	app := fiber.New()

	// Public routes
	app.Post("/auth/register", registerHandler)
	app.Post("/auth/login", loginHandler)
	app.Get("/health", healthHandler)

	// Protected routes (require JWT)
	protected := app.Group("/", jwtMiddleware())
	protected.Get("/profile", getProfileHandler)
	protected.Put("/profile", updateProfileHandler)
	protected.Post("/auth/refresh", refreshTokenHandler)

	// Admin routes (require admin role)
	admin := app.Group("/admin", jwtMiddleware(), adminMiddleware())
	admin.Get("/users", listUsersHandler)
	admin.Put("/users/:id/role", updateUserRoleHandler)

	return app
}

func TestHealthCheck(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestRegisterUser(t *testing.T) {
	app := setupTestApp()

	registerData := RegisterRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "StrongP@ss123",
	}

	body, _ := json.Marshal(registerData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var response AuthResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response.Success)
	assert.Equal(t, "newuser", response.User.Username)
	assert.Equal(t, "user", response.User.Role)
}

func TestRegisterValidation(t *testing.T) {
	app := setupTestApp()

	// Test weak password
	registerData := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "weak", // Too weak
	}

	body, _ := json.Marshal(registerData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestLogin(t *testing.T) {
	app := setupTestApp()

	// First register a user
	registerData := RegisterRequest{
		Username: "testlogin",
		Email:    "testlogin@example.com",
		Password: "StrongP@ss123",
	}

	body, _ := json.Marshal(registerData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req)

	// Now test login
	loginData := LoginRequest{
		Username: "testlogin",
		Password: "StrongP@ss123",
	}

	body, _ = json.Marshal(loginData)
	req = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response AuthResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response.Success)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "testlogin", response.User.Username)
}

func TestInvalidLogin(t *testing.T) {
	app := setupTestApp()

	loginData := LoginRequest{
		Username: "nonexistent",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestProtectedRoute(t *testing.T) {
	app := setupTestApp()

	// Test without token
	req := httptest.NewRequest("GET", "/profile", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	// Test with invalid token
	req = httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAdminRoute(t *testing.T) {
	app := setupTestApp()

	// Test without token
	req := httptest.NewRequest("GET", "/admin/users", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestPasswordHashing(t *testing.T) {
	password := "TestPassword123!"

	hash, err := hashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Test verification
	assert.True(t, verifyPassword(password, hash))
	assert.False(t, verifyPassword("wrongpassword", hash))
}

func TestPasswordValidation(t *testing.T) {
	// Test valid password
	assert.True(t, validatePassword("StrongP@ss123"))

	// Test invalid passwords
	assert.False(t, validatePassword("weak"))         // Too short
	assert.False(t, validatePassword("nouppercase"))  // No uppercase
	assert.False(t, validatePassword("NOLOWERCASE"))  // No lowercase
	assert.False(t, validatePassword("NoDigits!"))    // No digits
	assert.False(t, validatePassword("NoSpecial123")) // No special chars
}

func TestJWTGeneration(t *testing.T) {
	user := User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	token, err := generateJWT(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test validation
	claims, err := validateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.Role, claims.Role)
}

func TestDuplicateRegistration(t *testing.T) {
	app := setupTestApp()

	registerData := RegisterRequest{
		Username: "duplicate",
		Email:    "duplicate@example.com",
		Password: "StrongP@ss123",
	}

	body, _ := json.Marshal(registerData)

	// First registration should succeed
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	// Second registration with same username should fail
	req = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 409, resp.StatusCode) // Conflict
}
