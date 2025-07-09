'use client'

import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react'

interface User {
  id: string
  email: string
  name: string
  role: string
}

interface AuthState {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
}

interface AuthContextType extends AuthState {
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  refreshAuth: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const isAuthenticated = !!user

  const login = async (email: string, password: string) => {
    setIsLoading(true)
    try {
      // TODO: Implement actual login logic with API
      console.log('Login:', email, password)
      // Mock user for now
      setUser({
        id: '1',
        email,
        name: 'User',
        role: 'user'
      })
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }

  const logout = () => {
    setUser(null)
    // TODO: Clear tokens, etc.
  }

  const refreshAuth = async () => {
    setIsLoading(true)
    try {
      // TODO: Implement token refresh logic
      console.log('Refreshing auth...')
    } catch (error) {
      console.error('Auth refresh failed:', error)
      setUser(null)
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    // TODO: Check for existing auth tokens on mount
    setIsLoading(false)
  }, [])

  const value = {
    user,
    isLoading,
    isAuthenticated,
    login,
    logout,
    refreshAuth
  }

  return React.createElement(AuthContext.Provider, { value }, children)
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
