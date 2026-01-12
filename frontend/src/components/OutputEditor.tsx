import { useMemo } from 'react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { toast } from 'sonner'
import type { GeneratedFile } from '@/lib/api'
import { downloadAllAsZip } from '@/lib/zip'
import { Copy, Download, FileText, FolderCog, Webhook, RotateCcw, Package } from 'lucide-react'

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
      <Card key={file.path} className="border-border/50 bg-card/50 backdrop-blur mb-4">
        <CardHeader className="pb-3">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div className="flex items-center gap-2 min-w-0">
              <code className="text-sm font-mono text-muted-foreground truncate">{file.path}</code>
              {isEdited && (
                <span className="shrink-0 text-xs bg-primary/20 text-primary px-2 py-0.5 rounded-full font-medium">
                  Modified
                </span>
              )}
            </div>
            <div className="flex gap-2 shrink-0">
              {isEdited && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onReset(file.path)}
                  className="gap-1.5"
                >
                  <RotateCcw className="h-3.5 w-3.5" />
                  <span className="hidden sm:inline">Reset</span>
                </Button>
              )}
              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleCopy(file.path)}
                className="gap-1.5"
              >
                <Copy className="h-3.5 w-3.5" />
                <span className="hidden sm:inline">Copy</span>
              </Button>
              <Button
                variant="secondary"
                size="sm"
                onClick={() => handleDownload(file.path)}
                className="gap-1.5"
              >
                <Download className="h-3.5 w-3.5" />
                <span className="hidden sm:inline">Download</span>
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          <Textarea
            value={content}
            onChange={(e) => onEdit(file.path, e.target.value)}
            className="font-mono text-sm min-h-[300px] resize-y bg-background/50 border-border/50"
            aria-label={`Edit ${filename}`}
          />
        </CardContent>
      </Card>
    )
  }

  const tabCounts = {
    kickoff: groupedFiles.kickoff.length,
    steering: groupedFiles.steering.length,
    hook: groupedFiles.hook.length,
  }

  // Find first non-empty tab
  const defaultTab = tabCounts.kickoff > 0 ? 'kickoff' : tabCounts.steering > 0 ? 'steering' : 'hook'

  const tabIcons = {
    kickoff: <FileText className="h-4 w-4" />,
    steering: <FolderCog className="h-4 w-4" />,
    hook: <Webhook className="h-4 w-4" />,
  }

  return (
    <div className="py-8 space-y-6">
      <Card className="border-border/50 bg-card/50 backdrop-blur">
        <CardHeader>
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <CardTitle className="text-xl">Generated Files</CardTitle>
              <CardDescription className="mt-1">
                Edit files below, then download individually or as a ZIP
              </CardDescription>
            </div>
            <Button onClick={handleDownloadAll} className="gap-2 w-full sm:w-auto">
              <Package className="h-4 w-4" />
              Download All (ZIP)
            </Button>
          </div>
        </CardHeader>
      </Card>

      <Tabs defaultValue={defaultTab}>
        <TabsList className="w-full sm:w-auto grid grid-cols-3 sm:flex mb-4">
          {tabCounts.kickoff > 0 && (
            <TabsTrigger value="kickoff" className="gap-2">
              {tabIcons.kickoff}
              <span className="hidden sm:inline">Kickoff</span>
              <span className="text-xs text-muted-foreground">({tabCounts.kickoff})</span>
            </TabsTrigger>
          )}
          {tabCounts.steering > 0 && (
            <TabsTrigger value="steering" className="gap-2">
              {tabIcons.steering}
              <span className="hidden sm:inline">Steering</span>
              <span className="text-xs text-muted-foreground">({tabCounts.steering})</span>
            </TabsTrigger>
          )}
          {tabCounts.hook > 0 && (
            <TabsTrigger value="hook" className="gap-2">
              {tabIcons.hook}
              <span className="hidden sm:inline">Hooks</span>
              <span className="text-xs text-muted-foreground">({tabCounts.hook})</span>
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
