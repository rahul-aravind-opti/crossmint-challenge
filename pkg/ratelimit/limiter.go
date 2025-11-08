package ratelimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limiter provides rate limiting functionality
type Limiter struct {
	limiter *rate.Limiter
	mu      sync.Mutex
}

// NewLimiter creates a new rate limiter
func NewLimiter(rps float64) *Limiter {
	return &Limiter{
		limiter: rate.NewLimiter(rate.Limit(rps), int(rps)),
	}
}

// Wait blocks until the limiter permits an event to happen
func (l *Limiter) Wait(ctx context.Context) error {
	return l.limiter.Wait(ctx)
}

// Allow reports whether an event may happen now
func (l *Limiter) Allow() bool {
	return l.limiter.Allow()
}

// Reserve returns a reservation that indicates how long the caller must wait before event may happen
func (l *Limiter) Reserve() *rate.Reservation {
	return l.limiter.Reserve()
}

// SetRate updates the rate limit
func (l *Limiter) SetRate(rps float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.limiter.SetLimit(rate.Limit(rps))
	l.limiter.SetBurst(int(rps))
}

// TokenBucketLimiter implements a token bucket algorithm
type TokenBucketLimiter struct {
	capacity    int
	tokens      int
	refillRate  time.Duration
	mu          sync.Mutex
	lastRefill  time.Time
}

// NewTokenBucketLimiter creates a new token bucket limiter
func NewTokenBucketLimiter(capacity int, refillRate time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request can proceed
func (t *TokenBucketLimiter) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.refill()

	if t.tokens > 0 {
		t.tokens--
		return true
	}
	return false
}

// Wait blocks until a token is available
func (t *TokenBucketLimiter) Wait(ctx context.Context) error {
	for {
		if t.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(t.refillRate / time.Duration(t.capacity)):
			// Check again after a short wait
		}
	}
}

// refill adds tokens based on elapsed time
func (t *TokenBucketLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(t.lastRefill)
	tokensToAdd := int(elapsed / t.refillRate)

	if tokensToAdd > 0 {
		t.tokens = min(t.capacity, t.tokens+tokensToAdd)
		t.lastRefill = now
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
