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
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.New("could not connect to db")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products(
            id  INTEGER PRIMARY KEY AUTOINCREMENT,
            name VARCHAR,
            price REAL,
            quantity INTEGER,
            category VARCHAR )`)

	if err != nil {
		return nil, errors.New("could not create products table")
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	query := `INSERT INTO products (name, price, quantity, category)
              VALUES (?, ?, ?, ?)`

	rslt, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}

	id, err := rslt.LastInsertId()
	if err != nil {
		return err
	}
    
	product.ID = id 
	return nil
}

// GetProduct retrieves a product by ID
    func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	row := ps.db.QueryRow(`SELECT * FROM products WHERE id = ?`, id)

	var product Product
	err := row.Scan(&product.ID,
		&product.Name,
		&product.Price,
		&product.Quantity,
		&product.Category,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("could not found")
		}
		return nil, err
	}

	return &product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	query := `UPDATE products
             SET name = ?, price = ?, quantity = ?, category = ?
             WHERE id = ?`

	result, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("product does no exist")
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	result, err := ps.db.Exec(`DELETE FROM products WHERE id = ?`, id)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("product does no exist")
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var rows *sql.Rows
	var err error

	query := `SELECT * FROM products`
	if category != "" {
		query += ` WHERE category = ? ;`
		rows, err = ps.db.Query(query, category)
	} else {
		rows, err = ps.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category); err != nil {
			return nil, err
		}

		products = append(products, &product)
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
			return errors.New("could not found one of the products")
		}
	}

	return tx.Commit()
}


func main() {
}


