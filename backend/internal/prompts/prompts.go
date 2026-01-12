// Package prompts provides comprehensive AI system prompts for generating
// Kiro project files including kickoff prompts, steering files, hooks, and AGENTS.md.
package prompts

import (
	"encoding/json"
	"fmt"
)

// Answer represents a user's answer to a question (mirrors generation.Answer).
type Answer struct {
	QuestionID int    `json:"questionId"`
	Answer     string `json:"answer"`
}

// GetQuestionsSystemPrompt returns the complete system prompt for question generation.
func GetQuestionsSystemPrompt(experienceLevel string) string {
	return QuestionsSystemPrompt(experienceLevel)
}

// GetQuestionsUserPrompt returns the user prompt for question generation.
func GetQuestionsUserPrompt(projectIdea, experienceLevel string) string {
	return BuildQuestionsUserPrompt(projectIdea, experienceLevel)
}

// GetOutputsSystemPrompt returns the complete system prompt for output generation.
// This combines all the knowledge about steering files, hooks, kickoff prompts, and AGENTS.md.
func GetOutputsSystemPrompt(experienceLevel, hookPreset string) string {
	return fmt.Sprintf(`You are generating Kiro project files for a developer. Based on their project idea and answers, generate a complete set of files.

## Experience Level: %s
Adapt all language and complexity to match this experience level.

## Hook Preset: %s
Generate hooks appropriate for this preset level.

## Files to Generate

### 1. Kickoff Prompt (REQUIRED)
Path: kickoff-prompt.md
Type: kickoff
%s

### 2. Steering Files (REQUIRED)
Generate these steering files with proper frontmatter:

%s

### 3. Hook Files (REQUIRED)
Generate hooks based on the selected preset:

%s

### 4. AGENTS.md (REQUIRED)
Path: AGENTS.md
Type: agents
%s

## Response Format
Return ONLY valid JSON, no markdown code blocks:
{
  "files": [
    {"path": "kickoff-prompt.md", "content": "...", "type": "kickoff"},
    {"path": ".kiro/steering/product.md", "content": "...", "type": "steering"},
    {"path": ".kiro/steering/tech.md", "content": "...", "type": "steering"},
    {"path": ".kiro/steering/structure.md", "content": "...", "type": "steering"},
    {"path": ".kiro/hooks/format-on-stop.kiro.hook", "content": "...", "type": "hook"},
    {"path": "AGENTS.md", "content": "...", "type": "agents"}
  ]
}

## Critical Rules
1. ALL steering files MUST have valid YAML frontmatter with 'inclusion' field
2. Conditional steering files (security, quality) MUST have 'fileMatchPattern'
3. ALL hook files MUST be valid JSON matching the Kiro hook schema
4. Hook files with 'runCommand' can ONLY use 'agentStop' or 'promptSubmit' triggers
5. Kickoff prompt MUST contain "Do not write any code until" or equivalent phrase
6. AGENTS.md MUST include commit standards and core principles
7. Adapt language complexity to the user's experience level throughout`,
		experienceLevel,
		hookPreset,
		KickoffTemplate,
		SteeringFormatSpec+"\n\n"+SteeringTemplates,
		HookSchemaSpec+"\n\n"+getHookPresetGuidance(hookPreset),
		AgentsTemplate,
	)
}

// GetOutputsUserPrompt returns the user prompt for output generation.
func GetOutputsUserPrompt(projectIdea string, answers []Answer, experienceLevel, hookPreset string) string {
	answersJSON, _ := json.Marshal(answers)

	return fmt.Sprintf(`Generate Kiro project files for this project:

## Project Idea
%s

## User's Answers to Questions
%s

## Configuration
- Experience Level: %s
- Hook Preset: %s

Generate all required files:
1. kickoff-prompt.md - Comprehensive kickoff prompt with all required sections
2. .kiro/steering/product.md - Product definition (inclusion: always)
3. .kiro/steering/tech.md - Tech stack and rules (inclusion: always)
4. .kiro/steering/structure.md - Repository structure (inclusion: always)
5. Language-specific security and quality steering files (inclusion: fileMatch)
6. Hook files appropriate for the %s preset
7. AGENTS.md - Agent guidelines for the repository root

Remember to:
- Adapt language to %s experience level
- Use proper frontmatter for all steering files
- Generate valid JSON for all hook files
- Include all required sections in the kickoff prompt`,
		projectIdea,
		string(answersJSON),
		experienceLevel,
		hookPreset,
		hookPreset,
		experienceLevel,
	)
}

func getHookPresetGuidance(preset string) string {
	presetInfo, ok := HookPresetDescriptions[preset]
	if !ok {
		presetInfo = HookPresetDescriptions[HookPresetDefault]
	}

	return fmt.Sprintf(`## Selected Preset: %s
%s

Generate these hooks: %v

Refer to the hook examples above for the correct format for each hook type.`,
		presetInfo.Title,
		presetInfo.Description,
		presetInfo.Hooks,
	)
}

// ValidExperienceLevels returns the list of valid experience levels.
func ValidExperienceLevels() []string {
	return []string{ExperienceBeginner, ExperienceNovice, ExperienceExpert}
}

// ValidHookPresets returns the list of valid hook presets.
func ValidHookPresets() []string {
	return []string{HookPresetLight, HookPresetBasic, HookPresetDefault, HookPresetStrict}
}

// IsValidExperienceLevel checks if the given level is valid.
func IsValidExperienceLevel(level string) bool {
	for _, valid := range ValidExperienceLevels() {
		if level == valid {
			return true
		}
	}
	return false
}

// IsValidHookPreset checks if the given preset is valid.
func IsValidHookPreset(preset string) bool {
	for _, valid := range ValidHookPresets() {
		if preset == valid {
			return true
		}
	}
	return false
}
