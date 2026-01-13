const API_BASE = '/api'

// Default timeout for API requests (120 seconds as per Requirements 3.1)
const DEFAULT_TIMEOUT_MS = 120 * 1000

// Experience level type
export type ExperienceLevel = 'beginner' | 'novice' | 'expert'

// Hook preset type
export type HookPreset = 'light' | 'basic' | 'default' | 'strict'

// New AI-driven generation types
export interface Question {
  id: number
  text: string
  hint?: string
}

export interface GeneratedFile {
  path: string
  content: string
  type: 'kickoff' | 'steering' | 'hook' | 'agents'
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
  hookPreset: HookPreset
}

export interface GenerateOutputsResponse {
  files: GeneratedFile[]
}

export interface ErrorResponse {
  error: string
  retryAfter?: number
}

// Custom error class for API errors
export class ApiError extends Error {
  status: number
  retryAfter?: number
  isTimeout?: boolean

  constructor(message: string, status: number, retryAfter?: number, isTimeout?: boolean) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.retryAfter = retryAfter
    this.isTimeout = isTimeout
  }
}

// Check if an error is recoverable (should auto-retry)
function isRecoverableError(status: number): boolean {
  return status === 503 || status === 504
}

// Create an AbortController with timeout
function createTimeoutController(timeoutMs: number = DEFAULT_TIMEOUT_MS): { controller: AbortController; timeoutId: ReturnType<typeof setTimeout> } {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => {
    controller.abort()
  }, timeoutMs)
  return { controller, timeoutId }
}

// Generic fetch with automatic retry for recoverable errors and timeout support
async function fetchWithRetry<T>(
  url: string,
  options: RequestInit,
  errorMessage: string,
  timeoutMs: number = DEFAULT_TIMEOUT_MS
): Promise<T> {
  let lastError: ApiError | null = null
  
  // Try up to 2 times (initial + 1 retry)
  for (let attempt = 0; attempt < 2; attempt++) {
    const { controller, timeoutId } = createTimeoutController(timeoutMs)
    
    try {
      const res = await fetch(url, {
        ...options,
        signal: controller.signal,
      })
      
      clearTimeout(timeoutId)
      
      if (!res.ok) {
        const error: ErrorResponse = await res.json().catch(() => ({ error: errorMessage }))
        const apiError = new ApiError(error.error, res.status, error.retryAfter)
        
        // Only retry on recoverable errors and first attempt
        if (isRecoverableError(res.status) && attempt === 0) {
          lastError = apiError
          continue
        }
        
        throw apiError
      }
      
      return res.json()
    } catch (err) {
      clearTimeout(timeoutId)
      
      // Handle abort/timeout errors
      if (err instanceof DOMException && err.name === 'AbortError') {
        const timeoutError = new ApiError('Request timed out. Please try again.', 504, undefined, true)
        
        // Retry once on timeout
        if (attempt === 0) {
          lastError = timeoutError
          continue
        }
        
        throw timeoutError
      }
      
      // If it's already an ApiError, handle retry logic
      if (err instanceof ApiError) {
        if (isRecoverableError(err.status) && attempt === 0) {
          lastError = err
          continue
        }
        throw err
      }
      
      // Network errors (TypeError from fetch) - retry once
      if (err instanceof TypeError && attempt === 0) {
        lastError = new ApiError('Network error', 0)
        continue
      }
      
      throw err
    }
  }
  
  // If we exhausted retries, throw the last error
  if (lastError) {
    throw lastError
  }
  
  throw new ApiError(errorMessage, 500)
}

// API functions
export async function generateQuestions(projectIdea: string, experienceLevel: ExperienceLevel): Promise<GenerateQuestionsResponse> {
  return fetchWithRetry<GenerateQuestionsResponse>(
    `${API_BASE}/generate/questions`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ projectIdea, experienceLevel }),
    },
    'Failed to generate questions'
  )
}

export async function generateOutputs(projectIdea: string, answers: Answer[], experienceLevel: ExperienceLevel, hookPreset: HookPreset): Promise<GenerateOutputsResponse> {
  return fetchWithRetry<GenerateOutputsResponse>(
    `${API_BASE}/generate/outputs`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ projectIdea, answers, experienceLevel, hookPreset }),
    },
    'Failed to generate outputs'
  )
}
