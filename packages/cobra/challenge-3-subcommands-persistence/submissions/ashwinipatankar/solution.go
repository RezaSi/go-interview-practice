package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

const inventoryFile = "./inventory.json"

var categoryP Category
var productP Product
var searchNameF string
var searchMinPriceF float64
var searchMaxPriceF float64
var searchCategoryF string
var searchStockF bool

// Global inventory instance
var inventory *Inventory

// Create the root command for the inventory CLI
// Command name: "inventory"
// Description: "Inventory Management CLI - Manage your products and categories"
var rootCmd = &cobra.Command{
	// TODO: Implement root command
	Use:   "inventory",
	Short: "Inventory Management CLI - Manage your products and categories",
	Long:  "Inventory Management CLI - Manage your products and categories",
	Run: func(cmd *cobra.Command, args []string) {
		// Show help when no subcommand is provided
		cmd.Help()
	},
}

// Create product parent command
// Command name: "product"
// Description: "Manage products in inventory"
var productCmd = &cobra.Command{
	// TODO: Implement product command
	Use:   "product",
	Short: "Manage products in inventory",
	Long:  "Manage products in inventory",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// TODO: Create product add command
// Command name: "add"
// Description: "Add a new product to inventory"
// Flags: --name, --price, --category, --stock
var productAddCmd = &cobra.Command{
	// TODO: Implement product add command
	Use:   "add",
	Short: "Add a new product to inventory",
	Long:  "Add a new product to inventory",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("product add: ", productP.Name, productP.Price, productP.Category, productP.Stock)
		if productP.Name == "" && productP.Category == "" && productP.Price <= 0 && productP.Stock <= 0 {
			rootCmd.Println("required flag(s) \"category\", \"name\", \"price\", \"stock\" not set")
			return
		}

		// Get flag values and add product
		err := ValidateProduct(&productP)
		if err != nil {
			fmt.Println("Error adding product:", err)
			return
		}

		noOfProducts := 0
		if inventory.Products != nil && len(inventory.Products) != 0 {
			noOfProducts = len(inventory.Products)
		}

		fmt.Println("No of products: ", noOfProducts)
		fmt.Println("product name: ", productP.Name)
		productP.ID = noOfProducts + 1
		inventory.NextID++
		inventory.Products = append(inventory.Products, productP)

		//Save inventory to file
		err = SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}

		//Print success message
		rootCmd.Println("Product added successfully name: ", productP.Name, " price: ", productP.Price, " category: ", productP.Category, " stock: ", productP.Stock)

		productP = Product{}
	},
}

// TODO: Create product list command
// Command name: "list"
// Description: "List all products"
var productListCmd = &cobra.Command{
	// TODO: Implement product list command
	Use:   "list",
	Short: "List all products",
	Long:  "List all products",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Display products in table format
		rootCmd.Println(`📦 Inventory Products:
		ID  | Name          | Price    | Category     | Stock
		----|---------------|----------|--------------|-------`)
		for _, p := range inventory.Products {
			rootCmd.Printf("%d   | %s          | $%.2f    | %s      | %d\n", p.ID, p.Name, p.Price, p.Category, p.Stock)
		}
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
	Long:  "Get product by ID",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Parse ID and find product
		// TODO: Display product details

		productID, err := strconv.Atoi(args[0])
		if err != nil {
			rootCmd.Println("Invalid product ID")
			return
		}

		if productID < 0 {
			rootCmd.Println("Invalid product ID")
			return
		}

		for _, p := range inventory.Products {
			if p.ID == productID {
				rootCmd.Printf("ID: %d\nName: %s\nPrice: $%.2f\nCategory: %s\nStock: %d\n", p.ID, p.Name, p.Price, p.Category, p.Stock)
				return
			}
		}

		rootCmd.Println("Product not found")
	},
}

// TODO: Create product update command
// Command name: "update"
// Description: "Update an existing product"
// Args: product ID
// Flags: --name, --price, --category, --stock
var productUpdateCmd = &cobra.Command{
	// TODO: Implement product update command
	Use:   "update",
	Short: "Update an existing product",
	Long:  "Update an existing product",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Parse ID, update product fields
		idP, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error parsing ID:", err)
			return
		}

		if idP < 0 {
			fmt.Println("Invalid ID:", idP)
			return
		}

		product, id := FindProductByID(idP)
		if err != nil || id < 0 {
			fmt.Println("Error finding product:", err)
			fmt.Println("ID is : ", id)
			fmt.Println("IDP is: ", idP)
			return
		}

		product.Name = productP.Name
		product.Price = productP.Price
		product.Category = productP.Category
		product.Stock = productP.Stock

		inventory.Products[id] = *product

		// Save inventory to file
		err = SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}

		// Print success message
		rootCmd.Println("updated successfully name: ", productP.Name, " price: ", productP.Price, " category: ", productP.Category, " stock: ", productP.Stock)
		productP = Product{}
	},
}

