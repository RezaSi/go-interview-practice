package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID             int        `json:"id"`
	Username       string     `json:"username" binding:"required,min=3,max=30"`
	Email          string     `json:"email" binding:"required,email"`
	PasswordHash   string     `json:"-"`
	FirstName      string     `json:"first_name" binding:"required,min=2,max=50"`
	LastName       string     `json:"last_name" binding:"required,min=2,max=50"`
	Role           string     `json:"role"`
	IsActive       bool       `json:"is_active"`
	EmailVerified  bool       `json:"email_verified"`
	LastLogin      *time.Time `json:"last_login"`
	FailedAttempts int        `json:"-"`
	LockedUntil    *time.Time `json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=30"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	FirstName       string `json:"first_name" binding:"required,min=2,max=50"`
	LastName        string `json:"last_name" binding:"required,min=2,max=50"`
}

// TokenResponse represents JWT token response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshClaims represents JWT claims for refresh tokens
type RefreshClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// APIResponse represents standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Global data stores (in a real app, these would be databases)
var (
	usersMu           sync.RWMutex
	users             = []User{}
	nextUserID        = 1
	blacklistMutex    sync.RWMutex
	blacklistedTokens = make(map[string]bool) // Token blacklist for logout
	refreshTokenMu    sync.RWMutex
	refreshTokens     = make(map[string]int) // RefreshToken -> UserID mapping
)

// Configuration
var (
	jwtSecret         = []byte("your-super-secret-jwt-key")
	accessTokenTTL    = 15 * time.Minute   // 15 minutes
	refreshTokenTTL   = 7 * 24 * time.Hour // 7 days
	maxFailedAttempts = 5
	lockoutDuration   = 30 * time.Minute
	validRoles        = []string{RoleUser, RoleAdmin, RoleModerator}
)

// User roles
const (
	RoleUser      = "user"
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
)

// Context Keys
const (
	ClaimsKey = "claims"
	UserIDKey = "user_id"
	RoleKey   = "role"
)

// returns true if the password has an uppercase letter, lowercase letter, special character, number, and is at least 8 chars
func isStrongPassword(password string) bool {
	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Check for uppercase letter
	hasUpper := false
	for _, c := range password {
		if c >= 'A' && c <= 'Z' {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return false
	}

	// Check for lowercase letter
	hasLower := false
	for _, c := range password {
		if c >= 'a' && c <= 'z' {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return false
	}

	// Check for digit
	hasDigit := false
	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	// Check for special character
	hasSpecial := false
	for _, c := range password {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			hasSpecial = true
			break
		}
	}

	return hasSpecial
}

// hashPassword returns a hash of the password with bcrypt cost 12
func hashPassword(password string) (string, error) {
	// Use bcrypt to hash the password with cost 12
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// verifyPassword compares a password with a bcrypt hash
func verifyPassword(password, hash string) bool {
	// Use bcrypt to compare password with hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generate tokens creates a token response for a user
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	// Generate access token with 15 minute expiry
	now := time.Now()
	accessTokenClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// Generate refresh token with 7 day expiry
	refreshTokenClaims := RefreshClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// Store refresh token in memory store
	refreshTokenMu.Lock()
	refreshTokens[refreshTokenString] = userID
	refreshTokenMu.Unlock()

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    now.Add(accessTokenTTL),
	}, nil
}

// Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	// Check if token is blacklisted
	blacklistMutex.RLock()
	blocked := blacklistedTokens[tokenString]
	blacklistMutex.RUnlock()
	if blocked {
		return nil, errors.New("blocked jwt token")
	}

	// Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// Find user by username in users slice
func findUserByUsername(username string) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()
	for _, u := range users {
		if strings.EqualFold(u.Username, username) {
			userCopy := u
			return &userCopy
		}
	}
	return nil
}

// Find user by email in users slice
func findUserByEmail(email string) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()
	for _, u := range users {
		if strings.EqualFold(u.Email, email) {
			userCopy := u
			return &userCopy
		}
	}
	return nil
}

// Find user by ID in users slice
func findUserByID(id int) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()
	for _, u := range users {
		if u.ID == id {
			userCopy := u
			return &userCopy
		}
	}
	return nil
}

// isAccountLocked return true if the user is locked based on LockedUntil field
func isAccountLocked(user *User) bool {
	if user.LockedUntil == nil {
		return false
	}
	return user.LockedUntil.After(time.Now())
}

// recordFailedAttempt increments failed attempts, locks account if max attempts reached
func recordFailedAttempt(user *User) {
	usersMu.Lock()
	defer usersMu.Unlock()
	// Increment failed attempts counter
	user.FailedAttempts++
	// Lock account if max attempts reached
	if user.FailedAttempts >= maxFailedAttempts {
		lockedUntil := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockedUntil
	}
	putUser(*user)
}

// Reset failed attempts counter and unlock account
func resetFailedAttempts(user *User) {
	usersMu.Lock()
	defer usersMu.Unlock()
	user.FailedAttempts = 0
	user.LockedUntil = nil
	putUser(*user)
}

// POST /auth/register - User registration
func register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	// Validate password confirmation
	if req.Password != req.ConfirmPassword {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Passwords do not match",
		})
		return
	}

	// Validate password strength
	if !isStrongPassword(req.Password) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Password does not meet strength requirements",
		})
		return
	}

	// Check if username already exists
	if findUserByUsername(req.Username) != nil {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   "Username already exists",
		})
		return
	}

	// Check if email already exists
	if findUserByEmail(req.Email) != nil {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   "Email already exists",
		})
		return
	}

	// Hash password
	hash, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	// Create user and add to users slice
	now := time.Now()
	usersMu.Lock()
	inputUser := User{
		ID:           nextUserID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         RoleUser,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	nextUserID++
	users = append(users, inputUser)
	usersMu.Unlock()

	c.JSON(201, APIResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

// POST /auth/login - User login
func login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid credentials format",
		})
		return
	}

	// Find user by username
	user := findUserByUsername(req.Username)
	if user == nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// Check if account is locked
	if isAccountLocked(user) {
		c.JSON(423, APIResponse{
			Success: false,
			Error:   "Account is temporarily locked",
		})
		return
	}

	// Verify password
	if !verifyPassword(req.Password, user.PasswordHash) {
		recordFailedAttempt(user)
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// Reset failed attempts on successful login
	resetFailedAttempts(user)

	// Update last login time
	now := time.Now()
	usersMu.Lock()
	user.LastLogin = &now
	putUser(*user)
	usersMu.Unlock()

	// Generate tokens
	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to generate tokens",
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    tokens,
		Message: "Login successful",
	})
}

// POST /auth/logout - User logout
func logout(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}

	// Extract token from "Bearer <token>" format
	token, err := extractBearerToken(authHeader)
	if err != nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if _, err := validateToken(token); err != nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "invalid token",
		})
		return
	}

	// Add token to blacklist
	blacklistMutex.Lock()
	blacklistedTokens[token] = true
	blacklistMutex.Unlock()

	// Remove refresh token from store
	var req struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	c.ShouldBindJSON(&req)
	if req.RefreshToken != "" {
		refreshTokenMu.Lock()
		delete(refreshTokens, req.RefreshToken)
		refreshTokenMu.Unlock()
	}

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// POST /auth/refresh - Refresh access token
func refreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Refresh token required",
		})
		return
	}

	// Get user ID from refresh token store
	refreshTokenMu.RLock()
	userID, ok := refreshTokens[req.RefreshToken]
	refreshTokenMu.RUnlock()

	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}

	// Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Generate new access token
	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	// Rotate refresh token
	refreshTokenMu.Lock()
	delete(refreshTokens, req.RefreshToken)
	refreshTokens[tokens.RefreshToken] = user.ID
	refreshTokenMu.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data:    tokens,
	})
}

// Middleware: JWT Authentication
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		token, err := extractBearerToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		// Validate token using validateToken function
		claims, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "invalid token",
			})
			c.Abort()
			return
		}

		// Set user info in context for route handlers
		c.Set(ClaimsKey, claims)
		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)

		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
			Success: false,
			Message: "Forbidden",
		})
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// Get user ID from context
	userID := c.GetInt(UserIDKey)

	// Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "user not found",
		})
		return
	}

	// Return user profile (without sensitive data)
	safeUserCopy := safeUser(user)

	c.JSON(200, APIResponse{
		Success: true,
		Data:    safeUserCopy,
		Message: "Profile retrieved successfully",
	})
}

// PUT /user/profile - Update user profile
func updateUserProfile(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" binding:"required,min=2,max=50"`
		LastName  string `json:"last_name" binding:"required,min=2,max=50"`
		Email     string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	// Get user ID from context
	userID := c.GetInt(UserIDKey)

	// Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "user not found",
		})
		return
	}

	// Check if new email is already taken
	userWithEmail := findUserByEmail(req.Email)
	if userWithEmail != nil && userWithEmail.ID != user.ID {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   "email taken",
		})
		return
	}

	// Update user profile
	usersMu.Lock()
	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.UpdatedAt = time.Now()
	putUser(*user)
	usersMu.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Profile updated successfully",
	})
}

