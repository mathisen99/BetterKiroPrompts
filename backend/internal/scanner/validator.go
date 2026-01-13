// Package scanner provides security scanning functionality for repositories.
package scanner

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Validation errors for URL validation.
var (
	ErrEmptyURL           = errors.New("repository URL cannot be empty")
	ErrInvalidURLFormat   = errors.New("invalid repository URL format")
	ErrInvalidGitHubURL   = errors.New("URL must be a GitHub repository URL")
	ErrMissingOwner       = errors.New("repository URL missing owner")
	ErrMissingRepo        = errors.New("repository URL missing repository name")
	ErrInvalidOwnerFormat = errors.New("invalid owner format")
	ErrInvalidRepoFormat  = errors.New("invalid repository name format")
)

// ValidationError provides structured information about URL validation failures.
type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Example string `json:"example,omitempty"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// githubURLRegex matches valid GitHub repository URLs.
// Supports:
// - https://github.com/owner/repo
// - https://github.com/owner/repo.git
var githubURLRegex = regexp.MustCompile(`^https://github\.com/([a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?)/([a-zA-Z0-9._-]+?)(?:\.git)?/?$`)

// ownerRegex validates GitHub owner/organization names.
// Rules: alphanumeric, hyphens allowed (not at start/end), max 39 chars.
var ownerRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$`)

// repoRegex validates GitHub repository names.
// Rules: alphanumeric, hyphens, underscores, dots allowed.
var repoRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// ValidateGitHubURL validates a GitHub repository URL and returns structured errors.
// It accepts URLs in the format:
// - https://github.com/owner/repo
// - https://github.com/owner/repo.git
func ValidateGitHubURL(url string) *ValidationError {
	// Check for empty URL
	url = strings.TrimSpace(url)
	if url == "" {
		return &ValidationError{
			Code:    "EMPTY_URL",
			Message: ErrEmptyURL.Error(),
			Field:   "repo_url",
			Example: "https://github.com/owner/repo",
		}
	}

	// Check for basic URL structure
	if !strings.HasPrefix(url, "https://") {
		return &ValidationError{
			Code:    "INVALID_PROTOCOL",
			Message: "repository URL must use HTTPS protocol",
			Field:   "repo_url",
			Example: "https://github.com/owner/repo",
		}
	}

	// Check if it's a GitHub URL
	if !strings.HasPrefix(url, "https://github.com/") {
		return &ValidationError{
			Code:    "NOT_GITHUB",
			Message: ErrInvalidGitHubURL.Error(),
			Field:   "repo_url",
			Example: "https://github.com/owner/repo",
		}
	}

	// Extract owner and repo using regex
	matches := githubURLRegex.FindStringSubmatch(url)
	if matches == nil {
		return &ValidationError{
			Code:    "INVALID_FORMAT",
			Message: "invalid repository URL format. Use: https://github.com/owner/repo",
			Field:   "repo_url",
			Example: "https://github.com/owner/repo",
		}
	}

	owner := matches[1]
	repo := matches[2]

	// Validate owner format
	if !ownerRegex.MatchString(owner) {
		return &ValidationError{
			Code:    "INVALID_OWNER",
			Message: fmt.Sprintf("invalid owner format: %s", owner),
			Field:   "owner",
			Example: "Valid owner names are alphanumeric with optional hyphens",
		}
	}

	// Validate repo format
	if !repoRegex.MatchString(repo) {
		return &ValidationError{
			Code:    "INVALID_REPO",
			Message: fmt.Sprintf("invalid repository name format: %s", repo),
			Field:   "repo",
			Example: "Valid repo names are alphanumeric with optional hyphens, underscores, or dots",
		}
	}

	return nil
}

// ParseGitHubURL extracts owner and repo from a validated GitHub URL.
// Returns owner, repo, and any validation error.
func ParseGitHubURL(url string) (owner, repo string, err *ValidationError) {
	if validationErr := ValidateGitHubURL(url); validationErr != nil {
		return "", "", validationErr
	}

	url = strings.TrimSpace(url)
	matches := githubURLRegex.FindStringSubmatch(url)
	if matches == nil {
		return "", "", &ValidationError{
			Code:    "PARSE_ERROR",
			Message: "failed to parse repository URL",
			Field:   "repo_url",
		}
	}

	return matches[1], matches[2], nil
}

// NormalizeGitHubURL normalizes a GitHub URL to a consistent format.
// Removes trailing .git suffix and trailing slashes.
func NormalizeGitHubURL(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSuffix(url, ".git")
	return url
}

// IsValidGitHubURL returns true if the URL is a valid GitHub repository URL.
func IsValidGitHubURL(url string) bool {
	return ValidateGitHubURL(url) == nil
}
