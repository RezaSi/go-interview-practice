# Hints for Challenge 4: Advanced Features & Middleware

## Hint 1: Setting up the Root Command with Middleware

Configure the config-manager CLI with middleware support:

```go
var rootCmd = &cobra.Command{
    Use:   "config-manager",
    Short: "Configuration Management CLI - Advanced configuration management with plugins and middleware",
    Long:  "A powerful configuration management system supporting multiple formats, middleware, plugins, and environment integration.",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // Execute middleware before any command
        ApplyMiddleware(cmd, args)
    },
    PersistentPostRun: func(cmd *cobra.Command, args []string) {
        // Cleanup after command execution
        if err := SaveConfig(); err != nil {
            fmt.Printf("Warning: Failed to save config: %v\n", err)
        }
    },
}
```

## Hint 2: Implementing Middleware System

Create a middleware pipeline:

```go
type Middleware func(*cobra.Command, []string) error

var middlewares []Middleware

func ApplyMiddleware(cmd *cobra.Command, args []string) error {
    for _, middleware := range middlewares {
        if err := middleware(cmd, args); err != nil {
            return fmt.Errorf("middleware failed: %w", err)
        }
    }
    return nil
}

// Validation middleware
func ValidationMiddleware(cmd *cobra.Command, args []string) error {
    result := ValidateConfiguration()
    if !result.Valid && len(result.Errors) > 0 {
        fmt.Printf("⚠️  Configuration warnings: %v\n", result.Warnings)
    }
    return nil
}

// Audit middleware
func AuditMiddleware(cmd *cobra.Command, args []string) error {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("🔍 [%s] Executing: %s %v\n", timestamp, cmd.Name(), args)
    return nil
}
```

## Hint 3: Nested Key Access with Dot Notation

Implement configuration key access with dot notation:

```go
func GetNestedValue(key string) (interface{}, bool) {
    if config == nil || config.Data == nil {
        return nil, false
    }
    
    parts := strings.Split(key, ".")
    current := config.Data
    
    for i, part := range parts {
        if value, exists := current[part]; exists {
            if i == len(parts)-1 {
                // Last part, return the value
                return value, true
            }
            
            // Navigate deeper if it's a map
            if nestedMap, ok := value.(map[string]interface{}); ok {
                current = nestedMap
            } else {
                return nil, false
            }
        } else {
            return nil, false
        }
    }
    
    return nil, false
}

func SetNestedValue(key string, value interface{}) error {
    if config == nil {
        config = &Config{
            Data:     make(map[string]interface{}),
            Format:   "json",
            Version:  "1.0.0",
            Metadata: ConfigMetadata{},
        }
    }
    
    parts := strings.Split(key, ".")
    current := config.Data
    
    // Navigate to the parent of the target key
    for i, part := range parts[:len(parts)-1] {
        if _, exists := current[part]; !exists {
            current[part] = make(map[string]interface{})
        }
        
        if nestedMap, ok := current[part].(map[string]interface{}); ok {
            current = nestedMap
        } else {
            return fmt.Errorf("cannot set nested value: %s is not a map", strings.Join(parts[:i+1], "."))
        }
    }
    
    // Set the final value
    current[parts[len(parts)-1]] = value
    config.Metadata.Modified = time.Now()
    
    return nil
}
```

## Hint 4: Multi-Format Configuration Support

Implement format detection and conversion:

```go
func DetectFormat(filename string) string {
    ext := strings.ToLower(filepath.Ext(filename))
    switch ext {
    case ".yaml", ".yml":
        return "yaml"
    case ".toml":
        return "toml"
    case ".json":
        return "json"
    default:
        return "json" // Default fallback
    }
}

func ConvertFormat(targetFormat string) error {
    if config.Format == targetFormat {
        return nil // Already in target format
    }
    
    // Update format metadata
    config.Format = targetFormat
    config.Metadata.Modified = time.Now()
    
    return nil
}

func LoadConfigFromFile(filename string) error {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }
    
    format := DetectFormat(filename)
    
    switch format {
    case "json":
        err = json.Unmarshal(data, config)
    case "yaml":
        err = yaml.Unmarshal(data, config)
    case "toml":
        // Add TOML support if needed
        return fmt.Errorf("TOML format not implemented yet")
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
    
    if err != nil {
        return fmt.Errorf("failed to parse %s: %w", format, err)
    }
    
    config.Metadata.Source = filename
    config.Metadata.Modified = time.Now()
    
    return nil
}
```

## Hint 5: Plugin System Implementation

Create a basic plugin architecture:

```go
type PluginInterface interface {
    Initialize() error
    GetCommands() []*cobra.Command
    GetInfo() PluginInfo
    Cleanup() error
}

type PluginInfo struct {
    Name        string
    Version     string
    Description string
}

func RegisterPlugin(plugin Plugin) error {
    // Check if plugin already exists
    for _, existing := range plugins {
        if existing.Name == plugin.Name {
            return fmt.Errorf("plugin %s already registered", plugin.Name)
        }
    }
    
    // Add to plugin registry
    plugins = append(plugins, plugin)
    
    fmt.Printf("✅ Plugin '%s' v%s registered successfully\n", plugin.Name, plugin.Version)
    return nil
}

// Mock plugin installation
func InstallPlugin(name string) error {
    // In a real implementation, this would download and install a plugin
    plugin := Plugin{
        Name:        name,
        Version:     "1.0.0",
        Status:      "active",
        Description: fmt.Sprintf("Mock plugin: %s", name),
        Commands:    []PluginCommand{},
        Config:      make(map[string]string),
    }
    
    return RegisterPlugin(plugin)
}
```

