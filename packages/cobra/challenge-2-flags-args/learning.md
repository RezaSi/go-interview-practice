# Learning: Advanced Cobra CLI - Flags and Arguments

## 🌟 **What Are CLI Flags and Arguments?**

Command-line interfaces use **flags** and **arguments** to accept user input and configure program behavior. Understanding how to handle these properly is crucial for building professional CLI tools.

### **Flags vs Arguments**
- **Flags**: Optional named parameters that modify behavior (`--verbose`, `--format json`)
- **Arguments**: Positional parameters that provide data (`copy source.txt dest.txt`)

## 🏗️ **Core Concepts**

### **1. Flag Types**

Cobra supports various flag types to handle different data:

```go
// Boolean flags (true/false)
var verbose bool
cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

// String flags
var format string
cmd.Flags().StringVar(&format, "format", "table", "Output format")

// Integer flags
var size int
cmd.Flags().IntVar(&size, "size", 0, "File size in bytes")

// String slice flags (multiple values)
var tags []string
cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "File tags")
```

### **2. Global vs Command-Specific Flags**

**Global Flags** (available to all commands):
```go
// Add to root command
rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
```

**Command-Specific Flags** (only available to one command):
```go
// Add to specific command
listCmd.Flags().StringVar(&format, "format", "table", "Output format")
```

### **3. Required vs Optional Flags**

```go
// Required flag
createCmd.Flags().StringVar(&name, "name", "", "File name (required)")
createCmd.MarkFlagRequired("name")

// Optional flag with default
listCmd.Flags().StringVar(&format, "format", "table", "Output format")
```

## 📐 **Argument Validation**

Cobra provides built-in validators for command arguments:

### **Common Validators**
```go
// Exactly N arguments
var copyCmd = &cobra.Command{
    Use:  "copy <source> <destination>",
    Args: cobra.ExactArgs(2),
}

// At least N arguments
var processCmd = &cobra.Command{
    Use:  "process <file1> [file2...]",
    Args: cobra.MinimumNArgs(1),
}

// At most N arguments
var listCmd = &cobra.Command{
    Use:  "list [directory]",
    Args: cobra.MaximumNArgs(1),
}

// Range of arguments
var mergeCmd = &cobra.Command{
    Use:  "merge <files...>",
    Args: cobra.RangeArgs(2, 5),
}

// No arguments
var statusCmd = &cobra.Command{
    Use:  "status",
    Args: cobra.NoArgs,
}
```

### **Custom Argument Validation**
```go
var customCmd = &cobra.Command{
    Use: "custom",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 1 {
            return fmt.Errorf("requires at least 1 argument")
        }
        
        for _, arg := range args {
            if !strings.HasSuffix(arg, ".txt") {
                return fmt.Errorf("all arguments must be .txt files")
            }
        }
        
        return nil
    },
}
```

## 🎯 **Flag Binding and Variables**

### **Variable Binding**
Flags can be bound to variables for easy access:

```go
var config struct {
    Verbose bool
    Format  string
    Size    int
}

func init() {
    rootCmd.PersistentFlags().BoolVar(&config.Verbose, "verbose", false, "Verbose output")
    listCmd.Flags().StringVar(&config.Format, "format", "table", "Output format")
    createCmd.Flags().IntVar(&config.Size, "size", 0, "File size")
}
```

### **Flag Aliases (Short and Long Forms)**
```go
// Both --verbose and -v work
cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

// Both --format and -f work  
cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format")
```

## 🔧 **Advanced Flag Features**

### **Flag Dependencies**
```go
func init() {
    cmd.Flags().StringVar(&username, "username", "", "Username")
    cmd.Flags().StringVar(&password, "password", "", "Password")
    
    // If username is provided, password is required
    cmd.MarkFlagsRequiredTogether("username", "password")
    
    // These flags cannot be used together
    cmd.MarkFlagsMutuallyExclusive("username", "token")
}
```

### **Flag Value Validation**
```go
func init() {
    cmd.Flags().StringVar(&format, "format", "table", "Output format")
    
    // Validate flag values
    cmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
        return []string{"json", "table", "csv"}, cobra.ShellCompDirectiveNoFileComp
    })
}

func validateFormat(cmd *cobra.Command, args []string) error {
    validFormats := []string{"json", "table", "csv"}
    for _, valid := range validFormats {
        if format == valid {
            return nil
        }
    }
    return fmt.Errorf("invalid format: %s (valid: %s)", format, strings.Join(validFormats, ", "))
}
```

## 💡 **Command Structure Best Practices**

### **Consistent Command Handler Pattern**
```go
func commandHandler(cmd *cobra.Command, args []string) error {
    // 1. Validate inputs
    if err := validateInputs(args); err != nil {
        return err
    }
    
    // 2. Process global flags
    if verbose {
        fmt.Printf("Processing command with args: %v\n", args)
    }
    
    // 3. Execute main logic
    result, err := doWork(args)
    if err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    
    // 4. Format and output results
    return outputResult(result)
}
```

