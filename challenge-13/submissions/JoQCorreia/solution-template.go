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

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.New("Database failed to open")
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, errors.New("Db ping failed")
	}

	db.Exec(
		"CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	
	db := ps.db
	result, err := db.Exec(
		"INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
		product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}

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
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	db := ps.db
	row := db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE id = ?", id)

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

	test := ps.db.QueryRow("SELECT * FROM products WHERE id = ?", product.ID)

	if test.Err() != nil {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}

	ps.db.Exec("UPDATE products SET Name = ?, Price = ?, Quantity = ?, Category = ? WHERE ID = ?", product.Name, product.Price, product.Quantity, product.Category, product.ID)

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	test := ps.db.QueryRow("SELECT * FROM products WHERE id = ?", id)

	if test.Err() != nil {
		return fmt.Errorf("product with ID %d not found", id)
	}

	ps.db.Exec("DELETE FROM products WHERE ID = ?", id)

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var products []*Product

	var rows *sql.Rows
	var err error

	// if category is not empty, filter by category
	if category != "" {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products WHERE category = ?", category)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	} else {
		// if category is empty
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	}

	for rows.Next() {
		p := &Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
		tx, err := ps.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // Rollback on error
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

	// Commit the transaction
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
    

    

