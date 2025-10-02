package main

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// setupLogger configures the global logrus logger
func setupLogger(out io.Writer, level string) {
	// TODO: Set the logger's output to the `out` variable
	// Hint: logrus.SetOutput(out)

	// TODO: Set the logger's formatter to a new instance of logrus.JSONFormatter.
	// Hint: logrus.SetFormatter(&logrus.JSONFormatter{})

	// TODO: Parse the `level` string into a logrus.Level.
	// Use `logrus.ParseLevel()`
	// If there is an error during parsing, default the level to `logrus.InfoLevel`
	// Otherwise, use the parsed level
	// Hint:
	//   lvl, err := logrus.ParseLevel(level)
	//   if err != nil { ... }
	//   logrus.SetLevel(lvl)
}

// runLogbookOperations simulates the main logic of the logbook application
func runLogbookOperations() {
	// TODO: Add a Debug log: "Checking system status..."
	// Hint: logrus.Debug(...)

	// TODO: Add an Info log: "Logbook application starting up"
	// Hint: logrus.Info(...)

	// TODO: Add a Warn log: "Disk space is running low"
	// Hint: logrus.Warn(...)

	// TODO: Add an Error log: "Failed to connect to remote backup service"
	// Hint: logrus.Error(...)

	// TODO: Add a Fatal log: "Critical configuration file 'config.yml' not found"
	// This will terminate the application
	// Hint: logrus.Fatal(...)

	// TODO: Add a Panic log: "Unhandled database connection issue"
	// This will cause a panic
	// Hint: logrus.Panic(...)
}

func main() {
	// Default log level is "info"
	// If a command-line argument is provided, it's used as the log level
	logLevel := "info"
	if len(os.Args) > 1 {
		logLevel = os.Args[1]
	}

	setupLogger(os.Stdout, logLevel)

	// Add an informational log to show which level is currently active
	logrus.Infof("Log level set to '%s'", logrus.GetLevel().String())

	runLogbookOperations()
}