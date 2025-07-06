import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  // Only apply to admin routes
  if (request.nextUrl.pathname.startsWith('/admin')) {
    // Allow login page to be accessed without authentication
    if (request.nextUrl.pathname === '/admin/login') {
      return NextResponse.next()
    }

    // For cookieless auth, we can't reliably check JWT tokens in middleware
    // since localStorage is not accessible server-side
    // Let the client-side auth handle all authentication logic
    return NextResponse.next()
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    '/admin/:path*'
  ]
}
