package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	// Clean up any existing test files before and after tests
	os.Remove("inventory.json")
	code := m.Run()
	os.Remove("inventory.json")
	os.Exit(code)
}

func setupTest() {
	// Reset global inventory for each test
	inventory = &Inventory{
		Products:   []Product{},
		Categories: []Category{},
		NextID:     1,
	}
	os.Remove("inventory.json")
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	setupTest()
	output, err := executeCommand(rootCmd)

	if err != nil {
		t.Fatalf("Root command failed: %v", err)
	}

	if !strings.Contains(output, "Inventory Management CLI") {
		t.Error("Root command should contain 'Inventory Management CLI' in output")
	}

	if !strings.Contains(output, "product") {
		t.Error("Root command should show 'product' subcommand")
	}

	if !strings.Contains(output, "category") {
		t.Error("Root command should show 'category' subcommand")
	}
}

func TestProductAddCommand(t *testing.T) {
	setupTest()

	args := []string{"product", "add", "--name", "Test Product", "--price", "99.99", "--category", "Test", "--stock", "10"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Product add command failed: %v", err)
	}

	if !strings.Contains(output, "Product added successfully") {
		t.Error("Product add should show success message")
	}

	if !strings.Contains(output, "Test Product") {
		t.Error("Product add should show product name")
	}

	// Verify product was added to inventory
	if len(inventory.Products) != 1 {
		t.Error("Product should be added to inventory")
	}

	product := inventory.Products[0]
	if product.Name != "Test Product" || product.Price != 99.99 || product.Category != "Test" || product.Stock != 10 {
		t.Error("Product details should match input")
	}
}

func TestProductListCommand(t *testing.T) {
	setupTest()

	// Add test products
	inventory.Products = []Product{
		{ID: 1, Name: "Product 1", Price: 10.0, Category: "Cat1", Stock: 5},
		{ID: 2, Name: "Product 2", Price: 20.0, Category: "Cat2", Stock: 10},
	}

	output, err := executeCommand(rootCmd, "product", "list")

	if err != nil {
		t.Fatalf("Product list command failed: %v", err)
	}

	if !strings.Contains(output, "Product 1") || !strings.Contains(output, "Product 2") {
		t.Error("Product list should show all products")
	}

	if !strings.Contains(output, "ID") || !strings.Contains(output, "Name") {
		t.Error("Product list should have table headers")
	}
}

func TestProductGetCommand(t *testing.T) {
	setupTest()

	// Add test product
	inventory.Products = []Product{
		{ID: 1, Name: "Test Product", Price: 99.99, Category: "Test", Stock: 5},
	}

	output, err := executeCommand(rootCmd, "product", "get", "1")

	if err != nil {
		t.Fatalf("Product get command failed: %v", err)
	}

	if !strings.Contains(output, "Test Product") {
		t.Error("Product get should show product details")
	}

	// Test invalid ID
	output, err = executeCommand(rootCmd, "product", "get", "999")
	if err == nil && !strings.Contains(output, "not found") {
		t.Error("Product get should handle invalid ID")
	}
}

func TestProductUpdateCommand(t *testing.T) {
	setupTest()

	// Add test product
	inventory.Products = []Product{
		{ID: 1, Name: "Old Name", Price: 50.0, Category: "Old", Stock: 5},
	}
	inventory.NextID = 2

	args := []string{"product", "update", "1", "--name", "New Name", "--price", "75.0"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Product update command failed: %v", err)
	}

	if !strings.Contains(output, "updated successfully") {
		t.Error("Product update should show success message")
	}

	// Verify product was updated
	product := inventory.Products[0]
	if product.Name != "New Name" || product.Price != 75.0 {
		t.Error("Product should be updated with new values")
	}
}

