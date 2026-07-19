package main

import (
	"database/sql"
	//"errors"
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
	db, err :=sql.Open("sqlite3",dbPath)
	if err != nil{
	    return nil,err
	}
	
	if err = db.Ping(); err != nil {
        return nil, err
    }
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	_,err = db.Exec("CREATE TABLE IF NOT EXISTS products(id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")
	   
	  if err != nil{
	      return nil, err
	  }
	
	return db,err
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	query := `INSERT INTO products(name,price,quantity,category)
	          VALUES(?,?,?,?);`
	 
	result, err := ps.db.Exec(query,product.Name,product.Price,product.Quantity,product.Category)
	if err != nil{
	    return err
	}
	// TODO: Update the product.ID with the database-generated ID
	id, err := result.LastInsertId()
	if err != nil{
	    return err
	}
	product.ID = id
    return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	query := `
	SELECT *
	FROM products
	WHERE id = ?
	`
	var product Product
	err := ps.db.QueryRow(query,id).Scan(
	        &product.ID,
	        &product.Name,
	        &product.Price,
	        &product.Quantity,
	        &product.Category,
	    )
	 if err!= nil{
	     if err == sql.ErrNoRows {
            return nil, fmt.Errorf("product with ID %d not found", id)
        }
        return nil, fmt.Errorf("failed to get product: %w", err)
	 }
	// TODO: Return a Product struct populated with the data or an error if not found
	return &product,nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
    _,err := ps.GetProduct(product.ID)
    if err != nil{
        return err
    }
    
    query:= `
    UPDATE products
    SET name = ?,price = ?,quantity = ?,category=?
    WHERE id = ?
    `
    result,err:= ps.db.Exec(query,product.Name,product.Price,product.Quantity,product.Category,product.ID)
    
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
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	query := `DELETE FROM products WHERE id = ?`
	result , err := ps.db.Exec(query,id)
	if err != nil{
	    return err
	}
	
	rowsAffected , err := result.RowsAffected()
	if err != nil{
	    return err
	}
	
	if rowsAffected == 0{
	    return fmt.Errorf("no rows were deleted")
	}
	// TODO: Return an error if the product doesn't exist
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	query:= `SELECT *
	FROM products
	`
	if category != ""{
	    query = query + `WHERE category = ?`
	}
	
	rows, err := ps.db.Query(query,category)
	if err != nil{
	    return nil , err
	}
	
	var product []*Product
	
	for rows.Next(){
	    p:= &Product{}
	    err := rows.Scan(&p.ID,&p.Name,&p.Price,&p.Quantity,&p.Category)
	    if err!=nil{
	        return nil, err
	    }
	    product = append(product,p)
	}
	return product,nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction
	tx , err := ps.db.Begin()
	if err != nil{
	    return err
	}
	
	defer func(){
	    if r:= recover(); r!= nil{
	        tx.Rollback()
	        panic(r)
	    }
	}()
	
	query := `UPDATE products SET quantity = ? WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err!= nil{
	   tx.Rollback()
	   return err
	}
	
	defer stmt.Close()
	
	for id, qty := range updates{
	    result,err := stmt.Exec(qty,id)
	    if err != nil{
	        tx.Rollback()
	        return err
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
	
	if err:= tx.Commit();err!=nil{
	    return err
	}
	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
