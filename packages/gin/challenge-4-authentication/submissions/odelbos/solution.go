package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"slices"
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
var usersMutex sync.RWMutex
var blacklistedTokens = make(map[string]bool) // Token blacklist for logout
var blacklistMutex sync.RWMutex
var refreshTokens = make(map[string]int)      // RefreshToken -> UserID mapping
var refreshMutex sync.RWMutex
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

// ---------------------------------------------------------------
// Password security
// ---------------------------------------------------------------

func isStrongPassword(password string) bool {
	// At least 8 characters
	// Contains uppercase letter
	// Contains lowercase letter
	// Contains number
	// Contains special character
	upper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	lower := regexp.MustCompile(`[a-z]`).MatchString(password)
	number := regexp.MustCompile(`[0-9]`).MatchString(password)
	special := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)
	return len(password) > 7 && upper && lower && number && special
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func verifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// ---------------------------------------------------------------
// JWT functions
// ---------------------------------------------------------------

func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
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

	refreshMutex.Lock()
	refreshTokens[refreshToken] = userID
	refreshMutex.Unlock()

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
		ExpiresAt:    now.Add(accessTokenTTL),
	}, nil
}

func validateToken(tokenString string) (*JWTClaims, error) {
	// Parse and validate JWT token
	// Check if token is blacklisted
	// Return claims if valid
	blacklistMutex.RLock()
	if blacklistedTokens[tokenString] {
		blacklistMutex.RUnlock()
		return nil, fmt.Errorf("invalid token")
	}
	blacklistMutex.RUnlock()

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ! ok {
			return nil, fmt.Errorf("invalid token")
		}
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

// ---------------------------------------------------------------
// User functions
// ---------------------------------------------------------------

func findUserByUsername(username string) *User {
	usersMutex.RLock()
	defer usersMutex.RUnlock()
	for _, user := range(users) {
		if user.Username == username {
			return &user
		}
	}
	return nil
}

func findUserByEmail(email string) *User {
	usersMutex.RLock()
	defer usersMutex.RUnlock()
	for _, user := range(users) {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

func findUserByID(id int) *User {
	usersMutex.RLock()
	defer usersMutex.RUnlock()
	for _, user := range(users) {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

func isAccountLocked(user *User) bool {
	// Check if account is locked based on LockedUntil field
	return user.LockedUntil != nil && time.Now().Before(*user.LockedUntil)
}

func recordFailedAttempt(user *User) {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	user.FailedAttempts++
	if user.FailedAttempts >= maxFailedAttempts {
		lockTime := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockTime
		user.UpdatedAt = time.Now()
	}
}

func resetFailedAttempts(user *User) {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	user.FailedAttempts = 0
	user.LockedUntil = nil
	user.UpdatedAt = time.Now()
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ---------------------------------------------------------------
// Route handlers
// ---------------------------------------------------------------

func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}
	if req.Password != req.ConfirmPassword {
		errResponse(c, http.StatusBadRequest, "Not mathing password")
		return
	}
	if ! isStrongPassword(req.Password) {
		errResponse(c, http.StatusBadRequest, "Invalid password")
		return
	}

	if findUserByUsername(req.Username) != nil {
		errResponse(c, http.StatusConflict, "Username already exists")
		return
	}
	if findUserByEmail(req.Email) != nil {
		errResponse(c, http.StatusConflict, "Email already exists")
		return
	}

	pwdHash, err := hashPassword(req.Password)
	if err != nil {
		errResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	now := time.Now()
	user := User{
		ID:            nextUserID,
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  pwdHash,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          RoleUser,
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	users = append(users, user)
	nextUserID++
	okResponse(c, http.StatusCreated, "User registered successfully", nil)
}

func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	user := findUserByUsername(req.Username)
	if user == nil {
		errResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if isAccountLocked(user) {
		errResponse(c, http.StatusLocked, "Account is locked")
	}

	if ! verifyPassword(req.Password, user.PasswordHash) {
		recordFailedAttempt(user)
		errResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}

	resetFailedAttempts(user)

	usersMutex.Lock()
	defer usersMutex.Unlock()

	now := time.Now()
	user.LastLogin = &now

	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		errResponse(c, http.StatusInternalServerError, "Internal server error")
	}
	okResponse(c, http.StatusOK, "Login successful", tokens)
}

func logout(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		errResponse(c, http.StatusUnauthorized, "Authorization header required")
		return
	}
	if ! strings.HasPrefix(bearer, "Bearer ") {
		errResponse(c, http.StatusUnauthorized, "Invalid token")
		c.Abort()
		return
	}

	tokenStr := strings.TrimPrefix(bearer, "Bearer ")
	_, err := validateToken(tokenStr)
	if err != nil {
		errResponse(c, http.StatusUnauthorized, "Invalid token")
		c.Abort()
		return
	}

	blacklistMutex.Lock()
	blacklistedTokens[tokenStr] = true
	blacklistMutex.Unlock()

	// NOTE: Should remove all refresh tokens owned by the user

	okResponse(c, http.StatusOK, "Logout successful", nil)
}

func refreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	refreshMutex.Lock()
	userId, ok := refreshTokens[req.RefreshToken]
	refreshMutex.Unlock()
	if ! ok {
		errResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	user := findUserByID(userId)
	if user == nil {
		errResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		errResponse(c, http.StatusInternalServerError, "Internal server error")
	}
	okResponse(c, http.StatusOK, "Login successful", tokens)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			errResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}
		if ! strings.HasPrefix(bearer, "Bearer ") {
			errResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(bearer, "Bearer ")
		claims, err := validateToken(tokenStr)
		if err != nil {
			errResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		for _, r := range(roles) {
			if role == r {
				c.Next()
				return
			}
		}
		errResponse(c, http.StatusForbidden, "Forbidden")
		c.Abort()
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	userId, _ := c.Get("user_id")
	user := findUserByID(userId.(int))
	if user == nil {
		errResponse(c, http.StatusNotFound, "Not found")
		return
	}

	// Return user profile (without sensitive data)
	type safeUser struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Role      string    `json:"role"`
		IsActive  bool      `json:"is_active"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	usersMutex.RLock()
	defer usersMutex.RUnlock()
	result := safeUser{
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
	okResponse(c, http.StatusOK, "User profile", result)
}

func updateUserProfile(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" binding:"required,min=2,max=50"`
		LastName  string `json:"last_name" binding:"required,min=2,max=50"`
		Email     string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	userId, _ := c.Get("user_id")
	user := findUserByID(userId.(int))
	if user == nil {
		errResponse(c, http.StatusNotFound, "Not found")
		return
	}
	if findUserByEmail(req.Email) != nil {
		errResponse(c, http.StatusConflict, "Email already exists")
		return
	}

	usersMutex.Lock()
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.UpdatedAt = time.Now()
	usersMutex.Unlock()
	okResponse(c, http.StatusOK, "Profile updated successfully", nil)
}

func changePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	userId, _ := c.Get("user_id")
	user := findUserByID(userId.(int))
	if user == nil {
		errResponse(c, http.StatusNotFound, "Not found")
		return
	}

	if ! verifyPassword(req.CurrentPassword, user.PasswordHash) {
		errResponse(c, http.StatusBadRequest, "Incorrect password")
		return
	}
	if ! isStrongPassword(req.NewPassword) {
		errResponse(c, http.StatusBadRequest, "Invalid password")
		return
	}
	pwdHash, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to hash new password"})
		return
	}

	usersMutex.Lock()
	user.PasswordHash = pwdHash
	user.UpdatedAt = time.Now()
	usersMutex.Unlock()
	okResponse(c, http.StatusOK, "Password changed successfully", nil)
}

func listUsers(c *gin.Context) {
	type safeUser struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Role      string    `json:"role"`
		IsActive  bool      `json:"is_active"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	usersMutex.RLock()
	defer usersMutex.RUnlock()

	var results []safeUser
	for _, u := range(users) {
		results = append(results, safeUser{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	okResponse(c, http.StatusOK, "Users list", results)
}

func changeUserRole(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid Id")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	validRoles := []string{RoleUser, RoleAdmin, RoleModerator}
	if ! slices.Contains(validRoles, req.Role) {
		errResponse(c, http.StatusBadRequest, "Invalid role")
		return
	}

	user := findUserByID(userId)
	if user == nil {
		errResponse(c, http.StatusNotFound, "Not found")
		return
	}

	usersMutex.Lock()
	user.Role = req.Role
	user.UpdatedAt = time.Now()
	usersMutex.Unlock()
	okResponse(c, http.StatusOK, "User role updated successfully", nil)
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

// ---------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------

func okResponse(c *gin.Context, status int, msg string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func errResponse(c *gin.Context, status int, msg string) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: msg,
	})
}

// ---------------------------------------------------------------
// Main
// ---------------------------------------------------------------

func main() {
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
