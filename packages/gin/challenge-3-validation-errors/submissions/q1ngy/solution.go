package main

import (
	"errors"
	"slices"
	"strings"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Product represents a product in the catalog
type Product struct {
	ID          int                    `json:"id"`
	SKU         string                 `json:"sku" binding:"required"`
	Name        string                 `json:"name" binding:"required,min=3,max=100"`
	Description string                 `json:"description" binding:"max=1000"`
	Price       float64                `json:"price" binding:"required,min=0.01"`
	Currency    string                 `json:"currency" binding:"required"`
	Category    Category               `json:"category" binding:"required"`
	Tags        []string               `json:"tags"`
	Attributes  map[string]interface{} `json:"attributes"`
	Images      []Image                `json:"images"`
	Inventory   Inventory              `json:"inventory" binding:"required"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Category represents a product category
type Category struct {
	ID       int    `json:"id" binding:"required,min=1"`
	Name     string `json:"name" binding:"required"`
	Slug     string `json:"slug" binding:"required"`
	ParentID *int   `json:"parent_id,omitempty"`
}

// Image represents a product image
type Image struct {
	URL       string `json:"url" binding:"required,url"`
	Alt       string `json:"alt" binding:"required,min=5,max=200"`
	Width     int    `json:"width" binding:"min=100"`
	Height    int    `json:"height" binding:"min=100"`
	Size      int64  `json:"size"`
	IsPrimary bool   `json:"is_primary"`
}

// Inventory represents product inventory information
type Inventory struct {
	Quantity    int       `json:"quantity" binding:"required,min=0"`
	Reserved    int       `json:"reserved" binding:"min=0"`
	Available   int       `json:"available"` // Calculated field
	Location    string    `json:"location" binding:"required"`
	LastUpdated time.Time `json:"last_updated"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Tag     string      `json:"tag"`
	Message string      `json:"message"`
	Param   string      `json:"param,omitempty"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success   bool              `json:"success"`
	Data      interface{}       `json:"data,omitempty"`
	Message   string            `json:"message,omitempty"`
	Errors    []ValidationError `json:"errors,omitempty"`
	ErrorCode string            `json:"error_code,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// Global data stores (in a real app, these would be databases)
var products []Product
var categories = []Category{
	{ID: 1, Name: "Electronics", Slug: "electronics"},
	{ID: 2, Name: "Clothing", Slug: "clothing"},
	{ID: 3, Name: "Books", Slug: "books"},
	{ID: 4, Name: "Home & Garden", Slug: "home-garden"},
}
var validCurrencies = []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD"}
var validWarehouses = []string{"WH001", "WH002", "WH003", "WH004", "WH005"}
var nextProductID = 1

// Implement SKU format validator
// SKU format: ABC-123-XYZ (3 letters, 3 numbers, 3 letters)
func isValidSKU(sku string) bool {
	// Implement SKU validation
	// The SKU should match the pattern: ^[A-Z]{3}-\d{3}-[A-Z]{3}$
	pattern := `^[A-Z]{3}-\d{3}-[A-Z]{3}$`
	re := regexp.MustCompile(pattern, regexp.None)
	match, _ := re.MatchString(sku)
	if !match {
		return false
	}
	return true
}

// Implement currency validator
func isValidCurrency(currency string) bool {
	// Check if the currency is in the validCurrencies slice
	return slices.Contains(validCurrencies, currency)
}

// Implement category validator
func isValidCategory(categoryName string) bool {
	// Check if the category name exists in the categories slice
	for _, v := range categories {
		if categoryName == v.Name {
			return true
		}
	}
	return false
}

// Implement slug format validator
func isValidSlug(slug string) bool {
	// Implement slug validation
	// Slug should match: ^[a-z0-9]+(?:-[a-z0-9]+)*$
	pattern := `^[a-z0-9]+(?:-[a-z0-9]+)*$`
	re := regexp.MustCompile(pattern, regexp.None)
	b, _ := re.MatchString(slug)
	return b
}

// Implement warehouse code validator
func isValidWarehouseCode(code string) bool {
	// Check if warehouse code is in validWarehouses slice
	// Format should be WH### (e.g., WH001, WH002)
	return slices.Contains(validWarehouses, code)
}

// Implement comprehensive product validation
func validateProduct(product *Product) []ValidationError {
	var errs []ValidationError

	// Add custom validation logic:
	// - Validate SKU format and uniqueness
	// - Validate currency code
	// - Validate category exists
	// - Validate slug format
	// - Validate warehouse code
	// - Cross-field validations (reserved <= quantity, etc.)
	if !isValidSKU(product.SKU) {
		errs = append(errs, ValidationError{})
	}
	if !isValidCurrency(product.Currency) {
		errs = append(errs, ValidationError{})
	}
	if !isValidCategory(product.Category.Name) {
		errs = append(errs, ValidationError{})
	}
	if !isValidSlug(product.Category.Slug) {
		errs = append(errs, ValidationError{})
	}
	if !isValidWarehouseCode(product.Inventory.Location) {
		errs = append(errs, ValidationError{})
	}
	if product.Inventory.Reserved > product.Inventory.Quantity {
		errs = append(errs, ValidationError{})
	}

	return errs
}

// Implement input sanitization
func sanitizeProduct(product *Product) {
	// Sanitize input data:
	// - Trim whitespace from strings
	// - Convert currency to uppercase
	// - Convert slug to lowercase
	// - Calculate available inventory (quantity - reserved)
	// - Set timestamps

	product.SKU = strings.Trim(product.SKU, " ")
	product.Name = strings.Trim(product.Name, " ")
	product.Currency = strings.ToUpper(product.Currency)
	product.Category.Slug = strings.ToLower(product.Category.Slug)
	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
}

// POST /products - Create single product
func createProduct(c *gin.Context) {
	var product Product

	// Bind JSON and handle basic validation errors
	if err := c.ShouldBindJSON(&product); err != nil {
		var errs []ValidationError
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				errs = append(errs, ValidationError{
					Field: e.Field(),
				})
			}
		}

		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON or basic validation failed",
			Errors:  errs, // Convert gin validation errors
		})
		return
	}

	// Apply custom validation
	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}

	// Sanitize input data
	sanitizeProduct(&product)

	// TODO: Set ID and add to products slice
	product.ID = nextProductID
	nextProductID++
	products = append(products, product)

	c.JSON(201, APIResponse{
		Success: true,
		Data:    product,
		Message: "Product created successfully",
	})
}

