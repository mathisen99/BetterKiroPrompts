# Design Document: Comprehensive Logging

## Overview

This design implements a structured, file-based logging system across the entire application stack. The system uses Go's `log/slog` package (standard library) for structured JSON logging with colored console output, request ID correlation, and automatic log rotation. Frontend errors are captured and sent to a dedicated backend endpoint.

The logging covers EVERY component in the system:
- HTTP middleware (requests, responses, rate limiting)
- Generation service (question generation, output generation, AI calls)
- Gallery service (list, view, rate operations)
- Scanner service (clone, detect, scan, review, aggregate)
- OpenAI client (requests, responses, errors, retries)
- Database operations (queries, connections, migrations)
- Queue operations (acquire, release, stats)
- Rate limiter (allow, deny, remaining)
- Frontend errors (JS errors, API failures, React errors)

## Architecture

```mermaid
graph TB
    subgraph Frontend
        FE[React App] --> EB[ErrorBoundary]
        FE --> AC[API Client]
        FE --> QF[QuestionFlow]
        FE --> OE[OutputEditor]
        FE --> SP[ScanProgress]
        EB --> LC[LogCollector]
        AC --> LC
        LC -->|Batch POST| BE
    end
    
    subgraph "Backend - API Layer"
        BE[/api/logs/client] --> L[Logger]
        MW[Middleware] --> L
        GH[GenerateHandler] --> L
        GLH[GalleryHandler] --> L
        SH[ScanHandler] --> L
        HH[HealthHandler] --> L
    end
    
    subgraph "Backend - Services"
        GS[GenerationService] --> L
        GLS[GalleryService] --> L
        SS[ScannerService] --> L
        OAI[OpenAI Client] --> L
        Q[RequestQueue] --> L
        RL[RateLimiter] --> L
    end
    
    subgraph "Backend - Scanner Pipeline"
        CL[Cloner] --> L
        LD[LanguageDetector] --> L
        TR[ToolRunner] --> L
        AG[Aggregator] --> L
        CR[CodeReviewer] --> L
    end
    
    subgraph "Backend - Data Layer"
        DB[Database] --> L
        REPO[Repository] --> L
        MIG[Migrations] --> L
    end
    
    subgraph "Log Files (Host Mount ./logs)"
        L --> APP[YYYY-MM-DD-app.log]
        L --> HTTP[YYYY-MM-DD-http.log]
        L --> DBL[YYYY-MM-DD-db.log]
        L --> SCANL[YYYY-MM-DD-scanner.log]
        L --> CLIENT[YYYY-MM-DD-client.log]
    end
    
    subgraph Console
        L --> STDOUT[Colored stdout/stderr]
    end
```

## Components and Interfaces

### 1. Logger Package (`backend/internal/logger`)

The core logging infrastructure using Go's standard `log/slog` package.

```go
package logger

import (
    "context"
    "io"
    "log/slog"
    "os"
    "sync"
    "time"
)

// Level represents log severity
type Level = slog.Level

const (
    LevelDebug = slog.LevelDebug
    LevelInfo  = slog.LevelInfo
    LevelWarn  = slog.LevelWarn
    LevelError = slog.LevelError
)

// Config holds logger configuration
type Config struct {
    Level       Level
    LogDir      string
    MaxSizeMB   int
    MaxAgeDays  int
    EnableColor bool
}

// Logger wraps slog with file rotation and colored output
type Logger struct {
    config     Config
    handlers   map[string]*slog.Logger
    files      map[string]*RotatingFile
    mu         sync.RWMutex
    levelVar   *slog.LevelVar
}

// RotatingFile handles log rotation
type RotatingFile struct {
    path       string
    maxSize    int64
    file       *os.File
    size       int64
    mu         sync.Mutex
}

// New creates a new logger instance
func New(cfg Config) (*Logger, error)

// WithContext returns a logger with request context
func (l *Logger) WithContext(ctx context.Context) *slog.Logger

// App returns the application logger
func (l *Logger) App() *slog.Logger

// HTTP returns the HTTP request logger
func (l *Logger) HTTP() *slog.Logger

// DB returns the database logger
func (l *Logger) DB() *slog.Logger

// Scanner returns the scanner logger
func (l *Logger) Scanner() *slog.Logger

// Client returns the client error logger
func (l *Logger) Client() *slog.Logger

// SetLevel changes the log level at runtime
func (l *Logger) SetLevel(level Level)

// Close closes all log files
func (l *Logger) Close() error
```

