import { Spinner } from './Spinner'

interface LoadingStateProps {
  message: string
  estimatedTime?: string
}

export function LoadingState({ message, estimatedTime = 'up to 60 seconds' }: LoadingStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16 gap-4" role="status" aria-live="polite">
      <Spinner className="h-8 w-8" />
      <p className="text-muted-foreground">{message}</p>
      <p className="text-sm text-muted-foreground">
        This may take {estimatedTime}
      </p>
    </div>
  )
}
