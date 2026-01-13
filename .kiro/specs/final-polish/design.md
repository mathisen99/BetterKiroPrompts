# Design Document: Final Polish

## Overview

This design covers the final polish phase for BetterKiroPrompts, transforming the MVP into a production-ready application. The implementation spans frontend UI improvements, syntax highlighting, state persistence, database storage, a public gallery with ratings, and security/concurrency hardening.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ App Shell    │  │ Generation   │  │ Gallery              │  │
│  │ - Logo mgmt  │  │ Flow         │  │ - List view          │  │
│  │ - Transitions│  │ - Questions  │  │ - Detail modal       │  │
│  │ - Header     │  │ - Output     │  │ - Rating component   │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ LocalStorage │  │ Syntax       │  │ Error Recovery       │  │
│  │ Manager      │  │ Highlighter  │  │ - Retry logic        │  │
│  │ - Save/Load  │  │ - JSON/MD    │  │ - State restore      │  │
│  │ - Expiry     │  │ - Themes     │  │ - User feedback      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Backend (Go)                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ API Router   │  │ Generation   │  │ Gallery              │  │
│  │ - Middleware │  │ Service      │  │ Service              │  │
│  │ - Request ID │  │ - OpenAI     │  │ - List/Filter        │  │
│  │ - Logging    │  │ - Queue      │  │ - Ratings            │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ Rate Limiter │  │ Storage      │  │ Sanitizer            │  │
│  │ - Per-IP     │  │ Repository   │  │ - Input validation   │  │
│  │ - Per-action │  │ - Postgres   │  │ - XSS prevention     │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      PostgreSQL Database                         │
├─────────────────────────────────────────────────────────────────┤
│  generations          │  ratings              │  categories      │
│  - id (UUID)          │  - id (UUID)          │  - id (INT)      │
│  - project_idea       │  - generation_id (FK) │  - name          │
│  - experience_level   │  - score (1-5)        │  - keywords[]    │
│  - hook_preset        │  - voter_hash         │                  │
│  - files (JSONB)      │  - created_at         │                  │
│  - category_id (FK)   │                       │                  │
│  - avg_rating         │                       │                  │
│  - rating_count       │                       │                  │
│  - view_count         │                       │                  │
│  - created_at         │                       │                  │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Frontend Components

#### 1. App Shell Updates (`App.tsx`)

```typescript
interface AppState {
  showLargeLogo: boolean
  currentPhase: Phase
  isTransitioning: boolean
}

// Logo visibility based on phase
const shouldShowLargeLogo = (phase: Phase): boolean => {
  return phase === 'level-select'
}

// Compact header for non-landing phases
interface CompactHeaderProps {
  onStartOver: () => void
}
```

#### 2. LocalStorage Manager (`lib/storage.ts`)

```typescript
interface SessionState {
  phase: Phase
  experienceLevel: ExperienceLevel | null
  projectIdea: string
  hookPreset: HookPreset
  questions: Question[]
  answers: Record<number, string>
  currentQuestionIndex: number
  savedAt: number // Unix timestamp
}

interface StorageManager {
  save(state: SessionState): void
  load(): SessionState | null
  clear(): void
  isExpired(state: SessionState): boolean
}

const STORAGE_KEY = 'bkp_session'
const EXPIRY_MS = 24 * 60 * 60 * 1000 // 24 hours
```

#### 3. Syntax Highlighter (`components/SyntaxHighlighter.tsx`)

Using `react-syntax-highlighter` with `prism` for lightweight highlighting:

```typescript
interface SyntaxHighlighterProps {
  code: string
  language: 'json' | 'markdown' | 'yaml'
  editable?: boolean
  onChange?: (code: string) => void
}

// Detect language from file path
const detectLanguage = (path: string): 'json' | 'markdown' | 'yaml' => {
  if (path.endsWith('.json')) return 'json'
  if (path.endsWith('.md')) return 'markdown'
  if (path.endsWith('.yaml') || path.endsWith('.yml')) return 'yaml'
  return 'markdown' // default
}
```

#### 4. Gallery Components (`components/Gallery/`)

