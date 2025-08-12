# Learning: Input Validation & Error Handling

## 🌟 **What is Input Validation?**

Input validation ensures that data received by your API meets specific criteria before processing. It's your first line of defense against bad data, security vulnerabilities, and application crashes.

### **Why Validate Input?**
- **Security**: Prevent injection attacks and malicious data
- **Data Integrity**: Ensure data meets business requirements
- **User Experience**: Provide clear feedback on invalid input
- **System Stability**: Prevent crashes from unexpected data

## 🛠️ **Validation in Fiber**

Fiber doesn't include built-in validation, but integrates well with the `validator` package:

```go
import "github.com/go-playground/validator/v10"

type User struct {
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=0,lte=120"`
}
```

## 📏 **Built-in Validation Tags**

### **Required Fields**
```go
type Product struct {
    Name string `validate:"required"`        // Must be present
    SKU  string `validate:"required,min=5"`  // Required and min length
}
```

### **String Validation**
```go
type User struct {
    Username string `validate:"min=3,max=20"`           // Length constraints
    Email    string `validate:"email"`                  // Email format
    Website  string `validate:"url"`                    // URL format
    Phone    string `validate:"e164"`                   // Phone number format
}
```

### **Numeric Validation**
```go
type Product struct {
    Price    float64 `validate:"gt=0"`          // Greater than 0
    Quantity int     `validate:"gte=0,lte=1000"` // Range: 0-1000
    Rating   float64 `validate:"min=1,max=5"`   // Rating scale
}
```

### **Enum Validation**
```go
type Product struct {
    Category string `validate:"oneof=electronics clothing books home"`
    Status   string `validate:"oneof=active inactive pending"`
}
```

## 🔧 **Custom Validators**

Create validators for business-specific rules:

```go
func validateSKU(fl validator.FieldLevel) bool {
    sku := fl.Field().String()
    // Custom SKU format: PROD-12345
    matched, _ := regexp.MatchString(`^PROD-\d{5}$`, sku)
    return matched
}

// Register custom validator
validate := validator.New()
validate.RegisterValidation("sku", validateSKU)
```

### **Complex Custom Validators**
```go
func validateBusinessHours(fl validator.FieldLevel) bool {
    hour := fl.Field().Int()
    // Business hours: 9 AM to 5 PM
    return hour >= 9 && hour <= 17
}

func validatePasswordStrength(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    
    // Check length
    if len(password) < 8 {
        return false
    }
    
    // Check for uppercase, lowercase, digit, special char
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`\d`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
    
    return hasUpper && hasLower && hasDigit && hasSpecial
}
```

## 🚨 **Error Handling Patterns**

### **Structured Error Responses**
```go
type ValidationError struct {
    Field   string      `json:"field"`
    Tag     string      `json:"tag"`
    Value   interface{} `json:"value"`
    Message string      `json:"message"`
}

type ErrorResponse struct {
    Success bool              `json:"success"`
    Error   string            `json:"error"`
    Details []ValidationError `json:"details,omitempty"`
}
```

### **Converting Validator Errors**
```go
func formatValidationErrors(err error) []ValidationError {
    var errors []ValidationError
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            errors = append(errors, ValidationError{
                Field:   e.Field(),
                Tag:     e.Tag(),
                Value:   e.Value(),
                Message: getErrorMessage(e),
            })
        }
    }
    
    return errors
}

func getErrorMessage(e validator.FieldError) string {
    switch e.Tag() {
    case "required":
        return e.Field() + " is required"
    case "email":
        return e.Field() + " must be a valid email"
    case "min":
        return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
    case "max":
        return fmt.Sprintf("%s cannot exceed %s characters", e.Field(), e.Param())
    default:
        return e.Field() + " is invalid"
    }
}
```

## 🔍 **Advanced Validation Techniques**

### **Cross-Field Validation**
```go
type User struct {
    Password        string `validate:"required,min=8"`
    ConfirmPassword string `validate:"required,eqfield=Password"`
}

