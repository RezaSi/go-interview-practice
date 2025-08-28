package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
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

// SKU format: ABC-123-XYZ (3 letters, 3 numbers, 3 letters)
// The SKU should match the pattern: ^[A-Z]{3}-\d{3}-[A-Z]{3}$
func isValidSKU(sku string) bool {
	return regexp.MustCompile(`^[A-Z]{3}-\d{3}-[A-Z]{3}$`).MatchString(sku)
}

func isValidCurrency(currency string) bool {
	return slices.Contains(validCurrencies, currency)
}

func isValidCategory(categoryName string) bool {
	for _, c := range(categories) {
		if c.Name == categoryName {
			return true
		}
	}
	return false
}

// Slug should match: ^[a-z0-9]+(?:-[a-z0-9]+)*$
func isValidSlug(slug string) bool {
	return regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`).MatchString(slug)
}

// Format should be WH### (e.g., WH001, WH002)
func isValidWarehouseCode(code string) bool {
	return slices.Contains(validWarehouses, code)
}

func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError

	// Add custom validation logic:
	// - Validate SKU format and uniqueness
	// - Validate currency code
	// - Validate category exists
	// - Validate slug format
	// - Validate warehouse code
	// - Cross-field validations (reserved <= quantity, etc.)

	if ! isValidSKU(product.SKU) {
		errors = append(errors, ValidationError{Field: "sku", Message: "Invalid SKU"})
	}
	if ! isValidCurrency(product.Currency) {
		errors = append(errors, ValidationError{Field: "currency", Message: "Invalid currency"})
	}
	if ! isValidCategory(product.Category.Name) {
		errors = append(errors, ValidationError{Field: "category.name", Message: "Category does not exist"})
	}
	if ! isValidSlug(product.Category.Slug) {
		errors = append(errors, ValidationError{Field: "category.slug", Message: "Invalid slug"})
	}
	if ! isValidWarehouseCode(product.Inventory.Location) {
		errors = append(errors, ValidationError{Field: "inventory.location", Message: "Invalid warehouse code"})
	}
	if product.Inventory.Reserved > product.Inventory.Quantity {
		errors = append(errors, ValidationError{Field: "inventory.reserved", Message: "Reserved > quantity"})
	}
	for _, p := range(products) {
		if p.SKU == product.SKU {
			errors = append(errors, ValidationError{Field: "sku", Message: "SKU already exists"})
			break
		}
	}
	return errors
}

func sanitizeProduct(product *Product) {
	// Sanitize input data:
	// - Trim whitespace from strings
	// - Convert currency to uppercase
	// - Convert slug to lowercase
	// - Calculate available inventory (quantity - reserved)
	// - Set timestamps

	product.SKU = strings.TrimSpace(product.SKU)
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)
	product.Currency = strings.ToUpper(strings.TrimSpace(product.Currency))
	product.Category.Slug = strings.ToLower(strings.TrimSpace(product.Category.Slug))
	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved

	now := time.Now().UTC()
	product.CreatedAt = now
	product.UpdatedAt = now
	product.Inventory.LastUpdated = now
}

func formatBindingErrors(err error) []ValidationError {
	var ve validator.ValidationErrors
	var result []ValidationError
	if errors.As(err, &ve) {
		for _, fe := range(ve) {
			result = append(result, ValidationError{
				Field:   fe.Field(),
				Tag:     fe.Tag(),
				Value:   fe.Value(),
				Message: fmt.Sprintf("Validation failed on field '%s' - '%s'", fe.Field(), fe.Tag()),
				Param:   fe.Param(),
			})
		}
	}
	return result
}

// POST /products - Create single product
func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid JSON or basic validation failed",
			Errors:  formatBindingErrors(err),
		})
		return
	}

	// Sanitization must be done before validation
	sanitizeProduct(&product)

	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}

	product.ID = nextProductID
	nextProductID++
	products = append(products, product)

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    product,
		Message: "Product created successfully",
	})
}

// POST /products/bulk - Create multiple products
func createProductsBulk(c *gin.Context) {
	var inputProducts []Product

	if err := c.ShouldBindJSON(&inputProducts); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}

	type BulkResult struct {
		Index   int               `json:"index"`
		Success bool              `json:"success"`
		Product *Product          `json:"product,omitempty"`
		Errors  []ValidationError `json:"errors,omitempty"`
	}

	var results []BulkResult
	var successCount int

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

	c.JSON(http.StatusOK, APIResponse{
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
	if ! isValidSlug(category.Slug) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid slug",
		})
		return
	}

	for _, cat := range(categories) {
		if cat.Name == category.Name {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success:   false,
				Message:   "Category already exists",
			})
			return
		}
	}

	if category.ParentID != nil {
		exists := false
		for _, cat := range categories {
			if cat.ID == *category.ParentID {
				exists = true
				break
			}
		}
		if ! exists {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success:   false,
				Message:   "Parent category does not exists",
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "SKU is required",
		})
		return
	}

	if ! isValidSKU(request.SKU) {
		c.JSON(http.StatusOK, APIResponse{
			Success: false,
			Message: "Invalid SKU format",
		})
		return
	}

	for _, p := range(products) {
		if p.SKU == request.SKU {
			c.JSON(http.StatusOK, APIResponse{
				Success: false,
				Message: "SKU already exists",
			})
			return
		}
	}

	c.JSON(http.StatusOK, APIResponse{Success: true, Message: "SKU is valid"})
}

// POST /validate/product - Validate product without saving
func validateProductEndpoint(c *gin.Context) {
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}

	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
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
	}

	c.JSON(http.StatusOK, APIResponse{
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
