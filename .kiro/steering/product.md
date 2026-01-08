---
inclusion: always
---

# Product

## What We Are Building

A tool that generates better prompts, steering documents, and Kiro hooks to improve beginner thinking—not write applications for them.

## What We Are NOT Building

- Application code generator
- Replacement for Kiro specs
- Security guarantee system
- Programming tutorial

## Definition of Done

1. User can generate a full Kiro kickoff prompt
2. Steering files are generated, usable, correctly scoped
3. Hooks are generated, valid, and usable
4. (Optional) Repo scanning works end-to-end

## Absolute Rules

1. Do not generate application code unless explicitly instructed
2. Do not invent requirements—if input is missing, stop and ask
3. Do not skip the prompt flow—execute phases in order
4. Do not overbuild—prefer simplest working solution
5. Do not claim security guarantees—provide safer defaults only
6. Major tasks require: hooks, atomic commits, steering/doc updates
