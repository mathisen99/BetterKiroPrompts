# Design Document: UX Improvements

## Overview

This design addresses six key areas of improvement: experience-level-appropriate questions, clickable answer examples, loading feedback, increased timeouts, IP-based abuse prevention, and navigation visibility. The changes span backend prompts, API responses, frontend components, and database schema.

## Architecture

The improvements integrate into the existing architecture:

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │ QuestionFlow│  │ GalleryPage │  │ Navigation Components   │ │
│  │ + Examples  │  │ + IP Track  │  │ + Visible Links/Buttons │ │
│  │ + Loading   │  │ + Modal Fix │  │                         │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Backend API                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │ /generate/*     │  │ /gallery/*      │  │ Timeout Config  │ │
│  │ + Examples in   │  │ + IP Tracking   │  │ 180s everywhere │ │
│  │   response      │  │ + View Dedup    │  │                 │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Database                                 │
│  ┌─────────────────┐  ┌─────────────────┐                       │
│  │ generations     │  │ views (new)     │                       │
│  │                 │  │ - generation_id │                       │
│  │                 │  │ - ip_hash       │                       │
│  │                 │  │ - created_at    │                       │
│  └─────────────────┘  └─────────────────┘                       │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. Enhanced Question Response

The question generation response will include example answers:

```go
// Question with examples
type Question struct {
    ID       int      `json:"id"`
    Text     string   `json:"text"`
    Hint     string   `json:"hint,omitempty"`
    Examples []string `json:"examples"` // NEW: 3 clickable examples
}
```

### 2. Experience-Level Prompt Differentiation

Update `backend/internal/prompts/questions.go` to have drastically different prompts:

**Beginner Prompt Guidelines:**
- Use everyday language only
- Forbidden terms: API, database, schema, authentication, OAuth, microservices, CI/CD, containerization, deployment, backend, frontend, REST, GraphQL, SQL, NoSQL
- Ask about: what the app does, who uses it, what information it stores, where to save data (simple terms)
- Examples should use analogies to real-world concepts

**Novice Prompt Guidelines:**
- Can use basic technical terms with brief explanations
- Allowed: database (with explanation), API (with explanation), authentication
- Ask about: data storage approach, user authentication needs, basic architecture

**Expert Prompt Guidelines:**
- Full technical terminology
- Ask about: architecture patterns, scaling considerations, security model, data consistency, deployment strategy

### 3. View Tracking Table

New database migration for IP-based view deduplication:

```sql
CREATE TABLE views (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    generation_id UUID NOT NULL REFERENCES generations(id) ON DELETE CASCADE,
    ip_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(generation_id, ip_hash)
);

CREATE INDEX idx_views_generation_id ON views(generation_id);
```

### 4. Updated Rating Table

Ensure ratings table uses IP hash for deduplication (already has voter_hash, verify it's used correctly).

### 5. Timeout Configuration

```go
// backend/internal/openai/client.go
const defaultTimeout = 180 * time.Second // Changed from 120s

// frontend/src/lib/api.ts
const DEFAULT_TIMEOUT_MS = 180 * 1000 // Changed from 120s
```

### 6. Frontend Loading States

```typescript
// Loading messages
const LOADING_MESSAGES = {
  questions: "Generating questions... This may take up to 2 minutes",
  questionsLong: "Still working on your questions...",
  outputs: "Generating your files... This may take up to 3 minutes",
  outputsLong: "Still creating your files... Almost there!"
}
```

## Data Models

### Question (Updated)

```typescript
interface Question {
  id: number
  text: string
  hint?: string
  examples: string[] // NEW: exactly 3 examples
}
```

### View (New)

```go
type View struct {
    ID           string    `json:"id"`
    GenerationID string    `json:"generationId"`
    IPHash       string    `json:"ipHash"`
    CreatedAt    time.Time `json:"createdAt"`
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Beginner Questions Avoid Technical Jargon

*For any* project idea submitted with beginner experience level, the generated questions SHALL NOT contain any of the forbidden technical terms: "API", "database schema", "authentication flow", "microservices", "CI/CD", "containerization", "OAuth", "REST", "GraphQL", "SQL", "NoSQL", "backend", "frontend", "deployment".

**Validates: Requirements 1.1, 1.2**

### Property 2: Experience Levels Produce Different Questions

*For any* project idea, generating questions at beginner, novice, and expert levels SHALL produce three distinct question sets where no two sets are identical.

**Validates: Requirements 1.5**

### Property 3: Questions Include Exactly Three Examples

*For any* generated question response, each question SHALL contain exactly 3 example answers in the examples array.

**Validates: Requirements 2.1**

### Property 4: Idempotent View Counting

*For any* gallery item and IP address, multiple view requests SHALL result in exactly one view record in the database and the view_count SHALL increment by at most 1.

**Validates: Requirements 5.1, 5.3**

### Property 5: Vote Upsert Behavior

*For any* gallery item and IP address, submitting multiple votes SHALL result in exactly one rating record, with the score reflecting the most recent vote.

**Validates: Requirements 5.2, 5.4**

### Property 6: IP Addresses Are Hashed

*For any* view or rating record stored in the database, the IP identifier SHALL be a SHA-256 hash, not a raw IP address.

**Validates: Requirements 5.5**

## Error Handling

### Timeout Errors

- Display user-friendly message: "The request took too long. Please try again."
- Offer a "Retry" button
- Log timeout events for monitoring

### Generation Failures

- If question generation fails, show error with retry option
- If output generation fails after questions answered, preserve answers and allow retry

### View/Vote Failures

- Silently handle duplicate view attempts (no error to user)
- For vote updates, show success message regardless of create vs update

## Testing Strategy

### Unit Tests

- Test prompt generation for each experience level
- Test IP hashing function
- Test view deduplication logic
- Test timeout configuration values

### Property-Based Tests

Using Go's `testing/quick` package:

1. **Beginner Jargon Avoidance**: Generate questions for random project ideas at beginner level, verify no forbidden terms
2. **Experience Level Differentiation**: Generate questions for same project at all levels, verify differences
3. **Example Count**: Verify all questions have exactly 3 examples
4. **View Idempotence**: Simulate multiple views from same IP, verify count
5. **Vote Upsert**: Simulate multiple votes from same IP, verify single record with latest score
6. **IP Hashing**: Verify stored values match SHA-256 pattern

### Integration Tests

- End-to-end question generation flow
- Gallery view tracking
- Rating submission and update

### Frontend Tests

- Loading state visibility
- Example answer click behavior
- Modal close interactions (click outside, Escape key, X button)
- Navigation link visibility
