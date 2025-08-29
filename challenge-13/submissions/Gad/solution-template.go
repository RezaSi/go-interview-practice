package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Queries time-out
const TIME_OUT time.Duration = 3 * time.Second

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
	// TODO: Open a SQLite database connection
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.New("cannot open database")
	}

	createTable := `CREATE TABLE IF NOT EXISTS products(
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT,
						price REAL,
						quantity INTEGER,
						category TEXT
	)`

	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	if _, err := conn.ExecContext(ctx, createTable); err != nil {
		return nil, errors.New("cannot create table")
	}

	return conn, nil

}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	//  Insert the product into the database
	//  Update the product.ID with the database-generated ID
	queryCreateProduct := `INSERT INTO products(name, price, quantity, category)
	VALUES ($1, $2, $3, $4) RETURNING id;`

	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	err := ps.db.QueryRowContext(
		ctx,
		queryCreateProduct,
		product.Name,
		product.Price,
		product.Quantity,
		product.Category,
	).Scan(&product.ID)

	if err != nil {
		return errors.New("cannot create product")
	}

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// Query the database for a product with the given ID
	// Return a Product struct populated with the data or an error if not found

	var p *Product = &Product{
		ID: id,
	}

	queryGetProduct := `SELECT name, price, quantity, category 
	FROM products 
	WHERE id=$1;`

	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	err := ps.db.QueryRowContext(
		ctx,
		queryGetProduct,
		id,
	).Scan(&p.Name, &p.Price, &p.Quantity, &p.Category)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("cannot get product")
	}

	return p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// Update the product in the database
	// Return an error if the product doesn't exist

	queryUpdateProduct := `UPDATE products
	SET name = $1, price = $2, quantity = $3, category = $4 
	WHERE id=$5;`

	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // will be ignored if commit
	if _, err := ps.GetProduct(product.ID); err != nil {
		return errors.New("update error : product does not exist")
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	if _, err := ps.db.ExecContext(ctx,
		queryUpdateProduct,
		product.Name,
		product.Price,
		product.Quantity,
		product.Category,
		product.ID); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.New("update error: transaction failed")
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// Delete the product from the database
	// Return an error if the product doesn't exist
	queryDeleteProduct := `DELETE from products
								WHERE id=$5;`

	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // will be ignored if commit
	if _, err := ps.GetProduct(id); err != nil {
		return errors.New("delete error : product does not exist")
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	if _, err := ps.db.ExecContext(ctx,
		queryDeleteProduct,
		id,
	); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.New("delete error : transaction failed")
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// Query the database for products
	// If category is not empty, filter by category
	// Return a slice of Product pointers

	var productsList []*Product

	queryListProducts := `SELECT * from products
								WHERE category LIKE $0;`
	// select all products when category is empty
	if category == "" {

		category = "%"
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	rows, err := ps.db.QueryContext(ctx, queryListProducts, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category); err != nil {
			return nil, errors.New("ListProducts : error while scanning rows")
		}
		productsList = append(productsList, &p)
	}

	return productsList, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// Start a transaction
	// For each product ID in the updates map, update its quantity
	// If any update fails, roll back the transaction
	// Otherwise, commit the transaction

	queryUpdateProduct := `UPDATE products
	SET quantity = $1
	WHERE id=$2;`

	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()
	for id, newQuantity := range updates {
		var exists bool
		err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", id).Scan(&exists)
		if err != nil || !exists {
			return errors.New("product does not exist")
		}
		

		if _, err := tx.ExecContext(ctx,
			queryUpdateProduct,
			newQuantity,
			id,

		); err != nil {
			return err
		}

	}

	if err := tx.Commit(); err != nil {
		return errors.New("delete error : transaction failed")
	}
	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
