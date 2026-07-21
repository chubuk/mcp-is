package ratelimit

import (
	"testing"
	"time"
)

func TestPerKeyBreaker_Allow(t *testing.T) {
	pkb := NewPerKeyBreaker(3, 100*time.Millisecond)

	if !pkb.GetBreaker("key1").Allow() {
		t.Error("expected Allow to return true initially")
	}
}

func TestPerKeyBreaker_OpenAfterFailures(t *testing.T) {
	pkb := NewPerKeyBreaker(2, 100*time.Millisecond)
	cb := pkb.GetBreaker("key1")

	cb.RecordFailure()
	cb.RecordFailure()

	if cb.Allow() {
		t.Error("expected Allow to return false after 2 failures")
	}
}

func TestPerKeyBreaker_PerKeyIsolation(t *testing.T) {
	pkb := NewPerKeyBreaker(2, 100*time.Millisecond)

	// Open key1
	cb1 := pkb.GetBreaker("key1")
	cb1.RecordFailure()
	cb1.RecordFailure()

	// key2 should be unaffected
	cb2 := pkb.GetBreaker("key2")
	if !cb2.Allow() {
		t.Error("expected key2 to be allowed")
	}
}

func TestPerKeyBreaker_Recovery(t *testing.T) {
	pkb := NewPerKeyBreaker(2, 50*time.Millisecond)
	cb := pkb.GetBreaker("key1")

	cb.RecordFailure()
	cb.RecordFailure()

	time.Sleep(100 * time.Millisecond)

	if !cb.Allow() {
		t.Error("expected Allow to return true after timeout")
	}
}
