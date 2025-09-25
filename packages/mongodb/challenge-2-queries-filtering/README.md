# Challenge 2: Advanced Queries & Filtering

Build a **Product Search System** using MongoDB with advanced query operators, sorting, and pagination.

## Challenge Requirements

Implement a product search system with the following operations:

- **GetProductsByCategory** - Filter products by category with additional filters
- **GetProductsByPriceRange** - Find products within price range
- **SearchProductsByName** - Text search in name and description
- **GetProductsWithPagination** - Paginated results with filtering
- **GetProductsSorted** - Custom sorting with multiple criteria
- **FilterProducts** - Complex multi-field filtering with projection
- **GetProductsByTags** - Tag-based filtering (any/all matching)
- **GetTopRatedProducts** - Highest rated products by category

## Data Structure

```go
type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`
    Description string             `bson:"description" json:"description"`
    Category    string             `bson:"category" json:"category"`
    Price       float64            `bson:"price" json:"price"`
    Stock       int                `bson:"stock" json:"stock"`
    Tags        []string           `bson:"tags" json:"tags"`
    Brand       string             `bson:"brand" json:"brand"`
    Rating      float64            `bson:"rating" json:"rating"`
}
```

## Query Examples

**Price Range Query**
```go
// Products between $100 and $500
filter := bson.M{
    "price": bson.M{
        "$gte": 100.0,
        "$lte": 500.0,
    },
}
```

**Text Search Query**
```go
// Search in name and description
filter := bson.M{
    "$or": []bson.M{
        {"name": bson.M{"$regex": primitive.Regex{Pattern: "iPhone", Options: "i"}}},
        {"description": bson.M{"$regex": primitive.Regex{Pattern: "iPhone", Options: "i"}}},
    },
}
```

**Tag Filtering**
```go
// Products with any of these tags
filter := bson.M{"tags": bson.M{"$in": []string{"smartphone", "premium"}}}

// Products with all of these tags
filter := bson.M{"tags": bson.M{"$all": []string{"smartphone", "premium"}}}
```

## Testing Requirements

Your solution must pass tests for:
- Category filtering with additional filter criteria
- Price range queries with proper validation
- Text search with case-insensitive regex matching
- Pagination with correct skip/limit calculations
- Multi-field sorting with proper precedence
- Complex filtering with multiple conditions
- Tag-based filtering (any/all matching)
- Top-rated product queries with category filtering
