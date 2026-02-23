'use client'

import { useEffect, useState } from 'react'
import { useSSE } from '@/lib/hooks/useSSE'
import { useAuth } from '@/lib/hooks/useAuth'
import { Button } from '@/components/ui/Button'
import { Skeleton } from '@/components/ui/Skeleton'

export default function InsightsPage() {
  const { token } = useAuth()
  const [started, setStarted] = useState(false)
  const [url, setUrl] = useState<string | null>(null)

  useEffect(() => {
    if (started && token) {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
      setUrl(`${apiUrl}/api/ai/insights?token=${token}`)
    }
  }, [started, token])

  const { content, done, error } = useSSE(url)

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">AI Insights</h1>
          <p className="text-muted-foreground mt-1">
            Get AI-powered analysis of your collection patterns and habits.
          </p>
        </div>
        {!started && (
          <Button onClick={() => setStarted(true)}>Generate Insights</Button>
        )}
      </div>

      {started && (
        <div className="rounded-lg border bg-card p-6 min-h-[300px]">
          {error ? (
            <p className="text-destructive">Error: {error}</p>
          ) : content ? (
            <div className="prose prose-sm dark:prose-invert max-w-none whitespace-pre-wrap">
              {content}
              {!done && <span className="animate-pulse">▋</span>}
            </div>
          ) : (
            <div className="space-y-3">
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-5/6" />
              <Skeleton className="h-4 w-4/6" />
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-3/4" />
            </div>
          )}
        </div>
      )}

      {!started && (
        <div className="rounded-lg border border-dashed bg-muted/30 p-12 text-center">
          <p className="text-4xl mb-4">✨</p>
          <p className="text-lg font-medium">Discover patterns in your collection</p>
          <p className="text-sm text-muted-foreground mt-1">
            Claude will analyze your collection and share interesting observations.
          </p>
        </div>
      )}
    </div>
  )
}
