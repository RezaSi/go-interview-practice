# Hints for Challenge 2: Advanced Flags and Arguments

## Hint 1: Setting up the Root Command

Start with the basic Cobra root command structure:

```go
var rootCmd = &cobra.Command{
    Use:   "filecli",
    Short: "A file manager CLI tool",
    Long:  `A file manager CLI that demonstrates advanced flag and argument handling with Cobra.`,
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

## Hint 2: Adding Global Flags

Global flags are available to all commands. Add them to the root command using `PersistentFlags()`:

```go
func init() {
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
```

The `BoolVarP` function parameters:
- `&verbose`: pointer to the variable to store the flag value
- `"verbose"`: long flag name (--verbose)
- `"v"`: short flag name (-v)
- `false`: default value
- `"Enable verbose output"`: help text

## Hint 3: Command-Specific Flags

Add flags that are only available to specific commands using `Flags()`:

```go
func init() {
    listCmd.Flags().StringVar(&format, "format", "table", "Output format (json, table)")
}
```

## Hint 4: Required Flags

Make flags required using `MarkFlagRequired()`:

```go
func init() {
    deleteCmd.Flags().BoolVar(&force, "force", false, "Force deletion (required)")
    deleteCmd.MarkFlagRequired("force")
    
    createCmd.Flags().StringVar(&name, "name", "", "File name (required)")
    createCmd.MarkFlagRequired("name")
}
```

## Hint 5: Argument Validation

Use Cobra's built-in argument validators:

```go
var copyCmd = &cobra.Command{
    Use:  "copy <source> <destination>",
    Args: cobra.ExactArgs(2), // Requires exactly 2 arguments
    RunE: copyFile,
}

var listCmd = &cobra.Command{
    Use:  "list [directory]",
    Args: cobra.MaximumNArgs(1), // Optional argument (0 or 1)
    RunE: listFiles,
}

var deleteCmd = &cobra.Command{
    Use:  "delete <file>",
    Args: cobra.ExactArgs(1), // Requires exactly 1 argument
    RunE: deleteFile,
}
```

Common argument validators:
- `cobra.ExactArgs(n)`: Exactly n arguments
- `cobra.MinimumNArgs(n)`: At least n arguments
- `cobra.MaximumNArgs(n)`: At most n arguments
- `cobra.RangeArgs(min, max)`: Between min and max arguments
- `cobra.NoArgs`: No arguments allowed

## Hint 6: Command Implementation Pattern

Structure your command handlers consistently:

```go
func listFiles(cmd *cobra.Command, args []string) error {
    // Get directory from args or use default
    dir := "."
    if len(args) > 0 {
        dir = args[0]
    }

    // Use global verbose flag
    if verbose {
        fmt.Printf("Listing files in directory: %s\n", dir)
    }

    // Implement your logic here
    files, err := readDirectory(dir)
    if err != nil {
        return err
    }

    // Handle format flag
    if format == "json" {
        return formatAsJSON(files)
    } else {
        formatAsTable(files)
    }

    return nil
}
```

## Hint 7: Flag Types and Binding

Cobra supports various flag types:

```go
// String flags
cmd.Flags().StringVar(&stringVar, "name", "default", "help text")

// Integer flags
cmd.Flags().IntVar(&intVar, "size", 0, "help text")

// Boolean flags
cmd.Flags().BoolVar(&boolVar, "force", false, "help text")

// String slice flags
cmd.Flags().StringSliceVar(&sliceVar, "tags", []string{}, "help text")
```

## Hint 8: Adding Commands to Root

Don't forget to add your subcommands to the root command:

```go
func init() {
    rootCmd.AddCommand(listCmd)
    rootCmd.AddCommand(copyCmd)
    rootCmd.AddCommand(deleteCmd)
    rootCmd.AddCommand(createCmd)
}
```

## Hint 9: Error Handling Best Practices

Return meaningful errors and use `RunE` instead of `Run`:

```go
func deleteFile(cmd *cobra.Command, args []string) error {
    filename := args[0]

    // Check if file exists
    if !fileExists(filename) {
        return fmt.Errorf("file %s does not exist", filename)
    }

    // Perform operation
    if err := os.Remove(filename); err != nil {
        return fmt.Errorf("failed to delete file %s: %w", filename, err)
    }

    if verbose {
        fmt.Printf("Successfully deleted: %s\n", filename)
    }

    return nil
}
```

## Hint 10: Validation Helper Functions

Create helper functions for common validation patterns:

```go
func validateFileName(filename string) error {
    if filename == "" {
        return fmt.Errorf("filename cannot be empty")
    }
    
    if strings.ContainsAny(filename, "/\\:*?\"<>|") {
        return fmt.Errorf("filename contains invalid characters")
    }
    
    return nil
}

func fileExists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
}
```

## Hint 11: JSON vs Table Output

Implement flexible output formatting:

```go
func formatAsJSON(data interface{}) error {
    response := Response{
        Success: true,
        Data:    data,
    }
    
    jsonData, err := json.MarshalIndent(response, "", "  ")
    if err != nil {
        return err
    }
    
    fmt.Println(string(jsonData))
    return nil
}

func formatAsTable(files []FileInfo) {
    fmt.Printf("%-30s %-10s %-20s %s\n", "NAME", "SIZE", "MODIFIED", "TYPE")
    fmt.Println(strings.Repeat("-", 70))
    
    for _, file := range files {
        fileType := "FILE"
        if file.IsDir {
            fileType = "DIR"
        }
        
        fmt.Printf("%-30s %-10d %-20s %s\n", 
            file.Name, 
            file.Size, 
            file.ModTime.Format("2006-01-02 15:04:05"),
            fileType,
        )
    }
}
```

## Hint 12: Testing Your Implementation

Test your CLI with various flag combinations:

```bash
# Test global flag
./filecli --verbose list

# Test command-specific flags
./filecli list --format json
./filecli list --format table

# Test required flags
./filecli delete myfile.txt --force
./filecli create --name "newfile.txt" --size 100

# Test arguments
./filecli copy source.txt destination.txt
./filecli list /path/to/directory
```

## Common Pitfalls to Avoid

1. **Forgetting to add commands to root**: Use `rootCmd.AddCommand()`
2. **Not handling verbose flag**: Check `verbose` variable in handlers
3. **Incorrect argument validation**: Use appropriate `cobra.Args` validators
4. **Missing required flag validation**: Use `MarkFlagRequired()`
5. **Poor error messages**: Return descriptive errors with context 