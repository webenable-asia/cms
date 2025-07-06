import PostEditor from '@/components/post-editor'

interface EditPostPageProps {
  params: {
    id: string
  }
}

export default function EditPostPage({ params }: EditPostPageProps) {
  return <PostEditor postId={params.id} mode="edit" />
}
