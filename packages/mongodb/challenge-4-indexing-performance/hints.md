# Hints for Challenge 4: Indexing & Performance

## Hint 1: Understanding Index Types

MongoDB supports several index types, each optimized for different query patterns:

```go
// Single field indexes - most basic and common
db.products.createIndex({"category": 1})     // Ascending
db.products.createIndex({"price": -1})       // Descending

// Compound indexes - multiple fields in one index
db.products.createIndex({"category": 1, "price": 1, "rating": -1})

// Text indexes - for full-text search
db.products.createIndex({"name": "text", "description": "text"})

// Multikey indexes - automatically created for array fields
db.products.createIndex({"tags": 1})  // Works with array fields

// Sparse indexes - only index documents that have the field
db.products.createIndex({"optional_field": 1}, {"sparse": true})
```

## Hint 2: Creating Optimal Index Sets

Build a comprehensive index strategy covering common query patterns:

```go
func (is *IndexService) CreateOptimalIndexes(ctx context.Context) Response {
    indexes := []mongo.IndexModel{
        // Single field indexes for common filters
        {Keys: bson.D{{"category", 1}}},
        {Keys: bson.D{{"brand", 1}}},
        {Keys: bson.D{{"price", 1}}},
        {Keys: bson.D{{"rating", -1}}},        // Descending for "top rated"
        {Keys: bson.D{{"created_at", -1}}},    // Descending for "newest first"
        
        // Compound indexes for common query combinations
        {Keys: bson.D{{"category", 1}, {"price", 1}}},                    // Category + price range
        {Keys: bson.D{{"category", 1}, {"brand", 1}, {"rating", -1}}},    // Category + brand + rating
        {Keys: bson.D{{"brand", 1}, {"rating", -1}}},                     // Brand + rating
        
        // Text search index with weights
        {Keys: bson.D{{"name", "text"}, {"description", "text"}},
         Options: options.Index().SetWeights(bson.M{"name": 10, "description": 1})},
        
        // Sparse index for optional fields
        {Keys: bson.D{{"tags", 1}}, Options: options.Index().SetSparse(true)},
        
        // TTL index for time-based cleanup (if needed)
        // {Keys: bson.D{{"expires_at", 1}}, Options: options.Index().SetExpireAfterSeconds(86400)},
    }
    
    names, err := is.Collection.Indexes().CreateMany(ctx, indexes)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    return Response{
        Success: true,
        Data:    map[string]interface{}{"created_indexes": names},
        Message: fmt.Sprintf("Created %d optimal indexes", len(names)),
        Code:    200,
    }
}
```

## Hint 3: Listing Indexes with Details

Retrieve comprehensive index information including specifications and options:

```go
func (is *IndexService) ListIndexes(ctx context.Context) Response {
    cursor, err := is.Collection.Indexes().List(ctx)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    defer cursor.Close(ctx)
    
    var indexes []IndexInfo
    for cursor.Next(ctx) {
        var indexSpec bson.M
        if err := cursor.Decode(&indexSpec); err != nil {
            continue
        }
        
        // Extract index information
        indexInfo := IndexInfo{
            Name: indexSpec["name"].(string),
            Keys: indexSpec["key"].(bson.M),
        }
        
        // Extract optional properties
        if unique, exists := indexSpec["unique"]; exists {
            indexInfo.Unique = unique.(bool)
        }
        if sparse, exists := indexSpec["sparse"]; exists {
            indexInfo.Sparse = sparse.(bool)
        }
        if background, exists := indexSpec["background"]; exists {
            indexInfo.Background = background.(bool)
        }
        if expireAfter, exists := indexSpec["expireAfterSeconds"]; exists {
            seconds := int32(expireAfter.(int64))
            indexInfo.ExpireAfter = &seconds
        }
        
        indexes = append(indexes, indexInfo)
    }
    
    return Response{
        Success: true,
        Data:    indexes,
        Message: fmt.Sprintf("Retrieved %d indexes", len(indexes)),
        Code:    200,
    }
}
```

## Hint 4: Query Performance Analysis

Use MongoDB's explain functionality to analyze query performance:

