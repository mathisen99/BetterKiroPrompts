# Requirements Document

## Introduction

This specification covers two related features for the BetterKiroPrompts hackathon submission - the final phase of the project:

1. **Info Page** - A page explaining what the site does, who it's for, and how to use it. This complements the existing Landing Page and Gallery Page.
2. **Security Scanning** - Repository security scanning with local tools and AI-powered code review for flagged issues. This implements the "Future Module — Repo Scanning" from The plan.md.

The existing application already has:
- Landing page with AI-driven prompt generation (project idea → questions → outputs)
- Experience level selection (beginner, novice, expert)
- Gallery page with ratings and view counts
- Generation storage in PostgreSQL
- Rate limiting and security hardening

The security scanning feature is designed for the open-source self-hosted use case where users provide their own API keys. It runs local security tools first, then uses GPT-5.1-Codex-Max only when issues are found to provide targeted remediation guidance.

## Glossary

- **Info_Page**: A frontend page explaining the site's purpose, features, and target audience
- **Security_Scan_Page**: A frontend page for initiating and viewing repository security scans
- **Security_Scanner**: The backend service that orchestrates repository scanning
- **Security_Container**: A dedicated Docker container running local security scanning tools (Trivy, Semgrep, TruffleHog, Gitleaks, and language-specific tools)
- **Scan_Job**: A queued task representing a repository security scan request
- **Finding**: A security issue detected by local scanning tools with severity, file path, line number, and description
- **Code_Review**: AI-generated analysis of flagged files using GPT-5.1-Codex-Max
- **GitHub_Token**: Optional personal access token for scanning private repositories (configured in .env)
- **Language_Detector**: Component that identifies the primary programming language(s) in a repository

## Requirements

### Requirement 1: Info Page Content

**User Story:** As a visitor, I want to see an info page explaining what BetterKiroPrompts does, so that I can understand if this tool is right for me.

#### Acceptance Criteria

1. WHEN a user navigates to the info page, THE Info_Page SHALL display a clear explanation of the site's purpose: helping developers think through projects before coding
2. WHEN displaying the info page, THE Info_Page SHALL explain the problem it solves: beginners often "vibe-code" without understanding architecture, security, data, or concurrency
3. WHEN displaying the info page, THE Info_Page SHALL explain that the tool is primarily for beginners to avoid bad initial prompt creation
4. WHEN displaying the info page, THE Info_Page SHALL explain that experienced users can also benefit from the advanced options (expert-level questions)
5. WHEN displaying the info page, THE Info_Page SHALL list the main features with brief descriptions:
   - Kickoff prompt generation (forces answer-first, no-coding-first thinking)
   - Steering file generation (creates .kiro/steering/ files with proper frontmatter)
   - Hooks generation (creates .kiro/hooks/ files from presets)
   - Security scanning (scans repos for vulnerabilities with AI remediation)
6. WHEN displaying the info page, THE Info_Page SHALL explain that this is an open-source project users can self-host with their own API keys
7. WHEN displaying the info page, THE Info_Page SHALL include a "Get Started" call-to-action button linking to the main landing page

### Requirement 2: Info Page Navigation

**User Story:** As a user, I want easy navigation to and from the info page, so that I can access information without losing my place.

#### Acceptance Criteria

1. THE Info_Page SHALL be accessible from the main landing page via a visible "About" or "Info" link in the header area
2. THE Info_Page SHALL be accessible from the gallery page via the same navigation element
3. WHEN on the info page, THE Info_Page SHALL provide a way to return to the main application (logo click or explicit link)
4. WHEN on the info page, THE Info_Page SHALL provide a way to navigate to the gallery
5. THE Info_Page SHALL use the same visual styling (NightSkyBackground, dark theme, blue accents) as the rest of the application

### Requirement 3: Security Scan Page

**User Story:** As a user, I want a dedicated page for security scanning, so that I can scan repositories without interfering with the prompt generation flow.

#### Acceptance Criteria

