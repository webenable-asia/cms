import { Post, Category, Contact, User, LoginRequest } from '@/types/api'

// Use different API URLs for server-side vs client-side requests
const getApiBaseUrl = () => {
  // Check if we're on the server (Node.js environment)
  if (typeof window === 'undefined') {
    // Server-side: use internal Docker network URL
    return process.env.BACKEND_URL ? `${process.env.BACKEND_URL}/api` : 'http://backend:8080/api'
  }
  // Client-side: use full URL to ensure it works
  return `${window.location.protocol}//${window.location.host}/api`
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
      // Also set an expiry timestamp for automatic cleanup
      const expiryTime = Date.now() + (24 * 60 * 60 * 1000) // 24 hours
      localStorage.setItem('webenable_token_expiry', expiryTime.toString())
    }
  },
  
  getToken: (): string | null => {
    if (authToken) {
      // Check if token is expired
      if (typeof window !== 'undefined') {
        const expiry = localStorage.getItem('webenable_token_expiry')
        if (expiry && Date.now() > parseInt(expiry)) {
          // Token expired, remove it
          tokenManager.removeToken()
          return null
        }
      }
      return authToken
    }
    
    if (typeof window !== 'undefined') {
      authToken = localStorage.getItem('webenable_token')
      // Check expiry for stored token too
      const expiry = localStorage.getItem('webenable_token_expiry')
      if (expiry && Date.now() > parseInt(expiry)) {
        tokenManager.removeToken()
        return null
      }
    }
    return authToken
  },
  
  removeToken: () => {
    authToken = null
    if (typeof window !== 'undefined') {
      localStorage.removeItem('webenable_token')
      localStorage.removeItem('webenable_token_expiry')
      // Clear any other admin-related localStorage items
      const keysToRemove = []
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i)
        if (key && (key.startsWith('webenable_') || key.startsWith('admin_'))) {
          keysToRemove.push(key)
        }
      }
      keysToRemove.forEach(key => localStorage.removeItem(key))
    }
  },
  
  isTokenValid: (): boolean => {
    const token = tokenManager.getToken()
    return !!token
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
    'Cache-Control': 'no-cache, no-store, must-revalidate',
    'Pragma': 'no-cache',
    'Expires': '0',
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
    // Add cache-busting parameter to ensure fresh data
    const cacheBuster = `_t=${Date.now()}`
    const query = status ? `?status=${status}&${cacheBuster}` : `?${cacheBuster}`
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
    // Add cache-busting parameter to ensure fresh data
    const cacheBuster = `?_t=${Date.now()}`
    const result = await apiRequest<{ data: Contact[]; meta: any }>(`/contacts${cacheBuster}`)
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

// User Management API
export const usersApi = {
  // Get all users (admin only)
  getAll: async (): Promise<User[]> => {
    const result = await apiRequest<{ data: User[]; meta: any }>('/users')
    return result?.data || []
  },

  // Get single user by ID (admin only)
  getById: (id: string): Promise<User> => {
    return apiRequest<User>(`/users/${id}`)
  },

  // Create new user (admin only)
  create: (user: {
    username: string
    email: string
    password: string
    role: string
    active: boolean
  }): Promise<User> => {
    return apiRequest<User>('/users', {
      method: 'POST',
      body: JSON.stringify(user),
    })
  },

  // Update user (admin only)
  update: (id: string, updates: {
    username?: string
    email?: string
    password?: string
    role?: string
    active?: boolean
  }): Promise<User> => {
    return apiRequest<User>(`/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updates),
    })
  },

  // Delete user (admin only)
  delete: (id: string): Promise<void> => {
    return apiRequest<void>(`/users/${id}`, {
      method: 'DELETE',
    })
  },

  // Get user statistics (admin only)
  getStats: (): Promise<{
    total_users: number
    active_users: number
    admin_users: number
    editor_users: number
    author_users: number
    last_updated: string
  }> => {
    return apiRequest('/users/stats')
  },
}

// Health check
export const healthApi = {
  check: (): Promise<{ status: string; cache: string }> => {
    return apiRequest<{ status: string; cache: string }>('/health')
  },
}

export { ApiError }
