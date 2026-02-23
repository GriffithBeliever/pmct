import { Suspense } from 'react'
import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { fetchActivity } from '@/lib/api/activity'
import type { ActivityEvent } from '@/lib/types'

const EVENT_LABELS: Record<string, string> = {
  item_added: 'Added',
  item_updated: 'Updated',
  item_deleted: 'Deleted',
  status_changed: 'Status changed',
  rating_updated: 'Rating updated',
}

function EventRow({ event }: { event: ActivityEvent }) {
  const label = EVENT_LABELS[event.event_type] ?? event.event_type
  const title = (event.payload.title as string) ?? 'Unknown item'
  const detail =
    event.event_type === 'status_changed'
      ? `→ ${event.payload.new_status}`
      : event.event_type === 'rating_updated'
        ? `${event.payload.rating}/10`
        : null

  return (
    <div className="flex items-start gap-4 py-3 border-b border-border last:border-0">
      <div className="mt-0.5 h-2 w-2 rounded-full bg-primary shrink-0" />
      <div className="flex-1 min-w-0">
        <p className="text-sm">
          <span className="font-medium">{label}</span>
          {' — '}
          <span className="text-muted-foreground">{title}</span>
          {detail ? <span className="text-muted-foreground"> {detail}</span> : null}
        </p>
        <p className="text-xs text-muted-foreground mt-0.5">
          {new Date(event.created_at).toLocaleString()}
        </p>
      </div>
    </div>
  )
}

async function ActivityFeed({ token }: { token: string }) {
  const events = await fetchActivity(token)

  if (events.length === 0) {
    return (
      <p className="text-muted-foreground text-sm py-8 text-center">
        No activity yet. Start adding items to your collection.
      </p>
    )
  }

  return (
    <div className="rounded-lg border border-border bg-card p-4">
      {events.map((e) => (
        <EventRow key={e.id} event={e} />
      ))}
    </div>
  )
}

export default async function ActivityPage() {
  const cookieStore = await cookies()
  const token = cookieStore.get('ems_token')?.value
  if (!token) redirect('/login')

  return (
    <div className="space-y-6 max-w-2xl">
      <h1 className="text-2xl font-bold">Activity</h1>
      <Suspense
        fallback={
          <div className="rounded-lg border border-border bg-card p-4 space-y-3">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="flex items-start gap-4 py-3 border-b border-border last:border-0">
                <div className="mt-1 h-2 w-2 rounded-full bg-muted shrink-0" />
                <div className="flex-1 space-y-1.5">
                  <div className="h-3 w-48 rounded bg-muted animate-pulse" />
                  <div className="h-2.5 w-24 rounded bg-muted animate-pulse" />
                </div>
              </div>
            ))}
          </div>
        }
      >
        <ActivityFeed token={token} />
      </Suspense>
    </div>
  )
}
