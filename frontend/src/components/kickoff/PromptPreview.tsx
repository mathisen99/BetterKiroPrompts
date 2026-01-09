import { OutputPanel } from '../shared/OutputPanel'
import { CommitContract } from '../shared/CommitContract'

interface PromptPreviewProps {
  prompt: string
}

export function PromptPreview({ prompt }: PromptPreviewProps) {
  return (
    <div className="space-y-4">
      <OutputPanel content={prompt} filename="kickoff-prompt.md" />
      <CommitContract />
    </div>
  )
}
