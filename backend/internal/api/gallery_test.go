package api

import (
	"math/rand"
	"regexp"
	"testing"
	"testing/quick"
)

// Feature: ux-improvements, Property 6: IP Addresses Are Hashed
// **Validates: Requirements 5.5**
// For any view or rating record stored in the database, the IP identifier
// SHALL be a SHA-256 hash, not a raw IP address.

// sha256HexPattern matches a valid SHA-256 hex string (64 lowercase hex characters)
var sha256HexPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

// TestProperty6_IPAddressesAreHashed tests that IP addresses are hashed using SHA-256.
// Feature: ux-improvements, Property 6: IP Addresses Are Hashed
// **Validates: Requirements 5.5**
func TestProperty6_IPAddressesAreHashed(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// Generate a random IP address
		ip := generateRandomIP(r)

		// Hash the IP
		hash := hashIP(ip)

		// Property 1: Hash should be exactly 64 hex characters (SHA-256 = 256 bits = 64 hex chars)
		if len(hash) != 64 {
			t.Logf("Hash length should be 64, got %d for IP %s", len(hash), ip)
			return false
		}

		// Property 2: Hash should match SHA-256 hex pattern (lowercase hex)
		if !sha256HexPattern.MatchString(hash) {
			t.Logf("Hash should be lowercase hex, got %s for IP %s", hash, ip)
			return false
		}

		// Property 3: Hash should NOT contain the original IP
		if containsIP(hash, ip) {
			t.Logf("Hash should not contain original IP %s", ip)
			return false
		}

		// Property 4: Same IP should produce same hash (deterministic)
		hash2 := hashIP(ip)
		if hash != hash2 {
			t.Logf("Same IP should produce same hash: %s vs %s", hash, hash2)
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 6 (IP Addresses Are Hashed) failed: %v", err)
	}
}

// TestProperty6_DifferentIPsProduceDifferentHashes tests that different IPs produce different hashes.
func TestProperty6_DifferentIPsProduceDifferentHashes(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// Generate two different IP addresses
		ip1 := generateRandomIP(r)
		ip2 := generateRandomIP(r)

		// If IPs happen to be the same, skip this test case
		if ip1 == ip2 {
			return true
		}

		hash1 := hashIP(ip1)
		hash2 := hashIP(ip2)

		// Different IPs should produce different hashes
		if hash1 == hash2 {
			t.Logf("Different IPs should produce different hashes: %s and %s both hash to %s", ip1, ip2, hash1)
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 6 (Different IPs Produce Different Hashes) failed: %v", err)
	}
}

// TestProperty6_IPv6AddressesAreHashed tests that IPv6 addresses are also properly hashed.
func TestProperty6_IPv6AddressesAreHashed(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// Generate a random IPv6 address
		ip := generateRandomIPv6(r)

		// Hash the IP
		hash := hashIP(ip)

		// Property 1: Hash should be exactly 64 hex characters
		if len(hash) != 64 {
			t.Logf("Hash length should be 64, got %d for IPv6 %s", len(hash), ip)
			return false
		}

		// Property 2: Hash should match SHA-256 hex pattern
		if !sha256HexPattern.MatchString(hash) {
			t.Logf("Hash should be lowercase hex, got %s for IPv6 %s", hash, ip)
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 6 (IPv6 Addresses Are Hashed) failed: %v", err)
	}
}

// generateRandomIP generates a random IPv4 address string.
func generateRandomIP(r *rand.Rand) string {
	return string([]byte{
		byte('0' + r.Intn(10)),
		byte('0' + r.Intn(10)),
		byte('0' + r.Intn(10)),
		'.',
		byte('0' + r.Intn(10)),
		byte('0' + r.Intn(10)),
		byte('0' + r.Intn(10)),
		'.',
		byte('0' + r.Intn(10)),
		byte('0' + r.Intn(10)),
		byte('0' + r.Intn(10)),
		'.',
		byte('0' + r.Intn(10)),
		byte('0' + r.Intn(10)),
		byte('0' + r.Intn(10)),
	})
}

// generateRandomIPv6 generates a random IPv6 address string.
func generateRandomIPv6(r *rand.Rand) string {
	hexChars := "0123456789abcdef"
	result := make([]byte, 39) // xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx

	for i := 0; i < 39; i++ {
		if (i+1)%5 == 0 && i < 38 {
			result[i] = ':'
		} else {
			result[i] = hexChars[r.Intn(16)]
		}
	}

	return string(result)
}

// containsIP checks if the hash contains the original IP (which would be a security issue).
func containsIP(hash, ip string) bool {
	// Check if the raw IP string appears anywhere in the hash
	// This would indicate the IP wasn't properly hashed
	if len(ip) > 0 && len(hash) > 0 {
		// Check for direct substring match
		for i := 0; i <= len(hash)-len(ip); i++ {
			if hash[i:i+len(ip)] == ip {
				return true
			}
		}
	}
	return false
}
