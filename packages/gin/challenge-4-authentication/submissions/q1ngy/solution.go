package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	regexp "github.com/dlclark/regexp2"
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
var blacklistedTokens = make(map[string]bool) // Token blacklist for logout
var refreshTokens = make(map[string]int)      // RefreshToken -> UserID mapping
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

// Implement password strength validation
func isStrongPassword(password string) bool {
	// Validate password strength:
	// - At least 8 characters
	// - Contains uppercase letter
	// - Contains lowercase letter
	// - Contains number
	// - Contains special character
	pattern := `^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$`
	re := regexp.MustCompile(pattern, regexp.None)
	b, err := re.MatchString(password)
	if err != nil {
		return false
	} else {
		return b
	}
}

// Implement password hashing
func hashPassword(password string) (string, error) {
	// Use bcrypt to hash the password with cost 12
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

// Implement password verification
func verifyPassword(password, hash string) bool {
	// Use bcrypt to compare password with hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// Implement JWT token generation
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	// Generate access token with 15 minute expiry
	// Generate refresh token with 7 day expiry
	// Store refresh token in memory store
	access := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}
	refresh := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}
	tAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, access)
	tRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh)
	sAccess, _ := tAccess.SignedString(jwtSecret)
	sRefresh, _ := tRefresh.SignedString(jwtSecret)

	refreshTokens[sRefresh] = userID
	return &TokenResponse{
		AccessToken:  sAccess,
		RefreshToken: sRefresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    time.Now().Add(accessTokenTTL),
	}, nil
}

// Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	// Parse and validate JWT token
	var claims = &JWTClaims{}
	// var claims *JWTClaims 只是声明了一个空指针，而上面则声明了零值
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	if _, ok := blacklistedTokens[tokenString]; ok {
		return nil, err
	}

	// Return claims if valid
	return claims, nil
}

// Implement user lookup functions
func findUserByUsername(username string) *User {
	// Find user by username in users slice
	for i := range users {
		if username == users[i].Username {
			return &users[i]
		}
	}
	return nil
}

func findUserByEmail(email string) *User {
	// Find user by email in users slice
	for i := range users {
		if email == users[i].Email {
			return &users[i]
		}
	}
	return nil
}

func findUserByID(id int) *User {
	// Find user by ID in users slice
	for i := range users {
		if id == users[i].ID {
			return &users[i]
		}
	}
	return nil
}

// Implement account lockout check
func isAccountLocked(user *User) bool {
	// Check if account is locked based on LockedUntil field
	if user.LockedUntil != nil && time.Until(*user.LockedUntil) > 0 {
		return true
	}
	return false
}

// Implement failed attempt tracking
func recordFailedAttempt(user *User) {
	// Increment failed attempts counter
	user.FailedAttempts++
	// Lock account if max attempts reached
	if user.FailedAttempts >= maxFailedAttempts {
		lockUntil := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockUntil
	}
}

func resetFailedAttempts(user *User) {
	// Reset failed attempts counter and unlock account
	now := time.Now()
	user.LockedUntil = &now
	user.FailedAttempts = 0
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
	u := findUserByUsername(req.Username)
	if u != nil {
		c.JSON(409, APIResponse{
			Success: false,
		})
		return
	}

	// Check if email already exists
	u = findUserByEmail(req.Email)
	if u != nil {
		c.JSON(409, APIResponse{
			Success: false,
		})
		return
	}

	// Hash password
	encrypted, _ := hashPassword(req.Password)

	// Create user and add to users slice
	users = append(users, User{
		Username:     req.Username,
		Password:     req.Password,
		PasswordHash: encrypted,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	})

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
	user.LastLogin = &now

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
	segments := strings.Split(authHeader, " ")
	token := segments[1]

	// Add token to blacklist
	blacklistedTokens[token] = true

	// Remove refresh token from store
	delete(refreshTokens, token)

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

	// Validate refresh token
	claims := JWTClaims{}
	_, err := jwt.ParseWithClaims(req.RefreshToken, &claims, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Refresh token required",
		})
		return
	}

	// Get user ID from refresh token store
	userId := refreshTokens[req.RefreshToken]

	// Find user by ID
	u := findUserByID(userId)
	if u == nil {
	}

	// Generate new access token
	tokens, err := generateTokens(userId, u.Username, u.Role)

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

		// Extract token from "Bearer <token>" format
		segments := strings.Split(authHeader, " ")
		token := segments[1]

		// Validate token using validateToken function
		claims, err := validateToken(token)
		if claims == nil || err != nil {
			c.JSON(401, APIResponse{
				Success: false,
			})
			c.Abort()
			return
		}

		// Set user info in context for route handlers
		c.Set("id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by authMiddleware)
		role, _ := c.Get("role")

		// Check if user role is in allowed roles
		// Return 403 if not authorized
		if !slices.Contains(roles, role.(string)) {
			c.JSON(403, APIResponse{
				Success: false,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// Get user ID from context (set by authMiddleware)
	userId, _ := c.Get("id")

	// Find user by ID
	u := findUserByID(userId.(int))

	// Return user profile (without sensitive data
	c.JSON(200, APIResponse{
		Success: true,
		Data:    u, // Retuern user data
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
	userId, _ := c.Get("id")

	// Find user by ID
	u := findUserByID(userId.(int))
	if u == nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid credentials format",
		})
		return
	}

	// Check if new email is already taken
	userByEmail := findUserByEmail(req.Email)
	if userByEmail != nil && userId.(int) != userByEmail.ID {
		c.JSON(409, APIResponse{
			Success: false,
		})
		return
	}

	// Update user profile
	u.Email = req.Email
	u.FirstName = req.FirstName
	u.LastName = req.LastName

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

	// Get user ID from context
	userId, _ := c.Get("id")

	// Find user by ID
	u := findUserByID(userId.(int))
	if u == nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid credentials format",
		})
		return
	}

	// Verify current password
	if !verifyPassword(req.CurrentPassword, u.PasswordHash) {
		c.JSON(400, APIResponse{
			Success: false,
		})
		return
	}

	// Validate new password strength
	if !isStrongPassword(req.NewPassword) {
		c.JSON(400, APIResponse{
			Success: false,
		})
		return
	}

	// Hash new password and update user
	encrypted, _ := hashPassword(req.NewPassword)
	u.Password = req.NewPassword
	u.PasswordHash = encrypted

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	// Get pagination parameters
	// Return list of users (without sensitive data)

	c.JSON(200, APIResponse{
		Success: true,
		Data:    users, // Filter sensitive data
		Message: "Users retrieved successfully",
	})
}

// PUT /admin/users/:id/role - Change user role (admin only)
func changeUserRole(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	fmt.Println(id)
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

	// Find user by ID
	u := findUserByID(id)
	if u == nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid credentials format",
		})
		return
	}

	// Update user role
	u.Role = req.Role

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
	router.Run(":8080")
}
