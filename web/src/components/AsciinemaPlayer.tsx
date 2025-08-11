import { useEffect, useRef } from 'react'

declare global {
  interface Window { AsciinemaPlayer?: any }
}

type Props = { src: string, cols?: number, rows?: number, autoPlay?: boolean, fit?: 'none' | 'width' | 'height' | 'both' }

export default function AsciinemaPlayer({ src, cols = 80, rows = 24, autoPlay = true, fit = 'width' }: Props) {
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const container = ref.current
    if (!container || !window.AsciinemaPlayer) return
    // Clear any previous content to avoid duplicates (React StrictMode mounts twice in dev)
    container.innerHTML = ''
    const isDark = document.documentElement.classList.contains('dark')
    const el = window.AsciinemaPlayer.create(src, container, { cols, rows, autoplay: autoPlay, fit, theme: isDark ? 'asciinema' : 'asciinema' })
    return () => {
      try { el && (el.remove ? el.remove() : null) } catch { /* noop */ }
      try { container.innerHTML = '' } catch { /* noop */ }
    }
  }, [src, cols, rows, autoPlay, fit])

  return <div ref={ref} className="max-w-full" />
}
