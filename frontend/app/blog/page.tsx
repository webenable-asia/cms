'use client'

import Link from "next/link"
import { useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { CalendarDays, Clock, Search, Filter } from "lucide-react"
import { useBlogPosts } from "@/hooks/use-api"
import { 
  searchPosts, 
  filterPostsByCategory, 
  sortPostsByDate, 
  getUniqueCategories,
  formatDate 
} from "@/lib/blog-utils"

export default function BlogPage() {
  const { data: posts, loading, error } = useBlogPosts()
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedCategory, setSelectedCategory] = useState('')

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold mb-4">Our Blog</h1>
            <p className="text-xl text-muted-foreground">Insights, tips, and stories from our team</p>
          </div>
          <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <Card key={i} className="animate-pulse">
                <CardHeader>
                  <div className="h-4 bg-gray-300 rounded mb-2"></div>
                  <div className="h-6 bg-gray-300 rounded mb-2"></div>
                  <div className="h-4 bg-gray-300 rounded"></div>
                </CardHeader>
                <CardContent>
                  <div className="h-4 bg-gray-300 rounded"></div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-4xl font-bold mb-4">Our Blog</h1>
          <p className="text-xl text-red-600 mb-4">Error loading posts: {error}</p>
          <Button onClick={() => window.location.reload()}>
            Try Again
          </Button>
        </div>
      </div>
    )
  }

  if (!posts || posts.length === 0) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-4xl font-bold mb-4">Our Blog</h1>
          <p className="text-xl text-muted-foreground">No posts available at the moment.</p>
        </div>
      </div>
    )
  }

  // Filter and search posts
  let filteredPosts = posts
  if (searchQuery) {
    filteredPosts = searchPosts(filteredPosts, searchQuery)
  }
  if (selectedCategory) {
    filteredPosts = filterPostsByCategory(filteredPosts, selectedCategory)
  }
  filteredPosts = sortPostsByDate(filteredPosts)

  const categories = getUniqueCategories(posts)

  return (
    <div className="container mx-auto px-4 py-16">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold mb-4">Our Blog</h1>
          <p className="text-xl text-muted-foreground">Insights, tips, and stories from our team</p>
        </div>

        {/* Search and Filter */}
        <div className="mb-8 space-y-4 md:space-y-0 md:flex md:items-center md:gap-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
            <Input
              placeholder="Search posts..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
          <div className="flex items-center gap-2">
            <Filter className="w-4 h-4 text-muted-foreground" />
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="px-3 py-2 border border-input bg-background rounded-md"
            >
              <option value="">All Categories</option>
              {categories.map((category) => (
                <option key={category} value={category}>
                  {category}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Results count */}
        <div className="mb-6 text-sm text-muted-foreground">
          Showing {filteredPosts.length} of {posts.length} posts
          {searchQuery && ` for "${searchQuery}"`}
          {selectedCategory && ` in "${selectedCategory}"`}
        </div>

        {/* Posts grid */}
        {filteredPosts.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-xl text-muted-foreground">No posts found matching your criteria.</p>
            <Button 
              onClick={() => {
                setSearchQuery('')
                setSelectedCategory('')
              }}
              variant="outline"
              className="mt-4"
            >
              Clear Filters
            </Button>
          </div>
        ) : (
          <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
            {filteredPosts.map((post) => (
              <Card key={post.id} className="hover:shadow-lg transition-shadow">
                {post.featuredImage && (
                  <div className="aspect-video overflow-hidden rounded-t-lg">
                    <img
                      src={post.featuredImage}
                      alt={post.imageAlt || post.title}
                      className="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
                    />
                  </div>
                )}
                <CardHeader>
                  <div className="flex items-center justify-between mb-2">
                    <Badge variant="secondary">{post.category}</Badge>
                    <div className="flex items-center text-sm text-muted-foreground">
                      <Clock className="w-4 h-4 mr-1" />
                      {post.readTime}
                    </div>
                  </div>
                  <CardTitle className="line-clamp-2">
                    <Link href={`/blog/${post.id}`} className="hover:text-primary">
                      {post.title}
                    </Link>
                  </CardTitle>
                  <CardDescription className="line-clamp-3">{post.excerpt}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center text-sm text-muted-foreground">
                      <CalendarDays className="w-4 h-4 mr-1" />
                      {formatDate(post.publishedAt)}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      by {post.author}
                    </div>
                  </div>
                  {post.tags.length > 0 && (
                    <div className="mt-3 flex flex-wrap gap-1">
                      {post.tags.slice(0, 3).map((tag) => (
                        <span
                          key={tag}
                          className="text-xs px-2 py-1 bg-muted rounded"
                        >
                          {tag}
                        </span>
                      ))}
                      {post.tags.length > 3 && (
                        <span className="text-xs px-2 py-1 bg-muted rounded">
                          +{post.tags.length - 3}
                        </span>
                      )}
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
