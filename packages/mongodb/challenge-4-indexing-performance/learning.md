# Learning: MongoDB Indexing & Performance Mastery

## 🎯 **What You'll Learn**

This challenge teaches advanced MongoDB indexing strategies used by high-performance applications. You'll master the art of query optimization, index design, and performance tuning techniques used by major platforms like Amazon, Netflix, and Google.

## 📚 **Understanding Database Indexes**

### **What Are Indexes?**

Think of database indexes like the index in a book - they provide a fast way to find specific information without scanning every page:

```go
// Without Index (Collection Scan)
// MongoDB scans every document: O(n) time complexity
db.products.find({"category": "Electronics"})
// Examines: 1,000,000 documents to find 50 matches

// With Index (Index Scan)  
// MongoDB uses index: O(log n) time complexity
db.products.createIndex({"category": 1})
db.products.find({"category": "Electronics"})
// Examines: ~17 index entries to find 50 matches
```

### **Index Data Structure**

MongoDB uses B-tree indexes, which provide efficient operations:

```
Index on "category" field:
                    [Electronics]
                   /             \
            [Books]                 [Sports]
           /      \               /         \
    [Audio]      [Clothing]  [Games]    [Toys]
     |              |          |          |
   docs...       docs...    docs...   docs...
```

## 🔍 **Index Types and Use Cases**

### **1. Single Field Indexes**

The most basic and commonly used indexes:

```go
// Ascending index (1)
db.products.createIndex({"price": 1})
// Good for: price range queries, sorting by price ascending

// Descending index (-1)  
db.products.createIndex({"rating": -1})
// Good for: "highest rated first" queries

// Practical Go implementation
func CreateBasicIndexes(collection *mongo.Collection, ctx context.Context) error {
    indexes := []mongo.IndexModel{
        {Keys: bson.D{{"category", 1}}},     // Category filter
        {Keys: bson.D{{"price", 1}}},        // Price range/sort
        {Keys: bson.D{{"rating", -1}}},      // Top-rated products
        {Keys: bson.D{{"created_at", -1}}},  // Newest first
        {Keys: bson.D{{"stock", 1}}},        // Inventory queries
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

### **2. Compound Indexes**

Multiple fields in a single index - the powerhouse of query optimization:

```go
// Compound index following ESR rule (Equality, Sort, Range)
db.products.createIndex({"category": 1, "rating": -1, "price": 1})

// This index efficiently supports:
// 1. Category filter + rating sort + price range
db.products.find({"category": "Electronics"}).sort({"rating": -1})
db.products.find({"category": "Electronics", "price": {"$gte": 100, "$lte": 500}})
db.products.find({"category": "Electronics"}).sort({"rating": -1}).limit(10)

// ESR Rule Explanation:
// E (Equality): Exact match filters first
// S (Sort): Sort fields next  
// R (Range): Range queries last

