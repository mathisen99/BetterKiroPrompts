package ratelimit

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
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
	log    *slog.Logger
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

// NewLimiterWithLogger creates a new rate limiter with logging support.
func NewLimiterWithLogger(log *slog.Logger) *Limiter {
	return &Limiter{
		store:  make(map[string]*clientState),
		limit:  DefaultLimit,
		window: DefaultWindow,
		now:    time.Now,
		log:    log,
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

// NewLimiterWithConfigAndLogger creates a new rate limiter with custom settings and logging.
func NewLimiterWithConfigAndLogger(limit int, window time.Duration, log *slog.Logger) *Limiter {
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
		log:    log,
	}
}

// NewRatingLimiter creates a rate limiter configured for rating submissions (20/hour).
func NewRatingLimiter() *Limiter {
	return NewLimiterWithConfig(RatingLimit, RatingWindow)
}

// NewRatingLimiterWithLogger creates a rate limiter for rating submissions with logging.
func NewRatingLimiterWithLogger(log *slog.Logger) *Limiter {
	return NewLimiterWithConfigAndLogger(RatingLimit, RatingWindow, log)
}

// Allow checks if a request from the given IP is allowed.
// Returns true if allowed, false if rate limited.
// Also returns the duration until the rate limit resets.
func (l *Limiter) Allow(ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Hash IP for privacy in logs
	ipHash := hashIP(ip)

	now := l.now()
	state, exists := l.store[ip]

	if !exists {
		// First request from this IP
		l.store[ip] = &clientState{
			count:       1,
			windowStart: now,
		}
		if l.log != nil {
			l.log.Debug("rate_limit_allowed",
				slog.String("ip_hash", ipHash),
				slog.Int("remaining", l.limit-1),
			)
		}
		return true, 0
	}

	// Check if window has expired
	windowEnd := state.windowStart.Add(l.window)
	if now.After(windowEnd) {
		// Reset window
		state.count = 1
		state.windowStart = now
		if l.log != nil {
			l.log.Debug("rate_limit_allowed",
				slog.String("ip_hash", ipHash),
				slog.Int("remaining", l.limit-1),
			)
		}
		return true, 0
	}

	// Window still active
	if state.count >= l.limit {
		// Rate limited - return time until reset
		retryAfter := windowEnd.Sub(now)
		if l.log != nil {
			l.log.Warn("rate_limit_denied",
				slog.String("ip_hash", ipHash),
				slog.Int("count", state.count),
				slog.Int("limit", l.limit),
				slog.Duration("retry_after", retryAfter),
			)
		}
		return false, retryAfter
	}

	// Allow request and increment count
	state.count++
	if l.log != nil {
		l.log.Debug("rate_limit_allowed",
			slog.String("ip_hash", ipHash),
			slog.Int("remaining", l.limit-state.count),
		)
	}
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

// hashIP creates a privacy-preserving hash of an IP address for logging.
func hashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:8]) // First 8 bytes (16 hex chars)
}
