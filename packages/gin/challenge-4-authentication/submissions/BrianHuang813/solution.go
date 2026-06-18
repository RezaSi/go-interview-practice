package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors" 
	"fmt"
	"net/http"
	"strings"
	"time"
	"strconv"

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
	// - Contains uppercase letter
	// - Contains lowercase letter
	// - Contains number
	// - Contains special character
	if len(password) < 8 {
        return false
    }
    
    hasUpper := false
    hasLower := false
    hasDigit := false
    hasSpecial := false
    for _, char := range password {
        switch {
        case 'A' <= char && char <= 'Z':
            hasUpper = true
        case 'a' <= char && char <= 'z':
            hasLower = true
        case '0' <= char && char <= '9':
            hasDigit = true
        default:
            hasSpecial = true
        }
    }
    return hasUpper && hasLower && hasDigit && hasSpecial
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
	// TODO: Generate access token with 15 minute expiry
	// TODO: Generate refresh token with 7 day expiry
	// TODO: Store refresh token in memory store
	
    // Access Token
    accessClaims := &JWTClaims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "your-app",
        },
    }
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString(jwtSecret)
    if err != nil {
        return nil, err
    }
    // Refresh Token
    refreshToken, err := generateRandomToken()
    if err != nil {
        return nil, err
    }
    // Store refresh token
    refreshTokens[refreshToken] = userID
    return &TokenResponse{
        AccessToken:  accessTokenString,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    int64(accessTokenTTL.Seconds()),
        ExpiresAt:    time.Now().Add(accessTokenTTL),
    }, nil
}


// TODO: Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	// TODO: Parse and validate JWT token
	// TODO: Check if token is blacklisted
	// TODO: Return claims if valid
	if blacklistedTokens[tokenString] {
	    return nil, errors.New("the token is blacklisted")
	}
	
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}

// TODO: Implement user lookup functions
func findUserByUsername(username string) *User {
	// TODO: Find user by username in users slice
	if username == "" {
	    return nil
	}
	
	for i, user := range users {
	    if user.Username == username {
	        return &users[i]
	    }
	}
	return nil
}

func findUserByEmail(email string) *User {
	// TODO: Find user by email in users slice
	if email == "" {
	    return nil
	}
	
	for i, user := range users {
	    if user.Email == email {
	        return &users[i]
	    }
	}
	return nil
}

func findUserByID(id int) *User {
	// TODO: Find user by ID in users slice
	for i, user := range users {
	    if user.ID == id {
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
	return (*user.LockedUntil).After(time.Now())
}

// TODO: Implement failed attempt tracking
func recordFailedAttempt(user *User) {
	// TODO: Increment failed attempts counter
	// TODO: Lock account if max attempts reached
	user.FailedAttempts += 1
	
	if user.FailedAttempts >= maxFailedAttempts {
	    unlockTime := time.Now().Add(lockoutDuration)
	    user.LockedUntil = &unlockTime
	    user.FailedAttempts = 0 // reset the failed attempts after penalty
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
	// TODO: Check if email already exists
	if findUserByUsername(req.Username) != nil || findUserByEmail(req.Email) != nil {
	    c.JSON(http.StatusConflict, APIResponse{
	        Success: false,
	        Error: "Username or email is already taken",
	    })
	    return
	}
	// TODO: Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		// This is a server error, not the user's fault.
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to hash password",
		})
		return
	}
	// TODO: Create user and add to users slice
	newUser := User {
	    ID:            nextUserID,
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  string(passwordHash), 
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          "user",  
		IsActive:      true,  
		EmailVerified: false, 
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	users = append(users, newUser)
	nextUserID++

	c.JSON(201, APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data: newUser,
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
	userID, exist := refreshTokens[req.RefreshToken]
	if !exist {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid or expired refresh token",
		})
		return
	}
	// TODO: Get user ID from refresh token store
	// TODO: Find user by ID
	user := findUserByID(userID)
	if user == nil || !user.IsActive {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "User not found or is not active",
		})
		return
	}
	// TODO: Generate new access token
	delete(refreshTokens, req.RefreshToken)

    newTokens, err := generateTokens(user.ID, user.Username, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to generate new tokens"})
        return
    }
    
    c.JSON(http.StatusOK, APIResponse{
        Success: true,
        Message: "Token refreshed successfully",
        Data:    newTokens,
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
		    c.JSON(http.StatusUnauthorized, APIResponse{
		        Success: false,
		        Error: "invalid token",
		    })
		    return
		}
		// TODO: Set user info in context for route handlers
		c.Set("claims", claims)

		c.Set("user", &User{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		})
		
		c.Next()
	}
}

