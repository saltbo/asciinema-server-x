export type CastMetadata = {
  version: number
  width: number
  height: number
  timestamp: number
  title?: string
  duration?: number
}

export type CastItem = { 
  shortId: string
  sizeBytes: number
  mtime: string
  metadata?: CastMetadata
}

export function getAdminAuthHeader(): HeadersInit | undefined {
  const u = localStorage.getItem('admin_user')
  const p = localStorage.getItem('admin_pass')
  if (!u || !p) return undefined
  const token = btoa(`${u}:${p}`)
  return { Authorization: `Basic ${token}` }
}

export async function listUsers(): Promise<string[]> {
  const res = await fetch('/api/users', { headers: getAdminAuthHeader() })
  if (!res.ok) throw new Error(`listUsers failed: ${res.status}`)
  const data = await res.json()
  return data.items as string[]
}

export async function listUserCasts(username: string): Promise<CastItem[]> {
  const res = await fetch(`/api/users/${encodeURIComponent(username)}/casts`, { headers: getAdminAuthHeader() })
  if (!res.ok) throw new Error(`listUserCasts failed: ${res.status}`)
  const data = await res.json()
  // ensure mtime is ISO string
  return (data.items as any[]).map(i => ({ ...i, mtime: new Date(i.mtime).toISOString() }))
}

export function setAdminCreds(user: string, pass: string) {
  localStorage.setItem('admin_user', user)
  localStorage.setItem('admin_pass', pass)
}

// Build a URL for media requests with Basic auth embedded as query not possible;
// Instead, we rely on browser sending Authorization header via fetch/XHR.
// For <audio>/<video>/<img> tag you can't set headers, but asciinema-player fetches via XHR.
export function authorizedUrl(path: string): string {
  // We keep plain URL; the global fetch in asciinema-player will be intercepted by browser and include no headers.
  // So we recommend serving /api/casts/file behind admin only in UI context where fetch is used.
  return path
}
