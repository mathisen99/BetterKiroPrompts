import { useState, useCallback } from 'react'
import { LandingPage } from '@/pages/LandingPage'
import { NightSkyBackground } from '@/components/shared/NightSkyBackground'
import { CompactHeader } from '@/components/shared/CompactHeader'
import * as storage from '@/lib/storage'
import type { Phase } from '@/lib/storage'

function App() {
  const [currentPhase, setCurrentPhase] = useState<Phase>('level-select')
  const [resetKey, setResetKey] = useState(0)
  
  const handlePhaseChange = useCallback((phase: Phase) => {
    setCurrentPhase(phase)
  }, [])

  const handleStartOver = useCallback(() => {
    // Clear storage and force LandingPage to remount with fresh state
    storage.clear()
    setResetKey(k => k + 1)
    setCurrentPhase('level-select')
  }, [])

  // Show large logo only during level-select phase
  const showLargeLogo = currentPhase === 'level-select'
  const showCompactHeader = !showLargeLogo

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
        {/* Compact header for non-landing phases */}
        {showCompactHeader && (
          <div className="max-w-3xl mx-auto">
            <CompactHeader onStartOver={handleStartOver} />
          </div>
        )}
        {/* Big centered logo with fade-out animation */}
        <div
          className={`flex justify-center mb-6 transition-all duration-500 ease-out ${
            showLargeLogo
              ? 'opacity-100 max-h-96 scale-100'
              : 'opacity-0 max-h-0 scale-95 overflow-hidden'
          }`}
          aria-hidden={!showLargeLogo}
        >
          <img
            src="/logo.png"
            alt="BetterKiroPrompts"
            className="drop-shadow-[0_0_35px_rgba(99,102,241,0.5)]"
          />
        </div>
        <LandingPage key={resetKey} onPhaseChange={handlePhaseChange} />
      </main>
    </div>
  )
}

export default App
