package sanitize

import (
	"html"
	"math/rand"
	"strings"
	"testing"
	"testing/quick"
)

// Feature: final-polish, Property 13: Input Sanitization
// **Validates: Requirements 9.2, 9.3**
// For any user input containing HTML tags or JavaScript, the sanitized output
// SHALL not contain executable code, and the sanitized content SHALL be safe
// for database storage and HTML rendering.

// TestProperty13_InputSanitization tests that sanitized output does not contain
// executable code and is safe for storage and rendering.
// Feature: final-polish, Property 13: Input Sanitization
// **Validates: Requirements 9.2, 9.3**
func TestProperty13_InputSanitization(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		input := generateMaliciousInput(r)

		sanitized := Sanitize(input)

		// Property 1: No HTML tags in output
		if containsHTMLTags(sanitized) {
			t.Logf("Sanitized output contains HTML tags: %q -> %q", input, sanitized)
			return false
		}

		// Property 2: No script-like patterns in output
		if ContainsDangerousContent(sanitized) {
			t.Logf("Sanitized output contains dangerous content: %q -> %q", input, sanitized)
			return false
		}

		// Property 3: Output is safe for HTML rendering (properly escaped)
		// Re-escaping should not change the output (idempotent)
		doubleEscaped := html.EscapeString(sanitized)
		if doubleEscaped != sanitized {
			// This is expected - the output is already escaped
			// But we need to verify it doesn't contain unescaped special chars
			if strings.ContainsAny(sanitized, "<>\"'&") && !strings.Contains(sanitized, "&amp;") &&
				!strings.Contains(sanitized, "&lt;") && !strings.Contains(sanitized, "&gt;") &&
				!strings.Contains(sanitized, "&quot;") && !strings.Contains(sanitized, "&#39;") {
				t.Logf("Sanitized output contains unescaped special chars: %q", sanitized)
				return false
			}
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 13 (Input Sanitization) failed: %v", err)
	}
}

// generateMaliciousInput generates random input that may contain malicious content.
func generateMaliciousInput(r *rand.Rand) string {
	// Mix of safe and potentially dangerous content
	parts := []string{
		// Safe content
		"Build a REST API",
		"Create a web application",
		"Hello world",
		"Test project",
		// HTML tags
		"<script>alert('xss')</script>",
		"<img src=x onerror=alert(1)>",
		"<div onclick=alert(1)>click me</div>",
		"<a href='javascript:alert(1)'>link</a>",
		"<iframe src='evil.com'></iframe>",
		"<style>body{display:none}</style>",
		// Script patterns
		"javascript:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"vbscript:msgbox(1)",
		"onclick=alert(1)",
		"onerror=alert(1)",
		// Special characters
		"<>&\"'",
		"test & test",
		"a < b > c",
		// Unicode and encoding tricks
		"&#60;script&#62;",
		"%3Cscript%3E",
		// Nested/broken tags
		"<scr<script>ipt>",
		"<<script>script>",
		"<script",
		"script>",
	}

	// Build random combination
	numParts := 1 + r.Intn(5)
	var selected []string
	for i := 0; i < numParts; i++ {
		selected = append(selected, parts[r.Intn(len(parts))])
	}

	return strings.Join(selected, " ")
}

// containsHTMLTags checks if the string contains HTML tags.
func containsHTMLTags(s string) bool {
	return htmlTagRegex.MatchString(s)
}

// TestProperty13_SanitizePreservesLegitimateContent tests that legitimate
// content is preserved after sanitization.
func TestProperty13_SanitizePreservesLegitimateContent(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		input := generateLegitimateInput(r)

		sanitized := Sanitize(input)

		// The sanitized output should contain the core content
		// (though it may be escaped)
		if sanitized == "" && input != "" {
			t.Logf("Sanitization removed all content: %q -> %q", input, sanitized)
			return false
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 13 (Preserve Legitimate Content) failed: %v", err)
	}
}

// generateLegitimateInput generates random legitimate user input.
func generateLegitimateInput(r *rand.Rand) string {
	words := []string{
		"Build", "Create", "Develop", "Design", "Implement",
		"a", "an", "the", "my", "our",
		"REST", "API", "web", "mobile", "CLI",
		"application", "service", "tool", "system", "platform",
		"for", "with", "using", "that", "which",
		"users", "data", "files", "requests", "responses",
		"authentication", "authorization", "validation", "processing",
	}

	numWords := 3 + r.Intn(10)
	var selected []string
	for i := 0; i < numWords; i++ {
		selected = append(selected, words[r.Intn(len(words))])
	}

	return strings.Join(selected, " ")
}

// TestProperty13_LengthValidation tests that length validation works correctly.
func TestProperty13_LengthValidation(t *testing.T) {
	property := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		// Generate input of random length
		length := r.Intn(MaxProjectIdeaLength * 2)
		input := generateStringOfLength(r, length)

		result, err := ValidateAndSanitize(input, MaxProjectIdeaLength)

		if length > MaxProjectIdeaLength {
			// Should return error for too long input
			if err != ErrInputTooLong {
				t.Logf("Expected ErrInputTooLong for length %d, got %v", length, err)
				return false
			}
		} else if strings.TrimSpace(input) == "" {
			// Should return error for empty input
			if err != ErrEmptyInput {
				t.Logf("Expected ErrEmptyInput for empty input, got %v", err)
				return false
			}
		} else {
			// Should succeed
			if err != nil {
				t.Logf("Unexpected error for valid input: %v", err)
				return false
			}
			// Result should not be longer than input
			if len(result) > len(input) {
				t.Logf("Sanitized output longer than input: %d > %d", len(result), len(input))
				return false
			}
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 13 (Length Validation) failed: %v", err)
	}
}

