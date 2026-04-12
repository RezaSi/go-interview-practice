# Learning GORM CRUD Operations

## Overview

GORM (Go Object Relational Mapper) is a powerful ORM library for Go that simplifies database operations. This challenge focuses on mastering the fundamental CRUD (Create, Read, Update, Delete) operations using GORM.

## What is GORM?

GORM is a feature-rich ORM library for Go that provides:
- **Auto Migration**: Automatically create database tables from structs
- **CRUD Operations**: Simple methods for database operations
- **Hooks**: Lifecycle callbacks (BeforeCreate, AfterUpdate, etc.)
- **Associations**: Handle relationships between models
- **Validation**: Built-in validation support
- **Transactions**: Database transaction support

## Basic Setup

### 1. Installation
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite  # For SQLite
```

### 2. Database Connection
```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Auto migrate the schema
    err = db.AutoMigrate(&User{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}
```

## Defining Models

### Basic Model Structure
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

### GORM Tags
- `gorm:"primaryKey"` - Marks field as primary key
- `gorm:"not null"` - Makes field required
- `gorm:"unique"` - Makes field unique
- `gorm:"check:condition"` - Adds check constraint
- `gorm:"default:value"` - Sets default value
- `gorm:"index"` - Creates database index

## CRUD Operations

### 1. Create (C)
```go
// Create a single user
user := User{Name: "John Doe", Email: "john@example.com", Age: 25}
result := db.Create(&user)
if result.Error != nil {
    return result.Error
}

// Create multiple users
users := []User{
    {Name: "User 1", Email: "user1@example.com", Age: 25},
    {Name: "User 2", Email: "user2@example.com", Age: 30},
}
result = db.Create(&users)
```

### 2. Read (R)
```go
// Get first user
var user User
result := db.First(&user, 1) // Find by primary key
if result.Error != nil {
    return result.Error
}

// Get user by condition
var user User
result = db.Where("email = ?", "john@example.com").First(&user)

// Get all users
var users []User
result = db.Find(&users)

// Get users with conditions
var users []User
result = db.Where("age > ?", 18).Find(&users)
```

### 3. Update (U)
```go
// Update by primary key
user := User{ID: 1, Name: "Updated Name", Email: "updated@example.com", Age: 30}
result := db.Save(&user)

// Update specific fields
result = db.Model(&user).Update("Name", "New Name")

// Update multiple fields
result = db.Model(&user).Updates(User{Name: "New Name", Age: 31})

// Update with conditions
result = db.Model(&User{}).Where("age < ?", 18).Update("age", 18)
```

### 4. Delete (D)
```go
// Delete by primary key
result := db.Delete(&User{}, 1)

// Delete with conditions
result = db.Where("age < ?", 18).Delete(&User{})

// Soft delete (if model has DeletedAt field)
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

## Error Handling

### Common Error Patterns
```go
// Check for errors
result := db.Create(&user)
if result.Error != nil {
    // Handle error
    return result.Error
}

// Check for "not found" errors
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    // Handle not found
    return fmt.Errorf("user not found")
}

// Check for unique constraint violations
if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
    // Handle duplicate entry
    return fmt.Errorf("email already exists")
}
```

## Validation

### Built-in Validation
```go
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"not null"`
    Email string `gorm:"unique;not null"`
    Age   int    `gorm:"check:age > 0"`
}
```

### Custom Validation
```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Age < 0 {
        return fmt.Errorf("age cannot be negative")
    }
    if !strings.Contains(u.Email, "@") {
        return fmt.Errorf("invalid email format")
    }
    return nil
}
```

## Query Methods

### Where Clauses
```go
// Basic where
db.Where("name = ?", "John").Find(&users)

// Multiple conditions
db.Where("name = ? AND age > ?", "John", 18).Find(&users)

// IN clause
db.Where("name IN ?", []string{"John", "Jane"}).Find(&users)

// LIKE clause
db.Where("name LIKE ?", "%John%").Find(&users)
```

### Ordering and Limiting
```go
// Order by
db.Order("age DESC").Find(&users)

// Limit and offset
db.Limit(10).Offset(20).Find(&users)

// Select specific fields
db.Select("name, email").Find(&users)
```

## Transactions

### Basic Transaction
```go
func CreateUserWithProfile(db *gorm.DB, user *User, profile *Profile) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Create user
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        // Create profile
        profile.UserID = user.ID
        if err := tx.Create(profile).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### Manual Transaction
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Commit().Error; err != nil {
    return err
}
```

## Best Practices

### 1. Use Pointers for Models
```go
// Good
func GetUser(db *gorm.DB, id uint) (*User, error) {
    var user User
    result := db.First(&user, id)
    return &user, result.Error
}

// Bad
func GetUser(db *gorm.DB, id uint) (User, error) {
    var user User
    result := db.First(&user, id)
    return user, result.Error
}
```

### 2. Handle Errors Properly
```go
// Always check for errors
result := db.Create(&user)
if result.Error != nil {
    return result.Error
}
```

### 3. Use Appropriate Query Methods
```go
// Use First() for single records
var user User
db.First(&user, id)

// Use Find() for multiple records
var users []User
db.Find(&users)

// Use Take() when order doesn't matter
var user User
db.Take(&user)
```

### 4. Optimize Queries
```go
// Select only needed fields
db.Select("id, name").Find(&users)

// Use preloading for relationships
db.Preload("Posts").Find(&users)

// Use transactions for multiple operations
db.Transaction(func(tx *gorm.DB) error {
    // Multiple operations
    return nil
})
```

## Common Patterns

### CRUD Service Pattern
```go
type UserService struct {
    db *gorm.DB
}

func (s *UserService) Create(user *User) error {
    return s.db.Create(user).Error
}

func (s *UserService) GetByID(id uint) (*User, error) {
    var user User
    err := s.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *UserService) Update(user *User) error {
    return s.db.Save(user).Error
}

func (s *UserService) Delete(id uint) error {
    return s.db.Delete(&User{}, id).Error
}
```

### Repository Pattern
```go
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    GetAll() ([]User, error)
    Update(user *User) error
    Delete(id uint) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}
```

## Resources

- [GORM Documentation](https://gorm.io/docs/)
- [GORM CRUD Operations](https://gorm.io/docs/create.html)
- [GORM Query Interface](https://gorm.io/docs/query.html)
- [GORM Hooks](https://gorm.io/docs/hooks.html)
- [GORM Transactions](https://gorm.io/docs/transactions.html)

## Practice Exercises

1. Create a simple user management system
2. Implement CRUD operations for a blog post system
3. Build a product catalog with categories
4. Create a task management application
5. Implement a simple inventory system

These exercises will help you master GORM CRUD operations and understand database interactions in Go applications. 