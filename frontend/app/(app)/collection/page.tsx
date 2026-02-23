import { Suspense } from 'react'
import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { MediaGrid } from '@/components/media/MediaGrid'
import { CollectionFilters } from '@/components/media/CollectionFilters'
import { MediaCardSkeleton } from '@/components/ui/Skeleton'
import { fetchMediaList, searchMedia } from '@/lib/api/media'
import type { MediaType, MediaStatus } from '@/lib/types'

interface PageProps {
  searchParams: Promise<{
    type?: string
    status?: string
    q?: string
    page?: string
  }>
}

export default async function CollectionPage({ searchParams }: PageProps) {
  const cookieStore = await cookies()
  const token = cookieStore.get('ems_token')?.value
  if (!token) redirect('/login')

  const params = await searchParams

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">My Collection</h1>
      </div>
      <CollectionFilters />
      <Suspense
        key={JSON.stringify(params)}
        fallback={
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
            {Array.from({ length: 12 }).map((_, i) => (
              <MediaCardSkeleton key={i} />
            ))}
          </div>
        }
      >
        <CollectionContent token={token} params={params} />
      </Suspense>
    </div>
  )
}

async function CollectionContent({
  token,
  params,
}: {
  token: string
  params: { type?: string; status?: string; q?: string; page?: string }
}) {
  const page = params.page ? parseInt(params.page) : 1

  // Use search endpoint when query is present
  const data = params.q
    ? await searchMedia(token, params.q, params.type as MediaType | undefined)
    : await fetchMediaList(token, {
        type: params.type as MediaType | undefined,
        status: params.status as MediaStatus | undefined,
        page,
        page_size: 24,
      })

  return <MediaGrid items={data.items} />
}
