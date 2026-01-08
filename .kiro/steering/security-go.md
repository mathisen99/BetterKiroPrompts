---
inclusion: fileMatch
fileMatchPattern: "**/*.go"
---

# Security — Go

- No secrets in code—use environment variables
- Validate all input at API boundaries
- Auth boundaries must be explicit in handlers
- Least privilege for DB connections
- No `sql.Query` with string concatenation—use parameterized queries
- Log security events, never log secrets
