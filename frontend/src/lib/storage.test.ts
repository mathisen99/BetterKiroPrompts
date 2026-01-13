import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import * as fc from 'fast-check'
import { save, load, clear, isExpired, type SessionState, type Phase } from './storage'
import type { ExperienceLevel, HookPreset, Question } from './api'

/**
 * Property 1: State Persistence Consistency
 * For any user action that modifies session state (level selection, project idea submission,
 * answer submission), the LocalStorage_Manager SHALL save the updated state, and loading
 * that state SHALL restore the exact same values.
 *
 * Validates: Requirements 4.1, 4.3
 *
 * Feature: final-polish, Property 1: State Persistence Consistency
 */

// Mock localStorage for testing
const createMockStorage = () => {
  let store: Record<string, string> = {}
  return {
    getItem: vi.fn((key: string) => store[key] ?? null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      store = {}
    }),
    get length() {
      return Object.keys(store).length
    },
    key: vi.fn((index: number) => Object.keys(store)[index] ?? null),
  }
}

describe('Property 1: State Persistence Consistency', () => {
  let mockStorage: ReturnType<typeof createMockStorage>

  beforeEach(() => {
    mockStorage = createMockStorage()
    vi.stubGlobal('localStorage', mockStorage)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  // Arbitrary generators
  const phaseArb: fc.Arbitrary<Phase> = fc.constantFrom(
    'level-select',
    'input',
    'questions',
    'generating',
    'output',
    'error'
  )

  const experienceLevelArb: fc.Arbitrary<ExperienceLevel | null> = fc.oneof(
    fc.constant(null),
    fc.constantFrom('beginner', 'novice', 'expert')
  )

  const hookPresetArb: fc.Arbitrary<HookPreset> = fc.constantFrom(
    'light',
    'basic',
    'default',
    'strict'
  )

  const questionArb: fc.Arbitrary<Question> = fc.record({
    id: fc.integer({ min: 1, max: 1000 }),
    text: fc.string({ minLength: 1, maxLength: 200 }),
    hint: fc.option(fc.string({ minLength: 1, maxLength: 100 }), { nil: undefined }),
  })

  const answersArb: fc.Arbitrary<Record<number, string>> = fc
    .array(
      fc.tuple(fc.integer({ min: 1, max: 100 }), fc.string({ minLength: 0, maxLength: 500 })),
      { minLength: 0, maxLength: 10 }
    )
    .map((pairs) => Object.fromEntries(pairs))

  const sessionStateArb: fc.Arbitrary<Omit<SessionState, 'savedAt'>> = fc.record({
    phase: phaseArb,
    experienceLevel: experienceLevelArb,
    projectIdea: fc.string({ minLength: 0, maxLength: 1000 }),
    hookPreset: hookPresetArb,
    questions: fc.array(questionArb, { minLength: 0, maxLength: 10 }),
    answers: answersArb,
    currentQuestionIndex: fc.integer({ min: 0, max: 20 }),
  })

  it('save then load returns equivalent state values', () => {
    fc.assert(
      fc.property(sessionStateArb, (state) => {
        // Clear any previous state
        clear()

        // Save the state
        save(state)

        // Load the state
        const loaded = load()

        // Verify loaded state matches saved state (excluding savedAt timestamp)
        expect(loaded).not.toBeNull()
        expect(loaded!.phase).toBe(state.phase)
        expect(loaded!.experienceLevel).toBe(state.experienceLevel)
        expect(loaded!.projectIdea).toBe(state.projectIdea)
        expect(loaded!.hookPreset).toBe(state.hookPreset)
        expect(loaded!.questions).toEqual(state.questions)
        expect(loaded!.answers).toEqual(state.answers)
        expect(loaded!.currentQuestionIndex).toBe(state.currentQuestionIndex)
      }),
      { numRuns: 100 }
    )
  })

  it('clear removes saved state', () => {
    fc.assert(
      fc.property(sessionStateArb, (state) => {
        // Save state
        save(state)

        // Verify it was saved
        expect(load()).not.toBeNull()

        // Clear
        clear()

        // Verify it's gone
        expect(load()).toBeNull()
      }),
      { numRuns: 100 }
    )
  })

  it('multiple saves preserve only the latest state', () => {
    fc.assert(
      fc.property(
        fc.array(sessionStateArb, { minLength: 2, maxLength: 5 }),
        (states) => {
          // Save multiple states
          for (const state of states) {
            save(state)
          }

          // Load should return the last saved state
          const loaded = load()
          const lastState = states[states.length - 1]

          expect(loaded).not.toBeNull()
          expect(loaded!.phase).toBe(lastState.phase)
          expect(loaded!.experienceLevel).toBe(lastState.experienceLevel)
          expect(loaded!.projectIdea).toBe(lastState.projectIdea)
          expect(loaded!.hookPreset).toBe(lastState.hookPreset)
          expect(loaded!.questions).toEqual(lastState.questions)
          expect(loaded!.answers).toEqual(lastState.answers)
          expect(loaded!.currentQuestionIndex).toBe(lastState.currentQuestionIndex)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('savedAt timestamp is set on save', () => {
    fc.assert(
      fc.property(sessionStateArb, (state) => {
        const beforeSave = Date.now()
        save(state)
        const afterSave = Date.now()

        const loaded = load()

        expect(loaded).not.toBeNull()
        expect(loaded!.savedAt).toBeGreaterThanOrEqual(beforeSave)
        expect(loaded!.savedAt).toBeLessThanOrEqual(afterSave)
      }),
      { numRuns: 100 }
    )
  })
})

describe('isExpired', () => {
  it('returns false for recent state', () => {
    const state: SessionState = {
      phase: 'input',
      experienceLevel: 'beginner',
      projectIdea: 'test',
      hookPreset: 'default',
      questions: [],
      answers: {},
      currentQuestionIndex: 0,
      savedAt: Date.now(),
    }

    expect(isExpired(state)).toBe(false)
  })

  it('returns true for state older than 24 hours', () => {
    const state: SessionState = {
      phase: 'input',
      experienceLevel: 'beginner',
      projectIdea: 'test',
      hookPreset: 'default',
      questions: [],
      answers: {},
      currentQuestionIndex: 0,
      savedAt: Date.now() - 25 * 60 * 60 * 1000, // 25 hours ago
    }

    expect(isExpired(state)).toBe(true)
  })

  it('returns false for state exactly at 24 hour boundary', () => {
    const state: SessionState = {
      phase: 'input',
      experienceLevel: 'beginner',
      projectIdea: 'test',
      hookPreset: 'default',
      questions: [],
      answers: {},
      currentQuestionIndex: 0,
      savedAt: Date.now() - 24 * 60 * 60 * 1000, // exactly 24 hours ago
    }

    // At exactly 24 hours, it should not be expired (boundary is exclusive)
    expect(isExpired(state)).toBe(false)
  })
})

describe('load with invalid data', () => {
  let mockStorage: ReturnType<typeof createMockStorage>

  beforeEach(() => {
    mockStorage = createMockStorage()
    vi.stubGlobal('localStorage', mockStorage)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('returns null for invalid JSON', () => {
    mockStorage.setItem('bkp_session', 'not valid json')
    expect(load()).toBeNull()
  })

  it('returns null for missing required fields', () => {
    mockStorage.setItem('bkp_session', JSON.stringify({ phase: 'input' }))
    expect(load()).toBeNull()
  })

  it('returns null for expired state', () => {
    const expiredState: SessionState = {
      phase: 'input',
      experienceLevel: 'beginner',
      projectIdea: 'test',
      hookPreset: 'default',
      questions: [],
      answers: {},
      currentQuestionIndex: 0,
      savedAt: Date.now() - 25 * 60 * 60 * 1000,
    }
    mockStorage.setItem('bkp_session', JSON.stringify(expiredState))
    expect(load()).toBeNull()
  })
})
