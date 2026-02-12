package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// Product represents a product in the catalog
type Product struct {
	ID          int                    `json:"id"`
	SKU         string                 `json:"sku" binding:"required"`
	Name        string                 `json:"name" binding:"required,min=3,max=100"`
	Description string                 `json:"description" binding:"max=1000"`
	Price       float64                `json:"price" binding:"required,min=0.01"` // use integer cents or a decimal library in production
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
	Width     int    `json:"width" binding:"required,min=100"`
	Height    int    `json:"height" binding:"required,min=100"`
	Size      int64  `json:"size"`
	IsPrimary bool   `json:"is_primary"`
}

// Inventory represents product inventory information
type Inventory struct {
	Quantity    int       `json:"quantity" binding:"min=0"` //Removed required to distinguish "not provided" from "zero"
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

var (
	products      = []Product{}
	productsMutex sync.RWMutex
	nextProductID = 1
)

var (
	categories = []Category{
		{ID: 1, Name: "Electronics", Slug: "electronics"},
		{ID: 2, Name: "Clothing", Slug: "clothing"},
		{ID: 3, Name: "Books", Slug: "books"},
		{ID: 4, Name: "Home & Garden", Slug: "home-garden"},
	}
	// never request for productsMutex while holding categoriesMutex
	// the function signatures in the assignment lead us to get productsMutex and then categoriesMutex,
	// for production avoid holding two locks
	categoriesMutex sync.RWMutex
)
var validCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"GBP": true,
	"JPY": true,
}

var validWarehouses = []string{"WH001", "WH002", "WH003", "WH004", "WH005"}

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

// function signature is part of the assignment
func isValidCategory(cat string) bool {
	categoriesMutex.RLock()
	defer categoriesMutex.RUnlock()
	for _, c := range categories {
		if c.Name == cat {
			return true
		}
	}
	return false
}

// function signature is part of the assignment
// validateProduct validates a product's fields. Caller must hold productsMutex (at least RLock).
func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError
	if utf8.RuneCountInString(product.Name) < 3 || utf8.RuneCountInString(product.Name) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be between 3 and 100 characters after sanitization",
		})
	}
	if utf8.RuneCountInString(product.Description) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must not exceed 1000 characters after sanitization",
		})
	}

	if !isValidSKU(product.SKU) {
		errors = append(errors, ValidationError{
			Field:   "sku",
			Message: "SKU must follow ABC-123-XYZ format",
		})
	}

	for _, prod := range products {
		if product.SKU == prod.SKU {
			errors = append(errors, ValidationError{
				Field:   "sku",
				Message: "SKU already exists",
			})
			break
		}
	}

	if !isValidCurrency(product.Currency) {
		errors = append(errors, ValidationError{
			Field:   "currency",
			Message: "Must be a valid ISO 4217 currency code",
		})
	}
	if !isValidCategory(product.Category.Name) {
		errors = append(errors, ValidationError{
			Field:   "category",
			Message: "Category ID, name, and slug must match an existing category",
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
	if !isValidWarehouseCode(product.Inventory.Location) {
		errors = append(errors, ValidationError{
			Field:   "inventory.location",
			Value:   product.Inventory.Location,
			Tag:     "warehouse",
			Message: "Must be a valid warehouse code (e.g., WH001)",
		})
	}
	return errors
}

func sanitizeString(input string) string {
	return strings.TrimSpace(input)
}

func sanitizeProduct(product *Product) {
	product.Name = sanitizeString(product.Name)
	product.SKU = sanitizeString(product.SKU)
	product.Description = sanitizeString(product.Description)
	product.Currency = sanitizeString(product.Currency)

	product.Currency = strings.ToUpper(product.Currency)
	product.Category.Slug = strings.ToLower(product.Category.Slug)
	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved // required to be calculated in this function for this assignment
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
	sanitizeProduct(&product)
	productsMutex.Lock()
	defer productsMutex.Unlock()
	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}
	product.ID = nextProductID
	nextProductID++

	products = append(products, product)

	c.JSON(201, APIResponse{
		Success: true,
		Data:    product,
		Message: "Product created successfully",
	})
}

const maxBulkSize = 100

// POST /products/bulk - Create multiple products, non-atomic — partial inserts on mixed success/failure
func createProductsBulk(c *gin.Context) {
	// for production to lower lock duration, batch-validate first, then lock only for the insert phase
	var inputProducts []Product

	if err := c.ShouldBindJSON(&inputProducts); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}
	if len(inputProducts) > maxBulkSize {
		c.JSON(400, APIResponse{
			Success: false,
			Message: fmt.Sprintf("Bulk size cannot exceed %d items", maxBulkSize),
		})
		return
	}
	if len(inputProducts) == 0 {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "At least one product is required",
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
	productsMutex.Lock()
	defer productsMutex.Unlock()
	for i, prod := range inputProducts {
		product := prod // avoid pointing to the last item by pointer
		sanitizeProduct(&product)
		errors := validateProduct(&product)

		if len(errors) > 0 {
			results = append(results, BulkResult{
				Index:   i,
				Success: false,
				Errors:  errors,
			})
		} else {
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
			Message: "Invalid JSON",
		})
		return
	}
	category.Name = sanitizeString(category.Name)
	category.Slug = sanitizeString(category.Slug)
	if !isValidSlug(category.Slug) {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid category slug",
		})
		return
	}
	categoriesMutex.Lock()
	defer categoriesMutex.Unlock()
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
		if category.ID == ctg.ID {
			c.JSON(400, APIResponse{
				Success: false,
				Message: "category id should be unique",
			})
			return
		}
		if category.Slug == ctg.Slug {
			c.JSON(400, APIResponse{
				Success: false,
				Message: "category slug should be unique",
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
	if !isValidSKU(request.SKU) {
		c.JSON(200, APIResponse{
			Success: false,
			Message: "Invalid SKU",
		})
		return
	}
	productsMutex.RLock()
	defer productsMutex.RUnlock()
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
	sanitizeProduct(&product)
	productsMutex.RLock()
	defer productsMutex.RUnlock()
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
	rules := map[string]any{
		"sku": map[string]interface{}{
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
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
