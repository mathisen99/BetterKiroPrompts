import { useState } from 'react'
import { HooksPresetSelector } from '../components/hooks/HooksPresetSelector'
import { HookFilePreview } from '../components/hooks/HookFilePreview'
import { generateHooks, type HooksConfig, type GeneratedFile } from '../lib/api'

export function HooksPage() {
  const [files, setFiles] = useState<GeneratedFile[] | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleGenerate = async (config: HooksConfig) => {
    setLoading(true)
    setError(null)
    try {
      const response = await generateHooks(config)
      setFiles(response.files)
    } catch {
      setError('Failed to generate hooks. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="container mx-auto max-w-3xl px-4 py-8">
      <h1 className="mb-8 text-3xl font-bold">Hooks Generator</h1>
      {error && <p className="mb-4 text-sm text-destructive">{error}</p>}
      {files ? (
        <div className="space-y-6">
          <h2 className="text-xl font-semibold">Generated Hooks</h2>
          <HookFilePreview files={files} />
          <button
            onClick={() => setFiles(null)}
            className="rounded px-4 py-2 text-sm bg-secondary text-secondary-foreground hover:bg-secondary/80"
          >
            Back to Selector
          </button>
        </div>
      ) : (
        <HooksPresetSelector onGenerate={handleGenerate} loading={loading} />
      )}
    </main>
  )
}