### 2. Colored Console Handler

Custom slog handler for colored terminal output.

```go
// ColorHandler wraps slog.Handler with ANSI colors
type ColorHandler struct {
    handler slog.Handler
    writer  io.Writer
    noColor bool
}

// ANSI color codes
const (
    colorReset   = "\033[0m"
    colorRed     = "\033[31m"
    colorGreen   = "\033[32m"
    colorYellow  = "\033[33m"
    colorBlue    = "\033[34m"
    colorMagenta = "\033[35m"
    colorCyan    = "\033[36m"
    colorGray    = "\033[90m"
)

func NewColorHandler(w io.Writer, opts *slog.HandlerOptions) *ColorHandler
```

### 3. Context Keys and Helpers

```go
// Context keys for logging
type ctxKey string

const (
    RequestIDKey ctxKey = "request_id"
    ComponentKey ctxKey = "component"
    UserIPKey    ctxKey = "user_ip"
)

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, id string) context.Context

// WithComponent adds component name to context
func WithComponent(ctx context.Context, name string) context.Context

// GetRequestID retrieves request ID from context
func GetRequestID(ctx context.Context) string
```

### 4. HTTP Logging Middleware (Updated)

```go
// LoggingMiddleware logs HTTP requests with structured output
func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            requestID := GetRequestID(r.Context())
            
            // Log request start
            log.HTTP().Info("request_start",
                slog.String("request_id", requestID),
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.String("remote_addr", r.RemoteAddr),
                slog.String("user_agent", r.UserAgent()),
            )
            
            // Wrap response writer
            rw := newResponseWriter(w)
            next.ServeHTTP(rw, r)
            
            // Log request completion
            duration := time.Since(start)
            log.HTTP().Info("request_complete",
                slog.String("request_id", requestID),
                slog.Int("status", rw.statusCode),
                slog.Duration("duration", duration),
                slog.Int64("bytes", rw.bytesWritten),
            )
        })
    }
}
```

### 5. Database Logging Wrapper

```go
// LoggingDB wraps sql.DB with query logging
type LoggingDB struct {
    db  *sql.DB
    log *slog.Logger
}

func NewLoggingDB(db *sql.DB, log *slog.Logger) *LoggingDB

func (l *LoggingDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
    start := time.Now()
    requestID := logger.GetRequestID(ctx)
    
    rows, err := l.db.QueryContext(ctx, query, args...)
    
    l.log.Info("query",
        slog.String("request_id", requestID),
        slog.String("type", "SELECT"),
        slog.Duration("duration", time.Since(start)),
        slog.Bool("success", err == nil),
    )
    
    if err != nil {
        l.log.Error("query_error",
            slog.String("request_id", requestID),
            slog.String("error", err.Error()),
        )
    }
    
    return rows, err
}
```

### 6. Client Error Logging Endpoint

```go
// ClientLogRequest represents a batch of client logs
type ClientLogRequest struct {
    Logs []ClientLogEntry `json:"logs"`
}

type ClientLogEntry struct {
    Level     string    `json:"level"`
    Message   string    `json:"message"`
    Stack     string    `json:"stack,omitempty"`
    URL       string    `json:"url"`
    UserAgent string    `json:"user_agent"`
    Timestamp time.Time `json:"timestamp"`
}

// HandleClientLogs receives and logs frontend errors
func HandleClientLogs(log *logger.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req ClientLogRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid request", http.StatusBadRequest)
            return
        }
        
        for _, entry := range req.Logs {
            log.Client().Log(r.Context(), parseLevel(entry.Level),
                entry.Message,
                slog.String("url", entry.URL),
                slog.String("stack", entry.Stack),
                slog.String("user_agent", entry.UserAgent),
                slog.Time("client_time", entry.Timestamp),
            )
        }
        
        w.WriteHeader(http.StatusAccepted)
    }
}
```

### 7. Frontend Log Collector

