'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { nlSearch } from '@/lib/api/ai'

interface Props {
  token: string
}

export function NaturalSearchBar({ token }: Props) {
  const router = useRouter()
  const [query, setQuery] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    if (!query.trim() || !token) return

    setLoading(true)
    try {
      const result = await nlSearch(token, query)
      // Apply filters from NL parsing as query params
      const qs = new URLSearchParams()
      if (result.filters.type) qs.set('type', String(result.filters.type))
      if (result.filters.status) qs.set('status', String(result.filters.status))
      if (result.filters.q) qs.set('q', String(result.filters.q))
      else qs.set('q', query)
      router.push(`/collection?${qs.toString()}`)
    } catch {
      // Fall back to plain text search
      router.push(`/search?q=${encodeURIComponent(query)}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSearch} className="relative flex items-center">
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search naturally... e.g. 'action movies from the 90s'"
        className="w-full px-4 py-2 pr-10 text-sm rounded-md border border-input bg-background focus:outline-none focus:ring-2 focus:ring-ring"
        disabled={loading}
      />
      <button
        type="submit"
        disabled={loading || !query.trim()}
        className="absolute right-2 text-muted-foreground hover:text-foreground disabled:opacity-50"
      >
        {loading ? (
          <span className="animate-spin text-xs">⟳</span>
        ) : (
          <span>⏎</span>
        )}
      </button>
    </form>
  )
}
