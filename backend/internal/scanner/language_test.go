package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"testing/quick"
)

// =============================================================================
// Unit Tests for Language Detection
// =============================================================================

func TestLanguageDetector_GetLanguageForExtension(t *testing.T) {
	d := NewLanguageDetector()

	tests := []struct {
		ext  string
		want Language
	}{
		// Go
		{".go", LangGo},
		{"go", LangGo}, // Without dot

		// JavaScript
		{".js", LangJavaScript},
		{".jsx", LangJavaScript},
		{".mjs", LangJavaScript},
		{".cjs", LangJavaScript},

		// TypeScript
		{".ts", LangTypeScript},
		{".tsx", LangTypeScript},
		{".mts", LangTypeScript},
		{".cts", LangTypeScript},

		// Python
		{".py", LangPython},
		{".pyw", LangPython},
		{".pyi", LangPython},

		// Java
		{".java", LangJava},

		// Ruby
		{".rb", LangRuby},
		{".rake", LangRuby},
		{".gemspec", LangRuby},

		// PHP
		{".php", LangPHP},
		{".phtml", LangPHP},

		// C
		{".c", LangC},
		{".h", LangC},

		// C++
		{".cpp", LangCPP},
		{".cxx", LangCPP},
		{".cc", LangCPP},
		{".hpp", LangCPP},
		{".hxx", LangCPP},

		// Rust
		{".rs", LangRust},

		// Unknown
		{".txt", LangUnknown},
		{".md", LangUnknown},
		{".json", LangUnknown},
		{"", LangUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := d.GetLanguageForExtension(tt.ext)
			if got != tt.want {
				t.Errorf("GetLanguageForExtension(%q) = %v, want %v", tt.ext, got, tt.want)
			}
		})
	}
}

func TestLanguageDetector_GetLanguageForExtension_CaseInsensitive(t *testing.T) {
	d := NewLanguageDetector()

	tests := []struct {
		ext  string
		want Language
	}{
		{".GO", LangGo},
		{".Go", LangGo},
		{".JS", LangJavaScript},
		{".PY", LangPython},
		{".CPP", LangCPP},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := d.GetLanguageForExtension(tt.ext)
			if got != tt.want {
				t.Errorf("GetLanguageForExtension(%q) = %v, want %v", tt.ext, got, tt.want)
			}
		})
	}
}

