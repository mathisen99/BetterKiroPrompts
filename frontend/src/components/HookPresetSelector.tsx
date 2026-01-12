import { Card, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { cn } from '@/lib/utils'

export type HookPreset = 'light' | 'basic' | 'default' | 'strict'

interface HookPresetSelectorProps {
  onSelect: (preset: HookPreset) => void
  selected: HookPreset
}

interface PresetOption {
  id: HookPreset
  title: string
  description: string
  hooks: string[]
}

const PRESET_OPTIONS: PresetOption[] = [
  {
    id: 'light',
    title: 'Light',
    description: 'Minimum friction - just formatters on agent stop',
    hooks: ['format-on-stop'],
  },
  {
    id: 'basic',
    title: 'Basic',
    description: 'Daily discipline - formatters, linters, manual test runner',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual'],
  },
  {
    id: 'default',
    title: 'Default (Recommended)',
    description: 'Balanced safety - adds secret scanning and prompt guardrails',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual', 'secret-scan', 'prompt-guardrails'],
  },
  {
    id: 'strict',
    title: 'Strict',
    description: 'Maximum enforcement - adds static analysis and dependency scanning',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual', 'secret-scan', 'prompt-guardrails', 'static-analysis', 'dep-scan'],
  },
]

export function HookPresetSelector({ onSelect, selected }: HookPresetSelectorProps) {
  return (
    <div className="space-y-4">
      <div className="space-y-1">
        <h3 className="text-lg font-medium">Hook Preset</h3>
        <p className="text-sm text-muted-foreground">
          Choose how many automated hooks to include in your project.
        </p>
      </div>

      <div className="grid gap-3">
        {PRESET_OPTIONS.map((option) => (
          <Card
            key={option.id}
            className={cn(
              'cursor-pointer transition-all hover:border-primary/50',
              selected === option.id && 'border-primary ring-2 ring-primary/20'
            )}
            onClick={() => onSelect(option.id)}
          >
            <CardHeader className="py-3 px-4">
              <div className="flex items-start justify-between gap-4">
                <div className="space-y-1">
                  <CardTitle className="text-base">{option.title}</CardTitle>
                  <CardDescription className="text-sm">
                    {option.description}
                  </CardDescription>
                </div>
                <div className="flex items-center justify-center w-5 h-5 rounded-full border-2 border-muted-foreground/30 shrink-0 mt-0.5">
                  {selected === option.id && (
                    <div className="w-2.5 h-2.5 rounded-full bg-primary" />
                  )}
                </div>
              </div>
              <div className="flex flex-wrap gap-1.5 mt-2">
                {option.hooks.map((hook) => (
                  <span
                    key={hook}
                    className="text-xs px-2 py-0.5 rounded-full bg-muted text-muted-foreground"
                  >
                    {hook}
                  </span>
                ))}
              </div>
            </CardHeader>
          </Card>
        ))}
      </div>
    </div>
  )
}
