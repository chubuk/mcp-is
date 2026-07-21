package ratelimit

import (
	"testing"
	"time"
)

func TestRateTracker_RecordRequest(t *testing.T) {
	rt := NewRateTracker(100, 5*time.Minute)
	rt.RecordRequest("key1")

	count, limit, ratio := rt.GetUsage("key1")
	if count != 1 {
		t.Errorf("expected 1 request, got %d", count)
	}
	if limit != 100 {
		t.Errorf("expected limit 100, got %d", limit)
	}
	if ratio != 0.01 {
		t.Errorf("expected ratio 0.01, got %f", ratio)
	}
}

func TestRateTracker_ThrottleDelay(t *testing.T) {
	rt := NewRateTracker(100, 5*time.Minute)

	// No throttle at 0%
	if d := rt.GetThrottleDelay("key1"); d != 0 {
		t.Errorf("expected 0 delay at 0%%, got %v", d)
	}

	// Add 80 requests
	for i := 0; i < 80; i++ {
		rt.RecordRequest("key1")
	}

	// 80% = 200ms
	if d := rt.GetThrottleDelay("key1"); d != 200*time.Millisecond {
		t.Errorf("expected 200ms at 80%%, got %v", d)
	}

	// Add 10 more (total 90)
	for i := 0; i < 10; i++ {
		rt.RecordRequest("key1")
	}

	// 90% = 500ms
	if d := rt.GetThrottleDelay("key1"); d != 500*time.Millisecond {
		t.Errorf("expected 500ms at 90%%, got %v", d)
	}

	// Add 5 more (total 95)
	for i := 0; i < 5; i++ {
		rt.RecordRequest("key1")
	}

	// 95% = 2s
	if d := rt.GetThrottleDelay("key1"); d != 2*time.Second {
		t.Errorf("expected 2s at 95%%, got %v", d)
	}
}

func TestRateTracker_WindowReset(t *testing.T) {
	rt := NewRateTracker(100, 50*time.Millisecond)

	for i := 0; i < 100; i++ {
		rt.RecordRequest("key1")
	}

	// At 100% — should throttle
	if d := rt.GetThrottleDelay("key1"); d == 0 {
		t.Error("expected throttle at 100%%")
	}

	// Wait for window to expire
	time.Sleep(100 * time.Millisecond)

	// After window — should not throttle
	if d := rt.GetThrottleDelay("key1"); d != 0 {
		t.Errorf("expected 0 delay after window, got %v", d)
	}
}

func TestRateTracker_PerKeyIsolation(t *testing.T) {
	rt := NewRateTracker(100, 5*time.Minute)

	// Fill key1 to 95%
	for i := 0; i < 95; i++ {
		rt.RecordRequest("key1")
	}

	// key2 should be unaffected
	if d := rt.GetThrottleDelay("key2"); d != 0 {
		t.Errorf("expected 0 delay for key2, got %v", d)
	}
}