func TestLanguageDetector_Detect(t *testing.T) {
	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "lang-detect-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test files
	files := map[string]string{
		"main.go":       "package main",
		"util.go":       "package util",
		"app.ts":        "const x = 1",
		"component.tsx": "export default",
		"script.py":     "print('hello')",
	}

	for name, content := range files {
		path := filepath.Join(tempDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	d := NewLanguageDetector()
	results, err := d.Detect(tempDir)
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Verify results
	if len(results) != 3 {
		t.Errorf("Expected 3 languages, got %d", len(results))
	}

	// Go should be first (2 files)
	if results[0].Language != LangGo || results[0].FileCount != 2 {
		t.Errorf("Expected Go with 2 files first, got %v with %d files",
			results[0].Language, results[0].FileCount)
	}

	// TypeScript should be second (2 files)
	if results[1].Language != LangTypeScript || results[1].FileCount != 2 {
		t.Errorf("Expected TypeScript with 2 files second, got %v with %d files",
			results[1].Language, results[1].FileCount)
	}

	// Python should be third (1 file)
	if results[2].Language != LangPython || results[2].FileCount != 1 {
		t.Errorf("Expected Python with 1 file third, got %v with %d files",
			results[2].Language, results[2].FileCount)
	}
}

func TestLanguageDetector_Detect_SkipsDirectories(t *testing.T) {
	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "lang-detect-skip-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create a source file
	if err := os.WriteFile(filepath.Join(tempDir, "main.go"), []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create node_modules with JS files (should be skipped)
	nodeModules := filepath.Join(tempDir, "node_modules")
	if err := os.MkdirAll(nodeModules, 0755); err != nil {
		t.Fatalf("Failed to create node_modules: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nodeModules, "lib.js"), []byte("module.exports"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create vendor with Go files (should be skipped)
	vendor := filepath.Join(tempDir, "vendor")
	if err := os.MkdirAll(vendor, 0755); err != nil {
		t.Fatalf("Failed to create vendor: %v", err)
	}
	if err := os.WriteFile(filepath.Join(vendor, "dep.go"), []byte("package dep"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d := NewLanguageDetector()
	results, err := d.Detect(tempDir)
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should only detect the main.go file
	if len(results) != 1 {
		t.Errorf("Expected 1 language, got %d", len(results))
	}

	if results[0].Language != LangGo || results[0].FileCount != 1 {
		t.Errorf("Expected Go with 1 file, got %v with %d files",
			results[0].Language, results[0].FileCount)
	}
}

func TestLanguageDetector_DetectLanguages(t *testing.T) {
	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "lang-detect-langs-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test files
	files := []string{"main.go", "util.go", "app.ts"}
	for _, name := range files {
		path := filepath.Join(tempDir, name)
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	d := NewLanguageDetector()
	languages, err := d.DetectLanguages(tempDir)
	if err != nil {
		t.Fatalf("DetectLanguages() error = %v", err)
	}

	if len(languages) != 2 {
		t.Errorf("Expected 2 languages, got %d", len(languages))
	}

	// Go should be first (2 files)
	if languages[0] != LangGo {
		t.Errorf("Expected Go first, got %v", languages[0])
	}

	// TypeScript should be second (1 file)
	if languages[1] != LangTypeScript {
		t.Errorf("Expected TypeScript second, got %v", languages[1])
	}
}

func TestLanguageDetector_GetSupportedLanguages(t *testing.T) {
	d := NewLanguageDetector()
	languages := d.GetSupportedLanguages()

	expected := []Language{
		LangGo, LangJavaScript, LangTypeScript, LangPython,
		LangJava, LangRuby, LangPHP, LangC, LangCPP, LangRust,
	}

	if len(languages) != len(expected) {
		t.Errorf("Expected %d languages, got %d", len(expected), len(languages))
	}

	for i, lang := range expected {
		if languages[i] != lang {
			t.Errorf("Expected language %d to be %v, got %v", i, lang, languages[i])
		}
	}
}

// =============================================================================
// Property-Based Tests for Language Detection
// =============================================================================

// TestProperty6_LanguageDetectionAccuracy tests Property 6: Language Detection Accuracy
// Feature: info-and-security-scan, Property 6: Language Detection Accuracy
// **Validates: Requirements 6.1, 6.2, 6.3**
//
// Property: For any repository containing files with known extensions:
// - Files with .go extension SHALL be identified as Go
// - Files with .js extension SHALL be identified as JavaScript
// - Files with .ts extension SHALL be identified as TypeScript
// - Files with .py extension SHALL be identified as Python
// - Files with .java extension SHALL be identified as Java
// - Files with .rb extension SHALL be identified as Ruby
// - Files with .php extension SHALL be identified as PHP
// - Files with .c/.h extension SHALL be identified as C
// - Files with .cpp/.hpp extension SHALL be identified as C++
// - Files with .rs extension SHALL be identified as Rust
//
// For any repository with multiple languages, the results SHALL be sorted by file count in descending order.
func TestProperty6_LanguageDetectionAccuracy(t *testing.T) {
	d := NewLanguageDetector()

	// Sub-property 1: Extension to language mapping is correct
	t.Run("extension_to_language_mapping", func(t *testing.T) {
		extensionTests := map[string]Language{
			".go":   LangGo,
			".js":   LangJavaScript,
			".jsx":  LangJavaScript,
			".ts":   LangTypeScript,
			".tsx":  LangTypeScript,
			".py":   LangPython,
			".java": LangJava,
			".rb":   LangRuby,
			".php":  LangPHP,
			".c":    LangC,
			".h":    LangC,
			".cpp":  LangCPP,
			".hpp":  LangCPP,
			".rs":   LangRust,
		}

		for ext, expectedLang := range extensionTests {
			got := d.GetLanguageForExtension(ext)
			if got != expectedLang {
				t.Errorf("Extension %s: expected %v, got %v", ext, expectedLang, got)
			}
		}
	})

	// Sub-property 2: Results are sorted by file count (descending)
	t.Run("results_sorted_by_file_count", func(t *testing.T) {
		// Create a temporary directory with varying file counts
		tempDir, err := os.MkdirTemp("", "lang-sort-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		// Create files: 5 Go, 3 Python, 1 TypeScript
		for i := range 5 {
			path := filepath.Join(tempDir, "file"+string(rune('0'+i))+".go")
			if err := os.WriteFile(path, []byte("package main"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}
		for i := range 3 {
			path := filepath.Join(tempDir, "file"+string(rune('0'+i))+".py")
			if err := os.WriteFile(path, []byte("print()"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}
		if err := os.WriteFile(filepath.Join(tempDir, "app.ts"), []byte("const x = 1"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		results, err := d.Detect(tempDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Verify sorting
		for i := 1; i < len(results); i++ {
			if results[i].FileCount > results[i-1].FileCount {
				t.Errorf("Results not sorted: %v (%d) should come before %v (%d)",
					results[i].Language, results[i].FileCount,
					results[i-1].Language, results[i-1].FileCount)
			}
		}

		// Verify order: Go (5) > Python (3) > TypeScript (1)
		if results[0].Language != LangGo || results[0].FileCount != 5 {
			t.Errorf("Expected Go with 5 files first, got %v with %d", results[0].Language, results[0].FileCount)
		}
		if results[1].Language != LangPython || results[1].FileCount != 3 {
			t.Errorf("Expected Python with 3 files second, got %v with %d", results[1].Language, results[1].FileCount)
		}
		if results[2].Language != LangTypeScript || results[2].FileCount != 1 {
			t.Errorf("Expected TypeScript with 1 file third, got %v with %d", results[2].Language, results[2].FileCount)
		}
	})

	// Sub-property 3: Percentage calculation is correct
	t.Run("percentage_calculation", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "lang-pct-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		// Create 4 files: 2 Go, 1 Python, 1 TypeScript
		_ = os.WriteFile(filepath.Join(tempDir, "a.go"), []byte(""), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "b.go"), []byte(""), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "c.py"), []byte(""), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "d.ts"), []byte(""), 0644)

		results, err := d.Detect(tempDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		// Go should be 50%
		if results[0].Language != LangGo || results[0].Percentage != 50.0 {
			t.Errorf("Expected Go at 50%%, got %v at %.1f%%", results[0].Language, results[0].Percentage)
		}

		// Python and TypeScript should each be 25%
		for _, r := range results[1:] {
			if r.Percentage != 25.0 {
				t.Errorf("Expected %v at 25%%, got %.1f%%", r.Language, r.Percentage)
			}
		}
	})

	// Sub-property 4: Property test for extension mapping consistency
	t.Run("extension_mapping_property", func(t *testing.T) {
		property := func(ext string) bool {
			// Only test valid extensions (non-empty, starts with dot or we add it)
			if ext == "" {
				return true
			}

			lang := d.GetLanguageForExtension(ext)

			// Result should always be a valid language (including Unknown)
			validLanguages := map[Language]bool{
				LangGo:         true,
				LangJavaScript: true,
				LangTypeScript: true,
				LangPython:     true,
				LangJava:       true,
				LangRuby:       true,
				LangPHP:        true,
				LangC:          true,
				LangCPP:        true,
				LangRust:       true,
				LangUnknown:    true,
			}

			return validLanguages[lang]
		}

		if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 6 (extension mapping consistency) failed: %v", err)
		}
	})
}

// TestProperty6_LanguageDetectionAccuracy_EdgeCases tests edge cases for language detection.
// Feature: info-and-security-scan, Property 6: Language Detection Accuracy
// **Validates: Requirements 6.1, 6.2, 6.3**
func TestProperty6_LanguageDetectionAccuracy_EdgeCases(t *testing.T) {
	d := NewLanguageDetector()

	t.Run("empty_directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "lang-empty-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		results, err := d.Detect(tempDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 results for empty directory, got %d", len(results))
		}
	})

	t.Run("files_without_extensions", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "lang-noext-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		// Create files without extensions
		_ = os.WriteFile(filepath.Join(tempDir, "Makefile"), []byte(""), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "Dockerfile"), []byte(""), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "README"), []byte(""), 0644)

		results, err := d.Detect(tempDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 results for files without extensions, got %d", len(results))
		}
	})

	t.Run("mixed_case_extensions", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "lang-case-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		// Create files with mixed case extensions
		_ = os.WriteFile(filepath.Join(tempDir, "main.GO"), []byte(""), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "app.Ts"), []byte(""), 0644)
		_ = os.WriteFile(filepath.Join(tempDir, "script.PY"), []byte(""), 0644)

		results, err := d.Detect(tempDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if len(results) != 3 {
			t.Errorf("Expected 3 languages, got %d", len(results))
		}
	})

	t.Run("nested_directories", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "lang-nested-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		// Create nested structure
		srcDir := filepath.Join(tempDir, "src", "pkg")
		if err := os.MkdirAll(srcDir, 0755); err != nil {
			t.Fatalf("Failed to create nested dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(srcDir, "main.go"), []byte(""), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(srcDir, "util.go"), []byte(""), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		results, err := d.Detect(tempDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if len(results) != 1 || results[0].FileCount != 2 {
			t.Errorf("Expected 1 language with 2 files, got %d languages", len(results))
		}
	})
}
