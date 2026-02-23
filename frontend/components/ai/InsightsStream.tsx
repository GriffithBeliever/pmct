'use client'

import { useSSE } from '@/lib/hooks/useSSE'
import { Skeleton } from '@/components/ui/Skeleton'

interface Props {
  url: string | null
}

export function InsightsStream({ url }: Props) {
  const { content, done, error } = useSSE(url)

  if (error) {
    return <p className="text-destructive text-sm">Error: {error}</p>
  }

  if (!content) {
    return (
      <div className="space-y-3">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-5/6" />
        <Skeleton className="h-4 w-4/6" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
      </div>
    )
  }

  return (
    <div className="prose prose-sm dark:prose-invert max-w-none whitespace-pre-wrap">
      {content}
      {!done && <span className="animate-pulse">▋</span>}
    </div>
  )
}
