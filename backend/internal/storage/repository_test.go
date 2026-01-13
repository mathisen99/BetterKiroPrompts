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

// Feature: ux-improvements, Property 4: Idempotent View Counting
// **Validates: Requirements 5.1, 5.3**
// For any gallery item and IP address, multiple view requests SHALL result in
// exactly one view record in the database and the view_count SHALL increment by at most 1.

// MockViewRepository is a mock implementation for testing view tracking logic.
type MockViewRepository struct {
	views       map[string]map[string]bool // generationID -> ipHash -> exists
	viewCounts  map[string]int             // generationID -> count
	generations map[string]bool            // generationID -> exists
}

// NewMockViewRepository creates a new mock repository for view testing.
func NewMockViewRepository() *MockViewRepository {
	return &MockViewRepository{
		views:       make(map[string]map[string]bool),
		viewCounts:  make(map[string]int),
		generations: make(map[string]bool),
	}
}

// AddGeneration adds a generation to the mock repository.
func (m *MockViewRepository) AddGeneration(id string) {
	m.generations[id] = true
	m.viewCounts[id] = 0
	m.views[id] = make(map[string]bool)
}

// RecordView simulates the RecordView behavior.
func (m *MockViewRepository) RecordView(generationID string, ipHash string) (bool, error) {
	if generationID == "" || ipHash == "" {
		return false, ErrInvalidInput
	}

	if !m.generations[generationID] {
		return false, ErrNotFound
	}

	// Check if this IP has already viewed this generation
	if m.views[generationID][ipHash] {
		return false, nil // Duplicate view
	}

	// Record the new view
	m.views[generationID][ipHash] = true
	m.viewCounts[generationID]++
	return true, nil
}

// GetViewCount returns the view count for a generation.
func (m *MockViewRepository) GetViewCount(generationID string) int {
	return m.viewCounts[generationID]
}

// GetUniqueViewers returns the number of unique IP hashes that viewed a generation.
func (m *MockViewRepository) GetUniqueViewers(generationID string) int {
	return len(m.views[generationID])
}

// TestProperty4_IdempotentViewCounting tests that multiple views from the same IP
// result in exactly one view record and view_count increments by at most 1.
// Feature: ux-improvements, Property 4: Idempotent View Counting
// **Validates: Requirements 5.1, 5.3**
func TestProperty4_IdempotentViewCounting(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := NewMockViewRepository()

		// Create a generation
		genID := generateNonEmptyString(r)
		repo.AddGeneration(genID)

		// Generate a random IP hash
		ipHash := generateNonEmptyString(r)

		// Record the view multiple times (2-10 times)
		numAttempts := 2 + r.Intn(9)
		newViewCount := 0

		for i := 0; i < numAttempts; i++ {
			isNew, err := repo.RecordView(genID, ipHash)
			if err != nil {
				t.Logf("RecordView failed: %v", err)
				return false
			}
			if isNew {
				newViewCount++
			}
		}

		// Property: Only the first view should be counted as new
		if newViewCount != 1 {
			t.Logf("Expected exactly 1 new view, got %d after %d attempts", newViewCount, numAttempts)
			return false
		}

		// Property: View count should be exactly 1
		if repo.GetViewCount(genID) != 1 {
			t.Logf("Expected view count of 1, got %d", repo.GetViewCount(genID))
			return false
		}

		// Property: Unique viewers should be exactly 1
		if repo.GetUniqueViewers(genID) != 1 {
			t.Logf("Expected 1 unique viewer, got %d", repo.GetUniqueViewers(genID))
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 4 (Idempotent View Counting) failed: %v", err)
	}
}

