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
	// TODO: Implement the main circuit breaker logic
	// 0. Check canExecute() </
	// 1. Check current state and handle accordingly </
	// 2. For StateClosed: execute operation and track metrics
	// 3. For StateOpen: check if timeout has passed, transition to half-open or fail fast
	// 4. For StateHalfOpen: limit concurrent requests and handle state transitions
	// 5. Update metrics and state based on operation result
	err := cb.canExecute()
	if err != nil {
	    return nil, err
	}
	
	type response struct{
	    r interface{}
	    e error
	}
	resChan := make(chan *response)
	go func(o func() (interface{}, error), rc chan<-*response){
	    res, err := o()
	    rc <- &response{r: res, e: err}
	}(operation, resChan)
	
	if cb.GetState() == StateHalfOpen{
	    cb.mutex.Lock()
	    cb.halfOpenRequests += 1
	    cb.mutex.Unlock()
	}
	
	select{
    case <- ctx.Done():
	    return nil, ctx.Err()
    case res := <- resChan:
        if res.e != nil{
            cb.recordFailure()
        }else{
            cb.recordSuccess()
        }
        return res.r, res.e
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

// GetConfig returns the config of the circuit breaker
func (cb *circuitBreakerImpl) GetConfig() Config {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.config
}

// setState changes the circuit breaker state and triggers callbacks
func (cb *circuitBreakerImpl) setState(newState State) {
	curState := cb.GetState()
	if curState == newState {
	    return
	}
	
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = newState
	cb.lastStateChange = time.Now()
	
	switch newState{
    case StateClosed: //only from HalfOpen
        cb.metrics = Metrics{}
    case StateHalfOpen: //only from Open
        cb.halfOpenRequests = 0
    //case StateOpen: //from both Closed and HalfOpen
	}
	if cb.config.OnStateChange != nil{
	    cb.config.OnStateChange(cb.name, curState, newState)
	}
}

// canExecute determines if a request can be executed in the current state
func (cb *circuitBreakerImpl) canExecute() error {
	curState := cb.GetState()
	
	switch curState {
	case StateClosed:
	    return nil
	    
	case StateHalfOpen:
	    cb.mutex.RLock()
	    defer cb.mutex.RUnlock()
	    if cb.halfOpenRequests >= cb.config.MaxRequests {
	        return ErrTooManyRequests
	    }
	    return nil
	    
	case StateOpen:
	    if !cb.isReady(){
	        return ErrCircuitBreakerOpen
	    }
	    cb.setState(StateHalfOpen)
        return nil
    default:
        return errors.New("unknown state is not allowed to execute")
	}
}

// recordSuccess records a successful operation
func (cb *circuitBreakerImpl) recordSuccess() {
	cb.mutex.Lock()
	cb.metrics.Successes += 1
	cb.metrics.Requests += 1
	cb.metrics.ConsecutiveFailures = 0
	cb.mutex.Unlock()
	
	if cb.GetState() == StateHalfOpen { //transition after 1 success ?
	    cb.setState(StateClosed)
	}
}

// recordFailure records a failed operation
func (cb *circuitBreakerImpl) recordFailure() {
    cb.mutex.Lock()
	cb.metrics.Failures += 1
	cb.metrics.Requests += 1
	cb.metrics.ConsecutiveFailures += 1
	cb.metrics.LastFailureTime = time.Now()
	cb.mutex.Unlock()
	
	if cb.shouldTrip() || cb.GetState() == StateHalfOpen{
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
