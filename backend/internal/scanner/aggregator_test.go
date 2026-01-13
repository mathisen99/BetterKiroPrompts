package scanner

import (
	"testing"
	"testing/quick"
)

// =============================================================================
// Unit Tests for Aggregator
// =============================================================================

func TestAggregator_normalizeSeverity(t *testing.T) {
	a := NewAggregator()

	tests := []struct {
		input string
		want  string
	}{
		// Critical
		{"critical", SeverityCritical},
		{"CRITICAL", SeverityCritical},
		{"crit", SeverityCritical},

		// High
		{"high", SeverityHigh},
		{"HIGH", SeverityHigh},
		{"error", SeverityHigh},

		// Medium
		{"medium", SeverityMedium},
		{"MEDIUM", SeverityMedium},
		{"moderate", SeverityMedium},
		{"warning", SeverityMedium},
		{"warn", SeverityMedium},

		// Low
		{"low", SeverityLow},
		{"LOW", SeverityLow},

		// Info
		{"info", SeverityInfo},
		{"INFO", SeverityInfo},
		{"informational", SeverityInfo},
		{"note", SeverityInfo},

		// Unknown defaults to medium
		{"unknown", SeverityMedium},
		{"", SeverityMedium},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := a.normalizeSeverity(tt.input)
			if got != tt.want {
				t.Errorf("normalizeSeverity(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAggregator_Aggregate(t *testing.T) {
	a := NewAggregator()

	results := []ToolResult{
		{
			Tool: "trivy",
			Findings: []RawFinding{
				{FilePath: "main.go", LineNumber: 10, Description: "Issue 1", Severity: "high"},
				{FilePath: "util.go", LineNumber: 20, Description: "Issue 2", Severity: "medium"},
			},
		},
		{
			Tool: "semgrep",
			Findings: []RawFinding{
				{FilePath: "main.go", LineNumber: 15, Description: "Issue 3", Severity: "low"},
			},
		},
		{
			Tool:     "gitleaks",
			TimedOut: true, // Should be skipped
			Findings: []RawFinding{
				{FilePath: "secret.txt", Description: "Secret", Severity: "critical"},
			},
		},
	}

	findings := a.Aggregate(results)

	if len(findings) != 3 {
		t.Errorf("Expected 3 findings, got %d", len(findings))
	}

	// Verify tool attribution
	toolCounts := make(map[string]int)
	for _, f := range findings {
		toolCounts[f.Tool]++
	}

	if toolCounts["trivy"] != 2 {
		t.Errorf("Expected 2 trivy findings, got %d", toolCounts["trivy"])
	}
	if toolCounts["semgrep"] != 1 {
		t.Errorf("Expected 1 semgrep finding, got %d", toolCounts["semgrep"])
	}
	if toolCounts["gitleaks"] != 0 {
		t.Errorf("Expected 0 gitleaks findings (timed out), got %d", toolCounts["gitleaks"])
	}
}

func TestAggregator_Deduplicate(t *testing.T) {
	a := NewAggregator()

	lineNum := 10
	findings := []Finding{
		{ID: "1", FilePath: "main.go", LineNumber: &lineNum, Description: "Issue 1", Tool: "trivy"},
		{ID: "2", FilePath: "main.go", LineNumber: &lineNum, Description: "Issue 1", Tool: "semgrep"}, // Duplicate
		{ID: "3", FilePath: "main.go", LineNumber: &lineNum, Description: "Issue 2", Tool: "trivy"},   // Different description
		{ID: "4", FilePath: "util.go", LineNumber: &lineNum, Description: "Issue 1", Tool: "trivy"},   // Different file
	}

	unique := a.Deduplicate(findings)

	if len(unique) != 3 {
		t.Errorf("Expected 3 unique findings, got %d", len(unique))
	}
}

func TestAggregator_RankBySeverity(t *testing.T) {
	a := NewAggregator()

	findings := []Finding{
		{ID: "1", Severity: SeverityLow},
		{ID: "2", Severity: SeverityCritical},
		{ID: "3", Severity: SeverityMedium},
		{ID: "4", Severity: SeverityHigh},
		{ID: "5", Severity: SeverityInfo},
	}

	ranked := a.RankBySeverity(findings)

	expectedOrder := []string{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo}
	for i, expected := range expectedOrder {
		if ranked[i].Severity != expected {
			t.Errorf("Position %d: expected %s, got %s", i, expected, ranked[i].Severity)
		}
	}
}

func TestAggregator_GetUniqueFiles(t *testing.T) {
	a := NewAggregator()

	findings := []Finding{
		{FilePath: "main.go"},
		{FilePath: "util.go"},
		{FilePath: "main.go"}, // Duplicate
		{FilePath: "test.go"},
	}

	files := a.GetUniqueFiles(findings)

	if len(files) != 3 {
		t.Errorf("Expected 3 unique files, got %d", len(files))
	}
}

func TestAggregator_FilterBySeverity(t *testing.T) {
	a := NewAggregator()

	findings := []Finding{
		{ID: "1", Severity: SeverityCritical},
		{ID: "2", Severity: SeverityHigh},
		{ID: "3", Severity: SeverityMedium},
		{ID: "4", Severity: SeverityLow},
		{ID: "5", Severity: SeverityInfo},
	}

	// Filter for high and above
	filtered := a.FilterBySeverity(findings, SeverityHigh)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 findings (critical, high), got %d", len(filtered))
	}

	// Filter for medium and above
	filtered = a.FilterBySeverity(findings, SeverityMedium)
	if len(filtered) != 3 {
		t.Errorf("Expected 3 findings (critical, high, medium), got %d", len(filtered))
	}
}

func TestIsValidSeverity(t *testing.T) {
	validSeverities := []string{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo}
	for _, s := range validSeverities {
		if !IsValidSeverity(s) {
			t.Errorf("Expected %s to be valid", s)
		}
	}

	invalidSeverities := []string{"unknown", "severe", ""}
	for _, s := range invalidSeverities {
		if IsValidSeverity(s) {
			t.Errorf("Expected %s to be invalid", s)
		}
	}
}

// =============================================================================
// Property-Based Tests for Finding Aggregation
// =============================================================================

// TestProperty8_FindingAggregationCompleteness tests Property 8: Finding Aggregation Completeness
// Feature: info-and-security-scan, Property 8: Finding Aggregation Completeness
// **Validates: Requirements 8.1, 8.2, 8.3, 8.4**
//
// Property: For any set of tool results aggregated into findings:
// - Each finding SHALL have a non-empty file_path
// - Each finding SHALL have a non-empty description
// - Each finding SHALL have a valid severity (critical, high, medium, low, info)
// - Each finding SHALL have a tool source identifier
// - Duplicate findings (same file, line, description) SHALL be deduplicated
// - Findings SHALL be sorted by severity (critical first, info last)
func TestProperty8_FindingAggregationCompleteness(t *testing.T) {
	a := NewAggregator()

	// Sub-property 1: All findings have required fields
	t.Run("findings_have_required_fields", func(t *testing.T) {
		results := []ToolResult{
			{
				Tool: "trivy",
				Findings: []RawFinding{
					{FilePath: "main.go", Description: "Issue 1", Severity: "high"},
					{FilePath: "util.go", Description: "Issue 2", Severity: "medium"},
				},
			},
			{
				Tool: "semgrep",
				Findings: []RawFinding{
					{FilePath: "app.ts", Description: "Issue 3", Severity: "low"},
				},
			},
		}

		findings := a.AggregateAndProcess(results)

		for _, f := range findings {
			// Check non-empty file_path
			if f.FilePath == "" {
				t.Error("Finding has empty file_path")
			}

			// Check non-empty description
			if f.Description == "" {
				t.Error("Finding has empty description")
			}

			// Check valid severity
			if !IsValidSeverity(f.Severity) {
				t.Errorf("Finding has invalid severity: %s", f.Severity)
			}

			// Check tool source
			if f.Tool == "" {
				t.Error("Finding has empty tool")
			}

			// Check ID is set
			if f.ID == "" {
				t.Error("Finding has empty ID")
			}
		}
	})

	// Sub-property 2: Deduplication works correctly
	t.Run("deduplication_removes_duplicates", func(t *testing.T) {
		lineNum := 10
		findings := []Finding{
			{ID: "1", FilePath: "main.go", LineNumber: &lineNum, Description: "Same issue", Severity: SeverityHigh, Tool: "trivy"},
			{ID: "2", FilePath: "main.go", LineNumber: &lineNum, Description: "Same issue", Severity: SeverityHigh, Tool: "semgrep"},
			{ID: "3", FilePath: "main.go", LineNumber: &lineNum, Description: "Same issue", Severity: SeverityHigh, Tool: "gitleaks"},
		}

		unique := a.Deduplicate(findings)

		if len(unique) != 1 {
			t.Errorf("Expected 1 unique finding after deduplication, got %d", len(unique))
		}
	})

	// Sub-property 3: Severity sorting is correct
	t.Run("severity_sorting_correct", func(t *testing.T) {
		property := func(numCritical, numHigh, numMedium, numLow, numInfo uint8) bool {
			// Limit to reasonable numbers
			numCritical = numCritical % 10
			numHigh = numHigh % 10
			numMedium = numMedium % 10
			numLow = numLow % 10
			numInfo = numInfo % 10

			var findings []Finding
			for i := uint8(0); i < numInfo; i++ {
				findings = append(findings, Finding{Severity: SeverityInfo})
			}
			for i := uint8(0); i < numLow; i++ {
				findings = append(findings, Finding{Severity: SeverityLow})
			}
			for i := uint8(0); i < numMedium; i++ {
				findings = append(findings, Finding{Severity: SeverityMedium})
			}
			for i := uint8(0); i < numHigh; i++ {
				findings = append(findings, Finding{Severity: SeverityHigh})
			}
			for i := uint8(0); i < numCritical; i++ {
				findings = append(findings, Finding{Severity: SeverityCritical})
			}

			ranked := a.RankBySeverity(findings)

			// Verify ordering
			for i := 1; i < len(ranked); i++ {
				prevOrder := severityOrder[ranked[i-1].Severity]
				currOrder := severityOrder[ranked[i].Severity]
				if prevOrder > currOrder {
					return false
				}
			}

			return true
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 8 (severity sorting) failed: %v", err)
		}
	})

	// Sub-property 4: Severity normalization is consistent
	t.Run("severity_normalization_consistent", func(t *testing.T) {
		severityInputs := []string{
			"critical", "CRITICAL", "crit",
			"high", "HIGH", "error",
			"medium", "MEDIUM", "moderate", "warning",
			"low", "LOW",
			"info", "INFO", "informational",
		}

		for _, input := range severityInputs {
			normalized := a.normalizeSeverity(input)
			if !IsValidSeverity(normalized) {
				t.Errorf("normalizeSeverity(%q) = %q is not valid", input, normalized)
			}
		}
	})

	// Sub-property 5: Aggregation preserves all non-timed-out findings
	t.Run("aggregation_preserves_findings", func(t *testing.T) {
		property := func(numFindings uint8) bool {
			numFindings = numFindings % 50 // Limit to reasonable number

			var rawFindings []RawFinding
			for i := uint8(0); i < numFindings; i++ {
				rawFindings = append(rawFindings, RawFinding{
					FilePath:    "file.go",
					Description: "Issue",
					Severity:    "high",
				})
			}

			results := []ToolResult{
				{Tool: "test", Findings: rawFindings},
			}

			findings := a.Aggregate(results)

			// All findings should be preserved (before deduplication)
			return len(findings) == int(numFindings)
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 8 (aggregation preserves findings) failed: %v", err)
		}
	})

	// Sub-property 6: Timed out results are excluded
	t.Run("timed_out_results_excluded", func(t *testing.T) {
		results := []ToolResult{
			{
				Tool:     "trivy",
				TimedOut: true,
				Findings: []RawFinding{
					{FilePath: "main.go", Description: "Issue", Severity: "high"},
				},
			},
			{
				Tool: "semgrep",
				Findings: []RawFinding{
					{FilePath: "util.go", Description: "Issue", Severity: "medium"},
				},
			},
		}

		findings := a.Aggregate(results)

		// Only semgrep findings should be included
		if len(findings) != 1 {
			t.Errorf("Expected 1 finding (trivy timed out), got %d", len(findings))
		}

		if findings[0].Tool != "semgrep" {
			t.Errorf("Expected semgrep finding, got %s", findings[0].Tool)
		}
	})
}

// TestProperty8_FindingAggregationCompleteness_EdgeCases tests edge cases.
// Feature: info-and-security-scan, Property 8: Finding Aggregation Completeness
// **Validates: Requirements 8.1, 8.2, 8.3, 8.4**
func TestProperty8_FindingAggregationCompleteness_EdgeCases(t *testing.T) {
	a := NewAggregator()

	t.Run("empty_results", func(t *testing.T) {
		findings := a.AggregateAndProcess([]ToolResult{})
		if len(findings) != 0 {
			t.Errorf("Expected 0 findings for empty results, got %d", len(findings))
		}
	})

	t.Run("all_tools_timed_out", func(t *testing.T) {
		results := []ToolResult{
			{Tool: "trivy", TimedOut: true, Findings: []RawFinding{{FilePath: "a.go", Description: "x", Severity: "high"}}},
			{Tool: "semgrep", TimedOut: true, Findings: []RawFinding{{FilePath: "b.go", Description: "y", Severity: "high"}}},
		}

		findings := a.AggregateAndProcess(results)
		if len(findings) != 0 {
			t.Errorf("Expected 0 findings when all tools timed out, got %d", len(findings))
		}
	})

	t.Run("findings_without_line_numbers", func(t *testing.T) {
		results := []ToolResult{
			{
				Tool: "trivy",
				Findings: []RawFinding{
					{FilePath: "main.go", Description: "Issue", Severity: "high", LineNumber: 0},
				},
			},
		}

		findings := a.Aggregate(results)
		if len(findings) != 1 {
			t.Fatalf("Expected 1 finding, got %d", len(findings))
		}

		if findings[0].LineNumber != nil {
			t.Error("Expected nil LineNumber for line 0")
		}
	})

	t.Run("unknown_severity_defaults_to_medium", func(t *testing.T) {
		results := []ToolResult{
			{
				Tool: "custom",
				Findings: []RawFinding{
					{FilePath: "main.go", Description: "Issue", Severity: "unknown-severity"},
				},
			},
		}

		findings := a.Aggregate(results)
		if findings[0].Severity != SeverityMedium {
			t.Errorf("Expected medium severity for unknown, got %s", findings[0].Severity)
		}
	})
}