```typescript
// lib/logger.ts
interface LogEntry {
  level: 'debug' | 'info' | 'warn' | 'error';
  message: string;
  stack?: string;
  url: string;
  component?: string;
  userAgent: string;
  timestamp: string;
}

class LogCollector {
  private buffer: LogEntry[] = [];
  private flushInterval: number = 5000; // 5 seconds
  private maxBufferSize: number = 50;
  
  constructor() {
    this.setupFlushInterval();
    this.setupErrorHandlers();
  }
  
  private setupErrorHandlers(): void {
    window.onerror = (message, source, line, col, error) => {
      this.error(`${message} at ${source}:${line}:${col}`, error?.stack, 'window');
    };
    
    window.onunhandledrejection = (event) => {
      this.error(`Unhandled rejection: ${event.reason}`, undefined, 'promise');
    };
  }
  
  log(level: LogEntry['level'], message: string, component?: string, stack?: string): void {
    this.buffer.push({
      level,
      message,
      component,
      stack,
      url: window.location.href,
      userAgent: navigator.userAgent,
      timestamp: new Date().toISOString(),
    });
    
    // Console output with colors
    const colors = { debug: '#36D7B7', info: '#3498DB', warn: '#F39C12', error: '#E74C3C' };
    console.log(`%c[${level.toUpperCase()}]%c [${component || 'app'}] ${message}`, 
      `color: ${colors[level]}; font-weight: bold`, 'color: inherit');
    
    if (this.buffer.length >= this.maxBufferSize) {
      this.flush();
    }
  }
  
  debug(message: string, component?: string): void { this.log('debug', message, component); }
  info(message: string, component?: string): void { this.log('info', message, component); }
  warn(message: string, component?: string): void { this.log('warn', message, component); }
  error(message: string, stack?: string, component?: string): void { 
    this.log('error', message, component, stack); 
  }
  
  // Log API calls with timing
  logApiCall(method: string, url: string, status: number, duration: number): void {
    const level = status >= 400 ? 'error' : 'info';
    this.log(level, `${method} ${url} → ${status} (${duration}ms)`, 'api');
  }
  
  async flush(): Promise<void> {
    if (this.buffer.length === 0) return;
    
    const logs = [...this.buffer];
    this.buffer = [];
    
    try {
      await fetch('/api/logs/client', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ logs }),
      });
    } catch {
      // Re-add failed logs to buffer (up to max)
      this.buffer = [...logs.slice(-this.maxBufferSize / 2), ...this.buffer];
    }
  }
  
  private setupFlushInterval(): void {
    setInterval(() => this.flush(), this.flushInterval);
    window.addEventListener('beforeunload', () => this.flush());
  }
}

export const logger = new LogCollector();
```

### 8. Generation Service Logging

