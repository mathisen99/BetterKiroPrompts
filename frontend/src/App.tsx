import { LandingPage } from '@/pages/LandingPage'
import { Header } from '@/components/shared/Header'

function App() {
  return (
    <div className="min-h-screen bg-background">
      <a
        href="#main-content"
        className="sr-only focus:not-sr-only focus:absolute focus:top-2 focus:left-2 focus:z-50 focus:rounded focus:bg-primary focus:px-4 focus:py-2 focus:text-primary-foreground"
      >
        Skip to main content
      </a>
      <Header />
      <main id="main-content" className="container mx-auto px-4 py-8">
        <LandingPage />
      </main>
    </div>
  )
}

export default App
