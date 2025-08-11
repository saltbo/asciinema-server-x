import { useEffect, useState } from 'react'
import { listUsers, setAdminCreds } from '../api/client'
import { Link } from 'react-router-dom'
import { Skeleton, TextSkeleton } from '../components/Skeleton'
import { User } from 'lucide-react'

export default function Home() {
  const [users, setUsers] = useState<string[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [adminUser, setAdminUser] = useState(localStorage.getItem('admin_user') || '')
  const [adminPass, setAdminPass] = useState(localStorage.getItem('admin_pass') || '')

  useEffect(() => {
    listUsers().then(setUsers).catch(e => setError(String(e)))
  }, [])

  return (
    <div className="space-y-8">
      <section className="p-5 border rounded-xl bg-white dark:bg-gray-800 shadow-sm">
        <h2 className="font-semibold mb-3">管理员凭据（MVP 临时）</h2>
        <p className="text-sm text-gray-500 mb-3">用于前端调用受保护 API。保存在本地浏览器，仅当前设备生效。</p>
        <div className="flex flex-col md:flex-row gap-3">
          <input className="border border-gray-300 dark:border-gray-700 bg-transparent p-2 rounded md:flex-1" placeholder="ADMIN_BASIC_USER" value={adminUser} onChange={e => setAdminUser(e.target.value)} />
          <input className="border border-gray-300 dark:border-gray-700 bg-transparent p-2 rounded md:flex-1" placeholder="ADMIN_BASIC_PASS" type="password" value={adminPass} onChange={e => setAdminPass(e.target.value)} />
          <button className="px-3 py-2 bg-blue-600 hover:bg-blue-700 transition text-white rounded" onClick={() => { setAdminCreds(adminUser, adminPass); location.reload() }}>
            保存
          </button>
        </div>
      </section>

      <section>
        <div className="flex items-center justify-between mb-3">
          <h2 className="font-semibold">用户列表</h2>
        </div>
        {error && <div className="p-4 border border-red-300 text-red-700 rounded bg-red-50 dark:bg-red-900/20">{error}</div>}
        {!users && !error && (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="p-4 border rounded-xl bg-white dark:bg-gray-800">
                <Skeleton className="h-10 w-10 rounded-full mb-3" />
                <TextSkeleton lines={1} />
              </div>
            ))}
          </div>
        )}
        {users && users.length === 0 && (
          <div className="p-8 text-center border rounded-xl bg-white dark:bg-gray-800">
            <div className="mx-auto w-12 h-12 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center mb-2">
              <User className="text-gray-400" size={20} />
            </div>
            <div className="text-gray-500">暂无用户</div>
          </div>
        )}
        {users && users.length > 0 && (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            {users.map(u => (
              <Link key={u} to={`/users/${encodeURIComponent(u)}`} className="group p-4 border rounded-xl bg-white dark:bg-gray-800 hover:shadow transition">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 flex items-center justify-center">
                    <User size={18} />
                  </div>
                  <div>
                    <div className="font-medium group-hover:text-blue-600">{u}</div>
                    <div className="text-xs text-gray-500">点击查看该用户的 casts</div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </section>
    </div>
  )
}
