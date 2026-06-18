package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	// Open a SQLite database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create the products table if it doesn't exist
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
		return nil, fmt.Errorf("failed to create products table: %v", err)
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// Insert the product into the database
	result, err := ps.db.Exec(
		"INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
		product.Name, product.Price, product.Quantity, product.Category,
	)
	if err != nil {
		return fmt.Errorf("failed to create product: %v", err)
	}

	// Update the product.ID with the database-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}
	product.ID = id

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// Query the database for a product with the given ID
	row := ps.db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE id = ?", id)

	product := &Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// Update the product in the database
	result, err := ps.db.Exec(
		"UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?",
		product.Name, product.Price, product.Quantity, product.Category, product.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %v", err)
	}

	// Check if the product exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// Delete the product from the database
	result, err := ps.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}

	// Check if the product exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var rows *sql.Rows
	var err error

	// If category is not empty, filter by category
	if category != "" {
		rows, err = ps.db.Query(
			"SELECT id, name, price, quantity, category FROM products WHERE category = ?",
			category,
		)
	} else {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list products: %v", err)
	}
	defer rows.Close()

	// Return a slice of Product pointers
	products := []*Product{}
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over products: %v", err)
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// Start a transaction
	tx, err := ps.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// For each product ID in the updates map, update its quantity
	for productID, newQuantity := range updates {
		result, err := tx.Exec("UPDATE products SET quantity = ? WHERE id = ?", newQuantity, productID)
		if err != nil {
			// If any update fails, roll back the transaction
			tx.Rollback()
			return fmt.Errorf("failed to update product %d: %v", productID, err)
		}

		// Check if the product exists
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to check rows affected for product %d: %v", productID, err)
		}
		if rowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("product with ID %d not found", productID)
		}
	}

	// Otherwise, commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func main() {
	// Example usage of the product inventory system
	db, err := InitDB("inventory.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := NewProductStore(db)

	// Create some sample products
	products := []*Product{
		{Name: "Laptop", Price: 999.99, Quantity: 10, Category: "Electronics"},
		{Name: "Book", Price: 19.99, Quantity: 50, Category: "Books"},
		{Name: "Mouse", Price: 29.99, Quantity: 30, Category: "Electronics"},
	}

	for _, product := range products {
		err := store.CreateProduct(product)
		if err != nil {
			log.Printf("Failed to create product %s: %v", product.Name, err)
		} else {
			log.Printf("Created product: %s (ID: %d)", product.Name, product.ID)
		}
	}

	// List all electronics products
	electronics, err := store.ListProducts("Electronics")
	if err != nil {
		log.Printf("Failed to list electronics: %v", err)
	} else {
		log.Printf("Electronics products: %d", len(electronics))
		for _, product := range electronics {
			log.Printf("  - %s: $%.2f (Qty: %d)", product.Name, product.Price, product.Quantity)
		}
	}
}