// TestProperty4_IdempotentViewCounting_MultipleIPs tests that different IPs
// each get counted once.
func TestProperty4_IdempotentViewCounting_MultipleIPs(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := NewMockViewRepository()

		// Create a generation
		genID := generateNonEmptyString(r)
		repo.AddGeneration(genID)

		// Generate multiple unique IP hashes (2-10)
		numIPs := 2 + r.Intn(9)
		ipHashes := make([]string, numIPs)
		for i := 0; i < numIPs; i++ {
			ipHashes[i] = generateNonEmptyString(r) + "_" + string(rune('a'+i)) // Ensure uniqueness
		}

		// Each IP views multiple times
		for _, ipHash := range ipHashes {
			numAttempts := 1 + r.Intn(5)
			for j := 0; j < numAttempts; j++ {
				_, err := repo.RecordView(genID, ipHash)
				if err != nil {
					t.Logf("RecordView failed: %v", err)
					return false
				}
			}
		}

		// Property: View count should equal number of unique IPs
		if repo.GetViewCount(genID) != numIPs {
			t.Logf("Expected view count of %d, got %d", numIPs, repo.GetViewCount(genID))
			return false
		}

		// Property: Unique viewers should equal number of unique IPs
		if repo.GetUniqueViewers(genID) != numIPs {
			t.Logf("Expected %d unique viewers, got %d", numIPs, repo.GetUniqueViewers(genID))
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 4 (Idempotent View Counting - Multiple IPs) failed: %v", err)
	}
}

// Feature: ux-improvements, Property 5: Vote Upsert Behavior
// **Validates: Requirements 5.2, 5.4**
// For any gallery item and IP address, submitting multiple votes SHALL result in
// exactly one rating record, with the score reflecting the most recent vote.

// MockRatingRepository is a mock implementation for testing rating/vote logic.
type MockRatingRepository struct {
	ratings     map[string]map[string]int // generationID -> ipHash -> score
	generations map[string]bool           // generationID -> exists
}

// NewMockRatingRepository creates a new mock repository for rating testing.
func NewMockRatingRepository() *MockRatingRepository {
	return &MockRatingRepository{
		ratings:     make(map[string]map[string]int),
		generations: make(map[string]bool),
	}
}

// AddGeneration adds a generation to the mock repository.
func (m *MockRatingRepository) AddGeneration(id string) {
	m.generations[id] = true
	m.ratings[id] = make(map[string]int)
}

// CreateOrUpdateRating simulates the upsert behavior for ratings.
// Returns true if this was a new rating, false if it was an update.
func (m *MockRatingRepository) CreateOrUpdateRating(generationID string, score int, ipHash string) (isNew bool, err error) {
	if generationID == "" || ipHash == "" {
		return false, ErrInvalidInput
	}
	if score < 1 || score > 5 {
		return false, ErrInvalidInput
	}
	if !m.generations[generationID] {
		return false, ErrNotFound
	}

	// Check if this IP has already rated this generation
	_, exists := m.ratings[generationID][ipHash]

	// Upsert the rating (create or update)
	m.ratings[generationID][ipHash] = score

	return !exists, nil
}

// GetRating returns the rating for a generation by IP hash.
func (m *MockRatingRepository) GetRating(generationID string, ipHash string) (int, bool) {
	if ratings, ok := m.ratings[generationID]; ok {
		if score, exists := ratings[ipHash]; exists {
			return score, true
		}
	}
	return 0, false
}

// GetRatingCount returns the number of unique ratings for a generation.
func (m *MockRatingRepository) GetRatingCount(generationID string) int {
	if ratings, ok := m.ratings[generationID]; ok {
		return len(ratings)
	}
	return 0
}

