package main

import (
	"database/sql"
	"errors"
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

	// The table should have columns: id, name, price, quantity, category
	query := `CREATE TABLE IF NOT EXISTS products (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT NOT NULL,
    	price REAL NOT NULL,
		quantity INTEGER NOT NULL,
		category TEXT NOT NULL
	)`

	_, err = db.Exec(query)

	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	query := `INSERT INTO products(name, price, quantity, category) VALUES(?, ?, ?, ?)`

	res, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}

	if lastID, err := res.LastInsertId(); err == nil {
		product.ID = lastID
	}

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	query := `SELECT * FROM products WHERE id = ?`

	row := ps.db.QueryRow(query, id)
	p := Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	if _, err := ps.GetProduct(product.ID); err != nil {
		return err
	}

	query := `UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?`
	_, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	if _, err := ps.GetProduct(id); err != nil {
		return err
	}

	query := `DELETE FROM products WHERE id = ?`
	_, err := ps.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var query strings.Builder
	query.WriteString("SELECT * FROM products")
	if category != "" {
		query.WriteString(" WHERE category = ?")
	}

	rows, err := ps.db.Query(query.String(), category)
	if err != nil {
		return nil, err
	}

	var products []*Product
	p := Product{}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	fail := func(err error) error {
		return err
	}

	tx, err := ps.db.Begin()
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	for key, value := range updates {
		res, err1 := tx.Exec(`UPDATE products SET quantity = ? WHERE id = ?`, value, key)
		if err1 != nil {
			return fail(err1)
		}
		affected, err2 := res.RowsAffected()
		if err2 != nil {
			return fail(err2)
		}
		if affected == 0 {
			return errors.New("failed to update product: not found")
		}
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return err
}
