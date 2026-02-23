import { Suspense } from 'react'
import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { searchMedia } from '@/lib/api/media'
import { MediaGrid } from '@/components/media/MediaGrid'
import { MediaCardSkeleton } from '@/components/ui/Skeleton'
import type { MediaType } from '@/lib/types'

interface PageProps {
  searchParams: Promise<{ q?: string; type?: string }>
}

export default async function SearchPage({ searchParams }: PageProps) {
  const cookieStore = await cookies()
  const token = cookieStore.get('ems_token')?.value
  if (!token) redirect('/login')

  const params = await searchParams
  const query = params.q ?? ''

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Search Results</h1>
      {query ? (
        <p className="text-muted-foreground">
          Results for &ldquo;<span className="text-foreground font-medium">{query}</span>&rdquo;
        </p>
      ) : (
        <p className="text-muted-foreground">Enter a search query in the header to find items.</p>
      )}
      {query && (
        <Suspense
          key={query}
          fallback={
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
              {Array.from({ length: 8 }).map((_, i) => (
                <MediaCardSkeleton key={i} />
              ))}
            </div>
          }
        >
          <SearchResults query={query} type={params.type} token={token} />
        </Suspense>
      )}
    </div>
  )
}

async function SearchResults({
  query,
  type,
  token,
}: {
  query: string
  type?: string
  token: string
}) {
  const data = await searchMedia(token, query, type as MediaType | undefined)
  return <MediaGrid items={data.items} />
}
