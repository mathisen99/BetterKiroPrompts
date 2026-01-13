package gallery

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"sort"
	"testing"
	"testing/quick"
	"time"

	"better-kiro-prompts/internal/ratelimit"
	"better-kiro-prompts/internal/storage"
)

// mockRepository implements storage.Repository for testing.
type mockRepository struct {
	generations []storage.Generation
	categories  []storage.Category
	ratings     map[string]map[string]int // genID -> voterHash -> score
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		generations: []storage.Generation{},
		categories:  storage.DefaultCategories(),
		ratings:     make(map[string]map[string]int),
	}
}

func (m *mockRepository) CreateGeneration(_ context.Context, gen *storage.Generation) error {
	if gen == nil {
		return storage.ErrInvalidInput
	}
	gen.ID = generateID()
	gen.CreatedAt = time.Now()
	m.generations = append(m.generations, *gen)
	return nil
}

func (m *mockRepository) GetGeneration(_ context.Context, id string) (*storage.Generation, error) {
	for i := range m.generations {
		if m.generations[i].ID == id {
			return &m.generations[i], nil
		}
	}
	return nil, storage.ErrNotFound
}

func (m *mockRepository) ListGenerations(_ context.Context, filter storage.ListFilter) ([]storage.Generation, int, error) {
	// Apply category filter
	filtered := []storage.Generation{}
	for _, gen := range m.generations {
		if filter.CategoryID != nil && gen.CategoryID != *filter.CategoryID {
			continue
		}
		filtered = append(filtered, gen)
	}

	total := len(filtered)

	// Apply sorting
	switch filter.SortBy {
	case "highest_rated":
		sort.Slice(filtered, func(i, j int) bool {
			if filtered[i].AvgRating != filtered[j].AvgRating {
				return filtered[i].AvgRating > filtered[j].AvgRating
			}
			return filtered[i].RatingCount > filtered[j].RatingCount
		})
	case "most_viewed":
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].ViewCount > filtered[j].ViewCount
		})
	default: // "newest"
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
		})
	}

	// Apply pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	start := (filter.Page - 1) * filter.PageSize
	if start >= len(filtered) {
		return []storage.Generation{}, total, nil
	}

	end := start + filter.PageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

func (m *mockRepository) IncrementViewCount(_ context.Context, id string) error {
	for i := range m.generations {
		if m.generations[i].ID == id {
			m.generations[i].ViewCount++
			return nil
		}
	}
	return storage.ErrNotFound
}

func (m *mockRepository) RecordView(_ context.Context, generationID string, ipHash string) (bool, error) {
	if generationID == "" || ipHash == "" {
		return false, storage.ErrInvalidInput
	}

	// Check generation exists
	found := false
	genIndex := -1
	for i := range m.generations {
		if m.generations[i].ID == generationID {
			found = true
			genIndex = i
			break
		}
	}
	if !found {
		return false, storage.ErrNotFound
	}

	// Check if this IP has already viewed this generation
	// Use a simple map stored in ratings for simplicity (reusing the structure)
	viewKey := "view:" + generationID + ":" + ipHash
	if m.ratings[viewKey] != nil {
		return false, nil // Already viewed
	}

	// Record the view
	m.ratings[viewKey] = make(map[string]int)
	m.ratings[viewKey]["viewed"] = 1
	m.generations[genIndex].ViewCount++
	return true, nil
}

func (m *mockRepository) CreateOrUpdateRating(_ context.Context, genID string, score int, voterHash string) error {
	if score < 1 || score > 5 {
		return storage.ErrInvalidInput
	}

	// Check generation exists
	found := false
	for i := range m.generations {
		if m.generations[i].ID == genID {
			found = true
			break
		}
	}
	if !found {
		return storage.ErrNotFound
	}

	if m.ratings[genID] == nil {
		m.ratings[genID] = make(map[string]int)
	}
	m.ratings[genID][voterHash] = score

	// Update average rating
	for i := range m.generations {
		if m.generations[i].ID == genID {
			total := 0
			count := 0
			for _, s := range m.ratings[genID] {
				total += s
				count++
			}
			m.generations[i].AvgRating = float64(total) / float64(count)
			m.generations[i].RatingCount = count
			break
		}
	}

	return nil
}

