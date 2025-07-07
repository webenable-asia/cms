'use client'

import { use } from "react"
import { notFound } from "next/navigation"
import { Badge } from "@/components/ui/badge"
import { CalendarDays, Clock, ArrowLeft, Share2 } from "lucide-react"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { usePost, useBlogPosts } from "@/hooks/use-api"
import { formatDate, getRelatedPosts } from "@/lib/blog-utils"
import { MarkdownContent } from "@/components/markdown-content"

interface BlogPostPageProps {
  params: Promise<{
    slug: string
  }>
}

export default function BlogPostPage({ params }: BlogPostPageProps) {
  const { slug } = use(params)
  
  // Always call hooks in the same order
  const { data: post, loading, error } = usePost(slug)
  const { data: allPosts } = useBlogPosts()

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto">
          <div className="animate-pulse">
            <div className="h-8 bg-gray-300 rounded mb-8 w-32"></div>
            <div className="h-12 bg-gray-300 rounded mb-4"></div>
            <div className="h-6 bg-gray-300 rounded mb-8 w-3/4"></div>
            <div className="space-y-4">
              {Array.from({ length: 10 }).map((_, i) => (
                <div key={i} className="h-4 bg-gray-300 rounded"></div>
              ))}
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !post) {
    return notFound()
  }

  const relatedPosts = allPosts ? getRelatedPosts(allPosts, {
    id: post.id || '',
    title: post.title,
    excerpt: post.excerpt,
    content: post.content,
    author: post.author,
    category: post.categories?.[0] || 'General',
    tags: post.tags || [],
    publishedAt: post.published_at || post.created_at,
    readTime: post.reading_time ? `${post.reading_time} min read` : '5 min read',
    slug: slug,
    featuredImage: post.featured_image,
    imageAlt: post.image_alt,
  }) : []

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: post.title,
          text: post.excerpt,
          url: window.location.href,
        })
      } catch (err) {
        console.log('Error sharing:', err)
      }
    } else {
      navigator.clipboard.writeText(window.location.href)
    }
  }

  return (
    <div className="container mx-auto px-4 py-16">
      <div className="max-w-4xl mx-auto">
        <Link href="/blog">
          <Button variant="ghost" className="mb-8">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Blog
          </Button>
        </Link>

        <article>
          <header className="mb-8">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-4">
                {post.categories?.map((category) => (
                  <Badge key={category} variant="secondary">{category}</Badge>
                ))}
                <div className="flex items-center text-sm text-muted-foreground">
                  <CalendarDays className="w-4 h-4 mr-1" />
                  {formatDate(post.published_at || post.created_at)}
                </div>
                <div className="flex items-center text-sm text-muted-foreground">
                  <Clock className="w-4 h-4 mr-1" />
                  {post.reading_time ? `${post.reading_time} min read` : '5 min read'}
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Button variant="outline" size="sm" onClick={handleShare}>
                  <Share2 className="w-4 h-4 mr-2" />
                  Share
                </Button>
              </div>
            </div>
            
            <h1 className="text-4xl font-bold mb-4">{post.title}</h1>
            
            {post.excerpt && (
              <p className="text-xl text-muted-foreground mb-4">{post.excerpt}</p>
            )}
            
            <div className="flex items-center justify-between text-sm text-muted-foreground border-b pb-4">
              <div>By {post.author}</div>
              {post.view_count && (
                <div>{post.view_count} views</div>
              )}
            </div>
          </header>

          <MarkdownContent 
            content={post.content} 
            className="mb-12"
          />

          {post.tags && post.tags.length > 0 && (
            <div className="mb-8">
              <h3 className="text-lg font-semibold mb-3">Tags</h3>
              <div className="flex flex-wrap gap-2">
                {post.tags.map((tag) => (
                  <Badge key={tag} variant="outline">
                    {tag}
                  </Badge>
                ))}
              </div>
            </div>
          )}
        </article>
      </div>
    </div>
  )
}
