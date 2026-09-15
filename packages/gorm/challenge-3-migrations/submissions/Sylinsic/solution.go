package main

import (
	"errors"
	"slices"
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

type migration struct {
	version int
	up      func(*gorm.DB) error
	down    func(*gorm.DB) error
}

var (
	migrations []migration = []migration{
		{
			version: 1,
			up: func(tx *gorm.DB) error {
				var err error
				if applyMigration(&err, func() error {
					return tx.AutoMigrate(&MigrationVersion{})
				}) ||
					applyMigration(&err, func() error {
						return tx.Exec(`
							CREATE TABLE IF NOT EXISTS products (
								id INTEGER PRIMARY KEY,
								name TEXT NOT NULL,
								price REAL NOT NULL,
								description TEXT,
								created_at DATETIME,
								updated_at DATETIME 
							)
						`).Error
					}) {
					return err
				}

				return nil
			},
			down: func(tx *gorm.DB) error {
				var err error
				if applyMigration(&err, func() error {
					return tx.Migrator().DropTable(&MigrationVersion{}, &Product{})
				}) {
					return err
				}

				return nil
			},
		},
		{
			version: 2,
			up: func(tx *gorm.DB) error {
				var err error
				if applyMigration(&err, func() error {
					return tx.Migrator().CreateTable(&Category{})
				}) ||
					applyMigration(&err, func() error {
						return tx.Migrator().AddColumn(&Product{}, "CategoryID")
					}) {
					return err
				}

				return nil
			},
			down: func(tx *gorm.DB) error {
				var err error
				if applyMigration(&err, func() error {
					return tx.Migrator().DropTable(&Category{})
				}) ||
					applyMigration(&err, func() error {
						return tx.Migrator().DropColumn(&Product{}, "CategoryID")
					}) {
					return err
				}

				return nil
			},
		},
		{
			version: 3,
			up: func(tx *gorm.DB) error {
				var err error
				if applyMigration(&err, func() error {
					return tx.Migrator().AddColumn(&Product{}, "Stock")
				}) || applyMigration(&err, func() error {
					return tx.Exec(`
						ALTER TABLE products ADD COLUMN sku TEXT NOT NULL DEFAULT '';
						UPDATE products SET sku = "SKU-" || id;
						CREATE UNIQUE INDEX idx_products_sku ON products(sku);						
					`).Error
				}) || applyMigration(&err, func() error {
					return tx.Migrator().AddColumn(&Product{}, "IsActive")
				}) {
					return err
				}

				return nil
			},
			down: func(tx *gorm.DB) error {
				var err error
				if applyMigration(&err, func() error {
					return tx.Migrator().DropColumn(&Product{}, "Stock")
				}) || applyMigration(&err, func() error {
					return tx.Exec(`
						DROP INDEX IF EXISTS idx_products_sku
					`).Error
				}) || applyMigration(&err, func() error {
					return tx.Migrator().DropColumn(&Product{}, "SKU")
				}) || applyMigration(&err, func() error {
					return tx.Migrator().DropColumn(&Product{}, "IsActive")
				}) {
					return err
				}

				return nil
			},
		},
	}
)

var (
	migrationsNotYetApplied error = errors.New("Migrations haven't been applied yet")
)

func applyMigration(err *error, migration func() error) bool {
	if *err = migration(); err != nil {
		return false
	}

	return true
}

// ConnectDB establishes a connection to the SQLite database
func ConnectDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
}

func getMigrationByVersion(version int) (*migration, error) {
	migrationIndex := slices.IndexFunc(migrations, func(m migration) bool {
		return m.version == version
	})

	if migrationIndex == -1 {
		return nil, errors.New("No migration found")
	}

	return &migrations[migrationIndex], nil
}

// RunMigration runs a specific migration version
func RunMigration(db *gorm.DB, version int) error {
	tx := db.Begin()
	defer tx.Rollback()

	// Check migration exists
	_, err := getMigrationByVersion(version)
	if err != nil {
		return err
	}

	latestVersion, err := GetMigrationVersion(tx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// Already at same migration level
	if latestVersion == version {
		return nil
	}

	// If ahead, downgrade
	if latestVersion > version {
		if err := migrateDown(tx, latestVersion, version); err != nil {
			return err
		}
	}

	// If behind, upgrade
	if latestVersion < version {
		if err := migrateUp(tx, latestVersion, version); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// RollbackMigration rolls back to a specific migration version
func RollbackMigration(db *gorm.DB, version int) error {
	return RunMigration(db, version)
}

func migrateUp(tx *gorm.DB, latestVersion int, desiredVersion int) error {
	for i := latestVersion; i < desiredVersion; i++ {
		migration, err := getMigrationByVersion(i + 1)
		if err != nil {
			return err
		}

		err = migration.up(tx)
		if err != nil {
			return err
		}

		// Update db with latest migration
		tx.Create(&MigrationVersion{
			Version:   i + 1,
			AppliedAt: time.Now().UTC(),
		})
	}

	return nil
}

func migrateDown(tx *gorm.DB, latestVersion int, desiredVersion int) error {
	for i := latestVersion; i > desiredVersion; i-- {
		migration, err := getMigrationByVersion(i)
		if err != nil {
			return err
		}

		err = migration.down(tx)
		if err != nil {
			return err
		}

		// Update db with latest migration
		tx.Delete(&MigrationVersion{}, "Version = ?", i)
	}

	return nil
}

// GetMigrationVersion gets the current migration version
func GetMigrationVersion(db *gorm.DB) (int, error) {
	if !db.Migrator().HasTable("migration_versions") {
		db.AutoMigrate(&MigrationVersion{})
		db.Create(&MigrationVersion{
			Version: 0,
			AppliedAt: time.Now().UTC(),
		})
		return 0, nil
	}

	m := MigrationVersion{}
	if err := db.Last(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}

		return 0, err
	}
	return m.Version, nil
}

// SeedData populates the database with initial data
func SeedData(db *gorm.DB) error {
	db.Create(&Category{
		Name: "Computers",
		Description: "Computing devices",
		Products: []Product{
			{
				Name: "Dell laptop",
				Price: 300,
				Description: "Latest technology!",
				SKU: "ABC-DEF",
			},
			{
				Name: "Home-built PC",
				Price: 1500,
				Description: "Built to your specification!",
				SKU: "GHI-jkl",
			},
		},
	})
	return nil
}

// CreateProduct creates a new product with validation
func CreateProduct(db *gorm.DB, product *Product) error {
	if product.Price <= 0 {
		return errors.New("Product must be over £0")
	}

	return db.Create(product).Error
}

// GetProductsByCategory retrieves all products in a specific category
func GetProductsByCategory(db *gorm.DB, categoryID uint) ([]Product, error) {
	var category Category
	if err := db.Preload("Products").First(&category, categoryID).Error; err != nil {
		return nil, err
	}
	return category.Products, nil
}

// UpdateProductStock updates the stock quantity of a product
func UpdateProductStock(db *gorm.DB, productID uint, quantity int) error {
	if err := db.Model(&Product{}).
		Where("ID = ?", productID).
		UpdateColumn("Stock", quantity).
		Error; err != nil {
		return err
	}

	return nil
}
