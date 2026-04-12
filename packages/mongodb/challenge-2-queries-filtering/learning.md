# Learning: Advanced MongoDB Queries & Filtering

## 🎯 **What You'll Learn**

This challenge teaches advanced MongoDB query techniques essential for building sophisticated search and filtering systems. You'll master the query operators, sorting, pagination, and optimization techniques used by major e-commerce platforms.

## 🔍 **MongoDB Query Operators**

### **Comparison Operators**

MongoDB provides rich comparison operators for numeric and string fields:

```go
// Price range filtering
filter := bson.M{
    "price": bson.M{
        "$gte": 100,    // Greater than or equal to 100
        "$lte": 500,    // Less than or equal to 500
    },
}

// Age restrictions
filter := bson.M{
    "age": bson.M{
        "$gt": 18,      // Greater than 18
        "$lt": 65,      // Less than 65
    },
}

// Exact matches and exclusions
filter := bson.M{
    "status": bson.M{"$in": []string{"active", "pending"}},     // In array
    "type":   bson.M{"$nin": []string{"deleted", "archived"}}, // Not in array
    "rating": bson.M{"$ne": 0},                                // Not equal
}
```

### **Logical Operators**

Combine multiple conditions with logical operators:

```go
// AND conditions (default behavior)
filter := bson.M{
    "category": "Electronics",
    "price":    bson.M{"$lt": 1000},
    "stock":    bson.M{"$gt": 0},
}

// OR conditions
filter := bson.M{
    "$or": []bson.M{
        {"category": "Electronics"},
        {"category": "Computers"},
    },
}

// Complex combinations
filter := bson.M{
    "$and": []bson.M{
        {"price": bson.M{"$gte": 100}},
        {
            "$or": []bson.M{
                {"brand": "Apple"},
                {"rating": bson.M{"$gte": 4.5}},
            },
        },
    },
}
```

### **Array Operators**

Work with array fields using specialized operators:

```go
// Products with any of these tags
filter := bson.M{
    "tags": bson.M{"$in": []string{"smartphone", "premium", "5g"}},
}

// Products with ALL of these tags
filter := bson.M{
    "tags": bson.M{"$all": []string{"waterproof", "wireless"}},
}

// Array size
filter := bson.M{
    "tags": bson.M{"$size": 3}, // Exactly 3 tags
}

// Element match for complex arrays
filter := bson.M{
    "reviews": bson.M{
        "$elemMatch": bson.M{
            "rating": bson.M{"$gte": 4},
            "verified": true,
        },
    },
}
```

## 🔤 **Text Search and Pattern Matching**

### **Regex Search**

Use regular expressions for flexible text matching:

```go
// Case-insensitive search
regex := primitive.Regex{
    Pattern: "iphone",
    Options: "i", // Case insensitive
}

filter := bson.M{
    "name": bson.M{"$regex": regex},
}

// Search multiple fields
filter := bson.M{
    "$or": []bson.M{
        {"name": bson.M{"$regex": regex}},
        {"description": bson.M{"$regex": regex}},
        {"brand": bson.M{"$regex": regex}},
    },
}

// Pattern matching
phoneRegex := primitive.Regex{
    Pattern: "^\\+?[1-9]\\d{1,14}$", // Phone number pattern
    Options: "",
}
```

### **Text Indexes (Advanced)**

For production applications, use MongoDB text indexes:

```go
// Create text index (usually done once)
indexModel := mongo.IndexModel{
    Keys: bson.M{
        "name":        "text",
        "description": "text",
        "brand":       "text",
    },
    Options: options.Index().SetWeights(bson.M{
        "name":  10, // Higher weight for name matches
        "brand": 5,
        "description": 1,
    }),
}

// Use text search
filter := bson.M{
    "$text": bson.M{
        "$search": "smartphone apple",
    },
}

// Sort by text score
opts := options.Find().SetSort(bson.M{
    "score": bson.M{"$meta": "textScore"},
})
```

## 📄 **Sorting and Pagination**

### **Single Field Sorting**

```go
// Sort by price ascending
opts := options.Find().SetSort(bson.M{"price": 1})

// Sort by rating descending
opts := options.Find().SetSort(bson.M{"rating": -1})

// Sort by creation date (newest first)
opts := options.Find().SetSort(bson.M{"created_at": -1})
```

### **Multi-Field Sorting**

Use `bson.D` for ordered sort criteria:

