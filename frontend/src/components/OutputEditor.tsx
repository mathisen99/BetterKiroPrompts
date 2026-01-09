import { useMemo } from 'react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { toast } from 'sonner'
import type { GeneratedFile } from '@/lib/api'
import { downloadAllAsZip } from '@/lib/zip'

interface OutputEditorProps {
  files: GeneratedFile[]
  editedFiles: Map<string, string>
  onEdit: (path: string, content: string) => void
  onReset: (path: string) => void
  getFileContent: (path: string) => string
}

export function OutputEditor({
  files,
  editedFiles,
  onEdit,
  onReset,
  getFileContent,
}: OutputEditorProps) {
  const groupedFiles = useMemo(() => {
    const groups: Record<string, GeneratedFile[]> = {
      kickoff: [],
      steering: [],
      hook: [],
    }
    files.forEach((file) => {
      if (groups[file.type]) {
        groups[file.type].push(file)
      }
    })
    return groups
  }, [files])

  const handleCopy = async (path: string) => {
    const content = getFileContent(path)
    await navigator.clipboard.writeText(content)
    toast.success('Copied to clipboard')
  }

  const handleDownload = (path: string) => {
    const content = getFileContent(path)
    const filename = path.split('/').pop() ?? path
    const blob = new Blob([content], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
    toast.success('Downloaded successfully')
  }

  const handleDownloadAll = () => {
    const filesWithContent = files.map((file) => ({
      path: file.path,
      content: getFileContent(file.path),
    }))
    downloadAllAsZip(filesWithContent)
    toast.success('ZIP downloaded successfully')
  }

  const renderFileEditor = (file: GeneratedFile) => {
    const content = getFileContent(file.path)
    const isEdited = editedFiles.has(file.path)
    const filename = file.path.split('/').pop() ?? file.path

    return (
      <div key={file.path} className="rounded-lg border border-border bg-card mb-4">
        <div className="flex items-center justify-between border-b border-border px-4 py-2">
          <div className="flex items-center gap-2">
            <span className="text-sm font-mono text-muted-foreground">{file.path}</span>
            {isEdited && (
              <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                Modified
              </span>
            )}
          </div>
          <div className="flex gap-2">
            {isEdited && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onReset(file.path)}
              >
                Reset
              </Button>
            )}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleCopy(file.path)}
            >
              Copy
            </Button>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => handleDownload(file.path)}
            >
              Download
            </Button>
          </div>
        </div>
        <div className="p-4">
          <Textarea
            value={content}
            onChange={(e) => onEdit(file.path, e.target.value)}
            className="font-mono text-sm min-h-[300px] resize-y"
            aria-label={`Edit ${filename}`}
          />
        </div>
      </div>
    )
  }

  const tabCounts = {
    kickoff: groupedFiles.kickoff.length,
    steering: groupedFiles.steering.length,
    hook: groupedFiles.hook.length,
  }

  // Find first non-empty tab
  const defaultTab = tabCounts.kickoff > 0 ? 'kickoff' : tabCounts.steering > 0 ? 'steering' : 'hook'

  return (
    <div className="py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-xl font-semibold">Generated Files</h2>
          <p className="text-sm text-muted-foreground">
            Edit files below, then download individually or as a ZIP
          </p>
        </div>
        <Button onClick={handleDownloadAll}>
          Download All (ZIP)
        </Button>
      </div>

      <Tabs defaultValue={defaultTab}>
        <TabsList className="mb-4">
          {tabCounts.kickoff > 0 && (
            <TabsTrigger value="kickoff">
              Kickoff ({tabCounts.kickoff})
            </TabsTrigger>
          )}
          {tabCounts.steering > 0 && (
            <TabsTrigger value="steering">
              Steering ({tabCounts.steering})
            </TabsTrigger>
          )}
          {tabCounts.hook > 0 && (
            <TabsTrigger value="hook">
              Hooks ({tabCounts.hook})
            </TabsTrigger>
          )}
        </TabsList>

        {tabCounts.kickoff > 0 && (
          <TabsContent value="kickoff">
            {groupedFiles.kickoff.map(renderFileEditor)}
          </TabsContent>
        )}

        {tabCounts.steering > 0 && (
          <TabsContent value="steering">
            {groupedFiles.steering.map(renderFileEditor)}
          </TabsContent>
        )}

        {tabCounts.hook > 0 && (
          <TabsContent value="hook">
            {groupedFiles.hook.map(renderFileEditor)}
          </TabsContent>
        )}
      </Tabs>
    </div>
  )
}
