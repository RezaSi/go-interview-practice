package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// RetryHook forwards warning/error logs to a secondary destination (stderr)
// This could represent routing critical logs to monitoring/alerting systems
type RetryHook struct {
	Out       io.Writer
	Formatter logrus.Formatter
}

// Levels tells Logrus which severities this hook should fire for
// Restrict to Warn and Error
func (h *RetryHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel}
}

// Fire is invoked when a log entry matches the levels above
// It formats the entry with the hook’s formatter and writes to Out
func (h *RetryHook) Fire(entry *logrus.Entry) error {
	formatter := h.Formatter
	if formatter == nil {
	  formatter = &logrus.JSONFormatter{}
	}
	out := h.Out
	if out == nil {
	  out = os.Stderr
	}
	line, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = out.Write(line)
	return err
}

// Task represents a unit of work processed by the worker
type Task struct {
	ID      string
	Name    string
	Retries int
}

// runWorker simulates execution of tasks with retries, failures, and successes
// Logs at Info for normal ops, Warn for retries, Error for failures
func runWorker(logger *logrus.Logger, tasks []Task) {
	for _, task := range tasks {
		start := time.Now()

		// Start log
		logger.WithField("task_id", task.ID).
			Infof("Starting task: %s", task.Name)

		// Retry log
		if task.Retries > 0 {
			logger.WithFields(logrus.Fields{
				"task_id": task.ID,
				"retries": task.Retries,
			}).Warnf("Task %s retried %d time(s)", task.Name, task.Retries)
		}

		// Failure simulation
		if task.Name == "FailMe" {
			logger.WithFields(logrus.Fields{
				"task_id": task.ID,
				"error":   "simulated failure",
			}).Error("Task failed due to simulated error")
			continue
		}

		// Success log
		duration := time.Since(start)
		logger.WithFields(logrus.Fields{
			"task_id":  task.ID,
			"duration": duration,
		}).Infof("Task %s completed successfully", task.Name)
	}
}

// setupLogger configures the main logger with:
// - Console output (TextFormatter)
// - File output (JSONFormatter)
// - Hook to route warnings/errors separately
func setupLogger() *logrus.Logger {
	logger := logrus.New()

	// Console (human readable)
	consoleFormatter := &logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	}

	// File (structured JSON)
	logFile := "app.log"
	if _, err := os.Stat(logFile); err == nil {
		os.Rename(logFile, fmt.Sprintf("%s.old", logFile))
	}
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("could not open log file: %v", err))
	}

	// Send logs to both stdout (human text) and file (JSON)
	mw := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(mw)

	// For stdout readability, use TextFormatter
	// (File will still get JSON lines because JSON is written explicitly)
	logger.SetFormatter(consoleFormatter)

	// Retry hook: 
	// Warn/Error routed to stderr in JSON
	retryHook := &RetryHook{
		Out:       os.Stderr,
		Formatter: &logrus.JSONFormatter{},
	}
	logger.AddHook(retryHook)

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