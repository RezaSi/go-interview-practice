package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// hookLogEntry is used to unmarshal the JSON output from the ErrorHook
type hookLogEntry struct {
	Level        string `json:"level"`
	Msg          string `json:"msg"`
	RetryAttempt int    `json:"retry_attempt"`
	Time         string `json:"time"`
}

func TestTaskSchedulerLogging(t *testing.T) {
	// 1. Setup
	// Create separate buffers to capture the output of the main logger and the hook
	mainOut := &bytes.Buffer{}
	hookOut := &bytes.Buffer{}

	// Create and configure the main logger
	logger := logrus.New()
	logger.SetOutput(mainOut)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true, // Disable colors for consistent testing
		DisableTimestamp: true, // Disable timestamp for simpler string matching
	})
	// Set level to Info so all our messages are processed
	logger.SetLevel(logrus.InfoLevel)

	// Create and configure the custom hook
	hook := &ErrorHook{
		Out:       hookOut,
		Formatter: &logrus.JSONFormatter{},
	}
	logger.AddHook(hook)

	// 2. Execution
	// Run the scheduler simulation which will generate logs
	runTaskScheduler(logger)

	// 3. Assertions for the Main Logger Output (Text Format)
	mainLogStr := mainOut.String()

	// Verify that all expected log messages are present in the main output
	assert.Contains(t, mainLogStr, "level=info msg=\"Starting task scheduler...\"", "Main logger should contain the starting message")
	assert.Contains(t, mainLogStr, "msg=\"Task 'Process daily reports' completed successfully\"", "Main logger should contain the success message")
	assert.Contains(t, mainLogStr, "task_duration=250ms", "Main logger should contain the duration field")
	assert.Contains(t, mainLogStr, "msg=\"Task 'Sync user data': upstream API is slow\"", "Main logger should contain the warning message")
	assert.Contains(t, mainLogStr, "task_id=sync-001", "Main logger should contain the task_id field")
	assert.Contains(t, mainLogStr, "msg=\"Task 'Backup database' failed: connection timed out\"", "Main logger should contain the error message")
	assert.Contains(t, mainLogStr, "retry_attempt=3", "Main logger should contain the retry_attempt field")

	// 4. Assertions for the ErrorHook Output (JSON Format)
	hookLogStr := hookOut.String()

	// The hook's output should not be empty
	require.NotEmpty(t, hookLogStr, "Hook output should not be empty")

	// The hook should ONLY contain the error log, not info or warning logs
	assert.NotContains(t, hookLogStr, "Process daily reports", "Hook should not log info messages")
	assert.NotContains(t, hookLogStr, "Sync user data", "Hook should not log warning messages")

	// Unmarshal the JSON output from the hook to verify its content and structure
	var entry hookLogEntry
	err := json.Unmarshal([]byte(hookLogStr), &entry)
	require.NoError(t, err, "Hook output must be valid JSON")

	// Verify the fields of the JSON log entry
	assert.Equal(t, "error", entry.Level, "Log level in hook output should be 'error'")
	assert.Equal(t, "Task 'Backup database' failed: connection timed out", entry.Msg, "Message in hook output is incorrect")
	assert.Equal(t, 3, entry.RetryAttempt, "retry_attempt field in hook output is incorrect")
}

func TestErrorHook_Levels(t *testing.T) {
	hook := &ErrorHook{}
	levels := hook.Levels()
	assert.Contains(t, levels, logrus.ErrorLevel, "Hook should fire on ErrorLevel")
	assert.Contains(t, levels, logrus.FatalLevel, "Hook should fire on FatalLevel")
	assert.NotContains(t, levels, logrus.WarnLevel, "Hook should not fire on WarnLevel")
	assert.NotContains(t, levels, logrus.InfoLevel, "Hook should not fire on InfoLevel")
	assert.Len(t, levels, 2, "Hook should only register for 2 levels")
}