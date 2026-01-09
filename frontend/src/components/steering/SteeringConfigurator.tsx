import { useState } from 'react'
import type { SteeringConfig } from '../../lib/api'

const initialConfig: SteeringConfig = {
  projectName: '',
  projectDescription: '',
  techStack: { backend: '', frontend: '', database: '' },
  includeConditional: false,
  customRules: {},
}

interface SteeringConfiguratorProps {
  onGenerate: (config: SteeringConfig) => void
  loading?: boolean
}

export function SteeringConfigurator({ onGenerate, loading }: SteeringConfiguratorProps) {
  const [config, setConfig] = useState<SteeringConfig>(initialConfig)

  const updateConfig = <K extends keyof SteeringConfig>(key: K, value: SteeringConfig[K]) => {
    setConfig((prev) => ({ ...prev, [key]: value }))
  }

  const updateTechStack = (field: keyof typeof config.techStack, value: string) => {
    setConfig((prev) => ({ ...prev, techStack: { ...prev.techStack, [field]: value } }))
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onGenerate(config)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-4">
        <label className="block">
          <span className="text-sm font-medium">Project Name</span>
          <input type="text" value={config.projectName} onChange={(e) => updateConfig('projectName', e.target.value)} required className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
        </label>
        <label className="block">
          <span className="text-sm font-medium">Project Description</span>
          <textarea value={config.projectDescription} onChange={(e) => updateConfig('projectDescription', e.target.value)} rows={2} className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
        </label>
      </div>

      <div className="space-y-4 border-t border-border pt-4">
        <p className="text-sm font-medium">Tech Stack</p>
        <div className="grid gap-3">
          <label className="block">
            <span className="text-sm text-muted-foreground">Backend</span>
            <input type="text" value={config.techStack.backend} onChange={(e) => updateTechStack('backend', e.target.value)} placeholder="e.g., Go, Node.js" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
          </label>
          <label className="block">
            <span className="text-sm text-muted-foreground">Frontend</span>
            <input type="text" value={config.techStack.frontend} onChange={(e) => updateTechStack('frontend', e.target.value)} placeholder="e.g., React, Vue" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
          </label>
          <label className="block">
            <span className="text-sm text-muted-foreground">Database</span>
            <input type="text" value={config.techStack.database} onChange={(e) => updateTechStack('database', e.target.value)} placeholder="e.g., PostgreSQL, MongoDB" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
          </label>
        </div>
      </div>

      <label className="flex items-center gap-2 border-t border-border pt-4">
        <input type="checkbox" checked={config.includeConditional} onChange={(e) => updateConfig('includeConditional', e.target.checked)} className="rounded border-input" />
        <span className="text-sm">Include conditional steering files (security, quality)</span>
      </label>

      <button type="submit" disabled={!config.projectName || loading} className="w-full rounded px-4 py-2 text-sm bg-primary text-primary-foreground hover:bg-primary/80 disabled:opacity-50">
        {loading ? 'Generating...' : 'Generate Steering Files'}
      </button>
    </form>
  )
}
