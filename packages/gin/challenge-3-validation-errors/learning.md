# Learning Materials: JSON API with Validation & Error Handling

## 🎯 **What You'll Learn**

This challenge teaches you advanced validation patterns and error handling techniques that are essential for building robust APIs in production environments.

## 📚 **Core Concepts**

### **1. Input Validation Layers**

Modern APIs use multiple validation layers:

```go
// Layer 1: Basic JSON binding validation (built-in)
if err := c.ShouldBindJSON(&product); err != nil {
    // Handle basic format errors
}

// Layer 2: Custom business logic validation
errors := validateProduct(&product)

// Layer 3: Cross-field and contextual validation
errors = append(errors, validateBusinessRules(&product)...)
```

### **2. Custom Validation Functions**

Create reusable validation functions for complex business rules:

```go
// Regular expression validation
func isValidSKU(sku string) bool {
    pattern := `^[A-Z]{3}-\d{3}-[A-Z]{3}$`
    matched, _ := regexp.MatchString(pattern, sku)
    return matched
}

// List-based validation
func isValidCurrency(currency string) bool {
    validCurrencies := []string{"USD", "EUR", "GBP", "JPY"}
    for _, valid := range validCurrencies {
        if currency == valid {
            return true
        }
    }
    return false
}

// Cross-field validation
func validateInventoryRules(inventory Inventory) []ValidationError {
    var errors []ValidationError
    
    if inventory.Reserved > inventory.Quantity {
        errors = append(errors, ValidationError{
            Field:   "inventory.reserved",
            Message: "Reserved cannot exceed quantity",
        })
    }
    
    return errors
}
```

### **3. Error Response Standardization**

Consistent error responses improve API usability:

```go
type ValidationError struct {
    Field   string      `json:"field"`
    Value   interface{} `json:"value"`
    Tag     string      `json:"tag"`
    Message string      `json:"message"`
    Param   string      `json:"param,omitempty"`
}

type APIResponse struct {
    Success   bool              `json:"success"`
    Data      interface{}       `json:"data,omitempty"`
    Message   string            `json:"message,omitempty"`
    Errors    []ValidationError `json:"errors,omitempty"`
    ErrorCode string            `json:"error_code,omitempty"`
    RequestID string            `json:"request_id,omitempty"`
}
```

## 🔧 **Input Sanitization Patterns**

Always sanitize input before validation:

```go
func sanitizeProduct(product *Product) {
    // Remove leading/trailing whitespace
    product.Name = strings.TrimSpace(product.Name)
    product.SKU = strings.TrimSpace(product.SKU)
    
    // Normalize case
    product.Currency = strings.ToUpper(product.Currency)
    product.Category.Slug = strings.ToLower(product.Category.Slug)
    
    // Calculate computed fields
    product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
    
    // Set system fields
    if product.ID == 0 {
        product.CreatedAt = time.Now()
    }
    product.UpdatedAt = time.Now()
}
```

## 🏗️ **Regular Expression Patterns**

Common validation patterns:

```go
// SKU format: ABC-123-XYZ
var skuPattern = `^[A-Z]{3}-\d{3}-[A-Z]{3}$`

// URL-friendly slug: lowercase-with-hyphens
var slugPattern = `^[a-z0-9]+(?:-[a-z0-9]+)*$`

// Warehouse code: WH001, WH002, etc.
var warehousePattern = `^WH\d{3}$`

// Email validation (basic)
var emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

// Phone number (US format)
var phonePattern = `^\(\d{3}\) \d{3}-\d{4}$`
```

## 📊 **Bulk Operations Best Practices**

Handle bulk operations with detailed feedback:

```go
type BulkResult struct {
    Index   int               `json:"index"`
    Success bool              `json:"success"`
    Product *Product          `json:"product,omitempty"`
    Errors  []ValidationError `json:"errors,omitempty"`
}

func processBulkData(items []Product) []BulkResult {
    var results []BulkResult
    
    for i, item := range items {
        errors := validateProduct(&item)
        
        if len(errors) > 0 {
            results = append(results, BulkResult{
                Index:   i,
                Success: false,
                Errors:  errors,
            })
        } else {
            // Process successful item
            sanitizeProduct(&item)
            // Save to database/storage
            
            results = append(results, BulkResult{
                Index:   i,
                Success: true,
                Product: &item,
            })
        }
    }
    
    return results
}
```

## 🌐 **Localization and Error Messages**

Support multiple languages for error messages:

```go
var ErrorMessages = map[string]map[string]string{
    "en": {
        "required":      "This field is required",
        "min":           "Value must be at least %s",
        "max":           "Value must be at most %s",
        "sku_format":    "SKU must follow ABC-123-XYZ format",
        "invalid_currency": "Must be a valid currency code",
    },
    "es": {
        "required":      "Este campo es obligatorio",
        "min":           "El valor debe ser al menos %s",
        "sku_format":    "SKU debe seguir el formato ABC-123-XYZ",
    },
}

func getLocalizedMessage(lang, key string, params ...string) string {
    messages, exists := ErrorMessages[lang]
    if !exists {
        messages = ErrorMessages["en"] // Fallback to English
    }
    
    message, exists := messages[key]
    if !exists {
        return "Validation failed"
    }
    
    // Handle parameter substitution
    for i, param := range params {
        placeholder := fmt.Sprintf("%%s")
        if i == 0 {
            message = strings.Replace(message, placeholder, param, 1)
        }
    }
    
    return message
}
```

## 🔐 **Security Considerations**

### **1. Input Sanitization**
```go
// Remove dangerous characters
func sanitizeString(input string) string {
    // Remove HTML tags
    input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")
    
    // Remove SQL injection attempts
    dangerous := []string{"'", "\"", ";", "--", "/*", "*/"}
    for _, char := range dangerous {
        input = strings.ReplaceAll(input, char, "")
    }
    
    return strings.TrimSpace(input)
}
```

### **2. Rate Limiting for Bulk Operations**
```go
const maxBulkSize = 100

func createProductsBulk(c *gin.Context) {
    var products []Product
    
    if err := c.ShouldBindJSON(&products); err != nil {
        c.JSON(400, APIResponse{
            Success: false,
            Message: "Invalid JSON",
        })
        return
    }
    
    if len(products) > maxBulkSize {
        c.JSON(400, APIResponse{
            Success: false,
            Message: fmt.Sprintf("Bulk size cannot exceed %d items", maxBulkSize),
        })
        return
    }
    
    // Process bulk operation...
}
```

## 🧪 **Testing Validation Logic**

Test each validation function thoroughly:

```go
func TestSKUValidation(t *testing.T) {
    tests := []struct {
        name     string
        sku      string
        expected bool
    }{
        {"Valid SKU", "ABC-123-XYZ", true},
        {"Invalid - lowercase", "abc-123-xyz", false},
        {"Invalid - wrong format", "ABC123XYZ", false},
        {"Empty string", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := isValidSKU(tt.sku)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## 💡 **Performance Optimization**

### **1. Precompile Regular Expressions**
```go
var (
    skuRegex       = regexp.MustCompile(`^[A-Z]{3}-\d{3}-[A-Z]{3}$`)
    slugRegex      = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
    warehouseRegex = regexp.MustCompile(`^WH\d{3}$`)
)

func isValidSKU(sku string) bool {
    return skuRegex.MatchString(sku)
}
```

### **2. Use Maps for Lookups**
```go
var validCurrencies = map[string]bool{
    "USD": true,
    "EUR": true,
    "GBP": true,
    "JPY": true,
}

func isValidCurrency(currency string) bool {
    return validCurrencies[currency]
}
```

## 🏭 **Real-World Applications**

These patterns are used in:

- **E-commerce platforms** - Product catalog validation
- **Financial systems** - Transaction data validation
- **Healthcare APIs** - Patient data validation
- **SaaS platforms** - Multi-tenant data validation
- **Data import systems** - Bulk data processing

## 📈 **Advanced Topics**

### **1. Conditional Validation**
```go
func validateProductByCategory(product *Product) []ValidationError {
    var errors []ValidationError
    
    switch product.Category.Name {
    case "Electronics":
        // Electronics need warranty info
        if len(product.Images) == 0 {
            errors = append(errors, ValidationError{
                Field:   "images",
                Message: "Electronics must have product images",
            })
        }
    case "Clothing":
        // Clothing needs size information
        if _, hasSize := product.Attributes["size"]; !hasSize {
            errors = append(errors, ValidationError{
                Field:   "attributes.size",
                Message: "Clothing must specify size",
            })
        }
    }
    
    return errors
}
```

### **2. Async Validation**
```go
func validateSKUUniqueness(sku string) <-chan ValidationResult {
    result := make(chan ValidationResult, 1)
    
    go func() {
        defer close(result)
        
        // Check database for existing SKU
        exists := checkSKUInDatabase(sku)
        
        result <- ValidationResult{
            Valid: !exists,
            Error: func() string {
                if exists {
                    return "SKU already exists"
                }
                return ""
            }(),
        }
    }()
    
    return result
}
```

Understanding these concepts will help you build robust, production-ready APIs with comprehensive validation and error handling. 