```go
// GenerateQuestions with comprehensive logging
func (s *Service) GenerateQuestions(ctx context.Context, projectIdea string, experienceLevel string) ([]Question, error) {
    requestID := logger.GetRequestID(ctx)
    start := time.Now()
    
    s.log.Info("generate_questions_start",
        slog.String("request_id", requestID),
        slog.String("experience_level", experienceLevel),
        slog.Int("idea_length", len(projectIdea)),
    )
    
    // Validation
    if err := ValidateProjectIdea(projectIdea); err != nil {
        s.log.Warn("generate_questions_validation_failed",
            slog.String("request_id", requestID),
            slog.String("error", err.Error()),
        )
        return nil, err
    }
    
    // Queue acquisition
    if s.requestQueue != nil {
        s.log.Debug("queue_acquire_start", slog.String("request_id", requestID))
        if err := s.requestQueue.Acquire(ctx); err != nil {
            s.log.Error("queue_acquire_failed",
                slog.String("request_id", requestID),
                slog.String("error", err.Error()),
            )
            return nil, err
        }
        defer s.requestQueue.Release()
        s.log.Debug("queue_acquire_success", slog.String("request_id", requestID))
    }
    
    // OpenAI call
    s.log.Debug("openai_call_start",
        slog.String("request_id", requestID),
        slog.String("operation", "generate_questions"),
    )
    
    response, err := s.openaiClient.ChatCompletion(ctx, messages)
    if err != nil {
        s.log.Error("generate_questions_failed",
            slog.String("request_id", requestID),
            slog.String("error", err.Error()),
            slog.Duration("duration", time.Since(start)),
        )
        return nil, err
    }
    
    questions, err := parseQuestionsResponse(response)
    if err != nil {
        s.log.Error("generate_questions_parse_failed",
            slog.String("request_id", requestID),
            slog.String("error", err.Error()),
        )
        return nil, err
    }
    
    s.log.Info("generate_questions_complete",
        slog.String("request_id", requestID),
        slog.Int("question_count", len(questions)),
        slog.Duration("duration", time.Since(start)),
    )
    
    return questions, nil
}

// GenerateOutputs with comprehensive logging
func (s *Service) GenerateOutputs(ctx context.Context, projectIdea string, answers []Answer, 
    experienceLevel string, hookPreset string) ([]GeneratedFile, error) {
    requestID := logger.GetRequestID(ctx)
    start := time.Now()
    
    s.log.Info("generate_outputs_start",
        slog.String("request_id", requestID),
        slog.String("experience_level", experienceLevel),
        slog.String("hook_preset", hookPreset),
        slog.Int("answer_count", len(answers)),
    )
    
    // ... validation and queue ...
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        s.log.Debug("generate_outputs_attempt",
            slog.String("request_id", requestID),
            slog.Int("attempt", attempt+1),
            slog.Int("max_attempts", maxRetries+1),
        )
        
        response, err := s.openaiClient.ChatCompletion(ctx, messages)
        if err != nil {
            s.log.Error("generate_outputs_openai_failed",
                slog.String("request_id", requestID),
                slog.Int("attempt", attempt+1),
                slog.String("error", err.Error()),
            )
            return nil, err
        }
        
        files, err := parseOutputsResponse(response)
        if err != nil {
            s.log.Warn("generate_outputs_parse_failed",
                slog.String("request_id", requestID),
                slog.Int("attempt", attempt+1),
                slog.String("error", err.Error()),
            )
            if attempt < maxRetries {
                continue // Retry
            }
            return nil, err
        }
        
        s.log.Info("generate_outputs_complete",
            slog.String("request_id", requestID),
            slog.Int("file_count", len(files)),
            slog.Int("attempts_used", attempt+1),
            slog.Duration("duration", time.Since(start)),
        )
        
        return files, nil
    }
    
    return nil, lastErr
}
```

### 9. Gallery Service Logging

```go
// ListGenerations with logging
func (s *Service) ListGenerations(ctx context.Context, req ListRequest) (*ListResponse, error) {
    requestID := logger.GetRequestID(ctx)
    start := time.Now()
    
    s.log.Info("gallery_list_start",
        slog.String("request_id", requestID),
        slog.String("sort_by", req.SortBy),
        slog.Int("page", req.Page),
        slog.Int("page_size", req.PageSize),
        slog.Any("category_id", req.CategoryID),
    )
    
    items, total, err := s.repo.ListGenerations(ctx, filter)
    if err != nil {
        s.log.Error("gallery_list_failed",
            slog.String("request_id", requestID),
            slog.String("error", err.Error()),
        )
        return nil, err
    }
    
    s.log.Info("gallery_list_complete",
        slog.String("request_id", requestID),
        slog.Int("item_count", len(items)),
        slog.Int("total", total),
        slog.Duration("duration", time.Since(start)),
    )
    
    return response, nil
}

// RateGeneration with logging
func (s *Service) RateGeneration(ctx context.Context, genID string, score int, 
    voterHash string, clientIP string) (retryAfter int, err error) {
    requestID := logger.GetRequestID(ctx)
    
    s.log.Info("gallery_rate_start",
        slog.String("request_id", requestID),
        slog.String("generation_id", genID),
        slog.Int("score", score),
    )
    
    // Rate limit check
    if s.rateLimiter != nil {
        allowed, duration := s.rateLimiter.Allow(clientIP)
        if !allowed {
            s.log.Warn("gallery_rate_limited",
                slog.String("request_id", requestID),
                slog.Duration("retry_after", duration),
            )
            return int(duration.Seconds()), ErrRateLimited
        }
    }
    
    err = s.repo.CreateOrUpdateRating(ctx, genID, score, voterHash)
    if err != nil {
        s.log.Error("gallery_rate_failed",
            slog.String("request_id", requestID),
            slog.String("error", err.Error()),
        )
        return 0, err
    }
    
    s.log.Info("gallery_rate_complete",
        slog.String("request_id", requestID),
        slog.String("generation_id", genID),
    )
    
    return 0, nil
}
```

