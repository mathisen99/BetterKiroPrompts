package prompts

// KickoffTemplate contains the complete kickoff prompt template with all required sections.
const KickoffTemplate = `# Kickoff Prompt Template

## Purpose
The kickoff prompt enforces "thinking before coding" by requiring all key decisions to be documented before implementation begins.

## Required Sections
Every kickoff prompt MUST include these sections:

1. **Project Identity** - One sentence description
2. **Success Criteria** - Measurable outcomes
3. **Users & Roles** - Who uses this and what can they do
4. **Data Sensitivity** - What data is stored, sensitivity labels
5. **Data Lifecycle** - Retention, deletion, export, audit, backups
6. **Auth Model** - Authentication approach
7. **Concurrency Expectations** - Multi-user, background jobs, shared state
8. **Risks & Tradeoffs** - Top 3 risks with mitigations
9. **Boundaries** - Public vs private data boundaries
10. **Boundary Examples** - Concrete CAN/CANNOT statements
11. **Non-Goals** - What will NOT be built
12. **Constraints** - Time, simplicity, tech limits

## Template Structure
` + "```markdown" + `
# Project Kickoff: {Project Name}

> ⚠️ **IMPORTANT**: Do not write any code until all questions below are answered and reviewed.

## Project Identity
{One sentence description of what this project is and the problem it solves}

## Success Criteria
{What does "done" look like? List 3-5 measurable outcomes}
- [ ] {Criterion 1}
- [ ] {Criterion 2}
- [ ] {Criterion 3}

## Users & Roles
{Who uses this system? Define each role and their capabilities}

| Role | Description | Key Capabilities |
|------|-------------|------------------|
| {Role 1} | {Description} | {What they can do} |
| {Role 2} | {Description} | {What they can do} |

## Data Sensitivity
{What data is stored? Label each type with sensitivity level}

| Data Type | Sensitivity | Storage | Notes |
|-----------|-------------|---------|-------|
| {Type 1} | {Public/Internal/Confidential/Restricted} | {Where stored} | {Notes} |
| {Type 2} | {Sensitivity} | {Where stored} | {Notes} |

### Data Lifecycle
- **Retention**: {How long is data kept? Any legal requirements?}
- **Deletion**: {How can users delete their data? Soft vs hard delete?}
- **Export**: {Can users export their data? What format?}
- **Audit**: {What actions are logged? How long are logs kept?}
- **Backups**: {Backup strategy and recovery time objectives}

## Auth Model
{How do users authenticate? Choose one or describe custom approach}
- [ ] None (public access)
- [ ] Basic (username/password)
- [ ] External provider (OAuth: Google, GitHub, etc.)
- [ ] API keys
- [ ] Custom: {describe}

## Concurrency Expectations
{Answer these questions about concurrent access}
- **Multi-user**: {Can multiple users access simultaneously?}
- **Shared state**: {Is there shared state that needs synchronization?}
- **Background jobs**: {Are there async tasks or scheduled jobs?}
- **Real-time**: {Any real-time features (websockets, SSE)?}

## Risks & Tradeoffs
{Identify top 3 risks and how they'll be addressed}

### Risk 1: {Risk Name}
- **Description**: {What could go wrong}
- **Likelihood**: {Low/Medium/High}
- **Impact**: {Low/Medium/High}
- **Mitigation**: {How to reduce risk}
- **Accepted**: {What we're choosing not to handle}

### Risk 2: {Risk Name}
- **Description**: {What could go wrong}
- **Likelihood**: {Low/Medium/High}
- **Impact**: {Low/Medium/High}
- **Mitigation**: {How to reduce risk}
- **Accepted**: {What we're choosing not to handle}

### Risk 3: {Risk Name}
- **Description**: {What could go wrong}
- **Likelihood**: {Low/Medium/High}
- **Impact**: {Low/Medium/High}
- **Mitigation**: {How to reduce risk}
- **Accepted**: {What we're choosing not to handle}

## Boundaries
{Define what is public vs private, who can access what}

### Public
{Data and features accessible without authentication}

### Private
{Data and features requiring authentication}

### Boundary Examples
{Concrete examples of access control - use CAN/CANNOT format}
- {Role} CAN {action} on {resource}
- {Role} CANNOT {action} on {resource}
- {Role} CAN {action} on {resource} IF {condition}
- {Role} CAN {action} on their own {resource} but CANNOT {action} on others' {resource}

## Non-Goals
{Explicitly state what will NOT be built - prevents scope creep}
- NOT building: {feature 1}
- NOT building: {feature 2}
- NOT building: {feature 3}
- Out of scope: {thing}

## Constraints
{Technical, time, or resource constraints}
- **Timeline**: {Deadline or time budget}
- **Simplicity**: {Complexity limits}
- **Tech**: {Required or forbidden technologies}
- **Budget**: {Cost constraints if any}
- **Team**: {Team size or skill constraints}

---

## Next Steps
1. Review this document with stakeholders
2. Create specs for each major feature
3. Begin implementation only after specs are approved
` + "```" + `
`

// KickoffLanguageAdaptation contains guidance for adapting kickoff language to experience levels.
const KickoffLanguageAdaptation = `## Language Adaptation by Experience Level

### For Beginners
- Use simple, everyday language throughout
- Explain technical terms when they must be used
- Provide examples for each section
- Keep sentences short and clear
- Use analogies to familiar concepts
- Add helpful notes explaining why each section matters

Example adaptations:
- Instead of "Auth Model", use "How Users Log In"
- Instead of "Concurrency Expectations", use "Multiple Users at Once"
- Instead of "Data Sensitivity", use "What Information We Store"
- Add notes like: "This section helps you think about..."

### For Novice Users
- Use common technical terms but explain advanced concepts
- Balance between simple and technical language
- Include brief explanations of trade-offs
- Provide hints for common choices

### For Expert Users
- Use precise technical terminology
- Assume familiarity with architecture patterns
- Focus on edge cases and trade-offs
- Include technical considerations like:
  - CAP theorem implications
  - Consistency models
  - Failure modes
  - Observability requirements
`

// KickoffSystemPrompt returns the complete system prompt for kickoff prompt generation.
func KickoffSystemPrompt(experienceLevel string) string {
	basePrompt := `You are generating a project kickoff prompt that enforces "thinking before coding".

## Critical Requirements
1. The kickoff prompt MUST contain the phrase "Do not write any code until" or equivalent
2. ALL required sections must be present (see template)
3. Language complexity must match the user's experience level
4. Be specific and actionable - avoid vague placeholders
5. Fill in sections based on the project idea and user's answers

` + KickoffTemplate + "\n\n" + KickoffLanguageAdaptation

	levelNote := getKickoffLevelNote(experienceLevel)
	return basePrompt + "\n\n" + levelNote
}

func getKickoffLevelNote(level string) string {
	switch level {
	case ExperienceBeginner:
		return `## Current User: Beginner
Adapt ALL language to be beginner-friendly:
- Replace technical jargon with simple terms
- Add explanatory notes for each section
- Use concrete examples throughout
- Keep the structure but simplify the language`

	case ExperienceNovice:
		return `## Current User: Novice
Use moderate technical language:
- Common terms are fine, explain advanced concepts
- Include helpful hints for decisions
- Balance between accessibility and precision`

	case ExperienceExpert:
		return `## Current User: Expert
Use full technical terminology:
- Be precise and concise
- Include technical considerations
- Focus on edge cases and trade-offs
- No need to explain common concepts`

	default:
		return getKickoffLevelNote(ExperienceNovice)
	}
}
