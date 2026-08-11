package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Product represents a product in the inventory
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Stock    int     `json:"stock"`
}

// Category represents a product category
type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Inventory represents the complete inventory data
type Inventory struct {
	Products   []Product  `json:"products"`
	Categories []Category `json:"categories"`
	NextID     int        `json:"next_id"`
}

const inventoryFile = "inventory.json"

// Global inventory instance
var inventory *Inventory

// TODO: Create the root command for the inventory CLI
// Command name: "inventory"
// Description: "Inventory Management CLI - Manage your products and categories"
var rootCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Inventory Management CLI - Manage your products and categories",
	Long:  "Inventory Management CLI - A complete inventory management system with product and category management, data persistence, and search capabilities.",
}

// TODO: Create product parent command
// Command name: "product"
// Description: "Manage products in inventory"
var productCmd = &cobra.Command{
	// TODO: Implement product command
	Use:   "product",
	Short: "Manage products in inventory",
}

// TODO: Create product add command
// Command name: "add"
// Description: "Add a new product to inventory"
// Flags: --name, --price, --category, --stock
var productAddCmd = &cobra.Command{
	// TODO: Implement product add command
	Use:   "add",
	Short: "Add a new product to inventory",
	Run: func(cmd *cobra.Command, args []string) {
		defer resetProductAddFlags(cmd)
		data := make(map[string]interface{})

		if cmd.Flags().Changed("name") {
			if name, err := cmd.Flags().GetString("name"); err == nil {
				data["name"] = name
			}
		}

		if cmd.Flags().Changed("price") {
			if price, err := cmd.Flags().GetFloat64("price"); err == nil {
				data["price"] = price
			}
		}

		if cmd.Flags().Changed("category") {
			if category, err := cmd.Flags().GetString("category"); err == nil {
				data["category"] = category
			}
		}

		if cmd.Flags().Changed("stock") {
			if stock, err := cmd.Flags().GetInt("stock"); err == nil {
				data["stock"] = stock
			}
		}

		if len(data) == 0 {
			cmd.Println("Error: No fields to create. Specify at least one of: --name, --price, --category, --stock")
			return
		}

		name, err := createProduct(data)
		if err != nil {
			cmd.Printf("Error creating product: %v\n", err)
			return
		}

		cmd.Printf("Product added successfully <%s>.\n", name)

	},
}

func resetProductAddFlags(cmd *cobra.Command) {
	cmd.Flags().Set("name", "")
	cmd.Flags().Set("category", "")
	cmd.Flags().Set("price", "0")
	cmd.Flags().Set("stock", "0")

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// TODO: Create product list command
// Command name: "list"
// Description: "List all products"
var productListCmd = &cobra.Command{
	// TODO: Implement product list command
	Use:   "list",
	Short: "List all products",
	Run: func(cmd *cobra.Command, args []string) {
		displayProductsTable(cmd, inventory.Products)
	},
}

// TODO: Create product get command
// Command name: "get"
// Description: "Get product by ID"
// Args: product ID
var productGetCmd = &cobra.Command{
	// TODO: Implement product get command
	Use:   "get",
	Short: "Get product by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Error: Invalid product ID '%s'. Must be a number.\n", args[0])
			return
		}

		displayProductTable(cmd, id)
	},
}

// TODO: Create product update command
// Command name: "update"
// Description: "Update an existing product"
// Args: product ID
// Flags: --name, --price, --category, --stock
var productUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing product",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer resetProductUpdateFlags(cmd)

		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Error: Invalid product ID '%s'. Must be a number.\n", args[0])
			return
		}

		data := make(map[string]interface{})

		if cmd.Flags().Changed("name") {
			if name, err := cmd.Flags().GetString("name"); err == nil {
				data["name"] = name
			}
		}

		if cmd.Flags().Changed("price") {
			if price, err := cmd.Flags().GetFloat64("price"); err == nil {
				data["price"] = price
			}
		}

		if cmd.Flags().Changed("category") {
			if category, err := cmd.Flags().GetString("category"); err == nil {
				data["category"] = category
			}
		}

		if cmd.Flags().Changed("stock") {
			if stock, err := cmd.Flags().GetInt("stock"); err == nil {
				data["stock"] = stock
			}
		}

		if len(data) == 0 {
			cmd.Println("Error: No fields to update. Specify at least one of: --name, --price, --category, --stock")
			return
		}

		if err := updateProduct(id, data); err != nil {
			cmd.Printf("Error updating product: %v\n", err)
			return
		}

		cmd.Printf("Product %d updated successfully\n", id)
	},
}

func resetProductUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().Set("name", "")
	cmd.Flags().Set("category", "")
	cmd.Flags().Set("price", "0")
	cmd.Flags().Set("stock", "0")

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// TODO: Create product delete command
// Command name: "delete"
// Description: "Delete a product from inventory"
// Args: product ID
var productDeleteCmd = &cobra.Command{
	// TODO: Implement product delete command
	Use:   "delete",
	Short: "Delete a product from inventory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Error: Invalid product ID '%s'. Must be a number.\n", args[0])
			return
		}

		if err := deleteProduct(id); err != nil {
			cmd.Printf("Error delete product: %v\n", err)
			return
		}

		cmd.Printf("Product %d deleted successfully\n", id)

	},
}

// TODO: Create category parent command
// Command name: "category"
// Description: "Manage categories"
var categoryCmd = &cobra.Command{
	// TODO: Implement category command
	Use:   "category",
	Short: "Manage categories",
}

// TODO: Create category add command
// Command name: "add"
// Description: "Add a new category"
// Flags: --name, --description
var categoryAddCmd = &cobra.Command{
	// TODO: Implement category add command
	Use:   "add",
	Short: "Add a new category",
	Run: func(cmd *cobra.Command, args []string) {
		defer resetCategoryAddFlags(cmd)

		data := make(map[string]interface{})

		if cmd.Flags().Changed("name") {
			if name, err := cmd.Flags().GetString("name"); err == nil {
				data["name"] = name
			}
		}

		if cmd.Flags().Changed("description") {
			if description, err := cmd.Flags().GetString("description"); err == nil {
				data["description"] = description
			}
		}

		if len(data) == 0 {
			cmd.Println("Error: No fields to create. Specify at least one of: --name, --description")
			return
		}

		name, err := createCategory(data)
		if err != nil {
			cmd.Printf("Error creating category: %v\n", err)
			return
		}

		cmd.Printf("Category added successfully <%s>\n", name)
	},
}

