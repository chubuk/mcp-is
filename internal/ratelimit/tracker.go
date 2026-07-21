package ratelimit

import (
	"sync"
	"time"
)

type keyState struct {
	count     int
	windowStart time.Time
}

type RateTracker struct {
	states map[string]*keyState
	limit  int
	window time.Duration
	mu     sync.RWMutex
}

func NewRateTracker(limit int, window time.Duration) *RateTracker {
	return &RateTracker{
		states: make(map[string]*keyState),
		limit:  limit,
		window: window,
	}
}

func (rt *RateTracker) RecordRequest(key string) {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	now := time.Now()
	state, ok := rt.states[key]
	if !ok || now.Sub(state.windowStart) > rt.window {
		rt.states[key] = &keyState{count: 1, windowStart: now}
		return
	}
	state.count++
}

func (rt *RateTracker) GetThrottleDelay(key string) time.Duration {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	state, ok := rt.states[key]
	if !ok {
		return 0
	}

	now := time.Now()
	if now.Sub(state.windowStart) > rt.window {
		return 0
	}

	ratio := float64(state.count) / float64(rt.limit)

	switch {
	case ratio >= 0.95:
		return 2 * time.Second
	case ratio >= 0.90:
		return 500 * time.Millisecond
	case ratio >= 0.80:
		return 200 * time.Millisecond
	default:
		return 0
	}
}

func (rt *RateTracker) GetUsage(key string) (int, int, float64) {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	state, ok := rt.states[key]
	if !ok {
		return 0, rt.limit, 0
	}
	return state.count, rt.limit, float64(state.count) / float64(rt.limit)
}
