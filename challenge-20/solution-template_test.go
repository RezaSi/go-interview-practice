package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Test helper functions and mocks
type mockOperation struct {
	shouldFail bool
	delay      time.Duration
	callCount  int
	mutex      sync.Mutex
}

func (m *mockOperation) execute() (interface{}, error) {
	m.mutex.Lock()
	m.callCount++
	m.mutex.Unlock()

	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	if m.shouldFail {
		return nil, errors.New("operation failed")
	}
	return "success", nil
}

func (m *mockOperation) getCallCount() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.callCount
}

func (m *mockOperation) reset() {
	m.mutex.Lock()
	m.callCount = 0
	m.mutex.Unlock()
}

// Basic functionality tests
func TestNewCircuitBreaker(t *testing.T) {
	config := Config{
		MaxRequests: 3,
		Interval:    time.Minute,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 5
		},
	}

	cb := NewCircuitBreaker(config)
	if cb == nil {
		t.Fatal("NewCircuitBreaker should return a non-nil circuit breaker")
	}

	if cb.GetState() != StateClosed {
		t.Errorf("Expected initial state to be Closed, got %v", cb.GetState())
	}

	metrics := cb.GetMetrics()
	if metrics.Requests != 0 || metrics.Successes != 0 || metrics.Failures != 0 {
		t.Errorf("Expected initial metrics to be zero, got %+v", metrics)
	}
}

func TestCircuitBreakerDefaults(t *testing.T) {
	// Test with empty config to verify defaults
	cb := NewCircuitBreaker(Config{})
	if cb == nil {
		t.Fatal("NewCircuitBreaker should handle empty config with defaults")
	}
}

