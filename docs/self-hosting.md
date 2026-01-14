# Self-Hosting Guide

This guide covers everything you need to deploy and configure BetterKiroPrompts for your team or personal use.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start](#quick-start)
3. [Configuration Reference](#configuration-reference)
4. [Example Configurations](#example-configurations)
5. [Database Setup](#database-setup)
6. [OpenAI Configuration](#openai-configuration)
7. [Private Repository Scanning](#private-repository-scanning)
8. [Resource Requirements](#resource-requirements)
9. [Maintenance](#maintenance)
10. [Scanner Customization](#scanner-customization)
11. [Troubleshooting](#troubleshooting)

---

## Prerequisites

Before deploying BetterKiroPrompts, ensure you have:

| Requirement | Version | Notes |
|-------------|---------|-------|
| Docker Engine | 29.1.3+ | Required for containerized deployment |
| Docker Compose | 5.0.1+ | For orchestrating services |
| OpenAI API Key | - | Required for AI generation features |
| GitHub Token | - | Optional, only for private repo scanning |

### Hardware Requirements

See [Resource Requirements](#resource-requirements) for detailed recommendations.

---

## Quick Start

1. **Clone the repository**
   ```bash
   git clone <repo-url>
   cd Better-Kiro-Prompts
   ```

2. **Configure secrets**
   ```bash
   cp .env.example .env
   # Edit .env and set your OPENAI_API_KEY
   ```

3. **Configure settings (optional)**
   ```bash
   cp config.example.toml config.toml
   # Edit config.toml to customize settings
   ```

4. **Start the application**
   ```bash
   ./build.sh
   ```

5. **Access the application**
   
   Open http://localhost:8080 in your browser.

---

## Configuration Reference

BetterKiroPrompts uses a two-file configuration approach:

- **`.env`** - Secrets (API keys, tokens, database credentials)
- **`config.toml`** - Application settings (ports, timeouts, limits)

### Configuration Precedence

1. Environment variables (highest priority)
2. `config.toml` values
3. Built-in defaults (lowest priority)

### Server Configuration

| Option | Type | Default | Range | Description |
|--------|------|---------|-------|-------------|
| `server.port` | int | `8080` | 1-65535 | HTTP server port |
| `server.host` | string | `"0.0.0.0"` | - | Bind address (`0.0.0.0` for all interfaces) |
| `server.shutdown_timeout` | duration | `"30s"` | ≥1s | Graceful shutdown timeout |

**Environment overrides:** `PORT`

### OpenAI Configuration

| Option | Type | Default | Valid Values | Description |
|--------|------|---------|--------------|-------------|
| `openai.model` | string | `"gpt-5.2"` | Any OpenAI model | Model for question/output generation |
| `openai.code_review_model` | string | `"gpt-5.1-codex-max"` | Any OpenAI model | Model for security code review |
| `openai.base_url` | string | `"https://api.openai.com/v1"` | Valid URL | API endpoint (for Azure/proxies) |
| `openai.timeout` | duration | `"180s"` | ≥10s | Request timeout |
| `openai.reasoning_effort` | string | `"medium"` | `none`, `low`, `medium`, `high`, `xhigh` | AI reasoning depth |
| `openai.verbosity` | string | `"medium"` | `low`, `medium`, `high` | Output detail level |

**Environment overrides:** `OPENAI_MODEL`

> **Note:** `OPENAI_API_KEY` must be set in `.env`, not in `config.toml`.

### Rate Limiting Configuration

| Option | Type | Default | Range | Description |
|--------|------|---------|-------|-------------|
| `rate_limit.generation_limit_per_hour` | int | `10` | ≥1 | Max AI generations per IP per hour |
| `rate_limit.rating_limit_per_hour` | int | `20` | ≥1 | Max gallery ratings per IP per hour |
| `rate_limit.scan_limit_per_hour` | int | `10` | ≥1 | Max security scans per IP per hour |

**Environment overrides:** `RATE_LIMIT_GENERATION`, `RATE_LIMIT_RATING`, `RATE_LIMIT_SCAN`

### Logging Configuration

| Option | Type | Default | Valid Values | Description |
|--------|------|---------|--------------|-------------|
| `logging.level` | string | `"INFO"` | `DEBUG`, `INFO`, `WARN`, `ERROR` | Log level threshold |
| `logging.directory` | string | `"./logs"` | Valid path | Log file directory |
| `logging.max_size_mb` | int | `100` | ≥1 | Max log file size before rotation |
| `logging.max_age_days` | int | `7` | ≥1 | Days to retain log files |
| `logging.enable_color` | bool | `true` | - | Colored console output |

**Environment overrides:** `LOG_LEVEL`

### Scanner Configuration

| Option | Type | Default | Range | Description |
|--------|------|---------|-------|-------------|
| `scanner.max_repo_size_mb` | int | `500` | ≥1 | Max repository size to clone |
| `scanner.max_review_files` | int | `10` | ≥1 | Max files for AI code review |
| `scanner.tool_timeout_seconds` | int | `300` | ≥10 | Timeout per security tool |
| `scanner.retention_days` | int | `7` | ≥1 | Days to retain scan results |
| `scanner.clone_timeout` | duration | `"5m"` | ≥10s | Git clone timeout |

**Environment overrides:** `SCANNER_MAX_REPO_SIZE_MB`, `SCANNER_MAX_REVIEW_FILES`, `SCANNER_TOOL_TIMEOUT_SECONDS`, `SCANNER_RESULT_RETENTION_DAYS`

### Generation Configuration

| Option | Type | Default | Range | Description |
|--------|------|---------|-------|-------------|
| `generation.max_project_idea_length` | int | `2000` | ≥100 | Max project idea input length |
| `generation.max_answer_length` | int | `1000` | ≥100 | Max answer length per question |
| `generation.min_questions` | int | `5` | ≥1 | Minimum questions to generate |
| `generation.max_questions` | int | `10` | ≥min_questions | Maximum questions to generate |
| `generation.max_retries` | int | `1` | ≥0 | AI generation retry attempts |

### Gallery Configuration

| Option | Type | Default | Valid Values | Description |
|--------|------|---------|--------------|-------------|
| `gallery.page_size` | int | `20` | 1-100 | Items per page in listings |
| `gallery.default_sort` | string | `"newest"` | `newest`, `highest_rated`, `most_viewed` | Default sort order |

---

## Example Configurations

### Minimal Configuration (Generation Only)

For basic prompt generation without security scanning:

```toml
# config.toml - Minimal setup

[server]
port = 8080

[openai]
model = "gpt-5.2"
timeout = "120s"

[rate_limit]
generation_limit_per_hour = 20

[logging]
level = "INFO"
```

### Full-Featured Configuration

All features enabled with balanced settings:

```toml
# config.toml - Full-featured setup

[server]
port = 8080
host = "0.0.0.0"
shutdown_timeout = "30s"

[openai]
model = "gpt-5.2"
code_review_model = "gpt-5.1-codex-max"
base_url = "https://api.openai.com/v1"
timeout = "180s"
reasoning_effort = "medium"
verbosity = "medium"

[rate_limit]
generation_limit_per_hour = 10
rating_limit_per_hour = 20
scan_limit_per_hour = 10

[logging]
level = "INFO"
directory = "./logs"
max_size_mb = 100
max_age_days = 7
enable_color = false  # Disable for production log aggregation

[scanner]
max_repo_size_mb = 500
max_review_files = 10
tool_timeout_seconds = 300
retention_days = 7
clone_timeout = "5m"

[generation]
max_project_idea_length = 2000
max_answer_length = 1000
min_questions = 5
max_questions = 10
max_retries = 1

[gallery]
page_size = 20
default_sort = "newest"
```

### High-Security Configuration

Strict rate limits and conservative settings for public deployments:

```toml
# config.toml - High-security setup

[server]
port = 8080
host = "127.0.0.1"  # Local only, use reverse proxy
shutdown_timeout = "60s"

[openai]
model = "gpt-5.2"
code_review_model = "gpt-5.2"
timeout = "120s"
reasoning_effort = "low"  # Faster, cheaper
verbosity = "low"

[rate_limit]
generation_limit_per_hour = 3   # Very strict
rating_limit_per_hour = 10
scan_limit_per_hour = 2         # Expensive operation

[logging]
level = "WARN"  # Reduce log volume
directory = "./logs"
max_size_mb = 50
max_age_days = 30  # Longer retention for auditing
enable_color = false

[scanner]
max_repo_size_mb = 100   # Smaller repos only
max_review_files = 5     # Fewer files to review
tool_timeout_seconds = 120
retention_days = 3       # Shorter retention
clone_timeout = "2m"

[generation]
max_project_idea_length = 1000  # Shorter inputs
max_answer_length = 500
min_questions = 3
max_questions = 7
max_retries = 0  # No retries to save costs

[gallery]
page_size = 10
default_sort = "newest"
```

---

## Database Setup

BetterKiroPrompts uses PostgreSQL for persistent storage.

### Using Docker Compose (Recommended)

The included `docker-compose.yml` automatically sets up PostgreSQL:

```yaml
postgres:
  image: postgres:18.1
  environment:
    - POSTGRES_USER=user
    - POSTGRES_PASSWORD=pass
    - POSTGRES_DB=app
  volumes:
    - postgres_data:/var/lib/postgresql
```

Configure the connection in `.env`:

```bash
DATABASE_URL=postgres://user:pass@postgres:5432/app?sslmode=disable
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=app
```

### External PostgreSQL

To use an existing PostgreSQL instance:

1. Create a database:
   ```sql
   CREATE DATABASE betterkiro;
   CREATE USER betterkiro WITH PASSWORD 'your-secure-password';
   GRANT ALL PRIVILEGES ON DATABASE betterkiro TO betterkiro;
   ```

2. Update `.env`:
   ```bash
   DATABASE_URL=postgres://betterkiro:your-secure-password@your-host:5432/betterkiro?sslmode=require
   ```

3. Remove the `postgres` service from `docker-compose.yml` or use `docker-compose.prod.yml` with modifications.

### Migrations

Migrations run automatically on application startup. The backend reads SQL files from `backend/internal/db/migrations/` in alphabetical order.

Migration files follow the naming convention:
```
YYYYMMDDHHMMSS_description.sql
```

To manually run migrations:
```bash
# Connect to the database
docker exec -it better-kiro-prompts-postgres-1 psql -U user -d app

# View applied migrations (check schema_migrations table if exists)
```

---

## OpenAI Configuration

### Getting an API Key

1. Go to [OpenAI Platform](https://platform.openai.com/api-keys)
2. Create a new API key
3. Add it to your `.env` file:
   ```bash
   OPENAI_API_KEY=sk-your-api-key-here
   ```

### Model Selection

| Use Case | Recommended Model | Notes |
|----------|-------------------|-------|
| Standard | `gpt-5.2` | Current standard, recommended for most uses |
| Budget-conscious | `gpt-5-mini` | Good balance of quality and cost |
| Code review | `gpt-5.1-codex-max` | Optimized for code analysis |
| Fast responses | `gpt-5-nano` | Fastest, good for high-volume tasks |

Configure in `config.toml`:
```toml
[openai]
model = "gpt-5.2"                    # For generation
code_review_model = "gpt-5.2"        # For security scanning
```

### Using Azure OpenAI

To use Azure OpenAI instead of OpenAI directly:

1. Set the base URL in `config.toml`:
   ```toml
   [openai]
   base_url = "https://your-resource.openai.azure.com/openai/deployments/your-deployment"
   model = "your-deployment-name"
   ```

2. Set your Azure API key in `.env`:
   ```bash
   OPENAI_API_KEY=your-azure-api-key
   ```

### Cost Optimization

- Use `reasoning_effort = "low"` for faster, cheaper responses
- Use `verbosity = "low"` for shorter outputs
- Reduce `max_questions` to generate fewer questions
- Set `max_retries = 0` to avoid retry costs

---

## Private Repository Scanning

The security scanner can scan private GitHub repositories with proper authentication.

### Setting Up GitHub Token

1. Go to [GitHub Settings > Tokens](https://github.com/settings/tokens)
2. Generate a new token (classic) with these scopes:
   - `repo` - Full control of private repositories
3. Add to `.env`:
   ```bash
   GITHUB_TOKEN=ghp_your-token-here
   ```

### Required Permissions

| Scope | Required For |
|-------|--------------|
| `repo` | Cloning private repositories |
| `read:org` | Scanning organization repositories (optional) |

### Security Considerations

- The GitHub token is only used for cloning repositories
- Cloned repositories are stored temporarily and deleted after scanning
- Token is never logged or exposed in scan results
- Consider using a dedicated service account with minimal permissions

### Enabling the Scanner Container

The scanner runs in a separate container with security tools. It's included by default when using `./build.sh`:

```bash
# Build and start (includes scanner)
./build.sh
```

Or modify `docker-compose.prod.yml` to always include the scanner.

---

## Resource Requirements

### Minimum Requirements

| Component | CPU | Memory | Disk |
|-----------|-----|--------|------|
| Backend | 1 core | 512 MB | 1 GB |
| PostgreSQL | 1 core | 256 MB | 5 GB |
| Frontend (build) | 2 cores | 1 GB | 500 MB |
| **Total** | **2 cores** | **1.5 GB** | **6.5 GB** |

### Recommended (with Scanner)

| Component | CPU | Memory | Disk |
|-----------|-----|--------|------|
| Backend | 2 cores | 1 GB | 2 GB |
| PostgreSQL | 2 cores | 1 GB | 20 GB |
| Scanner | 2 cores | 2 GB | 10 GB |
| **Total** | **4 cores** | **4 GB** | **32 GB** |

### Scaling Considerations

**Horizontal Scaling:**
- Backend is stateless and can be load-balanced
- Use external PostgreSQL with connection pooling
- Scanner containers can run on separate hosts

**Vertical Scaling:**
- Increase scanner memory for large repositories
- Increase PostgreSQL memory for better query performance
- More CPU cores improve concurrent scan performance

**Storage:**
- Log files grow based on usage (~10-50 MB/day typical)
- Scan results are retained based on `scanner.retention_days`
- Database grows with gallery submissions

---

## Maintenance

### Backup Procedures

**Database Backup:**
```bash
# Create backup
docker exec better-kiro-prompts-postgres-1 \
  pg_dump -U user -d app > backup_$(date +%Y%m%d).sql

# Restore backup
docker exec -i better-kiro-prompts-postgres-1 \
  psql -U user -d app < backup_20260114.sql
```

**Automated Backups:**
```bash
# Add to crontab
0 2 * * * docker exec better-kiro-prompts-postgres-1 pg_dump -U user -d app | gzip > /backups/db_$(date +\%Y\%m\%d).sql.gz
```

### Log Rotation

Logs are automatically rotated based on `config.toml` settings:

```toml
[logging]
max_size_mb = 100    # Rotate when file exceeds 100 MB
max_age_days = 7     # Delete files older than 7 days
```

Log files are organized by category:
- `app.log` - Application events
- `http.log` - HTTP requests
- `db.log` - Database operations
- `scanner.log` - Security scan operations
- `client.log` - Client-side errors

**Manual cleanup:**
```bash
# Remove logs older than 7 days
find ./logs -name "*.log" -mtime +7 -delete
```

### Database Cleanup

**Scan Results:**
Scan results are automatically cleaned up based on `scanner.retention_days`.

**Manual cleanup:**
```sql
-- Delete old scan results
DELETE FROM scan_results WHERE created_at < NOW() - INTERVAL '30 days';

-- Delete orphaned scan files
DELETE FROM scan_files WHERE scan_id NOT IN (SELECT id FROM scan_results);

-- Vacuum to reclaim space
VACUUM ANALYZE;
```

**Gallery Cleanup:**
```sql
-- Remove generations with no views in 90 days
DELETE FROM generations 
WHERE id NOT IN (SELECT generation_id FROM views WHERE created_at > NOW() - INTERVAL '90 days')
AND created_at < NOW() - INTERVAL '90 days';
```

### Health Checks

The backend exposes a health endpoint:

```bash
curl http://localhost:8080/api/health
```

Response:
```json
{"status": "ok"}
```

### Updating

1. Pull latest changes:
   ```bash
   git pull origin main
   ```

2. Rebuild and restart:
   ```bash
   ./build.sh
   ```

3. Check logs for migration status:
   ```bash
   docker logs better-kiro-prompts-backend-1
   ```

---

## Scanner Customization

The security scanner runs multiple tools based on detected languages.

### Available Tools

**Universal Tools (always run):**

| Tool | Purpose | Documentation |
|------|---------|---------------|
| Trivy | Vulnerability scanning, secrets, misconfigs | [aquasecurity/trivy](https://github.com/aquasecurity/trivy) |
| Semgrep | SAST with OWASP rules | [semgrep.dev](https://semgrep.dev) |
| TruffleHog | Secret detection in git history | [trufflesecurity/trufflehog](https://github.com/trufflesecurity/trufflehog) |
| Gitleaks | Additional secret detection | [gitleaks/gitleaks](https://github.com/gitleaks/gitleaks) |

**Language-Specific Tools:**

| Language | Tools |
|----------|-------|
| Go | govulncheck |
| Python | bandit, pip-audit, safety |
| JavaScript/TypeScript | npm audit |
| Rust | cargo-audit |
| Ruby | bundler-audit, brakeman |

### Adding New Tools

1. **Update Dockerfile.scanner:**
   ```dockerfile
   # Add your tool installation
   RUN curl -sSfL https://example.com/install.sh | sh -s -- -b /usr/local/bin
   ```

2. **Add tool runner in `backend/internal/scanner/tools.go`:**
   ```go
   func (r *ToolRunner) RunMyTool(ctx context.Context, repoPath string) ToolResult {
       start := time.Now()
       result := ToolResult{Tool: "mytool"}
       
       args := []string{"scan", "--json", repoPath}
       output, timedOut, err := r.runTool(ctx, "mytool", args, repoPath)
       
       result.Duration = time.Since(start)
       result.TimedOut = timedOut
       result.Findings = parseMyToolOutput(output)
       
       return result
   }
   ```

3. **Register the tool:**
   ```go
   // In GetToolsForLanguages or RunToolByName
   case "mytool":
       return r.RunMyTool(ctx, repoPath)
   ```

4. **Rebuild the scanner container:**
   ```bash
   ./build.sh
   ```

### Removing Tools

To disable specific tools, modify `GetToolsForLanguages` in `tools.go`:

```go
func (r *ToolRunner) GetToolsForLanguages(languages []Language) []string {
    tools := []string{
        "trivy",
        "semgrep",
        // "trufflehog",  // Commented out to disable
        "gitleaks",
    }
    // ...
}
```

### Adjusting Tool Configurations

**Trivy:**
```go
args := []string{
    "fs",
    "--format", "json",
    "--scanners", "vuln,secret",  // Remove "misconfig" to skip
    "--severity", "CRITICAL,HIGH", // Only high severity
    repoPath,
}
```

**Semgrep:**
```go
args := []string{
    "scan",
    "--config", "p/security-audit",
    // "--config", "p/owasp-top-ten",  // Remove for faster scans
    "--json",
    repoPath,
}
```

---

## Troubleshooting

### Common Issues

**"OPENAI_API_KEY not set"**
```bash
# Ensure .env file exists and contains the key
cat .env | grep OPENAI_API_KEY

# Rebuild and restart to pick up changes
./build.sh
```

**"Database connection refused"**
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check PostgreSQL logs
docker logs better-kiro-prompts-postgres-1

# Verify DATABASE_URL in .env matches docker-compose settings
```

**"Scanner container not found"**
```bash
# Rebuild to ensure scanner is included
./build.sh

# Check scanner container status
docker ps -a | grep scanner
```

**"Rate limit exceeded"**
- Wait for the rate limit window to reset (1 hour)
- Increase limits in `config.toml` if self-hosting
- Check `rate_limit.*` settings

**"AI generation timeout"**
```toml
# Increase timeout in config.toml
[openai]
timeout = "300s"  # 5 minutes
```

**"Repository too large"**
```toml
# Increase limit in config.toml
[scanner]
max_repo_size_mb = 1000  # 1 GB
```

### Viewing Logs

```bash
# All container logs
docker compose -f docker-compose.prod.yml logs

# Specific service
docker logs better-kiro-prompts-backend-1
docker logs better-kiro-prompts-postgres-1

# Follow logs in real-time
docker compose -f docker-compose.prod.yml logs -f

# Application log files
ls -la ./logs/
tail -f ./logs/*-app.log
```

### Debug Mode

Enable debug logging for more detail:

```toml
[logging]
level = "DEBUG"
```

Or via environment variable:
```bash
LOG_LEVEL=DEBUG ./build.sh
```

### Getting Help

1. Check the [API documentation](api.md)
2. Review the [Developer Guide](developer.md)
3. Search existing issues in the repository
4. Open a new issue with:
   - Configuration (redact secrets)
   - Error messages
   - Steps to reproduce