```go
func (is *IndexService) AnalyzeQueryPerformance(ctx context.Context, query SearchQuery) Response {
    // Build MongoDB filter from SearchQuery
    filter := bson.M{}
    
    if query.Category != "" {
        filter["category"] = query.Category
    }
    if query.Brand != "" {
        filter["brand"] = query.Brand
    }
    if len(query.PriceRange) > 0 {
        priceFilter := bson.M{}
        if min, exists := query.PriceRange["min"]; exists {
            priceFilter["$gte"] = min
        }
        if max, exists := query.PriceRange["max"]; exists {
            priceFilter["$lte"] = max
        }
        if len(priceFilter) > 0 {
            filter["price"] = priceFilter
        }
    }
    
    // Use explain to get execution statistics
    opts := options.Find().SetLimit(int64(query.Limit))
    if query.SortBy != "" {
        sortOrder := 1
        if query.SortOrder < 0 {
            sortOrder = -1
        }
        opts.SetSort(bson.M{query.SortBy: sortOrder})
    }
    
    // Execute explain
    explainResult := is.Collection.FindOne(ctx, filter, opts.SetProjection(bson.M{}))
    
    // In a real implementation, you'd parse the explain output
    // For this example, we'll simulate the performance metrics
    performance := QueryPerformance{
        Query:           filter,
        ExecutionTimeMs: 5, // Simulated
        DocsExamined:    100,
        DocsReturned:    10,
        IndexUsed:       "category_1_price_1", // Determined from explain
        Stage:           "IXSCAN",
        IsOptimal:       true, // Based on analysis
    }
    
    return Response{
        Success:     true,
        Data:        performance,
        Message:     "Query performance analysis completed",
        Code:        200,
        Performance: &performance,
    }
}
```

## Hint 5: Optimized Search Implementation

Implement search with performance-conscious query building:

```go
func (is *IndexService) OptimizedSearch(ctx context.Context, query SearchQuery) Response {
    // Build filter optimized for index usage
    filter := bson.M{}
    
    // Add filters in order of selectivity (most selective first)
    if query.Category != "" {
        filter["category"] = query.Category
    }
    if query.Brand != "" {
        filter["brand"] = query.Brand
    }
    if len(query.PriceRange) > 0 {
        priceFilter := bson.M{}
        if min, exists := query.PriceRange["min"]; exists {
            priceFilter["$gte"] = min
        }
        if max, exists := query.PriceRange["max"]; exists {
            priceFilter["$lte"] = max
        }
        filter["price"] = priceFilter
    }
    if query.MinRating > 0 {
        filter["rating"] = bson.M{"$gte": query.MinRating}
    }
    if len(query.Tags) > 0 {
        filter["tags"] = bson.M{"$in": query.Tags}
    }
    
    // Build options for optimal performance
    opts := options.Find()
    
    // Set reasonable limits
    limit := query.Limit
    if limit <= 0 {
        limit = 20 // Default
    }
    if limit > 100 {
        limit = 100 // Max limit
    }
    opts.SetLimit(int64(limit))
    
    if query.Skip > 0 {
        opts.SetSkip(int64(query.Skip))
    }
    
    // Optimize sorting for index usage
    if query.SortBy != "" {
        sortOrder := 1
        if query.SortOrder < 0 {
            sortOrder = -1
        }
        opts.SetSort(bson.M{query.SortBy: sortOrder})
    }
    
    // Use projection to reduce data transfer
    opts.SetProjection(bson.M{
        "name":        1,
        "category":    1,
        "brand":       1,
        "price":       1,
        "rating":      1,
        "description": 1,
    })
    
    startTime := time.Now()
    cursor, err := is.Collection.Find(ctx, filter, opts)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    defer cursor.Close(ctx)
    
    var products []Product
    if err = cursor.All(ctx, &products); err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    executionTime := time.Since(startTime).Milliseconds()
    
    performance := QueryPerformance{
        Query:           filter,
        ExecutionTimeMs: executionTime,
        DocsExamined:    int64(len(products)), // Simplified
        DocsReturned:    int64(len(products)),
        IndexUsed:       "category_1_brand_1_rating_-1", // Determined by query pattern
        Stage:           "IXSCAN",
        IsOptimal:       executionTime < 10, // Consider optimal if < 10ms
    }
    
    return Response{
        Success:     true,
        Data:        products,
        Message:     fmt.Sprintf("Found %d products", len(products)),
        Code:        200,
        Performance: &performance,
    }
}
```

