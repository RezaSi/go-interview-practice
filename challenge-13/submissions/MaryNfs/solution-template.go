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
	// Open a SQLite database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS products(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, price REAL, quantity INTEGER, category TEXT)")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	res, err := ps.db.Exec("INSERT INTO products(name, price, quantity, category) VALUES(?,?,?,?)", product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}
	product.ID, err = res.LastInsertId()
	if err !=nil {
		return err
	}
	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// Query the database for a product with the given ID
	res := ps.db.QueryRow("SELECT * FROM products WHERE id = ?", id)

	p := &Product{}
	err := res.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, err
	}
	return p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// Update the product in the database
	res, err := ps.db.Exec("UPDATE products SET name=?, price=?, quantity=?, category=? WHERE id=?", product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("Product not found")
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	res, err := ps.db.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("Product not found")
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	//Query the database for products
	query := "SELECT id, name, price, quantity, category FROM products"
	if category != "" {
		query += " WHERE category = ?"
	}
	rows, err := ps.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		p := &Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// TODO: If category is not empty, filter by category
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
		res, err := stmt.Exec(quantity, id)
		if err != nil {
			return err
		}
		row, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if row != 1 {
			return fmt.Errorf("product with id %d not found", id)
		}
	}

	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
	db, err := InitDB("db.db")
	if err != nil {
		fmt.Println(err)
	}
	p := Product{
		Name:     "Nike",
		Price:    200,
		Quantity: 3,
		Category: "Shoe",
	}
	ps := NewProductStore(db)
	ps.CreateProduct(&p)
	fmt.Println(ps.GetProduct(1))
	p2 := Product{
		ID:       3,
		Name:     "Yonex",
		Price:    660,
		Quantity: 2,
		Category: "Shoe",
	}
	fmt.Println(ps.UpdateProduct(&p2))
	fmt.Println(ps.DeleteProduct(2))
	res, err := ps.ListProducts("Shoe")
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range res {
		fmt.Println(val)
	}
}
