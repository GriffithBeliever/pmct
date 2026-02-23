'use client'

import { useState, useTransition } from 'react'
import type { MediaType, MediaStatus } from '@/lib/types'

export interface MediaFilters {
  type?: MediaType
  status?: MediaStatus
  genre?: string
  page: number
}

export function useMediaFilters() {
  const [filters, setFilters] = useState<MediaFilters>({ page: 1 })
  const [isPending, startTransition] = useTransition()

  const setType = (type: MediaType | undefined) => {
    startTransition(() => setFilters((f) => ({ ...f, type, page: 1 })))
  }

  const setStatus = (status: MediaStatus | undefined) => {
    startTransition(() => setFilters((f) => ({ ...f, status, page: 1 })))
  }

  const setGenre = (genre: string | undefined) => {
    startTransition(() => setFilters((f) => ({ ...f, genre, page: 1 })))
  }

  const setPage = (page: number) => {
    startTransition(() => setFilters((f) => ({ ...f, page })))
  }

  const reset = () => {
    startTransition(() => setFilters({ page: 1 }))
  }

  return { filters, setType, setStatus, setGenre, setPage, reset, isPending }
}
