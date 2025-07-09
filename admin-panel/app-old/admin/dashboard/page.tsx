'use client'

import { useState } from 'react'
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
  Trash2
} from 'lucide-react'
import { usePosts, useContacts, useLogout } from '@/hooks/use-api'
import { usePosts as usePostsAPI } from '@/hooks/use-posts'
import { formatDate, formatRelativeTime } from '@/lib/blog-utils'
import Link from 'next/link'

export default function AdminDashboard() {
  const { data: posts, loading: postsLoading } = usePosts()
  const { data: contacts, loading: contactsLoading } = useContacts()
  const { mutate: logout } = useLogout()
  const { deletePost } = usePostsAPI()
  
  const [activeTab, setActiveTab] = useState('overview')

  const handleDeletePost = async (postId: string) => {
    if (window.confirm('Are you sure you want to delete this post?')) {
      try {
        await deletePost(postId)
      } catch (error) {
        console.error('Failed to delete post:', error)
        alert('Failed to delete post. Please try again.')
      }
    }
  }

  const handleLogout = async () => {
    await logout({})
    window.location.href = '/admin'
  }

  const stats = {
    totalPosts: posts?.length || 0,
    publishedPosts: posts?.filter(p => p.status === 'published').length || 0,
    draftPosts: posts?.filter(p => p.status === 'draft').length || 0,
    totalContacts: contacts?.length || 0,
    newContacts: contacts?.filter(c => c.status === 'new').length || 0,
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">CMS Dashboard</h1>
              <p className="text-gray-600">Manage your content and contacts</p>
            </div>
            <div className="flex items-center gap-4">
              <Link href="/" target="_blank">
                <Button variant="outline" size="sm">
                  <Eye className="w-4 h-4 mr-2" />
                  View Site
                </Button>
              </Link>
              <Button variant="outline" size="sm" onClick={handleLogout}>
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
          <TabsList className="grid w-full grid-cols-5">
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="posts">Posts</TabsTrigger>
            <TabsTrigger value="contacts">Contacts</TabsTrigger>
            <TabsTrigger value="users">Users</TabsTrigger>
            <TabsTrigger value="settings">Settings</TabsTrigger>
          </TabsList>

          {/* Overview Tab */}
          <TabsContent value="overview" className="space-y-6">
            {/* Stats Cards */}
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Total Posts</CardTitle>
                  <FileText className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stats.totalPosts}</div>
                  <p className="text-xs text-muted-foreground">
                    {stats.publishedPosts} published
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Draft Posts</CardTitle>
                  <Edit className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stats.draftPosts}</div>
                  <p className="text-xs text-muted-foreground">
                    Waiting to be published
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Total Contacts</CardTitle>
                  <Mail className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stats.totalContacts}</div>
                  <p className="text-xs text-muted-foreground">
                    All time messages
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">New Messages</CardTitle>
                  <Users className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stats.newContacts}</div>
                  <p className="text-xs text-muted-foreground">
                    Unread messages
                  </p>
                </CardContent>
              </Card>
            </div>

            {/* Recent Activity */}
            <div className="grid gap-6 lg:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle>Recent Posts</CardTitle>
                  <CardDescription>Your latest blog posts</CardDescription>
                </CardHeader>
                <CardContent>
                  {postsLoading ? (
                    <div className="space-y-3">
                      {Array.from({ length: 3 }).map((_, i) => (
                        <div key={i} className="animate-pulse">
                          <div className="h-4 bg-gray-300 rounded mb-2"></div>
                          <div className="h-3 bg-gray-200 rounded w-3/4"></div>
                        </div>
                      ))}
                    </div>
                  ) : posts && posts.length > 0 ? (
                    <div className="space-y-4">
                      {posts.slice(0, 5).map((post) => (
                        <div key={post.id} className="flex items-center justify-between">
                          <div>
                            <p className="font-medium truncate">{post.title}</p>
                            <p className="text-xs text-muted-foreground">
                              {formatRelativeTime(post.updated_at)}
                            </p>
                          </div>
                          <Badge variant={post.status === 'published' ? 'default' : 'secondary'}>
                            {post.status}
                          </Badge>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-muted-foreground">No posts yet</p>
                  )}
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Recent Messages</CardTitle>
                  <CardDescription>Latest contact form submissions</CardDescription>
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
              <Link href="/admin/posts/new">
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
                          <Link href={`/admin/posts/${post.id}/edit`}>
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
                    <Link href="/admin/posts/new">
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
              <Link href="/admin/users">
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
                  <Link href="/admin/users">
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
