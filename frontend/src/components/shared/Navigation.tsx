type Page = 'kickoff' | 'steering' | 'hooks'

interface NavigationProps {
  page: Page
  setPage: (page: Page) => void
}

export function Navigation({ page, setPage }: NavigationProps) {
  const baseClass = "rounded px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"

  return (
    <nav className="border-b border-border bg-card px-4 py-2" aria-label="Main navigation">
      <div className="container mx-auto flex gap-2">
        <button
          onClick={() => setPage('kickoff')}
          className={`${baseClass} ${page === 'kickoff' ? 'text-primary font-medium' : 'text-muted-foreground hover:text-foreground'}`}
          aria-current={page === 'kickoff' ? 'page' : undefined}
        >
          Kickoff
        </button>
        <button
          onClick={() => setPage('steering')}
          className={`${baseClass} ${page === 'steering' ? 'text-primary font-medium' : 'text-muted-foreground hover:text-foreground'}`}
          aria-current={page === 'steering' ? 'page' : undefined}
        >
          Steering
        </button>
        <button
          onClick={() => setPage('hooks')}
          className={`${baseClass} ${page === 'hooks' ? 'text-primary font-medium' : 'text-muted-foreground hover:text-foreground'}`}
          aria-current={page === 'hooks' ? 'page' : undefined}
        >
          Hooks
        </button>
      </div>
    </nav>
  )
}
