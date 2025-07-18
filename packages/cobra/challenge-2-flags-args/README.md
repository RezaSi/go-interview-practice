# Challenge 2: Flags and Arguments

Build a **File Manager CLI** that demonstrates advanced flag and argument handling with Cobra.

## Challenge Requirements

Create a CLI application that supports:

1. **Global Flags** - Available to all commands
2. **Command-specific Flags** - Available only to specific commands  
3. **Required and Optional Flags** - Different validation levels
4. **Arguments** - Positional parameters
5. **Flag Types** - String, int, bool flags

## Expected CLI Structure

```
filecli --verbose list                    # Global flag
filecli list --format json               # Command flag
filecli copy file1.txt file2.txt         # Arguments
filecli delete --force myfile.txt        # Required flag + argument
filecli create --name "test" --size 100  # Multiple flags
```

## Testing Requirements

Your solution must pass tests for:
- Global and command-specific flags
- Required flag validation
- Argument handling and validation
- Flag type conversion (string, int, bool)
- Help text generation for flags 