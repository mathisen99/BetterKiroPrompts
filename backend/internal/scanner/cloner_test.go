package scanner

import (
	"strings"
	"testing"
	"testing/quick"
)

// =============================================================================
// Unit Tests for Cloner
// =============================================================================

func TestNewCloner(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		c := NewCloner()
		if c.maxSizeMB != DefaultMaxRepoSizeMB {
			t.Errorf("maxSizeMB = %d, want %d", c.maxSizeMB, DefaultMaxRepoSizeMB)
		}
		if c.cloneTimeout != DefaultCloneTimeout {
			t.Errorf("cloneTimeout = %v, want %v", c.cloneTimeout, DefaultCloneTimeout)
		}
		if c.githubToken != "" {
			t.Error("githubToken should be empty by default")
		}
	})

	t.Run("with options", func(t *testing.T) {
		c := NewCloner(
			WithGitHubToken("test-token"),
			WithMaxSizeMB(100),
			WithTempDir("/tmp/test"),
		)
		if c.githubToken != "test-token" {
			t.Error("githubToken not set correctly")
		}
		if c.maxSizeMB != 100 {
			t.Errorf("maxSizeMB = %d, want 100", c.maxSizeMB)
		}
		if c.tempDir != "/tmp/test" {
			t.Errorf("tempDir = %s, want /tmp/test", c.tempDir)
		}
	})
}

