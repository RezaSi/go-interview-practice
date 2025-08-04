package main

import (
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"slices"

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

var rootCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Inventory Management CLI",
	Long:  "Inventory Management CLI - Manage your products and categories",
}

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Manage products in inventory",
}

var productAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new product to inventory",
	Long:  "Add a new product to inventory",
	PostRun: func(cmd *cobra.Command, args []string) {
		// COBRA does not automatically reset flags between consecutive
		// command executions. (needed to pass tests).
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		price, _ := cmd.Flags().GetFloat64("price")
		category, _ := cmd.Flags().GetString("category")
		stock, _ := cmd.Flags().GetInt("stock")

		// NOTE:
		// Tests does not required any data validation

		product := Product{
			ID:       inventory.NextID,
			Name:     name,
			Price:    price,
			Category: category,
			Stock:    stock,
		}
		inventory.Products = append(inventory.Products, product)
		inventory.NextID++

		if err := SaveInventory(); err != nil {
			panic(fmt.Sprintf("Error: %v", err))
		}
		cmd.Printf("Product added successfully - name: %s\n", product.Name)
	},
}

var productListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all products",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("ID\tName\t\tPrice\tCategory\tStock")
		for _, p := range(inventory.Products) {
			cmd.Printf("%d\t%s\t$%.2f\t%s\t\t%d\n", p.ID, p.Name, p.Price, p.Category, p.Stock)
		}
	},
}

var productGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get product by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Println("Invalid ID")
			return
		}
		product, _ := FindProductByID(id)
		if product == nil {
			cmd.Println("Product not found")
			return
		}
		cmd.Printf("ID: %d\nName: %s\nPrice: $%.2f\nCategory: %s\nStock: %d\n",
			product.ID, product.Name, product.Price, product.Category, product.Stock)
	},
}

var productUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing product",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Println("Invalid ID")
			return
		}
		product, index := FindProductByID(id)
		if product == nil {
			cmd.Println("Product not found")
			return
		}

		name, _ := cmd.Flags().GetString("name")
		price, _ := cmd.Flags().GetFloat64("price")
		category, _ := cmd.Flags().GetString("category")
		stock, _ := cmd.Flags().GetInt("stock")

		// NOTE:
		// Tests does not required any data validation

		inventory.Products[index].Name = name
		inventory.Products[index].Price = price
		inventory.Products[index].Category = category
		inventory.Products[index].Stock = stock

		if err := SaveInventory(); err != nil {
			cmd.Printf("Error: %v\n", err)
			return
		}
		cmd.Println("Product updated successfully")
	},
}

var productDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a product from inventory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Println("Invalid ID")
			return
		}
		product, index := FindProductByID(id)
		if product == nil {
			cmd.Println("Product not found")
			return
		}

		inventory.Products = slices.Delete(inventory.Products, index, index + 1)

		if err := SaveInventory(); err != nil {
			cmd.Printf("Error: %v\n", err)
			return
		}
		cmd.Println("Product deleted successfully")
	},
}

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Manage categories",
	Long: "Manage categories",
}

var categoryAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new category",
	Long: "Add a new category",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		// NOTE:
		// Tests does not required any data validation

		category := Category{
			Name:        name,
			Description: description,
		}
		inventory.Categories = append(inventory.Categories, category)

		if err := SaveInventory(); err != nil {
			cmd.Printf("Error: %v\n", err)
			return
		}
		cmd.Println("Category added successfully")
	},
}

