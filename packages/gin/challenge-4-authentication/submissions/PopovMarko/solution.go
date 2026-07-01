package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID             int        `json:"id"`
	Username       string     `json:"username" binding:"required,min=3,max=30"`
	Email          string     `json:"email" binding:"required,email"`
	Password       string     `json:"-"` // Never return in JSON
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
var usersMU sync.RWMutex
var blacklistedTokens = make(map[string]bool) // Token blacklist for logout

var refreshTokens = make(map[string]int) // RefreshToken -> UserID mapping
var tokenMU sync.RWMutex
var nextUserID = 1

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

func isStrongPassword(password string) bool {
	if utf8.RuneCountInString(password) < 8 {
		return false
	}
	var (
		hasUpper  bool
		hasLower  bool
		hasDigit  bool
		hasSymbol bool
	)
	runePassword := []rune(password)

	for _, r := range runePassword {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r):
			hasSymbol = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSymbol
}

func hashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(pass), nil
}

func verifyPassword(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}

func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	claimsForAccess := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsForAccess)
	accessTokenSring, err := accessToken.SignedString(jwtSecret)

	claimsForRefresh := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsForRefresh)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)

	if err != nil {
		return nil, err
	}
	response := TokenResponse{
		AccessToken:  accessTokenSring,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    time.Now().Add(accessTokenTTL),
	}

	tokenMU.Lock()
	refreshTokens[refreshTokenString] = userID
	tokenMU.Unlock()

	return &response, nil
}

func validateToken(tokenString string) (*JWTClaims, error) {
	claims := JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected jwt method")
		}

		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	tokenMU.RLock()
	if _, exists := blacklistedTokens[tokenString]; exists {
		tokenMU.RUnlock()
		return nil, fmt.Errorf("Token compromized")
	}
	tokenMU.RUnlock()

	return &claims, nil
}

func findUserByUsername(username string) *User {
	if username == "" {
		return nil
	}
	usersMU.RLock()
	defer usersMU.RUnlock()

	for i, user := range users {
		if user.Username == username {

			return &users[i]
		}
	}

	return nil
}

func findUserByEmail(email string) *User {
	if email == "" {
		return nil
	}
	usersMU.RLock()
	defer usersMU.RUnlock()

	for i, user := range users {
		if user.Email == email {

			return &users[i]
		}
	}

	return nil
}

func findUserByID(id int) *User {
	usersMU.RLock()
	defer usersMU.RUnlock()
	for i, user := range users {
		if user.ID == id {

			return &users[i]
		}
	}

	return nil
}

func isAccountLocked(user *User) bool {
	// TODO: Check if account is locked based on LockedUntil field
	usersMU.RLock()
	defer usersMU.RUnlock()
	return time.Until(*user.LockedUntil) > 0
}

func recordFailedAttempt(user *User) {
	usersMU.Lock()
	defer usersMU.Unlock()

	user.FailedAttempts++
	now := time.Now()
	if user.LockedUntil == nil {
		user.LockedUntil = &now
	}

	if user.FailedAttempts > maxFailedAttempts {
		*user.LockedUntil = time.Now().Add(lockoutDuration)
		return
	}
}

func resetFailedAttempts(user *User) {
	usersMU.Lock()
	defer usersMU.Unlock()

	user.FailedAttempts = 0
	now := time.Now()
	user.LockedUntil = &now
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Passwords do not match",
		})
		return
	}

	if !isStrongPassword(req.Password) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Password does not meet strength requirements",
		})
		return
	}

	if user := findUserByUsername(req.Username); user != nil {
		c.JSON(409, APIResponse{
			Success: false,
			Error:   "Username duplicated",
		})
		return
	}

	if user := findUserByEmail(req.Email); user != nil {
		c.JSON(200, APIResponse{
			Success: false,
			Error:   "Email already registered",
		})
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	now := time.Now()
	usersMU.Lock()
	id := nextUserID
	nextUserID++

	user := User{
		ID:            id,
		Username:      req.Username,
		Email:         req.Email,
		Password:      req.Password,
		PasswordHash:  hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          c.GetString("role"),
		IsActive:      true,
		EmailVerified: false,
		LockedUntil:   &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	users = append(users, user)
	usersMU.Unlock()
	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to genesate tokens",
		})
	}

	c.JSON(201, APIResponse{
		Success: true,
		Data:    tokens,
		Message: "User registered successfully",
	})
}

