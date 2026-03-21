import { useEffect, useRef } from 'react'

export function usePolling(fetchFn, intervalMs, deps = []) {
  const savedFn = useRef(fetchFn)

  useEffect(() => {
    savedFn.current = fetchFn
  }, [fetchFn])

  useEffect(() => {
    savedFn.current()
    const id = setInterval(() => savedFn.current(), intervalMs)
    return () => clearInterval(id)
  }, [intervalMs, ...deps])
}
