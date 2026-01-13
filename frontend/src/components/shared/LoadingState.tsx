import { useState, useEffect } from 'react'
import { Spinner } from './Spinner'
import { Card, CardContent } from '@/components/ui/card'

export type LoadingType = 'questions' | 'outputs'

interface LoadingStateProps {
  type: LoadingType
  startTime?: number // Unix timestamp when loading started
}

// Threshold after which we show "still working" message (30 seconds)
const STILL_WORKING_THRESHOLD_MS = 30 * 1000

// Loading messages per type
const LOADING_MESSAGES: Record<LoadingType, { initial: string; stillWorking: string }> = {
  questions: {
    initial: 'Generating questions... This may take up to 2 minutes',
    stillWorking: 'Still working on your questions...',
  },
  outputs: {
    initial: 'Generating your files... This may take up to 3 minutes',
    stillWorking: 'Still creating your files... Almost there!',
  },
}

export function LoadingState({ type, startTime }: LoadingStateProps) {
  const [elapsedSeconds, setElapsedSeconds] = useState(0)
  const [showStillWorking, setShowStillWorking] = useState(false)

  useEffect(() => {
    if (!startTime) return

    const updateElapsed = () => {
      const elapsed = Date.now() - startTime
      setElapsedSeconds(Math.floor(elapsed / 1000))
      setShowStillWorking(elapsed >= STILL_WORKING_THRESHOLD_MS)
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

  const messages = LOADING_MESSAGES[type]
  const displayMessage = showStillWorking ? messages.stillWorking : messages.initial

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
              <p className="text-lg font-medium">{displayMessage}</p>
              {startTime && (
                <p className="text-sm text-muted-foreground">
                  Elapsed: {formatElapsed(elapsedSeconds)}
                </p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