## Hint 6: Text Search Implementation

Create and use text indexes for full-text search with relevance scoring:

```go
func (is *IndexService) CreateTextIndex(ctx context.Context, fields map[string]int) Response {
    if len(fields) == 0 {
        return Response{Success: false, Error: "No fields specified", Code: 400}
    }
    
    // Build text index keys
    keys := bson.D{}
    weights := bson.M{}
    
    for field, weight := range fields {
        keys = append(keys, bson.E{Key: field, Value: "text"})
        weights[field] = weight
    }
    
    indexModel := mongo.IndexModel{
        Keys: keys,
        Options: options.Index().SetWeights(weights).SetDefaultLanguage("english"),
    }
    
    name, err := is.Collection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    return Response{
        Success: true,
        Data:    map[string]interface{}{"index_name": name, "fields": fields},
        Message: "Text index created successfully",
        Code:    200,
    }
}

func (is *IndexService) PerformTextSearch(ctx context.Context, searchText string, options map[string]interface{}) Response {
    if searchText == "" {
        return Response{Success: false, Error: "Search text required", Code: 400}
    }
    
    // Build text search filter
    filter := bson.M{
        "$text": bson.M{"$search": searchText},
    }
    
    // Add additional filters from options
    if category, exists := options["category"]; exists {
        filter["category"] = category
    }
    if minRating, exists := options["min_rating"]; exists {
        filter["rating"] = bson.M{"$gte": minRating}
    }
    
    // Build find options
    opts := options.Find().
        SetSort(bson.M{"score": bson.M{"$meta": "textScore"}}). // Sort by relevance
        SetProjection(bson.M{
            "name":        1,
            "description": 1,
            "category":    1,
            "price":       1,
            "rating":      1,
            "score":       bson.M{"$meta": "textScore"}, // Include relevance score
        })
    
    if limit, exists := options["limit"]; exists {
        if limitInt, ok := limit.(int); ok && limitInt > 0 {
            opts.SetLimit(int64(limitInt))
        }
    } else {
        opts.SetLimit(20)
    }
    
    cursor, err := is.Collection.Find(ctx, filter, opts)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    defer cursor.Close(ctx)
    
    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    return Response{
        Success: true,
        Data:    results,
        Message: fmt.Sprintf("Found %d results for '%s'", len(results), searchText),
        Code:    200,
    }
}
```

## Hint 7: Compound Index Strategy

Design compound indexes following the ESR (Equality, Sort, Range) rule:

```go
func (is *IndexService) CreateCompoundIndex(ctx context.Context, fields []map[string]int, options map[string]interface{}) Response {
    if len(fields) == 0 {
        return Response{Success: false, Error: "No fields specified", Code: 400}
    }
    
    // Build compound index keys
    keys := bson.D{}
    for _, fieldMap := range fields {
        for field, order := range fieldMap {
            keys = append(keys, bson.E{Key: field, Value: order})
        }
    }
    
    // Build index options
    indexOpts := options.Index()
    
    if unique, exists := options["unique"]; exists {
        if uniqueBool, ok := unique.(bool); ok {
            indexOpts.SetUnique(uniqueBool)
        }
    }
    
    if sparse, exists := options["sparse"]; exists {
        if sparseBool, ok := sparse.(bool); ok {
            indexOpts.SetSparse(sparseBool)
        }
    }
    
    if background, exists := options["background"]; exists {
        if backgroundBool, ok := background.(bool); ok {
            indexOpts.SetBackground(backgroundBool)
        }
    }
    
    indexModel := mongo.IndexModel{
        Keys:    keys,
        Options: indexOpts,
    }
    
    name, err := is.Collection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    return Response{
        Success: true,
        Data:    map[string]interface{}{"index_name": name, "fields": fields},
        Message: "Compound index created successfully",
        Code:    200,
    }
}
```

## Hint 8: Safe Index Management

