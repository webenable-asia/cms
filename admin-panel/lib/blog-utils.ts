import { Post, BlogPost } from '@/types/api'

// Convert backend Post to frontend BlogPost
export function postToBlogPost(post: Post): BlogPost {
  return {
    id: post.id || '',
    title: post.title,
    excerpt: post.excerpt,
    content: post.content,
    author: post.author,
    category: post.categories?.[0] || 'General',
    tags: post.tags || [],
    publishedAt: post.published_at || post.created_at,
    readTime: post.reading_time ? `${post.reading_time} min read` : '5 min read',
    slug: generateSlug(post.title),
    featuredImage: post.featured_image,
    imageAlt: post.image_alt,
  }
}

// Generate URL-friendly slug from title
export function generateSlug(title: string): string {
  return title
    .toLowerCase()
    .replace(/[^\w\s-]/g, '') // Remove special characters
    .replace(/[\s_-]+/g, '-') // Replace spaces and underscores with hyphens
    .replace(/^-+|-+$/g, '') // Remove leading/trailing hyphens
}

// Calculate reading time based on content
export function calculateReadingTime(content: string): number {
  const wordsPerMinute = 200
  const words = content.trim().split(/\s+/).length
  return Math.ceil(words / wordsPerMinute)
}

// Format date for display
export function formatDate(dateString: string): string {
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  } catch {
    return 'Unknown date'
  }
}

// Format relative time (e.g., "2 days ago")
export function formatRelativeTime(dateString: string): string {
  try {
    const date = new Date(dateString)
    const now = new Date()
    const diffInMs = now.getTime() - date.getTime()
    const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24))
    
    if (diffInDays === 0) return 'Today'
    if (diffInDays === 1) return 'Yesterday'
    if (diffInDays < 7) return `${diffInDays} days ago`
    if (diffInDays < 30) return `${Math.floor(diffInDays / 7)} weeks ago`
    if (diffInDays < 365) return `${Math.floor(diffInDays / 30)} months ago`
    return `${Math.floor(diffInDays / 365)} years ago`
  } catch {
    return 'Unknown'
  }
}

// Extract excerpt from content if not provided
export function extractExcerpt(content: string, maxLength: number = 160): string {
  // Remove HTML tags and extra whitespace
  const plainText = content
    .replace(/<[^>]*>/g, '')
    .replace(/\s+/g, ' ')
    .trim()
  
  if (plainText.length <= maxLength) {
    return plainText
  }
  
  // Find the last complete sentence within the limit
  const truncated = plainText.substring(0, maxLength)
  const lastSentence = truncated.lastIndexOf('.')
  
  if (lastSentence > maxLength * 0.8) {
    return truncated.substring(0, lastSentence + 1)
  }
  
  // Find the last complete word
  const lastSpace = truncated.lastIndexOf(' ')
  return truncated.substring(0, lastSpace) + '...'
}

// Filter posts by category
export function filterPostsByCategory(posts: BlogPost[], category?: string): BlogPost[] {
  if (!category) return posts
  return posts.filter(post => 
    post.category.toLowerCase() === category.toLowerCase()
  )
}

// Filter posts by tag
export function filterPostsByTag(posts: BlogPost[], tag?: string): BlogPost[] {
  if (!tag) return posts
  return posts.filter(post => 
    post.tags.some(postTag => 
      postTag.toLowerCase() === tag.toLowerCase()
    )
  )
}

// Search posts by title and content
export function searchPosts(posts: BlogPost[], query: string): BlogPost[] {
  if (!query.trim()) return posts
  
  const searchTerm = query.toLowerCase()
  return posts.filter(post => 
    post.title.toLowerCase().includes(searchTerm) ||
    post.excerpt.toLowerCase().includes(searchTerm) ||
    post.content.toLowerCase().includes(searchTerm) ||
    post.tags.some(tag => tag.toLowerCase().includes(searchTerm))
  )
}

// Sort posts by date (newest first)
export function sortPostsByDate(posts: BlogPost[]): BlogPost[] {
  return [...posts].sort((a, b) => 
    new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime()
  )
}

// Get unique categories from posts
export function getUniqueCategories(posts: BlogPost[]): string[] {
  const categories = posts.map(post => post.category)
  return Array.from(new Set(categories)).sort()
}

// Get unique tags from posts
export function getUniqueTags(posts: BlogPost[]): string[] {
  const tags = posts.flatMap(post => post.tags)
  return Array.from(new Set(tags)).sort()
}

// Get related posts based on category and tags
export function getRelatedPosts(posts: BlogPost[], currentPost: BlogPost, limit: number = 3): BlogPost[] {
  const related = posts
    .filter(post => post.id !== currentPost.id)
    .map(post => {
      let score = 0
      
      // Same category gets high score
      if (post.category === currentPost.category) {
        score += 3
      }
      
      // Shared tags get points
      const sharedTags = post.tags.filter(tag => 
        currentPost.tags.includes(tag)
      ).length
      score += sharedTags
      
      return { post, score }
    })
    .filter(item => item.score > 0)
    .sort((a, b) => b.score - a.score)
    .slice(0, limit)
    .map(item => item.post)
  
  return related
}
