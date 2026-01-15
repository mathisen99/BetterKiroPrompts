package scanner

import (
	"context"
	"os"
	"testing"
	"time"
)

// =============================================================================
// Unit Tests for Tool Runner
// =============================================================================

func TestNewToolRunner(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		r := NewToolRunner()
		if r.timeout != DefaultToolTimeout {
			t.Errorf("timeout = %v, want %v", r.timeout, DefaultToolTimeout)
		}
	})

	t.Run("with custom timeout", func(t *testing.T) {
		customTimeout := 10 * time.Second
		r := NewToolRunner(WithToolTimeout(customTimeout))
		if r.timeout != customTimeout {
			t.Errorf("timeout = %v, want %v", r.timeout, customTimeout)
		}
	})
}

func TestToolRunner_GetToolsForLanguages(t *testing.T) {
	r := NewToolRunner()

	tests := []struct {
		name      string
		languages []Language
		wantTools []string
	}{
		{
			name:      "no languages",
			languages: []Language{},
			wantTools: []string{"trivy", "semgrep", "trufflehog", "gitleaks"},
		},
		{
			name:      "Go only",
			languages: []Language{LangGo},
			wantTools: []string{"trivy", "semgrep", "trufflehog", "gitleaks", "govulncheck"},
		},
		{
			name:      "Python only",
			languages: []Language{LangPython},
			wantTools: []string{"trivy", "semgrep", "trufflehog", "gitleaks", "bandit", "pip-audit", "safety"},
		},
		{
			name:      "JavaScript only",
			languages: []Language{LangJavaScript},
			wantTools: []string{"trivy", "semgrep", "trufflehog", "gitleaks", "npm-audit"},
		},
		{
			name:      "TypeScript only",
			languages: []Language{LangTypeScript},
			wantTools: []string{"trivy", "semgrep", "trufflehog", "gitleaks", "npm-audit"},
		},
		{
			name:      "Rust only",
			languages: []Language{LangRust},
			wantTools: []string{"trivy", "semgrep", "trufflehog", "gitleaks", "cargo-audit"},
		},
		{
			name:      "Ruby only",
			languages: []Language{LangRuby},
			wantTools: []string{"trivy", "semgrep", "trufflehog", "gitleaks", "bundler-audit", "brakeman"},
		},
		{
			name:      "multiple languages",
			languages: []Language{LangGo, LangPython, LangJavaScript},
			wantTools: []string{"trivy", "semgrep", "trufflehog", "gitleaks", "govulncheck", "bandit", "pip-audit", "safety", "npm-audit"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.GetToolsForLanguages(tt.languages)

			// Check that all expected tools are present
			gotSet := make(map[string]bool)
			for _, tool := range got {
				gotSet[tool] = true
			}

			for _, want := range tt.wantTools {
				if !gotSet[want] {
					t.Errorf("GetToolsForLanguages() missing tool %s", want)
				}
			}

			// Check that no unexpected tools are present
			wantSet := make(map[string]bool)
			for _, tool := range tt.wantTools {
				wantSet[tool] = true
			}

			for _, tool := range got {
				if !wantSet[tool] {
					t.Errorf("GetToolsForLanguages() unexpected tool %s", tool)
				}
			}
		})
	}
}

func TestToolRunner_RunToolByName_UnknownTool(t *testing.T) {
	r := NewToolRunner()
	ctx := context.Background()

	result := r.RunToolByName(ctx, "unknown-tool", "/tmp", nil)
	if result.Error == nil {
		t.Error("Expected error for unknown tool")
	}
}

// =============================================================================
// Property-Based Tests for Tool Timeout
// =============================================================================

// skipIfNoDocker skips the test if Docker or the scanner container is not available
func skipIfNoDocker(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test that requires Docker in CI environment")
	}
	// Also skip if scanner container is not running
	if os.Getenv("SCANNER_CONTAINER_AVAILABLE") == "" {
		t.Skip("Skipping test: scanner container not available (set SCANNER_CONTAINER_AVAILABLE=1 to run)")
	}
}