// Go implementation
func CreateCompoundIndexes(collection *mongo.Collection, ctx context.Context) error {
    indexes := []mongo.IndexModel{
        // E-commerce search patterns
        {Keys: bson.D{{"category", 1}, {"brand", 1}, {"rating", -1}}},
        {Keys: bson.D{{"category", 1}, {"price", 1}}},
        {Keys: bson.D{{"brand", 1}, {"rating", -1}}},
        
        // Time-based queries with filters
        {Keys: bson.D{{"status", 1}, {"created_at", -1}}},
        {Keys: bson.D{{"user_id", 1}, {"created_at", -1}}},
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

### **3. Text Indexes**

Full-text search capabilities with relevance scoring:

```go
// Create text index with field weights
db.products.createIndex(
    {"name": "text", "description": "text", "tags": "text"},
    {
        "weights": {
            "name": 10,        // Name matches are most important
            "description": 5,   // Description matches are medium
            "tags": 1          // Tag matches are least important
        },
        "default_language": "english"
    }
)

// Text search queries
db.products.find({"$text": {"$search": "wireless bluetooth headphones"}})
db.products.find(
    {"$text": {"$search": "smartphone"}, "category": "Electronics"},
    {"score": {"$meta": "textScore"}}
).sort({"score": {"$meta": "textScore"}})

// Go implementation
func CreateTextIndex(collection *mongo.Collection, ctx context.Context) error {
    indexModel := mongo.IndexModel{
        Keys: bson.D{
            {"name", "text"},
            {"description", "text"},
            {"tags", "text"},
        },
        Options: options.Index().
            SetWeights(bson.M{
                "name":        10,
                "description": 5,
                "tags":        1,
            }).
            SetDefaultLanguage("english").
            SetLanguageOverride("language"),
    }
    
    _, err := collection.Indexes().CreateOne(ctx, indexModel)
    return err
}

// Text search implementation
func TextSearch(collection *mongo.Collection, ctx context.Context, searchTerm string) ([]bson.M, error) {
    filter := bson.M{
        "$text": bson.M{"$search": searchTerm},
    }
    
    opts := options.Find().
        SetProjection(bson.M{
            "name":        1,
            "description": 1,
            "price":       1,
            "score":       bson.M{"$meta": "textScore"},
        }).
        SetSort(bson.M{"score": bson.M{"$meta": "textScore"}}).
        SetLimit(20)
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    
    return results, nil
}
```

### **4. Specialized Indexes**

Advanced indexing for specific use cases:

```go
// Sparse Index - only indexes documents that have the field
db.products.createIndex({"tags": 1}, {"sparse": true})
// Saves space when many documents don't have the field

// Unique Index - enforces uniqueness
db.products.createIndex({"sku": 1}, {"unique": true})
// Prevents duplicate SKUs

// TTL Index - automatically expires documents
db.sessions.createIndex({"created_at": 1}, {"expireAfterSeconds": 3600})
// Documents expire after 1 hour

// Partial Index - only indexes documents matching a condition
db.products.createIndex(
    {"price": 1},
    {"partialFilterExpression": {"price": {"$gt": 100}}}
)
// Only indexes products over $100

// Multikey Index - automatically handles array fields
db.products.createIndex({"tags": 1})
// Works efficiently with array fields

// Go implementation
func CreateSpecializedIndexes(collection *mongo.Collection, ctx context.Context) error {
    indexes := []mongo.IndexModel{
        // Sparse index for optional fields
        {
            Keys:    bson.D{{"tags", 1}},
            Options: options.Index().SetSparse(true),
        },
        
        // Unique index for business keys
        {
            Keys:    bson.D{{"sku", 1}},
            Options: options.Index().SetUnique(true),
        },
        
        // TTL index for temporary data
        {
            Keys:    bson.D{{"expires_at", 1}},
            Options: options.Index().SetExpireAfterSeconds(86400), // 24 hours
        },
        
        // Partial index for conditional data
        {
            Keys: bson.D{{"premium_price", 1}},
            Options: options.Index().SetPartialFilterExpression(
                bson.M{"price": bson.M{"$gte": 1000}},
            ),
        },
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

## ⚡ **Query Performance Optimization**

### **Understanding Query Execution**

MongoDB's query planner chooses the best index for each query:

```go
// Query execution stages
// 1. COLLSCAN - Full collection scan (slow)
// 2. IXSCAN - Index scan (fast)
// 3. FETCH - Retrieve documents from index results
// 4. SORT - In-memory sorting (expensive)
// 5. LIMIT - Limit results

// Use explain() to analyze query performance
func AnalyzeQuery(collection *mongo.Collection, ctx context.Context, filter bson.M) {
    // Get execution statistics
    cursor := collection.FindOne(ctx, filter)
    
    // In practice, you'd use:
    // cursor, err := collection.Find(ctx, filter, options.Find().SetExplain(true))
    
    // Example explain output analysis:
    explainResult := bson.M{
        "executionStats": bson.M{
            "totalDocsExamined": 1000,    // Documents scanned
            "totalDocsReturned": 10,      // Documents returned
            "executionTimeMillis": 50,    // Query time
            "indexesUsed": []string{"category_1_price_1"},
            "stage": "IXSCAN",            // Used index scan
        },
    }
    
    // Performance indicators:
    // - Low docsExamined/docsReturned ratio = efficient
    // - IXSCAN stage = using index
    // - Low executionTimeMillis = fast query
}
```

### **Index Selection Strategy**

MongoDB automatically chooses indexes, but you can optimize:

```go
// Index hint - force specific index usage
func QueryWithHint(collection *mongo.Collection, ctx context.Context) ([]Product, error) {
    filter := bson.M{
        "category": "Electronics",
        "price": bson.M{"$gte": 100, "$lte": 500},
    }
    
    // Force use of specific index
    opts := options.Find().SetHint(bson.M{"category": 1, "price": 1})
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var products []Product
    if err = cursor.All(ctx, &products); err != nil {
        return nil, err
    }
    
    return products, nil
}

// Covered queries - all data comes from index
func CoveredQuery(collection *mongo.Collection, ctx context.Context) error {
    // Create covering index
    indexModel := mongo.IndexModel{
        Keys: bson.D{{"category", 1}, {"name", 1}, {"price", 1}},
    }
    collection.Indexes().CreateOne(ctx, indexModel)
    
    // Query that only needs indexed fields
    opts := options.Find().SetProjection(bson.M{
        "_id":      0, // Exclude _id
        "name":     1, // Include name (in index)
        "price":    1, // Include price (in index)
    })
    
    cursor, err := collection.Find(ctx, bson.M{"category": "Electronics"}, opts)
    // This query is "covered" - no document fetching needed!
    
    return err
}
```

## 🏗️ **Index Design Patterns**

### **E-commerce Search Optimization**

Real-world indexing strategy for product search:

```go
func CreateEcommerceIndexes(collection *mongo.Collection, ctx context.Context) error {
    indexes := []mongo.IndexModel{
        // 1. Category browsing
        {Keys: bson.D{{"category", 1}, {"featured", -1}, {"rating", -1}}},
        
        // 2. Brand + category filtering
        {Keys: bson.D{{"category", 1}, {"brand", 1}, {"price", 1}}},
        
        // 3. Price-focused queries
        {Keys: bson.D{{"price", 1}, {"rating", -1}}},
        
        // 4. New arrivals
        {Keys: bson.D{{"created_at", -1}, {"category", 1}}},
        
        // 5. Inventory management
        {Keys: bson.D{{"stock", 1}, {"category", 1}}},
        
        // 6. Full-text product search
        {
            Keys: bson.D{{"name", "text"}, {"description", "text"}},
            Options: options.Index().SetWeights(bson.M{"name": 10, "description": 1}),
        },
        
        // 7. User-specific queries
        {Keys: bson.D{{"seller_id", 1}, {"status", 1}, {"created_at", -1}}},
        
        // 8. Geographic queries (if applicable)
        {Keys: bson.D{{"location", "2dsphere"}}},
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}

// Optimized product search function
func OptimizedProductSearch(collection *mongo.Collection, ctx context.Context, params SearchParams) (*SearchResults, error) {
    filter := bson.M{}
    
    // Build filter in order of selectivity
    if params.Category != "" {
        filter["category"] = params.Category
    }
    if params.Brand != "" {
        filter["brand"] = params.Brand
    }
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
    if params.MinRating > 0 {
        filter["rating"] = bson.M{"$gte": params.MinRating}
    }
    if params.InStock {
        filter["stock"] = bson.M{"$gt": 0}
    }
    
    // Build options for optimal performance
    opts := options.Find()
    
    // Sorting strategy
    switch params.SortBy {
    case "price_asc":
        opts.SetSort(bson.M{"price": 1})
    case "price_desc":
        opts.SetSort(bson.M{"price": -1})
    case "rating":
        opts.SetSort(bson.M{"rating": -1, "review_count": -1})
    case "newest":
        opts.SetSort(bson.M{"created_at": -1})
    default:
        opts.SetSort(bson.M{"featured": -1, "rating": -1})
    }
    
    // Pagination
    opts.SetSkip(int64(params.Skip)).SetLimit(int64(params.Limit))
    
    // Projection for faster data transfer
    opts.SetProjection(bson.M{
        "name":         1,
        "price":        1,
        "rating":       1,
        "image_url":    1,
        "brand":        1,
        "category":     1,
        "stock":        1,
    })
    
    startTime := time.Now()
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var products []Product
    if err = cursor.All(ctx, &products); err != nil {
        return nil, err
    }
    
    executionTime := time.Since(startTime)
    
    return &SearchResults{
        Products:        products,
        Total:          len(products),
        ExecutionTimeMs: executionTime.Milliseconds(),
    }, nil
}
```

### **Time-Series Data Indexing**

Optimizing for time-based queries:

```go
func CreateTimeSeriesIndexes(collection *mongo.Collection, ctx context.Context) error {
    indexes := []mongo.IndexModel{
        // 1. Recent data queries
        {Keys: bson.D{{"timestamp", -1}}},
        
        // 2. User activity over time
        {Keys: bson.D{{"user_id", 1}, {"timestamp", -1}}},
        
        // 3. Event type analysis
        {Keys: bson.D{{"event_type", 1}, {"timestamp", -1}}},
        
        // 4. Compound time-series index
        {Keys: bson.D{{"device_id", 1}, {"metric_type", 1}, {"timestamp", -1}}},
        
        // 5. TTL for data retention
        {
            Keys:    bson.D{{"timestamp", 1}},
            Options: options.Index().SetExpireAfterSeconds(2592000), // 30 days
        },
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

## 📊 **Performance Monitoring and Optimization**

### **Index Usage Analysis**

Monitor which indexes are actually being used:

```go
func AnalyzeIndexUsage(collection *mongo.Collection, ctx context.Context) (*IndexUsageReport, error) {
    // Get index usage statistics
    pipeline := mongo.Pipeline{
        {{"$indexStats", bson.M{}}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var stats []bson.M
    if err = cursor.All(ctx, &stats); err != nil {
        return nil, err
    }
    
    report := &IndexUsageReport{
        TotalIndexes: len(stats),
        UnusedIndexes: []string{},
        HighUsageIndexes: []string{},
    }
    
    for _, stat := range stats {
        name := stat["name"].(string)
        accesses := stat["accesses"].(bson.M)
        ops := accesses["ops"].(int64)
        
        if ops == 0 {
            report.UnusedIndexes = append(report.UnusedIndexes, name)
        } else if ops > 10000 {
            report.HighUsageIndexes = append(report.HighUsageIndexes, name)
        }
    }
    
    return report, nil
}

// Slow query analysis
func AnalyzeSlowQueries(collection *mongo.Collection, ctx context.Context) error {
    // Enable profiling for slow queries
    // db.setProfilingLevel(2, {slowms: 100})
    
    // Query the profiler collection
    profilerCollection := collection.Database().Collection("system.profile")
    
    filter := bson.M{
        "millis": bson.M{"$gte": 100}, // Queries taking > 100ms
        "ts": bson.M{"$gte": time.Now().Add(-24 * time.Hour)}, // Last 24 hours
    }
    
    cursor, err := profilerCollection.Find(ctx, filter)
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    for cursor.Next(ctx) {
        var profile bson.M
        if err := cursor.Decode(&profile); err != nil {
            continue
        }
        
        // Analyze slow query patterns
        command := profile["command"].(bson.M)
        executionTime := profile["millis"].(int64)
        
        fmt.Printf("Slow query detected: %v (took %dms)\n", command, executionTime)
        
        // Check if query used indexes
        if executionStats, exists := profile["executionStats"]; exists {
            stats := executionStats.(bson.M)
            if stage, exists := stats["executionStages"]; exists {
                stageInfo := stage.(bson.M)
                if stageInfo["stage"].(string) == "COLLSCAN" {
                    fmt.Printf("WARNING: Query performed collection scan!\n")
                }
            }
        }
    }
    
    return nil
}
```

### **Index Maintenance**

Keep indexes optimized over time:

```go
func OptimizeIndexes(collection *mongo.Collection, ctx context.Context) error {
    // 1. Rebuild fragmented indexes
    cursor, err := collection.Indexes().List(ctx)
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    for cursor.Next(ctx) {
        var indexSpec bson.M
        if err := cursor.Decode(&indexSpec); err != nil {
            continue
        }
        
        indexName := indexSpec["name"].(string)
        
        // Skip system indexes
        if indexName == "_id_" {
            continue
        }
        
        // Reindex if needed (this is a simplified example)
        // In practice, you'd check index statistics first
        fmt.Printf("Considering reindex for: %s\n", indexName)
    }
    
    // 2. Remove unused indexes (be careful!)
    unusedIndexes := []string{"old_unused_index_1", "temporary_index_2"}
    
    for _, indexName := range unusedIndexes {
        fmt.Printf("Dropping unused index: %s\n", indexName)
        _, err := collection.Indexes().DropOne(ctx, indexName)
        if err != nil {
            fmt.Printf("Failed to drop index %s: %v\n", indexName, err)
        }
    }
    
    return nil
}
```

## 🌍 **Real-World Applications**

### **High-Traffic E-commerce Platform**

```go
// Amazon-scale product search optimization
func AmazonStyleIndexing(collection *mongo.Collection, ctx context.Context) error {
    // Indexes supporting millions of products and searches per second
    indexes := []mongo.IndexModel{
        // Primary search patterns
        {Keys: bson.D{{"category", 1}, {"subcategory", 1}, {"rating", -1}, {"price", 1}}},
        {Keys: bson.D{{"brand", 1}, {"category", 1}, {"rating", -1}}},
        
        // Text search with business logic
        {
            Keys: bson.D{{"title", "text"}, {"description", "text"}, {"keywords", "text"}},
            Options: options.Index().SetWeights(bson.M{
                "title": 100,      // Product title most important
                "keywords": 50,    // Keywords very important  
                "description": 1,  // Description least important
            }),
        },
        
        // Personalization and recommendations
        {Keys: bson.D{{"tags", 1}, {"rating", -1}, {"sales_rank", 1}}},
        
        // Inventory and fulfillment
        {Keys: bson.D{{"warehouse_id", 1}, {"stock", -1}, {"category", 1}}},
        
        // Pricing and promotions
        {Keys: bson.D{{"price", 1}, {"discount_percent", -1}}},
        
        // Geographic optimization
        {Keys: bson.D{{"shipping_zones", 1}, {"category", 1}, {"price", 1}}},
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

### **Social Media Feed Optimization**

```go
// Twitter/Facebook-style feed indexing
func SocialMediaIndexing(collection *mongo.Collection, ctx context.Context) error {
    indexes := []mongo.IndexModel{
        // Timeline queries
        {Keys: bson.D{{"user_id", 1}, {"created_at", -1}}},
        
        // Feed aggregation
        {Keys: bson.D{{"followers", 1}, {"created_at", -1}}},
        
        // Trending content
        {Keys: bson.D{{"engagement_score", -1}, {"created_at", -1}}},
        
        // Hashtag searches
        {Keys: bson.D{{"hashtags", 1}, {"created_at", -1}}},
        
        // Geographic feeds
        {Keys: bson.D{{"location", "2dsphere"}, {"created_at", -1}}},
        
        // Content moderation
        {Keys: bson.D{{"status", 1}, {"reported_count", -1}}},
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

## 📚 **Best Practices Summary**

### **Index Design Principles**

1. **Analyze Query Patterns First**: Create indexes based on actual query patterns, not assumptions
2. **Follow ESR Rule**: Equality, Sort, Range for compound indexes
3. **Monitor Performance**: Use explain() and profiling to verify index effectiveness
4. **Balance Read vs Write**: More indexes = faster reads but slower writes
5. **Consider Cardinality**: High-cardinality fields make better index candidates
6. **Use Covering Indexes**: Include all query fields in the index when possible
7. **Regular Maintenance**: Monitor usage and drop unused indexes

### **Common Pitfalls to Avoid**

```go
// ❌ Bad: Too many indexes
// Creates maintenance overhead and slows writes
db.products.createIndex({"field1": 1})
db.products.createIndex({"field2": 1})
db.products.createIndex({"field3": 1})
// ... 50+ indexes

// ✅ Good: Strategic compound indexes
db.products.createIndex({"field1": 1, "field2": 1, "field3": 1})

// ❌ Bad: Wrong field order in compound index
db.products.createIndex({"price": 1, "category": 1}) // Range field first

// ✅ Good: ESR rule applied
db.products.createIndex({"category": 1, "price": 1}) // Equality field first

// ❌ Bad: Unused indexes consuming resources
db.products.createIndex({"rarely_queried_field": 1})

// ✅ Good: Regular index usage analysis and cleanup
```

## 🚀 **Next Steps**

After mastering indexing and performance:

1. **Challenge 5**: Learn transactions for consistent multi-document operations
2. **Advanced Topics**: Sharding strategies, replica set optimization
3. **Production Skills**: Monitoring, alerting, and capacity planning

## 🔗 **Additional Resources**

- [MongoDB Indexing Strategies](https://docs.mongodb.com/manual/applications/indexes/)
- [Query Performance Analysis](https://docs.mongodb.com/manual/tutorial/analyze-query-performance/)
- [Index Build Performance](https://docs.mongodb.com/manual/tutorial/build-indexes-on-replica-sets/)
- [MongoDB Performance Best Practices](https://www.mongodb.com/blog/post/performance-best-practices-mongodb-data-modeling-and-memory-sizing)

Ready to build lightning-fast MongoDB applications? Master the art of indexing! ⚡📊🚀
