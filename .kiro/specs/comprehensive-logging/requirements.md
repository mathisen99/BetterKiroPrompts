# Requirements Document

## Introduction

This feature implements comprehensive file-based logging across the entire application stack (Go backend and React frontend). Logs will be written to a `logs/` directory outside Docker containers, enabling persistent troubleshooting without container access. The logging system will capture the full flow of operations including API requests, database operations, service calls, OpenAI interactions, scanner operations, and frontend errors.

## Glossary

- **Logger**: The centralized logging component that writes structured log entries to files
- **Log_Level**: Severity classification (DEBUG, INFO, WARN, ERROR) for filtering log output
- **Log_Entry**: A single structured log record containing timestamp, level, context, and message
- **Request_ID**: A unique identifier assigned to each HTTP request for tracing across components
- **Log_Rotation**: The process of archiving old log files and creating new ones based on size or time
- **Structured_Logging**: Logging format using key-value pairs (JSON) for machine-parseable output

## Requirements

### Requirement 1: Centralized Backend Logging Infrastructure

**User Story:** As a developer, I want a centralized logging system in the Go backend, so that all components log consistently to external files.

#### Acceptance Criteria

1. THE Logger SHALL write all log entries to files in a `logs/` directory mounted from the host filesystem
2. WHEN the Logger initializes, THE Logger SHALL create the log directory if it does not exist
3. THE Logger SHALL support four log levels: DEBUG, INFO, WARN, ERROR
4. WHEN a log entry is created, THE Logger SHALL include timestamp, log level, component name, and message
5. THE Logger SHALL output logs in JSON format for structured parsing
6. WHEN the log file exceeds 100MB, THE Logger SHALL rotate to a new file with timestamp suffix
7. THE Logger SHALL retain log files for 7 days before automatic cleanup

### Requirement 2: HTTP Request/Response Logging

**User Story:** As a developer, I want all HTTP requests and responses logged, so that I can trace API interactions end-to-end.

#### Acceptance Criteria

1. WHEN an HTTP request is received, THE Logger SHALL log the method, path, query parameters, and request headers
2. WHEN an HTTP request is received, THE Logger SHALL generate and attach a unique Request_ID
3. WHEN an HTTP response is sent, THE Logger SHALL log the status code, response time, and Request_ID
4. THE Logger SHALL propagate the Request_ID through all downstream operations for correlation
5. IF a request body contains sensitive data, THEN THE Logger SHALL redact passwords, tokens, and API keys

### Requirement 3: Database Operation Logging

**User Story:** As a developer, I want database operations logged, so that I can troubleshoot data access issues.

#### Acceptance Criteria

1. WHEN a database query executes, THE Logger SHALL log the query type (SELECT, INSERT, UPDATE, DELETE)
2. WHEN a database query executes, THE Logger SHALL log the execution duration in milliseconds
3. IF a database error occurs, THEN THE Logger SHALL log the error message and stack trace
4. WHEN a database connection is established or closed, THE Logger SHALL log the event
5. THE Logger SHALL NOT log raw SQL with user-provided values to prevent sensitive data exposure

### Requirement 4: Service Layer Logging

**User Story:** As a developer, I want service operations logged, so that I can understand business logic execution flow.

#### Acceptance Criteria

1. WHEN a service method is called, THE Logger SHALL log the method name and input parameters summary
2. WHEN a service method completes, THE Logger SHALL log the outcome (success/failure) and duration
3. WHEN the Generation_Service calls OpenAI, THE Logger SHALL log the prompt type and token usage
4. WHEN the Scanner_Service processes a repository, THE Logger SHALL log each scanning phase
5. WHEN the Gallery_Service retrieves or stores data, THE Logger SHALL log the operation type and item count

### Requirement 5: External API Logging (OpenAI)

**User Story:** As a developer, I want OpenAI API interactions logged, so that I can debug AI generation issues.

#### Acceptance Criteria