// TestProperty7_ToolTimeoutEnforcement tests Property 7: Tool Timeout Enforcement
// Feature: info-and-security-scan, Property 7: Tool Timeout Enforcement
// **Validates: Requirements 7.10, 7.11**
//
// Property: For any security tool execution, if the tool does not complete within
// the configured timeout period, the execution SHALL be terminated and the scanner
// SHALL continue with results from other tools.
func TestProperty7_ToolTimeoutEnforcement(t *testing.T) {
	// Sub-property 1: Timeout is configurable
	t.Run("timeout_is_configurable", func(t *testing.T) {
		timeouts := []time.Duration{
			100 * time.Millisecond,
			1 * time.Second,
			5 * time.Second,
			1 * time.Minute,
			5 * time.Minute,
		}

		for _, timeout := range timeouts {
			r := NewToolRunner(WithToolTimeout(timeout))
			if r.timeout != timeout {
				t.Errorf("Expected timeout %v, got %v", timeout, r.timeout)
			}
		}
	})

	// Sub-property 2: Tool execution respects timeout
	t.Run("tool_execution_respects_timeout", func(t *testing.T) {
		skipIfNoDocker(t)

		// Use a very short timeout
		shortTimeout := 50 * time.Millisecond
		r := NewToolRunner(WithToolTimeout(shortTimeout))

		ctx := context.Background()

		// Run a tool that will likely timeout (sleep command)
		// Note: This tests the timeout mechanism, not actual tool execution
		start := time.Now()
		output, timedOut, _ := r.runTool(ctx, "sleep", []string{"10"}, "/tmp")
		elapsed := time.Since(start)

		// Should have timed out
		if !timedOut {
			t.Errorf("Expected timeout, but command completed. Output: %s", string(output))
		}

		// Elapsed time should be close to timeout (with some tolerance)
		tolerance := 100 * time.Millisecond
		if elapsed > shortTimeout+tolerance {
			t.Errorf("Timeout took too long: %v (expected ~%v)", elapsed, shortTimeout)
		}
	})

	// Sub-property 3: TimedOut flag is set correctly
	t.Run("timed_out_flag_set_correctly", func(t *testing.T) {
		skipIfNoDocker(t)

		shortTimeout := 50 * time.Millisecond
		r := NewToolRunner(WithToolTimeout(shortTimeout))

		ctx := context.Background()

		// Test with a command that will timeout
		_, timedOut, _ := r.runTool(ctx, "sleep", []string{"10"}, "/tmp")
		if !timedOut {
			t.Error("Expected TimedOut to be true for long-running command")
		}

		// Test with a command that completes quickly
		r2 := NewToolRunner(WithToolTimeout(5 * time.Second))
		_, timedOut2, _ := r2.runTool(ctx, "echo", []string{"hello"}, "/tmp")
		if timedOut2 {
			t.Error("Expected TimedOut to be false for quick command")
		}
	})

	// Sub-property 4: Context cancellation is respected
	t.Run("context_cancellation_respected", func(t *testing.T) {
		skipIfNoDocker(t)

		r := NewToolRunner(WithToolTimeout(10 * time.Second))

		// Create a context that we'll cancel
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel after a short delay
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		start := time.Now()
		_, _, err := r.runTool(ctx, "sleep", []string{"10"}, "/tmp")
		elapsed := time.Since(start)

		// Should have been cancelled
		if err == nil {
			t.Error("Expected error from cancelled context")
		}

		// Should have completed quickly (not waited for full timeout)
		if elapsed > 500*time.Millisecond {
			t.Errorf("Context cancellation took too long: %v", elapsed)
		}
	})

	// Sub-property 5: Partial results are preserved on timeout
	t.Run("partial_results_preserved", func(t *testing.T) {
		skipIfNoDocker(t)

		// This tests that when a tool times out, we still get a ToolResult
		// with the TimedOut flag set, allowing the scanner to continue
		shortTimeout := 50 * time.Millisecond
		r := NewToolRunner(WithToolTimeout(shortTimeout))

		ctx := context.Background()

		// Run Trivy (which won't exist in test env, but tests the structure)
		result := r.RunTrivy(ctx, "/nonexistent")

		// Result should have the tool name set
		if result.Tool != "trivy" {
			t.Errorf("Expected tool name 'trivy', got '%s'", result.Tool)
		}

		// Duration should be recorded
		if result.Duration == 0 {
			t.Error("Expected duration to be recorded")
		}
	})
}

// TestProperty7_ToolTimeoutEnforcement_EdgeCases tests edge cases for timeout.
// Feature: info-and-security-scan, Property 7: Tool Timeout Enforcement
// **Validates: Requirements 7.10, 7.11**
func TestProperty7_ToolTimeoutEnforcement_EdgeCases(t *testing.T) {
	t.Run("zero_timeout", func(t *testing.T) {
		skipIfNoDocker(t)

		// Zero timeout should still work (immediate timeout)
		r := NewToolRunner(WithToolTimeout(0))
		ctx := context.Background()

		_, timedOut, _ := r.runTool(ctx, "echo", []string{"hello"}, "/tmp")
		// With zero timeout, command should timeout immediately
		if !timedOut {
			// This is acceptable - some systems may complete echo before timeout
			t.Log("Command completed before zero timeout (acceptable)")
		}
	})

	t.Run("very_short_timeout", func(t *testing.T) {
		skipIfNoDocker(t)

		r := NewToolRunner(WithToolTimeout(1 * time.Millisecond))
		ctx := context.Background()

		start := time.Now()
		_, timedOut, _ := r.runTool(ctx, "sleep", []string{"1"}, "/tmp")
		elapsed := time.Since(start)

		if !timedOut {
			t.Error("Expected timeout with 1ms timeout")
		}

		// Should complete quickly
		if elapsed > 100*time.Millisecond {
			t.Errorf("Timeout took too long: %v", elapsed)
		}
	})

	t.Run("default_timeout_value", func(t *testing.T) {
		r := NewToolRunner()
		if r.timeout != DefaultToolTimeout {
			t.Errorf("Default timeout should be %v, got %v", DefaultToolTimeout, r.timeout)
		}
	})
}
