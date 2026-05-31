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

func (m Metrics) Clone() Metrics {
	return Metrics{
		Requests:            m.Requests,
		Successes:           m.Successes,
		Failures:            m.Failures,
		ConsecutiveFailures: m.ConsecutiveFailures,
		LastFailureTime:     m.LastFailureTime,
	}
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
		cb.mutex.Lock()
		if err := cb.canExecute(); err != nil {
			cb.mutex.Unlock()
			return nil, err
		}
		cb.insureRecordHalfOpen()
		cb.mutex.Unlock()

		result, err := operation()

		cb.mutex.Lock()
		defer cb.mutex.Unlock()
		if err != nil {
			cb.recordFailure()
		} else {
			cb.recordSuccess()
		}
		return result, err
	}
}

// GetState returns the current state of the circuit breaker
func (cb *circuitBreakerImpl) GetState() State {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.state
}

// GetMetrics returns the current metrics of the circuit breaker
func (cb *circuitBreakerImpl) GetMetrics() Metrics {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.metrics.Clone()
}

// setState changes the circuit breaker state and triggers callbacks
func (cb *circuitBreakerImpl) setState(newState State) {
	if cb.state == newState {
		return
	}
	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()

	switch newState {
	case StateClosed:
		cb.metrics.ConsecutiveFailures = 0
		cb.metrics.Successes = 0
		cb.metrics.Requests = 0
	case StateHalfOpen:
		cb.halfOpenRequests = 0
	}

	if cb.config.OnStateChange != nil {
		cb.mutex.Unlock()
		cb.config.OnStateChange(cb.name, oldState, newState)
		cb.mutex.Lock()
	}
}

// canExecute determines if a request can be executed in the current state
func (cb *circuitBreakerImpl) canExecute() error {
	switch cb.state {
	case StateClosed:
		return nil
	case StateOpen:
		if cb.reachedTimeout() {
			cb.setState(StateHalfOpen)
			return nil
		}
		return ErrCircuitBreakerOpen
	case StateHalfOpen:
		if cb.reachedMaxRequests() {
			return ErrTooManyRequests
		}
		return nil
	}
	return nil
}

func (cb *circuitBreakerImpl) reachedTimeout() bool {
	return time.Since(cb.lastStateChange) > cb.config.Timeout
}

func (cb *circuitBreakerImpl) reachedMaxRequests() bool {
	return cb.halfOpenRequests >= cb.config.MaxRequests
}

// recordSuccess records a successful operation
func (cb *circuitBreakerImpl) recordSuccess() {
	if cb.state == StateClosed && time.Since(cb.lastStateChange) > cb.config.Interval {
		cb.resetMetricsInternal()
	}

	cb.metrics.recordSuccess()
	if cb.state == StateHalfOpen {
		cb.setState(StateClosed)
	}
}

func (m *Metrics) recordSuccess() {
	m.Successes++
	m.ConsecutiveFailures = 0
	m.Requests++
}

// recordFailure records a failed operation
func (cb *circuitBreakerImpl) recordFailure() {
	if cb.state == StateClosed && time.Since(cb.lastStateChange) > cb.config.Interval {
		cb.resetMetricsInternal()
	}
	cb.metrics.recordFailure()
	if cb.state == StateHalfOpen || cb.shouldTrip() {
		cb.setState(StateOpen)
	}
}

// Helper to keep logic DRY
func (cb *circuitBreakerImpl) resetMetricsInternal() {
	cb.metrics = Metrics{} // Fresh start for the new window
	cb.lastStateChange = time.Now()
}

func (cb *circuitBreakerImpl) insureRecordHalfOpen() {
	if cb.state == StateHalfOpen {
		cb.halfOpenRequests++
	}
}

func (m *Metrics) recordFailure() {
	m.ConsecutiveFailures++
	m.Requests++
	m.Failures++
	m.LastFailureTime = time.Now()
}

// shouldTrip determines if the circuit breaker should trip to open state
func (cb *circuitBreakerImpl) shouldTrip() bool {
	metric := cb.metrics
	return cb.config.ReadyToTrip(metric)
}

// isReady checks if the circuit breaker is ready to transition from open to half-open
func (cb *circuitBreakerImpl) isReady() bool {
	timeout := cb.metrics.LastFailureTime.Add(cb.config.Timeout)
	return timeout.Before(time.Now())
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
