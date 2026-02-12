# Challenge 4: Indexing & Performance

Build a **High-Performance Search System** using MongoDB indexing strategies to optimize query performance and create lightning-fast search experiences.

## Challenge Requirements

Implement an `IndexService` that manages MongoDB indexes and optimizes query performance. Your service should provide the following capabilities:

- **Index Management**: Create, list, and drop various types of indexes (single, compound, text, sparse, TTL).
- **Performance Analysis**: Analyze query execution plans and identify optimization opportunities.
- **Optimized Search**: Implement search functionality that leverages indexes for maximum performance.
- **Text Search**: Create and utilize text indexes for full-text search with relevance scoring.
- **Index Monitoring**: Track index usage statistics and provide optimization recommendations.

## Data Structures

```go
type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`
    Description string             `bson:"description" json:"description"`
    Category    string             `bson:"category" json:"category"`
    Brand       string             `bson:"brand" json:"brand"`
    Price       float64            `bson:"price" json:"price"`
    Stock       int                `bson:"stock" json:"stock"`
    Rating      float64            `bson:"rating" json:"rating"`
    Tags        []string           `bson:"tags" json:"tags"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
    Metadata    ProductMetadata    `bson:"metadata" json:"metadata"`
}

type IndexInfo struct {
    Name       string                 `bson:"name" json:"name"`
    Keys       map[string]interface{} `bson:"key" json:"keys"`
    Unique     bool                   `bson:"unique,omitempty" json:"unique,omitempty"`
    Sparse     bool                   `bson:"sparse,omitempty" json:"sparse,omitempty"`
    Background bool                   `bson:"background,omitempty" json:"background,omitempty"`
    ExpireAfter *int32                `bson:"expireAfterSeconds,omitempty" json:"expire_after,omitempty"`
}

type QueryPerformance struct {
    Query            map[string]interface{} `json:"query"`
    ExecutionTimeMs  int64                  `json:"execution_time_ms"`
    DocsExamined     int64                  `json:"docs_examined"`
    DocsReturned     int64                  `json:"docs_returned"`
    IndexUsed        string                 `json:"index_used"`
    Stage            string                 `json:"stage"`
    IsOptimal        bool                   `json:"is_optimal"`
}
```

## Index Types and Examples

### **Single Field Indexes**
```go
// Create index on category field
db.products.createIndex({"category": 1})

// Create descending index on price
db.products.createIndex({"price": -1})
```

### **Compound Indexes**
```go
// Compound index for category + price queries
db.products.createIndex({"category": 1, "price": 1})

// Multi-field compound index
db.products.createIndex({"category": 1, "brand": 1, "rating": -1})
```

### **Text Indexes**
```go
// Text search index with field weights
db.products.createIndex({
    "name": "text",
    "description": "text"
}, {
    "weights": {
        "name": 10,        // Higher weight for name matches
        "description": 1
    }
})
```

### **Specialized Indexes**
```go
// Sparse index (only indexes documents with the field)
db.products.createIndex({"tags": 1}, {"sparse": true})

// TTL index (auto-expires documents)
db.products.createIndex({"created_at": 1}, {"expireAfterSeconds": 86400})

// Unique index
db.products.createIndex({"sku": 1}, {"unique": true})
```

## Performance Optimization Examples

**Query Performance Analysis:**
```json
{
    "query": {"category": "Electronics", "price": {"$gte": 100, "$lte": 500}},
    "execution_time_ms": 2,
    "docs_examined": 150,
    "docs_returned": 25,
    "index_used": "category_1_price_1",
    "stage": "IXSCAN",
    "is_optimal": true
}
```

**Optimized Search Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": "65b23c2e01d2a3b4c5d6e7f8",
            "name": "iPhone 15 Pro",
            "category": "Electronics",
            "price": 999.99,
            "rating": 4.8
        }
    ],
    "message": "Found 25 products",
    "performance": {
        "execution_time_ms": 2,
        "index_used": "category_1_price_1",
        "is_optimal": true
    }
}
```

## Testing Requirements

Your solution must pass tests for:
- Creating optimal index sets for common query patterns
- Listing all indexes with detailed information and options
- Analyzing query performance and identifying optimization opportunities
- Performing optimized searches using appropriate indexes
- Creating and utilizing text indexes for full-text search
- Managing compound indexes for complex query patterns
- Safely dropping indexes with proper validation
- Retrieving index usage statistics and recommendations
- Handling index creation conflicts and errors
- Consistent response structure with performance metrics

## Key Indexing Concepts

- **Index Selectivity**: Choose fields with high cardinality for better performance
- **Compound Index Order**: Field order matters - most selective fields first
- **Index Intersection**: MongoDB can use multiple indexes for complex queries
- **Covered Queries**: Queries that can be satisfied entirely from index data
- **Index Hints**: Force specific index usage for query optimization
- **ESR Rule**: Equality, Sort, Range - optimal compound index field ordering
- **Write Performance**: Indexes improve read performance but slow down writes
- **Memory Usage**: Indexes consume RAM - balance performance vs. memory