func TestSuccessfulOperations(t *testing.T) {
	config := Config{
		MaxRequests: 3,
		Timeout:     100 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 3
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()
	op := &mockOperation{shouldFail: false}

	// Execute successful operations
	for i := 0; i < 5; i++ {
		result, err := cb.Call(ctx, op.execute)
		if err != nil {
			t.Errorf("Expected no error for successful operation, got %v", err)
		}
		if result != "success" {
			t.Errorf("Expected result 'success', got %v", result)
		}
	}

	// Verify state remains closed
	if cb.GetState() != StateClosed {
		t.Errorf("State should remain Closed after successful operations, got %v", cb.GetState())
	}

	// Verify metrics
	metrics := cb.GetMetrics()
	if metrics.Requests != 5 || metrics.Successes != 5 || metrics.Failures != 0 {
		t.Errorf("Expected 5 requests, 5 successes, 0 failures, got %+v", metrics)
	}
}

func TestCircuitOpening(t *testing.T) {
	config := Config{
		MaxRequests: 3,
		Timeout:     100 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 3
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()
	op := &mockOperation{shouldFail: true}

	// Execute failing operations to trip the circuit
	for i := 0; i < 3; i++ {
		result, _ := cb.Call(ctx, op.execute)
		if result != nil {
			t.Errorf("Expected nil result for failed operation, got %v", result)
		}
		// The actual error might be the operation error or wrapped
	}

	// Circuit should be open now
	if cb.GetState() != StateOpen {
		t.Errorf("Expected state to be Open after consecutive failures, got %v", cb.GetState())
	}

	// Verify metrics
	metrics := cb.GetMetrics()
	if metrics.ConsecutiveFailures < 3 {
		t.Errorf("Expected at least 3 consecutive failures, got %d", metrics.ConsecutiveFailures)
	}
}

func TestFastFailWhenOpen(t *testing.T) {
	config := Config{
		MaxRequests: 3,
		Timeout:     200 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 2
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()
	op := &mockOperation{shouldFail: true}

	// Trip the circuit
	for i := 0; i < 2; i++ {
		cb.Call(ctx, op.execute)
	}

	// Verify circuit is open
	if cb.GetState() != StateOpen {
		t.Errorf("Expected state to be Open, got %v", cb.GetState())
	}

	op.reset()
	startTime := time.Now()

	// This should fail fast without calling the operation
	result, err := cb.Call(ctx, op.execute)
	elapsed := time.Since(startTime)

	if err == nil {
		t.Error("Expected error when circuit is open")
	}
	if result != nil {
		t.Errorf("Expected nil result when circuit is open, got %v", result)
	}
	if elapsed > 10*time.Millisecond {
		t.Errorf("Expected fast fail, but took %v", elapsed)
	}
	if op.getCallCount() > 0 {
		t.Error("Operation should not be called when circuit is open")
	}
}

func TestHalfOpenTransition(t *testing.T) {
	config := Config{
		MaxRequests: 2,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 2
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()
	op := &mockOperation{shouldFail: true}

	// Trip the circuit
	for i := 0; i < 2; i++ {
		cb.Call(ctx, op.execute)
	}

	// Verify circuit is open
	if cb.GetState() != StateOpen {
		t.Errorf("Expected state to be Open, got %v", cb.GetState())
	}

	// Wait for timeout to elapse
	time.Sleep(60 * time.Millisecond)

	// Next call should transition to half-open
	op.shouldFail = false // Make operation succeed
	result, err := cb.Call(ctx, op.execute)

	// The state should be either half-open (during execution) or closed (after success)
	state := cb.GetState()
	if state != StateHalfOpen && state != StateClosed {
		t.Errorf("Expected state to be HalfOpen or Closed after timeout, got %v", state)
	}

	if err != nil {
		t.Errorf("Expected successful operation in half-open state, got error: %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success' result, got %v", result)
	}
}

func TestHalfOpenToClosedTransition(t *testing.T) {
	config := Config{
		MaxRequests: 2,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 2
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()
	op := &mockOperation{shouldFail: true}

	// Trip the circuit
	for i := 0; i < 2; i++ {
		cb.Call(ctx, op.execute)
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Make operation succeed to close the circuit
	op.shouldFail = false
	cb.Call(ctx, op.execute)

	// Circuit should be closed after successful operation in half-open state
	if cb.GetState() != StateClosed {
		t.Errorf("Expected state to be Closed after successful half-open operation, got %v", cb.GetState())
	}
}

func TestHalfOpenToOpenTransition(t *testing.T) {
	config := Config{
		MaxRequests: 2,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 2
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()
	op := &mockOperation{shouldFail: true}

	// Trip the circuit
	for i := 0; i < 2; i++ {
		cb.Call(ctx, op.execute)
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Keep operation failing - should go back to open
	cb.Call(ctx, op.execute)

	// Circuit should be open again after failed operation in half-open state
	if cb.GetState() != StateOpen {
		t.Errorf("Expected state to be Open after failed half-open operation, got %v", cb.GetState())
	}
}

func TestMaxRequestsInHalfOpen(t *testing.T) {
	config := Config{
		MaxRequests: 2,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 2
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()
	op := &mockOperation{shouldFail: true, delay: 100 * time.Millisecond}

	// Trip the circuit
	for i := 0; i < 2; i++ {
		cb.Call(ctx, op.execute)
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Start multiple concurrent requests in half-open state
	var wg sync.WaitGroup
	results := make([]error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := cb.Call(ctx, op.execute)
			results[index] = err
		}(i)
	}

	wg.Wait()

	// Some requests should be rejected due to MaxRequests limit
	rejectedCount := 0
	for _, err := range results {
		if err != nil && err.Error() == ErrTooManyRequests.Error() {
			rejectedCount++
		}
	}

	if rejectedCount == 0 {
		t.Error("Expected some requests to be rejected due to MaxRequests limit in half-open state")
	}
}

func TestConcurrentAccess(t *testing.T) {
	config := Config{
		MaxRequests: 10,
		Timeout:     100 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 5
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()

	var wg sync.WaitGroup
	numGoroutines := 100
	numCallsPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(shouldFail bool) {
			defer wg.Done()
			op := &mockOperation{shouldFail: shouldFail}
			for j := 0; j < numCallsPerGoroutine; j++ {
				cb.Call(ctx, op.execute)
			}
		}(i%2 == 0) // Half succeed, half fail
	}

	wg.Wait()

	// Verify metrics consistency
	metrics := cb.GetMetrics()
	if metrics.Requests != metrics.Successes+metrics.Failures {
		t.Errorf("Inconsistent metrics: requests=%d, successes=%d, failures=%d",
			metrics.Requests, metrics.Successes, metrics.Failures)
	}
}

func TestStateChangeCallback(t *testing.T) {
	var stateChanges []string
	config := Config{
		MaxRequests: 2,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 2
		},
		OnStateChange: func(name string, from State, to State) {
			stateChanges = append(stateChanges, fmt.Sprintf("%s->%s", from, to))
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()
	op := &mockOperation{shouldFail: true}

	// Trip the circuit (should trigger Closed->Open)
	for i := 0; i < 2; i++ {
		cb.Call(ctx, op.execute)
	}

	// Wait for timeout and make successful call (should trigger Open->HalfOpen->Closed)
	time.Sleep(60 * time.Millisecond)
	op.shouldFail = false
	cb.Call(ctx, op.execute)

	// Verify state change callbacks were called
	if len(stateChanges) == 0 {
		t.Error("Expected state change callbacks to be called")
	}
}

func TestContextCancellation(t *testing.T) {
	config := Config{
		ReadyToTrip: func(m Metrics) bool {
			return false // Never trip
		},
	}

	cb := NewCircuitBreaker(config)
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context immediately
	cancel()

	op := &mockOperation{shouldFail: false, delay: 100 * time.Millisecond}
	result, err := cb.Call(ctx, op.execute)

	// Should respect context cancellation
	if err == nil {
		t.Error("Expected error due to context cancellation")
	}
	if result != nil {
		t.Errorf("Expected nil result due to context cancellation, got %v", result)
	}
}

func TestMetricsAccuracy(t *testing.T) {
	config := Config{
		ReadyToTrip: func(m Metrics) bool {
			return false // Never trip to test metrics in closed state
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()

	// Execute 5 successful operations
	successOp := &mockOperation{shouldFail: false}
	for i := 0; i < 5; i++ {
		cb.Call(ctx, successOp.execute)
	}

	// Execute 3 failed operations
	failOp := &mockOperation{shouldFail: true}
	for i := 0; i < 3; i++ {
		cb.Call(ctx, failOp.execute)
	}

	metrics := cb.GetMetrics()
	if metrics.Requests != 8 {
		t.Errorf("Expected 8 total requests, got %d", metrics.Requests)
	}
	if metrics.Successes != 5 {
		t.Errorf("Expected 5 successes, got %d", metrics.Successes)
	}
	if metrics.Failures != 3 {
		t.Errorf("Expected 3 failures, got %d", metrics.Failures)
	}
	if metrics.ConsecutiveFailures != 3 {
		t.Errorf("Expected 3 consecutive failures, got %d", metrics.ConsecutiveFailures)
	}
}
