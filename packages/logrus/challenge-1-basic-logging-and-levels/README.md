# Challenge 1: Basic Logging & Levels

Build a simple **Logbook Application** that demonstrates fundamental logging concepts using the `logrus` library. This challenge will teach you how to set up a logger, use different log levels, and control log output format.

## Challenge Requirements

Create a simple Go application that simulates adding entries to a logbook. The application must:

1. **Initialize Logrus**: Set up a global `logrus` logger instance
2. **Use Different Log Levels**: Implement a function that logs messages at various severity levels: `Debug`, `Info`, `Warn`, `Error`, `Fatal`, and `Panic`
3. **Control Log Output Level**: The application should be configurable to show logs above a certain severity. For example, setting the level to `info` should hide `debug` messages
4. **Format Logs**: Configure the logger to output logs in a structured `JSON` format

## How It Should Work

You will build a simple `runLogbookOperations` function that simulates a logbook's daily operations. The application's logging output will be controlled by setting the log level.

### Sample Output

**When the log level is set to `info`:**  
The output should only include `info`, `warn`, `error`, and `fatal` messages. The `fatal` log will terminate the application, so `panic` will not be reached.

```json
{"level":"info","msg":"Logbook application starting up.","time":"2025-10-02T17:30:00+05:30"}
{"level":"info","msg":"Opening today's log entry.","time":"2025-10-02T17:30:00+05:30"}
{"level":"warning","msg":"Disk space is running low.","time":"2025-10-02T17:30:00+05:30"}
{"level":"error","msg":"Failed to connect to remote backup service.","time":"2025-10-02T17:30:00+05:30"}
{"level":"fatal","msg":"Critical configuration file 'config.yml' not found.","time":"2025-10-02T17:30:00+05:30"}
```
### When the log level is set to `debug`

The output should include all messages from `debug` up to `fatal`.

```json
{"level":"debug","msg":"Checking system status...","time":"2025-10-02T17:30:00+05:30"}
{"level":"debug","msg":"Memory usage: 256MB","time":"2025-10-02T17:30:00+05:30"}
{"level":"info","msg":"Logbook application starting up.","time":"2025-10-02T17:30:00+05:30"}
{"level":"info","msg":"Opening today's log entry.","time":"2025-10-02T17:30:00+05:30"}
{"level":"warning","msg":"Disk space is running low.","time":"2025-10-02T17:30:00+05:30"}
{"level":"error","msg":"Failed to connect to remote backup service.","time":"2025-10-02T17:30:00+05:30"}
{"level":"fatal","msg":"Critical configuration file 'config.yml' not found.","time":"2025-10-02T17:30:00+05:30"}
```

## Implementation Requirements

### Logger Configuration (`setupLogger` function)

- Set the log formatter to `logrus.JSONFormatter`
- Set the output to `os.Stdout`
- Set the log level based on a provided string (e.g., `"info"`, `"debug"`). If the string is invalid, default to `logrus.InfoLevel`

### Main Logic (`runLogbookOperations` function)

This function should contain at least one log statement for each of the six levels:

- **Debug**: Log verbose details useful for development (e.g., `"Checking system status..."`)
- **Info**: Log informational messages about application progress (e.g., `"Logbook application starting up."`)
- **Warn**: Log potential issues that don't prevent the application from running (e.g., `"Disk space is running low."`)
- **Error**: Log errors the application might recover from (e.g., `"Failed to connect to remote backup service."`)
- **Fatal**: Log a critical error that must terminate the application (e.g., `"Critical configuration file not found."`)
- `Fatal` calls `os.Exit(1)` after logging
- **Panic**: Log a message and then panic. This is for unrecoverable application states

---

## Testing Requirements

Your solution must pass tests that verify:

- The logger is correctly configured (formatter, output)
- Setting a specific log level (e.g., `Warn`) correctly filters out lower-level messages (`Info`, `Debug`)
- Messages are logged in the expected JSON format
- A **Fatal** level log correctly triggers an exit
- A **Panic** level log correctly causes a panic
---