package main

import (
	"context"
	"fmt"
	"time"
	"errors"
)

// ContextManager defines a simplified interface for basic context operations
type ContextManager interface {
	// Create a cancellable context from a parent context
	CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc)

	// Create a context with timeout
	CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)

	// Add a value to context
	AddValue(parent context.Context, key, value interface{}) context.Context

	// Get a value from context
	GetValue(ctx context.Context, key interface{}) (interface{}, bool)

	// Execute a task with context cancellation support
	ExecuteWithContext(ctx context.Context, task func() error) error

	// Wait for context cancellation or completion
	WaitForCompletion(ctx context.Context, duration time.Duration) error
}

// Simple context manager implementation
type simpleContextManager struct{}

// NewContextManager creates a new context manager
func NewContextManager() ContextManager {
	return &simpleContextManager{}
}

// CreateCancellableContext creates a cancellable context
func (cm *simpleContextManager) CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

// CreateTimeoutContext creates a context with timeout
func (cm *simpleContextManager) CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// AddValue adds a key-value pair to the context
func (cm *simpleContextManager) AddValue(parent context.Context, key, value interface{}) context.Context {
	return context.WithValue(parent, key, value)
}

// GetValue retrieves a value from the context
func (cm *simpleContextManager) GetValue(ctx context.Context, key interface{}) (interface{}, bool) {
    val := ctx.Value(key)
    if val != nil {
        return val, true
    }
	return nil, false
}

// ExecuteWithContext executes a task that can be cancelled via context
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {
	ch := make(chan error)
	
	go func() {
	    ch <- task()
	    close(ch)
	}()
	
	for {
	    select {
	        case <-ctx.Done():
	            return ctx.Err()
	        case err, ok := <-ch:
	            if !ok {
	                return nil
	            }
	            return err
	    }
	}
	
	return nil
}

// WaitForCompletion waits for a duration or until context is cancelled
func (cm *simpleContextManager) WaitForCompletion(ctx context.Context, duration time.Duration) error {
	for {
	    select {
	        case <-ctx.Done():
	            return ctx.Err()
	        case <-time.After(duration):
	            return nil
	    }
	}
	
	return nil
}

// Helper function - simulate work that can be cancelled
func SimulateWork(ctx context.Context, workDuration time.Duration, _ string) error {
	for {
	    select {
	        case <-ctx.Done():
	            return ctx.Err()
	        case <-time.After(workDuration):
	            return nil
	    }
	}
	
	return nil
}

// Helper function - process multiple items with context
func ProcessItems(ctx context.Context, items []string) ([]string, error) {
	// TODO: Implement batch processing with context awareness
	// Process each item but check for cancellation between items
	// Return partial results if cancelled
	var result []string
	
	ch := make(chan string)
	
	go func() {
	    for _, item := range items {
	        ch <- fmt.Sprintf("processed_%s", item)
	        time.Sleep(50 * time.Millisecond)
	    }
	    close(ch)
	}()
	
	for {
	    select {
	        case <-ctx.Done():
	            err := ctx.Err()

	            if errors.Is(err, context.Canceled) {
	                return result, err
	            }

	            return []string{}, err
            case item, ok := <-ch:
                if !ok {
                    return result, nil
                }
                result = append(result, item)
	    }
	}
	
	return result, nil
}

// Example usage
func main() {
	fmt.Println("Context Management Challenge")
	fmt.Println("Implement the context manager methods!")

	// Example of how the context manager should work:
	cm := NewContextManager()

	// Create a cancellable context
	ctx, cancel := cm.CreateCancellableContext(context.Background())
	defer cancel()

	// Add some values
	ctx = cm.AddValue(ctx, "user", "alice")
	ctx = cm.AddValue(ctx, "requestID", "12345")

	// Use the context
	fmt.Println("Context created with values!")
}
