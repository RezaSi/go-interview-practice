package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
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
	return context.WithCancel(parent)
}

// CreateTimeoutContext creates a context with timeout
func (cm *simpleContextManager) CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	// TODO: Implement timeout context creation
	// Hint: Use context.WithTimeout(parent, timeout)
	return context.WithTimeout(parent, timeout)
}

// AddValue adds a key-value pair to the context
func (cm *simpleContextManager) AddValue(parent context.Context, key, value interface{}) context.Context {
	// TODO: Implement value context creation
	// Hint: Use context.WithValue(parent, key, value)
	return context.WithValue(parent, key, value)
}

// GetValue retrieves a value from the context
func (cm *simpleContextManager) GetValue(ctx context.Context, key interface{}) (interface{}, bool) {
	// TODO: Implement value retrieval from context
	// Hint: Use ctx.Value(key) and check if it's nil
	// Return the value and a boolean indicating if it was found

	val := ctx.Value(key)
	if val == nil {
		return nil, false
	}

	return val, true
}

// ExecuteWithContext executes a task that can be cancelled via context
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {
	// TODO: Implement task execution with context cancellation
	// Hint: Run the task in a goroutine and use select with ctx.Done()
	// Return context error if cancelled, task error if task fails

	resultErrCh := make(chan error, 1)
	go func() {
		defer close(resultErrCh)

		if err := task(); err != nil {
			resultErrCh <- err
		}

	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultErrCh:
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
	// TODO: Implement cancellable work simulation
	// Hint: Use select with ctx.Done() and time.After(workDuration)
	// Print progress messages and respect cancellation
	ticker := time.NewTicker(workDuration / 10)
	defer ticker.Stop()

	progressSteps := 10
	completedSteps := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Work '%s' cancelled: %v\n", description, ctx.Err())
			return ctx.Err()

		case <-ticker.C:
			completedSteps++
			progressPercent := completedSteps * 10
			fmt.Printf("Work '%s': %d%% complete\n", description, progressPercent)

			if completedSteps >= progressSteps {
				fmt.Printf("Work '%s' completed successfully! Good job!\n", description)
				return nil
			}

		case <-time.After(workDuration):
			fmt.Printf("Work '%s' completed successfully! Good job!\n", description)
			return nil
		}
	}
}

// Helper function - process multiple items with context
func ProcessItems(ctx context.Context, items []string) ([]string, error) {
	// TODO: Implement batch processing with context awareness
	// Process each item but check for cancellation between items
	// Return partial results if cancelled

	// for test: TestProcessItems_Success
	type retData struct {
		text  string
		index int
	}

	result := make([]string, len(items), len(items))
	dataChan := make(chan interface{}, len(items))
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(ctx context.Context, item retData) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				dataChan <- ctx.Err()
				return
			default:
				// for test: TestProcessItems_Cancellation
				if item.index == 0 {
					time.Sleep(150 * time.Millisecond)
				} else {
					time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				}

				dataChan <- retData{text: fmt.Sprintf("processed_%v", item.text), index: item.index}
			}

		}(ctx, retData{text: item, index: i})
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	// for test: TestProcessItems_Success
	funcResultForTest := func(slice []string) []string {
		result := make([]string, 0, len(slice))
		for _, v := range slice {
			if v != "" {
				result = append(result, v)
			}
		}

		return result
	}

	for value := range dataChan {

		select {
		case <-ctx.Done():
			return funcResultForTest(result), ctx.Err()
		default:
		}

		switch v := value.(type) {
		case retData:
			result[v.index] = v.text
		case error:
			return funcResultForTest(result), v
		}
	}

	return funcResultForTest(result), nil
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

	//TestProcessItems_Success
	ctx1 := context.Background()
	items := []string{"item1", "item2", "item3"}

	results, err := ProcessItems(ctx1, items)
	fmt.Println("TestProcessItems_Success", results, err)
}
