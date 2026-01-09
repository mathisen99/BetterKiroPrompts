import { useState, useEffect } from 'react'

interface RateLimitCountdownProps {
  retryAfterSeconds: number
}

export function RateLimitCountdown({ retryAfterSeconds }: RateLimitCountdownProps) {
  const [secondsRemaining, setSecondsRemaining] = useState(retryAfterSeconds)

  useEffect(() => {
    if (secondsRemaining <= 0) return

    const timer = setInterval(() => {
      setSecondsRemaining(prev => {
        if (prev <= 1) {
          clearInterval(timer)
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(timer)
  }, [secondsRemaining])

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    if (mins > 0) {
      return `${mins}m ${secs}s`
    }
    return `${secs}s`
  }

  if (secondsRemaining <= 0) {
    return (
      <p className="text-sm text-muted-foreground text-center">
        You can now refresh the page to try again.
      </p>
    )
  }

  return (
    <div className="text-center">
      <p className="text-sm text-muted-foreground">
        Rate limit exceeded. Please wait before trying again.
      </p>
      <p className="mt-2 text-lg font-mono text-foreground">
        {formatTime(secondsRemaining)}
      </p>
    </div>
  )
}
