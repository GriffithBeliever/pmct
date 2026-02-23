'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Search } from 'lucide-react'
import { Input } from '@/components/ui/Input'
import { Select } from '@/components/ui/Select'
import { Button } from '@/components/ui/Button'
import { createMedia, updateMedia } from '@/lib/api/media'
import { fetchMetadataSearch } from '@/lib/api/profile'
import type { MediaItem, MediaType, MediaStatus, MetadataResult } from '@/lib/types'

interface MediaFormProps {
  token: string
  item?: MediaItem
}

const mediaTypeOptions = [
  { value: 'movie', label: 'Movie' },
  { value: 'music', label: 'Music' },
  { value: 'game', label: 'Game' },
]

const statusOptions = [
  { value: 'owned', label: 'Owned' },
  { value: 'wishlist', label: 'Wishlist' },
  { value: 'currently_using', label: 'In Progress' },
  { value: 'completed', label: 'Completed' },
]

export function MediaForm({ token, item }: MediaFormProps) {
  const router = useRouter()
  const isEdit = !!item

  const [title, setTitle] = useState(item?.title ?? '')
  const [mediaType, setMediaType] = useState<MediaType>(item?.media_type ?? 'movie')
  const [status, setStatus] = useState<MediaStatus>(item?.status ?? 'owned')
  const [creator, setCreator] = useState(item?.creator ?? '')
  const [genre, setGenre] = useState(item?.genre?.join(', ') ?? '')
  const [releaseYear, setReleaseYear] = useState(item?.release_year?.toString() ?? '')
  const [coverUrl, setCoverUrl] = useState(item?.cover_url ?? '')
  const [notes, setNotes] = useState(item?.notes ?? '')
  const [rating, setRating] = useState(item?.rating?.toString() ?? '')
  const [enriching, setEnriching] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [metaSuggestions, setMetaSuggestions] = useState<MetadataResult[]>([])

  const handleLookup = async () => {
    if (!title) return
    setEnriching(true)
    try {
      const results = (await fetchMetadataSearch(
        token,
        title,
        mediaType,
        releaseYear ? parseInt(releaseYear) : undefined,
      )) as MetadataResult[]
      setMetaSuggestions(results ?? [])
    } catch {
      // ignore
    } finally {
      setEnriching(false)
    }
  }

  const applyMeta = (result: MetadataResult) => {
    if (result.title) setTitle(result.title)
    if (result.creator) setCreator(result.creator)
    if (result.genres?.length) setGenre(result.genres.join(', '))
    if (result.cover_url) setCoverUrl(result.cover_url)
    if (result.release_year) setReleaseYear(String(result.release_year))
    setMetaSuggestions([])
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSubmitting(true)

    const payload = {
      title,
      media_type: mediaType,
      status,
      creator,
      genre: genre.split(',').map((g) => g.trim()).filter(Boolean),
      release_year: releaseYear ? parseInt(releaseYear) : undefined,
      cover_url: coverUrl,
      notes,
      rating: rating ? parseFloat(rating) : undefined,
    }

    try {
      if (isEdit) {
        await updateMedia(token, item.id, payload)
        router.push(`/media/${item.id}`)
      } else {
        const created = await createMedia(token, payload)
        router.push(`/media/${created.id}`)
      }
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-1">
          <Input
            id="title"
            label="Title *"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            placeholder="Enter title..."
          />
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleLookup}
            loading={enriching}
            disabled={!title}
            className="w-full mt-1"
          >
            <Search className="h-4 w-4 mr-1" />
            Look up metadata
          </Button>
        </div>

        <Select
          id="mediaType"
          label="Type *"
          value={mediaType}
          onChange={(e) => setMediaType(e.target.value as MediaType)}
          options={mediaTypeOptions}
        />
      </div>

      {metaSuggestions.length > 0 ? (
        <div className="rounded-md border border-border p-3 space-y-2">
          <p className="text-sm font-medium text-muted-foreground">Metadata suggestions:</p>
          {metaSuggestions.slice(0, 3).map((r, i) => (
            <button
              key={i}
              type="button"
              onClick={() => applyMeta(r)}
              className="w-full text-left flex gap-3 p-2 hover:bg-muted rounded-md transition-colors"
            >
              {r.cover_url ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img src={r.cover_url} alt={r.title} className="h-12 w-9 object-cover rounded shrink-0" />
              ) : null}
              <div>
                <p className="text-sm font-medium">{r.title}</p>
                <p className="text-xs text-muted-foreground">
                  {r.creator} {r.release_year ? `(${r.release_year})` : ''}
                </p>
              </div>
            </button>
          ))}
        </div>
      ) : null}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Select
          id="status"
          label="Status"
          value={status}
          onChange={(e) => setStatus(e.target.value as MediaStatus)}
          options={statusOptions}
        />
        <Input
          id="creator"
          label="Creator / Artist / Developer"
          value={creator}
          onChange={(e) => setCreator(e.target.value)}
          placeholder="e.g. Christopher Nolan"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Input
          id="genre"
          label="Genres (comma-separated)"
          value={genre}
          onChange={(e) => setGenre(e.target.value)}
          placeholder="e.g. Action, Thriller"
        />
        <Input
          id="releaseYear"
          label="Release Year"
          type="number"
          value={releaseYear}
          onChange={(e) => setReleaseYear(e.target.value)}
          placeholder="e.g. 2023"
          min={1900}
          max={2030}
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Input
          id="coverUrl"
          label="Cover Image URL"
          value={coverUrl}
          onChange={(e) => setCoverUrl(e.target.value)}
          placeholder="https://..."
          type="url"
        />
        <Input
          id="rating"
          label="Rating (0–10)"
          type="number"
          value={rating}
          onChange={(e) => setRating(e.target.value)}
          placeholder="e.g. 8.5"
          min={0}
          max={10}
          step={0.1}
        />
      </div>

      <div>
        <label htmlFor="notes" className="block text-sm font-medium text-foreground mb-1">
          Notes
        </label>
        <textarea
          id="notes"
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          rows={3}
          placeholder="Personal notes, thoughts..."
          className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      {error ? (
        <p className="text-sm text-red-500 bg-red-50 dark:bg-red-950/20 rounded-md p-3">{error}</p>
      ) : null}

      <div className="flex gap-3">
        <Button type="submit" loading={submitting}>
          {isEdit ? 'Save changes' : 'Add to collection'}
        </Button>
        <Button type="button" variant="outline" onClick={() => router.back()}>
          Cancel
        </Button>
      </div>
    </form>
  )
}
