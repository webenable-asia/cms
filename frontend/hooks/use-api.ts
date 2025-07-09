'use client'

import { useState, useEffect } from 'react'
import { postsApi, contactApi, authApi, usersApi } from '@/lib/api'
import { Post, Contact, LoginRequest, BlogPost } from '@/types/api'
import { postToBlogPost } from '@/lib/blog-utils'

// Generic hook for API calls with loading/error states
function useApiCall<T>(
  apiCall: () => Promise<T>,
  dependencies: unknown[] = []
) {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    // Only run API calls in the browser, not during SSR
    if (typeof window === 'undefined') {
      setLoading(false)
      return
    }

    let mounted = true

    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const result = await apiCall()
        if (mounted) {
          setData(result)
        }
      } catch (err) {
        if (mounted) {
          setError(err instanceof Error ? err.message : 'An error occurred')
        }
      } finally {
        if (mounted) {
          setLoading(false)
        }
      }
    }

    fetchData()

    return () => {
      mounted = false
    }
  }, dependencies)

  return { data, loading, error, refetch: () => {
    setLoading(true)
    setError(null)
    apiCall().then(setData).catch(err => 
      setError(err instanceof Error ? err.message : 'An error occurred')
    ).finally(() => setLoading(false))
  }}
}

// Posts hooks
export function usePosts(status?: string) {
  const result = useApiCall(
    () => status ? postsApi.getAll(status) : postsApi.getAll(),
    [status]
  )
  
  return {
    ...result,
    data: result.data || [] // Ensure we always return an array instead of null
  }
}

export function usePost(id: string) {
  // Skip the API call entirely if no ID is provided
  if (!id) {
    return { data: null, loading: false, error: null, refetch: () => {} }
  }
  
  const result = useApiCall(
    () => postsApi.getById(id),
    [id]
  )
  
  return result
}

export function useBlogPosts() {
  const result = useApiCall(
    async () => {
      const posts = await postsApi.getPublished()
      return posts ? posts.map(postToBlogPost) : []
    },
    []
  )

  return {
    ...result,
    data: result.data || [] // Ensure we always return an array
  }
}

// Mutation hook for creating/updating data
function useMutation<TData, TVariables>(
  mutationFn: (variables: TVariables) => Promise<TData>
) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const mutate = async (variables: TVariables): Promise<TData | null> => {
    try {
      setLoading(true)
      setError(null)
      const result = await mutationFn(variables)
      return result
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
      return null
    } finally {
      setLoading(false)
    }
  }

  return { mutate, loading, error }
}

export function useCreatePost() {
  return useMutation((post: Omit<Post, 'id' | 'rev' | 'created_at' | 'updated_at'>) =>
    postsApi.create(post)
  )
}

export function useUpdatePost() {
  return useMutation(({ id, post }: { id: string; post: Partial<Post> }) =>
    postsApi.update(id, post)
  )
}

export function useDeletePost() {
  return useMutation((id: string) => postsApi.delete(id))
}

// Contact hooks
export function useContacts() {
  const result = useApiCall(() => contactApi.getAll(), [])
  
  return {
    ...result,
    data: result.data || [] // Ensure we always return an array instead of null
  }
}

export function useContact(id: string) {
  return useApiCall(() => contactApi.getById(id), [id])
}

export function useSubmitContact() {
  return useMutation((contact: Omit<Contact, 'id' | 'rev' | 'status' | 'created_at' | 'read_at' | 'replied_at'>) =>
    contactApi.submit(contact)
  )
}

export function useUpdateContactStatus() {
  return useMutation(({ id, status }: { id: string; status: string }) =>
    contactApi.updateStatus(id, status)
  )
}

export function useReplyToContact() {
  return useMutation(({ id, subject, message }: { id: string; subject: string; message: string }) =>
    contactApi.reply(id, subject, message)
  )
}

export function useDeleteContact() {
  return useMutation((id: string) => contactApi.delete(id))
}

// Auth hooks
export function useLogin() {
  return useMutation((credentials: LoginRequest) => authApi.login(credentials))
}

export function useLogout() {
  return useMutation(() => authApi.logout())
}

// User Management hooks - Profile only (admin functionality in admin-panel)
export function useUserProfile() {
  return useApiCall(() => usersApi.getProfile())
}

export function useUpdateUserProfile() {
  return useMutation((updates: {
    username?: string
    email?: string
  }) => usersApi.updateProfile(updates))
}
