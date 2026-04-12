# Hints for GORM Migrations Challenge

## Hint 1: Understanding Migration Versions

Each migration should have a unique version number that represents the order of application. Use a `MigrationVersion` table to track which migrations have been applied.

## Hint 2: Database Connection Setup

Use `gorm.Open()` with SQLite driver. Don't auto-migrate all models here - let migrations handle schema changes:

```go
func ConnectDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Create migration tracking table
    db.AutoMigrate(&MigrationVersion{})
    return db, nil
}
```

## Hint 3: Running Migrations Safely

Check if migration already exists, use transactions, and record the migration version after successful application:

```go
func RunMigration(db *gorm.DB, version int) error {
    // Check if already applied
    var existing MigrationVersion
    if err := db.Where("version = ?", version).First(&existing).Error; err == nil {
        return nil // Already applied
    }
    
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Apply migration based on version
    switch version {
    case 1:
        if err := tx.AutoMigrate(&Product{}).Error; err != nil {
            tx.Rollback()
            return err
        }
    case 2:
        if err := tx.AutoMigrate(&Category{}).Error; err != nil {
            tx.Rollback()
            return err
        }
    case 3:
        if err := tx.Exec("ALTER TABLE products ADD COLUMN stock INTEGER DEFAULT 0").Error; err != nil {
            tx.Rollback()
            return err
        }
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

## Hint 4: Migration Rollback

Apply rollback operations in reverse order and remove the migration record:

```go
func RollbackMigration(db *gorm.DB, targetVersion int) error {
    current, err := GetMigrationVersion(db)
    if err != nil {
        return err
    }
    
    for version := current; version > targetVersion; version-- {
        tx := db.Begin()
        
        switch version {
        case 3:
            tx.Exec("ALTER TABLE products DROP COLUMN stock")
        case 2:
            tx.Migrator().DropTable(&Category{})
        case 1:
            tx.Migrator().DropTable(&Product{})
        }
        
        tx.Where("version = ?", version).Delete(&MigrationVersion{})
        tx.Commit()
    }
    return nil
}
```

## Hint 5: Tracking Migration Version

Query the `MigrationVersion` table for the highest version number:

```go
func GetMigrationVersion(db *gorm.DB) (int, error) {
    var migration MigrationVersion
    err := db.Order("version DESC").First(&migration).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return 0, nil
        }
        return 0, err
    }
    return migration.Version, nil
}
```

## Hint 6: Data Seeding

Create sample data with proper relationships and use transactions for consistency:

```go
func SeedData(db *gorm.DB) error {
    categories := []Category{
        {Name: "Technology", Description: "Tech products"},
        {Name: "Sports", Description: "Sports equipment"},
    }
    
    for _, cat := range categories {
        db.Create(&cat)
    }
    
    products := []Product{
        {Name: "Laptop", Price: 999.99, CategoryID: 1, Stock: 10, SKU: "LAP-001"},
        {Name: "Football", Price: 29.99, CategoryID: 2, Stock: 50, SKU: "SPT-001"},
    }
    
    for _, prod := range products {
        db.Create(&prod)
    }
    
    return nil
}
```

## Hint 7: Creating Products

Validate required fields and check if the category exists:

```go
func CreateProduct(db *gorm.DB, product *Product) error {
    // Validate required fields
    if product.Name == "" || product.Price <= 0 || product.SKU == "" {
        return errors.New("missing required fields")
    }
    
    // Check if category exists
    var category Category
    if err := db.First(&category, product.CategoryID).Error; err != nil {
        return errors.New("category not found")
    }
    
    return db.Create(product).Error
}
```

## Hint 8: Querying Products by Category

Use `Where()` to filter and implement pagination:

```go
func GetProductsByCategory(db *gorm.DB, categoryID uint, page, pageSize int) ([]Product, int64, error) {
    var products []Product
    var total int64
    
    query := db.Where("category_id = ?", categoryID)
    query.Model(&Product{}).Count(&total)
    
    offset := (page - 1) * pageSize
    err := query.Offset(offset).Limit(pageSize).Find(&products).Error
    
    return products, total, err
}
```

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