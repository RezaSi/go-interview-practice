package main

import (
	"database/sql"
	//"errors"

	_ "github.com/mattn/go-sqlite3"
	"fmt"
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
	// Open a SQLite database connection
	db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Test the connection
    if err = db.Ping(); err != nil {
        return nil, err
    }
    
	
	// Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")
    if err != nil {
        return nil, err
    }
	
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// Insert the product into the database
	result, err := ps.db.Exec("INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
        product.Name, product.Price, product.Quantity, product.Category)
    if err != nil {
        return err
    }
    
    // Update the product.ID with the database-generated ID
    // Get the ID of the inserted row
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
	row := ps.db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE id = ?", id)

    p := &Product{}
    err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("product with ID %d not found", id)
        }
        return nil, err
    }
    
    return p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	_, err := ps.db.Exec("UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?",
        product.Name, product.Price, product.Quantity, product.Category, product.ID)
    if err != nil {
        return err
    }
    return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	_, err := ps.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
        return err
    }
    return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// Initialize the base query and an empty arguments slice
	query := "SELECT id, name, price, quantity, category FROM products"
	var args []any

	// Dynamically append the WHERE clause
	if category != "" {
		query += " WHERE category = ?"
		args = append(args, category)
	}

	// Execute the safe parameterized query
	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to prevent connection leaks

	// Iterate over the database results
	var products []*Product
	for rows.Next() {
		p := &Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return an empty slice instead of nil if no products were found
	if products == nil {
		return []*Product{}, nil
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}

	// Track whether the transaction successfully committed
	var committed bool
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare("UPDATE products SET quantity = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, quantity := range updates {
		result, err := stmt.Exec(quantity, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("product with ID %d not found", id)
		}
	}

	// If we reach this line, everything succeeded
	if err = tx.Commit(); err != nil {
		return err
	}
	
	committed = true // Prevent the deferred rollback from firing
	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
