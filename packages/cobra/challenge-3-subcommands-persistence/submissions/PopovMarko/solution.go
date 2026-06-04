package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

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
	// TODO: Implement root command
	Use:   "inventory",
	Short: "Inventory Management CLI - Manage your products and categories",
	Long:  "Inventory Management CLI - system for manageing your products",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Fprintf(os.Stdout, "fatal error %v", err)
			os.Exit(1)
		}
	},
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
		// TODO: Get flag values and add product
		// TODO: Save inventory to file
		// TODO: Print success message
		if !cmd.Flags().HasFlags() {
			cmd.Println("all flags are required")
			return
		}
		id := inventory.NextID
		name := cmd.Flag("name").Value.String()
		price, _ := cmd.Flags().GetFloat64("price")
		category := cmd.Flag("category").Value.String()
		stock, _ := cmd.Flags().GetInt("stock")

		product := Product{
			ID:       id,
			Name:     name,
			Price:    price,
			Category: category,
			Stock:    stock,
		}
		inventory.NextID++

		if !CategoryExists(category) {
			newCategory := Category{
				Name: category,
			}
			inventory.Categories = append(inventory.Categories, newCategory)
		}

		inventory.Products = append(inventory.Products, product)
		if err := SaveInventory(); err != nil {
			fmt.Fprintf(os.Stdout, "fatal error %v", err)
			os.Exit(1)
		}
		cmd.Printf("Product added successfully with ID %d\n", product.ID)
		cmd.Printf("ID: %d, Name: %s, Price: %.2f, Category: %s, Stock: %d\n",
			product.ID, product.Name, product.Price, product.Category, product.Stock,
		)

	},
}

// TODO: Create product list command
// Command name: "list"
// Description: "List all products"
var productListCmd = &cobra.Command{
	// TODO: Implement product list command
	Use:   "list",
	Short: "list all products",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Display products in table format
		PrintProducts(cmd, inventory.Products)
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
		// TODO: Parse ID and find product
		// TODO: Display product details
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Invalid product ID: %s\n", args[0])
		}

		p, index := FindProductByID(id)

		if index == -1 {
			cmd.Printf("Product with ID %d not found\n", id)
			return
		}

		cmd.Printf("ID: %d, Name: %s, Price: %.2f, Catigory: %s, Stock: %d\n", p.ID, p.Name, p.Price, p.Category, p.Stock)
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Parse ID, update product fields
		// TODO: Save inventory to file
		// TODO: Print success message
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Invalid product ID: %s\n", args[0])
		}

		product, index := FindProductByID(id)
		if index == -1 {
			cmd.Printf("Product with ID %d not found\n", id)
			return
		}
		if cmd.HasFlags() && cmd.Flag("name").Changed {
			product.Name = cmd.Flag("name").Value.String()
		}
		if cmd.HasFlags() && cmd.Flag("price").Changed {
			price, _ := cmd.Flags().GetFloat64("price")
			product.Price = price
		}
		if cmd.HasFlags() && cmd.Flag("category").Changed {
			product.Category = cmd.Flag("category").Value.String()
		}
		if cmd.HasFlags() && cmd.Flag("stock").Changed {
			stock, _ := cmd.Flags().GetInt("stock")
			product.Stock = stock
		}
		inventory.Products[index] = *product
		if err := SaveInventory(); err != nil {
			fmt.Fprintf(os.Stdout, "fatal error %v", err)
		}
		cmd.Printf("Product with ID %d updated successfully\n", id)
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Parse ID and delete product
		// TODO: Save inventory to file
		// TODO: Print success message
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Invalid product ID: %s\n", args[0])
		}

		p, index := FindProductByID(id)

		if index == -1 {
			cmd.Printf("Product with ID %d not found\n", id)
			return
		}

		// inventory.Products[index] = nil
		inventory.Products = append(inventory.Products[:index], inventory.Products[index+1:]...)
		cmd.Printf("Product with ID: %d deleted successfully\n", p.ID)
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
		// TODO: Get flag values and add category
		// TODO: Save inventory to file
		// TODO: Print success message
		name := cmd.Flag("name").Value.String()
		description := cmd.Flag("description").Value.String()

		if CategoryExists(name) {
			cmd.Printf("Category with name: %s alreadi exists", name)
			return
		}

		category := Category{
			Name:        name,
			Description: description,
		}

		inventory.Categories = append(inventory.Categories, category)
		if err := SaveInventory(); err != nil {
			fmt.Fprintf(os.Stdout, "fatal error %v", err)
			os.Exit(1)
		}

		cmd.Printf("Category added successfully with name: %s\n", name)
	},
}

