# Learning: Advanced Cobra Subcommands & Data Persistence

## 🌟 **Advanced Cobra Concepts**

This challenge introduces sophisticated CLI patterns that are essential for building production-ready applications. You'll learn to create complex command hierarchies and implement data persistence strategies.

### **Why These Patterns Matter**
- **Scalability**: Nested commands allow for organized feature sets
- **User Experience**: Logical command grouping improves discoverability
- **Data Management**: Persistence enables stateful CLI applications
- **Real-world Usage**: Most production CLIs use these patterns

## 🏗️ **Command Hierarchies & Organization**

### **1. Nested Command Structure**

Modern CLIs organize functionality into logical groups:

```
inventory                    # Root command
├── product                  # Parent command for product operations
│   ├── add                 # Product management subcommands
│   ├── list
│   ├── get <id>
│   ├── update <id>
│   └── delete <id>
├── category                 # Parent command for category operations
│   ├── add
│   └── list
├── search                   # Standalone search command
└── stats                    # Standalone statistics command
```

### **2. Command Grouping Strategy**

**By Entity Type:**
```go
// Group all product-related commands under 'product'
var productCmd = &cobra.Command{
    Use:   "product",
    Short: "Manage products in inventory",
}

// Group all category-related commands under 'category'
var categoryCmd = &cobra.Command{
    Use:   "category", 
    Short: "Manage categories",
}
```

**By Action Type:**
```go
// Alternative: Group by CRUD operations
var createCmd = &cobra.Command{Use: "create", Short: "Create resources"}
var listCmd = &cobra.Command{Use: "list", Short: "List resources"}
var updateCmd = &cobra.Command{Use: "update", Short: "Update resources"}
var deleteCmd = &cobra.Command{Use: "delete", Short: "Delete resources"}
```

### **3. Command Registration Patterns**

**Hierarchical Registration:**
```go
func init() {
    // Build command tree bottom-up
    productCmd.AddCommand(productAddCmd)
    productCmd.AddCommand(productListCmd)
    productCmd.AddCommand(productGetCmd)
    productCmd.AddCommand(productUpdateCmd)
    productCmd.AddCommand(productDeleteCmd)
    
    categoryCmd.AddCommand(categoryAddCmd)
    categoryCmd.AddCommand(categoryListCmd)
    
    // Register parent commands with root
    rootCmd.AddCommand(productCmd)
    rootCmd.AddCommand(categoryCmd)
    rootCmd.AddCommand(searchCmd)
    rootCmd.AddCommand(statsCmd)
}
```

## 💾 **Data Persistence Strategies**

### **1. File-Based Persistence**

**JSON Storage Pattern:**
```go
type Inventory struct {
    Products   []Product  `json:"products"`
    Categories []Category `json:"categories"`
    NextID     int        `json:"next_id"`
}

const inventoryFile = "inventory.json"

func LoadInventory() error {
    if _, err := os.Stat(inventoryFile); os.IsNotExist(err) {
        return createDefaultInventory()
    }
    
    data, err := ioutil.ReadFile(inventoryFile)
    if err != nil {
        return fmt.Errorf("failed to read inventory file: %w", err)
    }
    
    if err := json.Unmarshal(data, &inventory); err != nil {
        return fmt.Errorf("failed to parse inventory data: %w", err)
    }
    
    return nil
}

func SaveInventory() error {
    data, err := json.MarshalIndent(inventory, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal inventory: %w", err)
    }
    
    if err := ioutil.WriteFile(inventoryFile, data, 0644); err != nil {
        return fmt.Errorf("failed to write inventory file: %w", err)
    }
    
    return nil
}
```

### **2. Atomic Operations**

