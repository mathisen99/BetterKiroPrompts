# Requirements Document

## Introduction

This feature creates comprehensive documentation and a centralized TOML configuration file for the BetterKiroPrompts project. The project is a complete AI-driven tool that generates kickoff prompts, steering files, hooks, and AGENTS.md for Kiro users. It includes:

- Single-page AI-driven generation flow (project idea → questions → outputs)
- Experience level selection (beginner, novice, expert)
- Hook preset selection (light, basic, default, strict)
- Public gallery with ratings and view counts
- Security scanning with local tools and AI-powered code review
- Comprehensive structured logging system

The goal is to eliminate hardcoded values from the codebase, provide clear documentation for developers and self-hosters, and make the application fully configurable without code changes.

## Glossary

- **Config_Loader**: The component responsible for reading and parsing the config.toml file at application startup
- **Config_File**: The config.toml file containing all configurable application settings
- **README**: The main project documentation file in the repository root
- **Developer_Docs**: Documentation in the docs/ folder explaining architecture and customization
- **Self_Hosting_Guide**: Documentation explaining how to deploy and configure the application
- **Environment_Variables**: Runtime configuration values passed via .env or container environment

## Requirements

### Requirement 1: Centralized Configuration File

**User Story:** As a self-hoster, I want all configurable values in a single config.toml file, so that I can customize the application without modifying source code.

#### Acceptance Criteria

1. THE Config_File SHALL contain server configuration (port, host, shutdown timeout)
2. THE Config_File SHALL contain OpenAI configuration (model name, base URL, timeout, reasoning effort, verbosity, code review model)
3. THE Config_File SHALL contain rate limiting configuration (generation limit per hour, rating limit per hour, scan limit per hour)
4. THE Config_File SHALL contain logging configuration (level, directory, max file size MB, max age days, enable color)
5. THE Config_File SHALL contain scanner configuration (max repo size MB, max review files, tool timeout seconds, result retention days, clone timeout)
6. THE Config_File SHALL contain generation configuration (max project idea length, max answer length, min/max questions, max retries)
7. THE Config_File SHALL contain gallery configuration (page size, default sort order)
8. WHEN the Config_File is missing THEN THE Config_Loader SHALL use sensible default values matching current hardcoded values
9. WHEN the Config_File contains invalid values THEN THE Config_Loader SHALL log an error and exit with a non-zero status

### Requirement 2: Configuration Loading

**User Story:** As a developer, I want the application to load configuration at startup, so that I can change settings without recompiling.

#### Acceptance Criteria

1. THE Config_Loader SHALL read config.toml from the application root directory by default
2. THE Config_Loader SHALL allow environment variables to override config.toml values
3. THE Config_Loader SHALL validate all configuration values at startup
4. THE Config_Loader SHALL log the loaded configuration (with sensitive values redacted)
5. WHEN environment variable CONFIG_PATH is set THEN THE Config_Loader SHALL read from that path instead

### Requirement 3: Docker Integration

**User Story:** As a DevOps engineer, I want the Docker containers to receive configuration correctly, so that I can deploy with custom settings.

#### Acceptance Criteria

1. THE docker-compose.yml files SHALL mount config.toml as a volume
2. THE Docker_Compose files SHALL support environment variable overrides for secrets (API keys, tokens)
3. THE Dockerfile.prod SHALL copy config.toml into the image at build time
4. WHEN config.toml is mounted as a volume THEN changes SHALL be reflected without rebuilding

### Requirement 4: README Documentation

**User Story:** As a new user, I want a comprehensive README, so that I can understand and use the project quickly.

#### Acceptance Criteria

1. THE README SHALL contain a project overview explaining what BetterKiroPrompts does (generates kickoff prompts, steering files, hooks, AGENTS.md)
2. THE README SHALL contain prerequisites with exact version requirements (Go 1.25.5, Node.js 24.12.0, PostgreSQL 18.1, Docker)
3. THE README SHALL contain quick start instructions with copy-paste commands
4. THE README SHALL contain a project structure overview showing backend/frontend/docs layout
5. THE README SHALL contain links to detailed documentation in docs/
6. THE README SHALL contain a configuration section explaining config.toml and .env
7. THE README SHALL contain a troubleshooting section for common issues (API key errors, database connection, timeouts)
8. THE README SHALL explain the three main features: Generation, Gallery, Security Scanning

### Requirement 5: Developer Documentation

**User Story:** As a developer wanting to contribute or customize, I want detailed architecture documentation, so that I can understand how the system works.

#### Acceptance Criteria

1. THE Developer_Docs SHALL contain an architecture overview with component diagram showing frontend, backend, database, scanner container
2. THE Developer_Docs SHALL explain each backend package (api, db, gallery, generation, logger, openai, prompts, queue, ratelimit, sanitize, scanner, storage)
3. THE Developer_Docs SHALL explain the frontend component structure (pages, components, lib)
4. THE Developer_Docs SHALL document all API endpoints with request/response examples
5. THE Developer_Docs SHALL explain the database schema and migrations
6. THE Developer_Docs SHALL document the logging system (app, http, db, scanner, client logs)
7. THE Developer_Docs SHALL explain how to add new features or modify existing ones
8. THE Developer_Docs SHALL explain the AI prompt system and how to customize prompts

### Requirement 6: Self-Hosting Guide

**User Story:** As a self-hoster, I want a deployment guide, so that I can run this tool for my team.

#### Acceptance Criteria

1. THE Self_Hosting_Guide SHALL explain all configuration options in config.toml with descriptions and defaults
2. THE Self_Hosting_Guide SHALL provide example configurations for different use cases (minimal, full-featured, high-security)
3. THE Self_Hosting_Guide SHALL explain how to set up the PostgreSQL database
4. THE Self_Hosting_Guide SHALL explain how to configure OpenAI API access and model selection
5. THE Self_Hosting_Guide SHALL explain how to enable private repository scanning with GitHub tokens
6. THE Self_Hosting_Guide SHALL document resource requirements (CPU, memory, disk) for different scales
7. THE Self_Hosting_Guide SHALL explain backup and maintenance procedures for the database
8. THE Self_Hosting_Guide SHALL explain how to customize the security scanner tools

### Requirement 7: Configuration Validation

**User Story:** As an operator, I want configuration errors to be caught early, so that I don't deploy with invalid settings.

#### Acceptance Criteria

1. WHEN a required configuration value is missing THEN THE Config_Loader SHALL report which value is missing
2. WHEN a numeric configuration value is out of range THEN THE Config_Loader SHALL report the valid range
3. WHEN a string configuration value has an invalid format THEN THE Config_Loader SHALL report the expected format
4. THE Config_Loader SHALL validate all values before the application starts serving requests
5. IF configuration validation fails THEN THE Config_Loader SHALL exit with status code 1

### Requirement 8: Backward Compatibility

**User Story:** As an existing user, I want my current .env configuration to continue working, so that I don't have to reconfigure everything.

#### Acceptance Criteria

1. THE Config_Loader SHALL continue to read DATABASE_URL from environment variables
2. THE Config_Loader SHALL continue to read OPENAI_API_KEY from environment variables
3. THE Config_Loader SHALL continue to read GITHUB_TOKEN from environment variables
4. THE Config_Loader SHALL continue to read LOG_LEVEL from environment variables
5. WHEN both config.toml and environment variables specify a value THEN environment variables SHALL take precedence for secrets
6. THE Config_Loader SHALL NOT log deprecation warnings for environment variables (they are the preferred method for secrets)
