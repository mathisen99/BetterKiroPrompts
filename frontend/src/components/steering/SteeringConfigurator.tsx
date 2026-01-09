import { useState } from 'react'
import type { SteeringConfig } from '../../lib/api'
import { Spinner } from '../shared/Spinner'

const initialConfig: SteeringConfig = {
  projectName: '',
  projectDescription: '',
  techStack: { backend: '', frontend: '', database: '' },
  includeConditional: false,
  includeManual: false,
  fileReferences: [],
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
        <div>
          <label htmlFor="project-name" className="block text-sm font-medium">Project Name</label>
          <input id="project-name" type="text" value={config.projectName} onChange={(e) => updateConfig('projectName', e.target.value)} required className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
        </div>
        <div>
          <label htmlFor="project-description" className="block text-sm font-medium">Project Description</label>
          <textarea id="project-description" value={config.projectDescription} onChange={(e) => updateConfig('projectDescription', e.target.value)} rows={2} className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
        </div>
      </div>

      <div className="space-y-4 border-t border-border pt-4">
        <p className="text-sm font-medium">Tech Stack</p>
        <div className="grid gap-3">
          <div>
            <label htmlFor="tech-backend" className="block text-sm text-muted-foreground">Backend</label>
            <input id="tech-backend" type="text" value={config.techStack.backend} onChange={(e) => updateTechStack('backend', e.target.value)} placeholder="e.g., Go, Node.js" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
          </div>
          <div>
            <label htmlFor="tech-frontend" className="block text-sm text-muted-foreground">Frontend</label>
            <input id="tech-frontend" type="text" value={config.techStack.frontend} onChange={(e) => updateTechStack('frontend', e.target.value)} placeholder="e.g., React, Vue" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
          </div>
          <div>
            <label htmlFor="tech-database" className="block text-sm text-muted-foreground">Database</label>
            <input id="tech-database" type="text" value={config.techStack.database} onChange={(e) => updateTechStack('database', e.target.value)} placeholder="e.g., PostgreSQL, MongoDB" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
          </div>
        </div>
      </div>

      <div className="flex items-center gap-2 border-t border-border pt-4">
        <input id="include-conditional" type="checkbox" checked={config.includeConditional} onChange={(e) => updateConfig('includeConditional', e.target.checked)} className="rounded border-input" />
        <label htmlFor="include-conditional" className="text-sm">Include conditional steering files (security, quality)</label>
      </div>

      <div className="flex items-center gap-2">
        <input id="include-manual" type="checkbox" checked={config.includeManual} onChange={(e) => updateConfig('includeManual', e.target.checked)} className="rounded border-input" />
        <label htmlFor="include-manual" className="text-sm">Include manual steering files (referenced via #steering-file-name)</label>
      </div>

      <div className="space-y-2">
        <label htmlFor="file-references" className="block text-sm font-medium">File References (optional)</label>
        <p className="text-xs text-muted-foreground">Add file paths to reference in steering files (one per line)</p>
        <textarea
          id="file-references"
          value={config.fileReferences.join('\n')}
          onChange={(e) => updateConfig('fileReferences', e.target.value.split('\n').map(s => s.trim()).filter(Boolean))}
          rows={3}
          placeholder=".env.example&#10;backend/migrations/README.md"
          className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
        />
      </div>

      <button type="submit" disabled={!config.projectName || loading} className="inline-flex items-center justify-center gap-2 w-full rounded px-4 py-2 text-sm bg-primary text-primary-foreground hover:bg-primary/80 disabled:opacity-50">
        {loading && <Spinner />}
        {loading ? 'Generating...' : 'Generate Steering Files'}
      </button>
    </form>
  )
}
