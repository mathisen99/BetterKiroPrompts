import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

interface RestorePromptProps {
  projectIdea: string
  onRestore: () => void
  onStartFresh: () => void
}

export function RestorePrompt({ projectIdea, onRestore, onStartFresh }: RestorePromptProps) {
  // Truncate project idea for display
  const displayIdea = projectIdea.length > 100 
    ? projectIdea.slice(0, 100) + '...' 
    : projectIdea

  return (
    <Card className="border-blue-500/30 bg-card/50 backdrop-blur">
      <CardHeader className="text-center">
        <CardTitle className="text-xl">Welcome Back!</CardTitle>
        <CardDescription>
          We found your previous session. Would you like to continue where you left off?
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {displayIdea && (
          <div className="p-4 rounded-lg bg-muted/50 border border-border">
            <p className="text-sm text-muted-foreground mb-1">Your project idea:</p>
            <p className="text-sm font-medium">{displayIdea}</p>
          </div>
        )}
        
        <div className="flex flex-col sm:flex-row gap-3">
          <Button 
            onClick={onRestore}
            className="flex-1"
            variant="default"
          >
            Continue Session
          </Button>
          <Button 
            onClick={onStartFresh}
            className="flex-1"
            variant="outline"
          >
            Start Fresh
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
