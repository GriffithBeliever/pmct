'use client'

import { useState, useEffect, useRef } from 'react'

export interface SSEState {
  content: string
  done: boolean
  error: string | null
}

export function useSSE(url: string | null): SSEState {
  const [state, setState] = useState<SSEState>({ content: '', done: false, error: null })
  const esRef = useRef<EventSource | null>(null)

  useEffect(() => {
    if (!url) return

    setState({ content: '', done: false, error: null })

    const es = new EventSource(url)
    esRef.current = es

    es.onmessage = (e) => {
      try {
        const token = JSON.parse(e.data) as string
        setState((prev) => ({ ...prev, content: prev.content + token }))
      } catch {
        // ignore parse error
      }
    }

    es.addEventListener('done', () => {
      setState((prev) => ({ ...prev, done: true }))
      es.close()
    })

    es.addEventListener('error', (e) => {
      const msgEvent = e as MessageEvent
      setState((prev) => ({ ...prev, error: msgEvent.data ?? 'Stream error', done: true }))
      es.close()
    })

    es.onerror = () => {
      setState((prev) => ({ ...prev, error: 'Connection error', done: true }))
      es.close()
    }

    return () => {
      es.close()
      esRef.current = null
    }
  }, [url])

  return state
}