type DateRange struct {
    StartDate time.Time `validate:"required"`
    EndDate   time.Time `validate:"required,gtfield=StartDate"`
}
```

### **Conditional Validation**
```go
type Product struct {
    Type        string  `validate:"required,oneof=physical digital"`
    Weight      float64 `validate:"required_if=Type physical"`
    FileSize    int64   `validate:"required_if=Type digital"`
    ShippingCost float64 `validate:"excluded_if=Type digital"`
}
```

### **Slice Validation**
```go
type Order struct {
    Items []Item `validate:"required,dive,required"`
    Tags  []string `validate:"dive,min=2,max=20"`
}

type Item struct {
    Name     string  `validate:"required"`
    Price    float64 `validate:"gt=0"`
    Quantity int     `validate:"gte=1"`
}
```

## 📊 **Filtering and Search Patterns**

### **Query Parameter Filtering**
```go
func applyFilters(products []Product, c *fiber.Ctx) []Product {
    var filtered []Product
    
    // Get filter parameters
    category := c.Query("category")
    minPrice := c.Query("min_price")
    maxPrice := c.Query("max_price")
    inStock := c.Query("in_stock")
    search := c.Query("search")
    
    for _, product := range products {
        // Apply filters
        if category != "" && product.Category != category {
            continue
        }
        
        if minPrice != "" {
            if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
                if product.Price < min {
                    continue
                }
            }
        }
        
        if search != "" {
            searchLower := strings.ToLower(search)
            nameMatch := strings.Contains(strings.ToLower(product.Name), searchLower)
            descMatch := strings.Contains(strings.ToLower(product.Description), searchLower)
            if !nameMatch && !descMatch {
                continue
            }
        }
        
        filtered = append(filtered, product)
    }
    
    return filtered
}
```

### **Pagination Support**
```go
func paginateResults(items []Product, c *fiber.Ctx) ([]Product, map[string]interface{}) {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    
    offset := (page - 1) * limit
    end := offset + limit
    
    if offset >= len(items) {
        return []Product{}, map[string]interface{}{
            "page":       page,
            "limit":      limit,
            "total":      len(items),
            "total_pages": (len(items) + limit - 1) / limit,
        }
    }
    
    if end > len(items) {
        end = len(items)
    }
    
    return items[offset:end], map[string]interface{}{
        "page":        page,
        "limit":       limit,
        "total":       len(items),
        "total_pages": (len(items) + limit - 1) / limit,
    }
}
```

## 🔒 **Security Considerations**

### **Input Sanitization**
```go
import "html"

func sanitizeInput(input string) string {
    // Remove potentially dangerous characters
    input = html.EscapeString(input)
    input = strings.TrimSpace(input)
    return input
}
```

### **Rate Limiting Validation**
```go
func validateRequestRate(c *fiber.Ctx) error {
    // Limit validation requests to prevent abuse
    // Implementation depends on your rate limiting strategy
    return nil
}
```

## 🧪 **Testing Validation**

### **Unit Testing Validators**
```go
func TestProductValidation(t *testing.T) {
    validate := validator.New()
    
    tests := []struct {
        name    string
        product Product
        wantErr bool
    }{
        {
            name: "valid product",
            product: Product{
                Name:        "Test Product",
                Description: "A test product description",
                Price:       99.99,
                Category:    "electronics",
                SKU:         "PROD-12345",
            },
            wantErr: false,
        },
        {
            name: "invalid price",
            product: Product{
                Name:        "Test Product",
                Description: "A test product description",
                Price:       -10.00,
                Category:    "electronics",
                SKU:         "PROD-12345",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate.Struct(tt.product)
            hasErr := err != nil
            assert.Equal(t, tt.wantErr, hasErr)
        })
    }
}
```

## 🎯 **Best Practices**

1. **Validate Early**: Check input as soon as it enters your application
2. **Clear Messages**: Provide specific, actionable error messages
3. **Consistent Format**: Use standard error response formats
4. **Security First**: Sanitize input and validate against business rules
5. **Performance**: Cache validators and avoid expensive operations
6. **Documentation**: Document validation rules in API documentation
7. **Testing**: Test both valid and invalid input scenarios

## 📚 **Next Steps**

After mastering validation and error handling:
1. **Authentication & Authorization** - Secure your APIs
2. **Database Integration** - Persist validated data
3. **API Documentation** - Document validation rules
4. **Advanced Patterns** - Async validation, custom middleware