import { LandingPage } from '@/pages/LandingPage'
import { NightSkyBackground } from '@/components/shared/NightSkyBackground'

function App() {
  return (
    <div className="min-h-screen">
      <NightSkyBackground />
      <a
        href="#main-content"
        className="sr-only focus:not-sr-only focus:absolute focus:top-2 focus:left-2 focus:z-50 focus:rounded focus:bg-primary focus:px-4 focus:py-2 focus:text-primary-foreground"
      >
        Skip to main content
      </a>
      <main id="main-content" className="container mx-auto px-4 py-12">
        {/* Big centered logo */}
        <div className="flex justify-center mb-6">
          <img
            src="/logo.png"
            alt="BetterKiroPrompts"
            className="drop-shadow-[0_0_35px_rgba(99,102,241,0.5)]"
          />
        </div>
        <LandingPage />
      </main>
    </div>
  )
}

export default App
