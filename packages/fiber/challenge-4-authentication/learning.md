# Learning: Authentication & Session Management

## 🌟 **What is Authentication?**

Authentication is the process of verifying the identity of a user or system. It answers the question "Who are you?" and is fundamental to securing web applications.

### **Authentication vs Authorization**
- **Authentication**: Verifying identity ("Who are you?")
- **Authorization**: Determining permissions ("What can you do?")

## 🔐 **Password Security**

### **Password Hashing with bcrypt**
Never store plain text passwords. Use bcrypt for secure hashing:

```go
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
    // Cost of 12 provides good security vs performance balance
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}

func verifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### **Password Strength Requirements**
Implement strong password policies:

```go
func validatePassword(password string) []string {
    var errors []string
    
    if len(password) < 8 {
        errors = append(errors, "Password must be at least 8 characters")
    }
    
    if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        errors = append(errors, "Password must contain uppercase letter")
    }
    
    if !regexp.MustCompile(`[a-z]`).MatchString(password) {
        errors = append(errors, "Password must contain lowercase letter")
    }
    
    if !regexp.MustCompile(`\d`).MatchString(password) {
        errors = append(errors, "Password must contain a digit")
    }
    
    if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
        errors = append(errors, "Password must contain special character")
    }
    
    return errors
}
```

## 🎫 **JWT (JSON Web Tokens)**

JWT is a stateless authentication method that encodes user information in a token.

### **JWT Structure**
A JWT consists of three parts separated by dots:
- **Header**: Algorithm and token type
- **Payload**: Claims (user data)
- **Signature**: Verification signature

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

### **Creating JWT Tokens**
```go
import "github.com/golang-jwt/jwt/v5"

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func generateJWT(user User, secret []byte) (string, error) {
    claims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "your-app-name",
            Subject:   strconv.Itoa(user.ID),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}
```

### **Validating JWT Tokens**
```go
func validateJWT(tokenString string, secret []byte) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    return claims, nil
}
```

## 🛡️ **Middleware for Authentication**

### **JWT Authentication Middleware**
```go
func jwtMiddleware(secret []byte) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract token from Authorization header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Missing authorization header",
            })
        }
        
        // Parse "Bearer <token>" format
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid authorization header format",
            })
        }
        
        // Validate token
        claims, err := validateJWT(parts[1], secret)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }
        
        // Store user info in context
        c.Locals("user_id", claims.UserID)
        c.Locals("username", claims.Username)
        c.Locals("role", claims.Role)
        
        return c.Next()
    }
}
```

### **Role-Based Access Control**
```go
func requireRole(requiredRole string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("role").(string)
        
        if userRole != requiredRole {
            return c.Status(403).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }
        
        return c.Next()
    }
}

// Usage
admin := app.Group("/admin", jwtMiddleware(secret), requireRole("admin"))
```

## 🔄 **Session Management**

### **Token Refresh Pattern**
Implement token refresh for better security:

```go
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

func generateTokenPair(user User) (*TokenPair, error) {
    // Short-lived access token (15 minutes)
    accessToken, err := generateJWT(user, 15*time.Minute)
    if err != nil {
        return nil, err
    }
    
    // Long-lived refresh token (7 days)
    refreshToken, err := generateRefreshToken(user, 7*24*time.Hour)
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}
```

### **Token Blacklisting**
Maintain a blacklist of revoked tokens:

```go
type TokenBlacklist struct {
    mu     sync.RWMutex
    tokens map[string]time.Time
}

func (tb *TokenBlacklist) Add(tokenID string, expiry time.Time) {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    tb.tokens[tokenID] = expiry
}

