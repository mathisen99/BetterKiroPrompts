import { useState } from 'react'
import { KickoffPage } from './pages/KickoffPage'
import { SteeringPage } from './pages/SteeringPage'

function App() {
  const [page, setPage] = useState<'kickoff' | 'steering'>('kickoff')

  return (
    <div>
      <nav className="border-b border-border bg-card px-4 py-2">
        <div className="container mx-auto flex gap-4">
          <button onClick={() => setPage('kickoff')} className={`text-sm ${page === 'kickoff' ? 'text-primary font-medium' : 'text-muted-foreground'}`}>Kickoff</button>
          <button onClick={() => setPage('steering')} className={`text-sm ${page === 'steering' ? 'text-primary font-medium' : 'text-muted-foreground'}`}>Steering</button>
        </div>
      </nav>
      {page === 'kickoff' && <KickoffPage />}
      {page === 'steering' && <SteeringPage />}
    </div>
  )
}

export default App
