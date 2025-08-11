import { Moon, Sun } from 'lucide-react'
import { useEffect, useState } from 'react'

export default function ThemeToggle() {
  const [dark, setDark] = useState<boolean>(() => {
    const stored = localStorage.getItem('theme')
    if (stored) return stored === 'dark'
    return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
  })

  useEffect(() => {
    const root = document.documentElement
    if (dark) {
      root.classList.add('dark')
      localStorage.setItem('theme', 'dark')
    } else {
      root.classList.remove('dark')
      localStorage.setItem('theme', 'light')
    }
  }, [dark])

  return (
    <button
      className="inline-flex items-center justify-center rounded-md p-1.5 text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
      onClick={() => setDark(v => !v)}
      aria-label="切换主题"
      title={dark ? '切换为浅色' : '切换为深色'}
    >
      {dark ? <Moon size={18} /> : <Sun size={18} />}
    </button>
  )
}
