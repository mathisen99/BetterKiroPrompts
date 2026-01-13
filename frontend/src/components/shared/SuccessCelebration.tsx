import { useEffect, useState } from 'react'

interface SuccessCelebrationProps {
  onComplete: () => void
  duration?: number
}

interface Particle {
  id: number
  x: number
  y: number
  color: string
  size: number
  delay: number
  borderRadius: string
  animationDuration: number
}

const COLORS = [
  'rgb(99, 102, 241)',   // Primary indigo
  'rgb(139, 92, 246)',   // Purple
  'rgb(59, 130, 246)',   // Blue
  'rgb(34, 197, 94)',    // Green
  'rgb(250, 204, 21)',   // Yellow
  'rgb(244, 114, 182)',  // Pink
]

function generateParticles(count: number): Particle[] {
  return Array.from({ length: count }, (_, i) => ({
    id: i,
    x: Math.random() * 100,
    y: Math.random() * 100,
    color: COLORS[Math.floor(Math.random() * COLORS.length)],
    size: Math.random() * 8 + 4,
    delay: Math.random() * 0.5,
    borderRadius: Math.random() > 0.5 ? '50%' : '2px',
    animationDuration: 1.5 + Math.random(),
  }))
}

export function SuccessCelebration({ onComplete, duration = 2000 }: SuccessCelebrationProps) {
  const [particles] = useState(() => generateParticles(50))
  const [isVisible, setIsVisible] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false)
      // Small delay before calling onComplete to allow fade out
      setTimeout(onComplete, 300)
    }, duration)

    return () => clearTimeout(timer)
  }, [duration, onComplete])

  return (
    <div
      className={`fixed inset-0 pointer-events-none z-50 transition-opacity duration-300 ${
        isVisible ? 'opacity-100' : 'opacity-0'
      }`}
      aria-hidden="true"
    >
      {/* Celebration message */}
      <div className="absolute inset-0 flex items-center justify-center">
        <div className="text-center animate-bounce-in">
          <div className="text-6xl mb-4">🎉</div>
          <h2 className="text-2xl font-bold text-foreground">Success!</h2>
          <p className="text-muted-foreground">Your Kiro files are ready</p>
        </div>
      </div>

      {/* Confetti particles */}
      {particles.map((particle) => (
        <div
          key={particle.id}
          className="absolute animate-confetti"
          style={{
            left: `${particle.x}%`,
            top: '-20px',
            width: `${particle.size}px`,
            height: `${particle.size}px`,
            backgroundColor: particle.color,
            borderRadius: particle.borderRadius,
            animationDelay: `${particle.delay}s`,
            animationDuration: `${particle.animationDuration}s`,
          }}
        />
      ))}
    </div>
  )
}
