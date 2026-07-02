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
	TokenUse string `json:"token_use"`
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

// isStrongPassword reports whether the password is at least 8 characters long
// and contains an uppercase letter, a lowercase letter, a digit, and a symbol.
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

// hashPassword returns the bcrypt hash of the given plaintext password.
func hashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(pass), nil
}

// verifyPassword reports whether the plaintext password matches the bcrypt hash.
func verifyPassword(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}

// generateTokens creates a signed access token and refresh token for the user,
// records the refresh token in the store, and returns them as a TokenResponse.
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
	claimsForAccess := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		TokenUse: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsForAccess)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	claimsForRefresh := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		TokenUse: "refresh",
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
		AccessToken:  accessTokenString,
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

// validateToken parses and verifies the signed JWT, ensuring it uses HMAC
// signing and is not blacklisted, and returns its claims.
func validateToken(tokenString string) (*JWTClaims, error) {
	claims := JWTClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected JWT signing method")
		}

		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	tokenMU.RLock()
	if _, exists := blacklistedTokens[tokenString]; exists {
		tokenMU.RUnlock()
		return nil, fmt.Errorf("token has been revoked")
	}
	tokenMU.RUnlock()

	return &claims, nil
}

// findUserByUsername returns a pointer to the user with the given username,
// or nil if no such user exists.
func findUserByUsername(username string) *User {
	if username == "" {
		return nil
	}
	usersMU.RLock()
	defer usersMU.RUnlock()

	for _, user := range users {
		if user.Username == username {

			userToReturn := user
			return &userToReturn
		}
	}

	return nil
}

// findUserByEmail returns a pointer to the user with the given email,
// or nil if no such user exists.
func findUserByEmail(email string) *User {
	if email == "" {
		return nil
	}
	usersMU.RLock()
	defer usersMU.RUnlock()

	for _, user := range users {
		if user.Email == email {

			userToReturn := user
			return &userToReturn
		}
	}

	return nil
}

// findUserByID returns a pointer to the user with the given ID,
// or nil if no such user exists.
func findUserByID(id int) *User {
	usersMU.RLock()
	defer usersMU.RUnlock()
	for _, user := range users {
		if user.ID == id {

			userToReturn := user
			return &userToReturn
		}
	}

	return nil
}

// isAccountLocked reports whether the user's account is currently locked,
// based on whether LockedUntil is set to a time in the future.
func isAccountLocked(user *User) bool {
	usersMU.RLock()
	defer usersMU.RUnlock()
	if user.LockedUntil == nil {
		return false
	}
	return time.Until(*user.LockedUntil) > 0
}

// recordFailedAttempt increments the user's failed login counter and locks the
// account for lockoutDuration once maxFailedAttempts is reached.
func recordFailedAttempt(user *User) {
	usersMU.Lock()
	defer usersMU.Unlock()

	user.FailedAttempts++
	now := time.Now()
	if user.LockedUntil == nil {
		user.LockedUntil = &now
	}

	if user.FailedAttempts >= maxFailedAttempts {
		*user.LockedUntil = time.Now().Add(lockoutDuration)
		return
	}
}

// resetFailedAttempts clears the user's failed login counter and unlocks the account.
func resetFailedAttempts(user *User) {
	usersMU.Lock()
	defer usersMU.Unlock()

	user.FailedAttempts = 0
	user.LockedUntil = nil
}

// generateRandomToken returns a cryptographically random 32-byte token,
// hex-encoded as a string.
func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// register handles POST /auth/register: it validates the registration payload,
// enforces password strength and uniqueness, creates the user, and returns
// freshly issued tokens.
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
		c.JSON(409, APIResponse{
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

	role := c.GetString("role")
	if role == "" {
		role = RoleUser
	}

	usersMU.Lock()
	defer usersMU.Unlock()
	now := time.Now()
	id := nextUserID
	nextUserID++

	user := User{
		ID:            id,
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          role,
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	users = append(users, user)
	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to generate tokens",
		})
		return
	}

	c.JSON(201, APIResponse{
		Success: true,
		Data:    tokens,
		Message: "User registered successfully",
	})
}

// login handles POST /auth/login: it authenticates the credentials, enforces
// account lockout on repeated failures, and returns access and refresh tokens.
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

