import { useMemo } from 'react'
import { AlertTriangle, AlertCircle, Info, CheckCircle, XCircle, FileCode, Wrench } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { SyntaxHighlighter } from '@/components/SyntaxHighlighter'
import type { ScanJob, Finding, FindingSeverity } from '@/lib/api'

interface ScanResultsProps {
  job: ScanJob
}

// Severity configuration for styling and icons
const severityConfig: Record<FindingSeverity, { 
  icon: typeof AlertTriangle
  color: string
  bgColor: string
  borderColor: string
  label: string
}> = {
  critical: {
    icon: XCircle,
    color: 'text-red-400',
    bgColor: 'bg-red-500/10',
    borderColor: 'border-red-500/30',
    label: 'Critical',
  },
  high: {
    icon: AlertTriangle,
    color: 'text-orange-400',
    bgColor: 'bg-orange-500/10',
    borderColor: 'border-orange-500/30',
    label: 'High',
  },
  medium: {
    icon: AlertCircle,
    color: 'text-yellow-400',
    bgColor: 'bg-yellow-500/10',
    borderColor: 'border-yellow-500/30',
    label: 'Medium',
  },
  low: {
    icon: Info,
    color: 'text-blue-400',
    bgColor: 'bg-blue-500/10',
    borderColor: 'border-blue-500/30',
    label: 'Low',
  },
  info: {
    icon: Info,
    color: 'text-slate-400',
    bgColor: 'bg-slate-500/10',
    borderColor: 'border-slate-500/30',
    label: 'Info',
  },
}

// Severity order for sorting
const severityOrder: FindingSeverity[] = ['critical', 'high', 'medium', 'low', 'info']

function SeverityBadge({ severity }: { severity: FindingSeverity }) {
  const config = severityConfig[severity]
  const Icon = config.icon
  
  return (
    <Badge 
      variant="outline" 
      className={`${config.bgColor} ${config.borderColor} ${config.color} gap-1`}
    >
      <Icon className="h-3 w-3" />
      {config.label}
    </Badge>
  )
}

