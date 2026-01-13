import { describe, it, expect, vi, afterAll } from 'vitest'
import * as fc from 'fast-check'
import type { ExperienceLevel } from './api'

/**
 * Property 15: Automatic Retry Behavior
 * For any recoverable error (network timeout, 503), the system SHALL retry exactly once
 * before surfacing the error to the user.
 *
 * Validates: Requirements 10.4
 *
 * Feature: final-polish, Property 15: Automatic Retry Behavior
 */

// We need to test the retry logic in isolation, so we'll create a testable version
// of the fetchWithRetry function that we can control

// Check if an error is recoverable (should auto-retry)
function isRecoverableError(status: number): boolean {
  return status === 503 || status === 504
}

// Testable fetch with retry logic
async function fetchWithRetryTestable<T>(
  mockFetch: () => Promise<{ ok: boolean; status: number; json: () => Promise<unknown> }>,
  errorMessage: string
): Promise<T> {
  let lastError: Error | null = null

  // Try up to 2 times (initial + 1 retry)
  for (let attempt = 0; attempt < 2; attempt++) {
    try {
      const res = await mockFetch()

      if (!res.ok) {
        const error = await res.json().catch(() => ({ error: errorMessage })) as { error: string }
        const apiError = new Error(error.error)
        ;(apiError as unknown as { status: number }).status = res.status

        // Only retry on recoverable errors and first attempt
        if (isRecoverableError(res.status) && attempt === 0) {
          lastError = apiError
          continue
        }

        throw apiError
      }

      return res.json() as T
    } catch (err) {
      // If it's already an error with status, handle retry logic
      if (err instanceof Error && 'status' in err) {
        const status = (err as unknown as { status: number }).status
        if (isRecoverableError(status) && attempt === 0) {
          lastError = err
          continue
        }
        throw err
      }

      // Network errors (TypeError from fetch) - retry once
      if (err instanceof TypeError && attempt === 0) {
        lastError = new Error('Network error')
        continue
      }

      throw err
    }
  }

  // If we exhausted retries, throw the last error
  if (lastError) {
    throw lastError
  }

  throw new Error(errorMessage)
}

