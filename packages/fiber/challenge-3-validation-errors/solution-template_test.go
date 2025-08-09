package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	// Reset products data for each test
	products = []Product{
		{ID: 1, Name: "Laptop", Description: "High-performance laptop for professionals", Price: 999.99, Category: "electronics", SKU: "PROD-12345", InStock: true, Tags: []string{"computer", "work"}},
		{ID: 2, Name: "T-Shirt", Description: "Comfortable cotton t-shirt", Price: 29.99, Category: "clothing", SKU: "PROD-67890", InStock: true, Tags: []string{"clothing", "casual"}},
	}
	nextID = 3

	setupCustomValidator()

	app := fiber.New()

	// Setup routes (implement these in solution)
	app.Get("/products", getProductsHandler)
	app.Get("/products/:id", getProductHandler)
	app.Post("/products", createProductHandler)
	app.Put("/products/:id", updateProductHandler)
	app.Post("/products/bulk", bulkCreateHandler)

	return app
}

func TestGetProducts(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/products", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetProductByID(t *testing.T) {
	app := setupTestApp()

	// Test existing product
	req := httptest.NewRequest("GET", "/products/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test non-existent product
	req = httptest.NewRequest("GET", "/products/999", nil)
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestCreateValidProduct(t *testing.T) {
	app := setupTestApp()

	productData := Product{
		Name:        "Valid Product",
		Description: "This is a valid product description with enough characters",
		Price:       99.99,
		Category:    "electronics",
		SKU:         "PROD-11111",
		InStock:     true,
		Tags:        []string{"test", "valid"},
	}

	body, _ := json.Marshal(productData)
	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestCreateProductValidation(t *testing.T) {
	app := setupTestApp()

	// Test invalid product (missing required fields)
	invalidProduct := Product{
		Name:     "X",       // Too short
		Price:    -10,       // Invalid price
		Category: "invalid", // Invalid category
		SKU:      "INVALID", // Invalid SKU format
	}

	body, _ := json.Marshal(invalidProduct)
	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var response ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.False(t, response.Success)
	assert.NotEmpty(t, response.Details)
}

func TestUpdateProduct(t *testing.T) {
	app := setupTestApp()

	updateData := Product{
		Name:        "Updated Product",
		Description: "Updated product description with enough characters",
		Price:       149.99,
		Category:    "electronics",
		SKU:         "PROD-99999",
		InStock:     false,
		Tags:        []string{"updated", "test"},
	}

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestBulkCreate(t *testing.T) {
	app := setupTestApp()

	products := []Product{
		{
			Name:        "Bulk Product 1",
			Description: "First bulk product with valid description",
			Price:       49.99,
			Category:    "electronics",
			SKU:         "PROD-88888",
			InStock:     true,
			Tags:        []string{"bulk", "test"},
		},
		{
			Name:        "X", // Invalid - too short
			Description: "Invalid product",
			Price:       -10,       // Invalid price
			Category:    "invalid", // Invalid category
			SKU:         "INVALID", // Invalid SKU
		},
	}

	body, _ := json.Marshal(products)
	req := httptest.NewRequest("POST", "/products/bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	results := response["results"].([]interface{})

	// First product should succeed
	firstResult := results[0].(map[string]interface{})
	assert.True(t, firstResult["success"].(bool))

	// Second product should fail
	secondResult := results[1].(map[string]interface{})
	assert.False(t, secondResult["success"].(bool))
}

func TestProductFiltering(t *testing.T) {
	app := setupTestApp()

	// Test category filtering
	req := httptest.NewRequest("GET", "/products?category=electronics", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test price filtering
	req = httptest.NewRequest("GET", "/products?min_price=50&max_price=1000", nil)
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test search
	req = httptest.NewRequest("GET", "/products?search=laptop", nil)
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSKUValidation(t *testing.T) {
	// Test the SKU validator directly
	assert.True(t, isValidSKU("PROD-12345"))
	assert.False(t, isValidSKU("INVALID"))
	assert.False(t, isValidSKU("PROD-ABC"))
	assert.False(t, isValidSKU("PROD-123"))
}

// Helper function to test SKU validation
func isValidSKU(sku string) bool {
	// This should match the implementation in the solution
	return validateSKU(sku)
}