```go
// Sort by category, then by price within category
sortDoc := bson.D{
    {"category", 1},  // Primary sort: category ascending
    {"price", -1},    // Secondary sort: price descending
}

opts := options.Find().SetSort(sortDoc)
```

### **Pagination Implementation**

```go
func PaginateProducts(collection *mongo.Collection, page, limit int) (*PaginatedResponse, error) {
    // Validate parameters
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 20 // Default limit
    }
    
    // Count total documents
    total, err := collection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    
    // Calculate pagination
    skip := (page - 1) * limit
    totalPages := int(math.Ceil(float64(total) / float64(limit)))
    
    // Execute paginated query
    opts := options.Find().
        SetSkip(int64(skip)).
        SetLimit(int64(limit)).
        SetSort(bson.M{"created_at": -1})
    
    cursor, err := collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var products []Product
    if err = cursor.All(ctx, &products); err != nil {
        return nil, err
    }
    
    return &PaginatedResponse{
        Data:       products,
        Total:      total,
        Page:       page,
        Limit:      limit,
        TotalPages: totalPages,
    }, nil
}
```

## 🎯 **Field Projection**

Optimize queries by selecting only needed fields:

```go
// Include specific fields
projection := bson.M{
    "name":     1,
    "price":    1,
    "rating":   1,
    "_id":      1, // Always included unless explicitly excluded
}

// Exclude specific fields
projection := bson.M{
    "description": 0,
    "internal_notes": 0,
    "cost": 0,
}

// Use projection in queries
opts := options.Find().SetProjection(projection)
cursor, err := collection.Find(ctx, filter, opts)
```

## 📊 **Aggregation Basics**

While Challenge 3 covers aggregation in detail, here are basics:

```go
// Count documents by category
pipeline := mongo.Pipeline{
    {{"$group", bson.M{
        "_id":   "$category",
        "count": bson.M{"$sum": 1},
    }}},
    {{"$sort", bson.M{"count": -1}}},
}

cursor, err := collection.Aggregate(ctx, pipeline)
```

## 🔧 **Query Optimization**

### **Index Usage**

```go
// Create indexes for frequently queried fields
indexes := []mongo.IndexModel{
    {Keys: bson.M{"category": 1}},
    {Keys: bson.M{"price": 1}},
    {Keys: bson.M{"rating": -1}},
    {Keys: bson.M{"created_at": -1}},
    
    // Compound indexes for common query patterns
    {Keys: bson.M{"category": 1, "price": 1}},
    {Keys: bson.M{"brand": 1, "rating": -1}},
}

collection.Indexes().CreateMany(ctx, indexes)
```

### **Query Performance Tips**

1. **Use Indexes**: Create indexes for frequently queried fields
2. **Limit Results**: Always use limits for large datasets
3. **Project Fields**: Only select needed fields
4. **Avoid Regex**: Use text indexes instead of regex when possible
5. **Compound Indexes**: Create indexes matching your query patterns

```go
// Good: Uses index on category, limits results
filter := bson.M{"category": "Electronics"}
opts := options.Find().SetLimit(20)

// Bad: No index, no limit, selects all fields
filter := bson.M{"description": bson.M{"$regex": primitive.Regex{Pattern: ".*phone.*", Options: "i"}}}
```

## 🏗️ **Building Complex Filters**

### **Dynamic Filter Building**

```go
func BuildProductFilter(params ProductFilter) bson.M {
    filter := bson.M{}
    
    // Category filter
    if params.Category != "" {
        filter["category"] = params.Category
    }
    
    // Price range
    if params.MinPrice > 0 || params.MaxPrice > 0 {
        priceFilter := bson.M{}
        if params.MinPrice > 0 {
            priceFilter["$gte"] = params.MinPrice
        }
        if params.MaxPrice > 0 {
            priceFilter["$lte"] = params.MaxPrice
        }
        filter["price"] = priceFilter
    }
    
    // Brand filter
    if params.Brand != "" {
        filter["brand"] = params.Brand
    }
    
    // Rating filter
    if params.MinRating > 0 {
        filter["rating"] = bson.M{"$gte": params.MinRating}
    }
    
    // Stock filter
    if params.InStock {
        filter["stock"] = bson.M{"$gt": 0}
    }
    
    // Tags filter
    if len(params.Tags) > 0 {
        filter["tags"] = bson.M{"$in": params.Tags}
    }
    
    // Text search
    if params.SearchTerm != "" {
        regex := primitive.Regex{
            Pattern: params.SearchTerm,
            Options: "i",
        }
        filter["$or"] = []bson.M{
            {"name": bson.M{"$regex": regex}},
            {"description": bson.M{"$regex": regex}},
        }
    }
    
    return filter
}
```

