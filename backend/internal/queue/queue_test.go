package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"testing/quick"
	"time"
)

// Property 11: Concurrent Request Handling
// For any set of concurrent generation requests, all requests SHALL complete
// (success or controlled failure) without deadlock, and no request SHALL block indefinitely.
// Validates: Requirements 8.1, 8.3
//
// Feature: final-polish, Property 11: Concurrent Request Handling

// TestAcquireRelease_Property_ConcurrentRequestHandling tests that concurrent requests
// complete without deadlock.
// Property: For any number of concurrent requests, all SHALL eventually complete
// (either by acquiring a slot or timing out) without deadlock.
func TestAcquireRelease_Property_ConcurrentRequestHandling(t *testing.T) {
	property := func(maxConcurrent uint8, numRequests uint8) bool {
		// Ensure reasonable bounds
		mc := int(maxConcurrent%10) + 1 // 1-10 concurrent
		nr := int(numRequests%50) + 1   // 1-50 requests

		q := NewRequestQueue(mc)
		var wg sync.WaitGroup
		var completed atomic.Int64
		var failed atomic.Int64

		// Use a timeout context to prevent indefinite blocking
		timeout := 5 * time.Second

		for i := 0; i < nr; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				if err := q.Acquire(ctx); err != nil {
					failed.Add(1)
					return
				}
				defer q.Release()

				// Simulate some work
				time.Sleep(time.Millisecond)
				completed.Add(1)
			}()
		}

		// Wait for all goroutines with a deadline
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// All requests completed (success or timeout)
			total := completed.Load() + failed.Load()
			return total == int64(nr)
		case <-time.After(timeout + time.Second):
			// Deadlock detected
			return false
		}
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: concurrent requests should complete without deadlock: %v", err)
	}
}

// Property 12: Request Queue Fairness
// For any request that acquires a queue slot, it SHALL eventually release the slot,
// and the queue SHALL process requests in approximate FIFO order.
// Validates: Requirements 8.3
//
// Feature: final-polish, Property 12: Request Queue Fairness

// TestAcquireRelease_Property_SlotRelease tests that acquired slots are always releasable.
// Property: For any request that successfully acquires a slot, calling Release SHALL
// make that slot available for subsequent requests.
func TestAcquireRelease_Property_SlotRelease(t *testing.T) {
	property := func(maxConcurrent uint8) bool {
		mc := int(maxConcurrent%10) + 1 // 1-10 concurrent

		q := NewRequestQueue(mc)

		// Acquire all slots
		for i := 0; i < mc; i++ {
			if err := q.Acquire(context.Background()); err != nil {
				return false
			}
		}

		// Queue should be full
		if !q.IsFull() {
			return false
		}

		// TryAcquire should fail
		if q.TryAcquire() {
			return false
		}

		// Release one slot
		q.Release()

		// Now TryAcquire should succeed
		if !q.TryAcquire() {
			return false
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: released slots should become available: %v", err)
	}
}

// TestAcquireRelease_Property_FIFOApproximate tests approximate FIFO ordering.
// Property: For requests waiting in queue, earlier waiters SHALL generally be served
// before later waiters (approximate FIFO due to Go scheduler).
func TestAcquireRelease_Property_FIFOApproximate(t *testing.T) {
	property := func(numWaiters uint8) bool {
		nw := int(numWaiters%10) + 2 // 2-11 waiters

		q := NewRequestQueue(1) // Single slot to force queuing

		// Acquire the only slot
		if err := q.Acquire(context.Background()); err != nil {
			return false
		}

		var mu sync.Mutex
		order := make([]int, 0, nw)
		var wg sync.WaitGroup

		// Start waiters
		for i := 0; i < nw; i++ {
			wg.Add(1)
			idx := i
			go func() {
				defer wg.Done()
				// Small stagger to establish order
				time.Sleep(time.Duration(idx) * time.Millisecond)

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := q.Acquire(ctx); err != nil {
					return
				}
				defer q.Release()

				mu.Lock()
				order = append(order, idx)
				mu.Unlock()
			}()
		}

		// Give waiters time to queue up
		time.Sleep(time.Duration(nw+5) * time.Millisecond)

		// Release the initial slot
		q.Release()

		// Wait for all to complete
		wg.Wait()

		// Check that we got all waiters
		if len(order) != nw {
			return false
		}

		// Check approximate FIFO: count inversions (out-of-order pairs)
		inversions := 0
		for i := 0; i < len(order)-1; i++ {
			if order[i] > order[i+1] {
				inversions++
			}
		}

		// Allow some inversions due to Go scheduler, but not too many
		// For approximate FIFO, inversions should be less than half the pairs
		maxAllowedInversions := nw / 2
		return inversions <= maxAllowedInversions
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: queue should maintain approximate FIFO order: %v", err)
	}
}

// TestStats_Property_Consistency tests that stats are consistent with operations.
// Property: For any sequence of Acquire/Release operations, Stats SHALL accurately
// reflect the current state of the queue.
func TestStats_Property_Consistency(t *testing.T) {
	property := func(maxConcurrent uint8, acquireCount uint8) bool {
		mc := int(maxConcurrent%10) + 1
		ac := int(acquireCount % uint8(mc+1)) // Don't exceed max

		q := NewRequestQueue(mc)

		// Acquire some slots
		for i := 0; i < ac; i++ {
			if err := q.Acquire(context.Background()); err != nil {
				return false
			}
		}

		stats := q.Stats()

		// Verify stats consistency
		if stats.MaxConcurrent != mc {
			return false
		}
		if stats.Active != ac {
			return false
		}
		if q.Available() != mc-ac {
			return false
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: stats should be consistent with operations: %v", err)
	}
}

// TestContextCancellation_Property tests that context cancellation works correctly.
// Property: For any Acquire call with a cancelled context, the call SHALL return
// immediately with the context error.
func TestContextCancellation_Property(t *testing.T) {
	property := func(maxConcurrent uint8) bool {
		mc := int(maxConcurrent%10) + 1

		q := NewRequestQueue(mc)

		// Fill the queue
		for i := 0; i < mc; i++ {
			_ = q.Acquire(context.Background())
		}

		// Create a cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Acquire should return immediately with context error
		start := time.Now()
		err := q.Acquire(ctx)
		elapsed := time.Since(start)

		// Should return quickly (< 100ms) with context.Canceled
		return err == context.Canceled && elapsed < 100*time.Millisecond
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: cancelled context should return immediately: %v", err)
	}
}
