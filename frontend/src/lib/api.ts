const API_BASE = '/api'

// Kickoff types
export interface DataLifecycle {
  retention: string
  deletion: string
  export: string
  auditLogging: string
  backups: string
}

export interface RisksAndTradeoffs {
  topRisks: string[]
  mitigations: string[]
  notHandled: string[]
}

export interface KickoffAnswers {
  projectIdentity: string
  successCriteria: string
  usersAndRoles: string
  dataSensitivity: string
  dataLifecycle: DataLifecycle
  authModel: 'none' | 'basic' | 'external'
  concurrency: string
  risksAndTradeoffs: RisksAndTradeoffs
  boundaries: string
  boundaryExamples: string[]
  nonGoals: string
  constraints: string
}

export interface KickoffResponse {
  prompt: string
}

// Steering types
export interface TechStack {
  backend: string
  frontend: string
  database: string
}

export interface SteeringConfig {
  projectName: string
  projectDescription: string
  techStack: TechStack
  includeConditional: boolean
  customRules: Record<string, string[]>
}

export interface GeneratedFile {
  path: string
  content: string
}

export interface SteeringResponse {
  files: GeneratedFile[]
}

// Hooks types
export interface HooksConfig {
  preset: 'light' | 'basic' | 'default' | 'strict'
  techStack: {
    hasGo: boolean
    hasTypeScript: boolean
    hasReact: boolean
  }
}

export interface HooksResponse {
  files: GeneratedFile[]
}

// API functions
export async function generateKickoff(answers: KickoffAnswers): Promise<KickoffResponse> {
  const res = await fetch(`${API_BASE}/kickoff/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ answers }),
  })
  if (!res.ok) throw new Error('Failed to generate kickoff prompt')
  return res.json()
}

export async function generateSteering(config: SteeringConfig): Promise<SteeringResponse> {
  const res = await fetch(`${API_BASE}/steering/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ config }),
  })
  if (!res.ok) throw new Error('Failed to generate steering files')
  return res.json()
}

export async function generateHooks(config: HooksConfig): Promise<HooksResponse> {
  const res = await fetch(`${API_BASE}/hooks/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(config),
  })
  if (!res.ok) throw new Error('Failed to generate hooks')
  return res.json()
}
