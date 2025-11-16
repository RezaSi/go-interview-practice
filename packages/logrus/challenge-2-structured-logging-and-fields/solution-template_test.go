package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// logEntry is a struct to unmarshal the JSON log output for verification
// We add all expected fields from both the middleware and the handler
type logEntry struct {
	Level     string `json:"level"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
	Method    string `json:"http_method"`
	URI       string `json:"uri"`
	UserAgent string `json:"user_agent"`
	UserID    string `json:"user_id"`
}

// isValidUUID checks if a string is a valid UUID.
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// NOTE: These tests assume the TODOs in solution-template.go
// are implemented. If run against the raw template, they will fail

// For JSON Format
func TestLoggingMiddleware_JSONFormatter(t *testing.T) {
	// Save original configuration
	originalFormatter := logrus.StandardLogger().Formatter
	originalOutput := logrus.StandardLogger().Out
	t.Cleanup(func() {
		logrus.SetFormatter(originalFormatter)
		logrus.SetOutput(originalOutput)
	})

	// Ensure the global logger is set to JSON format for the test
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Capture log output by redirecting the logger's output to a buffer
	var logBuffer bytes.Buffer
	logrus.SetOutput(&logBuffer)

	// Create a test server with the full handler chain
	handler := loggingMiddleware(http.HandlerFunc(helloHandler))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a new HTTP request to the test server
	req, err := http.NewRequest("GET", server.URL+"/hello", nil)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "Test-Client-1.0")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert basic HTTP response correctness
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Split the captured log output into individual JSON lines
	logLines := strings.Split(strings.TrimSpace(logBuffer.String()), "\n")
	require.Len(t, logLines, 2, "Expected two log entries: one from middleware, one from handler")

	// Unmarshal and verify the first log entry (from middleware)
	var entry1 logEntry
	err = json.Unmarshal([]byte(logLines[0]), &entry1)
	require.NoError(t, err, "First log entry should be valid JSON")

	assert.Equal(t, "info", entry1.Level)
	assert.Equal(t, "Request received", entry1.Msg)
	assert.True(t, isValidUUID(entry1.RequestID), "request_id should be a valid UUID")
	assert.Equal(t, "GET", entry1.Method)
	assert.Equal(t, "/hello", entry1.URI)
	assert.Equal(t, "Test-Client-1.0", entry1.UserAgent)
	assert.Empty(t, entry1.UserID, "user_id should not be set by the middleware")

	// Unmarshal and verify the second log entry (from handler)
	var entry2 logEntry
	err = json.Unmarshal([]byte(logLines[1]), &entry2)
	require.NoError(t, err, "Second log entry should be valid JSON")

	assert.Equal(t, "info", entry2.Level)
	assert.Equal(t, "Processing hello request", entry2.Msg)
	assert.Equal(t, "user-99", entry2.UserID, "user_id should be set by the handler")

	// Correlation verification
	assert.NotEmpty(t, entry1.RequestID, "request_id must not be empty")
	assert.Equal(t, entry1.RequestID, entry2.RequestID, "request_id must be the same in both log entries")

	// Ensure middleware fields are still present in the handler log
	assert.Equal(t, "GET", entry2.Method, "http_method should propagate to handler log")
	assert.Equal(t, "/hello", entry2.URI, "uri should propagate to handler log")
	assert.Equal(t, "Test-Client-1.0", entry2.UserAgent, "user_agent should propagate to handler log")
}

// For Text Format
func TestLoggingMiddleware_TextFormatter(t *testing.T) {
	// Save original configuration
	originalFormatter := logrus.StandardLogger().Formatter
	originalOutput := logrus.StandardLogger().Out
	t.Cleanup(func() {
		logrus.SetFormatter(originalFormatter)
		logrus.SetOutput(originalOutput)
	})

	// Ensure the global logger is set to Text format
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Capture log output by redirecting the logger's output to a buffer
	var logBuffer bytes.Buffer
	logrus.SetOutput(&logBuffer)

	// Create a test server with the full handler chain
	handler := loggingMiddleware(http.HandlerFunc(helloHandler))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a new HTTP request to the test server
	req, err := http.NewRequest("GET", server.URL+"/hello", nil)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "Test-Client-Text")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get the log output
	logOutput := logBuffer.String()

	// Checking if key substrings exist
	assert.Contains(t, logOutput, "level=info")
	assert.Contains(t, logOutput, "msg=\"Request received\"")
	assert.Contains(t, logOutput, "msg=\"Processing hello request\"")
	assert.Contains(t, logOutput, "request_id=")
	assert.Contains(t, logOutput, "http_method=GET")
	assert.Contains(t, logOutput, "uri=/hello")
	assert.Contains(t, logOutput, "user_agent=Test-Client-Text")
	assert.Contains(t, logOutput, "user_id=user-99")
}