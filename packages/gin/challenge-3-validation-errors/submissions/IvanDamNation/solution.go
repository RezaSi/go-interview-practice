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
	"github.com/go-playground/validator/v10"
)

var (
	skuRegex  = regexp.MustCompile(`^[A-Z]{3}-\d{3}-[A-Z]{3}$`)
	slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

// Product represents a product in the catalog
type Product struct {
	ID          int            `json:"id"`
	SKU         string         `json:"sku" binding:"required"`
	Name        string         `json:"name" binding:"required,min=3,max=100"`
	Description string         `json:"description" binding:"max=1000"`
	Price       float64        `json:"price" binding:"required,min=0.01"`
	Currency    string         `json:"currency" binding:"required"`
	Category    Category       `json:"category" binding:"required"`
	Tags        []string       `json:"tags"`
	Attributes  map[string]any `json:"attributes"`
	Images      []Image        `json:"images"`
	Inventory   Inventory      `json:"inventory" binding:"required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
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
	Quantity    int       `json:"quantity" binding:"min=0"`
	Reserved    int       `json:"reserved" binding:"min=0"`
	Available   int       `json:"available"` // Calculated field
	Location    string    `json:"location" binding:"required"`
	LastUpdated time.Time `json:"last_updated"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Value   any    `json:"value"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
	Param   string `json:"param,omitempty"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success   bool              `json:"success"`
	Data      any               `json:"data,omitempty"`
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

// isValidSKU SKU format validator
// SKU format: ABC-123-XYZ (3 letters, 3 numbers, 3 letters)
func isValidSKU(sku string) bool {
	// The SKU should match the pattern: ^[A-Z]{3}-\d{3}-[A-Z]{3}$
	return skuRegex.MatchString(sku)
}

// isValidCurrency currency validator
func isValidCurrency(currency string) bool {
	return slices.Contains(validCurrencies, currency)
}

// isValidCategory category validator
func isValidCategory(categoryName string) bool {
	for _, cat := range categories {
		if cat.Name == categoryName {
			return true
		}
	}
	return false
}

// isValidSlug slug format validator
func isValidSlug(slug string) bool {
	// Slug should match: ^[a-z0-9]+(?:-[a-z0-9]+)*$
	return slugRegex.MatchString(slug)
}

// isValidWarehouseCode warehouse code validator
func isValidWarehouseCode(code string) bool {
	// Format should be WH### (e.g., WH001, WH002)
	return slices.Contains(validWarehouses, code)
}

// validateProduct is aggregated product validation
func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError

	if len(product.Name) < 3 || len(product.Name) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Value:   product.Name,
			Message: "Name must be between 3 and 100 characters after sanitization",
		})
	}

	if !isValidSKU(product.SKU) {
		errors = append(errors, ValidationError{
			Field:   "sku",
			Value:   product.SKU,
			Message: "check SKU format (valid format: ABC-123-XYZ)",
		})
	}

	if !isValidCurrency(product.Currency) {
		errors = append(errors, ValidationError{
			Field:   "currency",
			Value:   product.Currency,
			Message: "unsupported currency",
		})
	}

	if !isValidCategory(product.Category.Name) {
		errors = append(errors, ValidationError{
			Field:   "category.name",
			Value:   product.Category,
			Message: "unexisting category",
		})
	}

	if !isValidSlug(product.Category.Slug) {
		errors = append(errors, ValidationError{
			Field:   "slug",
			Value:   product.Category.Slug,
			Message: "unsupported slug format",
		})
	}

	if !isValidWarehouseCode(product.Inventory.Location) {
		errors = append(errors, ValidationError{
			Field:   "location",
			Value:   product.Inventory.Location,
			Message: "invalid warehouse code (check format: WH### (e.g., WH001, WH002))",
		})
	}

	if product.Inventory.Quantity < product.Inventory.Reserved {
		errors = append(errors, ValidationError{
			Field:   "inventory.reserved",
			Message: "reserved cannot be more, than quantity",
		})
	}

	return errors
}

