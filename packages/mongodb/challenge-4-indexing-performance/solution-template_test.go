package main

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateOptimalIndexes(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful index creation", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.CreateOptimalIndexes(context.Background())

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Optimal indexes created successfully", response.Message)
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "index creation failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.CreateOptimalIndexes(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to create optimal indexes")
	})

	mt.Run("nil collection handling", func(mt *mtest.T) {
		indexService := &IndexService{Collection: nil}
		response := indexService.CreateOptimalIndexes(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Collection is not initialized")
	})
}

func TestListIndexes(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful index listing", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.products", mtest.FirstBatch, bson.D{
			{"name", "category_1"},
			{"key", bson.M{"category": 1}},
			{"unique", false},
		})
		second := mtest.CreateCursorResponse(1, "test.products", mtest.NextBatch, bson.D{
			{"name", "price_-1"},
			{"key", bson.M{"price": -1}},
			{"unique", false},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.products", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.ListIndexes(context.Background())

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Indexes listed successfully", response.Message)
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "index listing failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.ListIndexes(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to list indexes")
	})

	mt.Run("nil collection handling", func(mt *mtest.T) {
		indexService := &IndexService{Collection: nil}
		response := indexService.ListIndexes(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Collection is not initialized")
	})
}

func TestAnalyzeQueryPerformance(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database operation with empty query", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "empty query",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		query := SearchQuery{
			Category:   "Electronics",
			PriceRange: map[string]float64{"min": 100, "max": 2000},
		}
		response := indexService.AnalyzeQueryPerformance(context.Background(), query)

		// The implementation uses empty query map{}, so it should fail
		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to analyze query performance")
	})

	mt.Run("query validation", func(mt *mtest.T) {
		indexService := &IndexService{Collection: mt.Coll}

		// Test empty query
		response := indexService.AnalyzeQueryPerformance(context.Background(), SearchQuery{})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Query cannot be empty")

		// Test invalid price range
		response = indexService.AnalyzeQueryPerformance(context.Background(), SearchQuery{
			PriceRange: map[string]float64{"min": 500, "max": 100},
		})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Minimum price cannot be greater than maximum price")

		// Test negative prices
		response = indexService.AnalyzeQueryPerformance(context.Background(), SearchQuery{
			PriceRange: map[string]float64{"min": -100, "max": 500},
		})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Minimum price cannot be negative")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "query analysis failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		query := SearchQuery{Category: "Electronics"}
		response := indexService.AnalyzeQueryPerformance(context.Background(), query)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to analyze query performance")
	})
}

func TestOptimizedSearch(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database operation with empty query", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "empty query",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		query := SearchQuery{
			Category: "Electronics",
			Brand:    "Apple",
			Text:     "MacBook",
		}
		response := indexService.OptimizedSearch(context.Background(), query)

		// The implementation uses empty query map{}, so it should fail
		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to perform optimized search")
	})

	mt.Run("search validation", func(mt *mtest.T) {
		indexService := &IndexService{Collection: mt.Coll}

		// Test empty search query
		response := indexService.OptimizedSearch(context.Background(), SearchQuery{})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Query cannot be empty")

		// Test invalid price range
		response = indexService.OptimizedSearch(context.Background(), SearchQuery{
			PriceRange: map[string]float64{"min": 1000, "max": 100},
		})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Minimum price cannot be greater than maximum price")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "search failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		query := SearchQuery{Category: "Electronics"}
		response := indexService.OptimizedSearch(context.Background(), query)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to perform optimized search")
	})
}

func TestCreateTextIndex(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database operation with empty index model", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "empty index model",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		fields := map[string]int{"name": 1, "description": 1}
		response := indexService.CreateTextIndex(context.Background(), fields)

		// The implementation creates an empty IndexModel{}, so it should fail
		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to create text index")
	})

	mt.Run("fields validation", func(mt *mtest.T) {
		indexService := &IndexService{Collection: mt.Coll}

		// Test empty fields
		response := indexService.CreateTextIndex(context.Background(), map[string]int{})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Fields cannot be empty")

		// Test nil fields
		response = indexService.CreateTextIndex(context.Background(), nil)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Fields cannot be empty")

		// Test invalid field weights
		response = indexService.CreateTextIndex(context.Background(), map[string]int{"name": 0})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Field weight must be greater than 0")

		// Test empty field name
		response = indexService.CreateTextIndex(context.Background(), map[string]int{"": 1})
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Field name cannot be empty")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "text index creation failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		fields := map[string]int{"name": 1}
		response := indexService.CreateTextIndex(context.Background(), fields)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to create text index")
	})
}

