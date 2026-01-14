# Implementation Plan: Documentation and Configuration

## Overview

This plan implements a centralized TOML configuration system and comprehensive documentation. Tasks are ordered to build incrementally: config package first, then integration with existing services, then documentation.

## Tasks

- [x] 1. Create configuration package
  - [x] 1.1 Create config structs and default values
    - Create `backend/internal/config/config.go` with all config structs
    - Implement `DefaultConfig()` returning current hardcoded values
    - Add TOML struct tags for all fields
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8_

  - [x] 1.2 Implement config loading logic
    - Implement `Load()` function with file detection
    - Implement `LoadFromPath()` for custom paths
    - Implement `ApplyEnvironmentOverrides()` for env var precedence
    - Support CONFIG_PATH environment variable
    - _Requirements: 2.1, 2.2, 2.5, 8.1, 8.2, 8.3, 8.4, 8.5_

  - [x] 1.3 Implement config validation
    - Implement `Validate()` with range checks for all numeric fields
    - Implement enum validation for string fields (log level, reasoning effort, etc.)
    - Return detailed error messages identifying specific invalid fields
    - _Requirements: 1.9, 7.1, 7.2, 7.3, 7.4, 7.5_

  - [x] 1.4 Implement config logging with redaction
    - Implement `LogConfig()` method for startup logging
    - Ensure no sensitive values (API keys, tokens) appear in logs
    - _Requirements: 2.4_

  - [x] 1.5 Write property tests for config package
    - **Property 1: Config structure completeness**
    - **Property 2: Default value fallback**
    - **Property 3: Invalid config rejection**
    - **Property 4: Environment variable override**
    - **Property 6: Config round-trip**
    - **Validates: Requirements 1.1-1.9, 2.2, 7.1-7.3**

- [x] 2. Create example configuration file
  - Create `config.example.toml` with all options documented
  - Include comments explaining each setting
  - Show default values for reference
  - _Requirements: 1.1-1.7_

- [x] 3. Integrate config with main.go
  - [x] 3.1 Load config at startup
    - Call `config.Load()` before initializing services
    - Log loaded configuration
    - Exit with code 1 on validation failure
    - _Requirements: 2.1, 2.3, 7.4, 7.5_

  - [x] 3.2 Pass config to services
    - Update service initialization to accept config values
    - Remove hardcoded constants from main.go
    - _Requirements: 1.1-1.7_

- [x] 4. Integrate config with OpenAI client
  - Update `NewClient()` to accept config values
  - Use config for model, timeout, reasoning effort, verbosity
  - Remove hardcoded defaults from openai package
  - _Requirements: 1.2_

- [x] 5. Integrate config with rate limiter
  - Update `NewLimiter()` to accept config values
  - Use config for generation, rating, and scan limits
  - Remove hardcoded defaults from ratelimit package
  - _Requirements: 1.3_

- [x] 6. Integrate config with logger
  - Update `New()` to accept config values
  - Use config for level, directory, max size, max age, color
  - Remove hardcoded defaults from logger package
  - _Requirements: 1.4_

- [x] 7. Integrate config with scanner
  - [x] 7.1 Update scanner service
    - Use config for max repo size, max review files, retention days
    - Remove hardcoded defaults from scanner package
    - _Requirements: 1.5_

  - [x] 7.2 Update cloner
    - Use config for clone timeout, max repo size
    - Remove hardcoded defaults from cloner
    - _Requirements: 1.5_

  - [x] 7.3 Update code reviewer
    - Use config for max review files, code review model
    - Remove hardcoded defaults from reviewer
    - _Requirements: 1.5_

  - [x] 7.4 Update tool runner
    - Use config for tool timeout
    - Remove hardcoded defaults from tool runner
    - _Requirements: 1.5_

- [x] 8. Integrate config with generation service
  - Use config for max lengths, question counts, retries
  - Remove hardcoded constants from generation package
  - _Requirements: 1.6_

- [x] 9. Integrate config with gallery service
  - Use config for page size, default sort
  - Remove hardcoded defaults from gallery package
  - _Requirements: 1.7_

- [x] 10. Checkpoint - Ensure all tests pass
  - Run `go test ./...` in backend
  - Verify application starts with default config
  - Verify application starts with custom config.toml
  - Ensure all tests pass, ask the user if questions arise.

