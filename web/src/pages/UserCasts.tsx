import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { listUserCasts, CastItem } from '../api/client'

export default function UserCasts() {
  const { username = '' } = useParams()
  const [items, setItems] = useState<CastItem[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!username) return
    listUserCasts(username).then(setItems).catch(e => setError(String(e)))
  }, [username])

  return (
    <div className="space-y-4">
      <h2 className="font-semibold">{username} 的 asciicasts</h2>
      {error && <div className="text-red-600">{error}</div>}
      {!items && !error && <div>加载中…</div>}
      {items && (
        <ul className="divide-y">
          {items.map((it) => (
            <li key={it.relPath} className="py-2 flex items-center justify-between">
              <div>
                <div className="font-mono text-sm">{it.relPath}</div>
                <div className="text-xs text-gray-500">size {(it.sizeBytes/1024).toFixed(1)} KB · {new Date(it.mtime).toLocaleString()}</div>
              </div>
              <Link className="px-3 py-1 bg-blue-600 text-white rounded" to={`/play/${encodeURIComponent(it.relPath)}`}>播放</Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
