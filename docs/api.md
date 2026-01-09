# API Documentation

Base URL: `http://localhost:8080/api`

## Endpoints

### GET /health

Health check endpoint.

**Response:**
```json
{"status": "ok"}
```

**Example:**
```bash
curl http://localhost:8080/api/health
```

---

### POST /kickoff/generate

Generate a kickoff prompt from wizard answers.

**Request:**
```json
{
  "answers": {
    "projectIdentity": "string",
    "successCriteria": "string",
    "usersAndRoles": "string",
    "dataSensitivity": "string",
    "dataLifecycle": {
      "retention": "string",
      "deletion": "string",
      "export": "string",
      "auditLogging": "string",
      "backups": "string"
    },
    "authModel": "none" | "basic" | "external",
    "concurrency": "string",
    "risksAndTradeoffs": {
      "topRisks": ["string"],
      "mitigations": ["string"],
      "notHandled": ["string"]
    },
    "boundaries": "string",
    "boundaryExamples": ["string"],
    "nonGoals": "string",
    "constraints": "string"
  }
}
```

**Response:**
```json
{
  "prompt": "string"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/kickoff/generate \
  -H "Content-Type: application/json" \
  -d '{"answers":{"projectIdentity":"A todo app","successCriteria":"Users can manage tasks","usersAndRoles":"Authenticated users","dataSensitivity":"User emails","dataLifecycle":{"retention":"1 year","deletion":"On request","export":"JSON","auditLogging":"None","backups":"Daily"},"authModel":"basic","concurrency":"Single user","risksAndTradeoffs":{"topRisks":["Data loss"],"mitigations":["Backups"],"notHandled":[]},"boundaries":"All private","boundaryExamples":[],"nonGoals":"Mobile app","constraints":"Ship in 2 weeks"}}'
```

---

### POST /steering/generate

Generate steering files from configuration.

**Request:**
```json
{
  "config": {
    "projectName": "string",
    "projectDescription": "string",
    "techStack": {
      "backend": "string",
      "frontend": "string",
      "database": "string"
    },
    "includeConditional": boolean,
    "includeManual": boolean,
    "fileReferences": ["string"],
    "customRules": {}
  }
}
```

**Response:**
```json
{
  "files": [
    {
      "path": ".kiro/steering/product.md",
      "content": "string"
    }
  ]
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/steering/generate \
  -H "Content-Type: application/json" \
  -d '{"config":{"projectName":"MyApp","projectDescription":"A web app","techStack":{"backend":"Go","frontend":"React","database":"PostgreSQL"},"includeConditional":true,"includeManual":false,"fileReferences":[],"customRules":{}}}'
```

---

### POST /hooks/generate

Generate Kiro hook files from preset selection.

**Request:**
```json
{
  "preset": "light" | "basic" | "default" | "strict",
  "techStack": {
    "hasGo": boolean,
    "hasTypeScript": boolean,
    "hasReact": boolean
  }
}
```

**Response:**
```json
{
  "files": [
    {
      "path": ".kiro/hooks/format.kiro.hook",
      "content": "string"
    }
  ]
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/hooks/generate \
  -H "Content-Type: application/json" \
  -d '{"preset":"default","techStack":{"hasGo":true,"hasTypeScript":true,"hasReact":true}}'
```

---

## Error Response

All endpoints return errors in this format:

```json
{
  "error": "Error message"
}
```

HTTP status codes:
- `400` — Bad request (invalid input)
- `500` — Internal server error
