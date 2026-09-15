package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestFailedFieldFindsAnUnwrappedError(t *testing.T) {
	err := &ValidationError{Field: "email", Reason: "missing @"}

	got, ok := FailedField(err)
	if !ok {
		t.Fatal("FailedField missed a bare *ValidationError")
	}
	if got != "email" {
		t.Fatalf("FailedField = %q, want %q", got, "email")
	}
}

// TestFailedFieldWalksTheChain is the reason this is not a type assertion.
func TestFailedFieldWalksTheChain(t *testing.T) {
	err := fmt.Errorf("saving user: %w",
		fmt.Errorf("validating: %w", &ValidationError{Field: "age", Reason: "negative"}))

	got, ok := FailedField(err)
	if !ok {
		t.Fatal("FailedField did not look through two layers of wrapping; a type assertion only sees the outermost error")
	}
	if got != "age" {
		t.Fatalf("FailedField = %q, want %q", got, "age")
	}
}

func TestFailedFieldReportsNoMatch(t *testing.T) {
	err := fmt.Errorf("saving user: %w", ErrNotFound)

	if got, ok := FailedField(err); ok {
		t.Fatalf("FailedField found %q in a chain with no ValidationError", got)
	}
}

func TestFailedFieldHandlesNil(t *testing.T) {
	if _, ok := FailedField(nil); ok {
		t.Fatal("FailedField reported a match for a nil error")
	}
}

func TestRetryAfterFindsTheWait(t *testing.T) {
	err := fmt.Errorf("calling api: %w", &RateLimitError{RetryAfterSeconds: 30})

	got, ok := RetryAfter(err)
	if !ok {
		t.Fatal("RetryAfter missed a wrapped *RateLimitError")
	}
	if got != 30 {
		t.Fatalf("RetryAfter = %d, want 30", got)
	}
}

func TestRetryAfterIgnoresOtherTypes(t *testing.T) {
	err := fmt.Errorf("saving: %w", &ValidationError{Field: "name", Reason: "blank"})

	if got, ok := RetryAfter(err); ok {
		t.Fatalf("RetryAfter matched a ValidationError and returned %d", got)
	}
}

func TestClassifyNamesEachKind(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"validation", &ValidationError{Field: "email", Reason: "bad"}, "validation"},
		{"wrapped validation", fmt.Errorf("a: %w", &ValidationError{Field: "e"}), "validation"},
		{"ratelimit", &RateLimitError{RetryAfterSeconds: 5}, "ratelimit"},
		{"wrapped ratelimit", fmt.Errorf("b: %w", &RateLimitError{RetryAfterSeconds: 5}), "ratelimit"},
		{"sentinel", ErrNotFound, "notfound"},
		{"wrapped sentinel", fmt.Errorf("c: %w", ErrNotFound), "notfound"},
		{"plain", errors.New("boom"), "unknown"},
		{"nil", nil, "unknown"},
	}

	for _, tc := range cases {
		if got := Classify(tc.err); got != tc.want {
			t.Errorf("Classify(%s) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// TestClassifyUsesIsForTheSentinel checks the right tool for the right question:
// a sentinel is matched by value, not by type.
func TestClassifyUsesIsForTheSentinel(t *testing.T) {
	wrapped := fmt.Errorf("looking up user: %w", fmt.Errorf("querying: %w", ErrNotFound))

	if got := Classify(wrapped); got != "notfound" {
		t.Fatalf("Classify = %q, want %q; use errors.Is for a sentinel value", got, "notfound")
	}
}
