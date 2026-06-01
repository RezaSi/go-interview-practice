package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
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
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Fprintf(os.Stdout, "Error displaing help: %v\n", err)
			return
		}
	},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Show version information Display the current version of the Task Manager CLI application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("%s version %s\n", appName, version)
		cmd.Println("Built with ❤️ using Cobra")
	},
}

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "About this application",
	Long:  `About this application Display detailed information about the Task Manager CLI application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("Task Manager CLI v%s\n\nA simple and efficient task management tool built with Go and Cobra.\nPerfect for managing your daily tasks from the command line.\n\nAuthor: Your Name\nRepository: https://github.com/example/taskcli\nLicense: MIT", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(aboutCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "root cmd error %s\n", err.Error())
		os.Exit(1)
	}
}
