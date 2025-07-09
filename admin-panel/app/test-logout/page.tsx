'use client'

import { useAuth } from '@/hooks/use-auth-simple'
import { Button } from '@/components/ui/button'
import { useRouter } from 'next/navigation'

export default function TestLogout() {
  const { logout, isAuthenticated } = useAuth()
  const router = useRouter()

  const handleTestLogout = async () => {
    console.log('🔄 Testing logout...')
    try {
      await logout()
      console.log('✅ Logout successful')
      router.push('/login')
    } catch (error) {
      console.error('❌ Logout failed:', error)
      // Force cleanup
      localStorage.clear()
      router.push('/login')
    }
  }

  const checkLocalStorage = () => {
    console.log('📊 Current localStorage:')
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      console.log(`  ${key}: ${localStorage.getItem(key)}`)
    }
  }

  if (!isAuthenticated) {
    return (
      <div className="p-8">
        <h1 className="text-2xl font-bold mb-4">Not Authenticated</h1>
        <p>Please login first</p>
        <Button onClick={() => router.push('/login')} className="mt-4">
          Go to Login
        </Button>
      </div>
    )
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Logout Test Page</h1>
      
      <div className="bg-gray-100 p-4 rounded mb-4">
        <h2 className="font-semibold mb-2">Authentication Status:</h2>
        <p>Authenticated: {isAuthenticated ? 'Yes' : 'No'}</p>
      </div>

      <div className="space-y-4">
        <Button onClick={handleTestLogout} variant="destructive">
          🚪 Test Logout
        </Button>
        
        <Button onClick={checkLocalStorage} variant="outline">
          📊 Check localStorage
        </Button>
        
        <Button onClick={() => router.push('/dashboard')} variant="default">
          🏠 Back to Dashboard
        </Button>
      </div>

      <div className="mt-8 text-sm text-gray-600">
        <p>Open browser console to see logout process details</p>
      </div>
    </div>
  )
}
