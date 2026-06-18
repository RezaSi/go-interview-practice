package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Data     map[string]interface{} `json:"data" yaml:"data" toml:"data"`
	Format   string                 `json:"format" yaml:"format" toml:"format"`
	Version  string                 `json:"version" yaml:"version" toml:"version"`
	Metadata ConfigMetadata         `json:"metadata" yaml:"metadata" toml:"metadata"`
}

// ConfigMetadata holds metadata about the configuration
type ConfigMetadata struct {
	Created     time.Time         `json:"created" yaml:"created" toml:"created"`
	Modified    time.Time         `json:"modified" yaml:"modified" toml:"modified"`
	Source      string            `json:"source" yaml:"source" toml:"source"`
	Validation  ValidationResult  `json:"validation" yaml:"validation" toml:"validation"`
	Environment map[string]string `json:"environment" yaml:"environment" toml:"environment"`
}

// ValidationResult holds validation information
type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// Plugin represents a CLI plugin
type Plugin struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Status      string            `json:"status"`
	Description string            `json:"description"`
	Commands    []PluginCommand   `json:"commands"`
	Config      map[string]string `json:"config"`
}

// PluginCommand represents a command provided by a plugin
type PluginCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

// Middleware type for command middleware
type Middleware func(*cobra.Command, []string) error

// Global configuration instance
var config *Config
var middlewares []Middleware
var plugins []Plugin

// TODO: Create the root command for the config-manager CLI
// Command name: "config-manager"
// Description: "Configuration Management CLI - Advanced configuration management with plugins and middleware"
var rootCmd = &cobra.Command{
	// TODO: Implement root command with custom help template
	Use:   "config-manager",
	Short: "Configuration Management CLI - Advanced configuration management with plugins and middleware",
	Long:  "Long. Configuration Management CLI - Advanced configuration management with plugins and middleware",
	// TODO: Add PersistentPreRun for middleware execution
	// TODO: Add PersistentPostRun for cleanup
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ApplyMiddleware(cmd, args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if err := SaveConfig(); err != nil {
			fmt.Printf("Warning: Failed to save config: %v\n", err)
		}
	},
}

// TODO: Create config parent command
// Command name: "config"
// Description: "Manage configuration settings"
var configCmd = &cobra.Command{
	// TODO: Implement config command
	Use:   "config",
	Short: "Manage configuration settings",
}

// TODO: Create config get command
// Command name: "get"
// Description: "Get configuration value by key"
// Args: configuration key (supports nested keys like "database.host")
var configGetCmd = &cobra.Command{
	// TODO: Implement config get command
	Use:   "get",
	Short: "Get configuration value by key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Get configuration value by key
		// TODO: Support nested keys
		// TODO: Display value with metadata

		key := args[0]
		value, exists := GetNestedValue(key)
		if !exists {
			cmd.Printf("❌ Key '%s' not found\n", key)
			return
		}
		cmd.Printf("📋 Configuration Value:\n")
		cmd.Printf("Key: %s\n", key)
		cmd.Printf("Value: %v\n", value)
		cmd.Printf("Type: %T\n", value)
		cmd.Printf("Source: %s\n", config.Metadata.Source)
		cmd.Printf("Last Modified: %s\n", config.Metadata.Modified.Format("2006-01-02 15:04:05"))

	},
}

// TODO: Create config set command
// Command name: "set"
// Description: "Set configuration value"
// Args: key and value
var configSetCmd = &cobra.Command{
	// TODO: Implement config set command
	Use:   "set",
	Short: "Set configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Set configuration value
		// TODO: Update metadata
		// TODO: Save configuration
		// TODO: Print success message
		key := args[0]
		value := args[1]

		var typedValue interface{} = value
		if intVal, err := strconv.Atoi(value); err == nil {
			typedValue = intVal
		} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			typedValue = floatVal
		} else if boolVal, err := strconv.ParseBool(value); err == nil {
			typedValue = boolVal
		}
		if err := SetNestedValue(key, typedValue); err != nil {
			cmd.Printf("❌ Failed to set value: %v\n", err)
			return
		}

		cmd.Printf("🔧 Configuration updated successfully\n")
		cmd.Printf("Key: %s\n", key)
		cmd.Printf("Value: %v\n", typedValue)
		cmd.Printf("Type: %T\n", typedValue)
		cmd.Printf("Format: %s\n", config.Format)
	},
}

