# Hints for Challenge 4: Production Logging Patterns

## Hint 1: Defining the `RetryHook` Struct and Levels

First, define the `RetryHook` struct with exported fields for its output destination and formatter. Then, implement the `Levels` method to tell Logrus that this hook should only activate for warnings and errors

```go
import (
	"io"
	"github.com/sirupsen/logrus"
)

type RetryHook struct {
	Out       io.Writer
	Formatter logrus.Formatter
}

func (h *RetryHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel}
}
```

---

## Hint 2: Implementing the `Fire` Method

The `Fire` method formats the log entry using the hook's formatter and writes it to its dedicated output writer

```go
func (h *RetryHook) Fire(entry *logrus.Entry) error {
	line, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.Out.Write(append(line, '\n'))
	return err
}
```

---

## Hint 3: Setting Up Multiple Destinations

Use `io.MultiWriter` to log to both console and file simultaneously

```go
logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
	panic("Failed to open log file")
}

multiWriter := io.MultiWriter(os.Stdout, logFile)
logger.SetOutput(multiWriter)
```

---

## Hint 4: Logging the Task Lifecycle

Define a simple `Task` struct and use `WithFields` for structured logging

```go
type Task struct {
	ID      int
	Name    string
	Retries int
}

for _, task := range tasks {
	logger.WithFields(logrus.Fields{"task_id": task.ID}).Infof("Starting task: %s", task.Name)

	for i := 1; i <= task.Retries; i++ {
		logger.WithFields(logrus.Fields{"task_id": task.ID, "retries": i}).Warnf(
			"Task '%s' retried (attempt %d of %d)", task.Name, i, task.Retries)
	}

	if task.Name == "FailMe" {
		logger.WithFields(logrus.Fields{"task_id": task.ID, "error": "simulated failure"}).Error(
			"Task '" + task.Name + "' failed after all retries")
	} else {
		logger.WithFields(logrus.Fields{"task_id": task.ID, "duration": "15ms"}).Info(
			"Task '" + task.Name + "' completed successfully")
	}
}
```

---

## Hint 5: Finalizing the Logger Configuration

Set main formatter and add `RetryHook`

```go
logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

hook := &RetryHook{
	Out:       os.Stderr,
	Formatter: &logrus.JSONFormatter{},
}
logger.AddHook(hook)
```

---

## Hint 6: Summary & Tips

* Always define `Task` struct with ID, Name, and Retries
* Log Info-level messages for task start and success
* Log Warn-level for retries and Error-level for failures
* Use `io.MultiWriter` for console + file output
* Use `RetryHook` to filter and send Warn/Error logs to stderr (or alerting system)

---