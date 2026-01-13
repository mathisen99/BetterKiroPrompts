package scanner

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"better-kiro-prompts/internal/logger"
	"better-kiro-prompts/internal/openai"

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
	ID          string       `json:"id"`
	Status      string       `json:"status"`
	RepoURL     string       `json:"repo_url"`
	Languages   []string     `json:"languages"`
	Findings    []Finding    `json:"findings"`
	ReviewStats *ReviewStats `json:"review_stats,omitempty"`
	Error       string       `json:"error,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
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
	log        *slog.Logger
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

// WithServiceLogger sets the logger for the service.
func WithServiceLogger(log *slog.Logger) ServiceOption {
	return func(s *Service) {
		if log != nil {
			s.log = log
		}
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
		log:        slog.Default(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// SetLogger sets the logger for the service.
func (s *Service) SetLogger(log *slog.Logger) {
	if log != nil {
		s.log = log
	}
}

// StartScan initiates a new security scan.
func (s *Service) StartScan(ctx context.Context, req ScanRequest) (*ScanJob, error) {
	requestID := logger.GetRequestID(ctx)

	s.log.Info("scan_start_request",
		slog.String("request_id", requestID),
		slog.String("repo_url", req.RepoURL),
	)

	// Validate URL
	if err := ValidateGitHubURL(req.RepoURL); err != nil {
		s.log.Warn("scan_validation_failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
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
		s.log.Error("scan_create_job_failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	s.log.Info("scan_job_created",
		slog.String("request_id", requestID),
		slog.String("job_id", job.ID),
		slog.String("repo_url", job.RepoURL),
	)

	// Start scan in background
	go s.runScan(context.Background(), job.ID)

	return job, nil
}

// GetJob retrieves a scan job by ID.
func (s *Service) GetJob(ctx context.Context, jobID string) (*ScanJob, error) {
	requestID := logger.GetRequestID(ctx)

	s.log.Debug("scan_get_job_start",
		slog.String("request_id", requestID),
		slog.String("job_id", jobID),
	)

	job, err := s.loadJob(ctx, jobID)
	if err != nil {
		if errors.Is(err, ErrJobNotFound) {
			s.log.Debug("scan_get_job_not_found",
				slog.String("request_id", requestID),
				slog.String("job_id", jobID),
			)
		} else {
			s.log.Error("scan_get_job_failed",
				slog.String("request_id", requestID),
				slog.String("job_id", jobID),
				slog.String("error", err.Error()),
			)
		}
		return nil, err
	}

	s.log.Debug("scan_get_job_complete",
		slog.String("request_id", requestID),
		slog.String("job_id", jobID),
		slog.String("status", job.Status),
		slog.Int("finding_count", len(job.Findings)),
	)

	return job, nil
}

// HasPrivateRepoSupport returns true if private repo scanning is available.
func (s *Service) HasPrivateRepoSupport() bool {
	return s.cloner.HasToken()
}

// runScan executes the full scan pipeline.
func (s *Service) runScan(ctx context.Context, jobID string) {
	var repoPath string
	var err error
	start := time.Now()

	s.log.Info("scan_pipeline_start",
		slog.String("job_id", jobID),
	)

	defer func() {
		// Cleanup cloned repo
		if repoPath != "" {
			s.log.Debug("scan_cleanup_start",
				slog.String("job_id", jobID),
				slog.String("path", repoPath),
			)
			_ = s.cloner.Cleanup(repoPath)
		}
	}()

	// Load job
	job, err := s.loadJob(ctx, jobID)
	if err != nil {
		s.log.Error("scan_load_job_failed",
			slog.String("job_id", jobID),
			slog.String("error", err.Error()),
		)
		return
	}

	// Phase 1: Clone repository
	s.log.Info("scan_phase_clone_start",
		slog.String("job_id", jobID),
		slog.String("repo_url", job.RepoURL),
	)
	cloneStart := time.Now()
	_ = s.updateJobStatus(ctx, jobID, StatusCloning, "")
	cloneResult, err := s.cloner.Clone(ctx, job.RepoURL)
	if err != nil {
		s.log.Error("scan_phase_clone_failed",
			slog.String("job_id", jobID),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(cloneStart)),
		)
		_ = s.failJob(ctx, jobID, fmt.Sprintf("Clone failed: %v", err))
		return
	}
	repoPath = cloneResult.Path
	s.log.Info("scan_phase_clone_complete",
		slog.String("job_id", jobID),
		slog.String("path", repoPath),
		slog.Duration("duration", time.Since(cloneStart)),
	)

	// Phase 2: Detect languages
	s.log.Info("scan_phase_detect_start",
		slog.String("job_id", jobID),
	)
	detectStart := time.Now()
	languages, err := s.detector.DetectLanguages(repoPath)
	if err != nil {
		s.log.Error("scan_phase_detect_failed",
			slog.String("job_id", jobID),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(detectStart)),
		)
		_ = s.failJob(ctx, jobID, fmt.Sprintf("Language detection failed: %v", err))
		return
	}

	// Convert to string slice for storage and logging
	langStrings := make([]string, len(languages))
	for i, l := range languages {
		langStrings[i] = string(l)
	}
	_ = s.updateJobLanguages(ctx, jobID, langStrings)

	s.log.Info("scan_phase_detect_complete",
		slog.String("job_id", jobID),
		slog.Any("languages", langStrings),
		slog.Int("language_count", len(languages)),
		slog.Duration("duration", time.Since(detectStart)),
	)

	// Phase 3: Run security tools
	toolNames := s.toolRunner.GetToolsForLanguages(languages)
	s.log.Info("scan_phase_tools_start",
		slog.String("job_id", jobID),
		slog.Any("tools", toolNames),
		slog.Int("tool_count", len(toolNames)),
	)
	toolsStart := time.Now()
	_ = s.updateJobStatus(ctx, jobID, StatusScanning, "")

	var results []ToolResult
	for _, toolName := range toolNames {
		toolStart := time.Now()
		s.log.Debug("scan_tool_start",
			slog.String("job_id", jobID),
			slog.String("tool", toolName),
		)

		result := s.toolRunner.RunToolByName(ctx, toolName, repoPath, languages)

		s.log.Info("scan_tool_complete",
			slog.String("job_id", jobID),
			slog.String("tool", toolName),
			slog.Int("finding_count", len(result.Findings)),
			slog.Bool("timed_out", result.TimedOut),
			slog.Bool("success", result.Error == nil),
			slog.Duration("duration", time.Since(toolStart)),
		)

		if result.Error != nil {
			s.log.Warn("scan_tool_error",
				slog.String("job_id", jobID),
				slog.String("tool", toolName),
				slog.String("error", result.Error.Error()),
			)
		}

		results = append(results, result)
	}

	s.log.Info("scan_phase_tools_complete",
		slog.String("job_id", jobID),
		slog.Int("tool_count", len(toolNames)),
		slog.Duration("duration", time.Since(toolsStart)),
	)

	// Phase 4: Aggregate findings
	s.log.Info("scan_phase_aggregate_start",
		slog.String("job_id", jobID),
		slog.Int("result_count", len(results)),
	)
	aggStart := time.Now()
	findings := s.aggregator.AggregateAndProcess(results)

	// Count by severity
	severityCounts := map[string]int{"critical": 0, "high": 0, "medium": 0, "low": 0}
	for _, f := range findings {
		severityCounts[f.Severity]++
	}

	s.log.Info("scan_phase_aggregate_complete",
		slog.String("job_id", jobID),
		slog.Int("total_findings", len(findings)),
		slog.Int("critical", severityCounts["critical"]),
		slog.Int("high", severityCounts["high"]),
		slog.Int("medium", severityCounts["medium"]),
		slog.Int("low", severityCounts["low"]),
		slog.Duration("duration", time.Since(aggStart)),
	)

	// Phase 5: AI review (if findings exist and client available)
	var reviewStats *ReviewStats
	if len(findings) > 0 && s.reviewer.HasClient() {
		s.log.Info("scan_phase_review_start",
			slog.String("job_id", jobID),
			slog.Int("findings_to_review", len(findings)),
		)
		reviewStart := time.Now()
		_ = s.updateJobStatus(ctx, jobID, StatusReviewing, "")

		reviewResult, reviewErr := s.reviewer.Review(ctx, repoPath, findings)
		if reviewErr != nil {
			s.log.Warn("scan_phase_review_partial",
				slog.String("job_id", jobID),
				slog.String("error", reviewErr.Error()),
			)
		}
		findings = reviewResult.Findings
		reviewStats = &reviewResult.Stats

		s.log.Info("scan_phase_review_complete",
			slog.String("job_id", jobID),
			slog.Int("reviewed_findings", len(findings)),
			slog.Int("matched_findings", reviewStats.MatchedFindings),
			slog.Duration("duration", time.Since(reviewStart)),
		)
	} else {
		skipReason := "no_findings"
		if len(findings) > 0 {
			skipReason = "no_ai_client"
		}
		s.log.Debug("scan_phase_review_skipped",
			slog.String("job_id", jobID),
			slog.String("reason", skipReason),
			slog.Bool("has_findings", len(findings) > 0),
			slog.Bool("has_client", s.reviewer.HasClient()),
		)
	}

	// Complete job
	_ = s.completeJobWithStats(ctx, jobID, findings, reviewStats)

	s.log.Info("scan_pipeline_complete",
		slog.String("job_id", jobID),
		slog.Int("total_findings", len(findings)),
		slog.Duration("total_duration", time.Since(start)),
	)
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
		SELECT id, repo_url, status, languages, error, created_at, completed_at, review_stats
		FROM scan_jobs
		WHERE id = $1
	`

	var languagesJSON []byte
	var errorStr sql.NullString
	var completedAt sql.NullTime
	var reviewStatsJSON []byte

	err := s.db.QueryRowContext(ctx, query, jobID).Scan(
		&job.ID, &job.RepoURL, &job.Status, &languagesJSON,
		&errorStr, &job.CreatedAt, &completedAt, &reviewStatsJSON,
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
	if reviewStatsJSON != nil {
		var stats ReviewStats
		if json.Unmarshal(reviewStatsJSON, &stats) == nil {
			job.ReviewStats = &stats
		}
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

func (s *Service) completeJobWithStats(ctx context.Context, jobID string, findings []Finding, stats *ReviewStats) error {
	now := time.Now()

	// Update job status with optional review stats
	var err error
	if stats != nil {
		statsJSON, _ := json.Marshal(stats)
		query := `UPDATE scan_jobs SET status = $1, completed_at = $2, review_stats = $3 WHERE id = $4`
		_, err = s.db.ExecContext(ctx, query, StatusCompleted, now, statsJSON, jobID)
	} else {
		query := `UPDATE scan_jobs SET status = $1, completed_at = $2 WHERE id = $3`
		_, err = s.db.ExecContext(ctx, query, StatusCompleted, now, jobID)
	}
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
