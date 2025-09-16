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
	ctx, cancel := context.WithCancel(parent)
	// Hint: Use context.WithCancel(parent)
	return ctx, cancel
}

// CreateTimeoutContext creates a context with timeout
func (cm *simpleContextManager) CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	// TODO: Implement timeout context creation
	ctx, cancel := context.WithTimeout(parent, timeout)
	// Hint: Use context.WithTimeout(parent, timeout)
	return ctx, cancel
}

// AddValue adds a key-value pair to the context
func (cm *simpleContextManager) AddValue(parent context.Context, key, value interface{}) context.Context {
	// TODO: Implement value context creation
	ctx := context.WithValue(parent, key, value)
	// Hint: Use context.WithValue(parent, key, value)
	return ctx
}

// GetValue retrieves a value from the context
func (cm *simpleContextManager) GetValue(ctx context.Context, key interface{}) (interface{}, bool) {
	// TODO: Implement value retrieval from context
	// Hint: Use ctx.Value(key) and check if it's nil
	val := ctx.Value(key)
	if val != nil {
		return val, true
	}
	// Return the value and a boolean indicating if it was found
	return nil, false
}

// ExecuteWithContext executes a task that can be cancelled via context
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {
	// 创建一个 channel 来接收任务结果
	resultChan := make(chan error, 1)

	// 在 goroutine 中执行任务
	go func() {
		resultChan <- task()
	}()

	// 等待任务完成或 context 被取消
	select {
	case err := <-resultChan:
		// 任务完成，返回任务结果
		return err
	case <-ctx.Done():
		// context 被取消，返回 context 错误
		return ctx.Err()
	}
}

// WaitForCompletion waits for a duration or until context is cancelled
func (cm *simpleContextManager) WaitForCompletion(ctx context.Context, duration time.Duration) error {
	// TODO: Implement waiting with context awareness
	// Hint: Use select with ctx.Done() and time.After(duration)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
	// Return context error if cancelled, nil if duration completes
}

// Helper function - simulate work that can be cancelled
func SimulateWork(ctx context.Context, workDuration time.Duration, description string) error {
	// TODO: Implement cancellable work simulation
	// Hint: Use select with ctx.Done() and time.After(workDuration)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(workDuration):
		fmt.Println(description)
		return nil
	}
	// Print progress messages and respect cancellation
}

// Helper function - process multiple items with context
func ProcessItems(ctx context.Context, items []string) ([]string, error) {
	// TODO: Implement batch processing with context awareness
	// Process each item but check for cancellation between items
	processed := make([]string, 0)
	for _, item := range items {
		select {
		case <-ctx.Done():
			return processed, ctx.Err()
		default:
			time.Sleep(50 * time.Millisecond)
			processed = append(processed, "processed_"+item)
		}
	}
	return processed, nil
	// Return partial results if cancelled
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
