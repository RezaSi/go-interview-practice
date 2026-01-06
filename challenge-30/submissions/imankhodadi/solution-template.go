package main

import (
	"context"
	"fmt"
	"time"
	"sync"
	"strings"
)

type ContextManager interface {
	CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc)
	CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
	AddValue(parent context.Context, key, value interface{}) context.Context
	GetValue(ctx context.Context, key interface{}) (interface{}, bool)
	ExecuteWithContext(ctx context.Context, task func() error) error
	WaitForCompletion(ctx context.Context, duration time.Duration) error
}

type simpleContextManager struct{}

func NewContextManager() *simpleContextManager {
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

func (cm *simpleContextManager) GetValue(ctx context.Context, key interface{}) (interface{}, bool) {
    value := ctx.Value(key)
    if value == nil {
        return nil, false
    }
    return value, true
}

func getStringValue(ctx context.Context, key interface{}) (string, bool) {
    value := ctx.Value(key)
    if str, ok := value.(string); ok {
        return str, true
    }
    return "", false
}
// ExecuteWithContext executes a task that can be cancelled via context
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {
	// Channel to receive task result
    resultChan := make(chan error, 1)
    // Execute task in goroutine
    go func() {
        defer close(resultChan)
        if err := task(); err != nil {
            resultChan <- err
        }
    }()
    // Wait for either task completion or context cancellation
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
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(duration):
        return nil
    }
}
// Enhanced waiting with progress tracking
func (cm *simpleContextManager) WaitWithProgress(ctx context.Context, duration time.Duration, progressCallback func(elapsed time.Duration)) error {
    ticker := time.NewTicker(duration / 10) // 10% intervals
    defer ticker.Stop()
    start := time.Now()
    deadline := start.Add(duration)
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case now := <-ticker.C:
            if now.After(deadline) {
                return nil
            }
            if progressCallback != nil {
                progressCallback(now.Sub(start))
            }
        }
    }
}

// Helper function - simulate work that can be cancelled
func SimulateWork(ctx context.Context, workDuration time.Duration, description string) error {
    if description == "" {
        description = "work"
    }
    // Simulate work in small chunks to allow cancellation
    chunkDuration := time.Millisecond * 100
    chunks := int(workDuration / chunkDuration)
    remainder := workDuration % chunkDuration
    for i := 0; i < chunks; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(chunkDuration):
            // Continue working
        }
    }
    // Handle remainder duration
    if remainder > 0 {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(remainder):
            // Work completed
        }
    }
    return nil
}
// Simulate work with progress reporting
func SimulateWorkWithProgress(ctx context.Context, workDuration time.Duration, description string, progressFn func(float64)) error {
    start := time.Now()
    chunkDuration := time.Millisecond * 50
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(chunkDuration):
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
        processed := fmt.Sprintf("processed_%s",item)
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
    itemChan := make(chan struct{ index int; item string }, len(items))
    resultChan := make(chan result, len(items))
    // Send items to process
    for i, item := range items {
        itemChan <- struct{ index int; item string }{i, item}
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
                    resultChan <- result{work.index, "", ctx.Err()}
                    return
                default:
                    // Process item
                    processed := fmt.Sprintf("processed_%s", strings.ToUpper(work.item))
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
    for result := range resultChan {
        if result.err != nil {
            return nil, result.err
        }
        results[result.index] = result.value
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
// Chain multiple operations with context
func (cm *simpleContextManager) ChainOperations(ctx context.Context, operations []func(context.Context) error) error {
    for i, op := range operations {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := op(ctx); err != nil {
                return fmt.Errorf("operation %d failed: %w", i, err)
            }
        }
    }
    return nil
}
// Rate limited context operations
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

	// Use the context
	fmt.Println("Context created with values!")
}
