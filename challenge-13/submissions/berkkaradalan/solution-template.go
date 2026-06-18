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
	// The table should have columns: id, name, price, quantity, category

	db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
			return nil, err
    }

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price REAL NOT NULL,
        quantity INTEGER NOT NULL,
        category TEXT NOT NULL
    );`)

	if err != nil { 
		return nil, err
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	query := `
		INSERT INTO "products" (name, price, quantity, category) VALUES(?, ?, ?, ?)
	`
	result, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category)

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
	query := `
		SELECT id, name, price, quantity, category FROM products WHERE id = ?
	`
	
	selectedProduct := ps.db.QueryRow(query, id)

	selectedProductStruct := Product{}

	err := selectedProduct.Scan(&selectedProductStruct.ID, &selectedProductStruct.Name, &selectedProductStruct.Price, &selectedProductStruct.Quantity, &selectedProductStruct.Category)
	
	if err != nil {
		return nil, err
	}

	return &selectedProductStruct, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	query := `UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?`

	updateResult, err := ps.db.Exec(query, product.Name, product.Price, product.Quantity, product.Category, product.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := updateResult.RowsAffected()

	if err != nil{
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product doesn't exist")
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	query := `
		DELETE FROM products WHERE id = ?
	`

	deleteResult, err := ps.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := deleteResult.RowsAffected()

	if err != nil{
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product doesn't exist")
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var query string
	var rows *sql.Rows
	var err error

	if category == "" {
		query = `SELECT id, name, price, quantity, category FROM products`
		rows, err = ps.db.Query(query)
	} else {
		query = `SELECT id, name, price, quantity, category FROM products WHERE category = ?`
		rows, err = ps.db.Query(query, category)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := []*Product{}

	for rows.Next() {
		i := Product{}
		err = rows.Scan(&i.ID, &i.Name, &i.Price, &i.Quantity, &i.Category)

		if err != nil {
			return nil, err
		}

		data = append(data, &i)
	}
	
	return data, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	tx, err := ps.db.Begin()

	if err != nil {
		return err
	}

	for id, quantity := range updates {
		result, err := tx.Exec("UPDATE products SET quantity = ? WHERE id = ?", quantity, id)
		if err != nil {
			tx.Rollback()
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}

		if rowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("product not found: id=%d", id)
		}
	}

	return tx.Commit()
}

func main() {
	sqliteDB, err := InitDB("./test.sqlite")

	if err != nil { 
		fmt.Printf("Error :%v", err)
		return
	}
	defer sqliteDB.Close()

	pS := NewProductStore(sqliteDB)

	testProduct := Product{
		Name: "Test",
		Price: 2500,
		Quantity: 2500,
		Category: "TestItem",
	}

	err = pS.CreateProduct(&testProduct)

	if err != nil {
		fmt.Printf("Error :%v", err)
	}

	selectedProduct, err := pS.GetProduct(1)

	if err != nil {
		fmt.Printf("Error :%v", err)
	}

	fmt.Printf("Selected product : %v", selectedProduct)

	err = pS.DeleteProduct(8)

	if err != nil {
		fmt.Printf("Error :%v", err)
	}

	var ListedByCategory []*Product

	ListedByCategory, err = pS.ListProducts("")

	if err != nil {
		fmt.Printf("Error :%v", err)
	}

	for _, product := range ListedByCategory {
		fmt.Printf("Product : %v", product)
	}
}
