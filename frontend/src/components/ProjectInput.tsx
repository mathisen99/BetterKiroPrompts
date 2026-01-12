import { useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Spinner } from '@/components/shared/Spinner'
import { Lightbulb, Sparkles } from 'lucide-react'

interface ProjectInputProps {
  onSubmit: (idea: string) => void
  loading: boolean
  examples: string[]
}

export function ProjectInput({ onSubmit, loading, examples }: ProjectInputProps) {
  const [idea, setIdea] = useState('')

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    const trimmed = idea.trim()
    if (trimmed) {
      onSubmit(trimmed)
    }
  }

  const handleExampleClick = (example: string) => {
    setIdea(example)
  }

  return (
    <div className="py-8 space-y-8">
      <div className="text-center space-y-3">
        <h2 className="text-3xl font-bold tracking-tight">
          What are you building?
        </h2>
        <p className="text-muted-foreground text-lg max-w-xl mx-auto">
          Describe your project idea and we'll generate tailored Kiro prompts, steering files, and hooks.
        </p>
      </div>

      <Card className="border-border/50 bg-card/50 backdrop-blur">
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center gap-2 text-lg">
            <Sparkles className="h-5 w-5 text-primary" />
            Project Description
          </CardTitle>
          <CardDescription>
            Be specific about features, tech stack, and goals for better results.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <Textarea
              value={idea}
              onChange={(e) => setIdea(e.target.value)}
              placeholder="Example: A task management app with user authentication, project boards, and real-time collaboration using React and Node.js..."
              disabled={loading}
              className="min-h-[140px] text-base resize-none bg-background/50 border-border/50 focus:border-primary/50"
              aria-label="Project idea"
            />
            
            <Button
              type="submit"
              disabled={loading || !idea.trim()}
              className="w-full h-12 text-base font-medium"
              size="lg"
            >
              {loading ? (
                <>
                  <Spinner className="mr-2" />
                  Generating...
                </>
              ) : (
                <>
                  <Sparkles className="mr-2 h-4 w-4" />
                  Generate Questions
                </>
              )}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card className="border-border/30 bg-card/30">
        <CardHeader className="pb-3">
          <CardTitle className="flex items-center gap-2 text-base font-medium">
            <Lightbulb className="h-4 w-4 text-primary/70" />
            Need inspiration?
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            {examples.map((example) => (
              <Badge
                key={example}
                variant="secondary"
                className="cursor-pointer hover:bg-primary/20 hover:text-primary-foreground transition-all px-3 py-1.5 text-sm"
                onClick={() => handleExampleClick(example)}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault()
                    handleExampleClick(example)
                  }
                }}
              >
                {example}
              </Badge>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