**Safe Update Pattern:**
```go
func UpdateProduct(id int, updates map[string]interface{}) error {
    // Find product
    product, index := FindProductByID(id)
    if product == nil {
        return fmt.Errorf("product %d not found", id)
    }
    
    // Create backup for rollback
    backup := *product
    
    // Apply updates
    for field, value := range updates {
        switch field {
        case "name":
            if name, ok := value.(string); ok {
                product.Name = name
            }
        case "price":
            if price, ok := value.(float64); ok {
                product.Price = price
            }
        // ... other fields
        }
    }
    
    // Validate updated product
    if err := ValidateProduct(product); err != nil {
        // Rollback on validation failure
        inventory.Products[index] = backup
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Persist changes
    if err := SaveInventory(); err != nil {
        // Rollback on save failure
        inventory.Products[index] = backup
        return fmt.Errorf("failed to save: %w", err)
    }
    
    return nil
}
```

### **3. Data Validation**

**Input Validation Pipeline:**
```go
func ValidateProduct(product *Product) error {
    var errors []string
    
    if product.Name == "" {
        errors = append(errors, "name cannot be empty")
    }
    
    if product.Price <= 0 {
        errors = append(errors, "price must be positive")
    }
    
    if product.Stock < 0 {
        errors = append(errors, "stock cannot be negative")
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
    }
    
    return nil
}
```

## 🚩 **Advanced Flag Patterns**

### **1. Command-Specific Flags**

**Product Creation Flags:**
```go
func init() {
    productAddCmd.Flags().StringP("name", "n", "", "Product name (required)")
    productAddCmd.Flags().Float64P("price", "p", 0, "Product price (required)")
    productAddCmd.Flags().StringP("category", "c", "", "Product category (required)")
    productAddCmd.Flags().IntP("stock", "s", 0, "Stock quantity (required)")
    
    // Mark required flags
    productAddCmd.MarkFlagRequired("name")
    productAddCmd.MarkFlagRequired("price")
    productAddCmd.MarkFlagRequired("category")
    productAddCmd.MarkFlagRequired("stock")
}
```

**Search Flags with Optional Filters:**
```go
func init() {
    searchCmd.Flags().StringP("name", "n", "", "Filter by product name")
    searchCmd.Flags().StringP("category", "c", "", "Filter by category")
    searchCmd.Flags().Float64("min-price", 0, "Minimum price filter")
    searchCmd.Flags().Float64("max-price", 0, "Maximum price filter")
    searchCmd.Flags().BoolP("in-stock", "i", false, "Show only in-stock items")
}
```

### **2. Flag Validation**

**Custom Validation Logic:**
```go
func validateFlags(cmd *cobra.Command) error {
    minPrice, _ := cmd.Flags().GetFloat64("min-price")
    maxPrice, _ := cmd.Flags().GetFloat64("max-price")
    
    if minPrice > 0 && maxPrice > 0 && minPrice > maxPrice {
        return fmt.Errorf("min-price cannot be greater than max-price")
    }
    
    return nil
}
```

## 🔍 **Search & Filtering Implementation**

### **1. Multi-Criteria Search**

```go
type SearchCriteria struct {
    Name      string
    Category  string
    MinPrice  float64
    MaxPrice  float64
    InStock   bool
}

func SearchProducts(criteria SearchCriteria) []Product {
    var results []Product
    
    for _, product := range inventory.Products {
        if matchesCriteria(product, criteria) {
            results = append(results, product)
        }
    }
    
    return results
}

func matchesCriteria(product Product, criteria SearchCriteria) bool {
    // Name filter
    if criteria.Name != "" {
        if !strings.Contains(strings.ToLower(product.Name), strings.ToLower(criteria.Name)) {
            return false
        }
    }
    
    // Category filter
    if criteria.Category != "" {
        if strings.ToLower(product.Category) != strings.ToLower(criteria.Category) {
            return false
        }
    }
    
    // Price range filter
    if criteria.MinPrice > 0 && product.Price < criteria.MinPrice {
        return false
    }
    if criteria.MaxPrice > 0 && product.Price > criteria.MaxPrice {
        return false
    }
    
    // Stock filter
    if criteria.InStock && product.Stock <= 0 {
        return false
    }
    
    return true
}
```

### **2. Results Formatting**