func (m *mockRepository) GetUserRating(_ context.Context, genID string, voterHash string) (int, error) {
	if ratings, ok := m.ratings[genID]; ok {
		if score, ok := ratings[voterHash]; ok {
			return score, nil
		}
	}
	return 0, nil
}

func (m *mockRepository) GetCategoryByKeywords(_ context.Context, text string) (int, error) {
	return storage.MatchCategory(text), nil
}

func (m *mockRepository) GetCategories(_ context.Context) ([]storage.Category, error) {
	return m.categories, nil
}

// Helper functions for generating test data

var idCounter int

func generateID() string {
	idCounter++
	return "gen-" + string(rune('a'+idCounter%26)) + "-" + time.Now().Format("150405.000")
}

func generateRandomGeneration(r *rand.Rand, categoryID int) storage.Generation {
	experienceLevels := []string{"novice", "intermediate", "expert"}
	hookPresets := []string{"default", "minimal", "comprehensive"}

	files := []map[string]string{
		{"path": "kickoff-prompt.md", "content": "test content", "type": "kickoff"},
	}
	filesJSON, _ := json.Marshal(files)

	return storage.Generation{
		ID:              generateID(),
		ProjectIdea:     generateRandomString(r, 10, 100),
		ExperienceLevel: experienceLevels[r.Intn(len(experienceLevels))],
		HookPreset:      hookPresets[r.Intn(len(hookPresets))],
		Files:           filesJSON,
		CategoryID:      categoryID,
		CategoryName:    getCategoryName(categoryID),
		AvgRating:       float64(r.Intn(50)) / 10.0, // 0.0 to 5.0
		RatingCount:     r.Intn(100),
		ViewCount:       r.Intn(1000),
		CreatedAt:       time.Now().Add(-time.Duration(r.Intn(10000)) * time.Minute),
	}
}

func generateRandomString(r *rand.Rand, minLen, maxLen int) string {
	length := minLen + r.Intn(maxLen-minLen+1)
	chars := make([]byte, length)
	for i := range chars {
		chars[i] = byte('a' + r.Intn(26))
	}
	return string(chars)
}

func getCategoryName(id int) string {
	names := map[int]string{1: "API", 2: "CLI", 3: "Web App", 4: "Mobile", 5: "Other"}
	if name, ok := names[id]; ok {
		return name
	}
	return "Other"
}

// Feature: final-polish, Property 5: Gallery Filtering Correctness
// **Validates: Requirements 6.2**
// For any category filter applied to the gallery, all returned items SHALL have
// the matching category_id, and no items with different category_id SHALL be included.

func TestProperty5_GalleryFilteringCorrectness(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := newMockRepository()
		svc := NewService(repo, nil)

		// Generate random generations with various categories
		numGenerations := 10 + r.Intn(50)
		for i := 0; i < numGenerations; i++ {
			categoryID := 1 + r.Intn(5) // 1-5
			gen := generateRandomGeneration(r, categoryID)
			repo.generations = append(repo.generations, gen)
		}

		// Pick a random category to filter by
		filterCategoryID := 1 + r.Intn(5)

		// List with category filter
		resp, err := svc.ListGenerations(context.Background(), ListRequest{
			CategoryID: &filterCategoryID,
			Page:       1,
			PageSize:   100, // Get all
		})
		if err != nil {
			t.Logf("ListGenerations failed: %v", err)
			return false
		}

		// Verify all returned items have the correct category
		for _, item := range resp.Items {
			if item.CategoryID != filterCategoryID {
				t.Logf("Item %s has category %d, expected %d",
					item.ID, item.CategoryID, filterCategoryID)
				return false
			}
		}

		// Verify total count matches actual filtered count
		expectedCount := 0
		for _, gen := range repo.generations {
			if gen.CategoryID == filterCategoryID {
				expectedCount++
			}
		}
		if resp.Total != expectedCount {
			t.Logf("Total count %d doesn't match expected %d", resp.Total, expectedCount)
			return false
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 5 (Gallery Filtering Correctness) failed: %v", err)
	}
}

