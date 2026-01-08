---
inclusion: fileMatch
fileMatchPattern: "**/*.{ts,tsx}"
---

# Security — Web

- No secrets in frontend code
- Validate input before sending to API
- Escape user content in rendering (React handles by default)
- No `dangerouslySetInnerHTML` without explicit sanitization
- Use HTTPS URLs only for external resources
- No inline event handlers with user data
