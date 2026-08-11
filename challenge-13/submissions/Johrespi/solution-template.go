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
	// TODO: Open a SQLite database connection
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableQuery := ` CREATE TABLE IF NOT EXISTS products (
  		id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price REAL NOT NULL,
        quantity INTEGER NOT NULL,
        category TEXT NOT NULL
    );`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID
	stmt := `INSERT INTO products (name, price, quantity, category) VALUES(?, ?, ?, ?)`

	dbResult, err := ps.db.Exec(stmt, product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}

	id, err := dbResult.LastInsertId()
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

	stmt := `SELECT id, name, price, quantity, category FROM products WHERE id = ?`
	row := ps.db.QueryRow(stmt, id)

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil

}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist

	p, err := ps.GetProduct(product.ID)
	if err != nil {
		return err
	}

	p.ID = product.ID
	p.Name = product.Name
	p.Price = product.Price
	p.Quantity = product.Quantity
	p.Category = product.Category

	stmt := `UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?`

	_, err = ps.db.Exec(stmt, p.Name, p.Price, p.Quantity, p.Category, p.ID)
	if err != nil {
		return err
	}

	return nil

}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist

	stmt := `DELETE FROM products WHERE id = ?`

	result, err := ps.db.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil

}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	stmt := `SELECT id, name, price, quantity, category FROM products
		WHERE (? = '' OR category = ?)`

	var products []*Product

	rows, err := ps.db.Query(stmt, category, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil

}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction

	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt := `UPDATE products SET quantity = ? WHERE id = ?`

	for productID, newQuantity := range updates {
		result, err := tx.Exec(stmt, newQuantity, productID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("product with ID %d not found", productID)
		}

	}
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
