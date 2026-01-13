import { useState, useEffect, useCallback } from 'react'
import { ArrowLeft, Home } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { GalleryList } from '@/components/Gallery/GalleryList'
import { GalleryDetail } from '@/components/Gallery/GalleryDetail'
import { ErrorMessage } from '@/components/shared/ErrorMessage'
import {
  listGallery,
  getGalleryItem,
  rateGalleryItem,
  type GalleryItem,
  type GalleryDetail as GalleryDetailType,
  type GalleryFilters,
  ApiError,
} from '@/lib/api'
import { getVoterHashSync } from '@/lib/voter'
import { toast } from 'sonner'

interface GalleryPageProps {
  onBack: () => void
  initialItemId?: string | null
}

export function GalleryPage({ onBack, initialItemId }: GalleryPageProps) {
  const [items, setItems] = useState<GalleryItem[]>([])
  const [filters, setFilters] = useState<GalleryFilters>({
    sortBy: 'newest',
    page: 1,
  })
  const [totalPages, setTotalPages] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Detail modal state - initialize with initialItemId if provided
  const [selectedId, setSelectedId] = useState<string | null>(initialItemId ?? null)
  const [selectedGeneration, setSelectedGeneration] = useState<GalleryDetailType | null>(null)
  const [userRating, setUserRating] = useState<number | null>(null)
  const [isLoadingDetail, setIsLoadingDetail] = useState(false)
  const [isRating, setIsRating] = useState(false)

  const voterHash = getVoterHashSync()

  // Fetch gallery list
  const fetchGallery = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const response = await listGallery(filters)
      setItems(response.items)
      setTotalPages(response.totalPages)
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message)
      } else {
        setError('Failed to load gallery')
      }
    } finally {
      setIsLoading(false)
    }
  }, [filters])

  useEffect(() => {
    fetchGallery()
  }, [fetchGallery])

  // Fetch detail when item selected
  const fetchDetail = useCallback(async (id: string) => {
    setIsLoadingDetail(true)
    try {
      const response = await getGalleryItem(id, voterHash)
      setSelectedGeneration(response.generation)
      setUserRating(response.userRating || null)
    } catch (err) {
      if (err instanceof ApiError) {
        toast.error(err.message)
      } else {
        toast.error('Failed to load generation details')
      }
      setSelectedId(null)
    } finally {
      setIsLoadingDetail(false)
    }
  }, [voterHash])

  useEffect(() => {
    if (selectedId) {
      fetchDetail(selectedId)
    } else {
      setSelectedGeneration(null)
      setUserRating(null)
    }
  }, [selectedId, fetchDetail])

  const handleItemClick = (id: string) => {
    setSelectedId(id)
  }

  const handleCloseDetail = () => {
    setSelectedId(null)
    // Refresh list to update view counts
    fetchGallery()
  }

  const handleRate = async (score: number) => {
    if (!selectedId) return

    setIsRating(true)
    try {
      await rateGalleryItem(selectedId, score, voterHash)
      setUserRating(score)
      toast.success('Rating submitted!')
      // Refresh detail to get updated average
      await fetchDetail(selectedId)
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.status === 429) {
          toast.error('Too many ratings. Please wait before rating again.')
        } else {
          toast.error(err.message)
        }
      } else {
        toast.error('Failed to submit rating')
      }
    } finally {
      setIsRating(false)
    }
  }

  return (
    <div className="min-h-screen">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={onBack} aria-label="Go back">
              <ArrowLeft className="h-5 w-5" />
            </Button>
            <div>
              <h1 className="text-2xl font-bold">Gallery</h1>
              <p className="text-muted-foreground">
                Browse community generations for inspiration
              </p>
            </div>
          </div>
          <Button
            variant="outline"
            onClick={onBack}
            className="gap-2"
          >
            <Home className="h-4 w-4" />
            Back to Home
          </Button>
        </div>

        {/* Error state */}
        {error && (
          <ErrorMessage
            message={error}
            onRetry={fetchGallery}
            onStartOver={onBack}
          />
        )}

        {/* Gallery list */}
        {!error && (
          <GalleryList
            items={items}
            filters={filters}
            onFilterChange={setFilters}
            onItemClick={handleItemClick}
            totalPages={totalPages}
            isLoading={isLoading}
          />
        )}

        {/* Detail modal */}
        {selectedId && selectedGeneration && !isLoadingDetail && (
          <GalleryDetail
            generation={selectedGeneration}
            onClose={handleCloseDetail}
            onRate={handleRate}
            userRating={userRating}
            isRating={isRating}
          />
        )}

        {/* Loading detail overlay */}
        {selectedId && isLoadingDetail && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/80 backdrop-blur-sm">
            <div className="text-center">
              <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent mx-auto mb-4" />
              <p className="text-muted-foreground">Loading...</p>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
