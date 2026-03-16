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
	// Open a SQLite database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// Insert the product into the database
	result, err := ps.db.Exec(
		"INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
		product.Name, product.Price, product.Quantity, product.Category)
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
	row := ps.db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE id = ?", id)

	p := &Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	// Return a Product struct populated with the data or an error if not found
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
	// Update the product in the database
	result, err := ps.db.Exec(
		"UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?",
		product.Name,
		product.Price,
		product.Quantity,
		product.Category,
		product.ID,
	)

	// Return an error if the product doesn't exist
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// Delete the product from the database
	result, err := ps.db.Exec(
		"DELETE FROM products WHERE id = ?",
		id,
	)
	// Return an error if the product doesn't exist
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
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// Query the database for products
	// If category is not empty, filter by category
	var rows *sql.Rows
	var err error
	if category != "" {
		rows, err = ps.db.Query(
			"SELECT id, name, price, quantity, category FROM products WHERE category = ?",
			category,
		)
		if err != nil {
			return nil, err
		}
	} else if category == "" {
		// If category is empty, return all products
		rows, err = ps.db.Query(
			"SELECT id, name, price, quantity, category FROM products",
		)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()
	res := []*Product{}
	for rows.Next() {
		p := &Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// Return a slice of Product pointers
	return res, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// Start a transaction
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// For each product ID in the updates map, update its quantity
	stmt, err := tx.Prepare(
		"UPDATE products SET quantity = ? WHERE id = ?",
	)
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
		// If any update fails, roll back the transaction
		if rowsAffected == 0 {
			return fmt.Errorf("product with ID %d not found", id)
		}
	}
	// Otherwise, commit the transaction
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
