import { useState, useEffect } from 'react'
import { Spinner } from './Spinner'
import { Card, CardContent } from '@/components/ui/card'

interface LoadingStateProps {
  message: string
  estimatedTime?: string
  startTime?: number // Unix timestamp when loading started
}

// Threshold after which we show "taking longer than usual" message (90 seconds)
const LONG_WAIT_THRESHOLD_MS = 90 * 1000

export function LoadingState({ message, estimatedTime = 'up to 60 seconds', startTime }: LoadingStateProps) {
  const [elapsedSeconds, setElapsedSeconds] = useState(0)
  const [isTakingLong, setIsTakingLong] = useState(false)

  useEffect(() => {
    if (!startTime) return

    const updateElapsed = () => {
      const elapsed = Date.now() - startTime
      setElapsedSeconds(Math.floor(elapsed / 1000))
      setIsTakingLong(elapsed >= LONG_WAIT_THRESHOLD_MS)
    }

    // Initial update
    updateElapsed()

    // Update every second
    const interval = setInterval(updateElapsed, 1000)

    return () => clearInterval(interval)
  }, [startTime])

  const formatElapsed = (seconds: number): string => {
    if (seconds < 60) {
      return `${seconds}s`
    }
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}m ${secs}s`
  }

  return (
    <div className="py-12">
      <Card className="border-border/50 bg-card/50 backdrop-blur max-w-md mx-auto">
        <CardContent className="py-12">
          <div className="flex flex-col items-center justify-center gap-4" role="status" aria-live="polite">
            <div className="relative">
              <div className="absolute inset-0 rounded-full bg-primary/20 animate-ping" />
              <div className="relative p-4 rounded-full bg-primary/10">
                <Spinner className="h-8 w-8 text-primary" />
              </div>
            </div>
            <div className="text-center space-y-1">
              <p className="text-lg font-medium">{message}</p>
              {isTakingLong ? (
                <p className="text-sm text-amber-400">
                  Taking longer than usual ({formatElapsed(elapsedSeconds)})...
                </p>
              ) : startTime ? (
                <p className="text-sm text-muted-foreground">
                  Elapsed: {formatElapsed(elapsedSeconds)}
                </p>
              ) : (
                <p className="text-sm text-muted-foreground">
                  This may take {estimatedTime}
                </p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
