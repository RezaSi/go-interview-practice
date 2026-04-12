package main

import (
	"time"
	"fmt"
	"errors"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

// MigrationVersion tracks the current database schema version
type MigrationVersion struct {
	ID        uint `gorm:"primaryKey"`
	Version   int  `gorm:"unique;not null"`
	AppliedAt time.Time
}

// Product represents a product in the e-commerce system
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

// ----------------------------------------------------------------
// Degine all migrations
// ----------------------------------------------------------------

type Migration struct {
	Version int
	Up      func(*gorm.DB) error
	Down    func(*gorm.DB) error
}

var migrations = []Migration{
	{
		Version: 1,
		Up: func(tx *gorm.DB) error {
			return tx.Exec(`
				CREATE TABLE IF NOT EXISTS products (
					id integer PRIMARY KEY,
					name TEXT NOT NULL,
					price REAL NOT NULL,
					description TEXT,
					created_at DATETIME,
					updated_at DATETIME
				);
			`).Error
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("products")
		},
	},
	{
		Version: 2,
		Up: func(tx *gorm.DB) error {
			err := tx.Migrator().CreateTable(&Category{})
			if err != nil {
				return err
			}
			return tx.Exec("ALTER TABLE products ADD COLUMN category_id INTEGER NOT NULL DEFAULT 0").Error
		},
		Down: func(tx *gorm.DB) error {
			err := tx.Exec("ALTER TABLE products DROP COLUMN category_id").Error
			if err != nil {
				return err
			}
			return tx.Migrator().DropTable("categories")
		},
	},
	{
		Version: 3,
		Up: func(tx *gorm.DB) error {
			return tx.Exec(`
				ALTER TABLE products ADD COLUMN stock INTEGER NOT NULL DEFAULT 0;
				ALTER TABLE products ADD COLUMN sku TEXT NOT NULL DEFAULT '';
				ALTER TABLE products ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true;
			`).Error
		},
		Down: func(tx *gorm.DB) error {
			return tx.Exec(`
				ALTER TABLE products DROP COLUMN stock;
				ALTER TABLE products DROP COLUMN sku;
				ALTER TABLE products DROP COLUMN is_active;
			`).Error
		},
	},
}

// ConnectDB establishes a connection to the SQLite database
func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&MigrationVersion{}); err != nil {
		return nil, err
	}
	return db, nil
}

func GetMigrationVersion(db *gorm.DB) (int, error) {
	var mv MigrationVersion
	err := db.Order("version DESC").First(&mv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return mv.Version, err
}

func setMigrationVersion(db *gorm.DB, version int) error {
	mv := MigrationVersion{
		Version:   version,
		AppliedAt: time.Now(),
	}
	return db.Create(&mv).Error
}

func removeMigrationVersion(db *gorm.DB, version int) error {
	return db.Where("version=?", version).Delete(&MigrationVersion{}).Error
}

// ----------------------------------------------------------------
// Run and Rollback migrations
// ----------------------------------------------------------------

func RunMigration(db *gorm.DB, version int) error {
	current, err := GetMigrationVersion(db)
	if err != nil {
		return err
	}
	if version < current {
		return fmt.Errorf("version %d < current %d", version, current)
	}
	if version > len(migrations) {
		return fmt.Errorf("invalid version %d, (max = %d)", version, len(migrations))
	}
	if version == current {
		return nil
	}

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, m := range(migrations) {
		if m.Version > current && m.Version <= version {
			if err := m.Up(tx); err != nil {
				tx.Rollback()
				return err
			}
			if err := setMigrationVersion(tx, m.Version); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

func RollbackMigration(db *gorm.DB, version int) error {
	current, err := GetMigrationVersion(db)
	if err != nil {
		return err
	}
	if version >= current {
		return fmt.Errorf("nothing to rollback, version %d < current %d. ", version, current)
	}
	if version < 0 {
		return errors.New("version cannot be negative")
	}

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if m.Version <= current && m.Version > version {
			if err := m.Down(tx); err != nil {
				tx.Rollback()
				return err
			}
			if err := removeMigrationVersion(tx, m.Version); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

func SeedData(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Product{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // Already seeded
	}

	if err := db.Model(&Category{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		cat := Category{Name: "Category 1", Description: "Category 1"}
		if err := db.Create(&cat).Error; err != nil {
			return err
		}
	}
	cat := Category{}
	if err := db.First(&cat).Error; err != nil {
		return err
	}

	products := []Product{
		{
			Name:        "Product 1",
			Price:       421.39,
			Description: "Product 1",
			CategoryID:  cat.ID,
			Stock:       200,
			SKU:         "SKU-001",
			IsActive:    true,
		},
		{
			Name:        "Product 2",
			Price:       76.96,
			Description: "Product 2",
			CategoryID:  cat.ID,
			Stock:       1000,
			SKU:         "SKU-002",
			IsActive:    true,
		},
	}
	return db.Create(&products).Error
}

func CreateProduct(db *gorm.DB, product *Product) error {
	var count int64
	if err := db.Model(&Product{}).Where("sku=?", product.SKU).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("SKU %s already exists", product.SKU)
	}

	var cat Category
	if err := db.First(&cat, product.CategoryID).Error; err != nil {
		return err
	}
	return db.Create(product).Error
}

func GetProductsByCategory(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").
		Where("category_id=? AND is_active=?", categoryID, true).
		Find(&products).Error
	return products, err
}

func UpdateProductStock(db *gorm.DB, productID uint, quantity int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			return err
		}
		product.Stock = quantity
		return tx.Save(&product).Error
	})
}
