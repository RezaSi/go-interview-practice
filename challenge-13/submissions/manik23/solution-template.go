package main

import (
	"database/sql"
	"errors"
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
	// TODO: Open a SQLite database connection
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if db.Ping() != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error connecting database %s: %w", dbPath, err)
	}

	initTable := `create TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 0,
		category TEXT
	)`

	if _, err := db.Exec(initTable); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	query := `INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)`

	result, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category)
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
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found

	query := `SELECT id, name, price, quantity, category FROM products WHERE id = ?`

	p := &Product{}

	err := ps.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, err

	}
	return p, nil

}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist

	query := `UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?`

	result, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("update failed: product with id %d not found", product.ID)
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	query := `DELETE FROM products WHERE id = ?`

	// 1. Execute the deletion
	result, err := ps.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("delete failed: product with id %d not found", id)
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	query := `SELECT id, name, price, quantity, category FROM products WHERE (? = '' OR category = ?)`

	rows, err := ps.db.Query(query, category, category)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var products []*Product
	for rows.Next() {
		p := &Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE products SET quantity = ? WHERE id = ?`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, newQuantity := range updates {
		result, err := stmt.Exec(newQuantity, id)
		if err != nil {
			return fmt.Errorf("failed to update product %d: %w", id, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("product with id %d not found", id)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func main() {
	// Optional: you can write code here to test your implementation
}
