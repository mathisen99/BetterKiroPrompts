export function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="container mx-auto flex h-14 items-center px-4">
        <div className="flex items-center gap-2">
          <img
            src="/logo.png"
            alt="BetterKiroPrompts"
            className="h-8 w-8 rounded-lg"
          />
          <span className="text-lg font-semibold tracking-tight">
            BetterKiroPrompts
          </span>
        </div>
        <nav className="ml-auto flex items-center gap-4">
          <a
            href="https://github.com"
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            GitHub
          </a>
          <a
            href="https://kiro.dev"
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            Kiro Docs
          </a>
        </nav>
      </div>
    </header>
  )
}
