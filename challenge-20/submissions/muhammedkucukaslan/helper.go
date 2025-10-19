package main

import "time"

func (cb *circuitBreakerImpl) updateLastChange() {
	cb.mutex.Lock()
	cb.lastStateChange = time.Now()
	cb.mutex.Unlock()
}

func (cb *circuitBreakerImpl) resetMetrics() {
	cb.mutex.Lock()
	cb.metrics = Metrics{}
	cb.mutex.Unlock()
}

func (cb *circuitBreakerImpl) updateLastFailureTime() {
	cb.mutex.Lock()
	cb.metrics.LastFailureTime = time.Now()
	cb.mutex.Unlock()
}

func (cb *circuitBreakerImpl) canExecuteInHalfOpen() bool {
	return cb.GetMetrics().Requests < int64(cb.config.MaxRequests)
}

func (cb *circuitBreakerImpl) recordSuccessMetrics() {
	cb.mutex.Lock()
	cb.metrics.Requests++
	cb.metrics.Successes++
	cb.metrics.ConsecutiveFailures = 0
	cb.mutex.Unlock()
}
func (cb *circuitBreakerImpl) recordFailureMetrics() {
	cb.mutex.Lock()
	cb.metrics.Requests++
	cb.metrics.Failures++
	cb.metrics.ConsecutiveFailures++
	cb.mutex.Unlock()
}

func (cb *circuitBreakerImpl) updateState(newState State) {
	cb.mutex.Lock()
	cb.state = newState
	cb.mutex.Unlock()
}

func (cb *circuitBreakerImpl) isSameState(newState State) bool {
	return cb.GetState() == newState
}

func (cb *circuitBreakerImpl) isStateHalfOpen() bool {
	return cb.GetState() == StateHalfOpen
}

func (cb *circuitBreakerImpl) hasOnStateChangeCallback() bool {
	return cb.config.OnStateChange != nil
}

func (cb *circuitBreakerImpl) shouldSwichToOpen() bool {
	return cb.shouldTrip() || cb.GetState() == StateHalfOpen
}

func isStateClosed(state State) bool {
	return state == StateClosed
}

func isStateOpen(state State) bool {
	return state == StateOpen
}
