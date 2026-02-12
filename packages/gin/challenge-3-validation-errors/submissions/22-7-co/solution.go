package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
var products = []Product{}
var categories = []Category{
	{ID: 1, Name: "Electronics", Slug: "electronics"},
	{ID: 2, Name: "Clothing", Slug: "clothing"},
	{ID: 3, Name: "Books", Slug: "books"},
	{ID: 4, Name: "Home & Garden", Slug: "home-garden"},
}
var validCurrencies = []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD"}
var validWarehouses = []string{"WH001", "WH002", "WH003", "WH004", "WH005"}
var nextProductID = 1

// TODO: Implement SKU format validator
// SKU format: ABC-123-XYZ (3 letters, 3 numbers, 3 letters)
func isValidSKU(sku string) bool {
	// TODO: Implement SKU validation
	// The SKU should match the pattern: ^[A-Z]{3}-\d{3}-[A-Z]{3}$
	const skuRegPattern = `^[A-Z]{3}-\d{3}-[A-Z]{3}$`
	var skuReg = regexp.MustCompile(skuRegPattern)
	if skuReg.MatchString(sku) {
		return true
	}
	return false
}

// TODO: Implement currency validator
func isValidCurrency(currency string) bool {
	// TODO: Check if the currency is in the validCurrencies slice
	if currency == "" {
		return false
	}
	for _, validCurrency := range validCurrencies {
		if currency == validCurrency {
			return true
		}
	}
	return false
}

// TODO: Implement category validator
func isValidCategory(categoryName string) bool {
	// TODO: Check if the category name exists in the categories slice
	if categoryName == "" {
		return false
	}
	for _, category := range categories {
		if category.Name == categoryName {
			return true
		}
	}
	return false
}

// TODO: Implement slug format validator
func isValidSlug(slug string) bool {
	// TODO: Implement slug validation
	// Slug should match: ^[a-z0-9]+(?:-[a-z0-9]+)*$
	const slugRegPattern = `^[a-z0-9]+(?:-[a-z0-9]+)*$`
	var slugReg = regexp.MustCompile(slugRegPattern)
	if slugReg.MatchString(slug) {
		return true
	}
	return false
}

// TODO: Implement warehouse code validator
func isValidWarehouseCode(code string) bool {
	// TODO: Check if warehouse code is in validWarehouses slice
	// Format should be WH### (e.g., WH001, WH002)
	for _, warehouseCode := range validWarehouses {
		if warehouseCode == code {
			return true
		}
	}
	return false
}

// TODO: Implement comprehensive product validation
func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError

	// TODO: Add custom validation logic:
	// - Validate SKU format and uniqueness
	if !isValidSKU(product.SKU) {
		errors = append(errors, ValidationError{
			Field:   "SKU",
			Message: fmt.Sprintf("SKU %s is not valid", product.SKU),
		})
	} else {
		for _, existing := range products {
			if product.SKU == existing.SKU {
				errors = append(errors, ValidationError{
					Field:   "SKU",
					Message: fmt.Sprintf("SKU %s is not valid", product.SKU),
				})
				break
			}
		}
	}
	// - Validate currency code
	if !isValidCurrency(product.Currency) {
		errors = append(errors, ValidationError{
			Field:   "Currency",
			Message: fmt.Sprintf("Currency %s is not valid", product.Currency),
		})
	}
	// - Validate category exists
	if !isValidCategory(product.Category.Name) {
		errors = append(errors, ValidationError{
			Field:   "Category",
			Message: fmt.Sprintf("Category %s is not valid", product.Category.Name),
		})
	}
	// - Validate slug format
	if !isValidSlug(product.Category.Slug) {
		errors = append(errors, ValidationError{
			Field:   "Category",
			Message: fmt.Sprintf("Category %s is not valid", product.Category.Name),
		})
	}
	// - Validate warehouse code
	if !isValidWarehouseCode(product.Inventory.Location) {
		errors = append(errors, ValidationError{
			Field:   "Inventory",
			Message: fmt.Sprintf("Inventory %s is not valid", product.Inventory.Location),
		})
	}
	// - Cross-field validations (reserved <= quantity, etc.)
	if product.Inventory.Reserved > product.Inventory.Quantity {
		errors = append(errors, ValidationError{
			Field:   "Inventory",
			Message: fmt.Sprintf("Inventory %s is not valid", product.Inventory.Location),
		})
	}
	return errors
}

// TODO: Implement input sanitization
func sanitizeProduct(product *Product) {
	// TODO: Sanitize input data:
	// - Trim whitespace from strings
	product.SKU = strings.TrimSpace(product.SKU)
	product.Name = strings.TrimSpace(product.Name)
	product.Currency = strings.TrimSpace(product.Currency)
	product.Description = strings.TrimSpace(product.Description)
	product.Category.Name = strings.TrimSpace(product.Category.Name)
	product.Category.Slug = strings.TrimSpace(product.Category.Slug)
	for k, v := range product.Attributes {
		if s, ok := v.(string); ok {
			product.Attributes[k] = strings.TrimSpace(s)
		}
	}
	// - Convert SKU to uppercase
	product.SKU = strings.ToUpper(product.SKU)
	// - Convert currency to uppercase
	product.Currency = strings.ToUpper(product.Currency)
	// - Convert slug to lowercase
	product.Category.Slug = strings.ToLower(product.Category.Slug)
	// - Calculate available inventory (quantity - reserved)
	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
	// - Set timestamps
	now := time.Now()
	product.Inventory.LastUpdated = now
	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}
	product.UpdatedAt = now
}

// POST /products - Create single product
func createProduct(c *gin.Context) {
	var product Product

	// TODO: Bind JSON and handle basic validation errors
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON or basic validation failed",
			Errors: []ValidationError{
				ValidationError{Tag: "bind",
					Message: "Invalid JSON or basic validation failed"},
			}, // TODO: Convert gin validation errors
		})
		return
	}

	// TODO: Apply custom validation
	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}

	// TODO: Sanitize input data
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
	if !isValidSlug(category.Slug) {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid Slug",
		})
		return
	}
	// - Check parent category exists if specified

	if category.ParentID != nil {
		ok := false
		for _, existing := range categories {
			if existing.ID == *category.ParentID {
				ok = true
				break
			}
		}
		if !ok {
			c.JSON(400, APIResponse{
				Success: false,
				Message: "ParentCategory not found",
			})
			return
		}
	}

	// - Ensure category name is unique
	for _, existing := range categories {
		if existing.Name == category.Name {
			c.JSON(400, APIResponse{
				Success: false,
				Message: "Category already exists",
			})
			return
		}
	}
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

	// TODO: Implement SKU validation endpoint
	// - Check format using isValidSKU
	if !isValidSKU(request.SKU) {
		c.JSON(200, APIResponse{
			Success: false,
			Message: "Invalid SKU",
		})
		return
	}
	// - Check uniqueness against existing products
	for _, existing := range products {
		if existing.SKU == request.SKU {
			c.JSON(200, APIResponse{
				Success: false,
				Message: "SKU already exists",
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