### 10. Scanner Pipeline Logging (Full Flow)

```go
// runScan with comprehensive phase logging
func (s *Service) runScan(ctx context.Context, jobID string) {
    start := time.Now()
    
    s.log.Info("scan_pipeline_start",
        slog.String("job_id", jobID),
    )
    
    job, err := s.loadJob(ctx, jobID)
    if err != nil {
        s.log.Error("scan_load_job_failed",
            slog.String("job_id", jobID),
            slog.String("error", err.Error()),
        )
        return
    }
    
    // Phase 1: Clone
    s.log.Info("scan_phase_clone_start",
        slog.String("job_id", jobID),
        slog.String("repo_url", job.RepoURL),
    )
    cloneStart := time.Now()
    cloneResult, err := s.cloner.Clone(ctx, job.RepoURL)
    if err != nil {
        s.log.Error("scan_phase_clone_failed",
            slog.String("job_id", jobID),
            slog.String("error", err.Error()),
            slog.Duration("duration", time.Since(cloneStart)),
        )
        _ = s.failJob(ctx, jobID, fmt.Sprintf("Clone failed: %v", err))
        return
    }
    s.log.Info("scan_phase_clone_complete",
        slog.String("job_id", jobID),
        slog.String("path", cloneResult.Path),
        slog.Duration("duration", time.Since(cloneStart)),
    )
    repoPath := cloneResult.Path
    
    // Phase 2: Language Detection
    s.log.Info("scan_phase_detect_start",
        slog.String("job_id", jobID),
    )
    detectStart := time.Now()
    languages, err := s.detector.DetectLanguages(repoPath)
    if err != nil {
        s.log.Error("scan_phase_detect_failed",
            slog.String("job_id", jobID),
            slog.String("error", err.Error()),
        )
        _ = s.failJob(ctx, jobID, fmt.Sprintf("Language detection failed: %v", err))
        return
    }
    langStrings := make([]string, len(languages))
    for i, l := range languages {
        langStrings[i] = string(l)
    }
    s.log.Info("scan_phase_detect_complete",
        slog.String("job_id", jobID),
        slog.Any("languages", langStrings),
        slog.Int("language_count", len(languages)),
        slog.Duration("duration", time.Since(detectStart)),
    )
    
    // Phase 3: Tool Execution
    toolNames := s.toolRunner.GetToolsForLanguages(languages)
    s.log.Info("scan_phase_tools_start",
        slog.String("job_id", jobID),
        slog.Any("tools", toolNames),
        slog.Int("tool_count", len(toolNames)),
    )
    
    var results []ToolResult
    for _, toolName := range toolNames {
        toolStart := time.Now()
        s.log.Debug("scan_tool_start",
            slog.String("job_id", jobID),
            slog.String("tool", toolName),
        )
        
        result := s.toolRunner.RunToolByName(ctx, toolName, repoPath, languages)
        
        s.log.Info("scan_tool_complete",
            slog.String("job_id", jobID),
            slog.String("tool", toolName),
            slog.Int("finding_count", len(result.Findings)),
            slog.Bool("timed_out", result.TimedOut),
            slog.Bool("success", result.Error == nil),
            slog.Duration("duration", time.Since(toolStart)),
        )
        
        if result.Error != nil {
            s.log.Warn("scan_tool_error",
                slog.String("job_id", jobID),
                slog.String("tool", toolName),
                slog.String("error", result.Error.Error()),
            )
        }
        
        results = append(results, result)
    }
    
    // Phase 4: Aggregation
    s.log.Info("scan_phase_aggregate_start",
        slog.String("job_id", jobID),
        slog.Int("result_count", len(results)),
    )
    aggStart := time.Now()
    findings := s.aggregator.AggregateAndProcess(results)
    
    // Count by severity
    severityCounts := map[string]int{"critical": 0, "high": 0, "medium": 0, "low": 0}
    for _, f := range findings {
        severityCounts[f.Severity]++
    }
    
    s.log.Info("scan_phase_aggregate_complete",
        slog.String("job_id", jobID),
        slog.Int("total_findings", len(findings)),
        slog.Int("critical", severityCounts["critical"]),
        slog.Int("high", severityCounts["high"]),
        slog.Int("medium", severityCounts["medium"]),
        slog.Int("low", severityCounts["low"]),
        slog.Duration("duration", time.Since(aggStart)),
    )
    
    // Phase 5: AI Review (if applicable)
    if len(findings) > 0 && s.reviewer.HasClient() {
        s.log.Info("scan_phase_review_start",
            slog.String("job_id", jobID),
            slog.Int("findings_to_review", len(findings)),
        )
        reviewStart := time.Now()
        
        findings, err = s.reviewer.Review(ctx, repoPath, findings)
        if err != nil {
            s.log.Warn("scan_phase_review_partial",
                slog.String("job_id", jobID),
                slog.String("error", err.Error()),
            )
        }
        
        s.log.Info("scan_phase_review_complete",
            slog.String("job_id", jobID),
            slog.Int("reviewed_findings", len(findings)),
            slog.Duration("duration", time.Since(reviewStart)),
        )
    } else {
        s.log.Debug("scan_phase_review_skipped",
            slog.String("job_id", jobID),
            slog.Bool("has_findings", len(findings) > 0),
            slog.Bool("has_client", s.reviewer.HasClient()),
        )
    }
    
    // Complete
    _ = s.completeJob(ctx, jobID, findings)
    
    s.log.Info("scan_pipeline_complete",
        slog.String("job_id", jobID),
        slog.Int("total_findings", len(findings)),
        slog.Duration("total_duration", time.Since(start)),
    )
}
```

