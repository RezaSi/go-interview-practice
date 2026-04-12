package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
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

// Password validation regex
var (
	rePasswordUpper   = regexp.MustCompile(`[A-Z]`)
	rePasswordLower   = regexp.MustCompile(`[a-z]`)
	rePasswordNumber  = regexp.MustCompile(`[0-9]`)
	rePasswordSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

// TODO: Implement password strength validation
func isStrongPassword(password string) bool {
	// TODO: Validate password strength:
	// - At least 8 characters
	// - Contains uppercase letter
	// - Contains lowercase letter
	// - Contains number
	// - Contains special character
	if len(password) < 8 ||
		!rePasswordUpper.MatchString(password) ||
		!rePasswordLower.MatchString(password) ||
		!rePasswordNumber.MatchString(password) ||
		!rePasswordSpecial.MatchString(password) {
		return false
	}

	return true
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
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

// TODO: Implement JWT token generation
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	now := time.Now()

	// TODO: Generate access token with 15 minute expiry
	accessTokenClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   strconv.Itoa(userID),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// TODO: Generate refresh token with 7 day expiry
	refreshTokenClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   strconv.Itoa(userID),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// TODO: Store refresh token in memory store
	refreshTokens[refreshToken] = userID

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    now.Add(accessTokenTTL),
	}, nil
}

// TODO: Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	// TODO: Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// TODO: Check if token is blacklisted
	if blacklistedTokens[tokenString] {
		return nil, fmt.Errorf("token is blacklisted")
	}

	// TODO: Return claims if valid
	return claims, nil
}

// TODO: Implement user lookup functions
func findUserByUsername(username string) *User {
	// TODO: Find user by username in users slice
	uIndex := slices.IndexFunc(users, func(u User) bool {
		return u.Username == username
	})

	if uIndex == -1 {
		return nil
	}

	return &users[uIndex]
}

func findUserByEmail(email string) *User {
	// TODO: Find user by email in users slice
	uIndex := slices.IndexFunc(users, func(u User) bool {
		return u.Email == email
	})

	if uIndex == -1 {
		return nil
	}

	return &users[uIndex]
}

func findUserByID(id int) *User {
	// TODO: Find user by ID in users slice
	uIndex := slices.IndexFunc(users, func(u User) bool {
		return u.ID == id
	})

	if uIndex == -1 {
		return nil
	}

	return &users[uIndex]
}

// TODO: Implement account lockout check
func isAccountLocked(user *User) bool {
	// TODO: Check if account is locked based on LockedUntil field
	if user.LockedUntil != nil && user.LockedUntil.Before(time.Now()) {
		return true
	}

	return false
}

// TODO: Implement failed attempt tracking
func recordFailedAttempt(user *User) {
	// TODO: Increment failed attempts counter
	user.FailedAttempts++

	// TODO: Lock account if max attempts reached
	if user.FailedAttempts >= maxFailedAttempts {
		lockedUntil := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockedUntil
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
	if slices.ContainsFunc(users, func(u User) bool {
		return u.Username == req.Username
	}) {
		c.JSON(409, APIResponse{
			Success: false,
			Error:   "Username already exists",
		})
		return
	}

	// TODO: Check if email already exists
	if slices.ContainsFunc(users, func(u User) bool {
		return u.Email == req.Email
	}) {
		c.JSON(409, APIResponse{
			Success: false,
			Error:   "Email already exists",
		})
		return
	}

	// TODO: Hash password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to hash password",
		})
		return
	}

	// TODO: Create user and add to users slice
	user := User{
		ID:           nextUserID,
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
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
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// TODO: Add token to blacklist
	blacklistedTokens[token] = true

	// TODO: Remove refresh token from store
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

	// TODO: Validate refresh token
	_, err := validateToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}

	// TODO: Get user ID from refresh token store
	userID := refreshTokens[req.RefreshToken]

	// TODO: Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}

	// TODO: Generate new access token
	// TODO: Optionally rotate refresh token
	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to generate tokens",
		})
		return
	}

	delete(refreshTokens, req.RefreshToken)
	blacklistedTokens[req.RefreshToken] = true
	refreshTokens[tokens.RefreshToken] = user.ID

	c.Header("Authorization", "Bearer "+tokens.AccessToken)
	c.Header("X-Refresh-Token", tokens.RefreshToken)

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
	})
}

