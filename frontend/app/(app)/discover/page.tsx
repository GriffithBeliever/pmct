'use client'

import { useState } from 'react'
import { useAuth } from '@/lib/hooks/useAuth'
import { moodDiscovery } from '@/lib/api/ai'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'

type MoodItem = { title: string; media_type: string; reason: string }
type MoodSuggestion = { title: string; media_type: string; creator: string; reason: string }

interface MoodResult {
  mood: string
  interpretation: string
  from_collection: MoodItem[]
  new_suggestions: MoodSuggestion[]
}

const MOOD_PRESETS = [
  'Something chill and relaxing',
  'Action-packed and exciting',
  'Thought-provoking and deep',
  'Fun and lighthearted',
  'Classic and timeless',
  'Something new to explore',
]

export default function DiscoverPage() {
  const { token } = useAuth()
  const [mood, setMood] = useState('')
  const [result, setResult] = useState<MoodResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleDiscover(moodText: string) {
    if (!token || !moodText.trim()) return
    setLoading(true)
    setError(null)
    setResult(null)
    try {
      const data = await moodDiscovery(token, moodText)
      setResult(data)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Mood Discovery</h1>
        <p className="text-muted-foreground mt-1">
          Tell us how you&apos;re feeling and we&apos;ll suggest something from your collection.
        </p>
      </div>

      <div className="flex gap-2">
        <Input
          placeholder="I'm in the mood for..."
          value={mood}
          onChange={(e) => setMood(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleDiscover(mood)}
        />
        <Button onClick={() => handleDiscover(mood)} loading={loading} disabled={!mood.trim()}>
          Discover
        </Button>
      </div>

      <div className="flex flex-wrap gap-2">
        {MOOD_PRESETS.map((preset) => (
          <button
            key={preset}
            onClick={() => {
              setMood(preset)
              handleDiscover(preset)
            }}
            className="px-3 py-1.5 text-sm rounded-full border hover:bg-muted transition-colors"
          >
            {preset}
          </button>
        ))}
      </div>

      {error && (
        <div className="rounded-lg border border-destructive bg-destructive/10 p-4">
          <p className="text-destructive text-sm">{error}</p>
        </div>
      )}

      {result && (
        <div className="space-y-4">
          <div className="rounded-lg border bg-card p-4">
            <p className="text-sm font-medium text-muted-foreground mb-1">Mood interpretation</p>
            <p className="text-foreground">{result.interpretation}</p>
          </div>

          {result.from_collection && result.from_collection.length > 0 && (
            <div>
              <h2 className="font-semibold mb-3">From your collection</h2>
              <ul className="space-y-2">
                {result.from_collection.map((item, i) => (
                  <li key={i} className="flex items-center gap-3 p-3 rounded-lg border bg-card">
                    <span className="text-2xl">
                      {item.media_type === 'movie' ? '🎬' : item.media_type === 'music' ? '🎵' : '🎮'}
                    </span>
                    <div>
                      <p className="font-medium">{item.title}</p>
                      {item.reason && (
                        <p className="text-sm text-muted-foreground">{item.reason}</p>
                      )}
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {result.new_suggestions && result.new_suggestions.length > 0 && (
            <div>
              <h2 className="font-semibold mb-3">New suggestions to explore</h2>
              <ul className="space-y-2">
                {result.new_suggestions.map((item, i) => (
                  <li key={i} className="flex items-center gap-3 p-3 rounded-lg border bg-card border-dashed">
                    <span className="text-2xl">
                      {item.media_type === 'movie' ? '🎬' : item.media_type === 'music' ? '🎵' : '🎮'}
                    </span>
                    <div>
                      <p className="font-medium">{item.title}</p>
                      {item.reason && (
                        <p className="text-sm text-muted-foreground">{item.reason}</p>
                      )}
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