func (tb *TokenBlacklist) IsBlacklisted(tokenID string) bool {
    tb.mu.RLock()
    defer tb.mu.RUnlock()
    
    expiry, exists := tb.tokens[tokenID]
    if !exists {
        return false
    }
    
    // Clean up expired entries
    if time.Now().After(expiry) {
        delete(tb.tokens, tokenID)
        return false
    }
    
    return true
}
```

## 👤 **User Management Patterns**

### **User Registration**
```go
func registerUser(req RegisterRequest) (*User, error) {
    // Validate input
    if err := validateRegistration(req); err != nil {
        return nil, err
    }
    
    // Check if user exists
    if userExists(req.Username, req.Email) {
        return nil, errors.New("user already exists")
    }
    
    // Hash password
    hashedPassword, err := hashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &User{
        ID:       generateUserID(),
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
        Role:     "user",
        Active:   true,
        CreatedAt: time.Now(),
    }
    
    // Save to database/storage
    if err := saveUser(user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### **User Login**
```go
func loginUser(req LoginRequest) (*AuthResponse, error) {
    // Find user
    user, err := findUserByUsername(req.Username)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    // Verify password
    if !verifyPassword(req.Password, user.Password) {
        return nil, errors.New("invalid credentials")
    }
    
    // Check if account is active
    if !user.Active {
        return nil, errors.New("account is disabled")
    }
    
    // Generate tokens
    tokenPair, err := generateTokenPair(*user)
    if err != nil {
        return nil, err
    }
    
    return &AuthResponse{
        User:         *user,
        AccessToken:  tokenPair.AccessToken,
        RefreshToken: tokenPair.RefreshToken,
    }, nil
}
```

## 🔒 **Security Best Practices**

### **1. Secure Token Storage**
- Store JWT secret in environment variables
- Use strong, randomly generated secrets
- Rotate secrets periodically

```go
func getJWTSecret() []byte {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET environment variable is required")
    }
    return []byte(secret)
}
```

### **2. Rate Limiting**
Prevent brute force attacks:

```go
func loginRateLimit() fiber.Handler {
    // Allow 5 login attempts per minute per IP
    return limiter.New(limiter.Config{
        Max:        5,
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
    })
}
```

### **3. Input Validation**
Always validate and sanitize input:

```go
func validateLoginRequest(req LoginRequest) error {
    if req.Username == "" {
        return errors.New("username is required")
    }
    
    if req.Password == "" {
        return errors.New("password is required")
    }
    
    // Sanitize username
    req.Username = strings.TrimSpace(req.Username)
    
    return nil
}
```

### **4. HTTPS Only**
Always use HTTPS in production:

```go
app := fiber.New(fiber.Config{
    // Force HTTPS
    EnableTrustedProxyCheck: true,
    TrustedProxies: []string{"127.0.0.1"},
})

// Add security headers
app.Use(func(c *fiber.Ctx) error {
    c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
    return c.Next()
})
```

## 🧪 **Testing Authentication**

### **Testing JWT Functions**
```go
func TestJWTGeneration(t *testing.T) {
    user := User{
        ID:       1,
        Username: "testuser",
        Role:     "user",
    }
    
    secret := []byte("test-secret")
    
    token, err := generateJWT(user, secret)
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    
    // Validate the token
    claims, err := validateJWT(token, secret)
    assert.NoError(t, err)
    assert.Equal(t, user.ID, claims.UserID)
    assert.Equal(t, user.Username, claims.Username)
}
```

### **Testing Protected Endpoints**
```go
func TestProtectedEndpoint(t *testing.T) {
    app := setupTestApp()
    
    // Test without token
    req := httptest.NewRequest("GET", "/profile", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, 401, resp.StatusCode)
    
    // Test with valid token
    token := generateTestToken(t)
    req = httptest.NewRequest("GET", "/profile", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    resp, _ = app.Test(req)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## 🎯 **Best Practices Summary**

1. **Never store plain text passwords**
2. **Use strong password requirements**
3. **Implement proper JWT validation**
4. **Use HTTPS in production**
5. **Implement rate limiting**
6. **Validate and sanitize all input**
7. **Use environment variables for secrets**
8. **Implement proper error handling**
9. **Test authentication thoroughly**
10. **Follow the principle of least privilege**

## 📚 **Next Steps**

After mastering authentication:
1. **OAuth2 Integration** - Third-party authentication
2. **Multi-Factor Authentication** - Enhanced security
3. **Session Storage** - Redis/database sessions
4. **Audit Logging** - Track authentication events