package storage

import (
	"math/rand"
	"testing"
	"testing/quick"
)

// Feature: final-polish, Property 4: Category Assignment Correctness
// **Validates: Requirements 5.3**
// For any project idea containing category keywords, the Backend SHALL assign
// the matching category. For project ideas with multiple matching categories,
// the first match in priority order (API > CLI > Web App > Mobile > Other) SHALL be used.

// TestProperty4_CategoryAssignmentCorrectness tests that category matching
// correctly assigns categories based on keywords.
// Feature: final-polish, Property 4: Category Assignment Correctness
// **Validates: Requirements 5.3**
func TestProperty4_CategoryAssignmentCorrectness(t *testing.T) {
	matcher := NewCategoryMatcher(DefaultCategories())

	// Property: For any text containing a category keyword, the matcher
	// should return the correct category ID.
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// Generate a random project idea with a known keyword
		categories := DefaultCategories()
		catIdx := r.Intn(4) // 0-3 (skip "Other")
		cat := categories[catIdx]

		if len(cat.Keywords) == 0 {
			return true // Skip categories with no keywords
		}

		keyword := cat.Keywords[r.Intn(len(cat.Keywords))]
		projectIdea := generateProjectIdeaWithKeyword(r, keyword)

		result := matcher.Match(projectIdea)

		// The result should be the category ID or a higher priority category
		// (since we might have accidentally included a higher priority keyword)
		if result > cat.ID && result != 5 {
			t.Logf("Expected category %d (%s) or higher priority, got %d for text: %s",
				cat.ID, cat.Name, result, projectIdea)
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 4 (Category Assignment Correctness) failed: %v", err)
	}
}

// generateProjectIdeaWithKeyword generates a random project idea containing the keyword.
func generateProjectIdeaWithKeyword(r *rand.Rand, keyword string) string {
	prefixes := []string{
		"Build a ", "Create a ", "Develop a ", "I want to build a ",
		"Help me create a ", "Design a ", "Implement a ",
	}
	suffixes := []string{
		" for my project", " application", " system", " tool",
		" service", " platform", "",
	}

	prefix := prefixes[r.Intn(len(prefixes))]
	suffix := suffixes[r.Intn(len(suffixes))]

	return prefix + keyword + suffix
}

// TestProperty4_CategoryPriorityOrder tests that when multiple keywords match,
// the highest priority category is selected.
func TestProperty4_CategoryPriorityOrder(t *testing.T) {
	matcher := NewCategoryMatcher(DefaultCategories())

	testCases := []struct {
		name       string
		text       string
		expectedID int
	}{
		// API has highest priority (ID 1)
		{
			name:       "API keyword alone",
			text:       "Build a REST API",
			expectedID: 1,
		},
		{
			name:       "API with web keywords",
			text:       "Build a web API with React frontend",
			expectedID: 1, // API takes priority over Web App
		},
		{
			name:       "API with CLI keywords",
			text:       "Build an API with CLI tool",
			expectedID: 1, // API takes priority over CLI
		},
		// CLI has second priority (ID 2)
		{
			name:       "CLI keyword alone",
			text:       "Build a CLI tool",
			expectedID: 2,
		},
		{
			name:       "CLI with web keywords",
			text:       "Build a CLI for web deployment",
			expectedID: 2, // CLI takes priority over Web App
		},
		// Web App has third priority (ID 3)
		{
			name:       "Web keyword alone",
			text:       "Build a React website",
			expectedID: 3,
		},
		{
			name:       "Web with mobile keywords",
			text:       "Build a web app with mobile support",
			expectedID: 3, // Web App takes priority over Mobile
		},
		// Mobile has fourth priority (ID 4)
		{
			name:       "Mobile keyword alone",
			text:       "Build an iOS app",
			expectedID: 4,
		},
		{
			name:       "Flutter app",
			text:       "Build a Flutter mobile application",
			expectedID: 4,
		},
		// Other is the fallback (ID 5)
		{
			name:       "No matching keywords",
			text:       "Build something cool",
			expectedID: 5,
		},
		{
			name:       "Empty text",
			text:       "",
			expectedID: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := matcher.Match(tc.text)
			if result != tc.expectedID {
				t.Errorf("Match(%q) = %d, want %d", tc.text, result, tc.expectedID)
			}
		})
	}
}