### 11. OpenAI Client Logging

```go
// ChatCompletion with detailed logging
func (c *Client) ChatCompletion(ctx context.Context, messages []Message) (string, error) {
    requestID := logger.GetRequestID(ctx)
    start := time.Now()
    
    // Calculate prompt metrics
    promptLength := 0
    for _, m := range messages {
        promptLength += len(m.Content)
    }
    
    c.log.Info("openai_request_start",
        slog.String("request_id", requestID),
        slog.String("model", c.model),
        slog.Int("prompt_length", promptLength),
        slog.Int("message_count", len(messages)),
        slog.String("reasoning_effort", string(c.reasoningEffort)),
    )
    
    // Debug: truncated preview (first 500 chars of last message)
    if len(messages) > 0 {
        lastMsg := messages[len(messages)-1].Content
        preview := lastMsg
        if len(preview) > 500 {
            preview = preview[:500] + "..."
        }
        c.log.Debug("openai_request_preview",
            slog.String("request_id", requestID),
            slog.String("prompt_preview", preview),
        )
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        c.log.Error("openai_request_failed",
            slog.String("request_id", requestID),
            slog.String("error", err.Error()),
            slog.Duration("duration", time.Since(start)),
        )
        return "", err
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    if resp.StatusCode != http.StatusOK {
        c.log.Error("openai_response_error",
            slog.String("request_id", requestID),
            slog.Int("status_code", resp.StatusCode),
            slog.Duration("latency", time.Since(start)),
        )
        return "", fmt.Errorf("status %d", resp.StatusCode)
    }
    
    c.log.Info("openai_response_received",
        slog.String("request_id", requestID),
        slog.Int("status_code", resp.StatusCode),
        slog.Int("response_length", len(body)),
        slog.Duration("latency", time.Since(start)),
    )
    
    // Debug: truncated response preview
    text := extractTextFromResponse(responsesResp)
    if len(text) > 500 {
        c.log.Debug("openai_response_preview",
            slog.String("request_id", requestID),
            slog.String("response_preview", text[:500]+"..."),
        )
    }
    
    return text, nil
}
```

### 12. Queue and Rate Limiter Logging

