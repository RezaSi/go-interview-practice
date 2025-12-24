package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"

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
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")

	if err != nil {
		return nil, err
	}

	// The table should have columns: id, name, price, quantity, category
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	result, err := ps.db.Exec("insert into products (name, price, quantity, category) values(?,?,?,?)",
		product.Name, product.Price, product.Quantity, product.Category)

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

	row := ps.db.QueryRow("select id, name, price, quantity, category from products where id = ?", id)

	product := &Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)

	if err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {

	_, err := ps.db.Exec(
		`update products 
            set name = ?
            , price = ?
            , quantity = ?
            , category = ?
    where id = ?`, product.Name, product.Price, product.Quantity, product.Category, product.ID)

	if err != nil {
		return err
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {

	_, err := ps.db.Exec(`delete from products where id = ?`, id)

	if err != nil {
		return err
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	query := `
        select * 
        from products
        where 1 = 1`

	var args []interface{}

	if len(category) != 0 {
		query += ` and category = ?`
		args = append(args, category)
	}

	stmt, err := ps.db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	result, err := stmt.Query(args...)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	products := []*Product{}

	for result.Next() {
		product := &Product{}
		err := result.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Category)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
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
			tx.Rollback() // Rollback on error
		}
	}()

	for id, quantity := range updates {
		fmt.Printf("quantity: %v, id: %v", quantity, id)

		result, err := tx.Exec("update products set quantity = ? where id = ?", quantity, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("product with ID %d not found", id)
		}
	}

	// if err := tx.Commit(); err != nil {
	// 	return err
	// }

	return tx.Commit()
}

func main() {
	defer os.Remove("inventory.db")

	db, err := InitDB("inventory.db")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize DB: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	store := NewProductStore(db)

	for i := 0; i < 1000; i++ {
		qty := rand.Intn(100)
		price := rand.Float64() * 100

		product := &Product{
			Name:     "Product 1",
			Price:    price,
			Quantity: qty,
			Category: "Test",
		}

		store.CreateProduct(product)
	}

	randomId := rand.Int63n(1000)

	p, _ := store.GetProduct(randomId)

	fmt.Printf("Retrieved Product: %+v\n", p)

	p.Price = 12.99
	store.UpdateProduct(p)
	fmt.Printf("Updated  Product: %+v\n", p)

	p, _ = store.GetProduct(p.ID)
	fmt.Printf("Retrieved Product: %+v\n", p)

	products, _ := store.ListProducts("Test")

	for _, p := range products {
		fmt.Printf("Listed Product: %+v\n", p)
	}
}
