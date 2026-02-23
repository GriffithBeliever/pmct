'use client'

import { useState } from 'react'
import { updateStatus } from '@/lib/api/media'
import { formatStatus } from '@/lib/utils/format'
import type { MediaStatus } from '@/lib/types'

const STATUSES: MediaStatus[] = ['owned', 'wishlist', 'currently_using', 'completed']

interface StatusSelectProps {
  itemId: string
  currentStatus: MediaStatus
  token: string
  onUpdate?: (status: MediaStatus) => void
}

export function StatusSelect({ itemId, currentStatus, token, onUpdate }: StatusSelectProps) {
  const [status, setStatus] = useState<MediaStatus>(currentStatus)
  const [loading, setLoading] = useState(false)

  const handleChange = async (newStatus: MediaStatus) => {
    if (newStatus === status) return
    setLoading(true)
    try {
      await updateStatus(token, itemId, newStatus)
      setStatus(newStatus)
      onUpdate?.(newStatus)
    } catch {
      // revert on error
    } finally {
      setLoading(false)
    }
  }

  return (
    <select
      value={status}
      onChange={(e) => handleChange(e.target.value as MediaStatus)}
      disabled={loading}
      className="h-8 rounded-md border border-input bg-background px-2 text-xs focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50"
    >
      {STATUSES.map((s) => (
        <option key={s} value={s}>
          {formatStatus(s)}
        </option>
      ))}
    </select>
  )
}
