package main

import (
	"github.com/spf13/cobra"
)

var (
	// Application metadata
	version = "1.0.0"
	appName = "taskcli"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: "Task Manager CLI - Manage your tasks efficiently",
	Long: `Task Manager CLI - Manage your tasks efficiently

A simple and powerful command-line tool for managing your daily tasks.
Built with Go and Cobra for optimal performance and ease of use.`,
	// TODO: Implement Run function to show help by default
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Show help when no subcommand is provided
		// Hint: Use cmd.Help()
	},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of the Task Manager CLI application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Print version information
		// Expected format: "taskcli version 1.0.0\nBuilt with ❤️ using Cobra"
	},
}

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "About this application",
	Long:  `Display detailed information about the Task Manager CLI application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Print detailed application information
		// Include: version, description, author, repository, license
		// Use the expected format from the README
	},
}

func init() {
	// TODO: Add subcommands to the root command
	// Hint: Use rootCmd.AddCommand()

	// Add version command
	// TODO: rootCmd.AddCommand(versionCmd)

	// Add about command
	// TODO: rootCmd.AddCommand(aboutCmd)
}

func main() {
	// TODO: Execute the root command
	// Handle any errors that occur during execution
	// Hint: Use rootCmd.Execute()
}
