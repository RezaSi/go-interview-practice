# Hints for Challenge 20: Circuit Breaker Pattern

## Hint 1: Basic Circuit Breaker Structure
Start with the core structure and state management:
```go
type circuitBreaker struct {
    config     Config
    state      State
    metrics    Metrics
    mutex      sync.RWMutex
    lastTrip   time.Time
    requests   int64
}

func NewCircuitBreaker(config Config) CircuitBreaker {
    return &circuitBreaker{
        config: config,
        state:  StateClosed,
    }
}
```

## Hint 2: State Transition Logic
Implement the core state transition methods:
```go
func (cb *circuitBreaker) setState(newState State) {
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
```

## Hint 3: Call Method Implementation Pattern
Structure the main Call method with proper locking:
```go
func (cb *circuitBreaker) Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
    // Check current state
    state, err := cb.checkState()
    if err != nil {
        return nil, err
    }
    
    // Execute operation based on state
    switch state {
    case StateClosed:
        return cb.callClosed(operation)
    case StateHalfOpen:
        return cb.callHalfOpen(operation)
    case StateOpen:
        return nil, ErrCircuitBreakerOpen
    default:
        return nil, errors.New("unknown circuit breaker state")
    }
}
```

## Hint 4: Closed State Handling
Handle requests when circuit is closed:
```go
func (cb *circuitBreaker) callClosed(operation func() (interface{}, error)) (interface{}, error) {
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
            cb.lastTrip = time.Now()
            cb.setState(StateOpen)
        }
    } else {
        cb.metrics.Successes++
        cb.metrics.ConsecutiveFailures = 0
    }
    
    return result, err
}
```

## Hint 5: Half-Open State Management
Handle the testing phase when circuit is half-open:
```go
func (cb *circuitBreaker) callHalfOpen(operation func() (interface{}, error)) (interface{}, error) {
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
        cb.lastTrip = time.Now()
        cb.setState(StateOpen)
    } else {
        // Success in half-open, go to closed
        cb.setState(StateClosed)
    }
    
    return result, err
}
```

## Hint 6: State Checking with Timeout Logic
Implement state checking with automatic transitions:
```go
func (cb *circuitBreaker) checkState() (State, error) {
    cb.mutex.RLock()
    state := cb.state
    lastTrip := cb.lastTrip
    cb.mutex.RUnlock()
    
    // If open, check if timeout has passed
    if state == StateOpen {
        if time.Since(lastTrip) >= cb.config.Timeout {
            cb.mutex.Lock()
            // Double-check after acquiring write lock
            if cb.state == StateOpen && time.Since(cb.lastTrip) >= cb.config.Timeout {
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
```

## Hint 7: Thread-Safe Metrics Access
Implement safe metrics retrieval:
```go
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
```

## Hint 8: Error Definitions and Configuration Validation
Define required errors and validate configuration:
```go
var (
    ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
    ErrTooManyRequests   = errors.New("too many requests in half-open state")
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
```

## Key Circuit Breaker Concepts:
- **State Management**: Careful transitions between Closed/Open/Half-Open
- **Thread Safety**: Use RWMutex for concurrent access
- **Timeout Handling**: Automatic transition from Open to Half-Open
- **Metrics Tracking**: Count requests, failures, and consecutive failures
- **Fast Fail**: Immediate rejection when circuit is open
- **Recovery Testing**: Limited requests in half-open state
- **Configuration**: Flexible failure detection and state change callbacks 