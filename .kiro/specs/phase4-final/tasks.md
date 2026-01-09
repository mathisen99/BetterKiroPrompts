# Implementation Plan: Phase 4 - Final Polish

## Overview

This phase delivers professional UI design, smarter AI generation with experience-level adaptation, and complete Kiro file coverage with proper formats.

## Tasks

- [ ] 1. Add Experience Level Selection to Frontend
  - [ ] 1.1 Create ExperienceLevelSelector component with Beginner/Novice/Expert cards
    - Use shadcn Card components with icons
    - Dark theme with blue accents
    - _Requirements: 1.1, 2.1, 2.2, 2.3_
  - [ ] 1.2 Update LandingPage state to include experienceLevel and new phase 'level-select'
    - Add experienceLevel to state interface
    - Start flow at level-select phase
    - _Requirements: 1.5_
  - [ ] 1.3 Update API types to include experienceLevel in requests
    - Add to GenerateQuestionsRequest
    - Add to GenerateOutputsRequest
    - _Requirements: 1.2, 1.3, 1.4_

- [ ] 2. Add Hook Preset Selection to Frontend
  - [ ] 2.1 Create HookPresetSelector component with Light/Basic/Default/Strict options
    - Radio group with descriptions
    - Default to 'default' preset
    - _Requirements: 5.7, 10.6_
  - [ ] 2.2 Add hookPreset to LandingPage state and wire to output generation
    - _Requirements: 5.7_

- [ ] 3. Professional UI Redesign
  - [ ] 3.1 Create Header component with branding and navigation
    - Logo, title, dark theme
    - _Requirements: 2.6_
  - [ ] 3.2 Update index.css with refined shadcn blue dark theme variables
    - Match the reference image colors
    - _Requirements: 2.1, 2.2_
  - [ ] 3.3 Redesign ProjectInput with better card styling and spacing
    - Larger input, better examples display
    - _Requirements: 2.3, 2.4, 2.5_
  - [ ] 3.4 Redesign QuestionFlow with improved card layout
    - Better progress indicator, cleaner Q&A cards
    - _Requirements: 2.3, 2.4_
  - [ ] 3.5 Ensure responsive design works on mobile/tablet/desktop
    - _Requirements: 2.7_

- [ ] 4. Checkpoint - Frontend UI Complete
  - Verify all UI components render correctly
  - Test responsive behavior
  - Ensure dark theme with blue accents is consistent

- [ ] 5. Update Backend API to Accept Experience Level and Hook Preset
  - [ ] 5.1 Update GenerateQuestionsRequest struct to include ExperienceLevel
    - Add validation for valid values
    - _Requirements: 1.2, 1.3, 1.4_
  - [ ] 5.2 Update GenerateOutputsRequest struct to include ExperienceLevel and HookPreset
    - Add validation for valid values
    - _Requirements: 5.7, 6.6_
  - [ ] 5.3 Update API handlers to pass new fields to generation service
    - _Requirements: 3.1_

- [ ] 6. Create Comprehensive AI System Prompts
  - [ ] 6.1 Create prompts/questions.go with experience-level-aware question generation prompt
    - Include jargon avoidance for beginners
    - Include architecture questions for experts
    - Include question ordering rules
    - _Requirements: 1.2, 1.3, 1.4, 3.1, 3.2, 3.3, 3.6_
  - [ ] 6.2 Create prompts/steering.go with complete steering file format specification
    - Include frontmatter rules
    - Include all file type templates (product, tech, structure, security, quality)
    - Include best practices
    - _Requirements: 4.1-4.8, 8.1-8.6_
  - [ ] 6.3 Create prompts/hooks.go with complete hook schema and preset definitions
    - Include all valid when.type and then.type values
    - Include hook examples for each preset
    - Include runCommand restrictions
    - _Requirements: 5.1-5.6, 8.2, 8.5_
  - [ ] 6.4 Create prompts/kickoff.go with kickoff prompt template
    - Include all required sections
    - Include experience-level language adaptation
    - _Requirements: 6.1-6.6, 8.7_
  - [ ] 6.5 Create prompts/agents.go with AGENTS.md template
    - _Requirements: 10.5_

