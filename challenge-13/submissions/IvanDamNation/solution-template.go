package main

import (
	"database/sql"
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

type ProductNotFoundError struct {
    ID int64
}

func (e *ProductNotFoundError) Error() string {
    return fmt.Sprintf("product with ID %d not found", e.ID)
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
	
	if err = db.Ping(); err != nil {
	    db.Close()
	    return nil, err
	}
	
	createStmt := `
	CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	price REAL NOT NULL CHECK (price >= 0),
	quantity INTEGER NOT NULL CHECK (quantity >= 0),
	category TEXT NOT NULL
	);`
	
	_, err = db.Exec(createStmt)
	if err != nil {
	    db.Close()
	    return nil, err
	}
	
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	createProdSQL := `
	INSERT INTO products (
	    name, price, quantity, category
	) VALUES (?, ?, ?, ?);
	`
	
	result, err := ps.db.Exec(
	    createProdSQL, 
	    product.Name, product.Price, product.Quantity, product.Category,
	)
	
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
    getProductSQL := `
    SELECT id, name, price, quantity, category 
    FROM products 
    WHERE id = ?;
    `
    row := ps.db.QueryRow(getProductSQL, id)
    
    prod := &Product{}
    err := row.Scan(
        &prod.ID, &prod.Name, &prod.Price, &prod.Quantity, &prod.Category,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, &ProductNotFoundError{ID: id}
        }
        return nil, err
    }
	return prod, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	updateProductSQL := `
	UPDATE products
	SET name = ?, price = ?, quantity = ?, category = ?
	WHERE id = ?;
	`
	
	result, err := ps.db.Exec(
	    updateProductSQL,
	    product.Name, product.Price, product.Quantity, product.Category, product.ID,
	)
	if err != nil {
	    return err
	}
	
	num, err := result.RowsAffected()
	if err != nil {
	    return err
	}
	
	if num == 0 {
	    return &ProductNotFoundError{ID: product.ID}
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
    deleteProductSQL := `
    DELETE FROM products
    WHERE id = ?;
    `
    
    result, err := ps.db.Exec(deleteProductSQL, id)
    if err != nil {
        return err
    }
    
    num, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if num == 0 {
        return &ProductNotFoundError{ID: id}
    }
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
    listProductsSQL := `
    SELECT id, name, price, quantity, category
    FROM products
    `
    
    var args []any
    if category != "" {
        listProductsSQL += `WHERE category = ?`
        args = append(args, category)
    }
    
    rows, err := ps.db.Query(listProductsSQL, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var productsList []*Product
    for rows.Next() {
        p := &Product{}
        err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
        if err != nil {
            return nil, err
        }
        productsList = append(productsList, p)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
	return productsList, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
    batchUpdateInventorySQL := `
    UPDATE products
    SET quantity = ?
    WHERE id = ?;
    `
    
    tx, err := ps.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.Prepare(batchUpdateInventorySQL)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for id, quantity := range updates {
        res, err := stmt.Exec(quantity, id)
        if err != nil {
            return err
        }
        
        numRows, err := res.RowsAffected()
        if err != nil {
            return err
        }
        
        if numRows == 0 {
            return &ProductNotFoundError{ID: id}
        }
    }
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}