// TODO: Create config list command
// Command name: "list"
// Description: "List all configuration keys"
var configListCmd = &cobra.Command{
	// TODO: Implement config list command
	Use:   "list",
	Short: "List all configuration keys",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Display all configuration keys in tree format
		// TODO: Show metadata for each key
		if config == nil || config.Data == nil {
			cmd.Println("No configuration loaded")
			return
		}

		cmd.Println("Configuration keys:")
		displayConfigTree(cmd, config.Data, "", 0)

		cmd.Printf("\nMetadata:\n")
		cmd.Printf("  Format: %s\n", config.Format)
		cmd.Printf("  Version: %s\n", config.Version)
		cmd.Printf("  Created: %s\n", config.Metadata.Created.Format(time.RFC3339))
		cmd.Printf("  Modified: %s\n", config.Metadata.Modified.Format(time.RFC3339))
		cmd.Printf("  Source: %s\n", config.Metadata.Source)
	},
}

// TODO: Create config delete command
// Command name: "delete"
// Description: "Delete configuration key"
// Args: configuration key
var configDeleteCmd = &cobra.Command{
	// TODO: Implement config delete command
	Use:   "delete",
	Short: "Delete configuration key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Delete configuration key
		// TODO: Update metadata
		// TODO: Save configuration
		key := args[0]

		if config == nil || config.Data == nil {
			cmd.Println("No configuration loaded")
			return
		}

		keys := strings.Split(key, ".")
		if err := deleteNestedKey(config.Data, keys); err != nil {
			cmd.Printf("Error deleting key: %v\n", err)
			return
		}

		config.Metadata.Modified = time.Now()
		config.Metadata.Source = "manual deletion"

		if err := SaveConfig(); err != nil {
			cmd.Printf("Error saving configuration: %v\n", err)
			return
		}

		cmd.Printf("Key '%s' deleted successfully\n", key)
	},
}

// TODO: Create config load command
// Command name: "load"
// Description: "Load configuration from file"
// Args: file path
// Flags: --format, --merge, --validate
var configLoadCmd = &cobra.Command{
	// TODO: Implement config load command
	Use:   "load",
	Short: "Load configuration from file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Load configuration from file
		// TODO: Auto-detect or use specified format
		// TODO: Validate configuration
		// TODO: Merge or replace existing config
		filePath := args[0]

		format, _ := cmd.Flags().GetString("format")

		loadedData, loadedFormat, err := loadConfigFromFile(filePath, format)
		if err != nil {
			cmd.Printf("Error loading configuration: %v\n", err)
			return
		}

		if format == "" {
			format = loadedFormat
		}

		if config == nil {
			config = &Config{
				Metadata: ConfigMetadata{
					Created:  time.Now(),
					Modified: time.Now(),
				},
			}
		}
		config.Data = loadedData

		config.Format = format
		config.Metadata.Modified = time.Now()
		config.Metadata.Source = filePath

		cmd.Println("Successfully loaded")

		if err := SaveConfig(); err != nil {
			cmd.Printf("Error saving configuration: %v\n", err)
		}
	},
}

// TODO: Create config save command
// Command name: "save"
// Description: "Save configuration to file"
// Args: file path
// Flags: --format, --pretty
var configSaveCmd = &cobra.Command{
	// TODO: Implement config save command
	Use:   "save",
	Short: "Save configuration to file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Save configuration to file
		// TODO: Use specified or current format
		// TODO: Pretty print if requested
		filePath := args[0]

		if config == nil {
			cmd.Println("No configuration to save")
			return
		}

		format, _ := cmd.Flags().GetString("format")
		pretty, _ := cmd.Flags().GetBool("pretty")

		if format == "" {
			format = config.Format
		}

		if err := saveConfigToFile(filePath, config.Data, format, pretty); err != nil {
			cmd.Printf("Error saving configuration: %v\n", err)
			return
		}

		config.Metadata.Modified = time.Now()
		config.Metadata.Source = filePath

		cmd.Printf("Configuration saved successfully to %s in %s format\n", filePath, format)
	},
}

