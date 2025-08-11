import { useEffect, useMemo, useState } from 'react'
import { useLocation, Link } from 'react-router-dom'
import AsciinemaPlayer from '../components/AsciinemaPlayer'
import { CastItem, getAdminAuthHeader, listUserCasts } from '../api/client'
import { formatBytes, formatDate } from '../utils/format'
import { Copy, ArrowLeft } from 'lucide-react'

export default function Player() {
  const location = useLocation()
  const relPath = useMemo(() => decodeURIComponent(location.pathname.replace(/^\/play\//, '')), [location.pathname])
  const [blobUrl, setBlobUrl] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)
  const [meta, setMeta] = useState<Pick<CastItem, 'sizeBytes' | 'mtime'> | null>(() => {
    const st = (location as any).state as { item?: CastItem } | undefined
    if (st?.item) return { sizeBytes: st.item.sizeBytes, mtime: st.item.mtime }
    return null
  })

  useEffect(() => {
    const controller = new AbortController()
    const url = `/api/casts/file?path=${encodeURIComponent(relPath)}`
    fetch(url, { headers: getAdminAuthHeader(), signal: controller.signal })
      .then(async (res) => {
        if (!res.ok) throw new Error(`fetch cast failed: ${res.status}`)
        const blob = await res.blob()
        const objectUrl = URL.createObjectURL(blob)
        setBlobUrl(objectUrl)
        setError(null)
      })
      .catch((e) => { setBlobUrl(null); setError(String(e)) })
    return () => { controller.abort(); if (blobUrl) URL.revokeObjectURL(blobUrl) }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [relPath])

  useEffect(() => {
    if (meta) return
    // try fetch metadata from list endpoint using username
    const [username] = relPath.split('/')
    if (!username) return
    listUserCasts(username)
      .then(list => {
        const found = list.find(i => i.relPath === relPath)
        if (found) setMeta({ sizeBytes: found.sizeBytes, mtime: found.mtime })
      })
      .catch(() => { /* ignore meta errors */ })
  }, [relPath, meta])

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
        <div className="min-w-0">
          <div className="text-sm text-gray-500 break-all font-mono">{relPath}</div>
          <div className="mt-1 flex flex-wrap items-center gap-2 text-xs text-gray-500">
            {meta?.sizeBytes != null && (
              <span className="inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 px-2 py-0.5">{formatBytes(meta.sizeBytes)}</span>
            )}
            {meta?.mtime && (
              <span className="inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 px-2 py-0.5">{formatDate(meta.mtime)}</span>
            )}
          </div>
        </div>
        <div className="shrink-0 flex items-center gap-2">
          <Link to={-1 as any} className="inline-flex items-center gap-1.5 px-3 py-1.5 border rounded hover:bg-gray-50 dark:hover:bg-gray-700/50 text-sm"><ArrowLeft size={14} /> 返回</Link>
          <button
            className="inline-flex items-center gap-1.5 px-3 py-1.5 border rounded hover:bg-gray-50 dark:hover:bg-gray-700/50 text-sm"
            onClick={async () => {
              const url = `${window.location.origin}/play/${encodeURIComponent(relPath)}`
              await navigator.clipboard.writeText(url)
              setCopied(true); setTimeout(() => setCopied(false), 1000)
            }}
          >
            <Copy size={14} /> {copied ? '已复制' : '复制链接'}
          </button>
        </div>
      </div>

      <div className="p-4 border rounded-xl bg-white dark:bg-gray-800">
        {error && <div className="text-red-600 text-sm mb-3">{error}</div>}
    {blobUrl ? (
          <div className="w-full overflow-auto">
      <AsciinemaPlayer src={blobUrl} />
          </div>
        ) : (
          <div>加载中…</div>
        )}
      </div>
    </div>
  )
}
