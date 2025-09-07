package main

import (
	"context"
	"fmt"
	"time"
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
	// TODO: Implement cancellable context creation
	// Hint: Use context.WithCancel(parent)
	ctx, cancel := context.WithCancel(parent)
	return ctx, cancel
}

// CreateTimeoutContext creates a context with timeout
func (cm *simpleContextManager) CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	// TODO: Implement timeout context creation
	// Hint: Use context.WithTimeout(parent, timeout)
	ctx, cancel := context.WithTimeout(parent, timeout)
	return ctx, cancel
}

// AddValue adds a key-value pair to the context
func (cm *simpleContextManager) AddValue(parent context.Context, key, value interface{}) context.Context {
	// TODO: Implement value context creation
	// Hint: Use context.WithValue(parent, key, value)
	ctx := context.WithValue(parent, key, value)
	return ctx
}

// GetValue retrieves a value from the context
func (cm *simpleContextManager) GetValue(ctx context.Context, key interface{}) (interface{}, bool) {
	// TODO: Implement value retrieval from context
	// Hint: Use ctx.Value(key) and check if it's nil
	// Return the value and a boolean indicating if it was found
	if val := ctx.Value(key); val != nil {
		return val, true
	}
	return nil, false
}

// ExecuteWithContext executes a task that can be cancelled via context
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {
	// TODO: Implement task execution with context cancellation
	// Hint: Run the task in a goroutine and use select with ctx.Done()
	// Return context error if cancelled, task error if task fails
	errCh := make(chan error, 1)

	go func() {
		errCh <- task()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}

}

// WaitForCompletion waits for a duration or until context is cancelled
func (cm *simpleContextManager) WaitForCompletion(ctx context.Context, duration time.Duration) error {
	// TODO: Implement waiting with context awareness
	// Hint: Use select with ctx.Done() and time.After(duration)
	// Return context error if cancelled, nil if duration completes
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}

// Helper function - simulate work that can be cancelled
func SimulateWork(ctx context.Context, workDuration time.Duration, description string) error {
	// TODO: Implement cancellable work simulation
	// Hint: Use select with ctx.Done() and time.After(workDuration)
	// Print progress messages and respect cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(workDuration):
		fmt.Println(description)
		return nil
	}
}

// Helper function - process multiple items with context
func ProcessItems(ctx context.Context, items []string) ([]string, error) {
	// TODO: Implement batch processing with context awareness
	// Process each item but check for cancellation between items
	// Return partial results if cancelled
	result := make([]string, 0, len(items))

	for _, v := range items {
		select {
		case <-ctx.Done():
			return result, ctx.Err() 
		default: 
		}

		time.Sleep(50 * time.Millisecond)
		
		processedItem := fmt.Sprint("processed_", v)
		result = append(result, processedItem)
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