### **Error Handling Strategy**
```go
func processFile(cmd *cobra.Command, args []string) error {
    filename := args[0]
    
    // Check file exists
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return fmt.Errorf("file %s does not exist", filename)
    }
    
    // Process file
    if err := doProcessing(filename); err != nil {
        return fmt.Errorf("failed to process %s: %w", filename, err)
    }
    
    // Success message
    if verbose {
        fmt.Printf("Successfully processed: %s\n", filename)
    }
    
    return nil
}
```

## 📊 **Output Formatting Patterns**

### **JSON vs Human-Readable Output**
```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func outputResult(data interface{}) error {
    if format == "json" {
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
    
    // Human-readable format
    return outputTable(data)
}
```

### **Table Formatting**
```go
func outputTable(files []FileInfo) error {
    // Header
    fmt.Printf("%-30s %-10s %-20s %s\n", "NAME", "SIZE", "MODIFIED", "TYPE")
    fmt.Println(strings.Repeat("-", 75))
    
    // Rows
    for _, file := range files {
        fileType := "FILE"
        if file.IsDir {
            fileType = "DIR"
        }
        
        fmt.Printf("%-30s %-10d %-20s %s\n", 
            truncate(file.Name, 30),
            file.Size,
            file.ModTime.Format("2006-01-02 15:04:05"),
            fileType,
        )
    }
    
    return nil
}

func truncate(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-3] + "..."
}
```

## 🧪 **Testing CLI Applications**

### **Command Testing Pattern**
```go
func TestCommand(t *testing.T) {
    // Setup
    cmd := &cobra.Command{
        Use: "test",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Your command logic
            return nil
        },
    }
    
    // Capture output
    buf := new(bytes.Buffer)
    cmd.SetOut(buf)
    cmd.SetErr(buf)
    cmd.SetArgs([]string{"arg1", "arg2"})
    
    // Execute
    err := cmd.Execute()
    
    // Assert
    assert.NoError(t, err)
    assert.Contains(t, buf.String(), "expected output")
}
```

### **Flag Testing**
```go
func TestFlags(t *testing.T) {
    var verbose bool
    var format string
    
    cmd := &cobra.Command{Use: "test"}
    cmd.Flags().BoolVar(&verbose, "verbose", false, "Verbose output")
    cmd.Flags().StringVar(&format, "format", "table", "Output format")
    
    // Test flag parsing
    cmd.SetArgs([]string{"--verbose", "--format", "json"})
    err := cmd.Execute()
    
    assert.NoError(t, err)
    assert.True(t, verbose)
    assert.Equal(t, "json", format)
}
```

## 🚀 **Real-World Examples**

### **Professional CLI Tools**
Understanding how popular tools use flags and arguments:

**Docker:**
```bash
docker run -d --name myapp -p 8080:80 nginx:latest
# -d: detached mode (boolean flag)
# --name: container name (string flag)  
# -p: port mapping (string flag)
# nginx:latest: image argument
```

**kubectl:**
```bash
kubectl get pods --namespace production --output json
# get: subcommand
# pods: argument  
# --namespace: string flag
# --output: string flag
```

**git:**
```bash
git commit -m "message" --author "name <email>"
# commit: subcommand
# -m: message flag (string)
# --author: author flag (string)
```

### **File Manager CLI Example**
```bash
# Global verbose flag
filecli --verbose list

# Command-specific format flag
filecli list --format json /home/user

# Required safety flag
filecli delete --force important.txt

# Multiple flags and arguments
filecli copy --preserve-permissions source.txt backup.txt
```

## 🎨 **Advanced Patterns**

### **Command Chaining and Pipelines**
```go
// Support for command chaining
var chainCmd = &cobra.Command{
    Use: "chain",
    PreRunE: func(cmd *cobra.Command, args []string) error {
        // Validate prerequisites
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        // Main execution
        return nil
    },
    PostRunE: func(cmd *cobra.Command, args []string) error {
        // Cleanup or follow-up
        return nil
    },
}
```

### **Dynamic Command Generation**
```go
func generateCommands() {
    for _, service := range services {
        cmd := &cobra.Command{
            Use:   service.Name,
            Short: fmt.Sprintf("Manage %s service", service.Name),
            RunE:  createServiceHandler(service),
        }
        
        // Add service-specific flags
        for _, flag := range service.Flags {
            cmd.Flags().StringVar(&flag.Value, flag.Name, flag.Default, flag.Help)
        }
        
        rootCmd.AddCommand(cmd)
    }
}
```

## 📚 **Key Takeaways**

1. **Flag Design**: Use clear, consistent naming conventions
2. **Validation**: Validate inputs early and provide helpful error messages  
3. **Documentation**: Write descriptive help text for all flags and commands
4. **Testing**: Test all flag combinations and edge cases
5. **User Experience**: Provide both short and long flag forms when appropriate
6. **Error Handling**: Return meaningful errors with actionable guidance
7. **Output Formatting**: Support both machine-readable (JSON) and human-readable formats

## 🔗 **Further Reading**

- [Cobra Documentation](https://github.com/spf13/cobra)
- [12-Factor CLI Apps](https://medium.com/@jdxcode/12-factor-cli-apps-dd3c227a0e46)
- [Command Line Interface Guidelines](https://clig.dev/)
- [POSIX Utility Conventions](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html) 