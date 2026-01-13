package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// RetryHook forwards warning/error logs to a secondary destination (stderr)
// TODO: Implement the hook struct with output target and formatter
type RetryHook struct {
	// TODO: define fields (e.g., Out, Formatter)
}

// Levels tells Logrus which severities this hook should fire for
// TODO: Restrict to Warn and Error
func (h *RetryHook) Levels() []logrus.Level {
	// TODO: return the correct levels
	return nil
}

// Fire is invoked when a log entry matches the levels above
// TODO: Format the entry and write it to the secondary output
func (h *RetryHook) Fire(entry *logrus.Entry) error {
	// TODO: implement hook writing logic
	return nil
}

// Task represents a unit of work processed by the worker
type Task struct {
	ID      string
	Name    string
	Retries int
}

// runWorker simulates execution of tasks with retries, failures, and successes
// TODO: Log at Info for normal ops, Warn for retries, Error for failures
// Add fields like task_id, retries, error, duration
func runWorker(logger *logrus.Logger, tasks []Task) {
	for _, task := range tasks {
		start := time.Now()

		// TODO: Log starting task (Info)

		// TODO: Log retries (Warn)

		// TODO: Log simulated failure (Error)

		// TODO: Log simulated success with duration (Info)

		_ = start // remove once implemented
	}
}

// setupLogger configures the main logger with:
// - Console output (TextFormatter)
// - File output (JSONFormatter)
// - Hook to route warnings/errors
// TODO: Implement logger setup.
func setupLogger() *logrus.Logger {
	logger := logrus.New()

	// TODO: Configure console output

	// TODO: Configure file output with rotation

	// TODO: Configure retry hook (stderr + JSON formatter)

	return logger
}

func main() {
	logger := setupLogger()

	tasks := []Task{
		{ID: "1", Name: "SendEmail", Retries: 0},
		{ID: "2", Name: "FailMe", Retries: 2},
		{ID: "3", Name: "GenerateReport", Retries: 1},
	}

	runWorker(logger, tasks)
}