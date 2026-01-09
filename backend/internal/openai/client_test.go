package openai

import (
	"strings"
	"testing"
	"testing/quick"
	"unicode"
)

// Property 7: Input Validation
// For any request with empty or whitespace-only projectIdea, the Generation_Service
// SHALL reject the request with HTTP 400 before calling OpenAI API.
// Validates: Requirements 6.5
//
// Feature: ai-driven-generation, Property 7: Input Validation

// generateWhitespaceString generates strings containing only whitespace characters.
func generateWhitespaceString(length int) string {
	whitespaceChars := []rune{' ', '\t', '\n', '\r', '\v', '\f'}
	result := make([]rune, length)
	for i := range result {
		result[i] = whitespaceChars[i%len(whitespaceChars)]
	}
	return string(result)
}

// isWhitespaceOnly checks if a string contains only whitespace characters.
func isWhitespaceOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// TestValidateInput_Property_EmptyStringRejected tests that empty strings are always rejected.
func TestValidateInput_Property_EmptyStringRejected(t *testing.T) {
	err := ValidateInput("")
	if err != ErrEmptyInput {
		t.Errorf("expected ErrEmptyInput for empty string, got %v", err)
	}
}

// TestValidateInput_Property_WhitespaceOnlyRejected tests that whitespace-only strings are rejected.
// Property: For any string composed entirely of whitespace characters, ValidateInput SHALL return ErrEmptyInput.
func TestValidateInput_Property_WhitespaceOnlyRejected(t *testing.T) {
	// Test various lengths of whitespace-only strings
	property := func(length uint8) bool {
		// Generate whitespace string of given length (0-255)
		ws := generateWhitespaceString(int(length))
		err := ValidateInput(ws)
		return err == ErrEmptyInput
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: whitespace-only strings should be rejected: %v", err)
	}
}

// TestValidateInput_Property_NonWhitespaceAccepted tests that strings with non-whitespace content are accepted.
// Property: For any string containing at least one non-whitespace character, ValidateInput SHALL return nil.
func TestValidateInput_Property_NonWhitespaceAccepted(t *testing.T) {
	property := func(s string) bool {
		// Skip if string is empty or whitespace-only (those should be rejected)
		if strings.TrimSpace(s) == "" {
			return true // Skip this case, it's covered by other tests
		}
		err := ValidateInput(s)
		return err == nil
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: non-whitespace strings should be accepted: %v", err)
	}
}

// TestValidateInput_Property_TrimmedEquivalence tests that validation is consistent with TrimSpace.
// Property: For any string s, ValidateInput(s) returns ErrEmptyInput if and only if strings.TrimSpace(s) == "".
func TestValidateInput_Property_TrimmedEquivalence(t *testing.T) {
	property := func(s string) bool {
		err := ValidateInput(s)
		trimmed := strings.TrimSpace(s)

		if trimmed == "" {
			return err == ErrEmptyInput
		}
		return err == nil
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: validation should be equivalent to TrimSpace check: %v", err)
	}
}
