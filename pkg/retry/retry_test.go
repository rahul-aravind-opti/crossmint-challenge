package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	retrygo "github.com/avast/retry-go/v4"
	"github.com/stretchr/testify/require"
)

func TestDoRetriesUntilSuccess(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	err := Do(ctx, func(context.Context) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary failure")
		}
		return nil
	}, Config{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     time.Millisecond,
		Multiplier:   1.0,
	}, nil)

	require.NoError(t, err)
	require.Equal(t, 2, attempts)
}

func TestDoStopsOnUnrecoverableError(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	expectedErr := errors.New("stop now")
	err := Do(ctx, func(context.Context) error {
		attempts++
		return retrygo.Unrecoverable(expectedErr)
	}, Config{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     time.Millisecond,
		Multiplier:   1.0,
	}, func(err error) bool {
		return retrygo.IsRecoverable(err)
	})

	require.Error(t, err)
	require.Equal(t, 1, attempts)
	require.Contains(t, err.Error(), expectedErr.Error())
}
