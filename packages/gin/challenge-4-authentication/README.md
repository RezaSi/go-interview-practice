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
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
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
- Hash with bcrypt (cost factor 12)

### JWT Implementation
- Use HS256 algorithm
- Include user ID and role in claims
- 24-hour expiration
- Proper token verification middleware

### Role-Based Access
- Users: Access to profile endpoints only
- Admins: Access to all endpoints including user management

## Testing Requirements

Your solution must pass tests for:
- User registration with password validation
- Password hashing (bcrypt verification)
- User login with correct credentials
- JWT token generation and validation
- Protected endpoint access control
- Role-based authorization (admin vs user)
- Token refresh functionality
- Proper error handling for invalid credentials 