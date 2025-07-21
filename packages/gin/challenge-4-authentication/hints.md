# Hints for Challenge 4: Authentication & Session Management

## Hint 1: Password Strength Validation

Implement comprehensive password strength checking:

```go
func isStrongPassword(password string) bool {
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
```

## Hint 2: Password Hashing with Bcrypt

Use bcrypt for secure password hashing:

```go
func hashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func verifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

## Hint 3: JWT Token Generation

Generate secure JWT tokens with proper claims:

```go
func generateTokens(userID int, username, role string) (*TokenResponse, error) {
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
```

## Hint 4: JWT Token Validation

Validate and parse JWT tokens:

```go
func validateToken(tokenString string) (*JWTClaims, error) {
    // Check if token is blacklisted
    if blacklistedTokens[tokenString] {
        return nil, errors.New("token is blacklisted")
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
```

## Hint 5: Authentication Middleware

Create middleware to protect routes:

```go
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
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := validateToken(tokenString)
        if err != nil {
            c.JSON(401, APIResponse{
                Success: false,
                Error:   "Invalid token",
            })
            c.Abort()
            return
        }
        
        // Set user info in context
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

## Hint 6: Role-Based Authorization

Implement role-based access control:

```go
func requireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(401, APIResponse{
                Success: false,
                Error:   "Unauthorized",
            })
            c.Abort()
            return
        }
        
        roleStr := userRole.(string)
        for _, allowedRole := range roles {
            if roleStr == allowedRole {
                c.Next()
                return
            }
        }
        
        c.JSON(403, APIResponse{
            Success: false,
            Error:   "Insufficient permissions",
        })
        c.Abort()
    }
}
```

## Hint 7: Account Lockout Management

Handle failed login attempts and account lockout:

```go
func recordFailedAttempt(user *User) {
    user.FailedAttempts++
    if user.FailedAttempts >= maxFailedAttempts {
        lockUntil := time.Now().Add(lockoutDuration)
        user.LockedUntil = &lockUntil
    }
}

func isAccountLocked(user *User) bool {
    if user.LockedUntil == nil {
        return false
    }
    return time.Now().Before(*user.LockedUntil)
}

func resetFailedAttempts(user *User) {
    user.FailedAttempts = 0
    user.LockedUntil = nil
}
```

## Hint 8: Secure Token Logout

Implement secure logout with token blacklisting:

```go
func logout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(401, APIResponse{
            Success: false,
            Error:   "Authorization header required",
        })
        return
    }
    
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    
    // Add token to blacklist
    blacklistedTokens[tokenString] = true
    
    // Remove refresh token if provided
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
``` 