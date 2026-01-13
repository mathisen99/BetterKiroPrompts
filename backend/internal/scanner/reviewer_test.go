package scanner

import (
	"testing"
	"testing/quick"
)

// =============================================================================
// Unit Tests for Code Reviewer
// =============================================================================

func TestNewCodeReviewer(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		r := NewCodeReviewer(nil)
		if r.maxFiles != DefaultMaxFilesToReview {
			t.Errorf("maxFiles = %d, want %d", r.maxFiles, DefaultMaxFilesToReview)
		}
	})

	t.Run("with custom max files", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(5))
		if r.maxFiles != 5 {
			t.Errorf("maxFiles = %d, want 5", r.maxFiles)
		}
	})
}

func TestCodeReviewer_HasClient(t *testing.T) {
	t.Run("no client", func(t *testing.T) {
		r := NewCodeReviewer(nil)
		if r.HasClient() {
			t.Error("Expected HasClient() to be false")
		}
	})
}

func TestCodeReviewer_GetMaxFiles(t *testing.T) {
	r := NewCodeReviewer(nil, WithMaxFiles(15))
	if r.GetMaxFiles() != 15 {
		t.Errorf("GetMaxFiles() = %d, want 15", r.GetMaxFiles())
	}
}

func TestCodeReviewer_selectFilesToReview(t *testing.T) {
	r := NewCodeReviewer(nil, WithMaxFiles(3))

	lineNum := 10
	findings := []Finding{
		{FilePath: "critical.go", Severity: SeverityCritical, LineNumber: &lineNum},
		{FilePath: "high.go", Severity: SeverityHigh, LineNumber: &lineNum},
		{FilePath: "medium.go", Severity: SeverityMedium, LineNumber: &lineNum},
		{FilePath: "low.go", Severity: SeverityLow, LineNumber: &lineNum},
		{FilePath: "info.go", Severity: SeverityInfo, LineNumber: &lineNum},
	}

	files := r.selectFilesToReview(findings)

	// Should only return 3 files (maxFiles)
	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}

	// Should prioritize by severity (critical, high, medium)
	expectedFiles := map[string]bool{
		"critical.go": true,
		"high.go":     true,
		"medium.go":   true,
	}

	for _, f := range files {
		if !expectedFiles[f] {
			t.Errorf("Unexpected file in selection: %s", f)
		}
	}
}

func TestCodeReviewer_parseResponse(t *testing.T) {
	r := NewCodeReviewer(nil)

	tests := []struct {
		name     string
		response string
		wantErr  bool
		wantLen  int
	}{
		{
			name: "valid JSON",
			response: `{
				"findings": [
					{"file_path": "main.go", "line_number": 10, "remediation": "Fix it", "code_example": "// fixed"}
				]
			}`,
			wantErr: false,
			wantLen: 1,
		},
		{
			name:     "JSON in markdown code block",
			response: "```json\n{\"findings\": [{\"file_path\": \"main.go\", \"remediation\": \"Fix\"}]}\n```",
			wantErr:  false,
			wantLen:  1,
		},
		{
			name:     "JSON in plain code block",
			response: "```\n{\"findings\": [{\"file_path\": \"main.go\", \"remediation\": \"Fix\"}]}\n```",
			wantErr:  false,
			wantLen:  1,
		},
		{
			name:     "invalid JSON",
			response: "not json at all",
			wantErr:  true,
		},
		{
			name:     "empty findings",
			response: `{"findings": []}`,
			wantErr:  false,
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.parseResponse(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(result.Findings) != tt.wantLen {
				t.Errorf("parseResponse() findings count = %d, want %d", len(result.Findings), tt.wantLen)
			}
		})
	}
}

