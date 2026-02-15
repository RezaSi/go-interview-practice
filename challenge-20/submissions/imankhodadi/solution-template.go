package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func validateConfig(config Config) error {
	if config.MaxRequests == 0 {
		return errors.New("MaxRequests must be greater than 0")
	}
	if config.Timeout <= 0 {
		return errors.New("Timeout must be greater than 0")
	}
	if config.ReadyToTrip == nil {
		return errors.New("ReadyToTrip function is required")
	}
	return nil
}

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

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

type Metrics struct {
	Requests            int64
	Successes           int64
	Failures            int64
	ConsecutiveFailures int64
	LastFailureTime     time.Time
}

type Config struct {
	MaxRequests   uint32                                  // Max requests allowed in half-open state
	Interval      time.Duration                           // Statistical window for closed state
	Timeout       time.Duration                           // Time to wait before half-open
	ReadyToTrip   func(Metrics) bool                      // Function to determine when to trip
	OnStateChange func(name string, from State, to State) // State change callback
}

type CircuitBreaker interface { // defines the operations for a circuit breaker
	Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)
	GetState() State
	GetMetrics() Metrics
}

type circuitBreaker struct {
	name             string
	config           Config
	state            State
	metrics          Metrics
	mutex            sync.RWMutex
	requests         int64
	lastStateChange  time.Time
	halfOpenRequests uint32
}

var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	ErrTooManyRequests    = errors.New("too many requests in half-open state")
)

func NewCircuitBreaker(config Config) CircuitBreaker {
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
	return &circuitBreaker{
		name:            "circuit-breaker",
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// executes the given operation through the circuit breaker
func (cb *circuitBreaker) Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	// 1. Check current state and handle accordingly
	// 2. For StateClosed: execute operation and track metrics
	// 3. For StateOpen: check if timeout has passed, transition to half-open or fail fast
	// 4. For StateHalfOpen: limit concurrent requests and handle state transitions
	// 5. Update metrics and state based on operation result
	state, err := cb.checkState()
	if err != nil {
		return nil, err
	}
	switch state {
	case StateClosed:
		return cb.callClosed(ctx, operation)
	case StateHalfOpen:
		return cb.callHalfOpen(ctx, operation)
	case StateOpen:
		return nil, ErrCircuitBreakerOpen
	default:
		return nil, errors.New("unknown circuit breaker state")
	}
}

func (cb *circuitBreaker) checkState() (State, error) {
	cb.mutex.RLock()
	state := cb.state
	lastStateChange := cb.lastStateChange
	cb.mutex.RUnlock()
	// If open, check if timeout has passed
	if state == StateOpen {
		if time.Since(lastStateChange) >= cb.config.Timeout {
			cb.mutex.Lock()
			// Double-check after acquiring write lock
			if cb.state == StateOpen && time.Since(cb.lastStateChange) >= cb.config.Timeout {
				cb.setState(StateHalfOpen)
				state = StateHalfOpen
			} else {
				state = cb.state
			}
			cb.mutex.Unlock()
		}
	}
	return state, nil
}

func (cb *circuitBreaker) callHalfOpen(ctx context.Context, operation func() (any, error)) (any, error) {
	select {
	case <-ctx.Done():
		// Context was canceled or deadline exceeded
		return nil, ctx.Err()
	default:

		cb.mutex.Lock()
		// Check if we've exceeded max requests in half-open
		if cb.requests >= int64(cb.config.MaxRequests) {
			cb.mutex.Unlock()
			return nil, ErrTooManyRequests
		}
		cb.requests++
		cb.mutex.Unlock()
		// Execute operation
		result, err := operation()
		cb.mutex.Lock()
		defer cb.mutex.Unlock()
		if err != nil {
			// Failed in half-open, go back to open
			cb.lastStateChange = time.Now()
			cb.setState(StateOpen)
		} else {
			// Success in half-open, go to closed
			cb.setState(StateClosed)
		}
		return result, err
	}
}
func (cb *circuitBreaker) callClosed(ctx context.Context, operation func() (any, error)) (any, error) {
	select {
	case <-ctx.Done():
		// Context was canceled or deadline exceeded
		return nil, ctx.Err()
	default:
		result, err := operation()
		cb.mutex.Lock()
		defer cb.mutex.Unlock()
		cb.metrics.Requests++
		if err != nil {
			cb.metrics.Failures++
			cb.metrics.ConsecutiveFailures++
			cb.metrics.LastFailureTime = time.Now()
			// Check if we should trip to open
			if cb.config.ReadyToTrip(cb.metrics) {
				cb.lastStateChange = time.Now()
				cb.setState(StateOpen)
			}
		} else {
			cb.metrics.Successes++
			cb.metrics.ConsecutiveFailures = 0
		}
		return result, err
	}
}

func (cb *circuitBreaker) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}
func (cb *circuitBreaker) GetMetrics() Metrics {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	// Return a copy to avoid race conditions
	return Metrics{
		Requests:            cb.metrics.Requests,
		Successes:           cb.metrics.Successes,
		Failures:            cb.metrics.Failures,
		ConsecutiveFailures: cb.metrics.ConsecutiveFailures,
		LastFailureTime:     cb.metrics.LastFailureTime,
	}
}

// setState changes the circuit breaker state and triggers callbacks
func (cb *circuitBreaker) setState(newState State) {
	// 1. Check if state actually changed
	// 2. Update lastStateChange time
	// 3. Reset appropriate metrics based on new state
	// 4. Call OnStateChange callback if configured
	// 5. Handle half-open specific logic (reset halfOpenRequests)
	if cb.state == newState {
		return
	}
	oldState := cb.state
	cb.state = newState
	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange("circuit-breaker", oldState, newState)
	}
	// Reset metrics when transitioning to closed
	if newState == StateClosed {
		cb.resetMetrics()
	}
}
func (cb *circuitBreaker) resetMetrics() {
	cb.metrics = Metrics{}
	cb.requests = 0
}

// canExecute determines if a request can be executed in the current state
func (cb *circuitBreaker) canExecute() error {
	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastStateChange) < cb.config.Timeout {
			return ErrCircuitBreakerOpen
		}
	case StateHalfOpen:
		if cb.requests >= int64(cb.config.MaxRequests) {
			return ErrTooManyRequests
		}
	}
	return nil
}

// recordSuccess records a successful operation
func (cb *circuitBreaker) recordSuccess() {
	// 1. Increment success and request counters
	// 2. Reset consecutive failures
	// 3. In half-open state, consider transitioning to closed
	cb.metrics.Successes++
	cb.requests++
	cb.metrics.Requests++
	cb.metrics.Failures = 0
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
	}
}

// recordFailure records a failed operation
func (cb *circuitBreaker) recordFailure() {
	// 1. Increment failure and request counters
	// 2. Increment consecutive failures
	// 3. Update last failure time
	// 4. Check if circuit should trip (ReadyToTrip function)
	// 5. In half-open state, transition back to open

	cb.metrics.Failures++
	cb.metrics.Requests++
	cb.metrics.ConsecutiveFailures++
	cb.metrics.LastFailureTime = time.Now()
	res := cb.config.ReadyToTrip(cb.metrics)
	if cb.state == StateHalfOpen && res {
		cb.state = StateOpen
	}
}

// shouldTrip determines if the circuit breaker should trip to open state
func (cb *circuitBreaker) shouldTrip() bool {
	// Use the ReadyToTrip function from config with current metrics
	return false
}

// isReady checks if the circuit breaker is ready to transition from open to half-open
func (cb *circuitBreaker) isReady() bool {
	// Check if enough time has passed since last state change (Timeout duration)
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
