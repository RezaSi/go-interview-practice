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

type CircuitBreaker interface {
	Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)
	GetState() State
	GetMetrics() Metrics
}

type circuitBreakerImpl struct {
	name             string
	config           Config
	state            State
	metrics          Metrics
	lastStateChange  time.Time
	halfOpenRequests uint32
	mutex            sync.RWMutex
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

	return &circuitBreakerImpl{
		name:            "circuit-breaker",
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

func (cb *circuitBreakerImpl) Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := cb.canExecute(); err != nil {
		return nil, err
	}

	switch cb.GetState() {
	case StateClosed:
		res, err := operation()
		if err != nil {
			cb.recordFailure()
			if cb.shouldTrip() {
				cb.setState(StateOpen)
			}

			return nil, err
		}
		cb.recordSuccess()

		return res, nil
	case StateOpen, StateHalfOpen:
		cb.setState(StateHalfOpen)

		res, err := operation()
		if err != nil {
			cb.recordFailure()
			cb.setState(StateOpen)
			return nil, err
		}

		cb.recordSuccess()
		cb.setState(StateClosed)
		return res, nil
	}
	return nil, errors.New("not implemented")
}

func (cb *circuitBreakerImpl) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

func (cb *circuitBreakerImpl) GetMetrics() Metrics {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.metrics
}

func (cb *circuitBreakerImpl) setState(newState State) {
	state := cb.GetState()
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if state != newState {
		cb.lastStateChange = time.Now()
		cb.state = newState
		if newState == StateHalfOpen {
			cb.halfOpenRequests = 0
		}
		if cb.config.OnStateChange != nil {
			cb.config.OnStateChange(cb.name, state, newState)
		}
	}
}

func (cb *circuitBreakerImpl) canExecute() error {
	state := cb.GetState()

	switch state {
	case StateClosed:
		return nil
	case StateOpen:
		if cb.isReady() {
			return nil
		}
		return ErrCircuitBreakerOpen
	case StateHalfOpen:
		cb.mutex.Lock()
		defer cb.mutex.Unlock()
		cb.halfOpenRequests++

		if cb.halfOpenRequests > cb.config.MaxRequests {
			return ErrTooManyRequests
		}
	}
	return nil
}

func (cb *circuitBreakerImpl) recordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.metrics.Successes++
	cb.metrics.Requests++
	cb.metrics.ConsecutiveFailures = 0
}

func (cb *circuitBreakerImpl) recordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.metrics.Failures++
	cb.metrics.Requests++
	cb.metrics.ConsecutiveFailures++
	cb.metrics.LastFailureTime = time.Now()
}

func (cb *circuitBreakerImpl) shouldTrip() bool {
	return cb.config.ReadyToTrip(cb.metrics)
}

func (cb *circuitBreakerImpl) isReady() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return time.Since(cb.lastStateChange) >= cb.config.Timeout
}

func main() {
	fmt.Println("Circuit Breaker Pattern Example")

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

	ctx := context.Background()

	result, err := cb.Call(ctx, func() (interface{}, error) {
		return "success", nil
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	result, err = cb.Call(ctx, func() (interface{}, error) {
		return nil, errors.New("simulated failure")
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	fmt.Printf("Current state: %v\n", cb.GetState())
	fmt.Printf("Current metrics: %+v\n", cb.GetMetrics())
}
