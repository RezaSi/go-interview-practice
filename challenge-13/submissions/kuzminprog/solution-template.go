package main

import (
	"database/sql"
	"fmt"
	"strings"

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

	createTableStm := `
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name TEXT NOT NULL, 
			price REAL NOT NULL, 
			quantity INTEGER NOT NULL, 
			category TEXT NOT NULL
		)
	`

	_, err = db.Exec(createTableStm)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	insertStm := `
	INSERT INTO products (name, price, quantity, category)
	VALUES (?, ?, ?, ?)
	RETURNING id
	`

	res, err := ps.db.Exec(insertStm, product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = id

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	selectStm := `
	SELECT id, name, price, quantity, category
	FROM products
	WHERE id = ?
	`

	row := ps.db.QueryRow(selectStm, id)

	var p Product
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Quantity,
		&p.Category,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	updateStm := `
	UPDATE products
	SET name = ?, price = ?, quantity = ?, category = ?
	WHERE id = ?
	`

	row, err := ps.db.Exec(updateStm,
		product.Name,
		product.Price,
		product.Quantity,
		product.Category,
		product.ID,
	)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	deleteStm := `
	DELETE FROM products
	WHERE id = ?
	`

	row, err := ps.db.Exec(deleteStm, id)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	selectStm := `
	SELECT id, name, price, quantity, category
	FROM products
	`

	args := make([]interface{}, 0)
	if strings.TrimSpace(category) != "" {
		selectStm += " WHERE category = ?"
		args = append(args, category)
	}

	rows, err := ps.db.Query(selectStm, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*Product, 0)
	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Quantity,
			&p.Category,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	updateStm := `
	UPDATE products
	SET quantity = ?
	WHERE id = ?
	`

	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(updateStm)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, quantity := range updates {
		row, err := stmt.Exec(quantity, id)
		if err != nil {
			return err
		}

		rowAffected, err := row.RowsAffected()
		if err != nil {
			return err
		}
		if rowAffected == 0 {
			return fmt.Errorf("product with ID %d not found", id)
		}
	}
	return tx.Commit()
}
