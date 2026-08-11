# Hints for Challenge 2: Advanced Queries & Filtering

## Hint 1: Understanding Query Operators

MongoDB provides powerful query operators for complex filtering:

```go
// Comparison operators
filter := bson.M{
    "price": bson.M{
        "$gte": 100,  // Greater than or equal
        "$lte": 500,  // Less than or equal
        "$gt":  50,   // Greater than
        "$lt":  1000, // Less than
    },
}

// Array operators
filter := bson.M{
    "tags": bson.M{
        "$in":  []string{"smartphone", "premium"}, // Any of these values
        "$all": []string{"running", "comfort"},    // All of these values
    },
}
```

## Hint 2: Building Complex Filters

Combine multiple conditions for advanced filtering:

```go
func (ps *ProductService) GetProductsByCategory(ctx context.Context, category string, filter ProductFilter) Response {
    // Start with category filter
    queryFilter := bson.M{"category": category}
    
    // Add price range if specified
    if filter.MinPrice > 0 || filter.MaxPrice > 0 {
        priceFilter := bson.M{}
        if filter.MinPrice > 0 {
            priceFilter["$gte"] = filter.MinPrice
        }
        if filter.MaxPrice > 0 {
            priceFilter["$lte"] = filter.MaxPrice
        }
        queryFilter["price"] = priceFilter
    }
    
    // Add other filters
    if filter.Brand != "" {
        queryFilter["brand"] = filter.Brand
    }
    
    if filter.InStock {
        queryFilter["stock"] = bson.M{"$gt": 0}
    }
    
    // Execute query
    cursor, err := ps.Collection.Find(ctx, queryFilter)
    // ... handle results
}
```

## Hint 3: Text Search with Regex

Use regex for flexible text searching:

```go
func (ps *ProductService) SearchProductsByName(ctx context.Context, searchTerm string, caseSensitive bool) Response {
    // Create regex pattern
    regexOptions := ""
    if !caseSensitive {
        regexOptions = "i" // Case insensitive
    }
    
    regex := primitive.Regex{
        Pattern: searchTerm,
        Options: regexOptions,
    }
    
    // Search in multiple fields using $or
    filter := bson.M{
        "$or": []bson.M{
            {"name": bson.M{"$regex": regex}},
            {"description": bson.M{"$regex": regex}},
        },
    }
    
    cursor, err := ps.Collection.Find(ctx, filter)
    // ... handle results
}
```

## Hint 4: Implementing Pagination

Pagination requires counting total documents and using skip/limit:

```go
func (ps *ProductService) GetProductsWithPagination(ctx context.Context, pagination PaginationOptions, filter ProductFilter) Response {
    // Validate pagination parameters
    if pagination.Page < 1 {
        pagination.Page = 1
    }
    if pagination.Limit < 1 {
        pagination.Limit = 10
    }
    
    // Build filter from ProductFilter
    queryFilter := bson.M{}
    if filter.Category != "" {
        queryFilter["category"] = filter.Category
    }
    // ... add other filters
    
    // Count total documents
    total, err := ps.Collection.CountDocuments(ctx, queryFilter)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    // Calculate pagination values
    skip := (pagination.Page - 1) * pagination.Limit
    totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
    
    // Execute paginated query
    opts := options.Find().
        SetSkip(int64(skip)).
        SetLimit(int64(pagination.Limit)).
        SetSort(bson.M{"created_at": -1}) // Sort by newest first
    
    cursor, err := ps.Collection.Find(ctx, queryFilter, opts)
    // ... decode results and return PaginatedResponse
}
```

## Hint 5: Sorting with Multiple Criteria

Use bson.D for ordered sort criteria:

```go
func (ps *ProductService) GetProductsSorted(ctx context.Context, sortOptions []SortOptions, limit int) Response {
    // Build sort document - order matters!
    sortDoc := bson.D{}
    
    validFields := map[string]bool{
        "price": true, "rating": true, "name": true, "created_at": true,
    }
    
    for _, sortOpt := range sortOptions {
        // Validate field names
        if !validFields[sortOpt.Field] {
            return Response{Success: false, Error: "Invalid sort field", Code: 400}
        }
        
        // Normalize order (-1 or 1)
        order := 1
        if sortOpt.Order < 0 {
            order = -1
        }
        
        sortDoc = append(sortDoc, bson.E{Key: sortOpt.Field, Value: order})
    }
    
    // Apply sorting and limit
    opts := options.Find().SetSort(sortDoc)
    if limit > 0 {
        opts.SetLimit(int64(limit))
    }
    
    cursor, err := ps.Collection.Find(ctx, bson.M{}, opts)
    // ... handle results
}
```

## Hint 6: Field Projection

Optimize queries by selecting only needed fields:

