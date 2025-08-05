package main

import (
	"time"
    "regexp"
    "slices"
    "strings"
    
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

// SKU format: ABC-123-XYZ (3 letters, 3 numbers, 3 letters)
func isValidSKU(sku string) bool {
	pattern := `^[A-Z]{3}-\d{3}-[A-Z]{3}$`
	matched, _ := regexp.MatchString(pattern, sku)
	return matched
}

func isValidCurrency(currency string) bool {
	return slices.Contains(validCurrencies, currency)
}

func isValidCategory(categoryName string) bool {
	return slices.ContainsFunc(categories, func (category Category) bool {
	    return category.Name == categoryName
	})
}

func isValidSlug(slug string) bool {
	pattern := `^[a-z0-9]+(?:-[a-z0-9]+)*$`
	matched, _ := regexp.MatchString(pattern, slug)
	return matched
}

func isValidWarehouseCode(code string) bool {
	return slices.Contains(validWarehouses, code)
}

// TODO: Implement comprehensive product validation
func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError

	// TODO: Add custom validation logic:
	if product.Name == "" {
	    errors = append(errors, ValidationError {
	        Field: "Name",
	        Value: product.Name,
	        Tag: "not-found",
	        Message: "No name provided",
	    })
	}
	if product.Price == 0.0 {
	    errors = append(errors, ValidationError {
	        Field: "Price",
	        Value: product.Price,
	        Tag: "not-found",
	        Message: "No price provided",
	    })
	}
	if !isValidSKU(product.SKU) || product.SKU == "" {
	    errors = append(errors, ValidationError {
	        Field: "SKU",
	        Value: product.SKU,
	        Tag: "sku-format",
	        Message: "Invalid SKU",
	    })
	}
	if !isValidCurrency(product.Currency) || product.Currency == "" {
	    errors = append(errors, ValidationError {
	        Field: "Currency",
	        Value: product.Currency,
	        Tag: "currency",
	        Message: "Invalid Currency",
	    })
	}
	if !isValidCategory(product.Category.Name) || product.Category.Name == "" {
	    errors = append(errors, ValidationError {
	        Field: "Category",
	        Value: product.Category.Name,
	        Tag: "not-found",
	        Message: "Category not found",
	    })
	}
	if !isValidSlug(product.Category.Slug) || product.Category.Slug == "" {
	    errors = append(errors, ValidationError {
	        Field: "Slug",
	        Value: product.Category.Slug,
	        Tag: "slug-format",
	        Message: "Invalid slug",
	    })
	}
	if product.Inventory.Location == "" {
	    errors = append(errors, ValidationError {
	        Field: "product.Inventory",
	        Value: product.Inventory,
	        Tag: "not-found",
	        Message: "No inventory provided",
	    })
	}
	// - Validate warehouse code
	// - Cross-field validations (reserved <= quantity, etc.)
    if product.Inventory.Reserved > product.Inventory.Quantity {
        errors = append(errors, ValidationError {
            Field: "product.Inventory.Reserved",
            Value: product.Inventory.Reserved,
            Tag: "max",
            Message: "Reserved exceeds quantity",
        })
    }
	return errors
}

// TODO: Implement input sanitization
func sanitizeProduct(product *Product) {
	// TODO: Sanitize input data:
	product.SKU = strings.TrimSpace(product.SKU)
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)
	product.Currency = strings.TrimSpace(product.Currency)
	
	// - Convert currency to uppercase
	product.Currency = strings.ToUpper(product.Currency)
	// - Convert slug to lowercase
	product.Category.Slug = strings.ToLower(product.Category.Slug)
	// - Calculate available inventory (quantity - reserved)
	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
	// - Set timestamps
	if product.ID == 0 {
	    product.CreatedAt = time.Now()
	}
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
			    {
			        Message: "Missing Required Fields",
			    },
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

	// TODO: Implement SKU validation endpoint
	// - Check format using isValidSKU
	// - Check uniqueness against existing products
	if !isValidSKU(request.SKU) {
	    c.JSON(200, APIResponse {
	        Success: false,
	        Message: "Invalid SKU",
	    })
	    return
	}
	if slices.ContainsFunc(products, func (product Product) bool {
	    return product.SKU == request.SKU
	}) {
	    c.JSON(200, APIResponse {
	        Success: false,
	        Message: "SKU already exists",
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
