// Package challenge20 contains the implementation for Challenge 20: Circuit Breaker Pattern
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// String returns the string representation of the state
func (s State) String() string {
	switch s {
	case StateClosed:
		return "Closed"
	case StateOpen:
		return "Open"
	case StateHalfOpen:
		return "Half-Open"
	default:
		return "Unknown"
	}
}

// Metrics represents the circuit breaker metrics
type Metrics struct {
	Requests            int64
	Successes           int64
	Failures            int64
	ConsecutiveFailures int64
	LastFailureTime     time.Time
}

// Config represents the configuration for the circuit breaker
// Interval field reserved for future use
type Config struct {
	MaxRequests   uint32                                  // Max requests allowed in half-open state
	Interval      time.Duration                           // Statistical window for closed state
	Timeout       time.Duration                           // Time to wait before half-open
	ReadyToTrip   func(Metrics) bool                      // Function to determine when to trip
	OnStateChange func(name string, from State, to State) // State change callback
}

// CircuitBreaker interface defines the operations for a circuit breaker
type CircuitBreaker interface {
	Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)
	GetState() State
	GetMetrics() Metrics
}

// circuitBreakerImpl is the concrete implementation of CircuitBreaker
type circuitBreakerImpl struct {
	name             string
	config           Config
	state            State
	metrics          Metrics
	lastStateChange  time.Time
	halfOpenRequests uint32
	mutex            sync.RWMutex
}

// Error definitions
var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	ErrTooManyRequests    = errors.New("too many requests in half-open state")
)

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(config Config) CircuitBreaker {
	// Set default values if not provided
	if config.MaxRequests == 0 {
		config.MaxRequests = 1
	}
	if config.Interval == 0 {
		config.Interval = time.Minute
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.ReadyToTrip == nil {
		config.ReadyToTrip = func(m Metrics) bool {
			return m.ConsecutiveFailures >= 5
		}
	}

	return &circuitBreakerImpl{
		name:            "circuit-breaker",
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// Call executes the given operation through the circuit breaker
func (cb *circuitBreakerImpl) Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	// Check ctx cancel
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Check for err from canExecute to determine reject request or not
		if err := cb.canExecute(); err != nil {
			return nil, err
		}
		// Do operation and make a record of the request result
		res, err := operation()
		if err != nil {
			cb.recordFailure()
			return nil, err
		}
		cb.recordSuccess()
		return res, nil
	}
}

// GetState returns the current state of the circuit breaker
func (cb *circuitBreakerImpl) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetMetrics returns the current metrics of the circuit breaker
func (cb *circuitBreakerImpl) GetMetrics() Metrics {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.metrics
}

// setState changes the circuit breaker state and triggers callbacks
func (cb *circuitBreakerImpl) setState(newState State) {
	// Change the State
	cb.mutex.Lock()
	// Check new State not equal old State
	oldState := cb.state
	if oldState == newState {
		cb.mutex.Unlock()
		return
	}
	cb.lastStateChange = time.Now()
	switch newState {
	case StateClosed:
		cb.metrics.Failures = 0
		cb.metrics.Successes = 0
		cb.metrics.Requests = 0
	case StateOpen:
		cb.metrics.Failures = 0
		cb.metrics.Successes = 0
		cb.metrics.Requests = 0
	case StateHalfOpen:
		cb.halfOpenRequests = 0
	}
	cb.state = newState
	cb.mutex.Unlock()

	// If callback set in config - call it
	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.name, oldState, newState)
	}
}

// canExecute determines if a request can be executed in the current state
func (cb *circuitBreakerImpl) canExecute() error {
	// Check the current state and make design to reject request
	var onStateChange func(string, State, State)
	cb.mutex.Lock()

	switch cb.state {
	case StateClosed:
		cb.mutex.Unlock()
		return nil
	case StateOpen:
		if time.Since(cb.metrics.LastFailureTime) > cb.config.Timeout {
			cb.state = StateHalfOpen
			cb.lastStateChange = time.Now()
			cb.halfOpenRequests = 1
			onStateChange = cb.config.OnStateChange
			cb.mutex.Unlock()
			if onStateChange != nil {
				onStateChange(cb.name, StateOpen, StateHalfOpen)
			}
			return nil
		}
		cb.mutex.Unlock()
		return ErrCircuitBreakerOpen
	case StateHalfOpen:
		if cb.halfOpenRequests >= cb.config.MaxRequests {
			cb.mutex.Unlock()
			return ErrTooManyRequests
		}
		cb.halfOpenRequests++
		cb.mutex.Unlock()
		return nil
	}
	cb.mutex.Unlock()
	return nil
}

// recordSuccess records a successful operation
func (cb *circuitBreakerImpl) recordSuccess() {
	// Recording metrics after success request
	cb.mutex.Lock()
	cb.metrics.Successes++
	cb.metrics.Requests++
	cb.metrics.ConsecutiveFailures = 0
	shouldClose := cb.state == StateHalfOpen
	// && uint32(cb.metrics.Successes) >= cb.config.MaxRequests
	cb.mutex.Unlock()
	if shouldClose && cb.GetState() == StateHalfOpen {
		cb.setState(StateClosed)
	}
}

// recordFailure records a failed operation
func (cb *circuitBreakerImpl) recordFailure() {
	// Recording metrics after failure request
	cb.mutex.Lock()
	cb.metrics.Failures++
	cb.metrics.Requests++
	cb.metrics.ConsecutiveFailures++
	cb.metrics.LastFailureTime = time.Now()
	state := cb.state
	cb.mutex.Unlock()
	if state == StateHalfOpen || cb.shouldTrip() {
		cb.setState(StateOpen)
	}
}

// shouldTrip determines if the circuit breaker should trip to open state
func (cb *circuitBreakerImpl) shouldTrip() bool {
	return cb.config.ReadyToTrip(cb.GetMetrics())
}

// isReady checks if the circuit breaker is ready to transition from open to half-open
func (cb *circuitBreakerImpl) isReady() bool {
	if time.Since(cb.GetMetrics().LastFailureTime) > cb.config.Timeout {
		return true
	}
	return false
}

// Example usage and testing helper functions
func main() {
	// Example usage of the circuit breaker
	fmt.Println("Circuit Breaker Pattern Example")

	// Create a circuit breaker configuration
	config := Config{
		MaxRequests: 3,
		Interval:    time.Minute,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from State, to State) {
			fmt.Printf("Circuit breaker %s: %s -> %s\n", name, from, to)
		},
	}

	cb := NewCircuitBreaker(config)

	// Simulate some operations
	ctx := context.Background()

	// Successful operation
	result, err := cb.Call(ctx, func() (interface{}, error) {
		return "success", nil
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	// Failing operation
	result, err = cb.Call(ctx, func() (interface{}, error) {
		return nil, errors.New("simulated failure")
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	// Print current state and metrics
	fmt.Printf("Current state: %v\n", cb.GetState())
	fmt.Printf("Current metrics: %+v\n", cb.GetMetrics())
}
