import { useEffect, useMemo, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { listUserCasts, CastItem } from '../api/client'
import { formatBytes, formatDuration } from '../utils/format'
import { Skeleton, TextSkeleton } from '../components/Skeleton'
import { Play, Terminal, Copy } from 'lucide-react'

export default function UserCasts() {
  const { username = '' } = useParams()
  const [items, setItems] = useState<CastItem[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState<string | null>(null)

  useEffect(() => {
    if (!username) return
    listUserCasts(username).then(setItems).catch(e => setError(String(e)))
  }, [username])

  if (error) {
    return <div className="p-4 border border-red-300 text-red-700 rounded bg-red-50 dark:bg-red-900/20">{error}</div>
  }

  return (
    <div className="space-y-4">
      <div className="flex items-end justify-between gap-3">
        <h2 className="font-semibold text-lg">{username} 的 asciicasts</h2>
        {items && <div className="text-xs text-gray-500">共 {items.length} 条</div>}
      </div>

      {!items && (
        <div className="space-y-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="p-4 border rounded-xl bg-white dark:bg-gray-800">
              <div className="flex items-start gap-3">
                <Skeleton className="h-10 w-10 rounded-full" />
                <div className="flex-1">
                  <TextSkeleton lines={2} />
                  <div className="flex gap-2 mt-3">
                    <Skeleton className="h-6 w-20 rounded-full" />
                    <Skeleton className="h-6 w-24 rounded-full" />
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {items && items.length > 0 && <GroupedList items={items} copied={copied} setCopied={setCopied} />}
      {items && items.length === 0 && (
        <div className="p-8 text-center border rounded-xl bg-white dark:bg-gray-800 text-gray-500">暂无记录</div>
      )}
    </div>
  )
}

function GroupedList({ items, copied, setCopied }: { items: CastItem[]; copied: string | null; setCopied: (s: string | null) => void }) {
  const itemsSorted = useMemo(() =>
    [...items].sort((a, b) => new Date(b.mtime).getTime() - new Date(a.mtime).getTime())
  , [items])

  const groups = useMemo(() => {
    const m = new Map<string, CastItem[]>()
    for (const it of itemsSorted) {
      const ymd = new Date(it.mtime).toISOString().slice(0, 10)
      const arr = m.get(ymd) || []
      arr.push(it)
      m.set(ymd, arr)
    }
    return m
  }, [itemsSorted])

  const handleCopy = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(text)
      setTimeout(() => setCopied(null), 1200)
    } catch {}
  }

  return (
    <div className="space-y-6">
      {Array.from(groups.entries()).map(([ymd, list]) => (
        <section key={ymd} className="space-y-3">
          <h3 className="text-sm font-medium text-gray-600 dark:text-gray-300">{ymd}</h3>
          <div className="space-y-3">
            {list.map((it) => {
              const parts = it.relPath.split('/')
              const fileName = parts.pop() || it.relPath
              const dir = parts.join('/')
              const timeText = new Date(it.mtime).toLocaleTimeString()
              const url = `${window.location.origin}/play/${encodeURIComponent(it.relPath)}`
              return (
                <div key={it.relPath} className="group p-4 rounded-xl border bg-white dark:bg-gray-800 shadow-sm hover:shadow-md hover:-translate-y-px transition">
                  <div className="flex items-start gap-3">
                    <div className="w-10 h-10 rounded-full bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 flex items-center justify-center">
                      <Terminal size={18} />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-start justify-between gap-3">
                        <div className="min-w-0">
                          <div className="font-medium truncate">{fileName}</div>
                          <div className="text-xs text-gray-500 font-mono truncate">{dir}</div>
                        </div>
                        <div className="shrink-0 flex items-center gap-2">
                          <span className="inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 text-xs px-2 py-0.5">{formatBytes(it.sizeBytes)}</span>
                          <span className="inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 text-xs px-2 py-0.5">{timeText}</span>
                          {it.metadata?.duration && (
                            <span className="inline-flex items-center rounded-full bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 text-xs px-2 py-0.5">
                              {formatDuration(it.metadata.duration)}
                            </span>
                          )}
                        </div>
                      </div>
            <div className="mt-3 flex items-center gap-2">
                        <Link className="inline-flex items-center gap-2 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 transition text-white rounded" to={`/play/${encodeURIComponent(it.relPath)}`} state={{ item: it }}>
                          <Play size={16} /> 播放
                        </Link>
                        <button
                          className="inline-flex items-center gap-1.5 px-3 py-1.5 border rounded hover:bg-gray-50 dark:hover:bg-gray-700/50 text-sm"
              onClick={() => handleCopy(url)}
              aria-label="复制链接"
                        >
              <Copy size={14} /> {copied === url ? '已复制' : '复制链接'}
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        </section>
      ))}
    </div>
  )
}
