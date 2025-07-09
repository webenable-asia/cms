export interface Post {
  id?: string
  rev?: string
  title: string
  content: string
  excerpt: string
  author: string
  status: 'draft' | 'published' | 'scheduled'
  tags: string[]
  categories: string[]
  featured_image?: string
  image_alt?: string
  meta_title?: string
  meta_description?: string
  reading_time?: number
  is_featured: boolean
  view_count: number
  created_at: string
  updated_at: string
  published_at?: string
  scheduled_at?: string
}

export interface Category {
  id?: string
  rev?: string
  name: string
  slug: string
  description?: string
  color?: string
  icon?: string
  post_count: number
  created_at: string
  updated_at: string
}

export interface Contact {
  id?: string
  rev?: string
  name: string
  email: string
  company?: string
  phone?: string
  subject: string
  message: string
  status: 'new' | 'read' | 'replied'
  created_at: string
  read_at?: string
  replied_at?: string
}

export interface User {
  id?: string
  rev?: string
  username: string
  email: string
  role: 'user' // Frontend only handles regular users
  active: boolean
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

// Frontend-specific types
export interface BlogPost {
  id: string
  title: string
  excerpt: string
  content: string
  author: string
  category: string
  tags: string[]
  publishedAt: string
  readTime: string
  slug: string
  featuredImage?: string
  imageAlt?: string
}

export interface ContactForm {
  name: string
  email: string
  company?: string
  phone?: string
  subject: string
  message: string
}
