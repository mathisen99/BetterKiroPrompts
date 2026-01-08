---
inclusion: fileMatch
fileMatchPattern: "**/*.tsx"
---

# Quality — Web

## Formatting & Linting

- Run `pnpm lint` before commit
- Run `pnpm typecheck` for type errors

## Components

- One component per file
- Props interface above component
- Use shadcn/ui primitives, don't reinvent

## Testing

- Test user interactions, not implementation
- Name: `ComponentName.test.tsx`
