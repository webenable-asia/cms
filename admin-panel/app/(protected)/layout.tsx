'use client'

import { ProtectedRoute } from '@/components/auth/protected-route'
import { AdminNav } from '@/components/admin/admin-nav'

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode
}) {
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
