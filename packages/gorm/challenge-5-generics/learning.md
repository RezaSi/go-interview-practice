# Learning Guide: GORM Generics API

## Introduction

GORM v1.30.0+ introduced a revolutionary Generics API that provides type safety, better performance, and enhanced developer experience. This guide will help you understand and master the new generics-based approach to database operations.

## Why Generics?

### Problems with Traditional API
```go
// Traditional API - potential issues
var user User
db.Where("name = ?", name).First(&user) // No compile-time type checking
// Prone to SQL pollution from reused db instances
// Less IDE support
```

### Benefits of Generics API
```go
// Generics API - improved
user, err := gorm.G[User](db).Where("name = ?", name).First(ctx)
// ✅ Compile-time type safety
// ✅ Better performance (reduced SQL pollution)
// ✅ Enhanced IDE support
// ✅ Required context usage
```

## Core Concepts

### 1. Type-Safe Operations

The `gorm.G[T]` function creates a type-safe database instance:

```go
// Type-safe user operations
userDB := gorm.G[User](db)
user, err := userDB.Where("age > ?", 18).First(ctx)

// Type-safe company operations  
companyDB := gorm.G[Company](db)
companies, err := companyDB.Find(ctx)
```

### 2. Context-First Design

All generics operations require a context:

```go
ctx := context.Background()

// Basic operations
user, err := gorm.G[User](db).Create(ctx, &User{Name: "John"})
users, err := gorm.G[User](db).Find(ctx)
user, err := gorm.G[User](db).First(ctx)
err := gorm.G[User](db).Delete(ctx)
```

### 3. Enhanced Configuration

Pass configuration options as additional parameters:

```go
// OnConflict handling
err := gorm.G[User](db, clause.OnConflict{DoNothing: true}).Create(ctx, &user)

// Execution hints
users, err := gorm.G[User](db, 
    hints.New("USE_INDEX(users, idx_name)"),
).Find(ctx)

// Result metadata
result := gorm.WithResult()
err := gorm.G[User](db, result).Create(ctx, &user)
fmt.Printf("Rows affected: %d", result.RowsAffected)
```

## CRUD Operations

### Create Operations

```go
// Single create
user := &User{Name: "Alice", Email: "alice@example.com", Age: 25}
err := gorm.G[User](db).Create(ctx, user)

// Batch create
users := []User{
    {Name: "Bob", Email: "bob@example.com", Age: 30},
    {Name: "Charlie", Email: "charlie@example.com", Age: 35},
}
err := gorm.G[User](db).CreateInBatches(ctx, users, 10)
```

### Read Operations

```go
// Find by condition
users, err := gorm.G[User](db).Where("age >= ?", 18).Find(ctx)

// First record
user, err := gorm.G[User](db).Where("email = ?", email).First(ctx)

// Count records
count, err := gorm.G[User](db).Where("age >= ?", 18).Count(ctx)
```

### Update Operations

```go
// Update single field
err := gorm.G[User](db).Where("id = ?", userID).Update(ctx, "age", 26)

// Update multiple fields
err := gorm.G[User](db).Where("id = ?", userID).Updates(ctx, User{
    Name: "Updated Name",
    Age:  26,
})
```

### Delete Operations

```go
// Delete by condition
err := gorm.G[User](db).Where("age < ?", 18).Delete(ctx)

// Delete by ID
err := gorm.G[User](db).Where("id = ?", userID).Delete(ctx)
```

## Advanced Features

### OnConflict Handling

Handle duplicate key conflicts gracefully:

```go
// Do nothing on conflict
err := gorm.G[User](db, clause.OnConflict{DoNothing: true}).Create(ctx, &user)

// Update on conflict
err := gorm.G[User](db, clause.OnConflict{
    Columns:   []clause.Column{{Name: "email"}},
    DoUpdates: clause.AssignmentColumns([]string{"name", "age", "updated_at"}),
}).Create(ctx, &user)

// Custom conflict resolution
err := gorm.G[User](db, clause.OnConflict{
    Columns: []clause.Column{{Name: "email"}},
    DoUpdates: clause.Assignments(map[string]interface{}{
        "login_count": gorm.Expr("login_count + 1"),
        "last_login":  time.Now(),
    }),
}).Create(ctx, &user)
```

### Enhanced Joins

More flexible and powerful join operations:

```go
// Basic association join
users, err := gorm.G[User](db).Joins(clause.InnerJoin.Association("Company"), nil).Find(ctx)

// Join with custom conditions
users, err := gorm.G[User](db).Joins(clause.LeftJoin.Association("Company"), 
    func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
        db.Where("companies.industry = ?", "Technology")
        db.Where("companies.founded_year > ?", 2000)
        return nil
    }).Find(ctx)

// Subquery joins
users, err := gorm.G[User](db).Joins(
    clause.LeftJoin.AssociationFrom("Company", 
        gorm.G[Company](db).Select("id", "name").Where("active = ?", true)
    ).As("active_companies"),
    func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
        db.Where("?.industry = ?", joinTable, "Technology")
        return nil
    },
).Find(ctx)
```

### Enhanced Preloading

Improved association loading with more control:

```go
// Basic preload
users, err := gorm.G[User](db).Preload("Posts", nil).Find(ctx)

// Preload with conditions
users, err := gorm.G[User](db).Preload("Posts", func(db gorm.PreloadBuilder) error {
    db.Where("published = ?", true)
    db.Order("created_at DESC")
    return nil
}).Find(ctx)

// Limit per record
users, err := gorm.G[User](db).Preload("Posts", func(db gorm.PreloadBuilder) error {
    db.Order("created_at DESC").LimitPerRecord(5)
    return nil
}).Find(ctx)

// Nested preloading
users, err := gorm.G[User](db).
    Preload("Posts", func(db gorm.PreloadBuilder) error {
        db.Where("published = ?", true)
        return nil
    }).
    Preload("Posts.Comments", func(db gorm.PreloadBuilder) error {
        db.Where("approved = ?", true)
        db.LimitPerRecord(3)
        return nil
    }).Find(ctx)
```

## Performance Optimizations

### Batch Operations

```go
// Efficient batch inserts
users := make([]User, 1000)
// ... populate users
err := gorm.G[User](db).CreateInBatches(ctx, users, 100)

// Batch updates
err := gorm.G[User](db).Where("department = ?", "Engineering").Updates(ctx, map[string]interface{}{
    "bonus": gorm.Expr("salary * 0.1"),
})
```

### Query Optimization

```go
// Select specific fields
users, err := gorm.G[User](db).Select("id", "name", "email").Find(ctx)

// Use indexes effectively
users, err := gorm.G[User](db, 
    hints.New("USE_INDEX(users, idx_email)"),
).Where("email LIKE ?", "%@company.com").Find(ctx)
```

### Connection Pool Management

```go
// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

users, err := gorm.G[User](db).Find(ctx)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Handle timeout
    }
}
```

## Migration from Traditional API

### Before (Traditional API)
```go
func GetActiveUsers(db *gorm.DB) ([]User, error) {
    var users []User
    err := db.Where("active = ?", true).Find(&users).Error
    return users, err
}

func CreateUser(db *gorm.DB, user *User) error {
    return db.Create(user).Error
}
```

### After (Generics API)
```go
func GetActiveUsers(ctx context.Context, db *gorm.DB) ([]User, error) {
    return gorm.G[User](db).Where("active = ?", true).Find(ctx)
}

func CreateUser(ctx context.Context, db *gorm.DB, user *User) error {
    return gorm.G[User](db).Create(ctx, user)
}
```

## Error Handling

### Context-Aware Error Handling

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := gorm.G[User](db).Where("id = ?", id).First(ctx)
if err != nil {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        return nil, fmt.Errorf("query timeout: %w", err)
    case errors.Is(err, context.Canceled):
        return nil, fmt.Errorf("query canceled: %w", err)
    case errors.Is(err, gorm.ErrRecordNotFound):
        return nil, fmt.Errorf("user not found: %w", err)
    default:
        return nil, fmt.Errorf("database error: %w", err)
    }
}
```

## Best Practices

### 1. Always Use Context
```go
// ✅ Good - Always pass context
ctx := context.Background()
users, err := gorm.G[User](db).Find(ctx)

// ❌ Bad - Traditional API still works but misses benefits
var users []User
db.Find(&users)
```

### 2. Leverage Type Safety
```go
// ✅ Good - Type-safe operations
user, err := gorm.G[User](db).Where("age > ?", 18).First(ctx)

// ✅ Good - Compile-time type checking
users, err := gorm.G[User](db).Where("company_id = ?", companyID).Find(ctx)
```

### 3. Use Enhanced Features
```go
// ✅ Good - Use OnConflict for upserts
err := gorm.G[User](db, clause.OnConflict{
    Columns:   []clause.Column{{Name: "email"}},
    DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
}).Create(ctx, &user)

// ✅ Good - Use LimitPerRecord for efficient preloading
users, err := gorm.G[User](db).Preload("Posts", func(db gorm.PreloadBuilder) error {
    db.Order("created_at DESC").LimitPerRecord(5)
    return nil
}).Find(ctx)
```

### 4. Handle Errors Properly
```go
// ✅ Good - Comprehensive error handling
user, err := gorm.G[User](db).Where("id = ?", id).First(ctx)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

## Common Patterns

### Repository Pattern with Generics

```go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    return gorm.G[User](r.db).Create(ctx, user)
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*User, error) {
    return gorm.G[User](r.db).Where("id = ?", id).First(ctx)
}

func (r *UserRepository) GetActiveUsers(ctx context.Context) ([]User, error) {
    return gorm.G[User](r.db).Where("active = ?", true).Find(ctx)
}
```

### Service Layer with Generics

```go
type UserService struct {
    repo *UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    user := &User{
        Name:  req.Name,
        Email: req.Email,
        Age:   req.Age,
    }
    
    err := s.repo.Create(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

## Summary

GORM's Generics API represents a significant evolution in Go ORM design:

- **Type Safety**: Compile-time type checking prevents runtime errors
- **Performance**: Reduced SQL pollution and better connection reuse
- **Developer Experience**: Enhanced IDE support and cleaner APIs
- **Modern Patterns**: Context-first design and enhanced error handling
- **Backward Compatibility**: Works alongside traditional APIs

Start using the Generics API in new projects to experience these benefits firsthand! 