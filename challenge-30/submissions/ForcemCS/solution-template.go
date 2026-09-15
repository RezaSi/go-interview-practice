package main

import (
	"context"
	"fmt"
	"time"
)

// ContextManager defines a simplified interface for basic context operations
type ContextManager interface {
	CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc)
	CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
	AddValue(parent context.Context, key, value interface{}) context.Context
	GetValue(ctx context.Context, key interface{}) (interface{}, bool)
	ExecuteWithContext(ctx context.Context, task func() error) error
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
	if val == nil {
		return nil, false
	}
	return val, true
}

// ExecuteWithContext executes a task that can be cancelled via context
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {
	// 【重点注意】：这里就是我们之前讨论的精髓！必须使用容量为 1 的缓冲 channel
	// 防止主流程因超时退出后，后台的 task goroutine 发生内存泄漏（阻塞卡死）
	errCh := make(chan error, 1)

	go func() {
		errCh <- task()
	}()

	select {
	case <-ctx.Done():
		// 上下文被取消或超时，直接返回上下文的错误
		return ctx.Err()
	case err := <-errCh:
		// 任务正常或异常执行完毕，返回任务的错误
		return err
	}
}

// WaitForCompletion waits for a duration or until context is cancelled
func (cm *simpleContextManager) WaitForCompletion(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		// 如果在等待期间收到取消信号，立即返回
		return ctx.Err()
	case <-time.After(duration):
		// 正常等待时间结束
		return nil
	}
}

// Helper function - simulate work that can be cancelled
func SimulateWork(ctx context.Context, workDuration time.Duration, description string) error {
	fmt.Printf("开始工作: %s...\n", description)
	
	select {
	case <-ctx.Done():
		fmt.Printf("工作被强制中断: %s\n", description)
		return ctx.Err()
	case <-time.After(workDuration):
		fmt.Printf("工作顺利完成: %s\n", description)
		return nil
	}
}

// Helper function - process multiple items with context
func ProcessItems(ctx context.Context, items []string) ([]string, error) {
	results := make([]string, 0, len(items))
	
	for i, item := range items {
		// Check for cancellation before processing each item
		select {
		case <-ctx.Done():
			fmt.Printf("Processing cancelled after %d items\n", i)
			return results, ctx.Err()
		default:
			// Process the item
			processedItem := fmt.Sprintf("processed_%s", item)
			results = append(results, processedItem)
			
			// Sleep to simulate work, but only if not the last item
			// This prevents unnecessary delay in successful completion
			if i < len(items)-1 {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	return results, nil
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
