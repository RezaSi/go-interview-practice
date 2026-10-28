package main

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MigrationVersion tracks the current database schema version
type MigrationVersion struct {
	ID        uint `gorm:"primaryKey"`
	Version   int  `gorm:"unique;not null"`
	AppliedAt time.Time
}

// Product represents a product in the e-commerce system (Version 3)
type Product struct {
	ID          uint     `gorm:"primaryKey"`
	Name        string   `gorm:"not null"`
	Price       float64  `gorm:"not null"`
	Description string   `gorm:"type:text"`
	CategoryID  uint     `gorm:"not null"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Stock       int      `gorm:"default:0"`
	SKU         string   `gorm:"unique;not null"`
	IsActive    bool     `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Category represents a product category
type Category struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"type:text"`
	Products    []Product `gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductV1 struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Description string  `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ConnectDB establishes a connection to the SQLite database
func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&MigrationVersion{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// RunMigration runs a specific migration version
func RunMigration(db *gorm.DB, version int) error {
	if version < 1 || version > 3 {
		return fmt.Errorf("unregistered version: %d", version)
	}

	type ProductV1 struct {
		ID          uint      `gorm:"primaryKey"`
		Name        string    `gorm:"not null"`
		Price       float64   `gorm:"not null"`
		Description string    `gorm:"type:text"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	type ProductV3Tmp struct {
		ID          uint      `gorm:"primaryKey"`
		Name        string    `gorm:"not null"`
		Price       float64   `gorm:"not null"`
		Description string    `gorm:"type:text"`
		CategoryID  uint      `gorm:"not null"`
		Stock       int       `gorm:"default:0"`
		SKU         string    `gorm:"unique;not null"`
		IsActive    bool      `gorm:"default:true"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var err error
		switch version {
		case 1:
			err = tx.Table("products").AutoMigrate(&ProductV1{})

		case 2:
			err = tx.AutoMigrate(&Category{})

		case 3:
			tmpTableName := "_products_v3_tmp"

			if err = tx.Table(tmpTableName).AutoMigrate(&ProductV3Tmp{}); err != nil {
				return err
			}

			if tx.Migrator().HasTable("products") {
				var defaultCategory Category
				if err = tx.Where("name = ?", "Electronics").FirstOrCreate(&defaultCategory, Category{Name: "Electronics", Description: "Default"}).Error; err != nil {
					return err
				}

				transferSQL := fmt.Sprintf(`
					INSERT INTO %s (id, name, price, description, category_id, stock, sku, is_active, created_at, updated_at)
					SELECT id, name, price, description, ?, 0, 'GEN-SKU-' || id, 1, created_at, updated_at 
					FROM products;
				`, tmpTableName)

				if err = tx.Exec(transferSQL, defaultCategory.ID).Error; err != nil {
					return err
				}

				if err = tx.Migrator().DropTable("products"); err != nil {
					return err
				}
			}

			if err = tx.Migrator().RenameTable(tmpTableName, "products"); err != nil {
				return err
			}
			
			err = tx.AutoMigrate(&Product{})
		}

		if err != nil {
			return err
		}

		mv := MigrationVersion{ID: 1, Version: version, AppliedAt: time.Now()}
		return tx.Save(&mv).Error
	})
}

// RollbackMigration rolls back to a specific migration version
func RollbackMigration(db *gorm.DB, version int) error {
	if version < 1 || version > 3 {
		return fmt.Errorf("unsupported version: %d", version)
	}

	type ProductV1 struct {
		ID          uint      `gorm:"primaryKey"`
		Name        string    `gorm:"not null"`
		Price       float64   `gorm:"not null"`
		Description string    `gorm:"type:text"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	return db.Transaction(func(tx *gorm.DB) error {
		currentVersion, err := GetMigrationVersion(tx)
		if err != nil {
			return err
		}

		if currentVersion == 3 && version <= 2 {
			tmpTableName := "_products_v1_tmp"

			if err = tx.Table(tmpTableName).AutoMigrate(&ProductV1{}); err != nil {
				return err
			}

			if tx.Migrator().HasTable("products") {
				transferSQL := fmt.Sprintf(`
					INSERT INTO %s (id, name, price, description, created_at, updated_at)
					SELECT id, name, price, description, created_at, updated_at 
					FROM products;
				`, tmpTableName)

				if err = tx.Exec(transferSQL).Error; err != nil {
					return err
				}

				if err = tx.Migrator().DropTable("products"); err != nil {
					return err
				}
			}

			if err = tx.Migrator().RenameTable(tmpTableName, "products"); err != nil {
				return err
			}
		}

		if currentVersion >= 2 && version == 1 {
			if tx.Migrator().HasTable(&Category{}) {
				if err = tx.Migrator().DropTable(&Category{}); err != nil {
					return err
				}
			}
		}

		mv := MigrationVersion{ID: 1, Version: version, AppliedAt: time.Now()}
		return tx.Save(&mv).Error
	})
}

// GetMigrationVersion gets the current migration version
func GetMigrationVersion(db *gorm.DB) (int, error) {
	var mv MigrationVersion
	err := db.First(&mv, 1).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}

	return mv.Version, nil
}

// SeedData populates the database with initial data
func SeedData(db *gorm.DB) error {
	category := Category{
		Name:        "Electronics",
		Description: "Gadgets and devices",
	}

	err := db.Where("name = ?", category.Name).FirstOrCreate(&category).Error
	if err != nil {
		return err
	}

	product := Product{
		Name:        "Laptop Pro",
		Price:       1299.99,
		Description: "High-performance development laptop",
		CategoryID:  category.ID,
		Stock:       50,
		SKU:         "LT-PRO-01",
		IsActive:    true,
	}

	if db.Migrator().HasColumn(&Product{}, "SKU") {
		err = db.Where("sku = ?", product.SKU).FirstOrCreate(&product).Error
	} else {
		err = db.Where("name = ?", product.Name).FirstOrCreate(&product).Error
	}
	return err
}

// CreateProduct creates a new product with validation
func CreateProduct(db *gorm.DB, product *Product) error {
	if !db.Migrator().HasTable(&Product{}) || !db.Migrator().HasColumn(&Product{}, "SKU") {
		return fmt.Errorf("table products does not match version 3 schema")
	}
	return db.Create(product).Error
}

// GetProductsByCategory retrieves all products in a specific category
func GetProductsByCategory(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	err := db.Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProductStock updates the stock quantity of a product
func UpdateProductStock(db *gorm.DB, productID uint, quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("stock cannot be negative")
	}

	res := db.Model(&Product{}).Where("id = ?", productID).Update("stock", quantity)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
