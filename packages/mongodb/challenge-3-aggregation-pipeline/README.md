# Challenge 3: Aggregation Pipeline

Build an **Analytics Dashboard** using MongoDB's powerful aggregation pipeline to analyze e-commerce order data and generate business insights.

## Challenge Requirements

Implement an `AnalyticsService` that uses MongoDB aggregation pipelines to perform complex data analysis. Your service should provide the following analytics:

- **Sales by Category**: Calculate total sales, order count, and average order value per product category.
- **Top Selling Products**: Find the best-performing products by quantity sold and revenue generated.
- **Monthly Revenue**: Analyze revenue trends over time, grouped by year and month.
- **Customer Analytics**: Identify top customers by total spending and purchase behavior.
- **Product Performance**: Combine order and customer data using `$lookup` for comprehensive analysis.
- **Order Trends**: Track daily order patterns and growth rates.

## Data Structures

```go
type Order struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
    ProductID  primitive.ObjectID `bson:"product_id" json:"product_id"`
    Quantity   int                `bson:"quantity" json:"quantity"`
    Price      float64            `bson:"price" json:"price"`
    Total      float64            `bson:"total" json:"total"`
    Category   string             `bson:"category" json:"category"`
    OrderDate  time.Time          `bson:"order_date" json:"order_date"`
    Status     string             `bson:"status" json:"status"`
}

type Customer struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name     string             `bson:"name" json:"name"`
    Email    string             `bson:"email" json:"email"`
    Location string             `bson:"location" json:"location"`
    JoinDate time.Time          `bson:"join_date" json:"join_date"`
}

type SalesByCategory struct {
    Category      string  `bson:"_id" json:"category"`
    TotalSales    float64 `bson:"total_sales" json:"total_sales"`
    OrderCount    int     `bson:"order_count" json:"order_count"`
    AvgOrderValue float64 `bson:"avg_order_value" json:"avg_order_value"`
}
```

## Example Aggregation Pipeline

**Sales by Category Analysis:**
```go
pipeline := mongo.Pipeline{
    // Stage 1: Filter completed orders only
    {{"$match", bson.M{"status": "completed"}}},
    
    // Stage 2: Group by category and calculate metrics
    {{"$group", bson.M{
        "_id":         "$category",
        "total_sales": bson.M{"$sum": "$total"},
        "order_count": bson.M{"$sum": 1},
    }}},
    
    // Stage 3: Calculate average order value
    {{"$addFields", bson.M{
        "avg_order_value": bson.M{"$divide": []interface{}{"$total_sales", "$order_count"}},
    }}},
    
    // Stage 4: Sort by total sales descending
    {{"$sort", bson.M{"total_sales": -1}}},
}
```

## Response Examples

**Sales by Category:**
```json
{
    "success": true,
    "data": [
        {
            "category": "Electronics",
            "total_sales": 15999.95,
            "order_count": 12,
            "avg_order_value": 1333.33
        },
        {
            "category": "Footwear", 
            "total_sales": 2599.88,
            "order_count": 8,
            "avg_order_value": 324.99
        }
    ],
    "message": "Sales analytics by category retrieved",
    "code": 200
}
```

**Top Selling Products:**
```json
{
    "success": true,
    "data": [
        {
            "product_id": "65b23c2e01d2a3b4c5d6e7f8",
            "total_quantity": 45,
            "total_revenue": 44999.55,
            "order_count": 15
        }
    ],
    "message": "Top selling products retrieved",
    "code": 200
}
```

## Testing Requirements

Your solution must pass tests for:
- Aggregating sales data by category with proper calculations
- Finding top-selling products sorted by quantity and revenue
- Generating monthly revenue reports with time-based grouping
- Analyzing customer behavior and spending patterns
- Handling edge cases like empty datasets and invalid parameters
- Proper error handling for aggregation failures
- Consistent response structure for all analytics endpoints

## Key Aggregation Concepts

- **$match**: Filter documents (like WHERE in SQL)
- **$group**: Group documents and perform calculations
- **$sort**: Sort results by specified fields
- **$limit**: Limit the number of results
- **$project**: Select and transform fields
- **$addFields**: Add calculated fields
- **$lookup**: Join collections (like JOIN in SQL)
- **$unwind**: Deconstruct array fields
