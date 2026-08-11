# Hints for Challenge 3: Validation & Error Handling

## Hint 1: Setting Up Validator

Initialize the validator package and register custom validators:

```go
import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
    validate = validator.New()
    validate.RegisterValidation("sku_format", validateSKUFormat)
    validate.RegisterValidation("valid_category", validateCategory)
}
```

## Hint 2: Struct Validation Tags

Use validation tags in your Product struct:

```go
type Product struct {
    Name          string  `json:"name" validate:"required,min=2,max=100"`
    SKU           string  `json:"sku" validate:"required,sku_format"`
    Price         float64 `json:"price" validate:"required,min=0.01"`
    Category      string  `json:"category" validate:"required,valid_category"`
    StockQuantity int     `json:"stock_quantity" validate:"required,min=0"`
}
```

## Hint 3: Custom SKU Validator

Implement SKU format validation (XXX-NNNNN):

```go
import (
    "regexp"
    "github.com/go-playground/validator/v10"
)

func validateSKUFormat(fl validator.FieldLevel) bool {
    sku := fl.Field().String()
    
    // Use regex to match pattern: 3 uppercase letters, dash, 5 digits
    matched, _ := regexp.MatchString(`^[A-Z]{3}-[0-9]{5}$`, sku)
    return matched
}
```

## Hint 4: Category Validator

Validate against allowed categories:

```go
var validCategories = map[string]bool{
    "electronics": true,
    "clothing":    true,
    "books":       true,
    "home":        true,
    "sports":      true,
}

func validateCategory(fl validator.FieldLevel) bool {
    category := fl.Field().String()
    return validCategories[category]
}
```

## Hint 5: Validation Error Formatting

Convert validation errors to user-friendly messages:

```go
func formatValidationErrors(err error) map[string]string {
    errors := make(map[string]string)
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, fieldError := range validationErrors {
            field := fieldError.Field()
            tag := fieldError.Tag()
            
            switch tag {
            case "required":
                errors[field] = fmt.Sprintf("%s is required", field)
            case "min":
                errors[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
            case "sku_format":
                errors[field] = "SKU must be in format XXX-NNNNN (e.g., ABC-12345)"
            case "valid_category":
                errors[field] = "Invalid category. Must be one of: electronics, clothing, books, home, sports"
            default:
                errors[field] = fmt.Sprintf("%s is invalid", field)
            }
        }
    }
    
    return errors
}
```

## Hint 6: Product Creation with Validation

Implement create product with proper validation:

```go
func createProduct(c echo.Context) error {
    var product Product
    
    if err := c.Bind(&product); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  "error",
            Message: "Invalid JSON format",
        })
    }
    
    if err := validate.Struct(&product); err != nil {
        errors := formatValidationErrors(err)
        return c.JSON(http.StatusBadRequest, ValidationErrorResponse{
            Status:  "error",
            Message: "Validation failed",
            Errors:  errors,
        })
    }
    
    // Generate ID and timestamps
    product.ID = uuid.New().String()
    product.CreatedAt = time.Now().Format(time.RFC3339)
    product.UpdatedAt = product.CreatedAt
    
    products = append(products, product)
    
    return c.JSON(http.StatusCreated, ProductResponse{
        Status:  "success",
        Message: "Product created successfully",
        Product: &product,
    })
}
```

## Hint 7: Query Parameter Filtering

Handle multiple filter parameters:

```go
func getAllProducts(c echo.Context) error {
    category := c.QueryParam("category")
    inStockStr := c.QueryParam("in_stock")
    minPriceStr := c.QueryParam("min_price")
    maxPriceStr := c.QueryParam("max_price")
    search := c.QueryParam("search")

    filteredProducts := make([]Product, 0)

    for _, product := range products {
        // Filter by category
        if category != "" && product.Category != category {
            continue
        }
        
        // Filter by stock status
        if inStockStr != "" {
            inStock, err := strconv.ParseBool(inStockStr)
            if err == nil && product.InStock != inStock {
                continue
            }
        }
        
        // Add other filters...
        filteredProducts = append(filteredProducts, product)
    }
    
    return c.JSON(http.StatusOK, ProductListResponse{
        Status:   "success",
        Count:    len(filteredProducts),
        Products: filteredProducts,
    })
}
```

## Hint 8: Bulk Operations Structure

Set up bulk create request and response:

```go
type BulkCreateRequest struct {
    Products []Product `json:"products" validate:"required,min=1,max=10,dive"`
}

type BulkCreateResponse struct {
    Status    string             `json:"status"`
    Message   string             `json:"message"`
    Succeeded []Product          `json:"succeeded"`
    Failed    []BulkCreateError  `json:"failed"`
}

type BulkCreateError struct {
    Index   int     `json:"index"`
    Product Product `json:"product"`
    Error   string  `json:"error"`
}
```

## Hint 9: Bulk Create Implementation

Process bulk operations with partial success handling:

```go
func bulkCreateProducts(c echo.Context) error {
    var request BulkCreateRequest
    
    if err := c.Bind(&request); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  "error",
            Message: "Invalid JSON format",
        })
    }
    
    if err := validate.Struct(&request); err != nil {
        errors := formatValidationErrors(err)
        return c.JSON(http.StatusBadRequest, ValidationErrorResponse{
            Status:  "error",
            Message: "Validation failed",
            Errors:  errors,
        })
    }
    
    var succeeded []Product
    var failed []BulkCreateError
    
    for i, product := range request.Products {
        // Validate individual product
        if err := validate.Struct(&product); err != nil {
            errorMsg := formatValidationErrors(err)
            failed = append(failed, BulkCreateError{
                Index:   i,
                Product: product,
                Error:   fmt.Sprintf("Validation failed: %v", errorMsg),
            })
            continue
        }
        
        // Process successful product
        product.ID = uuid.New().String()
        product.CreatedAt = time.Now().Format(time.RFC3339)
        product.UpdatedAt = product.CreatedAt
        
        products = append(products, product)
        succeeded = append(succeeded, product)
    }
    
    status := "success"
    if len(failed) > 0 {
        status = "partial_success"
    }
    
    return c.JSON(http.StatusOK, BulkCreateResponse{
        Status:    status,
        Message:   fmt.Sprintf("Processed %d products: %d succeeded, %d failed", 
                              len(request.Products), len(succeeded), len(failed)),
        Succeeded: succeeded,
        Failed:    failed,
    })
}
```

## Hint 10: Debugging Tips

Common validation issues and solutions:

- **Validation not working**: Make sure you've registered custom validators with `validate.RegisterValidation()`
- **SKU format failing**: Check the regex pattern `^[A-Z]{3}-[0-9]{5}$` - requires exactly 3 uppercase letters, dash, 5 digits
- **Category validation failing**: Ensure the category exists in your `validCategories` map
- **Bulk operations errors**: Remember to validate each product individually and collect both successes and failures