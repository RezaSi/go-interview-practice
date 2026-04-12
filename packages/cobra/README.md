# Cobra CLI Development Challenges

Master the art of building powerful command-line applications in Go using the Cobra library. This package contains 4 progressive challenges that take you from basic CLI concepts to advanced production-ready patterns.

## Challenge Overview

### 🎯 [Challenge 1: Basic CLI Application](./challenge-1-basic-cli/)
**Difficulty:** Beginner | **Duration:** 30-45 minutes

Learn the fundamentals of Cobra by building a simple task manager CLI with basic commands, version information, and help systems.

**Key Skills:**
- Basic Cobra command structure
- Root command setup
- Version and about commands
- Auto-generated help text
- Command hierarchy basics

**Topics Covered:**
- `cobra.Command` basics
- Command descriptions and usage
- Help system fundamentals
- CLI application structure

---

### 🚀 [Challenge 2: Flags and Arguments](./challenge-2-flags-args/)
**Difficulty:** Intermediate | **Duration:** 45-60 minutes

Build a more sophisticated CLI with comprehensive flag handling, argument validation, and interactive features.

**Key Skills:**
- Flag types and validation
- Positional arguments
- Required vs optional flags
- Flag inheritance
- Custom validation functions

**Topics Covered:**
- Persistent and local flags
- Flag binding and validation
- Argument handling patterns
- Error handling and user feedback
- Interactive CLI features

---

### 📦 [Challenge 3: Subcommands & Data Persistence](./challenge-3-subcommands-persistence/)
**Difficulty:** Intermediate | **Duration:** 45-60 minutes

Create an inventory management CLI that demonstrates advanced subcommand organization and JSON data persistence.

**Key Skills:**
- Nested command hierarchies
- CRUD operations via CLI
- JSON data persistence
- Search and filtering
- File I/O operations

**Topics Covered:**
- Complex command structures
- Data persistence patterns
- JSON marshaling/unmarshaling
- Search functionality
- Error handling with file operations

---

### ⚡ [Challenge 4: Advanced Features & Middleware](./challenge-4-advanced-features/)
**Difficulty:** Advanced | **Duration:** 60-90 minutes

Build a configuration management CLI showcasing advanced Cobra patterns including middleware, plugins, and multi-format configuration support.

**Key Skills:**
- Middleware system implementation
- Plugin architecture
- Configuration management (JSON/YAML/TOML)
- Environment variable integration
- Custom help templates

**Topics Covered:**
- Command middleware patterns
- Plugin system design
- Multiple configuration formats
- Validation pipelines
- Advanced CLI UX patterns

## Learning Path

```
Challenge 1: Basic CLI
        ↓
Challenge 2: Flags & Args  
        ↓
Challenge 3: Data & Subcommands
        ↓  
Challenge 4: Advanced Features
```

### Recommended Prerequisites
- **Challenge 1:** Basic Go knowledge
- **Challenge 2:** Completion of Challenge 1, understanding of data types
- **Challenge 3:** Completion of Challenges 1-2, JSON/file handling experience
- **Challenge 4:** Completion of Challenges 1-3, advanced Go patterns knowledge

## Key Cobra Concepts Covered

### Core Concepts
- **Commands:** Building command hierarchies and structures
- **Flags:** Persistent, local, required, and custom validation
- **Arguments:** Positional arguments and validation
- **Help System:** Custom help templates and documentation

### Advanced Patterns
- **Middleware:** Pre/post command execution hooks
- **Plugins:** Dynamic command registration and plugin architecture
- **Configuration:** Multi-format config management with environment integration
- **Validation:** Input validation and custom validators

### Production Features
- **Error Handling:** Graceful error management and user feedback
- **Performance:** Optimization patterns for CLI applications
- **Security:** Input sanitization and secure practices
- **UX Design:** Creating intuitive and helpful CLI interfaces

## Challenge Structure

Each challenge follows a consistent structure:

```
challenge-X-name/
├── README.md              # Challenge description and requirements
├── solution-template.go   # Template with TODOs to implement
├── solution-template_test.go  # Comprehensive test suite
├── run_tests.sh          # Test runner script
├── go.mod                # Go module with dependencies
├── metadata.json         # Challenge metadata
├── SCOREBOARD.md         # Participant scores
├── hints.md              # Implementation hints (when available)
├── learning.md           # Additional learning resources (when available)
└── submissions/          # Participant submission directory
```

## Getting Started

1. **Choose your starting challenge** based on your experience level
2. **Read the README.md** in the challenge directory
3. **Implement the solution** in `solution-template.go`
4. **Test your solution** using `./run_tests.sh`
5. **Submit via PR** to the submissions directory

## Testing Your Solutions

Each challenge includes a comprehensive test suite. To test your solution:

```bash
cd packages/cobra/challenge-X-name/
./run_tests.sh
```

The test script will:
- Prompt for your GitHub username
- Copy your solution to a temporary environment
- Run all tests against your implementation
- Provide detailed feedback on test results

## Common Patterns and Best Practices

### Command Structure
```go
var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "Brief description",
    Long:  "Detailed description",
}

var subCmd = &cobra.Command{
    Use:   "subcmd",
    Short: "Subcommand description",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}
```

### Flag Handling
```go
// Persistent flags (available to all subcommands)
rootCmd.PersistentFlags().StringVar(&config, "config", "", "config file")

// Local flags (only for this command)
cmd.Flags().StringVarP(&name, "name", "n", "", "name flag")

// Required flags
cmd.MarkFlagRequired("name")
```

### Error Handling
```go
func runCommand(cmd *cobra.Command, args []string) error {
    if err := validateInput(args); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // Command logic
    return nil
}
```

## Resources

- [Cobra Documentation](https://cobra.dev/)
- [Cobra GitHub Repository](https://github.com/spf13/cobra)
- [CLI Design Guidelines](https://clig.dev/)
- [Go CLI Best Practices](https://blog.gopheracademy.com/advent-2017/cli-application/)

## Contributing

Found an issue or want to improve a challenge? Contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

---

**Happy CLI Building!** 🚀

Master these challenges to become proficient in building production-ready command-line applications with Go and Cobra. 