package resilience

import (
	"errors"
	"testing"
)

func TestWithRetry_Success(t *testing.T) {
	attempts := 0
	err := WithRetry(func() error {
		attempts++
		return nil
	}, 3)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestWithRetry_RetriesOnFailure(t *testing.T) {
	attempts := 0
	err := WithRetry(func() error {
		attempts++
		if attempts < 3 {
			return errors.New("fail")
		}
		return nil
	}, 3)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestWithRetry_FailsAfterMaxRetries(t *testing.T) {
	attempts := 0
	err := WithRetry(func() error {
		attempts++
		return errors.New("always fail")
	}, 2)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts (1 initial + 2 retries), got %d", attempts)
	}
}

func TestWithRetry_NoRetries(t *testing.T) {
	attempts := 0
	err := WithRetry(func() error {
		attempts++
		return errors.New("fail")
	}, 0)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestIsRetryable(t *testing.T) {
	if IsRetryable(nil) {
		t.Error("nil should not be retryable")
	}
	if IsRetryable(ErrCircuitOpen) {
		t.Error("ErrCircuitOpen should not be retryable")
	}
	if !IsRetryable(errors.New("some error")) {
		t.Error("regular error should be retryable")
	}
}