// Middleware: Role-based authorization
func requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get user role from context (set by authMiddleware)
		userCtx, exists := c.Get("user")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
                Success: false,
                Error:   "User data not found in context. Access denied.",
            })
            return
        }

        currentUser, ok := userCtx.(*User)
        if !ok {
            c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
                Success: false,
                Error:   "Invalid user type in context.",
            })
            return
        }
        
        
        isAllowed := false
        for _, allowedRole := range roles {
            if currentUser.Role == allowedRole {
                isAllowed = true
                break
            }
        }
        
		if isAllowed {
            c.Next()
        } else {
            c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
                Success: false,
                Error:   fmt.Sprintf("Access denied. Required roles: %v, your role: %s", roles, currentUser.Role),
            })
        }
	}
}

// GET /user/profile - Get current user profile
func getUserProfile(c *gin.Context) {
	// TODO: Get user ID from context (set by authMiddleware)
	userCtx, _ := c.Get("user")

	user, ok := userCtx.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Internal server error: user data in context is corrupted",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    user,
		Message: "Profile retrieved successfully",
	})
}

func updateUserProfile(c *gin.Context) {
	// 1. Get the current user from the context (set by AuthMiddleware)
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{Success: false, Error: "User not found in context"})
		return
	}
	// Type assertion to get the *User object
	currentUser, ok := userCtx.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: "Invalid user type in context"})
		return
	}

	// 2. Bind the incoming JSON data
	var req struct {
		FirstName string `json:"first_name" binding:"required,min=2,max=50"`
		LastName  string `json:"last_name" binding:"required,min=2,max=50"`
		Email     string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid input data: " + err.Error()})
		return
	}

	// 3. Check if the new email is already taken by ANOTHER user
	if currentUser.Email != req.Email {
		if existingUser := findUserByEmail(req.Email); existingUser != nil {
			c.JSON(http.StatusConflict, APIResponse{Success: false, Error: "Email is already taken"})
			return
		}
	}

	// 4. Update the user's data
	currentUser.FirstName = req.FirstName
	currentUser.LastName = req.LastName
	currentUser.Email = req.Email
	currentUser.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    currentUser,
		Message: "Profile updated successfully",
	})
}

func changePassword(c *gin.Context) {
	
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{Success: false, Error: "Unauthorized"})
		return
	}
	claims, ok := claimsVal.(*JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{Success: false, Error: "Invalid claims"})
		return
	}

	currentUser := findUserByID(claims.UserID)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, APIResponse{Success: false, Error: "User not found"})
		return
	} 

    
	// Bind the incoming JSON data
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid input data: " + err.Error()})
		return
	}

	// 3. Verify the user's CURRENT password
	//    This is crucial to ensure the person changing the password is the legitimate user.
	err := bcrypt.CompareHashAndPassword([]byte(currentUser.PasswordHash), []byte(req.CurrentPassword))
	fmt.Println("CurrentPasswordHash: ", currentUser.PasswordHash)
	if err != nil {
		// If err is not nil, it means the password does not match.
		c.JSON(400, APIResponse{Success: false, Error: "Incorrect current password"})
		return
	}

	// 4. (Optional) You can add extra validation for the new password here if needed
	// if !isStrongPassword(req.NewPassword) { ... }

	// 5. Hash the NEW password and update the user
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to process new password"})
		return
	}
	currentUser.PasswordHash = string(newPasswordHash)
	currentUser.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, APIResponse{Success: true, Message: "Password changed successfully"})
}

// listUsers handles GET /admin/users - Lists all users (admin only)
func listUsers(c *gin.Context) {
    
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "0")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }

	startIndex := (page - 1) * limit
	if startIndex >= len(users) {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: make([]User, 0)})
		return
	}
	endIndex := startIndex + limit
	if endIndex > len(users) {
		endIndex = len(users)
	}

	paginatedUsers := users[startIndex:endIndex]

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    paginatedUsers,
		Message: "Users retrieved successfully",
	})
}


// changeUserRole handles PUT /admin/users/:id/role - Changes a user's role (admin only)
func changeUserRole(c *gin.Context) {
	// 1. Get and validate the user ID from the URL path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid user ID"})
		return
	}

	// 2. Bind the new role from the request body
	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid role data"})
		return
	}

	// Validate that the provided role is one of the valid, predefined roles
	validRoles := map[string]bool{RoleAdmin: true, RoleModerator: true, RoleUser: true}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid role specified"})
		return
	}

	// Find the user to be updated
	user := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "User not found"})
		return
	}

	// Update the user's role and save it
	user.Role = req.Role
	user.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    user,
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
