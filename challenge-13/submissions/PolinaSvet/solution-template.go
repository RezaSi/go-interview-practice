package main

import (
	"database/sql"
	"fmt"
	"os"

	//_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
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
	// The table should have columns: id, name, price, quantity, category

	os.Remove(dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS products (
        id       INTEGER PRIMARY KEY AUTOINCREMENT,
        name     TEXT NOT NULL UNIQUE,
        price    REAL NOT NULL CHECK (price >= 0),
        quantity INTEGER NOT NULL CHECK (quantity >= 0),
		category TEXT NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	return db, nil

}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID

	query := `
    INSERT INTO products (name, price, quantity, category) 
    VALUES (?, ?, ?, ?)`
	result, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	product.ID = id
	//fmt.Println(product, id)
	return nil

}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found

	query := `
    SELECT id, name, price, quantity, category 
    FROM products 
    WHERE id = ?`
	var product Product
	err := ps.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Quantity,
		&product.Category,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &product, nil

}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist

	query := `
    UPDATE products 
    SET name = ?, price = ?, quantity = ?, category = ? 
    WHERE id = ?`
	result, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}
	return nil

}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist

	query := `DELETE FROM products WHERE id = ?`
	result, err := ps.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers

	query := `SELECT id, name, price, quantity, category FROM products WHERE category LIKE COALESCE(NULLIF(?, ''), '%')`
	rows, err := ps.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()
	var products []*Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price,
			&product.Quantity, &product.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
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
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Use defer to handle rollback on panic or error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	query := "UPDATE products SET quantity = ? WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	for id, quantity := range updates {
		result, err := stmt.Exec(quantity, id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update quantity for product %d: %w", id, err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("no rows were updated")
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil

}

func main() {
	// InitDB sets up a new SQLite database and creates the products table
	db, err := InitDB("inventory.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	// NewProductStore creates a new ProductStore with the given database connection
	store := NewProductStore(db)

	// CreateProduct adds a new product to the database
	if err = store.CreateProduct(&Product{Name: "Free Product1", Price: 0, Quantity: 10, Category: "Free"}); err != nil {
		fmt.Println(err)
		return
	}
	if err = store.CreateProduct(&Product{Name: "Free Product2", Price: 0, Quantity: 10, Category: "Free"}); err != nil {
		fmt.Println(err)
		return
	}

	// GetProduct retrieves a product by ID
	product, err := store.GetProduct(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(product)

	// UpdateProduct updates an existing product
	err = store.UpdateProduct(&Product{ID: 1, Name: "Free Product1", Price: 1, Quantity: 10, Category: "Free"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// DeleteProduct removes a product by ID
	err = store.DeleteProduct(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ListProducts returns all products with optional filtering by category
	products, err := store.ListProducts("Free")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(products))

	// BatchUpdateInventory updates the quantity of multiple products in a single transaction
	err = store.BatchUpdateInventory(map[int64]int{2: 100})
	if err != nil {
		fmt.Println(err)
		return
	}

	product, err = store.GetProduct(2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(product)

}