// TestProperty5_NoFilterReturnsAll tests that without a filter, all items are returned.
func TestProperty5_NoFilterReturnsAll(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := newMockRepository()
		svc := NewService(repo, nil)

		// Generate random generations
		numGenerations := 5 + r.Intn(20)
		for i := 0; i < numGenerations; i++ {
			categoryID := 1 + r.Intn(5)
			gen := generateRandomGeneration(r, categoryID)
			repo.generations = append(repo.generations, gen)
		}

		// List without category filter
		resp, err := svc.ListGenerations(context.Background(), ListRequest{
			Page:     1,
			PageSize: 100,
		})
		if err != nil {
			t.Logf("ListGenerations failed: %v", err)
			return false
		}

		// Verify total matches all generations
		if resp.Total != numGenerations {
			t.Logf("Total %d doesn't match expected %d", resp.Total, numGenerations)
			return false
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 5 (No Filter Returns All) failed: %v", err)
	}
}

// Feature: final-polish, Property 6: Gallery Sorting Correctness
// **Validates: Requirements 6.3**
// For any sort option (newest, highest_rated, most_viewed), the returned items
// SHALL be in descending order by the corresponding field.

func TestProperty6_GallerySortingCorrectness(t *testing.T) {
	sortOptions := []string{"newest", "highest_rated", "most_viewed"}

	for _, sortBy := range sortOptions {
		t.Run(sortBy, func(t *testing.T) {
			property := func(seed int64) bool {
				r := rand.New(rand.NewSource(seed))
				repo := newMockRepository()
				svc := NewService(repo, nil)

				// Generate random generations
				numGenerations := 10 + r.Intn(30)
				for i := 0; i < numGenerations; i++ {
					categoryID := 1 + r.Intn(5)
					gen := generateRandomGeneration(r, categoryID)
					repo.generations = append(repo.generations, gen)
				}

				// List with sort option
				resp, err := svc.ListGenerations(context.Background(), ListRequest{
					SortBy:   sortBy,
					Page:     1,
					PageSize: 100,
				})
				if err != nil {
					t.Logf("ListGenerations failed: %v", err)
					return false
				}

				// Verify items are sorted correctly
				for i := 1; i < len(resp.Items); i++ {
					prev := resp.Items[i-1]
					curr := resp.Items[i]

					switch sortBy {
					case "newest":
						if prev.CreatedAt.Before(curr.CreatedAt) {
							t.Logf("Items not sorted by newest: %v < %v",
								prev.CreatedAt, curr.CreatedAt)
							return false
						}
					case "highest_rated":
						if prev.AvgRating < curr.AvgRating {
							t.Logf("Items not sorted by highest_rated: %v < %v",
								prev.AvgRating, curr.AvgRating)
							return false
						}
					case "most_viewed":
						if prev.ViewCount < curr.ViewCount {
							t.Logf("Items not sorted by most_viewed: %d < %d",
								prev.ViewCount, curr.ViewCount)
							return false
						}
					}
				}

				return true
			}

			cfg := &quick.Config{MaxCount: 100}
			if err := quick.Check(property, cfg); err != nil {
				t.Errorf("Property 6 (Gallery Sorting - %s) failed: %v", sortBy, err)
			}
		})
	}
}

// TestProperty6_DefaultSortIsNewest tests that the default sort is "newest".
func TestProperty6_DefaultSortIsNewest(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	// Add generations with known timestamps
	now := time.Now()
	for i := 0; i < 5; i++ {
		gen := storage.Generation{
			ID:              generateID(),
			ProjectIdea:     "Test project",
			ExperienceLevel: "novice",
			HookPreset:      "default",
			Files:           json.RawMessage(`[]`),
			CategoryID:      1,
			CreatedAt:       now.Add(-time.Duration(i) * time.Hour),
		}
		repo.generations = append(repo.generations, gen)
	}

	// List without specifying sort
	resp, err := svc.ListGenerations(context.Background(), ListRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("ListGenerations failed: %v", err)
	}

	// Verify sorted by newest (descending created_at)
	for i := 1; i < len(resp.Items); i++ {
		if resp.Items[i-1].CreatedAt.Before(resp.Items[i].CreatedAt) {
			t.Errorf("Default sort not by newest: item %d (%v) before item %d (%v)",
				i-1, resp.Items[i-1].CreatedAt, i, resp.Items[i].CreatedAt)
		}
	}
}

