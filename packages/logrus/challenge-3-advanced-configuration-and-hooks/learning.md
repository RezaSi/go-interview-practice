# Learning: Advanced Configuration & Hooks

In the previous challenges, our logs had a single destination: the console. While great for development, real applications require a more sophisticated approach. Logs are not just messages to be seen; they are critical data streams that need to be persisted, routed, and analyzed.

This is where **logrus hooks** come in. They are the key to transforming a simple logger into a powerful data pipeline.

---

## 🪝 The `logrus.Hook` Interface: A Logger's Secret Weapon

A hook is a mechanism that allows you to intercept log entries and take custom actions before the log is written by the main logger. It lets you "hook into" the logging process.

Think of it like adding a **BCC to an email**.  
The primary recipient (the main log output) gets the message as intended. At the same time, the hook silently sends a copy of that message to another destination, potentially in a completely different format.

The `logrus.Hook` interface is simple but powerful, defined by two methods:

### `Levels() []logrus.Level`
- **What it is:** A filter. This method returns a slice of the log levels the hook cares about.  
- **Why it's important:** It's highly efficient. Logrus won't waste time calling your hook for Info or Debug messages if your hook only registered for `ErrorLevel` and `FatalLevel`.

### `Fire(*logrus.Entry) error`
- **What it is:** An action. This method is the core of the hook. It's executed only when a log entry matches one of the levels specified in `Levels()`.  
- **Why it's important:** The hook receives the complete `logrus.Entry`, which contains the message, timestamp, level, and all the structured data fields. The hook has full control to format this entry however it likes and send it anywhere—a file, a network socket, or an external monitoring service.

---

## 🔧 Implementing a Custom Hook: Step-by-Step

Let’s break down how to build the `ErrorHook` for our Task Scheduler challenge.

### Step 1: Define the Hook Struct
A well-designed hook should be configurable. Instead of hardcoding the output destination or format, we define them as fields.

```go
type ErrorHook struct {
    Out       io.Writer
    Formatter logrus.Formatter
}
````

By accepting an `io.Writer`, we can send error logs to a file, a network connection, or an in-memory buffer during tests.

---

### Step 2: Implement `Levels()`

We want our hook to trigger only for critical errors, so we return a slice containing just those levels.

```go
func (h *ErrorHook) Levels() []logrus.Level {
    return []logrus.Level{
        logrus.ErrorLevel,
        logrus.FatalLevel,
    }
}
```

---

### Step 3: Implement `Fire()`

This is where the magic happens. The main logger might be using a human-readable `TextFormatter`, but our hook outputs machine-readable JSON.

```go
func (h *ErrorHook) Fire(entry *logrus.Entry) error {
    // Use the hook's dedicated formatter to serialize the entry.
    lineBytes, err := h.Formatter.Format(entry)
    if err != nil {
        return err
    }
    
    // Write the formatted JSON to the hook's output.
    _, err = h.Out.Write(lineBytes)
    return err
}
```

---

## Managing Log Files in Go

When logging to files, you can't just let them grow forever. The practice of managing log file size is called **log rotation**.

### Opening a File for Appending

To write logs to a file, you open it with specific flags.

```go
file, err := os.OpenFile("scheduler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    // Handle error
}
defer file.Close()

// You can now use this 'file' as an io.Writer for logrus or a hook.
```

---

### Simple Startup Rotation

A basic rotation strategy is to archive the old log file every time the application starts.

```go
if _, err := os.Stat("scheduler.log"); err == nil {
    // The file exists, so rename it.
    os.Rename("scheduler.log", "scheduler.log.old")
}
```

---

## A Note on Performance: Async Logging

The hooks we've designed are synchronous. When `logger.Error()` is called, your application code pauses and waits for both the main logger and the error hook to finish writing. For most apps, this is fine.

However, in **high-throughput systems**, I/O waits can become a bottleneck. The solution: **asynchronous logging**.

* When `logger.Error()` is called, the hook quickly places the log entry into a buffered Go channel (fast in-memory op).
* A background goroutine continuously reads from the channel and performs the slower I/O (writing to file/network).

This decouples your application’s performance from the speed of your logging backend.

---

## Resources

* [Logrus Documentation on Hooks](https://github.com/sirupsen/logrus#hooks)
* [Go Docs for os.OpenFile](https://pkg.go.dev/os#OpenFile)

---