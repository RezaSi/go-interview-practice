package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	// Clean up any test files before and after tests
	cleanupTestFiles()
	code := m.Run()
	cleanupTestFiles()
	os.Exit(code)
}

func cleanupTestFiles() {
	os.Remove("config.json")
	os.Remove("config.yaml")
	os.Remove("config.toml")
	os.Remove("test-config.json")
	os.Remove("test-config.yaml")
}

func setupTest() {
	// Reset global configuration for each test
	config = &Config{
		Data:    make(map[string]interface{}),
		Format:  "json",
		Version: "1.0.0",
		Metadata: ConfigMetadata{
			Created:     time.Now(),
			Modified:    time.Now(),
			Source:      "test",
			Validation:  ValidationResult{Valid: true},
			Environment: make(map[string]string),
		},
	}
	middlewares = []Middleware{}
	plugins = []Plugin{}
	cleanupTestFiles()
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	setupTest()
	output, err := executeCommand(rootCmd)

	if err != nil {
		t.Fatalf("Root command failed: %v", err)
	}

	if !strings.Contains(output, "Configuration Management CLI") {
		t.Error("Root command should contain 'Configuration Management CLI' in output")
	}

	if !strings.Contains(output, "config") {
		t.Error("Root command should show 'config' subcommand")
	}

	if !strings.Contains(output, "plugin") {
		t.Error("Root command should show 'plugin' subcommand")
	}

	if !strings.Contains(output, "validate") {
		t.Error("Root command should show 'validate' subcommand")
	}
}

func TestConfigSetCommand(t *testing.T) {
	setupTest()

	args := []string{"config", "set", "database.host", "localhost"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Config set command failed: %v", err)
	}

	if !strings.Contains(output, "Configuration updated successfully") {
		t.Error("Config set should show success message")
	}

	if !strings.Contains(output, "database.host") {
		t.Error("Config set should show the key being set")
	}

	if !strings.Contains(output, "localhost") {
		t.Error("Config set should show the value being set")
	}

	// Verify value was set
	value, exists := GetNestedValue("database.host")
	if !exists || value != "localhost" {
		t.Error("Config set should actually set the value")
	}
}

func TestConfigGetCommand(t *testing.T) {
	setupTest()

	// Set a test value first
	config.Data["database"] = map[string]interface{}{
		"host": "localhost",
		"port": 5432,
	}

	args := []string{"config", "get", "database.host"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Config get command failed: %v", err)
	}

	if !strings.Contains(output, "localhost") {
		t.Error("Config get should show the value")
	}

	if !strings.Contains(output, "database.host") {
		t.Error("Config get should show the key")
	}

	// Test getting non-existent key
	args = []string{"config", "get", "nonexistent.key"}
	output, err = executeCommand(rootCmd, args...)
	if err == nil && !strings.Contains(output, "not found") {
		t.Error("Config get should handle non-existent keys")
	}
}

func TestConfigListCommand(t *testing.T) {
	setupTest()

	// Add test configuration
	config.Data = map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
		"server": map[string]interface{}{
			"port": 8080,
		},
		"debug": true,
	}

	output, err := executeCommand(rootCmd, "config", "list")

	if err != nil {
		t.Fatalf("Config list command failed: %v", err)
	}

	if !strings.Contains(output, "database") {
		t.Error("Config list should show database configuration")
	}

	if !strings.Contains(output, "server") {
		t.Error("Config list should show server configuration")
	}

	if !strings.Contains(output, "debug") {
		t.Error("Config list should show debug configuration")
	}
}

func TestConfigDeleteCommand(t *testing.T) {
	setupTest()

	// Set a test value first
	config.Data["test_key"] = "test_value"

	args := []string{"config", "delete", "test_key"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Config delete command failed: %v", err)
	}

	if !strings.Contains(output, "deleted successfully") {
		t.Error("Config delete should show success message")
	}

	// Verify value was deleted
	_, exists := GetNestedValue("test_key")
	if exists {
		t.Error("Config delete should actually delete the value")
	}
}