// Feature: final-polish, Property 7: Pagination Bounds
// **Validates: Requirements 6.5**
// For any gallery page request, the response SHALL contain at most 20 items,
// and the total_pages calculation SHALL equal ceil(total_items / page_size).

func TestProperty7_PaginationBounds(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := newMockRepository()
		svc := NewService(repo, nil)

		// Generate random number of generations
		numGenerations := r.Intn(100) // 0 to 99
		for i := 0; i < numGenerations; i++ {
			categoryID := 1 + r.Intn(5)
			gen := generateRandomGeneration(r, categoryID)
			repo.generations = append(repo.generations, gen)
		}

		// Random page size (use default if 0)
		pageSize := r.Intn(50) // 0 to 49
		if pageSize == 0 {
			pageSize = DefaultPageSize
		}

		// List first page
		resp, err := svc.ListGenerations(context.Background(), ListRequest{
			Page:     1,
			PageSize: pageSize,
		})
		if err != nil {
			t.Logf("ListGenerations failed: %v", err)
			return false
		}

		// Verify items count doesn't exceed page size
		normalizedPageSize := NormalizePageSize(pageSize)
		if len(resp.Items) > normalizedPageSize {
			t.Logf("Items count %d exceeds page size %d", len(resp.Items), normalizedPageSize)
			return false
		}

		// Verify total pages calculation
		expectedTotalPages := CalculateTotalPages(resp.Total, normalizedPageSize)
		if resp.TotalPages != expectedTotalPages {
			t.Logf("TotalPages %d doesn't match expected %d (total=%d, pageSize=%d)",
				resp.TotalPages, expectedTotalPages, resp.Total, normalizedPageSize)
			return false
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 7 (Pagination Bounds) failed: %v", err)
	}
}

// TestProperty7_PaginationDefaultPageSize tests that default page size is 20.
func TestProperty7_PaginationDefaultPageSize(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	// Add 50 generations
	for i := 0; i < 50; i++ {
		gen := storage.Generation{
			ID:              generateID(),
			ProjectIdea:     "Test project",
			ExperienceLevel: "novice",
			HookPreset:      "default",
			Files:           json.RawMessage(`[]`),
			CategoryID:      1,
			CreatedAt:       time.Now(),
		}
		repo.generations = append(repo.generations, gen)
	}

	// List without specifying page size
	resp, err := svc.ListGenerations(context.Background(), ListRequest{
		Page: 1,
	})
	if err != nil {
		t.Fatalf("ListGenerations failed: %v", err)
	}

	// Verify default page size is 20
	if len(resp.Items) != DefaultPageSize {
		t.Errorf("Default page size: got %d items, expected %d", len(resp.Items), DefaultPageSize)
	}
	if resp.PageSize != DefaultPageSize {
		t.Errorf("PageSize in response: got %d, expected %d", resp.PageSize, DefaultPageSize)
	}
}

// TestProperty7_PaginationMaxPageSize tests that page size is capped at 100.
func TestProperty7_PaginationMaxPageSize(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	// Add 150 generations
	for i := 0; i < 150; i++ {
		gen := storage.Generation{
			ID:              generateID(),
			ProjectIdea:     "Test project",
			ExperienceLevel: "novice",
			HookPreset:      "default",
			Files:           json.RawMessage(`[]`),
			CategoryID:      1,
			CreatedAt:       time.Now(),
		}
		repo.generations = append(repo.generations, gen)
	}

	// Request with page size > 100
	resp, err := svc.ListGenerations(context.Background(), ListRequest{
		Page:     1,
		PageSize: 200,
	})
	if err != nil {
		t.Fatalf("ListGenerations failed: %v", err)
	}

	// Verify page size is capped at 100
	if len(resp.Items) > MaxPageSize {
		t.Errorf("Page size not capped: got %d items, max is %d", len(resp.Items), MaxPageSize)
	}
	if resp.PageSize != MaxPageSize {
		t.Errorf("PageSize in response: got %d, expected %d", resp.PageSize, MaxPageSize)
	}
}

