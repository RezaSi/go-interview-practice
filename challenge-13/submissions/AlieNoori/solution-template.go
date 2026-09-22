package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = errors.New("product nof found")

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
	if dbPath == "" {
		return nil, errors.New("empty database path")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			category TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	query := `
	INSERT INTO products (name, price, quantity, category)
	VALUES (?, ?, ?, ?)
	RETURNING id
	`

	err := ps.db.QueryRow(query,
		product.Name,
		product.Price,
		product.Quantity,
		product.Category).
		Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	product := &Product{}

	query := `
	SELECT id,name,price,quantity,category
	FROM products
	WHERE id = ?`
	if err := ps.db.
		QueryRow(query, id).
		Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&product.Category,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	query := `
	UPDATE products SET name=?, price=?, quantity=?, category=?
	WHERE id = ?
	`
	result, err := ps.db.Exec(query,
		product.Name,
		product.Price,
		product.Quantity,
		product.Category,
		product.ID,
	)
	if err != nil {
		return err
	}

	rewsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rewsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	query := "DELETE FROM products WHERE id=?"
	result, err := ps.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	query := `
	SELECT id, name, price, quantity, category 
	FROM products
	`
	rows, err := ps.db.Query(query)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	products := []*Product{}
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&product.Category,
		)
		if err != nil {
			return nil, err
		}
		if category != "" {
			if product.Category == category {
				products = append(products, product)
			}
		} else {
			products = append(products, product)
		}
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}

	query := `
	UPDATE products SET quantity = ? WHERE id = ?
	`
	for id, quantity := range updates {
		result, err := tx.Exec(query, quantity, id)
		if err != nil {
			tx.Rollback()
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if rowsAffected <= 0 {
			return ErrNotFound
		}
	}

	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
