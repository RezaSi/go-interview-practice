# Learning: MongoDB Aggregation Pipeline Mastery

## 🎯 **What You'll Learn**

This challenge teaches MongoDB's most powerful feature: the aggregation pipeline. You'll learn to build complex data processing workflows that transform raw data into actionable business insights, just like major analytics platforms.

## 🔄 **Understanding Aggregation Pipelines**

### **Pipeline Concept**

Think of aggregation pipelines as an assembly line for data processing. Each stage transforms the data and passes it to the next stage:

```go
// Pipeline stages flow like an assembly line
Documents → $match → $group → $sort → $limit → Results

// Example: Sales analysis pipeline
Orders Collection
    ↓ $match (filter completed orders)
Completed Orders  
    ↓ $group (group by category, sum totals)
Category Totals
    ↓ $sort (sort by total sales desc)
Sorted Results
    ↓ $limit (top 10 categories)
Final Results
```

### **Basic Pipeline Structure**

```go
pipeline := mongo.Pipeline{
    // Stage 1: Filter documents
    {{"$match", bson.M{"status": "active"}}},
    
    // Stage 2: Group and calculate
    {{"$group", bson.M{
        "_id":    "$category",
        "total":  bson.M{"$sum": "$amount"},
        "count":  bson.M{"$sum": 1},
    }}},
    
    // Stage 3: Sort results
    {{"$sort", bson.M{"total": -1}}},
    
    // Stage 4: Limit output
    {{"$limit", 10}},
}

// Execute the pipeline
cursor, err := collection.Aggregate(ctx, pipeline)
```

## 🔍 **Core Aggregation Stages**

### **$match - Filtering Documents**

Like `WHERE` in SQL, `$match` filters documents early in the pipeline:

```go
// Filter by single field
{{"$match", bson.M{"status": "completed"}}}

// Filter by multiple conditions
{{"$match", bson.M{
    "status": "completed",
    "total": bson.M{"$gte": 100},
    "order_date": bson.M{
        "$gte": time.Now().AddDate(0, -1, 0), // Last month
    },
}}}

// Complex filters with $and, $or
{{"$match", bson.M{
    "$and": []bson.M{
        {"status": "completed"},
        {
            "$or": []bson.M{
                {"category": "Electronics"},
                {"total": bson.M{"$gte": 500}},
            },
        },
    },
}}}
```

### **$group - Grouping and Aggregating**

The heart of analytics - group documents and calculate metrics:

```go
// Basic grouping
{{"$group", bson.M{
    "_id":    "$category",           // Group by category
    "total":  bson.M{"$sum": "$amount"},    // Sum amounts
    "count":  bson.M{"$sum": 1},            // Count documents
    "avg":    bson.M{"$avg": "$amount"},    // Average amount
}}}

// Multiple grouping fields
{{"$group", bson.M{
    "_id": bson.M{
        "category": "$category",
        "year":     bson.M{"$year": "$order_date"},
    },
    "total_sales": bson.M{"$sum": "$total"},
}}}

// Advanced aggregation operators
{{"$group", bson.M{
    "_id": "$product_id",
    "total_revenue":  bson.M{"$sum": "$total"},
    "max_order":      bson.M{"$max": "$total"},
    "min_order":      bson.M{"$min": "$total"},
    "first_order":    bson.M{"$first": "$order_date"},
    "last_order":     bson.M{"$last": "$order_date"},
    "unique_customers": bson.M{"$addToSet": "$customer_id"},
    "all_categories":   bson.M{"$push": "$category"},
}}}
```

### **$project - Selecting and Transforming Fields**

Shape your output and create calculated fields:

```go
// Select specific fields
{{"$project", bson.M{
    "name":     1,  // Include field
    "email":    1,
    "password": 0,  // Exclude field (explicit)
}}}

// Create calculated fields
{{"$project", bson.M{
    "name":        1,
    "total":       1,
    "tax":         bson.M{"$multiply": []interface{}{"$total", 0.08}},
    "grand_total": bson.M{"$add": []interface{}{"$total", bson.M{"$multiply": []interface{}{"$total", 0.08}}}},
    "order_year":  bson.M{"$year": "$order_date"},
    "full_name":   bson.M{"$concat": []interface{}{"$first_name", " ", "$last_name"}},
}}}

// Conditional fields
{{"$project", bson.M{
    "name":   1,
    "total":  1,
    "status": bson.M{
        "$cond": bson.M{
            "if":   bson.M{"$gte": []interface{}{"$total", 1000}},
            "then": "VIP",
            "else": "Regular",
        },
    },
}}}
```

### **$sort - Ordering Results**

Sort documents by one or more fields:

