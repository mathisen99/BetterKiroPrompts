import type { ExperienceLevel, HookPreset, Question } from './api'

export type Phase = 'welcome' | 'level-select' | 'input' | 'questions' | 'generating' | 'output' | 'error'

export interface SessionState {
  phase: Phase
  experienceLevel: ExperienceLevel | null
  projectIdea: string
  hookPreset: HookPreset
  questions: Question[]
  answers: Record<number, string>
  currentQuestionIndex: number
  savedAt: number // Unix timestamp in milliseconds
}

const STORAGE_KEY = 'bkp_session'
const WELCOME_SEEN_KEY = 'bkp_welcome_seen'
const EXPIRY_MS = 24 * 60 * 60 * 1000 // 24 hours

/**
 * Check if localStorage is available
 */
function isStorageAvailable(): boolean {
  try {
    const test = '__storage_test__'
    localStorage.setItem(test, test)
    localStorage.removeItem(test)
    return true
  } catch {
    return false
  }
}

/**
 * Check if a session state has expired (older than 24 hours)
 */
export function isExpired(state: SessionState): boolean {
  const now = Date.now()
  return now - state.savedAt > EXPIRY_MS
}

/**
 * Save session state to localStorage
 * Silently fails if localStorage is unavailable
 */
export function save(state: Omit<SessionState, 'savedAt'>): void {
  if (!isStorageAvailable()) {
    return
  }

  try {
    const stateWithTimestamp: SessionState = {
      ...state,
      savedAt: Date.now(),
    }
    localStorage.setItem(STORAGE_KEY, JSON.stringify(stateWithTimestamp))
  } catch {
    // Silently fail - localStorage might be full or disabled
  }
}

/**
 * Load session state from localStorage
 * Returns null if no state exists, state is expired, or localStorage is unavailable
 */
export function load(): SessionState | null {
  if (!isStorageAvailable()) {
    return null
  }

  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (!stored) {
      return null
    }

    const state = JSON.parse(stored) as SessionState
    
    // Validate required fields exist
    if (
      typeof state.phase !== 'string' ||
      typeof state.projectIdea !== 'string' ||
      typeof state.hookPreset !== 'string' ||
      typeof state.currentQuestionIndex !== 'number' ||
      typeof state.savedAt !== 'number' ||
      !Array.isArray(state.questions) ||
      typeof state.answers !== 'object'
    ) {
      return null
    }

    // Check expiry
    if (isExpired(state)) {
      clear()
      return null
    }

    return state
  } catch {
    // Invalid JSON or other error
    return null
  }
}

/**
 * Clear saved session state from localStorage
 * Silently fails if localStorage is unavailable
 */
export function clear(): void {
  if (!isStorageAvailable()) {
    return
  }

  try {
    localStorage.removeItem(STORAGE_KEY)
  } catch {
    // Silently fail
  }
}

/**
 * Check if there is restorable state (exists and not expired)
 */
export function hasRestorableState(): boolean {
  return load() !== null
}

/**
 * Check if user has seen the welcome screen before
 */
export function hasSeenWelcome(): boolean {
  if (!isStorageAvailable()) {
    return false
  }
  try {
    return localStorage.getItem(WELCOME_SEEN_KEY) === 'true'
  } catch {
    return false
  }
}

/**
 * Mark that user has seen the welcome screen
 */
export function markWelcomeSeen(): void {
  if (!isStorageAvailable()) {
    return
  }
  try {
    localStorage.setItem(WELCOME_SEEN_KEY, 'true')
  } catch {
    // Silently fail
  }
}
