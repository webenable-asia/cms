'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/use-auth'

interface ProtectedRouteProps {
  children: React.ReactNode
  fallback?: React.ReactNode
  requiredRole?: string
}

export function ProtectedRoute({ children, fallback, requiredRole = 'admin' }: ProtectedRouteProps) {
  const { user, isAuthenticated, isLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading) {
      if (!isAuthenticated) {
        // Not authenticated - redirect to login
        router.push('/admin/login')
        return
      }
      
      if (requiredRole && user?.role !== requiredRole) {
        // Authenticated but wrong role - redirect to login
        router.push('/admin/login')
        return
      }
      
      if (user && !user.active) {
        // User exists but is inactive - redirect to login
        router.push('/admin/login')
        return
      }
    }
  }, [isAuthenticated, isLoading, user, requiredRole, router])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Verifying authentication...</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated || !user) {
    return fallback || (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-900 mb-2">Access Denied</h2>
          <p className="text-gray-600">Redirecting to login...</p>
        </div>
      </div>
    )
  }

  if (requiredRole && user.role !== requiredRole) {
    return fallback || (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-red-600 mb-2">Insufficient Permissions</h2>
          <p className="text-gray-600">You don't have the required permissions to access this area.</p>
        </div>
      </div>
    )
  }

  if (!user.active) {
    return fallback || (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-orange-600 mb-2">Account Inactive</h2>
          <p className="text-gray-600">Your account has been deactivated. Please contact support.</p>
        </div>
      </div>
    )
  }

  return <>{children}</>
}