func TestProductDeleteCommand(t *testing.T) {
	setupTest()

	// Add test products
	inventory.Products = []Product{
		{ID: 1, Name: "Product 1", Price: 10.0, Category: "Cat1", Stock: 5},
		{ID: 2, Name: "Product 2", Price: 20.0, Category: "Cat2", Stock: 10},
	}

	output, err := executeCommand(rootCmd, "product", "delete", "1")

	if err != nil {
		t.Fatalf("Product delete command failed: %v", err)
	}

	if !strings.Contains(output, "deleted successfully") {
		t.Error("Product delete should show success message")
	}

	// Verify product was deleted
	if len(inventory.Products) != 1 || inventory.Products[0].ID != 2 {
		t.Error("Product should be deleted from inventory")
	}
}

func TestCategoryAddCommand(t *testing.T) {
	setupTest()

	args := []string{"category", "add", "--name", "Electronics", "--description", "Electronic devices"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Category add command failed: %v", err)
	}

	if !strings.Contains(output, "Category added successfully") {
		t.Error("Category add should show success message")
	}

	// Verify category was added
	if len(inventory.Categories) != 1 {
		t.Error("Category should be added to inventory")
	}

	category := inventory.Categories[0]
	if category.Name != "Electronics" || category.Description != "Electronic devices" {
		t.Error("Category details should match input")
	}
}

func TestCategoryListCommand(t *testing.T) {
	setupTest()

	// Add test categories
	inventory.Categories = []Category{
		{Name: "Electronics", Description: "Electronic devices"},
		{Name: "Books", Description: "Reading materials"},
	}

	output, err := executeCommand(rootCmd, "category", "list")

	if err != nil {
		t.Fatalf("Category list command failed: %v", err)
	}

	if !strings.Contains(output, "Electronics") || !strings.Contains(output, "Books") {
		t.Error("Category list should show all categories")
	}
}

func TestSearchCommand(t *testing.T) {
	setupTest()

	// Add test products
	inventory.Products = []Product{
		{ID: 1, Name: "Laptop", Price: 999.99, Category: "Electronics", Stock: 5},
		{ID: 2, Name: "Phone", Price: 599.99, Category: "Electronics", Stock: 10},
		{ID: 3, Name: "Book", Price: 19.99, Category: "Education", Stock: 20},
	}

	// Test search by name
	output, err := executeCommand(rootCmd, "search", "--name", "Laptop")
	if err != nil {
		t.Fatalf("Search by name failed: %v", err)
	}
	if !strings.Contains(output, "Laptop") || strings.Contains(output, "Phone") {
		t.Error("Search by name should find only matching products")
	}

	// Test search by category
	output, err = executeCommand(rootCmd, "search", "--category", "Electronics")
	if err != nil {
		t.Fatalf("Search by category failed: %v", err)
	}
	if !strings.Contains(output, "Laptop") || !strings.Contains(output, "Phone") || strings.Contains(output, "Book") {
		t.Error("Search by category should find all products in category")
	}

	// Test search by price range
	output, err = executeCommand(rootCmd, "search", "--min-price", "500", "--max-price", "1000")
	if err != nil {
		t.Fatalf("Search by price range failed: %v", err)
	}
	if !strings.Contains(output, "Laptop") || !strings.Contains(output, "Phone") || strings.Contains(output, "Book") {
		t.Error("Search by price range should find products within range")
	}
}

func TestStatsCommand(t *testing.T) {
	setupTest()

	// Add test data
	inventory.Products = []Product{
		{ID: 1, Name: "Product 1", Price: 100.0, Category: "Cat1", Stock: 5},
		{ID: 2, Name: "Product 2", Price: 200.0, Category: "Cat2", Stock: 0}, // Out of stock
		{ID: 3, Name: "Product 3", Price: 50.0, Category: "Cat1", Stock: 2},  // Low stock
	}
	inventory.Categories = []Category{
		{Name: "Cat1", Description: "Category 1"},
		{Name: "Cat2", Description: "Category 2"},
	}

	output, err := executeCommand(rootCmd, "stats")

	if err != nil {
		t.Fatalf("Stats command failed: %v", err)
	}

	if !strings.Contains(output, "Total Products: 3") {
		t.Error("Stats should show correct total products")
	}

	if !strings.Contains(output, "Total Categories: 2") {
		t.Error("Stats should show correct total categories")
	}

	if !strings.Contains(output, "Total Value") {
		t.Error("Stats should show total value")
	}

	if !strings.Contains(output, "Low Stock") {
		t.Error("Stats should show low stock count")
	}

	if !strings.Contains(output, "Out of Stock") {
		t.Error("Stats should show out of stock count")
	}
}

