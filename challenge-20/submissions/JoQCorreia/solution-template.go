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

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 1. Check current state and handle accordingly
	state := cb.GetState()

	switch state {
	case StateClosed:

		cb.mutex.RLock()
		defer cb.mutex.RUnlock()

		res, err := operation()

		if err != nil {
			cb.recordFailure()

			if cb.shouldTrip() {
				cb.setState(StateOpen)
				return nil, ErrCircuitBreakerOpen
			}

			return res, err
		}

		cb.recordSuccess()

		return res, err

	case StateOpen:

		cb.mutex.RLock()
		defer cb.mutex.RUnlock()

		if cb.isReady() {
			cb.setState(StateHalfOpen)
			res, err := operation()

			if err != nil {
				cb.recordFailure()
				return nil, err
			}

			cb.recordSuccess()
			return res, err
		}

		return nil, ErrCircuitBreakerOpen

	case StateHalfOpen:

		cb.mutex.RLock()
		defer cb.mutex.RUnlock()

		cb.halfOpenRequests++

		if cb.canExecute() != nil {
			return nil, cb.canExecute()
		} else {
			if cb.halfOpenRequests >= cb.config.MaxRequests {
				return nil, ErrTooManyRequests
			}

			res, err := operation()
			if err != nil {
				cb.recordFailure()

				if cb.shouldTrip() {
					cb.setState(StateOpen)
					return nil, ErrCircuitBreakerOpen
				}

				return res, err
			}

			cb.recordSuccess()

			return res, err

		}

	}

	return nil, fmt.Errorf("Unknown circuit breaker state: %v", cb.state.String())

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
	// TODO: Implement state transition logic
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	state := newState.String()

	if state != "Open" && state != "Closed" && state != "Half-Open" {
		test := fmt.Errorf("Unknown breaker state %v", state)

		fmt.Print("Value error, state not valid", test)
	}

	oldState := cb.GetState()
	oldStateStr := cb.GetState().String()
	newStateStr := newState.String()

	if oldStateStr == newStateStr {
		return
	}

	cb.lastStateChange = time.Now()

	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.name, oldState, newState)
	}

	cb.state = newState

	switch newStateStr {
	case "Closed":

		cb.metrics = Metrics{}

		return

	case "Half-Open":
		cb.halfOpenRequests = 0

		return
	}

}

// canExecute determines if a request can be executed in the current state
func (cb *circuitBreakerImpl) canExecute() error {

	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	state := cb.GetState().String()

	switch state {
	case "Closed":
		return nil
	case "Open":
		if cb.config.Timeout >= time.Since(cb.lastStateChange) {
			cb.setState(StateHalfOpen)
			return nil
		} else {
			return ErrCircuitBreakerOpen
		}
	case "Half-Open":
		if cb.halfOpenRequests >= cb.config.MaxRequests {
			return ErrTooManyRequests
		}
		return nil
	}

	return fmt.Errorf("Unknown circuit breaker state: %v", state)
}

// recordSuccess records a successful operation
func (cb *circuitBreakerImpl) recordSuccess() {

	// 1. Increment success and request counters
	// 2. Reset consecutive failures
	// 3. In half-open state, consider transitioning to closed

	cb.metrics.Requests++
	cb.metrics.Successes++

	cb.metrics.ConsecutiveFailures = 0

	if cb.GetState() == StateHalfOpen {
		if cb.shouldTrip() == false && cb.metrics.Requests >= int64(cb.config.MaxRequests) {
			cb.setState(StateClosed)
			return
		}
	}

}

// recordFailure records a failed operation
func (cb *circuitBreakerImpl) recordFailure() {

	// 1. Increment failure and request counters
	// 2. Increment consecutive failures
	// 3. Update last failure time
	// 4. Check if circuit should trip (ReadyToTrip function)
	// 5. In half-open state, transition back to open

	cb.metrics.Requests++
	cb.metrics.Failures++

	cb.metrics.ConsecutiveFailures++

	cb.metrics.LastFailureTime = time.Now()

	if cb.GetState() == StateHalfOpen && cb.shouldTrip() {

		cb.setState(StateOpen)

		return
	}
}

// shouldTrip determines if the circuit breaker should trip to open state
func (cb *circuitBreakerImpl) shouldTrip() bool {

	// Use the ReadyToTrip function from config with current metrics

	if !cb.config.ReadyToTrip(cb.metrics) {
		return false
	}

	return true
}

// isReady checks if the circuit breaker is ready to transition from open to half-open
func (cb *circuitBreakerImpl) isReady() bool {

	// Check if enough time has passed since last state change (Timeout duration)
	if time.Since(cb.lastStateChange) >= cb.config.Timeout {
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

	for i := 0; i <= 10; i++ {
		result, err = cb.Call(ctx, func() (interface{}, error) {
			return nil, errors.New("simulated failure")
		})
	}
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	// Print current state and metrics
	fmt.Printf("Current state: %v\n", cb.GetState())
	fmt.Printf("Current metrics: %+v\n", cb.GetMetrics())
}
