# Challenge 3: JSON API with Validation & Error Handling

Build a **Product Catalog API** with comprehensive input validation, custom validators, and robust error handling.

## Challenge Requirements

Implement a JSON API with the following endpoints:

- `GET /health` - Health check endpoint
- `POST /products` - Create new product with validation
- `PUT /products/:id` - Update product with validation  
- `POST /products/bulk` - Create multiple products in one request
- `GET /products` - Get all products with optional filtering
- `GET /products/:id` - Get product by ID
- `DELETE /products/:id` - Delete product

## Data Structure

```go
type Product struct {
    ID            string   `json:"id"`
    Name          string   `json:"name" validate:"required,min=2,max=100"`
    Description   string   `json:"description"`
    SKU           string   `json:"sku" validate:"required,sku_format"`
    Price         float64  `json:"price" validate:"required,min=0.01"`
    Category      string   `json:"category" validate:"required,valid_category"`
    InStock       bool     `json:"in_stock"`
    StockQuantity int      `json:"stock_quantity" validate:"min=0"`
    Tags          []string `json:"tags"`
    CreatedAt     string   `json:"created_at"`
    UpdatedAt     string   `json:"updated_at"`
}

type ProductResponse struct {
    Status  string   `json:"status"`
    Message string   `json:"message,omitempty"`
    Product *Product `json:"product,omitempty"`
}

type ProductListResponse struct {
    Status   string    `json:"status"`
    Count    int       `json:"count"`
    Products []Product `json:"products"`
}

type BulkCreateRequest struct {
    Products []Product `json:"products" validate:"required,min=1,max=10"`
}

type BulkCreateResponse struct {
    Status    string            `json:"status"`
    Message   string            `json:"message"`
    Succeeded []Product         `json:"succeeded"`
    Failed    []BulkCreateError `json:"failed"`
}

type BulkCreateError struct {
    Index   int     `json:"index"`
    Product Product `json:"product"`
    Error   string  `json:"error"`
}

type ValidationErrorResponse struct {
    Status  string            `json:"status"`
    Message string            `json:"message"`
    Errors  map[string]string `json:"errors"`
}

type ErrorResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}
```

## Validation Requirements

### Built-in Validators
- **required**: Field must be present
- **min/max**: String length or numeric range validation
- **min=0.01**: Price must be greater than 0

### Custom Validators
- **sku_format**: SKU format validation (XXX-NNNNN, e.g., "ABC-12345")
- **valid_category**: Category must be one of: electronics, clothing, books, home, sports

## Request/Response Examples

**POST /products** (Request body)
```json
{
    "name": "Wireless Headphones",
    "description": "High-quality wireless headphones with noise cancellation",
    "sku": "ELE-12345",
    "price": 199.99,
    "category": "electronics",
    "in_stock": true,
    "stock_quantity": 50,
    "tags": ["wireless", "audio", "electronics"]
}
```

**Validation Error Response**
```json
{
    "status": "error",
    "message": "Validation failed",
    "errors": {
        "name": "name is required",
        "sku": "sku must match format XXX-NNNNN",
        "price": "price must be at least 0.01"
    }
}
```

## Testing Requirements

Your solution must pass tests for:
- Health check returns proper status
- Product creation with valid data succeeds
- Validation errors return proper error format
- SKU format validation works correctly
- Category validation accepts only valid categories
- Bulk creation handles partial success/failure
- Product filtering by category and stock status
- All CRUD operations work correctly
- Proper HTTP status codes for all scenarios