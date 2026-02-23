import { cookies } from 'next/headers'
import { redirect, notFound } from 'next/navigation'
import { fetchMediaItem } from '@/lib/api/media'
import { MediaForm } from '@/components/media/MediaForm'

interface PageProps {
  params: Promise<{ id: string }>
}

export default async function EditMediaPage({ params }: PageProps) {
  const cookieStore = await cookies()
  const token = cookieStore.get('ems_token')?.value
  if (!token) redirect('/login')

  const { id } = await params

  const item = await fetchMediaItem(token, id).catch(() => null)
  if (!item) notFound()

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Edit Item</h1>
      <MediaForm token={token} item={item} />
    </div>
  )
}
