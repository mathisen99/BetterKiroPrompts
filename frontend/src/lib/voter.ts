/**
 * Voter hash generation for anonymous rating identification.
 * Generates a consistent hash from localStorage + browser fingerprint
 * to prevent duplicate ratings while maintaining user anonymity.
 * 
 * Requirements: 7.4 - Use browser fingerprinting or localStorage to prevent
 * duplicate ratings from the same user.
 */

const VOTER_HASH_KEY = 'bkp_voter_hash'
const VOTER_FINGERPRINT_KEY = 'bkp_voter_fingerprint'

/**
 * Collects browser fingerprint data for consistent identification.
 * Uses non-invasive browser properties that are stable across sessions.
 */
function collectFingerprint(): string {
  const components: string[] = []

  // Screen properties
  if (typeof screen !== 'undefined') {
    components.push(`${screen.width}x${screen.height}`)
    components.push(`${screen.colorDepth}`)
    components.push(`${screen.pixelDepth}`)
  }

  // Timezone
  try {
    components.push(Intl.DateTimeFormat().resolvedOptions().timeZone)
  } catch {
    components.push('unknown-tz')
  }

  // Language
  if (typeof navigator !== 'undefined') {
    components.push(navigator.language || 'unknown-lang')
    components.push(String(navigator.hardwareConcurrency || 0))
    components.push(navigator.platform || 'unknown-platform')
  }

  // Canvas fingerprint (lightweight version)
  try {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    if (ctx) {
      ctx.textBaseline = 'top'
      ctx.font = '14px Arial'
      ctx.fillText('BKP', 2, 2)
      components.push(canvas.toDataURL().slice(-50))
    }
  } catch {
    components.push('no-canvas')
  }

  return components.join('|')
}

/**
 * Simple hash function for generating a consistent hash from a string.
 * Uses a basic djb2-like algorithm for browser compatibility.
 */
function simpleHash(str: string): string {
  let hash = 5381
  for (let i = 0; i < str.length; i++) {
    hash = ((hash << 5) + hash) ^ str.charCodeAt(i)
  }
  // Convert to hex string and ensure positive
  return (hash >>> 0).toString(16).padStart(8, '0')
}

/**
 * Generates a SHA-256 hash using the Web Crypto API.
 * Falls back to simple hash if crypto is unavailable.
 */
async function sha256Hash(str: string): Promise<string> {
  try {
    const encoder = new TextEncoder()
    const data = encoder.encode(str)
    const hashBuffer = await crypto.subtle.digest('SHA-256', data)
    const hashArray = Array.from(new Uint8Array(hashBuffer))
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
  } catch {
    // Fallback to simple hash if crypto API is unavailable
    return simpleHash(str) + simpleHash(str.split('').reverse().join(''))
  }
}

/**
 * Generates a UUID-like random string with fallback for older browsers.
 * Uses crypto.randomUUID if available, otherwise falls back to crypto.getRandomValues.
 */
function generateUUID(): string {
  // Try native randomUUID first (not available in all browsers/contexts)
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }

  // Fallback using crypto.getRandomValues
  if (typeof crypto !== 'undefined' && typeof crypto.getRandomValues === 'function') {
    const bytes = new Uint8Array(16)
    crypto.getRandomValues(bytes)
    // Set version (4) and variant bits
    bytes[6] = (bytes[6] & 0x0f) | 0x40
    bytes[8] = (bytes[8] & 0x3f) | 0x80
    const hex = Array.from(bytes, b => b.toString(16).padStart(2, '0')).join('')
    return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`
  }

  // Last resort fallback using Math.random (less secure but functional)
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = (Math.random() * 16) | 0
    const v = c === 'x' ? r : (r & 0x3) | 0x8
    return v.toString(16)
  })
}

/**
 * Gets or creates a stored random component for the voter hash.
 * This ensures consistency across sessions while adding entropy.
 */
function getStoredRandom(): string {
  try {
    let stored = localStorage.getItem(VOTER_FINGERPRINT_KEY)
    if (!stored) {
      // Generate a random component
      stored = generateUUID()
      localStorage.setItem(VOTER_FINGERPRINT_KEY, stored)
    }
    return stored
  } catch {
    // If localStorage is unavailable, generate a session-only random
    return generateUUID()
  }
}

/**
 * Generates a consistent voter hash combining localStorage ID and browser fingerprint.
 * The hash is:
 * - Consistent across page reloads for the same browser
 * - Different across different browsers/devices
 * - Anonymous (no PII is stored or transmitted)
 * 
 * @returns A 64-character hex string voter hash
 */
export async function generateVoterHash(): Promise<string> {
  // Check for cached hash first
  try {
    const cached = localStorage.getItem(VOTER_HASH_KEY)
    if (cached && cached.length === 64) {
      return cached
    }
  } catch {
    // localStorage unavailable, continue to generate
  }

  // Collect fingerprint and stored random
  const fingerprint = collectFingerprint()
  const storedRandom = getStoredRandom()

  // Combine and hash
  const combined = `${storedRandom}:${fingerprint}`
  const hash = await sha256Hash(combined)

  // Cache the result
  try {
    localStorage.setItem(VOTER_HASH_KEY, hash)
  } catch {
    // Ignore storage errors
  }

  return hash
}

/**
 * Synchronous version that returns cached hash or generates a simple one.
 * Use this when async is not possible (e.g., in component initialization).
 * 
 * @returns A voter hash string (may be shorter if generated synchronously)
 */
export function getVoterHashSync(): string {
  try {
    const cached = localStorage.getItem(VOTER_HASH_KEY)
    if (cached) {
      return cached
    }
  } catch {
    // localStorage unavailable
  }

  // Generate a simple hash synchronously
  const fingerprint = collectFingerprint()
  const storedRandom = getStoredRandom()
  const combined = `${storedRandom}:${fingerprint}`
  
  // Use simple hash for sync version
  const hash = simpleHash(combined) + simpleHash(combined.split('').reverse().join(''))
  
  // Pad to consistent length
  const paddedHash = hash.padEnd(64, '0').slice(0, 64)

  try {
    localStorage.setItem(VOTER_HASH_KEY, paddedHash)
  } catch {
    // Ignore storage errors
  }

  return paddedHash
}

/**
 * Clears the stored voter hash (useful for testing or privacy reset).
 */
export function clearVoterHash(): void {
  try {
    localStorage.removeItem(VOTER_HASH_KEY)
    localStorage.removeItem(VOTER_FINGERPRINT_KEY)
  } catch {
    // Ignore storage errors
  }
}
