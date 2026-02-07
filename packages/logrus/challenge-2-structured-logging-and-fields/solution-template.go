package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// loggerKey is a custom type used as a key for storing the logger in the context
type loggerKey int

const key loggerKey = 0

// loggingMiddleware creates a new logger instance for each HTTP request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Generate a new UUID for the request_id
		// Hint: Use uuid.New().String()
		requestID := "" // Replace with UUID generation

		// TODO: Create a new logrus.Entry using logrus.WithFields()
		// This logger should be pre-populated with the following fields:
		// - "request_id": The UUID generated above
		// - "http_method": The request's method (r.Method)
		// - "uri": The request's URI (r.RequestURI)
		// - "user_agent": The request's User-Agent header (r.UserAgent())
		logger := logrus.WithFields(logrus.Fields{
			// Add fields here
		})

		// TODO: Log an informational message "Request received"
		// The fields you added above will be automatically included


		// TODO: Create a new context from the request's context and add the
		// enriched logger to it using the `key`
		// Hint: ctx := context.WithValue(r.Context(), key, logger)
		ctx := r.Context() // Replace this line

		// TODO: Call the next handler in the chain, passing the new request
		// with the updated context
		// Hint: next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// helloHandler is the final handler that processes the request
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Retrieve the logger from the context using the `key`
	// The value from the context will be of type `interface{}`, so you'll need
	// to perform a type assertion to get a `*logrus.Entry`
	// Example: logger, ok := r.Context().Value(key).(*logrus.Entry)
	// If `ok` is false (logger not found), fall back to the global logger `logrus.StandardLogger()`
	var logger *logrus.Entry = logrus.NewEntry(logrus.StandardLogger()) // Replace this line

	// TODO: Add a new field "user_id" to the logger with a sample value "user-99"
	// Hint: logger = logger.WithField("user_id", "user-99")


	// TODO: Log an informational message "Processing hello request"
	// This log should include both the fields from the middleware and the "user_id" field


	// TODO: Write a "Hello, world!" response to the client
	// Hint: fmt.Fprintln(w, "...")
}

func main() {
	// TODO: Set the global logrus formatter to a new instance of logrus.JSONFormatter


	// Create the handler chain by wrapping the helloHandler with the loggingMiddleware
	finalHandler := loggingMiddleware(http.HandlerFunc(helloHandler))

	// Configure HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}

	logrus.Info("Server starting on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}