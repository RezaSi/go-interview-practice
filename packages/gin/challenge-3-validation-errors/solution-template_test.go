package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestSKUValidation(t *testing.T) {
	tests := []struct {
		name     string
		sku      string
		expected bool
	}{
		{"Valid SKU", "ABC-123-XYZ", true},
		{"Valid SKU 2", "DEF-456-GHI", true},
		{"Invalid - lowercase", "abc-123-xyz", false},
		{"Invalid - wrong format", "ABC123XYZ", false},
		{"Invalid - missing parts", "ABC-123", false},
		{"Invalid - too many letters", "ABCD-123-XYZ", false},
		{"Invalid - too many numbers", "ABC-1234-XYZ", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidSKU(tt.sku)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCurrencyValidation(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		expected bool
	}{
		{"Valid USD", "USD", true},
		{"Valid EUR", "EUR", true},
		{"Valid GBP", "GBP", true},
		{"Invalid currency", "XYZ", false},
		{"Lowercase", "usd", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCurrency(tt.currency)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCategoryValidation(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected bool
	}{
		{"Valid Electronics", "Electronics", true},
		{"Valid Clothing", "Clothing", true},
		{"Invalid category", "InvalidCategory", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCategory(tt.category)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlugValidation(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		expected bool
	}{
		{"Valid slug", "electronics", true},
		{"Valid slug with hyphens", "home-garden", true},
		{"Valid alphanumeric", "abc123", true},
		{"Invalid - uppercase", "Electronics", false},
		{"Invalid - spaces", "home garden", false},
		{"Invalid - special chars", "home_garden", false},
		{"Invalid - starting with hyphen", "-electronics", false},
		{"Invalid - ending with hyphen", "electronics-", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidSlug(tt.slug)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWarehouseCodeValidation(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"Valid WH001", "WH001", true},
		{"Valid WH002", "WH002", true},
		{"Invalid - not in list", "WH999", false},
		{"Invalid - wrong format", "W001", false},
		{"Invalid - lowercase", "wh001", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidWarehouseCode(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateProductSuccess(t *testing.T) {
	router := setupRouter()

	product := map[string]interface{}{
		"sku":      "ABC-123-XYZ",
		"name":     "Test Product",
		"price":    29.99,
		"currency": "USD",
		"category": map[string]interface{}{
			"id":   1,
			"name": "Electronics",
			"slug": "electronics",
		},
		"inventory": map[string]interface{}{
			"quantity": 100,
			"reserved": 10,
			"location": "WH001",
		},
	}

	jsonData, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Product created successfully", response.Message)
}

func TestCreateProductValidationErrors(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		product        map[string]interface{}
		expectedStatus int
		expectedErrors int
	}{
		{
			name: "Missing required fields",
			product: map[string]interface{}{
				"name": "Test Product",
			},
			expectedStatus: 400,
			expectedErrors: 1, // Should have validation errors
		},
		{
			name: "Invalid SKU format",
			product: map[string]interface{}{
				"sku":      "invalid-sku",
				"name":     "Test Product",
				"price":    29.99,
				"currency": "USD",
				"category": map[string]interface{}{
					"id":   1,
					"name": "Electronics",
					"slug": "electronics",
				},
				"inventory": map[string]interface{}{
					"quantity": 100,
					"reserved": 10,
					"location": "WH001",
				},
			},
			expectedStatus: 400,
			expectedErrors: 1,
		},
		{
			name: "Invalid currency",
			product: map[string]interface{}{
				"sku":      "ABC-123-XYZ",
				"name":     "Test Product",
				"price":    29.99,
				"currency": "INVALID",
				"category": map[string]interface{}{
					"id":   1,
					"name": "Electronics",
					"slug": "electronics",
				},
				"inventory": map[string]interface{}{
					"quantity": 100,
					"reserved": 10,
					"location": "WH001",
				},
			},
			expectedStatus: 400,
			expectedErrors: 1,
		},
		{
			name: "Invalid category",
			product: map[string]interface{}{
				"sku":      "ABC-123-XYZ",
				"name":     "Test Product",
				"price":    29.99,
				"currency": "USD",
				"category": map[string]interface{}{
					"id":   1,
					"name": "InvalidCategory",
					"slug": "invalid-category",
				},
				"inventory": map[string]interface{}{
					"quantity": 100,
					"reserved": 10,
					"location": "WH001",
				},
			},
			expectedStatus: 400,
			expectedErrors: 1,
		},
		{
			name: "Reserved exceeds quantity",
			product: map[string]interface{}{
				"sku":      "ABC-123-XYZ",
				"name":     "Test Product",
				"price":    29.99,
				"currency": "USD",
				"category": map[string]interface{}{
					"id":   1,
					"name": "Electronics",
					"slug": "electronics",
				},
				"inventory": map[string]interface{}{
					"quantity": 10,
					"reserved": 20, // More than quantity
					"location": "WH001",
				},
			},
			expectedStatus: 400,
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.product)
			req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.False(t, response.Success)

			if tt.expectedErrors > 0 {
				assert.GreaterOrEqual(t, len(response.Errors), tt.expectedErrors)
			}
		})
	}
}

func TestBulkProductCreation(t *testing.T) {
	router := setupRouter()

	products := []map[string]interface{}{
		{
			"sku":      "ABC-123-XYZ",
			"name":     "Valid Product 1",
			"price":    29.99,
			"currency": "USD",
			"category": map[string]interface{}{
				"id":   1,
				"name": "Electronics",
				"slug": "electronics",
			},
			"inventory": map[string]interface{}{
				"quantity": 100,
				"reserved": 10,
				"location": "WH001",
			},
		},
		{
			"sku":      "invalid-sku", // Invalid SKU
			"name":     "Invalid Product",
			"price":    19.99,
			"currency": "USD",
			"category": map[string]interface{}{
				"id":   1,
				"name": "Electronics",
				"slug": "electronics",
			},
			"inventory": map[string]interface{}{
				"quantity": 50,
				"reserved": 5,
				"location": "WH001",
			},
		},
		{
			"sku":      "DEF-456-GHI",
			"name":     "Valid Product 2",
			"price":    39.99,
			"currency": "EUR",
			"category": map[string]interface{}{
				"id":   2,
				"name": "Clothing",
				"slug": "clothing",
			},
			"inventory": map[string]interface{}{
				"quantity": 75,
				"reserved": 0,
				"location": "WH002",
			},
		},
	}

	jsonData, _ := json.Marshal(products)
	req, _ := http.NewRequest("POST", "/products/bulk", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should be partial success (2 out of 3 products valid)
	assert.False(t, response.Success)

	data := response.Data.(map[string]interface{})
	assert.Equal(t, float64(3), data["total"])
	assert.Equal(t, float64(2), data["successful"])
	assert.Equal(t, float64(1), data["failed"])

	results := data["results"].([]interface{})
	assert.Len(t, results, 3)

	// First product should be successful
	result1 := results[0].(map[string]interface{})
	assert.True(t, result1["success"].(bool))

	// Second product should fail (invalid SKU)
	result2 := results[1].(map[string]interface{})
	assert.False(t, result2["success"].(bool))
	assert.NotEmpty(t, result2["errors"])

	// Third product should be successful
	result3 := results[2].(map[string]interface{})
	assert.True(t, result3["success"].(bool))
}

func TestCreateCategory(t *testing.T) {
	router := setupRouter()

	category := map[string]interface{}{
		"id":   5,
		"name": "Sports",
		"slug": "sports",
	}

	jsonData, _ := json.Marshal(category)
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Category created successfully", response.Message)
}

func TestValidateSKUEndpoint(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		sku            string
		expectedStatus int
		expectedValid  bool
	}{
		{"Valid SKU", "ABC-123-XYZ", 200, true},
		{"Invalid SKU", "invalid-sku", 200, false}, // Endpoint returns 200 but validation fails
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestData := map[string]string{"sku": tt.sku}
			jsonData, _ := json.Marshal(requestData)
			req, _ := http.NewRequest("POST", "/validate/sku", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedValid, response.Success)
		})
	}
}

func TestValidateProductEndpoint(t *testing.T) {
	router := setupRouter()

	validProduct := map[string]interface{}{
		"sku":      "ABC-123-XYZ",
		"name":     "Test Product",
		"price":    29.99,
		"currency": "USD",
		"category": map[string]interface{}{
			"id":   1,
			"name": "Electronics",
			"slug": "electronics",
		},
		"inventory": map[string]interface{}{
			"quantity": 100,
			"reserved": 10,
			"location": "WH001",
		},
	}

	jsonData, _ := json.Marshal(validProduct)
	req, _ := http.NewRequest("POST", "/validate/product", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Product data is valid", response.Message)
}

func TestGetValidationRules(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/validation/rules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	rules := response.Data.(map[string]interface{})
	assert.Contains(t, rules, "sku")
	assert.Contains(t, rules, "name")
	assert.Contains(t, rules, "currency")
	assert.Contains(t, rules, "warehouse")
}

func TestProductValidation(t *testing.T) {
	// Test the validateProduct function directly
	tests := []struct {
		name          string
		product       Product
		expectedError bool
	}{
		{
			name: "Valid product",
			product: Product{
				SKU:      "ABC-123-XYZ",
				Name:     "Test Product",
				Price:    29.99,
				Currency: "USD",
				Category: Category{
					ID:   1,
					Name: "Electronics",
					Slug: "electronics",
				},
				Inventory: Inventory{
					Quantity: 100,
					Reserved: 10,
					Location: "WH001",
				},
			},
			expectedError: false,
		},
		{
			name: "Invalid SKU",
			product: Product{
				SKU:      "invalid-sku",
				Name:     "Test Product",
				Price:    29.99,
				Currency: "USD",
				Category: Category{
					ID:   1,
					Name: "Electronics",
					Slug: "electronics",
				},
				Inventory: Inventory{
					Quantity: 100,
					Reserved: 10,
					Location: "WH001",
				},
			},
			expectedError: true,
		},
		{
			name: "Reserved exceeds quantity",
			product: Product{
				SKU:      "ABC-123-XYZ",
				Name:     "Test Product",
				Price:    29.99,
				Currency: "USD",
				Category: Category{
					ID:   1,
					Name: "Electronics",
					Slug: "electronics",
				},
				Inventory: Inventory{
					Quantity: 10,
					Reserved: 20, // More than quantity
					Location: "WH001",
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validateProduct(&tt.product)
			if tt.expectedError {
				assert.NotEmpty(t, errors)
			} else {
				assert.Empty(t, errors)
			}
		})
	}
}

func TestInputSanitization(t *testing.T) {
	product := &Product{
		SKU:      " ABC-123-XYZ ",    // Extra spaces
		Name:     "  Test Product  ", // Extra spaces
		Currency: "usd",              // Lowercase should become uppercase
		Category: Category{
			Slug: "ELECTRONICS", // Uppercase should become lowercase
		},
		Inventory: Inventory{
			Quantity: 100,
			Reserved: 10,
			// Available should be calculated as 90
		},
	}

	sanitizeProduct(product)

	// Check that sanitization worked
	assert.Equal(t, "ABC-123-XYZ", product.SKU)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, "USD", product.Currency)
	assert.Equal(t, "electronics", product.Category.Slug)
	assert.Equal(t, 90, product.Inventory.Available) // quantity - reserved
	assert.False(t, product.CreatedAt.IsZero())
	assert.False(t, product.UpdatedAt.IsZero())
}

// Helper function to reset test data
func resetTestData() {
	products = []Product{}
	nextProductID = 1
}
