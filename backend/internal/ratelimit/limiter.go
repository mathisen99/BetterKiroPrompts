package ratelimit

import (
	"sync"
	"time"
)

const (
	DefaultLimit  = 10
	DefaultWindow = time.Hour

	// RatingLimit is the rate limit for rating submissions (20 per hour).
	RatingLimit = 20
	// RatingWindow is the time window for rating rate limiting.
	RatingWindow = time.Hour
)

// clientState tracks the request count and window start for a client.
type clientState struct {
	count       int
	windowStart time.Time
}

// Limiter implements an in-memory sliding window rate limiter.
type Limiter struct {
	store  map[string]*clientState
	mu     sync.RWMutex
	limit  int
	window time.Duration
	now    func() time.Time // for testing
}

// NewLimiter creates a new rate limiter with default settings (10 requests per hour).
func NewLimiter() *Limiter {
	return &Limiter{
		store:  make(map[string]*clientState),
		limit:  DefaultLimit,
		window: DefaultWindow,
		now:    time.Now,
	}
}

// NewLimiterWithConfig creates a new rate limiter with custom settings.
func NewLimiterWithConfig(limit int, window time.Duration) *Limiter {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if window <= 0 {
		window = DefaultWindow
	}

	return &Limiter{
		store:  make(map[string]*clientState),
		limit:  limit,
		window: window,
		now:    time.Now,
	}
}

// NewRatingLimiter creates a rate limiter configured for rating submissions (20/hour).
func NewRatingLimiter() *Limiter {
	return NewLimiterWithConfig(RatingLimit, RatingWindow)
}

// Allow checks if a request from the given IP is allowed.
// Returns true if allowed, false if rate limited.
// Also returns the duration until the rate limit resets.
func (l *Limiter) Allow(ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	state, exists := l.store[ip]

	if !exists {
		// First request from this IP
		l.store[ip] = &clientState{
			count:       1,
			windowStart: now,
		}
		return true, 0
	}

	// Check if window has expired
	windowEnd := state.windowStart.Add(l.window)
	if now.After(windowEnd) {
		// Reset window
		state.count = 1
		state.windowStart = now
		return true, 0
	}

	// Window still active
	if state.count >= l.limit {
		// Rate limited - return time until reset
		retryAfter := windowEnd.Sub(now)
		return false, retryAfter
	}

	// Allow request and increment count
	state.count++
	return true, 0
}

// Remaining returns the number of requests remaining for the given IP.
func (l *Limiter) Remaining(ip string) int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	state, exists := l.store[ip]
	if !exists {
		return l.limit
	}

	now := l.now()
	windowEnd := state.windowStart.Add(l.window)
	if now.After(windowEnd) {
		return l.limit
	}

	remaining := l.limit - state.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears the rate limit state for a given IP.
func (l *Limiter) Reset(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.store, ip)
}

// ResetAll clears all rate limit state.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.store = make(map[string]*clientState)
}

// setNow sets a custom time function (for testing).
func (l *Limiter) setNow(fn func() time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.now = fn
}
