package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID       int64
	Name     string
	Price    float64
	Quantity int
	Category string
}

type ProductStore struct {
	db *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (
                        id INTEGER PRIMARY KEY AUTOINCREMENT, 
                        name TEXT NOT NULL, 
                        price REAL NOT NULL, 
                        quantity INTEGER NOT NULL, 
                        category TEXT);`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (ps *ProductStore) CreateProduct(product *Product) error {
	result, err := ps.db.Exec(
		"INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
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

func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	row := ps.db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE id = ?", id)
	p := &Product{}

	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, err
	}
	return p, nil
}

func (ps *ProductStore) UpdateProduct(product *Product) error {
	result, err := ps.db.Exec("UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?",
		product.Name, product.Price, product.Quantity, product.Category, product.ID)
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
}

func (ps *ProductStore) DeleteProduct(id int64) error {
	result, err := ps.db.Exec("DELETE FROM products WHERE id = ?", id)
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
	return nil
}

func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var rows *sql.Rows
	var err error
	if len(category) > 0 {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products WHERE category = ?", category)
	} else {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*Product
	for rows.Next() {
		p := &Product{}

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func ListProducts2(db *sql.DB, filters map[string]interface{}) ([]Product, error) {
	query := "SELECT id, name, price, quantity, category FROM products"
	args := []interface{}{}
	conditions := []string{}
	// Add filters dynamically
	if name, ok := filters["name"]; ok {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+name.(string)+"%")
	}
	if minPrice, ok := filters["min_price"]; ok {
		conditions = append(conditions, "price >= ?")
		args = append(args, minPrice)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY id DESC"
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price,
			&product.Quantity, &product.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	if ps.db == nil {
		return fmt.Errorf("database is not initialized")
	}
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()
	stmt, err := tx.Prepare("UPDATE products SET quantity = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for id, quantity := range updates {
		result, execErr := stmt.Exec(quantity, id)
		if execErr != nil {
			return execErr
		}
		rowsAffected, raErr := result.RowsAffected()
		if raErr != nil {
			return raErr
		}
		if rowsAffected == 0 {
			return fmt.Errorf("product with ID %d not found", id)
		}
	}
	err = tx.Commit()
	if err == nil {
		committed = true
	}
	return err

}

func BulkUpdatePrices2(db *sql.DB, updates map[int64]float64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()
	query := "UPDATE products SET price = ? WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	for id, price := range updates {

		result, err := stmt.Exec(price, id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update price for product %d: %w", id, err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get rows affected for product %d: %w", id, err)
		}
		if rowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("product with ID %d not found", id)
		}
	}
	err = tx.Commit()
	if err == nil {
		committed = true
	}
	return err
}

func (ps *ProductStore) ShowData() {
	products, err := ps.ListProducts("")
	if err != nil {
		fmt.Println("Error reading the data", err)
		return
	}
	fmt.Println("products:")
	for _, prd := range products {
		fmt.Println(prd)
	}
}
func main() {
	db, err := InitDB("database.db")
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}
	defer db.Close()
	ps := NewProductStore(db)

	err = ps.CreateProduct(&Product{0, "Cacao", 10.5, 3, "junk"})
	if err != nil {
		fmt.Println("Failed to create product:", err)
		return
	}
	err = ps.CreateProduct(&Product{0, "milk", 20.5, 2, "dairy"})
	if err != nil {
		fmt.Println(err)
		return
	}
	ps.ShowData()

	if err := ps.BatchUpdateInventory(map[int64]int{1: 100, 2: 100}); err != nil {
		fmt.Println("Failed to batch update:", err)
		return
	}
	ps.ShowData()
	prd, err := ps.GetProduct(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(prd)

	if err := ps.UpdateProduct(&Product{2, "new milk", 30.0, 150, "new dairy"}); err != nil {
		fmt.Println("Failed to update product:", err)
		return
	}
	ps.ShowData()

	if err := ps.DeleteProduct(1); err != nil {
		fmt.Println("Failed to delete product:", err)
		return
	}
	ps.ShowData()
}