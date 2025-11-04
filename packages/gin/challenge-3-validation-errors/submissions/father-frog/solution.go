package main

import (
	"regexp"
	"strings"
	"sync"
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
	ID       int    `json:"id"`
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
var (
	productsMu    sync.RWMutex
	products      = []Product{}
	nextProductID = 1

	categoriesMu sync.RWMutex
	categories   = []Category{
		{ID: 1, Name: "Electronics", Slug: "electronics"},
		{ID: 2, Name: "Clothing", Slug: "clothing"},
		{ID: 3, Name: "Books", Slug: "books"},
		{ID: 4, Name: "Home & Garden", Slug: "home-garden"},
	}
	nextCategoryID = 5

	validCurrencies = []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD"}
	// valid warehouses : format should be WH### (e.g., WH001, WH002)
	validWarehouses = []string{"WH001", "WH002", "WH003", "WH004", "WH005"}

	// SKU format: ABC-123-XYZ (3 letters, 3 numbers, 3 letters)
	skuRegPattern = `^[A-Z]{3}-\d{3}-[A-Z]{3}$`
	skuReg        = regexp.MustCompile(skuRegPattern)

	// Slug format: ^[a-z0-9]+(?:-[a-z0-9]+)*$
	slugRegPattern = `^[a-z0-9]+(?:-[a-z0-9]+)*$`
	slugReg        = regexp.MustCompile(slugRegPattern)
)

// isValidSKU returns true if the provided sku matches the SKU format: ABC-123-XYZ (3 letters, 3 numbers, 3 letters)
func isValidSKU(sku string) bool {
	return skuReg.MatchString(sku)
}

// isValidCurrency returns true if the currency is in the validCurrencies slice
func isValidCurrency(currency string) bool {
	// Check if the currency is in the validCurrencies slice
	for _, c := range validCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}

// isValidCategory returns true if the categoryName is in the categories slice
func isValidCategory(categoryName string) bool {
	categoriesMu.RLock()
	defer categoriesMu.RUnlock()
	// Check if the category name exists in the categories slice
	for _, c := range categories {
		if c.Name == categoryName {
			return true
		}
	}
	return false
}

// isValidSlug returns true if the slug matches the Slug format: ^[a-z0-9]+(?:-[a-z0-9]+)*$
func isValidSlug(slug string) bool {
	return slugReg.MatchString(slug)
}

// returns true if the warehouse is in validWarehouses
func isValidWarehouseCode(code string) bool {
	// Check if warehouse code is in validWarehouses slice
	for _, w := range validWarehouses {
		if w == code {
			return true
		}
	}
	return false
}

// Implement comprehensive product validation
func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError

	// Validate SKU format
	if !isValidSKU(product.SKU) {
		errors = append(errors, ValidationError{
			Field:   "sku",
			Value:   product.SKU,
			Tag:     "sku_format",
			Message: "SKU must match the format: XXX-###-XXX (3 letters, 3 numbers, 3 letters)",
		})
	}

	// Validate SKU uniqueness (check against existing products)
	productsMu.RLock()
	for _, p := range products {
		if p.SKU == product.SKU {
			errors = append(errors, ValidationError{
				Field:   "sku",
				Value:   product.SKU,
				Tag:     "sku_unique",
				Message: "SKU must be unique",
			})
			break
		}
	}
	productsMu.RUnlock()

	// Validate currency
	if !isValidCurrency(product.Currency) {
		errors = append(errors, ValidationError{
			Field:   "currency",
			Value:   product.Currency,
			Tag:     "currency_valid",
			Message: "Currency must be one of: USD, EUR, GBP, JPY, CAD, AUD",
		})
	}

	// Validate category exists
	if !isValidCategory(product.Category.Name) {
		errors = append(errors, ValidationError{
			Field:   "category.name",
			Value:   product.Category.Name,
			Tag:     "category_exists",
			Message: "Category must be a valid existing category",
		})
	}

	// Validate slug format
	if !isValidSlug(product.Category.Slug) {
		errors = append(errors, ValidationError{
			Field:   "category.slug",
			Value:   product.Category.Slug,
			Tag:     "slug_format",
			Message: "Slug must match the format: lowercase letters and numbers separated by hyphens",
		})
	}

	// Validate warehouse code
	if !isValidWarehouseCode(product.Inventory.Location) {
		errors = append(errors, ValidationError{
			Field:   "inventory.location",
			Value:   product.Inventory.Location,
			Tag:     "warehouse_valid",
			Message: "Warehouse location must be one of: WH001, WH002, WH003, WH004, WH005",
		})
	}

	// Cross-field validations
	if product.Inventory.Reserved > product.Inventory.Quantity {
		errors = append(errors, ValidationError{
			Field:   "inventory.reserved",
			Value:   product.Inventory.Reserved,
			Tag:     "reserved_less_than_quantity",
			Message: "Reserved quantity cannot exceed total quantity",
		})
	}

	return errors
}

