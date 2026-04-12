# Hints for Challenge 3: Aggregation Pipeline

## Hint 1: Understanding Aggregation Pipelines

MongoDB aggregation pipelines process data through multiple stages, like an assembly line:

```go
// Basic pipeline structure
pipeline := mongo.Pipeline{
    // Stage 1: Filter documents
    {{"$match", bson.M{"status": "completed"}}},
    
    // Stage 2: Group and calculate
    {{"$group", bson.M{
        "_id":    "$category",
        "total":  bson.M{"$sum": "$amount"},
        "count":  bson.M{"$sum": 1},
    }}},
    
    // Stage 3: Sort results
    {{"$sort", bson.M{"total": -1}}},
}

// Execute pipeline
cursor, err := collection.Aggregate(ctx, pipeline)
```

## Hint 2: Sales by Category Aggregation

Build a pipeline to group orders by category and calculate metrics:

```go
func (as *AnalyticsService) GetSalesByCategory(ctx context.Context) Response {
    pipeline := mongo.Pipeline{
        // Only include completed orders
        {{"$match", bson.M{"status": "completed"}}},
        
        // Group by category
        {{"$group", bson.M{
            "_id":         "$category",
            "total_sales": bson.M{"$sum": "$total"},
            "order_count": bson.M{"$sum": 1},
        }}},
        
        // Calculate average order value
        {{"$addFields", bson.M{
            "avg_order_value": bson.M{
                "$divide": []interface{}{"$total_sales", "$order_count"},
            },
        }}},
        
        // Sort by total sales (highest first)
        {{"$sort", bson.M{"total_sales": -1}}},
    }
    
    cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    defer cursor.Close(ctx)
    
    var results []SalesByCategory
    if err = cursor.All(ctx, &results); err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    return Response{Success: true, Data: results, Code: 200}
}
```

## Hint 3: Top Selling Products

Find products with highest sales using grouping and sorting:

```go
func (as *AnalyticsService) GetTopSellingProducts(ctx context.Context, limit int) Response {
    if limit <= 0 {
        limit = 10 // Default limit
    }
    
    pipeline := mongo.Pipeline{
        // Filter completed orders
        {{"$match", bson.M{"status": "completed"}}},
        
        // Group by product_id
        {{"$group", bson.M{
            "_id":            "$product_id",
            "total_quantity": bson.M{"$sum": "$quantity"},
            "total_revenue":  bson.M{"$sum": "$total"},
            "order_count":    bson.M{"$sum": 1},
        }}},
        
        // Sort by quantity sold (descending)
        {{"$sort", bson.M{"total_quantity": -1}}},
        
        // Limit results
        {{"$limit", limit}},
    }
    
    cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
    // ... handle results
}
```

## Hint 4: Monthly Revenue Analysis

Use date operators to extract year and month from dates:

```go
func (as *AnalyticsService) GetRevenueByMonth(ctx context.Context, year int) Response {
    // Build match stage - filter by year if specified
    matchStage := bson.M{"status": "completed"}
    if year > 0 {
        matchStage["$expr"] = bson.M{
            "$eq": []interface{}{
                bson.M{"$year": "$order_date"},
                year,
            },
        }
    }
    
    pipeline := mongo.Pipeline{
        {{"$match", matchStage}},
        
        // Group by year and month
        {{"$group", bson.M{
            "_id": bson.M{
                "year":  bson.M{"$year": "$order_date"},
                "month": bson.M{"$month": "$order_date"},
            },
            "total_revenue": bson.M{"$sum": "$total"},
            "order_count":   bson.M{"$sum": 1},
        }}},
        
        // Reshape the output
        {{"$project", bson.M{
            "_id":           0,
            "year":          "$_id.year",
            "month":         "$_id.month",
            "total_revenue": 1,
            "order_count":   1,
        }}},
        
        // Sort by year and month
        {{"$sort", bson.M{"year": 1, "month": 1}}},
    }
    
    cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
    // ... handle results
}
```

## Hint 5: Customer Analytics with Advanced Grouping

Calculate customer metrics including first and last order dates:

```go
func (as *AnalyticsService) GetCustomerAnalytics(ctx context.Context, limit int) Response {
    pipeline := mongo.Pipeline{
        {{"$match", bson.M{"status": "completed"}}},
        
        // Group by customer_id with advanced calculations
        {{"$group", bson.M{
            "_id":         "$customer_id",
            "total_spent": bson.M{"$sum": "$total"},
            "order_count": bson.M{"$sum": 1},
            "first_order": bson.M{"$min": "$order_date"},
            "last_order":  bson.M{"$max": "$order_date"},
        }}},
        
        // Calculate average order value
        {{"$addFields", bson.M{
            "avg_order_value": bson.M{
                "$divide": []interface{}{"$total_spent", "$order_count"},
            },
        }}},
        
        // Sort by total spent (highest spenders first)
        {{"$sort", bson.M{"total_spent": -1}}},
        
        // Limit to top customers
        {{"$limit", limit}},
    }
    
    cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
    // ... handle results
}
```

## Hint 6: Using $lookup for Joins

Combine order data with customer information:

