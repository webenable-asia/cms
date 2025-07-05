'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Save, Eye, Calendar, Star, Settings, ArrowLeft } from 'lucide-react'
import { api } from '../../../../lib/api'
import RichTextEditor from '../../../../components/RichTextEditor'
import ImageUploader from '../../../../components/ImageUploader'
import TagInput from '../../../../components/TagInput'
import { ThemeToggle } from '@/components/ui/theme-toggle'

const predefinedCategories = [
  'Technology', 'Business', 'Marketing', 'Design', 'Development', 
  'Strategy', 'Innovation', 'Tutorials', 'Case Studies', 'News'
]

const predefinedTags = [
  'react', 'nextjs', 'javascript', 'typescript', 'web-development',
  'ui-ux', 'design-systems', 'seo', 'performance', 'accessibility',
  'marketing', 'strategy', 'business', 'startup', 'innovation'
]

export default function NewPost() {
  const [post, setPost] = useState({
    title: '',
    content: '',
    excerpt: '',
    status: 'draft',
    tags: [] as string[],
    categories: [] as string[],
    featured_image: '',
    image_alt: '',
    meta_title: '',
    meta_description: '',
    is_featured: false,
    scheduled_at: ''
  })
  const [activeTab, setActiveTab] = useState('content')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [saveStatus, setSaveStatus] = useState('')
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/admin')
    }
  }, [router])

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
        // Auto-generate meta fields if empty
        meta_title: post.meta_title || post.title,
        meta_description: post.meta_description || post.excerpt
      }

      await api.post('/posts', postData)
      setSaveStatus('Saved!')
      setTimeout(() => {
        router.push('/admin/dashboard')
      }, 1000)
    } catch (error: any) {
      console.error('Save error:', error)
      setError('Failed to save post')
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

  return (
    <div className="min-h-screen bg-background">
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
              <div className="h-6 w-px bg-border" />
              <h1 className="text-xl font-semibold text-foreground">Create New Post</h1>
            </div>
            
            <div className="flex items-center space-x-3">
              <ThemeToggle />
              {saveStatus && (
                <span className="text-sm text-green-600 dark:text-green-400">{saveStatus}</span>
              )}
              <button
                type="button"
                onClick={saveDraft}
                disabled={loading}
                className="px-4 py-2 text-sm border border-border rounded-md text-muted-foreground hover:bg-muted disabled:opacity-50 flex items-center"
              >
                <Save size={16} className="mr-2" />
                Save Draft
              </button>
              <button
                type="submit"
                form="post-form"
                disabled={loading || !post.title || !post.content}
                className="px-4 py-2 text-sm bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50 flex items-center"
              >
                <Eye size={16} className="mr-2" />
                Publish
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
              <div className="border-b border-gray-200">
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
                        className="w-full px-3 py-2 border border-border rounded-md focus:ring-ring focus:border-ring bg-background text-foreground placeholder:text-muted-foreground"
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
                      value={post.featured_image}
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
                        value={post.meta_title}
                        onChange={(e) => setPost({...post, meta_title: e.target.value})}
                        placeholder="SEO title (defaults to post title)"
                        className="w-full px-3 py-2 border border-border rounded-md focus:ring-ring focus:border-ring bg-background text-foreground placeholder:text-muted-foreground"
                        maxLength={60}
                      />
                      <p className="mt-1 text-xs text-muted-foreground">
                        {post.meta_title.length}/60 characters
                      </p>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-foreground mb-2">
                        Meta Description
                      </label>
                      <textarea
                        rows={3}
                        value={post.meta_description}
                        onChange={(e) => setPost({...post, meta_description: e.target.value})}
                        placeholder="SEO description (defaults to excerpt)"
                        className="w-full px-3 py-2 border border-border rounded-md focus:ring-ring focus:border-ring bg-background text-foreground placeholder:text-muted-foreground"
                        maxLength={160}
                      />
                      <p className="mt-1 text-xs text-muted-foreground">
                        {post.meta_description.length}/160 characters
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
                          className="w-full px-3 py-2 border border-border rounded-md focus:ring-ring focus:border-ring bg-background text-foreground placeholder:text-muted-foreground"
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
                            value={post.scheduled_at}
                            onChange={(e) => setPost({...post, scheduled_at: e.target.value})}
                            className="w-full px-3 py-2 border border-border rounded-md focus:ring-ring focus:border-ring bg-background text-foreground placeholder:text-muted-foreground"
                          />
                        </div>
                      )}
                    </div>

                    <div className="flex items-center">
                      <input
                        type="checkbox"
                        id="featured"
                        checked={post.is_featured}
                        onChange={(e) => setPost({...post, is_featured: e.target.checked})}
                        className="h-4 w-4 text-primary focus:ring-ring border-border rounded"
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
                <div className="bg-destructive/10 border border-destructive/20 rounded-md p-4">
                  <p className="text-destructive text-sm">{error}</p>
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
                  value={post.categories}
                  onChange={(categories) => setPost({...post, categories})}
                  label="Categories"
                  placeholder="Add categories..."
                  suggestions={predefinedCategories}
                />

                <TagInput
                  value={post.tags}
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
