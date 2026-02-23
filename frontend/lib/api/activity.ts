import { apiFetch } from './client'
import type { ActivityEvent } from '@/lib/types'

export async function fetchActivity(token: string): Promise<ActivityEvent[]> {
  return apiFetch<ActivityEvent[]>('/api/activity', { token })
}
