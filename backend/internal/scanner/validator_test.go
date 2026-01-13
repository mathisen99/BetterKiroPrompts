package scanner

import (
	"fmt"
	"strings"
	"testing"
	"testing/quick"
)

// =============================================================================
// Unit Tests for URL Validation
// =============================================================================

func TestValidateGitHubURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		errCode string
	}{
		// Valid URLs
		{
			name:    "valid basic URL",
			url:     "https://github.com/owner/repo",
			wantErr: false,
		},
		{
			name:    "valid URL with .git suffix",
			url:     "https://github.com/owner/repo.git",
			wantErr: false,
		},
		{
			name:    "valid URL with trailing slash",
			url:     "https://github.com/owner/repo/",
			wantErr: false,
		},
		{
			name:    "valid URL with hyphenated owner",
			url:     "https://github.com/my-org/repo",
			wantErr: false,
		},
		{
			name:    "valid URL with hyphenated repo",
			url:     "https://github.com/owner/my-repo",
			wantErr: false,
		},
		{
			name:    "valid URL with underscored repo",
			url:     "https://github.com/owner/my_repo",
			wantErr: false,
		},
		{
			name:    "valid URL with dotted repo",
			url:     "https://github.com/owner/my.repo",
			wantErr: false,
		},
		{
			name:    "valid URL with numbers",
			url:     "https://github.com/owner123/repo456",
			wantErr: false,
		},
		// Invalid URLs
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
			errCode: "EMPTY_URL",
		},
		{
			name:    "whitespace only",
			url:     "   ",
			wantErr: true,
			errCode: "EMPTY_URL",
		},
		{
			name:    "HTTP instead of HTTPS",
			url:     "http://github.com/owner/repo",
			wantErr: true,
			errCode: "INVALID_PROTOCOL",
		},
		{
			name:    "not GitHub",
			url:     "https://gitlab.com/owner/repo",
			wantErr: true,
			errCode: "NOT_GITHUB",
		},
		{
			name:    "missing repo",
			url:     "https://github.com/owner",
			wantErr: true,
			errCode: "INVALID_FORMAT",
		},
		{
			name:    "missing owner",
			url:     "https://github.com//repo",
			wantErr: true,
			errCode: "INVALID_FORMAT",
		},
		{
			name:    "too many path segments",
			url:     "https://github.com/owner/repo/extra/path",
			wantErr: true,
			errCode: "INVALID_FORMAT",
		},
		{
			name:    "owner starting with hyphen",
			url:     "https://github.com/-owner/repo",
			wantErr: true,
			errCode: "INVALID_FORMAT",
		},
		{
			name:    "owner ending with hyphen",
			url:     "https://github.com/owner-/repo",
			wantErr: true,
			errCode: "INVALID_FORMAT",
		},
		{
			name:    "random string",
			url:     "not a url at all",
			wantErr: true,
			errCode: "INVALID_PROTOCOL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGitHubURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGitHubURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
			if err != nil && tt.errCode != "" && err.Code != tt.errCode {
				t.Errorf("ValidateGitHubURL(%q) error code = %q, want %q", tt.url, err.Code, tt.errCode)
			}
		})
	}
}

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "basic URL",
			url:       "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "URL with .git suffix",
			url:       "https://github.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "complex names",
			url:       "https://github.com/my-org/my-awesome-repo",
			wantOwner: "my-org",
			wantRepo:  "my-awesome-repo",
			wantErr:   false,
		},
		{
			name:    "invalid URL",
			url:     "not a url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ParseGitHubURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGitHubURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if owner != tt.wantOwner {
					t.Errorf("ParseGitHubURL(%q) owner = %q, want %q", tt.url, owner, tt.wantOwner)
				}
				if repo != tt.wantRepo {
					t.Errorf("ParseGitHubURL(%q) repo = %q, want %q", tt.url, repo, tt.wantRepo)
				}
			}
		})
	}
}

func TestNormalizeGitHubURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "already normalized",
			url:  "https://github.com/owner/repo",
			want: "https://github.com/owner/repo",
		},
		{
			name: "with .git suffix",
			url:  "https://github.com/owner/repo.git",
			want: "https://github.com/owner/repo",
		},
		{
			name: "with trailing slash",
			url:  "https://github.com/owner/repo/",
			want: "https://github.com/owner/repo",
		},
		{
			name: "with both",
			url:  "https://github.com/owner/repo.git/",
			want: "https://github.com/owner/repo",
		},
		{
			name: "with whitespace",
			url:  "  https://github.com/owner/repo  ",
			want: "https://github.com/owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeGitHubURL(tt.url)
			if got != tt.want {
				t.Errorf("NormalizeGitHubURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestIsValidGitHubURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://github.com/owner/repo", true},
		{"https://github.com/owner/repo.git", true},
		{"http://github.com/owner/repo", false},
		{"https://gitlab.com/owner/repo", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := IsValidGitHubURL(tt.url)
			if got != tt.want {
				t.Errorf("IsValidGitHubURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

// =============================================================================
// Property-Based Tests for URL Validation
// =============================================================================

// TestProperty1_URLValidation tests Property 1: URL Validation
// Feature: info-and-security-scan, Property 1: URL Validation
// **Validates: Requirements 4.1, 4.4, 4.6**
//
// Property: For any string input to the URL validator:
//   - If it matches a valid GitHub repository URL pattern (https://github.com/owner/repo
//     or https://github.com/owner/repo.git), the validator SHALL accept it
//   - If it does not match a valid pattern, the validator SHALL reject it with an appropriate error
func TestProperty1_URLValidation(t *testing.T) {
	// Sub-property 1: All valid GitHub URLs are accepted
	t.Run("valid_urls_accepted", func(t *testing.T) {
		property := func(owner, repo string) bool {
			// Generate valid owner (alphanumeric, optional hyphens, not at start/end)
			owner = generateValidOwner(owner)
			if owner == "" {
				return true // Skip empty owners
			}

			// Generate valid repo (alphanumeric, hyphens, underscores, dots)
			repo = generateValidRepo(repo)
			if repo == "" {
				return true // Skip empty repos
			}

			// Test basic URL format
			url := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
			err := ValidateGitHubURL(url)
			if err != nil {
				t.Logf("Valid URL rejected: %s, error: %v", url, err)
				return false
			}

			// Test with .git suffix
			urlWithGit := url + ".git"
			err = ValidateGitHubURL(urlWithGit)
			if err != nil {
				t.Logf("Valid URL with .git rejected: %s, error: %v", urlWithGit, err)
				return false
			}

			return true
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 1 (valid URLs accepted) failed: %v", err)
		}
	})

	// Sub-property 2: Invalid URLs are rejected
	t.Run("invalid_urls_rejected", func(t *testing.T) {
		property := func(s string) bool {
			// Skip strings that might accidentally be valid
			if strings.HasPrefix(s, "https://github.com/") {
				return true
			}

			err := ValidateGitHubURL(s)
			// Invalid URLs should be rejected
			return err != nil
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 1 (invalid URLs rejected) failed: %v", err)
		}
	})

	// Sub-property 3: Non-GitHub URLs are rejected
	t.Run("non_github_urls_rejected", func(t *testing.T) {
		nonGitHubDomains := []string{
			"gitlab.com",
			"bitbucket.org",
			"codeberg.org",
			"sr.ht",
			"example.com",
		}

		for _, domain := range nonGitHubDomains {
			url := fmt.Sprintf("https://%s/owner/repo", domain)
			err := ValidateGitHubURL(url)
			if err == nil {
				t.Errorf("Non-GitHub URL should be rejected: %s", url)
			}
		}
	})

	// Sub-property 4: HTTP URLs are rejected (must be HTTPS)
	t.Run("http_urls_rejected", func(t *testing.T) {
		property := func(owner, repo string) bool {
			owner = generateValidOwner(owner)
			repo = generateValidRepo(repo)
			if owner == "" || repo == "" {
				return true
			}

			url := fmt.Sprintf("http://github.com/%s/%s", owner, repo)
			err := ValidateGitHubURL(url)
			return err != nil && err.Code == "INVALID_PROTOCOL"
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 1 (HTTP URLs rejected) failed: %v", err)
		}
	})

	// Sub-property 5: Parse and validate round-trip
	t.Run("parse_validate_roundtrip", func(t *testing.T) {
		property := func(owner, repo string) bool {
			owner = generateValidOwner(owner)
			repo = generateValidRepo(repo)
			if owner == "" || repo == "" {
				return true
			}

			url := fmt.Sprintf("https://github.com/%s/%s", owner, repo)

			// Validate should pass
			if err := ValidateGitHubURL(url); err != nil {
				return false
			}

			// Parse should return the same owner and repo
			parsedOwner, parsedRepo, err := ParseGitHubURL(url)
			if err != nil {
				return false
			}

			return parsedOwner == owner && parsedRepo == repo
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 1 (parse/validate round-trip) failed: %v", err)
		}
	})
}

// generateValidOwner transforms a random string into a valid GitHub owner name.
func generateValidOwner(s string) string {
	// Filter to alphanumeric and hyphens
	var result strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result.WriteRune(c)
		} else if c == '-' && result.Len() > 0 {
			// Only add hyphen if not at start
			result.WriteRune(c)
		}
	}

	owner := result.String()

	// Remove trailing hyphens
	owner = strings.TrimRight(owner, "-")

	// Limit length to 39 chars (GitHub limit)
	if len(owner) > 39 {
		owner = owner[:39]
	}

	// Must start with alphanumeric
	if len(owner) > 0 && owner[0] == '-' {
		owner = owner[1:]
	}

	return owner
}

// generateValidRepo transforms a random string into a valid GitHub repo name.
func generateValidRepo(s string) string {
	// Filter to alphanumeric, hyphens, underscores, and dots
	var result strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' {
			result.WriteRune(c)
		}
	}

	repo := result.String()

	// Limit length to 100 chars (reasonable limit)
	if len(repo) > 100 {
		repo = repo[:100]
	}

	return repo
}

// TestProperty1_URLValidation_EdgeCases tests edge cases for URL validation.
// Feature: info-and-security-scan, Property 1: URL Validation
// **Validates: Requirements 4.1, 4.4, 4.6**
func TestProperty1_URLValidation_EdgeCases(t *testing.T) {
	edgeCases := []struct {
		name    string
		url     string
		isValid bool
	}{
		// Edge cases for owner names
		{"single char owner", "https://github.com/a/repo", true},
		{"max length owner (39 chars)", "https://github.com/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/repo", true},
		{"owner with numbers", "https://github.com/user123/repo", true},
		{"owner all numbers", "https://github.com/123/repo", true},

		// Edge cases for repo names
		{"single char repo", "https://github.com/owner/r", true},
		{"repo with dots", "https://github.com/owner/my.repo.name", true},
		{"repo with underscores", "https://github.com/owner/my_repo_name", true},
		{"repo with mixed chars", "https://github.com/owner/My-Repo_v1.0", true},

		// Edge cases for URL format
		{"double .git suffix", "https://github.com/owner/repo.git.git", true}, // .git.git is valid repo name
		{"trailing slash after .git", "https://github.com/owner/repo.git/", true},

		// Invalid edge cases
		{"empty owner", "https://github.com//repo", false},
		{"empty repo", "https://github.com/owner/", false},
		{"only github.com", "https://github.com/", false},
		{"github.com no path", "https://github.com", false},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateGitHubURL(tc.url)
			isValid := err == nil
			if isValid != tc.isValid {
				t.Errorf("ValidateGitHubURL(%q) valid = %v, want %v (error: %v)", tc.url, isValid, tc.isValid, err)
			}
		})
	}
}