// TODO: Create product delete command
// Command name: "delete"
// Description: "Delete a product from inventory"
// Args: product ID
var productDeleteCmd = &cobra.Command{
	// TODO: Implement product delete command
	Use:   "delete",
	Short: "Delete a product from inventory",
	Long:  "Delete a product from inventory",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		idP, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error parsing ID:", err)
			return
		}

		if idP < 0 {
			fmt.Println("Invalid ID:", idP)
			return
		}

		product, id := FindProductByID(idP)
		if err != nil || id < 0 {
			fmt.Println("Error finding product:", err)
			return
		}

		inventory.Products = append(inventory.Products[:id], inventory.Products[id+1:]...)

		rootCmd.Println("deleted successfully name: ", product.Name, " price: ", product.Price, " category: ", product.Category, " stock: ", product.Stock)
	},
}

// TODO: Create category parent command
// Command name: "category"
// Description: "Manage categories"
var categoryCmd = &cobra.Command{
	// TODO: Implement category command
	Use:   "category",
	Short: "Manage categories",
	Long:  "Manage categories",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("category")
	},
}

// TODO: Create category add command
// Command name: "add"
// Description: "Add a new category"
// Flags: --name, --description
var categoryAddCmd = &cobra.Command{
	// TODO: Implement category add command
	Use:   "add",
	Short: "Add a new category",
	Long:  "Add a new category",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get flag values and add category
		inventory.Categories = append(inventory.Categories, categoryP)
		// Save inventory to file
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
		// Print success message
		rootCmd.Println("Category added successfully: ", categoryP.Name, categoryP.Description)
		categoryP = Category{}
	},
}

// TODO: Create category list command
// Command name: "list"
// Description: "List all categories"
var categoryListCmd = &cobra.Command{
	// TODO: Implement category list command
	Use:   "list",
	Short: "List all categories",
	Long:  "List all categories",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Display categories
		rootCmd.Println("Inventory Category List")
		rootCmd.Println("Inventory Category")
		rootCmd.Println(" Name          | Description")
		rootCmd.Println("---------------|--------------")
		for _, c := range inventory.Categories {
			rootCmd.Printf(" %s          | %s\n", c.Name, c.Description)
		}

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
	Long:  "Search products by various criteria",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Filter products based on flags
		// TODO: Display matching products
		var products []Product
		productsByID := make(map[int]Product)

		for _, p := range inventory.Products {
			if searchNameF != "" && strings.Contains(strings.ToLower(p.Name), strings.ToLower(searchNameF)) {
				productsByID[p.ID] = p
			}

			if searchCategoryF != "" && strings.Contains(strings.ToLower(p.Category), strings.ToLower(searchCategoryF)) {
				productsByID[p.ID] = p
			}

			if searchMinPriceF != searchMaxPriceF && p.Price >= searchMinPriceF && p.Price <= searchMaxPriceF {
				productsByID[p.ID] = p
			}

			if searchStockF && p.Stock > 0 {
				productsByID[p.ID] = p
			}
		}

		for _, p := range productsByID {
			products = append(products, p)
		}

		rootCmd.Printf("🔍 Found %d product(s) in category \"%s\":\n", len(products), searchCategoryF)
		rootCmd.Println("ID  | Name          | Price    | Stock")
		rootCmd.Println("----|---------------|----------|-------")
		for _, p := range products {
			rootCmd.Printf("%d   | %s          | $%.2f    | %d\n", p.ID, p.Name, p.Price, p.Stock)
		}
	},
}