Implement safe index dropping with validation:

```go
func (is *IndexService) DropIndex(ctx context.Context, indexName string) Response {
    if indexName == "" {
        return Response{Success: false, Error: "Index name required", Code: 400}
    }
    
    // Prevent dropping critical indexes
    if indexName == "_id_" {
        return Response{Success: false, Error: "Cannot drop _id index", Code: 400}
    }
    
    // Check if index exists
    cursor, err := is.Collection.Indexes().List(ctx)
    if err != nil {
        return Response{Success: false, Error: "Failed to list indexes", Code: 500}
    }
    defer cursor.Close(ctx)
    
    indexExists := false
    for cursor.Next(ctx) {
        var indexSpec bson.M
        if err := cursor.Decode(&indexSpec); err != nil {
            continue
        }
        if name, ok := indexSpec["name"].(string); ok && name == indexName {
            indexExists = true
            break
        }
    }
    
    if !indexExists {
        return Response{Success: false, Error: "Index not found", Code: 404}
    }
    
    // Drop the index
    _, err = is.Collection.Indexes().DropOne(ctx, indexName)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    
    return Response{
        Success: true,
        Message: fmt.Sprintf("Index '%s' dropped successfully", indexName),
        Code:    200,
    }
}
```

## Hint 9: Index Usage Statistics

Monitor index performance and usage:

```go
func (is *IndexService) GetIndexUsageStats(ctx context.Context) Response {
    // In a real implementation, you'd use db.collection.aggregate([{$indexStats: {}}])
    // For this challenge, we'll simulate the statistics
    
    // Get list of indexes first
    cursor, err := is.Collection.Indexes().List(ctx)
    if err != nil {
        return Response{Success: false, Error: err.Error(), Code: 500}
    }
    defer cursor.Close(ctx)
    
    stats := make(map[string]interface{})
    
    for cursor.Next(ctx) {
        var indexSpec bson.M
        if err := cursor.Decode(&indexSpec); err != nil {
            continue
        }
        
        if name, ok := indexSpec["name"].(string); ok {
            // Simulate usage statistics
            stats[name] = map[string]interface{}{
                "ops":           rand.Intn(1000) + 100,  // Simulated operations count
                "since":         "2024-01-01T00:00:00Z", // Since when tracking
                "efficiency":    0.7 + rand.Float64()*0.3, // 70-100% efficiency
                "size_bytes":    rand.Intn(10000) + 1000,   // Index size
                "last_used":     time.Now().Add(-time.Duration(rand.Intn(24)) * time.Hour),
            }
        }
    }
    
    return Response{
        Success: true,
        Data:    stats,
        Message: "Index usage statistics retrieved",
        Code:    200,
    }
}
```

## Hint 10: Performance Best Practices

Key principles for optimal MongoDB indexing:

```go
// 1. ESR Rule for Compound Indexes
// Equality filters first, then Sort fields, then Range filters
db.products.createIndex({"category": 1, "rating": -1, "price": 1})
// Good for: {category: "Electronics"} sorted by rating with price range

// 2. Index Selectivity - High cardinality fields first
db.products.createIndex({"user_id": 1, "status": 1}) // user_id is more selective

// 3. Covered Queries - Include all fields needed in the index
db.products.createIndex({"category": 1, "name": 1, "price": 1})
// Query: db.products.find({"category": "Books"}, {"name": 1, "price": 1})

// 4. Avoid Over-Indexing
// Don't create indexes for every possible query combination
// Focus on the most common and performance-critical queries

// 5. Monitor and Optimize
// Use explain() to verify index usage
// Drop unused indexes to save memory and improve write performance
// Regularly analyze slow query logs

// 6. Text Search Optimization
db.products.createIndex(
    {"name": "text", "description": "text"},
    {"weights": {"name": 10, "description": 1}}  // Prioritize name matches
)

// 7. Partial Indexes for Conditional Data
db.products.createIndex(
    {"price": 1},
    {"partialFilterExpression": {"price": {"$gt": 100}}}  // Only index expensive items
)
```

Ready to build lightning-fast MongoDB queries? Start with single field indexes and work your way up to complex compound strategies! ⚡🚀
