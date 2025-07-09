import { NextRequest, NextResponse } from "next/server"

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl
  
  // Add security headers for cache prevention
  const response = NextResponse.next()
  
  // Prevent caching of admin pages
  response.headers.set('Cache-Control', 'no-cache, no-store, must-revalidate, max-age=0')
  response.headers.set('Pragma', 'no-cache')
  response.headers.set('Expires', '0')
  
  // Simple authentication check
  const hasTokenInCookie = request.cookies.has('webenable_token')
  const hasTokenInStorage = request.headers.get('authorization')?.startsWith('Bearer ')
  
  // Check localStorage token via custom header (set by client)
  const hasClientToken = request.headers.get('x-auth-token') === 'true'
  
  const hasToken = hasTokenInCookie || hasTokenInStorage || hasClientToken

  // If accessing login page and already authenticated, redirect to dashboard
  if (pathname === '/login' && hasToken) {
    return NextResponse.redirect(new URL('/dashboard', request.url))
  }

  // If accessing root, redirect to dashboard or login
  if (pathname === '/') {
    if (hasToken) {
      return NextResponse.redirect(new URL('/dashboard', request.url))
    } else {
      return NextResponse.redirect(new URL('/login', request.url))
    }
  }

  // Protected routes that require authentication
  const protectedPaths = ['/dashboard', '/posts', '/users', '/contacts']
  const isProtected = protectedPaths.some(path => pathname.startsWith(path))
  
  if (isProtected && !hasToken) {
    // Clear any existing cookies and redirect to login
    const loginResponse = NextResponse.redirect(new URL('/login', request.url))
    loginResponse.cookies.delete('webenable_token')
    return loginResponse
  }

  return response
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    '/((?!api|_next/static|_next/image|favicon.ico|_next/webpack-hmr).*)',
  ],
}
