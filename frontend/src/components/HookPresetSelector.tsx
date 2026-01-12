import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import { Shield, ShieldCheck, ShieldAlert, Zap } from 'lucide-react'

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
  icon: React.ReactNode
}

const PRESET_OPTIONS: PresetOption[] = [
  {
    id: 'light',
    title: 'Light',
    description: 'Minimum friction - just formatters on agent stop',
    hooks: ['format-on-stop'],
    icon: <Zap className="h-5 w-5" />,
  },
  {
    id: 'basic',
    title: 'Basic',
    description: 'Daily discipline - formatters, linters, manual test runner',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual'],
    icon: <Shield className="h-5 w-5" />,
  },
  {
    id: 'default',
    title: 'Default (Recommended)',
    description: 'Balanced safety - adds secret scanning and prompt guardrails',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual', 'secret-scan', 'prompt-guardrails'],
    icon: <ShieldCheck className="h-5 w-5" />,
  },
  {
    id: 'strict',
    title: 'Strict',
    description: 'Maximum enforcement - adds static analysis and dependency scanning',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual', 'secret-scan', 'prompt-guardrails', 'static-analysis', 'dep-scan'],
    icon: <ShieldAlert className="h-5 w-5" />,
  },
]

export function HookPresetSelector({ onSelect, selected }: HookPresetSelectorProps) {
  return (
    <Card className="border-border/50 bg-card/50 backdrop-blur">
      <CardHeader>
        <CardTitle className="text-lg">Hook Preset</CardTitle>
        <CardDescription>
          Choose how many automated hooks to include in your project.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        {PRESET_OPTIONS.map((option) => (
          <button
            key={option.id}
            type="button"
            className={cn(
              'w-full text-left p-4 rounded-lg border transition-all',
              'hover:border-primary/50 hover:bg-background/50',
              selected === option.id
                ? 'border-primary bg-primary/5 ring-1 ring-primary/20'
                : 'border-border/30 bg-background/30'
            )}
            onClick={() => onSelect(option.id)}
          >
            <div className="flex items-start gap-3">
              <div className={cn(
                'shrink-0 p-2 rounded-lg transition-colors',
                selected === option.id ? 'bg-primary text-primary-foreground' : 'bg-primary/10 text-primary'
              )}>
                {option.icon}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between gap-2">
                  <span className="font-medium">{option.title}</span>
                  <div className={cn(
                    'shrink-0 w-4 h-4 rounded-full border-2 transition-colors',
                    selected === option.id ? 'border-primary bg-primary' : 'border-muted-foreground/30'
                  )}>
                    {selected === option.id && (
                      <div className="w-full h-full flex items-center justify-center">
                        <div className="w-1.5 h-1.5 rounded-full bg-primary-foreground" />
                      </div>
                    )}
                  </div>
                </div>
                <p className="text-sm text-muted-foreground mt-0.5">
                  {option.description}
                </p>
                <div className="flex flex-wrap gap-1.5 mt-2">
                  {option.hooks.map((hook) => (
                    <span
                      key={hook}
                      className="text-xs px-2 py-0.5 rounded-full bg-muted/50 text-muted-foreground"
                    >
                      {hook}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </button>
        ))}
      </CardContent>
    </Card>
  )
}
