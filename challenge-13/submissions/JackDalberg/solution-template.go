package main

import (
	"database/sql"

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
	   // log.Fatal(err)
	    return nil, err
	}
	if err = db.Ping(); err != nil {
        return nil, err
    }
	tableStmt := `
        CREATE TABLE IF NOT EXISTS 
        Products(id INTEGER PRIMARY KEY AUTOINCREMENT, 
            name TEXT, 
            price REAL, 
            quantity INTEGER, 
            category TEXT);
    `
    _, err = db.Exec(tableStmt)
    if err != nil {
        // log.Printf("%q: %s\n", err, tableStmt)
        return nil, err
    }
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(p *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID
	insertStmt := `INSERT INTO Products(name, price, quantity, category) VALUES (?,?,?,?)`
	res, err := ps.db.Exec(insertStmt, p.Name, p.Price, p.Quantity, p.Category)
	if err != nil {
	    return err
	}
	id, err := res.LastInsertId()
	if err != nil {
	    return err
	}
	p.ID = id
	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	selectStmt := `SELECT id, name, price, quantity, category FROM Products WHERE id = ?`
	row := ps.db.QueryRow(selectStmt, id)
	p := Product{}
    err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
    if err != nil {
        return nil, err
    }
	return &p, nil
	
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(p *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist
	updateStmt := `
	    UPDATE Products
	    SET name = ?, price = ?, quantity = ?, category = ?
	    WHERE id = ?
	`
	_, err := ps.GetProduct(p.ID)
    if err != nil{
        return err
    }
	_, err = ps.db.Exec(updateStmt, p.Name, p.Price, p.Quantity, p.Category, p.ID)
	return err
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	deleteStmt := `DELETE FROM Products WHERE id = ?`
	_, err := ps.db.Exec(deleteStmt, id)
	return err
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	listStmt := `SELECT id, name, price, quantity, category FROM Products`
	if category != ""{
	    listStmt += ` WHERE category = ?`
	}
	rows, err := ps.db.Query(listStmt, category)
	if err != nil{
	    return nil, err
	}
	prods := make([]*Product, 0)
	for rows.Next() {
	    var p Product
	    if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category); err != nil{
	        return nil, err
	    }
	    prods = append(prods, &p)
	}
	return prods, nil
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
	defer tx.Rollback()
	updateStmt := `
	    UPDATE Products
	    SET quantity = ?
	    WHERE id = ?
	`
	for id, v := range updates {
	    _, err = ps.GetProduct(id)
	    if err != nil{
	        return err
	    }
	     _, err = tx.Exec(updateStmt, v, id);
	    if err != nil{
	        return err
	    }
	}
	tx.Commit()
	return nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
