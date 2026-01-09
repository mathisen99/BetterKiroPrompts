import { useState } from 'react'
import { SteeringConfigurator } from '../components/steering/SteeringConfigurator'
import { FilePreview } from '../components/steering/FilePreview'
import { ErrorMessage } from '../components/shared/ErrorMessage'
import { generateSteering, type SteeringConfig, type GeneratedFile } from '../lib/api'

export function SteeringPage() {
  const [files, setFiles] = useState<GeneratedFile[] | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [lastConfig, setLastConfig] = useState<SteeringConfig | null>(null)

  const handleGenerate = async (config: SteeringConfig) => {
    setLoading(true)
    setError(null)
    setLastConfig(config)
    try {
      const response = await generateSteering(config)
      setFiles(response.files)
    } catch {
      setError('Failed to generate steering files. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const handleRetry = () => {
    if (lastConfig) handleGenerate(lastConfig)
  }

  return (
    <main className="container mx-auto max-w-3xl px-4 py-8">
      <h1 className="mb-8 text-3xl font-bold">Steering Document Generator</h1>
      {error && <div className="mb-4"><ErrorMessage message={error} onRetry={handleRetry} /></div>}
      {files ? (
        <div className="space-y-6">
          <h2 className="text-xl font-semibold">Generated Files</h2>
          <FilePreview files={files} />
          <button
            onClick={() => setFiles(null)}
            className="rounded px-4 py-2 text-sm bg-secondary text-secondary-foreground hover:bg-secondary/80"
          >
            Back to Configurator
          </button>
        </div>
      ) : (
        <SteeringConfigurator onGenerate={handleGenerate} loading={loading} />
      )}
    </main>
  )
}
