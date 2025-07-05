'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { api } from '../../lib/api'
import { Navigation } from '@/components/navigation'
import { 
  Calendar, 
  Clock, 
  User, 
  ArrowRight, 
  BookOpen,
  Mail,
  Phone,
  MapPin,
  Filter,
  Search
} from 'lucide-react'

interface Post {
  id: string
  title: string
  excerpt: string
  content: string
  author: string
  created_at: string
  published_at?: string
  tags: string[]
  categories?: string[]
  category?: string
  read_time?: string
  reading_time?: number
  featured_image?: string
  image_alt?: string
  is_featured?: boolean
  view_count?: number
  status: string
}

export default function BlogPage() {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const [categories, setCategories] = useState<string[]>(['all'])

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        // Fetch posts from our Next.js API route
        const response = await fetch('/api/posts?status=published')
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        
        const data = await response.json()
        const fetchedPosts = data.posts || []
        
        // Transform API posts to match our interface
        const transformedPosts = fetchedPosts.map((post: any) => ({
          id: post._id || post.id,
          title: post.title || '',
          excerpt: post.excerpt || '',
          content: post.content || '',
          author: post.author || 'Unknown',
          created_at: post.created_at || new Date().toISOString(),
          tags: post.tags || [],
          categories: post.categories || [],
          category: post.categories?.[0] || '',
          read_time: post.reading_time ? `${post.reading_time} min read` : '',
          reading_time: post.reading_time || Math.ceil((post.content || '').split(/\s+/).length / 200),
          featured_image: post.featured_image?.url || '',
          image_alt: post.featured_image?.alt || '',
          is_featured: post.is_featured || false,
          view_count: post.view_count || 0,
          status: post.status || 'published'
        }))
        
        setPosts(transformedPosts)
        
        // Extract unique categories from posts
        const uniqueCategories = ['all', ...Array.from(new Set(
          transformedPosts.flatMap((post: any) => post.categories || [])
        )).filter((cat): cat is string => typeof cat === 'string')]
        setCategories(uniqueCategories)
        
      } catch (error) {
        console.error('Error fetching posts:', error)
        setPosts([])
        setCategories(['all'])
      } finally {
        setLoading(false)
      }
    }

    fetchPosts()
  }, [])

  const filteredPosts = selectedCategory === 'all' 
    ? posts 
    : posts.filter(post => post.categories?.includes(selectedCategory))

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-xl">Loading posts...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <Navigation />

      {/* Hero Section */}
      <section className="pt-16 pb-20 lg:pt-24 lg:pb-28">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-4xl mx-auto">
            <h1 className="text-4xl md:text-6xl font-bold text-foreground mb-6">
              Our Blog
            </h1>
            <p className="text-xl md:text-2xl text-muted-foreground mb-8 max-w-3xl mx-auto leading-relaxed">
              Insights, tips, and stories from our team
            </p>
          </div>
        </div>
      </section>

      {/* Category Filter */}
      <section className="pb-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-wrap justify-center gap-4 mb-8">
            {categories.map((category) => (
              <Button
                key={category}
                variant={selectedCategory === category ? "default" : "outline"}
                onClick={() => setSelectedCategory(category)}
                className="capitalize"
              >
                {category}
              </Button>
            ))}
          </div>
        </div>
      </section>

      {/* Blog Posts */}
      <section className="pb-20">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          {filteredPosts.length === 0 ? (
            <div className="text-center py-20">
              <div className="w-24 h-24 bg-muted rounded-full flex items-center justify-center mx-auto mb-6">
                <BookOpen className="w-12 h-12 text-muted-foreground" />
              </div>
              <h3 className="text-2xl font-bold text-foreground mb-4">No posts found</h3>
              <p className="text-muted-foreground mb-8 max-w-md mx-auto">
                {selectedCategory === 'all' 
                  ? "We're working on some amazing content for you. Check back soon for our latest insights and tutorials!"
                  : `No posts found in the ${selectedCategory} category. Try selecting a different category.`
                }
              </p>
              <Button onClick={() => setSelectedCategory('all')}>
                View All Posts
              </Button>
            </div>
          ) : (
            <div className="grid lg:grid-cols-3 md:grid-cols-2 gap-8">
              {filteredPosts.map((post, index) => (
                <Card key={post.id || `post-${index}`} className="border-border hover:shadow-lg transition-all duration-300 hover:-translate-y-1 overflow-hidden">
                  <CardHeader className="pb-4">
                    {post.categories && post.categories.length > 0 && (
                      <Badge variant="secondary" className="w-fit mb-3">
                        {post.categories[0]}
                      </Badge>
                    )}
                    <h2 className="text-xl font-bold text-card-foreground mb-3 line-clamp-2">
                      <Link 
                        href={`/blog/${post.id}`}
                        className="hover:text-primary transition-colors"
                      >
                        {post.title}
                      </Link>
                    </h2>
                  </CardHeader>
                  
                  <CardContent className="pt-0">
                    <p className="text-muted-foreground mb-6 leading-relaxed line-clamp-3">
                      {post.excerpt}
                    </p>
                    
                    <div className="flex items-center justify-between text-sm text-muted-foreground mb-4">
                      <div className="flex items-center">
                        <User className="h-4 w-4 mr-2" />
                        <span>{post.author}</span>
                      </div>
                      {post.read_time && (
                        <div className="flex items-center">
                          <Clock className="h-4 w-4 mr-1" />
                          <span>{post.read_time}</span>
                        </div>
                      )}
                    </div>

                    <div className="flex items-center justify-between text-sm text-muted-foreground mb-6">
                      <div className="flex items-center">
                        <Calendar className="h-4 w-4 mr-2" />
                        <span>{new Date(post.created_at).toLocaleDateString('en-US', { 
                          month: 'short', 
                          day: 'numeric',
                          year: 'numeric'
                        })}</span>
                      </div>
                    </div>
                    
                    {post.tags && post.tags.length > 0 && (
                      <div className="flex flex-wrap gap-2 mb-6">
                        {post.tags.slice(0, 3).map((tag, index) => (
                          <Badge key={index} variant="outline" className="text-xs">
                            {tag}
                          </Badge>
                        ))}
                        {post.tags.length > 3 && (
                          <Badge variant="outline" className="text-xs">
                            +{post.tags.length - 3} more
                          </Badge>
                        )}
                      </div>
                    )}

                    <Link href={`/blog/${post.id}`}>
                      <Button variant="ghost" className="w-full group">
                        Read More
                        <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                      </Button>
                    </Link>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-secondary text-secondary-foreground py-16">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div className="lg:col-span-2">
              <h3 className="text-2xl font-bold mb-4">WebEnable</h3>
              <p className="text-muted-foreground mb-6 max-w-md">
                Building exceptional digital experiences that help businesses grow and succeed.
              </p>
              <div className="space-y-2">
                <div className="flex items-center text-muted-foreground">
                  <Mail className="h-4 w-4 mr-2" />
                  hello@webenable.asia
                </div>
                <div className="flex items-center text-muted-foreground">
                  <Phone className="h-4 w-4 mr-2" />
                  +66 (0) 123-456-789
                </div>
                <div className="flex items-center text-muted-foreground">
                  <MapPin className="h-4 w-4 mr-2" />
                  Bangkok, Thailand
                </div>
              </div>
            </div>
            
            <div>
              <h4 className="text-lg font-semibold mb-4">Company</h4>
              <ul className="space-y-2">
                <li><Link href="/about" className="text-muted-foreground hover:text-foreground transition-colors">About Us</Link></li>
                <li><Link href="/blog" className="text-muted-foreground hover:text-foreground transition-colors">Blog</Link></li>
                <li><Link href="/contact" className="text-muted-foreground hover:text-foreground transition-colors">Contact</Link></li>
              </ul>
            </div>
            
            <div>
              <h4 className="text-lg font-semibold mb-4">Services</h4>
              <ul className="space-y-2">
                <li className="text-muted-foreground">Web Development</li>
                <li className="text-muted-foreground">Mobile Apps</li>
                <li className="text-muted-foreground">UI/UX Design</li>
                <li className="text-muted-foreground">Consulting</li>
              </ul>
            </div>
          </div>
          
          <div className="border-t border-border mt-12 pt-8">
            <div className="flex flex-col md:flex-row justify-between items-center">
              <p className="text-muted-foreground mb-4 md:mb-0">
                © 2025 WebEnable. All rights reserved.
              </p>
              <div className="flex space-x-6">
                <Link href="#" className="text-muted-foreground hover:text-foreground transition-colors">Privacy Policy</Link>
                <Link href="#" className="text-muted-foreground hover:text-foreground transition-colors">Terms of Service</Link>
                <Link href="#" className="text-muted-foreground hover:text-foreground transition-colors">Cookie Policy</Link>
              </div>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
