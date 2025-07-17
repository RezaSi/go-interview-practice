package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetup creates a temporary directory for testing
func setupTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "filecli_test")
	require.NoError(t, err)

	// Create some test files
	testFile1 := filepath.Join(tmpDir, "test1.txt")
	testFile2 := filepath.Join(tmpDir, "test2.txt")
	testDir := filepath.Join(tmpDir, "subdir")

	err = os.WriteFile(testFile1, []byte("test content 1"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(testFile2, []byte("test content 2"), 0644)
	require.NoError(t, err)

	err = os.Mkdir(testDir, 0755)
	require.NoError(t, err)

	return tmpDir
}

// Helper function to execute commands and capture output
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf.String(), err
}

// Test root command help
func TestRootCommandHelp(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "filecli",
		Short: "A file manager CLI tool",
		Long:  `A file manager CLI that demonstrates advanced flag and argument handling with Cobra.`,
	}

	// Add global verbose flag
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Test that the command can be created and executed without error
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "filecli", rootCmd.Use)
	assert.Contains(t, rootCmd.Short, "file manager")
}

// Test global verbose flag
func TestGlobalVerboseFlag(t *testing.T) {
	rootCmd := &cobra.Command{
		Use: "filecli",
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				cmd.Println("Verbose mode enabled")
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Test verbose flag (long form)
	output, err := executeCommand(rootCmd, "--verbose")
	assert.NoError(t, err)
	assert.Contains(t, output, "Verbose mode enabled")

	// Reset verbose flag
	verbose = false

	// Test verbose flag (short form)
	output, err = executeCommand(rootCmd, "-v")
	assert.NoError(t, err)
	assert.Contains(t, output, "Verbose mode enabled")
}

// Test list command
func TestListCommand(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	// Reset flags
	format = "table"
	verbose = false

	listCmd := &cobra.Command{
		Use:  "list [directory]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			files, err := readDirectory(dir)
			if err != nil {
				return err
			}

			if format == "json" {
				response := Response{
					Success: true,
					Data:    files,
				}
				jsonData, err := json.MarshalIndent(response, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(jsonData))
				return nil
			} else {
				cmd.Printf("%-30s %-10s %-20s %s\n", "NAME", "SIZE", "MODIFIED", "TYPE")
				cmd.Println(strings.Repeat("-", 70))
				for _, file := range files {
					fileType := "FILE"
					if file.IsDir {
						fileType = "DIR"
					}
					displayName := file.Name
					if len(displayName) > 30 {
						displayName = displayName[:27] + "..."
					}
					cmd.Printf("%-30s %-10d %-20s %s\n",
						displayName,
						file.Size,
						file.ModTime.Format("2006-01-02 15:04:05"),
						fileType,
					)
				}
			}
			return nil
		},
	}

	listCmd.Flags().StringVar(&format, "format", "table", "Output format (json, table)")

	// Test list with table format
	output, err := executeCommand(listCmd, tmpDir)
	assert.NoError(t, err)
	assert.Contains(t, output, "test1.txt")
	assert.Contains(t, output, "test2.txt")
	assert.Contains(t, output, "subdir")
	assert.Contains(t, output, "NAME")

	// Test list with JSON format
	format = "json"
	output, err = executeCommand(listCmd, "--format", "json", tmpDir)
	assert.NoError(t, err)

	var response Response
	err = json.Unmarshal([]byte(strings.TrimSpace(output)), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// Test copy command
func TestCopyCommand(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	// Reset flags
	verbose = false

	copyCmd := &cobra.Command{
		Use:  "copy <source> <destination>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			source := args[0]
			destination := args[1]

			if verbose {
				cmd.Printf("Copying %s to %s\n", source, destination)
			}

			// Simple copy implementation for testing
			data, err := os.ReadFile(source)
			if err != nil {
				return err
			}

			return os.WriteFile(destination, data, 0644)
		},
	}

	sourceFile := filepath.Join(tmpDir, "test1.txt")
	destFile := filepath.Join(tmpDir, "copied.txt")

	// Test copy command
	_, err := executeCommand(copyCmd, sourceFile, destFile)
	assert.NoError(t, err)

	// Verify file was copied
	assert.True(t, fileExists(destFile))

	// Test copy with verbose flag
	verbose = true
	destFile2 := filepath.Join(tmpDir, "copied2.txt")
	output, err := executeCommand(copyCmd, sourceFile, destFile2)
	assert.NoError(t, err)
	assert.Contains(t, output, "Copying")
}

// Test copy command argument validation
func TestCopyCommandArgs(t *testing.T) {
	copyCmd := &cobra.Command{
		Use:  "copy <source> <destination>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// Test with wrong number of arguments
	_, err := executeCommand(copyCmd, "file1.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 2 arg(s)")

	_, err = executeCommand(copyCmd, "file1.txt", "file2.txt", "file3.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 2 arg(s)")
}

// Test delete command
func TestDeleteCommand(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	// Reset flags
	force = false
	verbose = false

	deleteCmd := &cobra.Command{
		Use:  "delete <file>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				return fmt.Errorf("--force flag is required for safety")
			}

			filename := args[0]
			if verbose {
				cmd.Printf("Deleting file: %s\n", filename)
			}

			return os.Remove(filename)
		},
	}

	deleteCmd.Flags().BoolVar(&force, "force", false, "Force deletion (required)")
	err := deleteCmd.MarkFlagRequired("force")
	require.NoError(t, err)

	testFile := filepath.Join(tmpDir, "test1.txt")

	// Test delete without force flag (should fail)
	_, err = executeCommand(deleteCmd, testFile)
	assert.Error(t, err)

	// Test delete with force flag
	_, err = executeCommand(deleteCmd, "--force", testFile)
	assert.NoError(t, err)

	// Verify file was deleted
	assert.False(t, fileExists(testFile))
}

