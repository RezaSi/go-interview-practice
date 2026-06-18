# Hints for Challenge 4: Authentication & Security

## Hint 1: JWT Setup

Set up JWT with the golang-jwt package:

```go
import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key") // Use environment variable in production
```

## Hint 2: Password Hashing

Use bcrypt for secure password hashing:

```go
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

func verifyPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
```

## Hint 3: JWT Token Generation

Generate access and refresh tokens:

```go
func generateTokens(user *User) (string, string, error) {
    // Access token (15 minutes)
    accessClaims := JWTClaims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "echo-auth-api",
            Subject:   user.ID,
        },
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString(jwtSecret)
    if err != nil {
        return "", "", err
    }
    
    // Refresh token (7 days)
    refreshClaims := JWTClaims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "echo-auth-api",
            Subject:   user.ID,
        },
    }
    
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString(jwtSecret)
    
    return accessTokenString, refreshTokenString, err
}
```

## Hint 4: JWT Token Validation

Validate and parse JWT tokens:

```go
func validateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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
```

## Hint 5: JWT Middleware

Create middleware to protect routes:

```go
func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return c.JSON(http.StatusUnauthorized, ErrorResponse{
                Status:  "error",
                Message: "Authorization header required",
            })
        }
        
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            return c.JSON(http.StatusUnauthorized, ErrorResponse{
                Status:  "error",
                Message: "Invalid authorization header format",
            })
        }
        
        tokenString := tokenParts[1]
        
        // Check if token is blacklisted
        if blacklistedTokens[tokenString] {
            return c.JSON(http.StatusUnauthorized, ErrorResponse{
                Status:  "error",
                Message: "Token has been revoked",
            })
        }
        
        claims, err := validateToken(tokenString)
        if err != nil {
            return c.JSON(http.StatusUnauthorized, ErrorResponse{
                Status:  "error",
                Message: "Invalid or expired token",
            })
        }
        
        // Set user info in context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Set("token", tokenString)
        
        return next(c)
    }
}
```

## Hint 6: Admin Role Middleware

Create middleware to check admin role:

```go
func adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        userRole := c.Get("role")
        if userRole == nil {
            return c.JSON(http.StatusUnauthorized, ErrorResponse{
                Status:  "error",
                Message: "Authentication required",
            })
        }
        
        if userRole.(string) != "admin" {
            return c.JSON(http.StatusForbidden, ErrorResponse{
                Status:  "error",
                Message: "Admin access required",
            })
        }
        
        return next(c)
    }
}
```

## Hint 7: User Registration

Implement user registration with validation:

```go
func registerHandler(c echo.Context) error {
    var req RegisterRequest
    
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  "error",
            Message: "Invalid request format",
        })
    }
    
    // Basic validation
    if strings.TrimSpace(req.Username) == "" {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  "error",
            Message: "Username is required",
        })
    }
    
    if len(req.Password) < 6 {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  "error",
            Message: "Password must be at least 6 characters",
        })
    }
    
    // Check if username already exists
    for _, user := range users {
        if user.Username == req.Username {
            return c.JSON(http.StatusConflict, ErrorResponse{
                Status:  "error",
                Message: "Username already exists",
            })
        }
    }
    
    // Hash password and create user
    hashedPassword, err := hashPassword(req.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:  "error",
            Message: "Failed to process password",
        })
    }
    
    user := User{
        ID:        uuid.New().String(),
        Username:  req.Username,
        Email:     req.Email,
        Password:  hashedPassword,
        Role:      "user",
        CreatedAt: time.Now().Format(time.RFC3339),
        UpdatedAt: time.Now().Format(time.RFC3339),
    }
    
    users = append(users, user)
    
    // Generate tokens
    accessToken, refreshToken, err := generateTokens(&user)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:  "error",
            Message: "Failed to generate tokens",
        })
    }
    
    user.Password = "" // Remove password from response
    
    return c.JSON(http.StatusCreated, AuthResponse{
        Status:       "success",
        Message:      "User registered successfully",
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         &user,
    })
}
```

## Hint 8: User Login

Implement login with credential verification:

```go
func loginHandler(c echo.Context) error {
    var req LoginRequest
    
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  "error",
            Message: "Invalid request format",
        })
    }
    
    // Find user by username
    var foundUser *User
    for i, user := range users {
        if user.Username == req.Username {
            foundUser = &users[i]
            break
        }
    }
    
    if foundUser == nil {
        return c.JSON(http.StatusUnauthorized, ErrorResponse{
            Status:  "error",
            Message: "Invalid username or password",
        })
    }
    
    // Verify password
    if !verifyPassword(foundUser.Password, req.Password) {
        return c.JSON(http.StatusUnauthorized, ErrorResponse{
            Status:  "error",
            Message: "Invalid username or password",
        })
    }
    
    // Generate tokens
    accessToken, refreshToken, err := generateTokens(foundUser)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:  "error",
            Message: "Failed to generate tokens",
        })
    }
    
    foundUser.Password = "" // Remove password from response
    
    return c.JSON(http.StatusOK, AuthResponse{
        Status:       "success",
        Message:      "Login successful",
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         foundUser,
    })
}
```

## Hint 9: Token Blacklisting

Implement logout with token blacklisting:

```go
var blacklistedTokens = make(map[string]bool)

func logoutHandler(c echo.Context) error {
    token := c.Get("token").(string)
    
    // Add token to blacklist
    blacklistedTokens[token] = true
    
    return c.JSON(http.StatusOK, MessageResponse{
        Status:  "success",
        Message: "Logged out successfully",
    })
}

func isTokenBlacklisted(token string) bool {
    return blacklistedTokens[token]
}
```

## Hint 10: Route Setup with Middleware

Set up protected routes with proper middleware:

```go
func setupRoutes(e *echo.Echo) {
    // Public authentication routes
    auth := e.Group("/auth")
    auth.POST("/register", registerHandler)
    auth.POST("/login", loginHandler)
    auth.POST("/refresh", refreshHandler)
    
    // Protected routes (require authentication)
    protected := e.Group("/users")
    protected.Use(jwtMiddleware)
    protected.GET("/profile", getProfileHandler)
    protected.PUT("/profile", updateProfileHandler)
    protected.POST("/logout", logoutHandler)
    
    // Admin routes (require admin role)
    admin := e.Group("/users")
    admin.Use(jwtMiddleware)
    admin.Use(adminMiddleware)
    admin.GET("", getAllUsersHandler)
    admin.PUT("/:id/role", updateUserRoleHandler)
    admin.DELETE("/:id", deleteUserHandler)
}
```