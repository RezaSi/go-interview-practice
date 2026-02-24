package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// cannot change these signatures, they are part of the assignment
type ContextManager interface {
	CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc)
	CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
	AddValue(parent context.Context, key, value interface{}) context.Context
	GetValue(ctx context.Context, key interface{}) (interface{}, bool)
	GetStringValue(ctx context.Context, key interface{}) (string, bool)
	ExecuteWithContext(ctx context.Context, task func() error) error
	ExecuteWithContextTimeout(ctx context.Context, task func() error, timeout time.Duration) error
	WaitForCompletion(ctx context.Context, duration time.Duration) error
	WaitWithProgress(ctx context.Context, duration time.Duration, progressCallback func(elapsed time.Duration)) error
	CreateContextWithMultipleValues(parent context.Context, values map[interface{}]interface{}) context.Context
	ExecuteWithCleanup(ctx context.Context, task func() error, cleanup func()) error
	ChainOperations(ctx context.Context, operations []func() error) error
	RateLimitedExecution(ctx context.Context, tasks []func() error, rate time.Duration) error
}

type simpleContextManager struct{}

func NewContextManager() ContextManager {
	return &simpleContextManager{}
}

func (cm *simpleContextManager) CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

func (cm *simpleContextManager) CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

func (cm *simpleContextManager) AddValue(parent context.Context, key, value interface{}) context.Context {
	return context.WithValue(parent, key, value)
}

// be careful, cannot distinguish between "key not found" and "value is nil"
// whether don't store nil as key or assume the key exists with nil value
func (cm *simpleContextManager) GetValue(ctx context.Context, key interface{}) (interface{}, bool) {
	value := ctx.Value(key)
	if value == nil {
		return nil, false
	}
	return value, true
}

func (cm *simpleContextManager) GetStringValue(ctx context.Context, key interface{}) (string, bool) {
	value := ctx.Value(key)
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// ExecuteWithContext executes a context-aware task that can be cancelled.
// The task function must respect context cancellation by checking ctx.Done().
// If the context is cancelled, this returns immediately, but the task goroutine
// continues until it checks ctx.Done() or completes.
// function signature is part of the assignment, cannot change it
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {
	resultChan := make(chan error, 1)
	go func() {
		resultChan <- task()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultChan:
		return err
	}
}

// Alternative implementation with timeout
func (cm *simpleContextManager) ExecuteWithContextTimeout(ctx context.Context, task func() error, timeout time.Duration) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return cm.ExecuteWithContext(timeoutCtx, task)
}

// WaitForCompletion waits for a duration or until context is cancelled
func (cm *simpleContextManager) WaitForCompletion(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// Enhanced waiting with progress tracking
func (cm *simpleContextManager) WaitWithProgress(ctx context.Context, duration time.Duration, progressCallback func(elapsed time.Duration)) error {
	interval := duration / 10
	if interval <= 0 {
		interval = time.Millisecond // minimum interval
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timer := time.NewTimer(duration)
	defer timer.Stop()
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if progressCallback != nil {
				progressCallback(duration)
			}
			return nil
		case <-ticker.C:
			if progressCallback != nil {
				progressCallback(time.Since(start))
			}
		}
	}
}

// Helper function - simulate work that can be cancelled
// Cannot change function signature, it is part of the assignment
func SimulateWork(ctx context.Context, workDuration time.Duration, description string) error {
	timer := time.NewTimer(workDuration)
	fmt.Println("Executing: ", description)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func SimulateWorkWithProgress(ctx context.Context, workDuration time.Duration, description string, progressFn func(float64)) error {
	start := time.Now()
	chunkDuration := time.Millisecond * 50
	ticker := time.NewTicker(chunkDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			elapsed := time.Since(start)
			if elapsed >= workDuration {
				if progressFn != nil {
					progressFn(1.0)
				}
				return nil
			}
			if progressFn != nil {
				progress := float64(elapsed) / float64(workDuration)
				progressFn(progress)
			}
		}
	}
}
func ProcessItems(ctx context.Context, items []string) ([]string, error) {
	if len(items) == 0 {
		return []string{}, nil
	}
	results := make([]string, 0, len(items))
	for i, item := range items {
		// Check for cancellation before processing each item
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
			// Continue processing
		}
		// Simulate item processing time
		processingTime := time.Millisecond * 50
		if err := SimulateWork(ctx, processingTime, fmt.Sprintf("processing item %d", i)); err != nil {
			return results, err
		}
		// Transform the item (example: convert to uppercase)
		processed := fmt.Sprintf("processed_%s", item)
		results = append(results, processed)
	}
	return results, nil
}

