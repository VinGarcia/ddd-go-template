package rest

import (
	"context"
	"math/rand"
	"time"
)

// Retry retries the fn callback using an exponential backoff strategy
// starting from `baseDelay` and performing at most `maxRetries`
//
// A retry is attempted only when the callback returns true.
//
// This implementation is based on the following guide:
//
// https://gitlab.wearelayer.com/V3/openbanking/-/merge_requests/7
//
func Retry(ctx context.Context, baseDelay time.Duration, maxDelay time.Duration, maxRetries int, fn func() bool) {
	for i := 0; i < maxRetries; i, baseDelay = i+1, minDuration(baseDelay*2+randMillis(), maxDelay) {
		shouldRetry := fn()
		if !shouldRetry {
			break
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(baseDelay):
			continue
		}
	}
}

var retryRand = rand.New(rand.NewSource(time.Now().Unix()))

var randMillis = DefaultRandMillis

// Calculates retry random factor based on the following guide:
//
// https://cloud.google.com/iot/docs/how-tos/exponential-backoff
//
func DefaultRandMillis() time.Duration {
	return time.Duration(retryRand.Intn(1000)) * time.Millisecond
}

func minDuration(d1 time.Duration, d2 time.Duration) time.Duration {
	if d1 < d2 {
		return d1
	}
	return d2
}
