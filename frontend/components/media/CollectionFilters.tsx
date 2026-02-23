'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { useTransition } from 'react'
import { Button } from '@/components/ui/Button'
import type { MediaType, MediaStatus } from '@/lib/types'

const typeOptions: { value: MediaType; label: string }[] = [
  { value: 'movie', label: 'Movies' },
  { value: 'music', label: 'Music' },
  { value: 'game', label: 'Games' },
]

const statusOptions: { value: MediaStatus; label: string }[] = [
  { value: 'owned', label: 'Owned' },
  { value: 'wishlist', label: 'Wishlist' },
  { value: 'currently_using', label: 'In Progress' },
  { value: 'completed', label: 'Completed' },
]

export function CollectionFilters() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [, startTransition] = useTransition()

  const currentType = searchParams.get('type') ?? ''
  const currentStatus = searchParams.get('status') ?? ''
  const hasFilters = currentType || currentStatus

  function updateParam(key: string, value: string) {
    const params = new URLSearchParams(searchParams.toString())
    if (value) {
      params.set(key, value)
    } else {
      params.delete(key)
    }
    params.delete('page')
    startTransition(() => router.push(`/collection?${params.toString()}`))
  }

  function clearFilters() {
    startTransition(() => router.push('/collection'))
  }

  return (
    <div className="flex flex-wrap gap-3 items-center">
      <select
        value={currentType}
        onChange={(e) => updateParam('type', e.target.value)}
        className="h-9 rounded-md border border-input bg-background px-3 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
      >
        <option value="">All types</option>
        {typeOptions.map((o) => (
          <option key={o.value} value={o.value}>{o.label}</option>
        ))}
      </select>

      <select
        value={currentStatus}
        onChange={(e) => updateParam('status', e.target.value)}
        className="h-9 rounded-md border border-input bg-background px-3 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
      >
        <option value="">All statuses</option>
        {statusOptions.map((o) => (
          <option key={o.value} value={o.value}>{o.label}</option>
        ))}
      </select>

      {hasFilters ? (
        <Button variant="outline" size="sm" onClick={clearFilters}>
          Clear filters
        </Button>
      ) : null}
    </div>
  )
}
