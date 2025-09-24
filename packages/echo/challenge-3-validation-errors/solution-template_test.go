package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductAPI(t *testing.T) {
	// Create a new Echo instance for testing
	e := echo.New()
	setupRoutes(e)

	t.Run("GET /health - should return health status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response HealthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "Product API is running", response.Message)
		assert.NotEmpty(t, response.Timestamp)
		assert.Equal(t, "1.0.0", response.Version)
	})

	t.Run("POST /products - should create a new product with valid data", func(t *testing.T) {
		productJSON := `{"name":"Test Product","description":"Test description","sku":"ABC-12345","price":29.99,"category":"electronics","in_stock":true,"stock_quantity":100,"tags":["test","product"]}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response ProductResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Product created successfully", response.Message)
		assert.NotNil(t, response.Product)
		assert.Equal(t, "Test Product", response.Product.Name)
		assert.Equal(t, "ABC-12345", response.Product.SKU)
		assert.Equal(t, 29.99, response.Product.Price)
		assert.Equal(t, "electronics", response.Product.Category)
		assert.NotEmpty(t, response.Product.ID)
		assert.NotEmpty(t, response.Product.CreatedAt)
	})

	t.Run("POST /products - should validate required fields", func(t *testing.T) {
		productJSON := `{"description":"Missing required fields"}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response ValidationErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Validation failed", response.Message)
		assert.NotEmpty(t, response.Errors)
		assert.Contains(t, response.Errors, "Name")
		assert.Contains(t, response.Errors, "SKU")
		assert.Contains(t, response.Errors, "Price")
		assert.Contains(t, response.Errors, "Category")
	})

	t.Run("POST /products - should validate SKU format", func(t *testing.T) {
		productJSON := `{"name":"Test Product","sku":"invalid-sku","price":29.99,"category":"electronics","stock_quantity":100}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response ValidationErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Errors, "SKU")
		assert.Contains(t, response.Errors["SKU"], "XXX-NNNNN")
	})

	t.Run("POST /products - should validate category", func(t *testing.T) {
		productJSON := `{"name":"Test Product","sku":"ABC-12345","price":29.99,"category":"invalid-category","stock_quantity":100}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response ValidationErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Errors, "Category")
	})

	t.Run("GET /products - should return all products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProductListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.GreaterOrEqual(t, response.Count, 0)
		assert.NotNil(t, response.Products)
	})

	t.Run("GET /products with filters", func(t *testing.T) {
		// Reset products for clean test
		resetProducts()

		// First create some test products
		products := []string{
			`{"name":"Electronics Product","sku":"ELE-12345","price":99.99,"category":"electronics","in_stock":true,"stock_quantity":50}`,
			`{"name":"Clothing Product","sku":"CLO-12345","price":49.99,"category":"clothing","in_stock":false,"stock_quantity":0}`,
			`{"name":"Book Product","sku":"BOO-12345","price":19.99,"category":"books","in_stock":true,"stock_quantity":25}`,
		}

		for _, productJSON := range products {
			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		// Test category filter
		req := httptest.NewRequest(http.MethodGet, "/products?category=electronics", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProductListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		for _, product := range response.Products {
			assert.Equal(t, "electronics", product.Category)
		}

		// Test price filter
		req = httptest.NewRequest(http.MethodGet, "/products?min_price=50", nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		for _, product := range response.Products {
			assert.GreaterOrEqual(t, product.Price, 50.0)
		}
	})

	t.Run("GET /products/:id - should return specific product", func(t *testing.T) {
		// First create a product
		productJSON := `{"name":"Specific Product","sku":"SPC-12345","price":39.99,"category":"electronics","stock_quantity":75}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse ProductResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		productID := createResponse.Product.ID

		// Now get the specific product
		req = httptest.NewRequest(http.MethodGet, "/products/"+productID, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProductResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, productID, response.Product.ID)
		assert.Equal(t, "Specific Product", response.Product.Name)
	})

	t.Run("PUT /products/:id - should update existing product", func(t *testing.T) {
		// First create a product
		productJSON := `{"name":"Original Product","sku":"ORI-12345","price":29.99,"category":"electronics","stock_quantity":50}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse ProductResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		productID := createResponse.Product.ID

		// Now update the product
		updateJSON := `{"name":"Updated Product","sku":"UPD-12345","price":39.99,"category":"clothing","stock_quantity":75}`
		req = httptest.NewRequest(http.MethodPut, "/products/"+productID, strings.NewReader(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProductResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Product updated successfully", response.Message)
		assert.Equal(t, "Updated Product", response.Product.Name)
		assert.Equal(t, "UPD-12345", response.Product.SKU)
		assert.Equal(t, 39.99, response.Product.Price)
	})

	t.Run("DELETE /products/:id - should delete existing product", func(t *testing.T) {
		// First create a product
		productJSON := `{"name":"Product to Delete","sku":"DEL-12345","price":19.99,"category":"books","stock_quantity":25}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var createResponse ProductResponse
		err := json.Unmarshal(rec.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		productID := createResponse.Product.ID

		// Now delete the product
		req = httptest.NewRequest(http.MethodDelete, "/products/"+productID, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Product deleted successfully", response.Message)

		// Verify product is deleted
		req = httptest.NewRequest(http.MethodGet, "/products/"+productID, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("POST /products/bulk - should handle bulk creation", func(t *testing.T) {
		bulkJSON := `{
			"products": [
				{"name":"Bulk Product 1","sku":"BLK-11111","price":10.99,"category":"electronics","stock_quantity":10},
				{"name":"Bulk Product 2","sku":"BLK-22222","price":20.99,"category":"clothing","stock_quantity":20},
				{"name":"","sku":"BLK-33333","price":30.99,"category":"books","stock_quantity":30}
			]
		}`
		req := httptest.NewRequest(http.MethodPost, "/products/bulk", strings.NewReader(bulkJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response BulkCreateResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "partial_success", response.Status)
		assert.Len(t, response.Succeeded, 2)
		assert.Len(t, response.Failed, 1)
		assert.Equal(t, 2, response.Failed[0].Index)
	})
}

func TestProductStructure(t *testing.T) {
	t.Run("Product struct should have correct fields", func(t *testing.T) {
		product := Product{
			ID:            "test-id",
			Name:          "Test Product",
			Description:   "Test Description",
			SKU:           "ABC-12345",
			Price:         29.99,
			Category:      "electronics",
			InStock:       true,
			StockQuantity: 100,
			Tags:          []string{"test", "product"},
			CreatedAt:     "2024-01-01T00:00:00Z",
			UpdatedAt:     "2024-01-01T00:00:00Z",
		}

		assert.Equal(t, "test-id", product.ID)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, "ABC-12345", product.SKU)
		assert.Equal(t, 29.99, product.Price)
		assert.Equal(t, "electronics", product.Category)
		assert.True(t, product.InStock)
		assert.Equal(t, 100, product.StockQuantity)
	})

	t.Run("Response structs should have correct structure", func(t *testing.T) {
		product := Product{ID: "1", Name: "Test"}

		productResponse := ProductResponse{
			Status:  "success",
			Message: "Product created",
			Product: &product,
		}

		assert.Equal(t, "success", productResponse.Status)
		assert.Equal(t, "Product created", productResponse.Message)
		assert.NotNil(t, productResponse.Product)

		listResponse := ProductListResponse{
			Status:   "success",
			Count:    1,
			Products: []Product{product},
		}

		assert.Equal(t, "success", listResponse.Status)
		assert.Equal(t, 1, listResponse.Count)
		assert.Len(t, listResponse.Products, 1)
	})
}
