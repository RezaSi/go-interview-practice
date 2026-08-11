package main

import (
	"errors"
	"fmt"
	"regexp"
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
	bcryptCost        = 12
)

// User roles
const (
	RoleUser      = "user"
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
)

// TODO: Implement password strength validation
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUppercase {
		return false
	}

	hasLowcase := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLowcase {
		return false
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return false
	}

	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return false
	}

	return true
}

// TODO: Implement password hashing
func hashPassword(password string) (string, error) {
	// TODO: Use bcrypt to hash the password with cost 12
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// TODO: Implement password verification
func verifyPassword(password, hash string) bool {
	// TODO: Use bcrypt to compare password with hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TODO: Implement JWT token generation
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	// TODO: Generate access token with 15 minute expiry
	accessTokenString, accessClaims, err := generateToken(userID, username, role, accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// TODO: Generate refresh token with 7 day expiry
	refreshTokenString, _, err := generateToken(userID, username, role, refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	// TODO: Store refresh token in memory store
	refreshTokens[refreshTokenString] = userID

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    accessClaims.ExpiresAt.Time, //time.Now().Add(accessTokenTTL),
	}, nil
}

func generateToken(userID int, username, role string, ttl time.Duration) (string, *JWTClaims, error) {
	expiresAt := time.Now().Add(ttl)

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, claims, nil
}

// TODO: Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {

	// TODO: Check if token is blacklisted
	if blacklistedTokens[tokenString] {
		return nil, errors.New("token is blacklisted")
	}

	// TODO: Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// TODO: Return claims if valid
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")

}

// TODO: Implement user lookup functions
func findUserByUsername(username string) *User {
	// TODO: Find user by username in users slice
	for _, user := range users {
		if user.Username == username {
			return &user
		}
	}
	return nil
}

func findUserByEmail(email string) *User {
	// TODO: Find user by email in users slice
	for _, user := range users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

func findUserByID(id int) *User {
	// TODO: Find user by ID in users slice
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

func createUser(req RegisterRequest, passwordHash string) User {
	user := User{
		ID:           nextUserID,
		Username:     req.FirstName,
		Email:        req.Email,
		Password:     "", //req.Password,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	nextUserID++
	return user
}

// TODO: Implement account lockout check
func isAccountLocked(user *User) bool {
	// TODO: Check if account is locked based on LockedUntil field
	if user.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*user.LockedUntil)
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

/*
// TODO: Generate secure random token
func generateRandomToken() (string, error) {
	// TODO: Generate cryptographically secure random token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}*/

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
	if user := findUserByUsername(req.Username); user != nil {
		c.JSON(409, APIResponse{
			Success: false,
			Error:   "username already exists",
		})
		return
	}
	// TODO: Check if email already exists
	if user := findUserByEmail(req.Email); user != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "email already exists",
		})
		return
	}
	// TODO: Hash password
	hashPass, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// TODO: Create user and add to users slice
	user := createUser(req, hashPass)
	users = append(users, user)

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
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	c.ShouldBindJSON(&req)
	if req.RefreshToken != "" {
		delete(refreshTokens, req.RefreshToken)
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
	_, err := validateToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// TODO: Get user ID from refresh token store
	userID, ok := refreshTokens[req.RefreshToken]
	if !ok {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}

	// TODO: Find user by ID
	user := findUserByID(userID)
	// TODO: Generate new access token
	tokenResponse, err := generateTokens(userID, user.Username, user.Role)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	// TODO: Optionally rotate refresh token

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data:    tokenResponse,
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
		claim, err := validateToken(tokenString)
		if err != nil {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		// TODO: Set user info in context for route handlers
		c.Set("claims", claim)
		c.Set("user_id", claim.UserID)
		c.Set("role", claim.Role)

		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get user role from context (set by authMiddleware)
		role, _ := c.Get("role")
		// TODO: Check if user role is in allowed roles
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		// TODO: Return 403 if not authorized
		c.JSON(403, APIResponse{
			Success: false,
			Error:   "not authorized",
		})
		c.Abort()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// TODO: Get user ID from context (set by authMiddleware)
	userIDval, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "error user ID from context",
		})
		return
	}
	userID := userIDval.(int)

	// TODO: Find user by ID
	user := findUserByID(userID)
	// TODO: Return user profile (without sensitive data)
	userProfile := User{
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

	c.JSON(200, APIResponse{
		Success: true,
		Data:    userProfile, // TODO: Return user data
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
	userIDval, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "error user ID from context",
		})
		return
	}
	userID := userIDval.(int)

	// TODO: Find user by ID
	user := findUserByID(userID)

	// TODO: Check if new email is already taken
	if userTaken := findUserByEmail(req.Email); userTaken != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "email already exists",
		})
		return
	}

	// TODO: Update user profile
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.UpdatedAt = time.Now()

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
	userIDval, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "error user ID from context",
		})
		return
	}
	userID := userIDval.(int)

	// TODO: Find user by ID
	user := findUserByID(userID)

	// TODO: Verify current password
	if !verifyPassword(req.CurrentPassword, user.PasswordHash) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "error current password",
		})
		return
	}
	// TODO: Validate new password strength
	if !isStrongPassword(req.NewPassword) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "error new password",
		})
		return
	}
	// TODO: Hash new password and update user
	passwordHash, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user.PasswordHash = passwordHash
	user.UpdatedAt = time.Now()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	// TODO: Get pagination parameters
	// TODO: Return list of users (without sensitive data)
	var usersList []User
	for _, user := range users {
		usersList = append(usersList,
			User{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
				IsActive:  user.IsActive,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    usersList, // TODO: Filter sensitive data
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
	// TODO: Update user role
	user.Role = req.Role
	user.UpdatedAt = time.Now()

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
	router.Run(":8085")
}