// TODO: Create category list command
// Command name: "list"
// Description: "List all categories"
var categoryListCmd = &cobra.Command{
	// TODO: Implement category list command
	Use:   "list",
	Short: "list all categories",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Display categories
		tw := tabwriter.NewWriter(cmd.OutOrStdout(), 6, 3, 3, ' ', tabwriter.AlignRight)
		fmt.Fprintln(tw, "Name |\tDescription\t")
		fmt.Fprintln(tw, "---- |\t-----------\t")
		for _, c := range inventory.Categories {
			fmt.Fprintf(tw, "%s |\t%s\t\n", c.Name, c.Description)
		}
		tw.Flush()
	},
}

// TODO: Create search command
// Command name: "search"
// Description: "Search products by various criteria"
// Flags: --name, --category, --min-price, --max-price
var searchCmd = &cobra.Command{
	// TODO: Implement search command
	Use:   "search",
	Short: "search products by various criteria",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Filter products based on flags
		// TODO: Display matching products
		defer resetSearchFlags(cmd)
		params := make(map[string]interface{})

		if cmd.Flag("name").Value.String() != "" && cmd.Flag("name").Changed {
			params["name"] = cmd.Flag("name").Value.String()
		}
		if cmd.Flag("category").Value.String() != "" && cmd.Flag("category").Changed {
			params["category"] = cmd.Flag("category").Value.String()
		}
		if f, _ := cmd.Flags().GetFloat64("min-price"); f != 0.0 && cmd.Flag("min-price").Changed {
			params["minPrice"], _ = cmd.Flags().GetFloat64("min-price")
		}
		if f, _ := cmd.Flags().GetFloat64("max-price"); f != 0.0 && cmd.Flag("max-price").Changed {
			params["maxPrice"], _ = cmd.Flags().GetFloat64("max-price")
		}

		matchingProducts := FilterProducts(params)
		PrintProducts(cmd, matchingProducts)
	},
}

func FilterProducts(params map[string]interface{}) []Product {
	var (
		results []Product
		match   bool
	)

	for _, p := range inventory.Products {
		match = true
		if name, ok := params["name"].(string); ok && p.Name != name {
			match = false
		}
		if category, ok := params["category"].(string); ok && p.Category != category {
			match = false
		}
		if minPrice, ok := params["minPrice"].(float64); ok && p.Price < minPrice {
			match = false
		}
		if maxPrice, ok := params["maxPrice"].(float64); ok && p.Price > maxPrice {
			match = false
		}

		// switch {
		// case params["name"] != nil && p.Name != params["name"].(string):
		// 	match = false
		// case params["category"] != nil && p.Category != params["category"].(string):
		// 	match = false
		// case params["minPrice"] != nil && p.Price < params["minPrice"].(float64):
		// 	match = false
		// case params["maxPrice"] != nil && p.Price > params["maxPrice"].(float64):
		// 	match = false
		// }
		if match {
			results = append(results, p)
		}
	}

	return results
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

func PrintProducts(cmd *cobra.Command, products []Product) {
	tw := tabwriter.NewWriter(cmd.OutOrStdout(), 6, 3, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "ID |\tName |\tPrice |\tCategory |\tStock |\t")
	fmt.Fprintln(tw, "-- |\t---- |\t----- |\t-------- |\t----- |\t")
	for _, p := range products {
		fmt.Fprintf(tw, "%d |\t%s |\t%.2f |\t%s |\t%d |\t\n", p.ID, p.Name, p.Price, p.Category, p.Stock)
	}
	tw.Flush()
}

// TODO: Create stats command
// Command name: "stats"
// Description: "Show inventory statistics"
var statsCmd = &cobra.Command{
	// TODO: Implement stats command
	Use:   "stats",
	Short: "Show inventory statistics",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Calculate and display statistics
		totalProducts := len(inventory.Products)
		totalCategories := len(inventory.Categories)
		totalValue := 0.0
		lowStockProducts := 0
		outOfStockProducts := 0
		for _, p := range inventory.Products {
			totalValue += p.Price * float64(p.Stock)
			if p.Stock < 5 {
				lowStockProducts++
			}
			if p.Stock == 0 {
				outOfStockProducts++
			}
		}
		cmd.Printf("Total Products: %d\n", totalProducts)
		cmd.Printf("Total Categories: %d\n", totalCategories)
		cmd.Printf("Total Value: $%.2f\n", totalValue)
		cmd.Printf("Low Stock Items: %d\n", lowStockProducts)
		cmd.Printf("Out of Stock Items: %d\n", outOfStockProducts)
	},
}