func resetCategoryAddFlags(cmd *cobra.Command) {
	cmd.Flags().Set("name", "")
	cmd.Flags().Set("description", "")

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// TODO: Create category list command
// Command name: "list"
// Description: "List all categories"
var categoryListCmd = &cobra.Command{
	// TODO: Implement category list command
	Use:   "list",
	Short: "List all categories",
	Run: func(cmd *cobra.Command, args []string) {
		displayCategoriesTable(cmd, inventory.Categories)
	},
}

// TODO: Create search command
// Command name: "search"
// Description: "Search products by various criteria"
// Flags: --name, --category, --min-price, --max-price
var searchCmd = &cobra.Command{
	// TODO: Implement search command
	Use:   "search",
	Short: "Search products by various criteria",
	Run: func(cmd *cobra.Command, args []string) {
		defer resetSearchFlags(cmd)

		data := make(map[string]interface{})

		if cmd.Flags().Changed("name") {
			if name, err := cmd.Flags().GetString("name"); err == nil {
				data["name"] = name
			}
		}

		if cmd.Flags().Changed("category") {
			if category, err := cmd.Flags().GetString("category"); err == nil {
				data["category"] = category
			}
		}

		if cmd.Flags().Changed("min-price") {
			if price, err := cmd.Flags().GetFloat64("min-price"); err == nil {
				data["min-price"] = price
			}
		}

		if cmd.Flags().Changed("max-price") {
			if price, err := cmd.Flags().GetFloat64("max-price"); err == nil {
				data["max-price"] = price
			}
		}

		if len(data) == 0 {
			cmd.Println("Error: No fields to update. Specify at least one of: --name, --category, --min-pric, --max-pric")
			return
		}

		prodacts := searchProduct(data)
		displayProductsTable(cmd, prodacts)
	},
}

func resetSearchFlags(cmd *cobra.Command) {
	cmd.Flags().Set("name", "")
	cmd.Flags().Set("category", "")
	cmd.Flags().Set("min-price", "0")
	cmd.Flags().Set("max-price", "0")

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// TODO: Create stats command
// Command name: "stats"
// Description: "Show inventory statistics"
var statsCmd = &cobra.Command{
	// TODO: Implement stats command
	Use:   "stats",
	Short: "Show inventory statistics",
	Run: func(cmd *cobra.Command, args []string) {
		inventoryStatistics(cmd)
	},
}

// LoadInventory loads inventory data from JSON file
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

// SaveInventory saves inventory data to JSON file
func SaveInventory() error {
	data, err := json.MarshalIndent(inventory, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(inventoryFile, data, 0644)
}

// FindProductByID finds a product by its ID
func FindProductByID(id int) (*Product, int) {
	for i, product := range inventory.Products {
		if product.ID == id {
			return &product, i
		}
	}
	return nil, -1
}

// CategoryExists checks if a category exists
func CategoryExists(name string) bool {

	if name == "" {
		return false
	}

	for _, category := range inventory.Categories {
		if category.Name == name {
			return true
		}
	}
	return false
}

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

	productUpdateCmd.Flags().StringP("name", "n", "", "Product name (required)")
	productUpdateCmd.Flags().Float64P("price", "p", 0, "Product price (required)")
	productUpdateCmd.Flags().StringP("category", "c", "", "Product category (required)")
	productUpdateCmd.Flags().IntP("stock", "s", 0, "Stock quantity (required)")
	// Add product subcommands
	productCmd.AddCommand(productAddCmd)
	productCmd.AddCommand(productListCmd)
	productCmd.AddCommand(productGetCmd)
	productCmd.AddCommand(productUpdateCmd)
	productCmd.AddCommand(productDeleteCmd)

	// Category add flags
	categoryAddCmd.Flags().StringP("name", "n", "", "Product name (required)")
	categoryAddCmd.Flags().StringP("description", "d", "", "Product description (required)")
	// Mark required flags
	categoryAddCmd.MarkFlagRequired("name")
	categoryAddCmd.MarkFlagRequired("description")
	// Add category subcommands
	categoryCmd.AddCommand(categoryAddCmd)
	categoryCmd.AddCommand(categoryListCmd)

	// Search flags
	searchCmd.Flags().StringP("name", "n", "", "Search by product name")
	searchCmd.Flags().StringP("category", "c", "", "Search by category")
	searchCmd.Flags().Float64("min-price", 0, "Minimum price")
	searchCmd.Flags().Float64("max-price", 0, "Maximum price")

	// Add all commands to root
	rootCmd.AddCommand(productCmd)
	rootCmd.AddCommand(categoryCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(statsCmd)

	if err := LoadInventory(); err != nil {
		panic(fmt.Sprintf("LoadInventory error: %v", err))
	}
}

func main() {
	// TODO: Execute root command and handle errors
	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}
}

// Product
// ===============================================================
func displayProductTable(cmd *cobra.Command, id int) {
	product, id := FindProductByID(id)
	if product == nil {
		cmd.Printf("product %d not found\n", id)
		return
	}

	cmd.Println("📦 Inventory Product:")
	cmd.Printf("%-4s | %-15s | %-8s | %-12s | %-5s\n", "ID", "Name", "Price", "Category", "Stock")
	cmd.Println("-----|-----------------|----------|--------------|-------")

	cmd.Printf("%-4d | %-15s | $%-7.2f | %-12s | %-5d\n",
		product.ID, product.Name, product.Price, product.Category, product.Stock)

}

func displayProductsTable(cmd *cobra.Command, products []Product) {
	cmd.Println("📦 Inventory Products:")
	cmd.Printf("%-4s | %-15s | %-8s | %-12s | %-5s\n", "ID", "Name", "Price", "Category", "Stock")
	cmd.Println("-----|-----------------|----------|--------------|-------")
	for _, product := range products {
		cmd.Printf("%-4d | %-15s | $%-7.2f | %-12s | %-5d\n",
			product.ID, product.Name, product.Price, product.Category, product.Stock)
	}
}

func createProduct(data map[string]interface{}) (string, error) {

	var product Product

	for field, value := range data {
		switch field {
		case "name":
			if name, ok := value.(string); ok && name != "" {
				product.Name = name
			}
		case "price":
			if price, ok := value.(float64); ok && price > 0 {
				product.Price = price
			}
		case "category":
			if category, ok := value.(string); ok && category != "" {
				product.Category = category
			}
		case "stock":
			if stock, ok := value.(int); ok && stock >= 0 {
				product.Stock = stock
			}
		}
	}

	if err := validateProduct(&product); err != nil {
		return "", fmt.Errorf("validation failed: %w", err)
	}
	product.ID = inventory.NextID
	inventory.Products = append(inventory.Products, product)
	inventory.NextID++

	if err := SaveInventory(); err != nil {
		return "", fmt.Errorf("failed to save: %w", err)
	}

	return product.Name, nil
}

func updateProduct(id int, data map[string]interface{}) error {

	product, index := FindProductByID(id)
	if product == nil {
		return fmt.Errorf("product %d not found", id)
	}

	backup := *product

	for field, value := range data {
		switch field {
		case "name":
			if name, ok := value.(string); ok && name != "" {
				product.Name = name
			}
		case "price":
			if price, ok := value.(float64); ok && price > 0 {
				product.Price = price
			}
		case "category":
			if category, ok := value.(string); ok && category != "" {
				product.Category = category
			}
		case "stock":
			if stock, ok := value.(int); ok && stock >= 0 {
				product.Stock = stock
			}
		}
	}

	if err := validateProduct(product); err != nil {
		inventory.Products[index] = backup
		return fmt.Errorf("validation failed: %w", err)
	}

	inventory.Products[index] = *product
	if err := SaveInventory(); err != nil {
		inventory.Products[index] = backup
		return fmt.Errorf("failed to save: %w", err)
	}

	return nil
}

func deleteProduct(id int) error {

	product, index := FindProductByID(id)
	if product == nil {
		return fmt.Errorf("product %d not found", id)
	}

	inventory.Products = append(inventory.Products[:index], inventory.Products[index+1:]...)

	if err := SaveInventory(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	return nil
}

func searchProduct(data map[string]interface{}) []Product {

	name := ""
	category := ""
	minPrice := 0.0
	maxPrice := 0.0

	for field, value := range data {

		switch field {
		case "name":
			if val, ok := value.(string); ok && val != "" {
				name = val
			}
		case "category":
			if val, ok := value.(string); ok && val != "" {
				category = val
			}
		case "min-price":
			if val, ok := value.(float64); ok && val > 0 {
				minPrice = val
			}
		case "max-price":
			if val, ok := value.(float64); ok && val > 0 {
				maxPrice = val
			}
		}
	}

	var results []Product
	for _, product := range inventory.Products {
		match := true
		if name != "" && product.Name != name {
			match = false
		}
		if category != "" && product.Category != category {
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

	return results

}

func validateProduct(product *Product) error {
	var errors []string

	if product.Name == "" {
		errors = append(errors, "name cannot be empty")
	}

	if product.Category == "" {
		errors = append(errors, "category cannot be empty")
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

// Category
// ===============================================================
func createCategory(data map[string]interface{}) (string, error) {

	var category Category

	for field, value := range data {
		switch field {
		case "name":
			if name, ok := value.(string); ok && name != "" {
				category.Name = name
			}
		case "description":
			if description, ok := value.(string); ok && description != "" {
				category.Description = description
			}
		}
	}

	if err := validateCategory(&category); err != nil {
		return "", fmt.Errorf("validation failed: %w", err)
	}

	if ok := CategoryExists(category.Name); ok {
		return "", fmt.Errorf("category with name %s exist", category.Name)
	}

	inventory.Categories = append(inventory.Categories, category)

	if err := SaveInventory(); err != nil {
		return "", fmt.Errorf("failed to save: %w", err)
	}

	return category.Name, nil
}

func validateCategory(category *Category) error {
	var errors []string

	if category.Name == "" {
		errors = append(errors, "name cannot be empty")
	}

	if category.Description == "" {
		errors = append(errors, "description cannot be empty")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
	}

	return nil
}

func displayCategoriesTable(cmd *cobra.Command, Categories []Category) {
	cmd.Println("📦 Category Products:")
	cmd.Printf("%-10s | %-20s\n", "Name", "Description")
	cmd.Println("-----------|---------------------")
	for _, category := range Categories {
		cmd.Printf("%-10s | %-20s\n",
			category.Name, category.Description)
	}
}

// Statistics
// ===============================================================
func inventoryStatistics(cmd *cobra.Command) {
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
	cmd.Println("📊 Inventory Statistics:")
	cmd.Printf("- Total Products: %d\n", totalProducts)
	cmd.Printf("- Total Categories: %d\n", totalCategories)
	cmd.Printf("- Total Value: $%.2f\n", totalValue)
	cmd.Printf("- Low Stock Items (< 5): %d\n", lowStockCount)
	cmd.Printf("- Out of Stock Items: %d\n", outOfStockCount)
}
