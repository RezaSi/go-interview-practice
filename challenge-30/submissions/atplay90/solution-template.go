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
    value := ctx.Value(key)
    if value != nil {
        return value, true
    }
    return nil, false
}

// ExecuteWithContext executes a task that can be cancelled via context
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {
	// TODO: Implement task execution with context cancellation
	// Hint: Run the task in a goroutine and use select with ctx.Done()
	// Return context error if cancelled, task error if task fails
	ch := make(chan error, 1)
	
	go func() {
		ch <- task()
	}()
	
	
	select {
	case <-ctx.Done():
		return ctx.Err() 
	case err := <-ch:
		return err
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
	fmt.Printf("[START] Starting work: %s (Duration: %v)\n", description, workDuration)

	select {
	case <-ctx.Done():
		// Context was already canceled or timed out before we started
		fmt.Printf("[CANCELED] Work aborted early: %s\n", description)
		return ctx.Err()

	case <-time.After(workDuration):
		// Work completed successfully
		fmt.Printf("[SUCCESS] Finished work: %s\n", description)
		return nil
	}
}

// Helper function - process multiple items with context
func ProcessItems(ctx context.Context, items []string) ([]string, error) {
	processedResults := make([]string, 0, len(items))

	for _, item := range items {
		// Non-blocking check for context cancellation before starting the next item
		select {
		case <-ctx.Done():
			fmt.Printf("[BATCH INTERRUPTED] Context canceled. Returning %d completed items.\n", len(processedResults))
			return processedResults, ctx.Err()
		default:
			// Context is still active, proceed to processing
		}

		// Simulate processing overhead per item
		fmt.Printf("[BATCH] Processing item: %s\n", item)
		time.Sleep(200 * time.Millisecond) // Replace with real work execution
		
		// Collect successful output
		processedResults = append(processedResults, "processed_"+item)
	}

	fmt.Printf("[BATCH SUCCESS] All %d items processed successfully.\n", len(processedResults))
	return processedResults, nil
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