// TODO: Create config format command
// Command name: "format"
// Description: "Change configuration format"
// Args: format (json/yaml/toml)
var configFormatCmd = &cobra.Command{
	// TODO: Implement config format command
	Use:   "format",
	Short: "Change configuration format",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Convert configuration to new format
		// TODO: Update metadata
		// TODO: Save configuration
		newFormat := strings.ToLower(args[0])

		if config == nil {
			cmd.Println("No configuration loaded")
			return
		}

		validFormats := map[string]bool{"json": true, "yaml": true, "toml": true}
		if !validFormats[newFormat] {
			cmd.Printf("Invalid format: %s. Supported formats: json, yaml, toml\n", newFormat)
			return
		}

		config.Format = newFormat

		config.Metadata.Modified = time.Now()
		config.Metadata.Source = "format conversion"

		if err := SaveConfig(); err != nil {
			cmd.Printf("Error saving configuration: %v\n", err)
			return
		}

		cmd.Printf("Configuration format changed to %s\n", newFormat)
	},
}

// TODO: Create plugin parent command
// Command name: "plugin"
// Description: "Manage CLI plugins"
var pluginCmd = &cobra.Command{
	// TODO: Implement plugin command
	Use:   "plugin",
	Short: "Manage CLI plugins",
}

// TODO: Create plugin install command
// Command name: "install"
// Description: "Install a plugin"
// Args: plugin name
var pluginInstallCmd = &cobra.Command{
	// TODO: Implement plugin install command
	Use:   "install",
	Short: "Install a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Install plugin
		// TODO: Register plugin commands
		// TODO: Update plugin registry
		pluginName := args[0]

		plugin, err := installPlugin(pluginName)
		if err != nil {
			cmd.Printf("Error installing plugin: %v\n", err)
			return
		}

		if err := RegisterPlugin(plugin); err != nil {
			cmd.Printf("Error registering plugin commands: %v\n", err)
			return
		}

		cmd.Printf("Plugin '%s' installed successfully\n", pluginName)
	},
}

// TODO: Create plugin list command
// Command name: "list"
// Description: "List installed plugins"
var pluginListCmd = &cobra.Command{
	// TODO: Implement plugin list command
	Use:   "list",
	Short: "List installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Display installed plugins in table format
		// TODO: Show status and version information
		if len(plugins) == 0 {
			cmd.Println("No plugins installed")
			return
		}

		cmd.Printf("Name | Version | Status | Description\n")
		for _, p := range plugins {
			cmd.Printf("%s | %s | %s | %s\n", p.Name, p.Version, p.Status, p.Description)
		}
	},
}

// TODO: Create validate command
// Command name: "validate"
// Description: "Validate current configuration"
var validateCmd = &cobra.Command{
	// TODO: Implement validate command
	Use:   "validate",
	Short: "Validate current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Run validation pipeline
		// TODO: Display validation results
		// TODO: Show errors and warnings
		if config == nil {
			cmd.Println("No configuration to validate")
			return
		}

		result := ValidateConfiguration()
		if result.Valid {
			cmd.Println("✓ validation: passed")
		} else {
			cmd.Printf("✗ validation: failed\n")

			if len(result.Errors) > 0 {
				cmd.Println("\nErrors:")
				for i, err := range result.Errors {
					cmd.Printf("  %d. %s\n", i+1, err)
				}
			}

			if len(result.Warnings) > 0 {
				cmd.Println("\nWarnings:")
				for i, warning := range result.Warnings {
					cmd.Printf("  %d. %s\n", i+1, warning)
				}
			}
		}

		if result.Valid {
			config.Metadata.Modified = time.Now()
			if err := SaveConfig(); err != nil {
				cmd.Printf("Error saving configuration: %v\n", err)
			}
		}

	},
}

// TODO: Create env parent command
// Command name: "env"
// Description: "Environment variable integration"
var envCmd = &cobra.Command{
	// TODO: Implement env command
	Use:   "env",
	Short: "Environment variable integration",
}

