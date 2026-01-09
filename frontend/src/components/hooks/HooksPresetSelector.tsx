import { useState } from 'react'
import type { HooksConfig } from '../../lib/api'
import { PresetCard } from './PresetCard'
import { Spinner } from '../shared/Spinner'

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
        <div className="flex items-center gap-2">
          <input id="has-go" type="checkbox" checked={techStack.hasGo} onChange={(e) => setTechStack((s) => ({ ...s, hasGo: e.target.checked }))} />
          <label htmlFor="has-go" className="text-sm">Go</label>
        </div>
        <div className="flex items-center gap-2">
          <input id="has-typescript" type="checkbox" checked={techStack.hasTypeScript} onChange={(e) => setTechStack((s) => ({ ...s, hasTypeScript: e.target.checked }))} />
          <label htmlFor="has-typescript" className="text-sm">TypeScript</label>
        </div>
        <div className="flex items-center gap-2">
          <input id="has-react" type="checkbox" checked={techStack.hasReact} onChange={(e) => setTechStack((s) => ({ ...s, hasReact: e.target.checked }))} />
          <label htmlFor="has-react" className="text-sm">React</label>
        </div>
      </fieldset>

      <button type="submit" disabled={loading} className="inline-flex items-center justify-center gap-2 w-full rounded px-4 py-2 text-sm bg-primary text-primary-foreground hover:bg-primary/80 disabled:opacity-50">
        {loading && <Spinner />}
        {loading ? 'Generating...' : 'Generate Hooks'}
      </button>
    </form>
  )
}
