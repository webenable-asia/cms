'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { api } from '../../../lib/api'
import { ThemeToggle } from '@/components/ui/theme-toggle'

interface Post {
  id: string
  title: string
  status: string
  author: string
  created_at: string
  tags?: string[]
  categories?: string[]
  is_featured?: boolean
  reading_time?: number
  view_count?: number
}

export default function AdminDashboard() {
  const [posts, setPosts] = useState<Post[]>([])
  const [contacts, setContacts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState<any>(null)
  const [activeTab, setActiveTab] = useState('posts')
  const [replyModal, setReplyModal] = useState<{open: boolean, contact: any}>({open: false, contact: null})
  const [replyForm, setReplyForm] = useState({subject: '', message: ''})
  const [sending, setSending] = useState(false)
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userData = localStorage.getItem('user')
    
    if (!token) {
      router.push('/admin')
      return
    }
    
    if (userData) {
      setUser(JSON.parse(userData))
    }

    fetchPosts()
    fetchContacts()
  }, [router])

  const fetchContacts = async () => {
    try {
      const response = await api.get('/contacts')
      setContacts(Array.isArray(response.data) ? response.data : [])
    } catch (error) {
      console.error('Error fetching contacts:', error)
      setContacts([])
    }
  }

  const fetchPosts = async () => {
    try {
      // Admin dashboard should show all posts (both draft and published)
      const response = await api.get('/posts')
      // Ensure we always set an array
      setPosts(Array.isArray(response.data) ? response.data : [])
    } catch (error) {
      console.error('Error fetching posts:', error)
      // Set empty array on error
      setPosts([])
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/admin')
  }

  const deletePost = async (id: string) => {
    if (!confirm('Are you sure you want to delete this post?')) return
    
    try {
      await api.delete(`/posts/${id}`)
      fetchPosts() // Refresh the list
    } catch (error: any) {
      console.error('Error deleting post:', error)
      console.error('Delete error response:', error.response?.data)
      console.error('Delete error status:', error.response?.status)
      
      let errorMessage = 'Failed to delete post'
      if (error.response?.status === 404) {
        errorMessage = 'Post not found'
      } else if (error.response?.status === 401) {
        errorMessage = 'Unauthorized to delete post'
      } else if (error.code === 'ECONNREFUSED' || error.message.includes('Network Error')) {
        errorMessage = 'Cannot connect to server'
      }
      
      alert(errorMessage)
    }
  }

  const updateContactStatus = async (id: string, status: string) => {
    try {
      await api.put(`/contacts/${id}`, { status })
      fetchContacts() // Refresh the list
    } catch (error: any) {
      console.error('Error updating contact status:', error)
      alert('Failed to update contact status')
    }
  }

  const deleteContact = async (id: string) => {
    if (!confirm('Are you sure you want to delete this contact?')) return
    
    try {
      await api.delete(`/contacts/${id}`)
      fetchContacts() // Refresh the list
    } catch (error: any) {
      console.error('Error deleting contact:', error)
      alert('Failed to delete contact')
    }
  }

  const openReplyModal = (contact: any) => {
    setReplyModal({open: true, contact})
    setReplyForm({
      subject: `Re: ${contact.subject}`,
      message: `Hi ${contact.name},\n\nThank you for contacting WebEnable. \n\n\n\nBest regards,\nWebEnable Team`
    })
  }

  const closeReplyModal = () => {
    setReplyModal({open: false, contact: null})
    setReplyForm({subject: '', message: ''})
    setSending(false)
  }

  const sendReply = async () => {
    if (!replyModal.contact) return
    
    setSending(true)
    try {
      await api.post(`/contacts/${replyModal.contact.id}/reply`, replyForm)
      alert('Reply sent successfully!')
      fetchContacts() // Refresh the list
      closeReplyModal()
    } catch (error: any) {
      console.error('Error sending reply:', error)
      alert('Failed to send reply. Please try again.')
    } finally {
      setSending(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-xl text-foreground">Loading dashboard...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <nav className="bg-card shadow-sm border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-semibold text-foreground">Admin Dashboard</h1>
            </div>
            <div className="flex items-center space-x-4">
              <ThemeToggle />
              <span className="text-muted-foreground">Welcome, {user?.username}</span>
              <button
                onClick={handleLogout}
                className="bg-destructive hover:bg-destructive/90 text-destructive-foreground px-4 py-2 rounded text-sm"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          {/* Tab Navigation */}
          <div className="border-b border-border mb-6">
            <nav className="-mb-px flex space-x-8">
              <button
                onClick={() => setActiveTab('posts')}
                className={`py-2 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'posts'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-muted-foreground'
                }`}
              >
                Posts ({posts.length})
              </button>
              <button
                onClick={() => setActiveTab('contacts')}
                className={`py-2 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'contacts'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-muted-foreground'
                }`}
              >
                Contacts ({contacts.length})
                {contacts.filter(c => c.status === 'new').length > 0 && (
                  <span className="ml-2 bg-destructive/10 text-destructive text-xs font-medium px-2.5 py-0.5 rounded-full">
                    {contacts.filter(c => c.status === 'new').length} new
                  </span>
                )}
              </button>
            </nav>
          </div>

          {activeTab === 'posts' && (
            <>
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold text-foreground">Posts</h2>
                <Link
                  href="/admin/posts/new"
                  className="bg-primary hover:bg-primary/90 text-primary-foreground font-bold py-2 px-4 rounded"
                >
                  Create New Post
                </Link>
              </div>

              {/* Posts content */}
              {!posts || posts.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-muted-foreground text-lg">No posts found.</p>
                  <Link
                    href="/admin/posts/new"
                    className="mt-4 inline-block bg-primary hover:bg-primary/90 text-primary-foreground font-bold py-2 px-4 rounded"
                  >
                    Create Your First Post
                  </Link>
                </div>
              ) : (
                <div className="bg-card shadow overflow-hidden sm:rounded-md border border-border">
                  <ul className="divide-y divide-border">{posts.map((post, index) => (
                      <li key={post.id || `post-${index}`}>
                        <div className="px-4 py-4 flex items-center justify-between">
                          <div className="flex-1">
                            <div className="flex items-center justify-between mb-2">
                              <div className="flex items-center space-x-2">
                                <p className="text-sm font-medium text-primary truncate">
                                  {post.title}
                                </p>
                                {post.is_featured && (
                                  <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200">
                                    ★ Featured
                                  </span>
                                )}
                              </div>
                              <div className="ml-2 flex-shrink-0 flex">
                                <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                  post.status === 'published' 
                                    ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' 
                                    : post.status === 'scheduled'
                                    ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
                                    : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                                }`}>
                                  {post.status}
                                </span>
                              </div>
                            </div>
                            
                            {/* Tags and Categories */}
                            {(post.tags?.length || post.categories?.length) && (
                              <div className="mb-2 flex flex-wrap gap-1">
                                {post.categories?.slice(0, 2).map((category, idx) => (
                                  <span key={idx} className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200">
                                    {category}
                                  </span>
                                ))}
                                {post.tags?.slice(0, 3).map((tag, idx) => (
                                  <span key={idx} className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-muted text-muted-foreground">
                                    #{tag}
                                  </span>
                                ))}
                                {(post.tags?.length || 0) + (post.categories?.length || 0) > 5 && (
                                  <span className="text-xs text-muted-foreground">
                                    +{(post.tags?.length || 0) + (post.categories?.length || 0) - 5} more
                                  </span>
                                )}
                              </div>
                            )}
                            
                            <div className="mt-2 sm:flex sm:justify-between">
                              <div className="sm:flex space-x-4">
                                <p className="flex items-center text-sm text-muted-foreground">
                                  By {post.author}
                                </p>
                                {post.reading_time && (
                                  <p className="flex items-center text-sm text-muted-foreground">
                                    {post.reading_time} min read
                                  </p>
                                )}
                                {post.view_count !== undefined && (
                                  <p className="flex items-center text-sm text-muted-foreground">
                                    {post.view_count} views
                                  </p>
                                )}
                              </div>
                              <div className="mt-2 flex items-center text-sm text-muted-foreground sm:mt-0">
                                <p>
                                  {new Date(post.created_at).toLocaleDateString()}
                                </p>
                              </div>
                            </div>
                          </div>
                          <div className="ml-4 flex space-x-2">
                            <Link
                              href={`/admin/posts/${post.id}/edit`}
                              className="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded text-sm"
                            >
                              Edit
                            </Link>
                            <button
                              onClick={() => deletePost(post.id)}
                              className="bg-destructive hover:bg-destructive/90 text-destructive-foreground px-3 py-1 rounded text-sm"
                            >
                              Delete
                            </button>
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </>
          )}

          {activeTab === 'contacts' && (
            <>
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold text-foreground">Contact Messages</h2>
                <div className="flex space-x-2">
                  <button
                    onClick={fetchContacts}
                    className="bg-secondary hover:bg-secondary/80 text-secondary-foreground font-bold py-2 px-4 rounded"
                  >
                    Refresh
                  </button>
                </div>
              </div>

              {/* Contacts content */}
              {!contacts || contacts.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-muted-foreground text-lg">No contact messages found.</p>
                </div>
              ) : (
                <div className="bg-card shadow overflow-hidden sm:rounded-md border border-border">
                  <ul className="divide-y divide-border">{contacts.map((contact, index) => (
                      <li key={contact.id || `contact-${index}`}>
                        <div className="px-4 py-4">
                          <div className="flex items-center justify-between mb-2">
                            <div className="flex items-center">
                              <h3 className="text-sm font-medium text-foreground">
                                {contact.name}
                              </h3>
                              <span className="ml-2 text-sm text-muted-foreground">
                                ({contact.email})
                              </span>
                              {contact.company && (
                                <span className="ml-2 text-sm text-muted-foreground">
                                  - {contact.company}
                                </span>
                              )}
                            </div>
                            <div className="flex items-center space-x-2">
                              <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                contact.status === 'new' 
                                  ? 'bg-destructive/10 text-destructive'
                                  : contact.status === 'read'
                                  ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                                  : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                              }`}>
                                {contact.status}
                              </span>
                              <span className="text-xs text-muted-foreground">
                                {new Date(contact.created_at).toLocaleDateString()}
                              </span>
                            </div>
                          </div>
                          <div className="mb-2">
                            <p className="text-sm font-medium text-foreground">
                              Subject: {contact.subject}
                            </p>
                          </div>
                          <div className="mb-4">
                            <p className="text-sm text-muted-foreground">
                              {contact.message}
                            </p>
                          </div>
                          <div className="flex justify-end space-x-2">
                            <Link
                              href={`/admin/contacts/${contact.id}/reply`}
                              className="bg-primary hover:bg-primary/90 text-primary-foreground px-3 py-1 rounded text-sm"
                            >
                              Reply
                            </Link>
                            {contact.status === 'new' && (
                              <button
                                onClick={() => updateContactStatus(contact.id, 'read')}
                                className="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded text-sm"
                              >
                                Mark as Read
                              </button>
                            )}
                            {contact.status === 'read' && (
                              <button
                                onClick={() => updateContactStatus(contact.id, 'replied')}
                                className="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm"
                              >
                                Mark as Replied
                              </button>
                            )}
                            <button
                              onClick={() => deleteContact(contact.id)}
                              className="bg-destructive hover:bg-destructive/90 text-destructive-foreground px-3 py-1 rounded text-sm"
                            >
                              Delete
                            </button>
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
