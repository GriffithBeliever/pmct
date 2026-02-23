'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Library, Search, Sparkles, Compass, Settings, Activity } from 'lucide-react'
import { cn } from '@/lib/utils/cn'

const navItems = [
  { href: '/collection', label: 'Collection', icon: Library },
  { href: '/search', label: 'Search', icon: Search },
  { href: '/insights', label: 'AI Insights', icon: Sparkles },
  { href: '/discover', label: 'Discover', icon: Compass },
  { href: '/activity', label: 'Activity', icon: Activity },
  { href: '/settings', label: 'Settings', icon: Settings },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <aside className="hidden md:flex flex-col w-56 shrink-0 border-r border-border bg-card h-screen sticky top-0">
      <div className="p-4 border-b border-border">
        <Link href="/collection" className="flex items-center gap-2 font-bold text-lg">
          <span className="text-2xl">📚</span>
          <span>PMCT</span>
        </Link>
      </div>
      <nav className="flex-1 p-2 space-y-1">
        {navItems.map(({ href, label, icon: Icon }) => (
          <Link
            key={href}
            href={href}
            className={cn(
              'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
              pathname === href || pathname.startsWith(href + '/')
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:bg-muted hover:text-foreground',
            )}
          >
            <Icon className="h-4 w-4 shrink-0" />
            {label}
          </Link>
        ))}
      </nav>
    </aside>
  )
}
