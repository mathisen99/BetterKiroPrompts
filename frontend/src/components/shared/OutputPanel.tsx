import { useState } from 'react'

interface OutputPanelProps {
  content: string
  filename: string
}

export function OutputPanel({ content, filename }: OutputPanelProps) {
  const [copied, setCopied] = useState(false)

  const handleCopy = async () => {
    await navigator.clipboard.writeText(content)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const handleDownload = () => {
    const blob = new Blob([content], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
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
            {copied ? 'Copied!' : 'Copy'}
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
