package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
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

// APIResponse represents standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Global data stores (in a real app, these would be databases)
var users = []User{}
var nextUserID = 1
var userLock sync.RWMutex

var blacklistedTokens = make(map[string]bool) // Token blacklist for logout
var blacklistedTokenLock sync.RWMutex

var refreshTokens = make(map[string]int)                // RefreshToken -> UserID mapping
var refreshTokensExpiresAt = make(map[string]time.Time) // RefreshToken -> ExpiresAt
var refreshTokenLock sync.RWMutex

// Configuration
var (
	jwtSecret         = []byte("your-super-secret-jwt-key") // fixed for assignment, move it to env in production
	accessTokenTTL    = 15 * time.Minute                    // 15 minutes
	refreshTokenTTL   = 7 * 24 * time.Hour                  // 7 days
	maxFailedAttempts = 5
	lockoutDuration   = 30 * time.Minute
)

// User roles
const (
	RoleUser      = "user"
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
)

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12) // cost = 12
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWT token generation
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	// Generate access token with 15 minute expiry
	// Generate refresh token with 7 day expiry
	// Store refresh token in memory store

	// Access Token
	accessClaims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	// Refresh Token
	refreshToken, err := generateRandomToken()
	if err != nil {
		return nil, err
	}
	// Store refresh token

	refreshTokenLock.Lock()
	refreshTokens[refreshToken] = userID
	refreshTokensExpiresAt[refreshToken] = time.Now().Add(refreshTokenTTL)
	refreshTokenLock.Unlock()

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    time.Now().Add(accessTokenTTL),
	}, nil
}

// JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	// Parse and validate JWT token
	// Check if token is blacklisted
	// Return claims if valid
	// Use jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})) or reject non-HMAC/HS256 in the key function before returning jwtSecret.
	blacklistedTokenLock.RLock()
	defer blacklistedTokenLock.RUnlock()
	if blacklistedTokens[tokenString] {
		return nil, errors.New("token is blacklisted")
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}

			return jwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func findUserByUsername(username string) *User {
	userLock.RLock()
	defer userLock.RUnlock()
	for i := range users {
		if users[i].Username == username {
			return &users[i]
		}
	}
	return nil
}

func findUserByEmail(email string) *User {
	userLock.RLock()
	defer userLock.RUnlock()
	for i := range users {
		if users[i].Email == email {
			return &users[i]
		}
	}
	return nil
}

func findUserByID(id int) *User {
	userLock.RLock()
	defer userLock.RUnlock()
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

func isAccountLocked(user *User) bool {
	if user.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*user.LockedUntil)
}

func recordFailedAttempt(user *User) {
	userLock.Lock()
	defer userLock.Unlock()
	user.FailedAttempts++
	if user.FailedAttempts >= maxFailedAttempts {
		lockUntil := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockUntil
	}
}

func resetFailedAttempts(user *User) {
	userLock.Lock()
	defer userLock.Unlock()
	user.FailedAttempts = 0
	user.LockedUntil = nil
}

func generateRandomToken() (string, error) {
	// Generate cryptographically secure random token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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
	userLock.Lock()
	defer userLock.Unlock()
	for _, user := range users {
		if user.Username == req.Username {
			c.JSON(409, APIResponse{
				Success: false,
				Error:   "Account already exists",
			})
			return
		}
		if user.Email == req.Email {
			c.JSON(409, APIResponse{
				Success: false,
				Error:   "Account already exists",
			})
			return
		}
	}

	hashedPass, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Internal error",
		})
		return
	}
	createdTime := time.Now()
	newUser := User{
		ID:             nextUserID,
		Username:       req.Username,
		Email:          req.Email,
		PasswordHash:   hashedPass,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Role:           RoleUser,
		IsActive:       true,
		EmailVerified:  false,
		LastLogin:      nil,
		FailedAttempts: 0,
		LockedUntil:    nil,
		CreatedAt:      createdTime,
		UpdatedAt:      createdTime,
	}
	nextUserID++
	users = append(users, newUser)

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

	if !user.IsActive {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Not active",
		})
		return
	}
	// Reset failed attempts on successful login
	resetFailedAttempts(user)

	// Update last login time
	now := time.Now()
	userLock.Lock()
	user.LastLogin = &now
	userLock.Unlock()

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
	//  Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}
	// Extract token from "Bearer <token>" format
	// Add token to blacklist
	// Remove refresh token from store
	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid authorization header",
		})
		return
	}

	tokenString := parts[1]
	// Add token to blacklist
	blacklistedTokenLock.Lock()
	blacklistedTokens[tokenString] = true
	blacklistedTokenLock.Unlock()
	// Remove refresh token if provided
	var req struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	c.ShouldBindJSON(&req)
	if req.RefreshToken != "" {
		refreshTokenLock.Lock()
		delete(refreshTokens, req.RefreshToken)
		delete(refreshTokensExpiresAt, req.RefreshToken)
		refreshTokenLock.Unlock()

	}
	c.JSON(200, APIResponse{
		Success: true,
		Message: "Logout successful",
	})

}

