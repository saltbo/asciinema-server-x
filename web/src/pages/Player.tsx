import { useEffect, useMemo, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import AsciinemaPlayer from '../components/AsciinemaPlayer'
import { CastItem, getAdminAuthHeader } from '../api/client'
import { formatBytes, formatDate } from '../utils/format'
import { Copy, ArrowLeft } from 'lucide-react'

export default function Player() {
  const { id } = useParams<{ id: string }>()
  const [blobUrl, setBlobUrl] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)
  const [meta, setMeta] = useState<Pick<CastItem, 'sizeBytes' | 'mtime' | 'metadata'> | null>(null)

  // Fetch metadata
  useEffect(() => {
    if (!id) return
    const controller = new AbortController()
    const url = `/api/casts/${encodeURIComponent(id)}`
    fetch(url, { headers: getAdminAuthHeader(), signal: controller.signal })
      .then(async (res) => {
        if (!res.ok) throw new Error(`fetch metadata failed: ${res.status}`)
        const data = await res.json()
        setMeta({
          sizeBytes: data.sizeBytes,
          mtime: data.mtime,
          metadata: data.metadata
        })
      })
      .catch((e) => { 
        if (!controller.signal.aborted) {
          setError(`Failed to load metadata: ${String(e)}`)
        }
      })
    return () => controller.abort()
  }, [id])

  // Fetch cast file
  useEffect(() => {
    if (!id) return
    const controller = new AbortController()
    const url = `/api/casts/${encodeURIComponent(id)}/file`
    fetch(url, { headers: getAdminAuthHeader(), signal: controller.signal })
      .then(async (res) => {
        if (!res.ok) throw new Error(`fetch cast failed: ${res.status}`)
        const blob = await res.blob()
        const objectUrl = URL.createObjectURL(blob)
        setBlobUrl(objectUrl)
        setError(null)
      })
      .catch((e) => { 
        if (!controller.signal.aborted) {
          setBlobUrl(null)
          setError(`Failed to load cast: ${String(e)}`)
        }
      })
    return () => { 
      controller.abort()
      if (blobUrl) URL.revokeObjectURL(blobUrl)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])

  if (!id) {
    return <div className="text-red-600">Missing cast ID</div>
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
        <div className="min-w-0">
          {meta?.metadata?.title ? (
            <>
              <h1 className="text-lg font-semibold text-gray-900 dark:text-gray-100 break-words">
                {meta.metadata.title}
              </h1>
              <div className="text-sm text-gray-500 break-all font-mono">ID: {id}</div>
            </>
          ) : (
            <div className="text-lg font-semibold text-gray-900 dark:text-gray-100 break-all font-mono">ID: {id}</div>
          )}
          <div className="mt-1 flex flex-wrap items-center gap-2 text-xs text-gray-500">
            {meta?.sizeBytes != null && (
              <span className="inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 px-2 py-0.5">{formatBytes(meta.sizeBytes)}</span>
            )}
            {meta?.mtime && (
              <span className="inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 px-2 py-0.5">{formatDate(meta.mtime)}</span>
            )}
            {meta?.metadata?.width && meta?.metadata?.height && (
              <span className="inline-flex items-center rounded-full bg-blue-100 dark:bg-blue-700 text-blue-700 dark:text-blue-200 px-2 py-0.5">
                {meta.metadata.width}×{meta.metadata.height}
              </span>
            )}
            {meta?.metadata?.version && (
              <span className="inline-flex items-center rounded-full bg-green-100 dark:bg-green-700 text-green-700 dark:text-green-200 px-2 py-0.5">
                v{meta.metadata.version}
              </span>
            )}
            {meta?.metadata?.duration && (
              <span className="inline-flex items-center rounded-full bg-purple-100 dark:bg-purple-700 text-purple-700 dark:text-purple-200 px-2 py-0.5">
                {Math.round(meta.metadata.duration)}s
              </span>
            )}
            {meta?.metadata?.timestamp && (
              <span className="inline-flex items-center rounded-full bg-orange-100 dark:bg-orange-700 text-orange-700 dark:text-orange-200 px-2 py-0.5">
                {formatDate(new Date(meta.metadata.timestamp * 1000).toISOString())}
              </span>
            )}
          </div>
        </div>
        <div className="shrink-0 flex items-center gap-2">
          <Link to={-1 as any} className="inline-flex items-center gap-1.5 px-3 py-1.5 border rounded hover:bg-gray-50 dark:hover:bg-gray-700/50 text-sm"><ArrowLeft size={14} /> 返回</Link>
          <button
            className="inline-flex items-center gap-1.5 px-3 py-1.5 border rounded hover:bg-gray-50 dark:hover:bg-gray-700/50 text-sm"
            onClick={async () => {
              const url = `${window.location.origin}/a/${encodeURIComponent(id)}`
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
            <AsciinemaPlayer 
              src={blobUrl} 
              cols={meta?.metadata?.width} 
              rows={meta?.metadata?.height}
            />
          </div>
        ) : (
          <div>加载中…</div>
        )}
      </div>
    </div>
  )
}
