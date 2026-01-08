# Agent Guidelines

## Core Rules

- Always follow steering files in `.kiro/steering/`
- Never invent requirements—ask if unclear
- Prefer small, reviewable changes
- Update docs when behavior changes

## Major Tasks

A major task affects: behavior, API, auth, DB schema, security, or deployment.

Major tasks require:
- Trigger relevant hooks
- Clean, atomic commits
- Documentation/steering updates when behavior changes

## Commit Format

- Prefix: `feat:`, `fix:`, `docs:`, `chore:`
- One concern per commit
- One-sentence summary
