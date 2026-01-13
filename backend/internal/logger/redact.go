package logger

import (
	"log/slog"
	"strings"
)

// RedactedValue is the replacement string for sensitive data
const RedactedValue = "[REDACTED]"

// sensitiveKeys is a list of keys that should have their values redacted
var sensitiveKeys = []string{
	"password",
	"token",
	"api_key",
	"apikey",
	"secret",
	"authorization",
	"auth",
	"credential",
	"private_key",
	"privatekey",
	"access_token",
	"refresh_token",
	"bearer",
}

// RedactSensitive replaces sensitive values in a string map
func RedactSensitive(data map[string]string) map[string]string {
	result := make(map[string]string, len(data))
	for k, v := range data {
		if isSensitiveKey(k) {
			result[k] = RedactedValue
		} else {
			result[k] = v
		}
	}
	return result
}

// RedactSensitiveAny replaces sensitive values in a map[string]any
func RedactSensitiveAny(data map[string]any) map[string]any {
	result := make(map[string]any, len(data))
	for k, v := range data {
		if isSensitiveKey(k) {
			result[k] = RedactedValue
		} else {
			// Recursively redact nested maps
			switch val := v.(type) {
			case map[string]any:
				result[k] = RedactSensitiveAny(val)
			case map[string]string:
				result[k] = RedactSensitive(val)
			default:
				result[k] = v
			}
		}
	}
	return result
}

// isSensitiveKey checks if a key name indicates sensitive data
func isSensitiveKey(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(lowerKey, sensitive) {
			return true
		}
	}
	return false
}

// redactAttr redacts sensitive slog attributes
func redactAttr(a slog.Attr) slog.Attr {
	if isSensitiveKey(a.Key) {
		return slog.String(a.Key, RedactedValue)
	}

	// Handle nested groups
	if a.Value.Kind() == slog.KindGroup {
		attrs := a.Value.Group()
		redactedAttrs := make([]slog.Attr, len(attrs))
		for i, attr := range attrs {
			redactedAttrs[i] = redactAttr(attr)
		}
		return slog.Group(a.Key, anySlice(redactedAttrs)...)
	}

	return a
}

// anySlice converts []slog.Attr to []any for slog.Group
func anySlice(attrs []slog.Attr) []any {
	result := make([]any, len(attrs))
	for i, a := range attrs {
		result[i] = a
	}
	return result
}

// RedactString redacts sensitive patterns in a string
// This is useful for redacting values in log messages
func RedactString(s string) string {
	// Common patterns to redact
	patterns := []struct {
		prefix string
		suffix string
	}{
		{"password=", "&"},
		{"password=", " "},
		{"token=", "&"},
		{"token=", " "},
		{"api_key=", "&"},
		{"api_key=", " "},
		{"secret=", "&"},
		{"secret=", " "},
		{"Authorization: Bearer ", "\n"},
		{"Authorization: Bearer ", " "},
	}

	result := s
	for _, p := range patterns {
		result = redactBetween(result, p.prefix, p.suffix)
	}
	return result
}

// redactBetween redacts content between a prefix and suffix
func redactBetween(s, prefix, suffix string) string {
	result := s
	for {
		start := strings.Index(strings.ToLower(result), strings.ToLower(prefix))
		if start == -1 {
			break
		}
		start += len(prefix)

		end := strings.Index(result[start:], suffix)
		if end == -1 {
			// Redact to end of string
			result = result[:start] + RedactedValue
			break
		}

		result = result[:start] + RedactedValue + result[start+end:]
	}
	return result
}
