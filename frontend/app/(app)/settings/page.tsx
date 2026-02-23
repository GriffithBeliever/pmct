'use client'

import { useState, useEffect } from 'react'
import { useAuth } from '@/lib/hooks/useAuth'
import { updateProfile } from '@/lib/api/profile'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'

export default function SettingsPage() {
  const { user, token } = useAuth()
  const [displayName, setDisplayName] = useState('')
  const [bio, setBio] = useState('')
  const [avatarUrl, setAvatarUrl] = useState('')
  const [isPublic, setIsPublic] = useState(true)
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (user) {
      setDisplayName(user.display_name ?? '')
      setBio(user.bio ?? '')
      setAvatarUrl(user.avatar_url ?? '')
      setIsPublic(user.is_public ?? true)
    }
  }, [user])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!token) return
    setLoading(true)
    setError(null)
    setSuccess(false)
    try {
      await updateProfile(token, {
        display_name: displayName,
        bio,
        avatar_url: avatarUrl,
        is_public: isPublic,
      })
      setSuccess(true)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to save')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-lg mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Settings</h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Display name"
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          placeholder="Your display name"
        />
        <div className="space-y-1">
          <label className="text-sm font-medium">Bio</label>
          <textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            rows={3}
            placeholder="Tell us about yourself"
            className="w-full px-3 py-2 text-sm rounded-md border border-input bg-background focus:outline-none focus:ring-2 focus:ring-ring resize-none"
          />
        </div>
        <Input
          label="Avatar URL"
          value={avatarUrl}
          onChange={(e) => setAvatarUrl(e.target.value)}
          placeholder="https://example.com/avatar.jpg"
        />
        <div className="flex items-center gap-3">
          <input
            type="checkbox"
            id="is_public"
            checked={isPublic}
            onChange={(e) => setIsPublic(e.target.checked)}
            className="w-4 h-4"
          />
          <label htmlFor="is_public" className="text-sm">
            Public profile (others can view your collection)
          </label>
        </div>

        {error && <p className="text-sm text-destructive">{error}</p>}
        {success && <p className="text-sm text-green-600">Saved successfully!</p>}

        <Button type="submit" loading={loading}>
          Save changes
        </Button>
      </form>

      <div className="border-t pt-4">
        <p className="text-sm text-muted-foreground">
          Username: <span className="font-mono font-medium">{user?.username}</span>
        </p>
        <p className="text-sm text-muted-foreground mt-1">
          Email: <span className="font-medium">{user?.email}</span>
        </p>
      </div>
    </div>
  )
}
