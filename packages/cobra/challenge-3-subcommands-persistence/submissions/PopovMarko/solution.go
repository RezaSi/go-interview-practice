package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Product represents a single product stored in the inventory.
type Product struct {
	ID       int     `json:"id"`       // Unique, auto-incremented identifier.
	Name     string  `json:"name"`     // Human-readable product name.
	Price    float64 `json:"price"`    // Unit price in the local currency.
	Category string  `json:"category"` // Name of the category the product belongs to.
	Stock    int     `json:"stock"`    // Number of units currently in stock.
}

// Category represents a product category used to group products.
type Category struct {
	Name        string `json:"name"`        // Unique category name.
	Description string `json:"description"` // Optional human-readable description.
}

// Inventory is the complete, persisted state of the application: all products,
// all categories and the next ID to assign to a newly created product.
type Inventory struct {
	Products   []Product  `json:"products"`
	Categories []Category `json:"categories"`
	NextID     int        `json:"next_id"`
}

// inventoryFile is the path of the JSON file used to persist the inventory.
const inventoryFile = "inventory.json"

// inventory is the in-memory inventory shared by every command. It is loaded
// from disk in init and saved after each mutating operation.
var inventory *Inventory

// rootCmd is the top-level "inventory" command. With no subcommand it prints
// the help text describing the available subcommands.
var rootCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Inventory Management CLI - Manage your products and categories",
	Long:  "Inventory Management CLI - system for managing your products",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			cmd.PrintErrf("fatal error %v", err)
			os.Exit(1)
		}
	},
}

// productCmd is the parent command grouping every product-related subcommand.
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Manage products in inventory",
}

// productAddCmd adds a new product from the --name, --price, --category and
// --stock flags, auto-creating the category when it does not yet exist, and
// persists the result.
var productAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new product to inventory",
	Run: func(cmd *cobra.Command, args []string) {
		defer resetProductAddUpdateFlags(cmd)

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
			cmd.PrintErrf("fatal error %v", err)
			os.Exit(1)
		}
		cmd.Printf("Product added successfully with ID %d\n", product.ID)
		cmd.Printf("ID: %d, Name: %s, Price: %.2f, Category: %s, Stock: %d\n",
			product.ID, product.Name, product.Price, product.Category, product.Stock,
		)
	},
}

// productListCmd prints every product in a table.
var productListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all products",
	Run: func(cmd *cobra.Command, args []string) {
		PrintProducts(cmd, inventory.Products)
	},
}

// productGetCmd looks up a single product by the ID passed as its argument and
// prints its details, or a not-found message when the ID is unknown.
var productGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get product by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Invalid product ID: %s\n", args[0])
			return
		}

		p, index := FindProductByID(id)

		if index == -1 {
			cmd.Printf("Product with ID %d not found\n", id)
			return
		}

		cmd.Printf("ID: %d, Name: %s, Price: %.2f, Category: %s, Stock: %d\n", p.ID, p.Name, p.Price, p.Category, p.Stock)
	},
}

// productUpdateCmd updates the product identified by its argument. Only the
// fields whose flags were explicitly set are changed; the result is persisted.
var productUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing product",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer resetProductAddUpdateFlags(cmd)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Invalid product ID: %s\n", args[0])
			return
		}

		product, index := FindProductByID(id)
		if index == -1 {
			cmd.Printf("Product with ID %d not found\n", id)
			return
		}
		if cmd.Flag("name").Changed {
			product.Name = cmd.Flag("name").Value.String()
		}
		if cmd.Flag("price").Changed {
			price, _ := cmd.Flags().GetFloat64("price")
			product.Price = price
		}
		if cmd.Flag("category").Changed {
			product.Category = cmd.Flag("category").Value.String()
		}
		if cmd.Flag("stock").Changed {
			stock, _ := cmd.Flags().GetInt("stock")
			product.Stock = stock
		}
		inventory.Products[index] = *product
		if err := SaveInventory(); err != nil {
			cmd.PrintErrf("fatal error %v", err)
			return
		}
		cmd.Printf("Product with ID %d updated successfully\n", id)
	},
}

// productDeleteCmd removes the product identified by its argument from the
// inventory and persists the change.
var productDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a product from inventory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Printf("Invalid product ID: %s\n", args[0])
			return
		}

		p, index := FindProductByID(id)

		if index == -1 {
			cmd.Printf("Product with ID %d not found\n", id)
			return
		}

		inventory.Products = append(inventory.Products[:index], inventory.Products[index+1:]...)
		if err := SaveInventory(); err != nil {
			cmd.PrintErrf("fatal error: %v", err)
			os.Exit(1)
		}
		cmd.Printf("Product with ID: %d deleted successfully\n", p.ID)
	},
}

// categoryCmd is the parent command grouping every category-related subcommand.
var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Manage categories",
}

// categoryAddCmd adds a new category from the --name and --description flags,
// rejecting duplicates, and persists the result.
var categoryAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new category",
	Run: func(cmd *cobra.Command, args []string) {
		defer resetCategoryAddFlags(cmd)

		name := cmd.Flag("name").Value.String()
		description := cmd.Flag("description").Value.String()

		if CategoryExists(name) {
			cmd.Printf("Category with name: %s already exists", name)
			return
		}

		category := Category{
			Name:        name,
			Description: description,
		}

		inventory.Categories = append(inventory.Categories, category)
		if err := SaveInventory(); err != nil {
			cmd.PrintErrf("fatal error %v", err)
			os.Exit(1)
		}

		cmd.Printf("Category added successfully with name: %s\n", name)
	},
}

