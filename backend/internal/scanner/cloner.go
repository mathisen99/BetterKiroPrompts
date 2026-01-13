package scanner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Cloner errors.
var (
	ErrCloneFailed     = errors.New("failed to clone repository")
	ErrRepoNotFound    = errors.New("repository not found")
	ErrPrivateRepo     = errors.New("private repository requires authentication")
	ErrRepoTooLarge    = errors.New("repository exceeds maximum size limit")
	ErrCloneTimeout    = errors.New("clone operation timed out")
	ErrInvalidRepoPath = errors.New("invalid repository path")
	ErrCleanupFailed   = errors.New("failed to cleanup repository")
	ErrAuthFailed      = errors.New("authentication failed")
	ErrNetworkError    = errors.New("network error during clone")
)

// Default configuration values.
const (
	DefaultMaxRepoSizeMB = 500             // 500 MB default max repo size
	DefaultCloneTimeout  = 5 * time.Minute // 5 minute default clone timeout
	DefaultTempDirPrefix = "scan-repo-"
)

// Cloner handles repository cloning operations.
type Cloner struct {
	// githubToken is the optional GitHub personal access token for private repos.
	// SECURITY: This value must NEVER be logged or included in error messages.
	githubToken string

	// maxSizeMB is the maximum repository size in megabytes.
	maxSizeMB int64

	// cloneTimeout is the maximum time allowed for a clone operation.
	cloneTimeout time.Duration

	// tempDir is the base directory for cloned repositories.
	tempDir string
}

// ClonerOption is a functional option for configuring a Cloner.
type ClonerOption func(*Cloner)

// WithGitHubToken sets the GitHub token for private repository access.
func WithGitHubToken(token string) ClonerOption {
	return func(c *Cloner) {
		c.githubToken = token
	}
}

// WithMaxSizeMB sets the maximum repository size in megabytes.
func WithMaxSizeMB(sizeMB int64) ClonerOption {
	return func(c *Cloner) {
		c.maxSizeMB = sizeMB
	}
}

// WithCloneTimeout sets the clone operation timeout.
func WithCloneTimeout(timeout time.Duration) ClonerOption {
	return func(c *Cloner) {
		c.cloneTimeout = timeout
	}
}

// WithTempDir sets the base directory for cloned repositories.
func WithTempDir(dir string) ClonerOption {
	return func(c *Cloner) {
		c.tempDir = dir
	}
}

// NewCloner creates a new Cloner with the given options.
func NewCloner(opts ...ClonerOption) *Cloner {
	c := &Cloner{
		maxSizeMB:    DefaultMaxRepoSizeMB,
		cloneTimeout: DefaultCloneTimeout,
		tempDir:      os.TempDir(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// CloneResult contains information about a successful clone operation.
type CloneResult struct {
	// Path is the local filesystem path to the cloned repository.
	Path string

	// Owner is the repository owner/organization.
	Owner string

	// Repo is the repository name.
	Repo string

	// CloneDuration is how long the clone operation took.
	CloneDuration time.Duration
}

// Clone clones a GitHub repository to a temporary directory.
// The caller is responsible for calling Cleanup when done.
func (c *Cloner) Clone(ctx context.Context, repoURL string) (*CloneResult, error) {
	// Validate the URL first
	owner, repo, validationErr := ParseGitHubURL(repoURL)
	if validationErr != nil {
		return nil, fmt.Errorf("%w: %s", ErrCloneFailed, validationErr.Message)
	}

	// Create a temporary directory for the clone
	tempDir, err := os.MkdirTemp(c.tempDir, DefaultTempDirPrefix)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create temp directory", ErrCloneFailed)
	}

	// Build the clone URL (with token if available for private repos)
	cloneURL := c.buildCloneURL(owner, repo)

	// Create context with timeout
	cloneCtx, cancel := context.WithTimeout(ctx, c.cloneTimeout)
	defer cancel()

	startTime := time.Now()

	// Execute git clone with shallow clone (depth=1) for efficiency
	// SECURITY: We use --depth=1 to minimize data transfer and avoid pulling full history
	cmd := exec.CommandContext(cloneCtx, "git", "clone", "--depth=1", "--single-branch", cloneURL, tempDir)

	// SECURITY: Capture stderr but sanitize any token references before logging
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Clean up the temp directory on failure
		_ = os.RemoveAll(tempDir)

		// Parse the error to provide a meaningful message
		// SECURITY: Sanitize output to remove any token references
		sanitizedOutput := c.sanitizeOutput(string(output))
		return nil, c.parseCloneError(cloneCtx, err, sanitizedOutput)
	}

	cloneDuration := time.Since(startTime)

	// Check repository size
	size, err := c.getDirectorySize(tempDir)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, fmt.Errorf("%w: failed to check repository size", ErrCloneFailed)
	}

	maxSizeBytes := c.maxSizeMB * 1024 * 1024
	if size > maxSizeBytes {
		_ = os.RemoveAll(tempDir)
		return nil, fmt.Errorf("%w: repository is %d MB, maximum allowed is %d MB",
			ErrRepoTooLarge, size/(1024*1024), c.maxSizeMB)
	}

	return &CloneResult{
		Path:          tempDir,
		Owner:         owner,
		Repo:          repo,
		CloneDuration: cloneDuration,
	}, nil
}