// TestProperty7_PaginationEmptyResult tests pagination with no results.
func TestProperty7_PaginationEmptyResult(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	// List with no generations
	resp, err := svc.ListGenerations(context.Background(), ListRequest{
		Page:     1,
		PageSize: 20,
	})
	if err != nil {
		t.Fatalf("ListGenerations failed: %v", err)
	}

	// Verify empty result
	if len(resp.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(resp.Items))
	}
	if resp.Total != 0 {
		t.Errorf("Expected total 0, got %d", resp.Total)
	}
	// Even with 0 items, total pages should be 1
	if resp.TotalPages != 1 {
		t.Errorf("Expected 1 total page for empty result, got %d", resp.TotalPages)
	}
}

// TestProperty7_PaginationPageBeyondTotal tests requesting a page beyond total.
func TestProperty7_PaginationPageBeyondTotal(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	// Add 5 generations
	for i := 0; i < 5; i++ {
		gen := storage.Generation{
			ID:              generateID(),
			ProjectIdea:     "Test project",
			ExperienceLevel: "novice",
			HookPreset:      "default",
			Files:           json.RawMessage(`[]`),
			CategoryID:      1,
			CreatedAt:       time.Now(),
		}
		repo.generations = append(repo.generations, gen)
	}

	// Request page 10 (beyond total)
	resp, err := svc.ListGenerations(context.Background(), ListRequest{
		Page:     10,
		PageSize: 20,
	})
	if err != nil {
		t.Fatalf("ListGenerations failed: %v", err)
	}

	// Verify empty items but correct total
	if len(resp.Items) != 0 {
		t.Errorf("Expected 0 items for page beyond total, got %d", len(resp.Items))
	}
	if resp.Total != 5 {
		t.Errorf("Expected total 5, got %d", resp.Total)
	}
}

// Additional service tests

func TestService_GetGeneration(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	// Add a generation
	gen := storage.Generation{
		ID:              "test-gen-1",
		ProjectIdea:     "Test project",
		ExperienceLevel: "novice",
		HookPreset:      "default",
		Files:           json.RawMessage(`[]`),
		CategoryID:      1,
		ViewCount:       5,
		CreatedAt:       time.Now(),
	}
	repo.generations = append(repo.generations, gen)

	// Get the generation
	result, err := svc.GetGeneration(context.Background(), "test-gen-1")
	if err != nil {
		t.Fatalf("GetGeneration failed: %v", err)
	}

	if result.ID != "test-gen-1" {
		t.Errorf("Expected ID test-gen-1, got %s", result.ID)
	}

	// View count should be incremented
	if repo.generations[0].ViewCount != 6 {
		t.Errorf("Expected view count 6, got %d", repo.generations[0].ViewCount)
	}
}

