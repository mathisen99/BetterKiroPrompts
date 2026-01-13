// Package queue provides a semaphore-based concurrency limiter for managing
// concurrent requests to external services like OpenAI.
package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// DefaultMaxConcurrent is the default maximum number of concurrent requests.
	DefaultMaxConcurrent = 5
	// DefaultAcquireTimeout is the default timeout for acquiring a slot.
	DefaultAcquireTimeout = 30 * time.Second
)

// RequestQueue implements a semaphore-based concurrency limiter.
// It ensures that no more than maxConcurrent requests are processed simultaneously.
type RequestQueue struct {
	maxConcurrent int
	semaphore     chan struct{}
	waiting       atomic.Int64
	processed     atomic.Int64
	mu            sync.RWMutex
}

// NewRequestQueue creates a new request queue with the specified maximum concurrency.
func NewRequestQueue(maxConcurrent int) *RequestQueue {
	if maxConcurrent <= 0 {
		maxConcurrent = DefaultMaxConcurrent
	}
	return &RequestQueue{
		maxConcurrent: maxConcurrent,
		semaphore:     make(chan struct{}, maxConcurrent),
	}
}

// Acquire attempts to acquire a slot in the queue.
// It blocks until a slot is available or the context is cancelled.
// Returns nil on success, or the context error if cancelled/timed out.
func (q *RequestQueue) Acquire(ctx context.Context) error {
	q.waiting.Add(1)
	defer q.waiting.Add(-1)

	select {
	case q.semaphore <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// AcquireWithTimeout attempts to acquire a slot with a timeout.
// Returns nil on success, context.DeadlineExceeded on timeout.
func (q *RequestQueue) AcquireWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return q.Acquire(ctx)
}

// Release releases a slot back to the queue.
// Must be called after Acquire returns successfully.
func (q *RequestQueue) Release() {
	select {
	case <-q.semaphore:
		q.processed.Add(1)
	default:
		// This should never happen if Acquire/Release are paired correctly
	}
}

// TryAcquire attempts to acquire a slot without blocking.
// Returns true if a slot was acquired, false otherwise.
func (q *RequestQueue) TryAcquire() bool {
	select {
	case q.semaphore <- struct{}{}:
		return true
	default:
		return false
	}
}

// Stats returns current queue statistics.
type Stats struct {
	MaxConcurrent int   // Maximum concurrent requests allowed
	Active        int   // Currently active requests
	Waiting       int64 // Requests waiting for a slot
	Processed     int64 // Total requests processed
}

// Stats returns the current queue statistics.
func (q *RequestQueue) Stats() Stats {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return Stats{
		MaxConcurrent: q.maxConcurrent,
		Active:        len(q.semaphore),
		Waiting:       q.waiting.Load(),
		Processed:     q.processed.Load(),
	}
}

// Available returns the number of available slots.
func (q *RequestQueue) Available() int {
	return q.maxConcurrent - len(q.semaphore)
}

// IsFull returns true if all slots are currently in use.
func (q *RequestQueue) IsFull() bool {
	return len(q.semaphore) >= q.maxConcurrent
}
