package util

import (
	"fmt"
	"math/rand"
	"time"
)

func CallWithRetries(maxRetries int, fn func() ([]byte, error)) ([]byte, error) {
	// Exponential backoff parameters
	initialDelay := 1 * time.Second
	maxDelay := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		// Calculate the current delay using exponential backoff
		delay := calculateDelay(i, initialDelay, maxDelay)

		// Call the callback function
		response, err := fn()
		if err == nil {
			// If the function call is successful, return nil
			return response, nil
		}

		// Print the error and retry after the calculated delay
		fmt.Printf("Attempt %d failed: %s. Retrying in %s\n", i+1, err, delay)
		time.Sleep(delay)
	}

	// If all retries are exhausted, return an error
	return nil, fmt.Errorf("all retries exhausted")
}

func calculateDelay(attempt int, initialDelay, maxDelay time.Duration) time.Duration {
	// Calculate the delay using exponential backoff with jitter
	delay := initialDelay * time.Duration(1<<attempt)
	if delay > maxDelay {
		delay = maxDelay
	}
	// Add jitter to the delay to avoid synchronization
	jitter := time.Duration(rand.Int63n(int64(delay)))
	delay = delay + jitter
	return delay
}
