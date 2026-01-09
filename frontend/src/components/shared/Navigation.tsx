type Page = 'kickoff' | 'steering' | 'hooks'

interface NavigationProps {
  page: Page
  setPage: (page: Page) => void
}

export function Navigation({ page, setPage }: NavigationProps) {
  return (
    <nav className="border-b border-border bg-card px-4 py-2">
      <div className="container mx-auto flex gap-4">
        <button onClick={() => setPage('kickoff')} className={`text-sm ${page === 'kickoff' ? 'text-primary font-medium' : 'text-muted-foreground'}`}>Kickoff</button>
        <button onClick={() => setPage('steering')} className={`text-sm ${page === 'steering' ? 'text-primary font-medium' : 'text-muted-foreground'}`}>Steering</button>
        <button onClick={() => setPage('hooks')} className={`text-sm ${page === 'hooks' ? 'text-primary font-medium' : 'text-muted-foreground'}`}>Hooks</button>
      </div>
    </nav>
  )
}
