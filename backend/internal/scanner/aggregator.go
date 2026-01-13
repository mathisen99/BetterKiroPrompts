package scanner

import (
	"sort"
	"strings"

	"github.com/google/uuid"
)

// Severity levels for findings.
const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityLow      = "low"
	SeverityInfo     = "info"
)

// severityOrder defines the sort order for severities (lower = more severe).
var severityOrder = map[string]int{
	SeverityCritical: 0,
	SeverityHigh:     1,
	SeverityMedium:   2,
	SeverityLow:      3,
	SeverityInfo:     4,
}

// Finding represents an aggregated security finding.
type Finding struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	Tool        string `json:"tool"`
	FilePath    string `json:"file_path"`
	LineNumber  *int   `json:"line_number,omitempty"`
	Description string `json:"description"`
	Remediation string `json:"remediation,omitempty"`
	CodeExample string `json:"code_example,omitempty"`
	RuleID      string `json:"rule_id,omitempty"`
}

// Aggregator aggregates and deduplicates findings from multiple tools.
type Aggregator struct{}

// NewAggregator creates a new Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{}
}

// Aggregate converts tool results into unified findings.
func (a *Aggregator) Aggregate(results []ToolResult) []Finding {
	var findings []Finding

	for _, result := range results {
		if result.TimedOut || result.Error != nil {
			continue
		}

		for _, raw := range result.Findings {
			finding := a.convertRawFinding(raw, result.Tool)
			findings = append(findings, finding)
		}
	}

	return findings
}

// convertRawFinding converts a RawFinding to a Finding.
func (a *Aggregator) convertRawFinding(raw RawFinding, tool string) Finding {
	finding := Finding{
		ID:          uuid.New().String(),
		Tool:        tool,
		FilePath:    raw.FilePath,
		Description: raw.Description,
		Severity:    a.normalizeSeverity(raw.Severity),
		RuleID:      raw.RuleID,
	}

	if raw.LineNumber > 0 {
		lineNum := raw.LineNumber
		finding.LineNumber = &lineNum
	}

	return finding
}

// normalizeSeverity normalizes severity strings to standard values.
func (a *Aggregator) normalizeSeverity(severity string) string {
	severity = strings.ToLower(strings.TrimSpace(severity))

	switch severity {
	case "critical", "crit":
		return SeverityCritical
	case "high", "error":
		return SeverityHigh
	case "medium", "moderate", "warning", "warn":
		return SeverityMedium
	case "low":
		return SeverityLow
	case "info", "informational", "note":
		return SeverityInfo
	default:
		// Default to medium if unknown
		return SeverityMedium
	}
}

// Deduplicate removes duplicate findings based on file, line, and description.
func (a *Aggregator) Deduplicate(findings []Finding) []Finding {
	seen := make(map[string]bool)
	var unique []Finding

	for _, f := range findings {
		key := a.dedupeKey(f)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, f)
		}
	}

	return unique
}

// dedupeKey generates a unique key for deduplication.
func (a *Aggregator) dedupeKey(f Finding) string {
	lineStr := "0"
	if f.LineNumber != nil {
		lineStr = string(rune(*f.LineNumber))
	}
	return f.FilePath + ":" + lineStr + ":" + f.Description
}

// RankBySeverity sorts findings by severity (critical first, info last).
func (a *Aggregator) RankBySeverity(findings []Finding) []Finding {
	sorted := make([]Finding, len(findings))
	copy(sorted, findings)

	sort.SliceStable(sorted, func(i, j int) bool {
		orderI := severityOrder[sorted[i].Severity]
		orderJ := severityOrder[sorted[j].Severity]
		return orderI < orderJ
	})

	return sorted
}

// AggregateAndProcess performs full aggregation: aggregate, dedupe, and rank.
func (a *Aggregator) AggregateAndProcess(results []ToolResult) []Finding {
	findings := a.Aggregate(results)
	findings = a.Deduplicate(findings)
	findings = a.RankBySeverity(findings)
	return findings
}

// GetUniqueFiles returns a list of unique file paths from findings.
func (a *Aggregator) GetUniqueFiles(findings []Finding) []string {
	seen := make(map[string]bool)
	var files []string

	for _, f := range findings {
		if f.FilePath != "" && !seen[f.FilePath] {
			seen[f.FilePath] = true
			files = append(files, f.FilePath)
		}
	}

	return files
}

// GetFindingsByFile groups findings by file path.
func (a *Aggregator) GetFindingsByFile(findings []Finding) map[string][]Finding {
	byFile := make(map[string][]Finding)

	for _, f := range findings {
		byFile[f.FilePath] = append(byFile[f.FilePath], f)
	}

	return byFile
}

// GetFindingsBySeverity groups findings by severity.
func (a *Aggregator) GetFindingsBySeverity(findings []Finding) map[string][]Finding {
	bySeverity := make(map[string][]Finding)

	for _, f := range findings {
		bySeverity[f.Severity] = append(bySeverity[f.Severity], f)
	}

	return bySeverity
}

// CountBySeverity returns counts of findings by severity.
func (a *Aggregator) CountBySeverity(findings []Finding) map[string]int {
	counts := make(map[string]int)

	for _, f := range findings {
		counts[f.Severity]++
	}

	return counts
}

// FilterBySeverity returns findings with severity at or above the given level.
func (a *Aggregator) FilterBySeverity(findings []Finding, minSeverity string) []Finding {
	minOrder := severityOrder[minSeverity]
	var filtered []Finding

	for _, f := range findings {
		if severityOrder[f.Severity] <= minOrder {
			filtered = append(filtered, f)
		}
	}

	return filtered
}

// IsValidSeverity checks if a severity string is valid.
func IsValidSeverity(severity string) bool {
	_, ok := severityOrder[severity]
	return ok
}
