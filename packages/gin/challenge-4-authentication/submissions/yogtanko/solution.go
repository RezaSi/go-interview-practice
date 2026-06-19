package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
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
	ID       int    `json:"id"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	// Password       string     `json:"-"` // Never return in JSON
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
var blacklistedTokens = make(map[string]bool) // Token blacklist for logout
var refreshTokens = make(map[string]int)      // RefreshToken -> UserID mapping
var nextUserID = 1
var mu sync.RWMutex

// Configuration
var (
	jwtSecret         = []byte("your-super-secret-jwt-key")
	accessTokenTTL    = 15 * time.Minute   // 15 minutes
	refreshTokenTTL   = 7 * 24 * time.Hour // 7 days
	maxFailedAttempts = 5
	lockoutDuration   = 30 * time.Minute
)

// User roles
const (
	RoleUser      = "user"
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
)

var passwordFormat = `^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[!@#$%^&*(),.?":{}|<>]).{8,}$`

// TODO: Implement password strength validation
func isStrongPassword(password string) bool {
	// TODO: Validate password strength:
	// Check length
	if len(password) < 8 {
		return false
	}

	// Check for uppercase
	hasUpper := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return false
	}

	// Check for lowercase
	hasLower := false
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return false
	}

	// Check for digit
	hasDigit := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	// Check for special character
	specialChars := "!@#$%^&*(),.?\":{}|<>"
	hasSpecial := false
	for _, char := range password {
		for _, special := range specialChars {
			if char == special {
				hasSpecial = true
				break
			}
		}
		if hasSpecial {
			break
		}
	}
	if !hasSpecial {
		return false
	}

	return true
}

// TODO: Implement password hashing
func hashPassword(password string) (string, error) {
	// TODO: Use bcrypt to hash the password with cost 12
	hashedPassowrd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassowrd), nil
}

// TODO: Implement password verification
func verifyPassword(password, hash string) bool {
	// TODO: Use bcrypt to compare password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

// TODO: Implement JWT token generation
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	// TODO: Generate access token with 15 minute expiry
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	// TODO: Generate refresh token with 7 day expiry
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// TODO: Store refresh token in memory store
	mu.Lock()
	refreshTokens[refreshTokenString] = userID
	mu.Unlock()

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    time.Now().Add(accessTokenTTL),
	}, nil
}

// TODO: Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	// TODO: Parse and validate JWT token
	claim := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	mu.RLock()
	defer mu.RUnlock()
	// TODO: Check if token is blacklisted
	if blacklistedTokens[token.Raw] {
		return nil, errors.New("token blacklisted")
	}
	// TODO: Return claims if valid
	return claim, nil
}

// TODO: Implement user lookup functions
func findUserByUsername(username string) *User {
	// TODO: Find user by username in users slice
	mu.RLock()
	defer mu.RUnlock()
	for i := range users {
		if strings.EqualFold(users[i].Username, username) {
			return &users[i]
		}
	}
	return nil
}

func findUserByEmail(email string) *User {
	// TODO: Find user by email in users slice
	mu.RLock()
	defer mu.RUnlock()
	for i := range users {
		if strings.EqualFold(users[i].Email, email) {
			return &users[i]
		}
	}
	return nil
}

func findUserByID(id int) *User {
	// TODO: Find user by ID in users slice
	mu.RLock()
	defer mu.RUnlock()
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

// TODO: Implement account lockout check
func isAccountLocked(user *User) bool {
	// TODO: Check if account is locked based on LockedUntil field
	if user.LockedUntil == nil {
		return false
	}
	if time.Now().Compare(*user.LockedUntil) == 1 {
		return false
	}
	return true
}

// TODO: Implement failed attempt tracking
func recordFailedAttempt(user *User) {
	// TODO: Increment failed attempts counter
	user.FailedAttempts++
	// TODO: Lock account if max attempts reached
	if user.FailedAttempts == maxFailedAttempts {
		expiry := time.Now().Add(lockoutDuration)
		user.LockedUntil = &expiry
	}
}

func resetFailedAttempts(user *User) {
	// TODO: Reset failed attempts counter and unlock account
	user.FailedAttempts = 0
	user.LockedUntil = nil
}

// TODO: Generate secure random token
func generateRandomToken() (string, error) {
	// TODO: Generate cryptographically secure random token
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	// TODO: Validate password confirmation
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Passwords do not match",
		})
		return
	}

	// TODO: Validate password strength
	if !isStrongPassword(req.Password) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Password does not meet strength requirements",
		})
		return
	}
	// TODO: Check if username already exists
	if u := findUserByUsername(req.Username); u != nil {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   "Username already exists",
		})
		return
	}
	// TODO: Check if email already exists
	if u := findUserByEmail(req.Email); u != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Email already exists",
		})
		return
	}

	// TODO: Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	mu.Lock()
	// TODO: Create user and add to users slice
	users = append(users, User{
		ID:       nextUserID,
		Username: req.Username,
		Email:    req.Email,
		// Password:       req.Password,
		PasswordHash:   hashedPassword,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Role:           RoleUser,
		IsActive:       true,
		EmailVerified:  false,
		LastLogin:      nil,
		FailedAttempts: 0,
		LockedUntil:    nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	nextUserID++
	mu.Unlock()

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

// POST /auth/login - User login
func login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid credentials format",
		})
		return
	}

	// TODO: Find user by username
	user := findUserByUsername(req.Username)
	if user == nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// TODO: Check if account is locked
	if isAccountLocked(user) {
		c.JSON(http.StatusForbidden, APIResponse{
			Success: false,
			Error:   "Account is temporarily locked",
		})
		return
	}

	// TODO: Verify password
	if !verifyPassword(req.Password, user.PasswordHash) {
		recordFailedAttempt(user)
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// TODO: Reset failed attempts on successful login
	resetFailedAttempts(user)

	// TODO: Update last login time
	now := time.Now()
	user.LastLogin = &now

	// TODO: Generate tokens
	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: true,
			Error:   "Failed to generate tokens",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    tokens,
		Message: "Login successful",
	})
}

