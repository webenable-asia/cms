#!/bin/bash

# Admin Panel Logout Fix - Testing & Verification Script
echo "🔧 Admin Panel Logout Fix - Testing & Verification"
echo "=================================================="

# Set working directory
cd /Users/tsaa/Workspace/projects/webenable/cms

# Step 1: Verify current service status
echo "📊 Current Service Status:"
echo "========================="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""

# Step 2: Test current admin panel accessibility
echo "🧪 Testing Current Admin Panel Access:"
echo "======================================"

echo "Testing /admin/ route..."
curl -I http://localhost/admin/ 2>/dev/null | grep -E "(HTTP|Cache-Control|Location)" || echo "❌ Admin root not accessible"

echo "Testing /admin/login route..."
curl -I http://localhost/admin/login 2>/dev/null | grep -E "(HTTP|Cache-Control)" || echo "❌ Admin login not accessible"

echo ""

# Step 3: Create a simple test without complex middleware
echo "🔧 Creating Simple Test Environment:"
echo "===================================="

# Create a minimal test page to verify logout functionality
mkdir -p admin-panel/app/test-logout

cat > admin-panel/app/test-logout/page.tsx << 'EOF'
'use client'

import { useAuth } from '@/hooks/use-auth'
import { Button } from '@/components/ui/button'
import { useRouter } from 'next/navigation'

export default function TestLogout() {
  const { user, logout, isAuthenticated } = useAuth()
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
        <h2 className="font-semibold mb-2">Current User:</h2>
        <p>Username: {user?.username}</p>
        <p>Email: {user?.email}</p>
        <p>Role: {user?.role}</p>
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
EOF

echo "✅ Created test logout page at /admin/test-logout"
echo ""
echo "🎯 Summary of Logout Fixes Applied:"
echo "==================================="
echo "✅ Fixed admin-nav.tsx - Correct route redirect (/login instead of /admin/auth/login)"
echo "✅ Enhanced use-auth.ts - Comprehensive token cleanup"
echo "✅ Improved token manager - Auto-expiry and complete localStorage cleanup"
echo "✅ Updated login page - Better state management and redirects"
echo "✅ Fixed app root page - Uses router.replace for clean navigation"
echo ""
echo "📋 Manual Testing Instructions:"
echo "=============================="
echo ""
echo "1. Open browser and navigate to: http://localhost/admin"
echo "   - Should redirect to login page"
echo ""
echo "2. Login with admin credentials:"
echo "   - Username: admin"
echo "   - Password: /juk+vfdbNk6TICg"
echo ""
echo "3. After login, test logout in two ways:"
echo ""
echo "   Method A - Test Page:"
echo "   - Navigate to: http://localhost/admin/test-logout"
echo "   - Use the 'Test Logout' button"
echo "   - Check browser console for logs"
echo ""
echo "   Method B - Main Navigation:"
echo "   - Go to dashboard: http://localhost/admin/dashboard"
echo "   - Click logout button in top-right dropdown"
echo "   - Should redirect to login page"
echo ""
echo "4. Verify complete logout:"
echo "   - Check that localStorage is cleared"
echo "   - Try accessing dashboard directly - should redirect to login"
echo "   - Login page should not flash if already logged out"
echo ""
echo "🎉 Logout functionality should now work correctly!"
echo ""
echo "Expected Results:"
echo "✅ Logout button redirects to /admin/login (not /admin/auth/login)"
echo "✅ localStorage is completely cleared after logout"
echo "✅ No authentication state persists after logout"
echo "✅ Router navigation works without page refresh"
