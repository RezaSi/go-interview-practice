package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Order represents an order document
type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	Total      float64            `bson:"total" json:"total"`
	Status     string             `bson:"status" json:"status"`
	Category   string             `bson:"category" json:"category"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

// Customer represents a customer document
type Customer struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
}

// SalesByCategory represents sales analytics by category
type SalesByCategory struct {
	Category      string  `bson:"_id" json:"category"`
	TotalSales    float64 `bson:"total_sales" json:"total_sales"`
	OrderCount    int     `bson:"order_count" json:"order_count"`
	AvgOrderValue float64 `bson:"avg_order_value" json:"avg_order_value"`
}

// TopProduct represents top-selling product analytics
type TopProduct struct {
	ProductName string  `bson:"_id" json:"product_name"`
	TotalSold   int     `bson:"total_sold" json:"total_sold"`
	Revenue     float64 `bson:"revenue" json:"revenue"`
}

// MonthlyRevenue represents monthly revenue analytics
type MonthlyRevenue struct {
	Month   int     `bson:"_id" json:"month"`
	Revenue float64 `bson:"revenue" json:"revenue"`
}

// CustomerAnalytics represents customer spending analytics
type CustomerAnalytics struct {
	CustomerID   primitive.ObjectID `bson:"_id" json:"customer_id"`
	TotalSpent   float64            `bson:"total_spent" json:"total_spent"`
	OrderCount   int                `bson:"order_count" json:"order_count"`
	AvgOrderSize float64            `bson:"avg_order_size" json:"avg_order_size"`
}

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// AnalyticsService handles analytics operations
type AnalyticsService struct {
	Collection *mongo.Collection
}

// GetSalesByCategory calculates total sales per category using aggregation
func (as *AnalyticsService) GetSalesByCategory(ctx context.Context) Response {
	if as.Collection == nil {
		return Response{
			Success: false,
			Error:   "Collection not initialized",
			Code:    500,
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"status": "completed"}}},
		{{"$group", bson.M{
			"_id":             "$category",
			"total_sales":     bson.M{"$sum": "$total"},
			"order_count":     bson.M{"$sum": 1},
			"avg_order_value": bson.M{"$avg": "$total"},
		}}},
		{{"$sort", bson.M{"total_sales": -1}}},
	}

	cursor, err := as.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to aggregate sales by category: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	return Response{
		Success: true,
		Message: "Sales by category retrieved successfully",
		Code:    200,
	}
}

// GetTopSellingProducts retrieves top-selling products
func (as *AnalyticsService) GetTopSellingProducts(ctx context.Context, limit int) Response {
	if limit <= 0 {
		return Response{
			Success: false,
			Error:   "Limit must be greater than 0",
			Code:    400,
		}
	}

	if as.Collection == nil {
		return Response{
			Success: false,
			Error:   "Collection not initialized",
			Code:    500,
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"status": "completed"}}},
		{{"$group", bson.M{
			"_id":        "$product_name",
			"total_sold": bson.M{"$sum": "$quantity"},
			"revenue":    bson.M{"$sum": "$total"},
		}}},
		{{"$sort", bson.M{"total_sold": -1}}},
		{{"$limit", limit}},
	}

	cursor, err := as.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get top selling products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	return Response{
		Success: true,
		Message: "Top selling products retrieved successfully",
		Code:    200,
	}
}

// GetRevenueByMonth calculates monthly revenue
func (as *AnalyticsService) GetRevenueByMonth(ctx context.Context, year int) Response {
	if year < 2000 || year >= 2100 {
		return Response{
			Success: false,
			Error:   "Year must be between 2000 and 2099",
			Code:    400,
		}
	}

	if as.Collection == nil {
		return Response{
			Success: false,
			Error:   "Collection not initialized",
			Code:    500,
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"status": "completed",
			"created_at": bson.M{
				"$gte": time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
				"$lt":  time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}}},
		{{"$group", bson.M{
			"_id":     bson.M{"$month": "$created_at"},
			"revenue": bson.M{"$sum": "$total"},
		}}},
		{{"$sort", bson.M{"_id": 1}}},
	}

	cursor, err := as.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get revenue by month: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	return Response{
		Success: true,
		Message: "Monthly revenue retrieved successfully",
		Code:    200,
	}
}

// GetCustomerAnalytics retrieves customer spending analytics
func (as *AnalyticsService) GetCustomerAnalytics(ctx context.Context, limit int) Response {
	if limit <= 0 {
		return Response{
			Success: false,
			Error:   "Limit must be greater than 0",
			Code:    400,
		}
	}

	if as.Collection == nil {
		return Response{
			Success: false,
			Error:   "Collection not initialized",
			Code:    500,
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"status": "completed"}}},
		{{"$group", bson.M{
			"_id":            "$customer_id",
			"total_spent":    bson.M{"$sum": "$total"},
			"order_count":    bson.M{"$sum": 1},
			"avg_order_size": bson.M{"$avg": "$total"},
		}}},
		{{"$sort", bson.M{"total_spent": -1}}},
		{{"$limit", limit}},
	}

	cursor, err := as.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get customer analytics: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	return Response{
		Success: true,
		Message: "Customer analytics retrieved successfully",
		Code:    200,
	}
}

// GetProductPerformance analyzes product performance
func (as *AnalyticsService) GetProductPerformance(ctx context.Context) Response {
	if as.Collection == nil {
		return Response{
			Success: false,
			Error:   "Collection not initialized",
			Code:    500,
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"status": "completed"}}},
		{{"$group", bson.M{
			"_id":        "$product_name",
			"total_sold": bson.M{"$sum": "$quantity"},
			"revenue":    bson.M{"$sum": "$total"},
		}}},
		{{"$sort", bson.M{"revenue": -1}}},
	}

	cursor, err := as.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get product performance: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	return Response{
		Success: true,
		Message: "Product performance retrieved successfully",
		Code:    200,
	}
}

// GetOrderTrends analyzes order trends over time
func (as *AnalyticsService) GetOrderTrends(ctx context.Context, days int) Response {
	if days <= 0 {
		return Response{
			Success: false,
			Error:   "Days must be greater than 0",
			Code:    400,
		}
	}

	if as.Collection == nil {
		return Response{
			Success: false,
			Error:   "Collection not initialized",
			Code:    500,
		}
	}

	startDate := time.Now().AddDate(0, 0, -days)
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"created_at": bson.M{"$gte": startDate},
		}}},
		{{"$group", bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$created_at"},
				"month": bson.M{"$month": "$created_at"},
				"day":   bson.M{"$dayOfMonth": "$created_at"},
			},
			"order_count": bson.M{"$sum": 1},
			"revenue":     bson.M{"$sum": "$total"},
		}}},
		{{"$sort", bson.M{"_id": 1}}},
	}

	cursor, err := as.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get order trends: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	return Response{
		Success: true,
		Message: "Order trends retrieved successfully",
		Code:    200,
	}
}

// ConnectMongoDB establishes connection to MongoDB
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	return nil, fmt.Errorf("ConnectMongoDB not implemented")
}

func main() {
	fmt.Println("MongoDB Aggregation Pipeline Challenge")
}
