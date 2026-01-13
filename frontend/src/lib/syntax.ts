export type SupportedLanguage = 'json' | 'markdown' | 'yaml'

/**
 * Detect language from file path for syntax highlighting
 */
export function detectLanguage(path: string): SupportedLanguage {
  const lower = path.toLowerCase()
  if (lower.endsWith('.json')) return 'json'
  if (lower.endsWith('.yaml') || lower.endsWith('.yml')) return 'yaml'
  return 'markdown' // default for .md and unknown
}
