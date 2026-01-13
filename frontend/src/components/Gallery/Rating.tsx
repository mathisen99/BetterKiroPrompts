import { useState } from 'react'
import { Star } from 'lucide-react'
import { cn } from '@/lib/utils'

interface RatingProps {
  value: number
  count: number
  userRating: number | null
  onRate: (score: number) => void
  disabled?: boolean
}

export function Rating({ value, count, userRating, onRate, disabled }: RatingProps) {
  const [hoverRating, setHoverRating] = useState<number | null>(null)

  const displayRating = hoverRating ?? userRating ?? 0

  return (
    <div className="flex items-center gap-3">
      <div
        className="flex gap-0.5"
        onMouseLeave={() => setHoverRating(null)}
      >
        {[1, 2, 3, 4, 5].map((star) => (
          <button
            key={star}
            type="button"
            disabled={disabled}
            className={cn(
              'p-0.5 transition-transform hover:scale-110 focus:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded',
              disabled && 'cursor-not-allowed opacity-50'
            )}
            onMouseEnter={() => !disabled && setHoverRating(star)}
            onClick={() => !disabled && onRate(star)}
            aria-label={`Rate ${star} star${star > 1 ? 's' : ''}`}
          >
            <Star
              className={cn(
                'h-6 w-6 transition-colors',
                star <= displayRating
                  ? 'fill-yellow-500 text-yellow-500'
                  : 'fill-transparent text-muted-foreground hover:text-yellow-500/50'
              )}
            />
          </button>
        ))}
      </div>
      <div className="text-sm text-muted-foreground">
        <span className="font-medium">{value.toFixed(1)}</span>
        <span className="mx-1">·</span>
        <span>{count} {count === 1 ? 'rating' : 'ratings'}</span>
        {userRating && (
          <>
            <span className="mx-1">·</span>
            <span className="text-primary">Your rating: {userRating}</span>
          </>
        )}
      </div>
    </div>
  )
}