func TestDataPersistence(t *testing.T) {
	setupTest()

	// Add test data
	inventory.Products = []Product{
		{ID: 1, Name: "Test Product", Price: 99.99, Category: "Test", Stock: 10},
	}
	inventory.Categories = []Category{
		{Name: "Test", Description: "Test category"},
	}
	inventory.NextID = 2

	// Save inventory
	err := SaveInventory()
	if err != nil {
		t.Fatalf("SaveInventory failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat("inventory.json"); os.IsNotExist(err) {
		t.Fatal("Inventory file should be created")
	}

	// Reset inventory and load from file
	inventory = &Inventory{}
	err = LoadInventory()
	if err != nil {
		t.Fatalf("LoadInventory failed: %v", err)
	}

	// Verify data was loaded correctly
	if len(inventory.Products) != 1 || inventory.Products[0].Name != "Test Product" {
		t.Error("Products should be loaded correctly from file")
	}

	if len(inventory.Categories) != 1 || inventory.Categories[0].Name != "Test" {
		t.Error("Categories should be loaded correctly from file")
	}

	if inventory.NextID != 2 {
		t.Error("NextID should be loaded correctly from file")
	}
}

func TestFindProductByID(t *testing.T) {
	setupTest()

	inventory.Products = []Product{
		{ID: 1, Name: "Product 1", Price: 10.0, Category: "Cat1", Stock: 5},
		{ID: 3, Name: "Product 3", Price: 30.0, Category: "Cat3", Stock: 15},
	}

	// Test finding existing product
	product, index := FindProductByID(1)
	if product == nil || product.Name != "Product 1" || index != 0 {
		t.Error("FindProductByID should find existing product")
	}

	product, index = FindProductByID(3)
	if product == nil || product.Name != "Product 3" || index != 1 {
		t.Error("FindProductByID should find existing product at correct index")
	}

	// Test finding non-existing product
	product, index = FindProductByID(999)
	if product != nil || index != -1 {
		t.Error("FindProductByID should return nil and -1 for non-existing product")
	}
}

func TestCategoryExists(t *testing.T) {
	setupTest()

	inventory.Categories = []Category{
		{Name: "Electronics", Description: "Electronic devices"},
		{Name: "Books", Description: "Reading materials"},
	}

	if !CategoryExists("Electronics") {
		t.Error("CategoryExists should return true for existing category")
	}

	if !CategoryExists("Books") {
		t.Error("CategoryExists should return true for existing category")
	}

	if CategoryExists("NonExistent") {
		t.Error("CategoryExists should return false for non-existing category")
	}
}

func TestErrorHandling(t *testing.T) {
	setupTest()

	// Test missing required flags for product add
	output, err := executeCommand(rootCmd, "product", "add")
	if err == nil && !strings.Contains(output, "required") {
		t.Error("Product add should require flags")
	}

	// Test invalid price format
	args := []string{"product", "add", "--name", "Test", "--price", "invalid", "--category", "Test", "--stock", "1"}
	output, err = executeCommand(rootCmd, args...)
	if err == nil && !strings.Contains(output, "invalid") {
		t.Error("Product add should handle invalid price format")
	}

	// Test invalid stock format
	args = []string{"product", "add", "--name", "Test", "--price", "10.0", "--category", "Test", "--stock", "invalid"}
	output, err = executeCommand(rootCmd, args...)
	if err == nil && !strings.Contains(output, "invalid") {
		t.Error("Product add should handle invalid stock format")
	}
}
