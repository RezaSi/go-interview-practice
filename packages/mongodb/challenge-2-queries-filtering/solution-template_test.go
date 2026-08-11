package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetProductsByCategory(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful category filtering", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.products", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"name", "iPhone 14"},
			{"category", "Electronics"},
			{"price", 999.99},
			{"rating", 4.8},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.products", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetProductsByCategory(context.Background(), "Electronics")

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Products retrieved successfully", response.Message)
	})

	mt.Run("empty category validation", func(mt *mtest.T) {
		productService := &ProductService{Collection: mt.Coll}

		// Test empty string
		response := productService.GetProductsByCategory(context.Background(), "")
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Category cannot be empty")

		// Test whitespace-only string
		response = productService.GetProductsByCategory(context.Background(), "   ")
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Category cannot be empty")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "database connection failed",
		}))

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetProductsByCategory(context.Background(), "Electronics")

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to retrieve products")
	})

	mt.Run("nil collection handling", func(mt *mtest.T) {
		productService := &ProductService{Collection: nil}
		response := productService.GetProductsByCategory(context.Background(), "Electronics")

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Collection not initialized")
	})
}

func TestGetProductsByPriceRange(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful price range filtering", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.products", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"name", "Mid-range Phone"},
			{"price", 599.99},
			{"rating", 4.5},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.products", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetProductsByPriceRange(context.Background(), 500.0, 700.0)

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Products retrieved successfully", response.Message)
	})

	mt.Run("price range validation", func(mt *mtest.T) {
		productService := &ProductService{Collection: mt.Coll}

		// Test both prices zero
		response := productService.GetProductsByPriceRange(context.Background(), 0.0, 0.0)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Price range cannot be zero for both min and max")

		// Test negative minPrice
		response = productService.GetProductsByPriceRange(context.Background(), -10.0, 100.0)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Price cannot be negative")

		// Test negative maxPrice
		response = productService.GetProductsByPriceRange(context.Background(), 10.0, -100.0)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Price cannot be negative")

		// Test minPrice > maxPrice
		response = productService.GetProductsByPriceRange(context.Background(), 100.0, 50.0)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Minimum price cannot be greater than maximum price")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "database connection failed",
		}))

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetProductsByPriceRange(context.Background(), 10.0, 100.0)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to retrieve products")
	})
}

func TestSearchProductsByName(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful name search", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.products", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"name", "iPhone 14 Pro"},
			{"category", "Electronics"},
			{"price", 1099.99},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.products", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		productService := &ProductService{Collection: mt.Coll}
		response := productService.SearchProductsByName(context.Background(), "iPhone")

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Products found successfully", response.Message)
	})

	mt.Run("search term validation", func(mt *mtest.T) {
		productService := &ProductService{Collection: mt.Coll}

		// Test empty string
		response := productService.SearchProductsByName(context.Background(), "")
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Search term cannot be empty")

		// Test whitespace-only string
		response = productService.SearchProductsByName(context.Background(), "   ")
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Search term cannot be empty")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "database connection failed",
		}))

		productService := &ProductService{Collection: mt.Coll}
		response := productService.SearchProductsByName(context.Background(), "iPhone")

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to search products")
	})
}

func TestGetProductsWithPagination(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful pagination", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.products", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"name", "Product 1"},
			{"price", 99.99},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.products", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetProductsWithPagination(context.Background(), PaginationRequest{
			Page:  1,
			Limit: 10,
		})

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Products retrieved successfully", response.Message)
	})

	mt.Run("pagination validation", func(mt *mtest.T) {
		productService := &ProductService{Collection: mt.Coll}

		// Test page < 1
		response := productService.GetProductsWithPagination(context.Background(), PaginationRequest{
			Page:  0,
			Limit: 10,
		})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Page must be greater than 0")

		// Test limit < 1
		response = productService.GetProductsWithPagination(context.Background(), PaginationRequest{
			Page:  1,
			Limit: 0,
		})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Limit must be greater than 0")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "database connection failed",
		}))

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetProductsWithPagination(context.Background(), PaginationRequest{
			Page:  1,
			Limit: 10,
		})

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to retrieve products")
	})
}

