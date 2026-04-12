# Challenge 5: The Generics Way

Build a **Modern User & Post Management System** using GORM's new Generics API (v1.30.0+) that demonstrates type-safe database operations and enhanced features.

## Challenge Requirements

Create a Go application that leverages GORM's generics API to implement:

1. **Context-Aware Operations** - All operations use context for better control
2. **Type-Safe CRUD** - Using `gorm.G[T]` for type safety and reduced SQL pollution
3. **Enhanced Joins & Preload** - Advanced association handling with the new APIs
4. **Advanced Features** - OnConflict handling, execution hints, and result metadata
5. **Performance Optimizations** - Batch operations and connection management

## Data Models

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    Age       int       `gorm:"check:age > 0"`
    CompanyID *uint     `gorm:"index"`
    Company   *Company  `gorm:"foreignKey:CompanyID"`
    Posts     []Post    `gorm:"foreignKey:UserID"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Company struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null;unique"`
    Industry    string    `gorm:"not null"`
    FoundedYear int       `gorm:"not null"`
    Users       []User    `gorm:"foreignKey:CompanyID"`
    CreatedAt   time.Time
}

type Post struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `gorm:"not null"`
    Content   string    `gorm:"type:text"`
    UserID    uint      `gorm:"not null;index"`
    User      User      `gorm:"foreignKey:UserID"`
    ViewCount int       `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Required Functions

Implement these functions using GORM's Generics API:

### Basic Operations
- `ConnectDB() (*gorm.DB, error)` - Database connection with auto-migration
- `CreateUser(ctx context.Context, db *gorm.DB, user *User) error` - Create user with generics
- `GetUserByID(ctx context.Context, db *gorm.DB, id uint) (*User, error)` - Get user by ID
- `UpdateUserAge(ctx context.Context, db *gorm.DB, userID uint, age int) error` - Update specific field
- `DeleteUser(ctx context.Context, db *gorm.DB, userID uint) error` - Delete user

### Batch Operations
- `CreateUsersInBatches(ctx context.Context, db *gorm.DB, users []User, batchSize int) error` - Batch creation
- `FindUsersByAgeRange(ctx context.Context, db *gorm.DB, minAge, maxAge int) ([]User, error)` - Range queries

### Advanced Features
- `UpsertUser(ctx context.Context, db *gorm.DB, user *User) error` - OnConflict handling
- `CreateUserWithResult(ctx context.Context, db *gorm.DB, user *User) (int64, error)` - Return metadata

### Enhanced Associations
- `GetUsersWithCompany(ctx context.Context, db *gorm.DB) ([]User, error)` - Enhanced joins
- `GetUsersWithPosts(ctx context.Context, db *gorm.DB, limit int) ([]User, error)` - Preload with limits
- `GetUserWithPostsAndCompany(ctx context.Context, db *gorm.DB, userID uint) (*User, error)` - Multiple preloads

### Complex Queries
- `SearchUsersInCompany(ctx context.Context, db *gorm.DB, companyName string) ([]User, error)` - Join with filters
- `GetTopActiveUsers(ctx context.Context, db *gorm.DB, limit int) ([]User, error)` - Users with most posts

## Key Generics Features to Demonstrate

### 1. Type-Safe Operations
```go
// Instead of: db.Where("name = ?", name).First(&user)
user, err := gorm.G[User](db).Where("name = ?", name).First(ctx)
```

### 2. Context Support
```go
// All operations require context
ctx := context.Background()
users, err := gorm.G[User](db).Find(ctx)
```

### 3. OnConflict Handling
```go
// Handle duplicate key conflicts
err := gorm.G[User](db, clause.OnConflict{DoNothing: true}).Create(ctx, &user)
```

### 4. Enhanced Joins
```go
// More flexible join conditions
users, err := gorm.G[User](db).Joins(clause.LeftJoin.Association("Company"), 
    func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
        db.Where("companies.industry = ?", "Technology")
        return nil
    }).Find(ctx)
```

### 5. Preload Enhancements
```go
// Limit per record and custom conditions
users, err := gorm.G[User](db).Preload("Posts", func(db gorm.PreloadBuilder) error {
    db.Order("created_at DESC").LimitPerRecord(3)
    return nil
}).Find(ctx)
```

## Testing Requirements

Your solution must pass tests for:
- Context-aware database operations
- Type-safe CRUD operations with generics
- Batch operations and performance optimizations
- OnConflict handling for duplicate data
- Enhanced joins with custom conditions
- Advanced preloading with limits and filters
- Complex queries combining multiple features
- Proper error handling and context cancellation

## Performance Benefits

The generics API provides:
- **Type Safety** - Compile-time type checking
- **Reduced SQL Pollution** - Better connection reuse
- **Enhanced Performance** - Optimized query building
- **Better Tooling** - IDE support and autocompletion

## Migration from Traditional API

If migrating existing code:
```go
// Traditional API
var user User
db.Where("id = ?", id).First(&user)

// Generics API
user, err := gorm.G[User](db).Where("id = ?", id).First(ctx)
```

## Requirements

- Go 1.18+ (for generics support)
- GORM v1.30.0+ (for generics API support)
- Context-aware programming patterns

Start implementing and experience the improved type safety and performance of GORM's generics API! 🚀 