func login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid credentials format",
		})
		return
	}

	user := findUserByUsername(req.Username)
	if user == nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	if isAccountLocked(user) {
		c.JSON(423, APIResponse{
			Success: false,
			Error:   "Account is temporarily locked",
		})
		return
	}

	usersMU.RLock()
	userPasswordHash := user.PasswordHash
	usersMU.RUnlock()

	if !verifyPassword(req.Password, userPasswordHash) {
		recordFailedAttempt(user)
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	resetFailedAttempts(user)

	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to generate tokens",
		})
		return
	}

	now := time.Now()
	usersMU.Lock()
	user.LastLogin = &now

	data := map[string]interface{}{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}

	usersMU.Unlock()
	c.JSON(200, APIResponse{
		Success: true,
		Data:    data,
		Message: "Login successful",
	})
}

func logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}

	tokenS := strings.SplitN(authHeader, " ", 2)
	if tokenS[0] != "Bearer" {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Wrong token type",
		})
		return
	}
	token := tokenS[1]
	claims, err := validateToken(token)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid token",
		})
		return
	}

	tokenMU.Lock()
	blacklistedTokens[token] = true

	for r, id := range refreshTokens {
		if id == claims.UserID {
			delete(refreshTokens, r)
			break
		}
	}
	tokenMU.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Logout successful",
	})
}

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

	// TODO: Validate refresh token
	claims, err := validateToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid  refresh token",
		})
	}
	userID := claims.UserID
	user := findUserByID(userID)
	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to generate new tokens",
		})
	}

	tokenMU.Lock()
	delete(refreshTokens, req.RefreshToken)
	refreshTokens[tokens.RefreshToken] = user.ID
	tokenMU.Unlock()
	// TODO: Get user ID from refresh token store
	// TODO: Find user by ID
	// TODO: Generate new access token
	// TODO: Optionally rotate refresh token

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
		tokenS := strings.SplitN(authHeader, " ", 2)
		if tokenS[0] != "Bearer" {
			c.JSON(400, APIResponse{
				Success: false,
				Error:   "Wrong token type",
			})
			return
		}
		token := tokenS[1]
		claims, err := validateToken(token)
		if err != nil {
			c.JSON(400, APIResponse{
				Success: false,
				Error:   "Invalid token",
			})
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		// TODO: Extract token from "Bearer <token>" format
		// TODO: Validate token using validateToken function
		// TODO: Set user info in context for route handlers

		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get user role from context (set by authMiddleware)
		role := c.GetString("role")
		// TODO: Check if user role is in allowed roles
		if !slices.Contains(roles, role) {
			c.JSON(400, APIResponse{
				Success: false,
				Error:   "nsupported role",
			})
			return
		}
		// TODO: Return 403 if not authorized
		if role != RoleAdmin {
			c.JSON(403, APIResponse{
				Success: false,
				Error:   "Forbidden: insufficient permissions",
			})
			return
		}

		c.Next()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// TODO: Get user ID from context (set by authMiddleware)
	// TODO: Find user by ID
	// TODO: Return user profile (without sensitive data)

	c.JSON(200, APIResponse{
		Success: true,
		Data:    nil, // TODO: Return user data
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

	// TODO: Get user ID from context
	// TODO: Find user by ID
	// TODO: Check if new email is already taken
	// TODO: Update user profile

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

	// TODO: Get user ID from context
	// TODO: Find user by ID
	// TODO: Verify current password
	// TODO: Validate new password strength
	// TODO: Hash new password and update user

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	// TODO: Get pagination parameters
	// TODO: Return list of users (without sensitive data)

	c.JSON(200, APIResponse{
		Success: true,
		Data:    users, // TODO: Filter sensitive data
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
	_ = findUserByID(id)
	// TODO: Update user role

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
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server")
	}
}