```go
// Single field sort
{{"$sort", bson.M{"total": -1}}}  // Descending
{{"$sort", bson.M{"name": 1}}}    // Ascending

// Multi-field sort (order matters!)
{{"$sort", bson.M{
    "category": 1,   // Primary sort: category ascending
    "total": -1,     // Secondary sort: total descending
}}}

// Sort by calculated field
{{"$sort", bson.M{"calculated_field": -1}}}
```

## 📊 **Advanced Aggregation Techniques**

### **Date Operations**

Extract and manipulate date components:

```go
// Extract date parts
{{"$group", bson.M{
    "_id": bson.M{
        "year":  bson.M{"$year": "$order_date"},
        "month": bson.M{"$month": "$order_date"},
        "day":   bson.M{"$dayOfMonth": "$order_date"},
        "weekday": bson.M{"$dayOfWeek": "$order_date"},
    },
    "daily_sales": bson.M{"$sum": "$total"},
}}}

// Date filtering with expressions
{{"$match", bson.M{
    "$expr": bson.M{
        "$and": []bson.M{
            {"$gte": []interface{}{bson.M{"$year": "$order_date"}, 2024}},
            {"$eq": []interface{}{bson.M{"$month": "$order_date"}, 12}},
        },
    },
}}}

// Create date from parts
{{"$addFields", bson.M{
    "month_start": bson.M{
        "$dateFromParts": bson.M{
            "year":  bson.M{"$year": "$order_date"},
            "month": bson.M{"$month": "$order_date"},
            "day":   1,
        },
    },
}}}
```

### **$lookup - Joining Collections**

Combine data from multiple collections (like SQL JOINs):

```go
// Basic lookup
{{"$lookup", bson.M{
    "from":         "customers",      // Collection to join
    "localField":   "customer_id",    // Field in current collection
    "foreignField": "_id",            // Field in foreign collection
    "as":           "customer_info",  // Output array field
}}}

// Lookup with pipeline (advanced)
{{"$lookup", bson.M{
    "from": "products",
    "let":  bson.M{"product_id": "$product_id"},
    "pipeline": mongo.Pipeline{
        {{"$match", bson.M{
            "$expr": bson.M{"$eq": []interface{}{"$_id", "$$product_id"}},
        }}},
        {{"$project", bson.M{"name": 1, "category": 1, "price": 1}}},
    },
    "as": "product_details",
}}}

// Unwind joined data (convert array to object)
{{"$unwind", "$customer_info"}},  // Convert array to single object
{{"$unwind", bson.M{              // With options
    "path": "$customer_info",
    "preserveNullAndEmptyArrays": true,  // Keep docs with no matches
}}},
```

### **$addFields vs $project**

Add calculated fields without removing existing ones:

```go
// $addFields - adds fields, keeps existing ones
{{"$addFields", bson.M{
    "total_with_tax": bson.M{"$multiply": []interface{}{"$total", 1.08}},
    "order_year":     bson.M{"$year": "$order_date"},
    "customer_tier":  bson.M{
        "$switch": bson.M{
            "branches": []bson.M{
                {"case": bson.M{"$gte": []interface{}{"$total", 1000}}, "then": "Gold"},
                {"case": bson.M{"$gte": []interface{}{"$total", 500}}, "then": "Silver"},
            },
            "default": "Bronze",
        },
    },
}}}

// $project - explicitly choose fields (removes others)
{{"$project", bson.M{
    "customer_id":    1,
    "total":          1,
    "total_with_tax": bson.M{"$multiply": []interface{}{"$total", 1.08}},
    // All other fields are excluded
}}}
```

## 🏗️ **Real-World Analytics Examples**

### **E-commerce Sales Dashboard**

```go
// Complete sales analytics pipeline
func GetSalesDashboard(collection *mongo.Collection, ctx context.Context) (*SalesDashboard, error) {
    pipeline := mongo.Pipeline{
        // Filter recent completed orders
        {{"$match", bson.M{
            "status": "completed",
            "order_date": bson.M{"$gte": time.Now().AddDate(0, -3, 0)}, // Last 3 months
        }}},
        
        // Add calculated fields
        {{"$addFields", bson.M{
            "order_month": bson.M{
                "$dateToString": bson.M{
                    "format": "%Y-%m",
                    "date":   "$order_date",
                },
            },
        }}},
        
        // Create multiple analytics in one pipeline using $facet
        {{"$facet", bson.M{
            // Sales by category
            "category_sales": mongo.Pipeline{
                {{"$group", bson.M{
                    "_id":         "$category",
                    "total_sales": bson.M{"$sum": "$total"},
                    "order_count": bson.M{"$sum": 1},
                }}},
                {{"$sort", bson.M{"total_sales": -1}}},
            },
            
            // Monthly trends
            "monthly_trends": mongo.Pipeline{
                {{"$group", bson.M{
                    "_id":           "$order_month",
                    "monthly_sales": bson.M{"$sum": "$total"},
                    "order_count":   bson.M{"$sum": 1},
                }}},
                {{"$sort", bson.M{"_id": 1}}},
            },
            
            // Top products
            "top_products": mongo.Pipeline{
                {{"$group", bson.M{
                    "_id":            "$product_id",
                    "total_quantity": bson.M{"$sum": "$quantity"},
                    "total_revenue":  bson.M{"$sum": "$total"},
                }}},
                {{"$sort", bson.M{"total_revenue": -1}}},
                {{"$limit", 10}},
            },
        }}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    // ... handle results
}
```

