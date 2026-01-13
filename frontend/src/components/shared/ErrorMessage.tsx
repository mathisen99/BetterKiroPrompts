import { Button } from '@/components/ui/button'

interface ErrorMessageProps {
  message: string
  canRetry?: boolean
  onRetry?: () => void
  onStartOver?: () => void
}

export function ErrorMessage({ message, canRetry = false, onRetry, onStartOver }: ErrorMessageProps) {
  return (
    <div className="rounded-md border border-destructive/50 bg-destructive/10 p-4 space-y-4" role="alert">
      <p className="text-sm text-destructive">{message}</p>
      
      {(canRetry || onStartOver) && (
        <div className="flex gap-3 justify-center">
          {canRetry && onRetry && (
            <Button 
              variant="default" 
              size="sm" 
              onClick={onRetry}
            >
              Try Again
            </Button>
          )}
          {onStartOver && (
            <Button 
              variant="outline" 
              size="sm" 
              onClick={onStartOver}
            >
              Start Over
            </Button>
          )}
        </div>
      )}
    </div>
  )
}
