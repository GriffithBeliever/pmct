export type MediaType = 'movie' | 'music' | 'game'
export type MediaStatus = 'owned' | 'wishlist' | 'currently_using' | 'completed'

export interface MediaItem {
  id: string
  user_id: string
  title: string
  media_type: MediaType
  status: MediaStatus
  creator: string
  genre: string[]
  release_year?: number
  cover_url: string
  notes: string
  rating?: number
  tmdb_id?: string
  musicbrainz_id?: string
  igdb_id?: string
  metadata: Record<string, unknown>
  created_at: string
  updated_at: string
}

export interface User {
  id: string
  username: string
  email: string
  display_name: string
  bio: string
  avatar_url: string
  is_public: boolean
}

export interface Profile {
  id: string
  username: string
  display_name: string
  bio: string
  avatar_url: string
  is_public: boolean
  created_at: string
}

export interface Recommendation {
  title: string
  media_type: string
  creator: string
  reason: string
  genre: string
  release_year?: number
}

export interface MetadataResult {
  external_id: string
  title: string
  creator: string
  genres: string[]
  cover_url: string
  release_year: number
  overview: string
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
}

export interface AuthResponse {
  user: User
  token: string
}
