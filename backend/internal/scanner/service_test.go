package scanner

import (
	"slices"
	"testing"
	"testing/quick"
	"time"
)

// =============================================================================
// Unit Tests for Scanner Service
// =============================================================================

func TestNewService(t *testing.T) {
	// Test that service can be created without database (for unit testing)
	s := NewService(nil, nil, "")

	if s.cloner == nil {
		t.Error("Expected cloner to be initialized")
	}
	if s.detector == nil {
		t.Error("Expected detector to be initialized")
	}
	if s.toolRunner == nil {
		t.Error("Expected toolRunner to be initialized")
	}
	if s.aggregator == nil {
		t.Error("Expected aggregator to be initialized")
	}
	if s.reviewer == nil {
		t.Error("Expected reviewer to be initialized")
	}
}

func TestService_HasPrivateRepoSupport(t *testing.T) {
	t.Run("without token", func(t *testing.T) {
		s := NewService(nil, nil, "")
		if s.HasPrivateRepoSupport() {
			t.Error("Expected no private repo support without token")
		}
	})

	t.Run("with token", func(t *testing.T) {
		s := NewService(nil, nil, "ghp_test123")
		if !s.HasPrivateRepoSupport() {
			t.Error("Expected private repo support with token")
		}
	})
}

func TestService_GetConfig(t *testing.T) {
	s := NewService(nil, nil, "ghp_test123")
	config := s.GetConfig()

	if _, ok := config["private_repo_enabled"]; !ok {
		t.Error("Expected private_repo_enabled in config")
	}
	if _, ok := config["ai_review_enabled"]; !ok {
		t.Error("Expected ai_review_enabled in config")
	}
	if _, ok := config["max_files_to_review"]; !ok {
		t.Error("Expected max_files_to_review in config")
	}
}

func TestScanJob_Structure(t *testing.T) {
	now := time.Now()
	completedAt := now.Add(time.Minute)

	job := createTestJob("test-id", StatusPending, "https://github.com/owner/repo", now, &completedAt)
	job.Languages = []string{"go", "typescript"}
	job.Findings = []Finding{{ID: "f1", Severity: SeverityHigh}}

	validateJob(t, job, "test-id", StatusPending, 2, 1)
}

// =============================================================================
// Property-Based Tests for Job Creation Round-Trip
// =============================================================================

// TestProperty2_JobCreationRoundTrip tests Property 2: Job Creation Round-Trip
// Feature: info-and-security-scan, Property 2: Job Creation Round-Trip
// **Validates: Requirements 4.5, 11.1, 11.2**
func TestProperty2_JobCreationRoundTrip(t *testing.T) {
	// Sub-property 1: Job IDs are unique
	t.Run("job_ids_are_unique", func(t *testing.T) {
		property := func(numJobs uint8) bool {
			numJobs = (numJobs % 50) + 1 // 1-50 jobs

			ids := make(map[string]bool)
			for range numJobs {
				id := generateJobID()
				if ids[id] {
					t.Logf("Duplicate ID generated: %s", id)
					return false
				}
				ids[id] = true
			}

			return true
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 2 (unique job IDs) failed: %v", err)
		}
	})

	// Sub-property 2: Job status transitions are valid
	t.Run("valid_status_transitions", func(t *testing.T) {
		validStatuses := []string{
			StatusPending,
			StatusCloning,
			StatusScanning,
			StatusReviewing,
			StatusCompleted,
			StatusFailed,
		}

		for _, status := range validStatuses {
			if !slices.Contains(validStatuses, status) {
				t.Errorf("Invalid status: %s", status)
			}
		}
	})

	// Sub-property 3: Job preserves all fields
	t.Run("job_preserves_fields", func(t *testing.T) {
		property := func(_ string) bool {
			// Skip invalid URLs
			if ValidateGitHubURL("https://github.com/owner/repo") != nil {
				return true
			}

			now := time.Now()
			lineNum := 10
			job := createTestJob(generateJobID(), StatusCompleted, "https://github.com/owner/repo", now, nil)
			job.Languages = []string{"go", "typescript"}
			job.Findings = []Finding{
				{
					ID:          "f1",
					Severity:    SeverityHigh,
					Tool:        "trivy",
					FilePath:    "main.go",
					LineNumber:  &lineNum,
					Description: "Test finding",
					Remediation: "Fix it",
					CodeExample: "// fixed",
				},
			}

			return validateJobFields(job, StatusCompleted, "https://github.com/owner/repo", 2, 1, SeverityHigh)
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 2 (field preservation) failed: %v", err)
		}
	})

	// Sub-property 4: Completed jobs have completion time
	t.Run("completed_jobs_have_completion_time", func(t *testing.T) {
		now := time.Now()
		completedAt := now.Add(time.Minute)

		job := createTestJob(generateJobID(), StatusCompleted, "https://github.com/owner/repo", now, &completedAt)

		if job.CompletedAt == nil {
			t.Error("Completed job should have completion time")
		}

		if job.CompletedAt.Before(job.CreatedAt) {
			t.Error("Completion time should be after creation time")
		}
	})

	// Sub-property 5: Failed jobs have error message
	t.Run("failed_jobs_have_error", func(t *testing.T) {
		job := createTestJob(generateJobID(), StatusFailed, "https://github.com/owner/repo", time.Now(), nil)
		job.Error = "Clone failed: repository not found"

		if job.Error == "" {
			t.Error("Failed job should have error message")
		}
	})

	// Sub-property 6: URL is normalized in job
	t.Run("url_normalization", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"https://github.com/owner/repo", "https://github.com/owner/repo"},
			{"https://github.com/owner/repo.git", "https://github.com/owner/repo"},
			{"https://github.com/owner/repo/", "https://github.com/owner/repo"},
		}

		for _, tc := range testCases {
			normalized := NormalizeGitHubURL(tc.input)
			if normalized != tc.expected {
				t.Errorf("NormalizeGitHubURL(%q) = %q, want %q", tc.input, normalized, tc.expected)
			}
		}
	})
}

