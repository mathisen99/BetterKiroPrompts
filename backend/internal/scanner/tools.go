package scanner

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Default tool configuration.
const (
	DefaultToolTimeout = 5 * time.Minute
)

// ToolRunner executes security scanning tools.
type ToolRunner struct {
	timeout time.Duration
}

// ToolRunnerOption is a functional option for configuring a ToolRunner.
type ToolRunnerOption func(*ToolRunner)

// WithToolTimeout sets the timeout for tool execution.
func WithToolTimeout(timeout time.Duration) ToolRunnerOption {
	return func(r *ToolRunner) {
		r.timeout = timeout
	}
}

// NewToolRunner creates a new ToolRunner with the given options.
func NewToolRunner(opts ...ToolRunnerOption) *ToolRunner {
	r := &ToolRunner{
		timeout: DefaultToolTimeout,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// ToolResult contains the result of a tool execution.
type ToolResult struct {
	Tool     string        `json:"tool"`
	Findings []RawFinding  `json:"findings"`
	Error    error         `json:"-"`
	TimedOut bool          `json:"timed_out"`
	Duration time.Duration `json:"duration"`
}

// RawFinding represents a finding from a security tool before aggregation.
type RawFinding struct {
	FilePath    string `json:"file_path"`
	LineNumber  int    `json:"line_number,omitempty"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	RuleID      string `json:"rule_id,omitempty"`
}

// scannerContainer is the name of the scanner container for docker exec.
// This can be overridden via environment variable SCANNER_CONTAINER.
var scannerContainer = "betterkiroprompts-scanner-1"

// SetScannerContainer sets the scanner container name.
func SetScannerContainer(name string) {
	if name != "" {
		scannerContainer = name
	}
}

// runTool executes a command inside the scanner container with timeout.
func (r *ToolRunner) runTool(ctx context.Context, name string, args []string, workDir string) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Build docker exec command to run tool in scanner container
	dockerArgs := []string{
		"exec",
		"-w", workDir,
		scannerContainer,
		name,
	}
	dockerArgs = append(dockerArgs, args...)

	log.Printf("[ToolRunner] Executing: docker %v", dockerArgs)
	cmd := exec.CommandContext(ctx, "docker", dockerArgs...)

	output, err := cmd.CombinedOutput()

	log.Printf("[ToolRunner] Tool %s output length: %d bytes, error: %v", name, len(output), err)
	if len(output) > 0 && len(output) < 500 {
		log.Printf("[ToolRunner] Tool %s output: %s", name, string(output))
	}

	if ctx.Err() == context.DeadlineExceeded {
		return output, true, ctx.Err()
	}

	return output, false, err
}

// RunTrivy executes Trivy for comprehensive vulnerability scanning.
func (r *ToolRunner) RunTrivy(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "trivy"}

	args := []string{
		"fs",
		"--format", "json",
		"--scanners", "vuln,secret,misconfig",
		"--severity", "CRITICAL,HIGH,MEDIUM,LOW",
		"--skip-dirs", ".git",
		repoPath,
	}

	output, timedOut, err := r.runTool(ctx, "trivy", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// Trivy may return non-zero exit code when findings exist
	// Try to parse output anyway
	_ = err

	result.Findings = parseTrivyOutput(output)
	return result
}

// RunSemgrep executes Semgrep with security rulesets.
func (r *ToolRunner) RunSemgrep(ctx context.Context, repoPath string, languages []string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "semgrep"}

	args := []string{
		"scan",
		"--config", "p/security-audit",
		"--config", "p/owasp-top-ten",
		"--json",
		repoPath,
	}

	output, timedOut, err := r.runTool(ctx, "semgrep", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// Semgrep may return non-zero exit code when findings exist
	_ = err

	result.Findings = parseSemgrepOutput(output)
	return result
}

// RunTruffleHog executes TruffleHog for secret detection in git history.
func (r *ToolRunner) RunTruffleHog(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "trufflehog"}

	args := []string{
		"filesystem",
		"--json",
		repoPath,
	}

	output, timedOut, err := r.runTool(ctx, "trufflehog", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// TruffleHog may return non-zero exit code when findings exist
	_ = err

	result.Findings = parseTruffleHogOutput(output)
	return result
}

// RunGitleaks executes Gitleaks for additional secret detection.
func (r *ToolRunner) RunGitleaks(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "gitleaks"}

	// Gitleaks writes JSON to report file, use a temp file
	reportPath := filepath.Join(repoPath, ".gitleaks-report.json")

	args := []string{
		"detect",
		"--source", repoPath,
		"--report-format", "json",
		"--report-path", reportPath,
		"--no-git",
	}

	_, timedOut, err := r.runTool(ctx, "gitleaks", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// Gitleaks returns exit code 1 when findings exist, that's expected
	_ = err

	// Read the report file from the scanner container
	catArgs := []string{"exec", scannerContainer, "cat", reportPath}
	cmd := exec.Command("docker", catArgs...)
	output, _ := cmd.Output()

	// Clean up report file
	rmArgs := []string{"exec", scannerContainer, "rm", "-f", reportPath}
	rmCmd := exec.Command("docker", rmArgs...)
	_ = rmCmd.Run()

	result.Findings = parseGitleaksOutput(output)
	return result
}

// RunGovulncheck executes govulncheck for Go vulnerability scanning.
func (r *ToolRunner) RunGovulncheck(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "govulncheck"}

	args := []string{
		"-json",
		"./...",
	}

	output, timedOut, err := r.runTool(ctx, "govulncheck", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// govulncheck may return non-zero when vulnerabilities found
	_ = err

	result.Findings = parseGovulncheckOutput(output)
	return result
}

// RunBandit executes Bandit for Python security analysis.
func (r *ToolRunner) RunBandit(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "bandit"}

	args := []string{
		"-r",
		"-f", "json",
		repoPath,
	}

	output, timedOut, err := r.runTool(ctx, "bandit", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// Bandit returns non-zero when findings exist
	_ = err

	result.Findings = parseBanditOutput(output)
	return result
}

// RunPipAudit executes pip-audit for Python dependency scanning.
func (r *ToolRunner) RunPipAudit(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "pip-audit"}

	// Check for requirements.txt or setup.py
	reqPath := filepath.Join(repoPath, "requirements.txt")
	args := []string{
		"-r", reqPath,
		"--format", "json",
	}

	output, timedOut, err := r.runTool(ctx, "pip-audit", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// pip-audit returns non-zero when vulnerabilities found
	_ = err

	result.Findings = parsePipAuditOutput(output)
	return result
}

// RunSafety executes Safety for Python dependency checking.
func (r *ToolRunner) RunSafety(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "safety"}

	reqPath := filepath.Join(repoPath, "requirements.txt")
	args := []string{
		"check",
		"-r", reqPath,
		"--json",
	}

	output, timedOut, err := r.runTool(ctx, "safety", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// Safety returns non-zero when vulnerabilities found
	_ = err

	result.Findings = parseSafetyOutput(output)
	return result
}

// RunNpmAudit executes npm audit for JavaScript/TypeScript dependency scanning.
func (r *ToolRunner) RunNpmAudit(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "npm-audit"}

	args := []string{
		"audit",
		"--json",
	}

	output, timedOut, err := r.runTool(ctx, "npm", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// npm audit returns non-zero when vulnerabilities found
	_ = err

	result.Findings = parseNpmAuditOutput(output)
	return result
}

// RunCargoAudit executes cargo audit for Rust dependency scanning.
func (r *ToolRunner) RunCargoAudit(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "cargo-audit"}

	args := []string{
		"audit",
		"--json",
	}

	output, timedOut, err := r.runTool(ctx, "cargo", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// cargo audit returns non-zero when vulnerabilities found
	_ = err

	result.Findings = parseCargoAuditOutput(output)
	return result
}

// RunBundlerAudit executes bundler-audit for Ruby dependency scanning.
func (r *ToolRunner) RunBundlerAudit(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "bundler-audit"}

	args := []string{
		"check",
		"--format", "json",
	}

	output, timedOut, err := r.runTool(ctx, "bundle-audit", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// bundler-audit returns non-zero when vulnerabilities found
	_ = err

	result.Findings = parseBundlerAuditOutput(output)
	return result
}

// RunBrakeman executes Brakeman for Ruby/Rails security scanning.
func (r *ToolRunner) RunBrakeman(ctx context.Context, repoPath string) ToolResult {
	start := time.Now()
	result := ToolResult{Tool: "brakeman"}

	args := []string{
		"-p", repoPath,
		"-f", "json",
		"--no-pager",
	}

	output, timedOut, err := r.runTool(ctx, "brakeman", args, repoPath)
	result.Duration = time.Since(start)
	result.TimedOut = timedOut

	if timedOut {
		return result
	}

	// Brakeman returns non-zero when findings exist
	_ = err

	result.Findings = parseBrakemanOutput(output)
	return result
}

// GetToolsForLanguages returns the list of tools to run for the given languages.
func (r *ToolRunner) GetToolsForLanguages(languages []Language) []string {
	tools := []string{
		// Universal tools (always run)
		"trivy",
		"semgrep",
		"trufflehog",
		"gitleaks",
	}

	// Language-specific tools
	langSet := make(map[Language]bool)
	for _, lang := range languages {
		langSet[lang] = true
	}

	if langSet[LangGo] {
		tools = append(tools, "govulncheck")
	}

	if langSet[LangPython] {
		tools = append(tools, "bandit", "pip-audit", "safety")
	}

	if langSet[LangJavaScript] || langSet[LangTypeScript] {
		tools = append(tools, "npm-audit")
	}

	if langSet[LangRust] {
		tools = append(tools, "cargo-audit")
	}

	if langSet[LangRuby] {
		tools = append(tools, "bundler-audit", "brakeman")
	}

	return tools
}

// RunToolByName runs a specific tool by name.
func (r *ToolRunner) RunToolByName(ctx context.Context, toolName string, repoPath string, languages []Language) ToolResult {
	switch toolName {
	case "trivy":
		return r.RunTrivy(ctx, repoPath)
	case "semgrep":
		return r.RunSemgrep(ctx, repoPath, nil)
	case "trufflehog":
		return r.RunTruffleHog(ctx, repoPath)
	case "gitleaks":
		return r.RunGitleaks(ctx, repoPath)
	case "govulncheck":
		return r.RunGovulncheck(ctx, repoPath)
	case "bandit":
		return r.RunBandit(ctx, repoPath)
	case "pip-audit":
		return r.RunPipAudit(ctx, repoPath)
	case "safety":
		return r.RunSafety(ctx, repoPath)
	case "npm-audit":
		return r.RunNpmAudit(ctx, repoPath)
	case "cargo-audit":
		return r.RunCargoAudit(ctx, repoPath)
	case "bundler-audit":
		return r.RunBundlerAudit(ctx, repoPath)
	case "brakeman":
		return r.RunBrakeman(ctx, repoPath)
	default:
		return ToolResult{
			Tool:  toolName,
			Error: ErrCloneFailed, // Using existing error for now
		}
	}
}

// =============================================================================
// Output Parsers
// =============================================================================

// trivyOutput represents Trivy JSON output structure.
type trivyOutput struct {
	Results []struct {
		Target          string `json:"Target"`
		Vulnerabilities []struct {
			VulnerabilityID string `json:"VulnerabilityID"`
			PkgName         string `json:"PkgName"`
			Severity        string `json:"Severity"`
			Title           string `json:"Title"`
			Description     string `json:"Description"`
		} `json:"Vulnerabilities"`
		Secrets []struct {
			RuleID    string `json:"RuleID"`
			Category  string `json:"Category"`
			Severity  string `json:"Severity"`
			Title     string `json:"Title"`
			Match     string `json:"Match"`
			StartLine int    `json:"StartLine"`
		} `json:"Secrets"`
	} `json:"Results"`
}

func parseTrivyOutput(output []byte) []RawFinding {
	var findings []RawFinding
	var result trivyOutput

	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, r := range result.Results {
		for _, v := range r.Vulnerabilities {
			findings = append(findings, RawFinding{
				FilePath:    r.Target,
				Description: v.Title + ": " + v.Description,
				Severity:    strings.ToLower(v.Severity),
				RuleID:      v.VulnerabilityID,
			})
		}
		for _, s := range r.Secrets {
			findings = append(findings, RawFinding{
				FilePath:    r.Target,
				LineNumber:  s.StartLine,
				Description: s.Title,
				Severity:    strings.ToLower(s.Severity),
				RuleID:      s.RuleID,
			})
		}
	}

	return findings
}

// semgrepOutput represents Semgrep JSON output structure.
type semgrepOutput struct {
	Results []struct {
		CheckID string `json:"check_id"`
		Path    string `json:"path"`
		Start   struct {
			Line int `json:"line"`
		} `json:"start"`
		Extra struct {
			Message  string `json:"message"`
			Severity string `json:"severity"`
		} `json:"extra"`
	} `json:"results"`
}

func parseSemgrepOutput(output []byte) []RawFinding {
	var findings []RawFinding
	var result semgrepOutput

	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, r := range result.Results {
		findings = append(findings, RawFinding{
			FilePath:    r.Path,
			LineNumber:  r.Start.Line,
			Description: r.Extra.Message,
			Severity:    strings.ToLower(r.Extra.Severity),
			RuleID:      r.CheckID,
		})
	}

	return findings
}

func parseTruffleHogOutput(output []byte) []RawFinding {
	var findings []RawFinding

	// TruffleHog outputs one JSON object per line
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		var result struct {
			SourceMetadata struct {
				Data struct {
					Filesystem struct {
						File string `json:"file"`
						Line int    `json:"line"`
					} `json:"Filesystem"`
				} `json:"Data"`
			} `json:"SourceMetadata"`
			DetectorName string `json:"DetectorName"`
			Raw          string `json:"Raw"`
		}

		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}

		// Skip if no detector name (likely an error/log message, not a finding)
		if result.DetectorName == "" {
			continue
		}

		// Skip .git directory findings - these are local clone artifacts, not repo issues
		filePath := result.SourceMetadata.Data.Filesystem.File
		if strings.Contains(filePath, "/.git/") || strings.HasSuffix(filePath, "/.git") {
			continue
		}

		findings = append(findings, RawFinding{
			FilePath:    filePath,
			LineNumber:  result.SourceMetadata.Data.Filesystem.Line,
			Description: "Secret detected: " + result.DetectorName,
			Severity:    "high",
			RuleID:      result.DetectorName,
		})
	}

	return findings
}

// gitleaksOutput represents Gitleaks JSON output structure.
type gitleaksOutput []struct {
	RuleID      string `json:"RuleID"`
	Description string `json:"Description"`
	File        string `json:"File"`
	StartLine   int    `json:"StartLine"`
	Secret      string `json:"Secret"`
}

func parseGitleaksOutput(output []byte) []RawFinding {
	var findings []RawFinding
	var results gitleaksOutput

	if err := json.Unmarshal(output, &results); err != nil {
		return findings
	}

	for _, r := range results {
		findings = append(findings, RawFinding{
			FilePath:    r.File,
			LineNumber:  r.StartLine,
			Description: r.Description,
			Severity:    "high",
			RuleID:      r.RuleID,
		})
	}

	return findings
}

func parseGovulncheckOutput(output []byte) []RawFinding {
	var findings []RawFinding

	// govulncheck outputs JSON lines
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		var result struct {
			Finding struct {
				OSV   string `json:"osv"`
				Trace []struct {
					Module   string `json:"module"`
					Package  string `json:"package"`
					Function string `json:"function"`
					Position struct {
						Filename string `json:"filename"`
						Line     int    `json:"line"`
					} `json:"position"`
				} `json:"trace"`
			} `json:"finding"`
		}

		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}

		if result.Finding.OSV == "" {
			continue
		}

		filePath := ""
		lineNum := 0
		if len(result.Finding.Trace) > 0 {
			filePath = result.Finding.Trace[0].Position.Filename
			lineNum = result.Finding.Trace[0].Position.Line
		}

		findings = append(findings, RawFinding{
			FilePath:    filePath,
			LineNumber:  lineNum,
			Description: "Go vulnerability: " + result.Finding.OSV,
			Severity:    "high",
			RuleID:      result.Finding.OSV,
		})
	}

	return findings
}

// banditOutput represents Bandit JSON output structure.
type banditOutput struct {
	Results []struct {
		Filename   string `json:"filename"`
		LineNumber int    `json:"line_number"`
		IssueText  string `json:"issue_text"`
		Severity   string `json:"issue_severity"`
		TestID     string `json:"test_id"`
	} `json:"results"`
}

func parseBanditOutput(output []byte) []RawFinding {
	var findings []RawFinding

	// Bandit outputs INFO lines before JSON, find the JSON start
	jsonStart := bytes.Index(output, []byte("{"))
	if jsonStart == -1 {
		return findings
	}
	output = output[jsonStart:]

	var result banditOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, r := range result.Results {
		findings = append(findings, RawFinding{
			FilePath:    r.Filename,
			LineNumber:  r.LineNumber,
			Description: r.IssueText,
			Severity:    strings.ToLower(r.Severity),
			RuleID:      r.TestID,
		})
	}

	return findings
}

func parsePipAuditOutput(output []byte) []RawFinding {
	var findings []RawFinding

	var result struct {
		Dependencies []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Vulns   []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
			} `json:"vulns"`
		} `json:"dependencies"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, dep := range result.Dependencies {
		for _, vuln := range dep.Vulns {
			findings = append(findings, RawFinding{
				FilePath:    "requirements.txt",
				Description: dep.Name + "@" + dep.Version + ": " + vuln.Description,
				Severity:    "high",
				RuleID:      vuln.ID,
			})
		}
	}

	return findings
}

func parseSafetyOutput(output []byte) []RawFinding {
	var findings []RawFinding

	// Safety outputs an array of vulnerability arrays
	var result [][]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, vuln := range result {
		if len(vuln) < 5 {
			continue
		}
		pkgName, _ := vuln[0].(string)
		vulnID, _ := vuln[4].(string)
		desc, _ := vuln[3].(string)

		findings = append(findings, RawFinding{
			FilePath:    "requirements.txt",
			Description: pkgName + ": " + desc,
			Severity:    "high",
			RuleID:      vulnID,
		})
	}

	return findings
}

func parseNpmAuditOutput(output []byte) []RawFinding {
	var findings []RawFinding

	var result struct {
		Vulnerabilities map[string]struct {
			Name     string        `json:"name"`
			Severity string        `json:"severity"`
			Via      []interface{} `json:"via"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, vuln := range result.Vulnerabilities {
		desc := "Vulnerability in " + vuln.Name
		findings = append(findings, RawFinding{
			FilePath:    "package.json",
			Description: desc,
			Severity:    strings.ToLower(vuln.Severity),
			RuleID:      vuln.Name,
		})
	}

	return findings
}

func parseCargoAuditOutput(output []byte) []RawFinding {
	var findings []RawFinding

	var result struct {
		Vulnerabilities struct {
			List []struct {
				Advisory struct {
					ID          string `json:"id"`
					Title       string `json:"title"`
					Description string `json:"description"`
				} `json:"advisory"`
				Package struct {
					Name    string `json:"name"`
					Version string `json:"version"`
				} `json:"package"`
			} `json:"list"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, vuln := range result.Vulnerabilities.List {
		findings = append(findings, RawFinding{
			FilePath:    "Cargo.toml",
			Description: vuln.Package.Name + "@" + vuln.Package.Version + ": " + vuln.Advisory.Title,
			Severity:    "high",
			RuleID:      vuln.Advisory.ID,
		})
	}

	return findings
}

func parseBundlerAuditOutput(output []byte) []RawFinding {
	var findings []RawFinding

	var result struct {
		Results []struct {
			Gem struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"gem"`
			Advisory struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			} `json:"advisory"`
		} `json:"results"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, r := range result.Results {
		findings = append(findings, RawFinding{
			FilePath:    "Gemfile.lock",
			Description: r.Gem.Name + "@" + r.Gem.Version + ": " + r.Advisory.Title,
			Severity:    "high",
			RuleID:      r.Advisory.ID,
		})
	}

	return findings
}

// brakemanOutput represents Brakeman JSON output structure.
type brakemanOutput struct {
	Warnings []struct {
		WarningType string `json:"warning_type"`
		Message     string `json:"message"`
		File        string `json:"file"`
		Line        int    `json:"line"`
		Confidence  string `json:"confidence"`
	} `json:"warnings"`
}

func parseBrakemanOutput(output []byte) []RawFinding {
	var findings []RawFinding
	var result brakemanOutput

	if err := json.Unmarshal(output, &result); err != nil {
		return findings
	}

	for _, w := range result.Warnings {
		var severity string
		switch w.Confidence {
		case "High":
			severity = "high"
		case "Weak":
			severity = "low"
		default:
			severity = "medium"
		}

		findings = append(findings, RawFinding{
			FilePath:    w.File,
			LineNumber:  w.Line,
			Description: w.WarningType + ": " + w.Message,
			Severity:    severity,
			RuleID:      w.WarningType,
		})
	}

	return findings
}
