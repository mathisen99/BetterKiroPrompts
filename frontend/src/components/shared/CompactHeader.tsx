import { Button } from '@/components/ui/button'
import { ImageIcon, Info } from 'lucide-react'

interface CompactHeaderProps {
  onStartOver: () => void
  onOpenGallery?: () => void
  onOpenInfo?: () => void
}

export function CompactHeader({ onStartOver, onOpenGallery, onOpenInfo }: CompactHeaderProps) {
  return (
    <header className="w-full mb-6">
      <div className="flex items-center justify-between">
        <a
          href="/"
          className="flex items-center gap-2 hover:opacity-80 transition-opacity"
          onClick={(e) => {
            e.preventDefault()
            onStartOver()
          }}
          aria-label="Go to start"
        >
          <img
            src="/logo.png"
            alt="BetterKiroPrompts"
            className="h-10 w-auto drop-shadow-[0_0_15px_rgba(99,102,241,0.4)]"
          />
        </a>
        <div className="flex items-center gap-2">
          {onOpenGallery && (
            <Button
              variant="ghost"
              size="sm"
              onClick={onOpenGallery}
              className="text-muted-foreground hover:text-foreground gap-1.5"
            >
              <ImageIcon className="h-4 w-4" />
              Gallery
            </Button>
          )}
          {onOpenInfo && (
            <Button
              variant="ghost"
              size="sm"
              onClick={onOpenInfo}
              className="text-muted-foreground hover:text-foreground gap-1.5"
            >
              <Info className="h-4 w-4" />
              About
            </Button>
          )}
          <Button
            variant="ghost"
            size="sm"
            onClick={onStartOver}
            className="text-muted-foreground hover:text-foreground"
          >
            Start Over
          </Button>
        </div>
      </div>
    </header>
  )
}
