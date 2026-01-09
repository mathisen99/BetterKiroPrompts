import { useState } from 'react'
import type { HooksConfig } from '../../lib/api'
import { PresetCard } from './PresetCard'

const presets = [
  { value: 'light', label: 'Light', desc: 'Formatters only' },
  { value: 'basic', label: 'Basic', desc: 'Formatters + linters + tests' },
  { value: 'default', label: 'Default', desc: 'Basic + secret scan + prompt guard' },
  { value: 'strict', label: 'Strict', desc: 'Default + static analysis + vuln scan' },
] as const

interface HooksPresetSelectorProps {
  onGenerate: (config: HooksConfig) => void
  loading?: boolean
}

export function HooksPresetSelector({ onGenerate, loading }: HooksPresetSelectorProps) {
  const [preset, setPreset] = useState<HooksConfig['preset']>('default')
  const [techStack, setTechStack] = useState({ hasGo: true, hasTypeScript: true, hasReact: true })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onGenerate({ preset, techStack })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <fieldset className="space-y-3">
        <legend className="text-sm font-medium">Select Preset</legend>
        {presets.map((p) => (
          <PresetCard
            key={p.value}
            value={p.value}
            label={p.label}
            description={p.desc}
            selected={preset === p.value}
            onSelect={() => setPreset(p.value)}
          />
        ))}
      </fieldset>

      <fieldset className="space-y-3 border-t border-border pt-4">
        <legend className="text-sm font-medium">Tech Stack</legend>
        <label className="flex items-center gap-2">
          <input type="checkbox" checked={techStack.hasGo} onChange={(e) => setTechStack((s) => ({ ...s, hasGo: e.target.checked }))} />
          <span className="text-sm">Go</span>
        </label>
        <label className="flex items-center gap-2">
          <input type="checkbox" checked={techStack.hasTypeScript} onChange={(e) => setTechStack((s) => ({ ...s, hasTypeScript: e.target.checked }))} />
          <span className="text-sm">TypeScript</span>
        </label>
        <label className="flex items-center gap-2">
          <input type="checkbox" checked={techStack.hasReact} onChange={(e) => setTechStack((s) => ({ ...s, hasReact: e.target.checked }))} />
          <span className="text-sm">React</span>
        </label>
      </fieldset>

      <button type="submit" disabled={loading} className="w-full rounded px-4 py-2 text-sm bg-primary text-primary-foreground hover:bg-primary/80 disabled:opacity-50">
        {loading ? 'Generating...' : 'Generate Hooks'}
      </button>
    </form>
  )
}
