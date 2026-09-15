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
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := cb.canExecute(); err != nil {
		return nil, err
	}

	result, err := operation()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err == nil {
		cb.recordSuccess()
	} else {
		cb.recordFailure()
	}
	return result, err
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
// NOTE: Caller MUST hold write lock (mutex.Lock)
func (cb *circuitBreakerImpl) setState(newState State) {
	// Check if state actually changed
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()

	// Reset appropriate metrics based on new state
	if newState == StateClosed {
		cb.metrics = Metrics{}
	}
	if newState == StateHalfOpen {
		cb.halfOpenRequests = 0
	}

	// Call OnStateChange callback if configured
	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.name, oldState, newState)
	}
}

// canExecute determines if a request can be executed in the current state
func (cb *circuitBreakerImpl) canExecute() error {
	cb.mutex.RLock()
	state := cb.state
	lastChange := cb.lastStateChange
	cb.mutex.RUnlock()

	// For StateClosed: always allow
	if state == StateClosed {
		return nil
	}

	// For StateOpen: check if timeout has passed for transition to half-open
	if state == StateOpen {
		if time.Since(lastChange) > cb.config.Timeout {
			// Need to transition to half-open - acquire write lock
			cb.mutex.Lock()
			// Double-check state hasn't changed while we were waiting for lock
			if cb.state == StateOpen {
				cb.setState(StateHalfOpen)
				cb.halfOpenRequests = 1 // Count this request
			} else if cb.state == StateHalfOpen {
				// Another goroutine already transitioned, increment counter
				if cb.halfOpenRequests >= cb.config.MaxRequests {
					cb.mutex.Unlock()
					return ErrTooManyRequests
				}
				cb.halfOpenRequests++
			}
			cb.mutex.Unlock()
			return nil
		}
		return ErrCircuitBreakerOpen
	}

	// For StateHalfOpen: check if we've exceeded MaxRequests
	if state == StateHalfOpen {
		cb.mutex.Lock()
		defer cb.mutex.Unlock()
		
		if cb.halfOpenRequests >= cb.config.MaxRequests {
			return ErrTooManyRequests
		}
		cb.halfOpenRequests++
		return nil
	}

	return nil
}

// recordSuccess records a successful operation
// NOTE: Caller MUST hold write lock (mutex.Lock)
func (cb *circuitBreakerImpl) recordSuccess() {
	// Increment success and request counters
	cb.metrics.Requests++
	cb.metrics.Successes++
	
	// Reset consecutive failures
	cb.metrics.ConsecutiveFailures = 0
	
	// In half-open state, consider transitioning to closed
	if cb.state == StateHalfOpen {
		cb.setState(StateClosed)
	}
}

// recordFailure records a failed operation
// NOTE: Caller MUST hold write lock (mutex.Lock)
func (cb *circuitBreakerImpl) recordFailure() {
	// Increment failure and request counters
	cb.metrics.Requests++
	cb.metrics.Failures++
	
	// Increment consecutive failures
	cb.metrics.ConsecutiveFailures++
	
	// Update last failure time
	cb.metrics.LastFailureTime = time.Now()
	
	// In half-open state, transition back to open immediately
	if cb.state == StateHalfOpen {
		cb.setState(StateOpen)
	} else if cb.config.ReadyToTrip(cb.metrics) {
		// Check if circuit should trip (ReadyToTrip function)
		cb.setState(StateOpen)
	}
}

// shouldTrip determines if the circuit breaker should trip to open state
func (cb *circuitBreakerImpl) shouldTrip() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.config.ReadyToTrip(cb.metrics)
}

// isReady checks if the circuit breaker is ready to transition from open to half-open
func (cb *circuitBreakerImpl) isReady() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return time.Since(cb.lastStateChange) > cb.config.Timeout
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