// Package challenge20 contains the implementation for Challenge 20: Circuit Breaker Pattern
package main

import (
	"context"
	"errors"
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

func (m *Metrics) Reset() {
    m.Requests = 0
    m.Successes = 0
    m.Failures = 0
    m.ConsecutiveFailures = 0
    m.LastFailureTime = time.Time{}
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
    
    reservedHalfOpen, err := cb.canExecute()
    if err != nil {
        return nil, err
    }
    
	defer func() {
    	if reservedHalfOpen {
    		cb.releaseHalfOpenPermit()
    	}
    }()
    
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        val, err := operation()
        if err != nil {
            cb.recordFailure()
        } else {
            cb.recordSuccess()
        }
        reservedHalfOpen = false
        return val, err
    }
}

func (cb *circuitBreakerImpl) releaseHalfOpenPermit() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == StateHalfOpen && cb.halfOpenRequests > 0 {
		cb.halfOpenRequests--
	}
}

func (cb *circuitBreakerImpl) runStateMachine(opSuccess bool) {
    switch cb.state {
	    case StateHalfOpen: {
	        if opSuccess {
	            cb.metrics.Reset()
	            cb.setState(StateClosed)
	        } else {
                cb.setState(StateOpen)
	        }
	    }
	    case StateClosed: {
	        if cb.config.ReadyToTrip(cb.metrics) {
	            cb.metrics.LastFailureTime = time.Now()
	            cb.setState(StateOpen)
	        } else {
	            if time.Since(cb.lastStateChange) > cb.config.Interval {
	                cb.metrics.Reset()
	            }
	        }
	    }
	    default:
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
    if cb.state == newState {
		return
	}
    
    prevState := cb.state
	cb.lastStateChange = time.Now()
	
	if newState != StateHalfOpen {
	    cb.halfOpenRequests = 0
	}
	
	cb.state = newState
	
	if cb.config.OnStateChange != nil {
	    go cb.config.OnStateChange(cb.name, prevState, newState)
	}
}

// canExecute determines if a request can be executed in the current state
func (cb *circuitBreakerImpl) canExecute() (bool, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
    switch cb.state {
        case StateClosed: {
            return false, nil
        }
        case StateHalfOpen: {
            if cb.halfOpenRequests >= cb.config.MaxRequests {
                return false, ErrTooManyRequests
            }
            cb.halfOpenRequests++
            return true, nil
        }
        case StateOpen: {
            if time.Since(cb.metrics.LastFailureTime) <= cb.config.Timeout {
                return false, ErrCircuitBreakerOpen
            }
            cb.halfOpenRequests = 0
            cb.setState(StateHalfOpen)
		    cb.halfOpenRequests++
		    return true, nil
        }
    }
	return false, nil
}

// recordSuccess records a successful operation
func (cb *circuitBreakerImpl) recordSuccess() {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    cb.metrics.Requests++
    cb.metrics.Successes++
    cb.metrics.ConsecutiveFailures = 0
    cb.runStateMachine(true)
}

// recordFailure records a failed operation
func (cb *circuitBreakerImpl) recordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.metrics.Requests++
	cb.metrics.Failures++
    cb.metrics.ConsecutiveFailures++
    cb.metrics.LastFailureTime = time.Now()
    cb.runStateMachine(false)
}

func main() {}