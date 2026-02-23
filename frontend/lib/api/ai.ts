import { apiFetch } from './client'
import type { Recommendation } from '@/lib/types'

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

export async function fetchRecommendations(token: string): Promise<Recommendation[]> {
  return apiFetch<Recommendation[]>('/api/ai/recommendations', { token })
}

export async function nlSearch(
  token: string,
  query: string,
): Promise<{ query: string; filters: Record<string, unknown> }> {
  return apiFetch('/api/ai/nl-search', {
    method: 'POST',
    body: JSON.stringify({ query }),
    token,
  })
}

export async function moodDiscovery(
  token: string,
  mood: string,
): Promise<{
  mood: string
  interpretation: string
  from_collection: Array<{ title: string; media_type: string; reason: string }>
  new_suggestions: Array<{ title: string; media_type: string; creator: string; reason: string }>
}> {
  return apiFetch('/api/ai/mood', {
    method: 'POST',
    body: JSON.stringify({ mood }),
    token,
  })
}

export async function detectDuplicates(
  token: string,
  title: string,
  media_type: string,
  creator: string,
): Promise<{ is_duplicate: boolean; reason?: string; match_title?: string }> {
  return apiFetch('/api/ai/duplicates', {
    method: 'POST',
    body: JSON.stringify({ title, media_type, creator }),
    token,
  })
}

export function createInsightsStream(token: string): EventSource {
  // EventSource doesn't support custom headers; pass token as query param
  return new EventSource(`${API_URL}/api/ai/insights?token=${encodeURIComponent(token)}`)
}
