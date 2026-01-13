import { useMemo, useEffect, useCallback, useRef } from 'react'
import { X, Star, Eye, Clock, Copy, Download, Package, FileText, FolderCog, Webhook, Bot } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { toast } from 'sonner'
import { SyntaxHighlighter } from '@/components/SyntaxHighlighter'
import { detectLanguage } from '@/lib/syntax'
import { downloadAllAsZip } from '@/lib/zip'
import type { GalleryDetail as GalleryDetailType, GeneratedFile } from '@/lib/api'
import { Rating } from './Rating'

interface GalleryDetailProps {
  generation: GalleryDetailType
  onClose: () => void
  onRate: (score: number) => void
  userRating: number | null
  isRating?: boolean
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export function GalleryDetail({
  generation,
  onClose,
  onRate,
  userRating,
  isRating,
}: GalleryDetailProps) {
  const files = generation.files as GeneratedFile[]
  const modalContentRef = useRef<HTMLDivElement>(null)

  // Handle Escape key to close modal
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose()
      }
    }
    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [onClose])

  // Handle click outside modal content to close
  const handleBackdropClick = useCallback((event: React.MouseEvent<HTMLDivElement>) => {
    // Only close if clicking directly on the backdrop, not on modal content
    if (event.target === event.currentTarget) {
      onClose()
    }
  }, [onClose])

  const groupedFiles = useMemo(() => {
    const groups: Record<string, GeneratedFile[]> = {
      kickoff: [],
      steering: [],
      hook: [],
      agents: [],
    }
    files.forEach((file) => {
      if (groups[file.type]) {
        groups[file.type].push(file)
      }
    })
    return groups
  }, [files])

  const handleCopy = async (content: string) => {
    await navigator.clipboard.writeText(content)
    toast.success('Copied to clipboard')
  }

  const handleDownload = (path: string, content: string) => {
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
      content: file.content,
    }))
    downloadAllAsZip(filesWithContent)
    toast.success('ZIP downloaded successfully')
  }

  const renderFile = (file: GeneratedFile) => {
    const language = detectLanguage(file.path)

    return (
      <Card key={file.path} className="border-border/50 bg-card/50 backdrop-blur mb-4">
        <CardHeader className="pb-3">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <code className="text-sm font-mono text-muted-foreground truncate">{file.path}</code>
            <div className="flex gap-2 shrink-0">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleCopy(file.content)}
                className="gap-1.5"
              >
                <Copy className="h-3.5 w-3.5" />
                <span className="hidden sm:inline">Copy</span>
              </Button>
              <Button
                variant="secondary"
                size="sm"
                onClick={() => handleDownload(file.path, file.content)}
                className="gap-1.5"
              >
                <Download className="h-3.5 w-3.5" />
                <span className="hidden sm:inline">Download</span>
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          <div className="max-h-[400px] overflow-auto rounded-lg border border-border/50">
            <SyntaxHighlighter code={file.content} language={language} />
          </div>
        </CardContent>
      </Card>
    )
  }

  const tabCounts = {
    kickoff: groupedFiles.kickoff.length,
    steering: groupedFiles.steering.length,
    hook: groupedFiles.hook.length,
    agents: groupedFiles.agents.length,
  }

  const defaultTab =
    tabCounts.kickoff > 0
      ? 'kickoff'
      : tabCounts.steering > 0
        ? 'steering'
        : tabCounts.hook > 0
          ? 'hook'
          : 'agents'

  const tabIcons = {
    kickoff: <FileText className="h-4 w-4" />,
    steering: <FolderCog className="h-4 w-4" />,
    hook: <Webhook className="h-4 w-4" />,
    agents: <Bot className="h-4 w-4" />,
  }

  return (
    <div 
      className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-background/80 backdrop-blur-sm p-4"
      onClick={handleBackdropClick}
      role="dialog"
      aria-modal="true"
      aria-labelledby="gallery-detail-title"
    >
      <div ref={modalContentRef} className="relative w-full max-w-4xl my-8">
        {/* Close button - positioned at top right of modal, outside card */}
        <Button
          variant="default"
          size="icon"
          className="absolute -top-2 -right-2 z-10 h-11 w-11 rounded-full bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg"
          onClick={onClose}
          aria-label="Close modal"
        >
          <X className="h-6 w-6" />
        </Button>

        {/* Header card */}
        <Card className="mb-6">
          <CardHeader className="pr-12">
            <div className="flex flex-col gap-4">
              <div className="flex items-start justify-between gap-4">
                <div className="space-y-2 flex-1 min-w-0">
                  <CardTitle id="gallery-detail-title" className="text-xl">{generation.projectIdea}</CardTitle>
                  <div className="flex flex-wrap items-center gap-2">
                    <Badge variant="secondary">{generation.category}</Badge>
                    <Badge variant="outline">{generation.experienceLevel}</Badge>
                    <Badge variant="outline">{generation.hookPreset} hooks</Badge>
                  </div>
                </div>
              </div>
              
              {/* Download button on its own row */}
              <Button onClick={handleDownloadAll} className="gap-2 w-fit">
                <Package className="h-4 w-4" />
                Download All
              </Button>

              <CardDescription className="flex flex-wrap items-center gap-4">
                <span className="flex items-center gap-1">
                  <Star className="h-4 w-4 fill-yellow-500 text-yellow-500" />
                  {generation.avgRating.toFixed(1)} ({generation.ratingCount} ratings)
                </span>
                <span className="flex items-center gap-1">
                  <Eye className="h-4 w-4" />
                  {generation.viewCount} views
                </span>
                <span className="flex items-center gap-1">
                  <Clock className="h-4 w-4" />
                  {formatDate(generation.createdAt)}
                </span>
              </CardDescription>

              {/* Rating section */}
              <div className="border-t pt-4">
                <div className="flex items-center gap-4">
                  <span className="text-sm text-muted-foreground">Rate this generation:</span>
                  <Rating
                    value={generation.avgRating}
                    count={generation.ratingCount}
                    userRating={userRating}
                    onRate={onRate}
                    disabled={isRating}
                  />
                </div>
              </div>
            </div>
          </CardHeader>
        </Card>

        {/* Files tabs */}
        <Tabs defaultValue={defaultTab}>
          <TabsList className="w-full sm:w-auto grid grid-cols-4 sm:flex mb-4">
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
            {tabCounts.agents > 0 && (
              <TabsTrigger value="agents" className="gap-2">
                {tabIcons.agents}
                <span className="hidden sm:inline">Agents</span>
                <span className="text-xs text-muted-foreground">({tabCounts.agents})</span>
              </TabsTrigger>
            )}
          </TabsList>

          {tabCounts.kickoff > 0 && (
            <TabsContent value="kickoff">{groupedFiles.kickoff.map(renderFile)}</TabsContent>
          )}
          {tabCounts.steering > 0 && (
            <TabsContent value="steering">{groupedFiles.steering.map(renderFile)}</TabsContent>
          )}
          {tabCounts.hook > 0 && (
            <TabsContent value="hook">{groupedFiles.hook.map(renderFile)}</TabsContent>
          )}
          {tabCounts.agents > 0 && (
            <TabsContent value="agents">{groupedFiles.agents.map(renderFile)}</TabsContent>
          )}
        </Tabs>
      </div>
    </div>
  )
}
