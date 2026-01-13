package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

// ctxKey is a type for context keys to avoid collisions
type ctxKey string

// Context keys for logging
const (
	RequestIDKey ctxKey = "request_id"
	ComponentKey ctxKey = "component"
	UserIPKey    ctxKey = "user_ip"
)

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

// WithComponent adds a component name to the context
func WithComponent(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, ComponentKey, name)
}

// WithUserIP adds a user IP to the context
func WithUserIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, UserIPKey, ip)
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// GetComponent retrieves the component name from the context
func GetComponent(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if name, ok := ctx.Value(ComponentKey).(string); ok {
		return name
	}
	return ""
}

// GetUserIP retrieves the user IP from the context
func GetUserIP(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if ip, ok := ctx.Value(UserIPKey).(string); ok {
		return ip
	}
	return ""
}

// GenerateRequestID generates a new unique request ID
func GenerateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a simple counter-based ID if random fails
		return "req-fallback"
	}
	return hex.EncodeToString(b)
}
