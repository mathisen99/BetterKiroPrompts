import { Spinner } from './Spinner'
import { Card, CardContent } from '@/components/ui/card'

interface LoadingStateProps {
  message: string
  estimatedTime?: string
}

export function LoadingState({ message, estimatedTime = 'up to 60 seconds' }: LoadingStateProps) {
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
              <p className="text-sm text-muted-foreground">
                This may take {estimatedTime}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
