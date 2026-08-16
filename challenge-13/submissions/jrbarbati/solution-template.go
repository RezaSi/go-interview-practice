package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

const (
    createProductsTable   = "CREATE TABLE IF NOT EXISTS product ( id INTEGER PRIMARY KEY, name VARCHAR(255) NOT NULL, price FLOAT NOT NULL, quantity INTEGER NOT NULL, category VARCHAR(255) )"
    selectAll             = "SELECT id, name, price, quantity, category FROM product"
    idClause              = " WHERE id = ?"
    categoryClause        = " WHERE category = ?"
    insertProduct         = "INSERT INTO product (name, price, quantity, category) VALUES (?, ?, ?, ?)"
    updateProduct         = "UPDATE product SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?"
    updateProductQuantity = "UPDATE product SET quantity = ? WHERE id = ?"
    deleteProduct         = "DELETE FROM product WHERE id = ?"
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

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    
    if err != nil {
        return nil, err
    }
    
    err = db.Ping()
    
    if err != nil {
        return nil, err
    }
    
    _, err = db.Exec(createProductsTable)
    
    if err != nil {
        return nil, err
    }
    
    return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
    result, err := ps.db.Exec(insertProduct, product.Name, product.Price, product.Quantity, product.Category)
    
    if err != nil {
        return err
    }
    
    id, err := result.LastInsertId()
    
    if err != nil {
        return err
    }
    
    product.ID = id
    return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
    var product Product
    
    row := ps.db.QueryRow(selectAll+idClause, id)
    
    err := row.Scan(
        &product.ID,
        &product.Name,
        &product.Price,
        &product.Quantity,
        &product.Category,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
    if _, err := ps.db.Exec(updateProduct, product.Name, product.Price, product.Quantity, product.Category, product.ID); err != nil {
        return err
    }
    
    return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
    if _, err := ps.db.Exec(deleteProduct, id); err != nil {
        return err
    }
    
    return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
    args := make([]any, 0)
    query := selectAll
    
    if category != "" {
        query += categoryClause
        args = append(args, category)
    }
    
    rows, err := ps.db.Query(query, args...)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    products := make([]*Product, 0)
    
    for rows.Next() {
        var product Product
        
        rowErr := rows.Scan(
            &product.ID,
            &product.Name,
            &product.Price,
            &product.Quantity,
            &product.Category,
        )
        
        if rowErr != nil {
            return nil, err
        }
        
        products = append(products, &product)
    }
    
	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
    tx, err := ps.db.Begin()
    
    if err != nil {
        return err
    }
    
    for id, quantity := range updates {
        result, err := tx.Exec(updateProductQuantity, quantity, id)
        
        if err != nil {
            _ = tx.Rollback()
            return err
        }
        
        if affected, err := result.RowsAffected(); affected == 0 || err != nil {
            _ = tx.Rollback()
            
            if err != nil {
                return err
            }
            
            return errors.New("invalid update, no rows affected")
        }
    }
    
    if err = tx.Commit(); err != nil {
        return err
    }
    
    return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
