import { Star, Eye, Clock, ChevronLeft, ChevronRight } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import type { GalleryItem, GalleryFilters } from '@/lib/api'

const CATEGORIES = [
  { id: undefined, name: 'All' },
  { id: 1, name: 'API' },
  { id: 2, name: 'CLI' },
  { id: 3, name: 'Web App' },
  { id: 4, name: 'Mobile' },
  { id: 5, name: 'Other' },
]

const SORT_OPTIONS: { value: GalleryFilters['sortBy']; label: string }[] = [
  { value: 'newest', label: 'Newest' },
  { value: 'highest_rated', label: 'Highest Rated' },
  { value: 'most_viewed', label: 'Most Viewed' },
]

interface GalleryListProps {
  items: GalleryItem[]
  filters: GalleryFilters
  onFilterChange: (filters: GalleryFilters) => void
  onItemClick: (id: string) => void
  totalPages: number
  isLoading?: boolean
}

function GalleryItemCard({ item, onClick }: { item: GalleryItem; onClick: () => void }) {
  return (
    <Card
      className="cursor-pointer transition-all hover:border-primary/50 hover:shadow-md"
      onClick={onClick}
    >
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between gap-2">
          <CardTitle className="line-clamp-2 text-base">{item.projectIdea}</CardTitle>
          <span className="shrink-0 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
            {item.category}
          </span>
        </div>
        <CardDescription className="line-clamp-2">{item.preview}</CardDescription>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="flex items-center gap-4 text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <Star className="size-4 fill-yellow-500 text-yellow-500" />
            <span>{item.avgRating.toFixed(1)}</span>
            <span className="text-xs">({item.ratingCount})</span>
          </div>
          <div className="flex items-center gap-1">
            <Eye className="size-4" />
            <span>{item.viewCount}</span>
          </div>
          <div className="ml-auto flex items-center gap-1">
            <Clock className="size-4" />
            <span>{formatDate(item.createdAt)}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function GalleryListSkeleton() {
  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {Array.from({ length: 6 }).map((_, i) => (
        <Card key={i}>
          <CardHeader className="pb-2">
            <Skeleton className="h-5 w-3/4" />
            <Skeleton className="h-4 w-full" />
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-4">
              <Skeleton className="h-4 w-16" />
              <Skeleton className="h-4 w-12" />
              <Skeleton className="ml-auto h-4 w-20" />
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return 'Yesterday'
  if (diffDays < 7) return `${diffDays}d ago`
  if (diffDays < 30) return `${Math.floor(diffDays / 7)}w ago`
  return date.toLocaleDateString()
}

export function GalleryList({
  items,
  filters,
  onFilterChange,
  onItemClick,
  totalPages,
  isLoading,
}: GalleryListProps) {
  const handleCategoryChange = (categoryId: number | undefined) => {
    onFilterChange({ ...filters, category: categoryId, page: 1 })
  }

  const handleSortChange = (sortBy: GalleryFilters['sortBy']) => {
    onFilterChange({ ...filters, sortBy, page: 1 })
  }

  const handlePageChange = (page: number) => {
    onFilterChange({ ...filters, page })
  }

  return (
    <div className="space-y-6">
      {/* Filters */}
      <div className="flex flex-wrap items-center gap-4">
        {/* Category filter */}
        <div className="flex flex-wrap gap-2">
          {CATEGORIES.map((cat) => (
            <Button
              key={cat.name}
              variant={filters.category === cat.id ? 'default' : 'outline'}
              size="sm"
              onClick={() => handleCategoryChange(cat.id)}
            >
              {cat.name}
            </Button>
          ))}
        </div>

        {/* Sort dropdown */}
        <div className="ml-auto flex items-center gap-2">
          <span className="text-sm text-muted-foreground">Sort by:</span>
          <select
            value={filters.sortBy}
            onChange={(e) => handleSortChange(e.target.value as GalleryFilters['sortBy'])}
            className="rounded-md border bg-background px-3 py-1.5 text-sm"
          >
            {SORT_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Items grid */}
      {isLoading ? (
        <GalleryListSkeleton />
      ) : items.length === 0 ? (
        <div className="py-12 text-center text-muted-foreground">
          No generations found. Be the first to create one!
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {items.map((item) => (
            <GalleryItemCard key={item.id} item={item} onClick={() => onItemClick(item.id)} />
          ))}
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="icon-sm"
            onClick={() => handlePageChange(filters.page - 1)}
            disabled={filters.page <= 1}
          >
            <ChevronLeft className="size-4" />
          </Button>
          <span className="px-4 text-sm">
            Page {filters.page} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="icon-sm"
            onClick={() => handlePageChange(filters.page + 1)}
            disabled={filters.page >= totalPages}
          >
            <ChevronRight className="size-4" />
          </Button>
        </div>
      )}
    </div>
  )
}
