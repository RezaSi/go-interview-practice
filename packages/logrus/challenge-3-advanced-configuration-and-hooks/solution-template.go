package main

import (
	"bytes"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// ErrorHook is a custom hook that sends logs of specified levels to a dedicated writer
// It must implement the logrus.Hook interface
type ErrorHook struct {
	Out       io.Writer
	Formatter logrus.Formatter
}

// Levels returns the log levels that this hook will be triggered for
func (h *ErrorHook) Levels() []logrus.Level {
	// TODO: Return a slice of logrus.Level containing only ErrorLevel and FatalLevel
	// Hint: return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel}
	return nil
}

// Fire is called when a log entry is fired for one of the specified Levels
func (h *ErrorHook) Fire(entry *logrus.Entry) error {
	// TODO: Use the hook's formatter to format the entry into bytes
	// Hint: line, err := h.Formatter.Format(entry)
	
	// TODO: Write the formatted line to the hook's output writer (h.Out)
	// Remember to add a newline character at the end to separate log entries
	// Hint: _, err = h.Out.Write(append(line, '\n'))
	
	return nil
}

// runTaskScheduler simulates a scheduler running various tasks and logging their outcomes
func runTaskScheduler(logger *logrus.Logger) {
	logger.Info("Starting task scheduler...")

	// TODO: Log a successful task using the Info level
	// Message: "Task 'Process daily reports' completed successfully"
	// Field: "task_duration" with a value like "250ms"
	// Hint: logger.WithField(...).Infof(...)
	
	// TODO: Log a task with a warning
	// Message: "Task 'Sync user data': upstream API is slow"
	// Field: "task_id" with a value like "sync-001"
	// Hint: logger.WithField(...).Warnf(...)
	
	// TODO: Log a failed task using the Error level
	// Message: "Task 'Backup database' failed: connection timed out"
	// Field: "retry_attempt" with an integer value like 3
	// Hint: logger.WithField(...).Errorf(...)
}

func main() {
	// TODO: Create a new instance of the logger
	// Hint: logger := logrus.New()
	logger := logrus.New() // Basic initialization - configure below

	// TODO: Set the logger's formatter to TextFormatter
	logger.SetFormatter(&logrus.TextFormatter{})

	// TODO: Set the logger's output to standard out
	logger.SetOutput(os.Stdout)

	// The hook will write to a separate buffer. In a real app, this could be a file
	var hookWriter io.Writer = &bytes.Buffer{}

	// TODO: Create an instance of your ErrorHook
	// It needs an output writer (hookWriter) and a JSON formatter
	// Hint: hook := &ErrorHook{ Out: hookWriter, Formatter: &logrus.JSONFormatter{} }
	hook := &ErrorHook{Out: hookWriter, Formatter: &logrus.JSONFormatter{}}

	// TODO: Add the hook to the logger
	logger.AddHook(hook)

	// Run the scheduler simulation
	runTaskScheduler(logger)

	// In a real application, you might want to inspect the hook's output.
	// For this challenge, the tests will handle that.
	// fmt.Println("\n--- Hook Output ---")
	// fmt.Println(hookWriter.(*bytes.Buffer).String())
}