```typescript
interface GalleryItem {
  id: string
  projectIdea: string
  category: string
  avgRating: number
  ratingCount: number
  viewCount: number
  createdAt: string
  preview: string // First 200 chars of kickoff prompt
}

interface GalleryFilters {
  category: string | null
  sortBy: 'newest' | 'highest_rated' | 'most_viewed'
  page: number
}

interface GalleryListProps {
  items: GalleryItem[]
  filters: GalleryFilters
  onFilterChange: (filters: GalleryFilters) => void
  onItemClick: (id: string) => void
  totalPages: number
}

interface GalleryDetailProps {
  generation: GenerationDetail
  onClose: () => void
  onRate: (score: number) => void
  userRating: number | null
}

interface RatingProps {
  value: number
  count: number
  userRating: number | null
  onRate: (score: number) => void
  disabled?: boolean
}
```

### Backend Components

#### 1. Database Schema (Migration)

```sql
-- 20260113000001_create_generations.sql

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    keywords TEXT[] NOT NULL DEFAULT '{}'
);

INSERT INTO categories (name, keywords) VALUES
    ('API', ARRAY['api', 'rest', 'graphql', 'endpoint', 'backend']),
    ('CLI', ARRAY['cli', 'command', 'terminal', 'shell', 'script']),
    ('Web App', ARRAY['web', 'frontend', 'react', 'vue', 'angular', 'website']),
    ('Mobile', ARRAY['mobile', 'ios', 'android', 'react native', 'flutter']),
    ('Other', ARRAY[]);

CREATE TABLE generations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_idea TEXT NOT NULL,
    experience_level VARCHAR(20) NOT NULL,
    hook_preset VARCHAR(20) NOT NULL,
    files JSONB NOT NULL,
    category_id INTEGER REFERENCES categories(id) DEFAULT 5,
    avg_rating DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    generation_id UUID NOT NULL REFERENCES generations(id) ON DELETE CASCADE,
    score SMALLINT NOT NULL CHECK (score >= 1 AND score <= 5),
    voter_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(generation_id, voter_hash)
);

CREATE INDEX idx_generations_category ON generations(category_id);
CREATE INDEX idx_generations_created_at ON generations(created_at DESC);
CREATE INDEX idx_generations_avg_rating ON generations(avg_rating DESC);
CREATE INDEX idx_generations_view_count ON generations(view_count DESC);
CREATE INDEX idx_ratings_generation ON ratings(generation_id);
```

#### 2. Storage Repository (`internal/storage/repository.go`)

```go
type Generation struct {
    ID              string          `json:"id"`
    ProjectIdea     string          `json:"projectIdea"`
    ExperienceLevel string          `json:"experienceLevel"`
    HookPreset      string          `json:"hookPreset"`
    Files           json.RawMessage `json:"files"`
    CategoryID      int             `json:"categoryId"`
    CategoryName    string          `json:"categoryName"`
    AvgRating       float64         `json:"avgRating"`
    RatingCount     int             `json:"ratingCount"`
    ViewCount       int             `json:"viewCount"`
    CreatedAt       time.Time       `json:"createdAt"`
}

type Repository interface {
    // Generations
    CreateGeneration(ctx context.Context, gen *Generation) error
    GetGeneration(ctx context.Context, id string) (*Generation, error)
    ListGenerations(ctx context.Context, filter ListFilter) ([]Generation, int, error)
    IncrementViewCount(ctx context.Context, id string) error
    
    // Ratings
    CreateOrUpdateRating(ctx context.Context, genID string, score int, voterHash string) error
    GetUserRating(ctx context.Context, genID string, voterHash string) (int, error)
    
    // Categories
    GetCategoryByKeywords(ctx context.Context, text string) (int, error)
}

type ListFilter struct {
    CategoryID *int
    SortBy     string // "newest", "highest_rated", "most_viewed"
    Page       int
    PageSize   int
}
```

#### 3. Gallery Service (`internal/gallery/service.go`)

```go
type Service struct {
    repo        storage.Repository
    rateLimiter *ratelimit.Limiter
}

func (s *Service) ListGenerations(ctx context.Context, filter ListFilter) (*ListResponse, error)
func (s *Service) GetGeneration(ctx context.Context, id string) (*Generation, error)
func (s *Service) RateGeneration(ctx context.Context, id string, score int, voterHash string) error
```

#### 4. Input Sanitizer (`internal/sanitize/sanitize.go`)

