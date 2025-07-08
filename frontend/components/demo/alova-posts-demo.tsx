'use client'

import React, { useState, useEffect } from 'react'
import { postsApi } from '@/lib/api/posts'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Loader2, RefreshCw } from 'lucide-react'

interface Post {
  id: string
  title: string
  content: string
  published: boolean
  created_at: string
  updated_at: string
}

export function AlovaPostsDemo() {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchPosts = async () => {
    try {
      setLoading(true)
      setError(null)
      
      // Use Alova to fetch posts
      const response = await postsApi.getAll().send()
      
      // Handle the response based on your API structure
      const postsData = (response as any)?.data || response || []
      setPosts(Array.isArray(postsData) ? postsData : [])
    } catch (err) {
      console.error('Error fetching posts:', err)
      setError(err instanceof Error ? err.message : 'Failed to fetch posts')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchPosts()
  }, [])

  if (loading) {
    return (
      <Card className="w-full max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle>Posts (Alova Demo)</CardTitle>
          <CardDescription>Loading posts using Alova data fetching...</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <Loader2 className="h-8 w-8 animate-spin text-blue-500" />
            <span className="ml-2 text-gray-600">Loading posts...</span>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (error) {
    return (
      <Card className="w-full max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle>Posts (Alova Demo)</CardTitle>
          <CardDescription>Error loading posts</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8">
            <p className="text-red-600 mb-4">{error}</p>
            <Button onClick={fetchPosts} variant="outline">
              <RefreshCw className="h-4 w-4 mr-2" />
              Try Again
            </Button>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="w-full max-w-4xl mx-auto">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Posts (Alova Demo)</CardTitle>
            <CardDescription>
              Demonstrating Alova data fetching with caching ({posts.length} posts)
            </CardDescription>
          </div>
          <Button onClick={fetchPosts} variant="outline" size="sm">
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {posts.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-600">No posts found</p>
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <Card key={post.id} className="border border-gray-200">
                <CardContent className="pt-6">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="font-semibold text-lg mb-2">{post.title}</h3>
                      <p className="text-gray-600 mb-3 line-clamp-2">
                        {post.content.substring(0, 150)}
                        {post.content.length > 150 ? '...' : ''}
                      </p>
                      <div className="flex items-center gap-2 text-sm text-gray-500">
                        <span>Created: {new Date(post.created_at).toLocaleDateString()}</span>
                        <span>•</span>
                        <span>Updated: {new Date(post.updated_at).toLocaleDateString()}</span>
                      </div>
                    </div>
                    <div className="ml-4">
                      <Badge variant={post.published ? 'default' : 'secondary'}>
                        {post.published ? 'Published' : 'Draft'}
                      </Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
        
        <div className="mt-6 p-4 bg-blue-50 rounded-lg">
          <h4 className="font-medium text-blue-900 mb-2">Alova Features Demonstrated:</h4>
          <ul className="text-sm text-blue-800 space-y-1">
            <li>✅ Automatic caching (60 seconds for posts list)</li>
            <li>✅ Global error handling and response interceptors</li>
            <li>✅ Automatic authentication header injection</li>
            <li>✅ TypeScript support with proper typing</li>
            <li>✅ Request deduplication (try clicking refresh quickly)</li>
          </ul>
        </div>
      </CardContent>
    </Card>
  )
}
