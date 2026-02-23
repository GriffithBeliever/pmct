import { Suspense } from 'react'
import { cookies } from 'next/headers'
import { redirect, notFound } from 'next/navigation'
import Image from 'next/image'
import Link from 'next/link'
import { fetchMediaItem } from '@/lib/api/media'
import { StatusBadge } from '@/components/media/StatusBadge'
import { StatusSelect } from '@/components/media/StatusSelect'
import { Button } from '@/components/ui/Button'
import { Skeleton } from '@/components/ui/Skeleton'
import dynamic from 'next/dynamic'
import { formatDate, formatMediaType } from '@/lib/utils/format'

const RecommendationPanel = dynamic(
  () => import('@/components/ai/RecommendationPanel'),
  { ssr: false, loading: () => <Skeleton className="h-64 w-full" /> },
)

interface PageProps {
  params: Promise<{ id: string }>
}

export default async function MediaDetailPage({ params }: PageProps) {
  const cookieStore = await cookies()
  const token = cookieStore.get('ems_token')?.value
  if (!token) redirect('/login')

  const { id } = await params

  const item = await fetchMediaItem(token, id).catch(() => null)
  if (!item) notFound()

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <div className="flex gap-6">
        {item.cover_url ? (
          <div className="relative w-40 h-56 flex-shrink-0 rounded-lg overflow-hidden">
            <Image
              src={item.cover_url}
              alt={item.title}
              fill
              className="object-cover"
            />
          </div>
        ) : (
          <div className="w-40 h-56 flex-shrink-0 rounded-lg bg-muted flex items-center justify-center text-5xl">
            {item.media_type === 'movie' ? '🎬' : item.media_type === 'music' ? '🎵' : '🎮'}
          </div>
        )}

        <div className="flex-1 space-y-3">
          <div className="flex items-start justify-between gap-4">
            <h1 className="text-3xl font-bold">{item.title}</h1>
            <div className="flex gap-2 flex-shrink-0">
              <Link href={`/media/${id}/edit`}>
                <Button variant="outline" size="sm">Edit</Button>
              </Link>
            </div>
          </div>

          {item.creator && (
            <p className="text-lg text-muted-foreground">{item.creator}</p>
          )}

          <div className="flex items-center gap-3 flex-wrap">
            <StatusBadge status={item.status} />
            <span className="text-sm text-muted-foreground capitalize">
              {formatMediaType(item.media_type)}
            </span>
            {item.release_year && (
              <span className="text-sm text-muted-foreground">{item.release_year}</span>
            )}
            {item.rating && (
              <span className="text-sm font-medium">⭐ {item.rating}/10</span>
            )}
          </div>

          {item.genre && item.genre.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {item.genre.map((g) => (
                <span
                  key={g}
                  className="px-2 py-0.5 text-xs rounded-full bg-muted text-muted-foreground"
                >
                  {g}
                </span>
              ))}
            </div>
          )}

          <div className="pt-2">
            <StatusSelect itemId={item.id} currentStatus={item.status} token={token} />
          </div>

          {item.notes && (
            <p className="text-sm text-muted-foreground whitespace-pre-wrap pt-2">
              {item.notes}
            </p>
          )}

          <p className="text-xs text-muted-foreground pt-2">
            Added {formatDate(item.created_at)}
          </p>
        </div>
      </div>

      <RecommendationPanel token={token} />
    </div>
  )
}
