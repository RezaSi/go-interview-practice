# Learning: Advanced Cobra Patterns & Enterprise CLI Architecture

## 🌟 **Enterprise CLI Patterns**

This challenge represents the pinnacle of CLI application development, introducing patterns used in production systems like Kubernetes, Docker, and Terraform. You'll master middleware systems, plugin architectures, and advanced configuration management.

### **Why These Patterns Are Critical**
- **Extensibility**: Plugin systems allow third-party extensions
- **Maintainability**: Middleware separates concerns cleanly
- **Scalability**: Configuration management supports complex deployments
- **Production Ready**: These patterns are battle-tested in major CLI tools

## 🏗️ **Middleware Architecture**

### **1. Middleware Concept**

Middleware provides a way to execute code before and after command execution, similar to web frameworks:

```go
type Middleware func(*cobra.Command, []string) error

// Middleware pipeline
var middlewares []Middleware

// Execute all middleware in order
func ApplyMiddleware(cmd *cobra.Command, args []string) error {
    for _, middleware := range middlewares {
        if err := middleware(cmd, args); err != nil {
            return fmt.Errorf("middleware failed: %w", err)
        }
    }
    return nil
}
```

### **2. Common Middleware Types**

**Validation Middleware:**
```go
func ValidationMiddleware(cmd *cobra.Command, args []string) error {
    // Validate configuration state
    result := ValidateConfiguration()
    if !result.Valid {
        return fmt.Errorf("configuration validation failed: %v", result.Errors)
    }
    return nil
}
```

**Audit Middleware:**
```go
func AuditMiddleware(cmd *cobra.Command, args []string) error {
    // Log command execution
    log.Printf("Command executed: %s with args: %v", cmd.Name(), args)
    return nil
}
```

**Authentication Middleware:**
```go
func AuthMiddleware(cmd *cobra.Command, args []string) error {
    // Check authentication state
    if !isAuthenticated() {
        return fmt.Errorf("authentication required")
    }
    return nil
}
```

### **3. Middleware Registration**

**Global Middleware (All Commands):**
```go
func init() {
    middlewares = append(middlewares, 
        ValidationMiddleware,
        AuditMiddleware,
        AuthMiddleware,
    )
    
    rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
        if err := ApplyMiddleware(cmd, args); err != nil {
            fmt.Printf("Middleware error: %v\n", err)
            os.Exit(1)
        }
    }
}
```

**Command-Specific Middleware:**
```go
var sensitiveCmd = &cobra.Command{
    Use: "delete",
    PreRun: func(cmd *cobra.Command, args []string) {
        // Additional validation for sensitive operations
        if !confirmDestructiveOperation() {
            os.Exit(1)
        }
    },
    Run: func(cmd *cobra.Command, args []string) {
        // Command implementation
    },
}
```

## 🔌 **Plugin Architecture**

### **1. Plugin Interface Design**

**Core Plugin Interface:**
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
    Author      string
    License     string
}
```

### **2. Plugin Registration System**

**Plugin Registry:**
```go
type PluginRegistry struct {
    plugins map[string]PluginInterface
    mutex   sync.RWMutex
}

func (r *PluginRegistry) Register(plugin PluginInterface) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    info := plugin.GetInfo()
    if _, exists := r.plugins[info.Name]; exists {
        return fmt.Errorf("plugin %s already registered", info.Name)
    }
    
    if err := plugin.Initialize(); err != nil {
        return fmt.Errorf("plugin initialization failed: %w", err)
    }
    
    r.plugins[info.Name] = plugin
    
    // Register plugin commands
    for _, cmd := range plugin.GetCommands() {
        rootCmd.AddCommand(cmd)
    }
    
    return nil
}
```

### **3. Dynamic Plugin Loading**

**Plugin Discovery:**
```go
func LoadPluginsFromDirectory(dir string) error {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        if filepath.Ext(file.Name()) == ".so" { // Linux shared library
            if err := loadPlugin(filepath.Join(dir, file.Name())); err != nil {
                log.Printf("Failed to load plugin %s: %v", file.Name(), err)
            }
        }
    }
    
    return nil
}
```

## ⚙️ **Advanced Configuration Management**

### **1. Multi-Format Support**

**Format Detection and Parsing:**
```go
type ConfigFormat string

const (
    JSON ConfigFormat = "json"
    YAML ConfigFormat = "yaml"
    TOML ConfigFormat = "toml"
)

func DetectFormat(filename string) ConfigFormat {
    ext := strings.ToLower(filepath.Ext(filename))
    switch ext {
    case ".yaml", ".yml":
        return YAML
    case ".toml":
        return TOML
    default:
        return JSON
    }
}

