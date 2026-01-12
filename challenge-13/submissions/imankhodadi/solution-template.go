package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID       int64
	Name     string
	Price    float64
	Quantity int
	Category string //sql.NullString // Can be NULL
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
		fmt.Println("aaaaaaaaaaaaaaaaaaa")
		return nil, err
	}
	if err = db.Ping(); err != nil {
		fmt.Println("bbbbbbbbbbbbbbbbbbbbbbb")
		return nil, err
	}
	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)")
	if err != nil {
		fmt.Println("ccccccccccccccccccccccccccc")
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
	// var category sql.NullString
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, err
	}
	// if category.Valid {
	// 	p.Category = category
	// } //else {
	// 	p.Category = "Null"
	// }
	return p, nil
}

func (ps *ProductStore) UpdateProduct(product *Product) error {
	result, err := ps.db.Exec("UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?",
		product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", product.ID)
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
	// var category sql.NullString
	for rows.Next() {
		p := &Product{}

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		// if category.Valid {
		//     fmt.Println(category.String)
		// } //else {
		//     fmt.Println("Category is NULL")
		// }
		products = append(products, p)
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
	// Ensure rollback on any early return
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	stmt, err := tx.Prepare("UPDATE products SET quantity = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for id, quantity := range updates {
		result, err := stmt.Exec(quantity, id)
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
	return tx.Commit()
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
	fmt.Println("-------------------")
}
func main() {
	db, err := InitDB("database.db")
	if err != nil {
		fmt.Println("------------", err)
		return
	}
	ps := NewProductStore(db)

	err = ps.CreateProduct(&Product{0, "Cacao", 10.5, 3, "junk"})
	if err != nil {

		fmt.Println("DDDDDDDDDDDDD", err)
		return
	}
	err = ps.CreateProduct(&Product{0, "milk", 20.5, 2, "dairy"})
	if err != nil {
		fmt.Println(err)
		return
	}
	ps.ShowData()
	ps.BatchUpdateInventory(map[int64]int{1: 100, 2: 100})
	ps.ShowData()
	prd, err := ps.GetProduct(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(prd)

	ps.UpdateProduct(&Product{2, "new milk", 30.0, 150, "new dairy"})
	ps.ShowData()

	ps.DeleteProduct(1)
	ps.ShowData()

}
