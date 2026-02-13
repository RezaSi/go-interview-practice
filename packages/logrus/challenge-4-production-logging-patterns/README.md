# Challenge 4: Production Logging Patterns

Build a robust **Task Worker** in Go that demonstrates advanced, production-ready logging patterns with multi-destination outputs and a custom alerting hook

---

## Challenge Requirements

Implement a worker system with structured logging configured for a real-world production setup

### Multi-Destination Logging
- **Console (`stdout`)**  
  - Human-readable, real-time monitoring
  - Formatted as plaintext

- **Log File (`app.log`)**  
  - Persistent, structured storage 
  - Formatted as JSON for later analysis 

- **Alerting Stream (`stderr`)**  
  - Critical logs routed to a custom hook 
  - Formatted as JSON for integration with alerting pipelines 

### Custom RetryHook
- Triggered only for `WarnLevel` and `ErrorLevel` logs 
- Formats entries as JSON 
- Writes to a separate stream (`stderr`) 
- Simulates sending events to alerting systems (e.g., PagerDuty, Sentry) 

### Task Worker Simulation
Implement a `runWorker` function that processes a list of tasks and logs their lifecycle with structured fields:

- **Start**: `Info` log → `"Starting task: <Task Name>"`  
- **Retry**: `Warn` log → `"Task '<Task Name>' retried (attempt X of Y)"`  
- **Failure**: `Error` log → `"Task '<Task Name>' failed after all retries"`  
- **Success**: `Info` log → `"Task '<Task Name>' completed successfully"`  

---

## Data Structures

```go
// Task represents a single unit of work in the worker queue
type Task struct {
    ID       int
    Name     string
    Retries  int
    Duration time.Duration
}

// RetryHook captures warnings and errors for alerting streams
type RetryHook struct {
    Out       io.Writer
    Formatter logrus.Formatter
}
````

---

## Implementation Requirements

### Logger Configuration (`setupLogger`)

* Use `logrus.New()` to create a logger instance
* Use `io.MultiWriter` to send logs to both console (`stdout`) and file (`app.log`)
* Set console formatter as `logrus.TextFormatter`
* Register a `RetryHook` with `JSONFormatter` that writes to `stderr`

### RetryHook

* `Levels()`: Return `WarnLevel` and `ErrorLevel`
* `Fire()`: Format the log entry as JSON and write it to `Out`

### Worker Logic (`runWorker`)

* Iterate through tasks.
* Log start, retries, success, and failure with structured fields
* Retries and errors should appear in **all three streams**
* Info logs should appear in console and file, but **not** in alert stream

---

## Testing Requirements

Your solution must pass tests that verify:

* Worker logs the correct lifecycle events with structured fields
* `RetryHook` is correctly registered on the logger
* Hook outputs JSON-formatted `Warn` and `Error` logs to `stderr`
* All `Warn`/`Error` logs appear in **console + file + hook**
* All `Info` logs appear in **console + file only**

---