# Challenge 3: Subcommands & Data Persistence

Build an **Inventory Management CLI** using Cobra that demonstrates advanced subcommand organization and JSON data persistence.

## Challenge Requirements

Create a CLI application called `inventory` that manages product inventory with:

1. **Product Management** - Add, list, update, delete products
2. **Category Management** - Organize products by categories  
3. **JSON Persistence** - Save/load data from JSON file
4. **Search & Filter** - Find products by various criteria
5. **Nested Subcommands** - Organized command hierarchy

## Expected CLI Structure

```
inventory                          # Root command
inventory product add              # Add a new product
inventory product list             # List all products
inventory product get <id>         # Get product by ID
inventory product update <id>      # Update product
inventory product delete <id>      # Delete product
inventory category add             # Add a new category
inventory category list            # List all categories
inventory search --name <name>     # Search products by name
inventory search --category <cat>  # Search by category
inventory stats                    # Show inventory statistics
```

## Sample Output

**Add Product (`inventory product add`):**
```
$ inventory product add --name "Laptop" --price 999.99 --category "Electronics" --stock 10
✅ Product added successfully!
ID: 1, Name: Laptop, Price: $999.99, Category: Electronics, Stock: 10
```

**List Products (`inventory product list`):**
```
$ inventory product list
📦 Inventory Products:
ID  | Name          | Price    | Category     | Stock
----|---------------|----------|--------------|-------
1   | Laptop        | $999.99  | Electronics  | 10
2   | Coffee Mug    | $12.99   | Kitchen      | 25
3   | Notebook      | $5.49    | Stationery   | 100
```

**Search Products (`inventory search --category "Electronics"`):**
```
$ inventory search --category "Electronics"
🔍 Found 1 product(s) in category "Electronics":
ID  | Name          | Price    | Stock
----|---------------|----------|-------
1   | Laptop        | $999.99  | 10
```

**Statistics (`inventory stats`):**
```
$ inventory stats
📊 Inventory Statistics:
- Total Products: 3
- Total Categories: 3
- Total Value: $1,018.47
- Low Stock Items (< 5): 0
- Out of Stock Items: 0
```

## Data Model

```go
type Product struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Category string  `json:"category"`
    Stock    int     `json:"stock"`
}

type Category struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type Inventory struct {
    Products   []Product  `json:"products"`
    Categories []Category `json:"categories"`
    NextID     int        `json:"next_id"`
}
```

## Implementation Requirements

### Product Subcommands
- `product add` - Add new product with flags: `--name`, `--price`, `--category`, `--stock`
- `product list` - Display all products in table format
- `product get <id>` - Show details of specific product
- `product update <id>` - Update product with same flags as add
- `product delete <id>` - Remove product from inventory

### Category Subcommands
- `category add` - Add new category with flags: `--name`, `--description`
- `category list` - Show all categories

### Search Command
- `search` - Search products with flags: `--name`, `--category`, `--min-price`, `--max-price`

### Data Persistence
- Store data in `inventory.json` file
- Auto-create file if it doesn't exist
- Load data on startup, save after modifications
- Handle file read/write errors gracefully

### Error Handling
- Validate required flags
- Check for duplicate product/category names
- Handle invalid IDs
- Provide helpful error messages

## Testing Requirements

Your solution must pass tests for:
- All subcommands execute correctly
- Data persistence works (save/load from JSON)
- Product CRUD operations function properly
- Category management works
- Search functionality with different filters
- Statistics calculation is accurate
- Error handling for invalid inputs
- Command structure matches expected hierarchy 