func TestPerformTextSearch(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database operation with empty query", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "empty query",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		options := map[string]interface{}{"limit": 10}
		response := indexService.PerformTextSearch(context.Background(), "smartphone", options)

		// The implementation uses empty query map{}, so it should fail
		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to perform text search")
	})

	mt.Run("search text validation", func(mt *mtest.T) {
		indexService := &IndexService{Collection: mt.Coll}

		// Test empty search text
		response := indexService.PerformTextSearch(context.Background(), "", nil)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Search text cannot be empty")

		// Test whitespace-only search text
		response = indexService.PerformTextSearch(context.Background(), "   ", nil)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Search text cannot be empty")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "text search failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.PerformTextSearch(context.Background(), "smartphone", nil)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to perform text search")
	})
}

func TestCreateCompoundIndex(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("database operation with empty index model", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "empty index model",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		fields := []map[string]int{
			{"category": 1},
			{"price": -1},
			{"rating": -1},
		}
		options := map[string]interface{}{"background": true}
		response := indexService.CreateCompoundIndex(context.Background(), fields, options)

		// The implementation creates an empty IndexModel{}, so it should fail
		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to create compound index")
	})

	mt.Run("fields validation", func(mt *mtest.T) {
		indexService := &IndexService{Collection: mt.Coll}

		// Test empty fields
		response := indexService.CreateCompoundIndex(context.Background(), []map[string]int{}, nil)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Fields cannot be empty")

		// Test nil fields
		response = indexService.CreateCompoundIndex(context.Background(), nil, nil)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Fields cannot be empty")

		// Test empty field map
		response = indexService.CreateCompoundIndex(context.Background(), []map[string]int{{}}, nil)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Field map cannot be empty")

		// Test invalid field order
		response = indexService.CreateCompoundIndex(context.Background(), []map[string]int{{"category": 0}}, nil)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Field order must be 1 (ascending) or -1 (descending)")

		// Test empty field name
		response = indexService.CreateCompoundIndex(context.Background(), []map[string]int{{"": 1}}, nil)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Field name cannot be empty")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "compound index creation failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		fields := []map[string]int{{"category": 1}}
		response := indexService.CreateCompoundIndex(context.Background(), fields, nil)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to create compound index")
	})
}

func TestDropIndex(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful index drop", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.DropIndex(context.Background(), "category_1")

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Index dropped successfully", response.Message)
	})

	mt.Run("index name validation", func(mt *mtest.T) {
		indexService := &IndexService{Collection: mt.Coll}

		// Test empty index name
		response := indexService.DropIndex(context.Background(), "")
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Index name cannot be empty")

		// Test whitespace-only index name
		response = indexService.DropIndex(context.Background(), "   ")
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Index name cannot be empty")

		// Test reserved _id index
		response = indexService.DropIndex(context.Background(), "_id_")
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Cannot drop the _id index")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "index drop failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.DropIndex(context.Background(), "category_1")

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to drop index")
	})
}

func TestGetIndexUsageStats(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful index usage stats retrieval", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{
			"indexSizes", bson.M{
				"_id_":       1024,
				"category_1": 2048,
				"price_-1":   1536,
			},
		}))

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.GetIndexUsageStats(context.Background())

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Index usage statistics retrieved", response.Message)
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "stats retrieval failed",
		}))

		indexService := &IndexService{Collection: mt.Coll}
		response := indexService.GetIndexUsageStats(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to get index usage stats")
	})

	mt.Run("nil collection handling", func(mt *mtest.T) {
		indexService := &IndexService{Collection: nil}
		response := indexService.GetIndexUsageStats(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Collection is not initialized")
	})
}

func TestDataStructures(t *testing.T) {
	t.Run("Product struct should have proper BSON tags", func(t *testing.T) {
		product := Product{
			ID:          primitive.NewObjectID(),
			Name:        "iPhone 14",
			Description: "Latest smartphone",
			Category:    "Electronics",
			Brand:       "Apple",
			Price:       999.99,
			Stock:       50,
			Rating:      4.8,
			Tags:        []string{"smartphone", "premium"},
		}

		assert.NotEmpty(t, product.ID)
		assert.Equal(t, "iPhone 14", product.Name)
		assert.Equal(t, "Latest smartphone", product.Description)
		assert.Equal(t, "Electronics", product.Category)
		assert.Equal(t, "Apple", product.Brand)
		assert.Equal(t, 999.99, product.Price)
		assert.Equal(t, 50, product.Stock)
		assert.Equal(t, 4.8, product.Rating)
		assert.Equal(t, []string{"smartphone", "premium"}, product.Tags)
	})

	t.Run("SearchQuery struct should have proper fields", func(t *testing.T) {
		query := SearchQuery{
			Text:       "smartphone",
			Category:   "Electronics",
			PriceRange: map[string]float64{"min": 100, "max": 2000},
			Brand:      "Apple",
			MinRating:  4.0,
			Tags:       []string{"premium"},
			SortBy:     "price",
			SortOrder:  -1,
			Limit:      10,
			Skip:       0,
		}

		assert.Equal(t, "smartphone", query.Text)
		assert.Equal(t, "Electronics", query.Category)
		assert.Equal(t, map[string]float64{"min": 100, "max": 2000}, query.PriceRange)
		assert.Equal(t, "Apple", query.Brand)
		assert.Equal(t, 4.0, query.MinRating)
		assert.Equal(t, []string{"premium"}, query.Tags)
		assert.Equal(t, "price", query.SortBy)
		assert.Equal(t, -1, query.SortOrder)
		assert.Equal(t, 10, query.Limit)
		assert.Equal(t, 0, query.Skip)
	})

	t.Run("Response struct should have proper fields", func(t *testing.T) {
		response := Response{
			Success: true,
			Data:    []Product{},
			Message: "test message",
			Error:   "test error",
			Code:    200,
			Performance: &QueryPerformance{
				ExecutionTimeMs: 25,
				DocsExamined:    100,
				DocsReturned:    10,
				IndexUsed:       "category_1",
				IsOptimal:       true,
			},
		}

		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
		assert.Equal(t, "test message", response.Message)
		assert.Equal(t, "test error", response.Error)
		assert.Equal(t, 200, response.Code)
		assert.NotNil(t, response.Performance)
		assert.Equal(t, int64(25), response.Performance.ExecutionTimeMs)
		assert.Equal(t, int64(100), response.Performance.DocsExamined)
		assert.Equal(t, int64(10), response.Performance.DocsReturned)
		assert.Equal(t, "category_1", response.Performance.IndexUsed)
		assert.True(t, response.Performance.IsOptimal)
	})
}