func ParseConfig(data []byte, format ConfigFormat) (*Config, error) {
    var config Config
    
    switch format {
    case JSON:
        err := json.Unmarshal(data, &config)
        return &config, err
    case YAML:
        err := yaml.Unmarshal(data, &config)
        return &config, err
    case TOML:
        err := toml.Unmarshal(data, &config)
        return &config, err
    default:
        return nil, fmt.Errorf("unsupported format: %s", format)
    }
}
```

### **2. Configuration Hierarchy**

**Precedence Order (Highest to Lowest):**
1. Command-line flags
2. Environment variables
3. Configuration files
4. Default values

```go
func LoadConfigurationHierarchy() error {
    // 1. Start with defaults
    config = NewDefaultConfig()
    
    // 2. Load from config files (multiple sources)
    configSources := []string{
        "/etc/app/config.yaml",
        "$HOME/.config/app/config.yaml",
        "./config.yaml",
    }
    
    for _, source := range configSources {
        if expanded := os.ExpandEnv(source); fileExists(expanded) {
            if err := mergeConfigFromFile(expanded); err != nil {
                log.Printf("Failed to load config from %s: %v", expanded, err)
            }
        }
    }
    
    // 3. Override with environment variables
    applyEnvironmentOverrides()
    
    // 4. Apply command-line flag overrides (handled by Cobra/Viper)
    
    return nil
}
```

### **3. Nested Configuration Access**

**Dot Notation Support:**
```go
func GetNestedValue(key string) (interface{}, bool) {
    parts := strings.Split(key, ".")
    current := config.Data
    
    for _, part := range parts {
        if value, exists := current[part]; exists {
            if nestedMap, ok := value.(map[string]interface{}); ok {
                current = nestedMap
            } else {
                // Reached a leaf value
                return value, true
            }
        } else {
            return nil, false
        }
    }
    
    return current, true
}

func SetNestedValue(key string, value interface{}) error {
    parts := strings.Split(key, ".")
    current := config.Data
    
    // Navigate to parent
    for _, part := range parts[:len(parts)-1] {
        if _, exists := current[part]; !exists {
            current[part] = make(map[string]interface{})
        }
        current = current[part].(map[string]interface{})
    }
    
    // Set the value
    current[parts[len(parts)-1]] = value
    return nil
}
```

## 🔐 **Environment Integration**

### **1. Environment Variable Mapping**

**Automatic Environment Binding:**
```go
func ConfigureEnvironmentIntegration() {
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    // Manual mappings for complex cases
    viper.BindEnv("database.host", "DATABASE_HOST")
    viper.BindEnv("database.port", "DATABASE_PORT")
    viper.BindEnv("api.key", "API_SECRET_KEY")
}
```

**Environment Variable Precedence:**
```go
type EnvironmentMapping struct {
    ConfigKey string
    EnvVar    string
    Required  bool
    Validator func(string) error
}

var envMappings = []EnvironmentMapping{
    {"server.port", "PORT", false, validatePort},
    {"database.url", "DATABASE_URL", true, validateURL},
    {"api.key", "API_KEY", true, validateAPIKey},
}

func ApplyEnvironmentOverrides() error {
    for _, mapping := range envMappings {
        if value := os.Getenv(mapping.EnvVar); value != "" {
            if mapping.Validator != nil {
                if err := mapping.Validator(value); err != nil {
                    return fmt.Errorf("invalid %s: %w", mapping.EnvVar, err)
                }
            }
            SetNestedValue(mapping.ConfigKey, value)
        } else if mapping.Required {
            return fmt.Errorf("required environment variable %s not set", mapping.EnvVar)
        }
    }
    return nil
}
```

## 🔍 **Validation Systems**

### **1. Schema Validation**

**Configuration Schema:**
```go
type ConfigSchema struct {
    Fields map[string]FieldSchema
}

type FieldSchema struct {
    Type        string
    Required    bool
    Default     interface{}
    Validator   func(interface{}) error
    Description string
}

var schema = ConfigSchema{
    Fields: map[string]FieldSchema{
        "app.name": {
            Type:        "string",
            Required:    true,
            Description: "Application name",
        },
        "server.port": {
            Type:        "int",
            Required:    false,
            Default:     8080,
            Validator:   validatePort,
            Description: "Server port number",
        },
    },
}