```go
func (ps *ProductService) FilterProducts(ctx context.Context, filter ProductFilter, projection []string) Response {
    // Build complex filter with $and
    conditions := []bson.M{}
    
    if filter.Category != "" {
        conditions = append(conditions, bson.M{"category": filter.Category})
    }
    
    if filter.MinPrice > 0 || filter.MaxPrice > 0 {
        priceFilter := bson.M{}
        if filter.MinPrice > 0 {
            priceFilter["$gte"] = filter.MinPrice
        }
        if filter.MaxPrice > 0 {
            priceFilter["$lte"] = filter.MaxPrice
        }
        conditions = append(conditions, bson.M{"price": priceFilter})
    }
    
    // Combine conditions
    var queryFilter bson.M
    if len(conditions) > 0 {
        queryFilter = bson.M{"$and": conditions}
    } else {
        queryFilter = bson.M{}
    }
    
    // Build projection
    opts := options.Find()
    if len(projection) > 0 {
        projectionDoc := bson.M{}
        for _, field := range projection {
            projectionDoc[field] = 1
        }
        opts.SetProjection(projectionDoc)
    }
    
    cursor, err := ps.Collection.Find(ctx, queryFilter, opts)
    // ... handle results
}
```

## Hint 7: Tag-Based Filtering

Handle array field queries with $in and $all:

```go
func (ps *ProductService) GetProductsByTags(ctx context.Context, tags []string, matchAll bool) Response {
    if len(tags) == 0 {
        return Response{Success: false, Error: "At least one tag required", Code: 400}
    }
    
    // Clean up tags
    cleanTags := []string{}
    for _, tag := range tags {
        if strings.TrimSpace(tag) != "" {
            cleanTags = append(cleanTags, strings.TrimSpace(tag))
        }
    }
    
    var filter bson.M
    if matchAll {
        // Must have ALL tags
        filter = bson.M{"tags": bson.M{"$all": cleanTags}}
    } else {
        // Must have ANY tag
        filter = bson.M{"tags": bson.M{"$in": cleanTags}}
    }
    
    cursor, err := ps.Collection.Find(ctx, filter)
    // ... handle results
}
```

## Hint 8: Top-Rated Products Query

Combine filtering, sorting, and limiting:

```go
func (ps *ProductService) GetTopRatedProducts(ctx context.Context, category string, limit int) Response {
    // Build filter
    filter := bson.M{"rating": bson.M{"$gt": 0}} // Only rated products
    if category != "" {
        filter["category"] = category
    }
    
    // Sort by rating descending, then by name for consistency
    opts := options.Find().
        SetSort(bson.M{"rating": -1, "name": 1}).
        SetLimit(int64(limit))
    
    cursor, err := ps.Collection.Find(ctx, filter, opts)
    // ... handle results
}
```

## Hint 9: Counting and Aggregation

Use CountDocuments and Distinct for analytics:

```go
func (ps *ProductService) CountProductsByCategory(ctx context.Context) Response {
    // Get all distinct categories
    categories, err := ps.Collection.Distinct(ctx, "category", bson.M{})
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    // Count products in each category
    categoryCounts := make(map[string]int64)
    for _, cat := range categories {
        if category, ok := cat.(string); ok {
            count, err := ps.Collection.CountDocuments(ctx, bson.M{"category": category})
            if err != nil {
                return Response{Success: false, Error: err.Error(), Code: 500}
            }
            categoryCounts[category] = count
        }
    }
    
    return Response{
        Success: true,
        Data:    categoryCounts,
        Message: "Category counts retrieved",
        Code:    200,
    }
}
```

## Hint 10: Error Handling and Validation

Always validate inputs and handle errors properly:

```go
func (ps *ProductService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) Response {
    // Validate input
    if minPrice < 0 {
        return Response{Success: false, Error: "Minimum price cannot be negative", Code: 400}
    }
    if maxPrice < 0 {
        return Response{Success: false, Error: "Maximum price cannot be negative", Code: 400}
    }
    if minPrice > maxPrice {
        return Response{Success: false, Error: "Min price cannot exceed max price", Code: 400}
    }
    
    filter := bson.M{
        "price": bson.M{
            "$gte": minPrice,
            "$lte": maxPrice,
        },
    }
    
    // Sort by price for better UX
    opts := options.Find().SetSort(bson.M{"price": 1})
    
    cursor, err := ps.Collection.Find(ctx, filter, opts)
    if err != nil {
        return Response{Success: false, Error: "Query failed: " + err.Error(), Code: 500}
    }
    defer cursor.Close(ctx)
    
    var products []Product
    if err = cursor.All(ctx, &products); err != nil {
        return Response{Success: false, Error: "Failed to decode: " + err.Error(), Code: 500}
    }
    
    if products == nil {
        products = []Product{} // Return empty array, not nil
    }
    
    return Response{
        Success: true,
        Data:    products,
        Message: fmt.Sprintf("Found %d products in price range", len(products)),
        Code:    200,
    }
}
```

Ready to implement advanced MongoDB queries? Start with simple filters and build up to complex combinations! 🚀
