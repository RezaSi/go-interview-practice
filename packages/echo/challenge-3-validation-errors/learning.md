# Learning Materials: JSON API with Validation & Error Handling

## 🎯 **What You'll Learn**

This challenge teaches you advanced validation patterns and error handling techniques that are essential for building robust APIs in production environments.

## 📚 **Core Concepts**

### **1. Input Validation Layers**

Modern APIs use multiple validation layers:

```go
// Layer 1: Basic JSON binding validation
var product Product
if err := c.Bind(&product); err != nil {
    return c.JSON(http.StatusBadRequest, ErrorResponse{
        Status:  "error",
        Message: "Invalid JSON format",
    })
}

// Layer 2: Struct tag validation
if err := validate.Struct(&product); err != nil {
    errors := formatValidationErrors(err)
    return c.JSON(http.StatusBadRequest, ValidationErrorResponse{
        Status:  "error",
        Message: "Validation failed",
        Errors:  errors,
    })
}

// Layer 3: Custom business logic validation
if !isValidBusinessRules(&product) {
    return c.JSON(http.StatusBadRequest, ErrorResponse{
        Status:  "error",
        Message: "Business rules validation failed",
    })
}
```

### **2. Validation Tags**

Echo works with the `go-playground/validator` package for struct validation:

```go
type Product struct {
    Name     string  `json:"name" validate:"required,min=2,max=100"`
    Price    float64 `json:"price" validate:"required,min=0.01"`
    Category string  `json:"category" validate:"required,valid_category"`
    SKU      string  `json:"sku" validate:"required,sku_format"`
}
```

**Common Validation Tags:**
- `required`: Field must be present and non-zero
- `min=N, max=N`: String length or numeric range
- `email`: Valid email format
- `url`: Valid URL format
- `oneof=a b c`: Value must be one of specified options

### **3. Custom Validators**

Create custom validation functions for business-specific rules:

```go
// Register custom validator
validate.RegisterValidation("sku_format", validateSKUFormat)

func validateSKUFormat(fl validator.FieldLevel) bool {
    sku := fl.Field().String()
    matched, _ := regexp.MatchString(`^[A-Z]{3}-[0-9]{5}$`, sku)
    return matched
}

func validateCategory(fl validator.FieldLevel) bool {
    category := fl.Field().String()
    validCategories := map[string]bool{
        "electronics": true,
        "clothing":    true,
        "books":       true,
        "home":        true,
        "sports":      true,
    }
    return validCategories[category]
}
```

### **4. Error Response Formatting**

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
                errors[field] = "SKU must match format XXX-NNNNN"
            case "valid_category":
                errors[field] = "Category must be one of: electronics, clothing, books, home, sports"
            default:
                errors[field] = fmt.Sprintf("%s is invalid", field)
            }
        }
    }
    
    return errors
}
```

## 🔧 **Error Handling Patterns**

### **Consistent Error Responses**

Define standard error response structures:

```go
type ErrorResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

type ValidationErrorResponse struct {
    Status  string            `json:"status"`
    Message string            `json:"message"`
    Errors  map[string]string `json:"errors"`
}
```

### **HTTP Status Codes**

Use appropriate status codes for different error types:

- **400 Bad Request**: Invalid input, validation errors
- **404 Not Found**: Resource doesn't exist
- **409 Conflict**: Duplicate resource (e.g., SKU already exists)
- **422 Unprocessable Entity**: Valid JSON but business logic errors
- **500 Internal Server Error**: Server-side errors

### **Bulk Operations Error Handling**

Handle partial success in bulk operations:

```go
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

func bulkCreateProducts(c echo.Context) error {
    var request BulkCreateRequest
    if err := c.Bind(&request); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  "error",
            Message: "Invalid JSON format",
        })
    }
    
    var succeeded []Product
    var failed []BulkCreateError
    
    for i, product := range request.Products {
        if err := validate.Struct(&product); err != nil {
            failed = append(failed, BulkCreateError{
                Index:   i,
                Product: product,
                Error:   "Validation failed",
            })
            continue
        }
        
        // Process successful product
        succeeded = append(succeeded, product)
    }
    
    status := "success"
    if len(failed) > 0 {
        status = "partial_success"
    }
    
    return c.JSON(http.StatusOK, BulkCreateResponse{
        Status:    status,
        Message:   fmt.Sprintf("Processed %d products", len(request.Products)),
        Succeeded: succeeded,
        Failed:    failed,
    })
}
```

## 🎯 **Best Practices**

### **Validation Strategy**
1. **Fail Fast**: Validate input as early as possible
2. **Clear Messages**: Provide specific, actionable error messages
3. **Consistent Format**: Use the same error response structure
4. **Field-Level Errors**: Return errors for specific fields when possible

### **Security Considerations**
- **Input Sanitization**: Clean user input to prevent injection attacks
- **Rate Limiting**: Prevent abuse of validation endpoints
- **Error Information**: Don't expose internal system details in errors

### **Performance Tips**
- **Validate Once**: Don't re-validate the same data multiple times
- **Batch Validation**: Use bulk operations for multiple items
- **Cache Validators**: Reuse validator instances when possible

## 🧪 **Testing Validation**

### **Unit Tests for Validators**
```go
func TestSKUValidation(t *testing.T) {
    tests := []struct {
        sku   string
        valid bool
    }{
        {"ABC-12345", true},
        {"XYZ-99999", true},
        {"abc-12345", false}, // lowercase
        {"AB-12345", false},  // too short
        {"ABCD-12345", false}, // too long
    }
    
    for _, test := range tests {
        result := validateSKUFormat(test.sku)
        assert.Equal(t, test.valid, result)
    }
}
```

### **Integration Tests**
```go
func TestProductValidation(t *testing.T) {
    e := echo.New()
    
    // Test invalid product
    invalidProduct := `{"name":"","price":-1,"category":"invalid"}`
    req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(invalidProduct))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    
    c := e.NewContext(req, rec)
    err := createProduct(c)
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusBadRequest, rec.Code)
    
    var response ValidationErrorResponse
    json.Unmarshal(rec.Body.Bytes(), &response)
    assert.Equal(t, "error", response.Status)
    assert.Contains(t, response.Errors, "name")
    assert.Contains(t, response.Errors, "price")
}
```

Proper validation and error handling are crucial for building reliable, user-friendly APIs that can handle real-world usage patterns and edge cases!