func TestCodeReviewer_mergeRemediation(t *testing.T) {
	r := NewCodeReviewer(nil)

	lineNum := 10
	findings := []Finding{
		{ID: "1", FilePath: "main.go", LineNumber: &lineNum, Description: "Issue 1"},
		{ID: "2", FilePath: "util.go", LineNumber: &lineNum, Description: "Issue 2"},
	}

	review := &ReviewResponse{
		Findings: []ReviewFinding{
			{FilePath: "main.go", LineNumber: 10, Remediation: "Fix for main.go", CodeExample: "// fixed"},
		},
	}

	merged := r.mergeRemediation(findings, review)

	// First finding should have remediation
	if merged[0].Remediation != "Fix for main.go" {
		t.Errorf("Expected remediation for main.go, got %q", merged[0].Remediation)
	}

	// Second finding should not have remediation
	if merged[1].Remediation != "" {
		t.Errorf("Expected no remediation for util.go, got %q", merged[1].Remediation)
	}
}

// =============================================================================
// Property-Based Tests for AI Review Scope
// =============================================================================

// TestProperty9_AIReviewScopeLimitation tests Property 9: AI Review Scope Limitation
// Feature: info-and-security-scan, Property 9: AI Review Scope Limitation
// **Validates: Requirements 9.2, 9.3, 9.7**
//
// Property: For any scan with findings:
// - The Code_Review SHALL only receive files that have at least one associated finding
// - The number of files sent to Code_Review SHALL NOT exceed the configured maximum (default 10)
// - If there are more flagged files than the maximum, only the files with highest-severity findings SHALL be reviewed
func TestProperty9_AIReviewScopeLimitation(t *testing.T) {
	// Sub-property 1: File count never exceeds maxFiles
	t.Run("file_count_never_exceeds_max", func(t *testing.T) {
		property := func(maxFiles uint8, numFiles uint8) bool {
			// Limit to reasonable values
			maxFiles = (maxFiles % 20) + 1 // 1-20
			numFiles = numFiles % 50       // 0-49

			r := NewCodeReviewer(nil, WithMaxFiles(int(maxFiles)))

			// Generate findings for numFiles different files
			var findings []Finding
			for i := uint8(0); i < numFiles; i++ {
				findings = append(findings, Finding{
					FilePath: "file" + string(rune('0'+i)) + ".go",
					Severity: SeverityMedium,
				})
			}

			files := r.GetFilesToReview(findings)

			// File count should never exceed maxFiles
			if len(files) > int(maxFiles) {
				t.Logf("File count %d exceeds maxFiles %d", len(files), maxFiles)
				return false
			}

			return true
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 9 (file count limit) failed: %v", err)
		}
	})

	// Sub-property 2: Only files with findings are selected
	t.Run("only_files_with_findings_selected", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(10))

		findings := []Finding{
			{FilePath: "has-finding-1.go", Severity: SeverityHigh},
			{FilePath: "has-finding-2.go", Severity: SeverityMedium},
		}

		files := r.GetFilesToReview(findings)

		// All selected files should be in the findings
		findingFiles := make(map[string]bool)
		for _, f := range findings {
			findingFiles[f.FilePath] = true
		}

		for _, file := range files {
			if !findingFiles[file] {
				t.Errorf("Selected file %s has no findings", file)
			}
		}
	})

	// Sub-property 3: Highest severity files are prioritized
	t.Run("highest_severity_files_prioritized", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(2))

		lineNum := 1
		findings := []Finding{
			{FilePath: "info.go", Severity: SeverityInfo, LineNumber: &lineNum},
			{FilePath: "critical.go", Severity: SeverityCritical, LineNumber: &lineNum},
			{FilePath: "low.go", Severity: SeverityLow, LineNumber: &lineNum},
			{FilePath: "high.go", Severity: SeverityHigh, LineNumber: &lineNum},
		}

		files := r.GetFilesToReview(findings)

		// Should select critical and high severity files
		if len(files) != 2 {
			t.Fatalf("Expected 2 files, got %d", len(files))
		}

		// First should be critical
		if files[0] != "critical.go" {
			t.Errorf("Expected critical.go first, got %s", files[0])
		}

		// Second should be high
		if files[1] != "high.go" {
			t.Errorf("Expected high.go second, got %s", files[1])
		}
	})

	// Sub-property 4: Multiple findings per file don't cause duplicates
	t.Run("no_duplicate_files", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(10))

		lineNum := 1
		findings := []Finding{
			{FilePath: "main.go", Severity: SeverityHigh, LineNumber: &lineNum},
			{FilePath: "main.go", Severity: SeverityMedium, LineNumber: &lineNum},
			{FilePath: "main.go", Severity: SeverityLow, LineNumber: &lineNum},
			{FilePath: "util.go", Severity: SeverityHigh, LineNumber: &lineNum},
		}

		files := r.GetFilesToReview(findings)

		// Should only have 2 unique files
		if len(files) != 2 {
			t.Errorf("Expected 2 unique files, got %d", len(files))
		}

		// Check for duplicates
		seen := make(map[string]bool)
		for _, f := range files {
			if seen[f] {
				t.Errorf("Duplicate file in selection: %s", f)
			}
			seen[f] = true
		}
	})

	// Sub-property 5: Empty findings returns empty file list
	t.Run("empty_findings_returns_empty", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(10))

		files := r.GetFilesToReview([]Finding{})

		if len(files) != 0 {
			t.Errorf("Expected 0 files for empty findings, got %d", len(files))
		}
	})

	// Sub-property 6: File selection is deterministic for same severity
	t.Run("selection_deterministic", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(3))

		lineNum := 1
		findings := []Finding{
			{FilePath: "a.go", Severity: SeverityHigh, LineNumber: &lineNum},
			{FilePath: "b.go", Severity: SeverityHigh, LineNumber: &lineNum},
			{FilePath: "c.go", Severity: SeverityHigh, LineNumber: &lineNum},
			{FilePath: "d.go", Severity: SeverityHigh, LineNumber: &lineNum},
			{FilePath: "e.go", Severity: SeverityHigh, LineNumber: &lineNum},
		}

		// Run multiple times and verify consistency
		firstRun := r.GetFilesToReview(findings)
		for i := 0; i < 5; i++ {
			run := r.GetFilesToReview(findings)
			if len(run) != len(firstRun) {
				t.Errorf("Selection not deterministic: different lengths")
			}
			for j, f := range run {
				if f != firstRun[j] {
					t.Errorf("Selection not deterministic: different files")
				}
			}
		}
	})
}

