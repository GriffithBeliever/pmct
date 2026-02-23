import { cn } from '@/lib/utils/cn'
import { formatStatus, statusColor } from '@/lib/utils/format'
import type { MediaStatus } from '@/lib/types'

interface StatusBadgeProps {
  status: MediaStatus
  className?: string
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
        statusColor(status),
        className,
      )}
    >
      {formatStatus(status)}
    </span>
  )
}
