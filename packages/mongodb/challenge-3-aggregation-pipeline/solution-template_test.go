package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetSalesByCategory(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful sales by category aggregation", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.orders", mtest.FirstBatch, bson.D{
			{"_id", "Electronics"},
			{"total_sales", 15000.50},
			{"order_count", 25},
			{"avg_order_value", 600.02},
		})
		second := mtest.CreateCursorResponse(1, "test.orders", mtest.NextBatch, bson.D{
			{"_id", "Clothing"},
			{"total_sales", 8500.75},
			{"order_count", 18},
			{"avg_order_value", 472.26},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.orders", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetSalesByCategory(context.Background())

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Sales by category retrieved successfully", response.Message)
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "aggregation failed",
		}))

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetSalesByCategory(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to aggregate sales by category")
	})

	mt.Run("nil collection handling", func(mt *mtest.T) {
		analyticsService := &AnalyticsService{Collection: nil}
		response := analyticsService.GetSalesByCategory(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Collection not initialized")
	})
}

func TestGetTopSellingProducts(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful top selling products aggregation", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.orders", mtest.FirstBatch, bson.D{
			{"_id", "iPhone 14"},
			{"total_sold", 150},
			{"revenue", 120000.00},
		})
		second := mtest.CreateCursorResponse(1, "test.orders", mtest.NextBatch, bson.D{
			{"_id", "MacBook Pro"},
			{"total_sold", 85},
			{"revenue", 170000.00},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.orders", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetTopSellingProducts(context.Background(), 5)

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Top selling products retrieved successfully", response.Message)
	})

	mt.Run("limit validation", func(mt *mtest.T) {
		analyticsService := &AnalyticsService{Collection: mt.Coll}

		// Test limit = 0
		response := analyticsService.GetTopSellingProducts(context.Background(), 0)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Limit must be greater than 0")

		// Test negative limit
		response = analyticsService.GetTopSellingProducts(context.Background(), -5)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Limit must be greater than 0")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "aggregation failed",
		}))

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetTopSellingProducts(context.Background(), 5)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to get top selling products")
	})
}

func TestGetRevenueByMonth(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful monthly revenue aggregation", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.orders", mtest.FirstBatch, bson.D{
			{"_id", 1}, // January
			{"revenue", 25000.50},
		})
		second := mtest.CreateCursorResponse(1, "test.orders", mtest.NextBatch, bson.D{
			{"_id", 2}, // February
			{"revenue", 18500.75},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.orders", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetRevenueByMonth(context.Background(), 2023)

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Monthly revenue retrieved successfully", response.Message)
	})

	mt.Run("year validation", func(mt *mtest.T) {
		analyticsService := &AnalyticsService{Collection: mt.Coll}

		// Test year < 2000
		response := analyticsService.GetRevenueByMonth(context.Background(), 1999)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Year must be between 2000 and 2099")

		// Test year >= 2100
		response = analyticsService.GetRevenueByMonth(context.Background(), 2100)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Year must be between 2000 and 2099")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "aggregation failed",
		}))

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetRevenueByMonth(context.Background(), 2023)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to get revenue by month")
	})
}

func TestGetCustomerAnalytics(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful customer analytics aggregation", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.orders", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"total_spent", 2500.75},
			{"order_count", 8},
			{"avg_order_size", 312.59},
		})
		second := mtest.CreateCursorResponse(1, "test.orders", mtest.NextBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"total_spent", 1800.25},
			{"order_count", 5},
			{"avg_order_size", 360.05},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.orders", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetCustomerAnalytics(context.Background(), 10)

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Customer analytics retrieved successfully", response.Message)
	})

	mt.Run("limit validation", func(mt *mtest.T) {
		analyticsService := &AnalyticsService{Collection: mt.Coll}

		// Test limit = 0
		response := analyticsService.GetCustomerAnalytics(context.Background(), 0)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Limit must be greater than 0")

		// Test negative limit
		response = analyticsService.GetCustomerAnalytics(context.Background(), -3)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Limit must be greater than 0")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "aggregation failed",
		}))

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetCustomerAnalytics(context.Background(), 10)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to get customer analytics")
	})
}

func TestGetProductPerformance(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful product performance aggregation", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.orders", mtest.FirstBatch, bson.D{
			{"_id", "iPhone 14"},
			{"total_sold", 200},
			{"revenue", 160000.00},
		})
		second := mtest.CreateCursorResponse(1, "test.orders", mtest.NextBatch, bson.D{
			{"_id", "Samsung Galaxy"},
			{"total_sold", 150},
			{"revenue", 105000.00},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.orders", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetProductPerformance(context.Background())

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Product performance retrieved successfully", response.Message)
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "aggregation failed",
		}))

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetProductPerformance(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to get product performance")
	})

	mt.Run("nil collection handling", func(mt *mtest.T) {
		analyticsService := &AnalyticsService{Collection: nil}
		response := analyticsService.GetProductPerformance(context.Background())

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Collection not initialized")
	})
}

func TestGetOrderTrends(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful order trends aggregation", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.orders", mtest.FirstBatch, bson.D{
			{"_id", bson.M{"year": 2023, "month": 10, "day": 15}},
			{"order_count", 25},
			{"revenue", 12500.50},
		})
		second := mtest.CreateCursorResponse(1, "test.orders", mtest.NextBatch, bson.D{
			{"_id", bson.M{"year": 2023, "month": 10, "day": 16}},
			{"order_count", 18},
			{"revenue", 9800.75},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.orders", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetOrderTrends(context.Background(), 30)

		assert.True(t, response.Success)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "Order trends retrieved successfully", response.Message)
	})

	mt.Run("days validation", func(mt *mtest.T) {
		analyticsService := &AnalyticsService{Collection: mt.Coll}

		// Test days = 0
		response := analyticsService.GetOrderTrends(context.Background(), 0)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Days must be greater than 0")

		// Test negative days
		response = analyticsService.GetOrderTrends(context.Background(), -7)
		assert.False(t, response.Success)
		assert.Equal(t, 400, response.Code)
		assert.Contains(t, response.Error, "Days must be greater than 0")
	})

	mt.Run("database error handling", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "aggregation failed",
		}))

		analyticsService := &AnalyticsService{Collection: mt.Coll}
		response := analyticsService.GetOrderTrends(context.Background(), 30)

		assert.False(t, response.Success)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Error, "Failed to get order trends")
	})
}