// Middleware: JWT Authentication
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, APIResponse{
				Success: false,
				Error:   "Authorization header required",
			})
			return
		}

		// TODO: Extract token from "Bearer <token>" format
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(401, APIResponse{
				Success: false,
				Error:   "Authorization header required",
			})
			return
		}

		// TODO: Validate token using validateToken function
		claims, err := validateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, APIResponse{
				Success: false,
				Error:   "Invalid token",
			})
			return
		}

		// TODO: Set user info in context for route handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get user role from context (set by authMiddleware)
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(401, APIResponse{
				Success: false,
				Error:   "Unauthorized",
			})
			return
		}

		// TODO: Check if user role is in allowed roles
		// TODO: Return 403 if not authorized
		if !slices.Contains(roles, role) {
			c.AbortWithStatusJSON(403, APIResponse{
				Success: false,
				Error:   "Unauthorized",
			})
			return
		}

		c.Next()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// TODO: Get user ID from context (set by authMiddleware)
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.AbortWithStatusJSON(401, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// TODO: Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.AbortWithStatusJSON(401, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// TODO: Return user profile (without sensitive data)
	userProfile := User{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    userProfile,
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
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.AbortWithStatusJSON(401, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// TODO: Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.AbortWithStatusJSON(401, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// TODO: Check if new email is already taken
	if slices.ContainsFunc(users, func(u User) bool {
		return u.Email == req.Email
	}) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Email already taken",
		})
		return
	}

	// TODO: Update user profile
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email

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
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.AbortWithStatusJSON(401, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// TODO: Find user by ID
	user := findUserByID(userID)
	if user == nil {
		c.AbortWithStatusJSON(401, APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// TODO: Verify current password
	if !verifyPassword(req.CurrentPassword, user.PasswordHash) {
		c.AbortWithStatusJSON(400, APIResponse{
			Success: false,
			Error:   "Invalid current password",
		})
		return
	}

	// TODO: Validate new password strength
	if !isStrongPassword(req.NewPassword) {
		c.AbortWithStatusJSON(400, APIResponse{
			Success: false,
			Error:   "New password is not strong enough",
		})
		return
	}

	// TODO: Hash new password and update user
	newPasswordHash, err := hashPassword(req.NewPassword)
	if err != nil {
		c.AbortWithStatusJSON(500, APIResponse{
			Success: false,
			Error:   "Failed to hash password",
		})
		return
	}

	user.PasswordHash = newPasswordHash

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	// TODO: Get pagination parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.AbortWithStatusJSON(400, APIResponse{
			Success: false,
			Error:   "Invalid page number",
		})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		c.AbortWithStatusJSON(400, APIResponse{
			Success: false,
			Error:   "Limit must be between 1 and 100",
		})
		return
	}

	// Calcular índices para paginação
	totalUsers := len(users)
	start := (pageInt - 1) * limitInt
	if start > totalUsers {
		start = totalUsers
	}

	end := start + limitInt
	if end > totalUsers {
		end = totalUsers
	}

	// Obter sub-fatia de usuários
	// Nota: O struct User já possui `json:"-"` para campos sensíveis
	paginatedUsers := users[start:end]

	c.JSON(200, APIResponse{
		Success: true,
		Data:    paginatedUsers,
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
	user := findUserByID(id)
	if user == nil {
		c.AbortWithStatusJSON(404, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// TODO: Update user role
	user.Role = req.Role

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