1. THE Security_Scan_Page SHALL be accessible from the info page and main navigation
2. WHEN on the security scan page, THE System SHALL display an input field for repository URL
3. WHEN on the security scan page, THE System SHALL display a brief explanation of what the scan does
4. WHEN on the security scan page, THE System SHALL indicate whether private repo scanning is available (based on GitHub_Token configuration)
5. THE Security_Scan_Page SHALL use the same visual styling as the rest of the application

### Requirement 4: Security Scan Initiation

**User Story:** As a user, I want to submit a repository URL for security scanning, so that I can identify potential security issues in my code.

#### Acceptance Criteria

1. WHEN a user provides a valid public GitHub repository URL, THE Security_Scanner SHALL accept the scan request
2. WHEN a user provides a private repository URL and a GitHub_Token is configured in .env, THE Security_Scanner SHALL accept the scan request
3. WHEN a user provides a private repository URL and no GitHub_Token is configured, THE Security_Scanner SHALL reject the request with a clear error message explaining how to configure the token
4. IF an invalid repository URL is provided, THEN THE Security_Scanner SHALL return a validation error with format guidance
5. WHEN a scan is initiated, THE Security_Scanner SHALL create a Scan_Job and return a job identifier for status polling
6. THE Security_Scanner SHALL support GitHub URLs in formats: https://github.com/owner/repo and https://github.com/owner/repo.git

### Requirement 5: Repository Cloning

**User Story:** As a system operator, I want repositories cloned securely and temporarily, so that scans don't persist sensitive code.

#### Acceptance Criteria

1. WHEN a Scan_Job starts, THE Security_Scanner SHALL clone the repository to a temporary directory
2. THE Security_Scanner SHALL clone repositories in read-only mode
3. WHEN using a GitHub_Token for private repos, THE Security_Scanner SHALL use the token securely without logging it
4. WHEN a scan completes or fails, THE Security_Scanner SHALL delete the cloned repository
5. THE Security_Scanner SHALL enforce a maximum repository size limit to prevent resource exhaustion

### Requirement 6: Language Detection

**User Story:** As a user, I want the scanner to detect what languages my repository uses, so that appropriate security tools can be applied.

#### Acceptance Criteria

1. WHEN a repository is cloned, THE Language_Detector SHALL analyze file extensions and content to identify primary languages
2. THE Language_Detector SHALL identify at least: Go, JavaScript, TypeScript, Python, Java, Ruby, PHP, C, C++, Rust
3. WHEN multiple languages are detected, THE Language_Detector SHALL rank them by prevalence in the repository

### Requirement 7: Local Security Tool Execution

**User Story:** As a user, I want local security tools to scan my repository, so that common vulnerabilities are detected without AI costs.

#### Acceptance Criteria

1. WHEN a Scan_Job runs, THE Security_Container SHALL execute Trivy for comprehensive vulnerability scanning (dependencies, secrets, misconfigurations)
2. WHEN a Scan_Job runs, THE Security_Container SHALL execute Semgrep with security rulesets for SAST analysis
3. WHEN a Scan_Job runs, THE Security_Container SHALL execute TruffleHog for secret detection in git history
4. WHEN a Scan_Job runs, THE Security_Container SHALL execute Gitleaks for additional secret detection
5. WHEN Go code is detected, THE Security_Container SHALL execute govulncheck
6. WHEN Python code is detected, THE Security_Container SHALL execute bandit and pip-audit
7. WHEN JavaScript/TypeScript code is detected, THE Security_Container SHALL execute npm audit
8. WHEN Rust code is detected, THE Security_Container SHALL execute cargo audit
9. WHEN Ruby code is detected, THE Security_Container SHALL execute bundler-audit and brakeman (if Rails detected)
10. THE Security_Container SHALL enforce hard timeouts on all tool executions (default 5 minutes per tool)
11. IF a tool times out, THEN THE Security_Scanner SHALL continue with partial results from other tools
12. THE Security_Container SHALL have no outbound network access during scans except for vulnerability database updates

### Requirement 8: Finding Aggregation

**User Story:** As a user, I want all security findings consolidated into a single report, so that I can review issues efficiently.

#### Acceptance Criteria

