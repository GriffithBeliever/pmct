import { apiFetch } from './client'
import type { AuthResponse, User } from '@/lib/types'

export async function loginAPI(
  email: string,
  password: string,
): Promise<AuthResponse> {
  return apiFetch<AuthResponse>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export async function registerAPI(
  username: string,
  email: string,
  password: string,
): Promise<AuthResponse> {
  return apiFetch<AuthResponse>('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify({ username, email, password }),
  })
}

export async function fetchMe(token: string): Promise<User> {
  return apiFetch<User>('/api/auth/me', { token })
}
