'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/use-auth'

export default function AdminRoot() {
  const router = useRouter()
  const { isAuthenticated } = useAuth()

  useEffect(() => {
    // Delay navigation to prevent SSR hydration issues
    const timeoutId = setTimeout(() => {
      if (isAuthenticated) {
        router.push('/admin/dashboard')
      } else {
        router.push('/admin/login')
      }
    }, 100)

    return () => clearTimeout(timeoutId)
  }, [isAuthenticated, router])

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
        <p>Loading WebEnable Admin...</p>
      </div>
    </div>
  )
}
