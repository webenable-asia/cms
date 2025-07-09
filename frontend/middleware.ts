import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  // Currently no middleware logic needed for public frontend
  // All admin functionality moved to separate admin-panel service
  return NextResponse.next()
}

export const config = {
  matcher: [
    // No protected routes in frontend - all admin routes moved to admin-panel
  ]
}
