package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	appName = "taskcli"
)

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: "Task Manager CLI - Manage your tasks efficiently",
	Long: `Task Manager CLI - Manage your tasks efficiently

A simple and powerful command-line tool for managing your daily tasks.
Built with Go and Cobra for optimal performance and ease of use.`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Show version information. Display the current version.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("%s version %s\n", appName, version)
		cmd.Println("Built with ❤️ using Cobra")
	},
}

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "About this application",
	Long:  `About this application. Display detailed information.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("Task Manager CLI v%s", version)
		cmd.Println("A simple and efficient task management tool")
		cmd.Println("Author: author")
		cmd.Println("Repository: repository")
		cmd.Println("License: license")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(aboutCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
