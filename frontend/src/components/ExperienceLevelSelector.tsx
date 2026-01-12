import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import type { ExperienceLevel } from '@/lib/api'
import { Sprout, Leaf, TreeDeciduous } from 'lucide-react'

export type { ExperienceLevel }

interface ExperienceLevelSelectorProps {
  onSelect: (level: ExperienceLevel) => void
  selected?: ExperienceLevel
}

interface LevelOption {
  id: ExperienceLevel
  title: string
  description: string
  icon: React.ReactNode
  details: string[]
}

const LEVEL_OPTIONS: LevelOption[] = [
  {
    id: 'beginner',
    title: 'Beginner',
    description: 'New to programming',
    icon: <Sprout className="h-8 w-8" />,
    details: ['Simple explanations', 'No jargon', 'Step-by-step guidance'],
  },
  {
    id: 'novice',
    title: 'Novice',
    description: 'Some experience',
    icon: <Leaf className="h-8 w-8" />,
    details: ['Moderate technical terms', 'Helpful hints', 'Best practices'],
  },
  {
    id: 'expert',
    title: 'Expert',
    description: 'Experienced developer',
    icon: <TreeDeciduous className="h-8 w-8" />,
    details: ['Technical details', 'Architecture focus', 'Advanced patterns'],
  },
]

export function ExperienceLevelSelector({ onSelect, selected }: ExperienceLevelSelectorProps) {
  return (
    <div className="py-8 space-y-8">
      <div className="text-center space-y-3">
        <h2 className="text-3xl font-bold tracking-tight">
          What's your experience level?
        </h2>
        <p className="text-muted-foreground text-lg max-w-xl mx-auto">
          This helps us tailor questions and suggestions to your skill level.
        </p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {LEVEL_OPTIONS.map((option) => (
          <Card
            key={option.id}
            className={cn(
              'cursor-pointer transition-all hover:border-primary/50 hover:shadow-lg hover:shadow-primary/5',
              selected === option.id && 'border-primary ring-2 ring-primary/20 bg-primary/5'
            )}
            onClick={() => onSelect(option.id)}
          >
            <CardHeader className="text-center pb-2">
              <div className={cn(
                'mx-auto mb-3 p-3 rounded-xl transition-colors',
                selected === option.id ? 'bg-primary text-primary-foreground' : 'bg-primary/10 text-primary'
              )}>
                {option.icon}
              </div>
              <CardTitle className="text-xl">{option.title}</CardTitle>
              <CardDescription className="text-base">{option.description}</CardDescription>
            </CardHeader>
            <CardContent className="pt-0">
              <ul className="space-y-1.5">
                {option.details.map((detail) => (
                  <li key={detail} className="flex items-center gap-2 text-sm text-muted-foreground">
                    <span className="h-1 w-1 rounded-full bg-primary/50" />
                    {detail}
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
