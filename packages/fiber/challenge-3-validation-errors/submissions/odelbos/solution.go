package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// Product represents a product in the catalog
type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description string   `json:"description" validate:"required,min=10,max=500"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	Category    string   `json:"category" validate:"required,oneof=electronics clothing books home"`
	SKU         string   `json:"sku" validate:"required,sku"`
	InStock     bool     `json:"in_stock"`
	Tags        []string `json:"tags" validate:"dive,min=2,max=20"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Details []ValidationError `json:"details,omitempty"`
}

type BulkResult struct {
	Success bool              `json:"success"`
	Product *Product          `json:"product,omitempty"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

type BulkResponse struct {
	Success bool         `json:"success"`
	Results []BulkResult `json:"results"`
}

// In-memory storage
var products = []Product{
	{ID: 1, Name: "Laptop", Description: "High-performance laptop for professionals", Price: 999.99, Category: "electronics", SKU: "PROD-12345", InStock: true, Tags: []string{"computer", "work"}},
	{ID: 2, Name: "T-Shirt", Description: "Comfortable cotton t-shirt", Price: 29.99, Category: "clothing", SKU: "PROD-67890", InStock: true, Tags: []string{"clothing", "casual"}},
}
var nextID = 3
var productsMu sync.RWMutex

var validate *validator.Validate

func main() {
	app := fiber.New()

	setupCustomValidator()

	app.Get("/products", getProductsHandler)
	app.Get("/products/:id", getProductHandler)
	app.Post("/products", createProductHandler)
	app.Put("/products/:id", updateProductHandler)
	app.Post("/products/bulk", bulkCreateHandler)

	app.Listen(":3000")
}

func setupCustomValidator() {
	// Use github.com/go-playground/validator/v10 package
	// Register custom validators (SKU format)
	validate = validator.New()
	err := validate.RegisterValidation("sku", func(fl validator.FieldLevel) bool {
		return validateSKU(fl.Field().String())
	})
	if err != nil {
		panic(fmt.Errorf("failed to register validator"))
	}
}

// -------------------------------------------------------------------
// Custom validators
// -------------------------------------------------------------------

func validateSKU(val string) bool {
	// Format: "PROD-" followed by 5 digits : PROD-12345
	re := regexp.MustCompile("^PROD-\\d{5}$")
	return re.MatchString(val)

}

// -------------------------------------------------------------------
// Validation helper functions
// -------------------------------------------------------------------

func validateProduct(product Product) []ValidationError {
	err := validate.Struct(product)
	if err == nil {
		return nil
	}
	var errors []ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, formatValidationError(e))
	}
	return errors
}

func formatValidationError(e validator.FieldError) ValidationError {
	field := e.Field()
	var msg string

	switch e.Tag() {
	case "required":
		msg = fmt.Sprintf("%s is required", field)
	case "min":
		msg = fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
	case "max":
		msg = fmt.Sprintf("%s cannot exceed %s characters", field, e.Param())
	case "gt":
		msg = fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "oneof":
		msg = fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "sku":
		msg = fmt.Sprintf("%s must be in format PROD-XXXXX", field)
	case "dive":
		msg = fmt.Sprintf("All elements in %s must be valid", field)
	default:
		msg = fmt.Sprintf("Validation failed on field '%s' with tag '%s'", field, e.Tag())
	}

	return ValidationError{
		Field:   field,
		Tag:     e.Tag(),
		Value:   e.Value(),
		Message: msg,
	}
}

// -------------------------------------------------------------------
// Route handlers
// -------------------------------------------------------------------

func getProductsHandler(c fiber.Ctx) error {
	filters := map[string]string{
		"category":  c.Query("category"),
		"in_stock":  c.Query("in_stock"),
		"min_price": c.Query("min_price"),
		"max_price": c.Query("max_price"),
		"search":    c.Query("search"),
	}

	productsMu.RLock()
	defer productsMu.RUnlock()

	filteredProducts := filterProducts(products, filters)
	return c.JSON(filteredProducts)
}

func getProductHandler(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid ID",
		})
	}

	product, _ := findProductByID(id)
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Success: false,
			Error:   "Not found",
		})
	}
	return c.JSON(product)
}

func createProductHandler(c fiber.Ctx) error {
	var product Product
	if err := c.Bind().Body(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}
	errors := validateProduct(product)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Validation failed",
			Details: errors,
		})
	}

	// NOTE: Should check for duplicate SKU

	productsMu.Lock()
	defer productsMu.Unlock()

	product.ID = nextID
	nextID++
	products = append(products, product)

	return c.Status(fiber.StatusCreated).JSON(product)
}

func updateProductHandler(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid ID",
		})
	}

	productsMu.Lock()
	defer productsMu.Unlock()

	product, idx := findProductByID(id)
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Success: false,
			Error:   "Not found",
		})
	}

	var updated Product
	if err := c.Bind().Body(&updated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	if updated.SKU == "" {
		updated.SKU = product.SKU
	}
	updated.ID = product.ID

	// NOTE: Should check for duplicate SKU

	errors := validateProduct(updated)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Validation failed",
			Details: errors,
		})
	}

	products[idx] = updated

	return c.JSON(updated)
}

func bulkCreateHandler(c fiber.Ctx) error {
	var bulkProducts []Product
	if err := c.Bind().Body(&bulkProducts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	productsMu.Lock()
	defer productsMu.Unlock()

	results := make([]BulkResult, 0, len(bulkProducts))

	for _, product := range bulkProducts {
		if errors := validateProduct(product); len(errors) > 0 {
			results = append(results, BulkResult{
				Success: false,
				Errors:  errors,
			})
			continue
		}

		if product.SKU == "" {
			product.SKU = generateSKU()
		}

		// NOTE: Should check for duplicate SKU

		product.ID = nextID
		nextID++
		products = append(products, product)

		results = append(results, BulkResult{
			Success: true,
			Product: &product,
		})
	}

	return c.Status(fiber.StatusOK).JSON(BulkResponse{
		Success: true,
		Results: results,
	})
}

// -------------------------------------------------------------------
// Helper functions
// -------------------------------------------------------------------

func findProductByID(id int) (*Product, int) {
	for i, product := range products {
		if product.ID == id {
			return &product, i
		}
	}
	return nil, -1
}

func generateSKU() string {
	return fmt.Sprintf("PROD-%05d", time.Now().UnixNano()%100000)
}

func filterProducts(products []Product, filters map[string]string) []Product {
	var results []Product
	for _, p := range products {
		match := true
		if category, ok := filters["category"]; ok && category != "" && p.Category != category {
			match = false
		}
		if inStock, ok := filters["in_stock"]; ok && inStock != "" {
			stockStatus, err := strconv.ParseBool(inStock)
			if err == nil && p.InStock != stockStatus {
				match = false
			}
		}
		if minPriceStr, ok := filters["min_price"]; ok && minPriceStr != "" {
			minPrice, err := strconv.ParseFloat(minPriceStr, 64)
			if err == nil && p.Price < minPrice {
				match = false
			}
		}
		if maxPriceStr, ok := filters["max_price"]; ok && maxPriceStr != "" {
			maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
			if err == nil && p.Price > maxPrice {
				match = false
			}
		}
		if searchTerm, ok := filters["search"]; ok && searchTerm != "" {
			searchTerm = strings.ToLower(searchTerm)
			if !strings.Contains(strings.ToLower(p.Name), searchTerm) && !strings.Contains(strings.ToLower(p.Description), searchTerm) {
				match = false
			}
		}
		if match {
			results = append(results, p)
		}
	}
	return results
}
