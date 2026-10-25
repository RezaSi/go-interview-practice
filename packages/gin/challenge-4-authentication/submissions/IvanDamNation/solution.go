package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrBlacklisted      = errors.New("token is in blacklist")
	ErrMailCollision    = errors.New("this email already taken by other user")
	ErrInvalidCreds     = errors.New("invalid credentials")
	ErrInvalidFormatID  = errors.New("ID invalid format")
	ErrInvalidInputData = errors.New("Invalid input data")
	ErrInvalidRole      = errors.New("Invalid role data")
	ErrInvalidToken     = errors.New("invalid token")
	ErrHashingError     = errors.New("Password hashing error")
	ErrNoAuthHeader     = errors.New("Authorization header required")
	ErrNoBearerPrefix   = errors.New("No bearer prefix")
	ErrNoUserID         = errors.New("user ID is not provided")
	ErrUserNotFound     = errors.New("user not found")
	ErrWeakPassword     = errors.New("Password does not meet strength requirements")
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
	Role           Roles      `json:"role"`
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
	Type     string `json:"typ,omitempty"`
	jwt.RegisteredClaims
}

// APIResponse represents standard API response
type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
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

// User string
type Roles string

const (
	RoleUser      = "user"
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
)

// Exponential work factor
const hashCostFactor = 12

// isStrongPassword validates password strength
func isStrongPassword(password string) bool {
	pasRune := []rune(password)
	if len(pasRune) < 8 {
		return false
	}

	upper, lower, num, sChar := false, false, false, false
	for _, r := range pasRune {
		if unicode.IsUpper(r) {
			upper = true
		}
		if unicode.IsLower(r) {
			lower = true
		}
		if unicode.IsDigit(r) {
			num = true
		}
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			sChar = true
		}
	}

	if !upper || !lower || !num || !sChar {
		return false
	}

	return true
}

// hashPassword returns hash of password
func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), hashCostFactor)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// verifyPassword verifying provided password and it's hash
func verifyPassword(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

// generateTokens generates new JWT tokens
func generateTokens(userID int, username string, role Roles) (*TokenResponse, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     string(role),
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        generateJTI(),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	atStr, err := at.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	claims.ID = generateJTI()
	claims.Type = "refresh"
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(refreshTokenTTL))

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rtStr, err := rt.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshTokens[rtStr] = userID

	return &TokenResponse{
		AccessToken:  atStr,
		RefreshToken: rtStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    time.Now().Add(accessTokenTTL),
	}, nil
}

func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// validateToken validates provided JWT token
func validateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if blacklistedTokens[tokenString] {
		return nil, ErrBlacklisted
	}

	return claims, nil
}

// sanitizeToken
func sanitizeToken(headerWithToken string) (string, error) {
	if !strings.HasPrefix(headerWithToken, "Bearer ") {
		return "", ErrNoBearerPrefix
	}

	token := strings.TrimPrefix(headerWithToken, "Bearer ")

	return token, nil
}

// findUserByUsername is lookup functions by username
func findUserByUsername(username string) *User {
	for i, u := range users {
		if u.Username == username {
			return &users[i]
		}
	}
	return nil
}

// findUserByEmail is lookup functions by email
func findUserByEmail(email string) *User {
	for i, u := range users {
		if u.Email == email {
			return &users[i]
		}
	}
	return nil
}

// findUserByID is lookup functions by ID
func findUserByID(id int) *User {
	for i, u := range users {
		if u.ID == id {
			return &users[i]
		}
	}
	return nil
}

// extractUserID returns user ID from context
func extractUserID(c *gin.Context) (int, error) {
	idVal, exists := c.Get("userID")
	if !exists {
		return -1, ErrNoUserID
	}

	userID, ok := idVal.(int)
	if !ok {
		return -1, ErrInvalidFormatID
	}

	return userID, nil
}

// isAccountLocked check account lockout
func isAccountLocked(user *User) bool {
	if user.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*user.LockedUntil)
}

// recordFailedAttempt tracks password checks failed attempts
func recordFailedAttempt(user *User) {
	user.FailedAttempts++
	if user.FailedAttempts >= maxFailedAttempts {
		lockTime := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockTime
		user.IsActive = false
	}
}

// resetFailedAttempts resets counter failed attempts and unlocks user
func resetFailedAttempts(user *User) {
	user.FailedAttempts = 0
	user.LockedUntil = nil
	user.IsActive = true
}

