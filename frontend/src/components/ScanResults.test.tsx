import { describe, it, expect, afterEach } from 'vitest'
import * as fc from 'fast-check'
import { render, screen, cleanup } from '@testing-library/react'
import { ScanResults } from './ScanResults'
import type { ScanJob, Finding, FindingSeverity, ScanStatus } from '@/lib/api'

/**
 * Property 11: Finding Display Completeness
 * For any finding displayed in the UI:
 * - The file path SHALL be visible
 * - The description SHALL be visible
 * - The tool source SHALL be visible
 * - If line_number is present, it SHALL be visible
 * - If remediation is present, it SHALL be visible
 * 
 * Validates: Requirements 10.2, 10.3, 10.4
 * 
 * Feature: info-and-security-scan, Property 11: Finding Display Completeness
 */

// Cleanup after each test
afterEach(() => {
  cleanup()
})

// Arbitrary generators for findings - using more realistic values
const severityArb: fc.Arbitrary<FindingSeverity> = fc.constantFrom(
  'critical', 'high', 'medium', 'low', 'info'
)

const toolNameArb = fc.constantFrom(
  'trivy', 'semgrep', 'trufflehog', 'gitleaks', 'govulncheck', 
  'bandit', 'npm-audit', 'cargo-audit'
)

// Generate unique file paths with counter to avoid duplicates
let fileCounter = 0
const filePathArb = fc.nat({ max: 1000 }).map(n => {
  fileCounter++
  return `src/file_${fileCounter}_${n}.ts`
})

// Generate realistic descriptions (alphanumeric with spaces)
const descriptionArb = fc.stringMatching(/^[A-Za-z][A-Za-z0-9 ]{10,50}$/)

const lineNumberArb = fc.option(
  fc.integer({ min: 1, max: 10000 }),
  { nil: undefined }
)

// Generate realistic remediation text
const remediationArb = fc.option(
  fc.stringMatching(/^[A-Za-z][A-Za-z0-9 .,]{15,100}$/),
  { nil: undefined }
)

const codeExampleArb = fc.option(
  fc.constant('// Example fix\nconst x = 1;'),
  { nil: undefined }
)

const findingArb: fc.Arbitrary<Finding> = fc.record({
  id: fc.uuid(),
  severity: severityArb,
  tool: toolNameArb,
  file_path: filePathArb,
  line_number: lineNumberArb,
  description: descriptionArb,
  remediation: remediationArb,
  code_example: codeExampleArb,
})

// Generate 1-3 findings with unique IDs and file paths
const findingsArb = fc.array(findingArb, { minLength: 1, maxLength: 3 })
  .map(findings => {
    const seenIds = new Set<string>()
    const seenPaths = new Set<string>()
    return findings.filter(f => {
      if (seenIds.has(f.id) || seenPaths.has(f.file_path)) return false
      seenIds.add(f.id)
      seenPaths.add(f.file_path)
      return true
    })
  })
  .filter(findings => findings.length > 0)

// Create a completed scan job with findings
function createCompletedJob(findings: Finding[]): ScanJob {
  return {
    id: 'test-job-id',
    status: 'completed' as ScanStatus,
    repo_url: 'https://github.com/test/repo',
    languages: ['typescript', 'javascript'],
    findings,
    created_at: new Date().toISOString(),
    completed_at: new Date().toISOString(),
  }
}

describe('Property 11: Finding Display Completeness', () => {
  it('file path is visible for all findings', () => {
    fc.assert(
      fc.property(
        findingsArb,
        (findings) => {
          cleanup() // Clean up before each iteration
          const job = createCompletedJob(findings)
          render(<ScanResults job={job} />)
          
          for (const finding of findings) {
            // File path should be visible in the document (may appear multiple times)
            const elements = screen.getAllByText(finding.file_path, { exact: false })
            expect(elements.length).toBeGreaterThan(0)
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('description is visible for all findings', () => {
    fc.assert(
      fc.property(
        findingsArb,
        (findings) => {
          cleanup()
          const job = createCompletedJob(findings)
          const { container } = render(<ScanResults job={job} />)
          
          for (const finding of findings) {
            // Check that description appears in the rendered HTML
            expect(container.innerHTML).toContain(finding.description)
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('tool source is visible for all findings', () => {
    fc.assert(
      fc.property(
        findingsArb,
        (findings) => {
          cleanup()
          const job = createCompletedJob(findings)
          render(<ScanResults job={job} />)
          
          for (const finding of findings) {
            const elements = screen.getAllByText(finding.tool)
            expect(elements.length).toBeGreaterThan(0)
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('line number is visible when present', () => {
    fc.assert(
      fc.property(
        findingsArb.filter(findings => findings.some(f => f.line_number !== undefined)),
        (findings) => {
          cleanup()
          const job = createCompletedJob(findings)
          const { container } = render(<ScanResults job={job} />)
          
          for (const finding of findings) {
            if (finding.line_number !== undefined) {
              // Line number should appear as :lineNumber in the HTML
              const html = container.innerHTML
              expect(html).toContain(`:${finding.line_number}`)
            }
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('remediation is visible when present', () => {
    fc.assert(
      fc.property(
        findingsArb.filter(findings => findings.some(f => f.remediation !== undefined)),
        (findings) => {
          cleanup()
          const job = createCompletedJob(findings)
          const { container } = render(<ScanResults job={job} />)
          
          for (const finding of findings) {
            if (finding.remediation) {
              // Check that remediation text appears in the rendered HTML
              expect(container.innerHTML).toContain(finding.remediation)
            }
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('severity badge is visible for all findings', () => {
    fc.assert(
      fc.property(
        findingsArb,
        (findings) => {
          cleanup()
          const job = createCompletedJob(findings)
          render(<ScanResults job={job} />)
          
          // Each severity should have at least one badge visible
          const severityCounts = new Map<string, number>()
          for (const finding of findings) {
            severityCounts.set(finding.severity, (severityCounts.get(finding.severity) || 0) + 1)
          }
          
          for (const [severity, count] of severityCounts) {
            // Capitalize first letter for badge label
            const label = severity.charAt(0).toUpperCase() + severity.slice(1)
            const badges = screen.getAllByText(label)
            // At least one badge per finding with this severity (may have more due to summary)
            expect(badges.length).toBeGreaterThanOrEqual(count)
          }
        }
      ),
      { numRuns: 100 }
    )
  })
})

describe('ScanResults edge cases', () => {
  it('displays success message when no findings', () => {
    const job = createCompletedJob([])
    render(<ScanResults job={job} />)
    
    expect(screen.getByText('No Security Issues Found')).toBeInTheDocument()
  })

  it('displays error message when scan failed', () => {
    const job: ScanJob = {
      id: 'test-job-id',
      status: 'failed',
      repo_url: 'https://github.com/test/repo',
      languages: [],
      findings: [],
      error: 'Clone failed: repository not found',
      created_at: new Date().toISOString(),
    }
    render(<ScanResults job={job} />)
    
    expect(screen.getByText('Scan Failed')).toBeInTheDocument()
    expect(screen.getByText('Clone failed: repository not found')).toBeInTheDocument()
  })

  it('displays languages when available', () => {
    const job = createCompletedJob([])
    render(<ScanResults job={job} />)
    
    expect(screen.getByText('typescript')).toBeInTheDocument()
    expect(screen.getByText('javascript')).toBeInTheDocument()
  })
})