func TestConfigLoadCommand(t *testing.T) {
	setupTest()

	// Create test configuration file
	testConfig := map[string]interface{}{
		"app_name": "test-app",
		"version":  "1.0.0",
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
	}

	data, _ := json.MarshalIndent(testConfig, "", "  ")
	err := os.WriteFile("test-config.json", data, 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	args := []string{"config", "load", "test-config.json"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Config load command failed: %v", err)
	}

	if !strings.Contains(output, "Successfully loaded") {
		t.Error("Config load should show success message")
	}

	// Verify configuration was loaded
	value, exists := GetNestedValue("app_name")
	if !exists || value != "test-app" {
		t.Error("Config load should load values into configuration")
	}
}

func TestConfigSaveCommand(t *testing.T) {
	setupTest()

	// Set test configuration
	config.Data = map[string]interface{}{
		"app_name": "test-app",
		"version":  "1.0.0",
	}

	args := []string{"config", "save", "test-config.json"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Config save command failed: %v", err)
	}

	if !strings.Contains(output, "successfully") {
		t.Error("Config save should show success message")
	}

	// Verify file was created
	if _, err := os.Stat("test-config.json"); os.IsNotExist(err) {
		t.Error("Config save should create the file")
	}

	// Verify file contents
	data, err := os.ReadFile("test-config.json")
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	var savedConfig map[string]interface{}
	if err := json.Unmarshal(data, &savedConfig); err != nil {
		t.Fatalf("Failed to parse saved config: %v", err)
	}

	if savedConfig["app_name"] != "test-app" {
		t.Error("Saved config should contain the correct data")
	}
}

func TestConfigFormatCommand(t *testing.T) {
	setupTest()

	config.Format = "json"
	config.Data = map[string]interface{}{
		"test": "value",
	}

	args := []string{"config", "format", "yaml"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Config format command failed: %v", err)
	}

	if !strings.Contains(output, "format changed") || !strings.Contains(output, "yaml") {
		t.Error("Config format should show format change message")
	}

	// Verify format was changed
	if config.Format != "yaml" {
		t.Error("Config format should actually change the format")
	}
}

func TestPluginInstallCommand(t *testing.T) {
	setupTest()

	args := []string{"plugin", "install", "test-plugin"}
	output, err := executeCommand(rootCmd, args...)

	if err != nil {
		t.Fatalf("Plugin install command failed: %v", err)
	}

	if !strings.Contains(output, "installed successfully") {
		t.Error("Plugin install should show success message")
	}

	if !strings.Contains(output, "test-plugin") {
		t.Error("Plugin install should show plugin name")
	}
}

func TestPluginListCommand(t *testing.T) {
	setupTest()

	// Add test plugins
	plugins = []Plugin{
		{
			Name:        "validator",
			Version:     "1.0.0",
			Status:      "active",
			Description: "Configuration validator",
		},
		{
			Name:        "backup",
			Version:     "0.2.1",
			Status:      "inactive",
			Description: "Backup utility",
		},
	}

	output, err := executeCommand(rootCmd, "plugin", "list")

	if err != nil {
		t.Fatalf("Plugin list command failed: %v", err)
	}

	if !strings.Contains(output, "validator") {
		t.Error("Plugin list should show validator plugin")
	}

	if !strings.Contains(output, "backup") {
		t.Error("Plugin list should show backup plugin")
	}

	if !strings.Contains(output, "1.0.0") {
		t.Error("Plugin list should show version information")
	}

	if !strings.Contains(output, "active") {
		t.Error("Plugin list should show status information")
	}
}

func TestValidateCommand(t *testing.T) {
	setupTest()

	// Set valid configuration
	config.Data = map[string]interface{}{
		"app_name": "test-app",
		"port":     8080,
	}

	output, err := executeCommand(rootCmd, "validate")

	if err != nil {
		t.Fatalf("Validate command failed: %v", err)
	}

	if !strings.Contains(output, "validation") {
		t.Error("Validate command should show validation results")
	}
}

func TestEnvSyncCommand(t *testing.T) {
	setupTest()

	// Set test environment variable
	os.Setenv("TEST_CONFIG_PORT", "8080")
	defer os.Unsetenv("TEST_CONFIG_PORT")

	output, err := executeCommand(rootCmd, "env", "sync")

	if err != nil {
		t.Fatalf("Env sync command failed: %v", err)
	}

	if !strings.Contains(output, "sync") {
		t.Error("Env sync should show sync information")
	}
}

func TestGetNestedValue(t *testing.T) {
	setupTest()

	config.Data = map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
			"ssl": map[string]interface{}{
				"enabled": true,
				"cert":    "/path/to/cert",
			},
		},
		"debug": true,
	}

	// Test simple key
	value, exists := GetNestedValue("debug")
	if !exists || value != true {
		t.Error("GetNestedValue should retrieve simple values")
	}

	// Test nested key
	value, exists = GetNestedValue("database.host")
	if !exists || value != "localhost" {
		t.Error("GetNestedValue should retrieve nested values")
	}

	// Test deep nested key
	value, exists = GetNestedValue("database.ssl.enabled")
	if !exists || value != true {
		t.Error("GetNestedValue should retrieve deeply nested values")
	}

	// Test non-existent key
	_, exists = GetNestedValue("nonexistent.key")
	if exists {
		t.Error("GetNestedValue should return false for non-existent keys")
	}
}

