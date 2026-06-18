package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

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

// TODO: Implement password strength validation
func isStrongPassword(password string) bool {
	// TODO: Validate password strength:
	// - At least 8 characters
	if len(password) < 8 {
		return false
	}
	// - Contains uppercase letter
	// - Contains lowercase letter
	// - Contains number
	// - Contains special character
	upper, lower, number, special := false, false, false, false
	for _, c := range password {
		switch {
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsNumber(c):
			number = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}

		if upper && lower && number && special {
			return true
		}
	}
	return false
}

// TODO: Implement password hashing
func hashPassword(password string) (string, error) {
	// TODO: Use bcrypt to hash the password with cost 12
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// TODO: Implement password verification
func verifyPassword(password, hash string) bool {
	// TODO: Use bcrypt to compare password with hash
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// TODO: Implement JWT token generation
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	// TODO: Generate access token with 15 minute expiry
	now := time.Now()
	accessExpiry := now.Add(accessTokenTTL)
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   strconv.Itoa(userID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateRandomToken()
	if err != nil {
		return nil, err
	}
	refreshTokens[refreshToken] = userID
	// TODO: Generate refresh token with 7 day expiry
	// TODO: Store refresh token in memory store

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    accessExpiry,
	}, nil
}

// TODO: Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	// TODO: Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	// TODO: Check if token is blacklisted
	// TODO: Return claims if valid
	return nil, jwt.ErrSignatureInvalid
}

// TODO: Implement user lookup functions
func findUserByUsername(username string) *User {
	// TODO: Find user by username in users slice
	for i := range users {
		if users[i].Username == username {
			return &users[i]
		}
	}
	return nil
}

func findUserByEmail(email string) *User {
	// TODO: Find user by email in users slice
	for i := range users {
		if users[i].Email == email {
			return &users[i]
		}
	}
	return nil
}

func findUserByID(id int) *User {
	// TODO: Find user by ID in users slice
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
	if user.LockedUntil != nil {
		return user.LockedUntil.After(time.Now())
	}
	return false
}

// TODO: Implement failed attempt tracking
func recordFailedAttempt(user *User) {
	// TODO: Increment failed attempts counter
	user.FailedAttempts++
	// TODO: Lock account if max attempts reached
	if user.FailedAttempts >= maxFailedAttempts {
		lockUntil := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockUntil
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
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	// TODO: Validate password confirmation
	if req.Password != req.ConfirmPassword {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Passwords do not match",
		})
		return
	}

	// TODO: Validate password strength
	if !isStrongPassword(req.Password) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Password does not meet strength requirements",
		})
		return
	}

	// TODO: Check if username already exists
	if findUserByUsername(req.Username) != nil {
		c.JSON(409, APIResponse{
			Success: false,
			Error:   "Username is already taken",
		})
		return
	}
	// TODO: Check if email already exists
	if findUserByEmail(req.Email) != nil {
		c.JSON(409, APIResponse{
			Success: false,
			Error:   "Email is already taken",
		})
		return
	}
	// TODO: Hash password
	hashed, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Error hashing password",
		})
		return
	}
	// TODO: Create user and add to users slice
	user := User{
		ID:            nextUserID,
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  hashed,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          RoleUser,
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	users = append(users, user)
	nextUserID++

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

	// TODO: Find user by username
	user := findUserByUsername(req.Username)
	if user == nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// TODO: Check if account is locked
	if isAccountLocked(user) {
		c.JSON(423, APIResponse{
			Success: false,
			Error:   "Account is temporarily locked",
		})
		return
	}

	// TODO: Verify password
	if !verifyPassword(req.Password, user.PasswordHash) {
		recordFailedAttempt(user)
		c.JSON(401, APIResponse{
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
	// TODO: Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}

	// TODO: Extract token from "Bearer <token>" format
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// TODO: Add token to blacklist
	blacklistedTokens[tokenString] = true
	// TODO: Remove refresh token from store
	var req struct {
		RefreshToken string `json:"refresh-token,omitempty"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return
	}
	if req.RefreshToken != "" {
		delete(blacklistedTokens, req.RefreshToken)
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

	// TODO: Validate refresh token
	userID, ok := refreshTokens[req.RefreshToken]
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Refresh token not found",
		})
		return
	}
	if blacklistedTokens[req.RefreshToken] {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Refresh token is blacklisted",
		})
		return
	}
	// TODO: Get user ID from refresh token store
	user := findUserByID(userID)
	if user == nil || !user.IsActive {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Refresh token is not active",
		})
		return
	}
	// TODO: Find user by ID
	// TODO: Generate new access token
	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to generate tokens",
		})
		return
	}
	// TODO: Optionally rotate refresh token
	blacklistedTokens[req.RefreshToken] = true
	delete(refreshTokens, req.RefreshToken)

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

		// TODO: Extract token from "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// TODO: Validate token using validateToken function
		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   "Invalid token",
			})
			c.Abort()
			return
		}
		// TODO: Set user info in context for route handlers
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get user role from context (set by authMiddleware)
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   "Role not found in context",
			})
			c.Abort()
			return
		}
		// TODO: Check if user role is in allowed roles
		roleStr := userRole.(string)
		for _, existing := range roles {
			if roleStr == existing {
				c.Next()
				return
			}
		}
		// TODO: Return 403 if not authorized
		c.JSON(403, APIResponse{
			Success: false,
			Error:   "Insufficient permissions",
		})
		c.Abort()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// TODO: Get user ID from context (set by authMiddleware)
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found in context",
		})
		return
	}
	id, ok := userIDVal.(int)
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found in context",
		})
		return
	}
	// TODO: Find user by ID
	user := findUserByID(id)
	if user == nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found",
		})
		return
	}
	// TODO: Return user profile (without sensitive data)
	c.JSON(200, APIResponse{
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
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid input data",
		})
		return
	}

	// TODO: Get user ID from context
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found in context",
		})
		return
	}
	id, ok := userIDVal.(int)
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found in context",
		})
		return
	}
	// TODO: Find user by ID
	idx := -1
	for i, user := range users {
		if user.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found",
		})
		return
	}
	// TODO: Check if new email is already taken
	for i := range users {
		if users[i].Email == req.Email && users[i].ID != id {
			c.JSON(400, APIResponse{
				Success: false,
				Error:   "User email already in use",
			})
			return
		}
	}
	// TODO: Update user profile
	users[idx].Email = req.Email
	users[idx].LastName = req.LastName
	users[idx].FirstName = req.FirstName
	users[idx].UpdatedAt = time.Now()

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
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found in context",
		})
		return
	}
	id, ok := userIDVal.(int)
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found in context",
		})
		return
	}
	// TODO: Find user by ID
	idx := -1
	for i, user := range users {
		if user.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "User ID not found",
		})
		return
	}
	// TODO: Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(users[idx].PasswordHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Password is incorrect",
		})
		return
	}
	// TODO: Validate new password strength
	if !isStrongPassword(req.NewPassword) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "New password is not strong enough",
		})
		return
	}
	// TODO: Hash new password and update user
	newHash, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Failed to hash password",
		})
		return
	}
	users[idx].PasswordHash = newHash
	users[idx].UpdatedAt = time.Now()
	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	// TODO: Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(sizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	total := len(users)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	items := make([]User, end-start)
	for i := start; i < end; i++ {
		items = append(items, users[i])
	}
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
	user := *findUserByID(id)
	// TODO: Update user role
	req.Role = user.Role

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
