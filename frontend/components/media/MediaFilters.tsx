'use client'

import { Select } from '@/components/ui/Select'
import { Button } from '@/components/ui/Button'
import type { MediaType, MediaStatus } from '@/lib/types'
import type { MediaFilters } from '@/lib/hooks/useMediaFilters'

interface MediaFiltersProps {
  filters: MediaFilters
  onTypeChange: (type: MediaType | undefined) => void
  onStatusChange: (status: MediaStatus | undefined) => void
  onReset: () => void
}

const typeOptions = [
  { value: 'movie', label: 'Movies' },
  { value: 'music', label: 'Music' },
  { value: 'game', label: 'Games' },
]

const statusOptions = [
  { value: 'owned', label: 'Owned' },
  { value: 'wishlist', label: 'Wishlist' },
  { value: 'currently_using', label: 'In Progress' },
  { value: 'completed', label: 'Completed' },
]

export function MediaFiltersBar({ filters, onTypeChange, onStatusChange, onReset }: MediaFiltersProps) {
  return (
    <div className="flex flex-wrap gap-3 items-center">
      <Select
        options={typeOptions}
        placeholder="All types"
        value={filters.type ?? ''}
        onChange={(e) => onTypeChange(e.target.value as MediaType || undefined)}
        className="w-36"
      />
      <Select
        options={statusOptions}
        placeholder="All statuses"
        value={filters.status ?? ''}
        onChange={(e) => onStatusChange(e.target.value as MediaStatus || undefined)}
        className="w-40"
      />
      {(filters.type ?? filters.status ?? filters.genre) ? (
        <Button variant="outline" size="sm" onClick={onReset}>
          Clear filters
        </Button>
      ) : null}
    </div>
  )
}
