package main

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// setupLogger configures the global logrus logger
func setupLogger(out io.Writer, level string) {
	logrus.SetOutput(out)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(lvl)
	}
}

// runLogbookOperations simulates the main logic of the logbook application
func runLogbookOperations() {
	logrus.Debug("Checking system status")
	logrus.Info("Logbook application starting up.")
	logrus.Warn("Disk space is running low.")
	logrus.Error("Failed to connect to remote backup service.")
	logrus.Fatal("Critical configuration file 'config.yml' not found.")
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