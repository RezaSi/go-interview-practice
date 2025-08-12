package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-"`    // Never return in JSON
	Role     string `json:"role"` // "user" or "admin"
	Active   bool   `json:"active"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
}

// In-memory storage (in production, use a database)
var users = []User{
	{ID: 1, Username: "admin", Email: "admin@example.com", Password: "$2a$12$...", Role: "admin", Active: true},
	{ID: 2, Username: "user1", Email: "user1@example.com", Password: "$2a$12$...", Role: "user", Active: true},
}
var nextUserID = 3

// JWT secret (in production, use environment variable)
var jwtSecret = []byte("your-secret-key")

var validate *validator.Validate

func main() {
	// TODO: Create Fiber app
	app := fiber.New()

    // Setup custom validator
	setupCustomValidator()

	// TODO: Setup routes
	// Public routes
	// app.Post("/auth/register", registerHandler)
	// app.Post("/auth/login", loginHandler)
	// app.Get("/health", healthHandler)

	// Protected routes (require JWT)
	// protected := app.Group("/", jwtMiddleware())
	// protected.Get("/profile", getProfileHandler)
	// protected.Put("/profile", updateProfileHandler)
	// protected.Post("/auth/refresh", refreshTokenHandler)

	// Admin routes (require admin role)
	// admin := app.Group("/admin", jwtMiddleware(), adminMiddleware())
	// admin.Get("/users", listUsersHandler)
	// admin.Put("/users/:id/role", updateUserRoleHandler)

	// TODO: Start server on port 3000
}

func setupCustomValidator() {
	// TODO: Setup validator with custom password validator
}

// TODO: Implement password security

// validatePassword validates password strength
func validatePassword(password string) bool {
	// TODO: Implement password validation
	// - Minimum 8 characters
	// - At least one uppercase letter
	// - At least one lowercase letter
	// - At least one digit
	// - At least one special character
	return false
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	// TODO: Use bcrypt to hash password with cost 12
	return "", nil
}

// verifyPassword compares a password with its hash
func verifyPassword(password, hash string) bool {
	// TODO: Use bcrypt to compare password with hash
	return false
}

// TODO: Implement JWT functions

// generateJWT creates a JWT token for a user
func generateJWT(user User) (string, error) {
	// TODO: Create JWT token with claims
	// - User ID, username, role
	// - Expiration: 1 hour from now
	// - Sign with HS256 algorithm
	return "", nil
}

// validateJWT validates and parses a JWT token
func validateJWT(tokenString string) (*JWTClaims, error) {
	// TODO: Parse and validate JWT token
	// - Verify signature
	// - Check expiration
	// - Return claims if valid
	return nil, nil
}

// TODO: Implement middleware

// jwtMiddleware validates JWT tokens
func jwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Extract token from Authorization header
		// - Expected format: "Bearer <token>"
		// - Validate token using validateJWT()
		// - Store user claims in context
		// - Return 401 if token is invalid or missing

		return c.Next()
	}
}

// adminMiddleware checks for admin role
func adminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Check if user has admin role
		// - Get user claims from context (set by jwtMiddleware)
		// - Return 403 if user is not admin

		return c.Next()
	}
}

// TODO: Implement route handlers

// registerHandler handles user registration
func registerHandler(c *fiber.Ctx) error {
	// TODO: Parse and validate registration request
	// - Validate input using struct validation
	// - Check if username/email already exists
	// - Validate password strength
	// - Hash password
	// - Create new user with role "user"
	// - Return success response (without password)
	return nil
}

// loginHandler handles user login
func loginHandler(c *fiber.Ctx) error {
	// TODO: Parse and validate login request
	// - Find user by username
	// - Verify password
	// - Check if user is active
	// - Generate JWT token
	// - Return auth response with token and user info
	return nil
}

// healthHandler returns API health status
func healthHandler(c *fiber.Ctx) error {
	// TODO: Return health check response
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now(),
	})
}

// getProfileHandler returns current user's profile
func getProfileHandler(c *fiber.Ctx) error {
	// TODO: Get user from JWT claims in context
	// - Find user by ID from claims
	// - Return user profile (without password)
	return nil
}

// updateProfileHandler updates current user's profile
func updateProfileHandler(c *fiber.Ctx) error {
	// TODO: Update user profile
	// - Get user ID from JWT claims
	// - Parse update request (email, username)
	// - Validate new data
	// - Check for conflicts (username/email uniqueness)
	// - Update user data
	// - Return updated profile
	return nil
}

// refreshTokenHandler generates a new JWT token
func refreshTokenHandler(c *fiber.Ctx) error {
	// TODO: Generate new token for current user
	// - Get user from JWT claims
	// - Generate new token with extended expiry
	// - Return new token
	return nil
}

// listUsersHandler returns all users (admin only)
func listUsersHandler(c *fiber.Ctx) error {
	// TODO: Return list of all users
	// - Remove password field from response
	// - Add pagination if needed
	// - Only accessible by admin users
	return nil
}

// updateUserRoleHandler updates a user's role (admin only)
func updateUserRoleHandler(c *fiber.Ctx) error {
	// TODO: Update user role
	// - Get user ID from URL parameter
	// - Parse new role from request body
	// - Validate role (must be "user" or "admin")
	// - Find and update user
	// - Return updated user
	return nil
}

// Helper functions

// findUserByUsername finds a user by username
func findUserByUsername(username string) (*User, int) {
	for i, user := range users {
		if user.Username == username {
			return &user, i
		}
	}
	return nil, -1
}

// findUserByID finds a user by ID
func findUserByID(id int) (*User, int) {
	for i, user := range users {
		if user.ID == id {
			return &user, i
		}
	}
	return nil, -1
}

// findUserByEmail finds a user by email
func findUserByEmail(email string) (*User, int) {
	for i, user := range users {
		if user.Email == email {
			return &user, i
		}
	}
	return nil, -1
}
