package scanner

import (
	"better-kiro-prompts/internal/openai"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Default configuration for code review.
const (
	DefaultMaxFilesToReview    = 10
	DefaultMaxFindingsToReview = 10
	DefaultMaxFileSize         = 50 * 1024 // 50KB max file size
)

// ReviewableSeverities defines which severities get AI review (high and medium only)
var ReviewableSeverities = map[string]bool{
	"critical": true,
	"high":     true,
	"medium":   true,
	// "low" and "info" are excluded
}

// CodeReviewer uses AI to provide remediation guidance for security findings.
type CodeReviewer struct {
	client   *openai.Client
	maxFiles int
	model    string
	log      *slog.Logger
}

// CodeReviewerOption is a functional option for configuring a CodeReviewer.
type CodeReviewerOption func(*CodeReviewer)

// WithMaxFiles sets the maximum number of files to review.
func WithMaxFiles(max int) CodeReviewerOption {
	return func(r *CodeReviewer) {
		r.maxFiles = max
	}
}

// WithModel sets the model to use for code review.
func WithModel(model string) CodeReviewerOption {
	return func(r *CodeReviewer) {
		r.model = model
	}
}

// NewCodeReviewer creates a new CodeReviewer.
func NewCodeReviewer(client *openai.Client, opts ...CodeReviewerOption) *CodeReviewer {
	r := &CodeReviewer{
		client:   client,
		maxFiles: DefaultMaxFilesToReview,
		model:    "gpt-5.1-codex-max", // Use codex model for security code review
		log:      slog.Default().With("component", "reviewer"),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// codeReviewSystemPrompt is the system prompt for the AI code reviewer.
const codeReviewSystemPrompt = `You are a security code reviewer. Your task is to analyze code files that have been flagged by security scanning tools and provide actionable remediation guidance.

For each finding:
1. Explain what the security issue is in plain language
2. Explain why it's a problem (potential impact)
3. Provide a concrete code fix with before/after examples
4. Keep explanations concise and actionable

Format your response as JSON:
{
  "findings": [
    {
      "file_path": "path/to/file",
      "line_number": 42,
      "remediation": "Clear explanation of the fix",
      "code_example": "// Before:\n...\n\n// After:\n..."
    }
  ]
}

Focus on practical fixes. Do not invent new vulnerabilities - only address the specific issues flagged.`

// ReviewResponse represents the AI's response structure.
type ReviewResponse struct {
	Findings []ReviewFinding `json:"findings"`
}

// ReviewFinding represents a single finding's remediation from the AI.
type ReviewFinding struct {
	FilePath    string `json:"file_path"`
	LineNumber  int    `json:"line_number,omitempty"`
	Remediation string `json:"remediation"`
	CodeExample string `json:"code_example"`
}

// ReviewResult contains the reviewed findings and statistics.
type ReviewResult struct {
	Findings []Finding
	Stats    ReviewStats
}

// ReviewStats tracks AI review statistics.
type ReviewStats struct {
	TotalFindings      int `json:"total_findings"`
	ReviewableFindings int `json:"reviewable_findings"` // high/medium/critical only
	ReviewedFindings   int `json:"reviewed_findings"`   // actually sent to AI (max 10)
	MatchedFindings    int `json:"matched_findings"`    // successfully matched with AI response
}

// Review analyzes findings and adds AI-generated remediation guidance.
// It only reviews high/medium severity findings, limited to DefaultMaxFindingsToReview.
func (r *CodeReviewer) Review(ctx context.Context, repoPath string, findings []Finding) (ReviewResult, error) {
	stats := ReviewStats{TotalFindings: len(findings)}

	if r.client == nil {
		// No AI client configured, return findings as-is
		r.log.Info("review_skipped", slog.String("reason", "no_ai_client"))
		return ReviewResult{Findings: findings, Stats: stats}, nil
	}

	if len(findings) == 0 {
		// No findings to review
		return ReviewResult{Findings: findings, Stats: stats}, nil
	}

	// Filter to only high/medium severity findings
	var reviewableFindings []Finding
	for _, f := range findings {
		if ReviewableSeverities[f.Severity] {
			reviewableFindings = append(reviewableFindings, f)
		}
	}
	stats.ReviewableFindings = len(reviewableFindings)

	r.log.Info("findings_filtered",
		slog.Int("total", len(findings)),
		slog.Int("reviewable", len(reviewableFindings)))

	if len(reviewableFindings) == 0 {
		r.log.Info("review_skipped", slog.String("reason", "no_high_medium_findings"))
		return ReviewResult{Findings: findings, Stats: stats}, nil
	}

	// Limit to max findings for review
	if len(reviewableFindings) > DefaultMaxFindingsToReview {
		reviewableFindings = reviewableFindings[:DefaultMaxFindingsToReview]
		r.log.Info("findings_limited", slog.Int("limit", DefaultMaxFindingsToReview))
	}
	stats.ReviewedFindings = len(reviewableFindings)

	// Get unique files with findings, prioritized by severity
	filesToReview := r.selectFilesToReview(reviewableFindings)
	r.log.Info("files_selected", slog.Int("count", len(filesToReview)))

	// Read file contents
	fileContents := make(map[string]string)
	for _, filePath := range filesToReview {
		// File paths from tools may be absolute or relative
		var fullPath string
		if strings.HasPrefix(filePath, repoPath) {
			fullPath = filePath
		} else if strings.HasPrefix(filePath, "/") {
			fullPath = filePath
		} else {
			fullPath = filepath.Join(repoPath, filePath)
		}

		content, err := r.readFileContent(fullPath)
		if err != nil {
			r.log.Warn("file_read_failed", slog.String("path", fullPath), slog.String("error", err.Error()))
			continue
		}
		// Store with relative path for cleaner prompts
		relPath := strings.TrimPrefix(filePath, repoPath+"/")
		fileContents[relPath] = content
	}

	r.log.Info("files_read", slog.Int("count", len(fileContents)))

	if len(fileContents) == 0 {
		// No files could be read
		return ReviewResult{Findings: findings, Stats: stats}, nil
	}

	// Build the review request
	userPrompt := r.buildUserPrompt(reviewableFindings, fileContents)

	// Call the AI with codex model
	messages := []openai.Message{
		{Role: "system", Content: codeReviewSystemPrompt},
		{Role: "user", Content: userPrompt},
	}

	response, err := r.client.ChatCompletionWithModel(ctx, messages, r.model)
	if err != nil {
		// AI review failed, log and return findings without remediation
		r.log.Error("ai_review_failed", slog.String("error", err.Error()))
		return ReviewResult{Findings: findings, Stats: stats}, nil
	}

	r.log.Info("ai_response_received", slog.Int("length", len(response)))

	// Parse the response
	reviewResponse, err := r.parseResponse(response)
	if err != nil {
		// Failed to parse, log and return findings without remediation
		r.log.Error("parse_failed", slog.String("error", err.Error()))
		return ReviewResult{Findings: findings, Stats: stats}, nil
	}

	r.log.Info("remediation_parsed", slog.Int("count", len(reviewResponse.Findings)))

	// Merge remediation into findings and get match count
	mergedFindings, matchCount := r.mergeRemediation(findings, reviewResponse)
	stats.MatchedFindings = matchCount

	return ReviewResult{Findings: mergedFindings, Stats: stats}, nil
}

// selectFilesToReview selects files to review, prioritizing by severity.
// Returns at most maxFiles files. When files have the same severity,
// they are sorted alphabetically by path for deterministic ordering.
func (r *CodeReviewer) selectFilesToReview(findings []Finding) []string {
	// Group findings by file
	fileFindings := make(map[string][]Finding)
	for _, f := range findings {
		fileFindings[f.FilePath] = append(fileFindings[f.FilePath], f)
	}

	// Score each file by highest severity finding
	type fileScore struct {
		path  string
		score int
	}

	var scores []fileScore
	for path, ff := range fileFindings {
		// Find highest severity (lowest score = most severe)
		minScore := 999
		for _, f := range ff {
			if s, ok := severityOrder[f.Severity]; ok && s < minScore {
				minScore = s
			}
		}
		scores = append(scores, fileScore{path: path, score: minScore})
	}

	// Sort by score (most severe first), then by path (alphabetically) for determinism
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].score != scores[j].score {
			return scores[i].score < scores[j].score
		}
		// Secondary sort by path for deterministic ordering
		return scores[i].path < scores[j].path
	})

	// Take top maxFiles
	var files []string
	for i, s := range scores {
		if i >= r.maxFiles {
			break
		}
		files = append(files, s.path)
	}

	return files
}

