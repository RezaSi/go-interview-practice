# Challenge 4: Advanced Features & Middleware

Build a **Configuration Management CLI** using Cobra that demonstrates advanced CLI patterns including middleware, plugins, configuration files, and custom help systems.

## Challenge Requirements

Create a CLI application called `config-manager` that manages application configurations with:

1. **Configuration Management** - Load/save configs from multiple formats (JSON, YAML, TOML)
2. **Middleware System** - Pre/post command execution hooks
3. **Plugin Architecture** - Support for custom command plugins
4. **Environment Integration** - Environment variable support
5. **Advanced Help** - Custom help templates and documentation
6. **Validation Pipeline** - Input validation with custom validators

## Expected CLI Structure

```
config-manager                           # Root command with custom help
config-manager config get <key>          # Get configuration value
config-manager config set <key> <value> # Set configuration value  
config-manager config list               # List all configurations
config-manager config delete <key>      # Delete configuration
config-manager config load <file>       # Load config from file
config-manager config save <file>       # Save config to file
config-manager config format <format>   # Change config format (json/yaml/toml)
config-manager plugin install <name>    # Install a plugin
config-manager plugin list              # List installed plugins
config-manager validate                 # Validate current configuration
config-manager env sync                 # Sync with environment variables
config-manager completion bash          # Generate bash completion
```

## Sample Output

**Set Configuration (`config-manager config set database.host localhost`):**
```
$ config-manager config set database.host localhost
🔧 Configuration updated successfully
Key: database.host
Value: localhost
Type: string
Format: json
```

**Get Configuration (`config-manager config get database.host`):**
```
$ config-manager config get database.host
📋 Configuration Value:
Key: database.host
Value: localhost
Type: string
Source: file
Last Modified: 2024-01-15 10:30:45
```

**Load Configuration (`config-manager config load app.yaml`):**
```
$ config-manager config load app.yaml
📁 Loading configuration from app.yaml...
✅ Successfully loaded 12 configuration keys
Format: yaml
Validation: passed
```

**Plugin System (`config-manager plugin list`):**
```
$ config-manager plugin list
🔌 Installed Plugins:
Name        | Version | Status  | Description
------------|---------|---------|----------------------------------
validator   | 1.0.0   | active  | Advanced configuration validation
backup      | 0.2.1   | active  | Automatic configuration backup
generator   | 1.1.0   | active  | Configuration template generator
```

## Data Models

```go
type Config struct {
    Data     map[string]interface{} `json:"data" yaml:"data" toml:"data"`
    Format   string                 `json:"format" yaml:"format" toml:"format"`
    Version  string                 `json:"version" yaml:"version" toml:"version"`
    Metadata ConfigMetadata         `json:"metadata" yaml:"metadata" toml:"metadata"`
}

type ConfigMetadata struct {
    Created      time.Time          `json:"created" yaml:"created" toml:"created"`
    Modified     time.Time          `json:"modified" yaml:"modified" toml:"modified"`
    Source       string             `json:"source" yaml:"source" toml:"source"`
    Validation   ValidationResult   `json:"validation" yaml:"validation" toml:"validation"`
    Environment  map[string]string  `json:"environment" yaml:"environment" toml:"environment"`
}

type Plugin struct {
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Status      string            `json:"status"`
    Description string            `json:"description"`
    Commands    []PluginCommand   `json:"commands"`
    Config      map[string]string `json:"config"`
}

type ValidationResult struct {
    Valid   bool     `json:"valid"`
    Errors  []string `json:"errors"`
    Warnings []string `json:"warnings"`
}
```

## Implementation Requirements

### Configuration Management
- Support JSON, YAML, and TOML formats
- Nested key access (e.g., `database.host`, `server.port`)
- Type preservation (string, int, bool, float)
- Atomic updates with rollback capability

### Middleware System
- Pre-command validation middleware
- Post-command audit logging middleware
- Configuration backup middleware
- Performance monitoring middleware

### Plugin Architecture
- Dynamic plugin loading
- Plugin command registration
- Plugin configuration management
- Plugin lifecycle management (install/uninstall/enable/disable)

### Environment Integration
- Environment variable mapping
- Variable precedence handling
- Auto-sync capabilities
- Environment validation

### Advanced Help System
- Custom help templates
- Interactive help mode
- Example generation
- Command completion with descriptions

### Validation Pipeline
- Schema validation
- Custom validator functions
- Dependency validation
- Environment-specific validation

## Technical Requirements

### Middleware Implementation
```go
type Middleware func(*cobra.Command, []string) error

// PreRun middleware that executes before command
func ValidationMiddleware(cmd *cobra.Command, args []string) error {
    // Validate configuration before command execution
}

// PostRun middleware that executes after command
func AuditMiddleware(cmd *cobra.Command, args []string) error {
    // Log command execution for audit
}
```

### Plugin Interface
```go
type PluginInterface interface {
    Initialize() error
    GetCommands() []*cobra.Command
    GetInfo() PluginInfo
    Cleanup() error
}
```

### Configuration Format Detection
- Auto-detect format from file extension
- Content-based format detection
- Format conversion utilities
- Migration between formats

## Testing Requirements

Your solution must pass tests for:
- Configuration CRUD operations across all formats
- Middleware execution order and functionality
- Plugin loading and command registration
- Environment variable integration
- Validation pipeline with custom validators
- Help system customization
- Format conversion and migration
- Error handling and recovery
- Concurrent access protection
- Performance benchmarks

## Advanced Features

### Custom Help Templates
- Rich formatting with colors
- Interactive examples
- Context-aware help
- Multi-language support

### Performance Optimization
- Lazy loading of configurations
- Caching mechanisms
- Streaming for large configs
- Memory-efficient operations

### Security Features
- Configuration encryption
- Access control
- Audit logging
- Secure plugin loading

This challenge tests mastery of advanced Cobra patterns and demonstrates production-ready CLI application architecture. 