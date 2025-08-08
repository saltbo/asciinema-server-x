import { useEffect, useState } from 'react'
import { listUsers, setAdminCreds } from '../api/client'
import { Link } from 'react-router-dom'

export default function Home() {
  const [users, setUsers] = useState<string[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [adminUser, setAdminUser] = useState(localStorage.getItem('admin_user') || '')
  const [adminPass, setAdminPass] = useState(localStorage.getItem('admin_pass') || '')

  useEffect(() => {
    listUsers().then(setUsers).catch(e => setError(String(e)))
  }, [])

  return (
    <div className="space-y-6">
      <section className="p-4 border rounded-lg bg-white dark:bg-gray-800">
        <h2 className="font-semibold mb-2">管理员凭据（MVP 临时）</h2>
        <div className="flex gap-2">
          <input className="border p-2 rounded flex-1" placeholder="ADMIN_BASIC_USER" value={adminUser} onChange={e => setAdminUser(e.target.value)} />
          <input className="border p-2 rounded flex-1" placeholder="ADMIN_BASIC_PASS" type="password" value={adminPass} onChange={e => setAdminPass(e.target.value)} />
          <button className="px-3 py-2 bg-blue-600 text-white rounded" onClick={() => { setAdminCreds(adminUser, adminPass); location.reload() }}>保存</button>
        </div>
      </section>

      <section className="p-4 border rounded-lg bg-white dark:bg-gray-800">
        <h2 className="font-semibold mb-3">用户列表</h2>
        {error && <div className="text-red-600">{error}</div>}
        {!users && !error && <div>加载中…</div>}
        {users && (
          <ul className="list-disc pl-5 space-y-1">
            {users.map(u => (
              <li key={u}><Link className="text-blue-600 hover:underline" to={`/users/${encodeURIComponent(u)}`}>{u}</Link></li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}
