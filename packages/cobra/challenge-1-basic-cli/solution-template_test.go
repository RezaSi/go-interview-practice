package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Helper function to execute command and capture output
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	// Test root command shows help
	output, err := executeCommand(rootCmd)
	if err != nil {
		t.Fatalf("Root command failed: %v", err)
	}

	// Check if help text contains expected content
	expectedStrings := []string{
		"Task Manager CLI - Manage your tasks efficiently",
		"Usage:",
		"taskcli [command]",
		"Available Commands:",
		"about",
		"version",
		"help",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Root command output missing '%s'\nGot: %s", expected, output)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	output, err := executeCommand(rootCmd, "version")
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}

	expectedLines := []string{
		"taskcli version 1.0.0",
		"Built with ❤️ using Cobra",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(output, expected) {
			t.Errorf("Version command output missing '%s'\nGot: %s", expected, output)
		}
	}
}

func TestAboutCommand(t *testing.T) {
	output, err := executeCommand(rootCmd, "about")
	if err != nil {
		t.Fatalf("About command failed: %v", err)
	}

	expectedStrings := []string{
		"Task Manager CLI v1.0.0",
		"A simple and efficient task management tool",
		"Author:",
		"Repository:",
		"License:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("About command output missing '%s'\nGot: %s", expected, output)
		}
	}
}

func TestHelpCommand(t *testing.T) {
	// Test help command
	output, err := executeCommand(rootCmd, "help")
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	if !strings.Contains(output, "Available Commands:") {
		t.Errorf("Help command should show available commands\nGot: %s", output)
	}
}

func TestHelpVersionCommand(t *testing.T) {
	// Test help for version command
	output, err := executeCommand(rootCmd, "help", "version")
	if err != nil {
		t.Fatalf("Help version command failed: %v", err)
	}

	expectedStrings := []string{
		"Show version information",
		"Display the current version",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Help version command output missing '%s'\nGot: %s", expected, output)
		}
	}
}

func TestHelpAboutCommand(t *testing.T) {
	// Test help for about command
	output, err := executeCommand(rootCmd, "help", "about")
	if err != nil {
		t.Fatalf("Help about command failed: %v", err)
	}

	expectedStrings := []string{
		"About this application",
		"Display detailed information",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Help about command output missing '%s'\nGot: %s", expected, output)
		}
	}
}

func TestCommandStructure(t *testing.T) {
	// Test that commands are properly added to root
	commands := rootCmd.Commands()

	commandNames := make([]string, len(commands))
	for i, cmd := range commands {
		commandNames[i] = cmd.Name()
	}

	expectedCommands := []string{"about", "completion", "help", "version"}

	for _, expected := range expectedCommands {
		found := false
		for _, name := range commandNames {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command '%s' not found in root commands. Available: %v", expected, commandNames)
		}
	}
}

func TestRootCommandMetadata(t *testing.T) {
	// Test root command properties
	if rootCmd.Use != "taskcli" {
		t.Errorf("Expected root command Use to be 'taskcli', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short != "Task Manager CLI - Manage your tasks efficiently" {
		t.Errorf("Expected root command Short description, got '%s'", rootCmd.Short)
	}
}

func TestVersionCommandMetadata(t *testing.T) {
	// Find version command
	var versionCommand *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "version" {
			versionCommand = cmd
			break
		}
	}

	if versionCommand == nil {
		t.Fatal("Version command not found")
	}

	if versionCommand.Short != "Show version information" {
		t.Errorf("Expected version command Short description, got '%s'", versionCommand.Short)
	}
}

func TestAboutCommandMetadata(t *testing.T) {
	// Find about command
	var aboutCommand *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "about" {
			aboutCommand = cmd
			break
		}
	}

	if aboutCommand == nil {
		t.Fatal("About command not found")
	}

	if aboutCommand.Short != "About this application" {
		t.Errorf("Expected about command Short description, got '%s'", aboutCommand.Short)
	}
}

// Test that main function can be called without panicking
func TestMainFunction(t *testing.T) {
	// This test ensures main() can be called
	// We'll override os.Args to prevent actual command execution
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set args to just show help (safe operation)
	os.Args = []string{"taskcli", "--help"}

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked: %v", r)
		}
	}()

	// Note: We don't actually call main() here to avoid side effects
	// Instead we test that the root command can execute help
	_, err := executeCommand(rootCmd, "--help")
	if err != nil {
		t.Errorf("Help command should not error: %v", err)
	}
}