// Process items concurrently with context
func ProcessItemsConcurrently(ctx context.Context, items []string, maxWorkers int) ([]string, error) {
	if len(items) == 0 {
		return []string{}, nil
	}
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	type result struct {
		index int
		value string
		err   error
	}
	itemChan := make(chan struct {
		index int
		item  string
	}, len(items))
	resultChan := make(chan result, len(items))

	// Track context cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Send items to process
	for i, item := range items {
		itemChan <- struct {
			index int
			item  string
		}{i, item}
	}
	close(itemChan)
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for work := range itemChan {
				select {
				case <-ctx.Done():
					// Don't process this item, just return
					return
				default:
					processingTime := time.Millisecond * 50
					if err := SimulateWork(ctx, processingTime, "work"); err != nil {
						resultChan <- result{work.index, "", err}
						return
					}
					processed := fmt.Sprintf("processed_%s", work.item)
					resultChan <- result{work.index, processed, nil}
				}
			}
		}()
	}
	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	// Collect results
	results := make([]string, len(items))
	processedIndices := make([]int, 0, len(items))
	for result := range resultChan {
		if result.err != nil {
			cancel() // Signal other workers to stop
			// Drain remaining results
			for range resultChan {
			}
			return nil, result.err
		}
		results[result.index] = result.value
		processedIndices = append(processedIndices, result.index)
	}
	// Check if context was cancelled
	if ctx.Err() != nil {
		// Return only successfully processed results in order
		sort.Ints(processedIndices)
		successfulResults := make([]string, len(processedIndices))
		for i, idx := range processedIndices {
			successfulResults[i] = results[idx]
		}
		return successfulResults, ctx.Err()
	}
	return results, nil
}

// Context with multiple values
func (cm *simpleContextManager) CreateContextWithMultipleValues(parent context.Context, values map[interface{}]interface{}) context.Context {
	ctx := parent
	for key, value := range values {
		ctx = context.WithValue(ctx, key, value)
	}
	return ctx
}

// Timeout with cleanup
func (cm *simpleContextManager) ExecuteWithCleanup(ctx context.Context, task func() error, cleanup func()) error {
	if cleanup != nil {
		defer cleanup()
	}
	return cm.ExecuteWithContext(ctx, task)
}

// cannot change these signatures, they are part of the assignment
// Chain multiple operations with context, later change to func(context.Context) in production
func (cm *simpleContextManager) ChainOperations(ctx context.Context, operations []func() error) error {
	for i, op := range operations {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := op(); err != nil {
				return fmt.Errorf("operation %d failed: %w", i, err)
			}
		}
	}
	return nil
}

// function signature is part of the assignment, later change to func(context.Context) in production
// Rate limited context operations, later change to func(context.Context) in production
func (cm *simpleContextManager) RateLimitedExecution(ctx context.Context, tasks []func() error, rate time.Duration) error {
	ticker := time.NewTicker(rate)
	defer ticker.Stop()
	for i, task := range tasks {
		if i > 0 { // Don't wait before first task
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				// Continue to next task
			}
		}
		if err := cm.ExecuteWithContext(ctx, task); err != nil {
			return fmt.Errorf("task %d failed: %w", i, err)
		}
	}
	return nil
}
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

	fmt.Println(cm.GetStringValue(ctx, "user"))
}
