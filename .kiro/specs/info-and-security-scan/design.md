# Design Document

## Overview

This design covers two features for the BetterKiroPrompts hackathon submission:

1. **Info Page** - A new frontend page explaining the site's purpose, features, and target audience
2. **Security Scanning** - A complete repository security scanning system with local tools and AI-powered code review

The security scanning system follows a cost-efficient approach:
1. Clone repository to temporary directory
2. Detect programming languages in the repository
3. Run universal security tools (Trivy, Semgrep, TruffleHog, Gitleaks)
4. Run language-specific tools based on detected languages (govulncheck, bandit, npm audit, cargo audit, etc.)
5. Aggregate and deduplicate findings from all tools
6. Only if findings exist, send flagged files to GPT-5.1-Codex-Max for remediation guidance
7. Present results with actionable fixes

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Frontend                                    │
├─────────────────┬─────────────────┬─────────────────┬───────────────────┤
│   LandingPage   │   GalleryPage   │    InfoPage     │  SecurityScanPage │
│   (existing)    │   (existing)    │     (new)       │      (new)        │
└────────┬────────┴────────┬────────┴────────┬────────┴─────────┬─────────┘
         │                 │                 │                  │
         └─────────────────┴─────────────────┴──────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           Backend (Go)                                   │
├─────────────────────────────────────────────────────────────────────────┤
│  /api/health          - Health check (existing)                         │
│  /api/generate/*      - Prompt generation (existing)                    │
│  /api/gallery/*       - Gallery endpoints (existing)                    │
│  /api/scan            - POST: Start scan job                            │
│  /api/scan/{id}       - GET: Get scan status/results                    │
│  /api/scan/config     - GET: Get scan configuration (private repo avail)│
└────────┬────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Security Scanner Service                          │
├─────────────────────────────────────────────────────────────────────────┤
│  RepoCloner       - Clone repos (public/private with token)             │
│  LanguageDetector - Detect languages by file extension                  │
│  ToolRunner       - Execute security tools via Security Container       │
│  FindingAggregator - Dedupe and rank findings from all tools            │
│  CodeReviewer     - GPT-5.1-Codex-Max for remediation                   │
└────────┬────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     Security Container (Docker)                          │
├─────────────────────────────────────────────────────────────────────────┤
│  UNIVERSAL TOOLS (always run):                                          │
│    Trivy      - Comprehensive vuln scanner (deps, secrets, IaC)         │
│    Semgrep    - SAST with security rulesets (30+ languages)             │
│    TruffleHog - Secret detection in git history                         │
│    Gitleaks   - Additional secret detection patterns                    │
├─────────────────────────────────────────────────────────────────────────┤
│  LANGUAGE-SPECIFIC TOOLS (run based on detected languages):             │
│    Go:         govulncheck                                              │
│    Python:     bandit, pip-audit, safety                                │
│    JavaScript: npm audit                                                │
│    Rust:       cargo audit                                              │
│    Ruby:       bundler-audit, brakeman                                  │
└─────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        PostgreSQL (existing)                             │
├─────────────────────────────────────────────────────────────────────────┤
│  scan_jobs      - Scan job status and metadata                          │
│  scan_findings  - Individual findings with remediation                  │
└─────────────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Frontend Components

#### InfoPage Component
Location: `frontend/src/pages/InfoPage.tsx`

```typescript
interface InfoPageProps {
  onNavigateHome: () => void;
  onNavigateGallery: () => void;
  onNavigateScan: () => void;
}
```

The InfoPage displays:
- Hero section with site purpose
- Problem statement (vibe-coding without thinking)
- Feature cards for each capability
- Self-hosting explanation
- Call-to-action buttons

#### SecurityScanPage Component
Location: `frontend/src/pages/SecurityScanPage.tsx`

```typescript
interface SecurityScanPageProps {
  onNavigateHome: () => void;
  onNavigateGallery: () => void;
}

interface ScanConfig {
  privateRepoEnabled: boolean;
}

interface ScanJob {
  id: string;
  status: 'pending' | 'cloning' | 'scanning' | 'reviewing' | 'completed' | 'failed';
  repoUrl: string;
  languages: string[];
  findings: Finding[];
  error?: string;
  createdAt: string;
  completedAt?: string;
}

interface Finding {
  id: string;
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info';
  tool: string;
  filePath: string;
  lineNumber?: number;
  description: string;
  remediation?: string;
  codeExample?: string;
}
```

### Backend Components

#### Scanner Service
Location: `backend/internal/scanner/service.go`

```go
type Service struct {
    db           *sql.DB
    openaiClient *openai.Client
    githubToken  string
    maxFileSize  int64
    maxFiles     int
    toolTimeout  time.Duration
}

type ScanRequest struct {
    RepoURL string `json:"repo_url"`
}

type ScanJob struct {
    ID          string    `json:"id"`
    Status      string    `json:"status"`
    RepoURL     string    `json:"repo_url"`
    Languages   []string  `json:"languages"`
    Findings    []Finding `json:"findings"`
    Error       string    `json:"error,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Finding struct {
    ID          string  `json:"id"`
    Severity    string  `json:"severity"`
    Tool        string  `json:"tool"`
    FilePath    string  `json:"file_path"`
    LineNumber  *int    `json:"line_number,omitempty"`
    Description string  `json:"description"`
    Remediation string  `json:"remediation,omitempty"`
    CodeExample string  `json:"code_example,omitempty"`
}
```

#### Repo Cloner
Location: `backend/internal/scanner/cloner.go`

```go
type Cloner struct {
    githubToken string
    maxSize     int64
    tempDir     string
}

func (c *Cloner) Clone(ctx context.Context, repoURL string) (string, error)
func (c *Cloner) Cleanup(path string) error
func (c *Cloner) ValidateURL(repoURL string) error
```

#### Language Detector
Location: `backend/internal/scanner/language.go`

```go
type LanguageDetector struct{}

type LanguageResult struct {
    Language string
    FileCount int
    Percentage float64
}

func (d *LanguageDetector) Detect(repoPath string) ([]LanguageResult, error)
```

#### Tool Runner
Location: `backend/internal/scanner/tools.go`

```go
type ToolRunner struct {
    timeout time.Duration
}

type ToolResult struct {
    Tool     string
    Findings []RawFinding
    Error    error
    TimedOut bool
}

type RawFinding struct {
    FilePath    string
    LineNumber  int
    Description string
    Severity    string
    RuleID      string
}

// Universal tools (always run)
func (r *ToolRunner) RunTrivy(ctx context.Context, repoPath string) ToolResult
func (r *ToolRunner) RunSemgrep(ctx context.Context, repoPath string, languages []string) ToolResult
func (r *ToolRunner) RunTruffleHog(ctx context.Context, repoPath string) ToolResult
func (r *ToolRunner) RunGitleaks(ctx context.Context, repoPath string) ToolResult

// Language-specific tools
func (r *ToolRunner) RunGovulncheck(ctx context.Context, repoPath string) ToolResult      // Go
func (r *ToolRunner) RunBandit(ctx context.Context, repoPath string) ToolResult           // Python
func (r *ToolRunner) RunPipAudit(ctx context.Context, repoPath string) ToolResult         // Python
func (r *ToolRunner) RunSafety(ctx context.Context, repoPath string) ToolResult           // Python
func (r *ToolRunner) RunNpmAudit(ctx context.Context, repoPath string) ToolResult         // JavaScript/TypeScript
func (r *ToolRunner) RunCargoAudit(ctx context.Context, repoPath string) ToolResult       // Rust
func (r *ToolRunner) RunBundlerAudit(ctx context.Context, repoPath string) ToolResult     // Ruby
func (r *ToolRunner) RunBrakeman(ctx context.Context, repoPath string) ToolResult         // Ruby/Rails

// Tool selection based on detected languages
func (r *ToolRunner) GetToolsForLanguages(languages []string) []string
```

#### Finding Aggregator
Location: `backend/internal/scanner/aggregator.go`

```go
type Aggregator struct{}

func (a *Aggregator) Aggregate(results []ToolResult) []Finding
func (a *Aggregator) Deduplicate(findings []Finding) []Finding
func (a *Aggregator) RankBySeverity(findings []Finding) []Finding
```

#### Code Reviewer
Location: `backend/internal/scanner/reviewer.go`

```go
type CodeReviewer struct {
    client   *openai.Client
    maxFiles int
}

func (r *CodeReviewer) Review(ctx context.Context, repoPath string, findings []Finding) ([]Finding, error)
```

### API Endpoints

#### POST /api/scan
Start a new security scan.

Request:
```json
{
  "repo_url": "https://github.com/owner/repo"
}
```

Response (202 Accepted):
```json
{
  "id": "scan_abc123",
  "status": "pending",
  "repo_url": "https://github.com/owner/repo"
}
```

#### GET /api/scan/{id}
Get scan status and results.

Response:
```json
{
  "id": "scan_abc123",
  "status": "completed",
  "repo_url": "https://github.com/owner/repo",
  "languages": ["go", "typescript"],
  "findings": [
    {
      "id": "finding_1",
      "severity": "high",
      "tool": "gitleaks",
      "file_path": "config/secrets.go",
      "line_number": 42,
      "description": "Hardcoded API key detected",
      "remediation": "Move the API key to environment variables...",
      "code_example": "// Before:\nconst apiKey = \"sk-...\"\n\n// After:\napiKey := os.Getenv(\"API_KEY\")"
    }
  ],
  "created_at": "2026-01-13T10:00:00Z",
  "completed_at": "2026-01-13T10:02:30Z"
}
```

#### GET /api/scan/config
Get scan configuration (whether private repos are supported).

Response:
```json
{
  "private_repo_enabled": true
}
```

## Data Models

### Database Schema

```sql
-- Migration: 005_create_scan_tables.sql

CREATE TABLE scan_jobs (
    id VARCHAR(36) PRIMARY KEY,
    repo_url TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    languages TEXT[], -- PostgreSQL array
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '7 days'
);

CREATE INDEX idx_scan_jobs_status ON scan_jobs(status);
CREATE INDEX idx_scan_jobs_expires_at ON scan_jobs(expires_at);

CREATE TABLE scan_findings (
    id VARCHAR(36) PRIMARY KEY,
    scan_job_id VARCHAR(36) NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
    severity VARCHAR(10) NOT NULL,
    tool VARCHAR(50) NOT NULL,
    file_path TEXT NOT NULL,
    line_number INTEGER,
    description TEXT NOT NULL,
    remediation TEXT,
    code_example TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_scan_findings_job_id ON scan_findings(scan_job_id);
CREATE INDEX idx_scan_findings_severity ON scan_findings(severity);
```

## Security Container

The security container includes a comprehensive set of tools to scan repositories in multiple languages. Tools are selected based on detected languages.

### Tool Matrix by Language

| Language | Dependency Scanner | SAST/Vulnerability | Secret Detection |
|----------|-------------------|-------------------|------------------|
| All | Trivy (universal) | Semgrep (universal) | TruffleHog, Gitleaks |
| Go | govulncheck | Semgrep Go rules | - |
| Python | pip-audit, safety | Bandit | - |
| JavaScript/TypeScript | npm audit | Semgrep JS rules | - |
| Rust | cargo audit | Semgrep Rust rules | - |
| Java | Trivy | Semgrep Java rules | - |
| Ruby | bundler-audit | Brakeman | - |
| PHP | Trivy | Semgrep PHP rules | - |
| C/C++ | Trivy | Semgrep C rules | - |

### Universal Tools (Always Run)

1. **Trivy** - Comprehensive vulnerability scanner for dependencies, containers, IaC, and secrets
   - Supports: Go, Python, Node.js, Ruby, PHP, Java, Rust, .NET, and more
   - Scans: Dependencies, misconfigurations, secrets, licenses
   
2. **Semgrep** - Fast, lightweight SAST with rules for 30+ languages
   - Uses community rulesets (p/security-audit, p/owasp-top-ten)
   - Language-agnostic pattern matching

3. **TruffleHog** - Secret detection across git history
   - Detects: API keys, passwords, tokens, certificates

4. **Gitleaks** - Additional secret detection
   - Complements TruffleHog with different detection patterns

### Language-Specific Tools

**Go:**
- `govulncheck` - Official Go vulnerability checker

**Python:**
- `bandit` - Python SAST for common security issues
- `pip-audit` - Dependency vulnerability scanner
- `safety` - Additional dependency checker

**JavaScript/TypeScript:**
- `npm audit` - Built-in npm vulnerability scanner

**Rust:**
- `cargo audit` - Rust dependency vulnerability scanner

**Ruby:**
- `bundler-audit` - Ruby gem vulnerability scanner
- `brakeman` - Rails-specific security scanner

### Dockerfile
Location: `backend/Dockerfile.scanner`

```dockerfile
FROM golang:1.25-alpine AS go-builder

# Install Go-based tools
RUN go install golang.org/x/vuln/cmd/govulncheck@latest
RUN go install github.com/google/osv-scanner/cmd/osv-scanner@latest

FROM python:3.12-alpine AS python-builder

# Install Python-based tools
RUN pip install --no-cache-dir bandit pip-audit safety

FROM rust:1.75-alpine AS rust-builder

# Install Rust-based tools
RUN cargo install cargo-audit

FROM alpine:3.19

# Install base dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    curl \
    nodejs \
    npm \
    ruby \
    ruby-bundler \
    python3 \
    py3-pip

# Install Trivy
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin

# Install Semgrep
RUN pip3 install --no-cache-dir semgrep

# Install TruffleHog
RUN curl -sSfL https://raw.githubusercontent.com/trufflesecurity/trufflehog/main/scripts/install.sh | sh -s -- -b /usr/local/bin

# Install Gitleaks
RUN curl -sSfL https://github.com/gitleaks/gitleaks/releases/download/v8.21.0/gitleaks_8.21.0_linux_amd64.tar.gz | tar xz -C /usr/local/bin

# Install bundler-audit and brakeman for Ruby
RUN gem install bundler-audit brakeman --no-document

# Copy Go tools
COPY --from=go-builder /go/bin/govulncheck /usr/local/bin/
COPY --from=go-builder /go/bin/osv-scanner /usr/local/bin/

# Copy Python tools
COPY --from=python-builder /usr/local/bin/bandit /usr/local/bin/
COPY --from=python-builder /usr/local/bin/pip-audit /usr/local/bin/
COPY --from=python-builder /usr/local/bin/safety /usr/local/bin/

# Copy Rust tools
COPY --from=rust-builder /usr/local/cargo/bin/cargo-audit /usr/local/bin/

# Create non-root user
RUN adduser -D -u 1000 scanner
USER scanner

WORKDIR /scan

ENTRYPOINT ["/bin/sh"]
```

### Docker Compose Configuration
Addition to `docker-compose.yml`:

```yaml
services:
  # ... existing services ...
  
  scanner:
    build:
      context: ./backend
      dockerfile: Dockerfile.scanner
    profiles:
      - scan
    volumes:
      - scan-repos:/scan/repos
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
    networks:
      - app-network

volumes:
  scan-repos:
```

## GPT-5.1-Codex-Max Integration

### Code Review Prompt

Based on the GPT-5.1-Codex-Max prompting guide, the code review uses a focused prompt:

```go
const codeReviewSystemPrompt = `You are a security code reviewer. Your task is to analyze code files that have been flagged by security scanning tools and provide actionable remediation guidance.

For each finding:
1. Explain what the security issue is in plain language
2. Explain why it's a problem (potential impact)
3. Provide a concrete code fix with before/after examples
4. Keep explanations concise and actionable

Format your response as JSON:
{
  "findings": [
    {
      "file_path": "path/to/file",
      "line_number": 42,
      "remediation": "Clear explanation of the fix",
      "code_example": "// Before:\n...\n\n// After:\n..."
    }
  ]
}

Focus on practical fixes. Do not invent new vulnerabilities - only address the specific issues flagged.`
```

### Request Structure

```go
func (r *CodeReviewer) buildRequest(findings []Finding, fileContents map[string]string) openai.ChatCompletionRequest {
    var userContent strings.Builder
    userContent.WriteString("Review these security findings and provide remediation:\n\n")
    
    for _, f := range findings {
        userContent.WriteString(fmt.Sprintf("## Finding in %s (line %d)\n", f.FilePath, f.LineNumber))
        userContent.WriteString(fmt.Sprintf("Tool: %s\n", f.Tool))
        userContent.WriteString(fmt.Sprintf("Issue: %s\n\n", f.Description))
        
        if content, ok := fileContents[f.FilePath]; ok {
            userContent.WriteString("```\n")
            userContent.WriteString(content)
            userContent.WriteString("\n```\n\n")
        }
    }
    
    return openai.ChatCompletionRequest{
        Model: "gpt-5.1-codex-max",
        Messages: []openai.ChatCompletionMessage{
            {Role: "system", Content: codeReviewSystemPrompt},
            {Role: "user", Content: userContent.String()},
        },
        Temperature: 0.2, // Low temperature for consistent, focused output
        MaxTokens:   4096,
    }
}
```


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: URL Validation

*For any* string input to the URL validator:
- If it matches a valid GitHub repository URL pattern (https://github.com/owner/repo or https://github.com/owner/repo.git), the validator SHALL accept it
- If it does not match a valid pattern, the validator SHALL reject it with an appropriate error

**Validates: Requirements 4.1, 4.4, 4.6**

### Property 2: Job Creation Round-Trip

*For any* valid scan request, the Security_Scanner SHALL:
- Create a Scan_Job with a unique identifier
- Return the job identifier to the caller
- The job SHALL be retrievable by that identifier

**Validates: Requirements 4.5, 11.1, 11.2**

### Property 3: Token Security

*For any* scan operation using a GitHub token, the token value SHALL NOT appear in:
- Application logs
- Error messages returned to users
- Database records

**Validates: Requirements 5.3**

### Property 4: Repository Cleanup

*For any* scan job that reaches a terminal state (completed or failed), the cloned repository directory SHALL be deleted within the cleanup phase.

**Validates: Requirements 5.4**

### Property 5: Repository Size Limit

*For any* repository that exceeds the configured maximum size limit, the Security_Scanner SHALL reject the scan request with an appropriate error before completing the clone.

**Validates: Requirements 5.5**

### Property 6: Language Detection Accuracy

*For any* repository containing files with known extensions:
- Files with .go extension SHALL be identified as Go
- Files with .js extension SHALL be identified as JavaScript
- Files with .ts extension SHALL be identified as TypeScript
- Files with .py extension SHALL be identified as Python
- Files with .java extension SHALL be identified as Java
- Files with .rb extension SHALL be identified as Ruby
- Files with .php extension SHALL be identified as PHP
- Files with .c/.h extension SHALL be identified as C
- Files with .cpp/.hpp extension SHALL be identified as C++
- Files with .rs extension SHALL be identified as Rust

*For any* repository with multiple languages, the results SHALL be sorted by file count in descending order.

**Validates: Requirements 6.1, 6.2, 6.3**

### Property 7: Tool Timeout Enforcement

*For any* security tool execution, if the tool does not complete within the configured timeout period, the execution SHALL be terminated and the scanner SHALL continue with results from other tools.

**Validates: Requirements 7.10, 7.11**

### Property 8: Finding Aggregation Completeness

*For any* set of tool results aggregated into findings:
- Each finding SHALL have a non-empty file_path
- Each finding SHALL have a non-empty description
- Each finding SHALL have a valid severity (critical, high, medium, low, info)
- Each finding SHALL have a tool source identifier
- Duplicate findings (same file, line, description) SHALL be deduplicated
- Findings SHALL be sorted by severity (critical first, info last)

**Validates: Requirements 8.1, 8.2, 8.3, 8.4**

### Property 9: AI Review Scope Limitation

*For any* scan with findings:
- The Code_Review SHALL only receive files that have at least one associated finding
- The number of files sent to Code_Review SHALL NOT exceed the configured maximum (default 10)
- If there are more flagged files than the maximum, only the files with highest-severity findings SHALL be reviewed

**Validates: Requirements 9.2, 9.3, 9.7**

### Property 10: AI Review Content Quality

*For any* finding that receives AI code review:
- The remediation field SHALL contain a non-empty explanation
- The code_example field SHALL contain a non-empty code snippet showing the fix

**Validates: Requirements 9.5, 9.6**

### Property 11: Finding Display Completeness

*For any* finding displayed in the UI:
- The file path SHALL be visible
- The description SHALL be visible
- The tool source SHALL be visible
- If line_number is present, it SHALL be visible
- If remediation is present, it SHALL be visible

**Validates: Requirements 10.2, 10.3, 10.4**

### Property 12: Scan Result Retention

*For any* completed scan job, the results SHALL remain retrievable until the configured expiry time (default 7 days) has passed.

**Validates: Requirements 11.3**

### Property 13: Error Recording

*For any* scan job that fails, the error field SHALL contain a non-empty description of the failure reason.

**Validates: Requirements 11.4**

## Error Handling

### Repository Cloning Errors

| Error Condition | Response | User Message |
|-----------------|----------|--------------|
| Invalid URL format | 400 Bad Request | "Invalid repository URL. Please use format: https://github.com/owner/repo" |
| Repository not found | 404 Not Found | "Repository not found. Please check the URL and try again." |
| Private repo without token | 403 Forbidden | "Private repository access requires a GitHub token. See documentation for setup." |
| Repository too large | 413 Payload Too Large | "Repository exceeds maximum size limit of {limit}MB." |
| Clone timeout | 504 Gateway Timeout | "Repository clone timed out. Please try again or use a smaller repository." |

### Tool Execution Errors

| Error Condition | Behavior |
|-----------------|----------|
| Tool not found | Log error, skip tool, continue with other tools |
| Tool timeout | Terminate tool, log timeout, continue with partial results |
| Tool crash | Log error with exit code, continue with other tools |
| All tools fail | Mark scan as failed with aggregated error message |

### AI Review Errors

| Error Condition | Behavior |
|-----------------|----------|
| OpenAI API key not configured | Skip AI review, return findings without remediation |
| OpenAI API error | Log error, return findings without remediation |
| OpenAI rate limit | Retry with exponential backoff (max 3 attempts) |
| Response parsing error | Log error, return findings without remediation for affected files |

## Testing Strategy

### Unit Tests

Unit tests focus on specific examples and edge cases:

1. **URL Validation**
   - Valid GitHub URLs (with/without .git suffix)
   - Invalid URLs (wrong domain, missing parts, malformed)
   - Edge cases (trailing slashes, query params)

2. **Language Detection**
   - Single-language repositories
   - Multi-language repositories
   - Empty repositories
   - Repositories with unknown extensions

3. **Finding Aggregation**
   - Empty results
   - Single tool results
   - Multiple tool results with duplicates
   - Severity sorting

4. **API Handlers**
   - Request validation
   - Error responses
   - Success responses

### Property-Based Tests

Property-based tests verify universal properties across many generated inputs using a PBT library (e.g., `gopter` for Go, `fast-check` for TypeScript):

1. **Property 1: URL Validation** - Generate random strings, verify valid patterns accepted, invalid rejected
2. **Property 6: Language Detection** - Generate file trees with known extensions, verify correct identification
3. **Property 8: Finding Aggregation** - Generate random tool results, verify output format and ordering
4. **Property 9: AI Review Scope** - Generate findings lists, verify file count limits respected

Configuration:
- Minimum 100 iterations per property test
- Each test tagged with: **Feature: info-and-security-scan, Property {number}: {property_text}**

### Integration Tests

1. **End-to-end scan flow** - Submit scan, poll status, verify results
2. **Container communication** - Backend to security container tool execution
3. **Database persistence** - Scan job and findings storage/retrieval

### Frontend Tests

1. **InfoPage rendering** - Verify all content sections present
2. **SecurityScanPage flow** - Submit URL, show progress, display results
3. **Navigation** - Links between pages work correctly
