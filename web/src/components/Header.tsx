import { Link, NavLink } from 'react-router-dom'
import ThemeToggle from './ThemeToggle'
import { Github, TerminalSquare } from 'lucide-react'

export default function Header() {
  return (
    <header className="sticky top-0 z-10 backdrop-blur supports-[backdrop-filter]:bg-white/70 dark:supports-[backdrop-filter]:bg-gray-900/70 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800">
      <div className="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 font-semibold tracking-tight">
          <span className="inline-flex items-center justify-center w-6 h-6 rounded bg-blue-600 text-white">
            <TerminalSquare size={16} />
          </span>
          <span>asciinema-server-x</span>
        </Link>
    <nav className="flex items-center gap-1.5">
          <a
            href="https://github.com/saltbo/asciinema-server-x"
            target="_blank"
            rel="noreferrer"
      className="inline-flex items-center gap-1.5 text-sm px-2 py-1 rounded text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
            title="GitHub 仓库"
          >
            <Github size={16} /> GitHub
          </a>
          <ThemeToggle />
        </nav>
      </div>
    </header>
  )
}