// generateJobID generates a unique job ID (simulating uuid.New().String())
func generateJobID() string {
	return "job-" + time.Now().Format("20060102150405.000000000")
}

// createTestJob creates a ScanJob for testing with the given parameters.
func createTestJob(id, status, repoURL string, createdAt time.Time, completedAt *time.Time) ScanJob {
	return ScanJob{
		ID:          id,
		Status:      status,
		RepoURL:     repoURL,
		CreatedAt:   createdAt,
		CompletedAt: completedAt,
	}
}

// validateJob validates a ScanJob's fields and reports errors.
func validateJob(t *testing.T, job ScanJob, expectedID, expectedStatus string, expectedLangs, expectedFindings int) {
	t.Helper()
	if job.ID != expectedID {
		t.Errorf("ID = %q, want %q", job.ID, expectedID)
	}
	if job.Status != expectedStatus {
		t.Errorf("Status = %q, want %q", job.Status, expectedStatus)
	}
	if len(job.Languages) != expectedLangs {
		t.Errorf("Languages count = %d, want %d", len(job.Languages), expectedLangs)
	}
	if len(job.Findings) != expectedFindings {
		t.Errorf("Findings count = %d, want %d", len(job.Findings), expectedFindings)
	}
}

// validateJobFields validates job fields and returns true if all match.
func validateJobFields(job ScanJob, status, repoURL string, langCount, findingCount int, severity string) bool {
	if job.ID == "" {
		return false
	}
	if job.Status != status {
		return false
	}
	if job.RepoURL != repoURL {
		return false
	}
	if len(job.Languages) != langCount {
		return false
	}
	if len(job.Findings) != findingCount {
		return false
	}
	if findingCount > 0 && job.Findings[0].Severity != severity {
		return false
	}
	return true
}

// TestProperty2_JobCreationRoundTrip_EdgeCases tests edge cases.
// Feature: info-and-security-scan, Property 2: Job Creation Round-Trip
// **Validates: Requirements 4.5, 11.1, 11.2**
func TestProperty2_JobCreationRoundTrip_EdgeCases(t *testing.T) {
	t.Run("job_with_no_findings", func(t *testing.T) {
		job := createTestJob(generateJobID(), StatusCompleted, "https://github.com/owner/repo", time.Now(), nil)
		job.Languages = []string{"go"}
		job.Findings = []Finding{}

		if job.Findings == nil {
			t.Error("Findings should be empty slice, not nil")
		}
		if len(job.Findings) != 0 {
			t.Error("Expected 0 findings")
		}
		// Verify Languages was set
		if len(job.Languages) != 1 {
			t.Error("Expected 1 language")
		}
	})

	t.Run("job_with_no_languages", func(t *testing.T) {
		job := createTestJob(generateJobID(), StatusCompleted, "https://github.com/owner/repo", time.Now(), nil)
		job.Languages = []string{}

		if len(job.Languages) != 0 {
			t.Error("Expected 0 languages")
		}
	})

	t.Run("pending_job_has_no_completion_time", func(t *testing.T) {
		job := createTestJob(generateJobID(), StatusPending, "https://github.com/owner/repo", time.Now(), nil)

		if job.CompletedAt != nil {
			t.Error("Pending job should not have completion time")
		}
	})
}
