import { OutputPanel } from '../shared/OutputPanel'

interface PromptPreviewProps {
  prompt: string
}

export function PromptPreview({ prompt }: PromptPreviewProps) {
  return <OutputPanel content={prompt} filename="kickoff-prompt.md" />
}
