'use client'

import { useState, useEffect } from 'react'
import { fetchRecommendations } from '@/lib/api/ai'
import { Skeleton } from '@/components/ui/Skeleton'
import type { Recommendation } from '@/lib/types'

interface Props {
  token: string
}

export default function RecommendationPanel({ token }: Props) {
  const [recs, setRecs] = useState<Recommendation[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchRecommendations(token)
      .then(setRecs)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [token])

  if (loading) {
    return (
      <div className="space-y-3">
        <Skeleton className="h-6 w-48" />
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={i} className="h-20 w-full" />
        ))}
      </div>
    )
  }

  if (error) {
    return (
      <div className="rounded-lg border border-dashed p-6 text-center text-muted-foreground text-sm">
        Could not load recommendations.
      </div>
    )
  }

  if (recs.length === 0) return null

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold">You might also like</h2>
      <ul className="grid gap-3 sm:grid-cols-2">
        {recs.map((rec, i) => (
          <li key={i} className="rounded-lg border bg-card p-4 space-y-1">
            <div className="flex items-center gap-2">
              <span className="text-xl">
                {rec.media_type === 'movie' ? '🎬' : rec.media_type === 'music' ? '🎵' : '🎮'}
              </span>
              <p className="font-medium leading-tight">{rec.title}</p>
            </div>
            {rec.creator && (
              <p className="text-sm text-muted-foreground">{rec.creator}</p>
            )}
            {rec.reason && (
              <p className="text-xs text-muted-foreground italic">{rec.reason}</p>
            )}
          </li>
        ))}
      </ul>
    </div>
  )
}
