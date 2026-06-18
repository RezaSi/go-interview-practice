package main

import (
	"database/sql"
	"context"
	"time"
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
	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	if err = db.PingContext(ctx); err != nil {
	    return nil, err
	}
	
	query := `CREATE TABLE IF NOT EXISTS products (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name VARCHAR(255),
	    price DOUBLE,
	    quantity INT,
	    category VARCHAR(255)
	    );`
	_, err = db.Exec(query)
	if err != nil {
	    return nil, err
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID
	query := "INSERT INTO products(name, price, quantity, category) VALUES ($1, $2, $3, $4) RETURNING id"
	err := ps.db.QueryRow(query, product.Name, product.Price, product.Quantity, product.Category).Scan(&product.ID)
	return err
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	query := "SELECT * from products WHERE id = $1"
	var p Product
	err := ps.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
	    return nil, err
	}
	return &p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist
	query := "UPDATE products SET name=$2, price=$3, quantity=$4, category=$5 WHERE id = $1"
	_, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category, product.ID)
	return err
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	query := "DELETE FROM products WHERE id = $1"
	_, err := ps.db.Exec(query, id)
	return err
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	query := "SELECT * from products WHERE category = $1"
	if category == "" {
	   query = "SELECT * from products" 
	}
	var res []*Product
	rows, err := ps.db.Query(query, category)
	if err != nil {
	    return nil, err
	}
	
	for rows.Next() {
	    p := &Product{}
	    err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	    if err != nil {
	        return nil, err
	    }
	    
	    res = append(res, p)
	}
	
	if err = rows.Err(); err != nil {
	    return nil, err
	}
	return res, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction
	tx, _ := ps.db.BeginTx(context.Background(), nil)
	defer tx.Rollback()
	for id := range updates {
	    query := "UPDATE products SET quantity=$1 WHERE id = $2"
    	res, err := tx.Exec(query, updates[id], id)
    	if err != nil {
    	    tx.Rollback()
    	    return err
    	}
    	
    	rows, err := res.RowsAffected()
    	if err != nil {
    	    tx.Rollback()
    	    return err
    	}
    	
    	if rows == 0 {
    	    tx.Rollback()
    	    return errors.New("product not found")
    	}
	}
	tx.Commit()
	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
