# Challenge 1: CRUD Operations

Build a **User Management System** using GORM that demonstrates fundamental database operations.

## Challenge Requirements

Create a Go application that supports:

1. **Create** - Add new users to database
2. **Read** - Query and retrieve users  
3. **Update** - Modify existing user data
4. **Delete** - Remove users from database
5. **Database Connection** - Connect to SQLite database

## Data Model

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    Age       int       `gorm:"check:age > 0"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Required Functions

Implement these functions:
- `ConnectDB() (*gorm.DB, error)` - Database connection
- `CreateUser(db *gorm.DB, user *User) error` - Create user
- `GetUserByID(db *gorm.DB, id uint) (*User, error)` - Get user
- `GetAllUsers(db *gorm.DB) ([]User, error)` - Get all users
- `UpdateUser(db *gorm.DB, user *User) error` - Update user
- `DeleteUser(db *gorm.DB, id uint) error` - Delete user

## Testing Requirements

Your solution must pass tests for:
- Database connection and table creation
- Creating users with validation
- Querying users by ID and retrieving all users
- Updating user information
- Deleting users from database
- Error handling for invalid operations 