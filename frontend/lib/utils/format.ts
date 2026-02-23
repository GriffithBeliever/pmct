import type { MediaStatus, MediaType } from '@/lib/types'

export function formatStatus(status: MediaStatus): string {
  const labels: Record<MediaStatus, string> = {
    owned: 'Owned',
    wishlist: 'Wishlist',
    currently_using: 'In Progress',
    completed: 'Completed',
  }
  return labels[status] ?? status
}

export function formatMediaType(type: MediaType): string {
  const labels: Record<MediaType, string> = {
    movie: 'Movie',
    music: 'Music',
    game: 'Game',
  }
  return labels[type] ?? type
}

export function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

export function statusColor(status: MediaStatus): string {
  const colors: Record<MediaStatus, string> = {
    owned: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
    wishlist: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
    currently_using: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
    completed: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300',
  }
  return colors[status] ?? 'bg-gray-100 text-gray-800'
}

export function mediaTypeIcon(type: MediaType): string {
  const icons: Record<MediaType, string> = {
    movie: '🎬',
    music: '🎵',
    game: '🎮',
  }
  return icons[type] ?? '📦'
}