// TODO: Create env sync command
// Command name: "sync"
// Description: "Sync configuration with environment variables"
var envSyncCmd = &cobra.Command{
	// TODO: Implement env sync command
	Use:   "sync",
	Short: "Sync configuration with environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Sync with environment variables
		// TODO: Apply precedence rules
		// TODO: Update configuration
		if config == nil {
			config = &Config{
				Data: make(map[string]interface{}),
				Metadata: ConfigMetadata{
					Created:  time.Now(),
					Modified: time.Now(),
				},
			}
		}

		envVars := os.Environ()
		config.Metadata.Environment = make(map[string]string)

		for _, env := range envVars {
			pair := strings.SplitN(env, "=", 2)
			if len(pair) == 2 {
				key, value := pair[0], pair[1]
				config.Metadata.Environment[key] = value

				if strings.HasPrefix(key, "CONFIG_") {
					configKey := strings.ToLower(strings.TrimPrefix(key, "CONFIG_"))
					config.Data[configKey] = value
				}
			}
		}

		config.Metadata.Modified = time.Now()
		config.Metadata.Source = "environment sync"

		if err := SaveConfig(); err != nil {
			cmd.Printf("Error saving configuration: %v\n", err)
			return
		}

		cmd.Println("Configuration synchronized with environment variables")
	},
}

// LoadConfig loads configuration from default location or creates new
func LoadConfig() error {
	// TODO: Implement loading configuration
	// TODO: Create default config if not exists
	// TODO: Handle different formats
	// Implement loading configuration
	configPath := getConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {

		config = &Config{
			Data:    make(map[string]interface{}),
			Format:  "yaml",
			Version: "1.0",
			Metadata: ConfigMetadata{
				Created:     time.Now(),
				Modified:    time.Now(),
				Source:      "default",
				Validation:  ValidationResult{Valid: true},
				Environment: make(map[string]string),
			},
		}
		return SaveConfig()
	}

	loadedData, format, err := loadConfigFromFile(configPath, "")
	if err != nil {
		return err
	}

	config = &Config{
		Data:   loadedData,
		Format: format,
		Metadata: ConfigMetadata{
			Modified: time.Now(),
			Source:   configPath,
		},
	}

	return nil

}

// SaveConfig saves configuration to default location
func SaveConfig() error {
	// TODO: Implement saving configuration
	// TODO: Use current format
	// TODO: Update metadata
	if config == nil {
		return fmt.Errorf("no configuration to save")
	}

	configPath := getConfigPath()
	return saveConfigToFile(configPath, config.Data, config.Format, true)
}

