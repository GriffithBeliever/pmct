'use client'

import { useRouter } from 'next/navigation'

const MOODS = [
  { emoji: '😌', label: 'Chill' },
  { emoji: '⚡', label: 'Exciting' },
  { emoji: '🧠', label: 'Thoughtful' },
  { emoji: '😂', label: 'Fun' },
  { emoji: '🎯', label: 'Classic' },
  { emoji: '🌟', label: 'Discover' },
]

export function MoodPicker() {
  const router = useRouter()

  return (
    <div className="flex flex-wrap gap-2">
      {MOODS.map(({ emoji, label }) => (
        <button
          key={label}
          onClick={() => router.push(`/discover?mood=${encodeURIComponent(label.toLowerCase())}`)}
          className="flex items-center gap-1.5 px-3 py-1.5 text-sm rounded-full border hover:bg-muted transition-colors"
        >
          <span>{emoji}</span>
          <span>{label}</span>
        </button>
      ))}
    </div>
  )
}
