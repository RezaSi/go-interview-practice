package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// logLine is a helper map type for unmarshalling JSON log lines produced by logrus.JSONFormatter
type logLine map[string]interface{}

// helper: read lines and unmarshal JSON lines into logLine maps
func parseJSONLines(r io.Reader) ([]logLine, error) {
	var out []logLine
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var m logLine
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			// If a line isn't valid JSON, skip it (some formatters may emit text lines)
			continue
		}
		out = append(out, m)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// findRetryHook searches logger hooks for a hook of the concrete type *RetryHook
func findRetryHook(logger *logrus.Logger) *RetryHook {
	// Look in Warn and Error levels (hooks for those levels)
	levels := []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel}
	for _, lvl := range levels {
		hooks := logger.Hooks[lvl]
		for _, h := range hooks {
			// Try to assert to *RetryHook
			if rh, ok := h.(*RetryHook); ok {
				return rh
			}
		}
	}
	return nil
}

// Test that runWorker emits expected Info/Warn/Error messages and fields when using JSONFormatter
func TestRunWorker_EmitsExpectedStructuredLogs(t *testing.T) {
	assert := assert.New(t)

	// Create a test logger and capture output
	logger := logrus.New()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Prepare tasks similar to the example
	tasks := []Task{
		{ID: "1", Name: "SendEmail", Retries: 0},
		{ID: "2", Name: "FailMe", Retries: 2},
		{ID: "3", Name: "GenerateReport", Retries: 1},
	}

	// Run the worker (should write JSON lines into buf)
	runWorker(logger, tasks)

	// Give a tiny margin for any buffering (not usually necessary)
	time.Sleep(10 * time.Millisecond)

	lines, err := parseJSONLines(&buf)
	assert.NoError(err, "failed to parse JSON log lines")
	assert.Greater(len(lines), 0, "expected some JSON log lines")

	// Helpers to find an entry by message substring
	findByMsg := func(substr string) (logLine, bool) {
		for _, l := range lines {
			if msgRaw, ok := l["msg"]; ok {
				if s, ok := msgRaw.(string); ok && strings.Contains(s, substr) {
					return l, true
				}
			}
		}
		return nil, false
	}

	// 1 - Starting task log for SendEmail
	entry, ok := findByMsg("Starting task: SendEmail")
	assert.True(ok, "expected a 'Starting task: SendEmail' log entry")
	if ok {
		// Check task_id exists and matches
		assert.Equal("1", entry["task_id"], "task_id should be '1' in start entry")
	}

	// 2 - Retry warning for GenerateReport (Retries:1) or a retry for some task
	entry, ok = findByMsg("retried")
	assert.True(ok, "expected a retry warning log entry (message contains 'retried')")
	if ok {
		// Ensure retries field exists and is a number >= 1
		rRaw, exists := entry["retries"]
		assert.True(exists, "expected 'retries' field in retry warning entry")
		if exists {
			// JSON numbers unmarshal to float64
			f, ok := rRaw.(float64)
			assert.True(ok, "retries field should be a number")
			assert.GreaterOrEqual(f, float64(1), "retries should be >= 1")
		}
	}

	// 3 - Error entry for FailMe
	entry, ok = findByMsg("Task failed")
	assert.True(ok, "expected an error log entry for simulated failure")
	if ok {
		// Check task_id and error field presence
		assert.Equal("2", entry["task_id"], "expected task_id '2' in error entry")
		errRaw, exists := entry["error"]
		assert.True(exists, "expected 'error' field in failure entry")
		if exists {
			s, ok := errRaw.(string)
			assert.True(ok && s != "", "unexpected 'error' field value")
		}
	}

	// 4 - Success entry contains duration
	entry, ok = findByMsg("completed successfully")
	assert.True(ok, "expected a success completion log entry")
	if ok {
		_, hasDuration := entry["duration"]
		assert.True(hasDuration, "expected 'duration' field in success entry")
	}
}

// Test that setupLogger registers a RetryHook that actually writes formatted entries
// for Warn and Error levels to its configured Out target
func TestRetryHook_WritesOutOnWarnAndError(t *testing.T) {
	assert := assert.New(t)

	// Call the participant's setupLogger to get a configured logger
	logger := setupLogger()
	assert.NotNil(logger, "setupLogger should return a logger")

	// Find the retry hook within logger hooks
	retry := findRetryHook(logger)
	assert.NotNil(retry, "expected a RetryHook to be registered for warn/error levels (found none). Ensure setupLogger adds the hook.")

	// Create a temporary file to capture hook output
	tmpFile, err := os.CreateTemp("", "retryhook-*.log")
	assert.NoError(err, "failed to create temp file")
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	// Use reflection to set Out and Formatter fields on the found hook if they exist and are settable
	rv := reflect.ValueOf(retry)
	// Expect a pointer to struct
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	outField := rv.FieldByName("Out")
	fmtField := rv.FieldByName("Formatter")

	// Validate fields exist
	assert.True(outField.IsValid(), "RetryHook must have an 'Out' field (io.Writer) for tests to redirect output")
	assert.True(fmtField.IsValid(), "RetryHook must have a 'Formatter' field (logrus.Formatter) for tests to set JSON formatting")

	// Attempt to set fields
	if outField.IsValid() && outField.CanSet() {
		outField.Set(reflect.ValueOf(tmpFile))
	} else {
		// If the field exists but is not settable, fail with helpful message
		if outField.IsValid() && !outField.CanSet() {
			t.Fatalf("RetryHook.Out exists but is not settable; ensure it is an exported field (capitalized) and addressable")
		}
	}

	if fmtField.IsValid() && fmtField.CanSet() {
		fmtField.Set(reflect.ValueOf(&logrus.JSONFormatter{}))
	} else {
		if fmtField.IsValid() && !fmtField.CanSet() {
			t.Fatalf("RetryHook.Formatter exists but is not settable; ensure it is an exported field (capitalized) and addressable")
		}
	}

	// Fire a Warn and Error log which should call the hook and write to our temp file
	logger.Warn("hook-test-warn")
	logger.Error("hook-test-error")

	// Flush/close so data is written
	tmpFile.Sync()
	tmpFile.Close()

	// Read contents and ensure our messages are present (as JSON text)
	content, err := os.ReadFile(tmpFile.Name())
	assert.NoError(err, "failed to read temp hook file")
	text := string(content)
	assert.Contains(text, "hook-test-warn", "expected hook to write warn message")
	assert.Contains(text, "hook-test-error", "expected hook to write error message")
}

// A small smoke test: 
// ensure setupLogger returns a logger that can be used by runWorker without panicking
func TestIntegration_RunWorkerWithSetupLogger(t *testing.T) {
	assert := assert.New(t)

	logger := setupLogger()
	assert.NotNil(logger, "setupLogger should return a logger")

	// Capture console output so test logs don't pollute test stdout
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Minimal tasks to run without assertions - ensure no panic and produce some output
	tasks := []Task{
		{ID: "x", Name: "SendEmail", Retries: 0},
	}
	runWorker(logger, tasks)

	// Expect at least one JSON log line produced
	lines, err := parseJSONLines(&buf)
	assert.NoError(err, "failed to parse JSON lines from logger output")
	assert.Greater(len(lines), 0, "expected at least one log line when running worker with logger from setupLogger")
}