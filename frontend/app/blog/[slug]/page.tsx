'use client'

import { notFound, useRouter } from "next/navigation"
import { Badge } from "@/components/ui/badge"
import { CalendarDays, Clock, ArrowLeft, Share2, BookmarkPlus } from "lucide-react"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { usePost } from "@/hooks/use-api"
import { formatDate, getRelatedPosts } from "@/lib/blog-utils"
import { useBlogPosts } from "@/hooks/use-api"

interface BlogPostPageProps {
  params: {
    slug: string
  }
}

export default function BlogPostPage({ params }: BlogPostPageProps) {
  const { data: post, loading, error } = usePost(params.slug)
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
    slug: params.slug,
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
      // Fallback to copying URL to clipboard
      navigator.clipboard.writeText(window.location.href)
      // You could show a toast notification here
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
          {/* Featured Image */}
          {post.featured_image && (
            <div className="aspect-video overflow-hidden rounded-lg mb-8">
              <img
                src={post.featured_image}
                alt={post.image_alt || post.title}
                className="w-full h-full object-cover"
              />
            </div>
          )}

          {/* Header */}
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

          {/* Content */}
          <div className="prose prose-lg max-w-none mb-12">
            {/* For now, render as plain text. In production, you'd want to use a markdown renderer */}
            <div style={{ whiteSpace: 'pre-wrap' }}>{post.content}</div>
          </div>

          {/* Tags */}
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

          {/* Meta Information */}
          {(post.meta_title || post.meta_description) && (
            <div className="border-t pt-8 mb-8">
              <h3 className="text-lg font-semibold mb-3">SEO Information</h3>
              {post.meta_title && (
                <div className="mb-2">
                  <strong>Meta Title:</strong> {post.meta_title}
                </div>
              )}
              {post.meta_description && (
                <div>
                  <strong>Meta Description:</strong> {post.meta_description}
                </div>
              )}
            </div>
          )}
        </article>

        {/* Related Posts */}
        {relatedPosts.length > 0 && (
          <section className="border-t pt-12">
            <h2 className="text-2xl font-bold mb-6">Related Posts</h2>
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {relatedPosts.map((relatedPost) => (
                <Card key={relatedPost.id} className="hover:shadow-lg transition-shadow">
                  {relatedPost.featuredImage && (
                    <div className="aspect-video overflow-hidden rounded-t-lg">
                      <img
                        src={relatedPost.featuredImage}
                        alt={relatedPost.imageAlt || relatedPost.title}
                        className="w-full h-full object-cover"
                      />
                    </div>
                  )}
                  <CardHeader>
                    <Badge variant="secondary" className="w-fit mb-2">
                      {relatedPost.category}
                    </Badge>
                    <CardTitle className="line-clamp-2">
                      <Link href={`/blog/${relatedPost.id}`} className="hover:text-primary">
                        {relatedPost.title}
                      </Link>
                    </CardTitle>
                    <CardDescription className="line-clamp-2">
                      {relatedPost.excerpt}
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="flex items-center text-sm text-muted-foreground">
                      <CalendarDays className="w-4 h-4 mr-1" />
                      {formatDate(relatedPost.publishedAt)}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </section>
        )}
      </div>
    </div>
  )
}