1. WHEN local tools complete, THE Security_Scanner SHALL aggregate all Findings into a unified format
2. THE Security_Scanner SHALL deduplicate findings that appear in multiple tools
3. THE Security_Scanner SHALL rank findings by severity (critical, high, medium, low, info)
4. WHEN aggregating findings, THE Security_Scanner SHALL include file path, line number, tool source, and description for each Finding

### Requirement 9: AI Code Review for Flagged Files

**User Story:** As a user, I want AI-powered code review only for files with security issues, so that I get actionable remediation guidance without excessive costs.

#### Acceptance Criteria

1. WHEN local tools find zero Findings, THE Security_Scanner SHALL skip AI code review entirely
2. WHEN local tools find Findings, THE Security_Scanner SHALL extract the unique files referenced in those findings
3. THE Code_Review SHALL only analyze files that have associated Findings
4. THE Code_Review SHALL use GPT-5.1-Codex-Max for analysis
5. WHEN reviewing a file, THE Code_Review SHALL explain what the issue is and how to fix it
6. THE Code_Review SHALL provide concrete code examples showing the recommended fix
7. THE Security_Scanner SHALL enforce a maximum number of files to review per scan (default 10 files) to control costs

### Requirement 10: Scan Results Display

**User Story:** As a user, I want to view my scan results in a clear interface, so that I can understand and act on the findings.

#### Acceptance Criteria

1. WHEN a scan completes, THE Security_Scan_Page SHALL display findings grouped by severity
2. WHEN displaying a Finding, THE Security_Scan_Page SHALL show the file path, line number, and description
3. WHEN AI code review is available for a Finding, THE Security_Scan_Page SHALL display the remediation guidance with syntax-highlighted code examples
4. WHEN displaying results, THE Security_Scan_Page SHALL indicate which tool detected each Finding
5. WHEN a scan is in progress, THE Security_Scan_Page SHALL show a progress indicator with current status (cloning, scanning, reviewing)

### Requirement 11: Scan Job Management

**User Story:** As a user, I want to track my scan jobs, so that I can see progress and access past results.

#### Acceptance Criteria

1. THE Security_Scanner SHALL persist Scan_Job status and results to the database
2. WHEN a user requests scan status, THE Security_Scanner SHALL return current progress and any available results
3. THE Security_Scanner SHALL retain scan results for a configurable period before cleanup (default 7 days)
4. WHEN a scan fails, THE Security_Scanner SHALL record the error reason and make it available to the user

### Requirement 12: Security Container Configuration

**User Story:** As a system operator, I want the security scanning tools in a dedicated container, so that they are isolated from the main application.

#### Acceptance Criteria

1. THE Security_Container SHALL be defined in docker-compose.yml as a separate service
2. THE Security_Container SHALL include universal tools: Trivy, Semgrep, TruffleHog, Gitleaks
3. THE Security_Container SHALL include Go tools: govulncheck
4. THE Security_Container SHALL include Python tools: bandit, pip-audit, safety
5. THE Security_Container SHALL include JavaScript/TypeScript tools: npm audit (via Node.js)
6. THE Security_Container SHALL include Rust tools: cargo audit
7. THE Security_Container SHALL include Ruby tools: bundler-audit, brakeman
8. THE Security_Container SHALL share a volume with the backend for repository access
9. THE Security_Container SHALL have resource limits (CPU, memory) configured
10. THE Security_Container SHALL only start when needed (not always running) using Docker Compose profiles

### Requirement 13: Environment Configuration

**User Story:** As a self-hosting user, I want to configure API keys and tokens via environment variables, so that I can use my own credentials.

#### Acceptance Criteria

1. THE Security_Scanner SHALL read the OpenAI API key from environment variables for Code_Review (existing OPENAI_API_KEY)
2. THE Security_Scanner SHALL read an optional GitHub_Token from environment variable GITHUB_TOKEN for private repo access
3. WHEN the OpenAI API key is not configured, THE Security_Scanner SHALL still run local tools but skip AI code review
4. THE .env.example file SHALL document all security scanning related environment variables with descriptions