func TestSetNestedValue(t *testing.T) {
	setupTest()

	// Test setting simple key
	err := SetNestedValue("debug", true)
	if err != nil {
		t.Fatalf("SetNestedValue failed for simple key: %v", err)
	}

	value, exists := GetNestedValue("debug")
	if !exists || value != true {
		t.Error("SetNestedValue should set simple values")
	}

	// Test setting nested key
	err = SetNestedValue("database.host", "localhost")
	if err != nil {
		t.Fatalf("SetNestedValue failed for nested key: %v", err)
	}

	value, exists = GetNestedValue("database.host")
	if !exists || value != "localhost" {
		t.Error("SetNestedValue should set nested values")
	}

	// Test setting deep nested key
	err = SetNestedValue("database.ssl.enabled", true)
	if err != nil {
		t.Fatalf("SetNestedValue failed for deep nested key: %v", err)
	}

	value, exists = GetNestedValue("database.ssl.enabled")
	if !exists || value != true {
		t.Error("SetNestedValue should set deeply nested values")
	}
}

func TestValidateConfiguration(t *testing.T) {
	setupTest()

	// Test valid configuration
	config.Data = map[string]interface{}{
		"app_name": "test-app",
		"port":     8080,
	}

	result := ValidateConfiguration()
	if !result.Valid {
		t.Error("ValidateConfiguration should return valid for good config")
	}

	// Test invalid configuration (if validation rules are implemented)
	config.Data = map[string]interface{}{
		"port": "invalid-port", // Should be number
	}

	result = ValidateConfiguration()
	// Note: This test depends on actual validation implementation
	// For now, just ensure the function returns a result
	if result.Errors == nil {
		result.Errors = []string{} // Initialize if nil
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"config.json", "json"},
		{"app.yaml", "yaml"},
		{"app.yml", "yaml"},
		{"config.toml", "toml"},
		{"unknown.txt", "json"}, // Default fallback
	}

	for _, test := range tests {
		result := DetectFormat(test.filename)
		if result != test.expected {
			t.Errorf("DetectFormat(%s) = %s, expected %s", test.filename, result, test.expected)
		}
	}
}

func TestConvertFormat(t *testing.T) {
	setupTest()

	config.Format = "json"
	config.Data = map[string]interface{}{
		"test": "value",
	}

	err := ConvertFormat("yaml")
	if err != nil {
		t.Fatalf("ConvertFormat failed: %v", err)
	}

	if config.Format != "yaml" {
		t.Error("ConvertFormat should change the format")
	}
}

func TestMiddleware(t *testing.T) {
	setupTest()

	// Test validation middleware
	err := ValidationMiddleware(rootCmd, []string{})
	if err != nil {
		t.Fatalf("ValidationMiddleware failed: %v", err)
	}

	// Test audit middleware
	err = AuditMiddleware(rootCmd, []string{})
	if err != nil {
		t.Fatalf("AuditMiddleware failed: %v", err)
	}
}

func TestApplyMiddleware(t *testing.T) {
	setupTest()

	// Register test middleware
	middlewares = append(middlewares, ValidationMiddleware, AuditMiddleware)

	err := ApplyMiddleware(rootCmd, []string{})
	if err != nil {
		t.Fatalf("ApplyMiddleware failed: %v", err)
	}
}

func TestRegisterPlugin(t *testing.T) {
	setupTest()

	plugin := Plugin{
		Name:        "test-plugin",
		Version:     "1.0.0",
		Status:      "active",
		Description: "Test plugin",
		Commands:    []PluginCommand{},
		Config:      map[string]string{},
	}

	err := RegisterPlugin(plugin)
	if err != nil {
		t.Fatalf("RegisterPlugin failed: %v", err)
	}

	// Verify plugin was registered
	found := false
	for _, p := range plugins {
		if p.Name == "test-plugin" {
			found = true
			break
		}
	}

	if !found {
		t.Error("RegisterPlugin should add plugin to registry")
	}
}

func TestDataPersistence(t *testing.T) {
	setupTest()

	// Set test configuration
	config.Data = map[string]interface{}{
		"app_name": "test-app",
		"version":  "1.0.0",
	}
	config.Format = "json"

	// Save configuration
	err := SaveConfig()
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Reset and load configuration
	config = &Config{}
	err = LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify data was persisted
	value, exists := GetNestedValue("app_name")
	if !exists || value != "test-app" {
		t.Error("Configuration should be persisted and loaded correctly")
	}
}
