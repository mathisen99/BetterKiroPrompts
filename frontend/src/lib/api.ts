import { logger } from './logger'

const API_BASE = '/api'

// Default timeout for API requests (180 seconds as per Requirements 4.2)
const DEFAULT_TIMEOUT_MS = 180 * 1000

// Experience level type
export type ExperienceLevel = 'beginner' | 'novice' | 'expert'

// Hook preset type
export type HookPreset = 'light' | 'basic' | 'default' | 'strict'

// New AI-driven generation types
export interface Question {
  id: number
  text: string
  hint?: string
  examples: string[] // 3 clickable example answers
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
  generationId?: string // ID of stored generation for gallery link
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
  const startTime = Date.now()
  const method = options.method || 'GET'
  
  // Try up to 2 times (initial + 1 retry)
  for (let attempt = 0; attempt < 2; attempt++) {
    const { controller, timeoutId } = createTimeoutController(timeoutMs)
    
    try {
      const res = await fetch(url, {
        ...options,
        signal: controller.signal,
      })
      
      clearTimeout(timeoutId)
      const duration = Date.now() - startTime
      
      if (!res.ok) {
        const error: ErrorResponse = await res.json().catch(() => ({ error: errorMessage }))
        const apiError = new ApiError(error.error, res.status, error.retryAfter)
        
        // Log the failed API call
        logger.logApiCall(method, url, res.status, duration)
        
        // Only retry on recoverable errors and first attempt
        if (isRecoverableError(res.status) && attempt === 0) {
          lastError = apiError
          continue
        }
        
        throw apiError
      }
      
      // Log successful API call
      logger.logApiCall(method, url, res.status, duration)
      
      return res.json()
    } catch (err) {
      clearTimeout(timeoutId)
      const duration = Date.now() - startTime
      
      // Handle abort/timeout errors
      if (err instanceof DOMException && err.name === 'AbortError') {
        const timeoutError = new ApiError('Request timed out. Please try again.', 504, undefined, true)
        
        // Log timeout
        logger.logApiCall(method, url, 504, duration)
        logger.error(`API timeout: ${method} ${url} after ${duration}ms`, undefined, 'api')
        
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
        logger.error(`Network error: ${method} ${url} - ${err.message}`, undefined, 'api')
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

// Gallery types
export interface GalleryItem {
  id: string
  projectIdea: string
  category: string
  avgRating: number
  ratingCount: number
  viewCount: number
  createdAt: string
  preview: string
}

export interface GalleryListResponse {
  items: GalleryItem[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export interface GalleryDetail {
  id: string
  projectIdea: string
  experienceLevel: string
  hookPreset: string
  files: GeneratedFile[]
  category: string
  avgRating: number
  ratingCount: number
  viewCount: number
  createdAt: string
}

export interface GalleryDetailResponse {
  generation: GalleryDetail
  userRating: number
}

export interface GalleryFilters {
  category?: number
  sortBy: 'newest' | 'highest_rated' | 'most_viewed'
  page: number
  pageSize?: number
}

export interface RateResponse {
  success: boolean
}

// Gallery API functions
export async function listGallery(filters: GalleryFilters): Promise<GalleryListResponse> {
  const params = new URLSearchParams()
  if (filters.category !== undefined) {
    params.set('category', String(filters.category))
  }
  params.set('sort', filters.sortBy)
  params.set('page', String(filters.page))
  if (filters.pageSize) {
    params.set('pageSize', String(filters.pageSize))
  }

  return fetchWithRetry<GalleryListResponse>(
    `${API_BASE}/gallery?${params.toString()}`,
    { method: 'GET' },
    'Failed to load gallery'
  )
}

export async function getGalleryItem(id: string, voterHash?: string): Promise<GalleryDetailResponse> {
  const params = new URLSearchParams()
  if (voterHash) {
    params.set('voterHash', voterHash)
  }
  const queryString = params.toString()
  const url = queryString ? `${API_BASE}/gallery/${id}?${queryString}` : `${API_BASE}/gallery/${id}`

  return fetchWithRetry<GalleryDetailResponse>(
    url,
    { method: 'GET' },
    'Failed to load gallery item'
  )
}

export async function rateGalleryItem(id: string, score: number, voterHash: string): Promise<RateResponse> {
  return fetchWithRetry<RateResponse>(
    `${API_BASE}/gallery/${id}/rate`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ score, voterHash }),
    },
    'Failed to submit rating'
  )
}

// Security Scan types
export type ScanStatus = 'pending' | 'cloning' | 'scanning' | 'reviewing' | 'completed' | 'failed'
export type FindingSeverity = 'critical' | 'high' | 'medium' | 'low' | 'info'

export interface Finding {
  id: string
  severity: FindingSeverity
  tool: string
  file_path: string
  line_number?: number
  description: string
  remediation?: string
  code_example?: string
}

export interface ScanJob {
  id: string
  status: ScanStatus
  repo_url: string
  languages: string[]
  findings: Finding[]
  error?: string
  created_at: string
  completed_at?: string
}

export interface ScanConfig {
  privateRepoEnabled: boolean
  aiReviewEnabled?: boolean
  maxFilesToReview?: number
}

interface ScanConfigResponse {
  private_repo_enabled: boolean
  ai_review_enabled?: boolean
  max_files_to_review?: number
}

// Security Scan API functions
export async function startScan(repoUrl: string): Promise<ScanJob> {
  return fetchWithRetry<ScanJob>(
    `${API_BASE}/scan`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ repo_url: repoUrl }),
    },
    'Failed to start scan'
  )
}

export async function getScanStatus(jobId: string): Promise<ScanJob> {
  return fetchWithRetry<ScanJob>(
    `${API_BASE}/scan/${jobId}`,
    { method: 'GET' },
    'Failed to get scan status'
  )
}

export async function getScanConfig(): Promise<ScanConfig> {
  const response = await fetchWithRetry<ScanConfigResponse>(
    `${API_BASE}/scan/config`,
    { method: 'GET' },
    'Failed to get scan configuration'
  )
  
  // Transform snake_case to camelCase
  return {
    privateRepoEnabled: response.private_repo_enabled,
    aiReviewEnabled: response.ai_review_enabled,
    maxFilesToReview: response.max_files_to_review,
  }
}
