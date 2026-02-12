package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID       int64
	Name     string
	Price    float64
	Quantity int
	Category string
}

type ProductStore struct {
	db *sql.DB
}

// NewProductStore creates a new ProductStore with the given database connection
func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dPath string) (*sql.DB, error) {
	sql, err := sql.Open("sqlite3", dPath)
	if err != nil {
		return nil, err
	}

	_, err = sql.Exec(
		"CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)",
	)
	if err != nil {
		return nil, err
	}

	return sql, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	result, err := ps.db.Exec(
		"INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
		product.Name, product.Price, product.Quantity, product.Category,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = id

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	row := ps.db.QueryRow(
		"SELECT id, name, price, quantity, category FROM products WHERE id = ?",
		id,
	)

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
	_, err := ps.db.Exec("UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?",
		product.Name, product.Price, product.Quantity, product.Category, product.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	_, err := ps.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var rows *sql.Rows
	var err error
	if category == "" {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products")
	} else {
		rows, err = ps.db.Query(
			"SELECT id, name, price, quantity, category FROM products WHERE category = ?",
			category,
		)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*Product
	for rows.Next() {
		p := &Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
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
	defer func() {
		if err != nil {
			tx.Rollback()
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

	return tx.Commit()
}

func main() {
    
}