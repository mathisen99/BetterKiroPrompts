package scanner

import (
	"better-kiro-prompts/internal/openai"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Scan job status values.
const (
	StatusPending   = "pending"
	StatusCloning   = "cloning"
	StatusScanning  = "scanning"
	StatusReviewing = "reviewing"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

// Service errors.
var (
	ErrJobNotFound = errors.New("scan job not found")
	ErrScanFailed  = errors.New("scan failed")
)

// ScanJob represents a security scan job.
type ScanJob struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	RepoURL     string     `json:"repo_url"`
	Languages   []string   `json:"languages"`
	Findings    []Finding  `json:"findings"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// ScanRequest represents a request to start a scan.
type ScanRequest struct {
	RepoURL string `json:"repo_url"`
}

// Service orchestrates security scanning operations.
type Service struct {
	db         *sql.DB
	cloner     *Cloner
	detector   *LanguageDetector
	toolRunner *ToolRunner
	aggregator *Aggregator
	reviewer   *CodeReviewer
}

// ServiceOption is a functional option for configuring a Service.
type ServiceOption func(*Service)

// WithCloner sets the cloner for the service.
func WithServiceCloner(c *Cloner) ServiceOption {
	return func(s *Service) {
		s.cloner = c
	}
}

// WithToolRunner sets the tool runner for the service.
func WithServiceToolRunner(r *ToolRunner) ServiceOption {
	return func(s *Service) {
		s.toolRunner = r
	}
}

// WithCodeReviewer sets the code reviewer for the service.
func WithServiceCodeReviewer(r *CodeReviewer) ServiceOption {
	return func(s *Service) {
		s.reviewer = r
	}
}

// NewService creates a new scanner service.
func NewService(db *sql.DB, openaiClient *openai.Client, githubToken string, opts ...ServiceOption) *Service {
	s := &Service{
		db:         db,
		cloner:     NewCloner(WithGitHubToken(githubToken)),
		detector:   NewLanguageDetector(),
		toolRunner: NewToolRunner(),
		aggregator: NewAggregator(),
		reviewer:   NewCodeReviewer(openaiClient),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// StartScan initiates a new security scan.
func (s *Service) StartScan(ctx context.Context, req ScanRequest) (*ScanJob, error) {
	// Validate URL
	if err := ValidateGitHubURL(req.RepoURL); err != nil {
		return nil, err
	}

	// Create job
	job := &ScanJob{
		ID:        uuid.New().String(),
		Status:    StatusPending,
		RepoURL:   NormalizeGitHubURL(req.RepoURL),
		CreatedAt: time.Now(),
	}

	// Persist job
	if err := s.createJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Start scan in background
	go s.runScan(context.Background(), job.ID)

	return job, nil
}

// GetJob retrieves a scan job by ID.
func (s *Service) GetJob(ctx context.Context, jobID string) (*ScanJob, error) {
	return s.loadJob(ctx, jobID)
}

// HasPrivateRepoSupport returns true if private repo scanning is available.
func (s *Service) HasPrivateRepoSupport() bool {
	return s.cloner.HasToken()
}

// runScan executes the full scan pipeline.
func (s *Service) runScan(ctx context.Context, jobID string) {
	var repoPath string
	var err error

	defer func() {
		// Cleanup cloned repo
		if repoPath != "" {
			_ = s.cloner.Cleanup(repoPath)
		}
	}()

	// Load job
	job, err := s.loadJob(ctx, jobID)
	if err != nil {
		return
	}

	// Clone repository
	_ = s.updateJobStatus(ctx, jobID, StatusCloning, "")
	cloneResult, err := s.cloner.Clone(ctx, job.RepoURL)
	if err != nil {
		_ = s.failJob(ctx, jobID, fmt.Sprintf("Clone failed: %v", err))
		return
	}
	repoPath = cloneResult.Path

	// Detect languages
	languages, err := s.detector.DetectLanguages(repoPath)
	if err != nil {
		_ = s.failJob(ctx, jobID, fmt.Sprintf("Language detection failed: %v", err))
		return
	}

	// Convert to string slice for storage
	langStrings := make([]string, len(languages))
	for i, l := range languages {
		langStrings[i] = string(l)
	}
	_ = s.updateJobLanguages(ctx, jobID, langStrings)

	// Run security tools
	_ = s.updateJobStatus(ctx, jobID, StatusScanning, "")
	toolNames := s.toolRunner.GetToolsForLanguages(languages)
	var results []ToolResult
	for _, toolName := range toolNames {
		result := s.toolRunner.RunToolByName(ctx, toolName, repoPath, languages)
		results = append(results, result)
	}

	// Aggregate findings
	findings := s.aggregator.AggregateAndProcess(results)

	// AI review (if findings exist and client available)
	if len(findings) > 0 && s.reviewer.HasClient() {
		_ = s.updateJobStatus(ctx, jobID, StatusReviewing, "")
		findings, _ = s.reviewer.Review(ctx, repoPath, findings)
	}

	// Complete job
	_ = s.completeJob(ctx, jobID, findings)
}

// Database operations

func (s *Service) createJob(ctx context.Context, job *ScanJob) error {
	query := `
		INSERT INTO scan_jobs (id, repo_url, status, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	expiresAt := job.CreatedAt.Add(7 * 24 * time.Hour) // 7 days retention

	_, err := s.db.ExecContext(ctx, query,
		job.ID, job.RepoURL, job.Status, job.CreatedAt, expiresAt)
	return err
}

func (s *Service) loadJob(ctx context.Context, jobID string) (*ScanJob, error) {
	job := &ScanJob{}

	query := `
		SELECT id, repo_url, status, languages, error, created_at, completed_at
		FROM scan_jobs
		WHERE id = $1
	`

	var languagesJSON []byte
	var errorStr sql.NullString
	var completedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, query, jobID).Scan(
		&job.ID, &job.RepoURL, &job.Status, &languagesJSON,
		&errorStr, &job.CreatedAt, &completedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrJobNotFound
	}
	if err != nil {
		return nil, err
	}

	if languagesJSON != nil {
		_ = json.Unmarshal(languagesJSON, &job.Languages)
	}
	if errorStr.Valid {
		job.Error = errorStr.String
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}

	// Load findings
	findings, err := s.loadFindings(ctx, jobID)
	if err == nil {
		job.Findings = findings
	}

	return job, nil
}

func (s *Service) loadFindings(ctx context.Context, jobID string) ([]Finding, error) {
	query := `
		SELECT id, severity, tool, file_path, line_number, description, remediation, code_example
		FROM scan_findings
		WHERE scan_job_id = $1
		ORDER BY 
			CASE severity 
				WHEN 'critical' THEN 0 
				WHEN 'high' THEN 1 
				WHEN 'medium' THEN 2 
				WHEN 'low' THEN 3 
				ELSE 4 
			END
	`

	rows, err := s.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var findings []Finding
	for rows.Next() {
		var f Finding
		var lineNumber sql.NullInt64
		var remediation, codeExample sql.NullString

		err := rows.Scan(
			&f.ID, &f.Severity, &f.Tool, &f.FilePath, &lineNumber,
			&f.Description, &remediation, &codeExample,
		)
		if err != nil {
			return nil, err
		}

		if lineNumber.Valid {
			ln := int(lineNumber.Int64)
			f.LineNumber = &ln
		}
		if remediation.Valid {
			f.Remediation = remediation.String
		}
		if codeExample.Valid {
			f.CodeExample = codeExample.String
		}

		findings = append(findings, f)
	}

	return findings, rows.Err()
}

func (s *Service) updateJobStatus(ctx context.Context, jobID, status, errorMsg string) error {
	query := `UPDATE scan_jobs SET status = $1, error = $2 WHERE id = $3`
	var errPtr *string
	if errorMsg != "" {
		errPtr = &errorMsg
	}
	_, err := s.db.ExecContext(ctx, query, status, errPtr, jobID)
	return err
}

func (s *Service) updateJobLanguages(ctx context.Context, jobID string, languages []string) error {
	languagesJSON, _ := json.Marshal(languages)
	query := `UPDATE scan_jobs SET languages = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, languagesJSON, jobID)
	return err
}

func (s *Service) failJob(ctx context.Context, jobID, errorMsg string) error {
	now := time.Now()
	query := `UPDATE scan_jobs SET status = $1, error = $2, completed_at = $3 WHERE id = $4`
	_, err := s.db.ExecContext(ctx, query, StatusFailed, errorMsg, now, jobID)
	return err
}

func (s *Service) completeJob(ctx context.Context, jobID string, findings []Finding) error {
	now := time.Now()

	// Update job status
	query := `UPDATE scan_jobs SET status = $1, completed_at = $2 WHERE id = $3`
	_, err := s.db.ExecContext(ctx, query, StatusCompleted, now, jobID)
	if err != nil {
		return err
	}

	// Insert findings
	for _, f := range findings {
		err := s.insertFinding(ctx, jobID, f)
		if err != nil {
			// Log but continue
			continue
		}
	}

	return nil
}

func (s *Service) insertFinding(ctx context.Context, jobID string, f Finding) error {
	query := `
		INSERT INTO scan_findings (id, scan_job_id, severity, tool, file_path, line_number, description, remediation, code_example)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	var lineNumber *int
	if f.LineNumber != nil {
		lineNumber = f.LineNumber
	}

	var remediation, codeExample *string
	if f.Remediation != "" {
		remediation = &f.Remediation
	}
	if f.CodeExample != "" {
		codeExample = &f.CodeExample
	}

	_, err := s.db.ExecContext(ctx, query,
		f.ID, jobID, f.Severity, f.Tool, f.FilePath, lineNumber,
		f.Description, remediation, codeExample,
	)
	return err
}

// GetConfig returns the scanner configuration.
func (s *Service) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"private_repo_enabled": s.HasPrivateRepoSupport(),
		"ai_review_enabled":    s.reviewer.HasClient(),
		"max_files_to_review":  s.reviewer.GetMaxFiles(),
	}
}