// TODO: Create stats command
// Command name: "stats"
// Description: "Show inventory statistics"
var statsCmd = &cobra.Command{
	// TODO: Implement stats command
	Use:   "stats",
	Short: "Show inventory statistics",
	Long:  "Show inventory statistics",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := LoadInventory()
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		err := SaveInventory()
		if err != nil {
			fmt.Println("Error saving inventory:", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Calculate and display statistics
		totalValue := 0.0
		lowStockItems := 0
		outOfStockItems := 0

		for _, p := range inventory.Products {
			totalValue += p.Price * float64(p.Stock)
			if p.Stock < 5 {
				lowStockItems++
			}
			if p.Stock == 0 {
				outOfStockItems++
			}
		}
		rootCmd.Println("Inventory Statistics:")
		rootCmd.Printf("Total Products: %d\n", len(inventory.Products))
		rootCmd.Printf("Total Categories: %d\n", len(inventory.Categories))
		rootCmd.Printf("Total Value: $%.2f\n", totalValue)
		rootCmd.Printf("Low Stock Items (< 5): %d\n", lowStockItems)
		rootCmd.Printf("Out of Stock Items: %d\n", outOfStockItems)

	},
}

// LoadInventory loads inventory data from JSON file
func LoadInventory() error {
	// TODO: Implement loading inventory from JSON file
	// TODO: Create default inventory if file doesn't exist
	// TODO: Handle file read errors
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

func createDefaultInventory() error {
	err := ioutil.WriteFile(inventoryFile, []byte("{}"), 0644)

	return err
}

// SaveInventory saves inventory data to JSON file
func SaveInventory() error {
	// TODO: Implement saving inventory to JSON file
	// TODO: Handle file write errors
	data, err := json.MarshalIndent(inventory, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}

	if err := ioutil.WriteFile(inventoryFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write inventory file: %w", err)
	}

	return nil
}

// FindProductByID finds a product by its ID
func FindProductByID(id int) (*Product, int) {
	// Implement finding product by ID
	if len(inventory.Products) < 1 {
		return nil, -1
	}

	for i, p := range inventory.Products {
		if p.ID == id {
			return &p, i
		}
	}

	return nil, -1
}

// CategoryExists checks if a category exists
func CategoryExists(name string) bool {
	// TODO: Implement checking if category exists
	for _, c := range inventory.Categories {
		if c.Name == name {
			return true
		}
	}

	return false
}

// AddProduct adds a new product to the inventory
func AddProduct(product *Product) error {
	// TODO: Implement adding a new product
	return nil
}

// UpdateProduct updates a product
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

func init() {
	// Add flags to product add command
	productAddCmd.Flags().StringVarP(&productP.Name, "name", "n", "", "Product name (required)")
	productAddCmd.Flags().Float64VarP(&productP.Price, "price", "p", 0, "Product price (required)")
	productAddCmd.Flags().StringVarP(&productP.Category, "category", "c", "", "Product category (required)")
	productAddCmd.Flags().IntVarP(&productP.Stock, "stock", "s", 0, "Stock quantity (required)")

	// Mark required flags
	productAddCmd.MarkFlagRequired("name")
	productAddCmd.MarkFlagRequired("price")
	productAddCmd.MarkFlagRequired("category")
	productAddCmd.MarkFlagRequired("stock")

	// Add flags to product update command
	productUpdateCmd.Flags().StringVarP(&productP.Name, "name", "n", "", "Product name (required)")
	productUpdateCmd.Flags().Float64VarP(&productP.Price, "price", "p", 0, "Product price (required)")
	productUpdateCmd.Flags().StringP("category", "c", "", "Product category (required)")
	productUpdateCmd.Flags().IntP("stock", "s", 0, "Stock quantity (required)")

	// Add flags to category add command
	categoryAddCmd.Flags().StringVarP(&categoryP.Name, "name", "n", "", "Category name (required)")
	categoryAddCmd.Flags().StringVarP(&categoryP.Description, "description", "d", "", "Category description")

	categoryAddCmd.MarkFlagRequired("name")
	categoryAddCmd.MarkFlagRequired("description")

	// Add flags to search command
	searchCmd.Flags().StringVarP(&searchNameF, "name", "n", "", "Filter by product name")
	searchCmd.Flags().StringVarP(&searchCategoryF, "category", "c", "", "Filter by category")
	searchCmd.Flags().Float64Var(&searchMinPriceF, "min-price", 0, "Minimum price filter")
	searchCmd.Flags().Float64Var(&searchMaxPriceF, "max-price", 0, "Maximum price filter")
	searchCmd.Flags().BoolVarP(&searchStockF, "in-stock", "i", false, "Show only in-stock items")

	// Add all commands to root command
	rootCmd.AddCommand(productCmd)
	rootCmd.AddCommand(categoryCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(statsCmd)

	// Add subcommands to product command
	productCmd.AddCommand(productAddCmd)
	productCmd.AddCommand(productListCmd)
	productCmd.AddCommand(productGetCmd)
	productCmd.AddCommand(productUpdateCmd)
	productCmd.AddCommand(productDeleteCmd)

	// Add subcommands to category command
	categoryCmd.AddCommand(categoryAddCmd)
	categoryCmd.AddCommand(categoryListCmd)

	// Load inventory on startup

	err := LoadInventory()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func main() {
	defer SaveInventory()

	// Execute root command and handle errors
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
