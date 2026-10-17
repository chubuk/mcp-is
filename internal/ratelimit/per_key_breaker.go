package ratelimit

import (
	"sync"
	"time"
)

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreaker struct {
	failures  int
	threshold int
	state     CircuitState
	timeout   time.Duration
	lastFail  time.Time
	mu        sync.Mutex
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		state:     CircuitClosed,
		timeout:   timeout,
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == CircuitOpen {
		if time.Since(cb.lastFail) > cb.timeout {
			cb.state = CircuitHalfOpen
		} else {
			return false
		}
	}
	return true
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = CircuitClosed
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	cb.lastFail = time.Now()
	if cb.failures >= cb.threshold {
		cb.state = CircuitOpen
	}
}

type PerKeyBreaker struct {
	breakers  map[string]*CircuitBreaker
	threshold int
	timeout   time.Duration
	mu        sync.RWMutex
}

func NewPerKeyBreaker(threshold int, timeout time.Duration) *PerKeyBreaker {
	return &PerKeyBreaker{
		breakers:  make(map[string]*CircuitBreaker),
		threshold: threshold,
		timeout:   timeout,
	}
}

func (pkb *PerKeyBreaker) GetBreaker(key string) *CircuitBreaker {
	pkb.mu.RLock()
	cb, ok := pkb.breakers[key]
	pkb.mu.RUnlock()
	if ok {
		return cb
	}

	pkb.mu.Lock()
	cb = NewCircuitBreaker(pkb.threshold, pkb.timeout)
	pkb.breakers[key] = cb
	pkb.mu.Unlock()
	return cb
}
