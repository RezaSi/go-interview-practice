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

const createTableSql = `CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name STRING NOT NULL,
            price REAL NOT NULL,
            quantity INTEGER NOT NULL,
            category STRING NOT NULL
        )`
const insertProductSql = `INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)`
const getProductSql = `SELECT * FROM products WHERE id = ?`
const updateProductSql = `UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?`
const deleteProductSql = `DELETE FROM products WHERE id = ?`
const listProductsSql = `SELECT * FROM products`
const updateQuantitySql = `UPDATE products SET quantity = ? WHERE id = ?`

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(createTableSql); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	result, err := ps.db.Exec(insertProductSql, product.Name, product.Price, product.Quantity, product.Category)
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
	row := ps.db.QueryRow(getProductSql, id)
	product := &Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
	return product, err
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	result, err := ps.db.Exec(updateProductSql, product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if updated == 0 {
		return fmt.Errorf("product with id = %d was not updated", product.ID)
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	result, err := ps.db.Exec(deleteProductSql, id)
	if err != nil {
		return err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return fmt.Errorf("product with id = %d was not deleted", id)
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var rows *sql.Rows
	var err error

	if len(category) > 0 {
		rows, err = ps.db.Query(listProductsSql+` WHERE category = ?`, category)
	} else {
		rows, err = ps.db.Query(listProductsSql)
	}
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()

	products := make([]*Product, 0)

	for rows.Next() {
		product := &Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
		if err != nil {
			return []*Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for id, quantity := range updates {
		result, err := tx.Exec(updateQuantitySql, quantity, id)
		if err != nil {
			return err
		}
		updated, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if updated == 0 {
			return fmt.Errorf("product with id = %d was not updated", id)
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