**Flexible Output Formatting:**
```go
func DisplaySearchResults(products []Product, criteria SearchCriteria) {
    fmt.Printf("🔍 Found %d product(s)", len(products))
    
    // Show active filters
    filters := []string{}
    if criteria.Name != "" {
        filters = append(filters, fmt.Sprintf("name contains '%s'", criteria.Name))
    }
    if criteria.Category != "" {
        filters = append(filters, fmt.Sprintf("category is '%s'", criteria.Category))
    }
    if criteria.MinPrice > 0 || criteria.MaxPrice > 0 {
        if criteria.MinPrice > 0 && criteria.MaxPrice > 0 {
            filters = append(filters, fmt.Sprintf("price between $%.2f and $%.2f", criteria.MinPrice, criteria.MaxPrice))
        } else if criteria.MinPrice > 0 {
            filters = append(filters, fmt.Sprintf("price >= $%.2f", criteria.MinPrice))
        } else {
            filters = append(filters, fmt.Sprintf("price <= $%.2f", criteria.MaxPrice))
        }
    }
    
    if len(filters) > 0 {
        fmt.Printf(" matching: %s", strings.Join(filters, ", "))
    }
    fmt.Println()
    
    if len(products) == 0 {
        fmt.Println("No products found matching the criteria.")
        return
    }
    
    displayProductsTable(products)
}
```

## 📊 **Statistics & Analytics**

### **1. Comprehensive Metrics**

```go
type InventoryStats struct {
    TotalProducts    int
    TotalCategories  int
    TotalValue       float64
    AveragePrice     float64
    LowStockCount    int
    OutOfStockCount  int
    TopCategory      string
    CategoryStats    map[string]CategoryStat
}

type CategoryStat struct {
    ProductCount int
    TotalValue   float64
    AveragePrice float64
}

func CalculateStats() InventoryStats {
    stats := InventoryStats{
        CategoryStats: make(map[string]CategoryStat),
    }
    
    stats.TotalProducts = len(inventory.Products)
    stats.TotalCategories = len(inventory.Categories)
    
    categoryProductCount := make(map[string]int)
    categoryValues := make(map[string]float64)
    
    for _, product := range inventory.Products {
        // Total value calculation
        productValue := product.Price * float64(product.Stock)
        stats.TotalValue += productValue
        
        // Stock analysis
        if product.Stock == 0 {
            stats.OutOfStockCount++
        } else if product.Stock < 5 {
            stats.LowStockCount++
        }
        
        // Category analysis
        categoryProductCount[product.Category]++
        categoryValues[product.Category] += productValue
    }
    
    // Calculate averages
    if stats.TotalProducts > 0 {
        stats.AveragePrice = stats.TotalValue / float64(stats.TotalProducts)
    }
    
    // Find top category
    maxProducts := 0
    for category, count := range categoryProductCount {
        if count > maxProducts {
            maxProducts = count
            stats.TopCategory = category
        }
        
        stats.CategoryStats[category] = CategoryStat{
            ProductCount: count,
            TotalValue:   categoryValues[category],
            AveragePrice: categoryValues[category] / float64(count),
        }
    }
    
    return stats
}
```

## 🎯 **Best Practices for CLI Data Management**

### **1. Error Recovery Patterns**

- **Graceful Degradation**: Continue operation with warnings when non-critical data is corrupted
- **Backup Strategies**: Maintain backup files before destructive operations
- **Validation Gates**: Validate data integrity before major operations

### **2. Performance Considerations**

- **Lazy Loading**: Load data only when needed
- **Indexing**: Create in-memory indexes for frequent lookups
- **Caching**: Cache computed statistics and search results

### **3. User Experience**

- **Progress Indicators**: Show progress for long-running operations
- **Confirmations**: Require confirmation for destructive operations
- **Helpful Messages**: Provide clear, actionable error messages

This challenge bridges the gap between simple CLI tools and production-ready applications by introducing enterprise patterns that scale with complexity and usage. 