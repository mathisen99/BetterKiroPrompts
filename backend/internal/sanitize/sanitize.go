// Package sanitize provides input sanitization functions to prevent XSS
// and ensure safe storage of user-generated content.
package sanitize

import (
	"errors"
	"html"
	"regexp"
	"strings"
	"unicode"
)

// Length limits for various input types
const (
	MaxProjectIdeaLength = 2000
	MaxAnswerLength      = 1000
	MaxGenericLength     = 5000
)

// Validation errors
var (
	ErrInputTooLong   = errors.New("input exceeds maximum length")
	ErrEmptyInput     = errors.New("input is empty")
	ErrInvalidContent = errors.New("input contains invalid content")
)

// htmlTagRegex matches HTML tags including self-closing tags
var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

// scriptRegex matches script-like patterns that could be dangerous
var scriptRegex = regexp.MustCompile(`(?i)(javascript:|data:|vbscript:|on\w+\s*=)`)

// Sanitize removes potentially dangerous content from user input.
// It removes HTML tags, escapes special characters, and normalizes whitespace.
func Sanitize(input string) string {
	if input == "" {
		return ""
	}

	// Remove HTML tags
	result := htmlTagRegex.ReplaceAllString(input, "")

	// Remove script-like patterns
	result = scriptRegex.ReplaceAllString(result, "")

	// Escape HTML entities for safe rendering
	result = html.EscapeString(result)

	// Normalize whitespace (collapse multiple spaces, trim)
	result = normalizeWhitespace(result)

	return result
}

// SanitizePreserveNewlines sanitizes input while preserving intentional newlines.
// Useful for multi-line content like project descriptions.
func SanitizePreserveNewlines(input string) string {
	if input == "" {
		return ""
	}

	// Split by newlines, sanitize each line, rejoin
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		// Remove HTML tags
		line = htmlTagRegex.ReplaceAllString(line, "")
		// Remove script-like patterns
		line = scriptRegex.ReplaceAllString(line, "")
		// Escape HTML entities
		line = html.EscapeString(line)
		// Trim each line
		lines[i] = strings.TrimSpace(line)
	}

	// Rejoin with single newlines, removing empty lines
	var nonEmpty []string
	for _, line := range lines {
		if line != "" {
			nonEmpty = append(nonEmpty, line)
		}
	}

	return strings.Join(nonEmpty, "\n")
}

// ValidateAndSanitize validates input length and sanitizes content.
// Returns the sanitized string or an error if validation fails.
func ValidateAndSanitize(input string, maxLength int) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", ErrEmptyInput
	}

	if len(trimmed) > maxLength {
		return "", ErrInputTooLong
	}

	return Sanitize(trimmed), nil
}

// ValidateProjectIdea validates and sanitizes a project idea.
func ValidateProjectIdea(idea string) (string, error) {
	return ValidateAndSanitize(idea, MaxProjectIdeaLength)
}

// ValidateAnswer validates and sanitizes an answer.
func ValidateAnswer(answer string) (string, error) {
	// Answers can be empty (user might skip)
	if strings.TrimSpace(answer) == "" {
		return "", nil
	}

	if len(answer) > MaxAnswerLength {
		return "", ErrInputTooLong
	}

	return Sanitize(strings.TrimSpace(answer)), nil
}

// StripHTML removes all HTML tags from input without escaping.
// Use when you need plain text without HTML encoding.
func StripHTML(input string) string {
	if input == "" {
		return ""
	}
	return htmlTagRegex.ReplaceAllString(input, "")
}

// ContainsDangerousContent checks if input contains potentially dangerous patterns.
func ContainsDangerousContent(input string) bool {
	return scriptRegex.MatchString(input)
}

// normalizeWhitespace collapses multiple whitespace characters into single spaces
// and trims leading/trailing whitespace.
func normalizeWhitespace(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	prevSpace := false
	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace {
				builder.WriteRune(' ')
				prevSpace = true
			}
		} else {
			builder.WriteRune(r)
			prevSpace = false
		}
	}

	return strings.TrimSpace(builder.String())
}

// TruncateWithEllipsis truncates input to maxLength and adds ellipsis if truncated.
func TruncateWithEllipsis(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	if maxLength <= 3 {
		return input[:maxLength]
	}
	return input[:maxLength-3] + "..."
}
