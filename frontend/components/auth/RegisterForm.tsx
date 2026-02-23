'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Input } from '@/components/ui/Input'
import { Button } from '@/components/ui/Button'
import { registerAPI } from '@/lib/api/auth'
import { useAuth } from '@/lib/hooks/useAuth'

export function RegisterForm() {
  const router = useRouter()
  const { login } = useAuth()
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const { user, token } = await registerAPI(username, email, password)
      login(token, user)
      router.push('/collection')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <Input
        id="username"
        label="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        required
        autoComplete="username"
        placeholder="yourusername"
        minLength={3}
      />
      <Input
        id="email"
        label="Email"
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
        autoComplete="email"
        placeholder="you@example.com"
      />
      <Input
        id="password"
        label="Password"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
        autoComplete="new-password"
        placeholder="••••••••"
        minLength={8}
      />
      {error ? (
        <p className="text-sm text-red-500 bg-red-50 dark:bg-red-950/20 rounded-md px-3 py-2">
          {error}
        </p>
      ) : null}
      <Button type="submit" loading={loading} className="w-full">
        Create account
      </Button>
    </form>
  )
}
