import { useEffect, useRef, useState } from 'react'

interface UsePollingOptions {
  interval?: number
  enabled?: boolean
}

export function usePolling(callback: () => void, options: UsePollingOptions = {}) {
  const { interval = 5000, enabled = true } = options
  const callbackRef = useRef(callback)
  const intervalRef = useRef<NodeJS.Timeout>()

  // Update callback ref when callback changes
  useEffect(() => {
    callbackRef.current = callback
  }, [callback])

  useEffect(() => {
    if (!enabled) {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
      return
    }

    // Start polling
    intervalRef.current = setInterval(() => {
      callbackRef.current()
    }, interval)

    // Cleanup on unmount or when dependencies change
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
    }
  }, [interval, enabled])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
    }
  }, [])
}

export function useRealtimeData<T>(
  fetchFn: () => Promise<T>,
  options: UsePollingOptions & { initialData?: T } = {}
) {
  const { initialData, ...pollingOptions } = options
  const [data, setData] = useState<T | undefined>(initialData)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  const fetchData = async () => {
    try {
      setLoading(true)
      const result = await fetchFn()
      setData(result)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'))
    } finally {
      setLoading(false)
    }
  }

  // Initial fetch
  useEffect(() => {
    fetchData()
  }, [])

  // Setup polling
  usePolling(fetchData, pollingOptions)

  return { data, loading, error, refetch: fetchData }
}
