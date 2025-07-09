'use client'

import { useState, useEffect, useRef } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  PlusCircle,
  FileText,
  Users,
  Mail,
  BarChart3,
  Settings,
  LogOut,
  Eye,
  Edit,
  Trash2,
  RefreshCw
} from 'lucide-react'
import { usePosts, useContacts, useLogout } from '@/hooks/use-api'
import { usePosts as usePostsAPI } from '@/hooks/use-posts'
import { formatDate, formatRelativeTime } from '@/lib/blog-utils'
import Link from 'next/link'

export default function AdminDashboard() {
  const { data: posts, loading: postsLoading, refetch: refetchPosts } = usePosts()
  const { data: contacts, loading: contactsLoading, refetch: refetchContacts } = useContacts()
  const { mutate: logout } = useLogout()
  const { deletePost } = usePostsAPI()
  
  const [activeTab, setActiveTab] = useState('overview')
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [lastRefresh, setLastRefresh] = useState(new Date())
  const intervalRef = useRef<NodeJS.Timeout>()

  // Auto-refresh functionality
  useEffect(() => {
    // Set up auto-refresh interval (30 seconds)
    const refreshInterval = setInterval(() => {
      if (document.visibilityState === 'visible') {
        handleRefresh()
      }
    }, 30000)

    intervalRef.current = refreshInterval

    // Listen for posts update events
    const handlePostsUpdated = (event: CustomEvent) => {
      console.log('Posts updated event received:', event.detail)
      // Force immediate refresh when posts are updated
      handleRefresh()
    }

    window.addEventListener('postsUpdated', handlePostsUpdated as EventListener)

    // Also check for URL refresh parameter
    const urlParams = new URLSearchParams(window.location.search)
    if (urlParams.get('refresh')) {
      console.log('Refresh parameter detected, forcing refresh')
      handleRefresh()
      // Clean up the URL
      const cleanUrl = window.location.pathname + (urlParams.get('tab') ? '?tab=' + urlParams.get('tab') : '')
      window.history.replaceState({}, '', cleanUrl)
    }

    // Cleanup interval and event listener on unmount
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
      window.removeEventListener('postsUpdated', handlePostsUpdated as EventListener)
    }
  }, [])

  // Manual refresh function
  const handleRefresh = async () => {
    setIsRefreshing(true)
    setLastRefresh(new Date())
    
    try {
      await Promise.all([
        refetchPosts(),
        refetchContacts()
      ])
    } catch (error) {
      console.error('Failed to refresh data:', error)
    } finally {
      setIsRefreshing(false)
    }
  }

  const handleDeletePost = async (postId: string) => {
    if (window.confirm('Are you sure you want to delete this post?')) {
      try {
        await deletePost(postId)
        // Refresh data after deletion
        await handleRefresh()
      } catch (error) {
        console.error('Failed to delete post:', error)
        alert('Failed to delete post. Please try again.')
      }
    }
  }

  const handleLogout = async () => {
    await logout({})
    window.location.href = '/login'
  }

  const stats = {
    totalPosts: posts?.length || 0,
    publishedPosts: posts?.filter(p => p.status === 'published').length || 0,
    draftPosts: posts?.filter(p => p.status === 'draft').length || 0,
    totalContacts: contacts?.length || 0,
    newContacts: contacts?.filter(c => c.status === 'new').length || 0,
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-indigo-50/30">
      {/* Header */}
      <header className="bg-white/80 backdrop-blur-sm shadow-lg border-b border-gray-200/50 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center space-x-4">
              <div className="h-10 w-10 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-lg flex items-center justify-center shadow-lg">
                <svg className="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <div>
                <h1 className="text-2xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                  WebEnable CMS
                </h1>
                <p className="text-gray-600 text-sm">Manage your content and contacts</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="hidden sm:flex items-center space-x-2 text-xs text-gray-500 bg-gray-100 px-3 py-1.5 rounded-full">
                <div className="h-2 w-2 bg-green-400 rounded-full animate-pulse"></div>
                <span>Last updated: {lastRefresh.toLocaleTimeString()}</span>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleRefresh}
                disabled={isRefreshing}
                className="bg-white/50 border-gray-200 hover:bg-white hover:shadow-md transition-all duration-200"
              >
                <RefreshCw className={`w-4 h-4 mr-2 ${isRefreshing ? 'animate-spin' : ''}`} />
                {isRefreshing ? 'Refreshing...' : 'Refresh'}
              </Button>
              <Link href="http://localhost/" target="_blank">
                <Button variant="outline" size="sm" className="bg-white/50 border-gray-200 hover:bg-white hover:shadow-md transition-all duration-200">
                  <Eye className="w-4 h-4 mr-2" />
                  View Site
                </Button>
              </Link>
              <Button 
                variant="outline" 
                size="sm" 
                onClick={handleLogout}
                className="bg-white/50 border-gray-200 hover:bg-red-50 hover:border-red-200 hover:text-red-600 transition-all duration-200"
              >
                <LogOut className="w-4 h-4 mr-2" />
                Logout
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-5 bg-white/70 backdrop-blur-sm border border-gray-200/50 shadow-lg rounded-xl p-1">
            <TabsTrigger value="overview" className="rounded-lg data-[state=active]:bg-blue-600 data-[state=active]:text-white data-[state=active]:shadow-lg transition-all duration-200">Overview</TabsTrigger>
            <TabsTrigger value="posts" className="rounded-lg data-[state=active]:bg-blue-600 data-[state=active]:text-white data-[state=active]:shadow-lg transition-all duration-200">Posts</TabsTrigger>
            <TabsTrigger value="contacts" className="rounded-lg data-[state=active]:bg-blue-600 data-[state=active]:text-white data-[state=active]:shadow-lg transition-all duration-200">Contacts</TabsTrigger>
            <TabsTrigger value="users" className="rounded-lg data-[state=active]:bg-blue-600 data-[state=active]:text-white data-[state=active]:shadow-lg transition-all duration-200">Users</TabsTrigger>
            <TabsTrigger value="settings" className="rounded-lg data-[state=active]:bg-blue-600 data-[state=active]:text-white data-[state=active]:shadow-lg transition-all duration-200">Settings</TabsTrigger>
          </TabsList>

          {/* Overview Tab */}
          <TabsContent value="overview" className="space-y-8 mt-8">
            {/* Stats Cards */}
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
              <Card className="bg-white/80 backdrop-blur-sm border-gray-200/50 shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-semibold text-gray-700">Total Posts</CardTitle>
                  <div className="h-8 w-8 bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg flex items-center justify-center">
                    <FileText className="h-4 w-4 text-white" />
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold text-gray-900 mb-1">{stats.totalPosts}</div>
                  <p className="text-sm text-gray-600 flex items-center">
                    <span className="h-2 w-2 bg-green-400 rounded-full mr-2"></span>
                    {stats.publishedPosts} published
                  </p>
                </CardContent>
              </Card>

              <Card className="bg-white/80 backdrop-blur-sm border-gray-200/50 shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-semibold text-gray-700">Draft Posts</CardTitle>
                  <div className="h-8 w-8 bg-gradient-to-br from-amber-500 to-orange-600 rounded-lg flex items-center justify-center">
                    <Edit className="h-4 w-4 text-white" />
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold text-gray-900 mb-1">{stats.draftPosts}</div>
                  <p className="text-sm text-gray-600">
                    Waiting to be published
                  </p>
                </CardContent>
              </Card>

              <Card className="bg-white/80 backdrop-blur-sm border-gray-200/50 shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-semibold text-gray-700">Total Contacts</CardTitle>
                  <div className="h-8 w-8 bg-gradient-to-br from-green-500 to-emerald-600 rounded-lg flex items-center justify-center">
                    <Mail className="h-4 w-4 text-white" />
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold text-gray-900 mb-1">{stats.totalContacts}</div>
                  <p className="text-sm text-gray-600">
                    All time messages
                  </p>
                </CardContent>
              </Card>

              <Card className="bg-white/80 backdrop-blur-sm border-gray-200/50 shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-semibold text-gray-700">New Messages</CardTitle>
                  <div className="h-8 w-8 bg-gradient-to-br from-purple-500 to-violet-600 rounded-lg flex items-center justify-center">
                    <Users className="h-4 w-4 text-white" />
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold text-gray-900 mb-1">{stats.newContacts}</div>
                  <p className="text-sm text-gray-600 flex items-center">
                    {stats.newContacts > 0 && <span className="h-2 w-2 bg-red-400 rounded-full mr-2 animate-pulse"></span>}
                    Unread messages
                  </p>
                </CardContent>
              </Card>
            </div>

            {/* Recent Activity */}
            <div className="grid gap-8 lg:grid-cols-2">
              <Card className="bg-white/80 backdrop-blur-sm border-gray-200/50 shadow-lg hover:shadow-xl transition-all duration-300">
                <CardHeader className="pb-4">
                  <div className="flex items-center space-x-3">
                    <div className="h-8 w-8 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center">
                      <FileText className="h-4 w-4 text-white" />
                    </div>
                    <div>
                      <CardTitle className="text-lg font-semibold text-gray-900">Recent Posts</CardTitle>
                      <CardDescription className="text-gray-600">Your latest blog posts</CardDescription>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  {postsLoading ? (
                    <div className="space-y-4">
                      {Array.from({ length: 3 }).map((_, i) => (
                        <div key={i} className="animate-pulse bg-gray-100 rounded-lg p-4">
                          <div className="h-4 bg-gray-300 rounded mb-2"></div>
                          <div className="h-3 bg-gray-200 rounded w-3/4"></div>
                        </div>
                      ))}
                    </div>
                  ) : posts && posts.length > 0 ? (
                    <div className="space-y-3">
                      {posts.slice(0, 5).map((post) => (
                        <div key={post.id} className="flex items-center justify-between p-3 bg-gray-50/80 rounded-lg hover:bg-gray-100/80 transition-colors duration-200">
                          <div className="flex-1 min-w-0">
                            <p className="font-medium text-gray-900 truncate">{post.title}</p>
                            <p className="text-sm text-gray-500 flex items-center mt-1">
                              <svg className="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                              </svg>
                              {formatRelativeTime(post.updated_at)}
                            </p>
                          </div>
                          <Badge 
                            variant={post.status === 'published' ? 'default' : 'secondary'}
                            className={post.status === 'published' ? 'bg-green-100 text-green-800 border-green-200' : 'bg-gray-100 text-gray-800 border-gray-200'}
                          >
                            {post.status}
                          </Badge>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <FileText className="h-12 w-12 text-gray-300 mx-auto mb-3" />
                      <p className="text-gray-500 mb-2">No posts yet</p>
                      <p className="text-sm text-gray-400">Create your first blog post to get started</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              <Card className="bg-white/80 backdrop-blur-sm border-gray-200/50 shadow-lg hover:shadow-xl transition-all duration-300">
                <CardHeader className="pb-4">
                  <div className="flex justify-between items-start">
                    <div className="flex items-center space-x-3">
                      <div className="h-8 w-8 bg-gradient-to-br from-green-500 to-emerald-600 rounded-lg flex items-center justify-center">
                        <Mail className="h-4 w-4 text-white" />
                      </div>
                      <div>
                        <CardTitle className="text-lg font-semibold text-gray-900">Recent Messages</CardTitle>
                        <CardDescription className="text-gray-600">Latest contact form submissions</CardDescription>
                      </div>
                    </div>
                    <Link href="/contacts">
                      <Button variant="outline" size="sm" className="bg-white/70 border-gray-200 hover:bg-white hover:shadow-md transition-all duration-200">
                        View All
                      </Button>
                    </Link>
                  </div>
                </CardHeader>
                <CardContent>
                  {contactsLoading ? (
                    <div className="space-y-3">
                      {Array.from({ length: 3 }).map((_, i) => (
                        <div key={i} className="animate-pulse">
                          <div className="h-4 bg-gray-300 rounded mb-2"></div>
                          <div className="h-3 bg-gray-200 rounded w-3/4"></div>
                        </div>
                      ))}
                    </div>
                  ) : contacts && contacts.length > 0 ? (
                    <div className="space-y-4">
                      {contacts.slice(0, 5).map((contact) => (
                        <div key={contact.id} className="flex items-center justify-between">
                          <div>
                            <p className="font-medium">{contact.name}</p>
                            <p className="text-xs text-muted-foreground truncate">
                              {contact.subject}
                            </p>
                          </div>
                          <Badge variant={contact.status === 'new' ? 'destructive' : 'secondary'}>
                            {contact.status}
                          </Badge>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-muted-foreground">No messages yet</p>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          {/* Posts Tab */}
          <TabsContent value="posts" className="space-y-6">
            <div className="flex justify-between items-center">
              <h2 className="text-2xl font-bold">Posts Management</h2>
              <Link href="/posts/new">
                <Button>
                  <PlusCircle className="w-4 h-4 mr-2" />
                  New Post
                </Button>
              </Link>
            </div>

            <Card>
              <CardContent className="p-0">
                {postsLoading ? (
                  <div className="p-6">
                    <div className="space-y-4">
                      {Array.from({ length: 5 }).map((_, i) => (
                        <div key={i} className="animate-pulse flex justify-between items-center p-4 border rounded">
                          <div className="flex-1">
                            <div className="h-5 bg-gray-300 rounded mb-2"></div>
                            <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                          </div>
                          <div className="w-20 h-6 bg-gray-300 rounded"></div>
                        </div>
                      ))}
                    </div>
                  </div>
                ) : posts && posts.length > 0 ? (
                  <div className="divide-y">
                    {posts.map((post) => (
                      <div key={post.id} className="p-6 flex items-center justify-between hover:bg-gray-50">
                        <div className="flex-1">
                          <h3 className="font-medium">{post.title}</h3>
                          <p className="text-sm text-muted-foreground mt-1">
                            By {post.author} • {formatDate(post.updated_at)}
                          </p>
                          {post.categories && post.categories.length > 0 && (
                            <div className="flex gap-1 mt-2">
                              {post.categories.map((category) => (
                                <Badge key={category} variant="outline" className="text-xs">
                                  {category}
                                </Badge>
                              ))}
                            </div>
                          )}
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge variant={post.status === 'published' ? 'default' : 'secondary'}>
                            {post.status}
                          </Badge>
                          <Link href={`/posts/${post.id}/edit`}>
                            <Button variant="ghost" size="sm">
                              <Edit className="w-4 h-4" />
                            </Button>
                          </Link>
                          <Button variant="ghost" size="sm" onClick={() => handleDeletePost(post.id)}>
                            <Trash2 className="w-4 h-4" />
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="p-6 text-center">
                    <FileText className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                    <p className="text-muted-foreground">No posts yet</p>
                    <Link href="/posts/new">
                      <Button className="mt-4">
                        <PlusCircle className="w-4 h-4 mr-2" />
                        Create your first post
                      </Button>
                    </Link>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          {/* Contacts Tab */}
          <TabsContent value="contacts" className="space-y-6">
            <div className="flex justify-between items-center">
              <h2 className="text-2xl font-bold">Contact Messages</h2>
              <Link href="/contacts">
                <Button>
                  <Mail className="w-4 h-4 mr-2" />
                  Manage Contacts
                </Button>
              </Link>
            </div>

            <Card>
              <CardContent className="p-0">
                {contactsLoading ? (
                  <div className="p-6">
                    <div className="space-y-4">
                      {Array.from({ length: 5 }).map((_, i) => (
                        <div key={i} className="animate-pulse flex justify-between items-center p-4 border rounded">
                          <div className="flex-1">
                            <div className="h-5 bg-gray-300 rounded mb-2"></div>
                            <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                          </div>
                          <div className="w-20 h-6 bg-gray-300 rounded"></div>
                        </div>
                      ))}
                    </div>
                  </div>
                ) : contacts && contacts.length > 0 ? (
                  <div className="divide-y">
                    {contacts.map((contact) => (
                      <div key={contact.id} className="p-6 flex items-center justify-between hover:bg-gray-50">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <h3 className="font-medium">{contact.name}</h3>
                            <span className="text-sm text-muted-foreground">({contact.email})</span>
                          </div>
                          <p className="font-medium text-sm">{contact.subject}</p>
                          <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                            {contact.message}
                          </p>
                          <p className="text-xs text-muted-foreground mt-2">
                            {formatDate(contact.created_at)}
                            {contact.company && ` • ${contact.company}`}
                            {contact.phone && ` • ${contact.phone}`}
                          </p>
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge variant={
                            contact.status === 'new' ? 'destructive' : 
                            contact.status === 'read' ? 'default' : 'secondary'
                          }>
                            {contact.status}
                          </Badge>
                          <Button variant="ghost" size="sm">
                            <Eye className="w-4 h-4" />
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="p-6 text-center">
                    <Mail className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                    <p className="text-muted-foreground">No messages yet</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          {/* Users Tab */}
          <TabsContent value="users" className="space-y-6">
            <div className="flex justify-between items-center">
              <h2 className="text-2xl font-bold">User Management</h2>
              <Link href="/users">
                <Button>
                  <Users className="w-4 h-4 mr-2" />
                  Manage Users
                </Button>
              </Link>
            </div>

            <Card>
              <CardHeader>
                <CardTitle>Quick Access</CardTitle>
                <CardDescription>
                  User management tools and overview
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-center py-8">
                  <Users className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <p className="text-muted-foreground mb-4">
                    Manage system users, roles, and permissions
                  </p>
                  <Link href="/users">
                    <Button>
                      <Users className="w-4 h-4 mr-2" />
                      Go to User Management
                    </Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Settings Tab */}
          <TabsContent value="settings" className="space-y-6">
            <div className="flex justify-between items-center">
              <h2 className="text-2xl font-bold">Settings</h2>
            </div>

            <Card>
              <CardHeader>
                <CardTitle>CMS Configuration</CardTitle>
                <CardDescription>
                  Manage your content management system settings
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <h3 className="font-medium mb-2">API Status</h3>
                    <Badge variant="default">Connected</Badge>
                  </div>
                  <div>
                    <h3 className="font-medium mb-2">Database</h3>
                    <Badge variant="default">CouchDB Connected</Badge>
                  </div>
                  <div>
                    <h3 className="font-medium mb-2">Cache</h3>
                    <Badge variant="default">Valkey Connected</Badge>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </main>
    </div>
  )
}
