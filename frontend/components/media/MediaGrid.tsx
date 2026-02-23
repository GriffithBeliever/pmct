import { MediaCard } from './MediaCard'
import { MediaCardSkeleton } from '@/components/ui/Skeleton'
import type { MediaItem } from '@/lib/types'

interface MediaGridProps {
  items: MediaItem[]
  loading?: boolean
}

export function MediaGrid({ items, loading }: MediaGridProps) {
  const safeItems = items ?? []

  if (loading) {
    return (
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 media-grid">
        {Array.from({ length: 10 }).map((_, i) => (
          <MediaCardSkeleton key={i} />
        ))}
      </div>
    )
  }

  if (safeItems.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-center">
        <p className="text-4xl mb-4">📦</p>
        <p className="text-lg font-medium text-muted-foreground">No items yet</p>
        <p className="text-sm text-muted-foreground mt-1">
          Add your first movie, album, or game to get started.
        </p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 media-grid">
      {safeItems.map((item) => (
        <MediaCard key={item.id} item={item} />
      ))}
    </div>
  )
}