1. WHEN an OpenAI request is made, THE Logger SHALL log the model, prompt length, and request timestamp
2. WHEN an OpenAI response is received, THE Logger SHALL log the response status, token usage, and latency
3. IF an OpenAI error occurs, THEN THE Logger SHALL log the error type, message, and retry attempt number
4. THE Logger SHALL NOT log full prompt or response content at INFO level to manage log size
5. WHEN DEBUG level is enabled, THE Logger SHALL log truncated prompt/response previews (first 500 chars)

### Requirement 6: Scanner Pipeline Logging

**User Story:** As a developer, I want the security scanner pipeline fully logged, so that I can troubleshoot scan failures.

#### Acceptance Criteria

1. WHEN a scan is initiated, THE Logger SHALL log the repository URL and scan configuration
2. WHEN the Cloner clones a repository, THE Logger SHALL log clone progress and completion status
3. WHEN the Language_Detector analyzes files, THE Logger SHALL log detected languages and file counts
4. WHEN the Reviewer processes files, THE Logger SHALL log files reviewed and findings count
5. WHEN the Aggregator combines results, THE Logger SHALL log the aggregation summary
6. IF any scanner component fails, THEN THE Logger SHALL log the failure point and error details

### Requirement 7: Frontend Error Logging

**User Story:** As a developer, I want frontend errors sent to the backend for logging, so that I can troubleshoot client-side issues.

#### Acceptance Criteria

1. WHEN a JavaScript error occurs, THE Frontend SHALL capture the error and stack trace
2. WHEN a React component throws an error, THE ErrorBoundary SHALL send error details to the backend
3. WHEN an API call fails, THE Frontend SHALL log the endpoint, status, and error message
4. THE Frontend SHALL batch error logs and send them every 5 seconds to reduce network overhead
5. THE Backend SHALL provide a `/api/logs/client` endpoint to receive frontend error logs

### Requirement 8: Log File Organization

**User Story:** As a developer, I want logs organized by type, so that I can quickly find relevant information.

#### Acceptance Criteria

1. THE Logger SHALL write to separate files: `app.log`, `http.log`, `db.log`, `scanner.log`, `client.log`
2. THE Logger SHALL prefix each log file with the date in YYYY-MM-DD format
3. WHEN Docker containers start, THE Logger SHALL mount the host `./logs` directory
4. THE `.gitignore` file SHALL include the `logs/` directory to prevent committing logs
5. THE Logger SHALL create a `logs/README.md` explaining log file structure and retention

### Requirement 9: Log Level Configuration

**User Story:** As a developer, I want to configure log verbosity, so that I can adjust detail level for different environments.

#### Acceptance Criteria

1. THE Logger SHALL read the log level from the `LOG_LEVEL` environment variable
2. WHEN `LOG_LEVEL` is not set, THE Logger SHALL default to INFO level
3. THE Logger SHALL support runtime log level changes without restart via `/api/admin/log-level` endpoint
4. WHEN log level is DEBUG, THE Logger SHALL include additional context like goroutine IDs
5. WHEN log level is ERROR, THE Logger SHALL only log errors and critical failures

### Requirement 10: Colored Console Output

**User Story:** As a developer, I want colored log output in the console, so that I can quickly identify log levels and important information.

#### Acceptance Criteria

1. WHEN logs are written to stdout/stderr, THE Logger SHALL apply ANSI color codes based on log level
2. THE Logger SHALL use RED for ERROR level logs
3. THE Logger SHALL use YELLOW for WARN level logs
4. THE Logger SHALL use GREEN for INFO level logs
5. THE Logger SHALL use CYAN for DEBUG level logs
6. THE Logger SHALL highlight Request_IDs in MAGENTA for easy correlation
7. THE Logger SHALL highlight timestamps in GRAY for visual separation
8. WHEN the `NO_COLOR` environment variable is set, THE Logger SHALL disable colored output
9. WHEN output is not a TTY (piped or redirected), THE Logger SHALL disable colored output automatically
