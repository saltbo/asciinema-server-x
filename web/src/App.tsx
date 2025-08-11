import { Outlet } from 'react-router-dom'
import Header from './components/Header'

export default function App() {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100">
      <Header />
      <main className="px-4 py-6 max-w-6xl mx-auto">
        <Outlet />
      </main>
    </div>
  )
}
