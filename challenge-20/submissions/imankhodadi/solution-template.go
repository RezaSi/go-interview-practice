package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

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
	Requests            int64 //total request count
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
	Call(ctx context.Context, operation func() (any, error)) (any, error)
	GetState() State
	GetMetrics() Metrics
}

type circuitBreaker struct {
	name             string
	config           Config
	state            State
	metrics          Metrics
	mutex            sync.RWMutex
	halfOpenRequests int64
	lastStateChange  time.Time
}

var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	ErrTooManyRequests    = errors.New("too many requests in half-open state")
)

func NewCircuitBreaker(config Config) CircuitBreaker {
	if config.MaxRequests == 0 {
		config.MaxRequests = 1
	}
	if config.Interval <= 0 {
		config.Interval = time.Minute
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}
	if config.ReadyToTrip == nil {
		config.ReadyToTrip = func(m Metrics) bool {
			return m.ConsecutiveFailures >= 5
		}
	}
	return &circuitBreaker{
		name:            "circuit-breaker", //TODO: get from function parameter (cannot change function signature foe the assignment)
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// executes the given operation through the circuit breaker
func (cb *circuitBreaker) Call(ctx context.Context, operation func() (any, error)) (any, error) {
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
			var changed bool
			var oldState State
			// Double-check after acquiring write lock
			if cb.state == StateOpen && time.Since(cb.lastStateChange) >= cb.config.Timeout {
				changed, oldState = cb.setState(StateHalfOpen)
				state = StateHalfOpen
			} else {
				state = cb.state
			}
			cb.mutex.Unlock()
			if changed && cb.config.OnStateChange != nil {
				cb.config.OnStateChange(cb.name, oldState, StateHalfOpen)
			}
		}
	}
	return state, nil
}

func (cb *circuitBreaker) callHalfOpen(ctx context.Context, operation func() (any, error)) (any, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	cb.mutex.Lock()
	// Check if we've exceeded max requests in half-open
	if cb.halfOpenRequests >= int64(cb.config.MaxRequests) {
		cb.mutex.Unlock()
		return nil, ErrTooManyRequests
	}
	cb.halfOpenRequests++
	cb.mutex.Unlock()
	// Execute operation
	result, err := operation() //pass context to operation in production

	cb.mutex.Lock()
	// Another probe may have already transitioned the state.
	if cb.state != StateHalfOpen {
		cb.mutex.Unlock()
		return result, err
	}
	cb.metrics.Requests++
	var changed bool
	var oldState State
	if err != nil {
		// Failed in half-open, go back to open
		cb.metrics.Failures++
		cb.metrics.ConsecutiveFailures++
		cb.metrics.LastFailureTime = time.Now()
		changed, oldState = cb.setState(StateOpen)
	} else {
		// Success in half-open
		cb.metrics.Successes++
		// ths assignment requires to transfer to Closed state after one success, but a better solution is
		// transition to Closed after all MaxRequests probes succeed, which is the typical pattern (e.g., Sony gobreaker).
		// change in production
		// resetting all metrics. Subsequent in-flight probes hit the guard and return their result/error to the caller —
		// but their outcome is never recorded in metrics.
		// This means Metrics.Requests/Successes/Failures will undercount under concurrency.
		// This is a known trade-off, metrics are best-effort during HalfOpen-to-Closed transitions,
		// so consumers don't rely on exact counts.
		changed, oldState = cb.setState(StateClosed)
	}
	newState := cb.state
	cb.mutex.Unlock()

	if changed && cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.name, oldState, newState)
	}
	return result, err
}
func (cb *circuitBreaker) callClosed(ctx context.Context, operation func() (any, error)) (any, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	result, err := operation() //pass context to operation in production
	cb.mutex.Lock()
	if cb.state != StateClosed {
		cb.mutex.Unlock()
		return result, err
	}
	// Roll the statistical window if the interval has elapsed
	if cb.lastStateChange.Add(cb.config.Interval).Before(time.Now()) {
		cb.metrics.Requests = 0
		cb.metrics.Successes = 0
		cb.metrics.Failures = 0
		cb.metrics.ConsecutiveFailures = 0
		cb.lastStateChange = time.Now()
	}
	cb.metrics.Requests++
	var changed bool
	var oldState State
	if err != nil {
		cb.metrics.Failures++
		cb.metrics.ConsecutiveFailures++
		cb.metrics.LastFailureTime = time.Now()
		if cb.config.ReadyToTrip(cb.metrics) {
			changed, oldState = cb.setState(StateOpen)
		}
	} else {
		cb.metrics.Successes++
		cb.metrics.ConsecutiveFailures = 0
	}
	newState := cb.state
	cb.mutex.Unlock()
	if changed && cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.name, oldState, newState)
	}
	return result, err
}

func (cb *circuitBreaker) GetState() State {
	cb.checkState() // to update for Open→HalfOpen
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
// setState changes state and returns (changed, oldState) so the caller
// can invoke the callback outside the lock.
func (cb *circuitBreaker) setState(newState State) (bool, State) {
	if cb.state == newState {
		return false, cb.state
	}
	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()

	// Reset metrics when transitioning to closed
	if newState == StateClosed {
		cb.resetMetrics()
	}
	// Reset half-open request counter when entering half-open
	if newState == StateHalfOpen {
		cb.halfOpenRequests = 0
	}
	return true, oldState
}
func (cb *circuitBreaker) resetMetrics() {
	cb.metrics = Metrics{}
	cb.halfOpenRequests = 0
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
	result, err := cb.Call(ctx, func() (any, error) {
		return "success", nil
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	// Failing operation
	result, err = cb.Call(ctx, func() (any, error) {
		return nil, errors.New("simulated failure")
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)
	fmt.Printf("Current state: %v\n", cb.GetState())
	fmt.Printf("Current metrics: %+v\n", cb.GetMetrics())
}