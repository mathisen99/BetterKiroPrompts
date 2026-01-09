# Phase 3: Polish & Testing — Design

## Context

Phase 3 adds polish, testing, and documentation to the Phase 1-2 foundation. Optional repo scanning is a separate module.

## Testing Strategy

### Unit Tests

Location: `backend/internal/generator/*_test.go`

| Package | Test Focus |
|---------|------------|
| generator/kickoff | Template rendering, answer validation |
| generator/steering | Frontmatter generation, file content |
| generator/hooks | Preset logic, hook schema validity |

Framework: Go standard `testing` package.

### Integration Tests

Location: `backend/internal/api/*_test.go`

| Endpoint | Test Cases |
|----------|------------|
| POST /api/kickoff/generate | Valid input, missing fields, malformed JSON |
| POST /api/steering/generate | All file types, conditional toggle |
| POST /api/hooks/generate | All presets, tech stack variations |

Framework: Go `httptest` package.

### E2E Tests

Location: `frontend/e2e/`

| Flow | Test Cases |
|------|------------|
| Kickoff | Complete wizard, validation blocking, preview, copy, download |
| Steering | Config form, conditional files, multi-file download |
| Hooks | Preset selection, preview, download |

Framework: Playwright.

## Documentation Structure

```
/
├── README.md                    # Setup, quick start, project overview
└── docs/
    ├── api.md                   # Endpoint documentation
    └── user-guide.md            # Generated output explanations
```

### README.md Content
- Project purpose (one paragraph)
- Prerequisites (Docker, Node.js versions)
- Quick start (`./build.sh up`)
- Development commands
- Project structure overview

### docs/api.md Content
- Endpoint list with methods
- Request/response schemas (from Phase 2 design)
- Example curl commands
- Error response format

### docs/user-guide.md Content
- What is a kickoff prompt and how to use it
- Steering file types and inclusion modes
- Hook presets and customization
- Commit message contract explanation

## UI/UX Polish Implementation

### Error Handling

```tsx
// Error boundary at app level
<ErrorBoundary fallback={<ErrorFallback />}>
  <App />
</ErrorBoundary>

// API error handling in components
const { error, retry } = useApiCall();
if (error) return <ErrorMessage message={error} onRetry={retry} />;
```

### Loading States

```tsx
// Skeleton loading for wizard steps
{isLoading ? <Skeleton className="h-32" /> : <QuestionStep />}

// Button loading state
<Button disabled={isSubmitting}>
  {isSubmitting ? <Spinner /> : "Generate"}
</Button>
```

### Toast Notifications

Use shadcn/ui `toast` component:
- Success: "Copied to clipboard", "Downloaded successfully"
- Error: "Failed to generate. Please try again."

### Accessibility

- All form inputs have associated labels
- Focus management in wizard (auto-focus next step)
- ARIA live regions for dynamic content
- Skip links for keyboard navigation
- Color contrast verified with axe-core

## Missing Features Implementation

### Manual Steering

Add to SteeringConfigurator:
```tsx
<Checkbox
  id="manual-steering"
  checked={includeManual}
  onCheckedChange={setIncludeManual}
/>
<Label htmlFor="manual-steering">
  Include manual steering files (referenced via #steering-file-name)
</Label>
```

Backend generates with frontmatter:
```yaml
---
inclusion: manual
---
```

### File References

Add to steering templates:
```markdown
## Related Files

#[[file:.env.example]]
#[[file:backend/migrations/README.md]]
```

UI allows adding custom file references.

### Commit Contract Display

Add `CommitContract` component shown on all output panels:
```tsx
<Card className="bg-muted">
  <CardHeader>Commit Message Contract</CardHeader>
  <CardContent>
    <ul>
      <li>Atomic: one concern per commit</li>
      <li>Prefixed: feat:, fix:, docs:, chore:</li>
      <li>One-sentence summary</li>
    </ul>
  </CardContent>
</Card>
```

## OPTIONAL: Repo Scanning Architecture

### Isolation Model

```
┌─────────────────────────────────────┐
│           Main Backend              │
│  POST /api/scan/start               │
│  GET  /api/scan/:id/status          │
│  GET  /api/scan/:id/results         │
└──────────────┬──────────────────────┘
               │ Job Queue
               ▼
┌─────────────────────────────────────┐
│         Scan Worker Container       │
│  - No outbound network              │
│  - Read-only repo mount             │
│  - Hard timeout (5 min default)     │
│  - Tools: TruffleHog, Gitleaks,     │
│           osv-scanner, govulncheck  │
└─────────────────────────────────────┘
```

### Tool Execution

```go
type ScanResult struct {
    Tool     string    `json:"tool"`
    Findings []Finding `json:"findings"`
    Timeout  bool      `json:"timeout"`
    Error    string    `json:"error,omitempty"`
}

type Finding struct {
    Severity string `json:"severity"` // critical, high, medium, low
    File     string `json:"file"`
    Line     int    `json:"line"`
    Message  string `json:"message"`
    Tool     string `json:"tool"`
}
```

### AI Summary

- Input: raw tool output JSON
- Output: prioritized summary, no invented findings
- Model: GPT-5.1-Codex-Max (per plan)

## Deployment Considerations

- All services containerized (existing docker-compose)
- Production build: `./build.sh --prod --build -d up`
- Environment variables documented in `.env.example`
- Health checks on all services
- Graceful shutdown handling
