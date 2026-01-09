import { useState } from 'react'
import { Navigation } from './components/shared/Navigation'
import { KickoffPage } from './pages/KickoffPage'
import { SteeringPage } from './pages/SteeringPage'
import { HooksPage } from './pages/HooksPage'

function App() {
  const [page, setPage] = useState<'kickoff' | 'steering' | 'hooks'>('kickoff')

  return (
    <div>
      <a
        href="#main-content"
        className="sr-only focus:not-sr-only focus:absolute focus:top-2 focus:left-2 focus:z-50 focus:rounded focus:bg-primary focus:px-4 focus:py-2 focus:text-primary-foreground"
      >
        Skip to main content
      </a>
      <Navigation page={page} setPage={setPage} />
      <div id="main-content">
        {page === 'kickoff' && <KickoffPage />}
        {page === 'steering' && <SteeringPage />}
        {page === 'hooks' && <HooksPage />}
      </div>
    </div>
  )
}

export default App
