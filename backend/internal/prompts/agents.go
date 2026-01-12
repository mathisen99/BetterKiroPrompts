package prompts

// AgentsTemplate contains the AGENTS.md template for repository root.
const AgentsTemplate = `# AGENTS.md Template

## Purpose
AGENTS.md is placed at the repository root to provide guidelines for AI agents working on the codebase.
It establishes consistent behavior, commit standards, and decision-making principles.

## Template
` + "```markdown" + `
# Agent Guidelines

## Core Principles
1. Always follow steering files in ` + "`.kiro/steering/`" + `
2. Never invent requirements - ask if unclear
3. Prefer small, reviewable changes
4. Update docs when behavior changes
5. Security is non-negotiable

## Before Coding
- [ ] Requirements are clear and documented
- [ ] Checked for existing patterns in codebase
- [ ] Considered security implications
- [ ] Planned for testing
- [ ] Identified affected documentation

## Code Standards

### General
- Write self-documenting code with clear names
- Keep functions small and focused
- Handle errors explicitly, never silently
- Add comments for non-obvious logic only

### Testing
- Write tests before or alongside code
- Test edge cases and error paths
- Don't mock what you don't own
- Integration tests for critical paths

### Security
- Validate all inputs
- Never log sensitive data
- Use parameterized queries
- Follow least privilege principle

## Commit Standards

### Format
` + "```" + `
<type>: <short summary>

<optional body>

<optional footer>
` + "```" + `

### Types
- ` + "`feat:`" + ` New feature
- ` + "`fix:`" + ` Bug fix
- ` + "`docs:`" + ` Documentation only
- ` + "`style:`" + ` Formatting, no code change
- ` + "`refactor:`" + ` Code change that neither fixes nor adds
- ` + "`test:`" + ` Adding or updating tests
- ` + "`chore:`" + ` Maintenance tasks

### Rules
- One concern per commit
- Summary under 50 characters
- Use imperative mood ("Add feature" not "Added feature")
- Reference issues when applicable

### Examples
` + "```" + `
feat: add user authentication endpoint

fix: prevent null pointer in user lookup

docs: update API documentation for auth routes

refactor: extract validation logic to separate module
` + "```" + `

## When Stuck
1. Re-read the requirements and steering files
2. Check for similar patterns in the codebase
3. Ask for clarification rather than guessing
4. Suggest alternatives if blocked
5. Document assumptions if proceeding

## File Organization
- Follow existing project structure
- One concern per file
- Keep related code together
- Avoid deep nesting

## Pull Request Guidelines
- Clear title describing the change
- Link to related issues/specs
- Include testing instructions
- Note any breaking changes
- Request review from appropriate team members
` + "```" + `
`

// AgentsSystemPrompt returns the system prompt for AGENTS.md generation.
func AgentsSystemPrompt() string {
	return `You are generating an AGENTS.md file for a project repository.

## Purpose
AGENTS.md provides guidelines for AI agents (like Kiro, Copilot, Cursor) working on the codebase.
It should be placed at the repository root.

## Requirements
1. Include core principles for agent behavior
2. Define commit standards with types and format
3. Include pre-coding checklist
4. Provide guidance for when agents are stuck
5. Adapt to the project's tech stack and conventions

` + AgentsTemplate
}
