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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	res, err := ps.db.Exec(
		"INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
		product.Name,
		product.Price,
		product.Quantity,
		product.Category)
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
	var p Product
	err := ps.db.QueryRow("SELECT * FROM products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product does not exists, id: %d", id)
	} else if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	_, err := ps.db.Exec(
		"UPDATE products SET name=?, price=?, quantity=?, category=? WHERE id=?",
		product.Name,
		product.Price,
		product.Quantity,
		product.Category,
		product.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	_, err := ps.db.Exec("DELETE FROM products WHERE id=?", id)
	return err
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var rows *sql.Rows
	var err error

	if category == "" {
		rows, err = ps.db.Query("SELECT * FROM products")
	} else {
		rows, err = ps.db.Query("SELECT * FROM products WHERE category=?", category)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		p := new(Product)
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
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
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`UPDATE products SET quantity=? WHERE id=?`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for id, val := range updates {
		res, err := stmt.Exec(val, id)
		if err != nil {
			tx.Rollback()
			return err
		}

		nb, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if nb == 0 {
			tx.Rollback()
			return fmt.Errorf("product does not exists, id: %d", id)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