// TestProperty9_AIReviewScopeLimitation_EdgeCases tests edge cases.
// Feature: info-and-security-scan, Property 9: AI Review Scope Limitation
// **Validates: Requirements 9.2, 9.3, 9.7**
func TestProperty9_AIReviewScopeLimitation_EdgeCases(t *testing.T) {
	t.Run("maxFiles_of_1", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(1))

		lineNum := 1
		findings := []Finding{
			{FilePath: "low.go", Severity: SeverityLow, LineNumber: &lineNum},
			{FilePath: "critical.go", Severity: SeverityCritical, LineNumber: &lineNum},
		}

		files := r.GetFilesToReview(findings)

		if len(files) != 1 {
			t.Errorf("Expected 1 file, got %d", len(files))
		}

		if files[0] != "critical.go" {
			t.Errorf("Expected critical.go, got %s", files[0])
		}
	})

	t.Run("exactly_maxFiles_findings", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(3))

		lineNum := 1
		findings := []Finding{
			{FilePath: "a.go", Severity: SeverityHigh, LineNumber: &lineNum},
			{FilePath: "b.go", Severity: SeverityHigh, LineNumber: &lineNum},
			{FilePath: "c.go", Severity: SeverityHigh, LineNumber: &lineNum},
		}

		files := r.GetFilesToReview(findings)

		if len(files) != 3 {
			t.Errorf("Expected 3 files, got %d", len(files))
		}
	})

	t.Run("fewer_files_than_maxFiles", func(t *testing.T) {
		r := NewCodeReviewer(nil, WithMaxFiles(10))

		lineNum := 1
		findings := []Finding{
			{FilePath: "a.go", Severity: SeverityHigh, LineNumber: &lineNum},
			{FilePath: "b.go", Severity: SeverityHigh, LineNumber: &lineNum},
		}

		files := r.GetFilesToReview(findings)

		if len(files) != 2 {
			t.Errorf("Expected 2 files (all available), got %d", len(files))
		}
	})
}