func TestGetProductsByTags(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful tag filtering", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.products", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"name", "Premium Phone"},
			{"tags", []string{"premium", "bestseller"}},
			{"price", 899.99},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.products", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetProductsByTags(context.Background(), []string{"premium", "electronics"})

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Products retrieved successfully", response.Message)
	})

	mt.Run("tags validation", func(mt *mtest.T) {
		productService := &ProductService{Collection: mt.Coll}

		// Test empty tags array
		response := productService.GetProductsByTags(context.Background(), []string{})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Tags cannot be empty")

		// Test tags with empty string
		response = productService.GetProductsByTags(context.Background(), []string{"premium", "", "bestseller"})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Tags cannot contain empty strings")

		// Test tags with whitespace-only string
		response = productService.GetProductsByTags(context.Background(), []string{"premium", "   ", "bestseller"})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Tags cannot contain empty strings")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "database connection failed",
		}))

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetProductsByTags(context.Background(), []string{"premium"})

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to retrieve products")
	})
}

func TestGetTopRatedProducts(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful top rated products retrieval", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.products", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"name", "Excellent Product"},
			{"rating", 4.9},
			{"price", 299.99},
		})
		second := mtest.CreateCursorResponse(1, "test.products", mtest.NextBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"name", "Great Product"},
			{"rating", 4.8},
			{"price", 199.99},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.products", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetTopRatedProducts(context.Background(), 5)

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Top rated products retrieved successfully", response.Message)
	})

	mt.Run("limit validation", func(mt *mtest.T) {
		productService := &ProductService{Collection: mt.Coll}

		// Test limit = 0
		response := productService.GetTopRatedProducts(context.Background(), 0)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Limit must be greater than 0")

		// Test negative limit
		response = productService.GetTopRatedProducts(context.Background(), -5)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Limit must be greater than 0")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "database connection failed",
		}))

		productService := &ProductService{Collection: mt.Coll}
		response := productService.GetTopRatedProducts(context.Background(), 5)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to retrieve products")
	})
}

func TestDataStructures(t *testing.T) {
	t.Run("Product struct should have proper BSON tags", func(t *testing.T) {
		product := Product{
			ID:          primitive.NewObjectID(),
			Name:        "Test Product",
			Description: "Test Description",
			Category:    "Electronics",
			Price:       199.99,
			Rating:      4.5,
			Tags:        []string{"premium", "bestseller"},
			InStock:     true,
		}

		assert.NotEmpty(t, product.ID)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, "Test Description", product.Description)
		assert.Equal(t, "Electronics", product.Category)
		assert.Equal(t, 199.99, product.Price)
		assert.Equal(t, 4.5, product.Rating)
		assert.Equal(t, []string{"premium", "bestseller"}, product.Tags)
		assert.True(t, product.InStock)
	})

	t.Run("PaginationRequest struct should have proper fields", func(t *testing.T) {
		pagination := PaginationRequest{
			Page:  1,
			Limit: 10,
		}

		assert.Equal(t, 1, pagination.Page)
		assert.Equal(t, 10, pagination.Limit)
	})

	t.Run("Response struct should have proper fields", func(t *testing.T) {
		response := Response{
			Success: true,
			Data:    []Product{},
			Message: "test message",
			Error:   "test error",
			Code:    200,
		}

		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
		assert.Equal(t, "test message", response.Message)
		assert.Equal(t, "test error", response.Error)
		assert.Equal(t, 200, response.Code)
	})
}

func TestFunctionSignatures(t *testing.T) {
	t.Run("All required functions should exist with correct signatures", func(t *testing.T) {
		var service ProductService
		var ctx context.Context
		var category string
		var minPrice, maxPrice float64
		var searchTerm string
		var pagination PaginationRequest
		var tags []string
		var limit int

		_ = service.GetProductsByCategory(ctx, category)
		_ = service.GetProductsByPriceRange(ctx, minPrice, maxPrice)
		_ = service.SearchProductsByName(ctx, searchTerm)
		_ = service.GetProductsWithPagination(ctx, pagination)
		_ = service.GetProductsByTags(ctx, tags)
		_ = service.GetTopRatedProducts(ctx, limit)
	})
}

func TestObjectIDHandling(t *testing.T) {
	t.Run("ObjectID operations should work correctly", func(t *testing.T) {
		// Test ObjectID creation
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		assert.False(t, id1.IsZero())
		assert.NotEqual(t, id1, id2)

		// Test ObjectID hex conversion
		hexString := id1.Hex()
		assert.Equal(t, 24, len(hexString))

		// Test ObjectID from hex
		id3, err := primitive.ObjectIDFromHex(hexString)
		assert.NoError(t, err)
		assert.Equal(t, id1, id3)

		// Test invalid ObjectID
		_, err = primitive.ObjectIDFromHex("invalid-objectid")
		assert.Error(t, err)
	})
}
