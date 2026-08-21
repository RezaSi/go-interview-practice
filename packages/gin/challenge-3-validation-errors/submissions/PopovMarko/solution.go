package main

import (
	"fmt"
	"log"
	"net/http"
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
var (
	rgSKUFormat, _  = regexp.Compile(`^[A-Z]{3}-\d{3}-[A-Z]{3}$`)
	rgSlugFormat, _ = regexp.Compile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

// isValidSKU reports whether sku matches the required SKU format
// "ABC-123-XYZ" (three uppercase letters, three digits, three uppercase letters).
func isValidSKU(sku string) bool {
	return rgSKUFormat.MatchString(sku)
}

// isValidCurrency reports whether currency is one of the supported ISO
// currency codes listed in validCurrencies.
func isValidCurrency(currency string) bool {
	return slices.Contains(validCurrencies, currency)
}

// isValidCategory reports whether categoryName matches the Name of a
// known category in the categories store.
func isValidCategory(categoryName string) bool {
	for _, c := range categories {
		if categoryName == c.Name {
			return true
		}
	}
	return false
}

// isValidSlug reports whether slug is a valid URL slug: lowercase
// alphanumeric segments separated by single hyphens, with no leading
// or trailing hyphen.
func isValidSlug(slug string) bool {
	return rgSlugFormat.MatchString(slug)
}

// isValidWarehouseCode reports whether code is one of the known
// warehouse codes listed in validWarehouses.
func isValidWarehouseCode(code string) bool {
	return slices.Contains(validWarehouses, code)
}

// validateProduct runs the custom business-rule validations on product
// that cannot be expressed with gin binding tags: SKU format, currency
// code, category existence, category slug format, warehouse code, and
// the cross-field rule that reserved inventory cannot exceed quantity.
// It returns one ValidationError per failed rule, or an empty slice if
// the product is valid.
func validateProduct(product *Product) []ValidationError {
	var errors []ValidationError

	// Validate SKU format.
	if !isValidSKU(product.SKU) {
		err := ValidationError{
			Field:   "product.sku",
			Value:   product.SKU,
			Tag:     "SKU",
			Message: "Wrong SKU format",
		}
		errors = append(errors, err)
	}
	// Validate currency code.
	if !isValidCurrency(product.Currency) {
		err := ValidationError{
			Field:   "product.currency",
			Value:   product.Currency,
			Tag:     "oneof",
			Message: "Unsupported currency",
		}
		errors = append(errors, err)
	}
	// Validate that the category exists.
	if !isValidCategory(product.Category.Name) {
		err := ValidationError{
			Field:   "product.category.name",
			Value:   product.Category.Name,
			Tag:     "required",
			Message: "No such category",
		}
		errors = append(errors, err)
	}
	// Validate the category slug format.
	if !isValidSlug(product.Category.Slug) {
		err := ValidationError{
			Field:   "product.category.slug",
			Value:   product.Category.Slug,
			Tag:     "format",
			Message: "Wrong slug format",
		}
		errors = append(errors, err)
	}
	// Validate the warehouse location code.
	if !isValidWarehouseCode(product.Inventory.Location) {
		err := ValidationError{
			Field:   "product.inventory.location",
			Value:   product.Inventory.Location,
			Tag:     "oneof",
			Message: "Wrong warehouse code",
		}
		errors = append(errors, err)
	}
	// Cross-field validation: reserved must not exceed quantity.
	if product.Inventory.Reserved > product.Inventory.Quantity {
		err := ValidationError{
			Field: "product.inventory.reserved, product.inventory.quantity",
			Value: fmt.Sprintf(
				"Quantity = %d, Reserved = %d",
				product.Inventory.Quantity,
				product.Inventory.Reserved,
			),
			Tag:     "gt",
			Message: "Reserved products can't be more than Quantity",
		}
		errors = append(errors, err)
	}

	return errors
}

// sanitizeProduct normalizes user-supplied product fields in place:
// it trims surrounding whitespace, upper-cases the currency, lower-cases
// the category slug, computes the available inventory (quantity minus
// reserved), and stamps the created/updated/last-updated timestamps.
func sanitizeProduct(product *Product) {
	now := time.Now()
	product.SKU = strings.TrimSpace(product.SKU)
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)
	product.Currency = strings.ToUpper(strings.TrimSpace(product.Currency))
	product.Category.Slug = strings.ToLower(strings.TrimSpace(product.Category.Slug))
	product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
	product.Inventory.LastUpdated = now
	product.CreatedAt = now
	product.UpdatedAt = now
}

// createProduct handles POST /products. It binds and validates the
// request body, runs the custom business-rule validations, sanitizes the
// input, assigns a new ID, stores the product, and returns it with a 201
// status. Binding or validation failures return a 400 with the errors.
func createProduct(c *gin.Context) {
	var product Product
	// Bind JSON and surface any binding/basic validation errors.
	if err := c.ShouldBindJSON(&product); err != nil {
		errs := []ValidationError{
			{
				Message: err.Error(),
			},
		}
		errorResponse(c, "Invalid JSON or basic validation failed", errs, "", http.StatusBadRequest)
		return
	}

	// Apply the custom business-rule validations.
	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		errorResponse(c, "Validation failed", validationErrors, "", http.StatusBadRequest)
		return
	}

	// Sanitize the input data.
	sanitizeProduct(&product)

	// Assign an ID and store the product.
	product.ID = nextProductID
	nextProductID++
	products = append(products, product)

	successResponse(c, product, "Product created successfully", http.StatusCreated)
}

