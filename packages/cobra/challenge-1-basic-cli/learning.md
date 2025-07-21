# Learning: Cobra CLI Framework Fundamentals

## 🌟 **What is Cobra?**

Cobra is a powerful library for creating modern command-line interfaces in Go. It's used by many popular CLI tools including Docker, Kubernetes, Hugo, and GitHub CLI.

### **Why Cobra?**
- **Powerful**: Easy to create complex CLI applications with subcommands
- **User-friendly**: Automatic help generation, shell completion, and man pages
- **Flexible**: Support for flags, arguments, and nested commands
- **Well-tested**: Battle-tested in production by major projects
- **POSIX-compliant**: Follows standard CLI conventions

## 🏗️ **Core Concepts**

### **1. Commands**
Commands are the core building blocks of a CLI application. Each command can have:
- A name (e.g., "version", "about")
- Short and long descriptions
- A function to execute
- Subcommands
- Flags and arguments

```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "A brief description",
    Long:  "A longer description explaining what this command does",
    Run: func(cmd *cobra.Command, args []string) {
        // Command implementation
    },
}
```

### **2. Root Command**
The root command is the main entry point of your CLI application:

```go
var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "My application does amazing things",
}
```

### **3. Subcommands**
You can add subcommands to create hierarchical CLI structures:

```go
rootCmd.AddCommand(versionCmd)
rootCmd.AddCommand(configCmd)
```

## 📖 **Command Structure**

### **Command Hierarchy**
```
myapp                    # Root command
├── version             # Subcommand
├── config              # Subcommand
│   ├── set            # Sub-subcommand
│   └── get            # Sub-subcommand
└── help               # Auto-generated
```

### **Command Properties**
- **Use**: The command name and syntax
- **Short**: Brief description (shown in command lists)
- **Long**: Detailed description (shown in help)
- **Example**: Usage examples
- **Run**: Function to execute when command is called

## 🔧 **Building Your First CLI**

### **Step 1: Create Root Command**
```go
var rootCmd = &cobra.Command{
    Use:   "taskcli",
    Short: "Task Manager CLI",
    Long:  "A powerful task management tool for the command line",
}
```

### **Step 2: Add Subcommands**
```go
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("taskcli version 1.0.0")
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
```

### **Step 3: Execute**
```go
func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

## 🎯 **Auto-Generated Features**

### **Help System**
Cobra automatically generates:
- `help` command
- `-h, --help` flags for all commands
- Formatted help text with descriptions
- Usage information

### **Completion**
Cobra provides shell completion for:
- Bash
- Zsh
- Fish
- PowerShell

### **Error Handling**
Cobra provides built-in error handling for:
- Unknown commands
- Invalid flags
- Missing required arguments

## 💡 **Best Practices**

### **1. Command Naming**
- Use clear, descriptive names
- Follow verb-noun pattern (e.g., `list tasks`, `create user`)
- Keep names short but meaningful

### **2. Descriptions**
- Write helpful short descriptions for command lists
- Provide detailed long descriptions with examples
- Include usage examples when helpful

### **3. Error Messages**
- Provide clear, actionable error messages
- Suggest correct usage when possible
- Use consistent error formatting

### **4. Output Formatting**
- Use consistent output formatting
- Consider structured output (JSON/YAML) for automation
- Provide human-readable output by default

## 🚀 **Advanced Features**

### **PreRun Hooks**
Execute code before command runs:
```go
PreRun: func(cmd *cobra.Command, args []string) {
    // Setup or validation code
},
```

### **Persistent Flags**
Flags available to command and all subcommands:
```go
rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
```

### **Required Commands**
Make subcommands required:
```go
cmd.MarkFlagRequired("name")
```

## 📚 **Real-World Examples**

### **Docker CLI Structure**
```
docker
├── build
├── run
├── ps
├── images
└── system
    ├── df
    └── prune
```

### **Kubernetes CLI Structure**
```
kubectl
├── get
├── create
├── apply
├── delete
└── config
    ├── view
    └── set-context
```

## 🔗 **Resources**

- [Official Cobra Documentation](https://cobra.dev/)
- [Cobra GitHub Repository](https://github.com/spf13/cobra)
- [Cobra Generator](https://github.com/spf13/cobra-cli)
- [CLI Design Guidelines](https://clig.dev/)

## 🎪 **Common Patterns**

### **Version Command**
Every CLI should have a version command:
```go
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("%s version %s\n", appName, version)
    },
}
```

### **Configuration Command**
Many CLIs need configuration management:
```go
var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage configuration",
}
```

### **List Commands**
Common pattern for listing resources:
```go
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List items",
    Run: func(cmd *cobra.Command, args []string) {
        // List implementation
    },
} 