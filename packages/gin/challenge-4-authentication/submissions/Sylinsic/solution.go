package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
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

var (
	upperCaseRegexp = regexp.MustCompile(`[A-Z]+`)
	lowerCaseRegexp = regexp.MustCompile(`[a-z]+`)
	digitRegexp = regexp.MustCompile(`[0-9]+`)
	specialRegexp = regexp.MustCompile(`[!"£$%^&*()\-_=+\[\]{};:'@#~,.<>/?\\\|]+`)
)
func isStrongPassword(password string) bool {

	switch {
	case password == "":
		return false
	case len(password) < 8:
		return false
	}
	if password == "" {
		return false
	}

	if upperCaseRegexp.Match([]byte(password)) &&
	   lowerCaseRegexp.Match([]byte(password)) &&
	   digitRegexp.Match([]byte(password)) &&
	   specialRegexp.Match([]byte(password)) {
		return true
	}
	return false
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func verifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateAccessToken(userID int, username string, role string) (*jwt.Token) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{
			"iss": "sylinsic",
			"aud": "sylinsic",
			"sub": username,
			"user_id": userID,
			"username": username,
			"role": role,
			"exp": time.Now().UTC().Add(accessTokenTTL).Unix(),
			"iat": time.Now().UTC().Unix(),
		})
}

func generateRefreshToken() (*jwt.Token) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{
			"iss": "sylinsic",
			"aud": "sylinsic",
			"exp": time.Now().UTC().Add(refreshTokenTTL).Unix(),
			"iat": time.Now().UTC().Unix(),
		})
}

// TODO: Implement JWT token generation
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	at := generateAccessToken(userID, username, role)
	accessToken,err := at.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	rt := generateRefreshToken()
	refreshToken,err := rt.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	
	refreshTokens[refreshToken] = userID

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    time.Now().Add(accessTokenTTL),
	}, nil
}

// TODO: Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	var claims JWTClaims
	
	parser := jwt.NewParser(jwt.WithAudience("sylinsic"), 
							jwt.WithExpirationRequired(),
							jwt.WithIssuedAt(),
							jwt.WithIssuer("sylinsic"),
							jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	token,err := parser.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return jwtSecret,nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	if _,exists := blacklistedTokens[tokenString]; exists {
		return nil, errors.New("Token has been blacklisted")
	}

	return &claims, nil
}

// TODO: Implement user lookup functions
func findUserByUsername(username string) *User {
	index := slices.IndexFunc(users, func(u User) bool {
		return strings.EqualFold(u.Username, username)
	})

	if index == -1 {
		return nil
	}

	return &users[index]
}

func findUserByEmail(email string) *User {
	index := slices.IndexFunc(users, func(u User) bool {
		return strings.EqualFold(u.Email, email)
	})

	if index == -1 {
		return nil
	}

	return &users[index]
}

func findUserByID(id int) *User {
	index := slices.IndexFunc(users, func(u User) bool {
		return u.ID == id
	})

	if index == -1 {
		return nil
	}

	return &users[index]
}

func isAccountLocked(user *User) bool {
	if user.LockedUntil == nil {
		return false
	}
	return !user.LockedUntil.Before(time.Now())
}

func recordFailedAttempt(user *User) {
	user.FailedAttempts += 1
	if user.FailedAttempts >= maxFailedAttempts {
		lockoutTime := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockoutTime
	}
}

func resetFailedAttempts(user *User) {
	user.FailedAttempts = 0
	now := time.Now()
	user.LockedUntil = &now
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

	if u := findUserByEmail(req.Email); u != nil {
		c.JSON(409, APIResponse{
			Success: false,
			Error:   "User already exists with that email",
		})
		return
	}
	if u := findUserByUsername(req.Username); u != nil {
		c.JSON(409, APIResponse{
			Success: false,
			Error:   "Username taken",
		})
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Error hashing password",
		})
		return
	}

	user := User{
		ID:           nextUserID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Role:         RoleUser,
	}
	nextUserID += 1
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

	user := findUserByUsername(req.Username)
	if user == nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	if !verifyPassword(req.Password, user.PasswordHash) {
		recordFailedAttempt(user)
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

	resetFailedAttempts(user)

	now := time.Now()
	user.LastLogin = &now

	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(401, APIResponse{
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
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	blacklistedTokens[token] = true

	// TODO: Remove refresh token from store

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

	parser := jwt.NewParser(jwt.WithAudience("sylinsic"),
							jwt.WithIssuer("sylinsic"),
							jwt.WithIssuedAt(),
							jwt.WithExpirationRequired())

	token, err := parser.Parse(req.RefreshToken, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}

	userId := refreshTokens[req.RefreshToken]

	user := findUserByID(userId)
	tokens,err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Error generating new tokens",
		})
		return
	}

	delete(refreshTokens, req.RefreshToken)

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: tokens,
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

		token := strings.TrimPrefix(authHeader, "Bearer ")
		jwtClaims, err := validateToken(token)
		if err != nil {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   "Invalid token provided",
			})
			c.Abort()
			return
		}

		c.Set("user_id", jwtClaims.UserID)
		c.Set("username", jwtClaims.Username)
		c.Set("role", jwtClaims.Role)

		// TODO: Extract token from "Bearer <token>" format
		// TODO: Validate token using validateToken function
		// TODO: Set user info in context for route handlers

		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if !slices.Contains(roles, role) {
			c.JSON(403, APIResponse{
				Success: false,
				Error:   "Unauthorised",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// TODO: Return user profile (without sensitive data)

	userId := c.GetInt("user_id")
	user := findUserByID(userId)

	if user == nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		c.Abort()
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data: User{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Role:           user.Role,
			IsActive:       user.IsActive,
			EmailVerified:  user.EmailVerified,
			LastLogin:      user.LastLogin,
			FailedAttempts: user.FailedAttempts,
			LockedUntil:    user.LockedUntil,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		},
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

	userId := c.GetInt("user_id")
	user := findUserByID(userId)
	if user == nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	userWithNewEmail := findUserByEmail(req.Email)
	if userWithNewEmail != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "New email already taken",
		})
		return
	}

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

	userId := c.GetInt("user_id")
	user := findUserByID(userId)
	if !verifyPassword(req.CurrentPassword, user.PasswordHash) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid password",
		})
		return
	}

	if !isStrongPassword(req.NewPassword) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "New password is too weak",
		})
		return
	}

	hash, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Error hashing new password",
		})
		return
	}
	user.PasswordHash = hash

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	rawPageSize := c.Query("page_size")
	if rawPageSize == "" {
		rawPageSize = "10"
	}
	pageSize,err := strconv.Atoi(rawPageSize)
	if err != nil || pageSize <= 0 {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid page size",
		})
		return
	}
	if pageSize > 100 {
		pageSize = 100
	}

	rawPage := c.Query("page")
	if rawPage == "" {
		rawPage = "1"
	}
	page,err := strconv.Atoi(rawPage)
	if err != nil || page <= 0 {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid page",
		})
		return
	}

	upper := page * pageSize
	lower := upper - pageSize
	if lower > len(users) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid page",
		})
		return
	}

	if upper > len(users) {
		upper = len(users) - 1 
	}

	var parsedUsers []User
	for _,user := range users[lower:upper] {
		parsedUsers = append(parsedUsers, User{
			ID:             user.ID,
			Username:       user.Username,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Role:           user.Role,
			IsActive:       user.IsActive,
			LastLogin:      user.LastLogin,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		},)
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    parsedUsers,
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
	if !slices.Contains(validRoles, req.Role) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid role",
		})
		return
	}

	user := findUserByID(id)
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
