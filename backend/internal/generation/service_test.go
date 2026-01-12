package generation

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// Feature: ai-driven-generation, Property 1: Question Plan Structure
// **Validates: Requirements 2.2**
// For any valid question generation response, the response SHALL contain
// between 5 and 10 questions, each with a non-empty text field.

// generateValidQuestionsResponse generates a random valid QuestionsResponse.
func generateValidQuestionsResponse(r *rand.Rand) QuestionsResponse {
	// Generate between 5 and 10 questions
	numQuestions := 5 + r.Intn(6) // 5 to 10 inclusive

	questions := make([]Question, numQuestions)
	for i := 0; i < numQuestions; i++ {
		questions[i] = Question{
			ID:   i + 1,
			Text: generateNonEmptyString(r),
			Hint: maybeGenerateString(r),
		}
	}

	return QuestionsResponse{Questions: questions}
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

// maybeGenerateString generates a string or empty string randomly.
func maybeGenerateString(r *rand.Rand) string {
	if r.Intn(2) == 0 {
		return ""
	}
	return generateNonEmptyString(r)
}

// TestProperty1_QuestionPlanStructure tests that parseQuestionsResponse
// correctly validates question plan structure.
// Feature: ai-driven-generation, Property 1: Question Plan Structure
// **Validates: Requirements 2.2**
func TestProperty1_QuestionPlanStructure(t *testing.T) {
	// Property: For any valid JSON response with 5-10 questions with non-empty text,
	// parseQuestionsResponse should return questions with the same structure.
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		validResponse := generateValidQuestionsResponse(r)

		// Serialize to JSON (simulating AI response)
		jsonBytes, err := json.Marshal(validResponse)
		if err != nil {
			return false
		}

		// Parse the response
		questions, err := parseQuestionsResponse(string(jsonBytes))
		if err != nil {
			t.Logf("Parse error: %v", err)
			return false
		}

		// Verify: between 5 and 10 questions
		if len(questions) < minQuestions || len(questions) > maxQuestions {
			t.Logf("Invalid question count: %d", len(questions))
			return false
		}

		// Verify: each question has non-empty text
		for i, q := range questions {
			if q.Text == "" {
				t.Logf("Question %d has empty text", i)
				return false
			}
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 1 (Question Plan Structure) failed: %v", err)
	}
}

// TestProperty1_QuestionPlanStructure_RejectsInvalid tests that invalid responses
// are properly rejected.
func TestProperty1_QuestionPlanStructure_RejectsInvalid(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name:     "empty questions array",
			response: `{"questions": []}`,
			wantErr:  true,
		},
		{
			name:     "question with empty text",
			response: `{"questions": [{"id": 1, "text": ""}, {"id": 2, "text": "valid"}, {"id": 3, "text": "valid"}, {"id": 4, "text": "valid"}, {"id": 5, "text": "valid"}]}`,
			wantErr:  true,
		},
		{
			name:     "invalid JSON",
			response: `not json`,
			wantErr:  true,
		},
		{
			name:     "valid response with 5 questions",
			response: `{"questions": [{"id": 1, "text": "q1"}, {"id": 2, "text": "q2"}, {"id": 3, "text": "q3"}, {"id": 4, "text": "q4"}, {"id": 5, "text": "q5"}]}`,
			wantErr:  false,
		},
		{
			name:     "valid response with 10 questions",
			response: `{"questions": [{"id": 1, "text": "q1"}, {"id": 2, "text": "q2"}, {"id": 3, "text": "q3"}, {"id": 4, "text": "q4"}, {"id": 5, "text": "q5"}, {"id": 6, "text": "q6"}, {"id": 7, "text": "q7"}, {"id": 8, "text": "q8"}, {"id": 9, "text": "q9"}, {"id": 10, "text": "q10"}]}`,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseQuestionsResponse(tc.response)
			if (err != nil) != tc.wantErr {
				t.Errorf("parseQuestionsResponse() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// Generate implements quick.Generator for QuestionsResponse.
func (QuestionsResponse) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(generateValidQuestionsResponse(rand))
}

// Feature: ai-driven-generation, Property 2: Generation Response Completeness
// **Validates: Requirements 3.2, 3.3, 3.4**
// For any valid output generation response, the response SHALL contain
// a non-empty kickoff prompt, at least one steering file, and at least one hook file.

// generateValidOutputsResponse generates a random valid OutputsResponse.
func generateValidOutputsResponse(r *rand.Rand) OutputsResponse {
	files := []GeneratedFile{
		// Always include at least one kickoff file
		{
			Path:    "kickoff-prompt.md",
			Content: generateNonEmptyString(r),
			Type:    "kickoff",
		},
		// Always include at least one steering file
		{
			Path:    ".kiro/steering/product.md",
			Content: generateNonEmptyString(r),
			Type:    "steering",
		},
		// Always include at least one hook file
		{
			Path:    ".kiro/hooks/format.kiro.hook",
			Content: generateNonEmptyString(r),
			Type:    "hook",
		},
		// Always include AGENTS.md
		{
			Path:    "AGENTS.md",
			Content: generateNonEmptyString(r),
			Type:    "agents",
		},
	}

	// Optionally add more steering files
	if r.Intn(2) == 1 {
		files = append(files, GeneratedFile{
			Path:    ".kiro/steering/tech.md",
			Content: generateNonEmptyString(r),
			Type:    "steering",
		})
	}

	// Optionally add more hook files
	if r.Intn(2) == 1 {
		files = append(files, GeneratedFile{
			Path:    ".kiro/hooks/lint.kiro.hook",
			Content: generateNonEmptyString(r),
			Type:    "hook",
		})
	}

	return OutputsResponse{Files: files}
}

// TestProperty2_GenerationResponseCompleteness tests that parseOutputsResponse
// correctly validates response completeness.
// Feature: ai-driven-generation, Property 2: Generation Response Completeness
// **Validates: Requirements 3.2, 3.3, 3.4**
func TestProperty2_GenerationResponseCompleteness(t *testing.T) {
	// Property: For any valid JSON response with kickoff, steering, hook, and agents files,
	// parseOutputsResponse should return files with the same structure.
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		validResponse := generateValidOutputsResponse(r)

		// Serialize to JSON (simulating AI response)
		jsonBytes, err := json.Marshal(validResponse)
		if err != nil {
			return false
		}

		// Parse the response
		files, err := parseOutputsResponse(string(jsonBytes))
		if err != nil {
			t.Logf("Parse error: %v", err)
			return false
		}

		// Verify: has at least one of each type
		hasKickoff := false
		hasSteering := false
		hasHook := false
		hasAgents := false

		for _, f := range files {
			// Verify non-empty content
			if f.Content == "" {
				t.Logf("File %s has empty content", f.Path)
				return false
			}
			if f.Path == "" {
				t.Logf("File has empty path")
				return false
			}

			switch f.Type {
			case "kickoff":
				hasKickoff = true
			case "steering":
				hasSteering = true
			case "hook":
				hasHook = true
			case "agents":
				hasAgents = true
			}
		}

		if !hasKickoff {
			t.Logf("Missing kickoff file")
			return false
		}
		if !hasSteering {
			t.Logf("Missing steering file")
			return false
		}
		if !hasHook {
			t.Logf("Missing hook file")
			return false
		}
		if !hasAgents {
			t.Logf("Missing agents file")
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 2 (Generation Response Completeness) failed: %v", err)
	}
}

// TestProperty2_GenerationResponseCompleteness_RejectsInvalid tests that incomplete
// responses are properly rejected.
func TestProperty2_GenerationResponseCompleteness_RejectsInvalid(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name:     "empty files array",
			response: `{"files": []}`,
			wantErr:  true,
		},
		{
			name:     "missing kickoff",
			response: `{"files": [{"path": ".kiro/steering/product.md", "content": "test", "type": "steering"}, {"path": ".kiro/hooks/format.kiro.hook", "content": "test", "type": "hook"}, {"path": "AGENTS.md", "content": "test", "type": "agents"}]}`,
			wantErr:  true,
		},
		{
			name:     "missing steering",
			response: `{"files": [{"path": "kickoff-prompt.md", "content": "test", "type": "kickoff"}, {"path": ".kiro/hooks/format.kiro.hook", "content": "test", "type": "hook"}, {"path": "AGENTS.md", "content": "test", "type": "agents"}]}`,
			wantErr:  true,
		},
		{
			name:     "missing hook",
			response: `{"files": [{"path": "kickoff-prompt.md", "content": "test", "type": "kickoff"}, {"path": ".kiro/steering/product.md", "content": "test", "type": "steering"}, {"path": "AGENTS.md", "content": "test", "type": "agents"}]}`,
			wantErr:  true,
		},
		{
			name:     "missing agents",
			response: `{"files": [{"path": "kickoff-prompt.md", "content": "test", "type": "kickoff"}, {"path": ".kiro/steering/product.md", "content": "test", "type": "steering"}, {"path": ".kiro/hooks/format.kiro.hook", "content": "test", "type": "hook"}]}`,
			wantErr:  true,
		},
		{
			name:     "file with empty content",
			response: `{"files": [{"path": "kickoff-prompt.md", "content": "", "type": "kickoff"}, {"path": ".kiro/steering/product.md", "content": "test", "type": "steering"}, {"path": ".kiro/hooks/format.kiro.hook", "content": "test", "type": "hook"}, {"path": "AGENTS.md", "content": "test", "type": "agents"}]}`,
			wantErr:  true,
		},
		{
			name:     "file with empty path",
			response: `{"files": [{"path": "", "content": "test", "type": "kickoff"}, {"path": ".kiro/steering/product.md", "content": "test", "type": "steering"}, {"path": ".kiro/hooks/format.kiro.hook", "content": "test", "type": "hook"}, {"path": "AGENTS.md", "content": "test", "type": "agents"}]}`,
			wantErr:  true,
		},
		{
			name:     "valid complete response",
			response: `{"files": [{"path": "kickoff-prompt.md", "content": "# Kickoff", "type": "kickoff"}, {"path": ".kiro/steering/product.md", "content": "# Product", "type": "steering"}, {"path": ".kiro/hooks/format.kiro.hook", "content": "{}", "type": "hook"}, {"path": "AGENTS.md", "content": "# Agents", "type": "agents"}]}`,
			wantErr:  false,
		},
		{
			name:     "invalid JSON",
			response: `not json`,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseOutputsResponse(tc.response)
			if (err != nil) != tc.wantErr {
				t.Errorf("parseOutputsResponse() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// Generate implements quick.Generator for OutputsResponse.
func (OutputsResponse) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(generateValidOutputsResponse(rand))
}
