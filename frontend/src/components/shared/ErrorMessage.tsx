import { Button } from '@/components/ui/button'

interface ErrorMessageProps {
  message: string
  onRetry?: () => void
}

export function ErrorMessage({ message, onRetry }: ErrorMessageProps) {
  return (
    <div className="rounded-md border border-destructive/50 bg-destructive/10 p-4" role="alert">
      <p className="text-sm text-destructive">{message}</p>
      {onRetry && (
        <Button variant="outline" size="sm" onClick={onRetry} className="mt-3">
          Try again
        </Button>
      )}
    </div>
  )
}