// POST /user/change-password - Change user password
func changePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	// Validate new password strength
	if !isStrongPassword(req.NewPassword) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "new password is not strong",
		})
		return
	}

	hash, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	// Get user ID from context
	userID := c.GetInt(UserIDKey)
	// Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "user not found",
		})
		return
	}

	// Verify current password
	if !verifyPassword(req.CurrentPassword, user.PasswordHash) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "invalid password",
		})
		return
	}

	// Update user
	usersMu.Lock()
	user.PasswordHash = hash
	user.UpdatedAt = time.Now()
	putUser(*user)
	usersMu.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	usersMu.RLock()
	defer usersMu.RUnlock()
	// Return list of users (without sensitive data)
	var results []User
	for _, u := range users {
		safeUser := safeUser(&u)
		results = append(results, *safeUser)
	}
	c.JSON(200, APIResponse{
		Success: true,
		Data:    results,
		Message: "Users retrieved successfully",
	})
}

// PUT /admin/users/:id/role - Change user role (admin only)
func changeUserRole(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid role data",
		})
		return
	}

	// Validate role value
	isValid := false
	for _, role := range validRoles {
		if req.Role == role {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid role",
		})
		return
	}

	// Find user by ID
	user := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Update user role
	usersMu.Lock()
	user.Role = req.Role
	user.UpdatedAt = time.Now()
	putUser(*user)
	usersMu.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "User role updated successfully",
	})
}

// Setup router with authentication routes
func setupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/logout", logout)
		auth.POST("/refresh", refreshToken)
	}

	// Protected user routes
	user := router.Group("/user")
	user.Use(authMiddleware())
	{
		user.GET("/profile", getUserProfile)
		user.PUT("/profile", updateUserProfile)
		user.POST("/change-password", changePassword)
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(authMiddleware())
	admin.Use(requireRole(RoleAdmin))
	{
		admin.GET("/users", listUsers)
		admin.PUT("/users/:id/role", changeUserRole)
	}

	return router
}

func main() {
	// Initialize with a default admin user
	adminHash, _ := hashPassword("Admin1234!")
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

	router := setupRouter()
	router.Run(":8080")
}

func extractBearerToken(token string) (string, error) {
	if !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("invalid bearer token")
	}
	tokenStr := strings.TrimPrefix(token, "Bearer ")
	return tokenStr, nil
}

func safeUser(user *User) *User {
	return &User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func putUser(user User) {
	for i, u := range users {
		if u.ID == user.ID {
			users[i] = user
		}
	}
}
