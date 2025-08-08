import { useEffect, useRef } from 'react'

declare global {
  interface Window { AsciinemaPlayer?: any }
}

type Props = { src: string, cols?: number, rows?: number, autoPlay?: boolean }

export default function AsciinemaPlayer({ src, cols = 80, rows = 24, autoPlay = true }: Props) {
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!ref.current || !window.AsciinemaPlayer) return
    const el = window.AsciinemaPlayer.create(src, ref.current, { cols, rows, autoplay: autoPlay })
    return () => { try { el && el.remove() } catch { /* noop */ } }
  }, [src, cols, rows, autoPlay])

  return <div ref={ref} />
}
