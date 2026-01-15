# API Documentation

Base URL: `http://localhost:8090/api`

## Overview

BetterKiroPrompts provides a REST API for AI-driven generation, gallery browsing, and security scanning.

## Authentication

No authentication is required. Rate limiting is applied per IP address.

## Common Response Formats

### Error Response
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 202 | Accepted (async operation started) |
| 400 | Bad request (invalid input) |
| 404 | Resource not found |
| 429 | Rate limited |
| 500 | Internal server error |
| 504 | Gateway timeout |

---

## Health Check

### GET /health

Check if the API is running.

**Response:**
```json
{"status": "ok"}
```

**Example:**
```bash
curl http://localhost:8090/api/health
```

---

## Generation Endpoints

### POST /generate/questions

Generate contextual questions based on a project idea.

**Request:**
```json
{
  "projectIdea": "A todo app with categories and due dates",
  "experienceLevel": "novice"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| projectIdea | string | Yes | Project description (max 2000 chars) |
| experienceLevel | string | Yes | beginner, novice, or expert |

**Response:**
```json
{
  "questions": [
    {
      "id": 1,
      "text": "What authentication method will you use?",
      "hint": "Consider OAuth, JWT, or session-based auth",
      "examples": ["JWT with refresh tokens", "OAuth 2.0 with Google", "No authentication needed"]
    }
  ]
}
```

**Errors:**
- 400 - Invalid project idea or experience level
- 429 - Rate limited (check Retry-After header)

---

### POST /generate/outputs

Generate kickoff prompt, steering files, hooks, and AGENTS.md.

**Request:**
```json
{
  "projectIdea": "A todo app with categories and due dates",
  "answers": [
    {"questionId": 1, "answer": "JWT authentication"},
    {"questionId": 2, "answer": "PostgreSQL database"}
  ],
  "experienceLevel": "novice",
  "hookPreset": "default"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| projectIdea | string | Yes | Project description |
| answers | array | Yes | Answers to generated questions |
| experienceLevel | string | Yes | beginner, novice, or expert |
| hookPreset | string | Yes | light, basic, default, or strict |

**Response:**
```json
{
  "files": [
    {"path": "kickoff-prompt.md", "content": "...", "type": "kickoff"},
    {"path": ".kiro/steering/product.md", "content": "...", "type": "steering"},
    {"path": ".kiro/hooks/format-on-stop.kiro.hook", "content": "...", "type": "hook"},
    {"path": "AGENTS.md", "content": "...", "type": "agents"}
  ],
  "generationId": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Errors:**
- 400 - Invalid input
- 429 - Rate limited
- 504 - Generation timeout


---

## Gallery Endpoints

### GET /gallery

List gallery items with pagination and filtering.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| pageSize | int | 20 | Items per page (max 100) |
| sort | string | newest | newest, highest_rated, or most_viewed |
| category | int | - | Filter by category ID |

**Response:**
```json
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "projectIdea": "A todo app with categories",
      "category": "Web App",
      "avgRating": 4.5,
      "ratingCount": 12,
      "viewCount": 156,
      "createdAt": "2026-01-14T10:30:00Z",
      "preview": "A todo app with categories..."
    }
  ],
  "total": 100,
  "page": 1,
  "pageSize": 20,
  "totalPages": 5
}
```

**Example:**
```bash
curl "http://localhost:8090/api/gallery?sort=highest_rated&page=1"
```

---

### GET /gallery/{id}

Get full details of a gallery item. Increments view count (deduplicated by IP).

**Response:**
```json
{
  "generation": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "projectIdea": "A todo app with categories",
    "experienceLevel": "novice",
    "hookPreset": "default",
    "files": [...],
    "category": "Web App",
    "avgRating": 4.5,
    "ratingCount": 12,
    "viewCount": 157,
    "createdAt": "2026-01-14T10:30:00Z"
  },
  "userRating": 5
}
```

**Errors:**
- 404 - Generation not found

---

### POST /gallery/{id}/rate

Rate a gallery item (1-5 stars). One rating per IP per generation.

**Request:**
```json
{
  "score": 5,
  "voterHash": "abc123..."
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| score | int | Yes | Rating 1-5 |
| voterHash | string | Yes | Client-generated voter identifier |

**Response:**
```json
{"success": true}
```

**Errors:**
- 400 - Invalid score (must be 1-5)
- 404 - Generation not found
- 429 - Rate limited

---

## Security Scan Endpoints

### POST /scan

Start a new security scan for a GitHub repository.

**Request:**
```json
{
  "repo_url": "https://github.com/owner/repo"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| repo_url | string | Yes | GitHub repository URL |

**Response (202 Accepted):**
```json
{
  "id": "scan-123",
  "status": "pending",
  "repo_url": "https://github.com/owner/repo",
  "created_at": "2026-01-14T10:30:00Z"
}
```

**Errors:**
- 400 - Invalid repository URL
- 429 - Rate limited

---

### GET /scan/{id}

Get scan status and results.

**Response:**
```json
{
  "id": "scan-123",
  "status": "completed",
  "repo_url": "https://github.com/owner/repo",
  "languages": ["go", "javascript"],
  "findings": [
    {
      "id": "finding-1",
      "severity": "high",
      "tool": "semgrep",
      "file_path": "src/auth.go",
      "line_number": 42,
      "description": "Hardcoded credentials detected",
      "remediation": "Use environment variables for secrets",
      "code_example": "password := os.Getenv(\"DB_PASSWORD\")"
    }
  ],
  "review_stats": {
    "files_reviewed": 5,
    "matched_findings": 3
  },
  "created_at": "2026-01-14T10:30:00Z",
  "completed_at": "2026-01-14T10:32:00Z"
}
```

**Scan Status Values:**
- pending - Scan queued
- cloning - Cloning repository
- scanning - Running security tools
- reviewing - AI code review in progress
- completed - Scan finished
- failed - Scan failed (check error field)

**Errors:**
- 404 - Scan job not found

---

### GET /scan/config

Get scanner configuration.

**Response:**
```json
{
  "private_repo_enabled": true,
  "ai_review_enabled": true,
  "max_files_to_review": 10
}
```


---

## Admin Endpoints

### GET /admin/log-level

Get current log level.

**Response:**
```json
{"level": "INFO"}
```

---

### POST /admin/log-level

Change log level at runtime.

**Request:**
```json
{"level": "DEBUG"}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| level | string | Yes | DEBUG, INFO, WARN, or ERROR |

**Response:**
```json
{"level": "DEBUG"}
```

---

## Client Logging

### POST /logs/client

Submit frontend error logs to the backend.

**Request:**
```json
{
  "logs": [
    {
      "level": "error",
      "message": "Failed to load gallery",
      "stack": "Error: Network error...",
      "url": "/gallery",
      "component": "GalleryPage",
      "user_agent": "Mozilla/5.0...",
      "timestamp": "2026-01-14T10:30:00Z"
    }
  ]
}
```

**Response:** 202 Accepted

---

## Rate Limiting

Rate limits are applied per IP address:

| Endpoint | Limit |
|----------|-------|
| Generation | 10/hour |
| Rating | 20/hour |
| Scanning | 10/hour |

When rate limited, the response includes:
- HTTP status 429 Too Many Requests
- Retry-After header with seconds until reset

**Example Response:**
```json
{
  "error": "Rate limit exceeded",
  "code": "RATE_LIMITED",
  "retryAfter": 3600
}
```
