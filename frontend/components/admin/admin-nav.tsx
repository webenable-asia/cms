'use client'

import { useAuth } from '@/hooks/use-auth'
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
  const { user, logout } = useAuth()
  const router = useRouter()

  const handleLogout = async () => {
    try {
      await logout()
      router.push('/admin/login')
    } catch (error) {
      console.error('Logout failed:', error)
      // Force redirect even if logout API fails
      router.push('/admin/login')
    }
  }

  if (!user) {
    return null
  }

  const initials = user.username.substring(0, 2).toUpperCase()

  return (
    <div className="border-b bg-white">
      <div className="flex h-16 items-center px-4 justify-between">
        <div className="flex items-center space-x-4">
          <Shield className="h-6 w-6 text-blue-600" />
          <h1 className="text-lg font-semibold">Admin Panel</h1>
        </div>
        
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="relative h-8 w-8 rounded-full">
              <Avatar className="h-8 w-8">
                <AvatarFallback className="bg-blue-100 text-blue-600">
                  {initials}
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-56" align="end" forceMount>
            <DropdownMenuLabel className="font-normal">
              <div className="flex flex-col space-y-1">
                <p className="text-sm font-medium leading-none">{user.username}</p>
                <p className="text-xs leading-none text-muted-foreground">
                  {user.email}
                </p>
                <div className="flex items-center mt-1">
                  <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    <Shield className="w-3 h-3 mr-1" />
                    {user.role}
                  </span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleLogout} className="text-red-600 focus:text-red-600">
              <LogOut className="mr-2 h-4 w-4" />
              <span>Log out</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  )
}
