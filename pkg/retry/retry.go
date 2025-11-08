package retry

import (
	"context"
	"time"

	retry "github.com/avast/retry-go/v4"
)

// Config holds the retry configuration
type Config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultConfig returns a default retry configuration
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// IsRetryable determines if an error should trigger a retry
type IsRetryable func(error) bool

// DefaultIsRetryable returns true for all non-nil errors
func DefaultIsRetryable(err error) bool {
	return err != nil
}

// Do executes the function with retry logic using retry-go
func Do(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryable) error {
	if isRetryable == nil {
		isRetryable = DefaultIsRetryable
	}

	// Convert our config to retry-go options
	opts := []retry.Option{
		retry.Context(ctx),
		retry.Attempts(uint(config.MaxAttempts)),
		retry.Delay(config.InitialDelay),
		retry.MaxDelay(config.MaxDelay),
		retry.DelayType(func(n uint, err error, retryConfig *retry.Config) time.Duration {
			// Exponential backoff with our multiplier
			delay := config.InitialDelay
			for i := uint(0); i < n; i++ {
				delay = time.Duration(float64(delay) * config.Multiplier)
				if delay > config.MaxDelay {
					delay = config.MaxDelay
				}
			}
			return delay
		}),
		retry.RetryIf(func(err error) bool {
			return isRetryable(err)
		}),
	}

	// Wrap the function to work with retry-go's signature
	retryableFunc := func() error {
		return fn(ctx)
	}

	return retry.Do(retryableFunc, opts...)
}