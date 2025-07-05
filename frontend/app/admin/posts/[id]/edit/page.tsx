'use client'

import { useState, useEffect } from 'react'
import { useRouter, useParams } from 'next/navigation'
import Link from 'next/link'
import { Save, Eye, Calendar, Star, Settings, ArrowLeft } from 'lucide-react'
import { api } from '../../../../../lib/api'
import RichTextEditor from '../../../../../components/RichTextEditor'
import ImageUploader from '../../../../../components/ImageUploader'
import TagInput from '../../../../../components/TagInput'

const predefinedCategories = [
  'Technology', 'Business', 'Marketing', 'Design', 'Development', 
  'Strategy', 'Innovation', 'Tutorials', 'Case Studies', 'News'
]

const predefinedTags = [
  'react', 'nextjs', 'javascript', 'typescript', 'web-development',
  'ui-ux', 'design-systems', 'seo', 'performance', 'accessibility',
  'marketing', 'strategy', 'business', 'startup', 'innovation'
]

interface Post {
  id: string
  title: string
  content: string
  excerpt: string
  status: string
  tags: string[]
  categories?: string[]
  featured_image?: string
  image_alt?: string
  meta_title?: string
  meta_description?: string
  is_featured?: boolean
  scheduled_at?: string
}

export default function EditPost() {
  const [post, setPost] = useState<Post>({
    id: '',
    title: '',
    content: '',
    excerpt: '',
    status: 'draft',
    tags: [],
    categories: [],
    featured_image: '',
    image_alt: '',
    meta_title: '',
    meta_description: '',
    is_featured: false,
    scheduled_at: ''
  })
  const [activeTab, setActiveTab] = useState('content')
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState('')
  const [saveStatus, setSaveStatus] = useState('')
  const router = useRouter()
  const params = useParams()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/admin')
      return
    }

    fetchPost()
  }, [router, params.id])

  const fetchPost = async () => {
    try {
      const response = await api.get(`/posts/${params.id}`)
      setPost({
        ...response.data,
        categories: response.data.categories || [],
        featured_image: response.data.featured_image || '',
        image_alt: response.data.image_alt || '',
        meta_title: response.data.meta_title || '',
        meta_description: response.data.meta_description || '',
        is_featured: response.data.is_featured || false,
        scheduled_at: response.data.scheduled_at || ''
      })
    } catch (error) {
      console.error('Error fetching post:', error)
      setError('Failed to load post')
    } finally {
      setFetching(false)
    }
  }

  const calculateReadingTime = (content: string) => {
    const wordsPerMinute = 200
    const words = content.trim().split(/\s+/).length
    return Math.ceil(words / wordsPerMinute)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await savePost('published')
  }

  const saveDraft = async () => {
    await savePost('draft')
  }

  const savePost = async (status: string) => {
    setLoading(true)
    setError('')
    setSaveStatus('Saving...')

    try {
      const postData = {
        ...post,
        status,
        reading_time: calculateReadingTime(post.content),
        meta_title: post.meta_title || post.title,
        meta_description: post.meta_description || post.excerpt
      }

      await api.put(`/posts/${params.id}`, postData)
      setSaveStatus('Saved!')
      setTimeout(() => {
        router.push('/admin/dashboard')
      }, 1000)
    } catch (error: any) {
      console.error('Save error:', error)
      setError('Failed to update post')
      setSaveStatus('')
    } finally {
      setLoading(false)
    }
  }

  const handleImageChange = (url: string, alt?: string) => {
    setPost({
      ...post,
      featured_image: url,
      image_alt: alt || ''
    })
  }

  const tabs = [
    { id: 'content', label: 'Content', icon: null },
    { id: 'media', label: 'Media', icon: null },
    { id: 'seo', label: 'SEO', icon: null },
    { id: 'settings', label: 'Settings', icon: Settings }
  ]

  if (fetching) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-xl">Loading post...</div>
      </div>
    )
  }

  if (error && !post.id) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">Post Not Found</h1>
          <Link 
            href="/admin/dashboard"
            className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
          >
            Back to Dashboard
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-card border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-4">
              <Link 
                href="/admin/dashboard" 
                className="flex items-center text-muted-foreground hover:text-foreground"
              >
                <ArrowLeft size={20} className="mr-2" />
                Back to Dashboard
              </Link>
              <div className="h-6 w-px bg-gray-300" />
              <h1 className="text-xl font-semibold text-foreground">Edit Post</h1>
            </div>
            
            <div className="flex items-center space-x-3">
              {saveStatus && (
                <span className="text-sm text-green-600">{saveStatus}</span>
              )}
              <button
                type="button"
                onClick={saveDraft}
                disabled={loading}
                className="px-4 py-2 text-sm border border-gray-300 rounded-md text-foreground hover:bg-gray-50 disabled:opacity-50 flex items-center"
              >
                <Save size={16} className="mr-2" />
                Save Draft
              </button>
              <button
                type="submit"
                form="post-form"
                disabled={loading || !post.title || !post.content}
                className="px-4 py-2 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 flex items-center"
              >
                <Eye size={16} className="mr-2" />
                Update & Publish
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Main Content */}
          <div className="lg:col-span-3">
            <form id="post-form" onSubmit={handleSubmit} className="space-y-6">
              {/* Title */}
              <div>
                <input
                  type="text"
                  placeholder="Enter post title..."
                  value={post.title}
                  onChange={(e) => setPost({...post, title: e.target.value})}
                  className="w-full text-3xl font-bold border-none outline-none placeholder-gray-400 bg-transparent"
                  required
                />
              </div>

              {/* Tabs */}
              <div className="border-b border-border">
                <nav className="-mb-px flex space-x-8">
                  {tabs.map((tab) => (
                    <button
                      key={tab.id}
                      type="button"
                      onClick={() => setActiveTab(tab.id)}
                      className={`py-2 px-1 border-b-2 font-medium text-sm flex items-center ${
                        activeTab === tab.id
                          ? 'border-blue-500 text-blue-600'
                          : 'border-transparent text-muted-foreground hover:text-foreground hover:border-gray-300'
                      }`}
                    >
                      {tab.icon && <tab.icon size={16} className="mr-2" />}
                      {tab.label}
                    </button>
                  ))}
                </nav>
              </div>

              {/* Tab Content */}
              <div className="space-y-6">
                {activeTab === 'content' && (
                  <>
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-2">
                        Excerpt
                      </label>
                      <textarea
                        rows={3}
                        value={post.excerpt}
                        onChange={(e) => setPost({...post, excerpt: e.target.value})}
                        placeholder="Brief description of your post..."
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-foreground mb-2">
                        Content
                      </label>
                      <RichTextEditor
                        value={post.content}
                        onChange={(content) => setPost({...post, content})}
                        placeholder="Start writing your post..."
                      />
                    </div>
                  </>
                )}

                {activeTab === 'media' && (
                  <div>
                    <ImageUploader
                      value={post.featured_image ?? ''}
                      onChange={handleImageChange}
                      label="Featured Image"
                    />
                  </div>
                )}

                {activeTab === 'seo' && (
                  <>
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-2">
                        Meta Title
                      </label>
                      <input
                        type="text"
                        value={post.meta_title || ''}
                        onChange={(e) => setPost({...post, meta_title: e.target.value})}
                        placeholder="SEO title (defaults to post title)"
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                        maxLength={60}
                      />
                      <p className="mt-1 text-xs text-muted-foreground">
                        {(post.meta_title || '').length}/60 characters
                      </p>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-foreground mb-2">
                        Meta Description
                      </label>
                      <textarea
                        rows={3}
                        value={post.meta_description || ''}
                        onChange={(e) => setPost({...post, meta_description: e.target.value})}
                        placeholder="SEO description (defaults to excerpt)"
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                        maxLength={160}
                      />
                      <p className="mt-1 text-xs text-muted-foreground">
                        {(post.meta_description || '').length}/160 characters
                      </p>
                    </div>
                  </>
                )}

                {activeTab === 'settings' && (
                  <>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div>
                        <label className="block text-sm font-medium text-foreground mb-2">
                          Status
                        </label>
                        <select
                          value={post.status}
                          onChange={(e) => setPost({...post, status: e.target.value})}
                          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                        >
                          <option value="draft">Draft</option>
                          <option value="published">Published</option>
                          <option value="scheduled">Scheduled</option>
                        </select>
                      </div>

                      {post.status === 'scheduled' && (
                        <div>
                          <label className="block text-sm font-medium text-foreground mb-2">
                            Publish Date
                          </label>
                          <input
                            type="datetime-local"
                            value={post.scheduled_at || ''}
                            onChange={(e) => setPost({...post, scheduled_at: e.target.value})}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                          />
                        </div>
                      )}
                    </div>

                    <div className="flex items-center">
                      <input
                        type="checkbox"
                        id="featured"
                        checked={post.is_featured || false}
                        onChange={(e) => setPost({...post, is_featured: e.target.checked})}
                        className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                      />
                      <label htmlFor="featured" className="ml-2 block text-sm text-foreground flex items-center">
                        <Star size={16} className="mr-1" />
                        Featured Post
                      </label>
                    </div>
                  </>
                )}
              </div>

              {error && (
                <div className="bg-red-50 border border-red-200 rounded-md p-4">
                  <p className="text-red-600 text-sm">{error}</p>
                </div>
              )}
            </form>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            <div className="bg-card p-6 rounded-lg border border-border">
              <h3 className="text-lg font-medium text-foreground mb-4">Post Details</h3>
              
              <div className="space-y-4">
                <TagInput
                  value={post.categories || []}
                  onChange={(categories) => setPost({...post, categories})}
                  label="Categories"
                  placeholder="Add categories..."
                  suggestions={predefinedCategories}
                />

                <TagInput
                  value={post.tags || []}
                  onChange={(tags) => setPost({...post, tags})}
                  label="Tags"
                  placeholder="Add tags..."
                  suggestions={predefinedTags}
                />
              </div>
            </div>

            <div className="bg-card p-6 rounded-lg border border-border">
              <h3 className="text-lg font-medium text-foreground mb-4">Statistics</h3>
              
              <div className="space-y-3">
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Word Count:</span>
                  <span className="font-medium">
                    {post.content.split(/\s+/).filter(word => word.length > 0).length}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Reading Time:</span>
                  <span className="font-medium">{calculateReadingTime(post.content)} min</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Characters:</span>
                  <span className="font-medium">{post.content.length}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