// generateStringOfLength generates a random string of approximately the given length.
func generateStringOfLength(r *rand.Rand, length int) string {
	if length <= 0 {
		return ""
	}
	chars := make([]byte, length)
	for i := range chars {
		// Mix of letters, numbers, and spaces
		switch r.Intn(3) {
		case 0:
			chars[i] = byte('a' + r.Intn(26))
		case 1:
			chars[i] = byte('0' + r.Intn(10))
		case 2:
			chars[i] = ' '
		}
	}
	return string(chars)
}

// TestSanitize_SpecificCases tests specific sanitization cases.
func TestSanitize_SpecificCases(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "script tag",
			input:    "<script>alert('xss')</script>",
			expected: "alert(&#39;xss&#39;)",
		},
		{
			name:     "img with onerror",
			input:    "<img src=x onerror=alert(1)>",
			expected: "",
		},
		{
			name:     "javascript protocol",
			input:    "javascript:alert(1)",
			expected: "alert(1)",
		},
		{
			name:     "data protocol",
			input:    "data:text/html,<script>alert(1)</script>",
			expected: "text/html,alert(1)",
		},
		{
			name:     "onclick handler",
			input:    "onclick=alert(1)",
			expected: "alert(1)",
		},
		{
			name:     "special characters",
			input:    "<>&\"'",
			expected: "&amp;&#34;&#39;",
		},
		{
			name:     "mixed content",
			input:    "Hello <script>evil</script> world",
			expected: "Hello evil world",
		},
		{
			name:     "multiple spaces",
			input:    "Hello    world",
			expected: "Hello world",
		},
		{
			name:     "leading/trailing spaces",
			input:    "  Hello world  ",
			expected: "Hello world",
		},
		{
			name:     "nested tags",
			input:    "<div><script>alert(1)</script></div>",
			expected: "alert(1)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Sanitize(tc.input)
			if result != tc.expected {
				t.Errorf("Sanitize(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// TestStripHTML tests the StripHTML function.
func TestStripHTML(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no tags",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "simple tag",
			input:    "<b>bold</b>",
			expected: "bold",
		},
		{
			name:     "self-closing tag",
			input:    "Hello<br/>world",
			expected: "Helloworld",
		},
		{
			name:     "multiple tags",
			input:    "<div><p>Hello</p></div>",
			expected: "Hello",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StripHTML(tc.input)
			if result != tc.expected {
				t.Errorf("StripHTML(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// TestContainsDangerousContent tests the danger detection function.
func TestContainsDangerousContent(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"plain text", "Hello world", false},
		{"javascript protocol", "javascript:alert(1)", true},
		{"data protocol", "data:text/html", true},
		{"vbscript protocol", "vbscript:msgbox", true},
		{"onclick handler", "onclick=alert(1)", true},
		{"onerror handler", "onerror=alert(1)", true},
		{"onload handler", "onload=alert(1)", true},
		{"safe url", "https://example.com", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ContainsDangerousContent(tc.input)
			if result != tc.expected {
				t.Errorf("ContainsDangerousContent(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

// TestValidateProjectIdea tests project idea validation.
func TestValidateProjectIdea(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectError error
	}{
		{"valid input", "Build a REST API", nil},
		{"empty input", "", ErrEmptyInput},
		{"whitespace only", "   ", ErrEmptyInput},
		{"too long", strings.Repeat("a", MaxProjectIdeaLength+1), ErrInputTooLong},
		{"max length", strings.Repeat("a", MaxProjectIdeaLength), nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateProjectIdea(tc.input)
			if err != tc.expectError {
				t.Errorf("ValidateProjectIdea(%q) error = %v, want %v", tc.input, err, tc.expectError)
			}
		})
	}
}

// TestValidateAnswer tests answer validation.
func TestValidateAnswer(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectError error
	}{
		{"valid input", "Yes, I want authentication", nil},
		{"empty input", "", nil}, // Empty answers are allowed
		{"whitespace only", "   ", nil},
		{"too long", strings.Repeat("a", MaxAnswerLength+1), ErrInputTooLong},
		{"max length", strings.Repeat("a", MaxAnswerLength), nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateAnswer(tc.input)
			if err != tc.expectError {
				t.Errorf("ValidateAnswer(%q) error = %v, want %v", tc.input, err, tc.expectError)
			}
		})
	}
}

// TestTruncateWithEllipsis tests the truncation function.
func TestTruncateWithEllipsis(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{"short string", "Hello", 10, "Hello"},
		{"exact length", "Hello", 5, "Hello"},
		{"needs truncation", "Hello world", 8, "Hello..."},
		{"very short max", "Hello", 3, "Hel"},
		{"empty string", "", 10, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := TruncateWithEllipsis(tc.input, tc.maxLength)
			if result != tc.expected {
				t.Errorf("TruncateWithEllipsis(%q, %d) = %q, want %q", tc.input, tc.maxLength, result, tc.expected)
			}
		})
	}
}

// TestSanitizePreserveNewlines tests newline preservation.
func TestSanitizePreserveNewlines(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single line",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "multiple lines",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "with empty lines",
			input:    "Line 1\n\nLine 2",
			expected: "Line 1\nLine 2",
		},
		{
			name:     "with HTML tags",
			input:    "Line 1\n<script>evil</script>\nLine 2",
			expected: "Line 1\nevil\nLine 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SanitizePreserveNewlines(tc.input)
			if result != tc.expected {
				t.Errorf("SanitizePreserveNewlines(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}
