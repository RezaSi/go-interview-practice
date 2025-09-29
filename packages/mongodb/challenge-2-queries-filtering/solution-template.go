package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Product represents a product document in MongoDB
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
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// ProductFilter represents filtering criteria for products
type ProductFilter struct {
	Category   string   `json:"category,omitempty"`
	MinPrice   float64  `json:"min_price,omitempty"`
	MaxPrice   float64  `json:"max_price,omitempty"`
	Brand      string   `json:"brand,omitempty"`
	MinRating  float64  `json:"min_rating,omitempty"`
	InStock    bool     `json:"in_stock,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	SearchTerm string   `json:"search_term,omitempty"`
}

// SortOptions represents sorting criteria
type SortOptions struct {
	Field string `json:"field"` // price, rating, name, created_at
	Order int    `json:"order"` // 1 for ascending, -1 for descending
}

// PaginationOptions represents pagination parameters
type PaginationOptions struct {
	Page  int `json:"page"`  // Page number (1-based)
	Limit int `json:"limit"` // Items per page
}

// PaginatedResponse represents a paginated result
type PaginatedResponse struct {
	Data       []Product `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// ProductService handles product-related database operations
type ProductService struct {
	Collection *mongo.Collection
}

// GetProductsByCategory retrieves products by category with optional filtering
func (ps *ProductService) GetProductsByCategory(ctx context.Context, category string, filter ProductFilter) Response {
	// TODO: Implement category-based product retrieval
	// Steps:
	// 1. Build filter starting with category
	// 2. Add additional filters from ProductFilter struct
	// 3. Execute query with Find()
	// 4. Return products matching the criteria

	// Hint: Use bson.M{"category": category} as base filter
	// Hint: Add price range with $gte and $lte operators
	// Hint: Use $in operator for tags filtering

	return Response{
		Success: false,
		Error:   "GetProductsByCategory not implemented",
		Code:    500,
	}
}

// GetProductsByPriceRange retrieves products within a price range
func (ps *ProductService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) Response {
	// TODO: Implement price range filtering
	// Steps:
	// 1. Create filter with $gte and $lte operators
	// 2. Handle edge cases (negative prices, min > max)
	// 3. Execute query and return results
	// 4. Sort by price ascending for better UX

	// Hint: Use bson.M{"price": bson.M{"$gte": minPrice, "$lte": maxPrice}}
	// Hint: Add sorting with options.Find().SetSort(bson.M{"price": 1})

	return Response{
		Success: false,
		Error:   "GetProductsByPriceRange not implemented",
		Code:    500,
	}
}

// SearchProductsByName searches products by name using regex
func (ps *ProductService) SearchProductsByName(ctx context.Context, searchTerm string, caseSensitive bool) Response {
	// TODO: Implement name-based search with regex
	// Steps:
	// 1. Create regex pattern for case-insensitive search
	// 2. Build filter using $regex operator
	// 3. Search in both name and description fields
	// 4. Return matching products sorted by relevance

	// Hint: Use bson.M{"$or": []bson.M{...}} for multiple field search
	// Hint: Use primitive.Regex{Pattern: searchTerm, Options: "i"} for case-insensitive
	// Hint: Consider using $text search for better performance (bonus)

	return Response{
		Success: false,
		Error:   "SearchProductsByName not implemented",
		Code:    500,
	}
}

// GetProductsWithPagination retrieves products with pagination support
func (ps *ProductService) GetProductsWithPagination(ctx context.Context, pagination PaginationOptions, filter ProductFilter) Response {
	// TODO: Implement pagination with filtering
	// Steps:
	// 1. Build filter from ProductFilter
	// 2. Count total documents matching filter
	// 3. Calculate skip value: (page - 1) * limit
	// 4. Execute paginated query with Skip() and Limit()
	// 5. Return PaginatedResponse with metadata

	// Hint: Use collection.CountDocuments() for total count
	// Hint: Use options.Find().SetSkip().SetLimit() for pagination
	// Hint: Calculate total pages: (total + limit - 1) / limit

	return Response{
		Success: false,
		Error:   "GetProductsWithPagination not implemented",
		Code:    500,
	}
}

// GetProductsSorted retrieves products with custom sorting
func (ps *ProductService) GetProductsSorted(ctx context.Context, sortOptions []SortOptions, limit int) Response {
	// TODO: Implement multi-field sorting
	// Steps:
	// 1. Build sort document from SortOptions array
	// 2. Handle multiple sort criteria (price, rating, name, etc.)
	// 3. Apply limit if specified
	// 4. Execute query with sorting options

	// Hint: Use bson.D for ordered sort criteria
	// Hint: bson.D{{"price", -1}, {"rating", -1}} for price desc, rating desc
	// Hint: Validate sort fields to prevent injection

	return Response{
		Success: false,
		Error:   "GetProductsSorted not implemented",
		Code:    500,
	}
}

// FilterProducts applies complex multi-field filtering
func (ps *ProductService) FilterProducts(ctx context.Context, filter ProductFilter, projection []string) Response {
	// TODO: Implement complex filtering with field projection
	// Steps:
	// 1. Build comprehensive filter from all ProductFilter fields
	// 2. Handle array fields (tags) with $in or $all operators
	// 3. Apply field projection to limit returned data
	// 4. Combine multiple conditions with $and operator

	// Hint: Use bson.M{"$and": []bson.M{...}} for multiple conditions
	// Hint: Use $in for tags: bson.M{"tags": bson.M{"$in": filter.Tags}}
	// Hint: Build projection: bson.M{"field1": 1, "field2": 1, "_id": 0}

	return Response{
		Success: false,
		Error:   "FilterProducts not implemented",
		Code:    500,
	}
}

// GetProductsByTags retrieves products that have any of the specified tags
func (ps *ProductService) GetProductsByTags(ctx context.Context, tags []string, matchAll bool) Response {
	// TODO: Implement tag-based filtering
	// Steps:
	// 1. Use $in operator for "any tag" matching
	// 2. Use $all operator for "all tags" matching
	// 3. Handle empty tags array
	// 4. Return products with tag information

	// Hint: $in matches any: bson.M{"tags": bson.M{"$in": tags}}
	// Hint: $all matches all: bson.M{"tags": bson.M{"$all": tags}}

	return Response{
		Success: false,
		Error:   "GetProductsByTags not implemented",
		Code:    500,
	}
}

// GetTopRatedProducts retrieves highest rated products in category
func (ps *ProductService) GetTopRatedProducts(ctx context.Context, category string, limit int) Response {
	// TODO: Implement top-rated products query
	// Steps:
	// 1. Filter by category if specified
	// 2. Sort by rating in descending order
	// 3. Apply limit for top N products
	// 4. Only include products with ratings > 0

	// Hint: Use bson.M{"rating": bson.M{"$gt": 0}} to exclude unrated
	// Hint: Sort by rating descending: bson.M{"rating": -1}

	return Response{
		Success: false,
		Error:   "GetTopRatedProducts not implemented",
		Code:    500,
	}
}

// CountProductsByCategory counts products in each category
func (ps *ProductService) CountProductsByCategory(ctx context.Context) Response {
	// TODO: Implement category counting
	// Steps:
	// 1. Use aggregation pipeline with $group stage
	// 2. Group by category and count documents
	// 3. Sort results by count or category name
	// 4. Return category counts as map or array

	// Hint: This is a preview of aggregation (Challenge 3)
	// Hint: For now, you can use multiple CountDocuments calls
	// Hint: Return as map[string]int64 for category counts

	return Response{
		Success: false,
		Error:   "CountProductsByCategory not implemented",
		Code:    500,
	}
}

// ConnectMongoDB establishes connection to MongoDB
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// SeedSampleProducts creates sample product data for testing
func (ps *ProductService) SeedSampleProducts(ctx context.Context) error {
	// Sample products for testing
	products := []Product{
		{
			ID:          primitive.NewObjectID(),
			Name:        "iPhone 15 Pro",
			Description: "Latest Apple smartphone with advanced features",
			Category:    "Electronics",
			Price:       999.99,
			Stock:       50,
			Tags:        []string{"smartphone", "apple", "premium"},
			Brand:       "Apple",
			Rating:      4.8,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Samsung Galaxy S24",
			Description: "High-performance Android smartphone",
			Category:    "Electronics",
			Price:       899.99,
			Stock:       75,
			Tags:        []string{"smartphone", "samsung", "android"},
			Brand:       "Samsung",
			Rating:      4.6,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Nike Air Max 270",
			Description: "Comfortable running shoes with air cushioning",
			Category:    "Footwear",
			Price:       129.99,
			Stock:       100,
			Tags:        []string{"shoes", "running", "nike", "comfort"},
			Brand:       "Nike",
			Rating:      4.4,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Adidas Ultraboost 22",
			Description: "Premium running shoes with boost technology",
			Category:    "Footwear",
			Price:       189.99,
			Stock:       60,
			Tags:        []string{"shoes", "running", "adidas", "premium"},
			Brand:       "Adidas",
			Rating:      4.7,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "MacBook Pro M3",
			Description: "Professional laptop with M3 chip",
			Category:    "Electronics",
			Price:       1999.99,
			Stock:       25,
			Tags:        []string{"laptop", "apple", "professional", "m3"},
			Brand:       "Apple",
			Rating:      4.9,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Convert to interface{} slice for InsertMany
	documents := make([]interface{}, len(products))
	for i, product := range products {
		documents[i] = product
	}

	_, err := ps.Collection.InsertMany(ctx, documents)
	return err
}

// Example usage and testing
func main() {
	mongoURI := "mongodb://localhost:27017"

	client, err := ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("product_catalog").Collection("products")
	productService := &ProductService{Collection: collection}

	ctx := context.Background()

	// Clear existing data and seed sample products
	collection.Drop(ctx)
	if err := productService.SeedSampleProducts(ctx); err != nil {
		log.Printf("Failed to seed sample data: %v", err)
	} else {
		fmt.Println("Sample products seeded successfully!")
	}

	// Example queries (uncomment after implementation)

	// Get products by category
	// fmt.Println("\n=== Products in Electronics ===")
	// resp := productService.GetProductsByCategory(ctx, "Electronics", ProductFilter{})
	// fmt.Printf("Response: %+v\n", resp)

	// Search products by name
	// fmt.Println("\n=== Search for 'iPhone' ===")
	// resp = productService.SearchProductsByName(ctx, "iPhone", false)
	// fmt.Printf("Response: %+v\n", resp)

	// Get products by price range
	// fmt.Println("\n=== Products under $200 ===")
	// resp = productService.GetProductsByPriceRange(ctx, 0, 200)
	// fmt.Printf("Response: %+v\n", resp)
}

/*
IMPLEMENTATION CHECKLIST:

Query Operators:
[ ] Use $gte and $lte for price range filtering
[ ] Use $in operator for array field matching (tags)
[ ] Use $regex for text search with case-insensitive option
[ ] Use $and and $or for complex query combinations
[ ] Use $gt for minimum rating filtering

Sorting and Pagination:
[ ] Implement sorting with options.Find().SetSort()
[ ] Calculate skip value for pagination: (page - 1) * limit
[ ] Use SetSkip() and SetLimit() for pagination
[ ] Count total documents with CountDocuments()
[ ] Return pagination metadata (total, pages, etc.)

Field Projection:
[ ] Use SetProjection() to limit returned fields
[ ] Build projection document with field inclusion/exclusion
[ ] Optimize data transfer by selecting only needed fields

Error Handling:
[ ] Validate input parameters (negative prices, empty search terms)
[ ] Handle empty result sets gracefully
[ ] Return appropriate error messages and status codes
[ ] Handle MongoDB query errors properly

Performance Optimization:
[ ] Use indexes for frequently queried fields (bonus)
[ ] Limit result sets to prevent memory issues
[ ] Use projection to reduce data transfer
[ ] Consider using aggregation for complex queries

Response Structure:
[ ] Use consistent Response struct for all operations
[ ] Include appropriate success/error status
[ ] Return meaningful data and messages
[ ] Handle pagination metadata properly

Testing:
[ ] All tests should pass
[ ] Handle edge cases (empty filters, invalid ranges)
[ ] Test pagination boundary conditions
[ ] Verify sorting works correctly
[ ] Test complex filter combinations

BONUS FEATURES (Optional):
[ ] Implement full-text search with text indexes
[ ] Add geospatial queries for location-based filtering
[ ] Create search suggestions and autocomplete
[ ] Implement faceted search with category counts
[ ] Add query performance monitoring
[ ] Optimize with proper indexing strategies
*/
