# Hints for Challenge 3: Subcommands & Data Persistence

## Hint 1: Setting up the Root Command

Configure the inventory CLI root command:

```go
var rootCmd = &cobra.Command{
    Use:   "inventory",
    Short: "Inventory Management CLI - Manage your products and categories",
    Long:  "A complete inventory management system with product and category management, data persistence, and search capabilities.",
}
```

## Hint 2: Creating Nested Command Structure

Use `AddCommand()` to create hierarchical commands:

```go
func init() {
    // Add product subcommands
    productCmd.AddCommand(productAddCmd)
    productCmd.AddCommand(productListCmd)
    productCmd.AddCommand(productGetCmd)
    productCmd.AddCommand(productUpdateCmd)
    productCmd.AddCommand(productDeleteCmd)
    
    // Add category subcommands
    categoryCmd.AddCommand(categoryAddCmd)
    categoryCmd.AddCommand(categoryListCmd)
    
    // Add all commands to root
    rootCmd.AddCommand(productCmd)
    rootCmd.AddCommand(categoryCmd)
    rootCmd.AddCommand(searchCmd)
    rootCmd.AddCommand(statsCmd)
}
```

## Hint 3: Adding Flags to Commands

Add flags for user input:

```go
func init() {
    // Product add flags
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

## Hint 4: JSON Data Persistence

Implement JSON file operations:

```go
func LoadInventory() error {
    if _, err := os.Stat(inventoryFile); os.IsNotExist(err) {
        // Create default inventory
        inventory = &Inventory{
            Products:   []Product{},
            Categories: []Category{},
            NextID:     1,
        }
        return SaveInventory()
    }
    
    data, err := ioutil.ReadFile(inventoryFile)
    if err != nil {
        return err
    }
    
    return json.Unmarshal(data, &inventory)
}

func SaveInventory() error {
    data, err := json.MarshalIndent(inventory, "", "  ")
    if err != nil {
        return err
    }
    
    return ioutil.WriteFile(inventoryFile, data, 0644)
}
```

## Hint 5: Getting Flag Values in Commands

Access flag values in command execution:

```go
Run: func(cmd *cobra.Command, args []string) {
    name, _ := cmd.Flags().GetString("name")
    price, _ := cmd.Flags().GetFloat64("price")
    category, _ := cmd.Flags().GetString("category")
    stock, _ := cmd.Flags().GetInt("stock")
    
    product := Product{
        ID:       inventory.NextID,
        Name:     name,
        Price:    price,
        Category: category,
        Stock:    stock,
    }
    
    inventory.Products = append(inventory.Products, product)
    inventory.NextID++
    
    SaveInventory()
    fmt.Printf("✅ Product added successfully!\n")
    fmt.Printf("ID: %d, Name: %s, Price: $%.2f, Category: %s, Stock: %d\n", 
        product.ID, product.Name, product.Price, product.Category, product.Stock)
},
```

## Hint 6: Implementing Nested Key Access

Support dot notation for configuration-like access:

```go
func GetNestedValue(key string) (interface{}, bool) {
    parts := strings.Split(key, ".")
    
    for _, product := range inventory.Products {
        if parts[0] == "product" && len(parts) > 1 {
            // Handle product.field access
            if fmt.Sprintf("%d", product.ID) == parts[1] {
                if len(parts) > 2 {
                    switch parts[2] {
                    case "name":
                        return product.Name, true
                    case "price":
                        return product.Price, true
                    // ... other fields
                    }
                }
                return product, true
            }
        }
    }
    
    return nil, false
}
```

## Hint 7: Implementing Search Functionality

Add search flags and filtering logic:

```go
func init() {
    searchCmd.Flags().StringP("name", "n", "", "Search by product name")
    searchCmd.Flags().StringP("category", "c", "", "Search by category")
    searchCmd.Flags().Float64("min-price", 0, "Minimum price")
    searchCmd.Flags().Float64("max-price", 0, "Maximum price")
}

// In search command Run function:
Run: func(cmd *cobra.Command, args []string) {
    name, _ := cmd.Flags().GetString("name")
    category, _ := cmd.Flags().GetString("category")
    minPrice, _ := cmd.Flags().GetFloat64("min-price")
    maxPrice, _ := cmd.Flags().GetFloat64("max-price")
    
    var results []Product
    
    for _, product := range inventory.Products {
        match := true
        
        if name != "" && !strings.Contains(strings.ToLower(product.Name), strings.ToLower(name)) {
            match = false
        }
        if category != "" && strings.ToLower(product.Category) != strings.ToLower(category) {
            match = false
        }
        if minPrice > 0 && product.Price < minPrice {
            match = false
        }
        if maxPrice > 0 && product.Price > maxPrice {
            match = false
        }
        
        if match {
            results = append(results, product)
        }
    }
    
    // Display results
    fmt.Printf("🔍 Found %d product(s):\n", len(results))
    // ... format and display results
},
```

## Hint 8: Calculating Statistics

Implement comprehensive statistics:

```go
Run: func(cmd *cobra.Command, args []string) {
    totalProducts := len(inventory.Products)
    totalCategories := len(inventory.Categories)
    
    var totalValue float64
    lowStockCount := 0
    outOfStockCount := 0
    
    for _, product := range inventory.Products {
        totalValue += product.Price * float64(product.Stock)
        
        if product.Stock == 0 {
            outOfStockCount++
        } else if product.Stock < 5 {
            lowStockCount++
        }
    }
    
    fmt.Println("📊 Inventory Statistics:")
    fmt.Printf("- Total Products: %d\n", totalProducts)
    fmt.Printf("- Total Categories: %d\n", totalCategories)
    fmt.Printf("- Total Value: $%.2f\n", totalValue)
    fmt.Printf("- Low Stock Items (< 5): %d\n", lowStockCount)
    fmt.Printf("- Out of Stock Items: %d\n", outOfStockCount)
},
```

## Hint 9: Error Handling

Add proper error handling throughout:

```go
func FindProductByID(id int) (*Product, int) {
    for i, product := range inventory.Products {
        if product.ID == id {
            return &product, i
        }
    }
    return nil, -1
}

// In commands:
product, index := FindProductByID(id)
if product == nil {
    fmt.Printf("❌ Product with ID %d not found\n", id)
    return
}
```

## Hint 10: Table Formatting

Create nicely formatted table output:

```go
func displayProductsTable(products []Product) {
    fmt.Println("📦 Inventory Products:")
    fmt.Printf("%-4s | %-15s | %-8s | %-12s | %-5s\n", "ID", "Name", "Price", "Category", "Stock")
    fmt.Println("-----|-----------------|----------|--------------|-------")
    
    for _, product := range products {
        fmt.Printf("%-4d | %-15s | $%-7.2f | %-12s | %-5d\n",
            product.ID, product.Name, product.Price, product.Category, product.Stock)
    }
}
```

Remember to call `LoadInventory()` in the `init()` function and handle all file operations with proper error checking! 