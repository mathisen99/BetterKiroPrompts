package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Language represents a detected programming language.
type Language string

// Supported languages for detection.
const (
	LangGo         Language = "go"
	LangJavaScript Language = "javascript"
	LangTypeScript Language = "typescript"
	LangPython     Language = "python"
	LangJava       Language = "java"
	LangRuby       Language = "ruby"
	LangPHP        Language = "php"
	LangC          Language = "c"
	LangCPP        Language = "cpp"
	LangRust       Language = "rust"
	LangUnknown    Language = "unknown"
)

// LanguageResult contains information about a detected language.
type LanguageResult struct {
	// Language is the detected programming language.
	Language Language `json:"language"`

	// FileCount is the number of files detected for this language.
	FileCount int `json:"file_count"`

	// Percentage is the percentage of total detected files.
	Percentage float64 `json:"percentage"`
}

// LanguageDetector detects programming languages in a repository.
type LanguageDetector struct {
	// extensionMap maps file extensions to languages.
	extensionMap map[string]Language
}

// NewLanguageDetector creates a new LanguageDetector.
func NewLanguageDetector() *LanguageDetector {
	return &LanguageDetector{
		extensionMap: map[string]Language{
			// Go
			".go": LangGo,

			// JavaScript
			".js":  LangJavaScript,
			".jsx": LangJavaScript,
			".mjs": LangJavaScript,
			".cjs": LangJavaScript,

			// TypeScript
			".ts":  LangTypeScript,
			".tsx": LangTypeScript,
			".mts": LangTypeScript,
			".cts": LangTypeScript,

			// Python
			".py":  LangPython,
			".pyw": LangPython,
			".pyi": LangPython,

			// Java
			".java": LangJava,

			// Ruby
			".rb":      LangRuby,
			".rake":    LangRuby,
			".gemspec": LangRuby,

			// PHP
			".php":   LangPHP,
			".phtml": LangPHP,
			".php3":  LangPHP,
			".php4":  LangPHP,
			".php5":  LangPHP,
			".php7":  LangPHP,
			".phps":  LangPHP,

			// C
			".c": LangC,
			".h": LangC,

			// C++
			".cpp": LangCPP,
			".cxx": LangCPP,
			".cc":  LangCPP,
			".hpp": LangCPP,
			".hxx": LangCPP,
			".hh":  LangCPP,
			".c++": LangCPP,
			".h++": LangCPP,

			// Rust
			".rs": LangRust,
		},
	}
}

// Detect analyzes a repository and returns detected languages sorted by file count.
func (d *LanguageDetector) Detect(repoPath string) ([]LanguageResult, error) {
	// Count files per language
	langCounts := make(map[Language]int)
	totalFiles := 0

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip common non-source directories
			name := info.Name()
			if d.shouldSkipDir(name) {
				return filepath.SkipDir
			}
			return nil
		}

		// Get file extension
		ext := strings.ToLower(filepath.Ext(path))
		if ext == "" {
			return nil
		}

		// Look up language
		if lang, ok := d.extensionMap[ext]; ok {
			langCounts[lang]++
			totalFiles++
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to results slice
	results := make([]LanguageResult, 0, len(langCounts))
	for lang, count := range langCounts {
		percentage := 0.0
		if totalFiles > 0 {
			percentage = float64(count) / float64(totalFiles) * 100
		}
		results = append(results, LanguageResult{
			Language:   lang,
			FileCount:  count,
			Percentage: percentage,
		})
	}

	// Sort by file count (descending), then by language name (ascending) for determinism
	sort.Slice(results, func(i, j int) bool {
		if results[i].FileCount != results[j].FileCount {
			return results[i].FileCount > results[j].FileCount
		}
		// Secondary sort by language name for deterministic ordering
		return results[i].Language < results[j].Language
	})

	return results, nil
}

// DetectLanguages returns just the language names sorted by prevalence.
func (d *LanguageDetector) DetectLanguages(repoPath string) ([]Language, error) {
	results, err := d.Detect(repoPath)
	if err != nil {
		return nil, err
	}

	languages := make([]Language, len(results))
	for i, r := range results {
		languages[i] = r.Language
	}

	return languages, nil
}

// GetLanguageForExtension returns the language for a given file extension.
func (d *LanguageDetector) GetLanguageForExtension(ext string) Language {
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	if lang, ok := d.extensionMap[ext]; ok {
		return lang
	}
	return LangUnknown
}

// shouldSkipDir returns true if the directory should be skipped during detection.
func (d *LanguageDetector) shouldSkipDir(name string) bool {
	skipDirs := map[string]bool{
		".git":         true,
		"node_modules": true,
		"vendor":       true,
		".venv":        true,
		"venv":         true,
		"__pycache__":  true,
		".cache":       true,
		"dist":         true,
		"build":        true,
		"target":       true, // Rust/Java build output
		".idea":        true,
		".vscode":      true,
		".kiro":        true,
	}
	return skipDirs[name]
}

// GetSupportedExtensions returns all supported file extensions.
func (d *LanguageDetector) GetSupportedExtensions() []string {
	extensions := make([]string, 0, len(d.extensionMap))
	for ext := range d.extensionMap {
		extensions = append(extensions, ext)
	}
	sort.Strings(extensions)
	return extensions
}

// GetSupportedLanguages returns all supported languages.
func (d *LanguageDetector) GetSupportedLanguages() []Language {
	return []Language{
		LangGo,
		LangJavaScript,
		LangTypeScript,
		LangPython,
		LangJava,
		LangRuby,
		LangPHP,
		LangC,
		LangCPP,
		LangRust,
	}
}

// String returns the string representation of a Language.
func (l Language) String() string {
	return string(l)
}
