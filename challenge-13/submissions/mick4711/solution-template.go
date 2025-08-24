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
		return nil, err
	}

	// Check if the database is accessible
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		quantity INTEGER NOT NULL,
		category TEXT NOT NULL
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// Insert the product into the database
	result, err := ps.db.Exec(`
		INSERT INTO products (name, price, quantity, category)
		VALUES (?, ?, ?, ?)
	`, product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}

	// Update the product.ID with the database-generated ID
	product.ID, err = result.LastInsertId()
	return err
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// Query the database for a product with the given ID
	row := ps.db.QueryRow(`
		SELECT id, name, price, quantity, category
		FROM products
		WHERE id = ?
	`, id)

	// Scan the result into a Product struct
	product := &Product{}
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, err
	}
	return product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// Return an error if the product doesn't exist
	_, errExists := ps.GetProduct(product.ID)
	if errExists != nil {
		return errExists
	}

	_, err := ps.db.Exec(`
		UPDATE products
		SET name = ?, price = ?, quantity = ?, category = ?
		WHERE id = ?
	`, product.Name, product.Price, product.Quantity, product.Category, product.ID)
	return err
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// Return an error if the product doesn't exist
	_, errExists := ps.GetProduct(id)
	if errExists != nil {
		return errExists
	}

	// Execute the delete statement
	_, err := ps.db.Exec(`
		DELETE FROM products
		WHERE id = ?
	`, id)
	return err
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// Query the database for products, if category is not empty, filter by category
	rows, err := ps.db.Query(`
		SELECT id, name, price, quantity, category
		FROM products
		WHERE category = ? OR ? = ''
	`, category, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and scan each row into a Product struct
	var products []*Product
	for rows.Next() {
		product := &Product{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// Start a transaction
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}

	// For each product ID in the updates map, update its quantity
	for id, quantity := range updates {
		// Check if the product exists before updating
		_, errExists := ps.GetProduct(id)
		if errExists != nil {
			// If the product does not exist, roll back the transaction and return an error
			tx.Rollback()
			return errExists
		}

		// Update the product's quantity
		_, err := tx.Exec(`
			UPDATE products
			SET quantity = ?
			WHERE id = ?
		`, quantity, id)
		if err != nil {
			// If any update fails, roll back the transaction
			tx.Rollback()
			return err
		}
	}

	// Otherwise, commit the transaction
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