// POST /auth/logout - User logout
func logout(c *gin.Context) {
	// TODO: Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}

	// TODO: Extract token from "Bearer <token>" format
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid authorization format",
		})
		return
	}
	mu.Lock()
	// TODO: Add token to blacklist
	blacklistedTokens[parts[1]] = true
	mu.Unlock()
	// TODO: Remove refresh token from store
	claims, err := validateToken(parts[1])
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Message: "Logout successful",
		})
		return
	}

	for k, v := range refreshTokens {
		if v == claims.UserID {
			delete(refreshTokens, k)
		}
	}

	c.JSON(http.StatusOK, APIResponse{
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Refresh token required",
		})
		return
	}

	// TODO: Validate refresh token
	_, err := validateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Refresh token invalid",
		})
		return
	}
	mu.RLock()
	// TODO: Get user ID from refresh token store
	userId, exists := refreshTokens[req.RefreshToken]
	mu.RUnlock()
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}
	// TODO: Find user by ID
	user := findUserByID(userId)
	// TODO: Generate new access token
	token, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to generate new token",
		})
		return
	}
	// TODO: Optionally rotate refresh token
	delete(refreshTokens, req.RefreshToken)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    token,
		Message: "Token refreshed successfully",
	})
}

// Middleware: JWT Authentication
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Authorization header required",
			})
			c.Abort()
			return
		}

		// TODO: Extract token from "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid authorization format",
			})
			c.Abort()
			return
		}
		bearerToken := parts[1]
		// TODO: Validate token using validateToken function
		claim, err := validateToken(bearerToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid Token",
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		// TODO: Set user info in context for route handlers
		c.Set("user_id", claim.UserID)
		c.Set("user_role", claim.Role)
		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get user role from context (set by authMiddleware)
		userRole := c.GetString("user_role")
		// TODO: Check if user role is in allowed roles
		for i := range roles {
			if roles[i] == userRole {
				c.Next()
				return
			}
		}
		// TODO: Return 403 if not authorized
		c.JSON(http.StatusForbidden, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// TODO: Get user ID from context (set by authMiddleware)
	userId := c.GetInt("user_id")
	if userId == 0 {
		c.JSON(http.StatusOK, APIResponse{
			Success: false,
			Error:   "Unauthorized",
			Message: "Failed to get user ID",
		})
		return
	}
	// TODO: Find user by ID
	user := findUserByID(userId)
	// TODO: Return user profile (without sensitive data)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    user, // TODO: Return user data
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	// TODO: Get user ID from context
	userIdContext := c.GetString("user_id")
	userId, err := strconv.Atoi(userIdContext)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Unauthorized",
			Message: "Failed to get user ID",
		})
		return
	}
	// TODO: Find user by ID
	user := findUserByID(userId)
	// TODO: Check if new email is already taken
	if e := findUserByEmail(req.Email); e != nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Message: "Email already registered",
		})
		return
	}
	// TODO: Update user profile
	mu.Lock()
	defer mu.Unlock()
	user.FirstName = req.FirstName
	user.Email = req.Email
	user.LastName = req.LastName

	c.JSON(http.StatusOK, APIResponse{
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	// TODO: Get user ID from context
	userId := c.GetInt("user_id")

	// TODO: Find user by ID
	user := findUserByID(userId)
	// TODO: Verify current password
	if !verifyPassword(req.CurrentPassword, user.PasswordHash) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: true,
			Message: "Wrong password",
		})
		return
	}
	// TODO: Validate new password strength
	if !isStrongPassword(req.NewPassword) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Password does not meet strength requirements",
		})
		return
	}
	// TODO: Hash new password and update user
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	// user.Password = req.NewPassword
	user.PasswordHash = hashedPassword

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	// TODO: Get pagination parameters
	// TODO: Return list of users (without sensitive data)
	mu.RLock()
	defer mu.RUnlock()
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    users, // TODO: Filter sensitive data
		Message: "Users retrieved successfully",
	})
}

// PUT /admin/users/:id/role - Change user role (admin only)
func changeUserRole(c *gin.Context) {
	userIDContext := c.Param("id")
	// id, err := strconv.Atoi(userID)
	userID, err := strconv.Atoi(userIDContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid role",
		})
		return
	}

	// TODO: Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Bad Request",
			Message: "User not found",
		})
		return
	}
	// TODO: Update user role
	mu.Lock()
	user.Role = req.Role
	mu.Unlock()
	c.JSON(http.StatusOK, APIResponse{
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

	router := setupRouter()
	router.Run(":8080")
}