## 🌍 **Real-World Applications**

### **E-commerce Product Search**

```go
// Amazon-style product search
func SearchProducts(searchTerm string, filters ProductFilter, sort SortOptions, page int) (*SearchResults, error) {
    // Build complex filter
    filter := bson.M{}
    
    // Text search across multiple fields
    if searchTerm != "" {
        filter["$text"] = bson.M{"$search": searchTerm}
    }
    
    // Apply filters
    if filters.Category != "" {
        filter["category"] = filters.Category
    }
    
    if filters.MinPrice > 0 || filters.MaxPrice > 0 {
        priceFilter := bson.M{}
        if filters.MinPrice > 0 {
            priceFilter["$gte"] = filters.MinPrice
        }
        if filters.MaxPrice > 0 {
            priceFilter["$lte"] = filters.MaxPrice
        }
        filter["price"] = priceFilter
    }
    
    // Build sort
    sortDoc := bson.M{}
    if searchTerm != "" {
        sortDoc["score"] = bson.M{"$meta": "textScore"} // Relevance first
    }
    
    switch sort.Field {
    case "price":
        sortDoc["price"] = sort.Order
    case "rating":
        sortDoc["rating"] = sort.Order
    case "popularity":
        sortDoc["sales_count"] = sort.Order
    default:
        sortDoc["created_at"] = -1 // Newest first
    }
    
    // Execute with pagination
    opts := options.Find().
        SetSort(sortDoc).
        SetSkip(int64((page-1)*20)).
        SetLimit(20)
    
    cursor, err := collection.Find(ctx, filter, opts)
    // ... handle results
}
```

### **Faceted Search**

```go
// Get filter counts (like Amazon's left sidebar)
func GetSearchFacets(baseFilter bson.M) (*SearchFacets, error) {
    facets := &SearchFacets{}
    
    // Category counts
    pipeline := mongo.Pipeline{
        {{"$match", baseFilter}},
        {{"$group", bson.M{
            "_id":   "$category",
            "count": bson.M{"$sum": 1},
        }}},
        {{"$sort", bson.M{"count": -1}}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    // ... decode category counts
    
    // Brand counts
    pipeline = mongo.Pipeline{
        {{"$match", baseFilter}},
        {{"$group", bson.M{
            "_id":   "$brand",
            "count": bson.M{"$sum": 1},
        }}},
        {{"$sort", bson.M{"count": -1}}},
        {{"$limit", 10}}, // Top 10 brands
    }
    
    cursor, err = collection.Aggregate(ctx, pipeline)
    // ... decode brand counts
    
    return facets, nil
}
```

## 📚 **Best Practices**

### **Query Design**
1. **Index First**: Design indexes before writing queries
2. **Limit Always**: Never query without limits in production
3. **Project Wisely**: Only select fields you need
4. **Cache Results**: Cache expensive queries when possible

### **Error Handling**
```go
func HandleQueryError(err error) Response {
    if err == nil {
        return Response{Success: true}
    }
    
    // Log error for debugging
    log.Printf("Query error: %v", err)
    
    // Return user-friendly error
    return Response{
        Success: false,
        Error:   "Search temporarily unavailable",
        Code:    500,
    }
}
```

### **Performance Monitoring**
```go
// Monitor query performance
func MonitoredFind(collection *mongo.Collection, filter bson.M) (*mongo.Cursor, error) {
    start := time.Now()
    
    cursor, err := collection.Find(ctx, filter)
    
    duration := time.Since(start)
    if duration > 100*time.Millisecond {
        log.Printf("Slow query detected: %v took %v", filter, duration)
    }
    
    return cursor, err
}
```

## 🚀 **Next Steps**

After mastering queries and filtering:

1. **Challenge 3**: Learn aggregation pipelines for complex analytics
2. **Challenge 4**: Master indexing and performance optimization
3. **Challenge 5**: Implement transactions for data consistency

## 🔗 **Additional Resources**

- [MongoDB Query Operators](https://docs.mongodb.com/manual/reference/operator/query/)
- [Query Optimization](https://docs.mongodb.com/manual/core/query-optimization/)
- [Text Search](https://docs.mongodb.com/manual/text-search/)
- [Pagination Best Practices](https://docs.mongodb.com/manual/reference/method/cursor.skip/#pagination-example)

Ready to build sophisticated search systems like the pros? 🎯
