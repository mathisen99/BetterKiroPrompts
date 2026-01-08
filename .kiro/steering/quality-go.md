---
inclusion: fileMatch
fileMatchPattern: "**/*.go"
---

# Quality — Go

## Formatting

- Run `go fmt` before commit

## Linting

- Run `go vet`
- Run `golangci-lint run`

## Testing

- Table-driven tests preferred
- Test public API, not internals
- Name: `TestFunctionName_Scenario`

## Documentation

- Exported functions need doc comments
- Package comment in one file per package