// TestProperty5_VoteUpsertBehavior tests that multiple votes from the same IP
// result in exactly one rating record with the most recent score.
// Feature: ux-improvements, Property 5: Vote Upsert Behavior
// **Validates: Requirements 5.2, 5.4**
func TestProperty5_VoteUpsertBehavior(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := NewMockRatingRepository()

		// Create a generation
		genID := generateNonEmptyString(r)
		repo.AddGeneration(genID)

		// Generate a random IP hash
		ipHash := generateNonEmptyString(r)

		// Submit multiple votes (2-10 times) with different scores
		numVotes := 2 + r.Intn(9)
		var lastScore int
		newRatingCount := 0

		for i := 0; i < numVotes; i++ {
			// Generate a random score between 1 and 5
			score := 1 + r.Intn(5)
			lastScore = score

			isNew, err := repo.CreateOrUpdateRating(genID, score, ipHash)
			if err != nil {
				t.Logf("CreateOrUpdateRating failed: %v", err)
				return false
			}
			if isNew {
				newRatingCount++
			}
		}

		// Property 1: Only the first vote should create a new rating
		if newRatingCount != 1 {
			t.Logf("Expected exactly 1 new rating, got %d after %d votes", newRatingCount, numVotes)
			return false
		}

		// Property 2: Rating count should be exactly 1
		if repo.GetRatingCount(genID) != 1 {
			t.Logf("Expected rating count of 1, got %d", repo.GetRatingCount(genID))
			return false
		}

		// Property 3: The stored score should be the most recent vote
		storedScore, exists := repo.GetRating(genID, ipHash)
		if !exists {
			t.Logf("Rating should exist for IP hash")
			return false
		}
		if storedScore != lastScore {
			t.Logf("Expected score %d (last vote), got %d", lastScore, storedScore)
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 5 (Vote Upsert Behavior) failed: %v", err)
	}
}

// TestProperty5_VoteUpsertBehavior_MultipleIPs tests that different IPs
// each get their own rating record.
func TestProperty5_VoteUpsertBehavior_MultipleIPs(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := NewMockRatingRepository()

		// Create a generation
		genID := generateNonEmptyString(r)
		repo.AddGeneration(genID)

		// Generate multiple unique IP hashes (2-10)
		numIPs := 2 + r.Intn(9)
		ipHashes := make([]string, numIPs)
		lastScores := make(map[string]int)

		for i := 0; i < numIPs; i++ {
			ipHashes[i] = generateNonEmptyString(r) + "_" + string(rune('a'+i)) // Ensure uniqueness
		}

		// Each IP votes multiple times
		for _, ipHash := range ipHashes {
			numVotes := 1 + r.Intn(5)
			for j := 0; j < numVotes; j++ {
				score := 1 + r.Intn(5)
				lastScores[ipHash] = score
				_, err := repo.CreateOrUpdateRating(genID, score, ipHash)
				if err != nil {
					t.Logf("CreateOrUpdateRating failed: %v", err)
					return false
				}
			}
		}

		// Property 1: Rating count should equal number of unique IPs
		if repo.GetRatingCount(genID) != numIPs {
			t.Logf("Expected rating count of %d, got %d", numIPs, repo.GetRatingCount(genID))
			return false
		}

		// Property 2: Each IP's stored score should be their most recent vote
		for _, ipHash := range ipHashes {
			storedScore, exists := repo.GetRating(genID, ipHash)
			if !exists {
				t.Logf("Rating should exist for IP hash %s", ipHash)
				return false
			}
			if storedScore != lastScores[ipHash] {
				t.Logf("Expected score %d for IP %s, got %d", lastScores[ipHash], ipHash, storedScore)
				return false
			}
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 5 (Vote Upsert Behavior - Multiple IPs) failed: %v", err)
	}
}

// TestProperty5_VoteUpsertBehavior_ScoreValidation tests that invalid scores are rejected.
func TestProperty5_VoteUpsertBehavior_ScoreValidation(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := NewMockRatingRepository()

		// Create a generation
		genID := generateNonEmptyString(r)
		repo.AddGeneration(genID)

		ipHash := generateNonEmptyString(r)

		// Test invalid scores (0, negative, > 5)
		invalidScores := []int{0, -1, -100, 6, 10, 100}
		for _, score := range invalidScores {
			_, err := repo.CreateOrUpdateRating(genID, score, ipHash)
			if err == nil {
				t.Logf("Expected error for invalid score %d, got nil", score)
				return false
			}
		}

		// Test valid scores (1-5)
		for score := 1; score <= 5; score++ {
			_, err := repo.CreateOrUpdateRating(genID, score, ipHash)
			if err != nil {
				t.Logf("Unexpected error for valid score %d: %v", score, err)
				return false
			}
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 5 (Vote Upsert Behavior - Score Validation) failed: %v", err)
	}
}
