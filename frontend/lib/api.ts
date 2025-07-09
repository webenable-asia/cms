import { Post, Category, Contact, User, LoginRequest } from '@/types/api'

// Use different API URLs for server-side vs client-side requests
const getApiBaseUrl = () => {
  // Check if we're on the server (Node.js environment)
  if (typeof window === 'undefined') {
    // Server-side: use internal Docker network URL
    return process.env.BACKEND_URL ? `${process.env.BACKEND_URL}/api` : 'http://backend:8080/api'
  }
  // Client-side: use public URL that browser can access
  return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'
}

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

// Token management
let authToken: string | null = null

export const tokenManager = {
  setToken: (token: string) => {
    authToken = token
    if (typeof window !== 'undefined') {
      localStorage.setItem('webenable_token', token)
    }
  },
  
  getToken: (): string | null => {
    if (authToken) return authToken
    if (typeof window !== 'undefined') {
      authToken = localStorage.getItem('webenable_token')
    }
    return authToken
  },
  
  removeToken: () => {
    authToken = null
    if (typeof window !== 'undefined') {
      localStorage.removeItem('webenable_token')
    }
  }
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${getApiBaseUrl()}${endpoint}`
  
  // Get auth token for protected requests
  const token = tokenManager.getToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  
  // Add Authorization header if token exists
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  
  const response = await fetch(url, {
    headers,
    ...options,
  })

  if (!response.ok) {
    // If unauthorized, remove invalid token
    if (response.status === 401) {
      tokenManager.removeToken()
    }
    const errorText = await response.text()
    throw new ApiError(response.status, errorText)
  }

  if (response.status === 204) {
    return {} as T
  }

  return response.json()
}

// Blog/Posts API
export const postsApi = {
  // Get all posts with optional status filter
  getAll: async (status?: string): Promise<Post[]> => {
    const query = status ? `?status=${status}` : ''
    const result = await apiRequest<{ data: Post[]; meta: any }>(`/posts${query}`)
    return result?.data || []
  },

  // Get published posts only (for public blog)
  getPublished: async (): Promise<Post[]> => {
    const result = await apiRequest<{ data: Post[]; meta: any }>('/posts?status=published')
    return result?.data || []
  },

  // Get single post by ID
  getById: (id: string): Promise<Post> => {
    return apiRequest<Post>(`/posts/${id}`)
  },

  // Create new post (requires auth)
  create: (post: Omit<Post, 'id' | 'rev' | 'created_at' | 'updated_at'>): Promise<Post> => {
    return apiRequest<Post>('/posts', {
      method: 'POST',
      body: JSON.stringify(post),
    })
  },

  // Update post (requires auth)
  update: (id: string, post: Partial<Post>): Promise<Post> => {
    return apiRequest<Post>(`/posts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(post),
    })
  },

  // Delete post (requires auth)
  delete: (id: string): Promise<void> => {
    return apiRequest<void>(`/posts/${id}`, {
      method: 'DELETE',
    })
  },
}

// Contact API
export const contactApi = {
  // Submit contact form (public)
  submit: (contact: Omit<Contact, 'id' | 'rev' | 'status' | 'created_at' | 'read_at' | 'replied_at'>): Promise<void> => {
    return apiRequest<void>('/contact', {
      method: 'POST',
      body: JSON.stringify(contact),
    })
  },

  // Get all contacts (requires auth)
  getAll: async (): Promise<Contact[]> => {
    const result = await apiRequest<{ data: Contact[]; meta: any }>('/contacts')
    return result?.data || []
  },

  // Get single contact (requires auth)
  getById: (id: string): Promise<Contact> => {
    return apiRequest<Contact>(`/contacts/${id}`)
  },

  // Update contact status (requires auth)
  updateStatus: (id: string, status: string): Promise<Contact> => {
    return apiRequest<Contact>(`/contacts/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
    })
  },

  // Reply to contact (requires auth)
  reply: (id: string, subject: string, message: string): Promise<void> => {
    return apiRequest<void>(`/contacts/${id}/reply`, {
      method: 'POST',
      body: JSON.stringify({ subject, message }),
    })
  },

  // Delete contact (requires auth)
  delete: (id: string): Promise<void> => {
    return apiRequest<void>(`/contacts/${id}`, {
      method: 'DELETE',
    })
  },
}

// Auth API
export const authApi = {
  // Login
  login: async (credentials: LoginRequest): Promise<{ token: string; user: User }> => {
    const result = await apiRequest<{ token: string; user: User }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    })
    
    // Store the token after successful login
    if (result.token) {
      tokenManager.setToken(result.token)
    }
    
    return result
  },

  // Get current user
  me: (): Promise<{ user: User }> => {
    return apiRequest<{ user: User }>('/auth/me')
  },

  // Logout
  logout: async (): Promise<void> => {
    try {
      await apiRequest<void>('/auth/logout', {
        method: 'POST',
      })
    } finally {
      // Always remove token on logout, even if API call fails
      tokenManager.removeToken()
    }
  },
}

// User Management API - Basic user operations only
// Admin user management moved to admin-panel service
export const usersApi = {
  // Get current user profile
  getProfile: async (): Promise<User | null> => {
    try {
      const result = await apiRequest<User>('/users/profile')
      return result
    } catch (error) {
      console.error('Failed to get user profile:', error)
      return null
    }
  },

  // Update current user profile
  updateProfile: (updates: {
    username?: string
    email?: string
  }): Promise<User> => {
    return apiRequest<User>('/users/profile', {
      method: 'PUT',
      body: JSON.stringify(updates),
    })
  },
}

// Health check
export const healthApi = {
  check: (): Promise<{ status: string; cache: string }> => {
    return apiRequest<{ status: string; cache: string }>('/health')
  },
}

export { ApiError }
