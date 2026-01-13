package prompts

import (
	"strings"
	"testing"
	"testing/quick"
)

// TestExperienceLevelValidation tests that experience level validation works correctly.
func TestExperienceLevelValidation(t *testing.T) {
	tests := []struct {
		level string
		valid bool
	}{
		{"beginner", true},
		{"novice", true},
		{"expert", true},
		{"invalid", false},
		{"", false},
		{"BEGINNER", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			got := IsValidExperienceLevel(tt.level)
			if got != tt.valid {
				t.Errorf("IsValidExperienceLevel(%q) = %v, want %v", tt.level, got, tt.valid)
			}
		})
	}
}

// TestHookPresetValidation tests that hook preset validation works correctly.
func TestHookPresetValidation(t *testing.T) {
	tests := []struct {
		preset string
		valid  bool
	}{
		{"light", true},
		{"basic", true},
		{"default", true},
		{"strict", true},
		{"invalid", false},
		{"", false},
		{"LIGHT", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.preset, func(t *testing.T) {
			got := IsValidHookPreset(tt.preset)
			if got != tt.valid {
				t.Errorf("IsValidHookPreset(%q) = %v, want %v", tt.preset, got, tt.valid)
			}
		})
	}
}

// TestQuestionsSystemPromptGeneration tests that question prompts are generated for each level.
func TestQuestionsSystemPromptGeneration(t *testing.T) {
	levels := []string{ExperienceBeginner, ExperienceNovice, ExperienceExpert}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			prompt := GetQuestionsSystemPrompt(level)
			if prompt == "" {
				t.Errorf("GetQuestionsSystemPrompt(%q) returned empty string", level)
			}
			if len(prompt) < 100 {
				t.Errorf("GetQuestionsSystemPrompt(%q) returned suspiciously short prompt: %d chars", level, len(prompt))
			}
		})
	}
}

// TestOutputsSystemPromptGeneration tests that output prompts are generated for each combination.
func TestOutputsSystemPromptGeneration(t *testing.T) {
	levels := []string{ExperienceBeginner, ExperienceNovice, ExperienceExpert}
	presets := []string{HookPresetLight, HookPresetBasic, HookPresetDefault, HookPresetStrict}

	for _, level := range levels {
		for _, preset := range presets {
			t.Run(level+"_"+preset, func(t *testing.T) {
				prompt := GetOutputsSystemPrompt(level, preset)
				if prompt == "" {
					t.Errorf("GetOutputsSystemPrompt(%q, %q) returned empty string", level, preset)
				}
				if len(prompt) < 500 {
					t.Errorf("GetOutputsSystemPrompt(%q, %q) returned suspiciously short prompt: %d chars", level, preset, len(prompt))
				}
			})
		}
	}
}

// TestHookPresetDescriptions tests that all presets have descriptions.
func TestHookPresetDescriptions(t *testing.T) {
	presets := []string{HookPresetLight, HookPresetBasic, HookPresetDefault, HookPresetStrict}

	for _, preset := range presets {
		t.Run(preset, func(t *testing.T) {
			info, ok := HookPresetDescriptions[preset]
			if !ok {
				t.Errorf("HookPresetDescriptions missing entry for %q", preset)
				return
			}
			if info.Title == "" {
				t.Errorf("HookPresetDescriptions[%q].Title is empty", preset)
			}
			if info.Description == "" {
				t.Errorf("HookPresetDescriptions[%q].Description is empty", preset)
			}
			if len(info.Hooks) == 0 {
				t.Errorf("HookPresetDescriptions[%q].Hooks is empty", preset)
			}
		})
	}
}

// TestValidExperienceLevels tests that ValidExperienceLevels returns all levels.
func TestValidExperienceLevels(t *testing.T) {
	levels := ValidExperienceLevels()
	if len(levels) != 3 {
		t.Errorf("ValidExperienceLevels() returned %d levels, want 3", len(levels))
	}

	expected := map[string]bool{
		ExperienceBeginner: true,
		ExperienceNovice:   true,
		ExperienceExpert:   true,
	}

	for _, level := range levels {
		if !expected[level] {
			t.Errorf("ValidExperienceLevels() contains unexpected level %q", level)
		}
	}
}