// sanitizeProduct sanitizing input data
func sanitizeProduct(product *Product) {
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)
	product.SKU = strings.TrimSpace(product.SKU)

	product.Category.Name = strings.TrimSpace(product.Category.Name)
	product.Category.Slug = strings.ToLower(product.Category.Slug)

	product.Currency = strings.ToUpper(product.Currency)
	product.Inventory.Location = strings.TrimSpace(product.Inventory.Location)

	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
}

// POST /products - Create single product
func createProduct(c *gin.Context) {
	var product Product
	var validationErrors []ValidationError

	if c.ContentType() != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, APIResponse{
			Message: "JSON expected",
		})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, e := range ve {
				field := strings.ToLower(e.Field())
				validationErrors = append(validationErrors, ValidationError{
					Field:   field,
					Message: fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", field, e.Tag()),
				})
			}
			c.JSON(http.StatusBadRequest, APIResponse{
				Message: "Validation failed",
				Errors:  validationErrors,
			})
			return
		}

		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Invalid JSON format",
		})
		return
	}

	sanitizeProduct(&product)

	validationErrors = validateProduct(&product)
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
			Message: "Invalid JSON format",
		})
		return
	}

	type BulkItemResult struct {
		Success bool              `json:"success"`
		Errors  []ValidationError `json:"errors,omitempty"`
	}

	var results []BulkItemResult
	var successCount int

	validate := validator.New()
	validate.SetTagName("binding")

	for _, product := range inputProducts {
		sanitizeProduct(&product)

		var itemErrors []ValidationError

		if err := validate.Struct(product); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				for _, e := range ve {
					field := strings.ToLower(e.Field())
					itemErrors = append(itemErrors, ValidationError{
						Field:   field,
						Message: fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", field, e.Tag()),
					})
				}
			}
		}

		businessErrors := validateProduct(&product)
		if len(businessErrors) > 0 {
			itemErrors = append(itemErrors, businessErrors...)
		}

		if len(itemErrors) > 0 {
			results = append(results, BulkItemResult{
				Errors: itemErrors,
			})
		} else {
			product.ID = nextProductID
			nextProductID++
			products = append(products, product)

			results = append(results, BulkItemResult{
				Success: true,
			})
			successCount++
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: successCount == len(inputProducts),
		Data: map[string]any{
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Invalid JSON or validation failed",
		})
		return
	}

	if !isValidSlug(category.Slug) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Invalid slug format",
		})
		return
	}

	hasParent := category.ParentID != nil
	parentExists := !hasParent

	uniq := true
	for _, cat := range categories {
		if cat.Name == category.Name {
			uniq = false
		}
		if hasParent && cat.ID == *category.ParentID {
			parentExists = true
		}
	}

	if !uniq {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Category name is not unique",
		})
		return
	}

	if !parentExists {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Parent category does not exist",
		})
		return
	}

	categories = append(categories, category)

	c.JSON(http.StatusCreated, APIResponse{
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
			Message: "SKU is required",
		})
		return
	}

	request.SKU = strings.TrimSpace(request.SKU)

	if !isValidSKU(request.SKU) {
		c.JSON(http.StatusOK, APIResponse{
			Message: "Invalid format SKU",
		})
		return
	}

	for _, p := range products {
		if p.SKU == request.SKU {
			c.JSON(http.StatusOK, APIResponse{
				Message: "SKU alredy exists",
			})
			return
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "SKU is valid",
	})
}

// POST /validate/product - Validate product without saving
func validateProductEndpoint(c *gin.Context) {
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Invalid JSON format",
		})
		return
	}

	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
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
	rules := map[string]any{
		"sku": map[string]any{
			"format":   "ABC-123-XYZ",
			"required": true,
			"unique":   true,
		},
		"name": map[string]any{
			"required": true,
			"min":      3,
			"max":      100,
		},
		"currency": map[string]any{
			"required": true,
			"valid":    validCurrencies,
		},
		"warehouse": map[string]any{
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

	router.POST("/products", createProduct)
	router.POST("/products/bulk", createProductsBulk)

	router.POST("/categories", createCategory)

	router.POST("/validate/sku", validateSKUEndpoint)
	router.POST("/validate/product", validateProductEndpoint)
	router.GET("/validation/rules", getValidationRules)

	return router
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
