# Challenge 1: Basic CLI Application

Build a **Task Manager CLI** using Cobra that demonstrates fundamental command-line application concepts.

## Challenge Requirements

Create a CLI application called `taskcli` that supports:

1. **Root Command** - Display welcome message and help
2. **Version Command** - Show application version
3. **About Command** - Display application information  
4. **Basic Help** - Auto-generated help text for all commands

## Expected CLI Structure

```
taskcli                    # Root command - shows help
taskcli version           # Shows version information
taskcli about             # Shows about information
taskcli help              # Shows help (auto-generated)
taskcli help version      # Shows help for version command
```

## Sample Output

**Root Command (`taskcli`):**
```
Task Manager CLI - Manage your tasks efficiently

Usage:
  taskcli [command]

Available Commands:
  about       About this application
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Show version information

Flags:
  -h, --help   help for taskcli

Use "taskcli [command] --help" for more information about a command.
```

**Version Command (`taskcli version`):**
```
taskcli version 1.0.0
Built with ❤️ using Cobra
```

**About Command (`taskcli about`):**
```
Task Manager CLI v1.0.0

A simple and efficient task management tool built with Go and Cobra.
Perfect for managing your daily tasks from the command line.

Author: Your Name
Repository: https://github.com/example/taskcli
License: MIT
```

## Implementation Requirements

### Root Command
- Use "taskcli" as the command name
- Include a description: "Task Manager CLI - Manage your tasks efficiently"  
- Show help by default when no subcommand is provided

### Version Command
- Command name: "version"
- Short description: "Show version information"
- Output format: "taskcli version 1.0.0\nBuilt with ❤️ using Cobra"

### About Command  
- Command name: "about"
- Short description: "About this application"
- Show detailed application information including version, description, author, etc.

## Testing Requirements

Your solution must pass tests for:
- Root command displays help when run without arguments
- Version command outputs correct format
- About command shows application information
- Help command works for all commands
- Command structure matches expected hierarchy
- All commands have proper descriptions 