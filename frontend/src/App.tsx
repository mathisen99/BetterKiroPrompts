import { useState, useCallback } from 'react'
import { LandingPage } from '@/pages/LandingPage'
import { GalleryPage } from '@/pages/GalleryPage'
import { InfoPage } from '@/pages/InfoPage'
import { SecurityScanPage } from '@/pages/SecurityScanPage'
import { NightSkyBackground } from '@/components/shared/NightSkyBackground'
import { CompactHeader } from '@/components/shared/CompactHeader'
import * as storage from '@/lib/storage'
import type { Phase } from '@/lib/storage'

type AppView = 'main' | 'gallery' | 'info' | 'scan'

function App() {
  const [currentPhase, setCurrentPhase] = useState<Phase>('level-select')
  const [currentView, setCurrentView] = useState<AppView>('main')
  const [resetKey, setResetKey] = useState(0)
  const [initialGalleryItemId, setInitialGalleryItemId] = useState<string | null>(null)
  
  const handlePhaseChange = useCallback((phase: Phase) => {
    setCurrentPhase(phase)
  }, [])

  const handleStartOver = useCallback(() => {
    // Clear storage and force LandingPage to remount with fresh state
    storage.clear()
    setResetKey(k => k + 1)
    setCurrentPhase('level-select')
    setCurrentView('main')
  }, [])

  const handleOpenGallery = useCallback(() => {
    setInitialGalleryItemId(null)
    setCurrentView('gallery')
  }, [])

  const handleViewInGallery = useCallback((generationId: string) => {
    setInitialGalleryItemId(generationId)
    setCurrentView('gallery')
  }, [])

  const handleCloseGallery = useCallback(() => {
    setInitialGalleryItemId(null)
    setCurrentView('main')
  }, [])

  const handleOpenInfo = useCallback(() => {
    setCurrentView('info')
  }, [])

  const handleCloseInfo = useCallback(() => {
    setCurrentView('main')
  }, [])

  const handleOpenScan = useCallback(() => {
    setCurrentView('scan')
  }, [])

  const handleCloseScan = useCallback(() => {
    setCurrentView('main')
  }, [])

  // Show large logo only during level-select phase on main view
  const showLargeLogo = currentPhase === 'level-select' && currentView === 'main'
  const showCompactHeader = !showLargeLogo && currentView === 'main'

  // Gallery view
  if (currentView === 'gallery') {
    return (
      <div className="min-h-screen">
        <NightSkyBackground />
        <GalleryPage onBack={handleCloseGallery} initialItemId={initialGalleryItemId} onOpenInfo={handleOpenInfo} />
      </div>
    )
  }

  // Info view
  if (currentView === 'info') {
    return (
      <div className="min-h-screen">
        <NightSkyBackground />
        <InfoPage
          onNavigateHome={handleCloseInfo}
          onNavigateGallery={handleOpenGallery}
          onNavigateScan={handleOpenScan}
        />
      </div>
    )
  }

  // Scan view
  if (currentView === 'scan') {
    return (
      <div className="min-h-screen">
        <NightSkyBackground />
        <SecurityScanPage
          onNavigateHome={handleCloseScan}
          onNavigateGallery={handleOpenGallery}
        />
      </div>
    )
  }

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
            <CompactHeader onStartOver={handleStartOver} onOpenGallery={handleOpenGallery} onOpenInfo={handleOpenInfo} />
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
        {/* Gallery and About links on landing page */}
        {showLargeLogo && (
          <div className="flex justify-center gap-4 mb-6">
            <button
              onClick={handleOpenGallery}
              className="inline-flex items-center gap-2 px-6 py-3 text-base font-medium rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors shadow-lg shadow-primary/25"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect width="18" height="18" x="3" y="3" rx="2" ry="2"/>
                <circle cx="9" cy="9" r="2"/>
                <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"/>
              </svg>
              Browse Gallery
            </button>
            <button
              onClick={handleOpenInfo}
              className="inline-flex items-center gap-2 px-6 py-3 text-base font-medium rounded-lg border border-border bg-background/50 hover:bg-accent hover:text-accent-foreground transition-colors"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <path d="M12 16v-4"/>
                <path d="M12 8h.01"/>
              </svg>
              About
            </button>
          </div>
        )}
        <LandingPage key={resetKey} onPhaseChange={handlePhaseChange} onViewInGallery={handleViewInGallery} />
      </main>
    </div>
  )
}

export default App
