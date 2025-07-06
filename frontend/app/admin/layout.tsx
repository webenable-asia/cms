'use client'

import { usePathname } from 'next/navigation'
import { ProtectedRoute } from '@/components/auth/protected-route'
import { AdminNav } from '@/components/admin/admin-nav'

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const pathname = usePathname()
  const isLoginPage = pathname === '/admin/login'

  if (isLoginPage) {
    return (
      <div className="min-h-screen bg-gray-50">
        {children}
      </div>
    )
  }

  return (
    <ProtectedRoute requiredRole="admin">
      <div className="min-h-screen bg-gray-50">
        <AdminNav />
        <main className="container mx-auto py-6">
          {children}
        </main>
      </div>
    </ProtectedRoute>
  )
}