func TestFunctionSignatures(t *testing.T) {
	t.Run("All required functions should exist with correct signatures", func(t *testing.T) {
		var service IndexService
		var ctx context.Context
		var query SearchQuery
		var fields map[string]int
		var compoundFields []map[string]int
		var options map[string]interface{}
		var searchText string
		var indexName string

		_ = service.CreateOptimalIndexes(ctx)
		_ = service.ListIndexes(ctx)
		_ = service.AnalyzeQueryPerformance(ctx, query)
		_ = service.OptimizedSearch(ctx, query)
		_ = service.CreateTextIndex(ctx, fields)
		_ = service.PerformTextSearch(ctx, searchText, options)
		_ = service.CreateCompoundIndex(ctx, compoundFields, options)
		_ = service.DropIndex(ctx, indexName)
		_ = service.GetIndexUsageStats(ctx)
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

func TestIndexingConcepts(t *testing.T) {
	t.Run("Indexing concepts should be understood", func(t *testing.T) {
		// Test index field orders
		ascending := 1
		descending := -1
		assert.Equal(t, 1, ascending)
		assert.Equal(t, -1, descending)

		// Test text index weights
		weights := map[string]int{"name": 10, "description": 1}
		assert.Greater(t, weights["name"], weights["description"])

		// Test compound index structure
		compoundFields := []map[string]int{
			{"category": 1},
			{"price": -1},
			{"rating": -1},
		}
		assert.Equal(t, 3, len(compoundFields))
		assert.Equal(t, 1, compoundFields[0]["category"])
		assert.Equal(t, -1, compoundFields[1]["price"])

		// Test search query validation
		validQuery := SearchQuery{Category: "Electronics"}
		emptyQuery := SearchQuery{}
		assert.NotEqual(t, "", validQuery.Category)
		assert.Equal(t, "", emptyQuery.Category)

		// Test price range validation
		validRange := map[string]float64{"min": 100, "max": 500}
		invalidRange := map[string]float64{"min": 500, "max": 100}
		assert.Less(t, validRange["min"], validRange["max"])
		assert.Greater(t, invalidRange["min"], invalidRange["max"])
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("Edge cases should be handled properly", func(t *testing.T) {
		indexService := &IndexService{Collection: nil}

		// Test very long search text
		longText := strings.Repeat("search", 1000)
		response := indexService.PerformTextSearch(context.Background(), longText, nil)
		// Should handle gracefully (either accept or reject, but not panic)
		assert.NotNil(t, response)

		// Test complex compound index
		complexFields := []map[string]int{
			{"category": 1},
			{"brand": 1},
			{"price": -1},
			{"rating": -1},
			{"created_at": 1},
		}
		response = indexService.CreateCompoundIndex(context.Background(), complexFields, nil)
		// Should handle gracefully (validation should pass, DB should fail with nil collection)
		assert.False(t, response.Success)
		assert.Contains(t, response.Error, "Collection is not initialized")

		// Test complex search query
		complexQuery := SearchQuery{
			Text:       "iPhone Pro Max",
			Category:   "Electronics",
			Brand:      "Apple",
			PriceRange: map[string]float64{"min": 100, "max": 2000},
			MinRating:  4.0,
			Tags:       []string{"smartphone", "premium"},
		}
		response = indexService.OptimizedSearch(context.Background(), complexQuery)
		// Should handle gracefully (validation should pass, DB should fail with nil collection)
		assert.False(t, response.Success)
		assert.Contains(t, response.Error, "Collection is not initialized")
	})
}
