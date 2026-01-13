export function NightSkyBackground() {
  return (
    <div className="fixed inset-0 -z-10 overflow-hidden">
      {/* YouTube video background */}
      <iframe
        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[180vw] h-[180vh] min-w-[180vw] min-h-[180vh] pointer-events-none"
        src="https://www.youtube.com/embed/kSZddHca0ME?autoplay=1&mute=1&loop=1&playlist=kSZddHca0ME&controls=0&showinfo=0&rel=0&modestbranding=1&playsinline=1&disablekb=1"
        title="Background"
        allow="autoplay; encrypted-media"
        allowFullScreen
      />
      {/* Dark overlay to blend with theme */}
      <div className="absolute inset-0 bg-black/40" />
    </div>
  )
}
