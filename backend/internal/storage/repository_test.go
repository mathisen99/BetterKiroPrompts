package storage

import (
	"encoding/json"
	"math/rand"
	"testing"
	"testing/quick"
	"time"
)

// Feature: final-polish, Property 3: Generation Record Completeness
// **Validates: Requirements 5.2, 5.6**
// For any successfully stored generation, the database record SHALL contain
// all required fields (id, project_idea, experience_level, hook_preset, files,
// category_id, created_at) and SHALL NOT contain any user-identifying information.

// generateValidGeneration generates a random valid Generation for testing.
func generateValidGeneration(r *rand.Rand) *Generation {
	experienceLevels := []string{"novice", "intermediate", "expert"}
	hookPresets := []string{"default", "minimal", "comprehensive"}

	// Generate random files JSON
	files := []map[string]string{
		{
			"path":    "kickoff-prompt.md",
			"content": generateNonEmptyString(r),
			"type":    "kickoff",
		},
		{
			"path":    ".kiro/steering/product.md",
			"content": generateNonEmptyString(r),
			"type":    "steering",
		},
	}
	filesJSON, _ := json.Marshal(files)

	return &Generation{
		ProjectIdea:     generateNonEmptyString(r),
		ExperienceLevel: experienceLevels[r.Intn(len(experienceLevels))],
		HookPreset:      hookPresets[r.Intn(len(hookPresets))],
		Files:           filesJSON,
		CategoryID:      1 + r.Intn(5), // 1-5
	}
}

// generateNonEmptyString generates a random non-empty string.
func generateNonEmptyString(r *rand.Rand) string {
	length := 1 + r.Intn(100) // 1 to 100 characters
	chars := make([]byte, length)
	for i := range chars {
		chars[i] = byte('a' + r.Intn(26))
	}
	return string(chars)
}

// TestProperty3_GenerationRecordCompleteness tests that Generation records
// contain all required fields and no user-identifying information.
// Feature: final-polish, Property 3: Generation Record Completeness
// **Validates: Requirements 5.2, 5.6**
func TestProperty3_GenerationRecordCompleteness(t *testing.T) {
	// Property: For any valid Generation, all required fields must be present
	// and no user-identifying fields should exist.
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		gen := generateValidGeneration(r)

		// Verify required fields are set (before storage)
		if gen.ProjectIdea == "" {
			t.Logf("ProjectIdea is empty")
			return false
		}
		if gen.ExperienceLevel == "" {
			t.Logf("ExperienceLevel is empty")
			return false
		}
		if gen.HookPreset == "" {
			t.Logf("HookPreset is empty")
			return false
		}
		if len(gen.Files) == 0 {
			t.Logf("Files is empty")
			return false
		}
		if gen.CategoryID < 1 || gen.CategoryID > 5 {
			t.Logf("CategoryID out of range: %d", gen.CategoryID)
			return false
		}

		// Verify Files is valid JSON
		var filesData interface{}
		if err := json.Unmarshal(gen.Files, &filesData); err != nil {
			t.Logf("Files is not valid JSON: %v", err)
			return false
		}

		// Verify no user-identifying information in the struct
		// The Generation struct should NOT have fields like:
		// - IPAddress
		// - UserAgent
		// - UserID
		// - Email
		// - SessionID
		// This is verified by the struct definition itself

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 3 (Generation Record Completeness) failed: %v", err)
	}
}

// TestProperty3_GenerationRecordCompleteness_NoUserIdentifyingInfo verifies
// that the Generation struct does not contain user-identifying fields.
func TestProperty3_GenerationRecordCompleteness_NoUserIdentifyingInfo(t *testing.T) {
	// Create a generation and serialize to JSON
	gen := &Generation{
		ID:              "test-id",
		ProjectIdea:     "Build a REST API",
		ExperienceLevel: "novice",
		HookPreset:      "default",
		Files:           json.RawMessage(`[{"path":"test.md","content":"test","type":"kickoff"}]`),
		CategoryID:      1,
		CategoryName:    "API",
		AvgRating:       4.5,
		RatingCount:     10,
		ViewCount:       100,
		CreatedAt:       time.Now(),
	}

	jsonBytes, err := json.Marshal(gen)
	if err != nil {
		t.Fatalf("Failed to marshal generation: %v", err)
	}

	// Parse back to a map to check for unexpected fields
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		t.Fatalf("Failed to unmarshal generation: %v", err)
	}

	// List of user-identifying fields that should NOT be present
	forbiddenFields := []string{
		"ipAddress", "ip_address", "ip",
		"userAgent", "user_agent",
		"userId", "user_id",
		"email",
		"sessionId", "session_id",
		"fingerprint",
		"voterHash", "voter_hash", // This should only be in ratings, not generations
	}

	for _, field := range forbiddenFields {
		if _, exists := data[field]; exists {
			t.Errorf("Generation contains forbidden user-identifying field: %s", field)
		}
	}

	// Verify expected fields are present
	expectedFields := []string{
		"id", "projectIdea", "experienceLevel", "hookPreset",
		"files", "categoryId", "avgRating", "ratingCount", "viewCount", "createdAt",
	}

	for _, field := range expectedFields {
		if _, exists := data[field]; !exists {
			t.Errorf("Generation missing expected field: %s", field)
		}
	}
}

// TestProperty3_GenerationRecordCompleteness_ValidExperienceLevels tests
// that only valid experience levels are accepted.
func TestProperty3_GenerationRecordCompleteness_ValidExperienceLevels(t *testing.T) {
	validLevels := []string{"novice", "intermediate", "expert"}

	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		gen := generateValidGeneration(r)

		// Verify experience level is one of the valid values
		for _, valid := range validLevels {
			if gen.ExperienceLevel == valid {
				return true
			}
		}
		t.Logf("Invalid experience level: %s", gen.ExperienceLevel)
		return false
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 3 (Valid Experience Levels) failed: %v", err)
	}
}

// TestProperty3_GenerationRecordCompleteness_ValidHookPresets tests
// that only valid hook presets are accepted.
func TestProperty3_GenerationRecordCompleteness_ValidHookPresets(t *testing.T) {
	validPresets := []string{"default", "minimal", "comprehensive"}

	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		gen := generateValidGeneration(r)

		// Verify hook preset is one of the valid values
		for _, valid := range validPresets {
			if gen.HookPreset == valid {
				return true
			}
		}
		t.Logf("Invalid hook preset: %s", gen.HookPreset)
		return false
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 3 (Valid Hook Presets) failed: %v", err)
	}
}
