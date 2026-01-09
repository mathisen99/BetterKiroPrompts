package ratelimit

import (
	"testing"
	"testing/quick"
	"time"
)

// Property 6: Rate Limiting Enforcement
// For any IP address making more than 10 generation requests within one hour,
// the 11th and subsequent requests SHALL receive HTTP 429 response.
// Validates: Requirements 6.3
//
// Feature: ai-driven-generation, Property 6: Rate Limiting Enforcement

// TestAllow_Property_RateLimitEnforcement tests that the 11th request is always denied.
// Property: For any IP address, after exactly 10 allowed requests within the window,
// the next request SHALL be denied.
func TestAllow_Property_RateLimitEnforcement(t *testing.T) {
	property := func(ip string) bool {
		if ip == "" {
			return true // Skip empty IPs
		}

		limiter := NewLimiter()
		// Fix time to ensure we stay within the same window
		fixedTime := time.Now()
		limiter.setNow(func() time.Time { return fixedTime })

		// First 10 requests should be allowed
		for i := 0; i < 10; i++ {
			allowed, _ := limiter.Allow(ip)
			if !allowed {
				return false // Should have been allowed
			}
		}

		// 11th request should be denied
		allowed, retryAfter := limiter.Allow(ip)
		if allowed {
			return false // Should have been denied
		}

		// retryAfter should be positive (time until window resets)
		if retryAfter <= 0 {
			return false
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: rate limit should be enforced after 10 requests: %v", err)
	}
}

// TestAllow_Property_WindowReset tests that rate limit resets after window expires.
// Property: For any IP address that was rate limited, after the window expires,
// the next request SHALL be allowed.
func TestAllow_Property_WindowReset(t *testing.T) {
	property := func(ip string) bool {
		if ip == "" {
			return true // Skip empty IPs
		}

		limiter := NewLimiterWithConfig(10, time.Hour)
		currentTime := time.Now()
		limiter.setNow(func() time.Time { return currentTime })

		// Exhaust the limit
		for i := 0; i < 10; i++ {
			limiter.Allow(ip)
		}

		// Verify rate limited
		allowed, _ := limiter.Allow(ip)
		if allowed {
			return false // Should be rate limited
		}

		// Advance time past the window
		currentTime = currentTime.Add(time.Hour + time.Second)
		limiter.setNow(func() time.Time { return currentTime })

		// Should be allowed again
		allowed, _ = limiter.Allow(ip)
		return allowed
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: rate limit should reset after window expires: %v", err)
	}
}

// TestAllow_Property_IPIsolation tests that rate limits are isolated per IP.
// Property: For any two distinct IP addresses, rate limiting one SHALL NOT affect the other.
func TestAllow_Property_IPIsolation(t *testing.T) {
	property := func(ip1, ip2 string) bool {
		if ip1 == "" || ip2 == "" || ip1 == ip2 {
			return true // Skip invalid or same IPs
		}

		limiter := NewLimiter()
		fixedTime := time.Now()
		limiter.setNow(func() time.Time { return fixedTime })

		// Exhaust limit for ip1
		for i := 0; i < 10; i++ {
			limiter.Allow(ip1)
		}

		// ip1 should be rate limited
		allowed1, _ := limiter.Allow(ip1)
		if allowed1 {
			return false
		}

		// ip2 should still be allowed
		allowed2, _ := limiter.Allow(ip2)
		return allowed2
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: rate limits should be isolated per IP: %v", err)
	}
}

// TestRemaining_Property_DecreasesWithRequests tests that remaining count decreases correctly.
// Property: For any IP address, after N requests (where N <= limit), Remaining SHALL return limit - N.
func TestRemaining_Property_DecreasesWithRequests(t *testing.T) {
	property := func(numRequests uint8) bool {
		// Limit to 10 requests max for this test
		n := int(numRequests % 11)

		limiter := NewLimiter()
		fixedTime := time.Now()
		limiter.setNow(func() time.Time { return fixedTime })

		ip := "test-ip"

		// Make n requests
		for i := 0; i < n; i++ {
			limiter.Allow(ip)
		}

		expected := 10 - n
		actual := limiter.Remaining(ip)
		return actual == expected
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: remaining should decrease with each request: %v", err)
	}
}
