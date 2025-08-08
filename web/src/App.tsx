import { Outlet, Link } from 'react-router-dom'

export default function App() {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100">
      <header className="border-b border-gray-200 dark:border-gray-800 p-4 flex items-center justify-between">
        <h1 className="font-semibold">asciinema-server-x</h1>
        <nav className="space-x-4">
          <Link className="text-blue-600 hover:underline" to="/">Home</Link>
        </nav>
      </header>
      <main className="p-4 max-w-5xl mx-auto">
        <Outlet />
      </main>
    </div>
  )
}