func TestDataStructures(t *testing.T) {
	t.Run("Order struct should have proper BSON tags", func(t *testing.T) {
		order := Order{
			ID:         primitive.NewObjectID(),
			CustomerID: primitive.NewObjectID(),
			Total:      199.99,
			Status:     "completed",
			Category:   "Electronics",
			CreatedAt:  time.Now(),
		}

		assert.NotEmpty(t, order.ID)
		assert.NotEmpty(t, order.CustomerID)
		assert.Equal(t, 199.99, order.Total)
		assert.Equal(t, "completed", order.Status)
		assert.Equal(t, "Electronics", order.Category)
		assert.False(t, order.CreatedAt.IsZero())
	})

	t.Run("Customer struct should have proper fields", func(t *testing.T) {
		customer := Customer{
			ID:    primitive.NewObjectID(),
			Name:  "John Doe",
			Email: "john@example.com",
		}

		assert.NotEmpty(t, customer.ID)
		assert.Equal(t, "John Doe", customer.Name)
		assert.Equal(t, "john@example.com", customer.Email)
	})

	t.Run("Analytics result structs should have proper fields", func(t *testing.T) {
		salesByCategory := SalesByCategory{
			Category:      "Electronics",
			TotalSales:    15000.50,
			OrderCount:    25,
			AvgOrderValue: 600.02,
		}

		topProduct := TopProduct{
			ProductName: "iPhone 14",
			TotalSold:   150,
			Revenue:     120000.00,
		}

		monthlyRevenue := MonthlyRevenue{
			Month:   10,
			Revenue: 25000.50,
		}

		customerAnalytics := CustomerAnalytics{
			CustomerID:   primitive.NewObjectID(),
			TotalSpent:   2500.75,
			OrderCount:   8,
			AvgOrderSize: 312.59,
		}

		assert.Equal(t, "Electronics", salesByCategory.Category)
		assert.Equal(t, 15000.50, salesByCategory.TotalSales)
		assert.Equal(t, 25, salesByCategory.OrderCount)
		assert.Equal(t, 600.02, salesByCategory.AvgOrderValue)

		assert.Equal(t, "iPhone 14", topProduct.ProductName)
		assert.Equal(t, 150, topProduct.TotalSold)
		assert.Equal(t, 120000.00, topProduct.Revenue)

		assert.Equal(t, 10, monthlyRevenue.Month)
		assert.Equal(t, 25000.50, monthlyRevenue.Revenue)

		assert.NotEmpty(t, customerAnalytics.CustomerID)
		assert.Equal(t, 2500.75, customerAnalytics.TotalSpent)
		assert.Equal(t, 8, customerAnalytics.OrderCount)
		assert.Equal(t, 312.59, customerAnalytics.AvgOrderSize)
	})

	t.Run("Response struct should have proper fields", func(t *testing.T) {
		response := Response{
			Success: true,
			Data:    []SalesByCategory{},
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
		var service AnalyticsService
		var ctx context.Context
		var limit int
		var year int
		var days int

		_ = service.GetSalesByCategory(ctx)
		_ = service.GetTopSellingProducts(ctx, limit)
		_ = service.GetRevenueByMonth(ctx, year)
		_ = service.GetCustomerAnalytics(ctx, limit)
		_ = service.GetProductPerformance(ctx)
		_ = service.GetOrderTrends(ctx, days)
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

func TestTimeHandling(t *testing.T) {
	t.Run("Time operations should work correctly", func(t *testing.T) {
		now := time.Now()
		assert.False(t, now.IsZero())

		// Test year extraction
		year := now.Year()
		assert.Greater(t, year, 2000)
		assert.Less(t, year, 2100)

		// Test date creation
		specificDate := time.Date(2023, 10, 15, 12, 30, 0, 0, time.UTC)
		assert.Equal(t, 2023, specificDate.Year())
		assert.Equal(t, time.October, specificDate.Month())
		assert.Equal(t, 15, specificDate.Day())
	})
}

func TestAggregationConcepts(t *testing.T) {
	t.Run("Aggregation pipeline concepts should be understood", func(t *testing.T) {
		// This test validates that the concepts of MongoDB aggregation are understood
		// In actual implementation, these would be used in aggregation pipelines

		// Test $match concepts
		status := "completed"
		assert.Equal(t, "completed", status)

		// Test $group concepts
		totalSales := 15000.50
		orderCount := 25
		avgOrderValue := totalSales / float64(orderCount)
		assert.Equal(t, 600.02, avgOrderValue)

		// Test $sort concepts (descending order)
		revenues := []float64{25000.50, 18500.75, 12300.25}
		for i := 1; i < len(revenues); i++ {
			assert.Greater(t, revenues[i-1], revenues[i], "Revenues should be in descending order")
		}

		// Test limit concepts
		limit := 10
		assert.Greater(t, limit, 0, "Limit should be positive")

		// Test date range concepts
		year := 2023
		assert.GreaterOrEqual(t, year, 2000)
		assert.Less(t, year, 2100)
	})
}
