import { useState } from 'react'
import { OutputPanel } from '../shared/OutputPanel'
import { downloadAsZip } from '../../lib/zip'
import type { GeneratedFile } from '../../lib/api'

interface HookFilePreviewProps {
  files: GeneratedFile[]
}

export function HookFilePreview({ files }: HookFilePreviewProps) {
  const [activeIndex, setActiveIndex] = useState(0)

  if (files.length === 0) return null

  const activeFile = files[activeIndex]

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between border-b border-border pb-2">
        <div className="flex flex-wrap gap-2" role="tablist">
          {files.map((file, i) => (
            <button
              key={file.path}
              onClick={() => setActiveIndex(i)}
              role="tab"
              aria-selected={i === activeIndex}
              className={`rounded px-3 py-1 text-sm ${i === activeIndex ? 'bg-primary text-primary-foreground' : 'bg-secondary text-secondary-foreground hover:bg-secondary/80'}`}
            >
              {file.path.split('/').pop()}
            </button>
          ))}
        </div>
        <button
          onClick={() => downloadAsZip(files, 'hooks.zip')}
          className="rounded px-3 py-1 text-sm bg-primary text-primary-foreground hover:bg-primary/80"
          aria-label="Download all files as zip"
        >
          Download All
        </button>
      </div>
      <OutputPanel content={activeFile.content} filename={activeFile.path.split('/').pop() || activeFile.path} />
    </div>
  )
}
