package main

import (
	"io"
    "os"

    "github.com/sirupsen/logrus"
)

// setupLogger configures the global logrus logger
func setupLogger(out io.Writer, level string) {
    // TODO: Set the logger's output.
    logrus.SetOutput(out)

    // TODO: Set the logger's formatter to JSON.
    logrus.SetFormatter(&logrus.JSONFormatter{})

    // TODO: Parse and set the log level, defaulting to InfoLevel on error.
    lvl, err := logrus.ParseLevel(level)
    if err != nil {
        logrus.SetLevel(logrus.InfoLevel)
    } else {
        logrus.SetLevel(lvl)
    }
}

// runLogbookOperations simulates the main logic of the logbook application
func runLogbookOperations() {
	// Important Note: The content inside "text of log outputs" should be exact as to not fail the test

    // TODO: Add a Debug log: "Checking system status."
    logrus.Debug("Checking system status.")

    // TODO: Add an Info log: "Logbook application starting up."
    logrus.Info("Logbook application starting up.")

    // TODO: Add a Warn log: "Disk space is running low."
    logrus.Warn("Disk space is running low.")

    // TODO: Add an Error log: "Failed to connect to remote backup service."
    logrus.Error("Failed to connect to remote backup service.")

    // TODO: Add a Fatal log: "Critical configuration file 'config.yml' not found."
    logrus.Fatal("Critical configuration file 'config.yml' not found.")

    // TODO: Add a Panic log: "Unhandled database connection issue."
    logrus.Panic("Unhandled database connection issue.")
}

func main() {
    logLevel := "info"
    if len(os.Args) > 1 {
        logLevel = os.Args[1]
    }
    setupLogger(os.Stdout, logLevel)
    logrus.Infof("Log level set to '%s'", logrus.GetLevel().String())
    runLogbookOperations()
}