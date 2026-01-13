import { Prism as PrismHighlighter } from 'react-syntax-highlighter'
import type { CSSProperties } from 'react'
import type { SupportedLanguage } from '@/lib/syntax'

interface SyntaxHighlighterProps {
  code: string
  language: SupportedLanguage
  className?: string
}

// Custom dark theme matching app colors (blue-based dark theme)
const customDarkTheme: { [key: string]: CSSProperties } = {
  'code[class*="language-"]': {
    color: '#e2e8f0', // foreground
    background: 'none',
    fontFamily: 'ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace',
    fontSize: '0.875rem',
    textAlign: 'left',
    whiteSpace: 'pre',
    wordSpacing: 'normal',
    wordBreak: 'normal',
    wordWrap: 'normal',
    lineHeight: '1.6',
    tabSize: 2,
    hyphens: 'none',
  },
  'pre[class*="language-"]': {
    color: '#e2e8f0',
    background: 'oklch(0.18 0 0)', // slightly lighter than card background
    fontFamily: 'ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace',
    fontSize: '0.875rem',
    textAlign: 'left',
    whiteSpace: 'pre',
    wordSpacing: 'normal',
    wordBreak: 'normal',
    wordWrap: 'normal',
    lineHeight: '1.6',
    tabSize: 2,
    hyphens: 'none',
    padding: '1rem',
    margin: '0',
    overflow: 'auto',
    borderRadius: '0.5rem',
  },
  // Comments
  comment: { color: '#64748b' },
  prolog: { color: '#64748b' },
  doctype: { color: '#64748b' },
  cdata: { color: '#64748b' },
  // Punctuation
  punctuation: { color: '#94a3b8' },
  // Properties, tags, symbols
  property: { color: '#7dd3fc' }, // sky-300 - for JSON keys
  tag: { color: '#7dd3fc' },
  symbol: { color: '#7dd3fc' },
  deleted: { color: '#f87171' },
  // Strings
  string: { color: '#86efac' }, // green-300
  char: { color: '#86efac' },
  'attr-value': { color: '#86efac' },
  inserted: { color: '#86efac' },
  // Numbers, booleans
  number: { color: '#fbbf24' }, // amber-400
  boolean: { color: '#c4b5fd' }, // violet-300
  constant: { color: '#c4b5fd' },
  // Keywords
  keyword: { color: '#c4b5fd' },
  // Functions
  function: { color: '#60a5fa' }, // blue-400 (primary)
  'class-name': { color: '#60a5fa' },
  // Operators
  operator: { color: '#f472b6' }, // pink-400
  entity: { color: '#f472b6' },
  url: { color: '#38bdf8' }, // sky-400
  // Regex
  regex: { color: '#fbbf24' },
  important: { color: '#fbbf24', fontWeight: 'bold' },
  // Variables
  variable: { color: '#e2e8f0' },
  // Markdown specific
  title: { color: '#60a5fa', fontWeight: 'bold' },
  bold: { fontWeight: 'bold' },
  italic: { fontStyle: 'italic' },
  // YAML specific
  atrule: { color: '#c4b5fd' },
  'attr-name': { color: '#7dd3fc' },
  selector: { color: '#86efac' },
}

/**
 * SyntaxHighlighter component for rendering code with syntax highlighting.
 * Supports JSON, Markdown, and YAML languages with a dark theme matching the app.
 */
export function SyntaxHighlighter({ code, language, className = '' }: SyntaxHighlighterProps) {
  // Handle empty or invalid input gracefully
  const safeCode = typeof code === 'string' ? code : String(code ?? '')

  return (
    <PrismHighlighter
      language={language}
      style={customDarkTheme}
      className={className}
      customStyle={{
        margin: 0,
        borderRadius: '0.5rem',
      }}
    >
      {safeCode}
    </PrismHighlighter>
  )
}

export default SyntaxHighlighter