func ValidateAgainstSchema(config *Config) ValidationResult {
    result := ValidationResult{Valid: true}
    
    for key, fieldSchema := range schema.Fields {
        value, exists := GetNestedValue(key)
        
        if !exists {
            if fieldSchema.Required {
                result.Valid = false
                result.Errors = append(result.Errors, 
                    fmt.Sprintf("required field %s is missing", key))
            } else if fieldSchema.Default != nil {
                SetNestedValue(key, fieldSchema.Default)
            }
            continue
        }
        
        if fieldSchema.Validator != nil {
            if err := fieldSchema.Validator(value); err != nil {
                result.Valid = false
                result.Errors = append(result.Errors, 
                    fmt.Sprintf("validation failed for %s: %v", key, err))
            }
        }
    }
    
    return result
}
```

### **2. Custom Validators**

**Common Validation Functions:**
```go
func validatePort(value interface{}) error {
    switch v := value.(type) {
    case int:
        if v < 1 || v > 65535 {
            return fmt.Errorf("port must be between 1 and 65535")
        }
    case string:
        port, err := strconv.Atoi(v)
        if err != nil {
            return fmt.Errorf("port must be a valid integer")
        }
        return validatePort(port)
    default:
        return fmt.Errorf("port must be an integer")
    }
    return nil
}

func validateURL(value interface{}) error {
    str, ok := value.(string)
    if !ok {
        return fmt.Errorf("URL must be a string")
    }
    
    if _, err := url.Parse(str); err != nil {
        return fmt.Errorf("invalid URL format: %w", err)
    }
    
    return nil
}
```

## 🎨 **Custom Help Systems**

### **1. Enhanced Help Templates**

**Rich Help Formatting:**
```go
func SetupCustomHelp() {
    // Add custom template functions
    cobra.AddTemplateFunc("StyleHeading", func(s string) string {
        return fmt.Sprintf("\033[1;36m%s\033[0m", s)
    })
    
    cobra.AddTemplateFunc("StyleCommand", func(s string) string {
        return fmt.Sprintf("\033[1;32m%s\033[0m", s)
    })
    
    customTemplate := `{{.Short | StyleHeading}}

{{.Long}}

{{if .HasExample}}{{.Example}}{{end}}

{{if .HasAvailableSubCommands}}Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{.Name | StyleCommand}} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableLocalFlags}}Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasAvailableInheritedFlags}}Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`
    
    rootCmd.SetHelpTemplate(customTemplate)
}
```

### **2. Interactive Help Mode**

**Context-Aware Assistance:**
```go
func InteractiveHelp(cmd *cobra.Command) {
    fmt.Printf("Interactive help for: %s\n", cmd.Name())
    
    // Show examples based on context
    if hasConfigFile() {
        fmt.Println("📋 Current configuration file detected")
        fmt.Printf("   Location: %s\n", getConfigFilePath())
        fmt.Printf("   Format: %s\n", config.Format)
    } else {
        fmt.Println("💡 No configuration file found. Consider running:")
        fmt.Println("   config-manager config save config.json")
    }
    
    // Show relevant next steps
    fmt.Println("\n🎯 Common next steps:")
    fmt.Println("   • config-manager config list     - View all settings")
    fmt.Println("   • config-manager validate        - Check configuration")
    fmt.Println("   • config-manager env sync        - Sync with environment")
}
```

## 🚀 **Performance Optimization**

### **1. Lazy Loading Patterns**

```go
type LazyConfig struct {
    loaded bool
    data   *Config
    mutex  sync.RWMutex
}

func (lc *LazyConfig) Get() (*Config, error) {
    lc.mutex.RLock()
    if lc.loaded {
        defer lc.mutex.RUnlock()
        return lc.data, nil
    }
    lc.mutex.RUnlock()
    
    lc.mutex.Lock()
    defer lc.mutex.Unlock()
    
    // Double-check locking
    if lc.loaded {
        return lc.data, nil
    }
    
    // Load configuration
    config, err := loadConfigFromSources()
    if err != nil {
        return nil, err
    }
    
    lc.data = config
    lc.loaded = true
    
    return lc.data, nil
}
```

### **2. Caching Strategies**

```go
type ConfigCache struct {
    cache  map[string]interface{}
    mutex  sync.RWMutex
    maxAge time.Duration
    lastLoad time.Time
}

func (cc *ConfigCache) Get(key string) (interface{}, bool) {
    cc.mutex.RLock()
    defer cc.mutex.RUnlock()
    
    // Check if cache is expired
    if time.Since(cc.lastLoad) > cc.maxAge {
        return nil, false
    }
    
    value, exists := cc.cache[key]
    return value, exists
}
```

This challenge represents the cutting edge of CLI application development, preparing you to build tools that rival the complexity and functionality of enterprise-grade CLI applications used in production environments worldwide. 