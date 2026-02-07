package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// logEntry is used to unmarshal the JSON log output for verification
type logEntry struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

// captureOutput runs a function while capturing logrus's output
func captureOutput(action func()) string {
	var buffer bytes.Buffer
	originalOutput := logrus.StandardLogger().Out
	logrus.SetOutput(&buffer)

	defer logrus.SetOutput(originalOutput)
	action()
	return buffer.String()
}

func TestSetupLogger(t *testing.T) {
	t.Run("sets JSON formatter", func(t *testing.T) {
		setupLogger(io.Discard, "info")
		assert.IsType(t, &logrus.JSONFormatter{}, logrus.StandardLogger().Formatter, "Formatter should be JSONFormatter")
	})

	t.Run("sets correct output", func(t *testing.T) {
		var buffer bytes.Buffer
		setupLogger(&buffer, "info")
		assert.Equal(t, &buffer, logrus.StandardLogger().Out, "Output writer should be set correctly")
	})

	t.Run("sets valid log level", func(t *testing.T) {
		setupLogger(io.Discard, "debug")
		assert.Equal(t, logrus.DebugLevel, logrus.GetLevel(), "Log level should be set to Debug")
	})

	t.Run("defaults to info for invalid level", func(t *testing.T) {
		setupLogger(io.Discard, "invalid-level")
		assert.Equal(t, logrus.InfoLevel, logrus.GetLevel(), "Log level should default to Info for invalid input")
	})
}

func TestLogLevelFiltering(t *testing.T) {
	testCases := []struct {
		name           string
		levelToSet     string
		expectedLogs   []string
		unexpectedLogs []string
	}{
		{
			name:           "Debug Level",
			levelToSet:     "debug",
			expectedLogs:   []string{"Checking system status", "Logbook application starting up", "Disk space is running low"},
			unexpectedLogs: []string{},
		},
		{
			name:           "Info Level",
			levelToSet:     "info",
			expectedLogs:   []string{"Logbook application starting up", "Disk space is running low"},
			unexpectedLogs: []string{"Checking system status"},
		},
		{
			name:           "Warn Level",
			levelToSet:     "warn",
			expectedLogs:   []string{"Disk space is running low", "Failed to connect to remote backup service"},
			unexpectedLogs: []string{"Logbook application starting up", "Checking system status"},
		},
		{
			name:           "Error Level",
			levelToSet:     "error",
			expectedLogs:   []string{"Failed to connect to remote backup service"},
			unexpectedLogs: []string{"Disk space is running low", "Logbook application starting up"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the exit function to prevent test termination and signal with timeout
			done := make(chan struct{}, 1)
			originalExitFunc := logrus.StandardLogger().ExitFunc
			logrus.StandardLogger().ExitFunc = func(int) {
				select {
				case done <- struct{}{}:
				default:
				}
			}
			defer func() { logrus.StandardLogger().ExitFunc = originalExitFunc }()

			setupLogger(os.Stdout, tc.levelToSet)

			output := captureOutput(func() {
				// Recover from panic to allow test to continue
				defer func() {
					if r := recover(); r != nil {
						// A panic is expected for levels below panic
					}
				}()
				runLogbookOperations()
			})

			// For fatal logs, wait for the mocked exit with timeout to avoid hangs
			if tc.levelToSet != "panic" {
				select {
				case <-done:
				case <-time.After(2 * time.Second):
					t.Fatal("timed out waiting for logrus.Fatal ExitFunc")
				}
			}

			lines := strings.Split(strings.TrimSpace(output), "\n")
			var loggedMessages []string

			for _, line := range lines {
				if line == "" {
					continue
				}
				var entry logEntry
				err := json.Unmarshal([]byte(line), &entry)
				require.NoError(t, err, "Log output should be valid JSON")
				loggedMessages = append(loggedMessages, entry.Msg)
			}

			for _, expected := range tc.expectedLogs {
				assert.Contains(t, loggedMessages, expected, "Expected log message not found")
			}
			for _, unexpected := range tc.unexpectedLogs {
				assert.NotContains(t, loggedMessages, unexpected, "Unexpected log message was found")
			}
		})
	}
}

func TestFatalLogsExit(t *testing.T) {
	var exited bool
	originalExitFunc := logrus.StandardLogger().ExitFunc
	logrus.StandardLogger().ExitFunc = func(code int) {
		exited = true
	}

	defer func() {
		logrus.StandardLogger().ExitFunc = originalExitFunc
	}()
    
    // Add a recover block to prevent a subsequent Panic call from crashing this test
	defer func() {
		if r := recover(); r != nil {
			// A panic might occur after Fatal in the test context, which we can ignore
		}
	}()

	setupLogger(io.Discard, "fatal")

	// We wrap this in a function to ensure defer is called even if Fatal exits
	func() {
		// We expect runLogbookOperations to call Fatal, which will call our mocked exit func
		runLogbookOperations()
	}()

	assert.True(t, exited, "logrus.Fatal should call the exit function")
}

func TestPanicLogsPanic(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "Expected a panic to occur")
	}()

	setupLogger(io.Discard, "panic")

	// Mock exit function to prevent it from terminating before panic
	originalExitFunc := logrus.StandardLogger().ExitFunc
	logrus.StandardLogger().ExitFunc = func(int) { /* Do nothing */ }

	defer func() { logrus.StandardLogger().ExitFunc = originalExitFunc }()

	runLogbookOperations()
}