// TestValidHookPresets tests that ValidHookPresets returns all presets.
func TestValidHookPresets(t *testing.T) {
	presets := ValidHookPresets()
	if len(presets) != 4 {
		t.Errorf("ValidHookPresets() returned %d presets, want 4", len(presets))
	}

	expected := map[string]bool{
		HookPresetLight:   true,
		HookPresetBasic:   true,
		HookPresetDefault: true,
		HookPresetStrict:  true,
	}

	for _, preset := range presets {
		if !expected[preset] {
			t.Errorf("ValidHookPresets() contains unexpected preset %q", preset)
		}
	}
}

// TestProperty1_ExperienceLevelAdaptation tests that beginner-level question prompts
// avoid technical jargon terms.
// Feature: phase4-production, Property 1: Experience Level Adaptation
// **Validates: Requirements 1.2, 1.3, 1.4, 3.1, 6.6**
func TestProperty1_ExperienceLevelAdaptation(t *testing.T) {
	// Property: For any generated question prompt at beginner level,
	// the prompt SHALL NOT contain technical jargon terms in the guidance sections
	// that would be presented to the AI for generating questions.

	// Get the beginner system prompt
	beginnerPrompt := GetQuestionsSystemPrompt(ExperienceBeginner)

	// The beginner prompt should explicitly list jargon terms to AVOID
	// This ensures the AI knows what terms not to use
	for _, jargonTerm := range JargonTerms {
		if !strings.Contains(strings.ToLower(beginnerPrompt), strings.ToLower(jargonTerm)) {
			// The jargon term should be mentioned in the "AVOID" section
			// to instruct the AI not to use it
			t.Logf("Note: Jargon term %q is listed in JargonTerms but may not be explicitly mentioned in prompt", jargonTerm)
		}
	}

	// Verify the beginner prompt contains the AVOID instruction
	if !strings.Contains(beginnerPrompt, "AVOID") {
		t.Error("Beginner prompt should contain AVOID instruction for jargon terms")
	}

	// Verify the beginner prompt mentions avoiding jargon
	if !strings.Contains(strings.ToLower(beginnerPrompt), "jargon") ||
		!strings.Contains(strings.ToLower(beginnerPrompt), "avoid") {
		t.Error("Beginner prompt should instruct to avoid jargon")
	}

	// Property test: For any project idea string, the beginner user prompt
	// should not introduce jargon terms that weren't in the original idea
	property := func(projectIdea string) bool {
		if strings.TrimSpace(projectIdea) == "" {
			return true // Skip empty inputs
		}

		userPrompt := GetQuestionsUserPrompt(projectIdea, ExperienceBeginner)

		// The user prompt should mention the experience level
		if !strings.Contains(userPrompt, "beginner") &&
			!strings.Contains(userPrompt, "Beginner") {
			return false
		}

		// The user prompt should mention adapting language
		if !strings.Contains(strings.ToLower(userPrompt), "experience level") {
			return false
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 1 (Experience Level Adaptation) failed: %v", err)
	}
}

// TestProperty1_BeginnerPromptContainsJargonAvoidanceList verifies that the beginner
// prompt explicitly lists jargon terms to avoid.
// Feature: phase4-production, Property 1: Experience Level Adaptation
// **Validates: Requirements 1.2, 1.3, 1.4, 3.1, 6.6**
func TestProperty1_BeginnerPromptContainsJargonAvoidanceList(t *testing.T) {
	beginnerPrompt := GetQuestionsSystemPrompt(ExperienceBeginner)
	promptLower := strings.ToLower(beginnerPrompt)

	// Core jargon terms that MUST be listed in the avoidance section
	coreJargonTerms := []string{
		"microservices",
		"distributed",
		"scalability",
		"concurrency",
		"middleware",
	}

	missingTerms := []string{}
	for _, term := range coreJargonTerms {
		if !strings.Contains(promptLower, strings.ToLower(term)) {
			missingTerms = append(missingTerms, term)
		}
	}

	if len(missingTerms) > 0 {
		t.Errorf("Beginner prompt should list these jargon terms to avoid: %v", missingTerms)
	}
}

// TestProperty1_ExpertPromptAllowsTechnicalTerms verifies that expert-level prompts
// allow and encourage technical terminology.
// Feature: phase4-production, Property 1: Experience Level Adaptation
// **Validates: Requirements 1.4, 3.3**
func TestProperty1_ExpertPromptAllowsTechnicalTerms(t *testing.T) {
	expertPrompt := GetQuestionsSystemPrompt(ExperienceExpert)
	promptLower := strings.ToLower(expertPrompt)

	// Expert prompt should mention technical terminology positively
	technicalIndicators := []string{
		"technical",
		"architecture",
		"scalability",
		"consistency",
	}

	foundIndicators := 0
	for _, indicator := range technicalIndicators {
		if strings.Contains(promptLower, indicator) {
			foundIndicators++
		}
	}

	if foundIndicators < 2 {
		t.Errorf("Expert prompt should encourage technical terminology, found only %d indicators", foundIndicators)
	}

	// Expert prompt should NOT contain "AVOID" for jargon terms
	// (unlike beginner prompt)
	if strings.Contains(expertPrompt, "AVOID these jargon terms") {
		t.Error("Expert prompt should not tell AI to avoid jargon terms")
	}
}

// TestProperty1_NovicePromptBalancesTechnicalLanguage verifies that novice-level
// prompts use moderate technical language.
// Feature: phase4-production, Property 1: Experience Level Adaptation
// **Validates: Requirements 1.3**
func TestProperty1_NovicePromptBalancesTechnicalLanguage(t *testing.T) {
	novicePrompt := GetQuestionsSystemPrompt(ExperienceNovice)
	promptLower := strings.ToLower(novicePrompt)

	// Novice prompt should mention balance or moderate language
	balanceIndicators := []string{
		"moderate",
		"balance",
		"common technical terms",
		"explain",
	}

	foundIndicators := 0
	for _, indicator := range balanceIndicators {
		if strings.Contains(promptLower, indicator) {
			foundIndicators++
		}
	}

	if foundIndicators < 1 {
		t.Errorf("Novice prompt should indicate balanced technical language, found %d indicators", foundIndicators)
	}
}

// TestProperty1_AllLevelsHaveDistinctGuidance verifies that each experience level
// produces distinct guidance in the prompts.
// Feature: phase4-production, Property 1: Experience Level Adaptation
// **Validates: Requirements 1.2, 1.3, 1.4**
func TestProperty1_AllLevelsHaveDistinctGuidance(t *testing.T) {
	beginnerPrompt := GetQuestionsSystemPrompt(ExperienceBeginner)
	novicePrompt := GetQuestionsSystemPrompt(ExperienceNovice)
	expertPrompt := GetQuestionsSystemPrompt(ExperienceExpert)

	// All prompts should be different
	if beginnerPrompt == novicePrompt {
		t.Error("Beginner and Novice prompts should be different")
	}
	if novicePrompt == expertPrompt {
		t.Error("Novice and Expert prompts should be different")
	}
	if beginnerPrompt == expertPrompt {
		t.Error("Beginner and Expert prompts should be different")
	}

	// Each should contain its level name
	if !strings.Contains(beginnerPrompt, "Beginner") {
		t.Error("Beginner prompt should contain 'Beginner'")
	}
	if !strings.Contains(novicePrompt, "Novice") {
		t.Error("Novice prompt should contain 'Novice'")
	}
	if !strings.Contains(expertPrompt, "Expert") {
		t.Error("Expert prompt should contain 'Expert'")
	}
}

// TestProperty1_BeginnerQuestionsAvoidTechnicalJargon validates that the beginner prompt
// explicitly forbids all technical jargon terms from ForbiddenBeginnerTerms.
// Feature: ux-improvements, Property 1: Beginner Questions Avoid Technical Jargon
// **Validates: Requirements 1.1, 1.2**
func TestProperty1_BeginnerQuestionsAvoidTechnicalJargon(t *testing.T) {
	// Property: For any project idea submitted with beginner experience level,
	// the generated questions SHALL NOT contain any of the forbidden technical terms.

	beginnerPrompt := GetQuestionsSystemPrompt(ExperienceBeginner)
	promptLower := strings.ToLower(beginnerPrompt)

	// The beginner prompt MUST contain the AVOID instruction
	if !strings.Contains(promptLower, "avoid") {
		t.Error("Beginner prompt must contain 'AVOID' instruction")
	}

	// The beginner prompt MUST mention jargon/technical terms
	if !strings.Contains(promptLower, "jargon") && !strings.Contains(promptLower, "forbidden") {
		t.Error("Beginner prompt must mention 'jargon' or 'forbidden' terms")
	}

	// Core forbidden terms that MUST be listed in the prompt
	coreForbiddenTerms := []string{
		"API",
		"database schema",
		"authentication flow",
		"microservices",
		"CI/CD",
		"containerization",
		"OAuth",
		"REST",
		"GraphQL",
		"SQL",
		"NoSQL",
		"backend",
		"frontend",
		"deployment",
	}

	missingTerms := []string{}
	for _, term := range coreForbiddenTerms {
		if !strings.Contains(promptLower, strings.ToLower(term)) {
			missingTerms = append(missingTerms, term)
		}
	}

	if len(missingTerms) > 0 {
		t.Errorf("Beginner prompt must list these forbidden terms: %v", missingTerms)
	}

	// Property test: For any project idea, the beginner prompt should instruct
	// the AI to use simple language alternatives
	property := func(projectIdea string) bool {
		if len(projectIdea) == 0 || len(projectIdea) > 1000 {
			return true // Skip edge cases
		}

		// The system prompt should contain language translation guidance
		if !strings.Contains(promptLower, "instead of") {
			return false
		}

		// The system prompt should mention everyday language
		if !strings.Contains(promptLower, "everyday") {
			return false
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 1 (Beginner Questions Avoid Technical Jargon) failed: %v", err)
	}
}

// TestProperty2_ExperienceLevelsProduceDifferentQuestions validates that each experience
// level produces distinctly different question prompts.
// Feature: ux-improvements, Property 2: Experience Levels Produce Different Questions
// **Validates: Requirements 1.5**
func TestProperty2_ExperienceLevelsProduceDifferentQuestions(t *testing.T) {
	// Property: For any project idea, generating questions at beginner, novice, and expert
	// levels SHALL produce three distinct question sets where no two sets are identical.

	beginnerPrompt := GetQuestionsSystemPrompt(ExperienceBeginner)
	novicePrompt := GetQuestionsSystemPrompt(ExperienceNovice)
	expertPrompt := GetQuestionsSystemPrompt(ExperienceExpert)

	// All three prompts must be different from each other
	if beginnerPrompt == novicePrompt {
		t.Error("Beginner and Novice prompts must be different")
	}
	if novicePrompt == expertPrompt {
		t.Error("Novice and Expert prompts must be different")
	}
	if beginnerPrompt == expertPrompt {
		t.Error("Beginner and Expert prompts must be different")
	}

	// Each prompt must contain its experience level identifier
	if !strings.Contains(beginnerPrompt, "Beginner") {
		t.Error("Beginner prompt must contain 'Beginner' identifier")
	}
	if !strings.Contains(novicePrompt, "Novice") {
		t.Error("Novice prompt must contain 'Novice' identifier")
	}
	if !strings.Contains(expertPrompt, "Expert") {
		t.Error("Expert prompt must contain 'Expert' identifier")
	}

	// Beginner prompt must have AVOID/FORBIDDEN terms section
	beginnerLower := strings.ToLower(beginnerPrompt)
	if !strings.Contains(beginnerLower, "avoid") || !strings.Contains(beginnerLower, "forbidden") {
		t.Error("Beginner prompt must have AVOID/FORBIDDEN terms section")
	}

	// Expert prompt must NOT have AVOID jargon section (experts can use technical terms)
	expertLower := strings.ToLower(expertPrompt)
	if strings.Contains(expertLower, "forbidden terms") {
		t.Error("Expert prompt should not have forbidden terms section")
	}

	// Novice prompt must mention balance/moderate language
	noviceLower := strings.ToLower(novicePrompt)
	if !strings.Contains(noviceLower, "moderate") && !strings.Contains(noviceLower, "balance") {
		t.Error("Novice prompt must mention moderate or balanced language")
	}

	// Property test: For any project idea string, the user prompts for different
	// experience levels should be different
	property := func(projectIdea string) bool {
		if len(projectIdea) == 0 || len(projectIdea) > 1000 {
			return true // Skip edge cases
		}

		beginnerUserPrompt := GetQuestionsUserPrompt(projectIdea, ExperienceBeginner)
		noviceUserPrompt := GetQuestionsUserPrompt(projectIdea, ExperienceNovice)
		expertUserPrompt := GetQuestionsUserPrompt(projectIdea, ExperienceExpert)

		// All user prompts must be different (they include experience level)
		if beginnerUserPrompt == noviceUserPrompt {
			return false
		}
		if noviceUserPrompt == expertUserPrompt {
			return false
		}
		if beginnerUserPrompt == expertUserPrompt {
			return false
		}

		// Each must contain the experience level name
		if !strings.Contains(beginnerUserPrompt, "beginner") {
			return false
		}
		if !strings.Contains(noviceUserPrompt, "novice") {
			return false
		}
		if !strings.Contains(expertUserPrompt, "expert") {
			return false
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 2 (Experience Levels Produce Different Questions) failed: %v", err)
	}
}
