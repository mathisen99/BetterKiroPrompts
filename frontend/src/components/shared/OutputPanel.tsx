import { toast } from 'sonner'

interface OutputPanelProps {
  content: string
  filename: string
}

export function OutputPanel({ content, filename }: OutputPanelProps) {
  const handleCopy = async () => {
    await navigator.clipboard.writeText(content)
    toast.success('Copied to clipboard')
  }

  const handleDownload = () => {
    const blob = new Blob([content], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
    toast.success('Downloaded successfully')
  }

  return (
    <div className="rounded-lg border border-border bg-card">
      <div className="flex items-center justify-between border-b border-border px-4 py-2">
        <span className="text-sm text-muted-foreground">{filename}</span>
        <div className="flex gap-2">
          <button
            onClick={handleCopy}
            className="rounded px-3 py-1 text-sm bg-secondary text-secondary-foreground hover:bg-secondary/80"
            aria-label="Copy to clipboard"
          >
            Copy
          </button>
          <button
            onClick={handleDownload}
            className="rounded px-3 py-1 text-sm bg-primary text-primary-foreground hover:bg-primary/80"
            aria-label="Download file"
          >
            Download
          </button>
        </div>
      </div>
      <pre className="overflow-auto p-4 text-sm">
        <code>{content}</code>
      </pre>
    </div>
  )
}
