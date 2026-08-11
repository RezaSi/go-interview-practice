# Challenge 4: Authentication & Session Management

Build a secure **User Authentication API** with JWT tokens, password hashing, and role-based access control.

## Challenge Requirements

Implement authentication system with these endpoints:

### Public Endpoints
- `POST /auth/register` - User registration with validation
- `POST /auth/login` - User login with JWT token generation

### Protected Endpoints (Require JWT)
- `GET /users/profile` - Get current user profile
- `PUT /users/profile` - Update user profile
- `POST /auth/refresh` - Refresh JWT token
- `POST /users/logout` - Logout user (blacklist token)

### Admin Endpoints (Require admin role)
- `GET /users` - List all users (admin only)
- `PUT /users/:id/role` - Update user role (admin only)
- `DELETE /users/:id` - Delete user (admin only)

## Data Structures

```go
type User struct {
    ID        string `json:"id"`
    Username  string `json:"username" validate:"required,min=3,max=20"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"-"` // Never return in JSON
    Role      string `json:"role"` // "user" or "admin"
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

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
    Message string `json:"message,omitempty"`
    User    *User  `json:"user,omitempty"`
}

type UserListResponse struct {
    Status string `json:"status"`
    Count  int    `json:"count"`
    Users  []User `json:"users"`
}

type JWTClaims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

## Security Requirements

### Password Security
- Hash passwords using bcrypt before storing
- Never return passwords in API responses
- Minimum password length: 6 characters

### JWT Token Management
- Access tokens expire in 15 minutes
- Refresh tokens expire in 7 days
- Include user ID, username, and role in token claims
- Implement token blacklisting for logout

### Role-Based Access Control
- Default role: "user"
- Admin role: "admin"
- Protect admin endpoints with role validation

## Request/Response Examples

**POST /auth/register** (Request body)
```json
{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepass123"
}
```

**POST /auth/login** (Request body)
```json
{
    "username": "johndoe",
    "password": "securepass123"
}
```

**Authentication Response**
```json
{
    "status": "success",
    "message": "Login successful",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": "uuid-string",
        "username": "johndoe",
        "email": "john@example.com",
        "role": "user",
        "created_at": "2024-01-15T10:30:00Z"
    }
}
```

## Testing Requirements

Your solution must pass tests for:
- User registration with password hashing
- User login with JWT token generation
- Protected endpoints require valid JWT tokens
- Admin endpoints require admin role
- Token refresh functionality works
- User logout blacklists tokens
- Profile update functionality
- User role management (admin only)
- Proper HTTP status codes and error handling
- Input validation for all endpoints