```go
// Sanitize removes potentially dangerous content from user input
func Sanitize(input string) string {
    // Remove HTML tags
    // Escape special characters
    // Trim to max length
}

// ValidateProjectIdea validates and sanitizes project idea
func ValidateProjectIdea(idea string) (string, error) {
    if len(idea) > MaxProjectIdeaLength {
        return "", ErrProjectIdeaTooLong
    }
    return Sanitize(strings.TrimSpace(idea)), nil
}
```

#### 5. Request Queue (`internal/queue/queue.go`)

```go
type RequestQueue struct {
    maxConcurrent int
    semaphore     chan struct{}
}

func NewRequestQueue(maxConcurrent int) *RequestQueue {
    return &RequestQueue{
        maxConcurrent: maxConcurrent,
        semaphore:     make(chan struct{}, maxConcurrent),
    }
}

func (q *RequestQueue) Acquire(ctx context.Context) error {
    select {
    case q.semaphore <- struct{}{}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (q *RequestQueue) Release() {
    <-q.semaphore
}
```

#### 6. Middleware (`internal/api/middleware.go`)

```go
// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler

// LoggingMiddleware logs requests with timing and status
func LoggingMiddleware(next http.Handler) http.Handler

// RecoveryMiddleware recovers from panics and returns 500
func RecoveryMiddleware(next http.Handler) http.Handler
```

## Data Models

### Frontend State

```typescript
// Extended LandingPageState with persistence
interface LandingPageState {
  phase: Phase
  experienceLevel: ExperienceLevel | null
  projectIdea: string
  hookPreset: HookPreset
  questions: Question[]
  answers: Map<number, string>
  currentQuestionIndex: number
  generatedFiles: GeneratedFile[]
  editedFiles: Map<string, string>
  error: string | null
  retryAfter: number | null
  // New fields
  generationId: string | null // Set after successful generation
  canRetry: boolean // True if we have enough state to retry
  loadingStartTime: number | null // For progress tracking
}
```

### API Responses