// logout handles POST /auth/logout: it blacklists the presented access token
// and removes the user's refresh token from the store.
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
	if len(tokenS) != 2 || tokenS[0] != "Bearer" {
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
	if claims.TokenUse != "access" {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}

	tokenMU.Lock()
	blacklistedTokens[token] = true

	for r, id := range refreshTokens {
		if id == claims.UserID {
			delete(refreshTokens, r)
		}
	}
	tokenMU.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// refreshToken handles POST /auth/refresh: it validates the supplied refresh
// token against the store, then issues a new token pair and rotates the refresh
// token.
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

	claims, err := validateToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}
	if claims.TokenUse != "refresh" {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}
	userID := claims.UserID
	user := findUserByID(userID)
	if user == nil {
		c.JSON(404, APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	tokenMU.Lock()
	if storedUserID, exists := refreshTokens[req.RefreshToken]; !exists || storedUserID != userID {
		c.JSON(401, APIResponse{
			Success: false,
			Error:   "Refresh token not recognized",
		})
		tokenMU.Unlock()
		return
	}
	delete(refreshTokens, req.RefreshToken)
	tokenMU.Unlock()

	tokens, err := generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to generate new tokens",
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    tokens,
		Message: "Token refreshed successfully",
	})
}

// authMiddleware returns Gin middleware that requires a valid "Bearer" JWT in
// the Authorization header and stores the user's ID and role in the context.
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
		if len(tokenS) != 2 || tokenS[0] != "Bearer" {
			c.JSON(400, APIResponse{
				Success: false,
				Error:   "Wrong token type",
			})
			c.Abort()
			return
		}
		token := tokenS[1]
		claims, err := validateToken(token)
		if err != nil || claims.TokenUse != "access" {
			c.JSON(401, APIResponse{
				Success: false,
				Error:   "Invalid token",
			})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// requireRole returns Gin middleware that authorizes the request only if the
// role stored in the context (by authMiddleware) is one of the allowed roles.
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if !slices.Contains(roles, role) {
			c.JSON(403, APIResponse{
				Success: false,
				Error:   "Forbidden: insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getUserProfile handles GET /user/profile: it returns the authenticated
// user's profile.
func getUserProfile(c *gin.Context) {
	id := c.GetInt("userID")
	user := findUserByID(id)

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

// updateUserProfile handles PUT /user/profile: it validates and applies updates
// to the authenticated user's profile.
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
	id := c.GetInt("userID")
	usersMU.Lock()
	for _, user := range users {
		if user.Email == req.Email {
			c.JSON(409, APIResponse{
				Success: false,
				Error:   "Email already in use",
			})
			usersMU.Unlock()
			return
		}
	}
	for i, user := range users {
		if user.ID == id {
			users[i].FirstName = req.FirstName
			users[i].LastName = req.LastName
			users[i].Email = req.Email
			users[i].UpdatedAt = time.Now()
			break
		}
	}
	usersMU.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Profile updated successfully",
	})
}

// changePassword handles POST /user/change-password: it verifies the current
// password, enforces new-password strength, and stores the new hash.
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

	id := c.GetInt("userID")
	user := findUserByID(id)
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
			Error:   "New password does not meet strength requirements",
		})
		return
	}
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   "Failed to hash new password",
		})
		return
	}
	usersMU.Lock()
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()
	usersMU.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// listUsers handles GET /admin/users: it returns a snapshot of all users.
func listUsers(c *gin.Context) {
	usersMU.RLock()
	usersList := make([]User, len(users))
	copy(usersList, users)
	usersMU.RUnlock()
	c.JSON(200, APIResponse{
		Success: true,
		Data:    usersList,
		Message: "Users retrieved successfully",
	})
}

// changeUserRole handles PUT /admin/users/:id/role: it validates the requested
// role and applies it to the target user.
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

	validRoles := []string{RoleUser, RoleAdmin, RoleModerator}
	if !slices.Contains(validRoles, req.Role) {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid role",
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
	usersMU.Lock()
	for i, u := range users {
		if u.ID == id {
			users[i].Role = req.Role
			users[i].UpdatedAt = time.Now()
			break
		}
	}
	usersMU.Unlock()

	c.JSON(200, APIResponse{
		Success: true,
		Message: "User role updated successfully",
	})
}

// setupRouter builds the Gin engine with the public auth routes and the
// authenticated user and admin route groups.
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

// main seeds a default admin user and starts the HTTP server on port 8080.
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
