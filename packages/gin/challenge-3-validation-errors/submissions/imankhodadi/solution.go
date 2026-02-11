package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var productsMutex sync.Mutex
var categoriesMutex sync.Mutex

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
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Tag     string      `json:"tag"`
	Message string      `json:"message"`
	Param   string      `json:"param,omitempty"`
}
type APIResponse struct {
	Success   bool              `json:"success"`
	Data      interface{}       `json:"data,omitempty"`
	Message   string            `json:"message,omitempty"`
	Errors    []ValidationError `json:"errors,omitempty"`
	ErrorCode string            `json:"error_code,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

var products = []Product{}
var categories = []Category{
	{ID: 1, Name: "Electronics", Slug: "electronics"},
	{ID: 2, Name: "Clothing", Slug: "clothing"},
	{ID: 3, Name: "Books", Slug: "books"},
	{ID: 4, Name: "Home & Garden", Slug: "home-garden"},
}
var validCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"GBP": true,
	"JPY": true,
}

var validWarehouses = []string{"WH001", "WH002", "WH003", "WH004", "WH005"}
var nextProductID = 1

var (
	skuRegex       = regexp.MustCompile(`^[A-Z]{3}-\d{3}-[A-Z]{3}$`) // SKU format: ABC-123-XYZ (3 letters, 3 numbers, 3 letters)
	slugRegex      = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	warehouseRegex = regexp.MustCompile(`^WH\d{3}$`)
)

func isValidSKU(sku string) bool {
	return skuRegex.MatchString(sku)
}

func isValidCurrency(currency string) bool {
	return validCurrencies[currency]
}

func isValidCategory(categoryName string) bool {
	for _, category := range categories {
		if categoryName == category.Name {
			return true
		}
	}
	return false
}

func isValidSlug(slug string) bool {
	return slugRegex.MatchString(slug)
}

func isValidWarehouseCode(code string) bool {
	matched := warehouseRegex.MatchString(code)
	if !matched {
		return false
	}
	for _, valid := range validWarehouses {
		if code == valid {
			return true
		}
	}
	return false
}

// func validateProductByCategory(product *Product) []ValidationError {
// 	var errors []ValidationError

// 	switch product.Category.Name {
// 	case "Electronics":
// 		// Electronics need warranty info
// 		if len(product.Images) == 0 {
// 			errors = append(errors, ValidationError{
// 				Field:   "images",
// 				Message: "Electronics must have product images",
// 			})
// 		}
// 	case "Clothing":
// 		// Clothing needs size information
// 		if _, hasSize := product.Attributes["size"]; !hasSize {
// 			errors = append(errors, ValidationError{
// 				Field:   "attributes.size",
// 				Message: "Clothing must specify size",
// 			})
// 		}
// 	}

// 	return errors
// }

func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError
	// categoryErrors := validateProductByCategory(product)
	// errors = append(errors, categoryErrors...)
	if !isValidSKU(product.SKU) {
		errors = append(errors, ValidationError{
			Field:   "sku",
			Message: "SKU must follow ABC-123-XYZ format",
		})
	}
	if !isValidCurrency(product.Currency) {
		errors = append(errors, ValidationError{
			Field:   "currency",
			Message: "Must be a valid ISO 4217 currency code",
		})
	}

	if !isValidCategory(product.Category.Name) {
		errors = append(errors, ValidationError{
			Field:   "product category name",
			Message: "invalid category name",
		})
	}
	if !isValidSlug(product.Category.Slug) {
		errors = append(errors, ValidationError{
			Field:   "product category slug",
			Message: "invalid category slug",
		})
	}

	if product.Inventory.Reserved > product.Inventory.Quantity {
		errors = append(errors, ValidationError{
			Field:   "inventory.reserved",
			Value:   product.Inventory.Reserved,
			Tag:     "max",
			Message: "Reserved inventory cannot exceed total quantity",
		})
	}
	return errors
}

func sanitizeString(input string) string {
	// Remove HTML tags
	input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")

	// Remove SQL injection attempts
	dangerous := []string{"'", "\"", ";", "--", "/*", "*/"}
	for _, char := range dangerous {
		input = strings.ReplaceAll(input, char, "")
	}

	return strings.TrimSpace(input)
}

func sanitizeProduct(product *Product) {
	product.Name = strings.TrimSpace(sanitizeString(product.Name))
	product.SKU = strings.TrimSpace(sanitizeString(product.SKU))
	product.Description = strings.TrimSpace(sanitizeString(product.Description))
	product.Currency = strings.TrimSpace(sanitizeString(product.Currency))

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
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON",
			Errors:  []ValidationError{{Message: err.Error()}},
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
	sanitizeProduct(&product)

	productsMutex.Lock()
	product.ID = nextProductID
	nextProductID++
	products = append(products, product)
	productsMutex.Unlock()

	c.JSON(201, APIResponse{
		Success: true,
		Data:    product,
		Message: "Product created successfully",
	})
}

const maxBulkSize = 100

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
	if len(products) > maxBulkSize {
		c.JSON(400, APIResponse{
			Success: false,
			Message: fmt.Sprintf("Bulk size cannot exceed %d items", maxBulkSize),
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
		errors := validateProduct(&product)

		if len(errors) > 0 {
			results = append(results, BulkResult{
				Index:   i,
				Success: false,
				Errors:  errors,
			})
		} else {
			sanitizeProduct(&product)

			productsMutex.Lock()
			product.ID = nextProductID
			nextProductID++
			products = append(products, product)
			productsMutex.Unlock()

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
			Message: "Invalid JSON",
		})
		return
	}

	if !isValidSlug(category.Slug) {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid category slug",
		})
		return
	}
	safeParent := true
	if category.ParentID != nil {
		safeParent = false
		for _, ctg := range categories {
			if *category.ParentID == ctg.ID {
				safeParent = true
				break
			}
		}
	}
	if !safeParent {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "category parent is not valid",
		})
		return
	}
	for _, ctg := range categories {
		if category.Name == ctg.Name {
			c.JSON(400, APIResponse{
				Success: false,
				Message: "category name should be unique",
			})
			return
		}
	}

	categoriesMutex.Lock()
	categories = append(categories, category)
	categoriesMutex.Unlock()
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
			Message: "Invalid SKU",
		})
		return
	}

	for _, product := range products {
		if product.SKU == request.SKU {
			c.JSON(200, APIResponse{
				Success: false,
				Message: "Already exists SKU",
			})
			return
		}
	}
	c.JSON(200, APIResponse{
		Success: true,
		Message: "Valid SKU",
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
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    rules,
		Message: "Validation rules retrieved",
	})
}

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
