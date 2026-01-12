const API_BASE = '/api'

// Experience level type
export type ExperienceLevel = 'beginner' | 'novice' | 'expert'

// New AI-driven generation types
export interface Question {
  id: number
  text: string
  hint?: string
}

export interface GeneratedFile {
  path: string
  content: string
  type: 'kickoff' | 'steering' | 'hook'
}

export interface Answer {
  questionId: number
  answer: string
}

// Request/Response types
export interface GenerateQuestionsRequest {
  projectIdea: string
  experienceLevel: ExperienceLevel
}

export interface GenerateQuestionsResponse {
  questions: Question[]
}

export interface GenerateOutputsRequest {
  projectIdea: string
  answers: Answer[]
  experienceLevel: ExperienceLevel
}

export interface GenerateOutputsResponse {
  files: GeneratedFile[]
}

export interface ErrorResponse {
  error: string
  retryAfter?: number
}

// API functions
export async function generateQuestions(projectIdea: string, experienceLevel: ExperienceLevel): Promise<GenerateQuestionsResponse> {
  const res = await fetch(`${API_BASE}/generate/questions`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ projectIdea, experienceLevel }),
  })
  
  if (!res.ok) {
    const error: ErrorResponse = await res.json().catch(() => ({ error: 'Failed to generate questions' }))
    throw new ApiError(error.error, res.status, error.retryAfter)
  }
  
  return res.json()
}

export async function generateOutputs(projectIdea: string, answers: Answer[], experienceLevel: ExperienceLevel): Promise<GenerateOutputsResponse> {
  const res = await fetch(`${API_BASE}/generate/outputs`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ projectIdea, answers, experienceLevel }),
  })
  
  if (!res.ok) {
    const error: ErrorResponse = await res.json().catch(() => ({ error: 'Failed to generate outputs' }))
    throw new ApiError(error.error, res.status, error.retryAfter)
  }
  
  return res.json()
}

// Custom error class for API errors
export class ApiError extends Error {
  status: number
  retryAfter?: number

  constructor(message: string, status: number, retryAfter?: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.retryAfter = retryAfter
  }
}
