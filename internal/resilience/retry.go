package resilience

import (
	"errors"
	"math/rand"
	"time"
)

func WithRetry(fn func() error, maxRetries int) error {
	var err error
	for i := 0; i <= maxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < maxRetries {
			backoff := time.Duration(1<<uint(i)) * 200 * time.Millisecond
			jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
			time.Sleep(backoff + jitter)
		}
	}
	return err
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrCircuitOpen) {
		return false
	}
	return true
}
