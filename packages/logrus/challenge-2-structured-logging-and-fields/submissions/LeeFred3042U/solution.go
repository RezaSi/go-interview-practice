package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type loggerKeyType int

const loggerKey loggerKeyType = 0

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new request ID
		requestID := uuid.New().String()

		// Create a logger with the middleware fields
		logger := logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"http_method": r.Method,
			"uri": r.RequestURI,
			"user_agent": r.UserAgent(),
		})

		// Log the initial middleware message
		logger.Info("Request received")

		// Add logger to context
		ctx := context.WithValue(r.Context(), loggerKey, logger)

		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Retrieve logger from context
func getLogger(r *http.Request) *logrus.Entry {
	logger, ok := r.Context().Value(loggerKey).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return logger
}

// Final handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	logger := getLogger(r).WithField("user_id", "user-99")
	logger.Info("Processing hello request")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, world!"))
}

func main() {
	// Set global logger formatter
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Wrap handler with middleware
	finalHandler := loggingMiddleware(http.HandlerFunc(helloHandler))

	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		logrus.WithError(err).Fatal("Server failed to start")
	}
}