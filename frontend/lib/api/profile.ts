import { apiFetch } from './client'
import type { Profile, MediaItem } from '@/lib/types'

export interface PublicProfileResponse {
  profile: Profile
  items: MediaItem[]
}

export async function fetchPublicProfile(username: string): Promise<PublicProfileResponse> {
  return apiFetch<PublicProfileResponse>(`/api/profile/${username}`)
}

export interface UpdateProfilePayload {
  display_name?: string
  bio?: string
  avatar_url?: string
  is_public?: boolean
}

export async function updateProfile(
  token: string,
  payload: UpdateProfilePayload,
): Promise<Profile> {
  return apiFetch<Profile>('/api/profile', {
    method: 'PUT',
    body: JSON.stringify(payload),
    token,
  })
}

export async function fetchMetadataSearch(
  token: string,
  title: string,
  media_type: string,
  year?: number,
) {
  return apiFetch('/api/metadata/search', {
    method: 'POST',
    body: JSON.stringify({ title, media_type, year }),
    token,
  })
}
