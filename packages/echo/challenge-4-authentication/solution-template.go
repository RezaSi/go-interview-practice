package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in our system
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"` // Never include in JSON responses
	Role      string `json:"role"` // "user" or "admin"
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Request/Response structures
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	User         *User  `json:"user,omitempty"`
}

type UserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

type UserListResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	Users  []User `json:"users"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// JWT Claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// In-memory storage
var users []User
var blacklistedTokens = make(map[string]bool)

// JWT secret key (in production, use environment variable)
var jwtSecret = []byte("your-secret-key")

func main() {
	// TODO: Create Echo instance, setup middleware and routes
	// TODO: Start server on port 8080
}

func setupMiddleware(e *echo.Echo) {
	// TODO: Setup basic middleware (logger, recovery, CORS)
}

func setupRoutes(e *echo.Echo) {
	// TODO: Setup authentication routes (no auth required)
	// POST /auth/register - User registration
	// POST /auth/login - User login
	// POST /auth/refresh - Refresh token
	// POST /auth/logout - User logout

	// TODO: Setup protected user routes (auth required)
	// GET /users/profile - Get current user profile
	// PUT /users/profile - Update user profile

	// TODO: Setup admin routes (admin role required)
	// GET /users - List all users
	// PUT /users/:id/role - Update user role
	// DELETE /users/:id - Delete user
}

// TODO: Implement JWT middleware for authentication
func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: Extract JWT token from Authorization header
		// TODO: Validate and parse JWT token
		// TODO: Check if token is blacklisted
		// TODO: Set user info in context
		// TODO: Call next handler
		return next(c)
	}
}

// TODO: Implement admin role middleware
func adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: Check if user has admin role
		// TODO: Return 403 if not admin
		// TODO: Call next handler if admin
		return next(c)
	}
}

// TODO: Implement password hashing
func hashPassword(password string) (string, error) {
	// TODO: Use bcrypt to hash password
	return "", nil
}

// TODO: Implement password verification
func verifyPassword(hashedPassword, password string) bool {
	// TODO: Use bcrypt to compare passwords
	return false
}

// TODO: Implement JWT token generation
func generateTokens(user *User) (string, string, error) {
	// TODO: Create access token (15 minutes expiry)
	// TODO: Create refresh token (7 days expiry)
	// TODO: Return both tokens
	return "", "", nil
}

// TODO: Implement JWT token validation
func validateToken(tokenString string) (*JWTClaims, error) {
	// TODO: Parse and validate JWT token
	// TODO: Return claims if valid
	return nil, nil
}

func registerHandler(c echo.Context) error {
	// TODO: Implement user registration
	// 1. Bind and validate request
	// 2. Check if username/email already exists
	// 3. Hash password
	// 4. Create user with default role "user"
	// 5. Generate tokens
	// 6. Return success response
}

func loginHandler(c echo.Context) error {
	// TODO: Implement user login
	// 1. Bind and validate request
	// 2. Find user by username
	// 3. Verify password
	// 4. Generate tokens
	// 5. Return success response
}

func refreshHandler(c echo.Context) error {
	// TODO: Implement token refresh
	// 1. Extract refresh token
	// 2. Validate refresh token
	// 3. Generate new access token
	// 4. Return new token
}

func logoutHandler(c echo.Context) error {
	// TODO: Implement user logout
	// 1. Extract access token
	// 2. Add token to blacklist
	// 3. Return success response
}

func getProfileHandler(c echo.Context) error {
	// TODO: Get current user profile from context
	// Return user information
}

func updateProfileHandler(c echo.Context) error {
	// TODO: Update current user profile
	// Allow updating email only (not username/password/role)
}

func getAllUsersHandler(c echo.Context) error {
	// TODO: Return all users (admin only)
}

func updateUserRoleHandler(c echo.Context) error {
	// TODO: Update user role (admin only)
	// Allow changing between "user" and "admin"
}

func deleteUserHandler(c echo.Context) error {
	// TODO: Delete user (admin only)
	// Don't allow deleting self
}

