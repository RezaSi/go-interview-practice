# Challenge 3: JSON API with Validation & Error Handling

Build a **Product Catalog API** with comprehensive input validation, custom validators, and robust error handling.

## Challenge Requirements

Implement a JSON API with the following endpoints:

- `POST /products` - Create new product with validation
- `PUT /products/:id` - Update product with validation  
- `POST /products/bulk` - Create multiple products in one request
- `GET /products` - Get all products with optional filtering
- `GET /products/:id` - Get product by ID

## Data Structure

```go
type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name" binding:"required,min=2,max=100"`
    Description string  `json:"description" binding:"required,min=10,max=500"`
    Price       float64 `json:"price" binding:"required,gt=0"`
    Category    string  `json:"category" binding:"required,oneof=electronics clothing books home"`
    SKU         string  `json:"sku" binding:"required,sku"`
    InStock     bool    `json:"in_stock"`
    Tags        []string `json:"tags" binding:"dive,min=2,max=20"`
}
```

## Validation Requirements

### Built-in Validators
- **required**: Field must be present
- **min/max**: String length or numeric range validation
- **gt**: Greater than (price > 0)
- **oneof**: Value must be one of specified options

### Custom Validators
- **sku**: SKU format validation (e.g., "PROD-12345")
- **dive**: Validate each element in slices

## Error Response Format

```json
{
    "success": false,
    "error": "Validation failed",
    "details": [
        {
            "field": "name",
            "message": "Name is required and must be 2-100 characters"
        },
        {
            "field": "price", 
            "message": "Price must be greater than 0"
        }
    ]
}
```

## Testing Requirements

Your solution must pass tests for:
- Field presence validation (required fields)
- String length validation (min/max)
- Numeric validation (price > 0)
- Enum validation (category values)
- Custom SKU format validation
- Array element validation (tags)
- Bulk operation validation with partial failures
- Proper error response formatting
- Content-Type validation (application/json required) 