func TestService_GetGeneration_NotFound(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	_, err := svc.GetGeneration(context.Background(), "nonexistent")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestService_GetGeneration_EmptyID(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	_, err := svc.GetGeneration(context.Background(), "")
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
}

func TestService_InvalidSortOption(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	_, err := svc.ListGenerations(context.Background(), ListRequest{
		SortBy: "invalid_sort",
	})
	if !errors.Is(err, ErrInvalidSort) {
		t.Errorf("Expected ErrInvalidSort, got %v", err)
	}
}

// Feature: final-polish, Property 8: Rating Storage and Calculation
// **Validates: Requirements 7.2, 7.3**
// For any rating submission, the rating SHALL be stored with the correct generation_id
// and score, and the generation's avg_rating SHALL equal the arithmetic mean of all
// ratings for that generation.

func TestProperty8_RatingStorageAndCalculation(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := newMockRepository()
		svc := NewService(repo, nil)

		// Create a generation to rate
		gen := generateRandomGeneration(r, 1)
		gen.AvgRating = 0
		gen.RatingCount = 0
		repo.generations = append(repo.generations, gen)

		// Generate random ratings from different voters
		numRatings := 1 + r.Intn(20)
		expectedScores := make(map[string]int)

		for i := 0; i < numRatings; i++ {
			voterHash := generateRandomString(r, 10, 20)
			score := 1 + r.Intn(5) // 1-5

			// Submit rating
			_, err := svc.RateGeneration(context.Background(), gen.ID, score, voterHash, "127.0.0.1")
			if err != nil {
				t.Logf("RateGeneration failed: %v", err)
				return false
			}

			expectedScores[voterHash] = score
		}

		// Verify each rating was stored correctly
		for voterHash, expectedScore := range expectedScores {
			actualScore, err := svc.GetUserRating(context.Background(), gen.ID, voterHash)
			if err != nil {
				t.Logf("GetUserRating failed: %v", err)
				return false
			}
			if actualScore != expectedScore {
				t.Logf("Rating mismatch for voter %s: expected %d, got %d",
					voterHash, expectedScore, actualScore)
				return false
			}
		}

		// Calculate expected average
		totalScore := 0
		for _, score := range expectedScores {
			totalScore += score
		}
		expectedAvg := float64(totalScore) / float64(len(expectedScores))

		// Get the generation and verify average
		updatedGen, err := svc.GetGeneration(context.Background(), gen.ID)
		if err != nil {
			t.Logf("GetGeneration failed: %v", err)
			return false
		}

		// Allow small floating point tolerance
		if diff := updatedGen.AvgRating - expectedAvg; diff > 0.01 || diff < -0.01 {
			t.Logf("Average rating mismatch: expected %.2f, got %.2f",
				expectedAvg, updatedGen.AvgRating)
			return false
		}

		// Verify rating count
		if updatedGen.RatingCount != len(expectedScores) {
			t.Logf("Rating count mismatch: expected %d, got %d",
				len(expectedScores), updatedGen.RatingCount)
			return false
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 8 (Rating Storage and Calculation) failed: %v", err)
	}
}

// TestProperty8_RatingValidation tests that invalid ratings are rejected.
func TestProperty8_RatingValidation(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	// Create a generation
	gen := storage.Generation{
		ID:              "test-gen-rating",
		ProjectIdea:     "Test project",
		ExperienceLevel: "novice",
		HookPreset:      "default",
		Files:           json.RawMessage(`[]`),
		CategoryID:      1,
		CreatedAt:       time.Now(),
	}
	repo.generations = append(repo.generations, gen)

	// Test invalid scores
	invalidScores := []int{0, -1, 6, 100, -100}
	for _, score := range invalidScores {
		_, err := svc.RateGeneration(context.Background(), gen.ID, score, "voter1", "127.0.0.1")
		if !errors.Is(err, ErrInvalidRating) {
			t.Errorf("Expected ErrInvalidRating for score %d, got %v", score, err)
		}
	}

	// Test valid scores
	validScores := []int{1, 2, 3, 4, 5}
	for _, score := range validScores {
		voterHash := "voter-" + string(rune('a'+score))
		_, err := svc.RateGeneration(context.Background(), gen.ID, score, voterHash, "127.0.0.1")
		if err != nil {
			t.Errorf("Unexpected error for valid score %d: %v", score, err)
		}
	}
}

// TestProperty8_RatingNonexistentGeneration tests rating a nonexistent generation.
func TestProperty8_RatingNonexistentGeneration(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	_, err := svc.RateGeneration(context.Background(), "nonexistent-id", 5, "voter1", "127.0.0.1")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// Feature: final-polish, Property 9: Duplicate Rating Prevention
// **Validates: Requirements 7.4**
// For any voter_hash that has already rated a generation, subsequent rating submissions
// SHALL update the existing rating rather than create a duplicate, and the rating count
// SHALL remain unchanged.

func TestProperty9_DuplicateRatingPrevention(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := newMockRepository()
		svc := NewService(repo, nil)

		// Create a generation
		gen := generateRandomGeneration(r, 1)
		gen.AvgRating = 0
		gen.RatingCount = 0
		repo.generations = append(repo.generations, gen)

		// Create initial ratings from different voters
		numVoters := 3 + r.Intn(5)
		voterHashes := make([]string, numVoters)
		for i := 0; i < numVoters; i++ {
			voterHashes[i] = generateRandomString(r, 10, 20)
			score := 1 + r.Intn(5)
			_, err := svc.RateGeneration(context.Background(), gen.ID, score, voterHashes[i], "127.0.0.1")
			if err != nil {
				t.Logf("Initial rating failed: %v", err)
				return false
			}
		}

		// Get initial rating count
		updatedGen, err := svc.GetGeneration(context.Background(), gen.ID)
		if err != nil {
			t.Logf("GetGeneration failed: %v", err)
			return false
		}
		initialRatingCount := updatedGen.RatingCount

		// Pick a random voter to update their rating
		voterToUpdate := voterHashes[r.Intn(len(voterHashes))]
		newScore := 1 + r.Intn(5)

		// Submit duplicate rating (update)
		_, err = svc.RateGeneration(context.Background(), gen.ID, newScore, voterToUpdate, "127.0.0.1")
		if err != nil {
			t.Logf("Update rating failed: %v", err)
			return false
		}

		// Verify rating count hasn't changed
		updatedGen, err = svc.GetGeneration(context.Background(), gen.ID)
		if err != nil {
			t.Logf("GetGeneration after update failed: %v", err)
			return false
		}

		if updatedGen.RatingCount != initialRatingCount {
			t.Logf("Rating count changed after duplicate: expected %d, got %d",
				initialRatingCount, updatedGen.RatingCount)
			return false
		}

		// Verify the rating was updated to the new score
		actualScore, err := svc.GetUserRating(context.Background(), gen.ID, voterToUpdate)
		if err != nil {
			t.Logf("GetUserRating failed: %v", err)
			return false
		}

		if actualScore != newScore {
			t.Logf("Rating not updated: expected %d, got %d", newScore, actualScore)
			return false
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 9 (Duplicate Rating Prevention) failed: %v", err)
	}
}

// TestProperty9_MultipleUpdatesFromSameVoter tests multiple updates from the same voter.
func TestProperty9_MultipleUpdatesFromSameVoter(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	// Create a generation
	gen := storage.Generation{
		ID:              "test-gen-dup",
		ProjectIdea:     "Test project",
		ExperienceLevel: "novice",
		HookPreset:      "default",
		Files:           json.RawMessage(`[]`),
		CategoryID:      1,
		CreatedAt:       time.Now(),
	}
	repo.generations = append(repo.generations, gen)

	voterHash := "consistent-voter"

	// Submit multiple ratings from the same voter
	scores := []int{1, 3, 5, 2, 4}
	for _, score := range scores {
		_, err := svc.RateGeneration(context.Background(), gen.ID, score, voterHash, "127.0.0.1")
		if err != nil {
			t.Fatalf("RateGeneration failed: %v", err)
		}
	}

	// Verify only one rating exists (count should be 1)
	updatedGen, err := svc.GetGeneration(context.Background(), gen.ID)
	if err != nil {
		t.Fatalf("GetGeneration failed: %v", err)
	}

	if updatedGen.RatingCount != 1 {
		t.Errorf("Expected rating count 1, got %d", updatedGen.RatingCount)
	}

	// Verify the final score is the last one submitted
	finalScore, err := svc.GetUserRating(context.Background(), gen.ID, voterHash)
	if err != nil {
		t.Fatalf("GetUserRating failed: %v", err)
	}

	expectedFinalScore := scores[len(scores)-1]
	if finalScore != expectedFinalScore {
		t.Errorf("Expected final score %d, got %d", expectedFinalScore, finalScore)
	}

	// Verify average equals the single rating
	if updatedGen.AvgRating != float64(expectedFinalScore) {
		t.Errorf("Expected average %.2f, got %.2f", float64(expectedFinalScore), updatedGen.AvgRating)
	}
}

// Feature: final-polish, Property 10: Rate Limit Enforcement
// **Validates: Requirements 7.6**
// For any IP address that exceeds the rating rate limit (20/hour), subsequent rating
// requests SHALL return ErrRateLimited with a retry-after duration.

func TestProperty10_RateLimitEnforcement(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		repo := newMockRepository()

		// Create a rate limiter with a small limit for testing
		testLimit := 3 + r.Intn(5) // 3-7 requests
		limiter := ratelimit.NewLimiterWithConfig(testLimit, time.Hour)
		svc := NewService(repo, limiter)

		// Create a generation
		gen := generateRandomGeneration(r, 1)
		gen.AvgRating = 0
		gen.RatingCount = 0
		repo.generations = append(repo.generations, gen)

		clientIP := "192.168.1." + string(rune('0'+r.Intn(10)))

		// Make requests up to the limit - should all succeed
		for i := 0; i < testLimit; i++ {
			voterHash := generateRandomString(r, 10, 20)
			score := 1 + r.Intn(5)
			_, err := svc.RateGeneration(context.Background(), gen.ID, score, voterHash, clientIP)
			if err != nil {
				t.Logf("Request %d should have succeeded but got: %v", i+1, err)
				return false
			}
		}

		// Next request should be rate limited
		voterHash := generateRandomString(r, 10, 20)
		retryAfter, err := svc.RateGeneration(context.Background(), gen.ID, 3, voterHash, clientIP)
		if !errors.Is(err, ErrRateLimited) {
			t.Logf("Expected ErrRateLimited after %d requests, got: %v", testLimit, err)
			return false
		}

		// Verify retry-after is positive
		if retryAfter <= 0 {
			t.Logf("Expected positive retry-after, got: %d", retryAfter)
			return false
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 10 (Rate Limit Enforcement) failed: %v", err)
	}
}

// TestProperty10_RateLimitIPIsolation tests that rate limits are per-IP.
func TestProperty10_RateLimitIPIsolation(t *testing.T) {
	repo := newMockRepository()

	// Create a rate limiter with limit of 2
	limiter := ratelimit.NewLimiterWithConfig(2, time.Hour)
	svc := NewService(repo, limiter)

	// Create a generation
	gen := storage.Generation{
		ID:              "test-gen-ratelimit",
		ProjectIdea:     "Test project",
		ExperienceLevel: "novice",
		HookPreset:      "default",
		Files:           json.RawMessage(`[]`),
		CategoryID:      1,
		CreatedAt:       time.Now(),
	}
	repo.generations = append(repo.generations, gen)

	// IP1 makes 2 requests (hits limit)
	for i := 0; i < 2; i++ {
		_, err := svc.RateGeneration(context.Background(), gen.ID, 5, "voter-ip1-"+string(rune('a'+i)), "192.168.1.1")
		if err != nil {
			t.Fatalf("IP1 request %d failed: %v", i+1, err)
		}
	}

	// IP1's next request should be rate limited
	_, err := svc.RateGeneration(context.Background(), gen.ID, 5, "voter-ip1-extra", "192.168.1.1")
	if !errors.Is(err, ErrRateLimited) {
		t.Errorf("Expected IP1 to be rate limited, got: %v", err)
	}

	// IP2 should still be able to make requests
	_, err = svc.RateGeneration(context.Background(), gen.ID, 5, "voter-ip2-a", "192.168.1.2")
	if err != nil {
		t.Errorf("IP2 should not be rate limited, got: %v", err)
	}
}

// TestProperty10_RateLimitWithoutLimiter tests that service works without rate limiter.
func TestProperty10_RateLimitWithoutLimiter(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil) // No rate limiter

	// Create a generation
	gen := storage.Generation{
		ID:              "test-gen-no-limit",
		ProjectIdea:     "Test project",
		ExperienceLevel: "novice",
		HookPreset:      "default",
		Files:           json.RawMessage(`[]`),
		CategoryID:      1,
		CreatedAt:       time.Now(),
	}
	repo.generations = append(repo.generations, gen)

	// Should be able to make many requests without rate limiting
	for i := 0; i < 50; i++ {
		voterHash := "voter-" + string(rune('a'+i%26)) + string(rune('0'+i/26))
		_, err := svc.RateGeneration(context.Background(), gen.ID, 1+i%5, voterHash, "192.168.1.1")
		if err != nil {
			t.Errorf("Request %d failed unexpectedly: %v", i+1, err)
		}
	}
}

// TestProperty10_RatingLimiterConfiguration tests the rating limiter configuration.
func TestProperty10_RatingLimiterConfiguration(t *testing.T) {
	limiter := ratelimit.NewRatingLimiter()

	// Verify the limiter allows 20 requests
	ip := "test-ip"
	for i := 0; i < 20; i++ {
		allowed, _ := limiter.Allow(ip)
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 21st request should be rate limited
	allowed, retryAfter := limiter.Allow(ip)
	if allowed {
		t.Error("21st request should be rate limited")
	}
	if retryAfter <= 0 {
		t.Error("Expected positive retry-after duration")
	}
}