// createProductsBulk handles POST /products/bulk. It accepts an array of
// products and validates each one independently, allowing partial success:
// valid products are sanitized, stored, and reported as succeeded, while
// invalid ones are reported with their validation errors. The response
// includes per-item results and total/successful/failed counts.
func createProductsBulk(c *gin.Context) {
	var (
		inputProducts    []Product
		validationErrors []ValidationError
	)

	if err := c.ShouldBindJSON(&inputProducts); err != nil {
		validationErrors = append(validationErrors, ValidationError{Message: err.Error()})
		errorResponse(c, "Invalid JSON format", validationErrors, "", http.StatusBadRequest)
		return
	}

	// BulkResult captures the outcome of validating a single product
	// within the bulk request, keyed by its index in the input array.
	type BulkResult struct {
		Index   int               `json:"index"`
		Success bool              `json:"success"`
		Product *Product          `json:"product,omitempty"`
		Errors  []ValidationError `json:"errors,omitempty"`
	}

	var results []BulkResult
	var successCount int

	// Process each product and populate the per-item results.
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

			p := product
			results = append(results, BulkResult{
				Index:   i,
				Success: true,
				Product: &p,
			})
			successCount++
		}
	}

	// Success is only true when every product in the batch was created;
	// a partial failure reports success=false while still returning the
	// per-item results.
	c.JSON(http.StatusOK, APIResponse{
		Success: successCount == len(inputProducts),
		Data: map[string]interface{}{
			"results":    results,
			"total":      len(inputProducts),
			"successful": successCount,
			"failed":     len(inputProducts) - successCount,
		},
		Message:   "Bulk operation completed",
		RequestID: c.GetHeader("X-Request-ID"),
	})
}

// createCategory handles POST /categories. It binds the request body,
// appends the new category to the in-memory store, and returns the full
// category list with a 201 status.
func createCategory(c *gin.Context) {
	var category Category

	if err := c.ShouldBindJSON(&category); err != nil {
		errorResponse(c, "Invalid JSON or validation failed", nil, "", http.StatusBadRequest)
		return
	}

	categories = append(categories, category)
	successResponse(c, categories, "Category created successfully", http.StatusCreated)
}

// validateSKUEndpoint handles POST /validate/sku. It checks that the
// supplied SKU is well-formed and not already used by an existing product.
// Format or uniqueness failures return success=false with a 200 status;
// only a missing SKU in the request body yields a 400.
func validateSKUEndpoint(c *gin.Context) {
	var request struct {
		SKU string `json:"sku" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, "SKU is required", nil, "", http.StatusBadRequest)
		return
	}

	if valid := isValidSKU(request.SKU); !valid {
		errorResponse(c, "Invalid SKU", nil, "", http.StatusOK)
		return
	}
	for _, product := range products {
		if request.SKU == product.SKU {
			errorResponse(c, "SKU already exists", nil, "", http.StatusOK)
			return
		}
	}
	successResponse(c, nil, "SKU is valid", http.StatusOK)
}

// validateProductEndpoint handles POST /validate/product. It validates a
// product against all rules without persisting it, returning success when
// the product is valid or the list of validation errors otherwise.
func validateProductEndpoint(c *gin.Context) {
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		errorResponse(c, "Invalid JSON format", nil, "", http.StatusBadRequest)
		return
	}

	validationErrors := validateProduct(&product)
	if len(validationErrors) > 0 {
		errorResponse(c, "Validation failed", validationErrors, "", http.StatusBadRequest)
		return
	}

	successResponse(c, nil, "Product data is valid", http.StatusOK)
}

// getValidationRules handles GET /validation/rules. It returns a
// machine-readable description of the validation constraints (SKU, name,
// currency, and warehouse rules) that clients can use to pre-validate input.
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

	successResponse(c, rules, "Validation rules retrieved", http.StatusOK)
}

// successResponse writes a standard success APIResponse with the given
// data, message, and HTTP status, echoing the X-Request-ID header.
func successResponse(c *gin.Context, data interface{}, msg string, status int) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Message:   msg,
		RequestID: c.GetHeader("X-Request-ID"),
	}
	c.JSON(status, response)
}

// errorResponse writes a standard failure APIResponse with the given
// message, validation errors, error code, and HTTP status, echoing the
// X-Request-ID header.
func errorResponse(c *gin.Context, msg string, errs []ValidationError, code string, status int) {
	response := APIResponse{
		Success:   false,
		Message:   msg,
		Errors:    errs,
		ErrorCode: code,
		RequestID: c.GetHeader("X-Request-ID"),
	}
	c.JSON(status, response)
}

// setupRouter builds and returns the gin engine with all product,
// category, and validation routes registered.
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

// main starts the HTTP server on port 8080, exiting fatally if the
// server fails to start.
func main() {
	router := setupRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server")
	}
}
