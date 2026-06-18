package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthAPI(t *testing.T) {
	// Create a new Echo instance for testing
	e := echo.New()
	setupMiddleware(e)
	setupRoutes(e)

	t.Run("POST /auth/register - should register new user", func(t *testing.T) {
		registerJSON := `{"username":"testuser","email":"test@example.com","password":"testpassword123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response AuthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "User registered successfully", response.Message)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotNil(t, response.User)
		assert.Equal(t, "testuser", response.User.Username)
		assert.Equal(t, "test@example.com", response.User.Email)
		assert.Equal(t, "user", response.User.Role)
		assert.Empty(t, response.User.Password) // Password should not be in response
	})

	t.Run("POST /auth/register - should validate required fields", func(t *testing.T) {
		registerJSON := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "Username")
	})

	t.Run("POST /auth/register - should reject duplicate username", func(t *testing.T) {
		// First registration
		registerJSON := `{"username":"duplicate","email":"first@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Second registration with same username
		registerJSON = `{"username":"duplicate","email":"second@example.com","password":"password123"}`
		req = httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "already exists")
	})

	t.Run("POST /auth/login - should login with valid credentials", func(t *testing.T) {
		// First register a user
		registerJSON := `{"username":"loginuser","email":"login@example.com","password":"loginpassword123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Now login
		loginJSON := `{"username":"loginuser","password":"loginpassword123"}`
		req = httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response AuthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Login successful", response.Message)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotNil(t, response.User)
		assert.Equal(t, "loginuser", response.User.Username)
	})

	t.Run("POST /auth/login - should reject invalid credentials", func(t *testing.T) {
		loginJSON := `{"username":"nonexistent","password":"wrongpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "Invalid username or password")
	})

	t.Run("GET /users/profile - should return user profile with valid token", func(t *testing.T) {
		// First register and get token
		registerJSON := `{"username":"profileuser","email":"profile@example.com","password":"profilepass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var authResponse AuthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authResponse)
		require.NoError(t, err)
		token := authResponse.AccessToken

		// Now get profile
		req = httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response UserResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.NotNil(t, response.User)
		assert.Equal(t, "profileuser", response.User.Username)
		assert.Equal(t, "profile@example.com", response.User.Email)
	})

	t.Run("GET /users/profile - should reject request without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "Authorization header required")
	})

	t.Run("PUT /users/profile - should update user profile", func(t *testing.T) {
		// First register and get token
		registerJSON := `{"username":"updateuser","email":"update@example.com","password":"updatepass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var authResponse AuthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authResponse)
		require.NoError(t, err)
		token := authResponse.AccessToken

		// Now update profile
		updateJSON := `{"email":"newemail@example.com"}`
		req = httptest.NewRequest(http.MethodPut, "/users/profile", strings.NewReader(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response UserResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Profile updated successfully", response.Message)
		assert.Equal(t, "newemail@example.com", response.User.Email)
	})

	t.Run("POST /auth/refresh - should refresh access token", func(t *testing.T) {
		// First register and get tokens
		registerJSON := `{"username":"refreshuser","email":"refresh@example.com","password":"refreshpass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var authResponse AuthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authResponse)
		require.NoError(t, err)
		refreshToken := authResponse.RefreshToken

		// Now refresh token
		refreshJSON := `{"refresh_token":"` + refreshToken + `"}`
		req = httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(refreshJSON))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response AuthResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Token refreshed successfully", response.Message)
		assert.NotEmpty(t, response.AccessToken)
	})

	t.Run("POST /users/logout - should logout user", func(t *testing.T) {
		// First register and get token
		registerJSON := `{"username":"logoutuser","email":"logout@example.com","password":"logoutpass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var authResponse AuthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authResponse)
		require.NoError(t, err)
		token := authResponse.AccessToken

		// Now logout
		req = httptest.NewRequest(http.MethodPost, "/users/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response MessageResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Logged out successfully", response.Message)

		// Verify token is blacklisted
		req = httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAdminAPI(t *testing.T) {
	// Create a new Echo instance for testing
	e := echo.New()
	setupMiddleware(e)
	setupRoutes(e)

	// Helper function to create admin user and get token
	createAdminUser := func() string {
		registerJSON := `{"username":"admin","email":"admin@example.com","password":"adminpass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Manually set admin role (in real app, this would be done differently)
		for i, user := range users {
			if user.Username == "admin" {
				users[i].Role = "admin"
				break
			}
		}

		// Login again to get a token with the admin role
		loginJSON := `{"username":"admin","password":"adminpass123"}`
		req = httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var authResponse AuthResponse
		json.Unmarshal(rec.Body.Bytes(), &authResponse)

		return authResponse.AccessToken
	}

	t.Run("GET /users - should return all users for admin", func(t *testing.T) {
		token := createAdminUser()

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response UserListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.GreaterOrEqual(t, response.Count, 1)
		assert.NotEmpty(t, response.Users)
	})

	t.Run("PUT /users/:id/role - should update user role for admin", func(t *testing.T) {
		// Reset users to avoid conflicts
		users = []User{}
		blacklistedTokens = make(map[string]bool)

		token := createAdminUser()

		// Create a regular user first
		registerJSON := `{"username":"regularuser","email":"regular@example.com","password":"regularpass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var authResponse AuthResponse
		json.Unmarshal(rec.Body.Bytes(), &authResponse)
		userID := authResponse.User.ID

		// Update role
		roleJSON := `{"role":"admin"}`
		req = httptest.NewRequest(http.MethodPut, "/users/"+userID+"/role", strings.NewReader(roleJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response UserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "User role updated successfully", response.Message)
		assert.Equal(t, "admin", response.User.Role)
	})

	t.Run("DELETE /users/:id - should delete user for admin", func(t *testing.T) {
		// Reset users to avoid conflicts
		users = []User{}
		blacklistedTokens = make(map[string]bool)

		token := createAdminUser()

		// Create a user to delete
		registerJSON := `{"username":"deleteuser","email":"delete@example.com","password":"deletepass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var authResponse AuthResponse
		json.Unmarshal(rec.Body.Bytes(), &authResponse)
		userID := authResponse.User.ID

		// Delete user
		req = httptest.NewRequest(http.MethodDelete, "/users/"+userID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response MessageResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "User deleted successfully", response.Message)
	})
}

func TestUserStructure(t *testing.T) {
	t.Run("User struct should have correct fields", func(t *testing.T) {
		user := User{
			ID:        "test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Role:      "user",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		}

		assert.Equal(t, "test-id", user.ID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "user", user.Role)
	})

	t.Run("Response structs should have correct structure", func(t *testing.T) {
		user := User{ID: "1", Username: "Test"}

		authResponse := AuthResponse{
			Status:      "success",
			Message:     "Login successful",
			AccessToken: "token",
			User:        &user,
		}

		assert.Equal(t, "success", authResponse.Status)
		assert.Equal(t, "Login successful", authResponse.Message)
		assert.Equal(t, "token", authResponse.AccessToken)
		assert.NotNil(t, authResponse.User)
	})
}
