package main

import (
    "database/sql"
    "fmt"

    _ "github.com/mattn/go-sqlite3"
)

// Product represents a product in the inventory system
type Product struct {
    ID       int64
    Name     string
    Price    float64
    Quantity int
    Category string
}

// ProductStore manages product operations
type ProductStore struct {
    db *sql.DB
}

// NewProductStore creates a new ProductStore with the given database connection
func NewProductStore(db *sql.DB) *ProductStore {
    return &ProductStore{db: db}
}

var DB *sql.DB

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dbPath string) (*sql.DB, error) {
    // Open a SQLite database connection
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    // Create the products table if it doesn't exist
    sqlStmt := `
    CREATE TABLE IF NOT EXISTS products (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price REAL NOT NULL,
        quantity INTEGER NOT NULL,
        category TEXT
    )
    `
    _, err = db.Exec(sqlStmt)
    if err != nil {
        return nil, err
    }
    
    return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
    // Insert the product into the database
    sqlQuery := `
        INSERT INTO products (name, price, quantity, category)
        VALUES (?, ?, ?, ?)
    `
    
    result, err := ps.db.Exec(sqlQuery, product.Name, product.Price, product.Quantity, product.Category)
    if err != nil {
        return err
    }
    
    // Update the product.ID with the database-generated ID
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    
    product.ID = id
    return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
    // Query the database for a product with the given ID
    sqlQuery := `
        SELECT id, name, price, quantity, category
        FROM products
        WHERE id = ?
    `
    
    var product Product
    err := ps.db.QueryRow(sqlQuery, id).Scan(
        &product.ID,
        &product.Name,
        &product.Price,
        &product.Quantity,
        &product.Category,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("product not found with ID: %d", id)
        }
        return nil, fmt.Errorf("error querying product: %v", err)
    }
    
    return &product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
    // Update the product in the database
    sqlQuery := `
        UPDATE products
        SET name = ?, price = ?, quantity = ?, category = ? 
        WHERE id = ?
    `
    
    result, err := ps.db.Exec(sqlQuery,
        product.Name,
        product.Price,
        product.Quantity,
        product.Category,
        product.ID,
    )
    
    if err != nil {
        return fmt.Errorf("error updating product: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error checking affected rows: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("product not found with ID: %d", product.ID)
    }
    
    return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
    // Delete the product from the database
    sqlQuery := `
        DELETE FROM products
        WHERE id = ?
    `
    
    result, err := ps.db.Exec(sqlQuery, id)
    if err != nil {
        return fmt.Errorf("error deleting product: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error checking affected rows: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("product not found with ID: %d", id)
    }
    
    return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
    sqlQuery := `
        SELECT id, name, price, quantity, category
        FROM products
    `
    
    var args []interface{}
    
    if category != "" {
        sqlQuery += " WHERE category = ?"
        args = append(args, category)
    }
    
    sqlQuery += " ORDER BY name ASC"
    
    rows, err := ps.db.Query(sqlQuery, args...)
    if err != nil {
        return nil, fmt.Errorf("error querying products: %v", err)
    }
    defer rows.Close()
    
    var products []*Product
    
    for rows.Next() {
        var product Product
        err := rows.Scan(
            &product.ID,
            &product.Name,
            &product.Price,
            &product.Quantity,
            &product.Category,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning product row: %v", err)
        }
        
        products = append(products, &product)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating product rows: %v", err)
    }
    
    return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
    // Start a transaction
    tx, err := ps.db.Begin()
    if err != nil {
        return fmt.Errorf("error beginning transaction: %v", err)
    }
    
    // Defer a rollback in case anything fails
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    // Prepare the update statement
    stmt, err := tx.Prepare(`
        UPDATE products 
        SET quantity = ? 
        WHERE id = ?
    `)
    if err != nil {
        return fmt.Errorf("error preparing update statement: %v", err)
    }
    defer stmt.Close()
    
    for productID, newQuantity := range updates {
        // Update the product with the new quantity (not adding to it)
        result, err := stmt.Exec(newQuantity, productID)
        if err != nil {
            return fmt.Errorf("error updating product %d: %v", productID, err)
        }
        
        rowsAffected, err := result.RowsAffected()
        if err != nil {
            return fmt.Errorf("error checking affected rows for product %d: %v", productID, err)
        }
        
        if rowsAffected == 0 {
            return fmt.Errorf("product not found with ID: %d", productID)
        }
    }
    
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("error committing transaction: %v", err)
    }
    
    return nil
}