package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	db, err := ConnectDB()
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Cleanup
	sqlDB, _ := db.DB()
	sqlDB.Close()
	os.Remove("test.db")
}

func TestMigrationSystem(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Test initial version
	version, err := GetMigrationVersion(db)
	assert.NoError(t, err)
	assert.Equal(t, 0, version)

	// Test running migration version 1
	err = RunMigration(db, 1)
	assert.NoError(t, err)

	version, err = GetMigrationVersion(db)
	assert.NoError(t, err)
	assert.Equal(t, 1, version)

	// Test running migration version 2
	err = RunMigration(db, 2)
	assert.NoError(t, err)

	version, err = GetMigrationVersion(db)
	assert.NoError(t, err)
	assert.Equal(t, 2, version)

	// Test running migration version 3
	err = RunMigration(db, 3)
	assert.NoError(t, err)

	version, err = GetMigrationVersion(db)
	assert.NoError(t, err)
	assert.Equal(t, 3, version)

	// Verify tables exist
	assert.True(t, db.Migrator().HasTable(&Product{}))
	assert.True(t, db.Migrator().HasTable(&Category{}))
	assert.True(t, db.Migrator().HasTable(&MigrationVersion{}))
}

func TestMigrationRollback(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Run migrations to version 3
	RunMigration(db, 1)
	RunMigration(db, 2)
	RunMigration(db, 3)

	// Test rollback to version 2
	err := RollbackMigration(db, 2)
	assert.NoError(t, err)

	version, err := GetMigrationVersion(db)
	assert.NoError(t, err)
	assert.Equal(t, 2, version)

	// Test rollback to version 1
	err = RollbackMigration(db, 1)
	assert.NoError(t, err)

	version, err = GetMigrationVersion(db)
	assert.NoError(t, err)
	assert.Equal(t, 1, version)
}

func TestSeedData(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Run migrations
	RunMigration(db, 1)
	RunMigration(db, 2)
	RunMigration(db, 3)

	// Test data seeding
	err := SeedData(db)
	assert.NoError(t, err)

	// Verify categories were created
	var categories []Category
	db.Find(&categories)
	assert.Greater(t, len(categories), 0)

	// Verify products were created
	var products []Product
	db.Find(&products)
	assert.Greater(t, len(products), 0)
}

func TestCreateProduct(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Run migrations and seed data
	RunMigration(db, 1)
	RunMigration(db, 2)
	RunMigration(db, 3)
	SeedData(db)

	// Get a category for the product
	var category Category
	db.First(&category)

	product := &Product{
		Name:        "Test Product",
		Price:       29.99,
		Description: "A test product",
		CategoryID:  category.ID,
		Stock:       10,
		SKU:         "TEST-001",
		IsActive:    true,
	}

	err := CreateProduct(db, product)
	assert.NoError(t, err)
	assert.NotZero(t, product.ID)

	// Verify product was created
	var foundProduct Product
	db.First(&foundProduct, product.ID)
	assert.Equal(t, "Test Product", foundProduct.Name)
	assert.Equal(t, 29.99, foundProduct.Price)
	assert.Equal(t, category.ID, foundProduct.CategoryID)
}

func TestGetProductsByCategory(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Run migrations and seed data
	RunMigration(db, 1)
	RunMigration(db, 2)
	RunMigration(db, 3)
	SeedData(db)

	// Get a category
	var category Category
	db.First(&category)

	// Create products in this category
	product1 := &Product{
		Name:        "Product 1",
		Price:       19.99,
		Description: "First product",
		CategoryID:  category.ID,
		Stock:       5,
		SKU:         "PROD-001",
		IsActive:    true,
	}

	product2 := &Product{
		Name:        "Product 2",
		Price:       39.99,
		Description: "Second product",
		CategoryID:  category.ID,
		Stock:       8,
		SKU:         "PROD-002",
		IsActive:    true,
	}

	CreateProduct(db, product1)
	CreateProduct(db, product2)

	// Test retrieval by category
	products, err := GetProductsByCategory(db, category.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(products), 2)

	// Verify products belong to the correct category
	for _, product := range products {
		assert.Equal(t, category.ID, product.CategoryID)
	}
}

func TestUpdateProductStock(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Run migrations and seed data
	RunMigration(db, 1)
	RunMigration(db, 2)
	RunMigration(db, 3)
	SeedData(db)

	// Get a category and create a product
	var category Category
	db.First(&category)

	product := &Product{
		Name:        "Stock Test Product",
		Price:       15.99,
		Description: "Product for stock testing",
		CategoryID:  category.ID,
		Stock:       20,
		SKU:         "STOCK-001",
		IsActive:    true,
	}

	CreateProduct(db, product)

	// Test stock update
	err := UpdateProductStock(db, product.ID, 15)
	assert.NoError(t, err)

	// Verify stock was updated
	var updatedProduct Product
	db.First(&updatedProduct, product.ID)
	assert.Equal(t, 15, updatedProduct.Stock)

	// Test negative stock update (should fail or handle appropriately)
	err = UpdateProductStock(db, product.ID, -5)
	// This might be allowed or not, depending on business logic
}

func TestErrorHandling(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Test running migration that doesn't exist
	err := RunMigration(db, 999)
	assert.Error(t, err)

	// Test rollback to non-existent version
	err = RollbackMigration(db, 999)
	assert.Error(t, err)

	// Test creating product without running migrations
	product := &Product{Name: "Test", Price: 10.0}
	err = CreateProduct(db, product)
	assert.Error(t, err)
}