```go
// RequestQueue logging
func (q *RequestQueue) Acquire(ctx context.Context) error {
    requestID := logger.GetRequestID(ctx)
    
    q.log.Debug("queue_acquire_start",
        slog.String("request_id", requestID),
        slog.Int("available", q.Available()),
        slog.Int64("waiting", q.waiting.Load()),
    )
    
    q.waiting.Add(1)
    defer q.waiting.Add(-1)
    
    select {
    case q.semaphore <- struct{}{}:
        q.log.Debug("queue_acquire_success",
            slog.String("request_id", requestID),
            slog.Int("available_after", q.Available()),
        )
        return nil
    case <-ctx.Done():
        q.log.Warn("queue_acquire_timeout",
            slog.String("request_id", requestID),
            slog.String("error", ctx.Err().Error()),
        )
        return ctx.Err()
    }
}

func (q *RequestQueue) Release() {
    q.log.Debug("queue_release",
        slog.Int("available_after", q.Available()+1),
        slog.Int64("processed", q.processed.Load()+1),
    )
    // ... release logic ...
}

// RateLimiter logging
func (l *Limiter) Allow(ip string) (bool, time.Duration) {
    // Hash IP for privacy
    ipHash := hashIP(ip)
    
    // ... existing logic ...
    
    if !allowed {
        l.log.Warn("rate_limit_denied",
            slog.String("ip_hash", ipHash),
            slog.Int("count", state.count),
            slog.Int("limit", l.limit),
            slog.Duration("retry_after", retryAfter),
        )
    } else {
        l.log.Debug("rate_limit_allowed",
            slog.String("ip_hash", ipHash),
            slog.Int("remaining", l.limit-state.count),
        )
    }
    
    return allowed, retryAfter
}
```

### 13. Application Startup Logging

```go
// main.go startup logging
func main() {
    // Initialize logger first
    logCfg := logger.Config{
        Level:       parseLogLevel(os.Getenv("LOG_LEVEL")),
        LogDir:      "./logs",
        MaxSizeMB:   100,
        MaxAgeDays:  7,
        EnableColor: true,
    }
    log, err := logger.New(logCfg)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
        os.Exit(1)
    }
    defer log.Close()
    
    log.App().Info("application_starting",
        slog.String("version", version),
        slog.String("log_level", logCfg.Level.String()),
        slog.String("log_dir", logCfg.LogDir),
    )
    
    // Database connection
    log.App().Info("database_connecting")
    if err := db.Connect(ctx); err != nil {
        log.App().Error("database_connection_failed", slog.String("error", err.Error()))
        os.Exit(1)
    }
    log.App().Info("database_connected",
        slog.Int("max_open_conns", defaultMaxOpenConns),
        slog.Int("max_idle_conns", defaultMaxIdleConns),
    )
    
    // Service initialization
    log.App().Info("services_initializing")
    // ... initialize services with logger ...
    log.App().Info("services_initialized",
        slog.Bool("generation_enabled", routerCfg.GenerationService != nil),
        slog.Bool("gallery_enabled", routerCfg.GalleryService != nil),
        slog.Bool("scanner_enabled", routerCfg.ScannerService != nil),
    )
    
    // Server start
    log.App().Info("server_starting", slog.String("port", port))
    
    // Graceful shutdown
    sig := <-shutdown
    log.App().Info("shutdown_signal_received", slog.String("signal", sig.String()))
    
    if err := server.Shutdown(shutdownCtx); err != nil {
        log.App().Error("shutdown_error", slog.String("error", err.Error()))
    } else {
        log.App().Info("server_stopped_gracefully")
    }
    
    if err := db.Close(); err != nil {
        log.App().Error("database_close_error", slog.String("error", err.Error()))
    } else {
        log.App().Info("database_connection_closed")
    }
    
    log.App().Info("application_stopped")
}
```

## Data Models

### Log Entry Structure (JSON)

```json
{
  "time": "2026-01-13T10:30:45.123Z",
  "level": "INFO",
  "msg": "request_complete",
  "request_id": "abc123def456",
  "component": "http",
  "method": "POST",
  "path": "/api/generate/questions",
  "status": 200,
  "duration_ms": 1523,
  "remote_addr": "192.168.1.100"
}
```

### Log File Structure

