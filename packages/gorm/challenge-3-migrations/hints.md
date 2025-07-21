# Hints for GORM Migrations Challenge

## General Tips

1. **Understand migration versions** - Each migration should have a unique version number that represents the order of application.

2. **Use transactions** - Always wrap migrations in transactions to ensure atomicity.

3. **Track migration state** - Use a `MigrationVersion` table to track which migrations have been applied.

4. **Make migrations reversible** - Each migration should have a corresponding rollback operation.

## Function-Specific Hints

### ConnectDB()
- Use `gorm.Open()` with SQLite driver
- Don't auto-migrate all models here - let migrations handle schema changes
- Return the database connection

### RunMigration()
- Check if migration already exists in `MigrationVersion` table
- Use a switch statement to handle different migration versions
- Wrap the entire operation in a transaction
- Record the migration version after successful application

### RollbackMigration()
- Find the current migration version
- Apply rollback operations in reverse order
- Remove the migration record from `MigrationVersion` table
- Use transactions for safety

### GetMigrationVersion()
- Query the `MigrationVersion` table for the highest version number
- Return 0 if no migrations have been applied
- Handle the case where the table doesn't exist yet

### SeedData()
- Create sample categories first
- Create sample products with proper category relationships
- Use transactions to ensure data consistency
- Check for existing data to avoid duplicates

### CreateProduct()
- Validate required fields (name, price, category_id, sku)
- Check if the category exists
- Handle unique constraint violations for SKU
- Use transactions if needed

### GetProductsByCategory()
- Use `Where()` to filter by category
- Implement pagination with `Offset()` and `Limit()`
- Count total records for pagination info
- Preload related data if needed

### UpdateProductStock()
- Find the product first
- Update the stock field
- Validate stock quantity (should be non-negative)
- Handle the case where product doesn't exist

## Migration Implementation Patterns

### Version 1: Basic Products
```go
func CreateProductsTable(db *gorm.DB) error {
    return db.AutoMigrate(&Product{})
}
```

### Version 2: Add Categories
```go
func AddCategoriesTable(db *gorm.DB) error {
    // Create categories table
    if err := db.AutoMigrate(&Category{}); err != nil {
        return err
    }
    
    // Add category_id to products table
    return db.Exec("ALTER TABLE products ADD COLUMN category_id INTEGER").Error
}
```

### Version 3: Add Inventory Fields
```go
func AddInventoryFields(db *gorm.DB) error {
    // Add new columns to products table
    return db.Exec(`
        ALTER TABLE products ADD COLUMN stock INTEGER DEFAULT 0;
        ALTER TABLE products ADD COLUMN sku VARCHAR(255) UNIQUE;
        ALTER TABLE products ADD COLUMN is_active BOOLEAN DEFAULT 1;
    `).Error
}
```

## Transaction Patterns

### Safe Migration Execution
```go
func RunMigration(db *gorm.DB, version int) error {
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Check if already applied
    var existing MigrationVersion
    if err := tx.Where("version = ?", version).First(&existing).Error; err == nil {
        return nil // Already applied
    }
    
    // Apply migration
    if err := applyMigration(tx, version); err != nil {
        tx.Rollback()
        return err
    }
    
    // Record migration
    migration := MigrationVersion{Version: version, AppliedAt: time.Now()}
    if err := tx.Create(&migration).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Error
}
```

## Error Handling

1. **Migration already applied** - Check before applying
2. **Invalid migration version** - Return error for unknown versions
3. **Database errors** - Handle SQL errors gracefully
4. **Rollback errors** - Ensure rollback operations are safe

## Testing Strategies

### Test Migration Sequence
```go
// Test running migrations in order
err := RunMigration(db, 1)
assert.NoError(t, err)

err = RunMigration(db, 2)
assert.NoError(t, err)

err = RunMigration(db, 3)
assert.NoError(t, err)

// Verify final state
assert.True(t, db.Migrator().HasTable(&Product{}))
assert.True(t, db.Migrator().HasTable(&Category{}))
```

### Test Rollback
```go
// Run migrations
RunMigration(db, 1)
RunMigration(db, 2)
RunMigration(db, 3)

// Rollback to version 2
err := RollbackMigration(db, 2)
assert.NoError(t, err)

// Verify state
version, _ := GetMigrationVersion(db)
assert.Equal(t, 2, version)
```

## Data Seeding Patterns

### Create Sample Data
```go
func SeedData(db *gorm.DB) error {
    // Create categories
    categories := []Category{
        {Name: "Technology", Description: "Tech products"},
        {Name: "Sports", Description: "Sports equipment"},
        {Name: "Food", Description: "Food items"},
    }
    
    for _, cat := range categories {
        if err := db.Create(&cat).Error; err != nil {
            return err
        }
    }
    
    // Create products
    products := []Product{
        {Name: "Laptop", Price: 999.99, CategoryID: 1, Stock: 10, SKU: "LAP-001"},
        {Name: "Football", Price: 29.99, CategoryID: 2, Stock: 50, SKU: "SPT-001"},
        {Name: "Coffee", Price: 5.99, CategoryID: 3, Stock: 100, SKU: "FOD-001"},
    }
    
    for _, prod := range products {
        if err := db.Create(&prod).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

## Common Mistakes to Avoid

1. **Not using transactions** - Migrations should be atomic
2. **Forgetting to track versions** - Always record applied migrations
3. **Not handling rollbacks** - Every migration should be reversible
4. **Not testing migrations** - Always test in development first
5. **Not handling errors** - Check for errors at each step

## SQLite Specific Notes

- SQLite has limited ALTER TABLE support
- Use `db.Exec()` for complex schema changes
- Some operations might require table recreation
- Be careful with foreign key constraints

## Debugging Tips

1. **Enable GORM logging**:
```go
db = db.Debug()
```

2. **Check migration state**:
```go
var versions []MigrationVersion
db.Find(&versions)
for _, v := range versions {
    fmt.Printf("Migration %d applied at %s\n", v.Version, v.AppliedAt)
}
```

3. **Verify table structure**:
```go
// Check if columns exist
columns, _ := db.Migrator().ColumnTypes(&Product{})
for _, col := range columns {
    fmt.Printf("Column: %s, Type: %s\n", col.Name(), col.DatabaseTypeName())
}
```

## Performance Considerations

1. **Batch operations** - Use transactions for multiple operations
2. **Index creation** - Add indexes after data migration
3. **Data validation** - Validate data before migration

## Useful GORM Methods

- `db.Begin()` - Start transaction
- `db.Commit()` - Commit transaction
- `db.Rollback()` - Rollback transaction
- `db.Exec()` - Execute raw SQL
- `db.Migrator()` - Access migration methods
- `db.AutoMigrate()` - Auto-migrate models

## Final Tips

1. **Start with simple migrations** - Get basic version tracking working first
2. **Test each migration** - Verify each step works before moving on
3. **Keep migrations small** - Each migration should do one thing well
4. **Document your migrations** - Add comments explaining what each migration does
5. **Use the learning resources** - Check GORM documentation for migration examples 