// Test create command
func TestCreateCommand(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	// Reset flags
	name = ""
	size = 0
	verbose = false

	createCmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name flag is required")
			}

			if err := validateFileName(name); err != nil {
				return err
			}

			if verbose {
				cmd.Printf("Creating file: %s with size: %d bytes\n", name, size)
			}

			filepath := filepath.Join(tmpDir, name)

			// Create file with specified size
			data := make([]byte, size)
			return os.WriteFile(filepath, data, 0644)
		},
	}

	createCmd.Flags().StringVar(&name, "name", "", "File name (required)")
	createCmd.Flags().IntVar(&size, "size", 0, "File size in bytes")
	err := createCmd.MarkFlagRequired("name")
	require.NoError(t, err)

	// Test create without name flag (should fail)
	_, err = executeCommand(createCmd)
	assert.Error(t, err)

	// Test create with name flag
	_, err = executeCommand(createCmd, "--name", "newfile.txt")
	assert.NoError(t, err)

	// Verify file was created
	newFile := filepath.Join(tmpDir, "newfile.txt")
	assert.True(t, fileExists(newFile))

	// Test create with name and size
	_, err = executeCommand(createCmd, "--name", "bigfile.txt", "--size", "100")
	assert.NoError(t, err)

	bigFile := filepath.Join(tmpDir, "bigfile.txt")
	assert.True(t, fileExists(bigFile))

	info, err := os.Stat(bigFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), info.Size())
}

// Test flag validation
func TestFlagValidation(t *testing.T) {
	// Test invalid filename
	err := validateFileName("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	err = validateFileName("file/with/slashes")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid characters")

	// Test valid filename
	err = validateFileName("validfile.txt")
	assert.NoError(t, err)
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	// Test readDirectory
	files, err := readDirectory(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, files, 3) // test1.txt, test2.txt, subdir

	// Test fileExists
	testFile := filepath.Join(tmpDir, "test1.txt")
	assert.True(t, fileExists(testFile))
	assert.False(t, fileExists("nonexistent.txt"))

	// Test formatAsJSON
	testData := []FileInfo{{Name: "test.txt", Size: 100, IsDir: false}}
	err = formatAsJSON(testData)
	assert.NoError(t, err)
}

// Test command integration
func TestCommandIntegration(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	// Create a full CLI setup
	rootCmd := &cobra.Command{Use: "filecli"}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Add all commands
	listCmd := &cobra.Command{
		Use:  "list [directory]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := tmpDir
			if len(args) > 0 {
				dir = args[0]
			}
			files, err := readDirectory(dir)
			if err != nil {
				return err
			}
			if format == "json" {
				response := Response{
					Success: true,
					Data:    files,
				}
				jsonData, err := json.MarshalIndent(response, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(jsonData))
				return nil
			}
			cmd.Printf("%-30s %-10s %-20s %s\n", "NAME", "SIZE", "MODIFIED", "TYPE")
			cmd.Println(strings.Repeat("-", 70))
			for _, file := range files {
				fileType := "FILE"
				if file.IsDir {
					fileType = "DIR"
				}
				displayName := file.Name
				if len(displayName) > 30 {
					displayName = displayName[:27] + "..."
				}
				cmd.Printf("%-30s %-10d %-20s %s\n",
					displayName,
					file.Size,
					file.ModTime.Format("2006-01-02 15:04:05"),
					fileType,
				)
			}
			return nil
		},
	}
	listCmd.Flags().StringVar(&format, "format", "table", "Output format")

	rootCmd.AddCommand(listCmd)

	// Test global flag with subcommand
	verbose = false
	format = "table"

	output, err := executeCommand(rootCmd, "-v", "list")
	assert.NoError(t, err)
	assert.Contains(t, output, "test1.txt")
}