```typescript
// Extended generation response
interface GenerateOutputsResponse {
  files: GeneratedFile[]
  generationId: string // New: ID of stored generation
}

// Gallery API
interface GalleryListResponse {
  items: GalleryItem[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

interface GalleryDetailResponse {
  generation: GenerationDetail
  userRating: number | null
}

// Structured error response
interface ErrorResponse {
  error: string
  code: string // e.g., "TIMEOUT", "RATE_LIMITED", "VALIDATION_ERROR"
  retryAfter?: number
  requestId: string
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: State Persistence Consistency
*For any* user action that modifies session state (level selection, project idea submission, answer submission), the LocalStorage_Manager SHALL save the updated state, and loading that state SHALL restore the exact same values.
**Validates: Requirements 4.1, 4.3**

### Property 2: Syntax Highlighter Robustness
*For any* string input (including malformed JSON, invalid Markdown, empty strings, and strings with special characters), the Syntax_Highlighter SHALL render without throwing an exception.
**Validates: Requirements 2.6**

### Property 3: Generation Record Completeness
*For any* successfully stored generation, the database record SHALL contain all required fields (id, project_idea, experience_level, hook_preset, files, category_id, created_at) and SHALL NOT contain any user-identifying information (IP address, user agent, etc.).
**Validates: Requirements 5.2, 5.6**

### Property 4: Category Assignment Correctness
*For any* project idea containing category keywords, the Backend SHALL assign the matching category. For project ideas with multiple matching categories, the first match in priority order (API > CLI > Web App > Mobile > Other) SHALL be used.
**Validates: Requirements 5.3**

### Property 5: Gallery Filtering Correctness
*For any* category filter applied to the gallery, all returned items SHALL have the matching category_id, and no items with different category_id SHALL be included.
**Validates: Requirements 6.2**

### Property 6: Gallery Sorting Correctness
*For any* sort option (newest, highest_rated, most_viewed), the returned items SHALL be in descending order by the corresponding field (created_at, avg_rating, view_count).
**Validates: Requirements 6.3**

### Property 7: Pagination Bounds
*For any* gallery page request, the response SHALL contain at most 20 items, and the total_pages calculation SHALL equal ceil(total_items / page_size).
**Validates: Requirements 6.5**

### Property 8: Rating Storage and Calculation
*For any* rating submission, the rating SHALL be stored with the correct generation_id and score, and the generation's avg_rating SHALL equal the arithmetic mean of all ratings for that generation.
**Validates: Requirements 7.2, 7.3**

### Property 9: Duplicate Rating Prevention
*For any* voter_hash that has already rated a generation, subsequent rating submissions SHALL update the existing rating rather than create a duplicate, and the rating count SHALL remain unchanged.
**Validates: Requirements 7.4**

### Property 10: Rate Limit Enforcement
*For any* IP address that exceeds the rating rate limit (20/hour), subsequent rating requests SHALL return HTTP 429 with a Retry-After header.
**Validates: Requirements 7.6**

### Property 11: Concurrent Request Handling
*For any* set of concurrent generation requests, all requests SHALL complete (success or controlled failure) without deadlock, and no request SHALL block indefinitely.
**Validates: Requirements 8.1**

### Property 12: Request Queue Fairness
*For any* request that acquires a queue slot, it SHALL eventually release the slot, and the queue SHALL process requests in approximate FIFO order.
**Validates: Requirements 8.3**

### Property 13: Input Sanitization
*For any* user input containing HTML tags or JavaScript, the sanitized output SHALL not contain executable code, and the sanitized content SHALL be safe for database storage and HTML rendering.
**Validates: Requirements 9.2, 9.3**

### Property 14: Structured Error Responses
*For any* API error, the response SHALL include error message, error code, and request_id fields. Client errors (4xx) SHALL have codes starting with "CLIENT_", and server errors (5xx) SHALL have codes starting with "SERVER_".
**Validates: Requirements 10.5, 10.6**

### Property 15: Automatic Retry Behavior
*For any* recoverable error (network timeout, 503), the system SHALL retry exactly once before surfacing the error to the user.
**Validates: Requirements 10.4**

## Error Handling

### Frontend Error States

| Error Type | Code | User Message | Actions |
|------------|------|--------------|---------|
| Timeout | `TIMEOUT` | "Request timed out. Your progress is saved." | Retry, Start Over |
| Rate Limited | `RATE_LIMITED` | "Too many requests. Please wait X seconds." | Countdown, Start Over |
| Server Error | `SERVER_ERROR` | "Something went wrong. Please try again." | Retry, Start Over |
| Validation | `VALIDATION_ERROR` | Specific message | Edit Input |
| Network | `NETWORK_ERROR` | "Unable to connect. Check your connection." | Retry |

### Backend Error Codes

```go
const (
    ErrCodeTimeout        = "SERVER_TIMEOUT"
    ErrCodeRateLimited    = "CLIENT_RATE_LIMITED"
    ErrCodeValidation     = "CLIENT_VALIDATION"
    ErrCodeNotFound       = "CLIENT_NOT_FOUND"
    ErrCodeServerError    = "SERVER_INTERNAL"
    ErrCodeServiceUnavail = "SERVER_UNAVAILABLE"
)
```

## Testing Strategy

### Unit Tests

1. **LocalStorage Manager**: Test save/load/clear/expiry logic
2. **Syntax Highlighter**: Test language detection, theme application
3. **Sanitizer**: Test XSS prevention, length limits
4. **Category Matcher**: Test keyword matching logic
5. **Rating Calculator**: Test average calculation

### Property-Based Tests

Using `fast-check` for TypeScript and `testing/quick` for Go:

1. **State Persistence Round-Trip** (Property 1)
2. **Highlighter Robustness** (Property 2)
3. **Category Assignment** (Property 4)
4. **Gallery Filtering** (Property 5)
5. **Gallery Sorting** (Property 6)
6. **Rating Calculation** (Property 8)
7. **Input Sanitization** (Property 13)

### Integration Tests

1. **Full Generation Flow**: End-to-end with database storage
2. **Gallery CRUD**: List, filter, sort, paginate
3. **Rating Flow**: Submit, update, prevent duplicates
4. **Concurrent Requests**: Multiple simultaneous generations
5. **Error Recovery**: Timeout handling, retry logic

### Configuration

- Property tests: minimum 100 iterations
- Integration tests: use test database with transactions
- Frontend tests: use MSW for API mocking
