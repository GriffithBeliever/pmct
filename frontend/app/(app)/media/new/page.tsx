import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { MediaForm } from '@/components/media/MediaForm'

export default async function NewMediaPage() {
  const cookieStore = await cookies()
  const token = cookieStore.get('ems_token')?.value
  if (!token) redirect('/login')

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Add to Collection</h1>
      <MediaForm token={token} />
    </div>
  )
}
