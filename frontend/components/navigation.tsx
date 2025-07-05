"use client"

import Link from 'next/link'
import { ThemeToggle } from '@/components/ui/theme-toggle-reference'
import { Button } from '@/components/ui/button'

export function Navigation() {
  return (
    <nav className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center space-x-8">
            <Link href="/" className="flex items-center space-x-2">
              <span className="text-2xl font-bold text-primary">WebEnable</span>
            </Link>
            
            <div className="hidden md:flex space-x-6">
              <Link 
                href="/" 
                className="text-foreground/80 hover:text-foreground transition-colors"
              >
                Home
              </Link>
              <Link 
                href="/blog" 
                className="text-foreground/80 hover:text-foreground transition-colors"
              >
                Blog
              </Link>
              <Link 
                href="/contact" 
                className="text-foreground/80 hover:text-foreground transition-colors"
              >
                Contact
              </Link>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {/* Theme Toggle from Reference Repository */}
            <ThemeToggle />
            
            <Link href="/admin">
              <Button variant="outline" size="sm">
                Admin
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </nav>
  )
}