### **Customer Behavior Analysis**

```go
// Analyze customer purchase patterns
func AnalyzeCustomerBehavior(ordersCollection, customersCollection *mongo.Collection, ctx context.Context) error {
    pipeline := mongo.Pipeline{
        // Join orders with customer data
        {{"$lookup", bson.M{
            "from":         "customers",
            "localField":   "customer_id",
            "foreignField": "_id",
            "as":           "customer",
        }}},
        
        {{"$unwind", "$customer"}},
        
        // Filter completed orders
        {{"$match", bson.M{"status": "completed"}}},
        
        // Group by customer and calculate behavior metrics
        {{"$group", bson.M{
            "_id": "$customer_id",
            "customer_name":     bson.M{"$first": "$customer.name"},
            "customer_location": bson.M{"$first": "$customer.location"},
            "total_spent":       bson.M{"$sum": "$total"},
            "order_count":       bson.M{"$sum": 1},
            "avg_order_value":   bson.M{"$avg": "$total"},
            "first_order":       bson.M{"$min": "$order_date"},
            "last_order":        bson.M{"$max": "$order_date"},
            "favorite_categories": bson.M{"$addToSet": "$category"},
            "order_frequency":   bson.M{
                "$avg": bson.M{
                    "$divide": []interface{}{
                        bson.M{"$subtract": []interface{}{"$order_date", "$first_order"}},
                        86400000, // Convert to days
                    },
                },
            },
        }}},
        
        // Calculate customer lifetime value and segment
        {{"$addFields", bson.M{
            "days_since_first_order": bson.M{
                "$divide": []interface{}{
                    bson.M{"$subtract": []interface{}{time.Now(), "$first_order"}},
                    86400000,
                },
            },
            "customer_segment": bson.M{
                "$switch": bson.M{
                    "branches": []bson.M{
                        {
                            "case": bson.M{"$and": []bson.M{
                                {"$gte": []interface{}{"$total_spent", 5000}},
                                {"$gte": []interface{}{"$order_count", 10}},
                            }},
                            "then": "VIP",
                        },
                        {
                            "case": bson.M{"$gte": []interface{}{"$total_spent", 1000}},
                            "then": "Premium",
                        },
                    },
                    "default": "Regular",
                },
            },
        }}},
        
        // Sort by total spent
        {{"$sort", bson.M{"total_spent": -1}}},
    }
    
    cursor, err := ordersCollection.Aggregate(ctx, pipeline)
    // ... process results
}
```

## 🚀 **Performance Optimization**

### **Index Strategy for Aggregation**

```go
// Create indexes to support aggregation stages
func CreateAggregationIndexes(collection *mongo.Collection, ctx context.Context) error {
    indexes := []mongo.IndexModel{
        // Support $match stages
        {Keys: bson.M{"status": 1, "order_date": -1}},
        {Keys: bson.M{"category": 1, "status": 1}},
        
        // Support $group stages
        {Keys: bson.M{"customer_id": 1, "order_date": -1}},
        {Keys: bson.M{"product_id": 1, "status": 1}},
        
        // Support $sort stages
        {Keys: bson.M{"total": -1}},
        {Keys: bson.M{"order_date": -1, "total": -1}},
        
        // Compound indexes for complex queries
        {Keys: bson.M{"status": 1, "category": 1, "order_date": -1}},
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

### **Pipeline Optimization Tips**

```go
// 1. Filter early with $match
pipeline := mongo.Pipeline{
    // Put $match as early as possible
    {{"$match", bson.M{"status": "completed"}}},
    
    // Then do expensive operations
    {{"$lookup", bson.M{...}}},
    {{"$group", bson.M{...}}},
}

// 2. Use $project to reduce data size
{{"$project", bson.M{
    "customer_id": 1,
    "total":       1,
    "order_date":  1,
    // Only include fields you need for subsequent stages
}}}

// 3. Use $limit early when possible
{{"$limit", 10000}}, // Limit before expensive operations

// 4. Use $sample for large datasets
{{"$sample", bson.M{"size": 1000}}}, // Random sample

