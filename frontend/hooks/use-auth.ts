'use client'

import React, { useState, useEffect, useContext, createContext } from 'react'
import { authApi, tokenManager } from '@/lib/api'
import { User } from '@/types/api'

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
  checkAuth: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const checkAuth = async () => {
    try {
      setIsLoading(true)
      const token = tokenManager.getToken()
      
      if (!token) {
        setUser(null)
        return
      }

      const result = await authApi.me()
      
      if (result?.user && result.user.role === 'admin' && result.user.active) {
        setUser(result.user)
      } else {
        setUser(null)
        tokenManager.removeToken()
      }
    } catch (error) {
      console.warn('Auth verification failed:', error)
      setUser(null)
      tokenManager.removeToken()
    } finally {
      setIsLoading(false)
    }
  }

  const login = async (username: string, password: string) => {
    try {
      const result = await authApi.login({ username, password })
      
      if (result?.user && result.user.role === 'admin' && result.user.active) {
        setUser(result.user)
      } else {
        setUser(null)
        tokenManager.removeToken()
        throw new Error('Invalid user role or inactive account')
      }
    } catch (error) {
      setUser(null)
      tokenManager.removeToken()
      throw error
    }
  }

  const logout = async () => {
    try {
      await authApi.logout()
    } catch (error) {
      console.warn('Logout API call failed:', error)
    } finally {
      setUser(null)
      tokenManager.removeToken()
    }
  }

  // Check for existing token on mount
  useEffect(() => {
    const initAuth = async () => {
      const token = tokenManager.getToken()
      
      if (token) {
        await checkAuth()
      } else {
        setUser(null)
        setIsLoading(false)
      }
    }
    
    initAuth()
  }, [])

  const value = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
    checkAuth
  }

  return React.createElement(
    AuthContext.Provider,
    { value },
    children
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
