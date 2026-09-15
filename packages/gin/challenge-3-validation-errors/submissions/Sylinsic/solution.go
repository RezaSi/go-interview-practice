package main

import (
	"regexp"
	"slices"
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

var skuRegex = regexp.MustCompile(`^[A-Z]{3}-\d{3}-[A-Z]{3}$`)
func isValidSKU(sku string) bool {
	return skuRegex.Match([]byte(sku))
}

func isValidCurrency(currency string) bool {
	return slices.Contains(validCurrencies, currency)
}

func isValidCategory(categoryName string) bool {
	if categoryName == "" { return false }

	return slices.ContainsFunc(categories, func(c Category) bool {
		return strings.EqualFold(c.Name, categoryName)
	})
}

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
func isValidSlug(slug string) bool {
	return slugRegex.Match([]byte(slug))
}

// TODO: Implement warehouse code validator
func isValidWarehouseCode(code string) bool {
	return slices.Contains(validWarehouses, code)
}

func skuExists(sku string) bool {
	return slices.ContainsFunc(products, func (p Product) bool {
		return p.SKU == sku
	})
}

// TODO: Implement comprehensive product validation
func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError

	if !isValidSKU(product.SKU) {
		errors = append(errors, ValidationError{
			Field: "SKU",
			Value: product.SKU,
			Message: "Invalid SKU format",
		})
	}else if skuExists(product.SKU) {
		errors = append(errors, ValidationError{
			Field: "SKU",
			Value: product.SKU,
			Message: "SKU already exists",
		})
	}

	if !isValidCurrency(product.Currency) {
		errors = append(errors, ValidationError{
			Field: "Currency",
			Value: product.Currency,
			Message: "Invalid currency specified",
		})
	}

	if !isValidCategory(product.Category.Name) {
		errors = append(errors, ValidationError{
			Field: "Category.Name",
			Value: product.Category.Name,
			Message: "Category doesn't exist",
		})
	}

	if !isValidSlug(product.Category.Slug) {
		errors = append(errors, ValidationError{
			Field: "Category.Slug",
			Value: product.Category.Slug,
			Message: "Invalid category slug",
		})
	}

	if !isValidWarehouseCode(product.Inventory.Location) {
		errors = append(errors, ValidationError{
			Field: "Inventory.Location",
			Value: product.Inventory.Location,
			Message: "Invalid inventory location",
		})
	}

	if product.Inventory.Reserved > product.Inventory.Quantity {
		errors = append(errors, ValidationError{
			Field: "Inventory.Reserved",
			Value: product.Inventory.Reserved,
			Message: "Cannot reserve more than in stock",
		})
	}
	// TODO: Add custom validation logic:
	// - Cross-field validations (reserved <= quantity, etc.)

	return errors
}

// TODO: Implement input sanitization
func sanitizeProduct(product *Product) {
	product.Currency = strings.ToUpper(strings.TrimSpace(product.Currency))
	product.Description = strings.TrimSpace(product.Description)
	product.Name = strings.TrimSpace(product.Name)
	product.SKU = strings.TrimSpace(product.SKU)
	
	product.Category.Name = strings.TrimSpace(product.Category.Name)
	product.Category.Slug = strings.ToLower(strings.TrimSpace(product.Category.Slug))

	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
}

// POST /products - Create single product
func createProduct(c *gin.Context) {
	var product Product

	// TODO: Bind JSON and handle basic validation errors
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON or basic validation failed",
			Errors:  []ValidationError{
				ValidationError{
					Message: err.Error(),
				},
			},
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
	if !isValidSlug(category.Slug) {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid slug",
		})
		return
	} else if category.ParentID != nil {
		if !slices.ContainsFunc(categories, func (c Category) bool {
			return c.ID == *category.ParentID
		}) {
			c.JSON(400, APIResponse{
				Success: false,
				Message: "Parent category does not exist",
			})
			return
		}
	} else if isValidCategory(category.Name) {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Category already exists",
		})
		return
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

	if !isValidSKU(request.SKU) {
		c.JSON(200, APIResponse{
			Success: false,
			Message: "Invalid SKU format",
			Errors: []ValidationError{
				{
					Field: "sku",
					Value: request.SKU,
					Message: "Invalid SKU format",
				},
			},
		})
		return
	} else if skuExists(request.SKU) {
		c.JSON(200, APIResponse{
			Success: false,
			Message: "SKU already used",
			Errors: []ValidationError{
				{
					Field: "sku",
					Value: request.SKU,
					Message: "SKU already used",
				},
			},
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
