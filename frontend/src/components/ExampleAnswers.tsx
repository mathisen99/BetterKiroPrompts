import { Button } from '@/components/ui/button'
import { Lightbulb } from 'lucide-react'

interface ExampleAnswersProps {
  examples: string[]
  onSelect: (example: string) => void
}

export function ExampleAnswers({ examples, onSelect }: ExampleAnswersProps) {
  if (!examples || examples.length === 0) {
    return null
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Lightbulb className="h-4 w-4" />
        <span>Click an example to use it as your answer:</span>
      </div>
      <div className="flex flex-wrap gap-2">
        {examples.map((example, index) => (
          <Button
            key={index}
            type="button"
            variant="outline"
            size="sm"
            onClick={() => onSelect(example)}
            className="text-left h-auto py-2 px-3 whitespace-normal max-w-full hover:bg-primary/10 hover:border-primary/50 transition-colors"
          >
            <span className="line-clamp-2">{example}</span>
          </Button>
        ))}
      </div>
    </div>
  )
}
