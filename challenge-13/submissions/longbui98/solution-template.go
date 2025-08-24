package main

import (
	"database/sql"
	"errors"
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

	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS `products`" +
			" (id INTEGER PRIMARY KEY AUTO_INCREMENT, name TEXT, price REAL, quantity INTEGER, category TEXT) ",
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID
	_, err := ps.db.Exec(
		"INSERT INTO products(name, price, quantity, category) VALUES(?, ?, ?, ?)",
		product.ID, product.Name, product.Quantity, product.Category)
	if err != nil {
		return err
	}
	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	rows, err := ps.db.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	product := &Product{}
	err = rows.Scan(product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)
	if err != nil {
		return nil, err
	}
	if product.ID == 0 {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist
	if product.ID == 0 {
		return errors.New("id is null")
	}
	result, err := ps.db.Exec(
		"UPDATE product SET name = ?, quantity = ?, category = ? WHERE id = ?",
		product.Name, product.Quantity, product.Category, product.ID,
	)
	if err != nil {
		return err
	}
	effectedRow, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if effectedRow == 0 {
		return errors.New("product not found")
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	if id == 0 {
		return errors.New("id is null")
	}
	result, err := ps.db.Exec(
		"DELETE FROM product WHERE id = ?)", id,
	)
	if err != nil {
		return err
	}
	effectedRow, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if effectedRow == 0 {
		return errors.New("product not found")
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	if category == "" {
		return nil, errors.New("category is empty")
	}
	rows, err := ps.db.Query("SELECT * product WHERE category = ?", category)
	if err != nil {
		return nil, err
	}
	var product []*Product
	for rows.Next() {
		prod := &Product{}
		err = rows.Scan(prod.ID, prod.Name, prod.Category, prod.Price)
		if err != nil {
			return nil, err
		}
		product = append(product, prod)
	}
	if len(product) == 0 {
		return nil, errors.New("no products found")
	}
	return product, nil
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

	stmt, err := tx.Prepare("UPDATE product SET quantity = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, quantity := range updates {
		result, err := stmt.Exec(quantity, id)
		if err != nil {
			return err
		}
		effectedRow, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if effectedRow == 0 {
			return errors.New("product not found")
		}
	}
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
