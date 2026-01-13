import { useState, useEffect } from 'react'
import { GitBranch, Search, Bot, Loader2 } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import type { ScanJob, ScanStatus } from '@/lib/api'

interface ScanProgressProps {
  job: ScanJob
}

// Status configuration for display
const statusConfig: Record<ScanStatus, {
  icon: typeof Loader2
  label: string
  description: string
  step: number
}> = {
  pending: {
    icon: Loader2,
    label: 'Queued',
    description: 'Waiting to start...',
    step: 0,
  },
  cloning: {
    icon: GitBranch,
    label: 'Cloning Repository',
    description: 'Downloading repository files...',
    step: 1,
  },
  scanning: {
    icon: Search,
    label: 'Running Security Tools',
    description: 'Analyzing code with Trivy, Semgrep, and more...',
    step: 2,
  },
  reviewing: {
    icon: Bot,
    label: 'AI Code Review',
    description: 'Generating remediation guidance...',
    step: 3,
  },
  completed: {
    icon: Loader2,
    label: 'Completed',
    description: 'Scan finished',
    step: 4,
  },
  failed: {
    icon: Loader2,
    label: 'Failed',
    description: 'Scan encountered an error',
    step: -1,
  },
}

const steps = ['Clone', 'Scan', 'Review', 'Done']

export function ScanProgress({ job }: ScanProgressProps) {
  const [elapsedSeconds, setElapsedSeconds] = useState(0)
  
  // Track elapsed time
  useEffect(() => {
    const startTime = new Date(job.created_at).getTime()
    
    const updateElapsed = () => {
      const elapsed = Date.now() - startTime
      setElapsedSeconds(Math.floor(elapsed / 1000))
    }
    
    updateElapsed()
    const interval = setInterval(updateElapsed, 1000)
    
    return () => clearInterval(interval)
  }, [job.created_at])

  const formatElapsed = (seconds: number): string => {
    if (seconds < 60) {
      return `${seconds}s`
    }
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}m ${secs}s`
  }

  const config = statusConfig[job.status]
  const Icon = config.icon
  const currentStep = config.step

  return (
    <Card className="bg-card/50 backdrop-blur border-primary/20">
      <CardContent className="py-8">
        <div className="flex flex-col items-center gap-6">
          {/* Animated Icon */}
          <div className="relative">
            <div className="absolute inset-0 rounded-full bg-primary/20 animate-ping" />
            <div className="relative p-4 rounded-full bg-primary/10">
              <Icon className="h-8 w-8 text-primary animate-spin" />
            </div>
          </div>

          {/* Status Text */}
          <div className="text-center space-y-1">
            <h3 className="text-lg font-semibold">{config.label}</h3>
            <p className="text-sm text-muted-foreground">{config.description}</p>
            <p className="text-xs text-muted-foreground">
              Elapsed: {formatElapsed(elapsedSeconds)}
            </p>
          </div>

          {/* Progress Steps */}
          <div className="w-full max-w-md">
            <div className="flex items-center justify-between">
              {steps.map((step, index) => {
                const isActive = index === currentStep
                const isCompleted = index < currentStep
                
                return (
                  <div key={step} className="flex flex-col items-center gap-1">
                    <div
                      className={`
                        w-8 h-8 rounded-full flex items-center justify-center text-xs font-medium
                        transition-colors duration-300
                        ${isCompleted 
                          ? 'bg-primary text-primary-foreground' 
                          : isActive 
                            ? 'bg-primary/20 text-primary border-2 border-primary' 
                            : 'bg-muted text-muted-foreground'
                        }
                      `}
                    >
                      {isCompleted ? '✓' : index + 1}
                    </div>
                    <span className={`text-xs ${isActive ? 'text-primary font-medium' : 'text-muted-foreground'}`}>
                      {step}
                    </span>
                  </div>
                )
              })}
            </div>
            
            {/* Progress Bar */}
            <div className="mt-4 h-1 bg-muted rounded-full overflow-hidden">
              <div 
                className="h-full bg-primary transition-all duration-500 ease-out"
                style={{ width: `${Math.max(0, (currentStep / (steps.length - 1)) * 100)}%` }}
              />
            </div>
          </div>

          {/* Repository URL */}
          <div className="text-center">
            <p className="text-xs text-muted-foreground">Scanning</p>
            <p className="font-mono text-sm truncate max-w-md" title={job.repo_url}>
              {job.repo_url}
            </p>
          </div>

          {/* Languages (if detected) */}
          {job.languages && job.languages.length > 0 && (
            <div className="flex items-center gap-2 flex-wrap justify-center">
              <span className="text-xs text-muted-foreground">Languages:</span>
              {job.languages.map((lang) => (
                <span 
                  key={lang}
                  className="px-2 py-0.5 text-xs rounded-full bg-secondary text-secondary-foreground"
                >
                  {lang}
                </span>
              ))}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