// GetNestedValue retrieves value from nested configuration key
func GetNestedValue(key string) (interface{}, bool) {
	// TODO: Implement nested key access
	// TODO: Support dot notation (e.g., "database.host")
	if config == nil || config.Data == nil {
		return nil, false
	}
	parts := strings.Split(key, ".")
	current := config.Data
	for i, part := range parts {
		if value, exists := current[part]; exists {
			if i == len(parts)-1 {

				return value, true
			}

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

// SetNestedValue sets value for nested configuration key
func SetNestedValue(key string, value interface{}) error {
	// TODO: Implement nested key setting
	// TODO: Create intermediate keys if needed
	// TODO: Update metadata
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

	current[parts[len(parts)-1]] = value
	config.Metadata.Modified = time.Now()

	return nil
}

// ValidateConfiguration runs validation pipeline
func ValidateConfiguration() ValidationResult {
	// TODO: Implement configuration validation
	// TODO: Run custom validators
	// TODO: Check dependencies
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

	requiredFields := []string{"app.name", "app.version"}
	for _, field := range requiredFields {
		if _, exists := GetNestedValue(field); !exists {
			result.Warnings = append(result.Warnings, fmt.Sprintf("recommended field %s is missing", field))
		}
	}

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

// ApplyMiddleware executes all registered middleware
func ApplyMiddleware(cmd *cobra.Command, args []string) error {
	// TODO: Execute all middleware in order
	// TODO: Handle middleware errors
	for _, middleware := range middlewares {
		if err := middleware(cmd, args); err != nil {
			return fmt.Errorf("middleware failed: %w", err)
		}
	}
	return nil
}

// RegisterPlugin registers a new plugin
func RegisterPlugin(plugin Plugin) error {
	// TODO: Register plugin commands
	// TODO: Add to plugin registry
	// TODO: Initialize plugin

	for _, existing := range plugins {
		if existing.Name == plugin.Name {
			return fmt.Errorf("plugin %s already registered", plugin.Name)
		}
	}

	plugins = append(plugins, plugin)
	fmt.Printf("✅ Plugin '%s' v%s registered successfully\n", plugin.Name, plugin.Version)
	return nil

}

// DetectFormat auto-detects configuration format
func DetectFormat(filename string) string {
	// TODO: Detect format from file extension
	// TODO: Fall back to content detection
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".json":
		return "json"
	default:
		return "json"
	}
}

// ConvertFormat converts configuration to specified format
func ConvertFormat(targetFormat string) error {
	// TODO: Convert configuration data
	// TODO: Update format metadata
	// TODO: Preserve data integrity
	if config.Format == targetFormat {
		return nil
	}

	config.Format = targetFormat
	config.Metadata.Modified = time.Now()
	return nil
}

// ValidationMiddleware validates configuration before command execution
func ValidationMiddleware(cmd *cobra.Command, args []string) error {
	// TODO: Implement validation middleware
	result := ValidateConfiguration()
	if !result.Valid && len(result.Errors) > 0 {
		fmt.Printf("⚠️  Configuration warnings: %v\n", result.Warnings)
	}
	return nil
}

// AuditMiddleware logs command execution for audit
func AuditMiddleware(cmd *cobra.Command, args []string) error {
	// TODO: Implement audit logging
	// TODO: Log command, args, timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("🔍 [%s] Executing: %s %v\n", timestamp, cmd.Name(), args)
	return nil
}

// SetCustomHelpTemplate sets up custom help formatting
func SetCustomHelpTemplate() {
	// TODO: Define custom help template with colors and formatting
	// TODO: Add examples and interactive elements
	helpTemplate := `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}
{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
	cobra.AddTemplateFunc("StyleHeading", func(s string) string {
		return fmt.Sprintf("\033[1;36m%s\033[0m", s) // Cyan bold
	})
	rootCmd.SetHelpTemplate(helpTemplate)
}

func init() {

	// Initialize viper for configuration management
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config-manager")
	viper.AutomaticEnv()

	// Register middleware
	middlewares = append(middlewares, ValidationMiddleware, AuditMiddleware)

	// Add flags to commands
	configLoadCmd.Flags().String("format", "", "Configuration format (json/yaml/toml)")
	configLoadCmd.Flags().Bool("merge", false, "Merge with existing configuration")
	configLoadCmd.Flags().Bool("validate", true, "Validate configuration after loading")

	configSaveCmd.Flags().String("format", "", "Configuration format (json/yaml/toml)")
	configSaveCmd.Flags().Bool("pretty", true, "Pretty print output")

	// Add subcommands to config command
	configCmd.AddCommand(configGetCmd, configSetCmd, configListCmd, configDeleteCmd, configLoadCmd, configSaveCmd, configFormatCmd)

	// Add subcommands to plugin command
	pluginCmd.AddCommand(pluginInstallCmd, pluginListCmd)

	// Add subcommands to env command
	envCmd.AddCommand(envSyncCmd)

	// Add all commands to root command
	rootCmd.AddCommand(configCmd, pluginCmd, validateCmd, envCmd)

	SetCustomHelpTemplate()

	// Load configuration on startup
	if err := LoadConfig(); err != nil {
		log.Printf("Warning: Could not load configuration: %v", err)
	}
}

func main() {

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if err := ApplyMiddleware(cmd, args); err != nil {
			fmt.Printf("Middleware error: %v\n", err)
			os.Exit(1)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Config
// =========================================================================================
func displayConfigTree(cmd *cobra.Command, data interface{}, prefix string, depth int) {

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			newPrefix := prefix
			if depth > 0 {
				newPrefix += "  "
			}
			cmd.Printf("%s%s:\n", newPrefix, key)
			displayConfigTree(cmd, value, newPrefix, depth+1)
		}
	case []interface{}:
		for i, item := range v {
			cmd.Printf("%s[%d]: %v\n", prefix, i, item)
		}
	default:
		cmd.Printf("%s%v\n", prefix, v)
	}

}

func deleteNestedKey(data map[string]interface{}, keys []string) error {
	if len(keys) == 1 {
		delete(data, keys[0])
		return nil
	}

	next, exists := data[keys[0]]
	if !exists {
		return fmt.Errorf("key not found: %s", keys[0])
	}

	nextMap, ok := next.(map[string]interface{})
	if !ok {
		return fmt.Errorf("key is not a map: %s", keys[0])
	}

	return deleteNestedKey(nextMap, keys[1:])
}

// loadConfigFromFile
func loadConfigFromFile(filePath, format string) (map[string]interface{}, string, error) {

	if format == "" {
		format = DetectFormat(strings.ToLower(filepath.Ext(filePath)))
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("error reading file: %w", err)
	}

	var configData map[string]interface{}

	switch format {
	case "json":
		if err := json.Unmarshal(fileData, &configData); err != nil {
			return nil, "", fmt.Errorf("error parsing JSON: %w", err)
		}
	case "yaml":
		if err := yaml.Unmarshal(fileData, &configData); err != nil {
			return nil, "", fmt.Errorf("error parsing YAML: %w", err)
		}
	case "toml":
		var tomlData interface{}
		if err := toml.Unmarshal(fileData, &tomlData); err != nil {
			return nil, "", fmt.Errorf("error parsing TOML: %w", err)
		}

		configData = convertTomlToMap(tomlData)
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}

	return configData, format, nil
}

func convertTomlToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if table, ok := data.(map[string]interface{}); ok {
		for k, v := range table {
			if subTable, ok := v.(map[string]interface{}); ok {
				result[k] = convertTomlToMap(subTable)
			} else {
				result[k] = v
			}
		}
	}
	return result
}

// saveConfigToFile
func saveConfigToFile(filePath string, data map[string]interface{}, format string, pretty bool) error {
	if data == nil {
		return fmt.Errorf("no data to save")
	}

	var fileData []byte
	var err error

	switch format {
	case "json":
		if pretty {
			fileData, err = json.MarshalIndent(data, "", "  ")
		} else {
			fileData, err = json.Marshal(data)
		}
		if err != nil {
			return fmt.Errorf("error marshaling JSON: %w", err)
		}
	case "yaml":
		fileData, err = yaml.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshaling YAML: %w", err)
		}
	case "toml":
		fileData, err = toml.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshaling TOML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

// installPlugin
func installPlugin(pluginName string) (Plugin, error) {

	for _, p := range plugins {
		if p.Name == pluginName {
			return Plugin{}, fmt.Errorf("plugin '%s' is already installed", pluginName)
		}
	}

	plugin := Plugin{
		Name:        pluginName,
		Version:     "1.0.0",
		Status:      "active",
		Description: fmt.Sprintf("Plugin for %s functionality", pluginName),
		Config: map[string]string{
			"enabled": "true",
			"path":    fmt.Sprintf("/plugins/%s", pluginName),
		},
		Commands: []PluginCommand{
			{
				Name:        fmt.Sprintf("%s-run", pluginName),
				Description: fmt.Sprintf("Run %s operation", pluginName),
				Usage:       fmt.Sprintf("config-manager %s-run [options]", pluginName),
			},
		},
	}

	pluginDir := fmt.Sprintf("./plugins/%s", pluginName)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return Plugin{}, fmt.Errorf("error creating plugin directory: %w", err)
	}

	manifestData, err := json.MarshalIndent(plugin, "", "  ")
	if err != nil {
		return Plugin{}, fmt.Errorf("error creating plugin manifest: %w", err)
	}

	if err := os.WriteFile(filepath.Join(pluginDir, "manifest.json"), manifestData, 0644); err != nil {
		return Plugin{}, fmt.Errorf("error writing plugin manifest: %w", err)
	}

	return plugin, nil
}

// getConfigPath
func getConfigPath() string {

	if envPath := os.Getenv("CONFIG_MANAGER_PATH"); envPath != "" {
		return envPath
	}

	if _, err := os.Stat("./config.yaml"); err == nil {
		return "./config.yaml"
	}
	if _, err := os.Stat("./config.json"); err == nil {
		return "./config.json"
	}
	if _, err := os.Stat("./config.toml"); err == nil {
		return "./config.toml"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./config.yaml"
	}

	return filepath.Join(homeDir, ".config-manager", "config.yaml")
}
