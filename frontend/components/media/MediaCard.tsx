import Link from 'next/link'
import Image from 'next/image'
import { Star } from 'lucide-react'
import { StatusBadge } from './StatusBadge'
import { mediaTypeIcon } from '@/lib/utils/format'
import type { MediaItem } from '@/lib/types'

interface MediaCardProps {
  item: MediaItem
}

export function MediaCard({ item }: MediaCardProps) {
  return (
    <Link href={`/media/${item.id}`} className="block group">
      <div className="rounded-lg border border-border bg-card overflow-hidden hover:shadow-md transition-shadow">
        <div className="relative h-48 bg-muted">
          {item.cover_url ? (
            <Image
              src={item.cover_url}
              alt={item.title}
              fill
              className="object-cover group-hover:scale-105 transition-transform duration-200"
              sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
            />
          ) : (
            <div className="flex h-full items-center justify-center text-5xl">
              {mediaTypeIcon(item.media_type)}
            </div>
          )}
          <div className="absolute top-2 right-2">
            <StatusBadge status={item.status} />
          </div>
        </div>
        <div className="p-4">
          <h3 className="font-semibold truncate text-sm leading-tight">{item.title}</h3>
          {item.creator ? (
            <p className="text-xs text-muted-foreground truncate mt-0.5">{item.creator}</p>
          ) : null}
          <div className="mt-2 flex items-center justify-between">
            <div className="flex flex-wrap gap-1">
              {item.genre.slice(0, 2).map((g) => (
                <span
                  key={g}
                  className="text-xs bg-secondary text-secondary-foreground rounded px-1.5 py-0.5"
                >
                  {g}
                </span>
              ))}
            </div>
            {item.rating != null ? (
              <div className="flex items-center gap-0.5 text-xs text-yellow-500">
                <Star className="h-3 w-3 fill-current" />
                <span>{item.rating.toFixed(1)}</span>
              </div>
            ) : null}
          </div>
          {item.release_year ? (
            <p className="text-xs text-muted-foreground mt-1">{item.release_year}</p>
          ) : null}
        </div>
      </div>
    </Link>
  )
}
