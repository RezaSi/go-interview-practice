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
	// TODO: Open a SQLite database connection
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    if err = db.Ping(); err != nil {
        return nil, err
    }
    
    if _, err := db.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)");err != nil {
        return nil,err
    }
    return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	db:= ps.db
	result, err := db.Exec(
    "INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",product.Name, product.Price, product.Quantity, product.Category)
    if err != nil {
        return err
    }

    // Get the ID of the inserted row
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    product.ID = id
    return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
    db:= ps.db
    row:= db.QueryRow("SELECT id, name, price, quantity, category  FROM products WHERE id = ?",id)
    p:= &Product{}
    err:= row.Scan(&p.ID,&p.Name,&p.Price,&p.Quantity,&p.Category)
    if err!=nil{
        return nil,err
    }
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	return p,nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist
	db:=ps.db
	_,err:=db.Exec("UPDATE products SET name=?,price=?,quantity=?,category=? WHERE id = ?",product.Name,product.Price,product.Quantity,product.Category,product.ID)
	if err!=nil{
	    return err
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	db:=ps.db
	_,err:=db.Exec("DELETE FROM products WHERE id = ?",id)
	if err!=nil{
	    return err
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	var allProducts []*Product
	db:= ps.db
	query:= "SELECT id, name, price, quantity, category FROM products"
	rows, err := (*sql.Rows)(nil), error(nil)

	if category !=""{
	    query+=" WHERE category= ?"
	    rows,err=db.Query(query,category)
	}else{
	    rows,err=db.Query(query)
	}
	if err!=nil{
	    return allProducts,err
	}
	defer rows.Close()
	
	
	for rows.Next(){
	    p:=&Product{}
	    if err:=rows.Scan(&p.ID,&p.Name,&p.Price,&p.Quantity,&p.Category);err!=nil{
	        return allProducts,err
	    }
	    allProducts = append(allProducts,p)
	    
	}
	if err:=rows.Err();err!=nil{
	    return allProducts,err
	}
	return allProducts,nil
	

}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction
	db:= ps.db
	tx,err:= db.Begin()
	if err!=nil{
	    return err
	}
	defer func(){
	   if err!=nil{
	       tx.Rollback()
	   } 
	}()
	
	query:= "UPDATE products SET quantity = ? WHERE id = ?"
	for id,quantity := range updates{
	    res,err:= tx.Exec(query,quantity,id)
	    if err!=nil{
	        return err
	    }
	    if rowsAffected,err:= res.RowsAffected();err!=nil{
	        return err
	    }else{
	        if rowsAffected ==0{
	            return errors.New("No row found")
	        }
	    }
	}
	
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