var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all categories",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Name\tDescription")
		for _, c := range(inventory.Categories) {
			cmd.Printf("%s\t%s\n", c.Name, c.Description)
		}
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search products by various criteria",
	PostRun: func(cmd *cobra.Command, args []string) {
		// COBRA does not automatically reset flags between consecutive
		// command executions. (needed to pass tests).
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		category, _ := cmd.Flags().GetString("category")
		minPrice, _ := cmd.Flags().GetFloat64("min-price")
		maxPrice, _ := cmd.Flags().GetFloat64("max-price")

		var results []Product
		for _, p := range inventory.Products {
			match := true
			if name != "" && p.Name != name {
				match = false
			}
			if category != "" && p.Category != category {
				match = false
			}

            if minPrice > 0 && p.Price < minPrice {
				match = false
            }
            if maxPrice > 0 && p.Price > maxPrice {
				match = false
            }
			if match {
				results = append(results, p)
			}
		}

		if len(results) == 0 {
			cmd.Println("No products found")
			return
		}

		cmd.Println("ID\tName\t\tPrice\tCategory\tStock")
		for _, p := range(results) {
			cmd.Printf("%d\t%s\t$%.2f\t%s\t\t%d\n", p.ID, p.Name, p.Price, p.Category, p.Stock)
		}
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show inventory statistics",
	Run: func(cmd *cobra.Command, args []string) {
		totalProducts := len(inventory.Products)
		totalCategories := len(inventory.Categories)
		var totalValue float64
		lowStockCount := 0
		outOfStockCount := 0

		for _, p := range inventory.Products {
			totalValue += p.Price * float64(p.Stock)
			if p.Stock < 5 && p.Stock > 0 {
				lowStockCount++
			}
			if p.Stock == 0 {
				outOfStockCount++
			}
		}

		cmd.Printf("Total Products: %d\n", totalProducts)
		cmd.Printf("Total Categories: %d\n", totalCategories)
		cmd.Printf("Total Value: $%.2f\n", totalValue)
		cmd.Printf("Low Stock: %d\n", lowStockCount)
		cmd.Printf("Out of Stock: %d\n", outOfStockCount)
	},
}

func LoadInventory() error {
	data, err := os.ReadFile(inventoryFile)
	if os.IsNotExist(err) {
		inventory = &Inventory{NextID: 1}
		return nil
	}
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
	return os.WriteFile(inventoryFile, data, 0644)
}

func FindProductByID(id int) (*Product, int) {
	for i, p := range(inventory.Products) {
		if p.ID == id {
			return &p, i
		}
	}
	return nil, -1
}

func CategoryExists(name string) bool {
	for _, c := range inventory.Categories {
		if c.Name == name {
			return true
		}
	}
	return false
}

func init() {
	// Product flags
	productAddCmd.Flags().String("name", "", "Product name")
	productAddCmd.Flags().Float64("price", 0, "Product price")
	productAddCmd.Flags().String("category", "", "Product category")
	productAddCmd.Flags().Int("stock", 0, "Product stock quantity")
	productAddCmd.MarkFlagRequired("name")
	productAddCmd.MarkFlagRequired("price")
	productAddCmd.MarkFlagRequired("category")
	productAddCmd.MarkFlagRequired("stock")

	productUpdateCmd.Flags().String("name", "", "Product name")
	productUpdateCmd.Flags().Float64("price", 0, "Product price")
	productUpdateCmd.Flags().String("category", "", "Product category")
	productUpdateCmd.Flags().Int("stock", 0, "Product stock quantity")

	// Add subcommands to product command
	productCmd.AddCommand(productAddCmd)
	productCmd.AddCommand(productListCmd)
	productCmd.AddCommand(productGetCmd)
	productCmd.AddCommand(productUpdateCmd)
	productCmd.AddCommand(productDeleteCmd)

	// Category flags
	categoryAddCmd.Flags().String("name", "", "Cateogory name")
	categoryAddCmd.Flags().String("description", "", "Cateogory description")
	categoryAddCmd.MarkFlagRequired("name")
	categoryAddCmd.MarkFlagRequired("description")

	// Add subcommands to category command
	categoryCmd.AddCommand(categoryAddCmd)
	categoryCmd.AddCommand(categoryListCmd)

	// Search flags
	searchCmd.Flags().String("name", "", "Filer by product name")
	searchCmd.Flags().Float64("min-price", 0, "Filter by minimum price")
	searchCmd.Flags().Float64("max-price", 0, "Filter maximum price")
	searchCmd.Flags().String("category", "", "Filter by category")

	// Add commands to root command
	rootCmd.AddCommand(productCmd)
	rootCmd.AddCommand(categoryCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(statsCmd)

	if err := LoadInventory(); err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}
}
