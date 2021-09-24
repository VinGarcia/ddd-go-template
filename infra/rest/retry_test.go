package rest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	// Simplify testing by removing randomness:
	randMillis = func() time.Duration {
		return 0
	}
	defer func() {
		randMillis = DefaultRandMillis
	}()

	t.Run("should retry up to max times", func(t *testing.T) {
		var calls []struct{}
		ctx := context.TODO()
		Retry(ctx, 1*time.Millisecond, 1*time.Millisecond, 3, func() bool {
			calls = append(calls, struct{}{})
			return true
		})

		assert.Equal(t, len(calls), 3)
	})

	t.Run("should not retry if not required", func(t *testing.T) {
		var calls []struct{}
		ctx := context.TODO()
		Retry(ctx, 1*time.Millisecond, 1*time.Millisecond, 3, func() bool {
			calls = append(calls, struct{}{})
			return false
		})

		assert.Equal(t, len(calls), 1)
	})

	t.Run("should not retry if the context is cancelled", func(t *testing.T) {
		var calls []struct{}
		ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Millisecond)
		defer cancel()
		Retry(ctx, 10*time.Millisecond, 10*time.Millisecond, 3, func() bool {
			calls = append(calls, struct{}{})
			return true
		})

		assert.Equal(t, len(calls), 1)
	})

	t.Run("should backoff exponentially with a max limit", func(t *testing.T) {
		startTime := time.Now()
		var pauses []time.Duration
		ctx := context.TODO()
		Retry(ctx, 10*time.Millisecond, 24*time.Millisecond, 4, func() bool {
			pauses = append(pauses, time.Since(startTime))
			startTime = time.Now()
			return true
		})

		assert.Equal(t, len(pauses), 4)

		tolerance := 4 * time.Millisecond
		assertApprox(t, tolerance, pauses[0], 0)
		assertApprox(t, tolerance, pauses[1], 10*time.Millisecond)
		assertApprox(t, tolerance, pauses[2], 20*time.Millisecond)
		assertApprox(t, tolerance, pauses[3], 24*time.Millisecond)
	})
}

func assertApprox(t *testing.T, tolerance time.Duration, d1 time.Duration, d2 time.Duration) {
	diff := d1 - d2
	if diff < 0 {
		diff = -diff
	}

	if diff >= tolerance {
		t.Fatalf("the values of d1 and d2 are not close, d1: %v, d2: %v", d1, d2)
	}
}
