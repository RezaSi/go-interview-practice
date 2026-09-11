// Package challenge20 contains the implementation for Challenge 20: Circuit Breaker Pattern.
// It provides a mechanism to prevent an application from repeatedly trying to execute an operation
// that's likely to fail.
package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	// StateClosed represents the state where requests are allowed to pass through.
	// If failures exceed a threshold, the state transitions to Open.
	StateClosed State = iota

	// StateOpen represents the state where requests are immediately rejected
	// without executing the operation, to allow the downstream system to recover.
	StateOpen

	// StateHalfOpen represents the state where a limited number of "probe" requests
	// are allowed to pass through to check if the system has recovered.
	StateHalfOpen
)

// String returns the string representation of the state.
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

// Metrics represents the circuit breaker runtime metrics used to determine state changes.
type Metrics struct {
	Requests            int64     // Total number of requests attempted
	Successes           int64     // Total number of successful requests
	Failures            int64     // Total number of failed requests
	ConsecutiveFailures int64     // Current count of failures in a row (resets on success)
	LastFailureTime     time.Time // Timestamp of the most recent failure
}

// Config represents the configuration options for the circuit breaker.
type Config struct {
	// MaxRequests is the maximum number of requests allowed to pass through
	// when the circuit is in the Half-Open state.
	MaxRequests uint32

	// Interval is the cyclic period of the closed state to clear the internal Counts.
	// If the interval is 0, the circuit breaker doesn't clear internal counts during the closed state.
	Interval time.Duration

	// Timeout is the period of the open state, after which the state of the circuit breaker becomes half-open.
	Timeout time.Duration

	// ReadyToTrip is called with a copy of Counts whenever a request fails in the closed state.
	// If ReadyToTrip returns true, the circuit breaker will be placed into the open state.
	ReadyToTrip func(Metrics) bool

	// OnStateChange is called whenever the state of the CircuitBreaker changes.
	// This function is executed in a separate goroutine.
	OnStateChange func(name string, from State, to State)
}

// CircuitBreaker interface defines the operations available for a circuit breaker.
type CircuitBreaker interface {
	// Call executes the given operation with circuit breaker protection.
	Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)
	// GetState returns the current state of the circuit breaker.
	GetState() State
	// GetMetrics returns the current internal metrics.
	GetMetrics() Metrics
}

// circuitBreakerImpl is the thread-safe, concrete implementation of CircuitBreaker.
type circuitBreakerImpl struct {
	name             string
	config           Config
	state            State
	metrics          Metrics
	lastStateChange  time.Time
	halfOpenRequests uint32
	mutex            sync.RWMutex
}

// Error definitions for circuit breaker specific scenarios.
var (
	// ErrCircuitBreakerOpen is returned when the circuit breaker is open and refusing requests.
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

	// ErrTooManyRequests is returned when the maximum number of requests in Half-Open state is exceeded.
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// NewCircuitBreaker creates a new circuit breaker with the given configuration.
// If configuration values are missing, defaults are applied:
//   - MaxRequests: 1
//   - Interval: 1 Minute
//   - Timeout: 30 Seconds
//   - ReadyToTrip: > 5 Consecutive Failures
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

// Call executes the given operation through the circuit breaker.
// It checks the circuit state before execution and records metrics after execution.
// It returns an error immediately if the circuit is Open or if the context is cancelled.
func (cb *circuitBreakerImpl) Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	// fail-fast in case of context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := cb.canExecute(); err != nil {
		return nil, err
	}

	// Note: The lock is released inside canExecute() before running operation(),
	// allowing the potentially slow operation to run without blocking other goroutines
	// from checking the state.
	resp, err := operation()

	if err != nil {
		cb.recordFailure()
		return nil, err
	}

	cb.recordSuccess()
	return resp, nil
}

// GetState returns the current state of the circuit breaker in a thread-safe manner.
func (cb *circuitBreakerImpl) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetMetrics returns the current metrics of the circuit breaker in a thread-safe manner.
func (cb *circuitBreakerImpl) GetMetrics() Metrics {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.metrics
}

// setState changes the circuit breaker state and triggers callbacks.
// It expects the caller to hold the mutex.
func (cb *circuitBreakerImpl) setState(newState State) {
	if cb.state == newState {
		return
	}

	if cb.config.OnStateChange != nil {
		// Run callback in a goroutine to prevent blocking the circuit breaker logic
		go cb.config.OnStateChange(cb.name, cb.state, newState)
	}

	cb.lastStateChange = time.Now()
	cb.state = newState

	switch newState {
	case StateClosed:
		// Reset metrics when closing the circuit to start a fresh window
		cb.metrics = Metrics{}
	case StateHalfOpen:
		cb.halfOpenRequests = 0
	}
}

// canExecute determines if a request can be executed in the current state.
// It handles state transitions (e.g., Open -> HalfOpen) based on time.
func (cb *circuitBreakerImpl) canExecute() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Periodic reset of stats in Closed state
	if cb.state == StateClosed && time.Since(cb.lastStateChange) > cb.config.Interval {
		cb.metrics = Metrics{}
		cb.lastStateChange = time.Now()
	}

	// Check if Open state has timed out and should transition to HalfOpen
	if cb.state == StateOpen && cb.isReady() {
		cb.setState(StateHalfOpen)
	}

	switch cb.state {
	case StateHalfOpen:
		if cb.halfOpenRequests >= cb.config.MaxRequests {
			return ErrTooManyRequests
		}
		cb.halfOpenRequests++
	case StateOpen:
		return ErrCircuitBreakerOpen
	}

	return nil
}

// recordSuccess records a successful operation.
// If the circuit was Half-Open, it transitions back to Closed.
func (cb *circuitBreakerImpl) recordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.metrics.Requests++
	cb.metrics.Successes++
	cb.metrics.ConsecutiveFailures = 0

	if cb.state == StateHalfOpen {
		cb.setState(StateClosed)
	}
}

// recordFailure records a failed operation.
// It tracks consecutive failures and trips the circuit to Open if the threshold is met.
func (cb *circuitBreakerImpl) recordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.metrics.Requests++
	cb.metrics.Failures++
	cb.metrics.ConsecutiveFailures++
	cb.metrics.LastFailureTime = time.Now()

	if cb.state == StateHalfOpen {
		cb.setState(StateOpen)
		return
	}

	if cb.state == StateClosed && cb.shouldTrip() {
		cb.setState(StateOpen)
	}
}

// shouldTrip determines if the circuit breaker should trip to open state
// based on the configuration logic.
func (cb *circuitBreakerImpl) shouldTrip() bool {
	return cb.config.ReadyToTrip(cb.metrics)
}

// isReady checks if the circuit breaker is ready to transition from Open to Half-Open
// by checking if the Timeout duration has passed since the last state change.
func (cb *circuitBreakerImpl) isReady() bool {
	return time.Since(cb.lastStateChange) > cb.config.Timeout
}