// POST /auth/register - User registration
func register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   ErrInvalidInputData.Error(),
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Passwords do not match",
		})
		return
	}

	if !isStrongPassword(req.Password) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   ErrWeakPassword.Error(),
		})
		return
	}

	var u *User
	if u = findUserByUsername(req.Username); u != nil {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   "User with this username already exists",
		})
		return
	}

	if u = findUserByEmail(req.Email); u != nil {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   "User with this email already exists",
		})
		return
	}

	h, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   ErrHashingError.Error(),
		})
		return
	}

	now := time.Now()

	newUser := User{
		ID:           nextUserID,
		Username:     req.Username,
		PasswordHash: h,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Role:         RoleUser,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	users = append(users, newUser)
	nextUserID++

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    newUser,
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

	user := findUserByUsername(req.Username)
	if user == nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   ErrInvalidCreds.Error(),
		})
		return
	}

	if isAccountLocked(user) {
		c.JSON(http.StatusLocked, APIResponse{
			Success: false,
			Error:   "Account is temporarily locked",
		})
		return
	}

	if !verifyPassword(req.Password, user.PasswordHash) {
		recordFailedAttempt(user)
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   ErrInvalidCreds.Error(),
		})
		return
	}

	resetFailedAttempts(user)

	now := time.Now()
	user.LastLogin = &now

	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
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
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   ErrNoAuthHeader.Error(),
		})
		return
	}

	token, err := sanitizeToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	blacklistedTokens[token] = true

	for k, v := range refreshTokens {
		if v == userID {
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

	claims, err := validateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if claims.Type != "refresh" {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   ErrInvalidToken.Error(),
		})
		return
	}

	userID, exists := refreshTokens[req.RefreshToken]
	if !exists || userID != claims.UserID {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   ErrInvalidToken.Error(),
		})
		return
	}
	delete(refreshTokens, req.RefreshToken)

	user := findUserByID(userID)
	if user == nil || isAccountLocked(user) {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   ErrInvalidToken.Error(),
		})
		return
	}

	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to generate tokens",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   ErrNoAuthHeader.Error(),
			})
			return
		}

		token, err := sanitizeToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		claims, err := validateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		if claims.Type != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   ErrInvalidToken.Error(),
			})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(allowedRoles ...Roles) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
				Success: false,
				Error:   "No role provided",
			})
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
				Success: false,
				Error:   "Invalid role format",
			})
			return
		}

		allowed := false
		for _, r := range allowedRoles {
			if string(r) == userRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
				Success: false,
				Error:   "Insufficient permissions",
			})
			return
		}

		c.Next()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		if errors.Is(err, ErrNoUserID) {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		if errors.Is(err, ErrInvalidFormatID) {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
	}

	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   ErrUserNotFound.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   ErrInvalidInputData.Error(),
		})
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   ErrUserNotFound.Error(),
		})
		return
	}

	user2 := findUserByEmail(req.Email)
	if user2 != nil && user2 != user {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   ErrMailCollision.Error(),
		})
		return
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    user,
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
			Error:   ErrInvalidInputData.Error(),
		})
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   ErrUserNotFound.Error(),
		})
		return
	}

	if !verifyPassword(req.CurrentPassword, user.PasswordHash) {
		recordFailedAttempt(user)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   ErrInvalidCreds.Error(),
		})
		return
	}

	if !isStrongPassword(req.NewPassword) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   ErrWeakPassword.Error(),
		})
		return
	}

	hash, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   ErrHashingError.Error(),
		})
		return
	}

	user.PasswordHash = hash
	user.UpdatedAt = time.Now()

	for k, v := range refreshTokens {
		if v == user.ID {
			delete(refreshTokens, k)
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// UserPublicResponse represents safe data return for public requests
type UserPublicResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      Roles     `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// GET /admin/users - List all users (admin only)
func listUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	startIdx := (page - 1) * limit
	endIdx := startIdx + limit

	totalUsers := len(users)

	if startIdx > totalUsers {
		startIdx = totalUsers
	}

	if endIdx > totalUsers {
		endIdx = totalUsers
	}

	paginatedUsers := users[startIdx:endIdx]

	pubUsers := make([]UserPublicResponse, 0, len(paginatedUsers))
	for _, u := range paginatedUsers {
		pubUsers = append(pubUsers, UserPublicResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"users": pubUsers,
			"pagination": gin.H{
				"current_page": page,
				"limit":        limit,
				"total":        totalUsers,
			},
		},
		Message: "Users retrieved successfully",
	})
}

// PUT /admin/users/:id/role - Change user role (admin only)
func changeUserRole(c *gin.Context) {
	userIDVal := c.Param("id")
	userID, err := strconv.Atoi(userIDVal)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   ErrInvalidFormatID.Error(),
		})
		return
	}

	var req struct {
		Role Roles `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   ErrInvalidRole.Error(),
		})
		return
	}

	if req.Role != RoleAdmin && req.Role != RoleModerator && req.Role != RoleUser {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   ErrInvalidRole.Error(),
		})
		return
	}

	user := findUserByID(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   ErrUserNotFound.Error(),
		})
		return
	}

	user.Role = req.Role
	user.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "User role updated successfully",
	})
}

// Setup router with authentication routes
func setupRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/logout", authMiddleware(), logout)
		auth.POST("/refresh", refreshToken)
	}

	user := router.Group("/user")
	user.Use(authMiddleware())
	{
		user.GET("/profile", getUserProfile)
		user.PUT("/profile", updateUserProfile)
		user.POST("/change-password", changePassword)
	}

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
