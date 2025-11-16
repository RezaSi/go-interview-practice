# Hints for Challenge 3: Advanced Configuration & Hooks

---

### Hint 1: Implementing the Levels Method
The **Levels** method tells logrus which severities your hook should be triggered for
Since the hook is only for errors, return a slice containing `ErrorLevel` and `FatalLevel`

```go
// In the ErrorHook's Levels method...
return []logrus.Level{
    logrus.ErrorLevel,
    logrus.FatalLevel,
}
```

---

### Hint 2: Formatting the Entry in the Fire Method

Inside the **Fire** method, you have access to the hook's configured `Formatter`
Use its `Format` method to turn the `logrus.Entry` into a byte slice

```go
// In the ErrorHook's Fire method...
line, err := h.Formatter.Format(entry)
if err != nil {
    // It's good practice to return the error if formatting fails
    return err
}
```

---

### Hint 3: Writing the Formatted Log

After formatting the entry, write the result to the hook’s `Out` writer
Tests expect each log on its own line, so **append a newline character**

```go
// In the ErrorHook's Fire method, after formatting...
_, err = h.Out.Write(append(line, '\n'))
return err
```

---

### Hint 4: Logging the Simulated Tasks

In the `runTaskScheduler` function, use chained calls with `WithField` for structured logs

The tests **check exact field names and message strings**
For example:

* Success: `Task '<NAME>' completed successfully`
* Warning: `Task '<NAME>' took longer than expected`
* Error: `Task '<NAME>' failed: <reason>`

```go
// For the failed task log...
logger.WithField("retry_attempt", 3).
    Errorf("Task '%s' failed: %s", "Backup database", "connection timed out")
```

---

### Hint 5: Assembling the Logger in main

The `main` function is where everything connects. Follow this sequence:

1. Create the logger instance
2. Configure the main logger’s **TextFormatter** and send output to **os.Stdout**
3. Create an `io.Writer` for the hook (e.g. a file or buffer)
4. Instantiate `ErrorHook` with its **JSONFormatter** and the writer
5. Add the hook to the logger with `logger.AddHook(hook)`

```go
// In the main function...
logger := logrus.New()
logger.SetFormatter(&logrus.TextFormatter{})
logger.SetOutput(os.Stdout)

hookWriter := &bytes.Buffer{}
hook := &ErrorHook{
    Out:       hookWriter,
    Formatter: &logrus.JSONFormatter{},
}
logger.AddHook(hook)
```

---

### Hint 6: Remember the Outputs

* **Main logger (console):** Human-readable `TextFormatter`
* **Hook (errors only):** Machine-readable `JSONFormatter`
* Info/Warning logs → appear **only in console**
* Error/Fatal logs → appear in **both console and hook output**

---