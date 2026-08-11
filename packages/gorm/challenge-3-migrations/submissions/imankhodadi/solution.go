package main

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DROP COLUMN requires SQLite ≥ 3.35.0.

type MigrationVersion struct {
	ID        uint `gorm:"primaryKey"`
	Version   int  `gorm:"unique;not null"`
	AppliedAt time.Time
}

type Product struct {
	ID          uint     `gorm:"primaryKey"`
	Name        string   `gorm:"not null"`
	Price       float64  `gorm:"not null"` //use an integer (cents) or a decimal type
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
// Define all migrations
// ----------------------------------------------------------------

type Migration struct {
	Version int
	Up      func(*gorm.DB) error
	Down    func(*gorm.DB) error
}

func init() {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
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
			err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS categories (
					id INTEGER PRIMARY KEY,
					name TEXT NOT NULL UNIQUE,
					description TEXT,
					created_at DATETIME,
					updated_at DATETIME
			);`).Error
			if err != nil {
				return err
			}
			// category_id = 0, which won't match any category, change in production
			// consider adding category_id from the begining (this is just for assignment)
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
			if err := tx.Exec(`ALTER TABLE products ADD COLUMN stock INTEGER NOT NULL DEFAULT 0`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE products ADD COLUMN sku TEXT NOT NULL DEFAULT ''`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE products ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true`).Error; err != nil {
				return err
			}
			// Backfill unique SKU values for existing rows to avoid unique index violation
			if err := tx.Exec(`UPDATE products SET sku = 'SKU-LEGACY-' || id WHERE sku = ''`).Error; err != nil {
				return err
			}
			return tx.Exec(`CREATE UNIQUE INDEX idx_products_sku ON products(sku)`).Error
		},
		Down: func(tx *gorm.DB) error {
			if err := tx.Exec(`DROP INDEX IF EXISTS idx_products_sku`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE products DROP COLUMN stock`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE products DROP COLUMN sku`).Error; err != nil {
				return err
			}
			return tx.Exec(`ALTER TABLE products DROP COLUMN is_active`).Error
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
	maxVersion := migrations[len(migrations)-1].Version
	if version > maxVersion {
		return fmt.Errorf("invalid version %d, (max = %d)", version, maxVersion)
	}
	if version == current {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, m := range migrations {
			if m.Version > current && m.Version <= version {
				if err := m.Up(tx); err != nil {
					return err
				}
				if err := setMigrationVersion(tx, m.Version); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func RollbackMigration(db *gorm.DB, version int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// TODO: don't use len(migrations) as the max version, which only works if migration versions are exactly 1, 2, …, N.
		// If a version is ever skipped or reordered, this check and the iteration logic will silently skip or misapply migrations.
		// consider sorting migrations by Version at init time and deriving the max from the last element for robustness.

		current, err := GetMigrationVersion(tx)
		if err != nil {
			return err
		}
		if version >= current {
			return fmt.Errorf("nothing to rollback: target version %d >= current version %d", version, current)
		}
		if version < 0 {
			return errors.New("version cannot be negative")
		}
		for i := len(migrations) - 1; i >= 0; i-- {
			m := migrations[i]
			if m.Version <= current && m.Version > version {
				if err := m.Down(tx); err != nil {
					return err
				}
				if err := removeMigrationVersion(tx, m.Version); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func SeedData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&Product{}).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil // Already seeded
		}

		if err := tx.Model(&Category{}).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			cat := Category{Name: "Category 1", Description: "Category 1"}
			if err := tx.Create(&cat).Error; err != nil {
				return err
			}
		}
		cat := Category{}
		if err := tx.First(&cat).Error; err != nil {
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
		return tx.Create(&products).Error
	})
}

func CreateProduct(db *gorm.DB, product *Product) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&Product{}).Where("sku=?", product.SKU).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("SKU %s already exists", product.SKU)
		}

		var cat Category
		if err := tx.First(&cat, product.CategoryID).Error; err != nil {
			return err
		}
		return tx.Create(product).Error
	})
}

func GetProductsByCategory(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").
		Where("category_id=? AND is_active=?", categoryID, true).
		Find(&products).Error
	return products, err
}

func UpdateProductStock(db *gorm.DB, productID uint, newStock int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if newStock < 0 {
			return errors.New("stock quantity cannot be negative")
		}
		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			return err
		}
		return tx.Model(&product).Update("stock", newStock).Error
	})
}
