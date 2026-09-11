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
	db, err := sql.Open("sqlite3", "test_inventory.db")

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	sql := `CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY, 
            name TEXT,
            price REAL,
            quantity INTEGER,
            category TEXT
        );`

	db.Exec(sql)

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	println(product.Name, product.Price, product.Quantity, product.Category)
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID
	sql := `INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?);`

	result, err := ps.db.Exec(sql, product.Name, product.Price, product.Quantity, product.Category)
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
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	s := `SELECT * FROM products WHERE id = ?;`
	row := ps.db.QueryRow(s, id)

	p := &Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, err
	}

	return p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist
	stmt, err := ps.db.Prepare("UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}

	return nil
	//return errors.New("not implemented")
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	_, err := ps.db.Exec("DELETE FROM products WHERE id = ?;", id)
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
	var query string
	if len(category) > 0 {
		query = `SELECT * FROM products WHERE category = ?;`
	} else {
		query = `SELECT * FROM products;`
	}

	rows, err := ps.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*Product
	for rows.Next() {
		p := &Product{}
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

	stmt, err := tx.Prepare("UPDATE products SET quantity = ? WHERE id = ?;")
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
			return fmt.Errorf("product with id %d not found", id)
		}
	}
	tx.Commit()
	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}