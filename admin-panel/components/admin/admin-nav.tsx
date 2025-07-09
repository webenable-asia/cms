'use client'

import { useAuth } from '@/hooks/use-auth-simple'
import { Button } from '@/components/ui/button'
import { 
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { LogOut, User, Shield } from 'lucide-react'
import { useRouter } from 'next/navigation'

export function AdminNav() {
  const { logout } = useAuth()
  const router = useRouter()

  const handleLogout = async () => {
    console.log('🔄 Logout button clicked')
    try {
      console.log('📤 Calling logout API...')
      await logout()
      console.log('✅ Logout API successful')
      
      // Force clear localStorage
      if (typeof window !== 'undefined') {
        console.log('🧹 Clearing localStorage...')
        localStorage.removeItem('webenable_token')
        localStorage.removeItem('webenable_token_expiry')
        localStorage.clear()
        console.log('✅ localStorage cleared')
      }
      
      // Use window.location for immediate redirect
      console.log('🔄 Redirecting to login...')
      window.location.href = '/admin/login'
      
    } catch (error) {
      console.error('❌ Logout failed:', error)
      // Force redirect even if logout API fails
      if (typeof window !== 'undefined') {
        localStorage.clear()
        console.log('🧹 Force cleared localStorage')
      }
      window.location.href = '/admin/login'
    }
  }

  return (
    <div className="border-b bg-white">
      <div className="flex h-16 items-center px-4 justify-between">
        <div className="flex items-center space-x-4">
          <Shield className="h-6 w-6 text-blue-600" />
          <h1 className="text-lg font-semibold">Admin Panel</h1>
        </div>
        
        <Button variant="outline" onClick={handleLogout}>
          <LogOut className="mr-2 h-4 w-4" />
          Log out
        </Button>
      </div>
    </div>
  )
}