func TestCloner_HasToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  bool
	}{
		{"no token", "", false},
		{"with token", "ghp_test123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCloner(WithGitHubToken(tt.token))
			if got := c.HasToken(); got != tt.want {
				t.Errorf("HasToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloner_buildCloneURL(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		owner    string
		repo     string
		wantURL  string
		wantAuth bool
	}{
		{
			name:     "public repo",
			token:    "",
			owner:    "owner",
			repo:     "repo",
			wantURL:  "https://github.com/owner/repo.git",
			wantAuth: false,
		},
		{
			name:     "private repo with token",
			token:    "ghp_test123",
			owner:    "owner",
			repo:     "repo",
			wantURL:  "https://x-access-token:ghp_test123@github.com/owner/repo.git",
			wantAuth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCloner(WithGitHubToken(tt.token))
			got := c.buildCloneURL(tt.owner, tt.repo)
			if got != tt.wantURL {
				t.Errorf("buildCloneURL() = %v, want %v", got, tt.wantURL)
			}
		})
	}
}

func TestCloner_sanitizeOutput(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		output string
		want   string
	}{
		{
			name:   "no token configured",
			token:  "",
			output: "some error message",
			want:   "some error message",
		},
		{
			name:   "token in output",
			token:  "ghp_secret123",
			output: "fatal: Authentication failed for 'https://x-access-token:ghp_secret123@github.com/owner/repo.git'",
			want:   "fatal: Authentication failed for 'https://[REDACTED_AUTH]@github.com/owner/repo.git'",
		},
		{
			name:   "token appears multiple times",
			token:  "secret",
			output: "error: secret was found in secret location",
			want:   "error: [REDACTED] was found in [REDACTED] location",
		},
		{
			name:   "no token in output",
			token:  "ghp_secret123",
			output: "Repository not found",
			want:   "Repository not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCloner(WithGitHubToken(tt.token))
			got := c.sanitizeOutput(tt.output)
			if got != tt.want {
				t.Errorf("sanitizeOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloner_Cleanup_InvalidPaths(t *testing.T) {
	c := NewCloner(WithTempDir("/tmp"))

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"empty path", "", true},
		{"path outside temp dir", "/etc/passwd", true},
		{"relative path outside temp", "../../../etc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Cleanup(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cleanup(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

// =============================================================================
// Property-Based Tests for Token Security
// =============================================================================

// TestProperty3_TokenSecurity tests Property 3: Token Security
// Feature: info-and-security-scan, Property 3: Token Security
// **Validates: Requirements 5.3**
//
// Property: For any scan operation using a GitHub token, the token value SHALL NOT appear in:
// - Application logs
// - Error messages returned to users
// - Database records
func TestProperty3_TokenSecurity(t *testing.T) {
	// Sub-property 1: Token is never exposed in sanitized output
	t.Run("token_never_in_sanitized_output", func(t *testing.T) {
		property := func(token, output string) bool {
			// Skip empty tokens
			if token == "" {
				return true
			}

			// Skip if token doesn't appear in output (nothing to sanitize)
			if !strings.Contains(output, token) {
				return true
			}

			c := NewCloner(WithGitHubToken(token))
			sanitized := c.sanitizeOutput(output)

			// Token should NOT appear in sanitized output
			if strings.Contains(sanitized, token) {
				t.Logf("Token leaked in sanitized output: token=%q, output=%q, sanitized=%q",
					token, output, sanitized)
				return false
			}

			return true
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 3 (token never in sanitized output) failed: %v", err)
		}
	})

	// Sub-property 2: HasToken does not expose token value
	t.Run("has_token_does_not_expose_value", func(t *testing.T) {
		property := func(token string) bool {
			c := NewCloner(WithGitHubToken(token))

			// HasToken should only return bool, not the token itself
			hasToken := c.HasToken()

			// Verify the return type is bool and matches expectation
			expectedHasToken := token != ""
			return hasToken == expectedHasToken
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 3 (HasToken does not expose value) failed: %v", err)
		}
	})

	// Sub-property 3: Clone URL with token is properly formatted but token is sanitizable
	t.Run("clone_url_token_sanitizable", func(t *testing.T) {
		property := func(token, owner, repo string) bool {
			// Skip invalid inputs
			if token == "" || owner == "" || repo == "" {
				return true
			}

			// Filter to valid characters
			owner = generateValidOwner(owner)
			repo = generateValidRepo(repo)
			if owner == "" || repo == "" {
				return true
			}

			c := NewCloner(WithGitHubToken(token))
			cloneURL := c.buildCloneURL(owner, repo)

			// The clone URL should contain the token (for git to use)
			if !strings.Contains(cloneURL, token) {
				t.Logf("Token not in clone URL: token=%q, url=%q", token, cloneURL)
				return false
			}

			// But sanitizing the URL should remove the token
			sanitized := c.sanitizeOutput(cloneURL)
			if strings.Contains(sanitized, token) {
				t.Logf("Token leaked after sanitization: token=%q, sanitized=%q", token, sanitized)
				return false
			}

			return true
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 3 (clone URL token sanitizable) failed: %v", err)
		}
	})

	// Sub-property 4: Error messages never contain token
	t.Run("error_messages_never_contain_token", func(t *testing.T) {
		// Test that all error types don't contain tokens
		tokens := []string{
			"ghp_1234567890abcdef",
			"github_pat_test123",
			"gho_secrettoken",
		}

		errorOutputs := []string{
			"fatal: Authentication failed for 'https://x-access-token:%s@github.com/owner/repo.git'",
			"error: could not read from remote repository with token %s",
			"remote: Invalid credentials. Token: %s",
		}

		for _, token := range tokens {
			c := NewCloner(WithGitHubToken(token))

			for _, outputTemplate := range errorOutputs {
				output := strings.ReplaceAll(outputTemplate, "%s", token)
				sanitized := c.sanitizeOutput(output)

				if strings.Contains(sanitized, token) {
					t.Errorf("Token %q leaked in error message: %q", token, sanitized)
				}
			}
		}
	})

	// Sub-property 5: Multiple token occurrences are all sanitized
	t.Run("multiple_token_occurrences_sanitized", func(t *testing.T) {
		property := func(token string) bool {
			if token == "" || len(token) < 3 {
				return true
			}

			c := NewCloner(WithGitHubToken(token))

			// Create output with multiple token occurrences
			output := strings.Repeat(token+" ", 5)
			sanitized := c.sanitizeOutput(output)

			// Count occurrences of token in sanitized output
			count := strings.Count(sanitized, token)
			if count > 0 {
				t.Logf("Token appeared %d times in sanitized output", count)
				return false
			}

			return true
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 3 (multiple token occurrences sanitized) failed: %v", err)
		}
	})
}

// TestProperty3_TokenSecurity_EdgeCases tests edge cases for token security.
// Feature: info-and-security-scan, Property 3: Token Security
// **Validates: Requirements 5.3**
func TestProperty3_TokenSecurity_EdgeCases(t *testing.T) {
	edgeCases := []struct {
		name   string
		token  string
		output string
	}{
		{
			name:   "token at start of output",
			token:  "secret123",
			output: "secret123: error occurred",
		},
		{
			name:   "token at end of output",
			token:  "secret123",
			output: "error with token secret123",
		},
		{
			name:   "token in URL",
			token:  "ghp_abc123",
			output: "https://x-access-token:ghp_abc123@github.com/owner/repo",
		},
		{
			name:   "token with special chars nearby",
			token:  "token123",
			output: "[token123] {token123} (token123)",
		},
		{
			name:   "token in multiline output",
			token:  "mytoken",
			output: "line1\nmytoken\nline3",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCloner(WithGitHubToken(tc.token))
			sanitized := c.sanitizeOutput(tc.output)

			if strings.Contains(sanitized, tc.token) {
				t.Errorf("Token %q leaked in edge case %q: output=%q, sanitized=%q",
					tc.token, tc.name, tc.output, sanitized)
			}
		})
	}
}
