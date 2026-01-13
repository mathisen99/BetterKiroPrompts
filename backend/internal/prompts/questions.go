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

// ForbiddenBeginnerTerms are technical terms that MUST NOT appear in beginner questions.
// These terms are too advanced for users new to programming.
var ForbiddenBeginnerTerms = []string{
	// Architecture terms
	"API", "REST", "GraphQL", "microservices", "monolith", "backend", "frontend",
	// Database terms
	"database schema", "SQL", "NoSQL", "ORM", "migration", "query", "index",
	// Authentication terms
	"authentication flow", "OAuth", "JWT", "token", "session", "middleware",
	// Infrastructure terms
	"CI/CD", "containerization", "Docker", "Kubernetes", "deployment", "DevOps",
	"orchestration", "load balancing", "scaling", "serverless",
	// Advanced patterns
	"microservices", "distributed", "scalability", "concurrency",
	"sharding", "replication", "eventual consistency",
	"CQRS", "event sourcing", "saga pattern", "circuit breaker",
	// Other technical terms
	"endpoint", "payload", "serialization", "deserialization", "webhook",
	"caching", "CDN", "SSL", "TLS", "HTTPS",
}

// JargonTerms is an alias for backward compatibility
var JargonTerms = ForbiddenBeginnerTerms

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
The user is completely new to programming. Use ONLY everyday language they would understand.

### AVOID Technical Jargon - FORBIDDEN TERMS (DO NOT USE):
API, REST, GraphQL, microservices, monolith, backend, frontend, database schema, SQL, NoSQL, ORM, migration, query, index, authentication flow, OAuth, JWT, token, session, middleware, CI/CD, containerization, Docker, Kubernetes, deployment, DevOps, orchestration, load balancing, scaling, serverless, distributed, scalability, concurrency, sharding, replication, eventual consistency, CQRS, event sourcing, saga pattern, circuit breaker, endpoint, payload, serialization, webhook, caching, CDN, SSL, TLS

### Language Translation Guide - Use These Instead:
- Instead of "API" → say "a way for apps to talk to each other"
- Instead of "database" → say "where your app saves information"
- Instead of "authentication" → say "how users log in"
- Instead of "backend/frontend" → say "the behind-the-scenes part / what users see"
- Instead of "deploy" → say "put your app online for others to use"
- Instead of "server" → say "a computer that runs your app"
- Instead of "endpoint" → say "a specific page or action in your app"

### Question Style Rules
- Use everyday words a non-programmer would understand
- Ask about WHAT they want, not HOW to build it
- One simple question at a time
- Use real-world analogies to explain concepts
- Provide detailed hints with concrete examples

### Real-World Analogies to Use:
- Saving data = "like writing in a notebook that remembers everything"
- User accounts = "like having a membership card at a store"
- Different user types = "like how a library has librarians and visitors with different abilities"
- App features = "like buttons and pages you can click on"

### Example Questions for Beginners:
- "What will people be able to do with your app? (For example: share photos like Instagram, chat with friends like WhatsApp, or make lists like a to-do app)"
- "Will people need to sign up with an email and password, or can anyone use it right away without an account?"
- "What information does your app need to remember? (For example: people's names, their posts, their favorite items)"
- "Who will use your app? Just you, your friends, or anyone on the internet?"
- "If your app has different types of users (like teachers and students), what can each type do?"`

	case ExperienceNovice:
		return `## Experience Level: Novice
The user has some programming experience but is not an expert. Use moderate technical language with explanations.

### Allowed Technical Terms (with brief explanations when first used):
- Database (where data is stored permanently)
- API (how different parts of your app communicate)
- Authentication (verifying who users are)
- Frontend/Backend (what users see vs. server-side logic)
- Framework (pre-built tools that help you build faster)
- Hosting (where your app runs online)

### Terms to Still Avoid or Explain:
- Microservices, distributed systems, event sourcing
- CQRS, saga patterns, circuit breakers
- Sharding, replication strategies
- Advanced caching strategies

### Question Style Rules
- Use common technical terms but explain advanced concepts
- Balance between simple and technical language
- Ask about both features AND basic technical choices
- Include questions about data structure and storage
- Consider security basics

### Hint Style
- Provide hints that suggest common approaches
- Include brief explanations of trade-offs
- Reference popular tools/frameworks as examples
- Mention 2-3 options with pros/cons

### Example Questions for Novice:
- "What type of database would work best for your data? (For example: PostgreSQL if your data has clear relationships like users-have-posts, or MongoDB if your data structure might change often)"
- "How should users log in? (Options: email/password you manage, or let them use their Google/GitHub account which is easier to set up)"
- "Where do you want to host your app? (Options: a cloud service like Heroku for simplicity, or AWS/GCP for more control)"
- "Will your app need to work offline, or is it okay if users need internet access?"
- "How important is it that your app can handle many users at once? (Just you and friends, or potentially thousands of people)"`

	case ExperienceExpert:
		return `## Experience Level: Expert
The user is an experienced developer. Use full technical terminology without explanations.

### Language Rules
- Use precise technical terminology freely
- Assume familiarity with architecture patterns, design patterns, and best practices
- No need to explain common concepts
- Discuss trade-offs at a technical level
- Reference specific technologies, protocols, and standards

### Question Focus Areas
- Architecture decisions and patterns (monolith vs microservices, event-driven, CQRS)
- Scalability and performance requirements (horizontal scaling, caching strategies, CDN)
- Data consistency requirements (ACID vs BASE, eventual consistency, conflict resolution)
- Distributed systems concerns (CAP theorem trade-offs, partition tolerance)
- Security model (authentication flows, authorization patterns, data encryption)
- Observability strategy (metrics, distributed tracing, log aggregation)
- CI/CD and deployment strategy (blue-green, canary, feature flags)
- Infrastructure choices (containerization, orchestration, serverless)

### Hint Style
- Reference specific patterns and their trade-offs concisely
- Mention relevant technologies and alternatives
- Keep hints brief - experts don't need hand-holding
- Include links to relevant RFCs or documentation when applicable

### Example Questions for Expert:
- "What consistency model fits your use case? (Strong consistency with performance trade-offs, or eventual consistency for better availability)"
- "How will services communicate? (Sync REST/gRPC for simplicity, async messaging via Kafka/RabbitMQ for decoupling, or event sourcing for audit trails)"
- "What's your observability strategy? (OpenTelemetry for tracing, Prometheus/Grafana for metrics, ELK/Loki for logs)"
- "How will you handle database migrations in production? (Blue-green deployments, backward-compatible migrations, feature flags)"
- "What's your caching strategy? (Redis for session/hot data, CDN for static assets, application-level caching)"
- "How will you handle authentication across services? (JWT with short expiry, OAuth2 with refresh tokens, mTLS for service-to-service)"
- "What's your approach to data partitioning if you need to scale? (Sharding strategy, read replicas, multi-region considerations)"`

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
