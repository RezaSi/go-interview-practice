# Challenge 3: Database Migrations

Build an **E-commerce System** using GORM that demonstrates database migrations, schema evolution, and version control for database changes.

## Challenge Requirements

Create a Go application that implements:

1. **Database Migrations** - Version-controlled schema changes
2. **Schema Evolution** - Adding, modifying, and removing database structures
3. **Migration Rollbacks** - Ability to revert schema changes
4. **Data Seeding** - Populate database with initial data

## Data Models

```go
// Version 1: Basic product system
type Product struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`
    Price       float64   `gorm:"not null"`
    Description string    `gorm:"type:text"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Version 2: Add categories
type Category struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"unique;not null"`
    Description string    `gorm:"type:text"`
    Products    []Product `gorm:"foreignKey:CategoryID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Version 3: Enhanced product with inventory
type Product struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`
    Price       float64   `gorm:"not null"`
    Description string    `gorm:"type:text"`
    CategoryID  uint      `gorm:"not null"`
    Category    Category  `gorm:"foreignKey:CategoryID"`
    Stock       int       `gorm:"default:0"`
    SKU         string    `gorm:"unique;not null"`
    IsActive    bool      `gorm:"default:true"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## Required Functions

Implement these functions:
- `ConnectDB() (*gorm.DB, error)` - Database connection
- `RunMigration(db *gorm.DB, version int) error` - Run specific migration version
- `RollbackMigration(db *gorm.DB, version int) error` - Rollback to specific version
- `GetMigrationVersion(db *gorm.DB) (int, error)` - Get current migration version
- `SeedData(db *gorm.DB) error` - Seed database with initial data
- `CreateProduct(db *gorm.DB, product *Product) error` - Create product with validation
- `GetProductsByCategory(db *gorm.DB, categoryID uint) ([]Product, error)` - Get products by category
- `UpdateProductStock(db *gorm.DB, productID uint, quantity int) error` - Update product stock

## Migration Versions

**Version 1**: Create basic products table
**Version 2**: Add categories table and foreign key relationship
**Version 3**: Add inventory fields (stock, SKU, is_active) to products

## Testing Requirements

Your solution must pass tests for:
- Running migrations in sequence
- Rolling back migrations
- Tracking migration version
- Seeding initial data
- Creating products with category relationships
- Querying products by category
- Updating product inventory
- Handling migration conflicts and errors 