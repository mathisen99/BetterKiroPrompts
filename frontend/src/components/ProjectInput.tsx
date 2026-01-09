import { useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Spinner } from '@/components/shared/Spinner'

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
    <div className="py-12">
      <h2 className="text-2xl font-semibold text-center mb-2">
        What project do you want to make?
      </h2>
      <p className="text-muted-foreground text-center mb-8">
        Describe your project idea and we'll generate tailored Kiro prompts, steering files, and hooks.
      </p>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          type="text"
          value={idea}
          onChange={(e) => setIdea(e.target.value)}
          placeholder="What project do you want to make?"
          disabled={loading}
          className="text-lg py-6"
          aria-label="Project idea"
        />
        
        <Button
          type="submit"
          disabled={loading || !idea.trim()}
          className="w-full"
          size="lg"
        >
          {loading ? (
            <>
              <Spinner className="mr-2" />
              Generating...
            </>
          ) : (
            'Generate Questions'
          )}
        </Button>
      </form>

      <div className="mt-8">
        <p className="text-sm text-muted-foreground mb-3">Or try one of these examples:</p>
        <div className="flex flex-wrap gap-2">
          {examples.map((example) => (
            <Badge
              key={example}
              variant="secondary"
              className="cursor-pointer hover:bg-secondary/80 transition-colors px-3 py-1.5"
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
      </div>
    </div>
  )
}
