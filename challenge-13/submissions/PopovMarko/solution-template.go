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
	// Open a SQLite database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the products table if it doesn't exist
	// The table have columns: id, name, price, quantity, category
	stmt := `
	CREATE TABLE IF NOT EXISTS  products (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	price REAL,
	quantity INTEGER,
	category TEXT
	)	
	`
	_, err = db.Exec(stmt)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("init products table: %w\n", err)
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// Insert the product into the database
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}

	// Defer rollback if transation not commited
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	// Prepare sql statement
	stmt, err := tx.Prepare("insert into products (name, price, quantity, category) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	// Defer statement close
	defer stmt.Close()
	result, err := stmt.Exec(product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}

	//  Update the product.ID with the database-generated ID
	product.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	// Commit transation and set commited flag to true
	err = tx.Commit()
	if err != nil {
		return err
	}
	commited = true
	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	//  Query the database for a product with the given ID
	row := ps.db.QueryRow("select * from products where id = ?", id)
	if err := row.Err(); err != nil {
		_ = ps.db.Close()
		return nil, err
	}
	res := Product{}
	err := row.Scan(&res.ID, &res.Name, &res.Price, &res.Quantity, &res.Category)
	if err != nil {
		_ = ps.db.Close()
		return nil, err
	}
	return &res, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// Update the product in the database
	res, err := ps.db.Exec(`UPDATE products 
		SET 
		name = ?, 
		price = ?, 
		quantity = ?, 
		category = ? 
		WHERE id = ?`,
		product.Name,
		product.Price,
		product.Quantity,
		product.Category,
		product.ID)
	if err != nil {
		_ = ps.db.Close()
		return err
	}

	// Return an error if the product doesn't exist
	rowN, err := res.RowsAffected()
	if err != nil {
		ps.db.Close()
		return err
	}
	if rowN == 0 {
		ps.db.Close()
		return errors.New("product does not exist")
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// Delete the product from the database
	res, err := ps.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		ps.db.Close()
		return err
	}

	// Return an error if the product doesn't exist
	rowN, err := res.RowsAffected()
	if err != nil {
		ps.db.Close()
		return err
	}
	// Return en error if the product does not exist
	if rowN == 0 {
		ps.db.Close()
		return errors.New("product does not exist")
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// Query the database for products
	// If category is not empty, filter by category
	var (
		rows *sql.Rows
		err  error
	)

	// Query database with filter or not due to category
	if category == "" {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products")
	} else {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products WHERE category = ?", category)
	}
	if err != nil {
		_ = ps.db.Close()
		return nil, err
	}

	defer rows.Close()
	res := []*Product{}

	// Map database rows to Product structures
	for rows.Next() {
		p := Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			_ = ps.db.Close()
			return nil, err
		}
		res = append(res, &p)
	}

	// Check for row error
	if err := rows.Err(); err != nil {
		_ = ps.db.Close()
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

	// For each product ID in the updates map, update its quantity
	// If any update fails, roll back the transaction
	// Otherwise, commit the transaction
	stmt, err := tx.Prepare("UPDATE products SET quantity = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for k, v := range updates {
		res, err := stmt.Exec(v, k)
		if err != nil {
			return err
		}
		r, _ := res.RowsAffected()
		if r == 0 {
			return fmt.Errorf("Product with ID %d does not exist", k)
		}
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = tx.Commit()
	return err
}
