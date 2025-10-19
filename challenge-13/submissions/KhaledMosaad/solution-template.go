package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Product represents a product in the inventory system
type Product struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
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

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS products (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT, price REAL, quantity INTEGER, category TEXT);`)
	if err != nil {
		return nil, err
	}

	// The table should have columns: id, name, price, quantity, category
	return conn, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID

	result, err := ps.db.Exec(`
	INSERT INTO products (name,price,quantity,category) values(?, ?, ?, ?)`,
		product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}

	product.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found

	row := ps.db.QueryRow(`SELECT * FROM products WHERE id=?`, id)

	product := &Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)

	if err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist

	_, err := ps.GetProduct(product.ID)
	if err != nil {
		return err
	}

	_, err = ps.db.Exec(`UPDATE products 
	SET name = ?, price = ?, quantity = ?, category = ?
	WHERE id = ?`, product.Name, product.Price, product.Quantity, product.Category, product.ID)

	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	_, err := ps.GetProduct(id)
	if err != nil {
		return err
	}

	_, err = ps.db.Exec(`DELETE FROM products WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers

	res := make([]*Product, 0)
	sqlQuery := `SELECT * FROM products WHERE category = ?`

	if category == "" {
		sqlQuery = `SELECT * FROM products`
	}

	rows, err := ps.db.Query(sqlQuery, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := &Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
		if err != nil {
			return nil, err
		}

		res = append(res, product)
	}
	return res, nil
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

	for key, value := range updates {
		result, err := tx.Exec(`UPDATE products SET quantity = ? WHERE id = ?`, value, key)
		if err != nil {
			tx.Rollback()
			return err
		}

		rowAffected, err := result.RowsAffected()

		if err != nil || rowAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("Product with id = %d not exist", key)
		}
	}

	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
