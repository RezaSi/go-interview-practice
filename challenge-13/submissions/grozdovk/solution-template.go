package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
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
	db, err := sql.Open("sqlite", dbPath)
	if err!=nil{
	    return nil, err
	}
	err = db.Ping()
	if err!=nil{
	    return nil, err
	}
	// TODO: Create the products table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")
    if err!=nil{
        return nil,err
    }
	// The table should have columns: id, name, price, quantity, category
	return db ,err
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	result, err:= ps.db.Exec("INSERT INTO products (name, price, quantity, category) VALUES (?,?,?,?)", 
	product.Name, product.Price, product.Quantity, product.Category)
	if err!= nil {
	    return err
	}
	// TODO: Update the product.ID with the database-generated ID
	id, err := result.LastInsertId()
	if err!= nil{
	    return err
	}
	product.ID = id 
	return err
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	row:= ps.db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE ID = ?", id)
	p:= &Product{}
	err:= row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err!= nil {
	     if err== sql.ErrNoRows{
	         return &Product{} ,fmt.Errorf("product not found")
	     }
	     return &Product{}, err
	}
	// TODO: Return a Product struct populated with the data or an error if not found
	return p, err
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	result, err:= ps.db.Exec("UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE ID = ?",
	product.Name, product.Price,product.Quantity, product.Category, product.ID)
	if err!=nil{
	    return err
	}
	affRows, err:= result.RowsAffected()
	if err!=nil{
	    return err
	}
	if affRows==0{
	    return errors.New("product not found")
	}
	// TODO: Return an error if the product doesn't exist
	return err
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	result, err:= ps.db.Exec("DELETE FROM products WHERE ID = ?", id)
	if err!=nil{
	    return err
	}
	deletedRows, err:= result.RowsAffected()
	if err!=nil{
	    return err
	}
	if deletedRows==0{
	    return errors.New("product not found")
	}
	// TODO: Return an error if the product doesn't exist
	return err
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	var rows *sql.Rows
	var err error
	if category ==""{
	    rows,err = ps.db.Query("SELECT id, name, price, quantity, category FROM products")
	} else{
	    rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products WHERE category = ?", category)
    }
    if err!=nil{
        return nil,err
    }
    // TODO: If category is not empty, filter by category
	defer rows.Close()
	var products []*Product
	for rows.Next(){
	    p:= &Product{}
	    err:= rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	    if err!=nil{
	        return nil,err
	    }
	    products = append(products, p)
	} 
	if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }
	// TODO: Return a slice of Product pointers
	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
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
    if err!=nil{
        return err
    }
    defer stmt.Close()
	// TODO: For each product ID in the updates map, update its quantity
	for id, quantity := range updates {
	    result,err := stmt.Exec(quantity, id)
	    if err!= nil {
	        return err
	    }
	    rowsAffected,err:= result.RowsAffected()
	    if err!= nil{
	        return err
	    }
	    if rowsAffected==0{
	        return fmt.Errorf("Product with id %d not found", id)
	    }
	}
	// TODO: If any update fails, roll back the transaction
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
