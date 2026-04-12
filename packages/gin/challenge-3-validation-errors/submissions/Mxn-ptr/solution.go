package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
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
	re := regexp.MustCompile(`^[A-Z]{3}-\d{3}-[A-Z]{3}$`)
	return re.MatchString(sku)
}

// TODO: Implement currency validator
func isValidCurrency(currency string) bool {
	for _, v := range validCurrencies {
		if currency == v {
			return true
		}
	}
	return false
}

// TODO: Implement category validator
func isValidCategory(categoryName string) bool {
	for _, v := range categories {
		if categoryName == v.Name {
			return true
		}
	}
	return false
}

// TODO: Implement slug format validator
func isValidSlug(slug string) bool {
	re := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return re.MatchString(slug)
}

// TODO: Implement warehouse code validator
func isValidWarehouseCode(code string) bool {
	for _, v := range validWarehouses {
		if code == v {
			return true
		}
	}
	return false
}

func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError

	// Validate SKU format
	if !isValidSKU(product.SKU) {
		errors = append(errors, ValidationError{
			Field:   "sku",
			Value:   product.SKU,
			Tag:     "sku_format",
			Message: "SKU must match format ABC-123-XYZ (3 uppercase letters, 3 digits, 3 uppercase letters)",
		})
	} else {
		// Validate SKU uniqueness
		for _, p := range products {
			if p.SKU == product.SKU {
				errors = append(errors, ValidationError{
					Field:   "sku",
					Value:   product.SKU,
					Tag:     "sku_unique",
					Message: "SKU already exists",
				})
				break
			}
		}
	}

	// Validate currency
	if !isValidCurrency(product.Currency) {
		errors = append(errors, ValidationError{
			Field:   "currency",
			Value:   product.Currency,
			Tag:     "currency_valid",
			Message: "Currency must be one of: USD, EUR, GBP, JPY, CAD, AUD",
		})
	}

	// Validate category name
	if !isValidCategory(product.Category.Name) {
		errors = append(errors, ValidationError{
			Field:   "category.name",
			Value:   product.Category.Name,
			Tag:     "category_exists",
			Message: "Category does not exist",
		})
	}

	// Validate category slug format
	if !isValidSlug(product.Category.Slug) {
		errors = append(errors, ValidationError{
			Field:   "category.slug",
			Value:   product.Category.Slug,
			Tag:     "slug_format",
			Message: "Slug must be lowercase alphanumeric with hyphens (e.g. home-garden)",
		})
	}

	// Validate warehouse code
	if !isValidWarehouseCode(product.Inventory.Location) {
		errors = append(errors, ValidationError{
			Field:   "inventory.location",
			Value:   product.Inventory.Location,
			Tag:     "warehouse_valid",
			Message: "Warehouse code must be one of: WH001, WH002, WH003, WH004, WH005",
		})
	}

	// Cross-field: reserved must not exceed quantity
	if product.Inventory.Reserved > product.Inventory.Quantity {
		errors = append(errors, ValidationError{
			Field:   "inventory.reserved",
			Value:   product.Inventory.Reserved,
			Tag:     "reserved_lte_quantity",
			Message: "Reserved quantity cannot exceed total quantity",
		})
	}

	return errors
}

// TODO: Implement input sanitization
func sanitizeProduct(product *Product) {

	product.Name = strings.TrimSpace(product.Name)
	product.SKU = strings.TrimSpace(product.SKU)

	product.Currency = strings.ToUpper(product.Currency)
	product.Category.Slug = strings.ToLower(product.Category.Slug)

	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved

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
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var validationErrors []ValidationError
			for _, fe := range ve {
				validationErrors = append(validationErrors, ValidationError{
					Field:   fe.Field(),
					Tag:     fe.Tag(),
					Value:   fe.Value(),
					Param:   fe.Param(),
					Message: fmt.Sprintf("'%s' validation returns error '%s' with value '%s' for the param '%s", fe.Field(), fe.Tag(), fe.Value(), fe.Param()),
				})
			}
			c.AbortWithStatusJSON(400, APIResponse{
				Success:   false,
				Errors:    validationErrors,
				ErrorCode: "400",
				Message:   "validation failed",
			})
		} else {
			c.AbortWithStatusJSON(400, APIResponse{
				Success: false,
				Message: "Invalid JSON format",
			})
		}
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

	var validationErrors []ValidationError

	if !isValidSlug(category.Slug) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "slug",
			Value:   category.Slug,
			Tag:     "slug_format",
			Message: "Slug must be lowercase alphanumeric with hyphens (e.g. home-garden)",
		})
	}

	if category.ParentID != nil && *category.ParentID != 0 {
		exists := false
		for _, v := range categories {
			if *category.ParentID == v.ID {
				exists = true
			}
		}
		if !exists {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "parentId",
				Value:   category.ParentID,
				Tag:     "parentId_format",
				Message: "ParentId must pointed to an category ID",
			})
		}
	}

	for _, v := range categories {
		if v.Name == category.Name {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "name",
				Value:   category.Name,
				Tag:     "name_format",
				Message: "Name must be unique",
			})
		}
	}

	if len(validationErrors) > 0 {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
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

	var validationErrors []ValidationError

	if !isValidSKU(request.SKU) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "sky",
			Value:   request.SKU,
			Tag:     "sku_format",
			Message: "SKU must looks like ABC-123-XYZ (3 letters, 3 numbers, 3 letters)",
		})
	}

	for _, v := range products {
		if v.SKU == request.SKU {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "sky",
				Value:   request.SKU,
				Tag:     "sku_format",
				Message: "SKU must be unique",
			})
		}
	}

	if len(validationErrors) > 0 {
		c.JSON(200, APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
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
		"price": map[string]interface{}{
			"required": true,
			"min":      0.1,
		},
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
