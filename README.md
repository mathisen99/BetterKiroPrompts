# BetterKiroPrompts

[![CI](https://github.com/mathisen99/BetterKiroPrompts/actions/workflows/ci.yml/badge.svg)](https://github.com/mathisen99/BetterKiroPrompts/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go)](https://go.dev/)
[![React Version](https://img.shields.io/badge/React-19.2.3-61DAFB?logo=react)](https://react.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18.1-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-Polyform_NC-blue)](LICENSE)
[![Release](https://img.shields.io/github/v/release/mathisen99/BetterKiroPrompts)](https://github.com/mathisen99/BetterKiroPrompts/releases)

A tool that generates better prompts, steering documents, and Kiro hooks to improve beginner thinking—not write applications for them.

## Why This Exists

This project was designed primarily for beginners. When you're new to Kiro (or AI-assisted development in general), you don't know what you don't know. You might build something that works but has security holes, missing error handling, or architectural problems that will bite you later.

BetterKiroPrompts helps in two ways:

1. **Better starts** — The generated prompts, steering files, and hooks force you to think through your project before coding. What are you building? What are the risks? What should the AI never do? Answering these questions upfront prevents the "I built something but now I'm stuck" problem.

2. **Security awareness** — The security scanner lets you check earlier projects for vulnerabilities. Maybe you built something six months ago without thinking about secrets in code, SQL injection, or dependency vulnerabilities. The scanner catches these issues and explains them in plain language. It's not about shaming past work — it's about learning what to watch for next time.

## How This Differs from Kiro's Plan Agent

Kiro has a built-in [Plan Agent](https://kiro.dev/docs/cli/chat/planning-agent/) that helps transform ideas into implementation plans. So why does this tool exist?

**They solve different problems:**

| Kiro Plan Agent | BetterKiroPrompts |
|-----------------|-------------------|
| Helps you *plan* what to build | Helps you *think* before you plan |
| Outputs a task breakdown | Outputs ready-to-use `.kiro/` config files |
| Assumes developer competence | Adapts to experience level (beginner/novice/expert) |
| Freeform conversation | Structured 12-section kickoff template |
| No guardrails by default | Curated hook presets with security defaults |

**What this tool adds:**

1. **Experience-level adaptation** — Beginners get questions in plain language ("How do users log in?") while experts get technical depth ("What consistency model fits your use case?"). The plan agent doesn't adjust its language.

2. **Ready-to-use Kiro files** — This tool outputs actual `.kiro/steering/*.md` files and `.kiro.hook` files you can drop into your project. The plan agent outputs a plan, not config files.

3. **Opinionated hook presets** — Light, Basic, Default, and Strict presets encode best practices. Beginners don't know they need a secret scan on `agentStop` — this tool gives it to them.

4. **Forbidden jargon for beginners** — The question generator actively avoids terms like "API", "middleware", "OAuth" when talking to beginners. It uses analogies instead ("like a membership card at a store").

5. **Community gallery** — Browse what others generated, learn from their project structures, rate and share configurations.

**Think of it this way:** The plan agent helps you plan *what* to build. This tool helps you configure *how* Kiro should help you build it — with guardrails appropriate to your skill level.

## What It Does

BetterKiroPrompts helps developers create high-quality Kiro configurations through an AI-driven workflow:

1. **Enter your project idea** — Describe what you want to build
2. **Answer contextual questions** — AI generates questions tailored to your experience level
3. **Get tailored outputs** — Receive kickoff prompts, steering files, hooks, and AGENTS.md

The generated files help Kiro understand your project context, enforce coding standards, and guide you through development without writing code for you.

### Target Audience

- **Beginners** learning to use Kiro effectively and wanting to avoid common pitfalls
- **Teams** wanting consistent project scaffolding
- **Self-hosters** running their own instance for private use

## Features

### AI-Driven Generation
Enter your project idea, select your experience level (beginner/novice/expert), and answer AI-generated questions. The system produces:
- **Kickoff Prompt** — Structured project context for Kiro
- **Steering Files** — Rules and guidelines in `.kiro/steering/`
- **Hooks** — Automated actions in `.kiro/hooks/`
- **AGENTS.md** — Agent behavior documentation

### Public Gallery
Browse community-generated outputs for inspiration. Rate and filter by category, popularity, or recency.

### Security Scanning
Scan GitHub repositories for vulnerabilities using multiple security tools (Trivy, Semgrep, TruffleHog, Gitleaks) plus AI-powered code review.

## Prerequisites

| Tool | Version | Required |
|------|---------|----------|
| Docker Engine | 29.1.3+ | Yes |
| Docker Compose | 5.0.1+ | Yes |
| Go | 1.25.5+ | For local development only |
| Node.js | 24.12.0+ | For local development only |
| pnpm | Latest | For local development only |
| PostgreSQL | 18.1+ | Included in Docker Compose |

## Quick Start

```bash
# Clone the repository
git clone https://github.com/mathisen99/BetterKiroPrompts.git
cd BetterKiroPrompts

# Configure secrets
cp .env.example .env
# Edit .env and set your OPENAI_API_KEY

# (Optional) Customize settings
cp config.example.toml config.toml

# Start all services
./build.sh
```

Open http://localhost:8090 in your browser.

### Production Deployment

The default `./build.sh` builds and runs production. Use `--restart` to restart without rebuilding:

```bash
./build.sh --restart
```

## Project Structure

```
/
├── backend/                    # Go backend
│   ├── cmd/server/main.go      # Application entry point
│   ├── internal/
│   │   ├── api/                # HTTP handlers and middleware
│   │   ├── config/             # Configuration loading
│   │   ├── db/                 # Database connection and migrations
│   │   ├── gallery/            # Gallery service
│   │   ├── generation/         # AI generation service
│   │   ├── logger/             # Structured logging
│   │   ├── openai/             # OpenAI API client
│   │   ├── prompts/            # AI prompt templates
│   │   ├── queue/              # Request queuing
│   │   ├── ratelimit/          # Rate limiting
│   │   ├── sanitize/           # Input sanitization
│   │   ├── scanner/            # Security scanning
│   │   └── storage/            # Data persistence
│   └── migrations/             # Database migrations
├── frontend/                   # React frontend
│   ├── src/
│   │   ├── components/         # React components
│   │   ├── pages/              # Page components
│   │   └── lib/                # API client, utilities
│   └── package.json
├── docs/                       # Documentation
│   ├── api.md                  # API reference
│   ├── developer.md            # Architecture and customization
│   ├── self-hosting.md         # Deployment guide
│   └── user-guide.md           # User documentation
├── .kiro/                      # Kiro configuration
│   ├── steering/               # Project steering files
│   ├── specs/                  # Feature specifications
│   └── hooks/                  # Kiro hooks
├── config.example.toml         # Configuration template
├── .env.example                # Environment variables template
├── docker-compose.yml          # Development stack
├── docker-compose.prod.yml     # Production stack
└── build.sh                    # Build and run script
```

## Configuration

BetterKiroPrompts uses a two-file configuration approach:

| File | Purpose | Contains |
|------|---------|----------|
| `.env` | Secrets | API keys, tokens, database credentials |
| `config.toml` | Settings | Ports, timeouts, limits, feature toggles |

### Quick Reference

**Required in `.env`:**
```bash
OPENAI_API_KEY=sk-your-api-key-here
DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable
POSTGRES_USER=bkp_user
POSTGRES_PASSWORD=CHANGE_ME_IN_PRODUCTION
POSTGRES_DB=betterkiro
```

**Common `config.toml` settings:**
```toml
[server]
port = 8090

[openai]
model = "gpt-5.2"            # Current standard model
timeout = "180s"

[rate_limit]
generation_limit_per_hour = 10

[logging]
level = "INFO"               # DEBUG, INFO, WARN, ERROR
```

See [docs/self-hosting.md](docs/self-hosting.md) for complete configuration reference.

## Development Commands

```bash
./build.sh           # Stop, rebuild everything, and start
./build.sh --restart # Restart without rebuilding
./build.sh --stop    # Stop all containers
```

## Documentation

| Document | Description |
|----------|-------------|
| [Development Log](docs/DEVLOG.md) | Project journey, decisions, and lessons learned |
| [API Reference](docs/api.md) | Complete API endpoint documentation |
| [Developer Guide](docs/developer.md) | Architecture, packages, and customization |
| [Self-Hosting Guide](docs/self-hosting.md) | Deployment and configuration |
| [User Guide](docs/user-guide.md) | How to use the generated outputs |

## Troubleshooting

### "OPENAI_API_KEY not set"
Ensure `.env` exists and contains your API key:
```bash
cat .env | grep OPENAI_API_KEY
./build.sh  # Rebuild and restart to pick up changes
```

### "Database connection refused"
Check PostgreSQL is running:
```bash
docker ps | grep postgres
docker logs better-kiro-prompts-postgres-1
```

### "Rate limit exceeded"
Wait for the rate limit window to reset (1 hour), or increase limits in `config.toml`:
```toml
[rate_limit]
generation_limit_per_hour = 20
```

### "AI generation timeout"
Increase the timeout in `config.toml`:
```toml
[openai]
timeout = "300s"
```

### "Repository too large" (Scanner)
Increase the size limit in `config.toml`:
```toml
[scanner]
max_repo_size_mb = 1000
```

### Debug Mode
Enable debug logging for more detail:
```toml
[logging]
level = "DEBUG"
```

Or via environment variable:
```bash
LOG_LEVEL=DEBUG ./build.sh
```

## Tech Stack

- **Backend:** Go 1.25.5
- **Frontend:** React 19.2.3 + Vite 7.3.1 + shadcn/ui 0.9.5
- **Database:** PostgreSQL 18.1
- **Containerization:** Docker Compose 5.0.1

## License

See [LICENSE](LICENSE) for details.

## Contact

- **GitHub**: [mathisen99/BetterKiroPrompts](https://github.com/mathisen99/BetterKiroPrompts)
- **Creator**: Tommy Mathisen
- **Discord**: Mathisen
- **IRC**: Mathisen (Libera.Chat)
- **Email**: tommy.mathisen@aland.net

## Disclaimer

This is a community project, not affiliated with or endorsed by AWS or Kiro. Visit [kiro.dev](https://kiro.dev) for the official Kiro IDE.
