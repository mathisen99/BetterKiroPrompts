import { useState, useEffect, useCallback } from 'react'
import { ArrowLeft, Home, ImageIcon, Shield, Lock, Unlock, Search, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ScanProgress } from '@/components/ScanProgress'
import { ScanResults } from '@/components/ScanResults'
import { startScan, getScanStatus, getScanConfig } from '@/lib/api'
import type { ScanJob, ScanConfig } from '@/lib/api'

interface SecurityScanPageProps {
  onNavigateHome: () => void
  onNavigateGallery: () => void
}

export function SecurityScanPage({ onNavigateHome, onNavigateGallery }: SecurityScanPageProps) {
  const [repoUrl, setRepoUrl] = useState('')
  const [config, setConfig] = useState<ScanConfig | null>(null)
  const [currentJob, setCurrentJob] = useState<ScanJob | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  // Load scan configuration on mount
  useEffect(() => {
    getScanConfig()
      .then(setConfig)
      .catch(() => {
        // Config endpoint might not be available, use defaults
        setConfig({ privateRepoEnabled: false })
      })
  }, [])

  // Poll for scan status when job is in progress
  useEffect(() => {
    if (!currentJob || currentJob.status === 'completed' || currentJob.status === 'failed') {
      return
    }

    const pollInterval = setInterval(async () => {
      try {
        const updatedJob = await getScanStatus(currentJob.id)
        setCurrentJob(updatedJob)
      } catch {
        // Ignore polling errors
      }
    }, 2000) // Poll every 2 seconds

    return () => clearInterval(pollInterval)
  }, [currentJob])

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setIsSubmitting(true)

    try {
      const job = await startScan(repoUrl)
      setCurrentJob(job)
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message)
      } else {
        setError('Failed to start scan. Please try again.')
      }
    } finally {
      setIsSubmitting(false)
    }
  }, [repoUrl])

  const handleNewScan = useCallback(() => {
    setCurrentJob(null)
    setRepoUrl('')
    setError(null)
  }, [])

  const isScanning = currentJob && currentJob.status !== 'completed' && currentJob.status !== 'failed'

  return (
    <div className="min-h-screen">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={onNavigateHome} aria-label="Go back">
              <ArrowLeft className="h-5 w-5" />
            </Button>
            <div>
              <h1 className="text-2xl font-bold flex items-center gap-2">
                <Shield className="h-6 w-6 text-primary" />
                Security Scan
              </h1>
              <p className="text-muted-foreground">
                Scan repositories for vulnerabilities
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={onNavigateGallery}
              className="gap-1.5"
            >
              <ImageIcon className="h-4 w-4" />
              Gallery
            </Button>
            <Button
              variant="outline"
              onClick={onNavigateHome}
              className="gap-2"
            >
              <Home className="h-4 w-4" />
              Back to Home
            </Button>
          </div>
        </div>

        {/* Main Content */}
        <div className="max-w-4xl mx-auto">
          {/* Scan Form - Show when no active job */}
          {!currentJob && (
            <>
              {/* Explanation Card */}
              <Card className="mb-6 bg-card/50 backdrop-blur border-primary/20">
                <CardHeader>
                  <CardTitle className="text-lg">How it works</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-sm space-y-2">
                    <p>
                      Enter a GitHub repository URL to scan for security vulnerabilities. 
                      The scanner runs multiple security tools including Trivy, Semgrep, 
                      TruffleHog, and language-specific analyzers.
                    </p>
                    <p>
                      If issues are found, AI-powered code review provides actionable 
                      remediation guidance with concrete code examples.
                    </p>
                  </CardDescription>
                </CardContent>
              </Card>

              {/* Private Repo Indicator */}
              <div className="mb-4 flex items-center gap-2 text-sm">
                {config?.privateRepoEnabled ? (
                  <>
                    <Unlock className="h-4 w-4 text-green-400" />
                    <span className="text-green-400">Private repository scanning enabled</span>
                  </>
                ) : (
                  <>
                    <Lock className="h-4 w-4 text-muted-foreground" />
                    <span className="text-muted-foreground">
                      Only public repositories supported (configure GITHUB_TOKEN for private repos)
                    </span>
                  </>
                )}
              </div>

              {/* Scan Form */}
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="flex gap-2">
                  <Input
                    type="url"
                    placeholder="https://github.com/owner/repo"
                    value={repoUrl}
                    onChange={(e) => setRepoUrl(e.target.value)}
                    disabled={isSubmitting}
                    className="flex-1"
                    aria-label="Repository URL"
                    required
                  />
                  <Button 
                    type="submit" 
                    disabled={isSubmitting || !repoUrl.trim()}
                    className="gap-2"
                  >
                    <Search className="h-4 w-4" />
                    {isSubmitting ? 'Starting...' : 'Start Scan'}
                  </Button>
                </div>

                {error && (
                  <div className="flex items-center gap-2 text-destructive text-sm">
                    <AlertCircle className="h-4 w-4" />
                    {error}
                  </div>
                )}
              </form>
            </>
          )}

          {/* Scan Progress - Show when scanning */}
          {isScanning && currentJob && (
            <ScanProgress job={currentJob} />
          )}

          {/* Scan Results - Show when completed or failed */}
          {currentJob && (currentJob.status === 'completed' || currentJob.status === 'failed') && (
            <div className="space-y-4">
              <ScanResults job={currentJob} />
              <div className="flex justify-center">
                <Button onClick={handleNewScan} variant="outline" className="gap-2">
                  <Search className="h-4 w-4" />
                  Start New Scan
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