```
logs/
├── README.md
├── 2026-01-13-app.log
├── 2026-01-13-http.log
├── 2026-01-13-db.log
├── 2026-01-13-scanner.log
└── 2026-01-13-client.log
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*



### Property 1: Log Entry JSON Structure
*For any* log entry written to a file, parsing it as JSON SHALL succeed and the resulting object SHALL contain `time`, `level`, `msg`, and `component` fields.
**Validates: Requirements 1.4, 1.5**

### Property 2: HTTP Log Completeness
*For any* HTTP request processed by the middleware, the resulting log entries SHALL contain `method`, `path`, `status`, `duration`, and `request_id` fields.
**Validates: Requirements 2.1, 2.3**

### Property 3: Request ID Uniqueness
*For any* two distinct HTTP requests, the generated request IDs SHALL be different.
**Validates: Requirements 2.2**

### Property 4: Request ID Propagation
*For any* operation triggered by an HTTP request, the log entry SHALL contain the same `request_id` as the originating request.
**Validates: Requirements 2.4**

### Property 5: Sensitive Data Redaction
*For any* log entry, if the original data contained fields named `password`, `token`, `api_key`, `secret`, or `authorization`, the logged values SHALL be replaced with `[REDACTED]`.
**Validates: Requirements 2.5, 3.5**

### Property 6: Database Log Completeness
*For any* database operation, the log entry SHALL contain `type` (SELECT/INSERT/UPDATE/DELETE), `duration`, and `success` fields.
**Validates: Requirements 3.1, 3.2**

### Property 7: Service Log Completeness
*For any* service method invocation, log entries SHALL be created for both start (with method name) and completion (with outcome and duration).
**Validates: Requirements 4.1, 4.2**

### Property 8: OpenAI Log Completeness
*For any* OpenAI API call, log entries SHALL contain `model`, `prompt_length` on request and `status`, `latency` on response.
**Validates: Requirements 5.1, 5.2**

### Property 9: Content Truncation at INFO Level
*For any* log entry at INFO level involving OpenAI prompts or responses, the full content SHALL NOT be included in the log.
**Validates: Requirements 5.4**

### Property 10: Scanner Phase Logging
*For any* complete scan operation, log entries SHALL exist for each phase: initiation, cloning, language detection, tool execution, and aggregation.
**Validates: Requirements 6.1, 6.2, 6.3, 6.4, 6.5**

### Property 11: Frontend Error Capture
*For any* JavaScript error or failed API call in the frontend, a log entry SHALL be created containing `message`, `url`, and `timestamp`.
**Validates: Requirements 7.1, 7.2, 7.3**

### Property 12: Log File Routing
*For any* log entry, it SHALL be written to the correct file based on its component: HTTP logs to `http.log`, database logs to `db.log`, scanner logs to `scanner.log`, client logs to `client.log`, and all others to `app.log`.
**Validates: Requirements 8.1, 8.2**

### Property 13: Color Mapping by Level
*For any* log entry written to a TTY, the ANSI color code SHALL match the level: ERROR→RED, WARN→YELLOW, INFO→GREEN, DEBUG→CYAN.
**Validates: Requirements 10.1, 10.2, 10.3, 10.4, 10.5**

### Property 14: Level Filtering
*For any* configured log level L, only log entries with severity >= L SHALL be written. DEBUG includes all; ERROR excludes DEBUG, INFO, WARN.
**Validates: Requirements 9.4, 9.5**

## Error Handling

### Logger Initialization Errors
- If log directory creation fails: Log to stderr and continue with console-only logging
- If file open fails: Retry with exponential backoff, fall back to console after 3 attempts

### Log Rotation Errors
- If rotation fails: Continue writing to current file, log warning to stderr
- If cleanup fails: Log warning but don't block normal operation

### Frontend Log Submission Errors
- If POST fails: Buffer logs locally (up to 100 entries), retry on next flush
- If buffer overflows: Drop oldest entries, keep most recent

### Sensitive Data Handling
- Redaction happens before any I/O operation
- Unknown sensitive fields: Log field names only, not values

## Testing Strategy

### Unit Tests
- Logger initialization with various configurations
- Log level filtering behavior
- JSON output format validation
- Color code application
- Sensitive data redaction
- Log rotation trigger conditions

### Property-Based Tests
Using Go's `testing/quick` package:

1. **JSON Structure Property**: Generate random log entries, verify JSON parsing succeeds
2. **Request ID Uniqueness**: Generate N request IDs, verify all unique
3. **Sensitive Data Redaction**: Generate entries with sensitive fields, verify redaction
4. **Level Filtering**: Generate entries at all levels, verify correct filtering
5. **Color Mapping**: Generate entries at each level, verify correct color codes

### Integration Tests
- Full request flow with log correlation
- Database operation logging
- Scanner pipeline logging
- Frontend error submission and logging

### Test Configuration
- Property tests: Minimum 100 iterations
- Each test tagged with: `Feature: comprehensive-logging, Property N: {description}`
- Use `testing/quick` for Go property tests
- Use `fast-check` for TypeScript property tests
