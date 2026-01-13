import { describe, it, expect } from 'vitest'
import * as fc from 'fast-check'
import { render } from '@testing-library/react'
import { SyntaxHighlighter } from './SyntaxHighlighter'
import { detectLanguage, type SupportedLanguage } from '@/lib/syntax'

/**
 * Property 2: Syntax Highlighter Robustness
 * For any string input (including malformed JSON, invalid Markdown, empty strings,
 * and strings with special characters), the Syntax_Highlighter SHALL render
 * without throwing an exception.
 *
 * Validates: Requirements 2.6
 *
 * Feature: final-polish, Property 2: Syntax Highlighter Robustness
 */

describe('Property 2: Syntax Highlighter Robustness', () => {
  const languageArb: fc.Arbitrary<SupportedLanguage> = fc.constantFrom('json', 'markdown', 'yaml')

  it('renders any string input without throwing', () => {
    fc.assert(
      fc.property(fc.string(), languageArb, (code, language) => {
        // Should not throw for any string input
        expect(() => {
          render(<SyntaxHighlighter code={code} language={language} />)
        }).not.toThrow()
      }),
      { numRuns: 100 }
    )
  })

  it('renders malformed JSON without throwing', () => {
    const malformedJsonArb = fc.oneof(
      fc.constant('{'),
      fc.constant('{"key":'),
      fc.constant('{"key": "value"'),
      fc.constant('[1, 2, 3'),
      fc.constant('{"nested": {"broken":}'),
      fc.constant('null null'),
      fc.constant('{"key": undefined}'),
      fc.string().map((s) => `{${s}}`),
    )

    fc.assert(
      fc.property(malformedJsonArb, (code) => {
        expect(() => {
          render(<SyntaxHighlighter code={code} language="json" />)
        }).not.toThrow()
      }),
      { numRuns: 100 }
    )
  })

  it('renders strings with special characters without throwing', () => {
    const specialCharsArb = fc.oneof(
      fc.string().map((s) => `<script>${s}</script>`),
      fc.string().map((s) => `\n\r\t${s}\0`),
      fc.constant('\u0000\u0001\u0002'),
      fc.constant('🎉🚀💻'),
      fc.constant('中文日本語한국어'),
      fc.constant('<div onclick="alert(1)">test</div>'),
      fc.constant('SELECT * FROM users; DROP TABLE users;--'),
      fc.array(fc.integer({ min: 0, max: 0xFFFF }), { minLength: 1, maxLength: 50 })
        .map((codes) => String.fromCharCode(...codes)),
    )

    fc.assert(
      fc.property(specialCharsArb, languageArb, (code, language) => {
        expect(() => {
          render(<SyntaxHighlighter code={code} language={language} />)
        }).not.toThrow()
      }),
      { numRuns: 100 }
    )
  })

  it('renders empty and whitespace strings without throwing', () => {
    const emptyArb = fc.oneof(
      fc.constant(''),
      fc.constant(' '),
      fc.constant('\n'),
      fc.constant('\t'),
      fc.constant('   \n\t   '),
    )

    fc.assert(
      fc.property(emptyArb, languageArb, (code, language) => {
        expect(() => {
          render(<SyntaxHighlighter code={code} language={language} />)
        }).not.toThrow()
      }),
      { numRuns: 100 }
    )
  })

  it('handles non-string inputs gracefully', () => {
    // Test that the component handles edge cases like null/undefined coercion
    const edgeCaseArb = fc.oneof(
      fc.constant(null as unknown as string),
      fc.constant(undefined as unknown as string),
      fc.integer().map((n) => n as unknown as string),
      fc.boolean().map((b) => b as unknown as string),
    )

    fc.assert(
      fc.property(edgeCaseArb, languageArb, (code, language) => {
        expect(() => {
          render(<SyntaxHighlighter code={code} language={language} />)
        }).not.toThrow()
      }),
      { numRuns: 100 }
    )
  })
})

describe('detectLanguage', () => {
  it('detects JSON files', () => {
    expect(detectLanguage('file.json')).toBe('json')
    expect(detectLanguage('path/to/file.JSON')).toBe('json')
    expect(detectLanguage('.kiro/hooks/test.json')).toBe('json')
  })

  it('detects YAML files', () => {
    expect(detectLanguage('file.yaml')).toBe('yaml')
    expect(detectLanguage('file.yml')).toBe('yaml')
    expect(detectLanguage('path/to/config.YAML')).toBe('yaml')
  })

  it('defaults to markdown for .md files', () => {
    expect(detectLanguage('README.md')).toBe('markdown')
    expect(detectLanguage('docs/guide.MD')).toBe('markdown')
  })

  it('defaults to markdown for unknown extensions', () => {
    expect(detectLanguage('file.txt')).toBe('markdown')
    expect(detectLanguage('file')).toBe('markdown')
    expect(detectLanguage('')).toBe('markdown')
  })
})
