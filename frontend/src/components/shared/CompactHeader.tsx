import { Button } from '@/components/ui/button'

interface CompactHeaderProps {
  onStartOver: () => void
}

export function CompactHeader({ onStartOver }: CompactHeaderProps) {
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
        <Button
          variant="ghost"
          size="sm"
          onClick={onStartOver}
          className="text-muted-foreground hover:text-foreground"
        >
          Start Over
        </Button>
      </div>
    </header>
  )
}
