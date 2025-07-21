package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Global flags
var (
	verbose bool
)

// Command-specific flags
var (
	format string
	force  bool
	name   string
	size   int
)

// FileInfo represents a file in our system
type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool      `json:"is_dir"`
}

// Response represents command output
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	// Create root command
	var rootCmd = &cobra.Command{
		Use:   "filecli",
		Short: "A file manager CLI tool",
		Long:  `A file manager CLI that demonstrates advanced flag and argument handling with Cobra.`,
	}

	// Add global flags to root command
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Add subcommands to root
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(createCmd)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// List command implementation
var listCmd = &cobra.Command{
	Use:   "list [directory]",
	Short: "List files in a directory",
	Long:  `List files and directories with optional formatting options.`,
	Args:  cobra.MaximumNArgs(1), // Optional directory argument
	RunE:  listFiles,
}

func init() {
	// Add command-specific flags to listCmd
	listCmd.Flags().StringVar(&format, "format", "table", "Output format (json, table)")
}

// Copy command implementation
var copyCmd = &cobra.Command{
	Use:   "copy <source> <destination>",
	Short: "Copy files or directories",
	Long:  `Copy a file from source to destination.`,
	Args:  cobra.ExactArgs(2), // Requires exactly 2 arguments
	RunE:  copyFile,
}

// Delete command implementation
var deleteCmd = &cobra.Command{
	Use:   "delete <file>",
	Short: "Delete a file",
	Long:  `Delete a file with safety confirmation via --force flag.`,
	Args:  cobra.ExactArgs(1), // Requires exactly 1 argument
	RunE:  deleteFile,
}

func init() {
	// Add required --force flag to deleteCmd
	deleteCmd.Flags().BoolVar(&force, "force", false, "Force deletion (required)")
	deleteCmd.MarkFlagRequired("force")
}

// Create command implementation
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new file",
	Long:  `Create a new file with specified name and size.`,
	RunE:  createFile,
}

func init() {
	// Add flags to createCmd
	createCmd.Flags().StringVar(&name, "name", "", "File name (required)")
	createCmd.Flags().IntVar(&size, "size", 0, "File size in bytes")
	createCmd.MarkFlagRequired("name")
}

// Command handler functions
func listFiles(cmd *cobra.Command, args []string) error {
	// Get directory from args or use current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	if verbose {
		fmt.Printf("Listing files in directory: %s\n", dir)
	}

	// Read directory contents
	files, err := readDirectory(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	// Format output based on --format flag
	if format == "json" {
		return formatAsJSON(files)
	} else {
		formatAsTable(files)
	}

	return nil
}

func copyFile(cmd *cobra.Command, args []string) error {
	source := args[0]
	destination := args[1]

	if verbose {
		fmt.Printf("Copying %s to %s\n", source, destination)
	}

	// Check if source file exists
	if !fileExists(source) {
		return fmt.Errorf("source file %s does not exist", source)
	}

	// Read source file
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", source, err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destination, err)
	}
	defer destFile.Close()

	// Copy contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	if verbose {
		fmt.Printf("Successfully copied %s to %s\n", source, destination)
	}

	return nil
}

func deleteFile(cmd *cobra.Command, args []string) error {
	filename := args[0]

	if verbose {
		fmt.Printf("Deleting file: %s\n", filename)
	}

	// Check if file exists
	if !fileExists(filename) {
		return fmt.Errorf("file %s does not exist", filename)
	}

	// Delete the file
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filename, err)
	}

	if verbose {
		fmt.Printf("Successfully deleted: %s\n", filename)
	}

	return nil
}

func createFile(cmd *cobra.Command, args []string) error {
	// Validate filename
	if err := validateFileName(name); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Creating file: %s with size: %d bytes\n", name, size)
	}

	// Create file with specified size
	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", name, err)
	}
	defer file.Close()

	// Write data to reach the specified size
	if size > 0 {
		data := make([]byte, size)
		// Fill with spaces to make it readable
		for i := range data {
			data[i] = ' '
		}

		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf("failed to write data to file %s: %w", name, err)
		}
	}

	if verbose {
		fmt.Printf("Successfully created file: %s (%d bytes)\n", name, size)
	}

	return nil
}

// Helper functions

// formatAsJSON formats the response as JSON
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

// formatAsTable formats file info as a table
func formatAsTable(files []FileInfo) {
	fmt.Printf("%-30s %-10s %-20s %s\n", "NAME", "SIZE", "MODIFIED", "TYPE")
	fmt.Println(strings.Repeat("-", 70))

	for _, file := range files {
		fileType := "FILE"
		if file.IsDir {
			fileType = "DIR"
		}

		// Truncate long names
		displayName := file.Name
		if len(displayName) > 30 {
			displayName = displayName[:27] + "..."
		}

		fmt.Printf("%-30s %-10d %-20s %s\n",
			displayName,
			file.Size,
			file.ModTime.Format("2006-01-02 15:04:05"),
			fileType,
		)
	}
}

// readDirectory reads files from directory and returns FileInfo slice
func readDirectory(dirPath string) ([]FileInfo, error) {
	var files []FileInfo

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// validateFileName checks if filename is valid
func validateFileName(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	if strings.ContainsAny(filename, "/\\:*?\"<>|") {
		return fmt.Errorf("filename contains invalid characters")
	}

	return nil
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
