'use client'

import { useState, useEffect, useCallback } from 'react'
import Cookies from 'js-cookie'
import { fetchMe } from '@/lib/api/auth'
import type { User } from '@/lib/types'

const TOKEN_KEY = 'ems_token'

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const storedToken = Cookies.get(TOKEN_KEY)
    if (!storedToken) {
      setLoading(false)
      return
    }
    setToken(storedToken)
    fetchMe(storedToken)
      .then(setUser)
      .catch(() => {
        Cookies.remove(TOKEN_KEY)
        setToken(null)
      })
      .finally(() => setLoading(false))
  }, [])

  const login = useCallback((newToken: string, newUser: User) => {
    Cookies.set(TOKEN_KEY, newToken, { expires: 7, sameSite: 'lax' })
    setToken(newToken)
    setUser(newUser)
  }, [])

  const logout = useCallback(() => {
    Cookies.remove(TOKEN_KEY)
    setToken(null)
    setUser(null)
  }, [])

  return { user, token, loading, login, logout }
}
