package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter() *gin.Engine {
	// Reset global state for each test
	users = []User{}
	blacklistedTokens = make(map[string]bool)
	refreshTokens = make(map[string]int)
	nextUserID = 1

	// Add default admin user
	adminHash, _ := hashPassword("admin123")
	users = append(users, User{
		ID:            nextUserID,
		Username:      "admin",
		Email:         "admin@example.com",
		PasswordHash:  adminHash,
		FirstName:     "Admin",
		LastName:      "User",
		Role:          RoleAdmin,
		IsActive:      true,
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
	nextUserID++

	return setupRouter()
}

func TestPasswordStrength(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"password", false},   // No uppercase, number, special char
		{"Password", false},   // No number, special char
		{"Password1", false},  // No special char
		{"Password1!", true},  // Valid strong password
		{"Pass1!", false},     // Too short
		{"PASSWORD1!", false}, // No lowercase
		{"password1!", false}, // No uppercase
		{"Password!", false},  // No number
		{"Passw0rd!", true},   // Valid strong password
	}

	for _, test := range tests {
		result := isStrongPassword(test.password)
		assert.Equal(t, test.expected, result, "Password: %s", test.password)
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"

	hash, err := hashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Test verification
	assert.True(t, verifyPassword(password, hash))
	assert.False(t, verifyPassword("wrongpassword", hash))
}

func TestUserRegistration(t *testing.T) {
	router := setupTestRouter()

	t.Run("Valid Registration", func(t *testing.T) {
		regData := RegisterRequest{
			Username:        "testuser",
			Email:           "test@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
			FirstName:       "Test",
			LastName:        "User",
		}

		jsonData, _ := json.Marshal(regData)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "User registered successfully", response.Message)
	})

	t.Run("Duplicate Username", func(t *testing.T) {
		regData := RegisterRequest{
			Username:        "admin", // Already exists
			Email:           "admin2@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
			FirstName:       "Admin",
			LastName:        "Two",
		}

		jsonData, _ := json.Marshal(regData)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Password Mismatch", func(t *testing.T) {
		regData := RegisterRequest{
			Username:        "testuser2",
			Email:           "test2@example.com",
			Password:        "Password123!",
			ConfirmPassword: "DifferentPassword123!",
			FirstName:       "Test",
			LastName:        "User",
		}

		jsonData, _ := json.Marshal(regData)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Weak Password", func(t *testing.T) {
		regData := RegisterRequest{
			Username:        "testuser3",
			Email:           "test3@example.com",
			Password:        "weak",
			ConfirmPassword: "weak",
			FirstName:       "Test",
			LastName:        "User",
		}

		jsonData, _ := json.Marshal(regData)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		regData := RegisterRequest{
			Username:        "testuser4",
			Email:           "invalid-email",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
			FirstName:       "Test",
			LastName:        "User",
		}

		jsonData, _ := json.Marshal(regData)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserLogin(t *testing.T) {
	router := setupTestRouter()

	t.Run("Valid Login", func(t *testing.T) {
		loginData := LoginRequest{
			Username: "admin",
			Password: "admin123",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		// Check if response contains token data
		tokenData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, tokenData, "access_token")
		assert.Contains(t, tokenData, "refresh_token")
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		loginData := LoginRequest{
			Username: "admin",
			Password: "wrongpassword",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Nonexistent User", func(t *testing.T) {
		loginData := LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestTokenGeneration(t *testing.T) {
	tokens, err := generateTokens(1, "testuser", RoleUser)
	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Equal(t, "Bearer", tokens.TokenType)
	assert.Greater(t, tokens.ExpiresIn, int64(0))
}

func TestTokenValidation(t *testing.T) {
	tokens, err := generateTokens(1, "testuser", RoleUser)
	assert.NoError(t, err)

	claims, err := validateToken(tokens.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, 1, claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, RoleUser, claims.Role)

	// Test invalid token
	_, err = validateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestProtectedRoutes(t *testing.T) {
	router := setupTestRouter()

	// Get valid token
	tokens, _ := generateTokens(1, "admin", RoleAdmin)

	t.Run("Access Protected Route with Valid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/user/profile", nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Access Protected Route without Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/user/profile", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Access Protected Route with Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/user/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid.token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAdminRoutes(t *testing.T) {
	router := setupTestRouter()

	adminTokens, _ := generateTokens(1, "admin", RoleAdmin)
	userTokens, _ := generateTokens(2, "user", RoleUser)

	t.Run("Admin Access to Admin Route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+adminTokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("User Access to Admin Route (Forbidden)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+userTokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestTokenRefresh(t *testing.T) {
	router := setupTestRouter()

	tokens, _ := generateTokens(1, "admin", RoleAdmin)

	t.Run("Valid Token Refresh", func(t *testing.T) {
		refreshData := map[string]string{
			"refresh_token": tokens.RefreshToken,
		}

		jsonData, _ := json.Marshal(refreshData)
		req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Invalid Refresh Token", func(t *testing.T) {
		refreshData := map[string]string{
			"refresh_token": "invalid.refresh.token",
		}

		jsonData, _ := json.Marshal(refreshData)
		req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestLogout(t *testing.T) {
	router := setupTestRouter()

	tokens, _ := generateTokens(1, "admin", RoleAdmin)

	t.Run("Valid Logout", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Token should be blacklisted
		assert.True(t, blacklistedTokens[tokens.AccessToken])
	})

	t.Run("Logout without Token", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/logout", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestPasswordChange(t *testing.T) {
	router := setupTestRouter()

	tokens, _ := generateTokens(1, "admin", RoleAdmin)

	t.Run("Valid Password Change", func(t *testing.T) {
		changeData := map[string]string{
			"current_password": "admin123",
			"new_password":     "NewPassword123!",
		}

		jsonData, _ := json.Marshal(changeData)
		req, _ := http.NewRequest("POST", "/user/change-password", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Wrong Current Password", func(t *testing.T) {
		changeData := map[string]string{
			"current_password": "wrongpassword",
			"new_password":     "NewPassword123!",
		}

		jsonData, _ := json.Marshal(changeData)
		req, _ := http.NewRequest("POST", "/user/change-password", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Weak New Password", func(t *testing.T) {
		changeData := map[string]string{
			"current_password": "admin123",
			"new_password":     "weak",
		}

		jsonData, _ := json.Marshal(changeData)
		req, _ := http.NewRequest("POST", "/user/change-password", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserLookupFunctions(t *testing.T) {
	setupTestRouter() // Initialize users

	t.Run("Find User by Username", func(t *testing.T) {
		user := findUserByUsername("admin")
		assert.NotNil(t, user)
		assert.Equal(t, "admin", user.Username)

		user = findUserByUsername("nonexistent")
		assert.Nil(t, user)
	})

	t.Run("Find User by Email", func(t *testing.T) {
		user := findUserByEmail("admin@example.com")
		assert.NotNil(t, user)
		assert.Equal(t, "admin@example.com", user.Email)

		user = findUserByEmail("nonexistent@example.com")
		assert.Nil(t, user)
	})
}

func TestRoleChange(t *testing.T) {
	router := setupTestRouter()

	adminTokens, _ := generateTokens(1, "admin", RoleAdmin)

	t.Run("Admin Changes User Role", func(t *testing.T) {
		roleData := map[string]string{
			"role": RoleModerator,
		}

		jsonData, _ := json.Marshal(roleData)
		req, _ := http.NewRequest("PUT", "/admin/users/1/role", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminTokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Role", func(t *testing.T) {
		roleData := map[string]string{
			"role": "invalid_role",
		}

		jsonData, _ := json.Marshal(roleData)
		req, _ := http.NewRequest("PUT", "/admin/users/1/role", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminTokens.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRateLimiting(t *testing.T) {
	// This test would require implementing rate limiting
	// For now, we'll skip it as it's not in the basic requirements
	t.Skip("Rate limiting test - implement if rate limiting is added")
}

// Test basic endpoint functionality
func TestBasicEndpoints(t *testing.T) {
	router := setupTestRouter()

	t.Run("Registration Endpoint Exists", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should not return 404 (endpoint exists)
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})

	t.Run("Login Endpoint Exists", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should not return 404 (endpoint exists)
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})

	t.Run("Profile Endpoint Exists (Protected)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/user/profile", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 401 (unauthorized) or other auth error, not 404
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})

	t.Run("Admin Endpoint Exists (Protected)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/users", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 401 (unauthorized) or other auth error, not 404
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})
}

// Test that validation functions are implemented
func TestValidationFunctions(t *testing.T) {
	t.Run("Password Strength Function Exists", func(t *testing.T) {
		// Test that function doesn't panic and returns a boolean
		result := isStrongPassword("test")
		assert.IsType(t, false, result)
	})

	t.Run("Hash Password Function Exists", func(t *testing.T) {
		// Test that function doesn't panic
		hash, err := hashPassword("testpassword")
		// Function should return something (even if empty) and not panic
		assert.IsType(t, "", hash)
		assert.IsType(t, (*error)(nil), &err)
	})

	t.Run("Verify Password Function Exists", func(t *testing.T) {
		// Test that function doesn't panic and returns a boolean
		result := verifyPassword("test", "hash")
		assert.IsType(t, false, result)
	})
}

// Test that helper functions are implemented
func TestHelperFunctions(t *testing.T) {
	setupTestRouter() // Initialize data

	t.Run("Find User by Username Function", func(t *testing.T) {
		user := findUserByUsername("admin")
		// Function should not panic and return *User or nil
		assert.IsType(t, (*User)(nil), user)
	})

	t.Run("Find User by Email Function", func(t *testing.T) {
		user := findUserByEmail("admin@example.com")
		// Function should not panic and return *User or nil
		assert.IsType(t, (*User)(nil), user)
	})

	t.Run("Generate Tokens Function", func(t *testing.T) {
		tokens, err := generateTokens(1, "testuser", RoleUser)
		// Function should not panic
		assert.IsType(t, (*TokenResponse)(nil), tokens)
		assert.IsType(t, (*error)(nil), &err)
	})

	t.Run("Validate Token Function", func(t *testing.T) {
		claims, err := validateToken("dummy.token.string")
		// Function should not panic
		assert.IsType(t, (*JWTClaims)(nil), claims)
		assert.IsType(t, (*error)(nil), &err)
	})
}

// Test middleware functions
func TestMiddleware(t *testing.T) {
	router := setupTestRouter()

	t.Run("Auth Middleware Blocks Unauthorized Access", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/user/profile", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 401 without token
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Admin Middleware Blocks Non-Admin Access", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/users", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 401 or 403 without proper token/role
		assert.Contains(t, []int{http.StatusUnauthorized, http.StatusForbidden}, w.Code)
	})
}

// Test JSON response structure
func TestResponseStructure(t *testing.T) {
	router := setupTestRouter()

	t.Run("Registration Returns Proper JSON Structure", func(t *testing.T) {
		regData := map[string]string{
			"username": "test",
			"password": "short", // This should fail validation
		}

		jsonData, _ := json.Marshal(regData)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Response should have success field
		assert.IsType(t, false, response.Success)
	})

	t.Run("Login Returns Proper JSON Structure", func(t *testing.T) {
		loginData := map[string]string{
			"username": "nonexistent",
			"password": "wrongpassword",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Response should have success field
		assert.IsType(t, false, response.Success)
	})
}