// Sanitize input data:
// - Trim whitespace from strings
// - Convert currency to uppercase
// - Convert slug to lowercase
// - Calculate available inventory (quantity - reserved)
// - Set timestamps
func sanitizeProduct(product *Product) {
	// - Trim whitespace from strings
	product.SKU = strings.TrimSpace(product.SKU)
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)
	product.Currency = strings.TrimSpace(product.Currency)
	product.Category.Name = strings.TrimSpace(product.Category.Name)
	product.Category.Slug = strings.TrimSpace(product.Category.Slug)
	product.Inventory.Location = strings.TrimSpace(product.Inventory.Location)
	// - Convert currency to uppercase
	product.Currency = strings.ToUpper(product.Currency)
	// - Convert slug to lowercase
	product.Category.Slug = strings.ToLower(product.Category.Slug)
	// - Calculate available inventory (quantity - reserved)
	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
	// - Set timestamps
	now := time.Now()
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}
	product.UpdatedAt = now
	product.Inventory.LastUpdated = now
}

// POST /products - Create single product
func createProduct(c *gin.Context) {
	var product Product

	// Bind JSON and handle basic validation errors
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON or basic validation failed",
			Errors: []ValidationError{
				{
					Tag:     "bind",
					Message: err.Error(),
				},
			},
		})
		return
	}

	// Sanitize input data
	sanitizeProduct(&product)

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

	// Set ID and add to products slice
	productsMu.Lock()
	defer productsMu.Unlock()
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

	// Implement bulk validation
	type BulkResult struct {
		Index   int               `json:"index"`
		Success bool              `json:"success"`
		Product *Product          `json:"product,omitempty"`
		Errors  []ValidationError `json:"errors,omitempty"`
	}

	var results []BulkResult
	var successCount int

	// Process each product and populate results
	for i, product := range inputProducts {
		// Sanitize products before validating
		sanitizeProduct(&product)

		// Now we have consistent data to check for duplicates
		validationErrors := validateProduct(&product)
		if len(validationErrors) > 0 {
			results = append(results, BulkResult{
				Index:   i,
				Success: false,
				Errors:  validationErrors,
			})
		} else {
			productsMu.Lock()
			product.ID = nextProductID
			nextProductID++
			products = append(products, product)
			productsMu.Unlock()

			productCopy := product
			results = append(results, BulkResult{
				Index:   i,
				Success: true,
				Product: &productCopy,
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

	categoriesMu.Lock()
	defer categoriesMu.Unlock()

	// - Validate slug format
	if !isValidSlug(category.Slug) {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Slug is invalid",
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

	category.ID = nextCategoryID
	nextCategoryID++
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

	var errors []ValidationError

	// Validate SKU format
	if !isValidSKU(request.SKU) {
		errors = append(errors, ValidationError{
			Field:   "sku",
			Value:   request.SKU,
			Tag:     "sku_format",
			Message: "SKU must match the format: XXX-###-XXX (3 letters, 3 numbers, 3 letters)",
		})
	}

	// Validate SKU uniqueness (check against existing products)
	productsMu.RLock()
	for _, p := range products {
		if p.SKU == request.SKU {
			errors = append(errors, ValidationError{
				Field:   "sku",
				Value:   request.SKU,
				Tag:     "sku_unique",
				Message: "SKU must be unique",
			})
			break
		}
	}
	productsMu.RUnlock()

	if len(errors) != 0 {
		c.JSON(200, APIResponse{
			Success: false,
			Message: "SKU is invalid",
			Errors:  errors,
		})
		return
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

	// sanitize before validating
	sanitizeProduct(&product)

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
