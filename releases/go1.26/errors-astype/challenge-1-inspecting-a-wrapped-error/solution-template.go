package main

import (
	"errors"
	"fmt"
)

// ValidationError carries the field that failed
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

// RateLimitError carries how long to wait
type RateLimitError struct {
	RetryAfterSeconds int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited, retry in %ds", e.RetryAfterSeconds)
}

// ErrNotFound is a sentinel, matched by value rather than by type
var ErrNotFound = errors.New("not found")

// FailedField returns the field name from a ValidationError anywhere in err's
// chain, and false when there is none.
func FailedField(err error) (string, bool) {
	// TODO: use errors.AsType to find a *ValidationError
	return "", false
}

// RetryAfter returns the wait from a RateLimitError anywhere in err's chain,
// and false when there is none.
func RetryAfter(err error) (int, bool) {
	// TODO: use errors.AsType to find a *RateLimitError
	return 0, false
}

// Classify names the first thing it recognises in the chain:
// "validation", "ratelimit", "notfound", or "unknown".
func Classify(err error) string {
	// TODO: check for a ValidationError, then a RateLimitError, then ErrNotFound
	return "unknown"
}

func main() {
	err := fmt.Errorf("saving user: %w",
		fmt.Errorf("validating: %w", &ValidationError{Field: "email", Reason: "missing @"}))

	if field, ok := FailedField(err); ok {
		fmt.Println("bad field:", field)
	}
	fmt.Println("class:", Classify(err))
}