// categoryListCmd prints every category in a table.
var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all categories",
	Run: func(cmd *cobra.Command, args []string) {
		tw := tabwriter.NewWriter(cmd.OutOrStdout(), 6, 3, 3, ' ', tabwriter.AlignRight)
		fmt.Fprintln(tw, "Name |\tDescription\t")
		fmt.Fprintln(tw, "---- |\t-----------\t")
		for _, c := range inventory.Categories {
			fmt.Fprintf(tw, "%s |\t%s\t\n", c.Name, c.Description)
		}
		tw.Flush()
	},
}

// searchCmd filters products by the --name, --category, --min-price and
// --max-price flags and prints the matches.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search products by various criteria",
	Run: func(cmd *cobra.Command, args []string) {
		defer resetSearchFlags(cmd)
		params := make(map[string]interface{})

		if cmd.Flag("name").Value.String() != "" && cmd.Flag("name").Changed {
			params["name"] = cmd.Flag("name").Value.String()
		}
		if cmd.Flag("category").Value.String() != "" && cmd.Flag("category").Changed {
			params["category"] = cmd.Flag("category").Value.String()
		}
		if cmd.Flag("min-price").Changed {
			params["minPrice"], _ = cmd.Flags().GetFloat64("min-price")

		}
		if cmd.Flag("max-price").Changed {
			params["maxPrice"], _ = cmd.Flags().GetFloat64("max-price")
		}

		matchingProducts := FilterProducts(params)
		PrintProducts(cmd, matchingProducts)
	},
}

// FilterProducts returns the products matching every criterion present in
// params. Supported keys are "name", "category" (exact match), "minPrice" and
// "maxPrice" (inclusive bounds). Absent keys are ignored, so an empty params
// map returns all products.
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

		if match {
			results = append(results, p)
		}
	}

	return results
}

// resetSearchFlags clears the search flags and their "changed" state so that
// values do not leak between successive invocations of the shared command.
func resetSearchFlags(cmd *cobra.Command) {
	for _, pair := range [][2]string{
		{"name", ""},
		{"category", ""},
		{"min-price", "0"},
		{"max-price", "0"},
	} {
		if err := cmd.Flags().Set(pair[0], pair[1]); err != nil {
			cmd.PrintErrf("fatal error: %v", err)
		}
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// resetProductAddUpdateFlags clears the product add/update flags and their
// "changed" state so that values do not leak between successive invocations of
// the shared command.
func resetProductAddUpdateFlags(cmd *cobra.Command) {
	for _, pair := range [][2]string{
		{"name", ""},
		{"category", ""},
		{"price", "0"},
		{"stock", "0"},
	} {
		if err := cmd.Flags().Set(pair[0], pair[1]); err != nil {
			cmd.PrintErrf("fatal error: %v", err)
		}
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// resetCategoryAddFlags clears the category add flags and their "changed" state
// so that values do not leak between successive invocations of the shared
// command.
func resetCategoryAddFlags(cmd *cobra.Command) {
	for _, pair := range [][2]string{
		{"name", ""},
		{"description", ""},
	} {
		if err := cmd.Flags().Set(pair[0], pair[1]); err != nil {
			cmd.PrintErrf("fatal error: %v", err)
		}
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// PrintProducts writes the given products as an aligned table to the command's
// output writer.
func PrintProducts(cmd *cobra.Command, products []Product) {
	tw := tabwriter.NewWriter(cmd.OutOrStdout(), 6, 3, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "ID |\tName |\tPrice |\tCategory |\tStock |\t")
	fmt.Fprintln(tw, "-- |\t---- |\t----- |\t-------- |\t----- |\t")
	for _, p := range products {
		fmt.Fprintf(tw, "%d |\t%s |\t%.2f |\t%s |\t%d |\t\n", p.ID, p.Name, p.Price, p.Category, p.Stock)
	}
	tw.Flush()
}

// statsCmd prints aggregate statistics about the inventory: product and
// category counts, total stock value, and low/out-of-stock counts.
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show inventory statistics",
	Run: func(cmd *cobra.Command, args []string) {
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

func newInventory() *Inventory {
	return &Inventory{
		Products:   []Product{},
		Categories: []Category{},
		NextID:     1,
	}
}

// LoadInventory loads the inventory from the JSON file into the global
// inventory variable. When the file does not exist it initializes an empty
// inventory instead of returning an error.
func LoadInventory() error {
	data, err := os.Open(inventoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			inventory = newInventory()
			return nil
		}
		return err
	}

	defer data.Close()
	if err := json.NewDecoder(data).Decode(&inventory); err != nil {
		if errors.Is(err, io.EOF) {
			inventory = newInventory()
			return nil
		}
		return err
	}
	return nil
}

// SaveInventory writes the current global inventory to the JSON file,
// truncating any previous contents.
func SaveInventory() error {
	data, err := os.Create(inventoryFile)
	if err != nil {
		return err
	}
	defer data.Close()

	if err := json.NewEncoder(data).Encode(inventory); err != nil {
		return err
	}
	return nil
}

// FindProductByID returns a pointer to the product with the given ID and its
// index in the slice, or nil and -1 when no product matches.
func FindProductByID(id int) (*Product, int) {
	for i := range inventory.Products {
		if inventory.Products[i].ID == id {
			return &inventory.Products[i], i
		}
	}
	return nil, -1
}

// CategoryExists reports whether a category with the given name exists.
func CategoryExists(name string) bool {
	for _, c := range inventory.Categories {
		if c.Name == name {
			return true
		}
	}
	return false
}

// init wires the command hierarchy, registers flags and loads the inventory
// from disk before any command runs.
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
	searchCmd.Flags().Float64("max-price", 0.0, "Max price to filter by")

	rootCmd.AddCommand(statsCmd)

	if err := LoadInventory(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error %v", err)
		os.Exit(1)
	}
}

// main executes the root command and exits with a non-zero status on failure.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v", err)
		os.Exit(1)
	}
}
