package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[90m"
)

// ColorHandler wraps slog.Handler with ANSI colors for terminal output
type ColorHandler struct {
	writer    io.Writer
	opts      *slog.HandlerOptions
	component string
	levelVar  *slog.LevelVar
	mu        sync.Mutex
	attrs     []slog.Attr
	groups    []string
}

// NewColorHandler creates a new color handler
func NewColorHandler(w io.Writer, opts *slog.HandlerOptions, component string) *ColorHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColorHandler{
		writer:    w,
		opts:      opts,
		component: component,
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle handles the Record
func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if colors should be disabled
	noColor := !isTerminal(h.writer) || os.Getenv("NO_COLOR") != ""

	var buf strings.Builder

	// Timestamp in gray
	timestamp := r.Time.Format("2006-01-02 15:04:05.000")
	if noColor {
		buf.WriteString(timestamp)
	} else {
		buf.WriteString(colorGray)
		buf.WriteString(timestamp)
		buf.WriteString(colorReset)
	}
	buf.WriteString(" ")

	// Level with color
	levelStr := levelString(r.Level)
	if noColor {
		buf.WriteString(fmt.Sprintf("%-5s", levelStr))
	} else {
		buf.WriteString(levelColor(r.Level))
		buf.WriteString(fmt.Sprintf("%-5s", levelStr))
		buf.WriteString(colorReset)
	}
	buf.WriteString(" ")

	// Component
	buf.WriteString("[")
	buf.WriteString(h.component)
	buf.WriteString("] ")

	// Message
	buf.WriteString(r.Message)

	// Attributes
	r.Attrs(func(a slog.Attr) bool {
		// Apply redaction
		a = redactAttr(a)

		// Skip component as it's already shown
		if a.Key == "component" {
			return true
		}

		buf.WriteString(" ")

		// Highlight request_id in magenta
		if a.Key == "request_id" && !noColor {
			buf.WriteString(a.Key)
			buf.WriteString("=")
			buf.WriteString(colorMagenta)
			buf.WriteString(a.Value.String())
			buf.WriteString(colorReset)
		} else {
			buf.WriteString(a.Key)
			buf.WriteString("=")
			buf.WriteString(formatValue(a.Value))
		}
		return true
	})

	// Add pre-set attrs
	for _, a := range h.attrs {
		a = redactAttr(a)
		if a.Key == "component" {
			continue
		}
		buf.WriteString(" ")
		buf.WriteString(a.Key)
		buf.WriteString("=")
		buf.WriteString(formatValue(a.Value))
	}

	buf.WriteString("\n")

	_, err := h.writer.Write([]byte(buf.String()))
	return err
}

// WithAttrs returns a new Handler with the given attributes added
func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &ColorHandler{
		writer:    h.writer,
		opts:      h.opts,
		component: h.component,
		levelVar:  h.levelVar,
		attrs:     make([]slog.Attr, len(h.attrs)+len(attrs)),
		groups:    h.groups,
	}
	copy(newHandler.attrs, h.attrs)
	copy(newHandler.attrs[len(h.attrs):], attrs)
	return newHandler
}

// WithGroup returns a new Handler with the given group appended to the receiver's existing groups
func (h *ColorHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newHandler := &ColorHandler{
		writer:    h.writer,
		opts:      h.opts,
		component: h.component,
		levelVar:  h.levelVar,
		attrs:     h.attrs,
		groups:    append(h.groups, name),
	}
	return newHandler
}

// levelColor returns the ANSI color code for a log level
func levelColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return colorRed
	case level >= slog.LevelWarn:
		return colorYellow
	case level >= slog.LevelInfo:
		return colorGreen
	default:
		return colorCyan
	}
}

// levelString returns a string representation of the log level
func levelString(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARN"
	case level >= slog.LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}

// formatValue formats a slog.Value for output
func formatValue(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		s := v.String()
		if strings.ContainsAny(s, " \t\n\"") {
			return fmt.Sprintf("%q", s)
		}
		return s
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format("2006-01-02T15:04:05.000Z07:00")
	default:
		return fmt.Sprintf("%v", v.Any())
	}
}

// isTerminal checks if the writer is a terminal
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		info, err := f.Stat()
		if err != nil {
			return false
		}
		return (info.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// ColorWriter wraps a writer to add color support detection
type ColorWriter struct {
	writer   io.Writer
	levelVar *slog.LevelVar
}

// NewColorWriter creates a new color-aware writer
func NewColorWriter(w io.Writer, levelVar *slog.LevelVar) *ColorWriter {
	return &ColorWriter{
		writer:   w,
		levelVar: levelVar,
	}
}

// Write implements io.Writer
func (cw *ColorWriter) Write(p []byte) (n int, err error) {
	return cw.writer.Write(p)
}
