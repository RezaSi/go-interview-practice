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
	DBConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// The table should have columns: id, name, price, quantity, category
	tableCreationSQL := `
	CREATE TABLE IF NOT EXISTS products (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT NOT NULL,
    	price REAL NOT NULL,
    	quantity INTEGER DEFAULT 0,
    	category TEXT NOT NULL
	);`
	_, err = DBConn.Exec(tableCreationSQL)
	if err != nil {
		return nil, err
	}

	return DBConn, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	insertSQL, err := ps.db.Prepare("INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?);")
	if err != nil {
		return err
	}

	result, e := insertSQL.Exec(product.Name, product.Price, product.Quantity, product.Category)
	if e != nil {
		return err
	}

	productID, getIDErr := result.LastInsertId()
	if getIDErr != nil {
		return getIDErr
	}
	product.ID = productID

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	fetchSQL, err := ps.db.Prepare("SELECT name, price, quantity, category FROM products WHERE id = ?;")
	if err != nil {
		return nil, err
	}

	result := fetchSQL.QueryRow(id)

	theProduct := Product{
		ID: id,
	}
	err = result.Scan(&theProduct.Name, &theProduct.Price, &theProduct.Quantity, &theProduct.Category)
	if err != nil {
		return nil, err
	}

	return &theProduct, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	updateSQL, e := ps.db.Prepare("UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?;")
	if e != nil {
		return e
	}

	_, err := updateSQL.Exec(product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	deleteSQL, e := ps.db.Prepare("DELETE FROM products WHERE id = ?;")
	if e != nil {
		return e
	}

	_, err := deleteSQL.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	pattern := "%" + category + "%"
	fetchSQL, err := ps.db.Prepare("SELECT id FROM products WHERE category LIKE ?;")
	if err != nil {
		return nil, err
	}

	rows, e := fetchSQL.Query(pattern)
	if e != nil {
		return nil, e
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	result := make([]*Product, 0)
	var product *Product
	var ID int64
	for rows.Next() {
		err = rows.Scan(&ID)
		if err != nil {
			return nil, err
		}
		product, e = ps.GetProduct(ID)
		if e != nil {
			return nil, e
		}
		result = append(result, product)
	}

	return result, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err = tx.Rollback()
			if err != nil {
				return
			}
			panic(r)
		} else if err != nil {
			err = tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	updateSQL, e := tx.Prepare("UPDATE products SET quantity = ? WHERE id = ?;")
	if e != nil {
		return e
	}
	defer func(updateSQL *sql.Stmt) {
		err = updateSQL.Close()
		if err != nil {
			panic(err)
		}
	}(updateSQL)

	for ID, newQuantity := range updates {
		result, updateErr := updateSQL.Exec(newQuantity, ID)
		if updateErr != nil {
			return updateErr
		}

		rows, getRowsError := result.RowsAffected()
		if getRowsError != nil {
			return getRowsError
		} else if rows == 0 {
			return errors.New("no product with this id exists")
		}
	}

	e = tx.Commit()
	if e != nil {
		return e
	}

	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
