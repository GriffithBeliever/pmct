import { apiFetch } from './client'
import type { MediaItem, MediaType, MediaStatus, PaginatedResponse } from '@/lib/types'

export interface MediaListParams {
  type?: MediaType
  status?: MediaStatus
  genre?: string
  page?: number
  page_size?: number
}

export async function fetchMediaList(
  token: string,
  params: MediaListParams = {},
): Promise<PaginatedResponse<MediaItem>> {
  const qs = new URLSearchParams()
  if (params.type) qs.set('type', params.type)
  if (params.status) qs.set('status', params.status)
  if (params.genre) qs.set('genre', params.genre)
  if (params.page) qs.set('page', String(params.page))
  if (params.page_size) qs.set('page_size', String(params.page_size))

  const query = qs.toString() ? `?${qs.toString()}` : ''
  return apiFetch<PaginatedResponse<MediaItem>>(`/api/media${query}`, { token })
}

export interface CreateMediaPayload {
  title: string
  media_type: MediaType
  status: MediaStatus
  creator?: string
  genre?: string[]
  release_year?: number
  cover_url?: string
  notes?: string
  rating?: number
  enrich_metadata?: boolean
}

export async function createMedia(
  token: string,
  payload: CreateMediaPayload,
): Promise<MediaItem> {
  return apiFetch<MediaItem>('/api/media', {
    method: 'POST',
    body: JSON.stringify(payload),
    token,
  })
}

export async function fetchMediaItem(token: string, id: string): Promise<MediaItem> {
  return apiFetch<MediaItem>(`/api/media/${id}`, { token })
}

export async function updateMedia(
  token: string,
  id: string,
  payload: Partial<CreateMediaPayload>,
): Promise<MediaItem> {
  return apiFetch<MediaItem>(`/api/media/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
    token,
  })
}

export async function deleteMedia(token: string, id: string): Promise<void> {
  return apiFetch<void>(`/api/media/${id}`, { method: 'DELETE', token })
}

export async function updateStatus(
  token: string,
  id: string,
  status: MediaStatus,
): Promise<MediaItem> {
  return apiFetch<MediaItem>(`/api/media/${id}/status`, {
    method: 'PATCH',
    body: JSON.stringify({ status }),
    token,
  })
}

export async function searchMedia(
  token: string,
  query: string,
  type?: MediaType,
): Promise<PaginatedResponse<MediaItem>> {
  const qs = new URLSearchParams({ q: query })
  if (type) qs.set('type', type)
  return apiFetch<PaginatedResponse<MediaItem>>(`/api/search?${qs.toString()}`, { token })
}
