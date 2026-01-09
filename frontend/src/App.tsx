import { useState } from 'react'
import { Navigation } from './components/shared/Navigation'
import { KickoffPage } from './pages/KickoffPage'
import { SteeringPage } from './pages/SteeringPage'
import { HooksPage } from './pages/HooksPage'

function App() {
  const [page, setPage] = useState<'kickoff' | 'steering' | 'hooks'>('kickoff')

  return (
    <div>
      <Navigation page={page} setPage={setPage} />
      {page === 'kickoff' && <KickoffPage />}
      {page === 'steering' && <SteeringPage />}
      {page === 'hooks' && <HooksPage />}
    </div>
  )
}

export default App