// 5. Use allowDiskUse for large aggregations
opts := options.Aggregate().SetAllowDiskUse(true)
cursor, err := collection.Aggregate(ctx, pipeline, opts)
```

## 📈 **Common Aggregation Patterns**

### **Top N Analysis**

```go
// Find top N items in each category
pipeline := mongo.Pipeline{
    {{"$match", bson.M{"status": "completed"}}},
    
    // Sort within each category
    {{"$sort", bson.M{"category": 1, "total": -1}}},
    
    // Group and collect top items
    {{"$group", bson.M{
        "_id": "$category",
        "top_orders": bson.M{
            "$push": bson.M{
                "order_id": "$_id",
                "total":    "$total",
                "customer": "$customer_id",
            },
        },
    }}},
    
    // Limit to top 5 per category
    {{"$addFields", bson.M{
        "top_orders": bson.M{"$slice": []interface{}{"$top_orders", 5}},
    }}},
}
```

### **Moving Averages**

```go
// Calculate 7-day moving average
pipeline := mongo.Pipeline{
    {{"$match", bson.M{"order_date": bson.M{"$gte": startDate}}}},
    
    // Group by date
    {{"$group", bson.M{
        "_id":          bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$order_date"}},
        "daily_sales":  bson.M{"$sum": "$total"},
        "daily_orders": bson.M{"$sum": 1},
    }}},
    
    {{"$sort", bson.M{"_id": 1}}},
    
    // Calculate moving average using $setWindowFields (MongoDB 5.0+)
    {{"$setWindowFields", bson.M{
        "sortBy": bson.M{"_id": 1},
        "output": bson.M{
            "moving_avg_sales": bson.M{
                "$avg": "$daily_sales",
                "window": bson.M{
                    "range": []int{-6, 0}, // 7-day window (6 days back + current)
                    "unit":  "position",
                },
            },
        },
    }}},
}
```

## 🔗 **Integration with Go Applications**

### **Structured Result Handling**

```go
type AnalyticsResult struct {
    CategorySales   []SalesByCategory   `json:"category_sales"`
    MonthlyTrends   []MonthlyRevenue    `json:"monthly_trends"`
    TopProducts     []TopProduct        `json:"top_products"`
    CustomerMetrics []CustomerAnalytics `json:"customer_metrics"`
}

func RunCompleteAnalytics(collection *mongo.Collection, ctx context.Context) (*AnalyticsResult, error) {
    // Use $facet to run multiple analytics in parallel
    pipeline := mongo.Pipeline{
        {{"$match", bson.M{"status": "completed"}}},
        
        {{"$facet", bson.M{
            "category_sales": mongo.Pipeline{
                // Category analysis pipeline
            },
            "monthly_trends": mongo.Pipeline{
                // Monthly trends pipeline  
            },
            "top_products": mongo.Pipeline{
                // Top products pipeline
            },
            "customer_metrics": mongo.Pipeline{
                // Customer analytics pipeline
            },
        }}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var results []AnalyticsResult
    if err = cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    
    if len(results) == 0 {
        return &AnalyticsResult{}, nil
    }
    
    return &results[0], nil
}
```

## 🌍 **Real-World Applications**

### **Business Intelligence Dashboards**
- Sales performance by region, product, time period
- Customer segmentation and lifetime value analysis
- Inventory turnover and demand forecasting
- Marketing campaign effectiveness measurement

### **Financial Analytics**
- Revenue recognition and reporting
- Profit margin analysis by product/customer
- Cash flow projections and trend analysis
- Risk assessment and fraud detection

### **Operational Analytics**
- Supply chain optimization
- Quality metrics and defect analysis
- Performance KPIs and SLA monitoring
- Resource utilization and capacity planning

## 📚 **Best Practices**

1. **Design for Performance**: Create indexes before building pipelines
2. **Filter Early**: Use `$match` as the first stage when possible
3. **Limit Data**: Use `$project` and `$limit` to reduce processing
4. **Test Incrementally**: Build pipelines stage by stage
5. **Monitor Performance**: Use `explain()` to analyze pipeline performance
6. **Handle Errors**: Always validate inputs and handle aggregation failures
7. **Document Pipelines**: Complex aggregations need clear documentation

## 🚀 **Next Steps**

After mastering aggregation pipelines:

1. **Challenge 4**: Learn indexing strategies for optimal aggregation performance
2. **Challenge 5**: Implement transactions for consistent multi-collection analytics
3. **Advanced Topics**: Explore time-series collections, search indexes, and Atlas Search

## 🔗 **Additional Resources**

- [MongoDB Aggregation Pipeline](https://docs.mongodb.com/manual/core/aggregation-pipeline/)
- [Aggregation Pipeline Operators](https://docs.mongodb.com/manual/reference/operator/aggregation/)
- [Aggregation Performance](https://docs.mongodb.com/manual/core/aggregation-pipeline-optimization/)
- [Real-Time Analytics with MongoDB](https://www.mongodb.com/solutions/real-time-analytics)

Ready to transform raw data into powerful insights? Master the aggregation pipeline! 📊🚀
