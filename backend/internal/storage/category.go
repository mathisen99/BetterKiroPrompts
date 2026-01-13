package storage

import (
	"context"
	"strings"
)

// CategoryMatcher provides keyword-based category detection.
type CategoryMatcher struct {
	categories []Category
}

// NewCategoryMatcher creates a new CategoryMatcher with the given categories.
func NewCategoryMatcher(categories []Category) *CategoryMatcher {
	return &CategoryMatcher{categories: categories}
}

// DefaultCategories returns the default category definitions.
// These match the database seed data in the migrations.
func DefaultCategories() []Category {
	return []Category{
		{ID: 1, Name: "API", Keywords: []string{"api", "rest", "graphql", "endpoint", "backend", "server"}},
		{ID: 2, Name: "CLI", Keywords: []string{"cli", "command", "terminal", "shell", "script", "console"}},
		{ID: 3, Name: "Web App", Keywords: []string{"web", "frontend", "react", "vue", "angular", "website", "webapp"}},
		{ID: 4, Name: "Mobile", Keywords: []string{"mobile", "ios", "android", "react native", "flutter", "app"}},
		{ID: 5, Name: "Other", Keywords: []string{}},
	}
}

// Match finds the best matching category for the given text.
// Categories are checked in priority order: API > CLI > Web App > Mobile > Other.
// Returns the category ID of the first match, or 5 (Other) if no match.
func (m *CategoryMatcher) Match(text string) int {
	if text == "" {
		return 5 // Other
	}

	lowerText := strings.ToLower(text)

	// Check categories in priority order (by ID: 1=API, 2=CLI, 3=Web App, 4=Mobile)
	for _, cat := range m.categories {
		if cat.ID == 5 { // Skip "Other" - it's the fallback
			continue
		}
		for _, keyword := range cat.Keywords {
			if containsWord(lowerText, strings.ToLower(keyword)) {
				return cat.ID
			}
		}
	}

	return 5 // Other (default)
}

// containsWord checks if the text contains the keyword as a word or phrase.
// This handles multi-word keywords like "react native" and ensures
// partial matches don't trigger (e.g., "application" shouldn't match "app").
func containsWord(text, keyword string) bool {
	// For multi-word keywords, just check if the phrase exists
	if strings.Contains(keyword, " ") {
		return strings.Contains(text, keyword)
	}

	// For single-word keywords, check word boundaries
	// Split text into words and check for exact match
	words := strings.Fields(text)
	for _, word := range words {
		// Clean punctuation from word
		cleanWord := strings.Trim(word, ".,;:!?()[]{}\"'")
		if cleanWord == keyword {
			return true
		}
	}

	// Also check if keyword appears as a substring with word boundaries
	// This handles cases like "api-server" or "rest_api"
	idx := strings.Index(text, keyword)
	for idx != -1 {
		// Check if it's at a word boundary
		atStart := idx == 0 || !isAlphaNum(text[idx-1])
		atEnd := idx+len(keyword) >= len(text) || !isAlphaNum(text[idx+len(keyword)])

		if atStart && atEnd {
			return true
		}

		// Look for next occurrence
		nextIdx := strings.Index(text[idx+1:], keyword)
		if nextIdx == -1 {
			break
		}
		idx = idx + 1 + nextIdx
	}

	return false
}

// isAlphaNum returns true if the byte is alphanumeric.
func isAlphaNum(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// GetCategoryByKeywords implements Repository.GetCategoryByKeywords using the matcher.
func (r *PostgresRepository) GetCategoryByKeywords(ctx context.Context, text string) (int, error) {
	// Load categories from database
	categories, err := r.GetCategories(ctx)
	if err != nil {
		// Fall back to default categories if database is unavailable
		categories = DefaultCategories()
	}

	matcher := NewCategoryMatcher(categories)
	return matcher.Match(text), nil
}

// MatchCategory is a convenience function that matches text against default categories.
// Use this when you don't have database access.
func MatchCategory(text string) int {
	matcher := NewCategoryMatcher(DefaultCategories())
	return matcher.Match(text)
}
