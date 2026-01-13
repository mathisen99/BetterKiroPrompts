import { useEffect, useRef } from 'react'

export function NightSkyBackground() {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    // Set canvas size
    const resize = () => {
      canvas.width = window.innerWidth
      canvas.height = window.innerHeight
    }
    resize()
    window.addEventListener('resize', resize)

    // Star properties
    interface Star {
      x: number
      y: number
      size: number
      speed: number
      opacity: number
      twinkleSpeed: number
      twinklePhase: number
    }

    // Create stars
    const stars: Star[] = []
    const starCount = Math.floor((canvas.width * canvas.height) / 3000)
    
    for (let i = 0; i < starCount; i++) {
      stars.push({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        size: Math.random() * 2 + 0.5,
        speed: Math.random() * 0.02 + 0.005,
        opacity: Math.random() * 0.5 + 0.5,
        twinkleSpeed: Math.random() * 0.02 + 0.01,
        twinklePhase: Math.random() * Math.PI * 2,
      })
    }

    // Shooting stars
    interface ShootingStar {
      x: number
      y: number
      length: number
      speed: number
      opacity: number
      active: boolean
    }

    const shootingStars: ShootingStar[] = []
    const maxShootingStars = 3

    const createShootingStar = () => {
      if (shootingStars.filter(s => s.active).length >= maxShootingStars) return
      if (Math.random() > 0.002) return // Rare occurrence

      shootingStars.push({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height * 0.5,
        length: Math.random() * 80 + 40,
        speed: Math.random() * 8 + 6,
        opacity: 1,
        active: true,
      })
    }

    // Animation
    let animationId: number
    let time = 0

    const animate = () => {
      time += 0.016 // ~60fps
      ctx.fillStyle = '#0a0a0f'
      ctx.fillRect(0, 0, canvas.width, canvas.height)

      // Draw gradient overlay for depth
      const gradient = ctx.createRadialGradient(
        canvas.width / 2,
        canvas.height / 2,
        0,
        canvas.width / 2,
        canvas.height / 2,
        canvas.width * 0.8
      )
      gradient.addColorStop(0, 'rgba(30, 41, 82, 0.3)')
      gradient.addColorStop(0.5, 'rgba(15, 23, 42, 0.2)')
      gradient.addColorStop(1, 'rgba(0, 0, 0, 0)')
      ctx.fillStyle = gradient
      ctx.fillRect(0, 0, canvas.width, canvas.height)

      // Draw stars
      stars.forEach(star => {
        // Twinkle effect
        const twinkle = Math.sin(time * star.twinkleSpeed + star.twinklePhase) * 0.3 + 0.7
        const currentOpacity = star.opacity * twinkle

        // Draw star glow
        const glowGradient = ctx.createRadialGradient(
          star.x, star.y, 0,
          star.x, star.y, star.size * 3
        )
        glowGradient.addColorStop(0, `rgba(147, 197, 253, ${currentOpacity * 0.8})`)
        glowGradient.addColorStop(0.5, `rgba(147, 197, 253, ${currentOpacity * 0.2})`)
        glowGradient.addColorStop(1, 'rgba(147, 197, 253, 0)')
        
        ctx.beginPath()
        ctx.arc(star.x, star.y, star.size * 3, 0, Math.PI * 2)
        ctx.fillStyle = glowGradient
        ctx.fill()

        // Draw star core
        ctx.beginPath()
        ctx.arc(star.x, star.y, star.size, 0, Math.PI * 2)
        ctx.fillStyle = `rgba(255, 255, 255, ${currentOpacity})`
        ctx.fill()

        // Slow drift
        star.y += star.speed
        if (star.y > canvas.height + 10) {
          star.y = -10
          star.x = Math.random() * canvas.width
        }
      })

      // Create and update shooting stars
      createShootingStar()
      
      shootingStars.forEach((star, index) => {
        if (!star.active) return

        // Draw shooting star trail
        const tailX = star.x - star.length * 0.7
        const tailY = star.y - star.length * 0.7

        const trailGradient = ctx.createLinearGradient(tailX, tailY, star.x, star.y)
        trailGradient.addColorStop(0, 'rgba(255, 255, 255, 0)')
        trailGradient.addColorStop(0.8, `rgba(147, 197, 253, ${star.opacity * 0.5})`)
        trailGradient.addColorStop(1, `rgba(255, 255, 255, ${star.opacity})`)

        ctx.beginPath()
        ctx.moveTo(tailX, tailY)
        ctx.lineTo(star.x, star.y)
        ctx.strokeStyle = trailGradient
        ctx.lineWidth = 2
        ctx.stroke()

        // Draw head glow
        const headGlow = ctx.createRadialGradient(star.x, star.y, 0, star.x, star.y, 4)
        headGlow.addColorStop(0, `rgba(255, 255, 255, ${star.opacity})`)
        headGlow.addColorStop(1, 'rgba(147, 197, 253, 0)')
        ctx.beginPath()
        ctx.arc(star.x, star.y, 4, 0, Math.PI * 2)
        ctx.fillStyle = headGlow
        ctx.fill()

        // Move shooting star
        star.x += star.speed
        star.y += star.speed
        star.opacity -= 0.015

        if (star.opacity <= 0 || star.x > canvas.width + 100 || star.y > canvas.height + 100) {
          shootingStars.splice(index, 1)
        }
      })

      animationId = requestAnimationFrame(animate)
    }

    animate()

    return () => {
      window.removeEventListener('resize', resize)
      cancelAnimationFrame(animationId)
    }
  }, [])

  return (
    <canvas
      ref={canvasRef}
      className="fixed inset-0 -z-10"
      aria-hidden="true"
    />
  )
}