func isRefreshTokenExpired(token string) bool {
	refreshTokenLock.Lock()
	defer refreshTokenLock.Unlock()

	expiry, ok := refreshTokensExpiresAt[token]
	if !ok {
		return true
	}

	if time.Now().After(expiry) {
		delete(refreshTokens, token)
		delete(refreshTokensExpiresAt, token)
		return true
	}

	return false
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

	if isRefreshTokenExpired(req.RefreshToken) {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Refresh token expired",
		})
		return
	}
	refreshTokenLock.Lock()
	userId, exists := refreshTokens[req.RefreshToken]
	refreshTokenLock.Unlock()
	if !exists {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}
	user := findUserByID(userId)
	if user == nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}
	tokens, err := generateTokens(userId, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Internal Error",
		})
		return
	}
	refreshTokenLock.Lock()
	delete(refreshTokens, req.RefreshToken)
	delete(refreshTokensExpiresAt, req.RefreshToken)
	refreshTokenLock.Unlock()
	c.JSON(200, APIResponse{
		Success: true,
		Data:    tokens,
		Message: "Token refreshed successfully",
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
		// Validate token using validateToken function
		// Set user info in context for route handlers
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   "Invalid authorization header",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   "Invalid token",
			})
			c.Abort()
			return
		}
		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by authMiddleware)
		// Check if user role is in allowed roles
		// Return 403 if not authorized
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   "Unauthorized",
			})
			c.Abort()
			return
		}
		roleStr := userRole.(string)
		for _, allowedRole := range roles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}
		c.JSON(403, APIResponse{
			Success: false,
			Error:   "Insufficient permissions",
		})
		c.Abort()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	userID := userIDVal.(int)

	user := findUserByID(userID)
	if user == nil {
		c.JSON(404, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    user,
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
	// Find user by ID
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid user",
		})
		return
	}

	userIdInt, ok := userId.(int)
	if !ok {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid user",
		})
		return
	}
	userLock.RLock()
	for i := range users {
		if users[i].Email == req.Email &&
			users[i].ID != userIdInt {
			userLock.RUnlock()
			c.JSON(409, APIResponse{
				Success: false,
				Error:   "Email already exists",
			})
			return
		}
	}
	userLock.RUnlock()
	user := findUserByID(userIdInt)

	if user == nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "user not found",
		})
		return
	}

	userLock.Lock()
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.UpdatedAt = time.Now()
	userLock.Unlock()

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

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	userID := userIDVal.(int)
	user := findUserByID(userID)

	if user == nil {
		c.JSON(404, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	if !verifyPassword(req.CurrentPassword, user.PasswordHash) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Current password is incorrect",
		})
		return
	}

	if !isStrongPassword(req.NewPassword) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Password does not meet strength requirements",
		})
		return
	}

	hash, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Internal error",
		})
		return
	}

	userLock.Lock()
	user.PasswordHash = hash
	user.UpdatedAt = time.Now()
	userLock.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	userLock.RLock()
	defer userLock.RUnlock()

	usersCopy := make([]User, len(users))
	copy(usersCopy, users)
	c.JSON(200, APIResponse{
		Success: true,
		Data:    usersCopy,
		Message: "Users retrieved successfully",
	})
}

// PUT /admin/users/:id/role - Change user role (admin only)
func changeUserRole(c *gin.Context) {
	userID := c.Param("id")
	_, err := strconv.Atoi(userID)
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

	// TODO: Validate role value
	validRoles := []string{RoleUser, RoleAdmin, RoleModerator}
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

	// TODO: Find user by ID
	// TODO: Update user role
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	user := findUserByID(id)
	if user == nil {
		c.JSON(404, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}
	userLock.Lock()
	user.Role = req.Role
	userLock.Unlock()
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

		auth.POST("/refresh", refreshToken)
	}

	// Protected user routes
	user := router.Group("/user")
	auth.POST("/logout", authMiddleware(), logout)
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
	adminHash, _ := hashPassword("Admin123!")
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
