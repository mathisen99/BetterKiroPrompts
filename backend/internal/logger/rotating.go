package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// RotatingFile handles log rotation based on file size
type RotatingFile struct {
	path       string
	maxSize    int64
	maxAgeDays int
	file       *os.File
	size       int64
	mu         sync.Mutex
}

// NewRotatingFile creates a new rotating file writer
func NewRotatingFile(path string, maxSize int64, maxAgeDays int) (*RotatingFile, error) {
	rf := &RotatingFile{
		path:       path,
		maxSize:    maxSize,
		maxAgeDays: maxAgeDays,
	}

	if err := rf.openFile(); err != nil {
		return nil, err
	}

	// Run initial cleanup
	go rf.Cleanup()

	return rf, nil
}

// openFile opens or creates the log file
func (rf *RotatingFile) openFile() error {
	// Ensure directory exists
	dir := filepath.Dir(rf.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(rf.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Get current file size
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	rf.file = file
	rf.size = info.Size()

	return nil
}

// Write implements io.Writer and tracks size for rotation
func (rf *RotatingFile) Write(p []byte) (n int, err error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// Check if rotation is needed
	if rf.size+int64(len(p)) > rf.maxSize {
		if err := rf.rotate(); err != nil {
			// Log rotation failed, but continue writing to current file
			fmt.Fprintf(os.Stderr, "log rotation failed: %v\n", err)
		}
	}

	n, err = rf.file.Write(p)
	rf.size += int64(n)

	return n, err
}

// Rotate renames the current file with a timestamp suffix and opens a new file
func (rf *RotatingFile) rotate() error {
	if rf.file != nil {
		_ = rf.file.Close()
	}

	// Generate rotated filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	ext := filepath.Ext(rf.path)
	base := strings.TrimSuffix(rf.path, ext)
	rotatedPath := fmt.Sprintf("%s.%s%s", base, timestamp, ext)

	// Rename current file
	if err := os.Rename(rf.path, rotatedPath); err != nil {
		// If rename fails, try to reopen the original file
		if openErr := rf.openFile(); openErr != nil {
			return fmt.Errorf("rotation failed and couldn't reopen: %w", openErr)
		}
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Open new file
	if err := rf.openFile(); err != nil {
		return fmt.Errorf("failed to open new log file after rotation: %w", err)
	}

	// Trigger cleanup in background
	go rf.Cleanup()

	return nil
}

// Cleanup removes log files older than MaxAgeDays
func (rf *RotatingFile) Cleanup() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	dir := filepath.Dir(rf.path)
	base := filepath.Base(rf.path)
	ext := filepath.Ext(base)
	prefix := strings.TrimSuffix(base, ext)

	// Find all rotated files for this log
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read log directory for cleanup: %v\n", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -rf.maxAgeDays)
	var toDelete []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Match rotated files: prefix.TIMESTAMP.ext
		if !strings.HasPrefix(name, prefix+".") || !strings.HasSuffix(name, ext) {
			continue
		}

		// Skip the current log file
		if name == base {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			toDelete = append(toDelete, filepath.Join(dir, name))
		}
	}

	// Sort by modification time (oldest first) and delete
	sort.Slice(toDelete, func(i, j int) bool {
		iInfo, _ := os.Stat(toDelete[i])
		jInfo, _ := os.Stat(toDelete[j])
		if iInfo == nil || jInfo == nil {
			return false
		}
		return iInfo.ModTime().Before(jInfo.ModTime())
	})

	for _, path := range toDelete {
		if err := os.Remove(path); err != nil {
			fmt.Fprintf(os.Stderr, "failed to remove old log file %s: %v\n", path, err)
		}
	}
}

// Close closes the underlying file
func (rf *RotatingFile) Close() error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.file != nil {
		return rf.file.Close()
	}
	return nil
}

// Sync flushes the file to disk
func (rf *RotatingFile) Sync() error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.file != nil {
		return rf.file.Sync()
	}
	return nil
}
