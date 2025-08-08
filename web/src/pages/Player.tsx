import { useEffect, useMemo, useState } from 'react'
import { useLocation } from 'react-router-dom'
import AsciinemaPlayer from '../components/AsciinemaPlayer'
import { getAdminAuthHeader } from '../api/client'

export default function Player() {
  const location = useLocation()
  const relPath = useMemo(() => decodeURIComponent(location.pathname.replace(/^\/play\//, '')), [location.pathname])
  const [blobUrl, setBlobUrl] = useState<string | null>(null)

  useEffect(() => {
    const controller = new AbortController()
    const url = `/api/casts/file?path=${encodeURIComponent(relPath)}`
    fetch(url, { headers: getAdminAuthHeader(), signal: controller.signal })
      .then(async (res) => {
        if (!res.ok) throw new Error(`fetch cast failed: ${res.status}`)
        const blob = await res.blob()
        const objectUrl = URL.createObjectURL(blob)
        setBlobUrl(objectUrl)
      })
      .catch(() => setBlobUrl(null))
    return () => { controller.abort(); if (blobUrl) URL.revokeObjectURL(blobUrl) }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [relPath])

  return (
    <div className="space-y-4">
      <div className="text-sm text-gray-500">{relPath}</div>
      {blobUrl ? <AsciinemaPlayer src={blobUrl} /> : <div>加载中…</div>}
    </div>
  )
}
