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

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite db: %w", err)
	}

	_, err = db.Exec(`
	  CREATE TABLE IF NOT EXISTS products (
	    id INTEGER PRIMARY KEY,
		name TEXT,
		price REAL,
		quantity INTEGER,
		category TEXT
	  )
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table products: %w", err)
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	result, err := ps.db.Exec(`
	  INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)
	`, product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return fmt.Errorf("failed to insert into product table: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	product.ID = id
	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	result := ps.db.QueryRow(`
	  SELECT * FROM products WHERE ID = ?
	`, id)
	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	product := &Product{}
	err := result.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	result, err := ps.db.Exec(`
	  UPDATE products
	  SET name = ?, price = ?, quantity = ?, category = ?
	  WHERE id = ?
	`, product.Name, product.Price, product.Quantity, product.Category, product.ID)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("no product with ID %d", product.ID)
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	result, err := ps.db.Exec(`
	  DELETE FROM products WHERE id = ?
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	} else if affected != 1 {
		return fmt.Errorf("wrong number of rows affected, want 1, got %d", affected)
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var results *sql.Rows
	var err error
	if category == "" {
		results, err = ps.db.Query(`
	  		SELECT * FROM products
		`)
	} else {
		results, err = ps.db.Query(`
			SELECT * FROM products WHERE category = ?
		`, category)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}

	defer results.Close()

	products := []*Product{}
	for results.Next() {
		p := &Product{}
		err := results.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to parse product: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare(`UPDATE products SET quantity = ? WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for id, quantity := range updates {
		result, err := stmt.Exec(quantity, id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update quantity for ID %d: %w", id, err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to retrieve rows affected for ID %d: %w", id, err)
		}
		if affected != 1 {
			tx.Rollback()
			return fmt.Errorf("wrong number of affected rows, want 1, got %d", affected)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