- [ ] 7. Update Generation Service with New Prompts
  - [ ] 7.1 Refactor GenerateQuestions to use experience-level-aware prompts
    - Select prompt variant based on level
    - _Requirements: 1.2, 1.3, 1.4, 3.1_
  - [ ] 7.2 Refactor GenerateOutputs to use comprehensive file format prompts
    - Include all steering files
    - Include hooks based on preset
    - Include AGENTS.md
    - _Requirements: 4.1-4.8, 5.1-5.7, 10.1-10.5_
  - [ ] 7.3 Add output validation for steering file frontmatter
    - Validate inclusion mode
    - Validate fileMatchPattern for fileMatch mode
    - _Requirements: 8.8_
  - [ ] 7.4 Add output validation for hook file schema
    - Validate required fields
    - Validate when.type and then.type values
    - Validate runCommand restrictions
    - _Requirements: 8.8_

- [ ] 8. Checkpoint - Backend Generation Complete
  - Test question generation for each experience level
  - Test output generation with all presets
  - Verify all file types are generated correctly

- [ ] 9. Add Retry Logic for Invalid AI Responses
  - [ ] 9.1 Implement single retry on validation failure in GenerateOutputs
    - _Requirements: 9.3_
  - [ ] 9.2 Add better error messages for validation failures
    - _Requirements: 9.2_

- [ ] 10. Update Frontend API Client
  - [ ] 10.1 Update generateQuestions to send experienceLevel
    - _Requirements: 1.2, 1.3, 1.4_
  - [ ] 10.2 Update generateOutputs to send experienceLevel and hookPreset
    - _Requirements: 5.7, 6.6_
  - [ ] 10.3 Add 'agents' file type to GeneratedFile type
    - _Requirements: 10.5_

- [ ] 11. Update OutputEditor for New File Types
  - [ ] 11.1 Add 'Agents' tab for AGENTS.md file
    - _Requirements: 10.5_
  - [ ] 11.2 Ensure ZIP download preserves correct directory structure
    - _Requirements: 7.7_

- [ ] 12. Checkpoint - Full Integration Complete
  - Test complete flow from level selection to file download
  - Verify all file types appear in correct tabs
  - Test ZIP download structure

- [ ]* 13. Property Tests for Experience Level Adaptation
  - **Property 1: Experience Level Adaptation**
  - Verify beginner questions avoid jargon terms
  - **Validates: Requirements 1.2, 1.3, 1.4, 3.1, 6.6**

- [ ]* 14. Property Tests for Steering File Validity
  - **Property 2: Core Steering Files Validity**
  - **Property 3: Conditional Steering Files Pattern**
  - Verify all steering files have valid frontmatter
  - **Validates: Requirements 4.1-4.8, 10.1-10.3**

- [ ]* 15. Property Tests for Hook Schema
  - **Property 4: Hook File Schema Validity**
  - Verify all hooks have required fields and valid values
  - **Validates: Requirements 5.1-5.5, 10.4**

- [ ]* 16. Property Tests for Kickoff Prompt
  - **Property 5: Kickoff Prompt Completeness**
  - Verify all required sections present
  - **Validates: Requirements 6.1, 6.2, 6.4, 6.5**

- [ ]* 17. Property Tests for Question Generation
  - **Property 6: Question Generation Constraints**
  - Verify count 5-10 and logical ordering
  - **Validates: Requirements 3.4, 3.6**

- [ ] 18. Final Verification
  - [ ] 18.1 Manual test: Complete flow as beginner user
    - _Requirements: 1.1-1.6_
  - [ ] 18.2 Manual test: Complete flow as expert user
    - _Requirements: 1.4, 3.3_
  - [ ] 18.3 Verify generated files work in Kiro IDE
    - _Requirements: 4.8, 5.1-5.5_
  - [ ] 18.4 Verify all hook presets generate correct hooks
    - _Requirements: 5.7, 10.6_

## Notes

- Tasks marked with `*` are optional property-based tests
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Focus on AI prompt quality - this is the key to generating useful files
