package main

import (
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

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
    
    // Check connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    createProductTable := `
    CREATE TABLE IF NOT EXISTS product (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        price REAL,
        quantity INTEGER,
        category TEXT
    );`
    
    if _, err := db.Exec(createProductTable); err != nil {
        return nil, err
    }
    
    return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
    tx, err := ps.db.Begin()
	if err != nil {
        return fail("CreateProduct", err)
    }
    // Defer a rollback in case anything fails. If the transaction is committed
    // successfully, this rollback will run as a no-op.
    defer tx.Rollback()
    
    createProduct := `
    INSERT INTO product (name, price, quantity, category)
    VALUES (?, ?, ?, ?);
    `
    res, err := tx.Exec(createProduct, product.Name, product.Price, product.Quantity, product.Category);
    if err != nil {
        return fail("CreateProduct", err)
    }
    id, err := res.LastInsertId()
    if err != nil {
        return fail("CreateProduct", err)
    }
    if err = tx.Commit(); err != nil {
        return fail("CreateProduct", err)
    }

    product.ID = id
	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	selectProduct := `
	SELECT id, name, price, quantity, category
	FROM product
	WHERE id=?;
	`
	var p Product
	if err := ps.db.QueryRow(selectProduct, id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category); err != nil {
        if err == sql.ErrNoRows {
            return nil, fail(fmt.Sprintf("GetProduct, unknown product ID %v", id), err)
        }
        return nil, fail("GetProduct", err)
    }
	
	return &p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	tx, err := ps.db.Begin()
	if err != nil {
        return fail("UpdateProduct", err)
    }
    // Defer a rollback in case anything fails. If the transaction is committed
    // successfully, this rollback will run as a no-op.
    defer tx.Rollback()
    
    exists, err := productExists(product.ID, tx)
    if err != nil || !exists {
        return fail("UpdateProduct, product does not exist", err)
    }
    
    updateProduct := `
    UPDATE product
    SET name=?, price=?, quantity=?, category=?
    WHERE id=?;
    `
    if _, err := tx.Exec(updateProduct, product.Name, product.Price, product.Quantity, product.Category, product.ID); err != nil {
        return fail("UpdateProduct", err)
    }
    if err = tx.Commit(); err != nil {
        return fail("UpdateProduct", err)
    }
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	tx, err := ps.db.Begin()
	if err != nil {
        return fail("DeleteProduct", err)
    }
    // Defer a rollback in case anything fails. If the transaction is committed
    // successfully, this rollback will run as a no-op.
    defer tx.Rollback()
    
    exists, err := productExists(id, tx)
    if err != nil || !exists {
        return fail("DeleteProduct, product does not exist", err)
    }
    
    deleteProduct := `
    DELETE FROM product
    WHERE id=?;
    `
    if _, err := tx.Exec(deleteProduct, id); err != nil {
        return fail("DeleteProduct", err)
    }
    if err = tx.Commit(); err != nil {
        return fail("DeleteProduct", err)
    }
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var rows *sql.Rows
	var err error
	selectProducts := `
	SELECT id, name, price, quantity, category
	FROM product
	`
    if category == "" {
    	selectProducts += ";"
	    if rows, err = ps.db.Query(selectProducts); err != nil {
            return nil, fail("ListProducts, no products stored", err)
        }
    } else {
    	selectProducts += "WHERE category=?;"
    	if rows, err = ps.db.Query(selectProducts, category); err != nil {
            return nil, fail(fmt.Sprintf("ListProducts, no products with category %v", category), err)
        }
    }
    defer rows.Close()
    
    products := []Product{}
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category); err != nil {
            return nil, fail("ListProducts", err)
        }
        products = append(products, p)
    }
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, fail("ListProducts", err)
	}
	
	pointers := []*Product{}
	for _, p := range products {
	    pointers = append(pointers, &p)
	}
	return pointers, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
    if len(updates) == 0 {
        return errors.New("BatchUpdateInventory, no updates to execute")
    }
    
	tx, err := ps.db.Begin()
	if err != nil {
        return fail("BatchUpdateInventory", err)
    }
    // Defer a rollback in case anything fails. If the transaction is committed
    // successfully, this rollback will run as a no-op.
    defer tx.Rollback()
    
    // Fail fast if one product ID does not exist in db.
    idsInt64 := slices.Collect(maps.Keys(updates))
    idsStr := []string{}
    for _, v := range idsInt64 {
        // Must stringify int64 value
        idsStr = append(idsStr, strconv.FormatInt(v, 10))
    }
    allExist, err := allProductsExist(idsStr, tx)
    if err != nil || !allExist {
        return fail("BatchUpdateInventory", err)
    }
    
    // Stringify updates into SQL-friendly format
    idsAndQuantities := ""
    total := len(updates)
    count := 1
    for id, q := range updates {
        idsAndQuantities += fmt.Sprintf("(%v, %v)", id, q)
        // All values except last should have trailing comma.
        if count < total {
            idsAndQuantities += ","
            count++
        }
    }

    batchUpdateInventory := `
    WITH updates(id, new_quantity) AS (
      VALUES ` +
        idsAndQuantities + `
    )
    UPDATE product
    SET quantity=(
        SELECT new_quantity
        FROM updates
        WHERE updates.id = product.id
    )
    WHERE id IN (
        SELECT id FROM updates
    );
    `
    if _, err := tx.Exec(batchUpdateInventory); err != nil {
        return fail("BatchUpdateInventory", err)
    }
    if err = tx.Commit(); err != nil {
        return fail("BatchUpdateInventory", err)
    }
	return nil
}

// Create a helper function for preparing failure results.
func fail(method string, err error) error {
    return fmt.Errorf("%v: %v", method, err)
}

func productExists(id int64, tx *sql.Tx) (bool, error) {
    existsStm := `
    SELECT EXISTS(
        SELECT 1
        FROM product
        WHERE id=?
    );
    `
    var exists bool
    if err := tx.QueryRow(existsStm, id).Scan(&exists); err != nil {
        return false, fail(fmt.Sprintf("Product with ID %v does not exist", id), err)
    }
    return exists, nil
}

func allProductsExist(ids []string, tx *sql.Tx) (bool, error) {
    total := len(ids)
    existsStm := `
    SELECT CASE 
        WHEN COUNT(DISTINCT id) = ` + 
        strconv.Itoa(total) + `
        THEN 1 
        ELSE 0 
    END AS all_exist
    FROM product
    WHERE id IN (` + 
    strings.Join(ids, ",") + 
    `);
    `
    var allExist bool
    if err := tx.QueryRow(existsStm).Scan(&allExist); err != nil {
        return false, fail("At least one product does not exist", err)
    }
    return allExist, nil
}

func main() {
	// Optional: you can write code here to test your implementation
}