// Cleanup removes a cloned repository directory.
func (c *Cloner) Cleanup(path string) error {
	if path == "" {
		return fmt.Errorf("%w: empty path", ErrInvalidRepoPath)
	}

	// Validate the path is within our temp directory to prevent accidental deletion
	// of important directories
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%w: invalid path", ErrInvalidRepoPath)
	}

	absTempDir, err := filepath.Abs(c.tempDir)
	if err != nil {
		return fmt.Errorf("%w: invalid temp directory", ErrInvalidRepoPath)
	}

	// Ensure the path is within the temp directory
	if !strings.HasPrefix(absPath, absTempDir) {
		return fmt.Errorf("%w: path is not within temp directory", ErrInvalidRepoPath)
	}

	// Remove the directory
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("%w: %v", ErrCleanupFailed, err)
	}

	return nil
}

// HasToken returns true if a GitHub token is configured.
// SECURITY: This method does NOT expose the token value.
func (c *Cloner) HasToken() bool {
	return c.githubToken != ""
}

// buildCloneURL constructs the clone URL, optionally with authentication.
// SECURITY: The token is embedded in the URL for git clone but never logged.
func (c *Cloner) buildCloneURL(owner, repo string) string {
	if c.githubToken != "" {
		// Use token authentication for private repos
		// Format: https://x-access-token:TOKEN@github.com/owner/repo.git
		return fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git", c.githubToken, owner, repo)
	}
	// Public repo URL
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
}

// sanitizeOutput removes any potential token references from output.
// SECURITY: This ensures tokens are never exposed in logs or error messages.
func (c *Cloner) sanitizeOutput(output string) string {
	if c.githubToken == "" {
		return output
	}

	// Replace any occurrence of the token with [REDACTED]
	sanitized := strings.ReplaceAll(output, c.githubToken, "[REDACTED]")

	// Also redact the x-access-token pattern
	sanitized = strings.ReplaceAll(sanitized, "x-access-token:[REDACTED]", "[REDACTED_AUTH]")

	return sanitized
}

// parseCloneError converts git clone errors into appropriate error types.
func (c *Cloner) parseCloneError(ctx context.Context, _ error, output string) error {
	// Check for context timeout/cancellation
	if ctx.Err() == context.DeadlineExceeded {
		return ErrCloneTimeout
	}
	if ctx.Err() == context.Canceled {
		return fmt.Errorf("%w: operation canceled", ErrCloneFailed)
	}

	outputLower := strings.ToLower(output)

	// Check for common error patterns
	switch {
	case strings.Contains(outputLower, "repository not found"):
		return ErrRepoNotFound
	case strings.Contains(outputLower, "could not read from remote repository"):
		if c.githubToken == "" {
			return ErrPrivateRepo
		}
		return ErrAuthFailed
	case strings.Contains(outputLower, "authentication failed"):
		return ErrAuthFailed
	case strings.Contains(outputLower, "could not resolve host"):
		return ErrNetworkError
	case strings.Contains(outputLower, "unable to access"):
		return ErrNetworkError
	default:
		// Generic clone failure - don't expose raw output
		return fmt.Errorf("%w: git clone failed", ErrCloneFailed)
	}
}

// getDirectorySize calculates the total size of a directory in bytes.
func (c *Cloner) getDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