// LoadInventory loads inventory data from JSON file
func LoadInventory() error {
	data, err := os.OpenFile(inventoryFile, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			inventory = &Inventory{
				Products:   []Product{},
				Categories: []Category{},
				NextID:     1,
			}
			return nil
		} else {
			return err
		}
	}
	defer data.Close()
	if err := json.NewDecoder(data).Decode(&inventory); err != nil {
		return err
	}
	return nil
}

// SaveInventory saves inventory data to JSON file
func SaveInventory() error {
	data, err := os.OpenFile(inventoryFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer data.Close()
	if err := json.NewEncoder(data).Encode(inventory); err != nil {
		return err
	}
	return nil
}

// FindProductByID finds a product by its ID
func FindProductByID(id int) (*Product, int) {
	// TODO: Implement finding product by ID
	// TODO: Return product and index, or nil and -1 if not found
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

func init() {
	rootCmd.AddCommand(productCmd)

	productCmd.AddCommand(productAddCmd)
	productAddCmd.Flags().String("name", "", "Name of the product")
	productAddCmd.MarkFlagRequired("name")
	productAddCmd.Flags().Float64("price", 0.0, "Price of the product")
	productAddCmd.Flags().String("category", "", "Category of the product")
	productAddCmd.MarkFlagRequired("category")
	productAddCmd.Flags().Int("stock", 0, "Stock quantity of the product")

	productCmd.AddCommand(productListCmd)

	productCmd.AddCommand(productGetCmd)
	productGetCmd.Args = cobra.ExactArgs(1)

	productCmd.AddCommand(productUpdateCmd)
	productUpdateCmd.Args = cobra.ExactArgs(1)
	productUpdateCmd.Flags().String("name", "", "Name of the product")
	productUpdateCmd.Flags().Float64("price", 0.0, "Price of the product")
	productUpdateCmd.Flags().String("category", "", "Category of the product")
	productUpdateCmd.Flags().Int("stock", 0, "Stock quantity of the product")

	productCmd.AddCommand(productDeleteCmd)
	productDeleteCmd.Args = cobra.ExactArgs(1)

	rootCmd.AddCommand(categoryCmd)
	categoryCmd.AddCommand(categoryAddCmd)
	categoryAddCmd.Flags().String("name", "", "Name of the category")
	categoryAddCmd.MarkFlagRequired("name")
	categoryAddCmd.Flags().String("description", "", "Description of the category")

	categoryCmd.AddCommand(categoryListCmd)

	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().String("name", "", "Name to filter by")
	searchCmd.Flags().String("category", "", "Category to filter by")
	searchCmd.Flags().Float64("min-price", 0.0, "Min price to filter by")
	searchCmd.Flags().Float64("max-price", 0.0, "Max price to fiter by")

	rootCmd.AddCommand(statsCmd)

	// TODO: Load inventory on startup
	if err := LoadInventory(); err != nil {
		fmt.Fprintf(os.Stdout, "fatal error %v", err)
		os.Exit(1)
	}
}

func main() {
	// TODO: Execute root command and handle errors
	rootCmd.Execute()
}
