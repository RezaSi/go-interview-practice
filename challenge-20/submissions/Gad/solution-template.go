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


func lockReading[T any](mutex *sync.RWMutex, rfunc func() T) T {
	mutex.RLock()
	defer mutex.RUnlock()
	return rfunc()

}

func lockWriting[T any](mutex *sync.RWMutex, wfunc func() T) T {
	mutex.Lock()
	defer mutex.Unlock()
	return wfunc()

}

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
	// 1. Check current state and handle accordingly
	// 2. For StateClosed: execute operation and track metrics
	// 3. For StateOpen: check if timeout has passed, transition to half-open or fail fast
	// 4. For StateHalfOpen: limit concurrent requests and handle state transitions
	// 5. Update metrics and state based on operation result
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	state := cb.GetState()
	switch state {
	case StateClosed:
		if err := cb.canExecute(); err != nil {
			cb.setState(StateOpen)
			return nil, err
		}
		return cb.do(operation)

	case StateOpen:

		if cb.canExecute() == nil {
			cb.setState(StateHalfOpen)
			return cb.do(operation)
		}
		return nil, ErrCircuitBreakerOpen

	case StateHalfOpen:
		err := cb.canExecute()
		if err == nil {
			cb.setState(StateClosed)
			return cb.do(operation)
		}
		return nil, err

	}

	return nil, errors.New("not implemented")
}

func (cb *circuitBreakerImpl) do(operation func() (any, error)) (any, error) {
	res, err := operation()
	if err != nil {
		cb.recordFailure()
		return nil, err
	}

	cb.recordSuccess()
	return res, err
}

// GetState returns the current state of the circuit breaker
func (cb *circuitBreakerImpl) GetState() State {

	return lockReading(&cb.mutex, func() State {
		return cb.state
	})
}

// GetMetrics returns the current metrics of the circuit breaker
func (cb *circuitBreakerImpl) GetMetrics() Metrics {
	return lockReading(&cb.mutex, func() Metrics {
		return cb.metrics
	})
}

// setState changes the circuit breaker state and triggers callbacks
func (cb *circuitBreakerImpl) setState(newState State) {
	// 1. Check if state actually changed
	// 2. Update lastStateChange time
	// 3. Reset appropriate metrics based on new state
	// 4. Call OnStateChange callback if configured
	// 5. Handle half-open specific logic (reset halfOpenRequests)

	currentState := cb.GetState()
	if currentState == newState {
		return
	}
	lockWriting(&cb.mutex, func() struct{} {
		oldState := currentState
		cb.state = newState
		cb.lastStateChange = time.Now()

		switch newState {
		case StateClosed:
			cb.metrics = Metrics{}
		case StateHalfOpen:
			cb.halfOpenRequests = 0
		}

		if cb.config.OnStateChange != nil {
			cb.config.OnStateChange(cb.name, oldState, newState)
		}
		return struct{}{}
	})

}


// canExecute determines if a request can be executed in the current state
func (cb *circuitBreakerImpl) canExecute() error {
	// 1. For StateClosed: always allow
	// 2. For StateOpen: check if timeout has passed for transition to half-open
	// 3. For StateHalfOpen: check if we've exceeded MaxRequests

	currentState := cb.GetState()
	switch currentState {
	case StateClosed:

		return nil
	case StateOpen:

		if cb.isReady() {
			return nil
		}
		return ErrCircuitBreakerOpen

	case StateHalfOpen:
		return lockReading(&cb.mutex, func() error {
			if cb.metrics.Requests >= int64(cb.config.MaxRequests) {
				return ErrTooManyRequests
			}
			return nil
		})

	default:
		return errors.New("unknown state, can't execute")
	}

}

// recordSuccess records a successful operation
func (cb *circuitBreakerImpl) recordSuccess() {
	// 1. Increment success and request counters
	// 2. Reset consecutive failures
	// 3. In half-open state, consider transitioning to closed
	lockWriting(&cb.mutex, func() struct{} {

		cb.metrics.Successes++
		cb.metrics.Requests++
		cb.metrics.ConsecutiveFailures = 0
		return struct{}{}
	})
	if cb.state == StateHalfOpen {
		cb.setState(StateClosed)
	}
}

// recordFailure records a failed operation
func (cb *circuitBreakerImpl) recordFailure() {
	// 1. Increment failure and request counters
	// 2. Increment consecutive failures
	// 3. Update last failure time
	// 4. Check if circuit should trip (ReadyToTrip function)
	// 5. In half-open state, transition back to open

	lockWriting(&cb.mutex, func() struct{} {
		cb.metrics.Failures++
		cb.metrics.Requests++
		cb.metrics.ConsecutiveFailures++
		cb.metrics.LastFailureTime = time.Now()
		return struct{}{}
	})

	// check if we should go from closed to open or from half-open to open
	if cb.shouldTrip() || (cb.state == StateHalfOpen) {
		cb.setState(StateOpen)
	}
}

// shouldTrip determines if the circuit breaker should trip to open state
func (cb *circuitBreakerImpl) shouldTrip() bool {
	// Use the ReadyToTrip function from config with current metrics

	return lockReading(&cb.mutex, func() bool {
		return cb.config.ReadyToTrip(cb.metrics)
	})

}

// isReady checks if the circuit breaker is ready to transition from open to half-open
func (cb *circuitBreakerImpl) isReady() bool {
	// Check if enough time has passed since last state change (Timeout duration)
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
		MaxRequests: 2,
		Interval:    time.Minute,
		Timeout:     100 * time.Millisecond,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 2
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

	result, err = cb.Call(ctx, func() (interface{}, error) {
		return nil, errors.New("simulated failure")
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	fmt.Printf("Result: %v, Error: %v\n", result, err)

	time.Sleep(150 * time.Millisecond)
	result, err = cb.Call(ctx, func() (interface{}, error) {
		return "success", nil
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	// Print current state and metrics
	fmt.Printf("Current state: %v\n", cb.GetState())
	fmt.Printf("Current metrics: %+v\n", cb.GetMetrics())
}
