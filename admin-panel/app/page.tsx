'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/use-auth-simple'

export default function AdminRoot() {
  const router = useRouter()
  const { isAuthenticated } = useAuth()

  useEffect(() => {
    // Use window.location for more reliable redirects
    if (isAuthenticated) {
      window.location.href = '/admin/dashboard'
    } else {
      window.location.href = '/admin/login'
    }
  }, [isAuthenticated])

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p className="text-gray-600">Loading WebEnable Admin...</p>
      </div>
    </div>
  )
}