// readFileContent reads a file's content, respecting size limits.
func (r *CodeReviewer) readFileContent(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.Size() > DefaultMaxFileSize {
		// File too large, read only the beginning
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer func() { _ = file.Close() }()

		buf := make([]byte, DefaultMaxFileSize)
		n, err := file.Read(buf)
		if err != nil {
			return "", err
		}
		return string(buf[:n]) + "\n... (truncated)", nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// buildUserPrompt builds the user prompt for the AI.
func (r *CodeReviewer) buildUserPrompt(findings []Finding, fileContents map[string]string) string {
	var sb strings.Builder
	sb.WriteString("Review these security findings and provide remediation:\n\n")

	// Group findings by file - need to match both absolute and relative paths
	byFile := make(map[string][]Finding)
	for _, f := range findings {
		// Try to find matching file content
		for contentPath := range fileContents {
			// Match if the finding path ends with the content path (handles absolute vs relative)
			if strings.HasSuffix(f.FilePath, "/"+contentPath) || f.FilePath == contentPath {
				byFile[contentPath] = append(byFile[contentPath], f)
				break
			}
		}
	}

	for filePath, ff := range byFile {
		sb.WriteString(fmt.Sprintf("## File: %s\n\n", filePath))

		// List findings for this file
		sb.WriteString("### Findings:\n")
		for _, f := range ff {
			lineInfo := ""
			if f.LineNumber != nil {
				lineInfo = fmt.Sprintf(" (line %d)", *f.LineNumber)
			}
			sb.WriteString(fmt.Sprintf("- [%s] %s%s: %s\n", f.Severity, f.Tool, lineInfo, f.Description))
		}
		sb.WriteString("\n")

		// Include file content
		sb.WriteString("### Code:\n```\n")
		sb.WriteString(fileContents[filePath])
		sb.WriteString("\n```\n\n")
	}

	return sb.String()
}

// parseResponse parses the AI's JSON response.
func (r *CodeReviewer) parseResponse(response string) (*ReviewResponse, error) {
	// Try to extract JSON from the response
	response = strings.TrimSpace(response)

	// Handle markdown code blocks
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		if idx := strings.Index(response, "```"); idx != -1 {
			response = response[:idx]
		}
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		if idx := strings.Index(response, "```"); idx != -1 {
			response = response[:idx]
		}
	}

	response = strings.TrimSpace(response)

	var result ReviewResponse
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

// mergeRemediation merges AI remediation into findings and returns match count.
func (r *CodeReviewer) mergeRemediation(findings []Finding, review *ReviewResponse) ([]Finding, int) {
	if review == nil {
		return findings, 0
	}

	r.log.Info("merge_start",
		slog.Int("findings", len(findings)),
		slog.Int("review_items", len(review.Findings)))

	// Log what the AI returned
	for i, rf := range review.Findings {
		r.log.Debug("review_item",
			slog.Int("index", i),
			slog.String("file", rf.FilePath),
			slog.Int("line", rf.LineNumber),
			slog.Bool("has_remediation", rf.Remediation != ""))
	}

	// Create a lookup map for remediation - use just filename for matching
	remediationMap := make(map[string]ReviewFinding)
	for _, rf := range review.Findings {
		// Store with full key
		key := fmt.Sprintf("%s:%d", rf.FilePath, rf.LineNumber)
		remediationMap[key] = rf
		// Also store file-only key
		fileKey := fmt.Sprintf("%s:0", rf.FilePath)
		if _, exists := remediationMap[fileKey]; !exists {
			remediationMap[fileKey] = rf
		}
	}

	// Merge into findings
	result := make([]Finding, len(findings))
	matchCount := 0
	for i, f := range findings {
		result[i] = f

		lineNum := 0
		if f.LineNumber != nil {
			lineNum = *f.LineNumber
		}

		// Extract just the filename from the finding path for matching
		findingFile := f.FilePath
		if idx := strings.LastIndex(findingFile, "/"); idx != -1 {
			findingFile = findingFile[idx+1:]
		}

		// Try exact match with line number first
		key := fmt.Sprintf("%s:%d", findingFile, lineNum)
		if rf, ok := remediationMap[key]; ok {
			result[i].Remediation = rf.Remediation
			result[i].CodeExample = rf.CodeExample
			matchCount++
			r.log.Debug("match_found", slog.String("file", findingFile), slog.Int("line", lineNum), slog.String("method", "filename:line"))
			continue
		}

		// Try with full path
		key = fmt.Sprintf("%s:%d", f.FilePath, lineNum)
		if rf, ok := remediationMap[key]; ok {
			result[i].Remediation = rf.Remediation
			result[i].CodeExample = rf.CodeExample
			matchCount++
			r.log.Debug("match_found", slog.String("file", f.FilePath), slog.Int("line", lineNum), slog.String("method", "fullpath:line"))
			continue
		}

		// Try file-only match (no line number)
		key = fmt.Sprintf("%s:0", findingFile)
		if rf, ok := remediationMap[key]; ok {
			result[i].Remediation = rf.Remediation
			result[i].CodeExample = rf.CodeExample
			matchCount++
			r.log.Debug("match_found", slog.String("file", findingFile), slog.String("method", "filename_only"))
			continue
		}

		// Try full path file-only match
		key = fmt.Sprintf("%s:0", f.FilePath)
		if rf, ok := remediationMap[key]; ok {
			result[i].Remediation = rf.Remediation
			result[i].CodeExample = rf.CodeExample
			matchCount++
			r.log.Debug("match_found", slog.String("file", f.FilePath), slog.String("method", "fullpath_only"))
			continue
		}

		// Only log no_match for reviewable severities (high/medium/critical)
		if ReviewableSeverities[f.Severity] {
			r.log.Warn("no_match",
				slog.String("file_path", f.FilePath),
				slog.String("extracted_file", findingFile),
				slog.Int("line", lineNum))
		}
	}

	r.log.Info("merge_complete",
		slog.Int("matched", matchCount),
		slog.Int("total", len(findings)))
	return result, matchCount
}

// GetMaxFiles returns the maximum number of files to review.
func (r *CodeReviewer) GetMaxFiles() int {
	return r.maxFiles
}

// HasClient returns true if an AI client is configured.
func (r *CodeReviewer) HasClient() bool {
	return r.client != nil
}

// GetFilesToReview returns the files that would be reviewed for the given findings.
// This is useful for testing the file selection logic.
func (r *CodeReviewer) GetFilesToReview(findings []Finding) []string {
	return r.selectFilesToReview(findings)
}
