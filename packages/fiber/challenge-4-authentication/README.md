# Challenge 4: Authentication & Session Management

Build a secure **User Authentication API** with JWT tokens, password hashing, and role-based access control.

## Challenge Requirements

Implement authentication system with these endpoints:

### Public Endpoints
- `POST /auth/register` - User registration with validation
- `POST /auth/login` - User login with JWT token generation
- `GET /health` - Public health check

### Protected Endpoints (Require JWT)
- `GET /profile` - Get current user profile
- `PUT /profile` - Update user profile
- `POST /auth/refresh` - Refresh JWT token

### Admin Endpoints (Require admin role)
- `GET /admin/users` - List all users (admin only)
- `PUT /admin/users/:id/role` - Update user role (admin only)

## Data Structures

```go
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"-"` // Never return in JSON
    Role     string `json:"role"` // "user" or "admin"
    Active   bool   `json:"active"`
}

type AuthResponse struct {
    Success bool   `json:"success"`
    Token   string `json:"token,omitempty"`
    User    User   `json:"user,omitempty"`
    Message string `json:"message,omitempty"`
}
```

## Security Requirements

### Password Security
- Minimum 8 characters
- Must contain uppercase, lowercase, number, and special character
- Hash passwords using bcrypt with cost 12

### JWT Token Security
- Use HS256 algorithm
- Include user ID, username, and role in claims
- Token expiry: 1 hour
- Refresh token expiry: 7 days

### Role-Based Access
- **user**: Can access own profile and update it
- **admin**: Can access all user data and modify user roles

## API Examples

**POST /auth/register**
```json
{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "SecureP@ss123"
}
```

**POST /auth/login**
```json
{
    "username": "john_doe",
    "password": "SecureP@ss123"
}
```

**Response:**
```json
{
    "success": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "role": "user",
        "active": true
    }
}
```

## Testing Requirements

Your solution must handle:
- User registration with password validation
- Password hashing and verification
- JWT token generation and validation
- Role-based access control
- Token expiration handling
- Secure password storage (never return in responses)
- User authentication and authorization
- Profile updates with validation