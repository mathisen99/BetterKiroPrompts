import { Card, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import type { ExperienceLevel } from '@/lib/api'

export type { ExperienceLevel }

interface ExperienceLevelSelectorProps {
  onSelect: (level: ExperienceLevel) => void
  selected?: ExperienceLevel
}

interface LevelOption {
  id: ExperienceLevel
  title: string
  description: string
  icon: string
}

const LEVEL_OPTIONS: LevelOption[] = [
  {
    id: 'beginner',
    title: 'Beginner',
    description: 'New to programming. I need guidance on basics.',
    icon: '🌱',
  },
  {
    id: 'novice',
    title: 'Novice',
    description: 'Some experience. I understand basic concepts.',
    icon: '🌿',
  },
  {
    id: 'expert',
    title: 'Expert',
    description: 'Experienced developer. Give me the technical details.',
    icon: '🌳',
  },
]

export function ExperienceLevelSelector({ onSelect, selected }: ExperienceLevelSelectorProps) {
  return (
    <div className="space-y-6">
      <div className="text-center space-y-2">
        <h2 className="text-2xl font-semibold tracking-tight">
          What's your experience level?
        </h2>
        <p className="text-muted-foreground">
          This helps us tailor questions and suggestions to your skill level.
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        {LEVEL_OPTIONS.map((option) => (
          <Card
            key={option.id}
            className={cn(
              'cursor-pointer transition-all hover:border-primary/50 hover:shadow-md',
              selected === option.id && 'border-primary ring-2 ring-primary/20'
            )}
            onClick={() => onSelect(option.id)}
          >
            <CardHeader className="text-center">
              <div className="text-4xl mb-2">{option.icon}</div>
              <CardTitle className="text-lg">{option.title}</CardTitle>
              <CardDescription>{option.description}</CardDescription>
            </CardHeader>
          </Card>
        ))}
      </div>
    </div>
  )
}
