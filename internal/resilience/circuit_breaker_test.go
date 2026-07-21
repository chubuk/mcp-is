package resilience

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_ClosedOnSuccess(t *testing.T) {
	cb := NewCircuitBreaker(3, 1*time.Second)

	err := cb.Execute(func() error { return nil })
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if cb.State() != Closed {
		t.Errorf("expected Closed state, got %v", cb.State())
	}
}

func TestCircuitBreaker_OpenAfterFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 10*time.Second)

	for i := 0; i < 3; i++ {
		cb.Execute(func() error { return errors.New("fail") })
	}

	if cb.State() != Open {
		t.Errorf("expected Open state, got %v", cb.State())
	}
}

func TestCircuitBreaker_RejectsWhenOpen(t *testing.T) {
	cb := NewCircuitBreaker(2, 10*time.Second)

	cb.Execute(func() error { return errors.New("fail") })
	cb.Execute(func() error { return errors.New("fail") })

	err := cb.Execute(func() error { return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)

	cb.Execute(func() error { return errors.New("fail") })
	cb.Execute(func() error { return errors.New("fail") })

	time.Sleep(100 * time.Millisecond)

	err := cb.Execute(func() error { return nil })
	if err != nil {
		t.Errorf("expected nil (half-open allowed), got %v", err)
	}
	if cb.State() != Closed {
		t.Errorf("expected Closed after success, got %v", cb.State())
	}
}

func TestCircuitBreaker_ResetOnSuccess(t *testing.T) {
	cb := NewCircuitBreaker(3, 1*time.Second)

	cb.Execute(func() error { return errors.New("fail") })
	cb.Execute(func() error { return errors.New("fail") })
	cb.Execute(func() error { return nil })

	if cb.State() != Closed {
		t.Errorf("expected Closed after success, got %v", cb.State())
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(2, 10*time.Second)

	cb.Execute(func() error { return errors.New("fail") })
	cb.Execute(func() error { return errors.New("fail") })

	cb.Reset()

	if cb.State() != Closed {
		t.Errorf("expected Closed after reset, got %v", cb.State())
	}
}