## Hint 6: Environment Variable Integration

Implement environment variable synchronization:

```go
func SyncWithEnvironment() error {
    if config.Metadata.Environment == nil {
        config.Metadata.Environment = make(map[string]string)
    }
    
    // Define environment variable prefixes to sync
    prefixes := []string{"CONFIG_", "APP_"}
    
    for _, prefix := range prefixes {
        for _, env := range os.Environ() {
            pair := strings.SplitN(env, "=", 2)
            if len(pair) == 2 && strings.HasPrefix(pair[0], prefix) {
                key := strings.TrimPrefix(pair[0], prefix)
                key = strings.ToLower(strings.ReplaceAll(key, "_", "."))
                
                // Store in environment tracking
                config.Metadata.Environment[pair[0]] = pair[1]
                
                // Set in configuration
                SetNestedValue(key, pair[1])
            }
        }
    }
    
    config.Metadata.Modified = time.Now()
    return nil
}
```

## Hint 7: Validation Pipeline

Implement comprehensive configuration validation:

```go
func ValidateConfiguration() ValidationResult {
    result := ValidationResult{
        Valid:    true,
        Errors:   []string{},
        Warnings: []string{},
    }
    
    if config == nil || config.Data == nil {
        result.Valid = false
        result.Errors = append(result.Errors, "configuration is empty")
        return result
    }
    
    // Validate required fields
    requiredFields := []string{"app.name", "app.version"}
    for _, field := range requiredFields {
        if _, exists := GetNestedValue(field); !exists {
            result.Warnings = append(result.Warnings, fmt.Sprintf("recommended field %s is missing", field))
        }
    }
    
    // Validate data types
    if port, exists := GetNestedValue("server.port"); exists {
        if portStr, ok := port.(string); ok {
            if _, err := strconv.Atoi(portStr); err != nil {
                result.Valid = false
                result.Errors = append(result.Errors, "server.port must be a valid integer")
            }
        }
    }
    
    return result
}
```

## Hint 8: Custom Help Templates

Set up custom help formatting:

```go
func SetCustomHelpTemplate() {
    helpTemplate := `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
    
    cobra.AddTemplateFunc("StyleHeading", func(s string) string {
        return fmt.Sprintf("\033[1;36m%s\033[0m", s) // Cyan bold
    })
    
    rootCmd.SetHelpTemplate(helpTemplate)
}
```

## Hint 9: Configuration Loading with Viper Integration

Use Viper for advanced configuration management:

```go
func LoadConfig() error {
    viper.SetConfigName("config")
    viper.SetConfigType("json")
    viper.AddConfigPath(".")
    viper.AddConfigPath("$HOME/.config-manager")
    
    // Environment variable support
    viper.AutomaticEnv()
    viper.SetEnvPrefix("CONFIG")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Config file not found; create default
            return createDefaultConfig()
        } else {
            return fmt.Errorf("error reading config file: %w", err)
        }
    }
    
    // Unmarshal into our config structure
    var tempData map[string]interface{}
    if err := viper.Unmarshal(&tempData); err != nil {
        return fmt.Errorf("error unmarshaling config: %w", err)
    }
    
    config = &Config{
        Data:     tempData,
        Format:   "json",
        Version:  "1.0.0",
        Metadata: ConfigMetadata{
            Created:  time.Now(),
            Modified: time.Now(),
            Source:   viper.ConfigFileUsed(),
        },
    }
    
    return nil
}
```

## Hint 10: Command Implementation Examples

Implement key commands:

```go
// Config get command
var configGetCmd = &cobra.Command{
    Use:   "get <key>",
    Short: "Get configuration value by key",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        key := args[0]
        value, exists := GetNestedValue(key)
        
        if !exists {
            fmt.Printf("❌ Key '%s' not found\n", key)
            return
        }
        
        fmt.Printf("📋 Configuration Value:\n")
        fmt.Printf("Key: %s\n", key)
        fmt.Printf("Value: %v\n", value)
        fmt.Printf("Type: %T\n", value)
        fmt.Printf("Source: %s\n", config.Metadata.Source)
        fmt.Printf("Last Modified: %s\n", config.Metadata.Modified.Format("2006-01-02 15:04:05"))
    },
}

// Config set command
var configSetCmd = &cobra.Command{
    Use:   "set <key> <value>",
    Short: "Set configuration value",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        key := args[0]
        value := args[1]
        
        // Try to infer type
        var typedValue interface{} = value
        if intVal, err := strconv.Atoi(value); err == nil {
            typedValue = intVal
        } else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
            typedValue = floatVal
        } else if boolVal, err := strconv.ParseBool(value); err == nil {
            typedValue = boolVal
        }
        
        if err := SetNestedValue(key, typedValue); err != nil {
            fmt.Printf("❌ Failed to set value: %v\n", err)
            return
        }
        
        fmt.Printf("🔧 Configuration updated successfully\n")
        fmt.Printf("Key: %s\n", key)
        fmt.Printf("Value: %v\n", typedValue)
        fmt.Printf("Type: %T\n", typedValue)
        fmt.Printf("Format: %s\n", config.Format)
    },
}
```

Remember to register all middleware in the `init()` function and implement proper error handling throughout the application! 