// Strip the temp scan path prefix from file paths
function cleanFilePath(filePath: string): string {
  // Remove /scan/repos/scan-repo-XXXXXXX/ prefix
  return filePath.replace(/^\/scan\/repos\/scan-repo-\d+\//, '')
}

function FindingCard({ finding }: { finding: Finding }) {
  const config = severityConfig[finding.severity]
  const cleanPath = cleanFilePath(finding.file_path)
  
  return (
    <Card className={`${config.bgColor} ${config.borderColor} border overflow-hidden`}>
      <CardHeader className="pb-2">
        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-2">
          <div className="flex items-start gap-2 min-w-0 flex-1">
            <FileCode className="h-4 w-4 text-muted-foreground shrink-0 mt-0.5" />
            <span className="font-mono text-sm break-all" title={cleanPath}>
              {cleanPath}
              {finding.line_number && (
                <span className="text-muted-foreground">:{finding.line_number}</span>
              )}
            </span>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <Badge variant="secondary" className="text-xs">
              {finding.tool}
            </Badge>
            <SeverityBadge severity={finding.severity} />
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm">{finding.description}</p>
        
        {finding.remediation && (
          <div className="space-y-2">
            <div className="flex items-center gap-2 text-sm font-medium text-primary">
              <Wrench className="h-4 w-4" />
              Remediation
            </div>
            <p className="text-sm text-muted-foreground">{finding.remediation}</p>
          </div>
        )}
        
        {finding.code_example && (
          <div className="space-y-2">
            <div className="text-sm font-medium text-primary">Code Fix</div>
            <SyntaxHighlighter 
              code={finding.code_example} 
              language="markdown"
              className="text-xs"
            />
          </div>
        )}
      </CardContent>
    </Card>
  )
}

export function ScanResults({ job }: ScanResultsProps) {
  // Group findings by severity
  const findingsBySeverity = useMemo(() => {
    const grouped: Record<FindingSeverity, Finding[]> = {
      critical: [],
      high: [],
      medium: [],
      low: [],
      info: [],
    }
    
    for (const finding of job.findings || []) {
      if (grouped[finding.severity]) {
        grouped[finding.severity].push(finding)
      }
    }
    
    return grouped
  }, [job.findings])

  // Count findings by severity
  const severityCounts = useMemo(() => {
    const counts: Record<FindingSeverity, number> = {
      critical: 0,
      high: 0,
      medium: 0,
      low: 0,
      info: 0,
    }
    
    for (const finding of job.findings || []) {
      if (counts[finding.severity] !== undefined) {
        counts[finding.severity]++
      }
    }
    
    return counts
  }, [job.findings])

  const totalFindings = job.findings?.length || 0

  // Handle failed scan
  if (job.status === 'failed') {
    return (
      <Card className="bg-destructive/10 border-destructive/30">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-destructive">
            <XCircle className="h-5 w-5" />
            Scan Failed
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            {job.error || 'An unknown error occurred during the scan.'}
          </p>
        </CardContent>
      </Card>
    )
  }

  // Handle no findings
  if (totalFindings === 0) {
    return (
      <Card className="bg-green-500/10 border-green-500/30">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-green-400">
            <CheckCircle className="h-5 w-5" />
            No Security Issues Found
          </CardTitle>
        </CardHeader>
        <CardContent>
          <CardDescription>
            Great news! The security scan completed successfully and no vulnerabilities 
            were detected in the repository.
          </CardDescription>
          {job.languages && job.languages.length > 0 && (
            <div className="mt-4 flex items-center gap-2 flex-wrap">
              <span className="text-sm text-muted-foreground">Languages scanned:</span>
              {job.languages.map((lang) => (
                <Badge key={lang} variant="secondary" className="text-xs">
                  {lang}
                </Badge>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Summary Card */}
      <Card className="bg-card/50 backdrop-blur">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-primary" />
            Scan Results
          </CardTitle>
          <CardDescription>
            Found {totalFindings} security {totalFindings === 1 ? 'issue' : 'issues'} in{' '}
            <span className="font-mono text-foreground">{job.repo_url}</span>
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Severity Summary */}
          <div className="flex flex-wrap gap-3">
            {severityOrder.map((severity) => {
              const count = severityCounts[severity]
              if (count === 0) return null
              
              const config = severityConfig[severity]
              const Icon = config.icon
              
              return (
                <div 
                  key={severity}
                  className={`flex items-center gap-2 px-3 py-1.5 rounded-md ${config.bgColor} ${config.borderColor} border`}
                >
                  <Icon className={`h-4 w-4 ${config.color}`} />
                  <span className={`text-sm font-medium ${config.color}`}>
                    {count} {config.label}
                  </span>
                </div>
              )
            })}
          </div>
          
          {/* Languages */}
          {job.languages && job.languages.length > 0 && (
            <div className="mt-4 flex items-center gap-2 flex-wrap">
              <span className="text-sm text-muted-foreground">Languages detected:</span>
              {job.languages.map((lang) => (
                <Badge key={lang} variant="secondary" className="text-xs">
                  {lang}
                </Badge>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Findings by Severity */}
      {severityOrder.map((severity) => {
        const findings = findingsBySeverity[severity]
        if (findings.length === 0) return null
        
        const config = severityConfig[severity]
        
        return (
          <div key={severity} className="space-y-3">
            <h3 className={`text-lg font-semibold flex items-center gap-2 ${config.color}`}>
              <config.icon className="h-5 w-5" />
              {config.label} ({findings.length})
            </h3>
            <div className="space-y-3">
              {findings.map((finding) => (
                <FindingCard key={finding.id} finding={finding} />
              ))}
            </div>
          </div>
        )
      })}
    </div>
  )
}