// POST /products/bulk - Create multiple products
func createProductsBulk(c *gin.Context) {
	var inputProducts []Product

	if err := c.ShouldBindJSON(&inputProducts); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}

	// TODO: Implement bulk validation
	type BulkResult struct {
		Index   int               `json:"index"`
		Success bool              `json:"success"`
		Product *Product          `json:"product,omitempty"`
		Errors  []ValidationError `json:"errors,omitempty"`
	}

	var results []BulkResult
	var successCount int

	// TODO: Process each product and populate results
	for i, product := range inputProducts {
		validationErrors := validateProduct(&product)
		if len(validationErrors) > 0 {
			results = append(results, BulkResult{
				Index:   i,
				Success: false,
				Errors:  validationErrors,
			})
		} else {
			sanitizeProduct(&product)
			product.ID = nextProductID
			nextProductID++
			products = append(products, product)

			results = append(results, BulkResult{
				Index:   i,
				Success: true,
				Product: &product,
			})
			successCount++
		}
	}

	c.JSON(200, APIResponse{
		Success: successCount == len(inputProducts),
		Data: map[string]interface{}{
			"results":    results,
			"total":      len(inputProducts),
			"successful": successCount,
			"failed":     len(inputProducts) - successCount,
		},
		Message: "Bulk operation completed",
	})
}

// POST /categories - Create category
func createCategory(c *gin.Context) {
	var category Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON or validation failed",
		})
		return
	}

	// TODO: Add category-specific validation
	// - Validate slug format
	// - Check parent category exists if specified
	// - Ensure category name is unique

	categories = append(categories, category)

	c.JSON(201, APIResponse{
		Success: true,
		Data:    category,
		Message: "Category created successfully",
	})
}

// POST /validate/sku - Validate SKU format and uniqueness
func validateSKUEndpoint(c *gin.Context) {
	var request struct {
		SKU string `json:"sku" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "SKU is required",
		})
		return
	}

	// Implement SKU validation endpoint
	// - Check format using isValidSKU
	if !isValidSKU(request.SKU) {
		c.JSON(200, APIResponse{
			Success: false,
		})
		return
	}

	// - Check uniqueness against existing products
	for _, p := range products {
		if p.SKU == request.SKU {
			c.JSON(200, APIResponse{
				Success: false,
			})
			return
		}

	}

	c.JSON(200, APIResponse{
		Success: true,
		Message: "SKU is valid",
	})
}

// POST /validate/product - Validate product without saving
func validateProductEndpoint(c *gin.Context) {
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}

	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Message: "Product data is valid",
	})
}

// GET /validation/rules - Get validation rules
func getValidationRules(c *gin.Context) {
	rules := map[string]interface{}{
		"sku": map[string]interface{}{
			"format":   "ABC-123-XYZ",
			"required": true,
			"unique":   true,
		},
		"name": map[string]interface{}{
			"required": true,
			"min":      3,
			"max":      100,
		},
		"currency": map[string]interface{}{
			"required": true,
			"valid":    validCurrencies,
		},
		"warehouse": map[string]interface{}{
			"format": "WH###",
			"valid":  validWarehouses,
		},
		// TODO: Add more validation rules
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    rules,
		Message: "Validation rules retrieved",
	})
}

// Setup router
func setupRouter() *gin.Engine {
	router := gin.Default()

	// Product routes
	router.POST("/products", createProduct)
	router.POST("/products/bulk", createProductsBulk)

	// Category routes
	router.POST("/categories", createCategory)

	// Validation routes
	router.POST("/validate/sku", validateSKUEndpoint)
	router.POST("/validate/product", validateProductEndpoint)
	router.GET("/validation/rules", getValidationRules)

	return router
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
