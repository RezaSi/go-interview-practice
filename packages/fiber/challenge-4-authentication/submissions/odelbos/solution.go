package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	// Exp      int64  `json:"exp"`
	jwt.RegisteredClaims
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

// --------------------------------------------------------------------
// Main
// --------------------------------------------------------------------

func main() {
	app := fiber.New()

	// Setup validator with custom password validator
	setupCustomValidator()

	// Public routes
	app.Post("/auth/register", registerHandler)
	app.Post("/auth/login", loginHandler)
	app.Get("/health", healthHandler)

	// Protected routes
	protected := app.Group("/", jwtMiddleware())
	protected.Get("/profile", getProfileHandler)
	protected.Put("/profile", updateProfileHandler)
	protected.Post("/auth/refresh", refreshTokenHandler)

	// Admin routes
	admin := app.Group("/admin", jwtMiddleware(), adminMiddleware())
	admin.Get("/users", listUsersHandler)
	admin.Put("/users/:id/role", updateUserRoleHandler)

	app.Listen(":3000")
}

func setupCustomValidator() {
	validate = validator.New()
	err := validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return validatePassword(fl.Field().String())
	})
	if err != nil {
		panic(fmt.Errorf("failed to register validator"))
	}
}

// --------------------------------------------------------------------
// Password security
// --------------------------------------------------------------------

func validatePassword(password string) bool {
	// Minimum 8 characters
	// At least one uppercase letter
	// At least one lowercase letter
	// At least one digit
	// At least one special character
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

// --------------------------------------------------------------------
// JWT functions
// --------------------------------------------------------------------

func generateJWT(user User) (string, error) {
	// Expiration: 1 hour from now
	// Sign with HS256 algorithm
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token")
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

// --------------------------------------------------------------------
// Middlewares
// --------------------------------------------------------------------

func jwtMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		bearer := c.Get("Authorization")
		if !strings.HasPrefix(bearer, "Bearer ") {
			return errResponse(c, fiber.StatusUnauthorized, "Invalid token")
		}
		tokenStr := strings.TrimPrefix(bearer, "Bearer ")
		claims, err := validateJWT(tokenStr)
		if err != nil {
			return errResponse(c, fiber.StatusUnauthorized, "Invalid token")
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}

func adminMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		claims := c.Locals("claims").(*JWTClaims)
		if claims.Role != "admin" {
			return errResponse(c, fiber.StatusForbidden, "Forbidden")
		}
		return c.Next()
	}
}

// --------------------------------------------------------------------
// Implement route handlers
// --------------------------------------------------------------------

func registerHandler(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, "Invalid request")
	}
	if err := validate.Struct(req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, err.Error())
	}
	if _, index := findUserByUsername(req.Username); index != -1 {
		return errResponse(c, fiber.StatusConflict, "Already exists")
	}
	if _, index := findUserByEmail(req.Email); index != -1 {
		return errResponse(c, fiber.StatusConflict, "Already exists")
	}

	pwdHash, err := hashPassword(req.Password)
	if err != nil {
		return errResponse(c, fiber.StatusInternalServerError, "Internal error")
	}

	user := User{
		ID:       nextUserID,
		Username: req.Username,
		Email:    req.Email,
		Password: pwdHash,
		Role:     "user",
		Active:   true,
	}
	users = append(users, user)
	nextUserID++

	token, err := generateJWT(user)
	if err != nil {
		return errResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}
	return okResponse(c, fiber.StatusCreated, token, safeUser(user), "")
}

func loginHandler(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, err.Error())
	}

	user, _ := findUserByUsername(req.Username)
	if user == nil || !verifyPassword(req.Password, user.Password) {
		return errResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}
	if !user.Active {
		return errResponse(c, fiber.StatusForbidden, "Account inactive")
	}

	token, err := generateJWT(*user)
	if err != nil {
		return errResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}
	return okResponse(c, fiber.StatusOK, token, safeUser(*user), "Login success")
}

func healthHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "timestamp": time.Now()})
}

func getProfileHandler(c fiber.Ctx) error {
	// Find user by ID from claims
	// Return user profile (without password)
	claims := c.Locals("claims").(*JWTClaims)
	// NOTE: Should use mutex
	user, _ := findUserByID(claims.UserID)
	if user == nil {
		return errResponse(c, fiber.StatusNotFound, "Not found")
	}
	return okResponse(c, fiber.StatusOK, "", safeUser(*user), "")
}

func updateProfileHandler(c fiber.Ctx) error {
	claims := c.Locals("claims").(*JWTClaims)
	user, idx := findUserByID(claims.UserID)
	if idx == -1 {
		return errResponse(c, fiber.StatusNotFound, "Not found")
	}

	type UpdateData struct {
		Username string `json:"username" validate:"required,min=3,max=20"`
		Email    string `json:"email" validate:"required,email"`
	}
	var req UpdateData

	if err := c.Bind().Body(&req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// NOTE: Should use mutex
	if u, _ := findUserByUsername(req.Username); u != nil && u.ID != user.ID {
		return errResponse(c, fiber.StatusConflict, "Already exists")
	}
	if u, _ := findUserByEmail(req.Email); u != nil && u.ID != user.ID {
		return errResponse(c, fiber.StatusConflict, "Already exists")
	}

	user.Username = req.Username
	user.Email = req.Email
	users[idx] = *user
	return c.JSON(safeUser(*user))
}

func refreshTokenHandler(c fiber.Ctx) error {
	claims := c.Locals("claims").(*JWTClaims)
	user, _ := findUserByID(claims.UserID)
	if user == nil {
		return errResponse(c, fiber.StatusNotFound, "Not found")
	}

	token, err := generateJWT(*user)
	if err != nil {
		return errResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}
	return okResponse(c, fiber.StatusOK, token, safeUser(*user), "")
}

func listUsersHandler(c fiber.Ctx) error {
	results := make([]User, len(users))
	for _, user := range users {
		results = append(results, safeUser(user))
	}
	return c.JSON(results)
}

func updateUserRoleHandler(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id == 0 {
		return errResponse(c, fiber.StatusBadRequest, "Invalid ID")
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, "Invalid request")
	}
	if req.Role != "user" && req.Role != "admin" {
		return errResponse(c, fiber.StatusBadRequest, "Invalid role")
	}

	user, _ := findUserByID(id)
	if user == nil {
		return errResponse(c, fiber.StatusNotFound, "Not Found")
	}
	user.Role = req.Role
	return c.JSON(safeUser(*user))
}

// --------------------------------------------------------------------
// Helper functions
// --------------------------------------------------------------------

func findUserByUsername(username string) (*User, int) {
	for i, user := range users {
		if user.Username == username {
			return &user, i
		}
	}
	return nil, -1
}

func findUserByID(id int) (*User, int) {
	for i, user := range users {
		if user.ID == id {
			return &user, i
		}
	}
	return nil, -1
}

func findUserByEmail(email string) (*User, int) {
	for i, user := range users {
		if user.Email == email {
			return &user, i
		}
	}
	return nil, -1
}

func safeUser(user User) User {
	user.Password = ""
	return user
}

func okResponse(c fiber.Ctx, status int, token string, user User, msg string) error {
	return c.Status(status).JSON(AuthResponse{
		Success: true,
		Token:   token,
		User:    user,
		Message: msg,
	})
}

func errResponse(c fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(AuthResponse{
		Success: false,
		Message: msg,
	})
}