- [x] 11. Update Docker configuration
  - [x] 11.1 Update docker-compose.yml
    - Add volume mount for config.toml
    - Ensure env vars still work for secrets
    - _Requirements: 3.1, 3.2_

  - [x] 11.2 Update docker-compose.prod.yml
    - Add volume mount for config.toml
    - Ensure production build includes config
    - _Requirements: 3.1, 3.3_

  - [x] 11.3 Update Dockerfile.prod
    - Copy config.example.toml as default config
    - _Requirements: 3.3_

- [x] 12. Update .env.example
  - Document relationship between .env and config.toml
  - Clarify which values go where (secrets in .env, settings in config.toml)
  - _Requirements: 8.1-8.6_

- [x] 13. Create developer documentation
  - [x] 13.1 Write architecture overview
    - Create `docs/developer.md`
    - Include architecture diagram (Mermaid)
    - Explain component relationships
    - _Requirements: 5.1_

  - [x] 13.2 Document backend packages
    - Explain each package in backend/internal/
    - Include purpose and key interfaces
    - _Requirements: 5.2_

  - [x] 13.3 Document frontend structure
    - Explain pages, components, and lib folders
    - Include component hierarchy
    - _Requirements: 5.3_

  - [x] 13.4 Document API endpoints
    - Update `docs/api.md` with complete endpoint list
    - Include request/response examples
    - _Requirements: 5.4_

  - [x] 13.5 Document database schema
    - Explain tables and relationships
    - Document migration process
    - _Requirements: 5.5_

  - [x] 13.6 Document logging system
    - Explain log categories and files
    - Document log format and levels
    - _Requirements: 5.6_

  - [x] 13.7 Document customization
    - Explain how to add features
    - Document AI prompt customization
    - _Requirements: 5.7, 5.8_

- [x] 14. Create self-hosting guide
  - [x] 14.1 Write configuration reference
    - Create `docs/self-hosting.md`
    - Document all config.toml options with descriptions
    - Include default values and valid ranges
    - _Requirements: 6.1_

  - [x] 14.2 Write example configurations
    - Minimal config (generation only)
    - Full-featured config (all features)
    - High-security config (strict limits)
    - _Requirements: 6.2_

  - [x] 14.3 Document database setup
    - PostgreSQL installation and configuration
    - Migration instructions
    - _Requirements: 6.3_

  - [x] 14.4 Document OpenAI setup
    - API key configuration
    - Model selection guidance
    - _Requirements: 6.4_

  - [x] 14.5 Document private repo scanning
    - GitHub token setup
    - Permissions required
    - _Requirements: 6.5_

  - [x] 14.6 Document resource requirements
    - CPU, memory, disk recommendations
    - Scaling considerations
    - _Requirements: 6.6_

  - [x] 14.7 Document maintenance
    - Backup procedures
    - Log rotation
    - Database cleanup
    - _Requirements: 6.7_

  - [x] 14.8 Document scanner customization
    - Adding/removing security tools
    - Adjusting tool configurations
    - _Requirements: 6.8_

- [x] 15. Update README
  - [x] 15.1 Write project overview
    - Clear explanation of what BetterKiroPrompts does
    - Target audience (beginners, self-hosters)
    - _Requirements: 4.1_

  - [x] 15.2 Write prerequisites section
    - Exact version requirements
    - Required tools (Docker, Node.js, Go)
    - _Requirements: 4.2_

  - [x] 15.3 Write quick start guide
    - Copy-paste commands
    - Minimal steps to get running
    - _Requirements: 4.3_

  - [x] 15.4 Write project structure
    - Directory layout
    - Key files explained
    - _Requirements: 4.4_

  - [x] 15.5 Add documentation links
    - Link to api.md, developer.md, self-hosting.md
    - _Requirements: 4.5_

  - [x] 15.6 Write configuration section
    - Explain .env vs config.toml
    - Quick reference for common settings
    - _Requirements: 4.6_

  - [x] 15.7 Write troubleshooting section
    - Common errors and solutions
    - API key issues, database connection, timeouts
    - _Requirements: 4.7_

  - [x] 15.8 Document main features
    - Generation, Gallery, Security Scanning
    - Brief description of each
    - _Requirements: 4.8_

- [x] 16. Final checkpoint
  - Run full test suite
  - Verify documentation accuracy
  - Test fresh install with documentation
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
