# Hints for Challenge 13: SQL Database Operations with Go

## Hint 1: Database Setup and Connection
Start by setting up the SQLite database and creating the products table:
```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Create products table
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT,
        price REAL NOT NULL,
        stock INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    
    _, err = db.Exec(createTableSQL)
    if err != nil {
        return nil, err
    }
    
    return db, nil
}
```

## Hint 2: Product Struct and Model
Define the Product struct that matches your database schema:
```go
type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    CreatedAt   string  `json:"created_at"`
}
```

## Hint 3: Create Product with Parameter Binding
Use parameter binding to prevent SQL injection:
```go
func CreateProduct(db *sql.DB, product Product) (int, error) {
    query := `
    INSERT INTO products (name, description, price, stock) 
    VALUES (?, ?, ?, ?)`
    
    result, err := db.Exec(query, product.Name, product.Description, product.Price, product.Stock)
    if err != nil {
        return 0, fmt.Errorf("failed to create product: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert ID: %w", err)
    }
    
    return int(id), nil
}
```

## Hint 4: Get Product by ID
Use QueryRow for retrieving a single record:
```go
func GetProduct(db *sql.DB, id int) (*Product, error) {
    query := `
    SELECT id, name, description, price, stock, created_at 
    FROM products 
    WHERE id = ?`
    
    var product Product
    err := db.QueryRow(query, id).Scan(
        &product.ID,
        &product.Name,
        &product.Description,
        &product.Price,
        &product.Stock,
        &product.CreatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("product with ID %d not found", id)
        }
        return nil, fmt.Errorf("failed to get product: %w", err)
    }
    
    return &product, nil
}
```

## Hint 5: Update Product with Validation
Check if the product exists before updating:
```go
func UpdateProduct(db *sql.DB, id int, product Product) error {
    // First check if product exists
    _, err := GetProduct(db, id)
    if err != nil {
        return err
    }
    
    query := `
    UPDATE products 
    SET name = ?, description = ?, price = ?, stock = ? 
    WHERE id = ?`
    
    result, err := db.Exec(query, product.Name, product.Description, product.Price, product.Stock, id)
    if err != nil {
        return fmt.Errorf("failed to update product: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("no rows were updated")
    }
    
    return nil
}
```

## Hint 6: Delete Product
Implement delete with existence check:
```go
func DeleteProduct(db *sql.DB, id int) error {
    query := `DELETE FROM products WHERE id = ?`
    
    result, err := db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete product: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("product with ID %d not found", id)
    }
    
    return nil
}
```

## Hint 7: List Products with Filtering
Use Query for multiple records and implement optional filtering:
```go
func ListProducts(db *sql.DB, filters map[string]interface{}) ([]Product, error) {
    query := "SELECT id, name, description, price, stock, created_at FROM products"
    args := []interface{}{}
    conditions := []string{}
    
    // Add filters dynamically
    if name, ok := filters["name"]; ok {
        conditions = append(conditions, "name LIKE ?")
        args = append(args, "%"+name.(string)+"%")
    }
    
    if minPrice, ok := filters["min_price"]; ok {
        conditions = append(conditions, "price >= ?")
        args = append(args, minPrice)
    }
    
    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }
    
    query += " ORDER BY created_at DESC"
    
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to list products: %w", err)
    }
    defer rows.Close()
    
    var products []Product
    for rows.Next() {
        var product Product
        err := rows.Scan(&product.ID, &product.Name, &product.Description, 
                        &product.Price, &product.Stock, &product.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan product: %w", err)
        }
        products = append(products, product)
    }
    
    return products, nil
}
```

## Hint 8: Transaction Support
Implement transactions for operations that modify multiple records:
```go
func BulkUpdatePrices(db *sql.DB, updates map[int]float64) error {
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    // Use defer to handle rollback on panic or error
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    
    query := "UPDATE products SET price = ? WHERE id = ?"
    stmt, err := tx.Prepare(query)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()
    
    for id, price := range updates {
        _, err := stmt.Exec(price, id)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("failed to update price for product %d: %w", id, err)
        }
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}
```

## Key Database Concepts:
- **Parameter Binding**: Always use ? placeholders to prevent SQL injection
- **Error Handling**: Check for sql.ErrNoRows for not found cases
- **Resource Cleanup**: Always defer Close() on rows and statements
- **Transactions**: Use Begin/Commit/Rollback for multi-operation consistency
- **Prepared Statements**: Use Prepare() for repeated queries with different parameters
- **Context**: Use context-aware methods (QueryContext, ExecContext) for cancellation support 