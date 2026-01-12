// Package prompts contains comprehensive AI system prompts for generating
// Kiro project files with experience-level adaptation.
package prompts

import "fmt"

// Experience level constants
const (
	ExperienceBeginner = "beginner"
	ExperienceNovice   = "novice"
	ExperienceExpert   = "expert"
)

// JargonTerms are technical terms to avoid for beginners
var JargonTerms = []string{
	"microservices", "distributed", "scalability", "concurrency", "middleware",
	"orchestration", "containerization", "sharding", "replication", "eventual consistency",
	"CQRS", "event sourcing", "saga pattern", "circuit breaker", "load balancing",
}

// QuestionsSystemPrompt returns the system prompt for question generation
// adapted to the user's experience level.
func QuestionsSystemPrompt(experienceLevel string) string {
	basePrompt := `You are helping a developer plan their project by generating thoughtful follow-up questions.

## Your Role
Generate 5-10 follow-up questions to understand the project requirements better. Questions should help clarify scope, users, data, authentication, tech stack, and constraints.

## Question Ordering Rules (CRITICAL)
Questions MUST follow this logical order:
1. **Identity/Scope** - What is this? What problem does it solve?
2. **Users & Roles** - Who uses it? What can each role do?
3. **Data** - What data is stored? How sensitive is it?
4. **Authentication** - How do users log in? What access control?
5. **Architecture** - How is it structured? What components?
6. **Constraints** - Time limits? Tech requirements? Budget?

## Response Format
Return ONLY valid JSON, no markdown code blocks:
{"questions": [{"id": 1, "text": "...", "hint": "..."}]}

Each question must have:
- id: Sequential number starting from 1
- text: The question itself
- hint: A helpful hint or example answer (optional but recommended)
`

	levelGuidance := getLevelGuidance(experienceLevel)
	return basePrompt + "\n" + levelGuidance
}

func getLevelGuidance(level string) string {
	switch level {
	case ExperienceBeginner:
		return `## Experience Level: Beginner
The user is new to programming. Adapt your questions accordingly:

### Language Rules
- Use simple, everyday language
- Explain any technical terms you must use
- AVOID these jargon terms entirely: microservices, distributed, scalability, concurrency, middleware, orchestration, containerization, sharding, replication, eventual consistency, CQRS, event sourcing, saga pattern, circuit breaker, load balancing
- Instead of "What's your auth strategy?", ask "How will users log in to your app?"
- Instead of "What's your data persistence layer?", ask "Where will your app save information?"

### Question Focus
- Focus on WHAT they want to build, not HOW
- Ask about features users will see and use
- Keep questions about one thing at a time
- Provide detailed hints with concrete examples
- If the project sounds complex, gently suggest starting simpler

### Hint Style
- Provide 2-3 concrete example answers
- Use familiar analogies (e.g., "like a to-do list" or "like Instagram")
- Explain why the question matters

### Example Questions for Beginners
- "What will users be able to do with your app? (e.g., create posts, send messages, track tasks)"
- "Will people need to create an account, or can anyone use it without signing up?"
- "What information will your app need to remember? (e.g., usernames, posts, settings)"`

	case ExperienceNovice:
		return `## Experience Level: Novice
The user has some programming experience. Use moderate technical language:

### Language Rules
- Use common technical terms but explain advanced concepts
- Balance between simple and technical language
- Include helpful hints that bridge basic and advanced concepts
- Can mention frameworks and tools by name

### Question Focus
- Ask about both features and basic technical choices
- Include questions about data structure and storage
- Ask about deployment preferences
- Consider security basics

### Hint Style
- Provide hints that suggest common approaches
- Include brief explanations of trade-offs
- Reference popular tools/frameworks as examples

### Example Questions for Novice
- "What database would you prefer? (e.g., PostgreSQL for relational data, MongoDB for flexible documents)"
- "How will you handle user authentication? (e.g., email/password, OAuth with Google/GitHub)"
- "What's your deployment target? (e.g., cloud hosting, self-hosted, serverless)"`

	case ExperienceExpert:
		return `## Experience Level: Expert
The user is an experienced developer. Use full technical terminology:

### Language Rules
- Use precise technical terminology
- Assume familiarity with architecture patterns
- No need to explain common concepts
- Can discuss trade-offs at a technical level

### Question Focus
- Ask about architecture decisions and patterns
- Include questions about scalability and performance
- Ask about data consistency requirements
- Consider distributed systems concerns if relevant
- Ask about observability and monitoring
- Include questions about CI/CD and deployment strategy

### Hint Style
- Reference specific patterns and their trade-offs
- Mention relevant technologies and alternatives
- Keep hints concise - experts don't need hand-holding

### Example Questions for Expert
- "What consistency model do you need? (strong consistency vs eventual consistency trade-offs)"
- "How will you handle service-to-service communication? (sync REST, async messaging, gRPC)"
- "What's your observability strategy? (metrics, tracing, logging aggregation)"
- "How will you manage database migrations in production?"`

	default:
		// Default to novice if unknown
		return getLevelGuidance(ExperienceNovice)
	}
}

// BuildQuestionsUserPrompt builds the user prompt for question generation.
func BuildQuestionsUserPrompt(projectIdea, experienceLevel string) string {
	levelDesc := getExperienceLevelDescription(experienceLevel)
	return fmt.Sprintf(`Project Idea: %s

User Experience Level: %s (%s)

Generate 5-10 follow-up questions to understand this project better. Remember to:
1. Follow the question ordering rules (identity → users → data → auth → architecture → constraints)
2. Adapt language complexity to the user's experience level
3. Provide helpful hints with each question`, projectIdea, experienceLevel, levelDesc)
}

func getExperienceLevelDescription(level string) string {
	switch level {
	case ExperienceBeginner:
		return "new to programming, needs simple language and guidance"
	case ExperienceNovice:
		return "some experience, understands basic concepts"
	case ExperienceExpert:
		return "experienced developer, comfortable with technical details"
	default:
		return "some experience, understands basic concepts"
	}
}
