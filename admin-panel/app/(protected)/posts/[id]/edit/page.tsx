import PostEditor from '@/components/post-editor'

interface EditPostPageProps {
  params: Promise<{
    id: string
  }>
}

export default async function EditPostPage({ params }: EditPostPageProps) {
  const { id } = await params
  return <PostEditor postId={id} mode="edit" />
}