```go
func (as *AnalyticsService) GetProductPerformance(ctx context.Context) Response {
    pipeline := mongo.Pipeline{
        // Match completed orders
        {{"$match", bson.M{"status": "completed"}}},
        
        // Join with customers collection
        {{"$lookup", bson.M{
            "from":         "customers",
            "localField":   "customer_id",
            "foreignField": "_id",
            "as":           "customer_info",
        }}},
        
        // Unwind customer info (convert array to object)
        {{"$unwind", "$customer_info"}},
        
        // Group by product and include customer demographics
        {{"$group", bson.M{
            "_id": "$product_id",
            "total_revenue": bson.M{"$sum": "$total"},
            "total_quantity": bson.M{"$sum": "$quantity"},
            "unique_customers": bson.M{"$addToSet": "$customer_id"},
            "customer_locations": bson.M{"$addToSet": "$customer_info.location"},
        }}},
        
        // Add calculated fields
        {{"$addFields", bson.M{
            "customer_count": bson.M{"$size": "$unique_customers"},
            "location_diversity": bson.M{"$size": "$customer_locations"},
        }}},
        
        {{"$sort", bson.M{"total_revenue": -1}}},
    }
    
    cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
    // ... handle results
}
```

## Hint 7: Time-Based Trend Analysis

Analyze order trends over recent days:

```go
func (as *AnalyticsService) GetOrderTrends(ctx context.Context, days int) Response {
    if days <= 0 {
        days = 30 // Default to 30 days
    }
    
    // Calculate date threshold
    cutoffDate := time.Now().AddDate(0, 0, -days)
    
    pipeline := mongo.Pipeline{
        // Filter recent orders
        {{"$match", bson.M{
            "status": "completed",
            "order_date": bson.M{"$gte": cutoffDate},
        }}},
        
        // Group by date (year-month-day)
        {{"$group", bson.M{
            "_id": bson.M{
                "year":  bson.M{"$year": "$order_date"},
                "month": bson.M{"$month": "$order_date"},
                "day":   bson.M{"$dayOfMonth": "$order_date"},
            },
            "daily_revenue": bson.M{"$sum": "$total"},
            "daily_orders":  bson.M{"$sum": 1},
            "avg_order_value": bson.M{"$avg": "$total"},
        }}},
        
        // Reshape and add date field
        {{"$addFields", bson.M{
            "date": bson.M{
                "$dateFromParts": bson.M{
                    "year":  "$_id.year",
                    "month": "$_id.month",
                    "day":   "$_id.day",
                },
            },
        }}},
        
        // Sort by date
        {{"$sort", bson.M{"date": 1}}},
        
        // Project final fields
        {{"$project", bson.M{
            "_id":             0,
            "date":            1,
            "daily_revenue":   1,
            "daily_orders":    1,
            "avg_order_value": 1,
        }}},
    }
    
    cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
    // ... handle results
}
```

## Hint 8: Error Handling and Validation

Always validate inputs and handle aggregation errors:

```go
func (as *AnalyticsService) GetTopSellingProducts(ctx context.Context, limit int) Response {
    // Validate and set defaults
    if limit <= 0 {
        limit = 10
    }
    if limit > 100 {
        limit = 100 // Prevent excessive results
    }
    
    pipeline := mongo.Pipeline{
        // ... your pipeline stages
    }
    
    cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
    if err != nil {
        // Log the error for debugging
        log.Printf("Aggregation failed: %v", err)
        return Response{
            Success: false,
            Error:   "Failed to retrieve top products",
            Code:    500,
        }
    }
    defer cursor.Close(ctx) // Always close cursor
    
    var results []TopProduct
    if err = cursor.All(ctx, &results); err != nil {
        return Response{
            Success: false,
            Error:   "Failed to decode results",
            Code:    500,
        }
    }
    
    // Handle empty results
    if results == nil {
        results = []TopProduct{}
    }
    
    return Response{
        Success: true,
        Data:    results,
        Message: fmt.Sprintf("Retrieved top %d products", len(results)),
        Code:    200,
    }
}
```

## Hint 9: Performance Optimization

Optimize aggregation pipelines for better performance:

```go
// 1. Use $match early to filter documents
pipeline := mongo.Pipeline{
    // Filter first to reduce data processed in later stages
    {{"$match", bson.M{
        "status": "completed",
        "order_date": bson.M{"$gte": startDate},
    }}},
    
    // Then group and calculate
    {{"$group", bson.M{...}}},
}

// 2. Use indexes for $match and $sort stages
// Create indexes on frequently filtered/sorted fields:
// db.orders.createIndex({"status": 1, "order_date": -1})
// db.orders.createIndex({"category": 1, "status": 1})

// 3. Use $project to reduce data transfer
{{"$project", bson.M{
    "category":    1,
    "total":       1,
    "order_date":  1,
    // Only include fields you need
}}},

// 4. Use $limit early when possible
{{"$limit", 1000}}, // Limit before expensive operations
```

## Hint 10: Testing Aggregation Results

Verify your aggregation logic with test data:

```go
func TestSalesByCategory(t *testing.T) {
    // Setup test data with known values
    testOrders := []Order{
        {Category: "Electronics", Total: 1000, Status: "completed"},
        {Category: "Electronics", Total: 500, Status: "completed"},
        {Category: "Books", Total: 25, Status: "completed"},
    }
    
    // Expected results
    expected := []SalesByCategory{
        {Category: "Electronics", TotalSales: 1500, OrderCount: 2, AvgOrderValue: 750},
        {Category: "Books", TotalSales: 25, OrderCount: 1, AvgOrderValue: 25},
    }
    
    // Run aggregation and verify results
    result := analyticsService.GetSalesByCategory(ctx)
    
    // Verify calculations are correct
    actual := result.Data.([]SalesByCategory)
    if len(actual) != len(expected) {
        t.Errorf("Expected %d categories, got %d", len(expected), len(actual))
    }
    
    // Verify specific calculations
    for i, category := range actual {
        if category.TotalSales != expected[i].TotalSales {
            t.Errorf("Category %s: expected total sales %.2f, got %.2f", 
                category.Category, expected[i].TotalSales, category.TotalSales)
        }
    }
}
```

Ready to build powerful analytics with MongoDB aggregation? Start with simple grouping and build up to complex multi-stage pipelines! 🚀
