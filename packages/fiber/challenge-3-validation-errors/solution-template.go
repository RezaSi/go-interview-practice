package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

// Product represents a product in the catalog
type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description string   `json:"description" validate:"required,min=10,max=500"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	Category    string   `json:"category" validate:"required,oneof=electronics clothing books home"`
	SKU         string   `json:"sku" validate:"required,sku"`
	InStock     bool     `json:"in_stock"`
	Tags        []string `json:"tags" validate:"dive,min=2,max=20"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Details []ValidationError `json:"details,omitempty"`
}

// In-memory storage
var products = []Product{
	{ID: 1, Name: "Laptop", Description: "High-performance laptop for professionals", Price: 999.99, Category: "electronics", SKU: "PROD-12345", InStock: true, Tags: []string{"computer", "work"}},
	{ID: 2, Name: "T-Shirt", Description: "Comfortable cotton t-shirt", Price: 29.99, Category: "clothing", SKU: "PROD-67890", InStock: true, Tags: []string{"clothing", "casual"}},
}
var nextID = 3

var validate *validator.Validate

func main() {
	// TODO: Create Fiber app
	app := fiber.New()

	// Setup custom validator
	setupCustomValidator()

	// TODO: Define routes
	// GET /products - get all products with filtering
	// GET /products/:id - get product by ID
	// POST /products - create new product
	// PUT /products/:id - update product
	// POST /products/bulk - create multiple products

	// TODO: Start server on port 3000
}

func setupCustomValidator() {
	// TODO: Setup custom validator
	// Use github.com/go-playground/validator/v10 package
	// Register custom validators (SKU format)
}

// TODO: Implement custom validators

// validateSKU validates SKU format (PROD-XXXXX)
func validateSKU(val string) bool {
	// TODO: Implement SKU validation
	// Format: "PROD-" followed by 5 digits
	// Example: PROD-12345
	return false
}

// TODO: Implement validation helper functions

// validateProduct validates a product using the validator
func validateProduct(product Product) []ValidationError {
	// TODO: Use validator to validate product struct
	// Convert validator errors to ValidationError slice
	// Return custom error messages for each validation failure
	return nil
}

// formatValidationError formats a validation error with custom message
func formatValidationError(field, tag string, value interface{}) ValidationError {
	// TODO: Create custom error messages based on validation tag
	// Examples:
	// - required: "{field} is required"
	// - min: "{field} must be at least {param} characters"
	// - max: "{field} cannot exceed {param} characters"
	// - gt: "{field} must be greater than {param}"
	// - oneof: "{field} must be one of: {param}"
	// - sku: "{field} must be in format PROD-XXXXX"

	return ValidationError{
		Field:   field,
		Tag:     tag,
		Value:   value,
		Message: "", // TODO: Generate appropriate message
	}
}

// TODO: Implement route handlers

// getProductsHandler returns all products with optional filtering
func getProductsHandler(c *fiber.Ctx) error {
	// TODO: Implement filtering by query parameters
	// ?category=electronics - filter by category
	// ?in_stock=true - filter by stock status
	// ?min_price=10&max_price=100 - filter by price range
	// ?search=laptop - search in name and description

	return c.JSON(products)
}

// getProductHandler returns a specific product by ID
func getProductHandler(c *fiber.Ctx) error {
	// TODO: Get product ID from URL parameter
	// Return 404 if product not found
	// Return product as JSON
	return nil
}

// createProductHandler creates a new product
func createProductHandler(c *fiber.Ctx) error {
	// TODO: Parse request body into Product struct
	// Validate product using validateProduct()
	// Return validation errors if validation fails
	// Generate unique SKU if not provided
	// Add product to storage with auto-increment ID
	// Return created product
	return nil
}

// updateProductHandler updates an existing product
func updateProductHandler(c *fiber.Ctx) error {
	// TODO: Get product ID from URL parameter
	// Parse request body for updates
	// Validate updated product
	// Update product if exists, return 404 if not found
	// Return updated product
	return nil
}

// bulkCreateHandler creates multiple products in one request
func bulkCreateHandler(c *fiber.Ctx) error {
	// TODO: Parse request body into []Product slice
	// Validate each product
	// Continue processing even if some products fail validation
	// Return results with success/failure for each product
	// Format: {
	//   "success": true,
	//   "results": [
	//     {"success": true, "product": {...}},
	//     {"success": false, "errors": [...]}
	//   ]
	// }
	return nil
}

// Helper functions

// findProductByID finds a product by ID
func findProductByID(id int) (*Product, int) {
	for i, product := range products {
		if product.ID == id {
			return &product, i
		}
	}
	return nil, -1
}

// generateSKU generates a unique SKU
func generateSKU() string {
	// TODO: Generate unique SKU in format PROD-XXXXX
	// Use timestamp or random number for uniqueness
	return ""
}

// filterProducts filters products based on query parameters
func filterProducts(products []Product, filters map[string]string) []Product {
	// TODO: Implement filtering logic
	// Support category, in_stock, min_price, max_price, search filters
	return products
}