describe('Property 15: Automatic Retry Behavior', () => {
  // Arbitrary generators
  const experienceLevelArb: fc.Arbitrary<ExperienceLevel> = fc.constantFrom(
    'beginner',
    'novice',
    'expert'
  )

  const projectIdeaArb = fc.string({ minLength: 1, maxLength: 200 })

  // Recoverable status codes that should trigger retry
  const recoverableStatusArb = fc.constantFrom(503, 504)

  // Non-recoverable status codes that should NOT trigger retry
  const nonRecoverableStatusArb = fc.constantFrom(400, 401, 403, 404, 429, 500)

  describe('Retry logic', () => {
    it('retries exactly once on recoverable errors then succeeds', async () => {
      await fc.assert(
        fc.asyncProperty(
          projectIdeaArb,
          experienceLevelArb,
          recoverableStatusArb,
          async (_projectIdea, _experienceLevel, status) => {
            // Create a fresh mock for each iteration
            let callCount = 0
            const mockFetch = vi.fn().mockImplementation(() => {
              callCount++
              if (callCount === 1) {
                // First call fails with recoverable error
                return Promise.resolve({
                  ok: false,
                  status,
                  json: () => Promise.resolve({ error: 'Service unavailable' }),
                })
              }
              // Second call succeeds
              return Promise.resolve({
                ok: true,
                json: () => Promise.resolve({ questions: [{ id: 1, text: 'Test?' }] }),
              })
            })

            const result = await fetchWithRetryTestable<{ questions: { id: number; text: string }[] }>(
              mockFetch,
              'Failed to generate questions'
            )

            // Should have called fetch exactly twice
            expect(callCount).toBe(2)
            // Should return successful result
            expect(result.questions).toHaveLength(1)
          }
        ),
        { numRuns: 100 }
      )
    })

    it('retries exactly once on recoverable errors then fails', async () => {
      await fc.assert(
        fc.asyncProperty(
          projectIdeaArb,
          experienceLevelArb,
          recoverableStatusArb,
          async (_projectIdea, _experienceLevel, status) => {
            // Create a fresh mock for each iteration
            let callCount = 0
            const mockFetch = vi.fn().mockImplementation(() => {
              callCount++
              // Both calls fail with recoverable error
              return Promise.resolve({
                ok: false,
                status,
                json: () => Promise.resolve({ error: 'Service unavailable' }),
              })
            })

            await expect(
              fetchWithRetryTestable(mockFetch, 'Failed to generate questions')
            ).rejects.toThrow()

            // Should have called fetch exactly twice (initial + 1 retry)
            expect(callCount).toBe(2)
          }
        ),
        { numRuns: 100 }
      )
    })

    it('does not retry on non-recoverable errors', async () => {
      await fc.assert(
        fc.asyncProperty(
          projectIdeaArb,
          experienceLevelArb,
          nonRecoverableStatusArb,
          async (_projectIdea, _experienceLevel, status) => {
            // Create a fresh mock for each iteration
            let callCount = 0
            const mockFetch = vi.fn().mockImplementation(() => {
              callCount++
              return Promise.resolve({
                ok: false,
                status,
                json: () => Promise.resolve({ error: 'Client error' }),
              })
            })

            await expect(
              fetchWithRetryTestable(mockFetch, 'Failed to generate questions')
            ).rejects.toThrow()

            // Should have called fetch exactly once (no retry)
            expect(callCount).toBe(1)
          }
        ),
        { numRuns: 100 }
      )
    })
  })

  describe('Network errors', () => {
    it('retries exactly once on network errors then succeeds', async () => {
      await fc.assert(
        fc.asyncProperty(projectIdeaArb, experienceLevelArb, async () => {
          // Create a fresh mock for each iteration
          let callCount = 0
          const mockFetch = vi.fn().mockImplementation(() => {
            callCount++
            if (callCount === 1) {
              // First call throws network error
              return Promise.reject(new TypeError('Failed to fetch'))
            }
            // Second call succeeds
            return Promise.resolve({
              ok: true,
              json: () => Promise.resolve({ questions: [{ id: 1, text: 'Test?' }] }),
            })
          })

          const result = await fetchWithRetryTestable<{ questions: { id: number; text: string }[] }>(
            mockFetch,
            'Failed to generate questions'
          )

          // Should have called fetch exactly twice
          expect(callCount).toBe(2)
          // Should return successful result
          expect(result.questions).toHaveLength(1)
        }),
        { numRuns: 100 }
      )
    })
  })

  describe('Integration with actual API functions', () => {
    // Store original fetch
    const originalFetch = globalThis.fetch

    afterAll(() => {
      globalThis.fetch = originalFetch
    })

    it('generateQuestions retries on 503', async () => {
      // Import the actual function
      const { generateQuestions } = await import('./api')

      let callCount = 0
      globalThis.fetch = vi.fn().mockImplementation(() => {
        callCount++
        if (callCount === 1) {
          return Promise.resolve({
            ok: false,
            status: 503,
            json: () => Promise.resolve({ error: 'Service unavailable' }),
          })
        }
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ questions: [{ id: 1, text: 'Test?' }] }),
        })
      }) as typeof fetch

      const result = await generateQuestions('test project', 'beginner')
      expect(callCount).toBe(2)
      expect(result.questions).toHaveLength(1)
    })

    it('generateQuestions does not retry on 400', async () => {
      // Import the actual function
      const { generateQuestions, ApiError } = await import('./api')

      let callCount = 0
      globalThis.fetch = vi.fn().mockImplementation(() => {
        callCount++
        return Promise.resolve({
          ok: false,
          status: 400,
          json: () => Promise.resolve({ error: 'Bad request' }),
        })
      }) as typeof fetch

      await expect(generateQuestions('test project', 'beginner')).rejects.toThrow(ApiError)
      expect(callCount).toBe(1)
    })
  })
})
