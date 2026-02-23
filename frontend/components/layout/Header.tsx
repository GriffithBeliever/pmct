'use client'

import Link from 'next/link'
import { Plus, LogOut } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { Avatar } from '@/components/ui/Avatar'
import { Button } from '@/components/ui/Button'
import { ThemeToggle } from '@/components/ui/ThemeToggle'
import { useAuth } from '@/lib/hooks/useAuth'

export function Header() {
  const { user, logout } = useAuth()
  const router = useRouter()

  const handleLogout = () => {
    logout()
    router.push('/login')
  }

  return (
    <header className="sticky top-0 z-40 flex h-14 items-center justify-between border-b border-border bg-background/95 backdrop-blur px-4 md:px-6">
      {/* Mobile logo */}
      <Link href="/collection" className="flex md:hidden items-center gap-2 font-bold">
        <span className="text-xl">📚</span>
        <span>PMCT</span>
      </Link>
      <div className="hidden md:block" />

      <div className="flex items-center gap-2">
        <ThemeToggle />
        <Link
          href="/media/new"
          className="inline-flex items-center justify-center gap-2 rounded-md text-sm font-medium transition-colors h-8 px-3 text-xs bg-primary text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="h-4 w-4" />
          <span className="hidden sm:inline">Add</span>
        </Link>
        {user ? (
          <div className="flex items-center gap-2">
            <Link href="/settings">
              <Avatar name={user.display_name || user.username} size="sm" />
            </Link>
            <Button variant="ghost" size="sm" onClick={handleLogout} title="Sign out">
              <LogOut className="h-4 w-4" />
            </Button>
          </div>
        ) : null}
      </div>
    </header>
  )
}