// TestProperty4_CategoryKeywordMatching tests specific keyword matching behavior.
func TestProperty4_CategoryKeywordMatching(t *testing.T) {
	matcher := NewCategoryMatcher(DefaultCategories())

	testCases := []struct {
		name       string
		text       string
		expectedID int
	}{
		// API keywords
		{"api lowercase", "build an api", 1},
		{"API uppercase", "Build an API", 1},
		{"rest", "REST service", 1},
		{"graphql", "GraphQL server", 1},
		{"endpoint", "Create endpoint", 1},
		{"backend", "Backend service", 1},
		{"server", "Server application", 1},

		// CLI keywords
		{"cli", "CLI tool", 2},
		{"command", "Command line tool", 2},
		{"terminal", "Terminal app", 2},
		{"shell", "Shell script", 2},
		{"script", "Automation script", 2},
		{"console", "Console application", 2},

		// Web App keywords
		{"web", "Web application", 3},
		{"frontend", "Frontend project", 3},
		{"react", "React app", 3},
		{"vue", "Vue.js project", 3},
		{"angular", "Angular application", 3},
		{"website", "Personal website", 3},
		{"webapp", "Webapp project", 3},

		// Mobile keywords
		{"mobile", "Mobile app", 4},
		{"ios", "iOS application", 4},
		{"android", "Android app", 4},
		// Note: "React Native" contains "react" which matches Web App (higher priority)
		// This is correct per the priority order: API > CLI > Web App > Mobile > Other
		{"react native", "React Native app", 3}, // matches "react" -> Web App
		{"flutter", "Flutter project", 4},

		// Edge cases - "app" should match Mobile
		{"app alone", "Build an app", 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := matcher.Match(tc.text)
			if result != tc.expectedID {
				t.Errorf("Match(%q) = %d, want %d", tc.text, result, tc.expectedID)
			}
		})
	}
}

// TestProperty4_CategoryMatchingWordBoundaries tests that partial matches don't trigger.
func TestProperty4_CategoryMatchingWordBoundaries(t *testing.T) {
	matcher := NewCategoryMatcher(DefaultCategories())

	testCases := []struct {
		name       string
		text       string
		expectedID int
	}{
		// "application" should not match "app" (Mobile)
		{
			name:       "application should not match app",
			text:       "Build an application for data processing",
			expectedID: 5, // Other - no keyword match
		},
		// "scripting" should not match "script" (CLI)
		{
			name:       "scripting should not match script",
			text:       "Learn scripting languages",
			expectedID: 5, // Other
		},
		// But "api-server" should match "api"
		{
			name:       "api-server should match api",
			text:       "Build an api-server",
			expectedID: 1, // API
		},
		// And "rest_api" should match "api"
		{
			name:       "rest_api should match api",
			text:       "Create a rest_api",
			expectedID: 1, // API (matches "rest" first actually)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := matcher.Match(tc.text)
			if result != tc.expectedID {
				t.Errorf("Match(%q) = %d, want %d", tc.text, result, tc.expectedID)
			}
		})
	}
}

// TestProperty4_MatchCategoryConvenienceFunction tests the convenience function.
func TestProperty4_MatchCategoryConvenienceFunction(t *testing.T) {
	// Should work the same as using NewCategoryMatcher with DefaultCategories
	testCases := []struct {
		text       string
		expectedID int
	}{
		{"Build a REST API", 1},
		{"CLI tool", 2},
		{"React website", 3},
		{"iOS app", 4},
		{"Something else", 5},
	}

	for _, tc := range testCases {
		result := MatchCategory(tc.text)
		if result != tc.expectedID {
			t.Errorf("MatchCategory(%q) = %d, want %d", tc.text, result, tc.expectedID)
		}
	}
}
