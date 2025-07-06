'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  ArrowLeft, 
  Save, 
  Eye, 
  Send, 
  X,
  Plus,
  AlertCircle,
  CheckCircle
} from 'lucide-react'
import { useCreatePost, useUpdatePost, usePost } from '@/hooks/use-api'
import { Post } from '@/types/api'
import Link from 'next/link'

interface PostEditorProps {
  postId?: string
  mode: 'create' | 'edit'
}

export default function PostEditor({ postId, mode }: PostEditorProps) {
  const router = useRouter()
  const shouldFetchPost = mode === 'edit' && postId
  const { data: existingPost, loading: loadingPost } = usePost(shouldFetchPost ? postId! : '')
  const { mutate: createPost, loading: creating, error: createError } = useCreatePost()
  const { mutate: updatePost, loading: updating, error: updateError } = useUpdatePost()
  
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    excerpt: '',
    author: 'Admin',
    status: 'draft' as 'draft' | 'published' | 'scheduled',
    tags: [] as string[],
    categories: [] as string[],
    featured_image: '',
    image_alt: '',
    meta_title: '',
    meta_description: '',
    is_featured: false,
  })

  const [newTag, setNewTag] = useState('')
  const [newCategory, setNewCategory] = useState('')
  const [saved, setSaved] = useState(false)

  // Load existing post data
  useEffect(() => {
    if (existingPost && mode === 'edit') {
      setFormData({
        title: existingPost.title || '',
        content: existingPost.content || '',
        excerpt: existingPost.excerpt || '',
        author: existingPost.author || 'Admin',
        status: existingPost.status || 'draft',
        tags: existingPost.tags || [],
        categories: existingPost.categories || [],
        featured_image: existingPost.featured_image || '',
        image_alt: existingPost.image_alt || '',
        meta_title: existingPost.meta_title || '',
        meta_description: existingPost.meta_description || '',
        is_featured: existingPost.is_featured || false,
      })
    }
  }, [existingPost, mode])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value
    }))
  }

  const addTag = () => {
    if (newTag.trim() && !formData.tags.includes(newTag.trim())) {
      setFormData(prev => ({
        ...prev,
        tags: [...prev.tags, newTag.trim()]
      }))
      setNewTag('')
    }
  }

  const removeTag = (tagToRemove: string) => {
    setFormData(prev => ({
      ...prev,
      tags: prev.tags.filter(tag => tag !== tagToRemove)
    }))
  }

  const addCategory = () => {
    if (newCategory.trim() && !formData.categories.includes(newCategory.trim())) {
      setFormData(prev => ({
        ...prev,
        categories: [...prev.categories, newCategory.trim()]
      }))
      setNewCategory('')
    }
  }

  const removeCategory = (categoryToRemove: string) => {
    setFormData(prev => ({
      ...prev,
      categories: prev.categories.filter(category => category !== categoryToRemove)
    }))
  }

  const calculateReadingTime = (content: string): number => {
    const wordsPerMinute = 200
    const words = content.trim().split(/\s+/).length
    return Math.ceil(words / wordsPerMinute)
  }

  const handleSave = async (status: 'draft' | 'published') => {
    const postData = {
      ...formData,
      status,
      reading_time: calculateReadingTime(formData.content),
      view_count: 0,
      published_at: status === 'published' ? new Date().toISOString() : undefined,
    }

    let result
    if (mode === 'create') {
      result = await createPost(postData)
    } else if (postId) {
      result = await updatePost({ id: postId, post: postData })
    }

    if (result) {
      setSaved(true)
      setTimeout(() => setSaved(false), 3000)
      if (status === 'published') {
        router.push('/admin/dashboard?tab=posts')
      }
    }
  }

  const error = createError || updateError
  const loading = creating || updating || loadingPost

  if (loadingPost && mode === 'edit') {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p>Loading post...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center gap-4">
              <Link href="/admin/dashboard">
                <Button variant="ghost" size="sm">
                  <ArrowLeft className="w-4 h-4 mr-2" />
                  Back to Dashboard
                </Button>
              </Link>
              <div>
                <h1 className="text-xl font-bold text-gray-900">
                  {mode === 'create' ? 'Create New Post' : 'Edit Post'}
                </h1>
                <p className="text-sm text-gray-600">
                  {formData.status === 'published' ? 'Published' : 'Draft'} • 
                  {formData.content ? ` ${calculateReadingTime(formData.content)} min read` : ' No content yet'}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              {saved && (
                <div className="flex items-center text-green-600 text-sm">
                  <CheckCircle className="w-4 h-4 mr-1" />
                  Saved
                </div>
              )}
              <Button 
                variant="outline" 
                onClick={() => handleSave('draft')}
                disabled={loading || !formData.title || !formData.content}
              >
                <Save className="w-4 h-4 mr-2" />
                Save Draft
              </Button>
              <Button 
                onClick={() => handleSave('published')}
                disabled={loading || !formData.title || !formData.content}
              >
                <Send className="w-4 h-4 mr-2" />
                Publish
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-md">
            <div className="flex">
              <AlertCircle className="w-5 h-5 text-red-400" />
              <div className="ml-3">
                <h3 className="text-sm font-medium text-red-800">Error</h3>
                <div className="mt-2 text-sm text-red-700">{error}</div>
              </div>
            </div>
          </div>
        )}

        <div className="grid gap-8 lg:grid-cols-3">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Post Content</CardTitle>
                <CardDescription>Write your blog post content</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="title">Title *</Label>
                  <Input
                    id="title"
                    name="title"
                    value={formData.title}
                    onChange={handleInputChange}
                    placeholder="Enter post title..."
                    required
                  />
                </div>

                <div>
                  <Label htmlFor="excerpt">Excerpt</Label>
                  <Textarea
                    id="excerpt"
                    name="excerpt"
                    value={formData.excerpt}
                    onChange={handleInputChange}
                    placeholder="Brief description of the post..."
                    rows={3}
                  />
                </div>

                <div>
                  <Label htmlFor="content">Content *</Label>
                  <Textarea
                    id="content"
                    name="content"
                    value={formData.content}
                    onChange={handleInputChange}
                    placeholder="Write your post content here..."
                    rows={20}
                    required
                  />
                  <p className="text-sm text-muted-foreground mt-1">
                    {formData.content.split(/\s+/).filter(word => word.length > 0).length} words
                  </p>
                </div>
              </CardContent>
            </Card>

            {/* SEO Settings */}
            <Card>
              <CardHeader>
                <CardTitle>SEO Settings</CardTitle>
                <CardDescription>Optimize your post for search engines</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="meta_title">Meta Title</Label>
                  <Input
                    id="meta_title"
                    name="meta_title"
                    value={formData.meta_title}
                    onChange={handleInputChange}
                    placeholder="SEO title (leave empty to use post title)"
                  />
                </div>

                <div>
                  <Label htmlFor="meta_description">Meta Description</Label>
                  <Textarea
                    id="meta_description"
                    name="meta_description"
                    value={formData.meta_description}
                    onChange={handleInputChange}
                    placeholder="Brief description for search engines..."
                    rows={3}
                  />
                  <p className="text-sm text-muted-foreground">
                    {formData.meta_description.length}/160 characters
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Post Settings */}
            <Card>
              <CardHeader>
                <CardTitle>Post Settings</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="author">Author</Label>
                  <Input
                    id="author"
                    name="author"
                    value={formData.author}
                    onChange={handleInputChange}
                  />
                </div>

                <div>
                  <Label htmlFor="status">Status</Label>
                  <select
                    id="status"
                    name="status"
                    value={formData.status}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-input bg-background rounded-md"
                  >
                    <option value="draft">Draft</option>
                    <option value="published">Published</option>
                    <option value="scheduled">Scheduled</option>
                  </select>
                </div>

                <div className="flex items-center space-x-2">
                  <input
                    id="is_featured"
                    name="is_featured"
                    type="checkbox"
                    checked={formData.is_featured}
                    onChange={handleInputChange}
                    className="rounded border-input"
                  />
                  <Label htmlFor="is_featured">Featured Post</Label>
                </div>
              </CardContent>
            </Card>

            {/* Featured Image */}
            <Card>
              <CardHeader>
                <CardTitle>Featured Image</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="featured_image">Image URL</Label>
                  <Input
                    id="featured_image"
                    name="featured_image"
                    value={formData.featured_image}
                    onChange={handleInputChange}
                    placeholder="https://example.com/image.jpg"
                  />
                </div>

                <div>
                  <Label htmlFor="image_alt">Alt Text</Label>
                  <Input
                    id="image_alt"
                    name="image_alt"
                    value={formData.image_alt}
                    onChange={handleInputChange}
                    placeholder="Describe the image..."
                  />
                </div>

                {formData.featured_image && (
                  <div className="mt-2">
                    <img
                      src={formData.featured_image}
                      alt={formData.image_alt}
                      className="w-full h-32 object-cover rounded"
                      onError={(e) => {
                        e.currentTarget.style.display = 'none'
                      }}
                    />
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Categories */}
            <Card>
              <CardHeader>
                <CardTitle>Categories</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex gap-2">
                  <Input
                    value={newCategory}
                    onChange={(e) => setNewCategory(e.target.value)}
                    placeholder="Add category..."
                    onKeyPress={(e) => e.key === 'Enter' && addCategory()}
                  />
                  <Button type="button" variant="outline" size="sm" onClick={addCategory}>
                    <Plus className="w-4 h-4" />
                  </Button>
                </div>
                <div className="flex flex-wrap gap-2">
                  {formData.categories.map((category) => (
                    <Badge key={category} variant="secondary" className="flex items-center gap-1">
                      {category}
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="h-auto p-0 ml-1"
                        onClick={() => removeCategory(category)}
                      >
                        <X className="w-3 h-3" />
                      </Button>
                    </Badge>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Tags */}
            <Card>
              <CardHeader>
                <CardTitle>Tags</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex gap-2">
                  <Input
                    value={newTag}
                    onChange={(e) => setNewTag(e.target.value)}
                    placeholder="Add tag..."
                    onKeyPress={(e) => e.key === 'Enter' && addTag()}
                  />
                  <Button type="button" variant="outline" size="sm" onClick={addTag}>
                    <Plus className="w-4 h-4" />
                  </Button>
                </div>
                <div className="flex flex-wrap gap-2">
                  {formData.tags.map((tag) => (
                    <Badge key={tag} variant="outline" className="flex items-center gap-1">
                      {tag}
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="h-auto p-0 ml-1"
                        onClick={() => removeTag(tag)}
                      >
                        <X className="w-3 h-3" />
                      </Button>
                    </Badge>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
    </div>
  )
}
