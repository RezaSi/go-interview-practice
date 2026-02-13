# Hints for Challenge 1: Basic Logging & Levels

## Hint 1: Configuring the Logger Output

In the `setupLogger` function, the first step is to tell Logrus where to send the logs. The function receives an `io.Writer` named `out`.

```go
// Tell the global logrus logger to write to the `out` variable.
logrus.SetOutput(out)
```

---

## Hint 2: Setting the JSON Formatter

To make logs structured, you need to set the formatter. Create a new instance of `logrus.JSONFormatter` and pass its address to `logrus.SetFormatter`.

```go
// This tells logrus to format all subsequent logs as JSON.
logrus.SetFormatter(&logrus.JSONFormatter{})
```

---

## Hint 3: Parsing and Setting the Log Level

You need to convert the `level` string (e.g., `"debug"`) into a `logrus.Level` type. The `logrus.ParseLevel` function does this for you. It returns the level and an error if the string is invalid.

```go
// Try to parse the level string.
lvl, err := logrus.ParseLevel(level)

// If the string is not a valid level, `err` will not be nil.
if err != nil {
    // In case of an error, we fall back to a sensible default.
    logrus.SetLevel(logrus.InfoLevel)
} else {
    // If parsing was successful, use the parsed level.
    logrus.SetLevel(lvl)
}
```

---

## Hint 4: Logging the First Message

In the `runLogbookOperations` function, you can use the package-level functions like `logrus.Debug()`, `logrus.Info()`, etc., to log messages.

```go
// For the first TODO, use the Debug function.
logrus.Debug("Checking system status...")
```

---

## Hint 5: Logging the Remaining Messages

Follow the same pattern for the other log levels in the specified order.

```go
func runLogbookOperations() {
    logrus.Debug("Checking system status...")
    logrus.Info("Logbook application starting up")
    logrus.Warn("Disk space is running low")
    logrus.Error("Failed to connect to remote backup service")
    logrus.Fatal("Critical configuration file 'config.yml' not found")
    logrus.Panic("Unhandled database connection issue")
}
```

---

## Hint 6: Understanding Fatal and Panic

Remember the special behavior of `Fatal` and `Panic`:

* `logrus.Fatal(...)` will log the message **and then immediately terminate** the program (by calling `os.Exit(1)`). No code after it will run.
* `logrus.Panic(...)` will log the message **and then cause a panic**.

This means in a normal run, you will never see the `Panic` log if the `Fatal` log comes before it. The tests are designed to handle these specific behaviors.

---

## Hint 7: Running and Testing Your Code

Once you've filled in the TODOs, you can run your `main.go` file from the terminal and pass a log level as an argument to see the effect.

```bash
# Run with the default "info" level (no debug messages)
go run .

# Run and show only messages from "warn" level and above
go run . warn

# Run with "debug" level to see all messages
go run . debug
```
---