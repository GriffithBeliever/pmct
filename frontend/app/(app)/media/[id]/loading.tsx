import { Skeleton } from '@/components/ui/Skeleton'

export default function Loading() {
  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <div className="flex gap-6">
        <Skeleton className="w-40 h-56 flex-shrink-0 rounded-lg" />
        <div className="flex-1 space-y-3">
          <Skeleton className="h-9 w-3/4" />
          <Skeleton className="h-6 w-1/2" />
          <Skeleton className="h-5 w-40" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-full" />
        </div>
      </div>
      <Skeleton className="h-64 w-full" />
    </div>
  )
}
