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
func InitDB(dbPath string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")
	if err != nil {
		db.Close()
		db = nil
	}
	return
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	if ps == nil || ps.db == nil {
		return sql.ErrConnDone
	}

	if product == nil {
		return errors.New("product cannot be nil")
	}
	res, err := ps.db.Exec("INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)", product.Name, product.Price, product.Quantity, product.Category)
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
	if ps == nil || ps.db == nil {
		return nil, sql.ErrConnDone
	}
	p := &Product{}
	err := ps.db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product of %d id not found", id)
		}
		return nil, err
	}
	return p, nil
}

func convert(rows *sql.Rows) (products []*Product, err error) {
	defer rows.Close()
	for rows.Next() {
		p := &Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	if ps == nil || ps.db == nil {
		return sql.ErrConnDone
	}

	if product == nil {
		return errors.New("product cannot be nil")
	}
	res, err := ps.db.Exec("UPDATE products set name = ?, price = ?, quantity = ?, category = ? WHERE id = ?", product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}

	changes, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if changes == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	if ps == nil || ps.db == nil {
		return sql.ErrConnDone
	}
	res, err := ps.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}
	changes, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if changes == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	if ps == nil || ps.db == nil {
		return nil, sql.ErrConnDone
	}
	var rows *sql.Rows
	var err error
	if category != "" {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products WHERE category = ?", category)
	} else {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products")
	}
	if err != nil {
		return nil, err
	}
	return convert(rows)
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	if ps == nil || ps.db == nil {
		return sql.ErrConnDone
	}
	tr, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tr.Rollback() // Safe to call after Commit, does nothing if already committed

	for key, quantity := range updates {
		res, err := tr.Exec("UPDATE products set quantity = ? WHERE id = ?", quantity, key)
		if err != nil {
			return err
		}

		changes, err := res.RowsAffected()
		if err != nil {
			return err
		} else if changes == 0 {
			return fmt.Errorf("product not found to update of id %d", key)
		}
	}

	return tr.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
